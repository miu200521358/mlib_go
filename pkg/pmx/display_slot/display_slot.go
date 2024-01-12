package display_slot

// 表示枠要素タイプ
type DisplayType int

const (
	// ボーン
	BONE DisplayType = 0
	// モーフ
	MORPH DisplayType = 1
)

type Reference struct {
	// 要素対象 0:ボーン 1:モーフ
	DisplayType DisplayType
	// ボーンIndex or モーフIndex
	DisplayIndex int
}

type SpecialFlag int

const (
	NORMAL  SpecialFlag = 0
	SPECIAL SpecialFlag = 1
)

type T struct {
	// 枠名
	Index int
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
func NewDisplaySlot(index int, name string, englishName string, specialFlag SpecialFlag) *T {
	return &T{
		Index:       index,
		Name:        name,
		EnglishName: englishName,
		SpecialFlag: specialFlag,
		References:  make([]Reference, 0),
		IsSystem:    false,
	}
}

// Copy
func (v *T) Copy() *T {
	copied := *v
	copied.References = make([]Reference, len(v.References))
	copy(copied.References, v.References)
	return &copied
}
