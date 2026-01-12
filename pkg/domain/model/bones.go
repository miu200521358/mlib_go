// 指示: miu200521358
package model

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
	modelerrors "github.com/miu200521358/mlib_go/pkg/domain/model/errors"
)

// BoneCollection は Insert 特例を持つボーンコレクション。
type BoneCollection struct {
	values     []*Bone
	nameIndex  *collection.NameIndex[*Bone]
	indexToPos map[int]int
}

// NewBoneCollection は capacity を指定して BoneCollection を生成する。
func NewBoneCollection(capacity int) *BoneCollection {
	return &BoneCollection{
		values:     make([]*Bone, 0, capacity),
		nameIndex:  collection.NewNameIndex[*Bone](),
		indexToPos: make(map[int]int),
	}
}

// Len は要素数を返す。
func (c *BoneCollection) Len() int {
	return len(c.values)
}

// Values は内部スライス順を返す。
func (c *BoneCollection) Values() []*Bone {
	return c.values
}

// Get は index のボーンを返す。
func (c *BoneCollection) Get(index int) (*Bone, error) {
	if index < 0 || index >= len(c.values) {
		return nil, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}
	pos, ok := c.indexToPos[index]
	if !ok {
		return nil, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}
	return c.values[pos], nil
}

// GetByName は name のボーンを返す。
func (c *BoneCollection) GetByName(name string) (*Bone, error) {
	idx, ok := c.nameIndex.GetByName(name)
	if !ok {
		return nil, modelerrors.NewNameNotFoundError(name)
	}
	return c.Get(idx)
}

// Append は末尾にボーンを追加する。
func (c *BoneCollection) Append(value *Bone) (int, collection.ReindexResult) {
	oldLen := len(c.values)
	value.SetIndex(oldLen)
	c.values = append(c.values, value)
	c.indexToPos[oldLen] = len(c.values) - 1
	if value.IsValid() {
		c.nameIndex.SetIfAbsent(value.Name(), value.Index())
	}
	oldToNew, newToOld := identityMappings(oldLen)
	return oldLen, collection.ReindexResult{
		Changed:  false,
		OldToNew: oldToNew,
		NewToOld: newToOld,
		Added:    []int{oldLen},
	}
}

// Insert は既存 index を変えずにボーンを挿入する。
func (c *BoneCollection) Insert(value *Bone, afterBoneIndex int) (int, collection.ReindexResult, error) {
	oldLen := len(c.values)
	value.SetIndex(oldLen)
	insertPos := oldLen
	if afterBoneIndex < 0 {
		for _, bone := range c.values {
			if bone != nil {
				bone.Layer++
			}
		}
		value.Layer = 0
	} else if pos, ok := c.indexToPos[afterBoneIndex]; ok {
		insertPos = pos + 1
		prevLayer := 0
		if c.values[pos] != nil {
			prevLayer = c.values[pos].Layer
		}
		if insertPos >= len(c.values) {
			value.Layer = prevLayer
		} else {
			nextLayer := prevLayer
			if c.values[insertPos] != nil {
				nextLayer = c.values[insertPos].Layer
			}
			if nextLayer <= prevLayer+1 {
				value.Layer = prevLayer
			} else {
				value.Layer = prevLayer + 1
			}
		}
	} else if oldLen > 0 {
		last := c.values[oldLen-1]
		if last != nil {
			value.Layer = last.Layer
		}
	}

	c.values = append(c.values, value)
	if insertPos < len(c.values)-1 {
		copy(c.values[insertPos+1:], c.values[insertPos:])
		c.values[insertPos] = value
	}
	for i := insertPos; i < len(c.values); i++ {
		bone := c.values[i]
		if bone != nil {
			c.indexToPos[bone.Index()] = i
		}
	}
	if value.IsValid() {
		c.nameIndex.SetIfAbsent(value.Name(), value.Index())
	}
	oldToNew, newToOld := identityMappings(oldLen)
	return value.Index(), collection.ReindexResult{
		Changed:  false,
		OldToNew: oldToNew,
		NewToOld: newToOld,
		Added:    []int{value.Index()},
	}, nil
}

// Remove は index のボーンを削除し、残りを再インデックスする。
func (c *BoneCollection) Remove(index int) (collection.ReindexResult, error) {
	if index < 0 || index >= len(c.values) {
		return collection.ReindexResult{}, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}
	pos, ok := c.indexToPos[index]
	if !ok {
		return collection.ReindexResult{}, modelerrors.NewIndexOutOfRangeError(index, len(c.values))
	}
	oldLen := len(c.values)
	copy(c.values[pos:], c.values[pos+1:])
	c.values = c.values[:oldLen-1]
	for _, bone := range c.values {
		if bone != nil && bone.Index() > index {
			bone.SetIndex(bone.Index() - 1)
		}
	}
	c.rebuildIndexToPos()
	c.rebuildNameIndex()

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
	return collection.ReindexResult{
		Changed:  true,
		OldToNew: oldToNew,
		NewToOld: newToOld,
		Removed:  []int{index},
	}, nil
}

// Update は名前を変えずにボーンを置き換える。
func (c *BoneCollection) Update(index int, value *Bone) (collection.ReindexResult, error) {
	current, err := c.Get(index)
	if err != nil {
		return collection.ReindexResult{}, err
	}
	if current.Name() != value.Name() {
		return collection.ReindexResult{}, modelerrors.NewNameMismatchError(index, current.Name(), value.Name())
	}
	pos := c.indexToPos[index]
	value.SetIndex(index)
	c.values[pos] = value
	oldToNew, newToOld := identityMappings(len(c.values))
	return collection.ReindexResult{
		Changed:  false,
		OldToNew: oldToNew,
		NewToOld: newToOld,
	}, nil
}

// Rename はボーン名を変更し、NameIndex を再構築する。
func (c *BoneCollection) Rename(index int, newName string) (bool, error) {
	value, err := c.Get(index)
	if err != nil {
		return false, err
	}
	oldName := value.Name()
	if oldName == newName {
		return false, nil
	}
	if idx, ok := c.nameIndex.GetByName(newName); ok && idx != index {
		return false, modelerrors.NewNameConflictError(newName)
	}
	value.SetName(newName)
	c.rebuildNameIndex()
	return true, nil
}

// Contains は index のボーンが有効か判定する。
func (c *BoneCollection) Contains(index int) bool {
	if index < 0 || index >= len(c.values) {
		return false
	}
	pos, ok := c.indexToPos[index]
	if !ok {
		return false
	}
	return c.values[pos].IsValid()
}

// rebuildIndexToPos は values の順序から indexToPos を再構築する。
func (c *BoneCollection) rebuildIndexToPos() {
	c.indexToPos = make(map[int]int, len(c.values))
	for pos, bone := range c.values {
		if bone == nil {
			continue
		}
		c.indexToPos[bone.Index()] = pos
	}
}

// rebuildNameIndex は現在の indexToPos を基に NameIndex を再構築する。
func (c *BoneCollection) rebuildNameIndex() {
	ordered := make([]*Bone, len(c.values))
	for idx := 0; idx < len(c.values); idx++ {
		pos, ok := c.indexToPos[idx]
		if !ok {
			continue
		}
		ordered[idx] = c.values[pos]
	}
	c.nameIndex.Rebuild(ordered)
}

// identityMappings は length 分の恒等マッピングを生成して返す。
func identityMappings(length int) ([]int, []int) {
	oldToNew := make([]int, length)
	newToOld := make([]int, length)
	for i := 0; i < length; i++ {
		oldToNew[i] = i
		newToOld[i] = i
	}
	return oldToNew, newToOld
}
