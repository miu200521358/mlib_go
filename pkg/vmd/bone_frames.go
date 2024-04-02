package vmd

import (
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type BoneFrames struct {
	Data map[string]*BoneNameFrames
}

func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		Data: make(map[string]*BoneNameFrames, 0),
	}
}

const (
	// 88.0f / 180.0f*3.14159265f
	GIMBAL_RAD  = math.Pi * 89.99 / 180
	GIMBAL2_RAD = math.Pi * 89.99 * 2 / 180
	QUARTER_RAD = math.Pi / 2
	HALF_RAD    = math.Pi
	FULL_RAD    = math.Pi * 2
)

func (bfs *BoneFrames) Contains(boneName string) bool {
	_, ok := bfs.Data[boneName]
	return ok
}

func (bfs *BoneFrames) Append(bnfs *BoneNameFrames) {
	bfs.Data[bnfs.Name] = bnfs
}

func (bfs *BoneFrames) GetItem(boneName string) *BoneNameFrames {
	if !bfs.Contains(boneName) {
		bfs.Append(NewBoneNameFrames(boneName))
	}
	return bfs.Data[boneName]
}

func (bfs *BoneFrames) GetNames() []string {
	names := make([]string, 0, len(bfs.Data))
	for name := range bfs.Data {
		names = append(names, name)
	}
	return names
}

func (bfs *BoneFrames) GetIndexes() []float32 {
	indexes := make([]float32, 0)
	for _, bnfs := range bfs.Data {
		for _, index := range bnfs.Indexes {
			if !slices.Contains(indexes, index) {
				indexes = append(indexes, index)
			}
		}
	}
	slices.Sort(indexes)
	return indexes
}

func (bfs *BoneFrames) GetRegisteredIndexes() []float32 {
	indexes := make([]float32, 0)
	for _, bnfs := range bfs.Data {
		for _, index := range bnfs.RegisteredIndexes {
			if !slices.Contains(indexes, index) {
				indexes = append(indexes, index)
			}
		}
	}
	slices.Sort(indexes)
	return indexes
}

func (bfs *BoneFrames) GetCount() int {
	count := 0
	for _, bnfs := range bfs.Data {
		count += len(bnfs.RegisteredIndexes)
	}
	return count
}

