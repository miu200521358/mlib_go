package vmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
	return fs.DeformByPhysicsFlag(frame, model, boneNames, isCalcIk, beforeBoneDeltas, nil, ikFrame, false)
}

func (fs *BoneFrames) DeformByPhysicsFlag(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	beforeBoneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
	ikFrame *IkFrame,
	isAfterPhysics bool,
) *BoneDeltas {
	// mlog.Memory(fmt.Sprintf("Deform 1)frame: %d", frame))
	deformBoneIndexes, boneDeltas := fs.prepareDeltas(frame, model, boneNames, isCalcIk,
		beforeBoneDeltas, morphDeltas, ikFrame, isAfterPhysics)
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
	morphDeltas *MorphDeltas,
	ikFrame *IkFrame,
	isAfterPhysics bool,
) ([]int, *BoneDeltas) {
	// mlog.Memory(fmt.Sprintf("prepareDeltas 1)frame: %d", frame))

	deformBoneIndexes, boneDeltas := fs.createBoneDeltas(frame, model, boneNames, beforeBoneDeltas, isAfterPhysics)

	// mlog.Memory(fmt.Sprintf("prepareDeltas 2)frame: %d", frame))

	// IK事前計算
	if isCalcIk {
		// ボーン変形行列操作
		boneDeltas = fs.prepareIk(frame, model, deformBoneIndexes, boneDeltas, morphDeltas, ikFrame, isAfterPhysics)
	}

	// mlog.Memory(fmt.Sprintf("prepareDeltas 3)frame: %d", frame))

	// ボーンデフォーム情報を埋める
	boneDeltas = fs.fillBoneDeform(frame, model, deformBoneIndexes, boneDeltas, morphDeltas)

	// mlog.Memory(fmt.Sprintf("prepareDeltas 4)frame: %d", frame))

	return deformBoneIndexes, boneDeltas
}

// IK事前計算処理
func (fs *BoneFrames) prepareIk(
	frame int,
	model *pmx.PmxModel,
	deformBoneIndexes []int,
	boneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
	ikFrame *IkFrame,
	isAfterPhysics bool,
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

				boneDeltas = fs.calcIk(frame, ikBone, model, boneDeltas, morphDeltas, isAfterPhysics, ikFrame, prefixPath)
			}
		}
	}

	return boneDeltas
}

