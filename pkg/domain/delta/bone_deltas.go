package delta

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

// BoneDeltas は BoneDelta の集合を管理する
type BoneDeltas struct {
	data  []*BoneDelta
	bones *pmx.Bones
}

// NewBoneDeltas はボーン数に合わせて BoneDelta のスライスを用意する
func NewBoneDeltas(bones *pmx.Bones) *BoneDeltas {
	return &BoneDeltas{
		data:  make([]*BoneDelta, bones.Length()),
		bones: bones,
	}
}

func (bds *BoneDeltas) Length() int {
	return len(bds.data)
}

// Get は boneIndex に対応する BoneDelta を返す
func (bds *BoneDeltas) Get(boneIndex int) *BoneDelta {
	if boneIndex < 0 || boneIndex >= len(bds.data) {
		return nil
	}
	return bds.data[boneIndex]
}

// GetByName はボーン名に対応する BoneDelta を返す
func (bds *BoneDeltas) GetByName(boneName string) *BoneDelta {
	if bone, err := bds.bones.GetByName(boneName); err == nil {
		return bds.Get(bone.Index())
	}

	return nil
}

// Update は BoneDelta をデータにセットする
func (bds *BoneDeltas) Update(bd *BoneDelta) {
	if bd == nil || bd.Bone == nil {
		return
	}
	idx := bd.Bone.Index()
	if idx >= 0 && idx < len(bds.data) {
		bds.data[idx] = bd
	}
}

// Contains は指定のインデックスに BoneDelta が存在するかを返す
func (bds *BoneDeltas) Contains(boneIndex int) bool {
	return boneIndex >= 0 && boneIndex < len(bds.data) && bds.data[boneIndex] != nil
}
func (bds *BoneDeltas) Iterator() <-chan struct {
	Index int
	Delta *BoneDelta
} {
	ch := make(chan struct {
		Index int
		Delta *BoneDelta
	})
	go func() {
		for i, bd := range bds.data {
			if bd != nil {
				ch <- struct {
					Index int
					Delta *BoneDelta
				}{Index: i, Delta: bd}
			}
		}
		close(ch)
	}()
	return ch
}
