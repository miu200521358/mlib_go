package vmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/miu200521358/mlib_go/pkg/mcore"
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
		lock: sync.RWMutex{},
	}
}

func (fs *BoneFrames) Contains(boneName string) bool {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	_, ok := fs.Data[boneName]
	return ok
}

func (fs *BoneFrames) Append(nfs *BoneNameFrames) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Data[nfs.Name] = nfs
}

func (fs *BoneFrames) Delete(boneName string) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	delete(fs.Data, boneName)
}

func (fs *BoneFrames) Get(boneName string) *BoneNameFrames {
	if !fs.Contains(boneName) {
		fs.Append(NewBoneNameFrames(boneName))
	}

	fs.lock.RLock()
	defer fs.lock.RUnlock()

	return fs.Data[boneName]
}

func (fs *BoneFrames) GetNames() []string {
	names := make([]string, 0, len(fs.Data))
	for name := range fs.Data {
		names = append(names, name)
	}
	return names
}

func (fs *BoneFrames) GetIndexes() []int {
	indexes := mcore.NewIntIndexes()
	for _, fs := range fs.Data {
		for _, f := range fs.List() {
			indexes.ReplaceOrInsert(f.Index)
		}
	}
	return indexes.List()
}

func (fs *BoneFrames) GetRegisteredIndexes() []int {
	indexes := mcore.NewIntIndexes()
	for _, fs := range fs.Data {
		for _, index := range fs.RegisteredIndexes.List() {
			indexes.ReplaceOrInsert(mcore.NewInt(index))
		}
	}
	return indexes.List()
}

func (fs *BoneFrames) Len() int {
	count := 0
	for _, fs := range fs.Data {
		count += fs.RegisteredIndexes.Len()
	}
	return count
}

