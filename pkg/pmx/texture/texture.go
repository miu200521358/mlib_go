package texture

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"

)

type TextureType int

const (
	TEXTURE TextureType = 0
	TOON    TextureType = 1
	SPHERE  TextureType = 2
)

type Texture struct {
	index_model.IndexModel
	Index       int
	Name        string
	TextureType TextureType
	Path        string
	Valid       bool
}

func NewTexture(index int, name string, textureType TextureType) *Texture {
	return &Texture{
		Index:       index,
		Name:        name,
		TextureType: textureType,
		Path:        "",
		Valid:       false,
	}
}

// テクスチャリスト
type Textures struct {
	index_model.IndexModelCorrection[*Texture]
}

func NewTextures(name string) *Textures {
	return &Textures{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Texture](),
	}
}

// 共有テクスチャ辞書
type ToonTextures struct {
	index_model.IndexModelCorrection[*Texture]
}

func NewToonTextures(name string) *ToonTextures {
	return &ToonTextures{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Texture](),
	}
}
