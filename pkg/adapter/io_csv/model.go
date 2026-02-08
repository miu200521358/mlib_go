// 指示: miu200521358
package io_csv

import "github.com/miu200521358/mlib_go/pkg/shared/hashable"

// CsvModel はCSV行列を保持するハッシュモデルを表す。
type CsvModel struct {
	*hashable.HashableBase
	records [][]string
}

// NewCsvModel はCsvModelを生成する。
func NewCsvModel(records [][]string) *CsvModel {
	base := hashable.NewHashableBase("", "")
	model := &CsvModel{
		HashableBase: base,
	}
	model.SetHashPartsFunc(model.GetHashParts)
	model.SetRecords(records)
	return model
}

// Records はCSVの行列データを返す。
func (m *CsvModel) Records() [][]string {
	if m == nil {
		return nil
	}
	return deepCopyCsvRecords(m.records)
}

// SetRecords はCSVの行列データを設定する。
func (m *CsvModel) SetRecords(records [][]string) {
	if m == nil {
		return
	}
	m.records = deepCopyCsvRecords(records)
}

// GetHashParts はCSVモデル追加ハッシュ要素を返す。
func (m *CsvModel) GetHashParts() string {
	return ""
}

// deepCopyCsvRecords はCSV行列を深いコピーで複製する。
func deepCopyCsvRecords(records [][]string) [][]string {
	if records == nil {
		return nil
	}
	copied := make([][]string, len(records))
	for i := range records {
		if records[i] == nil {
			continue
		}
		row := make([]string, len(records[i]))
		copy(row, records[i])
		copied[i] = row
	}
	return copied
}
