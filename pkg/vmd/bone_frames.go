package vmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
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
	ikFrame *IkFrame,
) *BoneDeltas {
	// mlog.Memory(fmt.Sprintf("Deform 1)frame: %d", frame))
	deformBoneIndexes, boneDeltas := fs.prepareDeltas(frame, model, boneNames, isCalcIk, beforeBoneDeltas, ikFrame)
	// mlog.Memory(fmt.Sprintf("Deform 2)frame: %d", frame))
	boneDeltas = fs.calcBoneDeltas(frame, model, deformBoneIndexes, boneDeltas)
	// mlog.Memory(fmt.Sprintf("Deform 3)frame: %d", frame))
	return boneDeltas
}

func (fs *BoneFrames) prepareDeltas(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	beforeBoneDeltas *BoneDeltas,
	ikFrame *IkFrame,
) ([]int, *BoneDeltas) {
	// mlog.Memory(fmt.Sprintf("prepareDeltas 1)frame: %d", frame))

	deformBoneIndexes, boneDeltas := fs.createBoneDeltas(frame, model, boneNames, beforeBoneDeltas)

	// mlog.Memory(fmt.Sprintf("prepareDeltas 2)frame: %d", frame))

	// IK事前計算
	if isCalcIk {
		// ボーン変形行列操作
		boneDeltas = fs.prepareIk(frame, model, deformBoneIndexes, boneDeltas, ikFrame)
	}

	// mlog.Memory(fmt.Sprintf("prepareDeltas 3)frame: %d", frame))

	// ボーンデフォーム情報を埋める
	boneDeltas = fs.fillBoneDeform(frame, model, deformBoneIndexes, boneDeltas)

	// mlog.Memory(fmt.Sprintf("prepareDeltas 4)frame: %d", frame))

	return deformBoneIndexes, boneDeltas
}

// IK事前計算処理
func (fs *BoneFrames) prepareIk(
	frame int,
	model *pmx.PmxModel,
	deformBoneIndexes []int,
	boneDeltas *BoneDeltas,
	ikFrame *IkFrame,
) *BoneDeltas {
	for i, boneIndex := range deformBoneIndexes {
		// ボーンIndexがIkTreeIndexesに含まれていない場合、スルー
		if _, ok := model.Bones.IkTreeIndexes[boneIndex]; !ok {
			continue
		}

		for m := range len(model.Bones.IkTreeIndexes[boneIndex]) {
			ikBone := model.Bones.Get(model.Bones.IkTreeIndexes[boneIndex][m])

			if ikFrame == nil || ikFrame.IsEnable(ikBone.Name) {
				var prefixPath string
				if mlog.IsIkVerbose() {
					// IK計算デバッグ用モーション
					dirPath := fmt.Sprintf("%s/IK_step", filepath.Dir(model.Path))
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						log.Fatal(err)
					}

					date := time.Now().Format("20060102_150405")
					prefixPath = fmt.Sprintf("%s/%04d_%s_%03d_%03d", dirPath, frame, date, i, m)
				}

				boneDeltas = fs.calcIk(frame, ikBone, model, boneDeltas, ikFrame, prefixPath)
			}
		}
	}

	mlog.IV("[IK計算終了][%04d] -----------------------------------------------", frame)

	return boneDeltas
}

