package repository

import (
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type VmdRepository struct {
	*baseRepository[*vmd.VmdMotion]
	isLog bool
}

func NewVmdRepository(isLog bool) *VmdRepository {
	return &VmdRepository{
		baseRepository: &baseRepository[*vmd.VmdMotion]{
			newFunc: func(path string) *vmd.VmdMotion {
				return vmd.NewVmdMotion(path)
			},
		},
		isLog: isLog,
	}
}
