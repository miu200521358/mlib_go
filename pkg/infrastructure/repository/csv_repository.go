package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mfile"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mcsv"
)

type CsvRepository struct {
	*baseRepository[*mcsv.CsvModel]
}

func NewCsvRepository() *CsvRepository {
	return &CsvRepository{
		baseRepository: &baseRepository[*mcsv.CsvModel]{
			newFunc: func(path string) *mcsv.CsvModel {
				return mcsv.NewCsvModel(make([][]string, 0))
			},
		},
	}
}

func (rep *CsvRepository) Save(path string, model core.IHashModel, includeSystem bool) error {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL("%s", mi18n.T("保存開始", map[string]interface{}{"Type": "Csv", "Path": path}))
	defer mlog.I("%s", mi18n.T("保存終了", map[string]interface{}{"Type": "Csv"}))

	// CSVファイルを開く
	file, err := os.Create(path)
	if err != nil {
		mlog.E("Save.Save error: %v", err)
		return err
	}
	defer file.Close()

	// CSVライターを作成
	writer := csv.NewWriter(file)

	// CSVファイルに書き込む
	for _, record := range model.(*mcsv.CsvModel).Records() {
		if err := writer.Write(record); err != nil {
			mlog.E("Save.Save error: %v", err)
			return err
		}
	}

	// ファイルに書き込む
	writer.Flush()

	return nil
}

// 指定されたパスのファイルからデータを読み込む
func (rep *CsvRepository) Load(path string) (core.IHashModel, error) {
	runtime.GOMAXPROCS(int(runtime.NumCPU()))
	defer runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	mlog.IL("%s", mi18n.T("読み込み開始", map[string]interface{}{"Type": "Csv", "Path": path}))
	defer mlog.I("%s", mi18n.T("読み込み終了", map[string]interface{}{"Type": "Csv"}))

	// CSVファイルを開く
	file, err := os.Open(path)
	if err != nil {
		mlog.E("Load.Load error: %v", err)
		return nil, err
	}
	defer file.Close()

	// CSVリーダーを作成
	reader := csv.NewReader(file)

	// CSVファイルを読み込む
	records, err := reader.ReadAll()
	if err != nil {
		mlog.E("Load.Load error: %v", err)
		return nil, err
	}

	return mcsv.NewCsvModel(records), nil
}

func (rep *CsvRepository) CanLoad(path string) (bool, error) {
	if isExist, err := mfile.ExistsFile(path); err != nil || !isExist {
		return false, fmt.Errorf("%s", mi18n.T("ファイル存在エラー", map[string]interface{}{"Path": path}))
	}

	_, _, ext := mfile.SplitPath(path)
	if strings.ToLower(ext) != ".csv" {
		return false, fmt.Errorf("%s", mi18n.T("拡張子エラー", map[string]interface{}{"Path": path, "Ext": ".csv"}))
	}

	return true, nil
}

func (rep *CsvRepository) LoadName(path string) string {
	return ""
}