func (fs *BoneFrames) GetMaxFrame() int {
	maxFno := int(0)
	for _, fs := range fs.Data {
		fno := fs.GetMaxFrame()
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (fs *BoneFrames) GetMinFrame() int {
	minFno := math.MaxInt
	for _, fs := range fs.Data {
		fno := fs.GetMinFrame()
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (fs *BoneFrames) Deform(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	beforeBoneDeltas *BoneDeltas,
) *BoneDeltas {
	boneDeformsMap := fs.prepareDeform(frame, model, boneNames, isCalcIk, beforeBoneDeltas)
	return fs.calcBoneDeltas(frame, boneDeformsMap, beforeBoneDeltas)
}

func (fs *BoneFrames) prepareDeform(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	beforeBoneDeltas *BoneDeltas,
) map[bool]*boneDeforms {
	isAfterPhysics := false
	if beforeBoneDeltas != nil {
		isAfterPhysics = true
	}

	boneDeformsMap := fs.prepareBoneDeforms(model, boneNames, isAfterPhysics)

	// IK事前計算
	if isCalcIk {
		// ボーン変形行列操作
		fs.prepareIk(frame, model, boneDeformsMap, beforeBoneDeltas)
	}

	// ボーンデフォーム情報を埋める
	fs.fillBoneDeform(frame, model, boneDeformsMap, isAfterPhysics)

	return boneDeformsMap
}

// IK事前計算処理
func (fs *BoneFrames) clearIk(
	frame int,
	model *pmx.PmxModel,
	boneDeforms *boneDeforms,
) {
	for _, boneDeform := range boneDeforms.deforms {
		// ボーンIndexがIkTreeIndexesに含まれていない場合、スルー
		if _, ok := model.Bones.IkTreeIndexes[boneDeform.bone.Index]; !ok {
			continue
		}

		for i := 0; i < len(model.Bones.IkTreeIndexes[boneDeform.bone.Index]); i++ {
			ikBone := model.Bones.Get(model.Bones.IkTreeIndexes[boneDeform.bone.Index][i])
			if _, ok := boneDeforms.names[ikBone.Name]; !ok {
				continue
			}

			for _, linkIndex := range ikBone.Ik.Links {
				// IKリンクボーンの回転量を初期化
				linkBone := model.Bones.Get(linkIndex.BoneIndex)
				linkBf := fs.Get(linkBone.Name).data[frame]
				if linkBf != nil {
					linkBf.IkRotation = nil

					// IK用なので登録フラグは既存のままで追加して補間曲線は分割しない
					fs.Get(linkBone.Name).Append(linkBf)
				}
			}
		}
	}
}

// IK事前計算処理
func (fs *BoneFrames) prepareIk(
	frame int,
	model *pmx.PmxModel,
	boneDeformsMap map[bool]*boneDeforms,
	beforeBoneDeltas *BoneDeltas,
) {
	isAfterPhysicsList := make([]bool, 0, 2)
	isAfterPhysicsList = append(isAfterPhysicsList, false)
	if beforeBoneDeltas != nil {
		isAfterPhysicsList = append(isAfterPhysicsList, true)
	}

	for _, isAfterPhysics := range isAfterPhysicsList {
		// IKクリア
		fs.clearIk(frame, model, boneDeformsMap[isAfterPhysics])
	}

	for _, isAfterPhysics := range isAfterPhysicsList {
		boneDeforms := boneDeformsMap[isAfterPhysics]

		for _, boneIndex := range boneDeforms.boneIndexes {
			bd := boneDeforms.deforms[boneIndex]

			// ボーンIndexがIkTreeIndexesに含まれていない場合、スルー
			if _, ok := model.Bones.IkTreeIndexes[bd.bone.Index]; !ok {
				continue
			}

			for m := range len(model.Bones.IkTreeIndexes[bd.bone.Index]) {
				ikBone := model.Bones.Get(model.Bones.IkTreeIndexes[bd.bone.Index][m])
				if _, ok := boneDeforms.names[ikBone.Name]; !ok {
					continue
				}

				// IK計算
				effectorBoneDeltas := fs.calcIk(frame, ikBone, model, beforeBoneDeltas)

				for _, linkIndex := range ikBone.Ik.Links {
					// IKリンクボーンの回転量を更新
					linkBone := model.Bones.Get(linkIndex.BoneIndex)
					effectorBoneDelta := effectorBoneDeltas.Get(linkBone.Index)
					if effectorBoneDelta != nil {
						linkBf := fs.Get(linkBone.Name).Get(frame)
						linkBf.IkRotation = effectorBoneDelta.FrameRotation()
						// IK用なので登録フラグは既存のままで追加して補間曲線は分割しない
						fs.Get(linkBone.Name).Append(linkBf)
					}
				}
			}
		}
	}

	mlog.V("[IK計算終了][%04d] -----------------------------------------------", frame)
}

// IK計算
func (fs *BoneFrames) calcIk(
	frame int,
	ikBone *pmx.Bone,
	model *pmx.PmxModel,
	beforeBoneDeltas *BoneDeltas,
) *BoneDeltas {
	boneDeltas := NewBoneDeltas()
	isAfterPhysics := beforeBoneDeltas != nil
	// IKターゲットボーン
	effectorBone := model.Bones.Get(ikBone.Ik.BoneIndex)
	// IK関連の行列を一括計算
	ikDeltas := fs.Deform(frame, model, []string{ikBone.Name}, false, beforeBoneDeltas)
	// エフェクタ関連情報取得
	effectorBoneDeformsMap := fs.prepareDeform(
		frame, model, []string{effectorBone.Name}, false, beforeBoneDeltas)

	var ikFile *os.File
	var ikMotion *VmdMotion
	count := int(1.0)
	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		dirPath := fmt.Sprintf("%s/IK_step", filepath.Dir(model.Path))
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		date := time.Now().Format("20060102_150405")
		ikMotionPath := fmt.Sprintf("%s/%04d_%s_%s.vmd", dirPath, frame, date, ikBone.Name)
		ikMotion = NewVmdMotion(ikMotionPath)

		ikLogPath := fmt.Sprintf("%s/%04d_%s_%s.log", dirPath, frame, date, ikBone.Name)
		ikFile, err = os.OpenFile(ikLogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(ikFile, "----------------------------------------")
		fmt.Println(ikFile, "[IK計算出力先][%04d][%s] %s", frame, ikMotionPath)
	}
	defer func() {
		mlog.V("[IK計算終了][%04d][%s]", frame, ikBone.Name)

		if ikMotion != nil {
			Write(ikMotion)
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

	loopLimitHalf := ikBone.Ik.LoopCount / 2

	// IK計算
ikLoop:
	for loop := 0; loop < max(ikBone.Ik.LoopCount, 1); loop++ {
		for lidx, ikLink := range ikBone.Ik.Links {
			// ikLink は末端から並んでる
			if !model.Bones.Contains(ikLink.BoneIndex) {
				continue
			}

			// 処理対象IKリンクボーン
			linkBone := model.Bones.Get(ikLink.BoneIndex)

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

			for _, l := range ikBone.Ik.Links {
				if effectorBoneDeformsMap[isAfterPhysics].deforms[l.BoneIndex] != nil {
					effectorBoneDeformsMap[isAfterPhysics].deforms[l.BoneIndex].unitMatrix = nil
				}
			}
			if effectorBoneDeformsMap[isAfterPhysics].deforms[ikBone.Ik.BoneIndex] != nil {
				effectorBoneDeformsMap[isAfterPhysics].deforms[ikBone.Ik.BoneIndex].unitMatrix = nil
			}

			// IK関連の行列を取得
			effectorDeltas := fs.calcBoneDeltas(frame, effectorBoneDeformsMap, beforeBoneDeltas)

			// IKボーンのグローバル位置
			ikGlobalPosition := ikDeltas.Get(ikBone.Index).GlobalPosition()

			// 現在のIKターゲットボーンのグローバル位置を取得
			effectorGlobalPosition := effectorDeltas.Get(effectorBone.Index).GlobalPosition()

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][00][Global] [%s]ikGlobalPosition: %s, "+
						"[%s]effectorGlobalPosition: %s, [%s]linkGlobalPosition: %s\n",
					frame, loop, linkBone.Name, count-1,
					ikBone.Name, ikGlobalPosition.MMD().String(),
					effectorBone.Name, effectorGlobalPosition.MMD().String(),
					linkBone.Name, effectorDeltas.Get(linkBone.Index).GlobalPosition().MMD().String())
			}

			// fmt.Fprintf(ikFile,
			// 	"[%04d][%03d][%s][%05d][01][グローバル位置終了判定] %sと%sの距離: %v(%0.6f)\n",
			// 	frame, loop, linkBone.Name, count-1, ikBone.Name, effectorBone.Name,
			// 	ikGlobalPosition.Distance(effectorGlobalPosition) < 1e-6,
			// 	ikGlobalPosition.Distance(effectorGlobalPosition))

			// // 位置の差がほとんどない場合、終了
			// if ikGlobalPosition.Distance(effectorGlobalPosition) < 1e-6 {
			// 	break ikLoop
			// }

			// 注目ノード（実際に動かすボーン=リンクボーン）
			linkMatrix := effectorDeltas.Get(linkBone.Index).GlobalMatrix()
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := linkMatrix.Inverse()

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][00][linkMatrix] linkMatrix: %s, linkMatrix.Inverse: %s\n",
					frame, loop, linkBone.Name, count-1,
					linkMatrix.String(), linkInvMatrix.String())
			}

			// 注目ノードを起点とした、エフェクタのローカル位置
			effectorLocalPosition := linkInvMatrix.MulVec3(effectorGlobalPosition)
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				{
					bf := NewBoneFrame(count)
					bf.Position = ikDeltas.Get(ikBone.Index).framePosition
					bf.Rotation = ikDeltas.Get(ikBone.Index).FrameRotation()
					ikMotion.AppendRegisteredBoneFrame(ikBone.Name, bf)
					count++
				}
				{
					bf := NewBoneFrame(count)
					bf.Position = effectorDeltas.Get(linkBone.Index).framePosition
					bf.Rotation = effectorDeltas.Get(linkBone.Index).FrameRotation()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}
			}
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][00][Local] effectorLocalPosition: %s, ikLocalPosition: %s\n",
					frame, loop, linkBone.Name, count-1,
					effectorLocalPosition.MMD().String(), ikLocalPosition.MMD().String())
			}

			effectorLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			// ベクトル (1) を (2) に一致させるための最短回転量（Axis-Angle）
			// 回転軸
			linkAxis := effectorLocalPosition.Cross(ikLocalPosition).Normalize()
			// 回転角(ラジアン)
			linkAngle := math.Acos(mmath.ClampFloat(ikLocalPosition.Dot(effectorLocalPosition), -1, 1))

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][01][回転軸・角度] linkAxis: %s, linkAngle: %.5f\n",
					frame, loop, linkBone.Name, count-1, linkAxis.MMD().String(), mmath.ToDegree(linkAngle),
				)

				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][01][回転角度終了判定] linkAngle: %v(%0.6f)\n",
					frame, loop, linkBone.Name, count-1, linkAngle < 1e-6, linkAngle)
			}

			// 角度がほとんどない場合、終了
			if linkAngle < 1e-7 {
				break ikLoop
			}

			// 単位角を超えないようにする
			linkAngle = mmath.ClampFloat(linkAngle, -unitRad, unitRad)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][02][単位角制限] linkAngle: %.5f\n",
					frame, loop, linkBone.Name, count-1, mmath.ToDegree(linkAngle),
				)
			}

			// リンクボーンの角度を取得
			linkDeform := getBoneDeform(effectorBoneDeformsMap, linkBone)
			var linkQuat *mmath.MQuaternion
			if linkDeform != nil && linkDeform.rotation != nil {
				linkQuat = linkDeform.rotation.Copy()
				if linkDeform.effectRotation != nil {
					linkQuat = linkDeform.rotation.Muled(linkDeform.effectRotation)
				}
			} else {
				linkQuat = mmath.NewMQuaternion()
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := NewBoneFrame(count)
				bf.Rotation = linkQuat
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][03][linkQuat] linkQuat: %s(%s)\n",
					frame, loop, linkBone.Name, count-1, linkQuat.String(), linkQuat.ToMMDDegrees().String(),
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
					totalActualIkQuat = fs.calcIkSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetX(),
						ikLink.MaxAngleLimit.GetRadians().GetX(),
						linkQuat, linkAxis, linkAngle, 0, axisVector, loopLimitHalf,
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
					totalActualIkQuat = fs.calcIkSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetY(),
						ikLink.MaxAngleLimit.GetRadians().GetY(),
						linkQuat, linkAxis, linkAngle, 1, axisVector, loopLimitHalf,
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
					totalActualIkQuat = fs.calcIkSingleAxisRad(
						ikLink.MinAngleLimit.GetRadians().GetZ(),
						ikLink.MaxAngleLimit.GetRadians().GetZ(),
						linkQuat, linkAxis, linkAngle, 2, axisVector, loopLimitHalf,
						frame, count-1, loop, linkBone.Name, ikMotion, ikFile)
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = totalActualIkQuat
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][11][角度制限結果] totalActualIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToMMDDegrees().String(),
					)
				}
			} else {
				if linkBone.HasFixedAxis() {
					if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
						quat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle).Shorten()
						bf := NewBoneFrame(count)
						bf.Rotation = quat
						ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
						count++

						if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
							fmt.Fprintf(ikFile,
								"[%04d][%03d][%s][%05d][04][軸制限][理想回転] quat: %s(%s)\n",
								frame, loop, linkBone.Name, count-1, quat.String(), quat.ToDegrees().String(),
							)
						}
					}

					if linkAxis.Dot(linkBone.NormalizedFixedAxis) < 0 {
						linkAngle = -linkAngle
					}

					// 軸制限ありの場合、軸にそった理想回転量とする
					linkAxis = linkBone.NormalizedFixedAxis

					if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
						quat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle).Shorten()
						bf := NewBoneFrame(count)
						bf.Rotation = quat
						ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
						count++

						if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
							fmt.Fprintf(ikFile,
								"[%04d][%03d][%s][%05d][04][軸制限][理想軸制限回転] quat: %s(%s)\n",
								frame, loop, linkBone.Name, count-1, quat.String(), quat.ToDegrees().String(),
							)
						}
					}
				}

				correctIkQuat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = correctIkQuat
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][13][角度制限なし] correctIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, correctIkQuat.String(), correctIkQuat.ToMMDDegrees().String())
				}

				// 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
				totalActualIkQuat = linkQuat.Muled(correctIkQuat)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = totalActualIkQuat
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][14][角度制限なし] totalActualIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToMMDDegrees().String())
				}
			}

			if linkBone.HasFixedAxis() {
				// 軸制限回転を求める
				totalActualIkQuat = totalActualIkQuat.ToFixedAxisRotation(linkBone.NormalizedFixedAxis)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = totalActualIkQuat
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++
				}

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][15][軸制限後処理] totalActualIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToMMDDegrees().String())
				}
			}

			// if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			// 	fmt.Fprintf(ikFile,
			// 		"[%04d][%03d][%s][%05d][15] 前回差分中断判定: %v(%0.6f) 前回: %s 今回: %s\n",
			// 		frame, loop, linkBone.Name, count-1,
			// 		1-quatsWithoutEffect[linkIndex].Dot(totalActualIkQuat) < 1e-6, 1-quatsWithoutEffect[linkIndex].Dot(totalActualIkQuat),
			// 		quatsWithoutEffect[linkIndex].ToDegrees().String(), totalActualIkQuat.ToMMDDegrees().String())
			// }

			// // 前回（既存）とほぼ同じ回転量の場合、中断FLGを立てる
			// if 1-quatsWithoutEffect[linkIndex].Dot(totalActualIkQuat) < 1e-5 {
			// 	// 反対側に曲げる
			// 	aborts[lidx] += 1
			// } else {
			// 	aborts[lidx] = 0
			// }

			// if slices.Index(aborts, 0) == -1 {
			// 	// すべてのリンクボーンで中断FLG > 0の場合、反対側に曲げる
			// 	totalActualIkQuat.Invert()
			// }

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := NewBoneFrame(count)
				bf.Rotation = totalActualIkQuat
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][16][結果] totalActualIkQuat: %s(%s)\n",
					frame, loop, linkBone.Name, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToMMDDegrees().String())
			}

			// IKの結果を更新
			linkDeform.rotation = totalActualIkQuat
			linkDeform.effectRotation = nil
			linkDeform.unitMatrix = nil
			var delta *BoneDelta
			if beforeBoneDeltas != nil {
				delta = beforeBoneDeltas.Get(linkBone.Index)
			}
			if delta == nil {
				delta = &BoneDelta{Bone: linkBone, Frame: frame}
			}
			delta.frameRotation = totalActualIkQuat
			boneDeltas.Append(delta)
		}
	}

	return boneDeltas
}

