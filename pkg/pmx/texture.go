package pmx

import (
	"image"

	"github.com/jinzhu/copier"
	"github.com/miu200521358/mlib_go/pkg/mcore"
)

// テクスチャ種別
type TextureType int

const (
	TEXTURE_TYPE_TEXTURE TextureType = 0 // テクスチャ
	TEXTURE_TYPE_TOON    TextureType = 1 // Toonテクスチャ
	TEXTURE_TYPE_SPHERE  TextureType = 2 // スフィアテクスチャ
)

type Texture struct {
	*mcore.IndexModel
	Name              string       // テクスチャ名
	TextureType       TextureType  // テクスチャ種別
	Path              string       // テクスチャフルパス
	Valid             bool         // テクスチャフルパスが有効であるか否か
	glId              uint32       // OpenGLテクスチャID
	Initialized       bool         // 描画初期化済みフラグ
	Image             *image.NRGBA // テクスチャイメージ
	textureUnitId     uint32       // テクスチャ種類別描画先ユニットID
	textureUnitNo     uint32       // テクスチャ種類別描画先ユニット番号
	IsGeneratedMipmap bool         // ミップマップが生成されているか否か
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

func (t *Texture) Copy() mcore.IIndexModel {
	copied := NewTexture()
	copier.CopyWithOption(copied, t, copier.Option{DeepCopy: true})
	return copied
}

// テクスチャリスト
type Textures struct {
	*mcore.IndexModels[*Texture]
}

func NewTextures(count int) *Textures {
	return &Textures{
		IndexModels: mcore.NewIndexModels[*Texture](count, func() *Texture { return nil }),
	}
}

// 共有テクスチャ辞書
type ToonTextures struct {
	*mcore.IndexModels[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModels: mcore.NewIndexModels[*Texture](10, func() *Texture { return nil }),
	}
}