func (bfs *BoneFrames) GetMaxFrame() float32 {
	maxFno := float32(0)
	for _, bnfs := range bfs.Data {
		fno := bnfs.GetMaxFrame()
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (bfs *BoneFrames) Animate(
	frame float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isCalcMorph bool,
) *BoneDeltas {
	// 処理対象ボーン一覧取得
	targetBoneNames, targetBoneIndexes := bfs.getAnimatedBoneNames(model, boneNames)

	// IK事前計算
	if isCalcIk {
		bfs.prepareIkSolvers(frame, model, targetBoneNames, isCalcMorph)
	}

	// ボーン変形行列操作
	positions, rotations, scales, _ :=
		bfs.getBoneMatrixes(frame, model, targetBoneNames, targetBoneIndexes, isCalcMorph)

	// ボーン行列計算
	return bfs.calcBoneMatrixes(
		frame,
		model,
		targetBoneNames,
		targetBoneIndexes,
		positions,
		rotations,
		scales,
	)
}

// IK事前計算処理
func (bfs *BoneFrames) prepareIkSolvers(
	frame float32,
	model *pmx.PmxModel,
	targetBoneNames map[string]int,
	isCalcMorph bool,
) {
	var wg sync.WaitGroup
	for boneName := range targetBoneNames {
		bone := model.Bones.GetItemByName(boneName)
		// ボーンIndexがIkTreeIndexesに含まれていない場合、スルー
		if _, ok := model.Bones.IkTreeIndexes[bone.Index]; !ok {
			continue
		}

		wg.Add(1)
		go func(bone *pmx.Bone) {
			defer wg.Done()
			for i := 0; i < len(model.Bones.IkTreeIndexes[bone.Index]); i++ {
				ikBone := model.Bones.GetItem(model.Bones.IkTreeIndexes[bone.Index][i])
				// IK計算
				quats, effectorTargetBoneNames :=
					bfs.calcIk(frame, ikBone, model, isCalcMorph)

				for _, linkIndex := range ikBone.Ik.Links {
					// IKリンクボーンの回転量を更新
					linkBone := model.Bones.GetItem(linkIndex.BoneIndex)
					linkBf := bfs.GetItem(linkBone.Name).GetItem(frame)
					linkIndex := effectorTargetBoneNames[linkBone.Name]
					linkBf.IkRotation = mmath.NewRotationModelByQuaternion(quats[linkIndex])

					// IK用なので登録フラグは既存のままで追加して補間曲線は分割しない
					bfs.GetItem(linkBone.Name).Append(linkBf)
				}
			}
		}(bone)
	}
	wg.Wait()
}

// IK計算
func (bfs *BoneFrames) calcIk(
	frame float32,
	ikBone *pmx.Bone,
	model *pmx.PmxModel,
	isisCalcMorph bool,
) ([]*mmath.MQuaternion, map[string]int) {
	// IKターゲットボーン
	effectorBone := model.Bones.GetItem(ikBone.Ik.BoneIndex)
	// IK関連の行列を一括計算
	ikMatrixes := bfs.Animate(frame, model, []string{ikBone.Name}, false, isisCalcMorph)
	// 処理対象ボーン名取得
	effectorTargetBoneNames, effectorTargetBoneIndexes := bfs.getAnimatedBoneNames(model, []string{effectorBone.Name})
	// エフェクタボーンの関連ボーンの初期値を取得
	positions, rotations, scales, quats :=
		bfs.getBoneMatrixes(frame, model, effectorTargetBoneNames, effectorTargetBoneIndexes, isisCalcMorph)
	// 中断FLGが入ったか否か
	aborts := make([]bool, len(ikBone.Ik.Links))

	// IK計算でバッグ用モーション
	ikMotion := NewVmdMotion(fmt.Sprintf("E:/MMD_E/サイジング/足IK/IK_step_go/%s_%s.vmd",
		time.Now().Format("20060102_150405"), ikBone.Name))
	count := float32(1.0)
	defer Write(ikMotion)

	// IK計算
ikLoop:
	for loop := 0; loop < max(ikBone.Ik.LoopCount, 1); loop++ {
		for lidx, ikLink := range ikBone.Ik.Links {
			// ikLink は末端から並んでる
			if !model.Bones.Contains(ikLink.BoneIndex) {
				continue
			}

			// 処理対象IKリンクボーン
			linkBone := model.Bones.GetItem(ikLink.BoneIndex)
			linkIndex := effectorTargetBoneNames[linkBone.Name]

			// 角度制限があってまったく動かさない場合、IK計算しないで次に行く
			if (linkBone.AngleLimit &&
				linkBone.MinAngleLimit.GetRadians().IsZero() &&
				linkBone.MaxAngleLimit.GetRadians().IsZero()) ||
				(linkBone.LocalAngleLimit &&
					linkBone.LocalMinAngleLimit.GetRadians().IsZero() &&
					linkBone.LocalMaxAngleLimit.GetRadians().IsZero()) {
				continue
			}

			// 単位角
			unitRad := ikBone.Ik.UnitRotation.GetRadians().GetX() // * float64(lidx+1) * 2
			// ループ閾値
			loopThreshold := ikBone.Ik.LoopCount / 2

			// IK関連の行列を取得
			linkMatrixes := bfs.calcBoneMatrixes(
				frame,
				model,
				effectorTargetBoneNames,
				effectorTargetBoneIndexes,
				positions,
				rotations,
				scales,
			)

			// IKボーンのグローバル位置
			ikGlobalPosition := ikMatrixes.GetItem(ikBone.Name, frame).Position

			// 現在のIKターゲットボーンのグローバル位置を取得
			effectorGlobalPosition := linkMatrixes.GetItem(effectorBone.Name, frame).Position

			// 位置の差がほとんどない場合、終了
			if ikGlobalPosition.Distance(effectorGlobalPosition) < 1e-8 {
				break ikLoop
			}

			// 注目ノード（実際に動かすボーン=リンクボーン）
			linkMatrix := linkMatrixes.GetItem(linkBone.Name, frame).GlobalMatrix
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := linkMatrix.Inverse()

			// 注目ノードを起点とした、エフェクタのローカル位置
			effectorLocalPosition := linkInvMatrix.MulVec3(effectorGlobalPosition)
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			{
				bf := NewBoneFrame(count)
				bf.Position = ikMatrixes.GetItem(ikBone.Name, frame).FramePosition
				bf.Rotation.SetQuaternion(ikMatrixes.GetItem(ikBone.Name, frame).FrameRotation)
				ikMotion.AppendRegisteredBoneFrame(ikBone.Name, bf)
				count++
			}
			{
				bf := NewBoneFrame(count)
				bf.Rotation.SetQuaternion(linkMatrixes.GetItem(effectorBone.Name, frame).FrameRotation)
				ikMotion.AppendRegisteredBoneFrame(effectorBone.Name, bf)
				count++
			}
			{
				bf := NewBoneFrame(count)
				bf.Rotation.SetQuaternion(linkMatrixes.GetItem(linkBone.Name, frame).FrameRotation)
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}

			{
				mlog.I("[%d][%s][Local] effectorLocalPosition: %s, ikLocalPosition: %s\n", loop, linkBone.Name,
					effectorLocalPosition.String(), ikLocalPosition.String())
			}

			effectorLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			{
				mlog.I("[%d][%s][Normal] effectorLocalPosition: %s, ikLocalPosition: %s\n", loop, linkBone.Name,
					effectorLocalPosition.String(), ikLocalPosition.String())
			}

			// ベクトル (1) を (2) に一致させるための最短回転量（Axis-Angle）
			// 回転軸
			linkAxis := effectorLocalPosition.Cross(ikLocalPosition).Normalize()
			// 回転角(ラジアン)
			linkAngle := math.Acos(mmath.ClampFloat(effectorLocalPosition.Dot(ikLocalPosition), -1, 1))

			// リンクボーンの角度を取得
			linkQuat := quats[linkIndex]

			{
				bf := NewBoneFrame(count)
				bf.Rotation.SetQuaternion(linkQuat)
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}
			{
				bf := NewBoneFrame(count)
				bf.Rotation.SetQuaternion(mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle))
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}
			{
				mlog.I("[%d][%s][最短] linkAngle: %.5f, linkAxis: %s\n", loop, linkBone.Name,
					180.0*linkAngle/math.Pi, linkAxis.String())
			}

			var totalActualIkQuat *mmath.MQuaternion
			if ikLink.AngleLimit || ikLink.LocalAngleLimit {
				// 角度制限が入ってる場合
				if ikLink.MinAngleLimit.GetRadians().GetX() != 0 ||
					ikLink.MaxAngleLimit.GetRadians().GetX() != 0 {
					// グローバルX: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat, count = bfs.calcSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetX(),
						ikLink.MaxAngleLimit.GetRadians().GetX(),
						unitRad, loop < loopThreshold, linkQuat, linkAxis, linkAngle, 0,
						&mmath.MVec3{1, 0, 0}, lidx, ikMotion, linkBone, count)
					// 	} else if ikLink.MinAngleLimit.GetRadians().GetY() != 0 ||
					// 		ikLink.MaxAngleLimit.GetRadians().GetY() != 0 {
					// 		// グローバルY: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					// 		linkAngle, count = bfs.calcSingleAxisRad(
					// 			ikLink.MinAngleLimit.GetRadians().GetY(),
					// 			ikLink.MaxAngleLimit.GetRadians().GetY(),
					// 			linkQuat, linkAxis, linkAngle, 1, lidx, ikMotion, linkBone, count)
					// 		linkAxis = &mmath.MVec3{0, 1, 0}
					// 	} else if ikLink.MinAngleLimit.GetRadians().GetZ() != 0 ||
					// 		ikLink.MaxAngleLimit.GetRadians().GetZ() != 0 {
					// 		// グローバルZ: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					// 		linkAngle, count = bfs.calcSingleAxisRad(
					// 			ikLink.MinAngleLimit.GetRadians().GetZ(),
					// 			ikLink.MaxAngleLimit.GetRadians().GetZ(),
					// 			linkQuat, linkAxis, linkAngle, 2, lidx, ikMotion, linkBone, count)
					// 		linkAxis = &mmath.MVec3{0, 0, 1}
					// 	}
					// } else if ikLink.LocalAngleLimit {
					// 	// ローカル軸角度制限が入っている場合、ローカル軸に合わせて理想回転を求める
					// 	if ikLink.LocalMinAngleLimit.GetRadians().GetX() != 0 ||
					// 		ikLink.LocalMaxAngleLimit.GetRadians().GetX() != 0 {
					// 		// ローカルX: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					// 		linkAngle, count = bfs.calcSingleAxisRad(
					// 			ikLink.LocalMinAngleLimit.GetRadians().GetX(),
					// 			ikLink.LocalMaxAngleLimit.GetRadians().GetX(),
					// 			linkQuat, linkAxis, linkAngle, 0, lidx, ikMotion, linkBone, count)
					// 		linkAxis = &mmath.MVec3{1, 0, 0}
					// 	} else if ikLink.LocalMinAngleLimit.GetRadians().GetY() != 0 ||
					// 		ikLink.LocalMaxAngleLimit.GetRadians().GetY() != 0 {
					// 		// ローカルY: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					// 		linkAngle, count = bfs.calcSingleAxisRad(
					// 			ikLink.LocalMinAngleLimit.GetRadians().GetY(),
					// 			ikLink.LocalMaxAngleLimit.GetRadians().GetY(),
					// 			linkQuat, linkAxis, linkAngle, 1, lidx, ikMotion, linkBone, count)
					// 		linkAxis = &mmath.MVec3{0, 1, 0}
					// 	} else if ikLink.LocalMinAngleLimit.GetRadians().GetZ() != 0 ||
					// 		ikLink.LocalMaxAngleLimit.GetRadians().GetZ() != 0 {
					// 		// ローカルZ: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					// 		linkAngle, count = bfs.calcSingleAxisRad(
					// 			ikLink.LocalMinAngleLimit.GetRadians().GetZ(),
					// 			ikLink.LocalMaxAngleLimit.GetRadians().GetZ(),
					// 			linkQuat, linkAxis, linkAngle, 2, lidx, ikMotion, linkBone, count)
					// 		linkAxis = &mmath.MVec3{0, 0, 1}
				}
			} else {
				if linkBone.HasFixedAxis() {
					// 軸制限ありの場合、軸にそった理想回転量とする
					linkAxis = linkBone.NormalizedFixedAxis
					if linkBone.NormalizedFixedAxis.Dot(linkAxis) < 0 {
						linkAngle *= -1
					}
				}

				{
					bf := NewBoneFrame(count)
					bf.Rotation.SetQuaternion(linkQuat.Muled(mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle)))
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				// if lidx > 0 {
				// 	// 根元に行くほど回転角を半分にする
				// 	linkAngle /= float64(lidx * 2)
				// }

				// {
				// 	mlog.I("[%d][%s][半分] linkAngle: %.5f\n", loop, linkBone.Name, 180.0*linkAngle/math.Pi)
				// }

				{
					bf := NewBoneFrame(count)
					bf.Rotation.SetQuaternion(linkQuat.Muled(mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle)))
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if loop < loopThreshold {
					// 単位角を超えないようにする
					linkAngle = mmath.ClampFloat(linkAngle, -unitRad, unitRad)
				}

				{
					mlog.I("[%d][%s][単位角] linkAngle: %.5f\n", loop, linkBone.Name, 180.0*linkAngle/math.Pi)
				}

				{
					bf := NewBoneFrame(count)
					bf.Rotation.SetQuaternion(linkQuat.Muled(mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle)))
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				correctIkQuat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle)

				{
					bf := NewBoneFrame(count)
					bf.Rotation.SetQuaternion(correctIkQuat)
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				// 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
				totalActualIkQuat = linkQuat.Muled(correctIkQuat)
			}

			if linkBone.HasFixedAxis() {
				// 軸制限回転を求める
				totalActualIkQuat = totalActualIkQuat.ToFixedAxisRotation(linkBone.NormalizedFixedAxis)
			}

			{
				bf := NewBoneFrame(count)
				bf.Rotation.SetQuaternion(totalActualIkQuat)
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}

			// 前回（既存）とほぼ同じ回転量の場合、中断FLGを立てる
			if 1-quats[linkIndex].Dot(totalActualIkQuat) < 1e-8 {
				aborts[lidx] = true
			}

			// IKの結果を更新
			quats[linkIndex] = totalActualIkQuat
			rotations[linkIndex] = totalActualIkQuat.ToMat4()
		}

		// すべてのリンクボーンで中断FLG = ONの場合、ループ終了
		if slices.Index(aborts, false) == -1 {
			break
		}
	}

	return quats, effectorTargetBoneNames
}

