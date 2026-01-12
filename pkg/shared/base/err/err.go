// 指示: miu200521358
package err

import (
	"embed"
	"encoding/csv"
	"io"
	"strings"
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

// CommonError は ErrorID 付きの共通エラー。
type CommonError struct {
	ID      string
	Kind    ErrorKind
	Message string
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
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + ": " + e.Cause.Error()
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
func NewCommonError(id string, kind ErrorKind, message string, cause error) *CommonError {
	return &CommonError{
		ID:      id,
		Kind:    kind,
		Message: message,
		Cause:   cause,
	}
}

// NewOsPackageError は os パッケージ由来の共通委譲エラーを生成する。
func NewOsPackageError(message string, cause error) *CommonError {
	return NewCommonError(OsPackageErrorID, ErrorKindExternal, message, cause)
}

// NewJsonPackageError は json パッケージ由来の共通委譲エラーを生成する。
func NewJsonPackageError(message string, cause error) *CommonError {
	return NewCommonError(JsonPackageErrorID, ErrorKindExternal, message, cause)
}

// NewImagePackageError は image パッケージ由来の共通委譲エラーを生成する。
func NewImagePackageError(message string, cause error) *CommonError {
	return NewCommonError(ImagePackageErrorID, ErrorKindExternal, message, cause)
}

// NewFsPackageError は io/fs パッケージ由来の共通委譲エラーを生成する。
func NewFsPackageError(message string, cause error) *CommonError {
	return NewCommonError(FsPackageErrorID, ErrorKindExternal, message, cause)
}

// NewDeepcopyPackageError は go-deepcopy パッケージ由来の共通委譲エラーを生成する。
func NewDeepcopyPackageError(message string, cause error) *CommonError {
	return NewCommonError(DeepcopyPackageErrorID, ErrorKindExternal, message, cause)
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
