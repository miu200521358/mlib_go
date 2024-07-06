package vmd

import (
	"slices"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type BoneDelta struct {
	Bone                *pmx.Bone          // ボーン
	Frame               int                // キーフレーム
	globalMatrix        *mmath.MMat4       // グローバル行列
	localMatrix         *mmath.MMat4       // ローカル行列
	globalPosition      *mmath.MVec3       // グローバル位置
	framePosition       *mmath.MVec3       // キーフレ位置の変動量
	frameMorphPosition  *mmath.MVec3       // モーフ位置の変動量
	frameEffectPosition *mmath.MVec3       // キーフレ位置の変動量(付与親のみ)
	frameRotation       *mmath.MQuaternion // キーフレ回転の変動量
	frameMorphRotation  *mmath.MQuaternion // モーフ回転の変動量
	frameEffectRotation *mmath.MQuaternion // キーフレ回転の変動量(付与親のみ)
	frameScale          *mmath.MVec3       // キーフレスケールの変動量
	frameMorphScale     *mmath.MVec3       // モーフスケールの変動量
	unitMatrix          *mmath.MMat4
	*MorphFrameDelta
}

func (bd *BoneDelta) GlobalMatrix() *mmath.MMat4 {
	if bd.globalMatrix == nil {
		bd.globalMatrix = mmath.NewMMat4()
	}
	return bd.globalMatrix
}

func (bd *BoneDelta) LocalMatrix() *mmath.MMat4 {
	if bd.localMatrix == nil {
		// BOf行列: 自身のボーンのボーンオフセット行列をかけてローカル行列
		bd.localMatrix = bd.Bone.OffsetMatrix.Muled(bd.globalMatrix)
	}
	return bd.localMatrix
}

func (bd *BoneDelta) GlobalPosition() *mmath.MVec3 {
	if bd.globalPosition == nil {
		bd.globalPosition = bd.globalMatrix.Translation()
	}
	return bd.globalPosition
}

func (bd *BoneDelta) GlobalRotation() *mmath.MQuaternion {
	return bd.GlobalMatrix().Quaternion()
}

func (bd *BoneDelta) LocalPosition() *mmath.MVec3 {
	pos := bd.FramePosition().Copy()

	if bd.frameMorphPosition != nil && !bd.frameMorphPosition.IsZero() {
		pos.Add(bd.frameMorphPosition)
	}

	if bd.frameEffectPosition != nil && !bd.frameEffectPosition.IsZero() {
		pos.Add(bd.frameEffectPosition)
	}

	return pos
}

func (bd *BoneDelta) FramePosition() *mmath.MVec3 {
	if bd.framePosition == nil {
		bd.framePosition = mmath.NewMVec3()
	}
	return bd.framePosition
}

func (bd *BoneDelta) FrameMorphPosition() *mmath.MVec3 {
	if bd.frameMorphPosition == nil {
		bd.frameMorphPosition = mmath.NewMVec3()
	}
	return bd.frameMorphPosition
}

func (bd *BoneDelta) FrameEffectPosition() *mmath.MVec3 {
	if bd.frameEffectPosition == nil {
		bd.frameEffectPosition = mmath.NewMVec3()
	}
	return bd.frameEffectPosition
}

func (bd *BoneDelta) LocalRotation() *mmath.MQuaternion {
	rot := bd.FrameRotation().Copy()

	if bd.frameMorphRotation != nil && !bd.frameMorphRotation.IsIdent() {
		// rot = bd.frameMorphRotation.Muled(rot)
		rot.Mul(bd.frameMorphRotation)
	}

	if bd.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	if bd.frameEffectRotation != nil && !bd.frameEffectRotation.IsIdent() {
		// rot = bd.frameEffectRotation.Muled(rot)
		rot.Mul(bd.frameEffectRotation)
	}

	return rot
}

func (bd *BoneDelta) FrameRotation() *mmath.MQuaternion {
	if bd.frameRotation == nil {
		bd.frameRotation = mmath.NewMQuaternion()
	}
	return bd.frameRotation
}

func (bd *BoneDelta) FrameEffectRotation() *mmath.MQuaternion {
	if bd.frameEffectRotation == nil {
		bd.frameEffectRotation = mmath.NewMQuaternion()
	}
	return bd.frameEffectRotation
}

func (bd *BoneDelta) FrameMorphRotation() *mmath.MQuaternion {
	if bd.frameMorphRotation == nil {
		bd.frameMorphRotation = mmath.NewMQuaternion()
	}
	return bd.frameMorphRotation
}

func (bd *BoneDelta) LocalScale() *mmath.MVec3 {
	pos := bd.FrameScale().Copy()

	if bd.frameMorphScale != nil && !bd.frameMorphScale.IsZero() {
		pos.Add(bd.frameMorphScale)
	}

	return pos
}

func (bd *BoneDelta) FrameScale() *mmath.MVec3 {
	if bd.frameScale == nil {
		bd.frameScale = &mmath.MVec3{1, 1, 1}
	}
	return bd.frameScale
}

func (bd *BoneDelta) FrameMorphScale() *mmath.MVec3 {
	if bd.frameMorphScale == nil {
		bd.frameMorphScale = mmath.NewMVec3()
	}
	return bd.frameMorphScale
}

func (bd *BoneDelta) Copy() *BoneDelta {
	return &BoneDelta{
		Bone:                bd.Bone,
		Frame:               bd.Frame,
		globalMatrix:        bd.GlobalMatrix().Copy(),
		localMatrix:         bd.LocalMatrix().Copy(),
		globalPosition:      bd.GlobalPosition().Copy(),
		framePosition:       bd.FramePosition().Copy(),
		frameMorphPosition:  bd.FrameMorphPosition().Copy(),
		frameEffectPosition: bd.FrameEffectPosition().Copy(),
		frameRotation:       bd.FrameRotation().Copy(),
		frameMorphRotation:  bd.FrameMorphRotation().Copy(),
		frameEffectRotation: bd.FrameEffectRotation().Copy(),
		frameScale:          bd.FrameScale().Copy(),
		frameMorphScale:     bd.FrameMorphScale().Copy(),
		unitMatrix:          bd.unitMatrix.Copy(),
		MorphFrameDelta:     bd.MorphFrameDelta.Copy(),
	}
}

func NewBoneDelta(
	bone *pmx.Bone,
	frame int,
) *BoneDelta {
	return &BoneDelta{
		Bone:            bone,
		Frame:           frame,
		MorphFrameDelta: NewMorphFrameDelta(),
	}
}

type BoneDeltas struct {
	Data  map[int]*BoneDelta
	Names map[string]int
	mu    sync.RWMutex
}

func NewBoneDeltas() *BoneDeltas {
	return &BoneDeltas{
		Data:  make(map[int]*BoneDelta, 0),
		Names: make(map[string]int, 0),
	}
}

func (bds *BoneDeltas) Get(boneIndex int) *BoneDelta {
	bds.mu.RLock()
	defer bds.mu.RUnlock()

	if _, ok := bds.Data[boneIndex]; ok {
		return bds.Data[boneIndex]
	}
	return nil
}

func (bds *BoneDeltas) GetByName(boneName string) *BoneDelta {
	bds.mu.RLock()
	defer bds.mu.RUnlock()

	if _, ok := bds.Names[boneName]; ok {
		return bds.Data[bds.Names[boneName]]
	}
	return nil
}

func (bds *BoneDeltas) Append(boneDelta *BoneDelta) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	bds.Data[boneDelta.Bone.Index] = boneDelta
	bds.Names[boneDelta.Bone.Name] = boneDelta.Bone.Index
}

