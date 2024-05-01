package vmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/miu200521358/mlib_go/pkg/deform"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type BoneFrames struct {
	Data map[string]*BoneNameFrames
	lock sync.RWMutex // マップアクセス制御用
}

func NewBoneFrames() *BoneFrames {
	return &BoneFrames{
		Data: make(map[string]*BoneNameFrames, 0),
	}
}

func (bfs *BoneFrames) Contains(boneName string) bool {
	bfs.lock.RLock()
	defer bfs.lock.RUnlock()

	_, ok := bfs.Data[boneName]
	return ok
}

func (bfs *BoneFrames) Append(bnfs *BoneNameFrames) {
	bfs.lock.Lock()
	defer bfs.lock.Unlock()

	bfs.Data[bnfs.Name] = bnfs
}

func (bfs *BoneFrames) GetItem(boneName string) *BoneNameFrames {
	if !bfs.Contains(boneName) {
		bfs.Append(NewBoneNameFrames(boneName))
	}

	bfs.lock.RLock()
	defer bfs.lock.RUnlock()

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

func (bfs *BoneFrames) GetMinFrame() float32 {
	minFno := float32(math.MaxFloat32)
	for _, bnfs := range bfs.Data {
		fno := bnfs.GetMaxFrame()
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (bfs *BoneFrames) Animate(
	frame float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isCalcMorph bool,
) *deform.BoneDeltas {
	// 処理対象ボーン一覧取得
	targetBoneNames, targetBoneIndexes := bfs.getAnimatedBoneNames(model, boneNames)

	// IK事前計算
	if isCalcIk {
		bfs.prepareIkSolvers(frame, model, targetBoneNames, isCalcMorph)
	}

	// ボーン変形行列操作
	positions, rotations, scales, quatsWithoutEffect :=
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
		quatsWithoutEffect,
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

	mlog.V("[IK計算終了][%.2f] -----------------------------------------------", frame)
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
	positions, rotations, scales, quatsWithoutEffect :=
		bfs.getBoneMatrixes(frame, model, effectorTargetBoneNames, effectorTargetBoneIndexes, isisCalcMorph)
	// 中断FLGが入ったか否か
	aborts := make([]bool, len(ikBone.Ik.Links))

	var ikFile *os.File
	var ikMotion *VmdMotion
	count := float32(1.0)
	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		dirPath := fmt.Sprintf("%s/IK_step", filepath.Dir(model.Path))
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		date := time.Now().Format("20060102_150405")
		ikMotionPath := fmt.Sprintf("%s/%.2f_%s_%s.vmd", dirPath, frame, date, ikBone.Name)
		ikMotion = NewVmdMotion(ikMotionPath)

		ikLogPath := fmt.Sprintf("%s/%.2f_%s_%s.log", dirPath, frame, date, ikBone.Name)
		ikFile, err = os.OpenFile(ikLogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(ikFile, "----------------------------------------")
		fmt.Println(ikFile, "[IK計算出力先][%.2f][%s] %s", frame, ikMotionPath)
	}
	defer func() {
		mlog.V("[IK計算終了][%.2f][%s]", frame, ikBone.Name)

		if ikMotion != nil {
			Write(ikMotion)
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

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
			unitRad := ikBone.Ik.UnitRotation.GetRadians().GetX() * float64(lidx+1)

			// IK関連の行列を取得
			linkMatrixes := bfs.calcBoneMatrixes(
				frame,
				model,
				effectorTargetBoneNames,
				effectorTargetBoneIndexes,
				positions,
				rotations,
				scales,
				quatsWithoutEffect,
			)

			// IKボーンのグローバル位置
			ikGlobalPosition := ikMatrixes.GetItem(ikBone.Name, frame).Position

			// 現在のIKターゲットボーンのグローバル位置を取得
			effectorGlobalPosition := linkMatrixes.GetItem(effectorBone.Name, frame).Position

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.2f][%03d][%s][%05.0f][00][Global] ikGlobalPosition: %s, effectorGlobalPosition: %s\n",
					frame, loop, linkBone.Name, count-1,
					ikGlobalPosition.String(), effectorGlobalPosition.String())
			}

			// fmt.Fprintf(ikFile,
			// 	"[%.2f][%03d][%s][%05.0f][01][グローバル位置終了判定] %sと%sの距離: %v(%0.6f)\n",
			// 	frame, loop, linkBone.Name, count-1, ikBone.Name, effectorBone.Name,
			// 	ikGlobalPosition.Distance(effectorGlobalPosition) < 1e-6,
			// 	ikGlobalPosition.Distance(effectorGlobalPosition))

			// // 位置の差がほとんどない場合、終了
			// if ikGlobalPosition.Distance(effectorGlobalPosition) < 1e-6 {
			// 	break ikLoop
			// }

			// 注目ノード（実際に動かすボーン=リンクボーン）
			linkMatrix := linkMatrixes.GetItem(linkBone.Name, frame).GlobalMatrix
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := linkMatrix.Inverse()

			// 注目ノードを起点とした、エフェクタのローカル位置
			effectorLocalPosition := linkInvMatrix.MulVec3(effectorGlobalPosition)
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				{
					bf := deform.NewBoneFrame(count)
					bf.Position = ikMatrixes.GetItem(ikBone.Name, frame).FramePosition
					bf.Rotation.SetQuaternion(ikMatrixes.GetItem(ikBone.Name, frame).FrameRotation)
					ikMotion.AppendRegisteredBoneFrame(ikBone.Name, bf)
					count++
				}
				{
					bf := deform.NewBoneFrame(count)
					bf.Position = linkMatrixes.GetItem(linkBone.Name, frame).FramePosition
					bf.Rotation.SetQuaternion(linkMatrixes.GetItem(linkBone.Name, frame).FrameRotation)
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}
			}
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.2f][%03d][%s][%05.0f][00][Local] effectorLocalPosition: %s, ikLocalPosition: %s\n",
					frame, loop, linkBone.Name, count-1,
					effectorLocalPosition.String(), ikLocalPosition.String())
			}

			effectorLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			// ベクトル (1) を (2) に一致させるための最短回転量（Axis-Angle）
			// 回転軸
			linkAxis := effectorLocalPosition.Cross(ikLocalPosition).Normalize()
			// 回転角(ラジアン)
			linkAngle := math.Acos(mmath.ClampFloat(effectorLocalPosition.Dot(ikLocalPosition), -1, 1))

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.2f][%03d][%s][%05.0f][01][回転軸・角度] linkAxis: %s, linkAngle: %.5f\n",
					frame, loop, linkBone.Name, count-1, linkAxis.String(), mmath.ToDegree(linkAngle),
				)
			}

			fmt.Fprintf(ikFile,
				"[%.2f][%03d][%s][%05.0f][01][回転角度終了判定] linkAngle: %v(%0.6f)\n",
				frame, loop, linkBone.Name, count-1, linkAngle < 1e-6, linkAngle)

			// 角度がほとんどない場合、終了
			if linkAngle < 1e-7 || (linkBone.HasFixedAxis() && linkAngle < 1e-2) {
				break ikLoop
			}

			// 単位角を超えないようにする
			linkAngle = mmath.ClampFloat(linkAngle, -unitRad, unitRad)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.2f][%03d][%s][%05.0f][02][単位角制限] linkAngle: %.5f\n",
					frame, loop, linkBone.Name, count-1, mmath.ToDegree(linkAngle),
				)
			}

			// リンクボーンの角度を取得
			linkQuat := quatsWithoutEffect[linkIndex]

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := deform.NewBoneFrame(count)
				bf.Rotation.SetQuaternion(linkQuat)
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.2f][%03d][%s][%05.0f][03][linkQuat] linkQuat: %s(%s)\n",
					frame, loop, linkBone.Name, count-1, linkQuat.String(), linkQuat.ToDegrees().String(),
				)
			}

			var totalActualIkQuat *mmath.MQuaternion
			if ikLink.AngleLimit || ikLink.LocalAngleLimit {
				// 角度制限が入ってる場合
				if ikLink.MinAngleLimit.GetRadians().GetX() != 0 ||
					ikLink.MaxAngleLimit.GetRadians().GetX() != 0 ||
					ikLink.LocalMinAngleLimit.GetRadians().GetX() != 0 ||
					ikLink.LocalMaxAngleLimit.GetRadians().GetX() != 0 {
					axisVector := &mmath.MVec3{1, 0, 0}
					if ikLink.LocalAngleLimit {
						axisVector = ikBone.NormalizedLocalAxisX
					}
					// グローバルX or ローカルX
					totalActualIkQuat = bfs.calcSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetX(),
						ikLink.MaxAngleLimit.GetRadians().GetX(),
						linkQuat, linkAxis, linkAngle, 0, axisVector,
						frame, count-1, loop, linkBone.Name, ikMotion, ikFile)
				} else if ikLink.MinAngleLimit.GetRadians().GetY() != 0 ||
					ikLink.MaxAngleLimit.GetRadians().GetY() != 0 ||
					ikLink.LocalMinAngleLimit.GetRadians().GetY() != 0 ||
					ikLink.LocalMaxAngleLimit.GetRadians().GetY() != 0 {
					axisVector := &mmath.MVec3{0, 1, 0}
					if ikLink.LocalAngleLimit {
						axisVector = ikBone.NormalizedLocalAxisY
					}
					// グローバルY or ローカルY
					totalActualIkQuat = bfs.calcSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetY(),
						ikLink.MaxAngleLimit.GetRadians().GetY(),
						linkQuat, linkAxis, linkAngle, 1, axisVector,
						frame, count-1, loop, linkBone.Name, ikMotion, ikFile)
				} else if ikLink.MinAngleLimit.GetRadians().GetZ() != 0 ||
					ikLink.MaxAngleLimit.GetRadians().GetZ() != 0 ||
					ikLink.LocalMinAngleLimit.GetRadians().GetZ() != 0 ||
					ikLink.LocalMaxAngleLimit.GetRadians().GetZ() != 0 {
					axisVector := &mmath.MVec3{0, 0, 1}
					if ikLink.LocalAngleLimit {
						axisVector = ikBone.NormalizedLocalAxisZ
					}
					// グローバルZ or ローカルZ
					totalActualIkQuat = bfs.calcSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetZ(),
						ikLink.MaxAngleLimit.GetRadians().GetZ(),
						linkQuat, linkAxis, linkAngle, 2, axisVector,
						frame, count-1, loop, linkBone.Name, ikMotion, ikFile)
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := deform.NewBoneFrame(count)
					bf.Rotation.SetQuaternion(totalActualIkQuat)
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%.2f][%03d][%s][%05.0f][11][角度制限結果] totalActualIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToDegrees().String(),
					)
				}
			} else {
				if linkBone.HasFixedAxis() {
					if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
						quat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle).Shorten()
						bf := deform.NewBoneFrame(count)
						bf.Rotation.SetQuaternion(quat)
						ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
						count++

						if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
							fmt.Fprintf(ikFile,
								"[%.2f][%03d][%s][%05.0f][04][軸制限][理想回転] quat: %s(%s)\n",
								frame, loop, linkBone.Name, count-1, quat.String(), quat.ToDegrees().String(),
							)
						}
					}

					// 軸制限ありの場合、軸にそった理想回転量とする
					linkAxis = linkBone.NormalizedFixedAxis

					if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
						quat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle).Shorten()
						bf := deform.NewBoneFrame(count)
						bf.Rotation.SetQuaternion(quat)
						ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
						count++

						if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
							fmt.Fprintf(ikFile,
								"[%.2f][%03d][%s][%05.0f][04][軸制限][理想軸制限回転] quat: %s(%s)\n",
								frame, loop, linkBone.Name, count-1, quat.String(), quat.ToDegrees().String(),
							)
						}
					}
				}

				correctIkQuat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle).Shorten()

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := deform.NewBoneFrame(count)
					bf.Rotation.SetQuaternion(correctIkQuat)
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%.2f][%03d][%s][%05.0f][13][角度制限なし] correctIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, correctIkQuat.String(), correctIkQuat.ToDegrees().String())
				}

				// 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
				totalActualIkQuat = linkQuat.Muled(correctIkQuat)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := deform.NewBoneFrame(count)
					bf.Rotation.SetQuaternion(totalActualIkQuat)
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%.2f][%03d][%s][%05.0f][14][角度制限なし] totalActualIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToDegrees().String())
				}
			}

			if linkBone.HasFixedAxis() {
				// 軸制限回転を求める
				totalActualIkQuat = totalActualIkQuat.ToFixedAxisRotation(linkBone.NormalizedFixedAxis)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := deform.NewBoneFrame(count)
					bf.Rotation.SetQuaternion(totalActualIkQuat)
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%.2f][%03d][%s][%05.0f][15][軸制限後処理] totalActualIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToDegrees().String())
				}
			}

			fmt.Fprintf(ikFile,
				"[%.2f][%03d][%s][%05.0f][15] 前回差分中断判定: %v(%0.6f) 前回: %s 今回: %s\n",
				frame, loop, linkBone.Name, count-1,
				1-quatsWithoutEffect[linkIndex].Dot(totalActualIkQuat) < 1e-6, 1-quatsWithoutEffect[linkIndex].Dot(totalActualIkQuat),
				quatsWithoutEffect[linkIndex].ToDegrees().String(), totalActualIkQuat.ToDegrees().String())

			// 前回（既存）とほぼ同じ回転量の場合、中断FLGを立てる
			if 1-quatsWithoutEffect[linkIndex].Dot(totalActualIkQuat) < 1e-8 {
				aborts[lidx] = true
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := deform.NewBoneFrame(count)
				bf.Rotation.SetQuaternion(totalActualIkQuat)
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%.2f][%03d][%s][%05.0f][16][結果] totalActualIkQuat: %s(%s)\n",
					frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToDegrees().String())
			}

			// IKの結果を更新
			quatsWithoutEffect[linkIndex] = totalActualIkQuat
			rotations[linkIndex] = totalActualIkQuat.ToMat4()
		}

		// すべてのリンクボーンで中断FLG = ONの場合、ループ終了
		if slices.Index(aborts, false) == -1 {
			break
		}
	}

	return quatsWithoutEffect, effectorTargetBoneNames
}

