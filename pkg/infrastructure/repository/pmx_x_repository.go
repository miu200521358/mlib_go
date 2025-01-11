package repository

import (
	"fmt"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mfile"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

// VMDリーダー
type PmxXRepository struct {
	pmxRepository *PmxRepository
	xRepository   *XRepository
}

func NewPmxXRepository() *PmxXRepository {
	rep := new(PmxXRepository)
	rep.pmxRepository = NewPmxRepository()
	rep.xRepository = NewXRepository()
	return rep
}

func (rep *PmxXRepository) Save(overridePath string, data core.IHashModel, includeSystem bool) error {
	return nil
}

func (rep *PmxXRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf("%s", mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".x" && strings.ToLower(ext) != ".pmx" {
		return false, fmt.Errorf("%s", mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".x, .pmx"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *PmxXRepository) Load(path string) (core.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".x") {
		return rep.xRepository.Load(path)
	} else {
		return rep.pmxRepository.Load(path)
	}
}

func (rep *PmxXRepository) LoadName(path string) string {
	if ok, err := rep.CanLoad(path); !ok || err != nil {
		return mi18n.T("読み込み失敗")
	}

	if strings.HasSuffix(strings.ToLower(path), ".x") {
		return rep.xRepository.LoadName(path)
	} else {
		return rep.pmxRepository.LoadName(path)
	}
}