// IK計算
func (fs *BoneFrames) calcIk(
	frame int,
	ikBone *pmx.Bone,
	model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	ikFrame *IkFrame,
	prefixPath string,
) *BoneDeltas {
	var err error
	var ikFile *os.File
	var ikMotion *VmdMotion
	count := int(1.0)

	if mlog.IsIkVerbose() {
		// IK計算デバッグ用モーション
		ikMotionPath := fmt.Sprintf("%s_%s.vmd", prefixPath, ikBone.Name)
		ikMotion = NewVmdMotion(ikMotionPath)

		ikLogPath := fmt.Sprintf("%s_%s.log", prefixPath, ikBone.Name)
		ikFile, err = os.OpenFile(ikLogPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(ikFile, "----------------------------------------")
		fmt.Println(ikFile, "[IK計算出力先][%04d][%s] %s", frame, ikMotionPath)
	}
	defer func() {
		mlog.IV("[IK計算終了][%04d][%s]", frame, ikBone.Name)

		if ikMotion != nil {
			Write(ikMotion)
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

	// IKターゲットボーン
	effectorBone := model.Bones.Get(ikBone.Ik.BoneIndex)
	// IK関連の行列を一括計算
	ikDeltas := fs.Deform(frame, model, []string{ikBone.Name}, false, boneDeltas, ikFrame)
	// エフェクタ関連情報取得
	effectorDeformBoneIndexes, boneDeltas :=
		fs.prepareDeltas(frame, model, []string{effectorBone.Name}, false, boneDeltas, ikFrame)
	// 中断FLGが入ったか否か
	aborts := make([]bool, len(ikBone.Ik.Links))

	// // リンクボーン全体の長さ
	// ikBoneTotalLength := 0.0
	// for _, l := range ikBone.Ik.Links {
	// 	if model.Bones.Contains(l.BoneIndex) {
	// 		ikBoneTotalLength += model.Bones.Get(l.BoneIndex).ParentRelativePosition.Length()
	// 	}
	// }

	// 一段IKであるか否か
	isOneLinkIk := len(ikBone.Ik.Links) == 1
	// ループ回数
	loopCount := max(ikBone.Ik.LoopCount, 1)

	// IK計算
ikLoop:
	for loop := 0; loop < loopCount; loop++ {
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

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d] -------------------------------------------- \n",
					frame, loop, linkBone.Name, count-1)
			}

			for _, l := range ikBone.Ik.Links {
				if boneDeltas.Get(l.BoneIndex) != nil {
					boneDeltas.Get(l.BoneIndex).unitMatrix = nil
				}
			}
			if boneDeltas.Get(ikBone.Ik.BoneIndex) != nil {
				boneDeltas.Get(ikBone.Ik.BoneIndex).unitMatrix = nil
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := NewBoneFrame(count)
				bf.Position = ikDeltas.Get(ikBone.Index).FramePosition()
				bf.Rotation = ikDeltas.Get(ikBone.Index).LocalRotation()
				ikMotion.AppendRegisteredBoneFrame(ikBone.Name, bf)
				count++

				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Local] ikGlobalPosition: %s\n",
					frame, loop, linkBone.Name, count-1, bf.Position.MMD().String())
			}

			// IK関連の行列を取得
			boneDeltas = fs.calcBoneDeltas(frame, model, effectorDeformBoneIndexes, boneDeltas)

			// リンクボーンのIK角度を取得
			linkDelta := boneDeltas.Get(linkBone.Index)
			if linkDelta == nil {
				linkDelta = &BoneDelta{Bone: linkBone, Frame: frame}
			}
			linkQuat := linkDelta.LocalRotation()

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := NewBoneFrame(count)
				bf.Rotation = linkQuat.Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++

				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][linkQuat] %s(%s)\n",
					frame, loop, linkBone.Name, count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String(),
				)
			}

			// IKボーンのグローバル位置
			ikGlobalPosition := ikDeltas.Get(ikBone.Index).GlobalPosition()

			// 現在のIKターゲットボーンのグローバル位置を取得
			effectorGlobalPosition := boneDeltas.Get(effectorBone.Index).GlobalPosition()

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][Global] [%s]ikGlobalPosition: %s, "+
						"[%s]effectorGlobalPosition: %s, [%s]linkGlobalPosition: %s\n",
					frame, loop, linkBone.Name, count-1,
					ikBone.Name, ikGlobalPosition.MMD().String(),
					effectorBone.Name, effectorGlobalPosition.MMD().String(),
					linkBone.Name, boneDeltas.Get(linkBone.Index).GlobalPosition().MMD().String())
			}

			// 注目ノード（実際に動かすボーン=リンクボーン）
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := boneDeltas.Get(linkBone.Index).GlobalMatrix().Inverted()
			// 注目ノードを起点とした、エフェクタのローカル位置
			effectorLocalPosition := linkInvMatrix.MulVec3(effectorGlobalPosition)
			// 注目ノードを起点とした、IK目標のグローバル差分
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			// // 注目ノード（実際に動かすボーン=リンクボーン）
			// linkGlobalPosition := boneDeltas.Get(linkBone.Index).GlobalPosition()
			// // 注目ノードを起点とした、エフェクタのグローバル差分
			// effectorLocalPosition := linkGlobalPosition.Subed(effectorGlobalPosition)
			// // 注目ノードを起点とした、IK目標のグローバル差分
			// ikLocalPosition := linkGlobalPosition.Subed(ikGlobalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][Local] effectorLocalPosition: %s, ikLocalPosition: %s\n",
					frame, loop, linkBone.Name, count-1,
					effectorLocalPosition.MMD().String(), ikLocalPosition.MMD().String())
			}

			effectorLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			if effectorLocalPosition.Distance(ikLocalPosition) < 1e-7 {
				break ikLoop
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][Local] effectorLocalPositionNorm: %s, ikLocalPositionNorm: %s\n",
					frame, loop, linkBone.Name, count-1,
					effectorLocalPosition.MMD().String(), ikLocalPosition.MMD().String())
			}

			// 単位角
			unitRad := ikBone.Ik.UnitRotation.GetRadians().GetX() * float64(lidx+1)
			linkDot := ikLocalPosition.Dot(effectorLocalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][回転角度] unitRad: %.8f (%.5f), linkDot: %.8f\n",
					frame, loop, linkBone.Name, count-1, unitRad, mmath.ToDegree(unitRad), linkDot,
				)
			}

			// 回転角(ラジアン)
			// 単位角を超えないようにする
			originalLinkAngle := math.Acos(mmath.ClampFloat(linkDot, -1, 1))
			linkAngle := mmath.ClampFloat(originalLinkAngle, -unitRad, unitRad)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][単位角制限] linkAngle: %.8f(%.5f), originalLinkAngle: %.8f(%.5f)\n",
					frame, loop, linkBone.Name, count-1, linkAngle, mmath.ToDegree(linkAngle),
					originalLinkAngle, mmath.ToDegree(originalLinkAngle),
				)
			}

			// 角度がほとんどない場合
			if math.Abs(linkAngle) < 1e-8 {
				if isOneLinkIk {
					// 一段IKは単位角度を回す
					linkAngle = unitRad // * float64(loop+1)
				} else {
					// 多段IKは終了
					break ikLoop
				}
			}

			// 回転軸
			var originalLinkAxis, linkAxis *mmath.MVec3
			if (!isOneLinkIk ||
				ikLink.MinAngleLimit.GetRadians().IsOnlyY() || ikLink.MaxAngleLimit.GetRadians().IsOnlyY() ||
				ikLink.MinAngleLimit.GetRadians().IsOnlyZ() || ikLink.MaxAngleLimit.GetRadians().IsOnlyZ()) &&
				ikLink.AngleLimit {
				// グローバル軸制限
				linkAxis, originalLinkAxis = fs.getLinkAxis(
					ikLink.MinAngleLimit.GetRadians(),
					ikLink.MaxAngleLimit.GetRadians(),
					effectorLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name, ikMotion, ikFile,
				)
			} else if !isOneLinkIk && ikLink.LocalAngleLimit {
				// ローカル軸制限
				linkAxis, originalLinkAxis = fs.getLinkAxis(
					ikLink.LocalMinAngleLimit.GetRadians(),
					ikLink.LocalMaxAngleLimit.GetRadians(),
					effectorLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name, ikMotion, ikFile,
				)
			} else {
				// 軸制限なし or 初回
				linkAxis, originalLinkAxis = fs.getLinkAxis(
					mmath.MVec3MinVal,
					mmath.MVec3MaxVal,
					effectorLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name, ikMotion, ikFile,
				)
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][回転軸] linkAxis: %s, originalLinkAxis: %s\n",
					frame, loop, linkBone.Name, count-1, linkAxis.String(), originalLinkAxis.String(),
				)
			}

			if linkBone.HasFixedAxis() {
				if linkAxis.Dot(linkBone.NormalizedFixedAxis) < 0 {
					linkAngle = -linkAngle
				}

				// 軸制限ありの場合、軸にそった理想回転量とする
				linkAxis = linkBone.NormalizedFixedAxis

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][軸制限] linkAxis: %s\n",
						frame, loop, linkBone.Name, count-1, linkAxis.String(),
					)
				}
			}

			originalIkQuat := mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, originalLinkAngle)
			ikQuat := mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, linkAngle)

			originalTotalIkQuat := linkQuat.Muled(originalIkQuat)
			totalIkQuat := linkQuat.Muled(ikQuat)

			if isOneLinkIk && loop == 1 {
				totalIkQuat = ikQuat.Muled(linkQuat)
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				{
					bf := NewBoneFrame(count)
					bf.Rotation = originalTotalIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][originalTotalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name, count-1, originalTotalIkQuat.String(), originalTotalIkQuat.ToMMDDegrees().String())

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][originalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name, count-1, originalIkQuat.String(), originalIkQuat.ToMMDDegrees().String())
				}
				{
					bf := NewBoneFrame(count)
					bf.Rotation = totalIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][totalIkQuat] %s(%s)\n",
						frame, loop, linkBone.Name, count-1, totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String())

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][ikQuat] %s(%s)\n",
						frame, loop, linkBone.Name, count-1, ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			}

			var resultIkQuat *mmath.MQuaternion
			if ikLink.AngleLimit {
				// 角度制限が入ってる場合
				resultIkQuat, count = fs.calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.MinAngleLimit.GetRadians(),
					ikLink.MaxAngleLimit.GetRadians(),
					mmath.MVec3UnitX, mmath.MVec3UnitY, mmath.MVec3UnitZ,
					loop, loopCount,
					frame, count, linkBone.Name, ikMotion, ikFile,
				)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][角度制限後] resultIkQuat: %s(%s), totalIkQuat: %s(%s), ikQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String(),
						totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
						ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			} else if ikLink.LocalAngleLimit {
				// ローカル角度制限が入ってる場合
				resultIkQuat, count = fs.calcIkLimitQuaternion(
					totalIkQuat,
					ikLink.LocalMinAngleLimit.GetRadians(),
					ikLink.LocalMaxAngleLimit.GetRadians(),
					linkBone.NormalizedLocalAxisX, linkBone.NormalizedLocalAxisY, linkBone.NormalizedLocalAxisZ,
					loop, loopCount,
					frame, count, linkBone.Name, ikMotion, ikFile,
				)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][ローカル角度制限後] resultIkQuat: %s(%s), totalIkQuat: %s(%s), ikQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String(),
						totalIkQuat.String(), totalIkQuat.ToMMDDegrees().String(),
						ikQuat.String(), ikQuat.ToMMDDegrees().String())
				}
			} else {
				// 角度制限なしの場合
				resultIkQuat = totalIkQuat
			}

			if linkBone.HasFixedAxis() {
				// 軸制限ありの場合、軸にそった理想回転量とする
				resultIkQuat = resultIkQuat.ToFixedAxisRotation(linkBone.FixedAxis)

				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					bf := NewBoneFrame(count)
					bf.Rotation = resultIkQuat.Copy()
					ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
					count++

					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][軸制限後] resultIkQuat: %s(%s)\n",
						frame, loop, linkBone.Name, count-1, resultIkQuat.String(), resultIkQuat.ToMMDDegrees().String())
				}
			}

			// IKの結果を更新
			linkDelta.frameRotation = resultIkQuat
			boneDeltas.Append(linkDelta)

			// 前回（既存）とほぼ同じ回転量の場合、中断FLGを立てる
			isAbort := linkQuat != nil && 1-totalIkQuat.Dot(linkQuat) < 1e-10

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil && linkQuat != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d] 前回差分中断判定: %v(dot: %0.6f, distance: %0.8f) 前回: %s 今回: %s\n",
					frame, loop, linkBone.Name, count-1,
					isAbort, 1-totalIkQuat.Dot(linkQuat), effectorLocalPosition.Distance(ikLocalPosition),
					linkQuat.ToDegrees().String(), totalIkQuat.ToMMDDegrees().String())
			}

			if isAbort {
				aborts[lidx] = true
			}

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := NewBoneFrame(count)
				bf.Rotation = linkDelta.LocalRotation().Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++

				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][結果] bf.Rotation: %s(%s)\n",
					frame, loop, linkBone.Name, count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String())
			}

			// remainAngle := resultIkQuat.ToSignedRadian() - originalTotalIkQuat.ToSignedRadian()
			// remainingQuat := mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, remainAngle)

			// remainingQuat := mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, remainAngle)
			// var remainingQuat *mmath.MQuaternion
			// if isIkFar {
			// } else {
			// 	remainAngle := resultIkQuat.ToSignedRadian() - originalTotalIkQuat.ToSignedRadian()
			// 	remainingQuat = mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, remainAngle)
			// }
			// remainingQuat := resultIkQuat.Muled(originalTotalIkQuat.Inverted())
			// remainingQuat := mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, resultIkQuat.Muled(originalTotalIkQuat.Inverted()).ToSignedRadian())
			// remainingQuat := resultIkQuat .Muled(totalIkQuat.Inverted())

			// remainingQuat := resultIkQuat.Muled(originalTotalIkQuat.Inverted())

			// if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			// 	fmt.Fprintf(ikFile,
			// 		"[%04d][%03d][%s][%05d][残存] remainingQuat: %s(%s)\n",
			// 		frame, loop, linkBone.Name, count-1, remainingQuat.String(), remainingQuat.ToMMDDegrees().String())
			// }

			// if !remainingQuat.IsIdent() && false {
			// 	// IKターゲットの回転に残回転量を加算
			// 	effectorDelta := boneDeltas.Get(effectorBone.Index)
			// 	if effectorDelta == nil {
			// 		effectorDelta = &BoneDelta{Bone: effectorBone, Frame: frame}
			// 	}

			// 	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			// 		bf := NewBoneFrame(count)
			// 		bf.Rotation = effectorDelta.frameRotation.Copy()
			// 		ikMotion.AppendRegisteredBoneFrame(effectorBone.Name, bf)
			// 		count++

			// 		fmt.Fprintf(ikFile,
			// 			"[%04d][%03d][%s][%05d][残存] effectorRot: %s(%s)\n",
			// 			frame, loop, linkBone.Name, count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String())
			// 	}

			// 	effectorDelta.frameRotation = remainingQuat.Muled(effectorDelta.FrameRotation())
			// 	// effectorDelta.frameRotation = effectorDelta.FrameRotation().Muled(remainingQuat)
			// 	boneDeltas.Append(effectorDelta)

			// 	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			// 		bf := NewBoneFrame(count)
			// 		bf.Rotation = effectorDelta.frameRotation.Copy()
			// 		ikMotion.AppendRegisteredBoneFrame(effectorBone.Name, bf)
			// 		count++

			// 		fmt.Fprintf(ikFile,
			// 			"[%04d][%03d][%s][%05d][残存加算] effectorRot: %s(%s)\n",
			// 			frame, loop, linkBone.Name, count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String())
			// 	}
			// }
		}

		// if slices.Index(aborts, false) == -1 {
		// 	// すべてのリンクボーンで中断FLG = true の場合、終了
		// 	break ikLoop
		// }
	}

	ikDeltas = nil
	return boneDeltas
}

