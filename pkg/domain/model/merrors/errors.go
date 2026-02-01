// 指示: miu200521358
package merrors

import (
	"errors"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

const (
	invalidIndexErrorID     = "12201"
	nameNotFoundErrorID     = "20001"
	parentNotFoundErrorID   = "22201"
	indexOutOfRangeErrorID  = "92202"
	nameConflictErrorID     = "92203"
	nameMismatchErrorID     = "92204"
	modelCopyFailedErrorID  = "92201"
)

// IndexOutOfRangeError はインデックス範囲外エラーを表す。
type IndexOutOfRangeError struct {
	*merr.CommonError
	Index  int
	Length int
}

// NewIndexOutOfRangeError は IndexOutOfRangeError を生成する。
func NewIndexOutOfRangeError(index, length int) *IndexOutOfRangeError {
	message := "インデックスが範囲外です: index=%d length=%d"
	return &IndexOutOfRangeError{
		CommonError: merr.NewCommonError(indexOutOfRangeErrorID, merr.ErrorKindInternal, message, nil, index, length),
		Index:       index,
		Length:      length,
	}
}

// Error はエラーメッセージを返す。
func (e *IndexOutOfRangeError) Error() string {
	if e == nil || e.CommonError == nil {
		return ""
	}
	return e.CommonError.Error()
}

// IsIndexOutOfRangeError は err が IndexOutOfRangeError か判定する。
func IsIndexOutOfRangeError(err error) bool {
	var target *IndexOutOfRangeError
	return errors.As(err, &target)
}

// NameNotFoundError は名前未検出エラーを表す。
type NameNotFoundError struct {
	*merr.CommonError
	Name string
}

// NewNameNotFoundError は NameNotFoundError を生成する。
func NewNameNotFoundError(name string) *NameNotFoundError {
	message := "名前が見つかりません: %s"
	return &NameNotFoundError{
		CommonError: merr.NewCommonError(nameNotFoundErrorID, merr.ErrorKindNotFound, message, nil, name),
		Name:        name,
	}
}

// Error はエラーメッセージを返す。
func (e *NameNotFoundError) Error() string {
	if e == nil || e.CommonError == nil {
		return ""
	}
	return e.CommonError.Error()
}

// IsNameNotFoundError は err が NameNotFoundError か判定する。
func IsNameNotFoundError(err error) bool {
	var target *NameNotFoundError
	return errors.As(err, &target)
}

// NameConflictError は名前衝突エラーを表す。
type NameConflictError struct {
	*merr.CommonError
	Name string
}

// NewNameConflictError は NameConflictError を生成する。
func NewNameConflictError(name string) *NameConflictError {
	message := "名称が既存要素と衝突しました: %s"
	return &NameConflictError{
		CommonError: merr.NewCommonError(nameConflictErrorID, merr.ErrorKindInternal, message, nil, name),
		Name:        name,
	}
}

// Error はエラーメッセージを返す。
func (e *NameConflictError) Error() string {
	if e == nil || e.CommonError == nil {
		return ""
	}
	return e.CommonError.Error()
}

// IsNameConflictError は err が NameConflictError か判定する。
func IsNameConflictError(err error) bool {
	var target *NameConflictError
	return errors.As(err, &target)
}

// NameMismatchError は更新時の名前不一致エラーを表す。
type NameMismatchError struct {
	*merr.CommonError
	Index    int
	Expected string
	Actual   string
}

// NewNameMismatchError は NameMismatchError を生成する。
func NewNameMismatchError(index int, expected, actual string) *NameMismatchError {
	message := "名称が一致しません: index=%d expected=%s actual=%s"
	return &NameMismatchError{
		CommonError: merr.NewCommonError(nameMismatchErrorID, merr.ErrorKindInternal, message, nil, index, expected, actual),
		Index:       index,
		Expected:    expected,
		Actual:      actual,
	}
}

// Error はエラーメッセージを返す。
func (e *NameMismatchError) Error() string {
	if e == nil || e.CommonError == nil {
		return ""
	}
	return e.CommonError.Error()
}

// IsNameMismatchError は err が NameMismatchError か判定する。
func IsNameMismatchError(err error) bool {
	var target *NameMismatchError
	return errors.As(err, &target)
}

// ParentNotFoundError は親ボーン未検出エラーを表す。
type ParentNotFoundError struct {
	*merr.CommonError
	Parent     string
	Candidates []string
}

// NewParentNotFoundError は ParentNotFoundError を生成する。
func NewParentNotFoundError(parent string, candidates []string) *ParentNotFoundError {
	message := "親要素が見つかりません: %s"
	detail := parent
	if len(candidates) > 0 {
		detail = strings.Join(candidates, ",")
	}
	return &ParentNotFoundError{
		CommonError: merr.NewCommonError(parentNotFoundErrorID, merr.ErrorKindNotFound, message, nil, detail),
		Parent:      parent,
		Candidates:  candidates,
	}
}

// Error はエラーメッセージを返す。
func (e *ParentNotFoundError) Error() string {
	if e == nil || e.CommonError == nil {
		return ""
	}
	return e.CommonError.Error()
}

// IsParentNotFoundError は err が ParentNotFoundError か判定する。
func IsParentNotFoundError(err error) bool {
	var target *ParentNotFoundError
	return errors.As(err, &target)
}

// NewInvalidIndexError は無効なインデックスエラーを生成する。
func NewInvalidIndexError(index int) *merr.CommonError {
	return merr.NewCommonError(invalidIndexErrorID, merr.ErrorKindValidate, "インデックスが無効です: %d", nil, index)
}

// ModelCopyFailed はモデルコピー失敗を表す。
type ModelCopyFailed struct {
	*merr.CommonError
}

// NewModelCopyFailed は ModelCopyFailed を生成する。
func NewModelCopyFailed(cause error) *ModelCopyFailed {
	message := "モデルコピーに失敗"
	return &ModelCopyFailed{CommonError: merr.NewCommonError(modelCopyFailedErrorID, merr.ErrorKindInternal, message, cause)}
}

// Error はエラーメッセージを返す。
func (e *ModelCopyFailed) Error() string {
	if e == nil || e.CommonError == nil {
		return ""
	}
	return e.CommonError.Error()
}

// Unwrap は元の原因エラーを返す。
func (e *ModelCopyFailed) Unwrap() error {
	if e == nil || e.CommonError == nil {
		return nil
	}
	return e.CommonError.Unwrap()
}

// IsModelCopyFailed は err が ModelCopyFailed か判定する。
func IsModelCopyFailed(err error) bool {
	var target *ModelCopyFailed
	return errors.As(err, &target)
}