// 全ての角度をラジアン角度に分割して、そのうちのひとつの軸だけを動かす回転を取得する
// minAngleLimit: 最小軸制限（ラジアン）
// maxAngleLimit: 最大軸制限（ラジアン）
// linkQuat: 現在のリンクボーンの回転量
// quatAxis: 現在のIK回転の回転軸
// quatAngle: 現在のIK回転の回転角度（ラジアン）
// axisIndex: 制限軸INDEX
func (bfs *BoneFrames) calcSingleAxisRad(
	minAngleLimit, maxAngleLimit, unitRad float64,
	overLoopThreshold bool,
	linkQuat *mmath.MQuaternion,
	quatAxis *mmath.MVec3,
	quatAngle float64,
	axisIndex int,
	axisVector *mmath.MVec3,
	lidx int,
	ikMotion *VmdMotion,
	linkBone *pmx.Bone,
	count float32,
) (*mmath.MQuaternion, float32) {
	ikQuat := mmath.NewMQuaternionFromAxisAngles(quatAxis, quatAngle)

	{
		bf := NewBoneFrame(count)
		bf.Rotation.SetQuaternion(ikQuat)
		ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
		count++
	}

	// 軸別の角度を取得
	axisRads := ikQuat.ToEulerAngles()
	axisRad := axisRads.Vector()[axisIndex]

	// if axisRad < minAngleLimit || maxAngleLimit < axisRad {
	// 	// 角度制限をオーバーしている場合、反対側に曲げる
	// 	ikQuat.Invert()
	// 	axisRads = ikQuat.ToEulerAngles()
	// 	axisRad = axisRads.Vector()[axisIndex]
	// }

	{
		bf := NewBoneFrame(count)
		bf.Rotation.SetQuaternion(linkQuat.Muled(mmath.NewMQuaternionFromAxisAngles(axisVector, axisRad)))
		ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
		count++
	}

	// {
	// 	bf := NewBoneFrame(count)
	// 	bf.Rotation.SetQuaternion(linkQuat.Muled(mmath.NewMQuaternionFromAxisAngles(axisVector, axisRad)))
	// 	ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
	// 	count++
	// }

	// if lidx > 0 {
	// 	// 根元に行くほど回転角を半分にする
	// 	axisRad /= float64(lidx * 2)
	// }

	// {
	// 	bf := NewBoneFrame(count)
	// 	bf.Rotation.SetQuaternion(linkQuat.Muled(mmath.NewMQuaternionFromAxisAngles(axisVector, axisRad)))
	// 	ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
	// 	count++
	// }

	// {
	// 	mlog.I("[%s][軸制限半分] axisRad: %.5f\n", linkBone.Name, 180.0*axisRad/math.Pi)
	// }

	// if !overLoopThreshold {
	// 	// 単位角を超えないようにする
	// 	axisRad = mmath.ClampFloat(axisRad, -unitRad, unitRad)
	// }

	// {
	// 	mlog.I("[%s][軸制限単位角] axisRad: %.5f\n", linkBone.Name, 180.0*axisRad/math.Pi)
	// }

	// 調整した軸角度からクォータニオンを生成
	axisIkQuat := mmath.NewMQuaternionFromAxisAngles(axisVector, axisRad)

	{
		bf := NewBoneFrame(count)
		bf.Rotation.SetQuaternion(axisIkQuat)
		ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
		count++
	}

	// axisIkRad := axisIkQuat.ToEulerAngles().Vector()[axisIndex]
	// if axisIkRad < minAngleLimit || maxAngleLimit < axisIkRad {
	// 	// 角度制限をオーバーしている場合、反対側に曲げる
	// 	axisIkQuat = mmath.NewMQuaternionFromAxisAngles(axisVector, -axisRad)
	// }

	// 現在IKリンクに入る可能性のあるすべての角度
	totalIkQuat := linkQuat.Muled(axisIkQuat)

	{
		bf := NewBoneFrame(count)
		bf.Rotation.SetQuaternion(totalIkQuat)
		ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
		count++
	}

	// // 全体の軸角度を取得
	// totalAxisRad := totalIkQuat.ToSignedRadian(axisIndex)
	// if totalIkQuat.ToEulerAngles().Vector()[axisIndex] < 0 {
	// 	totalAxisRad *= -1
	// }

	// axisRad := ikQuat.ToRadian()
	// if unitRad <= axisRad {
	// 	if axisRads.Vector()[axisIndex] < 0 {
	// 		axisRad *= -1
	// 	}
	// } else {
	// 	axisRad = axisRads.Vector()[axisIndex]
	// }

	// axisRad := axisRads.Vector()[axisIndex]

	// axisRad := ikQuat.ToRadian()
	// if unitRad >= axisRad {
	// 	if axisRads.Vector()[axisIndex] < 0 {
	// 		axisRad *= -1
	// 	}
	// } else {
	// 	axisRad = axisRads.Vector()[axisIndex]
	// }

	// axisRad := quatAngle
	// totalAxisRad := totalIkQuat.ToEulerAngles().Vector()[axisIndex]

	// {
	// 	mlog.I("[%s][制限] ikQuat: %s, totalIkQuat: %s\n", linkBone.Name, ikQuat.String(), totalIkQuat.String())
	// 	mlog.I("[%s][制限] axisRad: %.5f(%s), totalAxisRad: %.5f(%s)\n", linkBone.Name,
	// 		180.0*axisRad/math.Pi, ikQuat.ToEulerAnglesDegrees().String(),
	// 		180.0*totalAxisRad/math.Pi, totalIkQuat.ToEulerAnglesDegrees().String())
	// }

	// if lidx > 0 {
	// 	// 根元に行くほど回転角度を半分にする
	// 	axisRad /= float64(lidx * 2)
	// }

	// {
	// 	mlog.I("[%s][制限根元] axisRad: %.5f(%s), totalAxisRad: %.5f(%s)\n", linkBone.Name,
	// 		180.0*axisRad/math.Pi, ikQuat.ToEulerAnglesDegrees().String(),
	// 		180.0*totalAxisRad/math.Pi, totalIkQuat.ToEulerAnglesDegrees().String())
	// }

	// 全体の軸角度を取得
	totalAxisRad := totalIkQuat.ToEulerAngles().Vector()[axisIndex]

	{
		mlog.I("[%s][制限] ikQuat: %s, axisIkQuat: %s, totalIkQuat: %s\n", linkBone.Name,
			ikQuat.String(), axisIkQuat.String(), totalIkQuat.String())
		mlog.I("[%s][制限] axisRad: %.5f(%s), totalAxisRad: %.5f(%s)\n", linkBone.Name,
			180.0*axisRad/math.Pi, ikQuat.ToEulerAnglesDegrees().String(),
			180.0*totalAxisRad/math.Pi, totalIkQuat.ToEulerAnglesDegrees().String())
	}

	if totalAxisRad < minAngleLimit || maxAngleLimit < totalAxisRad {
		// 角度制限をオーバーしている場合、反対側に曲げる

		// 	// invertedIkQuat := mmath.NewMQuaternionByValues(
		// 	// 	-ikQuat.GetX(), -ikQuat.GetY(), -ikQuat.GetZ(), ikQuat.GetW()).Normalize()
		// 	// invertedIkQuat := mmath.NewMQuaternionFromEulerAngles(-axisRads.GetX(), -axisRads.GetY(), -axisRads.GetZ())
		// 	// totalIkQuat = mmath.NewMQuaternionFromAxisAngles(axisVector, axisRad).Invert().Normalize()
		// totalAxisRad = totalIkQuat.Invert().ToEulerAngles().Vector()[axisIndex]

		totalIkQuat = linkQuat.Muled(axisIkQuat.Invert())

		// 	// invertedIkQuat.Normalize()
		// totalIkQuat.Invert()
		// 	totalAxisRad *= -1

		// {
		// 	mlog.I("[%s][制限逆] totalAxisRad: %.5f\n", linkBone.Name, 180.0*totalAxisRad/math.Pi)
		// }

		// 	// {
		// 	// 	bf := NewBoneFrame(count)
		// 	// 	bf.Rotation.SetQuaternion(invertedIkQuat)
		// 	// 	ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
		// 	// 	count++
		// 	// }

		// totalIkQuat = mmath.NewMQuaternionFromAxisAngles(axisVector, totalAxisRad)

		// 	// totalAxisRad = totalIkQuat.ToEulerAngles().Vector()[axisIndex]
		// 	// invertedAxisRads := invertedIkQuat.ToEulerAngles()

		// 	// // axisRad = invertedIkQuat.ToRadian()
		// 	// // if unitRad >= axisRad {
		// 	// // 	if invertedAxisRads.Vector()[axisIndex] < 0 {
		// 	// // 		axisRad *= -1
		// 	// // 	}
		// 	// // } else {
		// 	// axisRad = invertedAxisRads.Vector()[axisIndex]
		// 	// // }

		// 	// axisRad = invertedIkQuat.ToRadian()
		// 	// if invertedAxisRads.Vector()[axisIndex] < 0 {
		// 	// 	axisRad *= -1
		// 	// }

		// 	// axisRad = invertedAxisRads.Vector()[axisIndex]
		// 	// if invertedAxisRads.Vector()[axisIndex] < 0 || isGimbal {
		// 	// 	axisRad *= -1
		// 	// }

		// 	// if unitRad <= axisRad {
		// 	// 	} else {
		// 	// 		axisRad = invertedAxisRads.Vector()[axisIndex]
		// 	// 	}

		{
			mlog.I("[%s][制限逆] totalIkQuat: %s\n", linkBone.Name, totalIkQuat.String())
			// mlog.I("[%s][制限逆] axisRad: %.5f(%s), totalAxisRad: %.5f(%s)\n",
			// 	linkBone.Name, 180.0*axisRad/math.Pi, invertedIkQuat.ToEulerAnglesDegrees().String(),
			// 	180.0*totalAxisRad/math.Pi, totalIkQuat.ToEulerAnglesDegrees().String())
		}
	}

	// axisRad = mmath.ClampFloat(axisRad, -unitRad, unitRad)

	// {
	// 	mlog.I("[%s][制限Clamp] axisRad: %.5f(%s), totalAxisRad: %.5f(%s)\n", linkBone.Name,
	// 		180.0*axisRad/math.Pi, ikQuat.ToEulerAnglesDegrees().String(),
	// 		180.0*totalAxisRad/math.Pi, totalIkQuat.ToEulerAnglesDegrees().String())
	// }

	// totalIkQuat = mmath.NewMQuaternionFromAxisAngles(axisVector, totalAxisRad)

	{
		bf := NewBoneFrame(count)
		bf.Rotation.SetQuaternion(totalIkQuat)
		ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
		count++
	}

	return totalIkQuat, count

	// // 角度制限がある場合、全体の角度をその角度内に収める
	// totalLimitAxisRad := mmath.ClampFloat(
	// 	totalAxisRad,
	// 	minAngleLimit.GetRadians().Vector()[axisIndex],
	// 	maxAngleLimit.GetRadians().Vector()[axisIndex],
	// )

	// axisRad := math.Abs(ikRads.Vector()[axisIndex])
	// var limitRad float64
	// if GIMBAL_RAD < quatAngle && quatAngle < GIMBAL2_RAD {
	// 	limitRad = axisRad + HALF_RAD
	// } else {
	// 	limitRad = axisRad
	// }

	// // Calculate the maximum angle in radians
	// maxRad := math.Max(
	// 	math.Abs(minAngleLimit.GetRadians().Vector()[axisIndex]),
	// 	math.Abs(maxAngleLimit.GetRadians().Vector()[axisIndex]),
	// )

	// // 最大ラジアンが制限最大角度と等しくない場合、軸の符号を逆にする
	// var axisSign float64
	// if maxRad != math.Abs(maxAngleLimit.GetRadians().Vector()[axisIndex]) {
	// 	axisSign = -1
	// } else {
	// 	axisSign = 1
	// }

	// // 単位角で制限する
	// var limitAxisRad float64
	// if unitRadian != 0 {
	// 	limitAxisRad = math.Min(unitRadian, limitRad)
	// } else {
	// 	limitAxisRad = limitRad
	// }

	// // 単位角で制限した角度に基づいたクォータニオン
	// correctLimitIkQuat := mmath.NewMQuaternionFromAxisAngles(quatAxis, limitAxisRad)

	// // 現在IKリンクに入る可能性のあるすべての角度
	// totalIkQuat := linkQuat.Muled(correctLimitIkQuat)

	// // 全体の角度を計算する
	// totalAxisIkRad := totalIkQuat.ToRadian()
	// var totalAxisIkRads *mmath.MVec3
	// if isLocal {
	// 	// ローカル軸の場合、一旦グローバル軸に直す
	// 	totalAxisIkAxis := totalIkQuat.GetXYZ().Normalize()
	// 	totalAxisIkRad = totalIkQuat.ToRadian()
	// 	var totalAxisIkSign float64
	// 	if axisVector.Dot(totalAxisIkAxis) >= 0 {
	// 		totalAxisIkSign = 1
	// 	} else {
	// 		totalAxisIkSign = -1
	// 	}

	// 	var globalAxisVec *mmath.MVec3
	// 	if axisIndex == 0 {
	// 		globalAxisVec = &mmath.MVec3{1, 0, 0}
	// 	} else if axisIndex == 1 {
	// 		globalAxisVec = &mmath.MVec3{0, 1, 0}
	// 	} else {
	// 		globalAxisVec = &mmath.MVec3{0, 0, 1}
	// 	}

	// 	totalAxisIkQuat := mmath.NewMQuaternionFromAxisAngles(globalAxisVec, totalAxisIkRad*totalAxisIkSign)
	// 	totalAxisIkRads = totalAxisIkQuat.ToEulerAngles().MMD()
	// } else {
	// 	// MMD上でのIKリンクの角度
	// 	totalAxisIkRads = totalIkQuat.ToEulerAngles().MMD()
	// }

	// var totalAxisRad float64
	// if unitRadian > quatAngle && QUARTER_RAD > totalAxisIkRad && unitRadian > totalAxisIkRad {
	// 	// トータルが制限角度以内であれば全軸の角度を使う
	// 	totalIkQq := linkQuat.Muled(ikQuat)
	// 	totalAxisRad = totalIkQq.ToRadian() * axisSign
	// } else if GIMBAL_RAD > quatAngle && QUARTER_RAD > totalAxisIkRad && unitRadian > totalAxisIkRad {
	// 	// トータルが88度以内で、軸分け後が制限角度以内であれば制限角度を使う
	// 	totalAxisRad = unitRadian * axisSign
	// } else if HALF_RAD > totalAxisIkRad {
	// 	// トータルが180度以内であれば一軸の角度を全部使う
	// 	totalAxisRad = totalAxisIkRad * axisSign
	// } else {
	// 	// 180度を超えている場合、軸の値だけ使用する
	// 	totalAxisRad = math.Abs(totalAxisIkRads.Vector()[axisIndex]) * axisSign
	// }

	// // 角度制限がある場合、全体の角度をその角度内に収める
	// totalLimitAxisRad := mmath.ClampFloat(
	// 	totalAxisRad,
	// 	minAngleLimit.GetRadians().Vector()[axisIndex],
	// 	maxAngleLimit.GetRadians().Vector()[axisIndex],
	// )

	// // 単位角とジンバルロックの整合性を取る
	// var resultAxisRad float64
	// if GIMBAL2_RAD < totalAxisIkRad && !isLocal {
	// 	resultAxisRad = HALF_RAD + math.Abs(totalLimitAxisRad)
	// } else if GIMBAL_RAD < totalAxisIkRad && !isLocal {
	// 	resultAxisRad = FULL_RAD + totalLimitAxisRad
	// } else {
	// 	resultAxisRad = totalLimitAxisRad
	// }

	// // 指定の軸方向に回す
	// resultLinkQuat := mmath.NewMQuaternionFromAxisAngles(axisVector, resultAxisRad)
	// return resultLinkQuat
}