func (fs *BoneFrames) getLinkAxis(
	minAngleLimitRadians *mmath.MVec3,
	maxAngleLimitRadians *mmath.MVec3,
	effectorLocalPosition, ikLocalPosition *mmath.MVec3,
	frame int,
	count int,
	loop int,
	linkBoneName string,
	ikMotion *VmdMotion,
	ikFile *os.File,
) (*mmath.MVec3, *mmath.MVec3) {
	// 回転軸
	linkAxis := effectorLocalPosition.Cross(ikLocalPosition).Normalize()

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile,
			"[%04d][%03d][%s][%05d][linkAxis] %s\n",
			frame, loop, linkBoneName, count-1, linkAxis.MMD().String(),
		)
	}

	// linkMat := linkQuat.ToMat4()
	// if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
	// 	fmt.Fprintf(ikFile,
	// 		"[%04d][%03d][%s][%05d][linkMat] %s (x: %s, y: %s, z: %s)\n",
	// 		frame, loop, linkBoneName, count-1, linkMat.String(), linkMat.AxisX().String(), linkMat.AxisY().String(), linkMat.AxisZ().String())
	// }

	if minAngleLimitRadians.IsOnlyX() || maxAngleLimitRadians.IsOnlyX() {
		// X軸のみの制限の場合
		vv := linkAxis.GetX()

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile,
				"[%04d][%03d][%s][%05d][linkAxis(X軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitXInv, linkAxis
		}
		return mmath.MVec3UnitX, linkAxis
	} else if minAngleLimitRadians.IsOnlyY() || maxAngleLimitRadians.IsOnlyY() {
		// Y軸のみの制限の場合
		vv := linkAxis.GetY()

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile,
				"[%04d][%03d][%s][%05d][linkAxis(Y軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitYInv, linkAxis
		}
		return mmath.MVec3UnitY, linkAxis
	} else if minAngleLimitRadians.IsOnlyZ() || maxAngleLimitRadians.IsOnlyZ() {
		// Z軸のみの制限の場合
		vv := linkAxis.GetZ()

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile,
				"[%04d][%03d][%s][%05d][linkAxis(Z軸制限)] vv: %.8f\n",
				frame, loop, linkBoneName, count-1, vv)
		}

		if vv < 0 {
			return mmath.MVec3UnitZInv, linkAxis
		}
		return mmath.MVec3UnitZ, linkAxis
	}

	return linkAxis, linkAxis
}

