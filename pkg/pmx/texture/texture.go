package texture

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_name_model"

)

type TextureType int

const (
	TEXTURE TextureType = 0
	TOON    TextureType = 1
	SPHERE  TextureType = 2
)

type T struct {
	index_name_model.T
	Index       int
	Name        string
	TextureType TextureType
	Path        string
	Valid       bool
}

func NewTexture(index int, name string, textureType TextureType) *T {
	return &T{
		Index:       index,
		Name:        name,
		TextureType: textureType,
		Path:        "",
		Valid:       false,
	}
}

// Copy
func (v *T) Copy() *T {
	copied := *v
	return &copied
}

// テクスチャリスト
type C struct {
	index_name_model.C
	Name    string
	Indexes []int
	data    map[int]T
	names   map[string]int
}

func NewTextures(name string) *C {
	return &C{
		Name:    name,
		Indexes: make([]int, 0),
		data:    make(map[int]T),
		names:   make(map[string]int),
	}
}

// 共有テクスチャ辞書
type ToonC struct {
	index_name_model.C
	Name    string
	Indexes []int
	data    map[int]T
	names   map[string]int
}

func NewToonTextures(name string) *ToonC {
	return &ToonC{
		Name:    name,
		Indexes: make([]int, 0),
		data:    make(map[int]T),
		names:   make(map[string]int),
	}
}
