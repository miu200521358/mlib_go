package vmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sort"
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
	fs.fillBoneDeform(frame, model, boneDeformsMap, beforeBoneDeltas, isAfterPhysics)

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
			if _, ok := boneDeforms.deforms[ikBone.Index]; !ok {
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
				if _, ok := boneDeforms.deforms[ikBone.Index]; !ok {
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
	// 中断FLGが入ったか否か
	aborts := make([]bool, len(ikBone.Ik.Links))

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
					"[%04d][%03d][%s][%05d][01][回転角度終了判定] originalLinkAngle: %v(%0.6f)\n",
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
			linkDeform := findBoneDeform(effectorBoneDeformsMap, linkBone)
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
			if ikLink.AngleLimit {
				// 角度制限が入ってる場合
				totalActualIkQuat, count = fs.calcIkLimitQuaternion(
					ikLink.MinAngleLimit.GetRadians(),
					ikLink.MaxAngleLimit.GetRadians(),
					linkQuat,
					linkAxis,
					linkAngle,
					mmath.MVec3UnitX,
					mmath.MVec3UnitY,
					mmath.MVec3UnitZ,
					frame,
					count,
					loop,
					linkBone.Name,
					ikMotion,
					ikFile,
				)
			} else if ikLink.LocalAngleLimit {
				// ローカル角度制限が入ってる場合
				totalActualIkQuat, count = fs.calcIkLimitQuaternion(
					ikLink.LocalMinAngleLimit.GetRadians(),
					ikLink.LocalMaxAngleLimit.GetRadians(),
					linkQuat,
					linkAxis,
					linkAngle,
					linkBone.NormalizedLocalAxisX,
					linkBone.NormalizedLocalAxisY,
					linkBone.NormalizedLocalAxisZ,
					frame,
					count,
					loop,
					linkBone.Name,
					ikMotion,
					ikFile,
				)
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

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil && linkDeform.rotation != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][15] 前回差分中断判定: %v(%0.6f) 前回: %s 今回: %s\n",
					frame, loop, linkBone.Name, count-1,
					1-totalActualIkQuat.Dot(linkDeform.rotation) < 1e-8, 1-totalActualIkQuat.Dot(linkDeform.rotation),
					linkDeform.rotation.ToDegrees().String(), totalActualIkQuat.ToMMDDegrees().String())
			}

			// 前回（既存）とほぼ同じ回転量の場合、中断FLGを立てる
			if linkDeform.rotation != nil && 1-totalActualIkQuat.Dot(linkDeform.rotation) < 1e-8 {
				aborts[lidx] = true
			}

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

		if slices.Index(aborts, false) == -1 {
			// すべてのリンクボーンで中断FLG = true の場合、終了
			break ikLoop
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
func (fs *BoneFrames) calcIkLimitQuaternion(
	minAngleLimitRadians *mmath.MVec3,
	maxAngleLimitRadians *mmath.MVec3,
	linkQuat *mmath.MQuaternion,
	quatAxis *mmath.MVec3,
	quatAngle float64,
	xAxisVector *mmath.MVec3,
	yAxisVector *mmath.MVec3,
	zAxisVector *mmath.MVec3,
	frame int,
	count int,
	loop int,
	linkBoneName string,
	ikMotion *VmdMotion,
	ikFile *os.File,
) (*mmath.MQuaternion, int) {
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

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][04][角度制限] minAngleLimitRadians: %s, maxAngleLimitRadians:%s\n",
			frame, loop, linkBoneName, count-1, minAngleLimitRadians.String(), maxAngleLimitRadians.String())
	}

	// 現在IKリンクに入る可能性のあるすべての角度
	totalIkQuat := quat.Muled(linkQuat)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		bf := NewBoneFrame(count)
		bf.Rotation = totalIkQuat
		ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
		count++
	}

	totalIkRads := totalIkQuat.ToRadians()
	totalIkRad := totalIkQuat.ToRadian()
	isXGimbal := totalIkRad > mmath.GIMBAL1_RAD ||
		(math.Abs(totalIkRads.GetY()) >= mmath.GIMBAL2_RAD && math.Abs(totalIkRads.GetZ()) >= mmath.GIMBAL2_RAD)
	isYGimbal := totalIkRad > mmath.GIMBAL1_RAD ||
		(math.Abs(totalIkRads.GetX()) >= mmath.GIMBAL2_RAD && math.Abs(totalIkRads.GetZ()) >= mmath.GIMBAL2_RAD)
	isZGimbal := totalIkRad > mmath.GIMBAL1_RAD ||
		(math.Abs(totalIkRads.GetX()) >= mmath.GIMBAL2_RAD && math.Abs(totalIkRads.GetY()) >= mmath.GIMBAL2_RAD)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][05][角度制限] totalIkQuat: %s(%s), isXGimbal: %v, "+
			"isYGimbal: %v, isZGimbal: %v\n",
			frame, loop, linkBoneName, count-1, totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
			isXGimbal, isYGimbal, isZGimbal)
	}

	tX := totalIkRads.GetX()

	if minAngleLimitRadians.GetY() == 0 && maxAngleLimitRadians.GetY() == 0 &&
		minAngleLimitRadians.GetZ() == 0 && maxAngleLimitRadians.GetZ() == 0 {
		// Xのみ動かせる場合、全部の値を引き取る
		tX = totalIkRad
		if totalIkQuat.GetXYZ().Dot(xAxisVector) < 0 {
			tX *= -1
		}
	}

	tY := totalIkRads.GetY()

	// if minAngleLimitRadians.GetX() == 0 && maxAngleLimitRadians.GetX() == 0 &&
	// 	minAngleLimitRadians.GetZ() == 0 && maxAngleLimitRadians.GetZ() == 0 {
	// 	// Yのみ動かせる場合、全部の値を引き取る
	// 	tY = totalIkRad
	// 	if totalIkQuat.GetXYZ().Dot(yAxisVector) < 0 {
	// 		tY *= -1
	// 	}
	// }

	tZ := totalIkRads.GetZ()

	// if minAngleLimitRadians.GetX() == 0 && maxAngleLimitRadians.GetX() == 0 &&
	// 	minAngleLimitRadians.GetY() == 0 && maxAngleLimitRadians.GetY() == 0 {
	// 	// Zのみ動かせる場合、全部の値を引き取る
	// 	tZ = totalIkRad
	// 	if totalIkQuat.GetXYZ().Dot(zAxisVector) < 0 {
	// 		tZ *= -1
	// 	}
	// }

	fSX := math.Sin(tX)  // sin(θx)
	fX := math.Asin(fSX) // 一軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限X] fSX: %.5f, tX: %.5f, fX: %.5f\n",
			frame, loop, linkBoneName, count-1, fSX, tX, fX)
	}

	if isXGimbal {
		// ジンバルロック回避
		if fX < 0 {
			fX = -(math.Pi - fX)
		} else {
			fX = math.Pi - fX
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][08][角度制限X][ジンバルロック回避] "+
				"fSX: %.5f, tX: %.5f, fX: %.5f\n", frame, loop, linkBoneName, count-1, fSX, tX, fX)
		}
	}

	fX = fs.judgeIkAngle(minAngleLimitRadians.GetX(), maxAngleLimitRadians.GetX(), fX,
		frame, count, loop, linkBoneName, ikMotion, ikFile)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限X][judge] fSX: %.5f, tX: %.5f, fX: %.5f\n",
			frame, loop, linkBoneName, count-1, fSX, tX, fX)
	}

	if minAngleLimitRadians.GetX() != 0 || maxAngleLimitRadians.GetX() != 0 {
		// X軸のみ回れるIK制限の場合、ここで終了(足IK想定)
		return mmath.NewMQuaternionFromAxisAngles(xAxisVector, fX).Shorten(), count
	}

	// Y軸回り
	fCY := math.Cos(tY)
	fY := math.Acos(fCY) // Y軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限Y] tY: %.5f, fY: %.5f\n",
			frame, loop, linkBoneName, count-1, tY, fY)
	}

	fY = fs.judgeIkAngle(minAngleLimitRadians.GetY(), maxAngleLimitRadians.GetY(), fY,
		frame, count, loop, linkBoneName, ikMotion, ikFile)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限Y][judge] tY: %.5f, fY: %.5f\n",
			frame, loop, linkBoneName, count-1, tY, fY)
	}

	// Z軸回り
	fCZ := math.Sin(tZ)
	fZ := math.Asin(fCZ) // Z軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限Z] tZ: %.5f, fZ: %.5f\n",
			frame, loop, linkBoneName, count-1, tZ, fZ)
	}

	fZ = fs.judgeIkAngle(minAngleLimitRadians.GetZ(), maxAngleLimitRadians.GetZ(), fZ,
		frame, count, loop, linkBoneName, ikMotion, ikFile)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][07][角度制限Z][judge] tZ: %.5f, fZ: %.5f\n",
			frame, loop, linkBoneName, count-1, tZ, fZ)
	}

	xQuat := mmath.NewMQuaternionFromAxisAngles(xAxisVector, fX)
	yQuat := mmath.NewMQuaternionFromAxisAngles(yAxisVector, -fY)
	zQuat := mmath.NewMQuaternionFromAxisAngles(zAxisVector, -fZ)
	totalActualIkQuat := yQuat.Muled(xQuat).Muled(zQuat).Shorten()

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		{
			{
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][11][X角度制限結果] xQuat: %s(%s)\n",
					frame, loop, linkBoneName, count-1, xQuat.String(), xQuat.ToMMDDegrees().String(),
				)
			}

			bf := NewBoneFrame(count)
			bf.Rotation = xQuat
			ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
			count++
		}
		{
			{
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][11][Y角度制限結果] yQuat: %s(%s)\n",
					frame, loop, linkBoneName, count-1, yQuat.String(), yQuat.ToMMDDegrees().String(),
				)
			}

			bf := NewBoneFrame(count)
			bf.Rotation = yQuat
			ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
			count++
		}
		{
			{
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][11][Z角度制限結果] zQuat: %s(%s)\n",
					frame, loop, linkBoneName, count-1, zQuat.String(), zQuat.ToMMDDegrees().String(),
				)
			}

			bf := NewBoneFrame(count)
			bf.Rotation = zQuat
			ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
			count++
		}
		{
			{
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][11][total角度制限結果] totalActualIkQuat: %s(%s)\n",
					frame, loop, linkBoneName, count-1, totalActualIkQuat.String(), totalActualIkQuat.ToMMDDegrees().String(),
				)
			}

			bf := NewBoneFrame(count)
			bf.Rotation = totalActualIkQuat
			ikMotion.AppendRegisteredBoneFrame(linkBoneName, bf)
			count++
		}
	}

	return totalActualIkQuat, count
}

