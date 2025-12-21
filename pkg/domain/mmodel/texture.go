package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
)

// Texture はテクスチャを表します。
type Texture struct {
	mcore.IndexModel        // インデックス
	Path             string // テクスチャファイルパス（相対パス）
}

// NewTexture は新しいTextureを生成します。
func NewTexture() *Texture {
	return &Texture{
		IndexModel: *mcore.NewIndexModel(-1),
		Path:       "",
	}
}

// NewTextureByPath はパスを指定してTextureを生成します。
func NewTextureByPath(path string) *Texture {
	return &Texture{
		IndexModel: *mcore.NewIndexModel(-1),
		Path:       path,
	}
}

// IsValid はTextureが有効かどうかを返します。
func (t *Texture) IsValid() bool {
	return t != nil && t.IndexModel.IsValid()
}

// Copy は深いコピーを作成します。
func (t *Texture) Copy() (*Texture, error) {
	return &Texture{
		IndexModel: *mcore.NewIndexModel(t.Index()),
		Path:       t.Path,
	}, nil
}
