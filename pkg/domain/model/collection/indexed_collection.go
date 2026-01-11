package collection

import (
	modelerrors "github.com/miu200521358/mlib_go/pkg/domain/model/errors"
)

// IndexedCollection は index ベースの参照を提供する。
type IndexedCollection[T IIndexable] struct {
	values []T
}

// NewIndexedCollection は capacity を指定して IndexedCollection を生成する。
func NewIndexedCollection[T IIndexable](capacity int) *IndexedCollection[T] {
	return &IndexedCollection[T]{values: make([]T, 0, capacity)}
}

// Len は要素数を返す。
func (c *IndexedCollection[T]) Len() int {
	return len(c.values)
}

// Values は内部スライスを返す。
func (c *IndexedCollection[T]) Values() []T {
	return c.values
}

// Get は index の要素を返す。
func (c *IndexedCollection[T]) Get(index int) (T, error) {
	var zero T
	if index < 0 || index >= len(c.values) {
		return zero, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}
	return c.values[index], nil
}

// Append は末尾に追加し、新しい index を付与する。
func (c *IndexedCollection[T]) Append(value T) (int, ReindexResult) {
	oldLen := len(c.values)
	value.SetIndex(oldLen)
	c.values = append(c.values, value)

	oldToNew, newToOld := identityMappings(oldLen)
	return oldLen, ReindexResult{
		Changed:  false,
		OldToNew: oldToNew,
		NewToOld: newToOld,
		Added:    []int{oldLen},
	}
}

// Insert は insertIndex に挿入し、以降の要素を再インデックスする。
func (c *IndexedCollection[T]) Insert(value T, insertIndex int) (int, ReindexResult, error) {
	if insertIndex < 0 || insertIndex > len(c.values) {
		return 0, ReindexResult{}, modelerrors.NewIndexOutOfRangeError(insertIndex, len(c.values))
	}
	if insertIndex == len(c.values) {
		idx, res := c.Append(value)
		return idx, res, nil
	}

	oldLen := len(c.values)
	c.values = append(c.values, value)
	copy(c.values[insertIndex+1:], c.values[insertIndex:])
	c.values[insertIndex] = value
	for i := insertIndex; i < len(c.values); i++ {
		c.values[i].SetIndex(i)
	}

	oldToNew := make([]int, oldLen)
	newToOld := make([]int, oldLen)
	for i := 0; i < oldLen; i++ {
		if i < insertIndex {
			oldToNew[i] = i
		} else {
			oldToNew[i] = i + 1
		}
	}
	for i := 0; i < oldLen; i++ {
		switch {
		case i == insertIndex:
			newToOld[i] = -1
		case i < insertIndex:
			newToOld[i] = i
		default:
			newToOld[i] = i - 1
		}
	}

	return insertIndex, ReindexResult{
		Changed:  true,
		OldToNew: oldToNew,
		NewToOld: newToOld,
		Added:    []int{insertIndex},
	}, nil
}

// Remove は要素を削除し、残りを再インデックスする。
func (c *IndexedCollection[T]) Remove(index int) (ReindexResult, error) {
	if index < 0 || index >= len(c.values) {
		return ReindexResult{}, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}

	oldLen := len(c.values)
	copy(c.values[index:], c.values[index+1:])
	c.values = c.values[:oldLen-1]
	for i := index; i < len(c.values); i++ {
		c.values[i].SetIndex(i)
	}

	oldToNew := make([]int, oldLen)
	newToOld := make([]int, oldLen)
	for i := 0; i < oldLen; i++ {
		switch {
		case i < index:
			oldToNew[i] = i
		case i == index:
			oldToNew[i] = -1
		default:
			oldToNew[i] = i - 1
		}
	}
	newLen := oldLen - 1
	for i := 0; i < oldLen; i++ {
		if i >= newLen {
			newToOld[i] = -1
			continue
		}
		if i < index {
			newToOld[i] = i
		} else {
			newToOld[i] = i + 1
		}
	}

	return ReindexResult{
		Changed:  true,
		OldToNew: oldToNew,
		NewToOld: newToOld,
		Removed:  []int{index},
	}, nil
}

// Update は index を変えずに要素を置き換える。
func (c *IndexedCollection[T]) Update(index int, value T) (ReindexResult, error) {
	if index < 0 || index >= len(c.values) {
		return ReindexResult{}, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}
	value.SetIndex(index)
	c.values[index] = value

	oldToNew, newToOld := identityMappings(len(c.values))
	return ReindexResult{
		Changed:  false,
		OldToNew: oldToNew,
		NewToOld: newToOld,
	}, nil
}

// Contains は index の要素が有効か判定する。
func (c *IndexedCollection[T]) Contains(index int) bool {
	if index < 0 || index >= len(c.values) {
		return false
	}
	return c.values[index].IsValid()
}

func identityMappings(length int) ([]int, []int) {
	oldToNew := make([]int, length)
	newToOld := make([]int, length)
	for i := 0; i < length; i++ {
		oldToNew[i] = i
		newToOld[i] = i
	}
	return oldToNew, newToOld
}
