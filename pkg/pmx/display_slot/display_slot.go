package display_slot

import "github.com/miu200521358/mlib_go/pkg/core/index_name_model"

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
	index_name_model.T
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

// 表示枠リスト
type C struct {
	index_name_model.C
	Name    string
	Indexes []int
	data    map[int]T
	names   map[string]int
}

func NewDisplaySlots(name string) *C {
	return &C{
		Name:    name,
		Indexes: make([]int, 0),
		data:    make(map[int]T),
		names:   make(map[string]int),
	}
}
