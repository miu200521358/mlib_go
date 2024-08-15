package repository

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// VMDリーダー
type PmxPmxJsonRepository struct {
	pmxRepository     *PmxRepository
	pmxJsonRepository *PmxJsonRepository
}

func NewPmxPmxJsonRepository() *PmxPmxJsonRepository {
	rep := new(PmxPmxJsonRepository)
	rep.pmxRepository = NewPmxRepository()
	rep.pmxJsonRepository = NewPmxJsonRepository()
	return rep
}

func (rep *PmxPmxJsonRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	return nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *PmxPmxJsonRepository) Load(path string) (core.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		return rep.pmxJsonRepository.Load(path)
	} else {
		return rep.pmxRepository.Load(path)
	}
}

func (rep *PmxPmxJsonRepository) LoadName(path string) (string, error) {
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		return rep.pmxJsonRepository.LoadName(path)
	} else {
		return rep.pmxRepository.LoadName(path)
	}
}
