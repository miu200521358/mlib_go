package core

import (
	"errors"
	"reflect"
)

// IIndexModel はインデックス操作を行うモデルのインターフェースを定義します
type IIndexModel interface {
	IsValid() bool
	Index() int
	SetIndex(index int)
}

type IndexModels[T IIndexModel] struct {
	values []T // モデルのスライス
}

// NewIndexModels は指定された初期個数を持つ IndexModels インスタンスを作成します
func NewIndexModels[T IIndexModel](capacity int) *IndexModels[T] {
	values := make([]T, capacity)
	return &IndexModels[T]{
		values: values,
	}
}

// Get は指定されたインデックスの値を取得します
func (im *IndexModels[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(im.values) {
		return *new(T), errors.New("index out of range")
	}
	return im.values[index], nil
}

// Update は指定された値をそのインデックスに基づいて更新します
func (im *IndexModels[T]) Update(value T) error {
	if value.Index() < 0 || value.Index() >= len(im.values) {
		// Update は index指定必須
		return errors.New("invalid index")
	}
	im.values[value.Index()] = value
	return nil
}

// Append は新しい値をコレクションに追加し、必要に応じてインデックスを設定します
func (im *IndexModels[T]) Append(value T) error {
	if value.Index() < 0 {
		// Append は index指定が無い場合、自動で設定
		value.SetIndex(len(im.values))
	}
	im.values = append(im.values, value)
	return nil
}

// Remove は指定されたインデックスの値を削除します
func (im *IndexModels[T]) Remove(index int) error {
	if index < 0 || index >= len(im.values) {
		return errors.New("index out of range") // インデックスが範囲外の場合にエラーを返します
	}
	im.values = append(im.values[:index], im.values[index+1:]...)
	return nil
}

// Length はコレクション内の要素数を返します
func (im *IndexModels[T]) Length() int {
	return len(im.values)
}

// Contains は指定されたインデックスに値が存在するかを確認します
func (im *IndexModels[T]) Contains(index int) bool {
	if index < 0 || index >= len(im.values) {
		return false
	}
	return !reflect.ValueOf(im.values[index]).IsNil()
}

// ForEach は全ての値をコールバック関数に渡します
func (im *IndexModels[T]) ForEach(callback func(index int, value T)) {
	for i, v := range im.values {
		callback(i, v)
	}
}
