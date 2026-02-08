// 指示: miu200521358
package io_csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
)

const csvExtension = ".csv"

// CsvRepository はCSV入出力の具象実装を表す。
type CsvRepository struct {
	mu      sync.RWMutex
	profile *CsvProfile
}

// NewCsvRepository はCsvRepositoryを生成する。
func NewCsvRepository() *CsvRepository {
	return &CsvRepository{}
}

// NewCsvRepositoryWithProfile は検証プロファイル付きCsvRepositoryを生成する。
func NewCsvRepositoryWithProfile(profile CsvProfile) *CsvRepository {
	return &CsvRepository{
		profile: cloneCsvProfilePtr(&profile),
	}
}

// SetProfile はCSV検証プロファイルを設定する。
func (r *CsvRepository) SetProfile(profile CsvProfile) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profile = cloneCsvProfilePtr(&profile)
}

// ClearProfile はCSV検証プロファイルを解除する。
func (r *CsvRepository) ClearProfile() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profile = nil
}

// CanLoad は読み込み可能な拡張子か判定する。
func (r *CsvRepository) CanLoad(path string) bool {
	return hasCsvExtension(path)
}

// InferName はパスから表示名を推定する。
func (r *CsvRepository) InferName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

// Load はCSVを読み込み、必要ならプロファイル検証を行う。
func (r *CsvRepository) Load(path string) (hashable.IHashable, error) {
	if !hasCsvExtension(path) {
		return nil, io_common.NewIoExtInvalid(path, nil)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, io_common.NewIoFileNotFound(path, err)
		}
		return nil, io_common.NewIoFileNotFound(path, err)
	}

	records, err := readCsvRecords(path)
	if err != nil {
		return nil, io_common.NewIoParseFailed("CSV解析に失敗しました: %s", err, filepath.Base(path))
	}

	profile := r.profileSnapshot()
	if err := validateCsvRecords(records, profile); err != nil {
		return nil, io_common.NewIoParseFailed("CSV検証に失敗しました: %s", err, filepath.Base(path))
	}

	model := NewCsvModel(records)
	model.SetPath(path)
	model.SetName(r.InferName(path))
	model.SetFileModTime(fileInfo.ModTime().UnixNano())
	model.UpdateHash()
	return model, nil
}

// Save はCSVモデルを指定パスへ保存する。
func (r *CsvRepository) Save(path string, data hashable.IHashable, opts io_common.SaveOptions) error {
	_ = opts
	if !hasCsvExtension(path) {
		return io_common.NewIoEncodeFailed("CSV拡張子ではないため保存できません: %s", nil, filepath.Base(path))
	}

	model, ok := data.(*CsvModel)
	if !ok || model == nil {
		return io_common.NewIoEncodeFailed("CSVモデル型ではないため保存できません", nil)
	}

	file, err := os.Create(path)
	if err != nil {
		return io_common.NewIoSaveFailed("CSVファイル保存に失敗しました: %s", err, filepath.Base(path))
	}
	defer func() {
		_ = file.Close()
	}()

	writer := csv.NewWriter(file)
	for rowIndex, row := range model.Records() {
		if err := writer.Write(row); err != nil {
			return io_common.NewIoEncodeFailed("CSV書き込みに失敗しました(行:%d)", err, rowIndex+1)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return io_common.NewIoEncodeFailed("CSV書き込み確定に失敗しました", err)
	}

	return nil
}

// profileSnapshot は検証プロファイルのスナップショットを返す。
func (r *CsvRepository) profileSnapshot() *CsvProfile {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return cloneCsvProfilePtr(r.profile)
}

// hasCsvExtension はCSV拡張子か判定する。
func hasCsvExtension(path string) bool {
	return strings.EqualFold(filepath.Ext(path), csvExtension)
}

// readCsvRecords はCSVレコードを読み込み、空行を除去して返す。
func readCsvRecords(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return normalizeCsvRecords(records), nil
}

// normalizeCsvRecords は空行を除去したCSVレコードを返す。
func normalizeCsvRecords(records [][]string) [][]string {
	if len(records) == 0 {
		return [][]string{}
	}
	normalized := make([][]string, 0, len(records))
	for _, row := range records {
		if isEmptyCsvRow(row) {
			continue
		}
		copiedRow := make([]string, len(row))
		copy(copiedRow, row)
		normalized = append(normalized, copiedRow)
	}
	return normalized
}

// isEmptyCsvRow は空白のみの行か判定する。
func isEmptyCsvRow(row []string) bool {
	if len(row) == 0 {
		return true
	}
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// validateCsvRecords はプロファイル条件でCSVレコードを検証する。
func validateCsvRecords(records [][]string, profile *CsvProfile) error {
	if profile == nil {
		return nil
	}

	startIndex := 0
	if profile.HasHeader {
		if len(records) == 0 {
			return fmt.Errorf("ヘッダ行が存在しません")
		}
		if err := validateCsvColumns(records[0], profile, 1); err != nil {
			return err
		}
		if len(profile.Header) > 0 {
			if err := validateCsvHeader(records[0], profile.Header, profile.AllowExtraColumns); err != nil {
				return err
			}
		}
		startIndex = 1
	}

	for i := startIndex; i < len(records); i++ {
		if err := validateCsvColumns(records[i], profile, i+1); err != nil {
			return err
		}
	}

	return nil
}

// validateCsvColumns は列数制約を検証する。
func validateCsvColumns(row []string, profile *CsvProfile, rowIndex int) error {
	columns := len(row)
	if profile.ExactColumns > 0 {
		if columns < profile.ExactColumns {
			return fmt.Errorf("列数不足です(行:%d 期待:%d 実際:%d)", rowIndex, profile.ExactColumns, columns)
		}
		if !profile.AllowExtraColumns && columns != profile.ExactColumns {
			return fmt.Errorf("列数不一致です(行:%d 期待:%d 実際:%d)", rowIndex, profile.ExactColumns, columns)
		}
		return nil
	}

	if profile.MinColumns > 0 && columns < profile.MinColumns {
		return fmt.Errorf("最小列数を満たしていません(行:%d 期待:%d 実際:%d)", rowIndex, profile.MinColumns, columns)
	}

	if len(profile.Header) > 0 {
		expected := len(profile.Header)
		if columns < expected {
			return fmt.Errorf("ヘッダ定義より列数が不足しています(行:%d 期待:%d 実際:%d)", rowIndex, expected, columns)
		}
		if !profile.AllowExtraColumns && columns != expected {
			return fmt.Errorf("ヘッダ定義と列数が一致しません(行:%d 期待:%d 実際:%d)", rowIndex, expected, columns)
		}
	}

	return nil
}

// validateCsvHeader はヘッダ行の列名一致を検証する。
func validateCsvHeader(header []string, expected []string, allowExtraColumns bool) error {
	if len(header) < len(expected) {
		return fmt.Errorf("ヘッダ列数が不足しています(期待:%d 実際:%d)", len(expected), len(header))
	}
	if !allowExtraColumns && len(header) != len(expected) {
		return fmt.Errorf("ヘッダ列数が一致しません(期待:%d 実際:%d)", len(expected), len(header))
	}
	for i := range expected {
		if strings.TrimSpace(header[i]) != strings.TrimSpace(expected[i]) {
			return fmt.Errorf("ヘッダ不一致です(列:%d 期待:%s 実際:%s)", i+1, expected[i], header[i])
		}
	}
	return nil
}
