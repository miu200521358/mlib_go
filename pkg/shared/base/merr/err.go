// 指示: miu200521358
package merr

import (
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

//go:embed error_registry.csv
var registryFiles embed.FS
var openRegistryFile = registryFiles.Open

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

// ExtractErrorID はエラーからErrorIDを取得する。
func ExtractErrorID(err error) string {
	if err == nil {
		return ""
	}
	var provider interface {
		ErrorID() string
	}
	if errors.As(err, &provider) {
		return provider.ErrorID()
	}
	return ""
}

// CommonError は ErrorID 付きの共通エラー。
type CommonError struct {
	ID      string
	Kind    ErrorKind
	Message string
	Params  []any
	Cause   error
}

// Error はエラーメッセージを返す。
func (e *CommonError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" && e.Cause == nil {
		return ""
	}
	if e.Message == "" {
		return e.Cause.Error()
	}
	msg := e.Message
	if len(e.Params) > 0 {
		msg = fmt.Sprintf(msg, e.Params...)
	}
	if e.Cause == nil {
		return msg
	}
	return msg + ": " + e.Cause.Error()
}

// Unwrap は原因エラーを返す。
func (e *CommonError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// ErrorID はエラーIDを返す。
func (e *CommonError) ErrorID() string {
	if e == nil {
		return ""
	}
	return e.ID
}

// ErrorKind はエラー種別を返す。
func (e *CommonError) ErrorKind() ErrorKind {
	if e == nil {
		return ""
	}
	return e.Kind
}

// MessageKey はエラーメッセージのキーを返す。
func (e *CommonError) MessageKey() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// MessageParams はエラーメッセージのパラメータを返す。
func (e *CommonError) MessageParams() []any {
	if e == nil {
		return nil
	}
	return e.Params
}

const (
	// OsPackageErrorID は os パッケージの共通委譲エラーID。
	OsPackageErrorID = "99001"
	// JsonPackageErrorID は json パッケージの共通委譲エラーID。
	JsonPackageErrorID = "99002"
	// ImagePackageErrorID は image パッケージの共通委譲エラーID。
	ImagePackageErrorID = "99003"
	// FsPackageErrorID は io/fs パッケージの共通委譲エラーID。
	FsPackageErrorID = "99004"
	// DeepcopyPackageErrorID は go-deepcopy パッケージの共通委譲エラーID。
	DeepcopyPackageErrorID = "99005"
)

// NewCommonError は ErrorID 付きの共通エラーを生成する。
func NewCommonError(id string, kind ErrorKind, message string, cause error, params ...any) *CommonError {
	return &CommonError{
		ID:      id,
		Kind:    kind,
		Message: message,
		Params:  params,
		Cause:   cause,
	}
}

// NewOsPackageError は os パッケージ由来の共通委譲エラーを生成する。
func NewOsPackageError(message string, cause error, params ...any) *CommonError {
	return NewCommonError(OsPackageErrorID, ErrorKindExternal, message, cause, params...)
}

// NewJsonPackageError は json パッケージ由来の共通委譲エラーを生成する。
func NewJsonPackageError(message string, cause error, params ...any) *CommonError {
	return NewCommonError(JsonPackageErrorID, ErrorKindExternal, message, cause, params...)
}

// NewImagePackageError は image パッケージ由来の共通委譲エラーを生成する。
func NewImagePackageError(message string, cause error, params ...any) *CommonError {
	return NewCommonError(ImagePackageErrorID, ErrorKindExternal, message, cause, params...)
}

// NewFsPackageError は io/fs パッケージ由来の共通委譲エラーを生成する。
func NewFsPackageError(message string, cause error, params ...any) *CommonError {
	return NewCommonError(FsPackageErrorID, ErrorKindExternal, message, cause, params...)
}

// NewDeepcopyPackageError は go-deepcopy パッケージ由来の共通委譲エラーを生成する。
func NewDeepcopyPackageError(message string, cause error, params ...any) *CommonError {
	return NewCommonError(DeepcopyPackageErrorID, ErrorKindExternal, message, cause, params...)
}

// ErrorRegistryPath は埋め込みCSVのパス。
const ErrorRegistryPath = "error_registry.csv"

// LoadDefaultRegistry は埋め込みCSVから読み込む。
func LoadDefaultRegistry() ([]ErrorRecord, error) {
	file, err := openRegistryFile(ErrorRegistryPath)
	if err != nil {
		return nil, err
	}
	recs, err := LoadRegistry(file)
	if closeErr := file.Close(); err == nil && closeErr != nil {
		err = closeErr
	}
	return recs, err
}

// LoadRegistry はCSVからエラー管理表を読み込む。
func LoadRegistry(r io.Reader) ([]ErrorRecord, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

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

var (
	defaultRegistryOnce sync.Once
	defaultRegistryMap  map[string]ErrorRecord
	defaultRegistryErr  error
)

// DefaultRegistryMap は埋め込みCSVから読み込んだエラー管理表を返す。
func DefaultRegistryMap() (map[string]ErrorRecord, error) {
	defaultRegistryOnce.Do(func() {
		recs, err := LoadDefaultRegistry()
		if err != nil {
			defaultRegistryErr = err
			return
		}
		registry := make(map[string]ErrorRecord, len(recs))
		for _, rec := range recs {
			if rec.ID == "" {
				continue
			}
			registry[rec.ID] = rec
		}
		defaultRegistryMap = registry
	})
	return defaultRegistryMap, defaultRegistryErr
}

// FindRecord は ErrorID からエラー管理表のレコードを取得する。
func FindRecord(id string) (*ErrorRecord, error) {
	if id == "" {
		return nil, nil
	}
	registry, err := DefaultRegistryMap()
	if err != nil {
		return nil, err
	}
	rec, ok := registry[id]
	if !ok {
		return nil, nil
	}
	out := rec
	return &out, nil
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
