package pmx

import (
	"github.com/jinzhu/copier"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 表示枠要素タイプ
type DisplayType int

const (
	DISPLAY_TYPE_BONE  DisplayType = 0 // ボーン
	DISPLAY_TYPE_MORPH DisplayType = 1 // モーフ
)

type Reference struct {
	DisplayType  DisplayType // 要素対象 0:ボーン 1:モーフ
	DisplayIndex int         // ボーンIndex or モーフIndex
}

func NewDisplaySlotReference() *Reference {
	return &Reference{
		DisplayType:  0,
		DisplayIndex: -1,
	}
}

// 特殊枠フラグ - 0:通常枠 1:特殊枠
type SpecialFlag int

const (
	SPECIAL_FLAG_OFF SpecialFlag = 0 // 通常枠
	SPECIAL_FLAG_ON  SpecialFlag = 1 // 特殊枠（Rootと表情）
)

type DisplaySlot struct {
	*core.IndexNameModel
	SpecialFlag SpecialFlag // 特殊枠フラグ - 0:通常枠 1:特殊枠
	References  []Reference // 表示枠要素
}

// NewDisplaySlot
func NewDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		IndexNameModel: core.NewIndexNameModel(-1, "", ""),
		SpecialFlag:    SPECIAL_FLAG_OFF,
		References:     make([]Reference, 0),
	}
}

// Copy
func (v *DisplaySlot) Copy() core.IIndexModel {
	copied := NewDisplaySlot()
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return copied
}

// 表示枠リスト
type DisplaySlots struct {
	*core.IndexModels[*DisplaySlot]
}

func NewDisplaySlots(count int) *DisplaySlots {
	return &DisplaySlots{
		IndexModels: core.NewIndexModels[*DisplaySlot](count, func() *DisplaySlot { return nil }),
	}
}
