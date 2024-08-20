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
}

func NewTexture() *Texture {
	return &Texture{
		index:       -1,
		name:        "",
		englishName: "",
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
	return tex != nil && tex.index >= 0
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

func NewTextures(count int) *Textures {
	return &Textures{
		IndexModels: core.NewIndexModels[*Texture](count, func() *Texture { return nil }),
	}
}

func (textures *Textures) Copy() *Textures {
	copied := NewTextures(len(textures.Data))
	for i, texture := range textures.Data {
		copied.SetItem(i, texture.Copy().(*Texture))
	}
	return copied
}

// 共有テクスチャ辞書
type ToonTextures struct {
	*core.IndexModels[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModels: core.NewIndexModels[*Texture](10, func() *Texture { return nil }),
	}
}
