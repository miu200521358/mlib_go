package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// モーフリスト
type Morphs struct {
	*core.IndexNameModels[*Morph]
}

func NewMorphs(capacity int) *Morphs {
	return &Morphs{
		IndexNameModels: core.NewIndexNameModels[*Morph](capacity),
	}
}