func (fs *BoneFrames) calcIkLimitQuaternion(
	totalIkQuat *mmath.MQuaternion, // リンクボーンの全体回転量
	minAngleLimitRadians *mmath.MVec3, // 最小軸制限（ラジアン）
	maxAngleLimitRadians *mmath.MVec3, // 最大軸制限（ラジアン）
	xAxisVector *mmath.MVec3, // X軸ベクトル
	yAxisVector *mmath.MVec3, // Y軸ベクトル
	zAxisVector *mmath.MVec3, // Z軸ベクトル
	loop int, // ループ回数
	loopCount int, // ループ総回数
	frame int, // キーフレーム
	count int, // デバッグ用: キーフレ位置
	linkBoneName string, // デバッグ用: リンクボーン名
	ikMotion *VmdMotion, // デバッグ用: IKモーション
	ikFile *os.File, // デバッグ用: IKファイル
) (*mmath.MQuaternion, int) {
	ikMat := totalIkQuat.ToMat4()
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile,
			"[%04d][%03d][%s][%05d][ikMat] %s (x: %s, y: %s, z: %s)\n",
			frame, loop, linkBoneName, count-1, ikMat.String(), ikMat.AxisX().String(), ikMat.AxisY().String(), ikMat.AxisZ().String())
	}

	// 軸回転角度を算出
	if minAngleLimitRadians.GetX() > -mmath.HALF_RAD && maxAngleLimitRadians.GetX() < mmath.HALF_RAD {
		// Z*X*Y順
		// X軸回り
		fSX := -ikMat.AxisZ().GetY() // sin(θx) = -m32
		fX := math.Asin(fSX)         // X軸回り決定
		fCX := math.Cos(fX)          // cos(θx)

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限] fSX: %f, fX: %f, fCX: %f\n",
				frame, loop, linkBoneName, count-1, fSX, fX, fCX)
		}

		// ジンバルロック回避
		if math.Abs(fX) > mmath.GIMBAL1_RAD {
			if fX < 0 {
				fX = -mmath.GIMBAL1_RAD
			} else {
				fX = mmath.GIMBAL1_RAD
			}
			fCX = math.Cos(fX)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-ジンバル] fSX: %f, fX: %f, fCX: %f\n",
					frame, loop, linkBoneName, count-1, fSX, fX, fCX)
			}
		}

		// Y軸回り
		fSY := ikMat.AxisZ().GetX() / fCX // sin(θy) = m31 / cos(θx)
		fCY := ikMat.AxisZ().GetZ() / fCX // cos(θy) = m33 / cos(θx)
		fY := math.Atan2(fSY, fCY)        // Y軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-Y軸回り] fSY: %f, fCY: %f, fY: %f\n",
				frame, loop, linkBoneName, count-1, fSY, fCY, fY)
		}

		// Z軸周り
		fSZ := ikMat.AxisX().GetY() / fCX // sin(θz) = m12 / cos(θx)
		fCZ := ikMat.AxisY().GetY() / fCX // cos(θz) = m22 / cos(θx)
		fZ := math.Atan2(fSZ, fCZ)        // Z軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-Z軸回り] fSZ: %f, fCZ: %f, fZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fCZ, fZ)
		}

		// 角度の制限
		fX = fs.getIkAxisValue(fX, minAngleLimitRadians.GetX(), maxAngleLimitRadians.GetX(), loop, loopCount,
			frame, count, "X軸制限-X", linkBoneName, ikMotion, ikFile)
		fY = fs.getIkAxisValue(fY, minAngleLimitRadians.GetY(), maxAngleLimitRadians.GetY(), loop, loopCount,
			frame, count, "X軸制限-Y", linkBoneName, ikMotion, ikFile)
		fZ = fs.getIkAxisValue(fZ, minAngleLimitRadians.GetZ(), maxAngleLimitRadians.GetZ(), loop, loopCount,
			frame, count, "X軸制限-Z", linkBoneName, ikMotion, ikFile)

		// 決定した角度でベクトルを回転
		xQuat := mmath.NewMQuaternionFromAxisAnglesRotate(xAxisVector, fX)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-xQuat] fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-yQuat] fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-zQuat] fZ: %f, zQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
		}

		return yQuat.Muled(xQuat).Muled(zQuat), count
	} else if minAngleLimitRadians.GetY() > -mmath.HALF_RAD && maxAngleLimitRadians.GetY() < mmath.HALF_RAD {
		// X*Y*Z順
		// Y軸回り
		fSY := -ikMat.AxisX().GetZ() // sin(θy) = m13
		fY := math.Asin(fSY)         // Y軸回り決定
		fCY := math.Cos(fY)          // cos(θy)

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限] fSY: %f, fY: %f, fCY: %f\n",
				frame, loop, linkBoneName, count-1, fSY, fY, fCY)
		}

		// ジンバルロック回避
		if math.Abs(fY) > mmath.GIMBAL1_RAD {
			if fY < 0 {
				fY = -mmath.GIMBAL1_RAD
			} else {
				fY = mmath.GIMBAL1_RAD
			}
			fCY = math.Cos(fY)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-ジンバル] fSY: %f, fY: %f, fCY: %f\n",
					frame, loop, linkBoneName, count-1, fSY, fY, fCY)
			}
		}

		// X軸回り
		fSX := ikMat.AxisY().GetZ() / fCY // sin(θx) = m23 / cos(θy)
		fCX := ikMat.AxisZ().GetZ() / fCY // cos(θx) = m33 / cos(θy)
		fX := math.Atan2(fSX, fCX)        // X軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-X軸回り] fSX: %f, fCX: %f, fX: %f\n",
				frame, loop, linkBoneName, count-1, fSX, fCX, fX)
		}

		// Z軸周り
		fSZ := ikMat.AxisX().GetY() / fCY // sin(θz) = m12 / cos(θy)
		fCZ := ikMat.AxisX().GetX() / fCY // cos(θz) = m11 / cos(θy)
		fZ := math.Atan2(fSZ, fCZ)        // Z軸回り決定

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-Z軸回り] fSZ: %f, fCZ: %f, fZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fCZ, fZ)
		}

		// 角度の制限
		fX = fs.getIkAxisValue(fX, minAngleLimitRadians.GetX(), maxAngleLimitRadians.GetX(), loop, loopCount,
			frame, count, "Y軸制限-X", linkBoneName, ikMotion, ikFile)
		fY = fs.getIkAxisValue(fY, minAngleLimitRadians.GetY(), maxAngleLimitRadians.GetY(), loop, loopCount,
			frame, count, "Y軸制限-Y", linkBoneName, ikMotion, ikFile)
		fZ = fs.getIkAxisValue(fZ, minAngleLimitRadians.GetZ(), maxAngleLimitRadians.GetZ(), loop, loopCount,
			frame, count, "Y軸制限-Z", linkBoneName, ikMotion, ikFile)

		// 決定した角度でベクトルを回転
		xQuat := mmath.NewMQuaternionFromAxisAnglesRotate(xAxisVector, fX)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-xQuat] fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-yQuat] fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-zQuat] fZ: %f, zQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
		}

		return zQuat.Muled(yQuat).Muled(xQuat), count
	}

	// Y*Z*X順
	// Z軸回り
	fSZ := -ikMat.AxisY().GetX() // sin(θz) = m21
	fZ := math.Asin(fSZ)         // Z軸回り決定
	fCZ := math.Cos(fZ)          // cos(θz)

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限] fSZ: %f, fZ: %f, fCZ: %f\n",
			frame, loop, linkBoneName, count-1, fSZ, fZ, fCZ)
	}

	// ジンバルロック回避
	if math.Abs(fZ) > mmath.GIMBAL1_RAD {
		if fZ < 0 {
			fZ = -mmath.GIMBAL1_RAD
		} else {
			fZ = mmath.GIMBAL1_RAD
		}
		fCZ = math.Cos(fZ)

		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-ジンバル] fSZ: %f, fZ: %f, fCZ: %f\n",
				frame, loop, linkBoneName, count-1, fSZ, fZ, fCZ)
		}
	}

	// X軸回り
	fSX := ikMat.AxisY().GetZ() / fCZ // sin(θx) = m23 / cos(θz)
	fCX := ikMat.AxisY().GetY() / fCZ // cos(θx) = m22 / cos(θz)
	fX := math.Atan2(fSX, fCX)        // X軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-X軸回り] fSX: %f, fCX: %f, fX: %f\n",
			frame, loop, linkBoneName, count-1, fSX, fCX, fX)
	}

	// Y軸周り
	fSY := ikMat.AxisZ().GetX() / fCZ // sin(θy) = m31 / cos(θz)
	fCY := ikMat.AxisX().GetX() / fCZ // cos(θy) = m11 / cos(θz)
	fY := math.Atan2(fSY, fCY)        // Y軸回り決定

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-Y軸回り] fSY: %f, fCY: %f, fY: %f\n",
			frame, loop, linkBoneName, count-1, fSY, fCY, fY)
	}

	// 角度の制限
	fX = fs.getIkAxisValue(fX, minAngleLimitRadians.GetX(), maxAngleLimitRadians.GetX(), loop, loopCount,
		frame, count, "Z軸制限-X", linkBoneName, ikMotion, ikFile)
	fY = fs.getIkAxisValue(fY, minAngleLimitRadians.GetY(), maxAngleLimitRadians.GetY(), loop, loopCount,
		frame, count, "Z軸制限-Y", linkBoneName, ikMotion, ikFile)
	fZ = fs.getIkAxisValue(fZ, minAngleLimitRadians.GetZ(), maxAngleLimitRadians.GetZ(), loop, loopCount,
		frame, count, "Z軸制限-Z", linkBoneName, ikMotion, ikFile)

	// 決定した角度でベクトルを回転
	xQuat := mmath.NewMQuaternionFromAxisAnglesRotate(xAxisVector, fX)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-xQuat] fX: %f, xQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, fX, xQuat.String(), xQuat.ToMMDDegrees().String())
	}

	yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-yQuat] fY: %f, yQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, fY, yQuat.String(), yQuat.ToMMDDegrees().String())
	}

	zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-zQuat] fZ: %f, zQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
	}

	return xQuat.Muled(zQuat).Muled(yQuat), count
}

