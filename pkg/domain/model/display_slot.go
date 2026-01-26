// 指示: miu200521358
package model

// DisplayType は表示枠参照の種別を表す。
type DisplayType int

const (
	// DISPLAY_TYPE_BONE はボーン参照。
	DISPLAY_TYPE_BONE DisplayType = iota
	// DISPLAY_TYPE_MORPH はモーフ参照。
	DISPLAY_TYPE_MORPH
)

// Reference は表示枠の参照を表す。
type Reference struct {
	DisplayType  DisplayType
	DisplayIndex int
}

// SpecialFlag は表示枠の特殊フラグを表す。
type SpecialFlag int

const (
	// SPECIAL_FLAG_OFF は通常枠。
	SPECIAL_FLAG_OFF SpecialFlag = iota
	// SPECIAL_FLAG_ON は特殊枠。
	SPECIAL_FLAG_ON
)

// DisplaySlot は表示枠要素を表す。
type DisplaySlot struct {
	index       int
	name        string
	EnglishName string
	SpecialFlag SpecialFlag
	References  []Reference
}

// Index は表示枠 index を返す。
func (d *DisplaySlot) Index() int {
	return d.index
}

// SetIndex は表示枠 index を設定する。
func (d *DisplaySlot) SetIndex(index int) {
	d.index = index
}

// Name は表示枠名を返す。
func (d *DisplaySlot) Name() string {
	return d.name
}

// SetName は表示枠名を設定する。
func (d *DisplaySlot) SetName(name string) {
	d.name = name
}

func NewRootDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		index:       0,
		name:        "Root",
		EnglishName: "Root",
		SpecialFlag: SPECIAL_FLAG_ON,
		References:  make([]Reference, 0),
	}
}

func NewMorphDisplaySlot() *DisplaySlot {
	return &DisplaySlot{
		index:       1,
		name:        "表情",
		EnglishName: "Exp",
		SpecialFlag: SPECIAL_FLAG_ON,
		References:  make([]Reference, 0),
	}
}

// IsValid は表示枠が有効か判定する。
func (d *DisplaySlot) IsValid() bool {
	return d != nil && d.index >= 0
}