func (bfs *BoneFrames) calcBoneMatrixes(
	frame float32,
	model *pmx.PmxModel,
	targetBoneNames map[string]int,
	targetBoneIndexes map[int]string,
	positions, rotations, scales []*mmath.MMat4,
) *BoneDeltas {
	matrixes := make([]*mmath.MMat4, 0, len(targetBoneIndexes))
	resultMatrixes := make([]*mmath.MMat4, 0, len(targetBoneIndexes))
	boneCount := len(targetBoneNames)

	// 最初にフレーム数*ボーン数分のスライスを確保
	for i := 0; i < len(targetBoneIndexes); i++ {
		matrixes = append(matrixes, mmath.NewMMat4())
		resultMatrixes = append(resultMatrixes, mmath.NewMMat4())
	}

	// ボーンを一定件数ごとに並列処理（件数は変数保持）
	count := 100

	var wg1 sync.WaitGroup
	for i := 0; i < boneCount; i += count {
		wg1.Add(1)
		go func(i int) {
			defer wg1.Done()
			for j := i; j < i+count; j++ {
				if j >= boneCount {
					break
				}
				// 各ボーンの座標変換行列×逆BOf行列
				boneName := targetBoneIndexes[j]
				bone := model.Bones.GetItemByName(boneName)
				// 逆BOf行列(初期姿勢行列)
				matrixes[j].Mul(bone.RevertOffsetMatrix)
				// 位置
				matrixes[j].Mul(positions[j])
				// 回転
				matrixes[j].Mul(rotations[j])
				// スケール
				matrixes[j].Mul(scales[j])
			}
		}(i)
	}
	wg1.Wait()

	boneDeltas := NewBoneDeltas()

	var wg2 sync.WaitGroup
	// ボーンを一定件数ごとに並列処理（件数は変数保持）
	for i := 0; i < boneCount; i += count {
		wg2.Add(1)
		go func(i int) {
			defer wg2.Done()
			for j := i; j < i+count; j++ {
				if j >= boneCount {
					break
				}
				boneName := targetBoneIndexes[j]
				bone := model.Bones.GetItemByName(boneName)
				localMatrix := mmath.NewMMat4()
				for _, l := range bone.ParentBoneIndexes {
					// 親ボーンの変形行列を掛ける(親->子の順で掛ける)
					parentName := model.Bones.GetItem(l).Name
					// targetBoneNames の中にある parentName のINDEXを取得
					parentIndex := targetBoneNames[parentName]
					localMatrix.Mul(matrixes[parentIndex])
				}
				// 最後に対象ボーン自身の行列をかける
				localMatrix.Mul(matrixes[j])
				// BOf行列: 自身のボーンのボーンオフセット行列
				localMatrix.Mul(bone.OffsetMatrix)
				resultMatrixes[j] = localMatrix
			}
		}(i)
	}
	wg2.Wait()

	for i := 0; i < len(targetBoneIndexes); i++ {
		bone := model.Bones.GetItemByName(targetBoneIndexes[i])
		localMatrix := resultMatrixes[i]
		// 初期位置行列を掛けてグローバル行列を作成
		boneDeltas.SetItem(bone.Name, frame, NewBoneDelta(
			bone.Name,
			frame,
			localMatrix.Muled(bone.Position.ToMat4()), // グローバル行列
			localMatrix,                // ローカル行列はそのまま
			positions[i].Translation(), // 移動
			rotations[i].Quaternion(),  // 回転
			scales[i].Scaling(),        // 拡大率
		))
	}

	return boneDeltas
}

