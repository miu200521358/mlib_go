package mcore

import (
	"reflect"

	"github.com/miu200521358/mlib_go/pkg/domain/merr"
)

// IIndexModels はインデックスベースコレクションのインターフェースです。
type IIndexModels[T IIndexModel] interface {
	// Get は指定されたインデックスの値を取得します。
	Get(index int) (T, error)
	// Update は指定された値をそのインデックスに基づいて更新します。
	Update(value T) error
	// Append は新しい値をコレクションに追加します。
	Append(value T) error
	// Remove は指定されたインデックスの値を削除します。
	Remove(index int) error
	// Length はコレクション内の要素数を返します。
	Length() int
	// Contains は指定されたインデックスに有効な値が存在するかを確認します。
	Contains(index int) bool
	// ForEach は全ての値をコールバック関数に渡します。
	ForEach(callback func(index int, value T) bool)
	// Values は全要素のスライスを返します。
	Values() []T
	// First は最初の要素を返します。
	First() (T, error)
	// Last は最後の要素を返します。
	Last() (T, error)
	// Clear は全要素をクリアします。
	Clear()
	// IsEmpty はコレクションが空かどうかを返します。
	IsEmpty() bool
	// IsNotEmpty はコレクションが空でないかどうかを返します。
	IsNotEmpty() bool
}

// IndexModels はインデックスベースのコレクション実装です。
type IndexModels[T IIndexModel] struct {
	values []T
}

// NewIndexModels は指定された要素数を持つ IndexModels インスタンスを作成します。
func NewIndexModels[T IIndexModel](length int) (*IndexModels[T], error) {
	if length < 0 {
		return nil, merr.NewInvalidArgumentError("length", "must be non-negative")
	}
	return &IndexModels[T]{
		values: make([]T, length),
	}, nil
}

// NewIndexModelsWithCapacity は指定された要素数と容量を持つ IndexModels インスタンスを作成します。
func NewIndexModelsWithCapacity[T IIndexModel](length, capacity int) (*IndexModels[T], error) {
	if length < 0 {
		return nil, merr.NewInvalidArgumentError("length", "must be non-negative")
	}
	if capacity < length {
		return nil, merr.NewInvalidArgumentError("capacity", "must be greater than or equal to length")
	}
	return &IndexModels[T]{
		values: make([]T, length, capacity),
	}, nil
}

// Get は指定されたインデックスの値を取得します。
func (im *IndexModels[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(im.values) {
		var zero T
		return zero, merr.NewIndexOutOfRangeError(index, len(im.values)-1)
	}
	return im.values[index], nil
}

// Update は指定された値をそのインデックスに基づいて更新します。
func (im *IndexModels[T]) Update(value T) error {
	index := value.Index()
	if index < 0 || index >= len(im.values) {
		return merr.NewIndexOutOfRangeError(index, len(im.values)-1)
	}
	im.values[index] = value
	return nil
}

// Append は新しい値をコレクションに追加します。
// valueのインデックスが未設定（-1）の場合は自動で設定されます。
func (im *IndexModels[T]) Append(value T) error {
	if value.Index() < 0 {
		value.SetIndex(len(im.values))
	} else if value.Index() < len(im.values) {
		return merr.NewInvalidOperationError("use Update() for existing index")
	}
	im.values = append(im.values, value)
	return nil
}

// Remove は指定されたインデックスの値を削除します。
func (im *IndexModels[T]) Remove(index int) error {
	if index < 0 || index >= len(im.values) {
		return merr.NewIndexOutOfRangeError(index, len(im.values)-1)
	}
	im.values = append(im.values[:index], im.values[index+1:]...)
	return nil
}

// Length はコレクション内の要素数を返します。
func (im *IndexModels[T]) Length() int {
	return len(im.values)
}

// Contains は指定されたインデックスに有効な値が存在するかを確認します。
func (im *IndexModels[T]) Contains(index int) bool {
	if index < 0 || index >= len(im.values) {
		return false
	}
	v := reflect.ValueOf(im.values[index])
	if !v.IsValid() || v.IsNil() {
		return false
	}
	return im.values[index].IsValid()
}

// ForEach は全ての値をコールバック関数に渡します。
// callbackがfalseを返すと処理を中断します。
func (im *IndexModels[T]) ForEach(callback func(index int, value T) bool) {
	for i, v := range im.values {
		if !callback(i, v) {
			break
		}
	}
}

// Values は全要素のスライスを返します。
func (im *IndexModels[T]) Values() []T {
	return im.values
}

// First は最初の要素を返します。
func (im *IndexModels[T]) First() (T, error) {
	if len(im.values) == 0 {
		var zero T
		return zero, merr.NewIndexOutOfRangeError(0, -1)
	}
	return im.values[0], nil
}

// Last は最後の要素を返します。
func (im *IndexModels[T]) Last() (T, error) {
	if len(im.values) == 0 {
		var zero T
		return zero, merr.NewIndexOutOfRangeError(0, -1)
	}
	return im.values[len(im.values)-1], nil
}

// Clear は全要素をクリアします。
func (im *IndexModels[T]) Clear() {
	im.values = im.values[:0]
}

// IsEmpty はコレクションが空かどうかを返します。
func (im *IndexModels[T]) IsEmpty() bool {
	return len(im.values) == 0
}

// IsNotEmpty はコレクションが空でないかどうかを返します。
func (im *IndexModels[T]) IsNotEmpty() bool {
	return len(im.values) > 0
}
