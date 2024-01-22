package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
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
	*mcore.IndexModel
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
		IndexModel:  &mcore.IndexModel{Index: -1},
		Name:        "",
		TextureType: TEXTURE_TYPE_TEXTURE,
		Path:        "",
		Valid:       false,
	}
}

// テクスチャリスト
type Textures struct {
	*mcore.IndexModelCorrection[*Texture]
}

func NewTextures() *Textures {
	return &Textures{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Texture](),
	}
}

// 共有テクスチャ辞書
type ToonTextures struct {
	*mcore.IndexModelCorrection[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Texture](),
	}
}
