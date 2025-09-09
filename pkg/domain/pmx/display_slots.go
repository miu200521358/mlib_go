package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 表示枠リスト
type DisplaySlots struct {
	*core.IndexNameModels[*DisplaySlot]
}

func NewDisplaySlots(capacity int) *DisplaySlots {
	return &DisplaySlots{
		IndexNameModels: core.NewIndexNameModels[*DisplaySlot](capacity),
	}
}

func NewInitialDisplaySlots() *DisplaySlots {
	displaySlots := &DisplaySlots{
		IndexNameModels: core.NewIndexNameModels[*DisplaySlot](0),
	}

	displaySlots.Append(NewRootDisplaySlot())
	displaySlots.Append(NewMorphDisplaySlot())

	return displaySlots
}

func (displaySlots *DisplaySlots) GetByBoneIndex(boneIndex int) *DisplaySlot {
	var result *DisplaySlot
	displaySlots.ForEach(func(index int, value *DisplaySlot) bool {
		for _, reference := range value.References {
			if reference.DisplayType == DISPLAY_TYPE_BONE && reference.DisplayIndex == boneIndex {
				result = value
				return false
			}
		}
		return true
	})
	return result
}

func (displaySlots *DisplaySlots) GetRootDisplaySlot() (*DisplaySlot, error) {
	return displaySlots.GetByName("Root")
}

func (displaySlots *DisplaySlots) GetMorphDisplaySlot() (*DisplaySlot, error) {
	return displaySlots.GetByName("表情")
}

func (displaySlots *DisplaySlots) Setup(bones *Bones) {
	bones.ForEach(func(boneIndex int, bone *Bone) bool {
		displaySlots.ForEach(func(index int, slot *DisplaySlot) bool {
			for _, reference := range slot.References {
				if reference.DisplayType == DISPLAY_TYPE_BONE && reference.DisplayIndex == boneIndex {
					// 該当ボーンの表示枠を設定
					bone.DisplaySlotIndex = slot.Index()
					return false
				}
			}
			return true
		})

		return true
	})
}
