package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// 共有テクスチャ辞書
type ToonTextures struct {
	*core.IndexModels[*Texture]
}

func NewToonTextures() *ToonTextures {
	return &ToonTextures{
		IndexModels: core.NewIndexModels[*Texture](10),
	}
}
