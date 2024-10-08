package deform

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/miter"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

func DeformModel(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	frame int,
) *pmx.PmxModel {
	vmdDeltas := delta.NewVmdDeltas(float32(frame), model.Bones, "", "")
	vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, float32(frame), nil)
	vmdDeltas = DeformBoneByPhysicsFlag(model, motion, vmdDeltas, false, float32(frame), nil, false)

	// 頂点にボーン変形を適用
	for _, vertex := range model.Vertices.Data {
		mat := &mmath.MMat4{}
		for j := range vertex.Deform.AllIndexes() {
			boneIndex := vertex.Deform.AllIndexes()[j]
			weight := vertex.Deform.AllWeights()[j]
			mat.Add(vmdDeltas.Bones.Get(boneIndex).FilledLocalMatrix().MuledScalar(weight))
		}

		var morphDelta *delta.VertexMorphDelta
		if vmdDeltas.Morphs != nil && vmdDeltas.Morphs.Vertices != nil {
			morphDelta = vmdDeltas.Morphs.Vertices.Get(vertex.Index())
		}

		// 頂点変形
		if morphDelta == nil {
			vertex.Position = mat.MulVec3(vertex.Position)
		} else {
			vertex.Position = mat.MulVec3(vertex.Position.Added(morphDelta.Position))
		}

		// 法線変形
		vertex.Normal = mat.MulVec3(vertex.Normal).Normalized()

		// SDEFの場合、パラメーターを再計算
		switch sdef := vertex.Deform.(type) {
		case *pmx.Sdef:
			// SDEF-C: ボーンのベクトルと頂点の交点
			sdef.SdefC = mmath.IntersectLinePoint(
				vmdDeltas.Bones.Get(sdef.AllIndexes()[0]).GlobalPosition,
				vmdDeltas.Bones.Get(sdef.AllIndexes()[1]).GlobalPosition,
				vertex.Position,
			)

			// SDEF-R0: 0番目のボーンとSDEF-Cの中点
			sdef.SdefR0 = vmdDeltas.Bones.Get(sdef.AllIndexes()[0]).GlobalPosition.Added(sdef.SdefC).MuledScalar(0.5)

			// SDEF-R1: 1番目のボーンとSDEF-Cの中点
			sdef.SdefR1 = vmdDeltas.Bones.Get(sdef.AllIndexes()[1]).GlobalPosition.Added(sdef.SdefC).MuledScalar(0.5)
		}
	}

	// ボーンの位置を更新
	for i, bone := range model.Bones.Data {
		if vmdDeltas.Bones.Get(i) != nil {
			bone.Position = vmdDeltas.Bones.Get(i).FilledGlobalPosition()
		}
	}

	return model
}

