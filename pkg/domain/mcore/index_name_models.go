package mcore

import (
	"reflect"

	"github.com/miu200521358/mlib_go/pkg/domain/merr"
)

// IIndexNameModels は名前付きコレクションのインターフェースです。
type IIndexNameModels[T IIndexNameModel] interface {
	IIndexModels[T]
	// GetByName は指定された名前の値を取得します。
	GetByName(name string) (T, error)
	// ContainsByName は指定された名前に値が存在するかを確認します。
	ContainsByName(name string) bool
	// RemoveByName は指定された名前の値を削除します。
	RemoveByName(name string) error
	// Names は全要素の名前のスライスを返します。
	Names() []string
	// Indexes は全要素のインデックスのスライスを返します。
	Indexes() []int
}

// IndexNameModels は名前とインデックスベースのコレクション実装です。
type IndexNameModels[T IIndexNameModel] struct {
	values      []T
	nameIndexes map[string]int
}

// NewIndexNameModels は指定された要素数を持つ IndexNameModels インスタンスを作成します。
func NewIndexNameModels[T IIndexNameModel](length int) (*IndexNameModels[T], error) {
	if length < 0 {
		return nil, merr.NewInvalidArgumentError("length", "must be non-negative")
	}
	return &IndexNameModels[T]{
		values:      make([]T, length),
		nameIndexes: make(map[string]int),
	}, nil
}

// NewIndexNameModelsWithCapacity は指定された要素数と容量を持つ IndexNameModels インスタンスを作成します。
func NewIndexNameModelsWithCapacity[T IIndexNameModel](length, capacity int) (*IndexNameModels[T], error) {
	if length < 0 {
		return nil, merr.NewInvalidArgumentError("length", "must be non-negative")
	}
	if capacity < length {
		return nil, merr.NewInvalidArgumentError("capacity", "must be greater than or equal to length")
	}
	return &IndexNameModels[T]{
		values:      make([]T, length, capacity),
		nameIndexes: make(map[string]int),
	}, nil
}

// Get は指定されたインデックスの値を取得します。
func (im *IndexNameModels[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(im.values) {
		var zero T
		return zero, merr.NewIndexOutOfRangeError(index, len(im.values)-1)
	}
	return im.values[index], nil
}

// GetByName は指定された名前の値を取得します。
func (im *IndexNameModels[T]) GetByName(name string) (T, error) {
	if index, ok := im.nameIndexes[name]; ok {
		return im.values[index], nil
	}
	var zero T
	return zero, merr.NewNameNotFoundError(name)
}

// Update は指定された値をそのインデックスに基づいて更新します。
func (im *IndexNameModels[T]) Update(value T) error {
	index := value.Index()
	if index < 0 || index >= len(im.values) {
		return merr.NewIndexOutOfRangeError(index, len(im.values)-1)
	}
	im.values[index] = value
	// 名前は先勝ち（既存の名前がない場合のみ登録）
	if _, ok := im.nameIndexes[value.Name()]; !ok {
		im.nameIndexes[value.Name()] = index
	}
	return nil
}

// Append は新しい値をコレクションに追加します。
// valueのインデックスが未設定（-1）の場合は自動で設定されます。
func (im *IndexNameModels[T]) Append(value T) error {
	if value.Index() < 0 {
		value.SetIndex(len(im.values))
	} else if value.Index() < len(im.values) {
		return merr.NewInvalidOperationError("use Update() for existing index")
	}
	im.values = append(im.values, value)
	// 名前は先勝ち
	if _, ok := im.nameIndexes[value.Name()]; !ok {
		im.nameIndexes[value.Name()] = value.Index()
	}
	return nil
}

// Remove は指定されたインデックスの値を削除します。
func (im *IndexNameModels[T]) Remove(index int) error {
	if index < 0 || index >= len(im.values) {
		return merr.NewIndexOutOfRangeError(index, len(im.values)-1)
	}
	name := im.values[index].Name()
	delete(im.nameIndexes, name)
	im.values = append(im.values[:index], im.values[index+1:]...)
	return nil
}

// RemoveByName は指定された名前の値を削除します。
func (im *IndexNameModels[T]) RemoveByName(name string) error {
	index, ok := im.nameIndexes[name]
	if !ok {
		return merr.NewNameNotFoundError(name)
	}
	delete(im.nameIndexes, name)
	im.values = append(im.values[:index], im.values[index+1:]...)
	return nil
}

// Length はコレクション内の要素数を返します。
func (im *IndexNameModels[T]) Length() int {
	return len(im.values)
}

// Contains は指定されたインデックスに有効な値が存在するかを確認します。
func (im *IndexNameModels[T]) Contains(index int) bool {
	if index < 0 || index >= len(im.values) {
		return false
	}
	v := reflect.ValueOf(im.values[index])
	if !v.IsValid() || v.IsNil() {
		return false
	}
	return im.values[index].IsValid()
}

// ContainsByName は指定された名前に値が存在するかを確認します。
func (im *IndexNameModels[T]) ContainsByName(name string) bool {
	_, ok := im.nameIndexes[name]
	return ok
}

// ForEach は全ての値をコールバック関数に渡します。
// callbackがfalseを返すと処理を中断します。
func (im *IndexNameModels[T]) ForEach(callback func(index int, value T) bool) {
	for i, v := range im.values {
		if !callback(i, v) {
			break
		}
	}
}

// Values は全要素のスライスを返します。
func (im *IndexNameModels[T]) Values() []T {
	return im.values
}

// Names は全要素の名前のスライスを返します。
func (im *IndexNameModels[T]) Names() []string {
	names := make([]string, len(im.values))
	for i := range im.values {
		if im.Contains(i) {
			names[i] = im.values[i].Name()
		}
	}
	return names
}

// Indexes は全要素のインデックスのスライスを返します。
func (im *IndexNameModels[T]) Indexes() []int {
	indexes := make([]int, len(im.values))
	for i := range im.values {
		indexes[i] = i
	}
	return indexes
}

// First は最初の要素を返します。
func (im *IndexNameModels[T]) First() (T, error) {
	if len(im.values) == 0 {
		var zero T
		return zero, merr.NewIndexOutOfRangeError(0, -1)
	}
	return im.values[0], nil
}

// Last は最後の要素を返します。
func (im *IndexNameModels[T]) Last() (T, error) {
	if len(im.values) == 0 {
		var zero T
		return zero, merr.NewIndexOutOfRangeError(0, -1)
	}
	return im.values[len(im.values)-1], nil
}

// Clear は全要素をクリアします。
func (im *IndexNameModels[T]) Clear() {
	im.values = im.values[:0]
	im.nameIndexes = make(map[string]int)
}

// IsEmpty はコレクションが空かどうかを返します。
func (im *IndexNameModels[T]) IsEmpty() bool {
	return len(im.values) == 0
}

// IsNotEmpty はコレクションが空でないかどうかを返します。
func (im *IndexNameModels[T]) IsNotEmpty() bool {
	return len(im.values) > 0
}

// UpdateNameIndexes は名前インデックスを再構築します。
func (im *IndexNameModels[T]) UpdateNameIndexes() {
	im.nameIndexes = make(map[string]int, len(im.values))
	for i := range im.values {
		if im.Contains(i) {
			im.nameIndexes[im.values[i].Name()] = i
		}
	}
}
