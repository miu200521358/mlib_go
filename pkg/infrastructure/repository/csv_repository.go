package repository

import (
	"encoding/csv"
	"os"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type CsvRepository struct {
	*baseRepository[*core.CsvModel]
}

func NewCsvRepository() *CsvRepository {
	return &CsvRepository{
		baseRepository: &baseRepository[*core.CsvModel]{
			newFunc: func(path string) *core.CsvModel {
				return core.NewCsvModel(make([][]string, 0))
			},
		},
	}
}

func (rep *CsvRepository) Save(path string, model core.IHashModel, includeSystem bool) error {
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
	for _, record := range model.(*core.CsvModel).Records() {
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

	return core.NewCsvModel(records), nil
}

func (rep *CsvRepository) LoadName(path string) (string, error) {
	return "", nil
}