func (fs *BoneFrames) getIkAxisValue(
	fV, minAngleLimit, maxAngleLimit float64,
	loop, loopCount int,
	frame int,
	count int,
	axisName, linkBoneName string,
	ikMotion *VmdMotion,
	ikFile *os.File,
) float64 {
	isLoopOver := float64(loop) < float64(loopCount)/2.0

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][%s-getIkAxisValue] loop: %d, loopCount: %d, float64(loopCount)/2.0: %f, isLoopOver: %v\n",
			frame, loop, linkBoneName, count-1, axisName, loop, loopCount, float64(loopCount)/2.0, isLoopOver)
	}

	if fV < minAngleLimit {
		tf := 2*minAngleLimit - fV
		if tf <= maxAngleLimit && isLoopOver {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][%s-最小角度(loop内)] minAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, minAngleLimit, fV, tf)
			}

			fV = tf
		} else {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][%s-最小角度(loop外)] minAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, minAngleLimit, fV, tf)
			}

			fV = minAngleLimit
		}
	}

	if fV > maxAngleLimit {
		tf := 2*maxAngleLimit - fV
		if tf >= minAngleLimit && isLoopOver {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][%s-最大角度(loop内)] maxAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, maxAngleLimit, fV, tf)
			}

			fV = tf
		} else {
			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][%s-最大角度(loop外)] maxAngleLimit: %f, fV: %f, tf: %f\n",
					frame, loop, linkBoneName, count-1, axisName, maxAngleLimit, fV, tf)
			}

			fV = maxAngleLimit
		}
	}

	return fV
}

