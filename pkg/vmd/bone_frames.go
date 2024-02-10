package vmd

import (
	"math"
	"slices"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mmath"
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
	GIMBAL_RAD  = math.Pi * 88 / 180
	GIMBAL2_RAD = math.Pi * 88 * 2 / 180
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

func (bfs *BoneFrames) Animate(
	frames []float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// 処理対象ボーン一覧取得
	targetBoneNames := bfs.getAnimatedBoneNames(model, boneNames)

	// IK事前計算
	if isCalcIk {
		bfs.prepareIkSolvers(frames, model, targetBoneNames, isOutLog, description)
	}

	// ボーン変形行列操作
	positions, rotations, scales, _ :=
		bfs.getBoneMatrixes(frames, model, targetBoneNames, isCalcIk, isOutLog, description)

	// ボーン行列計算
	return bfs.calcBoneMatrixes(
		frames,
		model,
		targetBoneNames,
		positions,
		rotations,
		scales,
		isOutLog,
		description,
	)
}

// IK事前計算処理
func (bfs *BoneFrames) prepareIkSolvers(
	frames []float32,
	model *pmx.PmxModel,
	targetBoneNames []string,
	isOutLog bool,
	description string,
) {
	// IKリンクに設定されているボーンが対象になっている場合、該当のIKボーンINDEXをリストに追加する
	ikBoneIndexes := make([]int, 0)
	for _, boneName := range targetBoneNames {
		bone := model.Bones.GetItemByName(boneName)
		for _, linkBoneIndex := range bone.IkLinkBoneIndexes {
			if !slices.Contains(ikBoneIndexes, linkBoneIndex) {
				ikBoneIndexes = append(ikBoneIndexes, linkBoneIndex)
			}
		}
	}

	// 念のためソート
	slices.Sort(ikBoneIndexes)

	for _, frame := range frames {
		for _, ikBoneIndex := range ikBoneIndexes {
			// 各フレームでIK計算
			quats, effectorTargetBoneNames := bfs.calcIk(frame, ikBoneIndex, model, isOutLog, description)

			for _, ikLink := range model.Bones.GetItem(ikBoneIndex).Ik.Links {
				// IKリンクボーンの回転量を更新
				linkBf := bfs.GetItem(model.Bones.GetItem(ikLink.BoneIndex).Name).GetItem(frame)
				linkIndex := slices.Index(effectorTargetBoneNames, model.Bones.GetItem(ikLink.BoneIndex).Name)
				linkBf.IkRotation.SetQuaternion(quats[linkIndex])

				// IK用なので登録フラグは既存のままで追加して補間曲線は分割しない
				bfs.GetItem(model.Bones.GetItem(ikLink.BoneIndex).Name).Append(linkBf)
			}
		}
	}
}

// IK計算
func (bfs *BoneFrames) calcIk(
	frame float32,
	boneIndex int,
	model *pmx.PmxModel,
	isOutLog bool,
	description string,
) ([]*mmath.MQuaternion, []string) {
	// IKボーン
	ikBone := model.Bones.GetItem(boneIndex)
	// IKターゲットボーン
	effectorBone := model.Bones.GetItem(ikBone.Ik.BoneIndex)
	// IK関連の行列を一括計算
	ikMatrixes := bfs.Animate([]float32{frame}, model, []string{ikBone.Name}, false, false, "")
	// 処理対象ボーン名取得
	effectorTargetBoneNames := bfs.getAnimatedBoneNames(model, []string{effectorBone.Name})
	// エフェクタボーンの関連ボーンの初期値を取得
	positions, rotations, scales, quats :=
		bfs.getBoneMatrixes([]float32{frame}, model, effectorTargetBoneNames, true, false, "")

		// IK計算
ikLoop:
	for loop := 0; loop < ikBone.Ik.LoopCount; loop++ {
		for _, ikLink := range ikBone.Ik.Links {
			// ikLink は末端から並んでる
			if !model.Bones.Contains(ikLink.BoneIndex) {
				continue
			}

			// 処理対象IKリンクボーン
			linkBone := model.Bones.GetItem(ikLink.BoneIndex)
			linkIndex := slices.Index(effectorTargetBoneNames, linkBone.Name)

			// 角度制限があってまったく動かさない場合、IK計算しないで次に行く
			if (linkBone.AngleLimit &&
				linkBone.MinAngleLimit.GetRadians().IsZero() &&
				linkBone.MaxAngleLimit.GetRadians().IsZero()) ||
				(linkBone.LocalAngleLimit &&
					linkBone.LocalMinAngleLimit.GetRadians().IsZero() &&
					linkBone.LocalMaxAngleLimit.GetRadians().IsZero()) {
				continue
			}

			// IK関連の行列を取得
			linkMatrixes := bfs.calcBoneMatrixes(
				[]float32{frame},
				model,
				effectorTargetBoneNames,
				positions,
				rotations,
				scales,
				isOutLog,
				description)

			// IKボーンのグローバル位置
			ikGlobalPosition := ikMatrixes.GetItem(ikBone.Name, frame).Position

			// 現在のIKターゲットボーンのグローバル位置を取得
			effectorGlobalPosition := linkMatrixes.GetItem(effectorBone.Name, frame).Position

			// 注目ノード（実際に動かすボーン=リンクボーン）
			linkMatrix := linkMatrixes.GetItem(linkBone.Name, frame).GlobalMatrix
			// ワールド座標系から注目ノードの局所座標系への変換
			linkInvMatrix := linkMatrix.Inverse()

			// 注目ノードを起点とした、エフェクタのローカル位置
			effectorLocalPosition := linkInvMatrix.MulVec3(effectorGlobalPosition)
			// 注目ノードを起点とした、IK目標のローカル位置
			ikLocalPosition := linkInvMatrix.MulVec3(ikGlobalPosition)

			// 位置の差がほとんどない場合、終了
			if ikLocalPosition.Distance(effectorLocalPosition) < 1e-8 {
				break ikLoop
			}

			normalizedEffectorLocalPosition := effectorLocalPosition.Normalized()
			normalizedIkLocalPosition := ikLocalPosition.Normalized()

			// ベクトル (1) を (2) に一致させるための最短回転量（Axis-Angle）
			// 回転軸
			axis := normalizedEffectorLocalPosition.Cross(normalizedIkLocalPosition).Normalize()
			// 回転角(ラジアン)
			angle := math.Acos(mmath.ClampFloat(
				normalizedIkLocalPosition.Dot(normalizedEffectorLocalPosition)/
					(normalizedIkLocalPosition.Length()*normalizedEffectorLocalPosition.Length()), 0, 1))

			// リンクボーンの角度を取得
			linkQuat := quats[0][linkIndex]
			var totalActualIkQuat *mmath.MQuaternion

			if ikLink.AngleLimit {
				// 角度制限が入ってる場合
				if ikLink.MinAngleLimit.GetRadians().GetX() != 0 ||
					ikLink.MaxAngleLimit.GetRadians().GetX() != 0 {
					// グローバルX: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat = bfs.calculateSingleAxisRadRotation(
						ikLink.MinAngleLimit, ikLink.MaxAngleLimit,
						linkQuat, axis, angle, 0, &mmath.MVec3{1, 0, 0},
						ikBone.Ik.UnitRotation.GetRadians().GetX(), false)
				} else if ikLink.MinAngleLimit.GetRadians().GetY() != 0 ||
					ikLink.MaxAngleLimit.GetRadians().GetY() != 0 {
					// グローバルY: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat = bfs.calculateSingleAxisRadRotation(
						ikLink.MinAngleLimit, ikLink.MaxAngleLimit,
						linkQuat, axis, angle, 1, &mmath.MVec3{0, 1, 0},
						ikBone.Ik.UnitRotation.GetRadians().GetY(), false)
				} else if ikLink.MinAngleLimit.GetRadians().GetZ() != 0 ||

					ikLink.MaxAngleLimit.GetRadians().GetZ() != 0 {
					// グローバルZ: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat = bfs.calculateSingleAxisRadRotation(
						ikLink.MinAngleLimit, ikLink.MaxAngleLimit,
						linkQuat, axis, angle, 2, &mmath.MVec3{0, 0, 1},
						ikBone.Ik.UnitRotation.GetRadians().GetZ(), false)
				}
			} else if ikLink.LocalAngleLimit {
				// ローカル軸角度制限が入っている場合、ローカル軸に合わせて理想回転を求める
				if ikLink.LocalMinAngleLimit.GetRadians().GetX() != 0 ||
					ikLink.LocalMaxAngleLimit.GetRadians().GetX() != 0 {
					// ローカルX: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat = bfs.calculateSingleAxisRadRotation(
						ikLink.LocalMinAngleLimit, ikLink.LocalMaxAngleLimit,
						linkQuat, axis, angle, 0, linkBone.NormalizedLocalAxisX,
						ikBone.Ik.UnitRotation.GetRadians().GetX(), true)
				} else if ikLink.LocalMinAngleLimit.GetRadians().GetY() != 0 ||
					ikLink.LocalMaxAngleLimit.GetRadians().GetY() != 0 {
					// ローカルY: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat = bfs.calculateSingleAxisRadRotation(
						ikLink.LocalMinAngleLimit, ikLink.LocalMaxAngleLimit,
						linkQuat, axis, angle, 1, linkBone.NormalizedLocalAxisY,
						ikBone.Ik.UnitRotation.GetRadians().GetY(), true)
				} else if ikLink.LocalMinAngleLimit.GetRadians().GetZ() != 0 ||
					ikLink.LocalMaxAngleLimit.GetRadians().GetZ() != 0 {
					// ローカルZ: 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
					totalActualIkQuat = bfs.calculateSingleAxisRadRotation(
						ikLink.LocalMinAngleLimit, ikLink.LocalMaxAngleLimit,
						linkQuat, axis, angle, 2, linkBone.NormalizedLocalAxisZ,
						ikBone.Ik.UnitRotation.GetRadians().GetZ(), true)
				}
			} else if linkBone.HasFixedAxis() {
				// 軸制限ありの場合、軸にそった理想回転量とする

				// 制限角で最大変位量を制限する
				limitRotationRad := math.Min(ikBone.Ik.UnitRotation.GetRadians().GetX(), angle)
				limitQuat := mmath.NewMQuaternionFromAxisAngles(axis, limitRotationRad)
				correctIkQuat := limitQuat

				actualIkQuat := linkQuat.Muled(correctIkQuat)
				linkAxis := actualIkQuat.GetXYZ().Normalized()
				linkRad := actualIkQuat.ToRadian()
				var linkSign float64
				if linkBone.NormalizedFixedAxis.Dot(linkAxis) >= 0 {
					linkSign = 1
				} else {
					linkSign = -1
				}

				// 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
				totalActualIkQuat = mmath.NewMQuaternionFromAxisAngles(linkBone.NormalizedFixedAxis, linkRad*linkSign)
			} else {
				// 制限が無い場合、制限角の制限だけ入れる

				// 制限角で最大変位量を制限する
				limitRotationRad := math.Min(ikBone.Ik.UnitRotation.GetRadians().GetX(), angle)
				correctIkQuat := mmath.NewMQuaternionFromAxisAngles(axis, limitRotationRad)

				// 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
				totalActualIkQuat = linkQuat.Muled(correctIkQuat)
			}

			if linkBone.HasFixedAxis() {
				// 軸制限回転を求める
				totalActualIkQuat = totalActualIkQuat.ToFixedAxisRotation(linkBone.NormalizedFixedAxis)
			}

			// IKの結果を更新
			quats[0][linkIndex] = totalActualIkQuat
			rotations[0][linkIndex] = totalActualIkQuat.ToMat4()
		}
	}

	return quats[0], effectorTargetBoneNames
}

