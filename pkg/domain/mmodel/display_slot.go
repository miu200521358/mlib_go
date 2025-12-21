package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/tiendc/go-deepcopy"
)

// DisplayType は表示枠要素タイプを表します。
type DisplayType int

const (
	DISPLAY_TYPE_BONE  DisplayType = 0 // ボーン
	DISPLAY_TYPE_MORPH DisplayType = 1 // モーフ
)

// SpecialFlag は特殊枠フラグを表します。
type SpecialFlag int

const (
	SPECIAL_FLAG_OFF SpecialFlag = 0 // 通常枠
	SPECIAL_FLAG_ON  SpecialFlag = 1 // 特殊枠（Rootと表情）
)

// Reference は表示枠要素への参照を表します。
type Reference struct {
	DisplayType  DisplayType // 要素対象（0:ボーン, 1:モーフ）
	DisplayIndex int         // ボーンIndex or モーフIndex
}

// NewReference は新しいReferenceを生成します。
func NewReference() *Reference {
	return &Reference{
		DisplayType:  DISPLAY_TYPE_BONE,
		DisplayIndex: -1,
	}
}

// NewReferenceByValues は指定値で新しいReferenceを生成します。
func NewReferenceByValues(displayType DisplayType, displayIndex int) *Reference {
	return &Reference{
		DisplayType:  displayType,
		DisplayIndex: displayIndex,
	}
}

// DisplaySlot は表示枠を表します。
type DisplaySlot struct {
	mcore.IndexNameModel
	SpecialFlag SpecialFlag  // 特殊枠フラグ（0:通常枠, 1:特殊枠）
	References  []*Reference // 表示枠要素
}

// NewDisplaySlot は新しいDisplaySlotを生成します。
func NewDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		IndexNameModel: *mcore.NewIndexNameModel(-1, "", ""),
		SpecialFlag:    SPECIAL_FLAG_OFF,
		References:     make([]*Reference, 0),
	}
}

// NewRootDisplaySlot はRoot表示枠を生成します。
func NewRootDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		IndexNameModel: *mcore.NewIndexNameModel(0, "Root", "Root"),
		SpecialFlag:    SPECIAL_FLAG_ON,
		References:     make([]*Reference, 0),
	}
}

// NewMorphDisplaySlot は表情表示枠を生成します。
func NewMorphDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		IndexNameModel: *mcore.NewIndexNameModel(1, "表情", "Exp"),
		SpecialFlag:    SPECIAL_FLAG_ON,
		References:     make([]*Reference, 0),
	}
}

// IsValid はDisplaySlotが有効かどうかを返します。
func (d *DisplaySlot) IsValid() bool {
	return d != nil && d.Index() >= 0
}

// Copy は深いコピーを作成します。
func (d *DisplaySlot) Copy() (*DisplaySlot, error) {
	cp := &DisplaySlot{}
	if err := deepcopy.Copy(cp, d); err != nil {
		return nil, err
	}
	return cp, nil
}
