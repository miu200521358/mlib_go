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

			effectorLocalPosition.Normalize()
			ikLocalPosition.Normalize()

			// ベクトル (1) を (2) に一致させるための最短回転量（Axis-Angle）
			// 回転軸
			linkAxis := effectorLocalPosition.Cross(ikLocalPosition).Normalize()
			// 回転角(ラジアン)
			linkAngle := math.Acos(mmath.ClampFloat(effectorLocalPosition.Dot(ikLocalPosition), -1, 1))

			// 単位角を超えないようにする
			linkAngle = mmath.ClampFloat(linkAngle, -unitRad, unitRad)

			// リンクボーンの角度を取得
			linkQuat := quats[linkIndex]

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
						linkQuat, linkAxis, linkAngle, 0, axisVector)
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
						linkQuat, linkAxis, linkAngle, 1, axisVector)
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
						linkQuat, linkAxis, linkAngle, 2, axisVector)
				}
			} else {
				if linkBone.HasFixedAxis() {
					// 軸制限ありの場合、軸にそった理想回転量とする
					linkAxis = linkBone.NormalizedFixedAxis
					if linkBone.NormalizedFixedAxis.Dot(linkAxis) < 0 {
						linkAngle *= -1
					}
				}

				correctIkQuat := mmath.NewMQuaternionFromAxisAngles(linkAxis, linkAngle)

				// 既存のFK回転・IK回転・今回の計算をすべて含めて実際回転を求める
				totalActualIkQuat = linkQuat.Muled(correctIkQuat)
			}

			if linkBone.HasFixedAxis() {
				// 軸制限回転を求める
				totalActualIkQuat = totalActualIkQuat.ToFixedAxisRotation(linkBone.NormalizedFixedAxis)
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
	minAngleLimit float64,
	maxAngleLimit float64,
	linkQuat *mmath.MQuaternion,
	quatAxis *mmath.MVec3,
	quatAngle float64,
	axisIndex int,
	axisVector *mmath.MVec3,
) *mmath.MQuaternion {
	quat := mmath.NewMQuaternionFromAxisAngles(quatAxis, quatAngle)

	// 現在IKリンクに入る可能性のあるすべての角度
	totalIkQuat := linkQuat.Muled(quat)

	totalIkRad := totalIkQuat.ToRadian()
	if quatAxis.Dot(axisVector) < 0 {
		totalIkRad *= -1
	}

	fSX := math.Sin(totalIkRad) // sin(θ)
	fX := math.Asin(fSX)        // 一軸回り決定

	// ジンバルロック回避
	totalIkRads, isGimbal := totalIkQuat.ToRadiansWithGimbal(axisIndex)
	if isGimbal || math.Abs(totalIkRad) > math.Pi {
		fX = totalIkRads.Vector()[axisIndex]
		if fX < 0 {
			fX = -(math.Pi - fX)
		} else {
			fX = math.Pi - fX
		}
	}

	// 角度の制限
	if fX < minAngleLimit {
		tf := 2*minAngleLimit - fX

		fX = mmath.ClampFloat(tf, minAngleLimit, maxAngleLimit)
	}
	if fX > maxAngleLimit {
		tf := 2*maxAngleLimit - fX

		fX = mmath.ClampFloat(tf, minAngleLimit, maxAngleLimit)
	}

	return mmath.NewMQuaternionFromAxisAngles(axisVector, fX)
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