// 全ての角度をラジアン角度に分割して、そのうちのひとつの軸だけを動かす回転を取得する
// minAngleLimit: 最小軸制限
// maxAngleLimit: 最大軸制限
// linkQuat: 現在のリンクボーンの回転量
// quatAxis: 現在のIK回転の回転軸
// quatAngle: 現在のIK回転の回転角度（ラジアン）
// axisIndex: 制限軸INDEX
// axisVector: 制限軸ベクトル
// unitRadian: IKループ計算時の1回あたりの制限角度
func (bfs *BoneFrames) calculateSingleAxisRadRotation(
	minAngleLimit, maxAngleLimit *mmath.MRotation,
	linkQuat *mmath.MQuaternion,
	quatAxis *mmath.MVec3,
	quatAngle float64,
	axisIndex int,
	axisVector *mmath.MVec3,
	unitRadian float64,
	isLocal bool,
) *mmath.MQuaternion {
	// 現在調整予定角度の全ての軸の角度
	ikQuat := mmath.NewMQuaternionFromAxisAngles(quatAxis, quatAngle)

	// 全ての角度をラジアン角度に分解する
	ikRads := ikQuat.ToEulerAngles()

	axisRad := math.Abs(ikRads.Vector()[axisIndex])
	var limitRad float64
	if GIMBAL_RAD < quatAngle && quatAngle < GIMBAL2_RAD {
		limitRad = axisRad + HALF_RAD
	} else {
		limitRad = axisRad
	}

	// Calculate the maximum angle in radians
	maxRad := math.Max(
		math.Abs(minAngleLimit.GetRadians().Vector()[axisIndex]),
		math.Abs(maxAngleLimit.GetRadians().Vector()[axisIndex]),
	)

	// 最大ラジアンが制限最大角度と等しくない場合、軸の符号を逆にする
	var axisSign float64
	if maxRad != math.Abs(maxAngleLimit.GetRadians().Vector()[axisIndex]) {
		axisSign = -1
	} else {
		axisSign = 1
	}

	// 単位角で制限する
	var limitAxisRad float64
	if unitRadian != 0 {
		limitAxisRad = math.Min(unitRadian, limitRad)
	} else {
		limitAxisRad = limitRad
	}

	// 単位角で制限した角度に基づいたクォータニオン
	correctLimitIkQuat := mmath.NewMQuaternionFromAxisAngles(quatAxis, limitAxisRad)

	// 現在IKリンクに入る可能性のあるすべての角度
	totalIkQuat := linkQuat.Muled(correctLimitIkQuat)

	// 全体の角度を計算する
	totalAxisIkRad := totalIkQuat.ToRadian()
	var totalAxisIkRads *mmath.MVec3
	if isLocal {
		// ローカル軸の場合、一旦グローバル軸に直す
		totalAxisIkAxis := totalIkQuat.GetXYZ().Normalize()
		totalAxisIkRad = totalIkQuat.ToRadian()
		var totalAxisIkSign float64
		if axisVector.Dot(totalAxisIkAxis) >= 0 {
			totalAxisIkSign = 1
		} else {
			totalAxisIkSign = -1
		}

		var globalAxisVec *mmath.MVec3
		if axisIndex == 0 {
			globalAxisVec = &mmath.MVec3{1, 0, 0}
		} else if axisIndex == 1 {
			globalAxisVec = &mmath.MVec3{0, 1, 0}
		} else {
			globalAxisVec = &mmath.MVec3{0, 0, 1}
		}

		totalAxisIkQuat := mmath.NewMQuaternionFromAxisAngles(globalAxisVec, totalAxisIkRad*totalAxisIkSign)
		totalAxisIkRads = totalAxisIkQuat.ToEulerAngles().MMD()
	} else {
		// MMD上でのIKリンクの角度
		totalAxisIkRads = totalIkQuat.ToEulerAngles().MMD()
	}

	var totalAxisRad float64
	if unitRadian > quatAngle && QUARTER_RAD > totalAxisIkRad && unitRadian > totalAxisIkRad {
		// トータルが制限角度以内であれば全軸の角度を使う
		totalIkQq := linkQuat.Muled(ikQuat)
		totalAxisRad = totalIkQq.ToRadian() * axisSign
	} else if GIMBAL_RAD > quatAngle && QUARTER_RAD > totalAxisIkRad && unitRadian > totalAxisIkRad {
		// トータルが88度以内で、軸分け後が制限角度以内であれば制限角度を使う
		totalAxisRad = unitRadian * axisSign
	} else if HALF_RAD > totalAxisIkRad {
		// トータルが180度以内であれば一軸の角度を全部使う
		totalAxisRad = totalAxisIkRad * axisSign
	} else {
		// 180度を超えている場合、軸の値だけ使用する
		totalAxisRad = math.Abs(totalAxisIkRads.Vector()[axisIndex]) * axisSign
	}

	// 角度制限がある場合、全体の角度をその角度内に収める
	totalLimitAxisRad := mmath.ClampFloat(
		totalAxisRad,
		minAngleLimit.GetRadians().Vector()[axisIndex],
		maxAngleLimit.GetRadians().Vector()[axisIndex],
	)

	// 単位角とジンバルロックの整合性を取る
	var resultAxisRad float64
	if GIMBAL2_RAD < totalAxisIkRad && !isLocal {
		resultAxisRad = HALF_RAD + math.Abs(totalLimitAxisRad)
	} else if GIMBAL_RAD < totalAxisIkRad && !isLocal {
		resultAxisRad = FULL_RAD + totalLimitAxisRad
	} else {
		resultAxisRad = totalLimitAxisRad
	}

	// 指定の軸方向に回す
	resultLinkQuat := mmath.NewMQuaternionFromAxisAngles(axisVector, resultAxisRad)
	return resultLinkQuat
}

