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

func (rep *VmdVpdRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	return nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *VmdVpdRepository) Load(path string) (core.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return rep.vpdRepository.Load(path)
	} else {
		return rep.vmdRepository.Load(path)
	}
}

func (rep *VmdVpdRepository) LoadName(path string) (string, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return rep.vpdRepository.LoadName(path)
	} else {
		return rep.vmdRepository.LoadName(path)
	}
}