// アニメーション対象ボーン一覧取得
func (bfs *BoneFrames) getAnimatedBoneNames(
	model *pmx.PmxModel,
	boneNames []string,
) (map[string]int, map[int]string) {
	// ボーン名の存在チェック用マップ
	exists := make(map[string]struct{})

	// 条件分岐の最適化
	if len(boneNames) > 0 {
		for _, boneName := range boneNames {
			// ボーン名の追加
			exists[boneName] = struct{}{}

			// 関連するボーンの追加
			relativeBoneIndexes := model.Bones.GetItemByName(boneName).RelativeBoneIndexes
			for _, index := range relativeBoneIndexes {
				relativeBoneName := model.Bones.GetItem(index).Name
				exists[relativeBoneName] = struct{}{}
			}
		}

		resultBoneNames := make(map[string]int)
		resultBoneIndexes := make(map[int]string)

		// 変形階層・ボーンINDEXでソート
		n := 0
		for _, boneIndex := range model.Bones.GetLayerIndexes() {
			bone := model.Bones.GetItem(boneIndex)
			if _, ok := exists[bone.Name]; ok {
				resultBoneNames[bone.Name] = n
				resultBoneIndexes[n] = bone.Name
				n++
			}
		}

		return resultBoneNames, resultBoneIndexes
	}

	// 全ボーンが対象の場合
	return model.Bones.LayerSortedNames, model.Bones.LayerSortedIndexes
}

