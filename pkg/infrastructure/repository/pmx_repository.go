package repository

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type PmxRepository struct {
	*baseRepository[*pmx.PmxModel]
}

func NewPmxRepository() *PmxRepository {
	return &PmxRepository{
		baseRepository: &baseRepository[*pmx.PmxModel]{
			newFunc: func(path string) *pmx.PmxModel {
				return &pmx.PmxModel{
					HashModel: core.NewHashModel(path),
				}
			},
		},
	}
}