// calcIkSingleAxisRad は単一軸の制限を持つ回転クォータニオンを計算し、その結果を返します。
// minAngleLimit: 最小軸制限（ラジアン）
// maxAngleLimit: 最大軸制限（ラジアン）
// linkQuat: 現在のリンクボーンの回転量
// quatAxis: 現在のIK回転の回転軸
// quatAngle: 現在のIK回転の回転角度（ラジアン）
// axisIndex: 制限軸INDEX
// axisVector: 回転軸ベクトル
// frame, count, loop: デバッグ用カウンター
// linkBoneName: ボーン名
// ikMotion: IKモーション
// ikFile: IKファイル
func (fs *BoneFrames) calcIkSingleAxisRad(
	minAngleLimit float64,
	maxAngleLimit float64,
	linkQuat *mmath.MQuaternion,
	quatAxis *mmath.MVec3,
	quatAngle float64,
	axisIndex int,
	axisVector *mmath.MVec3,
	loopLimitHalf int,
	frame int,
	count int,
	loop int,
	linkBoneName string,
	ikMotion *VmdMotion,
	ikFile *os.File,
) *mmath.MQuaternion {
	quat := mmath.NewMQuaternionFromAxisAngles(quatAxis, quatAngle).Shorten()

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		bf := NewBoneFrame(count)
		bf.Rotation = quat
		ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
		count++
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][04][角度制限] quat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, quat.String(), quat.ToMMDDegrees().String())
	}

	// 現在IKリンクに入る可能性のあるすべての角度
	totalIkQuat := quat.Muled(linkQuat)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		bf := NewBoneFrame(count)
		bf.Rotation = totalIkQuat
		ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
		count++
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][05][角度制限] totalIkQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String())
	}

	totalIkRads, isGimbal := totalIkQuat.ToRadiansWithGimbal(axisIndex)
	linkRads, isLinkGimbal := quat.ToRadiansWithGimbal(axisIndex)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][ジンバルロック判定] totalIkRads: %s(%s), isGimbal: %v, linkRads: %s(%s), isLinkGimbal: %v\n",
			frame, loop, linkBoneName, count-1, totalIkRads.String(), totalIkQuat.ToDegrees().String(),
			isGimbal, linkRads.String(), quat.ToDegrees().String(), isLinkGimbal)
	}

	totalIkRad := totalIkQuat.ToRadian()
	// TODO ローカル軸ベースの分割の場合、ローカル軸に合わせる
	if totalIkQuat.GetXYZ().Dot(axisVector) < 0 {
		totalIkRad *= -1
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][06][角度制限] totalIkRad: %.5f(%.5f)\n",
			frame, loop, linkBoneName, count-1, totalIkRad, mmath.ToDegree(totalIkRad))
	}

	fSX := math.Sin(totalIkRad) // sin(θ)
	fX := math.Asin(fSX)        // 一軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限] fSX: %.5f, fX: %.5f(%.5f)\n",
			frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
	}

	if isGimbal || math.Abs(totalIkRad) > mmath.GIMBAL1_RAD {
		// ジンバルロック回避
		if fX < 0 {
			fX = -(math.Pi - fX)
		} else {
			fX = math.Pi - fX
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][08][角度制限][ジンバルロック回避] fSX: %.5f, fX: %.5f(%.5f)\n",
				frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
		}
	}

	// 角度の制限
	if fX < minAngleLimit {
		tf := 2*minAngleLimit - fX

		if tf <= maxAngleLimit {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][09][角度制限][負角度制限Over] fSX: %.5f, "+
					"fX: %.5f(%.5f), tf: %.5f(%.5f)\n",
					frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX), tf, mmath.ToDegree(tf))
			}

			fX = tf
		} else {
			fX = mmath.ClampFloat(fX, minAngleLimit, maxAngleLimit)
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][09][角度制限][負角度制限] fSX: %.5f, fX: %.5f(%.5f)\n",
				frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
		}
	} else if fX > maxAngleLimit {
		tf := 2*maxAngleLimit - fX

		if tf >= minAngleLimit {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][10][角度制限][正角度制限Over] fSX: %.5f, "+
					"fX: %.5f(%.5f), tf: %.5f(%.5f)\n",
					frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX), tf, mmath.ToDegree(tf))
			}

			fX = tf
		} else {
			fX = mmath.ClampFloat(fX, minAngleLimit, maxAngleLimit)
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][10][角度制限][正角度制限] fSX: %.5f, fX: %.5f(%.5f)\n",
				frame, loop, linkBoneName, count-1, fSX, fX, mmath.ToDegree(fX))
		}
	}

	return mmath.NewMQuaternionFromAxisAngles(axisVector, fX).Shorten()
}