func (bfs *BoneFrames) calcBoneMatrixes(
	frames []float32,
	model *pmx.PmxModel,
	targetBoneNames []string,
	positions, rotations, scales [][]*mmath.MMat4,
	isOutLog bool,
	description string,
) BoneTrees {
	// 各ボーンの座標変換行列×逆BOf行列
	matrixes := make([][]*mmath.MMat4, 0, len(frames))
	for i := range frames {
		matrixes = append(matrixes, make([]*mmath.MMat4, 0, len(targetBoneNames)))
		for j, boneName := range targetBoneNames {
			bone := model.Bones.GetItemByName(boneName)
			matrixes[i] = append(matrixes[i], mmath.NewMMat4())
			// 逆BOf行列(初期姿勢行列)
			matrixes[i][j].Mul(bone.RevertOffsetMatrix)
			// 位置
			matrixes[i][j].Mul(positions[i][j])
			// 回転
			matrixes[i][j].Mul(rotations[i][j])
			// スケール
			matrixes[i][j].Mul(scales[i][j])
		}
	}

	boneTrees := NewBoneTrees()

	for i, frame := range frames {
		for j, boneName := range targetBoneNames {
			bone := model.Bones.GetItemByName(boneName)
			localMatrix := mmath.NewMMat4()
			for _, k := range bone.ParentBoneIndexes {
				// 親ボーンの変形行列を掛ける(親->子の順で掛ける)
				parentName := model.Bones.GetItem(k).Name
				// targetBoneNames の中にある parentName のINDEXを取得
				parentIndex := slices.Index(targetBoneNames, parentName)
				localMatrix.Mul(matrixes[i][parentIndex])
			}
			// 最後に対象ボーン自身の行列をかける
			localMatrix.Mul(matrixes[i][j])
			// BOf行列: 自身のボーンのボーンオフセット行列
			localMatrix.Mul(bone.OffsetMatrix)
			// 初期位置行列を掛けてグローバル行列を作成
			boneTrees.SetItem(bone.Name, frame, NewBoneTree(
				bone.Name,
				frame,
				localMatrix.Muled(bone.Position.ToMat4()), // グローバル行列
				localMatrix,                   // ローカル行列はそのまま
				positions[i][j].Translation(), // 移動
				rotations[i][j].Quaternion(),  // 回転
				scales[i][j].Scaling(),        // 拡大率
			))
		}
	}

	return *boneTrees
}

