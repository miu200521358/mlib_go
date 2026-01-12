// 指示: miu200521358
package delta

import "github.com/miu200521358/mlib_go/pkg/domain/model"

// BoneDeltas はボーン差分の集合を表す。
type BoneDeltas struct {
	data  []*BoneDelta
	bones *model.BoneCollection
}

// NewBoneDeltas はボーンコレクションに合わせて生成する。
func NewBoneDeltas(bones *model.BoneCollection) *BoneDeltas {
	length := 0
	if bones != nil {
		length = bones.Len()
	}
	return &BoneDeltas{
		data:  make([]*BoneDelta, length),
		bones: bones,
	}
}

// Len は要素数を返す。
func (b *BoneDeltas) Len() int {
	if b == nil {
		return 0
	}
	return len(b.data)
}

// Get はindexのBoneDeltaを返す。
func (b *BoneDeltas) Get(index int) *BoneDelta {
	if b == nil || index < 0 || index >= len(b.data) {
		return nil
	}
	return b.data[index]
}

// GetByName は名前に対応するBoneDeltaを返す。
func (b *BoneDeltas) GetByName(name string) *BoneDelta {
	if b == nil || b.bones == nil {
		return nil
	}
	bone, err := b.bones.GetByName(name)
	if err != nil {
		return nil
	}
	return b.Get(bone.Index())
}

// Update はBoneDeltaを更新する。
func (b *BoneDeltas) Update(delta *BoneDelta) {
	if b == nil || delta == nil || delta.Bone == nil {
		return
	}
	idx := delta.Bone.Index()
	if idx < 0 || idx >= len(b.data) {
		return
	}
	b.data[idx] = delta
}

// Delete はindexのBoneDeltaを削除する。
func (b *BoneDeltas) Delete(index int) {
	if b == nil || index < 0 || index >= len(b.data) {
		return
	}
	b.data[index] = nil
}

// Contains はindexの差分が存在するか判定する。
func (b *BoneDeltas) Contains(index int) bool {
	if b == nil || index < 0 || index >= len(b.data) {
		return false
	}
	return b.data[index] != nil
}

// ContainsByName は名前に対応する差分が存在するか判定する。
func (b *BoneDeltas) ContainsByName(name string) bool {
	if b == nil || b.bones == nil {
		return false
	}
	bone, err := b.bones.GetByName(name)
	if err != nil {
		return false
	}
	return b.Contains(bone.Index())
}

// ForEach は全要素を走査する。
func (b *BoneDeltas) ForEach(fn func(index int, delta *BoneDelta) bool) {
	if b == nil || fn == nil {
		return
	}
	for i, v := range b.data {
		if !fn(i, v) {
			return
		}
	}
}