func (fs *BoneFrames) calcBoneDeltas(
	frame int,
	boneDeformsMap map[bool]*boneDeforms,
	beforeBoneDeltas *BoneDeltas,
) *BoneDeltas {
	if beforeBoneDeltas != nil {
		return fs.calcBoneDeltasByIsAfterPhysics(frame, boneDeformsMap[true], beforeBoneDeltas)
	} else {
		return fs.calcBoneDeltasByIsAfterPhysics(frame, boneDeformsMap[false], NewBoneDeltas())
	}
}

func (fs *BoneFrames) calcBoneDeltasByIsAfterPhysics(
	frame int,
	boneDeforms *boneDeforms,
	boneDeltas *BoneDeltas,
) *BoneDeltas {
	for _, boneIndex := range boneDeforms.boneIndexes {
		deform := boneDeforms.deforms[boneIndex]
		if deform.unitMatrix != nil {
			// 既にグローバル行列が計算済みの場合、スルー
			continue
		}

		deform.unitMatrix = mmath.NewMMat4()

		// スケール
		if deform.scale != nil && !deform.scale.IsOne() {
			deform.unitMatrix.Mul(deform.scale.ToScaleMat4())
		}

		// 回転
		rot := deform.rotation
		isEffectorRot := false
		if deform.effectRotation != nil {
			rot.Mul(deform.effectRotation)
			isEffectorRot = true
		}
		if isEffectorRot && deform.bone.HasFixedAxis() {
			// 軸制限回転を求める
			rot = rot.ToFixedAxisRotation(deform.bone.NormalizedFixedAxis)
		}
		if rot != nil && !rot.IsIdent() {
			deform.unitMatrix.Mul(rot.ToMat4())
		}

		// 移動
		pos := deform.position
		if deform.effectPosition != nil && !deform.effectPosition.IsZero() {
			pos = pos.Add(deform.effectPosition)
		}
		if pos != nil && !pos.IsZero() {
			deform.unitMatrix.Mul(pos.ToMat4())
		}

		// 逆BOf行列(初期姿勢行列)
		deform.unitMatrix.Mul(deform.bone.RevertOffsetMatrix)
	}

	for _, boneIndex := range boneDeforms.boneIndexes {
		var globalMatrix *mmath.MMat4
		deform := boneDeforms.deforms[boneIndex]
		if _, ok := boneDeforms.deforms[deform.bone.ParentIndex]; ok {
			// 直近の親ボーンの変形行列を元にする
			parentDelta := boneDeltas.Get(deform.bone.ParentIndex)
			globalMatrix = deform.unitMatrix.Muled(parentDelta.GlobalMatrix())
		} else {
			// 対象ボーン自身の行列をかける
			globalMatrix = deform.unitMatrix.Copy()
		}

		// 初期位置行列を掛けてグローバル行列を作成
		boneDeltas.Append(NewBoneDelta(
			deform.bone,
			frame,
			globalMatrix,          // グローバル行列
			deform.position,       // キーフレの移動量
			deform.effectPosition, // キーフレの付与移動量
			deform.rotation,       // キーフレの回転量
			deform.effectRotation, // キーフレの付与回転量
			deform.scale,          // 拡大率
		))
	}

	return boneDeltas
}