func (fs *BoneFrames) calcBoneDeltas(
	frame int,
	model *pmx.PmxModel,
	deformBoneIndexes []int,
	boneDeltas *BoneDeltas,
) *BoneDeltas {
	for _, boneIndex := range deformBoneIndexes {
		delta := boneDeltas.Get(boneIndex)
		bone := model.Bones.Get(boneIndex)
		if delta == nil {
			delta = &BoneDelta{Bone: bone, Frame: frame}
		}

		delta.unitMatrix = mmath.NewMMat4()

		// スケール
		if delta.frameScale != nil && !delta.frameScale.IsOne() {
			delta.unitMatrix.Mul(delta.frameScale.ToScaleMat4())
		}

		// 回転
		rot := delta.LocalRotation()
		// rot := delta.frameRotation.Copy()
		// if rot == nil {
		// 	rot = mmath.NewMQuaternion()
		// }

		// isEffectorRot := false
		// isIkRot := false
		// if delta.frameIkRotation != nil && !delta.frameIkRotation.IsIdent() {
		// 	rot = delta.frameIkRotation.Muled(rot)
		// 	isIkRot = true
		// }
		// if delta.frameEffectRotation != nil && !delta.frameEffectRotation.IsIdent() {
		// 	rot = rot.Muled(delta.frameEffectRotation)
		// 	isEffectorRot = true
		// }
		// if (isIkRot || isEffectorRot) && delta.Bone.HasFixedAxis() {
		// 	// 軸制限回転を求める
		// 	rot = rot.ToFixedAxisRotation(delta.Bone.NormalizedFixedAxis)
		// }

		if rot != nil && !rot.IsIdent() {
			delta.unitMatrix.Mul(rot.ToMat4())
		}

		// 移動
		pos := delta.framePosition
		if pos == nil {
			pos = mmath.NewMVec3()
		}
		if delta.frameEffectPosition != nil && !delta.frameEffectPosition.IsZero() {
			pos.Add(delta.frameEffectPosition)
		}
		if pos != nil && !pos.IsZero() {
			delta.unitMatrix.Mul(pos.ToMat4())
		}

		// 逆BOf行列(初期姿勢行列)
		delta.unitMatrix.Mul(delta.Bone.RevertOffsetMatrix)
	}

	for _, boneIndex := range deformBoneIndexes {
		delta := boneDeltas.Get(boneIndex)
		parentDelta := boneDeltas.Get(delta.Bone.ParentIndex)
		if parentDelta != nil && parentDelta.globalMatrix != nil {
			delta.globalMatrix = delta.unitMatrix.Muled(parentDelta.globalMatrix)
		} else {
			// 対象ボーン自身の行列をかける
			delta.globalMatrix = delta.unitMatrix.Copy()
		}

		// BOf行列: 自身のボーンのボーンオフセット行列をかけてローカル行列
		delta.localMatrix = delta.Bone.OffsetMatrix.Muled(delta.globalMatrix)
		delta.globalPosition = delta.globalMatrix.Translation()
		boneDeltas.Append(delta)
	}

	return boneDeltas
}