func (bds *BoneDeltas) GetBoneIndexes() []int {
	bds.mu.RLock()
	defer bds.mu.RUnlock()

	boneIndexes := make([]int, 0)
	for key := range bds.Data {
		if !slices.Contains(boneIndexes, key) {
			boneIndexes = append(boneIndexes, key)
		}
	}
	return boneIndexes
}

func (bds *BoneDeltas) Contains(boneIndex int) bool {
	bds.mu.RLock()
	defer bds.mu.RUnlock()

	_, ok := bds.Data[boneIndex]
	return ok
}

func (bds *BoneDeltas) SetGlobalMatrix(bone *pmx.Bone, globalMatrix *mmath.MMat4) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.globalMatrix = globalMatrix

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetLocalMatrix(bone *pmx.Bone, localMatrix *mmath.MMat4) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.localMatrix = localMatrix

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetUnitMatrix(bone *pmx.Bone, unitMatrix *mmath.MMat4) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.unitMatrix = unitMatrix

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFramePosition(bone *pmx.Bone, framePosition *mmath.MVec3) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.framePosition = framePosition

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameMorphPosition(bone *pmx.Bone, frameMorphPosition *mmath.MVec3) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameMorphPosition = frameMorphPosition

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameEffectPosition(bone *pmx.Bone, frameEffectPosition *mmath.MVec3) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameEffectPosition = frameEffectPosition

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameRotation(bone *pmx.Bone, frameRotation *mmath.MQuaternion) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameRotation = frameRotation

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameMorphRotation(bone *pmx.Bone, frameMorphRotation *mmath.MQuaternion) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameMorphRotation = frameMorphRotation

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameEffectRotation(bone *pmx.Bone, frameEffectRotation *mmath.MQuaternion) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameEffectRotation = frameEffectRotation

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameScale(bone *pmx.Bone, frameScale *mmath.MVec3) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameScale = frameScale

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

