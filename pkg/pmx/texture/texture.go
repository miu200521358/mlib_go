package texture

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
)

// テクスチャ種別
type TextureType int

const (
	// テクスチャ
	TEXTURE_TYPE_TEXTURE TextureType = 0
	// Toonテクスチャ
	TEXTURE_TYPE_TOON TextureType = 1
	// スフィアテクスチャ
	TEXTURE_TYPE_SPHERE TextureType = 2
)

type Texture struct {
	*index_model.IndexModel
	// テクスチャ名
	Name string
	// テクスチャ種別
	TextureType TextureType
	// テクスチャフルパス
	Path string
	// テクスチャフルパスが有効であるか否か
	Valid bool
}

func NewTexture() *Texture {
	return &Texture{
		IndexModel:  &index_model.IndexModel{Index: -1},
		Name:        "",
		TextureType: TEXTURE_TYPE_TEXTURE,
		Path:        "",
		Valid:       false,
	}
}

// テクスチャリスト
type Textures struct {
	*index_model.IndexModelCorrection[*Texture]
}

func NewTextures() *Textures {
	return &Textures{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Texture](),
	}
}

// 共有テクスチャ辞書
type ToonTextures struct {
	*index_model.IndexModelCorrection[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Texture](),
	}
}