// 全ての角度をラジアン角度に分割して、そのうちのひとつの軸だけを動かす回転を取得する
// minAngleLimit: 最小軸制限（ラジアン）
// maxAngleLimit: 最大軸制限（ラジアン）
// linkQuat: 現在のリンクボーンの回転量
// quatAxis: 現在のIK回転の回転軸
// quatAngle: 現在のIK回転の回転角度（ラジアン）
// axisIndex: 制限軸INDEX
func (bfs *BoneFrames) calcSingleAxisRad(
	minAngleLimit float64,
	maxAngleLimit float64,
	linkQuat *mmath.MQuaternion,
	quatAxis *mmath.MVec3,
	quatAngle float64,
	axisIndex int,
	axisVector *mmath.MVec3,
	frame float32,
	count float32,
	loop int,
	linkBoneName string,
	ikMotion *VmdMotion,
	ikFile *os.File,
) *mmath.MQuaternion {
	quat := mmath.NewMQuaternionFromAxisAngles(quatAxis, quatAngle).Shorten()

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		bf := deform.NewBoneFrame(count)
		bf.Rotation.SetQuaternion(quat)
		ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
		count++
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][04][角度制限] quat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, quat.String(), quat.ToDegrees().String())
	}

	// 現在IKリンクに入る可能性のあるすべての角度
	totalIkQuat := linkQuat.Muled(quat)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		bf := deform.NewBoneFrame(count)
		bf.Rotation.SetQuaternion(totalIkQuat)
		ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
		count++
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][05][角度制限] totalIkQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, totalIkQuat.String(), totalIkQuat.ToDegrees().String())
	}

	totalIkRad := totalIkQuat.ToRadian()
	// TODO ローカル軸ベースの分割の場合、ローカル軸に合わせる
	if quatAxis.Dot(axisVector) < 0 {
		totalIkRad *= -1
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][06][角度制限] totalIkRad: %.5f(%.5f)\n",
			frame, loop, linkBoneName, count-1, totalIkRad, mmath.ToDegree(totalIkRad))
	}

	fSX := math.Sin(totalIkRad) // sin(θ)
	fX := math.Asin(fSX)        // 一軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][07][角度制限] fSX: %.5f, fX: %.5f(%.5f)\n",
			frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
	}

	// ジンバルロック回避
	totalIkRads, isGimbal := totalIkQuat.ToRadiansWithGimbal(axisIndex)
	if isGimbal || math.Abs(totalIkRad) > math.Pi {
		fX = totalIkRads.Vector()[axisIndex]
		if fX < 0 {
			fX = -(math.Pi - fX)
		} else {
			fX = math.Pi - fX
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][08][角度制限][ジンバルロック回避] fSX: %.5f, fX: %.5f(%.5f)\n",
				frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
		}
	}

	// 角度の制限
	if fX < minAngleLimit {
		tf := 2*minAngleLimit - fX

		fX = mmath.ClampFloat(tf, minAngleLimit, maxAngleLimit)

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][09][角度制限][負角度制限] fSX: %.5f, fX: %.5f(%.5f)\n",
				frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
		}
	}
	if fX > maxAngleLimit {
		tf := 2*maxAngleLimit - fX

		fX = mmath.ClampFloat(tf, minAngleLimit, maxAngleLimit)

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%.2f][%03d][%s][%05.0f][10][角度制限][正角度制限] fSX: %.5f, fX: %.5f(%.5f)\n",
				frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
		}
	}

	return mmath.NewMQuaternionFromAxisAngles(axisVector, fX).Shorten()
}

func (bfs *BoneFrames) calcBoneMatrixes(
	frame float32,
	model *pmx.PmxModel,
	targetBoneNames map[string]int,
	targetBoneIndexes map[int]string,
	positions, rotations, scales []*mmath.MMat4,
	quatsWithoutEffect []*mmath.MQuaternion,
) *deform.BoneDeltas {
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

	boneDeltas := deform.NewBoneDeltas()

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
		boneDeltas.SetItem(bone.Name, frame, deform.NewBoneDelta(
			bone.Name,
			frame,
			localMatrix.Muled(bone.Position.ToMat4()), // グローバル行列
			localMatrix,                // ローカル行列はそのまま
			positions[i].Translation(), // 移動
			rotations[i].Quaternion(),  // 回転
			quatsWithoutEffect[i],      // 回転(付与親なし)
			scales[i].Scaling(),        // 拡大率
			matrixes[i],                // ボーン変形行列
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

	return rotWithEffect.Shorten(), rot.Shorten()
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

	return rotWithEffect.MulScalar(bone.EffectFactor).Shorten()
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