// デフォーム対象ボーン情報一覧取得
func (fs *BoneFrames) createBoneDeltas(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	boneDeltas *BoneDeltas,
) ([]int, *BoneDeltas) {
	// mlog.Memory("createBoneDeltas 1)")

	isAfterPhysics := (boneDeltas != nil)

	// ボーン名の存在チェック用マップ
	targetSortedBones := model.Bones.LayerSortedBones[isAfterPhysics]

	if boneDeltas == nil {
		boneDeltas = NewBoneDeltas()
	}

	// mlog.Memory("createBoneDeltas 3)")
	deformBoneIndexes := make([]int, 0, len(targetSortedBones))

	if len(boneNames) > 0 {
		// 指定ボーンに関連するボーンのみ対象とする
		layerIndexes := make(pmx.LayerIndexes, 0, len(targetSortedBones))

		for _, boneName := range boneNames {
			if !model.Bones.ContainsName(boneName) {
				continue
			}

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
		// mlog.Memory("createBoneDeltas 4)")

		for _, layerIndex := range layerIndexes {
			bone := model.Bones.Get(layerIndex.Index)
			deformBoneIndexes = append(deformBoneIndexes, bone.Index)
			if !boneDeltas.Contains(bone.Index) {
				boneDeltas.Append(&BoneDelta{Bone: bone, Frame: frame})
			}
		}
	} else {
		// 変形階層・ボーンINDEXでソート
		for k := range len(targetSortedBones) {
			bone := targetSortedBones[k]
			deformBoneIndexes = append(deformBoneIndexes, bone.Index)
			boneDeltas.Append(&BoneDelta{Bone: bone, Frame: frame})
			if !boneDeltas.Contains(bone.Index) {
				boneDeltas.Append(&BoneDelta{Bone: bone, Frame: frame})
			}
		}
	}

	// mlog.Memory("createBoneDeltas 5)")

	return deformBoneIndexes, boneDeltas
}

// デフォーム情報を求めて設定
func (fs *BoneFrames) fillBoneDeform(
	frame int,
	model *pmx.PmxModel,
	deformBoneIndexes []int,
	boneDeltas *BoneDeltas,
) *BoneDeltas {
	for _, boneIndex := range deformBoneIndexes {
		delta := boneDeltas.Get(boneIndex)
		bone := model.Bones.Get(boneIndex)

		if delta == nil {
			delta = NewBoneDelta(bone, frame)
		}

		bf := fs.Get(bone.Name).Get(frame)
		if bf != nil {
			// ボーンの移動位置、回転角度、拡大率を取得
			delta.framePosition, delta.frameEffectPosition = fs.getPosition(bf, frame, bone, model, boneDeltas, 0)
			delta.frameRotation, delta.frameEffectRotation = fs.getRotation(bf, frame, bone, model, boneDeltas, 0)
			delta.frameScale = fs.getScale(bf, bone, boneDeltas)
		}
		boneDeltas.Append(delta)
	}

	return boneDeltas
}

// 該当キーフレにおけるボーンの移動位置
func (fs *BoneFrames) getPosition(
	bf *BoneFrame,
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	loop int,
) (*mmath.MVec3, *mmath.MVec3) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMVec3(), nil
	}

	var pos *mmath.MVec3
	if boneDeltas != nil && boneDeltas.Get(bone.Index) != nil && boneDeltas.Get(bone.Index).framePosition != nil {
		pos = boneDeltas.Get(bone.Index).framePosition.Copy()
	} else if bf.Position != nil && !bf.Position.IsZero() {
		pos = bf.Position.Copy()
	} else {
		pos = mmath.NewMVec3()
	}

	if bf.MorphPosition != nil && !bf.MorphPosition.IsZero() {
		pos.Add(bf.MorphPosition)
	}

	if bone.IsEffectorTranslation() {
		// 付与親ありの場合、外部親位置を取得する
		effectPos := fs.getEffectPosition(frame, bone, model, boneDeltas, loop+1)
		return pos, effectPos
	}

	return pos, nil
}