func DeformIk(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	frame float32,
	ikBone *pmx.Bone,
	ikGlobalPosition *mmath.MVec3,
) *delta.VmdDeltas {
	ikTargetBone := model.Bones.Get(ikBone.Ik.BoneIndex)

	ikTargetDeformBoneIndexes, deltas := newVmdDeltas(model, motion, nil, frame, []string{ikTargetBone.Name()}, false)
	deltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)

	// IKターゲットをIKの現在情報として埋める
	deltas = DeformBoneByPhysicsFlag(model, motion, deltas, false, frame, []string{ikTargetBone.Name()}, false)
	var prefixPath string
	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		dirPath := fmt.Sprintf("%s/IK_step", filepath.Dir(model.Path()))
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		date := time.Now().Format("20060102_150405")
		prefixPath = fmt.Sprintf("%s/%.3f_%s_%03d_%03d", dirPath, frame, date, 0, 0)
	}

	var err error
	var ikFile *os.File
	var ikMotion *vmd.VmdMotion
	var globalMotion *vmd.VmdMotion
	count := 1

	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		ikMotionPath := fmt.Sprintf("%s_%s.vmd", prefixPath, ikBone.Name())
		ikMotion = vmd.NewVmdMotion(ikMotionPath)

		globalMotionPath := fmt.Sprintf("%s_%s_global.vmd", prefixPath, ikBone.Name())
		globalMotion = vmd.NewVmdMotion(globalMotionPath)

		ikLogPath := fmt.Sprintf("%s_%s.log", prefixPath, ikBone.Name())
		ikFile, err = os.OpenFile(ikLogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(ikFile, "----------------------------------------")
		fmt.Println(ikFile, "[IK計算出力先][%.3f][%s] %s", frame, ikMotionPath)
	}
	defer func() {
		mlog.IV("[IK計算終了][%.3f][%s]", frame, ikBone.Name())

		if ikMotion != nil {
			r := repository.NewVmdRepository()
			r.Save("", ikMotion, true)
		}
		if globalMotion != nil {
			r := repository.NewVmdRepository()
			r.Save("", globalMotion, true)
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

	// IK計算
ikLoop:
	for loop := 0; loop < ikBone.Ik.LoopCount; loop++ {
		for lidx, ikLink := range ikBone.Ik.Links {
			// ikLink は末端から並んでる
			if !model.Bones.Contains(ikLink.BoneIndex) {
				continue
			}

			// 処理対象IKリンクボーン
			linkBone := model.Bones.Get(ikLink.BoneIndex)

			// 角度制限があってまったく動かさない場合、IK計算しないで次に行く
			if (linkBone.Extend.AngleLimit &&
				linkBone.Extend.MinAngleLimit.Radians().IsZero() &&
				linkBone.Extend.MaxAngleLimit.Radians().IsZero()) ||
				(linkBone.Extend.LocalAngleLimit &&
					linkBone.Extend.LocalMinAngleLimit.Radians().IsZero() &&
					linkBone.Extend.LocalMaxAngleLimit.Radians().IsZero()) {
				continue
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d] -------------------------------------------- \n",
					frame, loop, linkBone.Name(), count-1)
			}

			for _, l := range ikBone.Ik.Links {
				if deltas.Bones.Get(l.BoneIndex) != nil {
					deltas.Bones.Get(l.BoneIndex).UnitMatrix = nil
				}
			}
			if deltas.Bones.Get(ikBone.Ik.BoneIndex) != nil {
				deltas.Bones.Get(ikBone.Ik.BoneIndex).UnitMatrix = nil
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := vmd.NewBoneFrame(float32(count))
				bf.Position = ikGlobalPosition
				ikMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
				count++

				fmt.Fprintf(ikFile, "[%.3f][%03d][%s][%05d][Local] ikGlobalPosition: %s\n",
					frame, loop, linkBone.Name(), count-1, bf.Position.MMD().String())
			}

			// IK関連の行列を取得
			deltas.Bones = calcBoneDeltas(frame, model, ikTargetDeformBoneIndexes, deltas.Bones, false)

			// リンクボーンの変形情報を取得
			linkDelta := deltas.Bones.Get(linkBone.Index())
			if linkDelta == nil {
				linkDelta = &delta.BoneDelta{Bone: linkBone, Frame: frame}
			}
			linkQuat := linkDelta.FilledTotalRotation()

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := vmd.NewBoneFrame(float32(count))
				bf.Rotation = linkQuat.Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
				count++

				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][linkQuat] %s(%s)\n",
					frame, loop, linkBone.Name(), count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String(),
				)
			}

			// 現在のIKターゲットボーンのグローバル位置を取得
			ikTargetGlobalPosition := deltas.Bones.Get(ikTargetBone.Index()).FilledGlobalPosition()

			// 注目ノード（実際に動かすボーン=リンクボーン）
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := deltas.Bones.Get(linkBone.Index()).FilledGlobalMatrix().Inverted()

			if mlog.IsIkVerbose() && globalMotion != nil && ikFile != nil {
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Position = ikGlobalPosition
					globalMotion.AppendRegisteredBoneFrame(ikBone.Name(), bf)
				}
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Position = ikTargetGlobalPosition
					globalMotion.AppendRegisteredBoneFrame(ikTargetBone.Name(), bf)
				}
				count++
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Global] [%s]ikGlobalPosition: %s, "+
						"[%s]ikTargetGlobalPosition: %s, [%s]linkGlobalPosition: %s\n",
					frame, loop, linkBone.Name(), count-1,
					ikBone.Name(), ikGlobalPosition.MMD().String(),
					ikTargetBone.Name(), ikTargetGlobalPosition.MMD().String(),
					linkBone.Name(), deltas.Bones.Get(linkBone.Index()).FilledGlobalPosition().MMD().String())
			}

			// 注目ノードを起点とした、IKターゲットのローカル位置
			ikTargetLocalPosition := linkInvMatrix.MulVec3(ikTargetGlobalPosition)
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Local] ikTargetLocalPosition: %s, ikLocalPosition: %s (%f)\n",
					frame, loop, linkBone.Name(), count-1,
					ikTargetLocalPosition.MMD().String(), ikLocalPosition.MMD().String(),
					ikTargetLocalPosition.Distance(ikLocalPosition))
			}

			ikTargetLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			distanceThreshold := ikTargetLocalPosition.Distance(ikLocalPosition)
			if distanceThreshold < 1e-5 {
				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][Local] ***BREAK*** distanceThreshold: %f\n",
						frame, loop, linkBone.Name(), count-1, distanceThreshold)
				}

				break ikLoop
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][Local] ikTargetLocalPositionNorm: %s, ikLocalPositionNorm: %s\n",
					frame, loop, linkBone.Name(), count-1,
					ikTargetLocalPosition.MMD().String(), ikLocalPosition.MMD().String())
			}

			// 単位角
			unitRad := ikBone.Ik.UnitRotation.Radians().X * float64(lidx+1)
			linkDot := ikLocalPosition.Dot(ikTargetLocalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][回転角度] unitRad: %.8f (%.5f), linkDot: %.8f\n",
					frame, loop, linkBone.Name(), count-1, unitRad, mmath.ToDegree(unitRad), linkDot,
				)
			}

			// 回転角(ラジアン)
			// 単位角を超えないようにする
			originalLinkAngle := math.Acos(mmath.ClampedFloat(linkDot, -1, 1))
			linkAngle := mmath.ClampedFloat(originalLinkAngle, -unitRad, unitRad)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][単位角制限] linkAngle: %.8f(%.5f), originalLinkAngle: %.8f(%.5f)\n",
					frame, loop, linkBone.Name(), count-1, linkAngle, mmath.ToDegree(linkAngle),
					originalLinkAngle, mmath.ToDegree(originalLinkAngle),
				)
			}

			// 回転軸
			var originalLinkAxis, linkAxis *mmath.MVec3
			// 一段IKでない場合、または一段IKでかつ回転角が88度以上の場合
			if ikLink.AngleLimit {
				// グローバル軸制限
				linkAxis, originalLinkAxis = getLinkAxis(
					ikLink.MinAngleLimit.Radians(),
					ikLink.MaxAngleLimit.Radians(),
					ikTargetLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			} else if ikLink.LocalAngleLimit {
				// ローカル軸制限
				linkAxis, originalLinkAxis = getLinkAxis(
					ikLink.LocalMinAngleLimit.Radians(),
					ikLink.LocalMaxAngleLimit.Radians(),
					ikTargetLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			} else {
				// 軸制限なし or 一段IKでかつ回転角が88度未満の場合
				linkAxis, originalLinkAxis = getLinkAxis(
					mmath.MVec3MinVal,
					mmath.MVec3MaxVal,
					ikTargetLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name(), ikMotion, ikFile,
				)
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][回転軸] linkAxis: %s, originalLinkAxis: %s\n",
					frame, loop, linkBone.Name(), count-1, linkAxis.String(), originalLinkAxis.String(),
				)
			}

			originalIkQuat := mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, originalLinkAngle)
			ikQuat := mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, linkAngle)

			originalTotalIkQuat := linkQuat.Muled(originalIkQuat)
			totalIkQuat := linkQuat.Muled(ikQuat)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = originalTotalIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][originalTotalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, originalTotalIkQuat.String(), originalTotalIkQuat.ToMMDDegrees().String())

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][originalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, originalIkQuat.String(), originalIkQuat.ToMMDDegrees().String())
				}
				{
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = totalIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][totalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String())

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][ikQuat] %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			}

			var resultIkQuat *mmath.MQuaternion
			if ikLink.AngleLimit {
				// 角度制限が入ってる場合
				resultIkQuat, count = calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.MinAngleLimit.Radians(),
					ikLink.MaxAngleLimit.Radians(),
					mmath.MVec3UnitX, mmath.MVec3UnitY, mmath.MVec3UnitZ,
					loop, ikBone.Ik.LoopCount,
					frame, count, linkBone.Name(), ikMotion, ikFile,
				)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][角度制限後] resultIkQuat: %s(%s), totalIkQuat: %s(%s), ikQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String(),
						totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
						ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			} else if ikLink.LocalAngleLimit {
				// ローカル角度制限が入ってる場合
				resultIkQuat, count = calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.LocalMinAngleLimit.Radians(),
					ikLink.LocalMaxAngleLimit.Radians(),
					linkBone.Extend.NormalizedLocalAxisX,
					linkBone.Extend.NormalizedLocalAxisY,
					linkBone.Extend.NormalizedLocalAxisZ,
					loop, ikBone.Ik.LoopCount,
					frame, count, linkBone.Name(), ikMotion, ikFile,
				)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][ローカル角度制限後] resultIkQuat: %s(%s), totalIkQuat: %s(%s), ikQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String(),
						totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
						ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			} else {
				// 角度制限なしの場合
				resultIkQuat = totalIkQuat
			}

			if loop == 0 && deltas.Morphs != nil && deltas.Morphs.Bones != nil &&
				deltas.Morphs.Bones.Get(linkBone.Index()) != nil &&
				deltas.Morphs.Bones.Get(linkBone.Index()).FrameRotation != nil {
				// モーフ変形がある場合、モーフ変形を追加適用
				resultIkQuat = resultIkQuat.Muled(deltas.Morphs.Bones.Get(linkBone.Index()).FrameRotation)
			}

			if linkBone.HasFixedAxis() {
				// 軸制限ありの場合、軸にそった理想回転量とする
				resultIkQuat = resultIkQuat.ToFixedAxisRotation(linkBone.FixedAxis)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := vmd.NewBoneFrame(float32(count))
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
					count++

					fmt.Fprintf(ikFile,
						"[%.3f][%03d][%s][%05d][軸制限後] resultIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name(), count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String())
				}
			}

			// IKの結果を更新
			linkDelta.FrameRotation = resultIkQuat
			deltas.Bones.Update(linkDelta)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := vmd.NewBoneFrame(float32(count))
				bf.Rotation = linkDelta.FilledTotalRotation().Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name(), bf)
				count++

				fmt.Fprintf(ikFile,
					"[%.3f][%03d][%s][%05d][結果] bf.Rotation: %s(%s)\n",
					frame, loop, linkBone.Name(), count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String())
			}
		}
	}

	// ボーンデフォーム情報を埋める
	deltas.Bones = fillBoneDeform(model, motion, deltas, frame, ikTargetDeformBoneIndexes)
	deltas.Bones = calcBoneDeltas(frame, model, ikTargetDeformBoneIndexes, deltas.Bones, true)

	return deltas
}