// デフォーム対象ボーン情報一覧取得
func (fs *BoneFrames) prepareBoneDeforms(
	model *pmx.PmxModel,
	boneNames []string,
	isAfterPhysics bool,
) map[bool]*boneDeforms {
	boneDeformsMap := make(map[bool]*boneDeforms)
	isAfterPhysicsList := make([]bool, 0, 2)
	isAfterPhysicsList = append(isAfterPhysicsList, false)
	if isAfterPhysics {
		isAfterPhysicsList = append(isAfterPhysicsList, true)
	} else {
		boneDeformsMap[true] = &boneDeforms{
			deforms:     make(map[int]*BoneDeform),
			names:       make(map[string]int),
			boneIndexes: make([]int, 0),
		}
	}

	for _, ap := range isAfterPhysicsList {
		// ボーン名の存在チェック用マップ
		exists := make(map[string]struct{})
		boneDeforms := &boneDeforms{
			deforms:     make(map[int]*BoneDeform),
			names:       make(map[string]int),
			boneIndexes: make([]int, 0),
		}

		if len(boneNames) > 0 {
			for _, boneName := range boneNames {
				// ボーン名の追加
				exists[boneName] = struct{}{}

				// 関連するボーンの追加
				relativeBoneIndexes := model.Bones.GetByName(boneName).RelativeBoneIndexes
				for _, index := range relativeBoneIndexes {
					relativeBoneName := model.Bones.Get(index).Name
					exists[relativeBoneName] = struct{}{}
				}
			}
		}

		// 変形階層・ボーンINDEXでソート
		targetSortedBones := model.Bones.LayerSortedBones[ap]
		for k := range len(targetSortedBones) {
			bone := targetSortedBones[k]
			if len(boneNames) > 0 {
				if _, ok := exists[bone.Name]; ok {
					boneDeforms.deforms[bone.Index] = &BoneDeform{bone: bone}
					boneDeforms.names[bone.Name] = bone.Index
					boneDeforms.boneIndexes = append(boneDeforms.boneIndexes, bone.Index)
				}
			} else {
				boneDeforms.deforms[bone.Index] = &BoneDeform{bone: bone}
				boneDeforms.names[bone.Name] = bone.Index
				boneDeforms.boneIndexes = append(boneDeforms.boneIndexes, bone.Index)
			}
		}

		boneDeformsMap[ap] = boneDeforms
	}

	return boneDeformsMap
}