// ボーン変形行列を求める
func (bfs *BoneFrames) getBoneMatrixes(
	frame float32,
	model *pmx.PmxModel,
	targetBoneNames map[string]int,
	targetBoneIndexes map[int]string,
	isCalcMorph bool,
) ([]*mmath.MMat4, []*mmath.MMat4, []*mmath.MMat4, []*mmath.MQuaternion) {
	boneCount := len(targetBoneNames)
	positions := make([]*mmath.MMat4, boneCount)
	rotations := make([]*mmath.MMat4, boneCount)
	scales := make([]*mmath.MMat4, boneCount)
	quats := make([]*mmath.MQuaternion, boneCount)

	for i := 0; i < boneCount; i++ {
		positions = append(positions, mmath.NewMMat4())
		rotations = append(rotations, mmath.NewMMat4())
		scales = append(scales, mmath.NewMMat4())
		quats = append(quats, mmath.NewMQuaternion())
	}

	// ボーンを一定件数ごとに並列処理
	count := 100

	var boneWg sync.WaitGroup
	for i := 0; i < boneCount; i += count {
		boneWg.Add(1)
		go func(i int) {
			defer boneWg.Done()
			for j := i; j < i+count; j++ {
				if j >= boneCount {
					break
				}
				boneName := targetBoneIndexes[j]
				// ボーンの移動位置、回転角度、拡大率を取得
				positions[j] = bfs.getPosition(frame, boneName, model, isCalcMorph, 0)
				rotWithEffect, rotFk := bfs.getRotation(frame, boneName, model, isCalcMorph, 0)
				rotations[j] = rotWithEffect.ToMat4()
				quats[j] = rotFk
				scales[j] = bfs.getScale(frame, boneName, model, isCalcMorph)
			}
		}(i)
	}
	boneWg.Wait()

	return positions, rotations, scales, quats
}

