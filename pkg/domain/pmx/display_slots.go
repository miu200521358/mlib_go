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

func (displaySlots *DisplaySlots) GetByBoneIndex(boneIndex int) *DisplaySlot {
	var result *DisplaySlot
	displaySlots.ForEach(func(index int, value *DisplaySlot) {
		for _, reference := range value.References {
			if reference.DisplayType == DISPLAY_TYPE_BONE && reference.DisplayIndex == boneIndex {
				result = value
				return
			}
		}
	})
	return result
}

func (displaySlots *DisplaySlots) GetRootDisplaySlot() (*DisplaySlot, error) {
	return displaySlots.GetByName("Root")
}

func (displaySlots *DisplaySlots) GetMorphDisplaySlot() (*DisplaySlot, error) {
	return displaySlots.GetByName("表情")
}