// アニメーション対象ボーン一覧取得
func (bfs *BoneFrames) getAnimatedBoneNames(
	model *pmx.PmxModel,
	boneNames []string,
) []string {
	if boneNames != nil || len(boneNames) > 0 {
		targetBoneNames := make([]string, 0)
		for _, boneName := range boneNames {
			if !slices.Contains(targetBoneNames, boneName) {
				targetBoneNames = append(targetBoneNames, boneName)
			}
			relativeBoneIndexes := model.Bones.GetItemByName(boneName).RelativeBoneIndexes
			for _, index := range relativeBoneIndexes {
				relativeBoneName := model.Bones.GetItem(index).Name
				if !slices.Contains(targetBoneNames, relativeBoneName) {
					targetBoneNames = append(targetBoneNames, relativeBoneName)
				}
			}
		}

		resultBoneNames := make([]string, 0)

		// ボーンINDEXでソート
		for _, bone := range model.Bones.GetSortedData() {
			if slices.Contains(targetBoneNames, bone.Name) {
				resultBoneNames = append(resultBoneNames, bone.Name)
			}
		}

		return resultBoneNames
	}

	return model.Bones.GetNames()
}

// ボーン変形行列を求める
func (bfs *BoneFrames) getBoneMatrixes(
	frames []float32,
	model *pmx.PmxModel,
	targetBoneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) ([][]*mmath.MMat4, [][]*mmath.MMat4, [][]*mmath.MMat4, [][]*mmath.MQuaternion) {
	var frameWg sync.WaitGroup
	positions := make([][]*mmath.MMat4, 0, len(frames))
	rotations := make([][]*mmath.MMat4, 0, len(frames))
	scales := make([][]*mmath.MMat4, 0, len(frames))
	quats := make([][]*mmath.MQuaternion, 0, len(frames))
	boneCount := len(targetBoneNames)

	// 最初にフレーム数*ボーン数分のスライスを確保
	for i := range frames {
		positions = append(positions, make([]*mmath.MMat4, 0, boneCount))
		rotations = append(rotations, make([]*mmath.MMat4, 0, boneCount))
		scales = append(scales, make([]*mmath.MMat4, 0, boneCount))
		quats = append(quats, make([]*mmath.MQuaternion, 0, boneCount))
		for j := 0; j < boneCount; j++ {
			positions[i] = append(positions[i], mmath.NewMMat4())
			rotations[i] = append(rotations[i], mmath.NewMMat4())
			scales[i] = append(scales[i], mmath.NewMMat4())
			quats[i] = append(quats[i], mmath.NewMQuaternion())
		}
	}

	for i, frame := range frames {
		frameWg.Add(1)
		go func(i int, frame float32) {
			defer frameWg.Done()

			var boneWg sync.WaitGroup
			// ボーンを一定件数ごとに並列処理（件数は変数保持）
			count := 100
			if boneCount < count {
				count = boneCount
			}
			for j := 0; j < boneCount; j += count {
				boneWg.Add(1)
				go func(i, j int, frame float32) {
					defer boneWg.Done()
					for k := j; k < j+count; k++ {
						if k >= boneCount {
							break
						}
						bone := model.Bones.GetItemByName(targetBoneNames[k])
						// ボーンが対象の場合、そのボーンの移動位置、回転角度、拡大率を取得
						positions[i][k] = bfs.getPosition(frame, bone.Name, model)
						rotWithEffect, rotFk := bfs.getRotation(frame, bone.Name, model, isCalcIk)
						rotations[i][k] = rotWithEffect.ToMat4()
						quats[i][k] = rotFk
						scales[i][k] = bfs.getScale(frame, bone.Name, model)
					}
				}(i, j, frame)
			}
			boneWg.Wait()
		}(i, frame)
	}
	frameWg.Wait()

	return positions, rotations, scales, quats
}