func (bds *BoneDeltas) SetFrameMorphScale(bone *pmx.Bone, frameMorphScale *mmath.MVec3) {
	bds.mu.Lock()
	defer bds.mu.Unlock()

	var bd *BoneDelta
	if _, ok := bds.Data[bone.Index]; ok {
		bd = bds.Data[bone.Index]
	} else {
		bd = NewBoneDelta(bone, 0)
	}

	bd.frameMorphScale = frameMorphScale

	bds.Data[bd.Bone.Index] = bd
	bds.Names[bd.Bone.Name] = bd.Bone.Index
}

// FillLocalMatrix 物理演算後にグローバル行列を埋め終わった後に呼び出して、ローカル行列を計算する
func (bds *BoneDeltas) FillLocalMatrix(physicsBoneIndexes []int) {
	for _, boneIndex := range physicsBoneIndexes {
		bd := bds.Get(boneIndex)
		if bd == nil {
			continue
		}

		var parentGlobalMatrix *mmath.MMat4
		if bd.Bone.ParentIndex >= 0 {
			parentGlobalMatrix = bds.Get(bd.Bone.ParentIndex).GlobalMatrix()
		} else {
			parentGlobalMatrix = mmath.NewMMat4()
		}
		bds.mu.Lock()
		unitMatrix := parentGlobalMatrix.Muled(bd.globalMatrix.Inverted())

		bd.localMatrix = bd.Bone.OffsetMatrix.Muled(bd.globalMatrix)
		bd.globalPosition = nil
		bd.framePosition = unitMatrix.Translation()
		bd.frameMorphPosition = nil
		bd.frameEffectPosition = nil
		bd.frameRotation = unitMatrix.Quaternion()
		bd.frameMorphRotation = nil
		bd.frameEffectRotation = nil
		bds.mu.Unlock()
	}
}

func (bds *BoneDeltas) GetNearestBones(worldPos *mmath.MVec3) []*pmx.Bone {
	bds.mu.RLock()
	defer bds.mu.RUnlock()

	bones := make([]*pmx.Bone, 0)
	distances := make([]float64, len(bds.Data))
	for i := range len(bds.Data) {
		bd := bds.Get(i)
		distances[i] = worldPos.Distance(bd.GlobalPosition())
	}
	sortedDistances := mmath.Float64Slice(distances)
	sortedIndexes := mmath.ArgSort(sortedDistances)
	for i := range sortedIndexes {
		if len(bones) > 0 {
			if !mmath.NearEquals(distances[sortedIndexes[i]], distances[sortedIndexes[0]], 0.01) {
				break
			}
		}
		bones = append(bones, bds.Data[sortedIndexes[i]].Bone)
	}
	return bones
}