// 該当キーフレにおけるボーンの移動位置
func (bfs *BoneFrames) getPosition(
	frame float32,
	boneName string,
	model *pmx.PmxModel,
	isCalcMorph bool,
	loop int,
) *mmath.MMat4 {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMMat4()
	}

	bone := model.Bones.GetItemByName(boneName)
	bf := bfs.GetItem(boneName).GetItem(frame)

	mat := mmath.NewMMat4()
	if isCalcMorph {
		mat.Mul(bf.MorphPosition.ToMat4())
	}
	mat.Mul(bf.Position.ToMat4())

	if bone.IsEffectorTranslation() {
		// 外部親変形ありの場合、外部親変形行列を掛ける
		effectPosMat := bfs.getPositionWithEffect(frame, bone.Index, model, isCalcMorph, loop+1)
		mat.Mul(effectPosMat)
	}

	return mat
}

// 付与親を加味した移動位置
func (bfs *BoneFrames) getPositionWithEffect(
	frame float32,
	boneIndex int,
	model *pmx.PmxModel,
	isCalcMorph bool,
	loop int,
) *mmath.MMat4 {
	bone := model.Bones.GetItem(boneIndex)

	if bone.EffectFactor == 0 || loop > 20 {
		// 付与率が0の場合、常に0になる
		// MMDエンジン対策で無限ループを避ける
		return mmath.NewMMat4()
	}

	if !(bone.EffectIndex > 0 && model.Bones.Contains(bone.EffectIndex)) {
		// 付与親が存在しない場合、常に0になる
		return mmath.NewMMat4()
	}

	// 付与親が存在する場合、付与親の回転角度を掛ける
	effectBone := model.Bones.GetItem(bone.EffectIndex)
	posMat := bfs.getPosition(frame, effectBone.Name, model, isCalcMorph, loop+1)

	posMat[0][3] *= bone.EffectFactor
	posMat[1][3] *= bone.EffectFactor
	posMat[2][3] *= bone.EffectFactor

	return posMat
}

