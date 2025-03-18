package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// テクスチャ種別
type TextureType int

const (
	TEXTURE_TYPE_TEXTURE TextureType = 0 // テクスチャ
	TEXTURE_TYPE_TOON    TextureType = 1 // Toonテクスチャ
	TEXTURE_TYPE_SPHERE  TextureType = 2 // スフィアテクスチャ
)

type Texture struct {
	index       int         // テクスチャINDEX
	name        string      // テクスチャ名
	englishName string      // テクスチャ英名
	TextureType TextureType // テクスチャ種別
	valid       bool        // 有効フラグ
}

func NewTexture() *Texture {
	return &Texture{
		index:       -1,
		name:        "",
		englishName: "",
		valid:       false,
	}
}

func (tex *Texture) Index() int {
	return tex.index
}

func (tex *Texture) SetIndex(index int) {
	tex.index = index
}

func (tex *Texture) Name() string {
	return tex.name
}

func (tex *Texture) SetName(name string) {
	tex.name = name
}

func (tex *Texture) EnglishName() string {
	return tex.englishName
}

func (tex *Texture) SetEnglishName(englishName string) {
	tex.englishName = englishName
}

func (tex *Texture) IsValid() bool {
	return tex.valid
}

func (tex *Texture) SetValid(valid bool) {
	tex.valid = valid
}

func (tex *Texture) Copy() core.IIndexModel {
	return &Texture{
		index:       tex.index,
		name:        tex.name,
		englishName: tex.englishName,
		TextureType: tex.TextureType,
	}
}

// テクスチャリスト
type Textures struct {
	*core.IndexModels[*Texture]
}

func NewTextures(capacity int) *Textures {
	return &Textures{
		IndexModels: core.NewIndexModels[*Texture](capacity),
	}
}