// デフォーム情報を求めて設定
func (fs *BoneFrames) fillBoneDeform(
	frame int,
	model *pmx.PmxModel,
	boneDeforms map[bool]*boneDeforms,
	isAfterPhysics bool,
) {
	for _, boneIndex := range boneDeforms[isAfterPhysics].boneIndexes {
		boneDeform := boneDeforms[isAfterPhysics].deforms[boneIndex]

		if boneDeform.unitMatrix == nil {
			bf := fs.Get(boneDeform.bone.Name).Get(frame)
			// ボーンの移動位置、回転角度、拡大率を取得
			boneDeform.position, boneDeform.effectPosition =
				fs.getPosition(bf, frame, boneDeform.bone, model, isAfterPhysics, 0)
			boneDeform.rotation, boneDeform.effectRotation =
				fs.getRotation(bf, frame, boneDeform.bone, model, isAfterPhysics, 0)
			boneDeform.scale = fs.getScale(bf, frame, boneDeform.bone, model)
		}
	}
}

// 該当キーフレにおけるボーンの移動位置
func (fs *BoneFrames) getPosition(
	bf *BoneFrame,
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	isAfterPhysics bool,
	loop int,
) (*mmath.MVec3, *mmath.MVec3) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMVec3(), mmath.NewMVec3()
	}

	vec := mmath.NewMVec3()
	if bf.MorphPosition != nil {
		vec.Add(bf.MorphPosition)
	}
	vec.Add(bf.Position)

	if bone.IsEffectorTranslation() {
		// 外部親変形ありの場合、外部親変形行列を掛ける
		effectPos := fs.getEffectPosition(frame, bone, model, isAfterPhysics, loop+1)
		return vec, effectPos
	}

	return vec, mmath.NewMVec3()
}