// 該当キーフレにおけるボーンの回転角度
func (bfs *BoneFrames) getRotation(
	frame float32,
	boneName string,
	model *pmx.PmxModel,
	isCalcMorph bool,
	loop int,
) (*mmath.MQuaternion, *mmath.MQuaternion) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMQuaternion(), mmath.NewMQuaternion()
	}

	bone := model.Bones.GetItemByName(boneName)

	// FK(捩り) > IK(捩り) > 付与親(捩り)
	bf := bfs.GetItem(boneName).GetItem(frame)
	var rot *mmath.MQuaternion
	if bf.IkRotation != nil && !bf.IkRotation.GetRadians().IsZero() {
		// IK用回転を持っている場合、置き換え
		if isCalcMorph {
			rot = bf.MorphRotation.GetQuaternion().Copy()
			rot.Mul(bf.IkRotation.GetQuaternion())
		} else {
			rot = bf.IkRotation.GetQuaternion().Copy()
		}
	} else {
		if isCalcMorph {
			rot = bf.MorphRotation.GetQuaternion().Copy()
			rot.Mul(bf.Rotation.GetQuaternion())
		} else {
			rot = bf.Rotation.GetQuaternion().Copy()
		}

		if bone.HasFixedAxis() {
			rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
		}
	}

	var rotWithEffect *mmath.MQuaternion
	if bone.IsEffectorRotation() {
		// 外部親変形ありの場合、外部親変形行列を掛ける
		effectQ := rot.Muled(bfs.getRotationWithEffect(frame, bone.Index, model, isCalcMorph, loop+1))
		rotWithEffect = effectQ
	} else {
		rotWithEffect = rot
	}

	if bone.HasFixedAxis() {
		// 軸制限回転を求める
		rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
	}

	return rotWithEffect, rot
}

// 付与親を加味した回転角度
func (bfs *BoneFrames) getRotationWithEffect(
	frame float32,
	boneIndex int,
	model *pmx.PmxModel,
	isCalcMorph bool,
	loop int,
) *mmath.MQuaternion {
	bone := model.Bones.GetItem(boneIndex)

	if bone.EffectFactor == 0 || loop > 20 {
		// 付与率が0の場合、常に0になる
		// MMDエンジン対策で無限ループを避ける
		return mmath.NewMQuaternion()
	}

	if !(bone.EffectIndex > 0 && model.Bones.Contains(bone.EffectIndex)) {
		// 付与親が存在しない場合、常に0になる
		return mmath.NewMQuaternion()
	}

	// 付与親が存在する場合、付与親の回転角度を掛ける
	effectBone := model.Bones.GetItem(bone.EffectIndex)
	rotWithEffect, _ := bfs.getRotation(frame, effectBone.Name, model, isCalcMorph, loop+1)

	if bone.EffectFactor >= 0 {
		// 正の付与親
		effectQ := rotWithEffect.MulFactor(bone.EffectFactor)
		return effectQ
	} else {
		// 負の付与親の場合、逆回転
		effectQ := rotWithEffect.MulFactor(-bone.EffectFactor)
		effectQ.Invert()
		return effectQ
	}
}

// 該当キーフレにおけるボーンの拡大率
func (bfs *BoneFrames) getScale(
	frame float32,
	boneName string,
	model *pmx.PmxModel,
	isCalcMorph bool,
) *mmath.MMat4 {
	bf := bfs.GetItem(boneName).GetItem(frame)
	mat := mmath.NewMMat4()

	if isCalcMorph {
		mat.ScaleVec3(bf.MorphScale.AddedScalar(1))
	}
	mat.ScaleVec3(bf.Scale.AddedScalar(1))

	return mat
}
