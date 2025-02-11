package repository

import (
	"fmt"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
)

// VMDリーダー
type VmdVpdRepository struct {
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

func (rep *VmdVpdRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf("%s", mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".vmd" && strings.ToLower(ext) != ".vpd" {
		return false, fmt.Errorf("%s", mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".vmd, .vpd"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *VmdVpdRepository) Load(path string) (core.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return rep.vpdRepository.Load(path)
	} else {
		return rep.vmdRepository.Load(path)
	}
}

func (rep *VmdVpdRepository) LoadName(path string) string {
	if strings.HasSuffix(strings.ToLower(path), ".vpd") {
		return rep.vpdRepository.LoadName(path)
	} else {
		return rep.vmdRepository.LoadName(path)
	}
}