func (fs *BoneFrames) judgeIkAngle(
	minAngleLimit float64,
	maxAngleLimit float64,
	fX float64,
	frame int,
	count int,
	loop int,
	linkBoneName string,
	ikMotion *VmdMotion,
	ikFile *os.File,
) float64 {
	// 角度の制限
	if fX < minAngleLimit {
		tf := 2*minAngleLimit - fX

		if tf <= maxAngleLimit {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][09][角度制限][負角度制限Over] fX: %.5f, tf: %.5f\n",
					frame, loop, linkBoneName, count-1, fX, tf)
			}

			fX = tf
		} else {
			fX = mmath.ClampFloat(fX, minAngleLimit, maxAngleLimit)
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][09][角度制限][負角度制限] fX: %.5f, tf: %.5f\n",
				frame, loop, linkBoneName, count-1, fX, tf)
		}
	} else if fX > maxAngleLimit {
		tf := 2*maxAngleLimit - fX

		if tf >= minAngleLimit {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][10][角度制限][正角度制限Over] fX: %.5f, tf: %.5f\n",
					frame, loop, linkBoneName, count-1, fX, tf)
			}

			fX = tf
		} else {
			fX = mmath.ClampFloat(fX, minAngleLimit, maxAngleLimit)
		}

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][10][角度制限][正角度制限] fX: %.5f, tf: %.5f\n",
				frame, loop, linkBoneName, count-1, fX, tf)
		}
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][10][角度制限][角度制限結果] fX: %.5f\n",
			frame, loop, linkBoneName, count-1, fX)
	}

	return fX
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
		delta := boneDeltas.Get(boneIndex)
		if delta != nil && delta.globalMatrix != nil {
			// 既にグローバル行列が計算済みの場合、スルー
			continue
		}

		deform := boneDeforms.deforms[boneIndex]
		deform.unitMatrix = mmath.NewMMat4()

		// スケール
		if deform.scale != nil && !deform.scale.IsOne() {
			deform.unitMatrix.Mul(deform.scale.ToScaleMat4())
		}

		// 回転
		rot := deform.rotation
		isEffectorRot := false
		if deform.effectRotation != nil && !deform.effectRotation.IsIdent() {
			rot = deform.effectRotation.Muled(rot)
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
			pos = deform.effectPosition.Added(pos)
		}
		if pos != nil && !pos.IsZero() {
			deform.unitMatrix.Mul(pos.ToMat4())
		}

		// 逆BOf行列(初期姿勢行列)
		deform.unitMatrix.Mul(deform.bone.RevertOffsetMatrix)
		continue
	}

	for _, boneIndex := range boneDeforms.boneIndexes {
		delta := boneDeltas.Get(boneIndex)
		if delta != nil && delta.globalMatrix != nil {
			// 既にグローバル行列が計算済みの場合、スルー
			continue
		}

		var globalMatrix *mmath.MMat4
		deform := boneDeforms.deforms[boneIndex]
		parentDelta := boneDeltas.Get(deform.bone.ParentIndex)
		if parentDelta != nil {
			globalMatrix = deform.unitMatrix.Muled(parentDelta.GlobalMatrix())
		} else {
			// 対象ボーン自身の行列をかける
			globalMatrix = deform.unitMatrix.Copy()
		}

		delta = NewBoneDelta(
			deform.bone,
			frame,
			globalMatrix,          // グローバル行列
			deform.position,       // キーフレの移動量
			deform.effectPosition, // キーフレの付与移動量
			deform.rotation,       // キーフレの回転量
			deform.effectRotation, // キーフレの付与回転量
			deform.scale,          // 拡大率
		)

		boneDeltas.Append(delta)
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
			boneIndexes: make([]int, 0),
		}
	}

	for _, ap := range isAfterPhysicsList {
		// ボーン名の存在チェック用マップ
		targetSortedBones := model.Bones.LayerSortedBones[ap]

		boneDeforms := &boneDeforms{
			deforms:     make(map[int]*BoneDeform),
			boneIndexes: make([]int, 0, len(targetSortedBones)),
		}

		if len(boneNames) > 0 {
			// 指定ボーンに関連するボーンのみ対象とする
			layerIndexes := make(pmx.LayerIndexes, 0, len(targetSortedBones))

			for _, boneName := range boneNames {
				// ボーン名の追加
				bone := model.Bones.GetByName(boneName)
				layerIndexes = append(layerIndexes, pmx.LayerIndex{Layer: bone.Layer, Index: bone.Index})

				// 関連するボーンの追加
				relativeBoneIndexes := bone.RelativeBoneIndexes
				for _, index := range relativeBoneIndexes {
					bone := model.Bones.Get(index)
					if !layerIndexes.Contains(bone.Index) {
						layerIndexes = append(layerIndexes, pmx.LayerIndex{Layer: bone.Layer, Index: bone.Index})
					}
				}
			}
			sort.Sort(layerIndexes)

			for _, layerIndex := range layerIndexes {
				bone := model.Bones.Get(layerIndex.Index)
				boneDeforms.deforms[bone.Index] = &BoneDeform{bone: bone}
				boneDeforms.boneIndexes = append(boneDeforms.boneIndexes, bone.Index)
			}
		} else {
			// 変形階層・ボーンINDEXでソート
			for k := range len(targetSortedBones) {
				bone := targetSortedBones[k]
				boneDeforms.deforms[bone.Index] = &BoneDeform{bone: bone}
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
	boneDeformsMap map[bool]*boneDeforms,
	beforeBoneDeltas *BoneDeltas,
	isAfterPhysics bool,
) {
	for _, boneIndex := range boneDeformsMap[isAfterPhysics].boneIndexes {
		boneDeform := boneDeformsMap[isAfterPhysics].deforms[boneIndex]
		var delta *BoneDelta
		if beforeBoneDeltas != nil {
			delta = beforeBoneDeltas.Get(boneIndex)
		}

		if boneDeform.unitMatrix == nil && (delta == nil || delta.globalMatrix == nil) {
			bf := fs.Get(boneDeform.bone.Name).Get(frame)
			// ボーンの移動位置、回転角度、拡大率を取得
			boneDeform.position, boneDeform.effectPosition =
				fs.getPosition(bf, frame, boneDeform.bone, model, beforeBoneDeltas, isAfterPhysics, 0)
			boneDeform.rotation, boneDeform.effectRotation =
				fs.getRotation(bf, frame, boneDeform.bone, model, beforeBoneDeltas, isAfterPhysics, 0)
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
	beforeBoneDeltas *BoneDeltas,
	isAfterPhysics bool,
	loop int,
) (*mmath.MVec3, *mmath.MVec3) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMVec3(), mmath.NewMVec3()
	}

	var vec *mmath.MVec3
	if bf.MorphPosition != nil {
		vec = bf.MorphPosition.Add(bf.Position)
	} else {
		vec = bf.Position
	}

	if bone.IsEffectorTranslation() && bone.CanTranslate() {
		// 外部親変形ありの場合、外部親位置を取得する
		effectPos := fs.getEffectPosition(frame, bone, model, beforeBoneDeltas, isAfterPhysics, loop+1)
		return vec, effectPos
	}

	return vec, mmath.NewMVec3()
}

// 付与親を加味した移動位置
func (fs *BoneFrames) getEffectPosition(
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	beforeBoneDeltas *BoneDeltas,
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

	if beforeBoneDeltas != nil {
		effectDelta := beforeBoneDeltas.Get(effectBone.Index)
		if effectDelta != nil && effectDelta.framePosition != nil {
			return effectDelta.framePosition.MulScalar(bone.EffectFactor)
		}
	}

	bf := fs.Get(effectBone.Name).Get(frame)
	pos, effectPos := fs.getPosition(bf, frame, effectBone, model, beforeBoneDeltas, isAfterPhysics, loop+1)

	return pos.Add(effectPos).MulScalar(bone.EffectFactor)
}

// 該当キーフレにおけるボーンの回転角度
func (fs *BoneFrames) getRotation(
	bf *BoneFrame,
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	beforeBoneDeltas *BoneDeltas,
	isAfterPhysics bool,
	loop int,
) (*mmath.MQuaternion, *mmath.MQuaternion) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMQuaternion(), nil
	}

	// FK(捩り) > IK(捩り) > 付与親(捩り)
	var rot *mmath.MQuaternion
	if bf.IkRotation != nil && !bf.IkRotation.IsIdent() {
		// IK用回転を持っている場合、置き換え
		rot = bf.IkRotation
	} else {
		rot = bf.Rotation
	}

	if bf.MorphRotation != nil {
		rot = bf.MorphRotation.Mul(rot)
	}

	if bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
	}

	if bone.IsEffectorRotation() && bone.CanRotate() {
		// 外部親変形ありの場合、外部親回転を取得する
		effectQuat := fs.getEffectRotation(frame, bone, model, beforeBoneDeltas, isAfterPhysics, loop+1)
		return rot.Shorten(), effectQuat.Shorten()
	}

	return rot.Shorten(), nil
}

// 付与親を加味した回転角度
func (fs *BoneFrames) getEffectRotation(
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	beforeBoneDeltas *BoneDeltas,
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

	if beforeBoneDeltas != nil {
		effectDelta := beforeBoneDeltas.Get(effectBone.Index)
		if effectDelta != nil && effectDelta.frameRotation != nil {
			return effectDelta.frameRotation.MuledScalar(bone.EffectFactor).Shorten()
		}
	}

	bf := fs.Get(effectBone.Name).Get(frame)
	rot, effectRot := fs.getRotation(bf, frame, effectBone, model, beforeBoneDeltas, isAfterPhysics, loop+1)

	if effectRot != nil {
		rot.Mul(effectRot)
	}

	return rot.MuledScalar(bone.EffectFactor).Shorten()
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
