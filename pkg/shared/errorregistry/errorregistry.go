// 指示: miu200521358
package errorregistry

import (
	"embed"
	"encoding/csv"
	"io"
	"strings"
)

//go:embed error_registry.csv
var registryFiles embed.FS

// ErrorKind はエラー種別を表す。
type ErrorKind string

const (
	// ErrorKindValidate は検証エラー。
	ErrorKindValidate ErrorKind = "Validate"
	// ErrorKindNotFound は存在しないエラー。
	ErrorKindNotFound ErrorKind = "NotFound"
	// ErrorKindNotSupported は未対応エラー。
	ErrorKindNotSupported ErrorKind = "NotSupported"
	// ErrorKindExternal は外部要因エラー。
	ErrorKindExternal ErrorKind = "External"
	// ErrorKindInternal は内部エラー。
	ErrorKindInternal ErrorKind = "Internal"
)

// ErrorRecord はエラー管理表の1行を表す。
type ErrorRecord struct {
	ID          string
	Kind        ErrorKind
	Layer       string
	Module      string
	ErrorName   string
	Summary     string
	Remedy      string
	SourcePaths []string
}

// ErrorRegistryPath は埋め込みCSVのパス。
const ErrorRegistryPath = "error_registry.csv"

// LoadDefaultRegistry は埋め込みCSVから読み込む。
func LoadDefaultRegistry() ([]ErrorRecord, error) {
	file, err := registryFiles.Open(ErrorRegistryPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadRegistry(file)
}

// LoadRegistry はCSVからエラー管理表を読み込む。
func LoadRegistry(r io.Reader) ([]ErrorRecord, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []ErrorRecord{}, nil
	}

	start := 0
	if len(records[0]) > 0 && strings.EqualFold(records[0][0], "ID") {
		start = 1
	}

	out := make([]ErrorRecord, 0, len(records)-start)
	for i := start; i < len(records); i++ {
		row := records[i]
		if len(row) < 8 {
			continue
		}
		out = append(out, ErrorRecord{
			ID:          row[0],
			Kind:        ErrorKind(row[1]),
			Layer:       row[2],
			Module:      row[3],
			ErrorName:   row[4],
			Summary:     row[5],
			Remedy:      row[6],
			SourcePaths: splitPaths(row[7]),
		})
	}
	return out, nil
}

// splitPaths はセミコロン区切りの参照パスを分割する。
func splitPaths(raw string) []string {
	if raw == "" || raw == "-" {
		return nil
	}
	parts := strings.Split(raw, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
