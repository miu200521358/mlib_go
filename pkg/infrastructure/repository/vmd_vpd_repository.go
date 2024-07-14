package repository

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

// VMDリーダー
type VmdVpdRepository struct {
	baseRepository[*vmd.VmdMotion]
	vmdRepository *VmdRepository
	vpdRepository *VpdRepository
}

func NewVmdVpdRepository() *VmdVpdRepository {
	rep := new(VmdVpdRepository)
	rep.vmdRepository = NewVmdRepository()
	rep.vpdRepository = NewVpdRepository()
	return rep
}

// 指定されたパスのファイルからデータを読み込む
func (r *VmdVpdRepository) Load(path string) (core.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return r.vpdRepository.Load(path)
	} else {
		return r.vmdRepository.Load(path)
	}
}

func (r *VmdVpdRepository) LoadName(path string) (string, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return r.vpdRepository.LoadName(path)
	} else {
		return r.vmdRepository.LoadName(path)
	}
}
