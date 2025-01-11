package repository

import (
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type VmdRepository struct {
	*baseRepository[*vmd.VmdMotion]
}

func NewVmdRepository() *VmdRepository {
	return &VmdRepository{
		baseRepository: &baseRepository[*vmd.VmdMotion]{
			newFunc: func(path string) *vmd.VmdMotion {
				return vmd.NewVmdMotion(path)
			},
		},
	}
}