// 付与親を加味した移動位置
func (fs *BoneFrames) getEffectPosition(
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
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

	if boneDeltas != nil {
		effectDelta := boneDeltas.Get(effectBone.Index)
		if effectDelta != nil && effectDelta.framePosition != nil {
			return effectDelta.framePosition.MuledScalar(bone.EffectFactor)
		}
	}

	bf := fs.Get(effectBone.Name).Get(frame)
	if bf == nil {
		return mmath.NewMVec3()
	}

	pos, effectPos := fs.getPosition(bf, frame, effectBone, model, boneDeltas, loop+1)

	if effectPos == nil {
		return pos.MuledScalar(bone.EffectFactor)
	}

	return pos.Added(effectPos).MulScalar(bone.EffectFactor)
}

// 該当キーフレにおけるボーンの回転角度
func (fs *BoneFrames) getRotation(
	bf *BoneFrame,
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	loop int,
) (*mmath.MQuaternion, *mmath.MQuaternion) {
	if loop > 20 {
		// 無限ループを避ける
		return mmath.NewMQuaternion(), nil
	}

	// FK(捩り) > IK(捩り) > 付与親(捩り)
	var rot *mmath.MQuaternion
	if boneDeltas != nil && boneDeltas.Get(bone.Index) != nil && boneDeltas.Get(bone.Index).frameRotation != nil {
		rot = boneDeltas.Get(bone.Index).frameRotation.Copy()
	} else if bf.Rotation != nil && !bf.Rotation.IsIdent() {
		rot = bf.Rotation.Copy()
	} else {
		rot = mmath.NewMQuaternion()
	}

	if bf.MorphRotation != nil {
		rot.Mul(bf.MorphRotation)
	}

	if bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
	}

	if bone.IsEffectorRotation() {
		// 付与親ありの場合、外部親回転を取得する
		effectRot := fs.getEffectRotation(frame, bone, model, boneDeltas, loop+1)
		return rot, effectRot
	}

	return rot, nil
}

// 付与親を加味した回転角度
func (fs *BoneFrames) getEffectRotation(
	frame int,
	bone *pmx.Bone,
	model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	loop int,
) *mmath.MQuaternion {
	if bone.EffectFactor == 0 || loop > 20 {
		// 付与率が0の場合、常に0になる
		// MMDエンジン対策で無限ループを避ける
		return nil
	}

	if !(bone.EffectIndex > 0 && model.Bones.Contains(bone.EffectIndex)) {
		// 付与親が存在しない場合、常に0になる
		return nil
	}

	// 付与親が存在する場合、付与親の回転角度を掛ける
	effectBone := model.Bones.Get(bone.EffectIndex)

	bf := fs.Get(effectBone.Name).Get(frame)
	if bf == nil {
		return nil
	}

	rot, effectRot := fs.getRotation(bf, frame, effectBone, model, boneDeltas, loop+1)

	// 付与に対する付与は出来ない
	if effectRot != nil {
		// rot.Mul(effectRot)
		rot = effectRot.Muled(rot)
	}

	return rot.MuledScalar(bone.EffectFactor)
}

// 該当キーフレにおけるボーンの拡大率
func (fs *BoneFrames) getScale(
	bf *BoneFrame,
	bone *pmx.Bone,
	boneDeltas *BoneDeltas,
) *mmath.MVec3 {

	scale := &mmath.MVec3{1, 1, 1}
	if boneDeltas != nil && boneDeltas.Get(bone.Index) != nil && boneDeltas.Get(bone.Index).frameScale != nil {
		scale = boneDeltas.Get(bone.Index).frameScale
	} else if bf.Scale != nil && !bf.Scale.IsZero() {
		scale.Add(bf.Scale)
	}

	if bf.MorphScale != nil && !bf.MorphScale.IsZero() {
		return scale.Add(bf.MorphScale)
	}

	return scale
}
