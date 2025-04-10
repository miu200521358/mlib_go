package pmx

import (
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

func NewDisplaySlotReferenceByValues(displayType DisplayType, displayIndex int) *Reference {
	return &Reference{
		DisplayType:  displayType,
		DisplayIndex: displayIndex,
	}
}

// 特殊枠フラグ - 0:通常枠 1:特殊枠
type SpecialFlag int

const (
	SPECIAL_FLAG_OFF SpecialFlag = 0 // 通常枠
	SPECIAL_FLAG_ON  SpecialFlag = 1 // 特殊枠（Rootと表情）
)

type DisplaySlot struct {
	index       int          // 表示枠INDEX
	name        string       // 表示枠名
	englishName string       // 表示枠英名
	SpecialFlag SpecialFlag  // 特殊枠フラグ - 0:通常枠 1:特殊枠
	References  []*Reference // 表示枠要素
}

// NewDisplaySlot
func NewDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		index:       -1,
		name:        "",
		englishName: "",
		SpecialFlag: SPECIAL_FLAG_OFF,
		References:  make([]*Reference, 0),
	}
}

func NewRootDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		index:       0,
		name:        "Root",
		englishName: "Root",
		SpecialFlag: SPECIAL_FLAG_ON,
		References:  make([]*Reference, 0),
	}
}

func NewMorphDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		index:       1,
		name:        "表情",
		englishName: "Exp",
		SpecialFlag: SPECIAL_FLAG_ON,
		References:  make([]*Reference, 0),
	}
}

func (displaySlot *DisplaySlot) Index() int {
	return displaySlot.index
}

func (displaySlot *DisplaySlot) SetIndex(index int) {
	displaySlot.index = index
}

func (displaySlot *DisplaySlot) Name() string {
	return displaySlot.name
}

func (displaySlot *DisplaySlot) SetName(name string) {
	displaySlot.name = name
}

func (displaySlot *DisplaySlot) EnglishName() string {
	return displaySlot.englishName
}

func (displaySlot *DisplaySlot) SetEnglishName(englishName string) {
	displaySlot.englishName = englishName
}

func (displaySlot *DisplaySlot) IsValid() bool {
	return displaySlot != nil && displaySlot.index >= 0
}

// Copy
func (displaySlot *DisplaySlot) Copy() core.IIndexNameModel {
	copiedReference := make([]*Reference, len(displaySlot.References))
	copy(copiedReference, displaySlot.References)

	return &DisplaySlot{
		index:       displaySlot.index,
		name:        displaySlot.name,
		englishName: displaySlot.englishName,
		SpecialFlag: displaySlot.SpecialFlag,
		References:  copiedReference,
	}
}
