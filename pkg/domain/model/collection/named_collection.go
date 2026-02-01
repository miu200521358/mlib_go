// 指示: miu200521358
package collection

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model/merrors"
)

// NamedCollection は name と index の参照を提供する。
type NamedCollection[T INameable] struct {
	indexed   *IndexedCollection[T]
	nameIndex *NameIndex[T]
}

// NewNamedCollection は capacity を指定して NamedCollection を生成する。
func NewNamedCollection[T INameable](capacity int) *NamedCollection[T] {
	return &NamedCollection[T]{
		indexed:   NewIndexedCollection[T](capacity),
		nameIndex: NewNameIndex[T](),
	}
}

// Len は要素数を返す。
func (c *NamedCollection[T]) Len() int {
	return c.indexed.Len()
}

// Values は内部スライスを返す。
func (c *NamedCollection[T]) Values() []T {
	return c.indexed.values
}

// Get は index の要素を返す。
func (c *NamedCollection[T]) Get(index int) (T, error) {
	return c.indexed.Get(index)
}

// GetByName は name の要素を返す。
func (c *NamedCollection[T]) GetByName(name string) (T, error) {
	var zero T
	idx, ok := c.nameIndex.GetByName(name)
	if !ok {
		return zero, merrors.NewNameNotFoundError(name)
	}
	return c.Get(idx)
}

// Append は要素を追加し、NameIndex を更新する。
func (c *NamedCollection[T]) Append(value T) (int, ReindexResult) {
	idx, res := c.indexed.Append(value)
	if value.IsValid() {
		c.nameIndex.SetIfAbsent(value.Name(), idx)
	}
	return idx, res
}

// AppendRaw は再インデックス情報を作らずに追加し、NameIndex を更新する。
func (c *NamedCollection[T]) AppendRaw(value T) int {
	if c == nil {
		return -1
	}
	idx := c.indexed.AppendRaw(value)
	if value.IsValid() {
		c.nameIndex.SetIfAbsent(value.Name(), idx)
	}
	return idx
}

// Insert は insertIndex に挿入し、NameIndex を再構築する。
func (c *NamedCollection[T]) Insert(value T, insertIndex int) (int, ReindexResult, error) {
	idx, res, err := c.indexed.Insert(value, insertIndex)
	if err != nil {
		return 0, ReindexResult{}, err
	}
	c.nameIndex.Rebuild(c.indexed.values)
	return idx, res, nil
}

// Remove は要素を削除し、NameIndex を再構築する。
func (c *NamedCollection[T]) Remove(index int) (ReindexResult, error) {
	res, err := c.indexed.Remove(index)
	if err != nil {
		return ReindexResult{}, err
	}
	c.nameIndex.Rebuild(c.indexed.values)
	return res, nil
}

// Update は名前を変えずに要素を置き換える。
func (c *NamedCollection[T]) Update(index int, value T) (ReindexResult, error) {
	existing, err := c.Get(index)
	if err != nil {
		return ReindexResult{}, err
	}
	if existing.Name() != value.Name() {
		return ReindexResult{}, merrors.NewNameMismatchError(index, existing.Name(), value.Name())
	}
	return c.indexed.Update(index, value)
}

// Rename は要素名を変更し、NameIndex を更新する。
func (c *NamedCollection[T]) Rename(index int, newName string) (bool, error) {
	value, err := c.Get(index)
	if err != nil {
		return false, err
	}
	oldName := value.Name()
	if oldName == newName {
		return false, nil
	}
	if idx, ok := c.nameIndex.GetByName(newName); ok && idx != index {
		return false, merrors.NewNameConflictError(newName)
	}
	value.SetName(newName)
	c.nameIndex.Rebuild(c.indexed.values)
	return true, nil
}

// Contains は index の要素が有効か判定する。
func (c *NamedCollection[T]) Contains(index int) bool {
	return c.indexed.Contains(index)
}
