package pmx

import (
	"github.com/jinzhu/copier"
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
	*core.IndexModel
	Name        string      // テクスチャ名
	TextureType TextureType // テクスチャ種別
}

func NewTexture() *Texture {
	return &Texture{
		IndexModel: &core.IndexModel{Index: -1},
		Name:       "",
	}
}

func (t *Texture) Copy() core.IIndexModel {
	copied := NewTexture()
	copier.CopyWithOption(copied, t, copier.Option{DeepCopy: true})
	return copied
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

// 共有テクスチャ辞書
type ToonTextures struct {
	*core.IndexModels[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModels: core.NewIndexModels[*Texture](10, func() *Texture { return nil }),
	}
}
