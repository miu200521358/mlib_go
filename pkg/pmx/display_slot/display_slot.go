package display_slot

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"

)

// 表示枠要素タイプ
type DisplayType int

const (
	// ボーン
	DISPLAY_TYPE_BONE DisplayType = 0
	// モーフ
	DISPLAY_TYPE_MORPH DisplayType = 1
)

type Reference struct {
	// 要素対象 0:ボーン 1:モーフ
	DisplayType DisplayType
	// ボーンIndex or モーフIndex
	DisplayIndex int
}

// 特殊枠フラグ - 0:通常枠 1:特殊枠
type SpecialFlag int

const (
	// 通常枠
	SPECIAL_FLAG_OFF SpecialFlag = 0
	// 特殊枠（Rootと表情）
	SPECIAL_FLAG_ON SpecialFlag = 1
)

type DisplaySlot struct {
	*index_model.IndexModel
	// 枠名
	Name string
	// 枠名英
	EnglishName string
	// 特殊枠フラグ - 0:通常枠 1:特殊枠
	SpecialFlag SpecialFlag
	// 表示枠要素
	References []Reference
	// ツール側で追加した表示枠
	IsSystem bool
}

// NewDisplaySlot
func NewDisplaySlot(name string, englishName string, specialFlag SpecialFlag) *DisplaySlot {
	return &DisplaySlot{
		IndexModel:  &index_model.IndexModel{Index: -1},
		Name:        name,
		EnglishName: englishName,
		SpecialFlag: specialFlag,
		References:  make([]Reference, 0),
		IsSystem:    false,
	}
}

// Copy
func (v *DisplaySlot) Copy() index_model.IndexModelInterface {
	copied := *v
	copied.References = make([]Reference, len(v.References))
	copy(copied.References, v.References)
	return &copied
}

// 表示枠リスト
type DisplaySlots struct {
	*index_model.IndexModelCorrection[*DisplaySlot]
}

func NewDisplaySlots(name string) *DisplaySlots {
	return &DisplaySlots{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*DisplaySlot](),
	}
}