// DeformBone 前回情報なしでボーンデフォーム処理を実行する
func DeformBone(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	isCalcIk bool,
	frame int,
	boneNames []string,
) *delta.BoneDeltas {
	return DeformBoneByPhysicsFlag(model, motion, nil, isCalcIk, float32(frame), boneNames, false).Bones
}

func deformBeforePhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion, vmdDeltas *delta.VmdDeltas,
) *delta.VmdDeltas {
	frame := appState.Frame()

	if vmdDeltas == nil || vmdDeltas.Frame() != frame ||
		vmdDeltas.ModelHash() != model.Hash() || vmdDeltas.MotionHash() != motion.Hash() {
		vmdDeltas = delta.NewVmdDeltas(frame, model.Bones, model.Hash(), motion.Hash())
		vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)
		vmdDeltas = DeformBoneByPhysicsFlag(model, motion, vmdDeltas, true, frame, nil, false)
	}

	return vmdDeltas
}

func DeformPhysicsByBone(
	appState state.IAppState, model *pmx.PmxModel, vmdDeltas *delta.VmdDeltas, physics *mbt.MPhysics,
) {
	// 物理剛体位置を更新
	processFunc := func(i int) {
		rigidBody := model.RigidBodies.Get(i)

		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}

		if rigidBodyBone == nil || vmdDeltas.Bones.Get(rigidBodyBone.Index()) == nil {
			return
		}

		if (appState.IsEnabledPhysics() && rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC) ||
			appState.IsPhysicsReset() {
			// 通常はボーン追従剛体・物理＋ボーン剛体だけ。物理リセット時は全部更新
			physics.UpdateTransform(model.Index(), rigidBodyBone,
				vmdDeltas.Bones.Get(rigidBodyBone.Index()).FilledGlobalMatrix(), rigidBody)
		}
	}

	// 100件ずつ処理
	miter.IterParallelByCount(model.RigidBodies.Len(), 100, processFunc)
}

func DeformBonePyPhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion,
	vmdDeltas *delta.VmdDeltas, physics *mbt.MPhysics,
) *delta.VmdDeltas {
	if model != nil && appState.IsEnabledPhysics() && !appState.IsPhysicsReset() {
		// 物理剛体位置を更新
		processFunc := func(data, index int) {
			bone := model.Bones.Get(data)
			if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
				return
			}
			bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(model.Index(), bone.Extend.RigidBody)
			if vmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
				bd := delta.NewBoneDeltaByGlobalMatrix(bone, appState.Frame(),
					bonePhysicsGlobalMatrix, vmdDeltas.Bones.Get(bone.ParentIndex))
				vmdDeltas.Bones.Update(bd)
			}
		}

		// 100件ずつ処理
		miter.IterParallelByList(model.Bones.LayerSortedIndexes, 100, processFunc)
	}

	// 物理後のデフォーム情報
	return DeformBoneByPhysicsFlag(model, motion, vmdDeltas, true, appState.Frame(), nil, true)
}

func DeformForReset(
	physics *mbt.MPhysics, appState state.IAppState, timeStep float32,
	models []*pmx.PmxModel, motions []*vmd.VmdMotion, vmdDeltas []*delta.VmdDeltas,
) []*delta.VmdDeltas {
	// 物理前デフォーム
	for i := range models {
		if models[i] == nil || motions[i] == nil {
			continue
		}
		for i >= len(vmdDeltas) {
			vmdDeltas = append(vmdDeltas, nil)
		}
		// 物理前
		vmdDeltas[i] = DeformBoneByPhysicsFlag(models[i], motions[i], vmdDeltas[i], true, appState.Frame(), nil, false)
		// 物理後
		vmdDeltas[i] = DeformBoneByPhysicsFlag(models[i], motions[i], vmdDeltas[i], true, appState.Frame(), nil, true)
	}

	return vmdDeltas
}

