package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
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
	*mcore.IndexModel
	Name        string      // 枠名
	EnglishName string      // 枠名英
	SpecialFlag SpecialFlag // 特殊枠フラグ - 0:通常枠 1:特殊枠
	References  []Reference // 表示枠要素
	IsSystem    bool        // ツール側で追加した表示枠
}

// NewDisplaySlot
func NewDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Name:        "",
		EnglishName: "",
		SpecialFlag: SPECIAL_FLAG_OFF,
		References:  make([]Reference, 0),
		IsSystem:    false,
	}
}

// Copy
func (v *DisplaySlot) Copy() mcore.IIndexModel {
	copied := NewDisplaySlot()
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return copied
}

// 表示枠リスト
type DisplaySlots struct {
	*mcore.IndexModels[*DisplaySlot]
}

func NewDisplaySlots() *DisplaySlots {
	return &DisplaySlots{
		IndexModels: mcore.NewIndexModelCorrection[*DisplaySlot](),
	}
}