// 該当キーフレにおけるボーンの移動位置
func (bfs *BoneFrames) getPosition(frame float32, boneName string, model *pmx.PmxModel) *mmath.MMat4 {
	bone := model.Bones.GetItemByName(boneName)
	bf := bfs.Data[boneName].GetItem(frame)

	mat := mmath.NewMMat4()
	mat[0][3] = bf.Position.GetX()
	mat[1][3] = bf.Position.GetY()
	mat[2][3] = bf.Position.GetZ()

	if bone.IsEffectorTranslation() {
		// 外部親変形ありの場合、外部親変形行列を掛ける
		effectPosMat := bfs.getPositionWithEffect(frame, bone.Index, model, 0)
		mat.Mul(effectPosMat)
	}

	return mat
}

// 付与親を加味した移動位置
func (bfs *BoneFrames) getPositionWithEffect(frame float32, boneIndex int, model *pmx.PmxModel, loop int) *mmath.MMat4 {
	bone := model.Bones.GetItem(boneIndex)

	if bone.EffectFactor == 0 && loop > 20 {
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
	posMat := bfs.getPosition(frame, effectBone.Name, model)

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
	isCalcIk bool,
) (*mmath.MQuaternion, *mmath.MQuaternion) {
	bone := model.Bones.GetItemByName(boneName)

	// FK(捩り) > IK(捩り) > 付与親(捩り)
	bf := bfs.Data[boneName].GetItem(frame)
	rot := bf.Rotation.GetQuaternion().Copy()
	if isCalcIk && bf.IkRotation != nil && !bf.IkRotation.GetRadians().IsZero() {
		// IK用回転を持っている場合、置き換え
		rot = bf.IkRotation.GetQuaternion().Copy()
	} else {
		if bone.HasFixedAxis() {
			rot = rot.ToFixedAxisRotation(bone.NormalizedFixedAxis)
		}
	}

	var rotWithEffect *mmath.MQuaternion
	if bone.IsEffectorRotation() {
		// 外部親変形ありの場合、外部親変形行列を掛ける
		effectQ := rot.Muled(bfs.getRotationWithEffect(frame, bone.Index, model, isCalcIk, 0))
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
	isCalcIk bool,
	loop int,
) *mmath.MQuaternion {
	bone := model.Bones.GetItem(boneIndex)

	if bone.EffectFactor == 0 && loop > 20 {
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
	rotWithEffect, _ := bfs.getRotation(frame, effectBone.Name, model, isCalcIk)

	if bone.EffectFactor >= 0 {
		// 正の付与親
		effectQ := rotWithEffect.MulFactor(bone.EffectFactor)
		effectQ.Normalize()
		return effectQ
	} else {
		// 負の付与親の場合、逆回転
		effectQ := rotWithEffect.MulFactor(-bone.EffectFactor)
		effectQ.Invert()
		effectQ.Normalize()
		return effectQ
	}
}

// 該当キーフレにおけるボーンの拡大率
func (bfs *BoneFrames) getScale(frame float32, boneName string, model *pmx.PmxModel) *mmath.MMat4 {
	bf := bfs.Data[boneName].GetItem(frame)
	mat := mmath.NewMMat4()
	mat[0][0] += bf.Scale.GetX()
	mat[1][1] += bf.Scale.GetY()
	mat[2][2] += bf.Scale.GetZ()
	return mat
}