func Deform(
	physics *mbt.MPhysics, appState state.IAppState, timeStep float32,
	models []*pmx.PmxModel, motions []*vmd.VmdMotion, vmdDeltas []*delta.VmdDeltas,
) []*delta.VmdDeltas {
	// 物理前デフォーム
	for i := range models {
		if models[i] == nil || motions[i] == nil {
			continue
		}
		for i >= len(vmdDeltas) {
			vmdDeltas = append(vmdDeltas, nil)
		}
		vmdDeltas[i] = deformBeforePhysics(appState, models[i], motions[i], vmdDeltas[i])
	}

	return vmdDeltas
}

func DeformPhysics(
	physics *mbt.MPhysics, appState state.IAppState, timeStep float32,
	models []*pmx.PmxModel, motions []*vmd.VmdMotion, vmdDeltas []*delta.VmdDeltas,
) []*delta.VmdDeltas {
	// 物理デフォーム
	for i := range models {
		if models[i] == nil || vmdDeltas[i] == nil {
			continue
		}
		DeformPhysicsByBone(appState, models[i], vmdDeltas[i], physics)
	}

	if appState.IsEnabledPhysics() || appState.IsPhysicsReset() {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	for i := range models {
		if models[i] == nil || motions[i] == nil || vmdDeltas[i] == nil {
			continue
		}
		vmdDeltas[i] = DeformBonePyPhysics(appState, models[i], motions[i], vmdDeltas[i], physics)
	}

	return vmdDeltas
}