// IK計算
func (fs *BoneFrames) calcIk(
	frame int,
	ikBone *pmx.Bone,
	model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
	isAfterPhysics bool,
	ikFrame *IkFrame,
	prefixPath string,
) *BoneDeltas {
	if len(ikBone.Ik.Links) < 1 {
		// IKリンクが無ければスルー
		return boneDeltas
	}

	var err error
	var ikFile *os.File
	var ikMotion *VmdMotion
	count := 1

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
			ikMotion.Save("", "")
		}
		if ikFile != nil {
			ikFile.Close()
		}
	}()

	// つま先ＩＫであるか
	isToeIk := strings.Contains(ikBone.Name, "つま先ＩＫ")
	// 一段IKであるか
	isSingleIk := len(ikBone.Ik.Links) == 1

	// ループ回数
	loopCount := max(ikBone.Ik.LoopCount, 1)
	if isToeIk {
		// つま先IKの場合、初回に足首位置に向けるのでループ1回分余分に回す
		loopCount += 1
	}

	// IKターゲットボーン
	effectorBone := model.Bones.Get(ikBone.Ik.BoneIndex)
	// IK関連の行列を一括計算
	ikDeltas := fs.DeformByPhysicsFlag(frame, model, []string{ikBone.Name, effectorBone.Name}, false,
		boneDeltas, nil, ikFrame, false)
	if isAfterPhysics {
		// 物理後の場合は物理後のも取得する
		ikDeltas = fs.DeformByPhysicsFlag(frame, model, []string{ikBone.Name, effectorBone.Name}, false,
			ikDeltas, nil, ikFrame, true)
	}
	if !ikDeltas.Contains(ikBone.Index) {
		// IKボーンが存在しない場合、スルー
		return boneDeltas
	}

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		ikOffMotion := NewVmdMotion(fmt.Sprintf("%s_0_%s.vmd", prefixPath, ikBone.Name))

		bif := NewIkFrame(0)
		bif.Registered = true

		for _, bone := range model.Bones.Data {
			if bone.IsIK() {
				ef := NewIkEnableFrame(0)
				ef.Registered = true
				ef.BoneName = bone.Name
				ef.Enabled = false

				bif.IkList = append(bif.IkList, ef)
			}
		}

		ikOffMotion.AppendIkFrame(bif)

		for _, ikDelta := range ikDeltas.Data {
			if ikDelta == nil {
				continue
			}
			bf := NewBoneFrame(0)
			bf.Position = ikDelta.framePosition
			bf.Rotation = ikDelta.frameRotation
			ikOffMotion.AppendRegisteredBoneFrame(ikDelta.Bone.Name, bf)
		}

		ikOffMotion.Save("IK OFF", "")
	}

	var ikOffDeltas *BoneDeltas
	if isToeIk {
		ikOffDeltas = fs.DeformByPhysicsFlag(frame, model, []string{effectorBone.Name}, false,
			nil, nil, ikFrame, isAfterPhysics)
		if !ikOffDeltas.Contains(effectorBone.Index) {
			// IK OFFボーンが存在しない場合、スルー
			return boneDeltas
		}
	}

	// エフェクタ関連情報取得
	effectorDeformBoneIndexes, boneDeltas :=
		fs.prepareDeltas(frame, model, []string{effectorBone.Name}, false, boneDeltas, nil, ikFrame, false)
	if isAfterPhysics {
		// 物理後の場合は物理後のも取得する
		effectorDeformBoneIndexes, boneDeltas =
			fs.prepareDeltas(frame, model, []string{effectorBone.Name}, false, boneDeltas, nil, ikFrame, true)
	}
	if !boneDeltas.Contains(effectorBone.Index) || !boneDeltas.Contains(ikBone.Index) ||
		!boneDeltas.Contains(ikBone.Ik.BoneIndex) {
		// エフェクタボーンが存在しない場合、スルー
		return boneDeltas
	}

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

			// リンクボーンの変形情報を取得
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

			// 初回にIK事前計算
			if loop == 0 && isToeIk && ikOffDeltas != nil {
				// IK OFF 時の IKターゲットボーンのグローバル位置を取得
				ikGlobalPosition = ikOffDeltas.Get(effectorBone.Index).GlobalPosition()
				// 現在のIKターゲットボーンのグローバル位置を取得
				effectorGlobalPosition = ikDeltas.Get(effectorBone.Index).GlobalPosition()
			}

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
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][Local] effectorLocalPosition: %s, ikLocalPosition: %s (%f)\n",
					frame, loop, linkBone.Name, count-1,
					effectorLocalPosition.MMD().String(), ikLocalPosition.MMD().String(),
					effectorLocalPosition.Distance(ikLocalPosition))
			}

			distanceThreshold := effectorLocalPosition.Distance(ikLocalPosition)
			if distanceThreshold < 1e-5 {
				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][Local] ***BREAK*** distanceThreshold: %f\n",
						frame, loop, linkBone.Name, count-1, distanceThreshold)
				}

				break ikLoop
			}

			effectorLocalPosition.Normalize()
			ikLocalPosition.Normalize()

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
			originalLinkAngle := math.Acos(mmath.ClampedFloat(linkDot, -1, 1))
			linkAngle := mmath.ClampedFloat(originalLinkAngle, -unitRad, unitRad)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][単位角制限] linkAngle: %.8f(%.5f), originalLinkAngle: %.8f(%.5f)\n",
					frame, loop, linkBone.Name, count-1, linkAngle, mmath.ToDegree(linkAngle),
					originalLinkAngle, mmath.ToDegree(originalLinkAngle),
				)
			}

			// 角度がほとんどない場合
			angleThreshold := math.Abs(linkAngle)
			if angleThreshold < 1e-5 {
				if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
					fmt.Fprintf(ikFile,
						"[%04d][%03d][%s][%05d][Local] ***BREAK*** angleThreshold: %f\n",
						frame, loop, linkBone.Name, count-1, angleThreshold)
				}

				break ikLoop
			}

			// 回転軸
			var originalLinkAxis, linkAxis *mmath.MVec3
			// 一段IKでない場合、または一段IKでかつ回転角が88度以上の場合
			if !isSingleIk || (isSingleIk && linkAngle > mmath.GIMBAL1_RAD) && ikLink.AngleLimit {
				// グローバル軸制限
				linkAxis, originalLinkAxis = fs.getLinkAxis(
					ikLink.MinAngleLimit.GetRadians(),
					ikLink.MaxAngleLimit.GetRadians(),
					effectorLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name, ikMotion, ikFile,
				)
			} else if !isSingleIk || (isSingleIk && linkAngle > mmath.GIMBAL1_RAD) && ikLink.LocalAngleLimit {
				// ローカル軸制限
				linkAxis, originalLinkAxis = fs.getLinkAxis(
					ikLink.LocalMinAngleLimit.GetRadians(),
					ikLink.LocalMaxAngleLimit.GetRadians(),
					effectorLocalPosition, ikLocalPosition,
					frame, count, loop, linkBone.Name, ikMotion, ikFile,
				)
			} else {
				// 軸制限なし or 一段IKでかつ回転角が88度未満の場合
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

			originalIkQuat := mmath.NewMQuaternionFromAxisAnglesRotate(originalLinkAxis, originalLinkAngle)
			ikQuat := mmath.NewMQuaternionFromAxisAnglesRotate(linkAxis, linkAngle)

			originalTotalIkQuat := linkQuat.Muled(originalIkQuat)
			totalIkQuat := linkQuat.Muled(ikQuat)

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

			if loop == 0 && morphDeltas != nil && morphDeltas.Bones != nil &&
				morphDeltas.Bones.Get(linkBone.Index) != nil &&
				morphDeltas.Bones.Get(linkBone.Index).frameRotation != nil {
				// モーフ変形がある場合、モーフ変形を追加適用
				resultIkQuat = resultIkQuat.Muled(morphDeltas.Bones.Get(linkBone.Index).frameRotation)
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
			boneDeltas.Update(linkDelta)

			if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
				bf := NewBoneFrame(count)
				bf.Rotation = linkDelta.LocalRotation().Copy()
				ikMotion.AppendRegisteredBoneFrame(linkBone.Name, bf)
				count++

				fmt.Fprintf(ikFile,
					"[%04d][%03d][%s][%05d][結果] bf.Rotation: %s(%s)\n",
					frame, loop, linkBone.Name, count-1, bf.Rotation.String(), bf.Rotation.ToMMDDegrees().String())
			}
		}
	}

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
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][X軸制限-zQuat] zAxisVector: %s, fZ: %f, zQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, zAxisVector.String(), fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
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
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
		}

		yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
		}

		zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
		if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
			fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Y軸制限-zQuat] zAxisVector: %s, fZ: %f, zQuat: %s(%s)\n",
				frame, loop, linkBoneName, count-1, zAxisVector.String(), fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
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
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-xQuat] xAxisVector: %s, fX: %f, xQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, xAxisVector.String(), fX, xQuat.String(), xQuat.ToMMDDegrees().String())
	}

	yQuat := mmath.NewMQuaternionFromAxisAnglesRotate(yAxisVector, fY)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-yQuat] yAxisVector: %s, fY: %f, yQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, yAxisVector.String(), fY, yQuat.String(), yQuat.ToMMDDegrees().String())
	}

	zQuat := mmath.NewMQuaternionFromAxisAnglesRotate(zAxisVector, fZ)
	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][Z軸制限-zQuat] zAxisVector: %s, fZ: %f, zQuat: %s(%s)\n",
			frame, loop, linkBoneName, count-1, zAxisVector.String(), fZ, zQuat.String(), zQuat.ToMMDDegrees().String())
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
	isInLoop := float64(loop) < float64(loopCount)/2.0

	if mlog.IsIkVerbose() && ikMotion != nil && ikFile != nil {
		fmt.Fprintf(ikFile, "[%04d][%03d][%s][%05d][%s-getIkAxisValue] loop: %d, isInLoop: %v\n",
			frame, loop, linkBoneName, count-1, axisName, loop, isInLoop)
	}

	if fV < minAngleLimit {
		tf := 2*minAngleLimit - fV
		if tf <= maxAngleLimit && isInLoop {
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
		if tf >= minAngleLimit && isInLoop {
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
		delta.globalMatrix = nil
		delta.localMatrix = nil
		delta.globalPosition = nil

		// スケール
		if delta.frameScale != nil && !delta.frameScale.IsOne() {
			delta.unitMatrix.Mul(delta.frameScale.ToScaleMat4())
		}

		// 回転
		rot := boneDeltas.LocalRotation(bone.Index, 0)
		if rot != nil && !rot.IsIdent() {
			delta.unitMatrix.Mul(rot.ToMat4())
		}

		// 移動
		pos := boneDeltas.LocalPosition(bone.Index, 0)
		if pos != nil && !pos.IsZero() {
			delta.unitMatrix.Mul(pos.ToMat4())
		}

		// x := math.Abs(rot.GetX())

		// if bone.Name == "左袖_後_赤_04_04" {
		// 	mlog.I("[%s][%04d]: pos: %s, rot: %s(%s), x: %f\n", bone.Name, frame, pos.String(), rot.String(), rot.ToMMDDegrees().String(), x)
		// }

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

		// delta.localMatrix = delta.Bone.OffsetMatrix.Muled(delta.globalMatrix)
		// delta.globalPosition = delta.globalMatrix.Translation()
		boneDeltas.Update(delta)
	}

	return boneDeltas
}

// デフォーム対象ボーン情報一覧取得
func (fs *BoneFrames) createBoneDeltas(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	boneDeltas *BoneDeltas,
	isAfterPhysics bool,
) ([]int, *BoneDeltas) {
	// ボーン名の存在チェック用マップ
	targetSortedBones := model.Bones.LayerSortedBones[isAfterPhysics]

	if boneDeltas == nil {
		boneDeltas = NewBoneDeltas(model.Bones)
	}

	// 変形階層順ボーンIndexリスト
	deformBoneIndexes := make([]int, 0, len(targetSortedBones))

	// 関連ボーンINDEXリスト（順不同）
	relativeBoneIndexes := make([]int, 0)

	if len(boneNames) > 0 {
		// 指定ボーンに関連するボーンのみ対象とする

		for _, boneName := range boneNames {
			if !model.Bones.ContainsName(boneName) {
				continue
			}

			// ボーン
			bone := model.Bones.GetByName(boneName)

			// 対象のボーンは常に追加
			if !slices.Contains(relativeBoneIndexes, bone.Index) {
				relativeBoneIndexes = append(relativeBoneIndexes, bone.Index)
			}

			// 関連するボーンの追加
			for _, index := range bone.RelativeBoneIndexes {
				if !slices.Contains(relativeBoneIndexes, index) {
					relativeBoneIndexes = append(relativeBoneIndexes, index)
				}
			}
		}
	} else {
		// ボーン名の指定が無い場合、全ボーンを対象とする
		for _, bone := range targetSortedBones {
			// 対象のボーンは常に追加
			if !slices.Contains(relativeBoneIndexes, bone.Index) {
				relativeBoneIndexes = append(relativeBoneIndexes, bone.Index)
			}

			// 関連するボーンの追加
			for _, index := range bone.RelativeBoneIndexes {
				if !slices.Contains(relativeBoneIndexes, index) {
					relativeBoneIndexes = append(relativeBoneIndexes, index)
				}
			}
		}
	}

	// 変形階層・ボーンINDEXでソート
	for _, ap := range []bool{false, true} {
		for _, bone := range model.Bones.LayerSortedBones[ap] {
			if slices.Contains(relativeBoneIndexes, bone.Index) {
				deformBoneIndexes = append(deformBoneIndexes, bone.Index)
				if !boneDeltas.Contains(bone.Index) {
					boneDeltas.Update(&BoneDelta{Bone: bone, Frame: frame})
				}
			}
		}
	}

	return deformBoneIndexes, boneDeltas
}

// デフォーム情報を求めて設定
func (fs *BoneFrames) fillBoneDeform(
	frame int,
	model *pmx.PmxModel,
	deformBoneIndexes []int,
	boneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
) *BoneDeltas {
	for _, boneIndex := range deformBoneIndexes {
		bone := model.Bones.Get(boneIndex)
		delta := boneDeltas.Get(boneIndex)

		var bf *BoneFrame
		if bone.IsAfterPhysicsDeform() || boneDeltas == nil || boneDeltas.Get(bone.Index) == nil ||
			boneDeltas.Get(bone.Index).framePosition == nil ||
			boneDeltas.Get(bone.Index).frameRotation == nil ||
			boneDeltas.Get(bone.Index).frameScale == nil {
			bf = fs.Get(bone.Name).Get(frame)
		}
		// ボーンの移動位置、回転角度、拡大率を取得
		delta.framePosition, delta.frameMorphPosition = fs.getPosition(bf, bone, boneDeltas, morphDeltas)
		delta.frameRotation, delta.frameMorphRotation = fs.getRotation(bf, bone, boneDeltas, morphDeltas)
		delta.frameScale = fs.getScale(bf, bone, boneDeltas, morphDeltas)
		boneDeltas.Update(delta)
	}

	return boneDeltas
}

// 該当キーフレにおけるボーンの移動位置
func (fs *BoneFrames) getPosition(
	bf *BoneFrame,
	// frame int,
	bone *pmx.Bone,
	// model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
) (*mmath.MVec3, *mmath.MVec3) {
	var pos *mmath.MVec3
	if boneDeltas != nil && boneDeltas.Get(bone.Index) != nil && boneDeltas.Get(bone.Index).framePosition != nil {
		pos = boneDeltas.Get(bone.Index).framePosition.Copy()
	} else if bf != nil && bf.Position != nil && !bf.Position.IsZero() {
		pos = bf.Position.Copy()
	} else {
		pos = mmath.NewMVec3()
	}

	var morphPos *mmath.MVec3
	if morphDeltas != nil && morphDeltas.Bones.Get(bone.Index) != nil &&
		morphDeltas.Bones.Get(bone.Index).framePosition != nil {
		morphPos = morphDeltas.Bones.Get(bone.Index).framePosition
	}

	return pos, morphPos
}

// 該当キーフレにおけるボーンの回転角度
func (fs *BoneFrames) getRotation(
	bf *BoneFrame,
	// frame int,
	bone *pmx.Bone,
	// model *pmx.PmxModel,
	boneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
) (*mmath.MQuaternion, *mmath.MQuaternion) {
	// FK(捩り) > IK(捩り) > 付与親(捩り)
	var rot *mmath.MQuaternion
	var morphRot *mmath.MQuaternion
	if boneDeltas != nil && boneDeltas.Get(bone.Index) != nil && boneDeltas.Get(bone.Index).frameRotation != nil {
		rot = boneDeltas.Get(bone.Index).frameRotation.Copy()
	} else {
		if bf != nil && bf.Rotation != nil && !bf.Rotation.IsIdent() {
			rot = bf.Rotation.Copy()
		} else {
			rot = mmath.NewMQuaternion()

			if morphDeltas != nil && morphDeltas.Bones.Get(bone.Index) != nil &&
				morphDeltas.Bones.Get(bone.Index).frameRotation != nil {
				// IKの場合はIK計算時に組み込まれているので、まだframeRotationが無い場合のみ加味
				morphRot = morphDeltas.Bones.Get(bone.Index).frameRotation
				// mlog.I("[%s][%04d][%d]: rot: %s(%s), morphRot: %s(%s)\n", bone.Name, frame, loop,
				// 	rot.String(), rot.ToMMDDegrees().String(), morphRot.String(), morphRot.ToMMDDegrees().String())
			}
		}
	}

	if bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
	}

	return rot, morphRot
}

// 該当キーフレにおけるボーンの拡大率
func (fs *BoneFrames) getScale(
	bf *BoneFrame,
	bone *pmx.Bone,
	boneDeltas *BoneDeltas,
	morphDeltas *MorphDeltas,
) *mmath.MVec3 {

	scale := &mmath.MVec3{1, 1, 1}
	if boneDeltas != nil && boneDeltas.Get(bone.Index) != nil &&
		boneDeltas.Get(bone.Index).frameScale != nil {
		scale = boneDeltas.Get(bone.Index).frameScale
	} else if bf != nil && bf.Scale != nil && !bf.Scale.IsZero() {
		scale.Add(bf.Scale)
	}

	if morphDeltas != nil && morphDeltas.Bones.Get(bone.Index) != nil &&
		morphDeltas.Bones.Get(bone.Index).frameScale != nil {
		return scale.Add(morphDeltas.Bones.Get(bone.Index).frameScale)
	}

	return scale
}