// 付与親を加味した移動位置
func (fs *BoneFrames) getEffectPosition(
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	isAfterPhysics bool,
	loop int,
) *mmath.MVec3 {
	if bone.EffectFactor == 0 || loop > 20 {
		// 付与率が0の場合、常に0になる
		// MMDエンジン対策で無限ループを避ける
		return mmath.NewMVec3()
	}

	if !(bone.EffectIndex > 0 && model.Bones.Contains(bone.EffectIndex)) {
		// 付与親が存在しない場合、常に0になる
		return mmath.NewMVec3()
	}

	// 付与親が存在する場合、付与親の回転角度を掛ける
	effectBone := model.Bones.Get(bone.EffectIndex)
	bf := fs.Get(effectBone.Name).Get(frame)
	pos, effectPos := fs.getPosition(bf, frame, effectBone, model, isAfterPhysics, loop+1)

	return pos.Add(effectPos).MulScalar(bone.EffectFactor)
}

// 該当キーフレにおけるボーンの回転角度
func (fs *BoneFrames) getRotation(
	bf *BoneFrame,
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	isAfterPhysics bool,
	loop int,
) (*mmath.MQuaternion, *mmath.MQuaternion) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMQuaternion(), nil
	}

	// FK(捩り) > IK(捩り) > 付与親(捩り)
	rot := mmath.NewMQuaternion()
	if bf.MorphRotation != nil {
		rot.Mul(bf.MorphRotation)
	}

	if bf.IkRotation != nil && !bf.IkRotation.IsIdent() {
		// IK用回転を持っている場合、置き換え
		rot.Mul(bf.IkRotation)
	} else {
		rot.Mul(bf.Rotation)

		if bone.HasFixedAxis() {
			rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
		}
	}

	if bone.IsEffectorRotation() {
		// 外部親変形ありの場合、外部親変形行列を掛ける
		effectQuat := fs.getEffectRotation(frame, bone, model, isAfterPhysics, loop+1)
		return rot.Shorten(), effectQuat.Shorten()
	}

	return rot.Shorten(), nil
}

// 付与親を加味した回転角度
func (fs *BoneFrames) getEffectRotation(
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	isAfterPhysics bool,
	loop int,
) *mmath.MQuaternion {
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
	effectBone := model.Bones.Get(bone.EffectIndex)
	bf := fs.Get(effectBone.Name).Get(frame)
	rot, effectRot := fs.getRotation(bf, frame, effectBone, model, isAfterPhysics, loop+1)

	if effectRot != nil {
		rot.Mul(effectRot)
	}

	return rot.MulScalar(bone.EffectFactor).Shorten()
}

// 該当キーフレにおけるボーンの拡大率
func (fs *BoneFrames) getScale(
	bf *BoneFrame,
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
) *mmath.MVec3 {

	scale := &mmath.MVec3{1, 1, 1}
	if bf.MorphScale != nil {
		scale.Add(bf.MorphScale)
	}

	if bf.Scale != nil {
		scale.Add(bf.Scale)
	}

	return scale
}
