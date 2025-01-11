package repository

import (
	"fmt"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mfile"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
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

func (rep *PmxPmxJsonRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf(mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".json" && strings.ToLower(ext) != ".pmx" {
		return false, fmt.Errorf(mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".json, .pmx"}))
	}

	return true, nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *PmxPmxJsonRepository) Load(path string) (core.IHashModel, error) {
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		return rep.pmxJsonRepository.Load(path)
	} else {
		return rep.pmxRepository.Load(path)
	}
}

func (rep *PmxPmxJsonRepository) LoadName(path string) string {
	if ok, err := rep.CanLoad(path); !ok || err != nil {
		return mi18n.T("読み込み失敗")
	}

	if strings.HasSuffix(strings.ToLower(path), ".json") {
		return rep.pmxJsonRepository.LoadName(path)
	} else {
		return rep.pmxRepository.LoadName(path)
	}
}
