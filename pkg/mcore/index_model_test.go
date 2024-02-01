package mcore

import (
	"testing"

)

func TestIndexModel_GetIndex(t *testing.T) {
	model := &IndexModel{Index: 5}
	index := model.GetIndex()
	if index != 5 {
		t.Errorf("Expected index to be 5, but got %d", index)
	}
}

func TestIndexModel_SetIndex(t *testing.T) {
	model := &IndexModel{}
	model.SetIndex(10)
	if model.Index != 10 {
		t.Errorf("Expected index to be 10, but got %d", model.Index)
	}
}

func TestIndexModel_IsValid(t *testing.T) {
	model := &IndexModel{Index: 3}
	valid := model.IsValid()
	if !valid {
		t.Error("Expected IsValid to return true, but got false")
	}

	model.SetIndex(-1)
	valid = model.IsValid()
	if valid {
		t.Error("Expected IsValid to return false, but got true")
	}
}

func TestIndexModel_Copy(t *testing.T) {
	model := &IndexModel{Index: 7}
	copied := model.Copy().(*IndexModel)
	if copied.Index != 7 {
		t.Errorf("Expected copied index to be 7, but got %d", copied.Index)
	}

	// Modify the copied model
	copied.SetIndex(9)
	if model.Index == copied.Index {
		t.Errorf("Expected copied model to be a separate instance, but both have the same index %d", model.Index)
	}
}

type Face struct {
	IndexModel
	Index         int
	VertexIndexes [3]int
}

func NewFace(index, vertexIndex0, vertexIndex1, vertexIndex2 int) *Face {
	return &Face{
		Index:         index,
		VertexIndexes: [3]int{vertexIndex0, vertexIndex1, vertexIndex2},
	}
}

// 面リスト
type Faces struct {
	IndexModelCorrection[*Face]
	Data    map[int]*Face
	Indexes []int
}

func NewFaces() *Faces {
	return &Faces{
		IndexModelCorrection: *NewIndexModelCorrection[*Face](),
	}
}

func TestIndexModelCorrection_GetItem(t *testing.T) {
	model := NewFaces()
	item := NewFace(0, 0, 0, 0)
	model.SetItem(0, item)

	result := model.GetItem(0)
	if result != item {
		t.Errorf("Expected GetItem to return the item, but got %v", result)
	}

	// Test out of range index
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected GetItem to panic with out of range index")
		}
	}()
	model.GetItem(1)
}

func TestIndexModelCorrection_SetItem(t *testing.T) {
	model := NewFaces()
	item := NewFace(0, 0, 0, 0)
	model.SetItem(0, item)

	result := model.GetItem(0)
	if result != item {
		t.Errorf("Expected SetItem to set the item, but got %v", result)
	}
}

func TestIndexModelCorrection_Append(t *testing.T) {
	model := NewFaces()
	item := NewFace(0, 0, 0, 0)
	model.Append(item)

	result := model.GetItem(0)
	if result != item {
		t.Errorf("Expected Append to add the item, but got %v", result)
	}

	item2 := NewFace(1, 0, 0, 0)
	// Test sorting
	model.Append(item2)
	result = model.GetItem(0)
	if result == item {
		t.Errorf("Expected Append to sort the items, but got %v", result)
	}
	if result != item2 {
		t.Errorf("Expected Append to sort the items, but got %v", result)
	}
}

func TestIndexModelCorrection_DeleteItem(t *testing.T) {
	model := NewFaces()
	item := NewFace(0, 0, 0, 0)
	model.SetItem(0, item)

	model.DeleteItem(0)

	// Test deleted item
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected GetItem to panic with deleted item")
		}
	}()
	model.GetItem(0)
}

func TestIndexModelCorrection_Len(t *testing.T) {
	model := NewFaces()
	item := NewFace(0, 0, 0, 0)
	model.SetItem(0, item)

	result := model.Len()
	if result != 1 {
		t.Errorf("Expected Len to return 1, but got %d", result)
	}
}

func TestIndexModelCorrection_Contains(t *testing.T) {
	model := NewFaces()
	item := NewFace(0, 0, 0, 0)
	model.SetItem(0, item)

	result := model.Contains(0)
	if !result {
		t.Error("Expected Contains to return true, but got false")
	}

	result = model.Contains(1)
	if result {
		t.Error("Expected Contains to return false, but got true")
	}
}

func TestIndexModelCorrection_IsEmpty(t *testing.T) {
	model := NewFaces()

	result := model.IsEmpty()
	if !result {
		t.Error("Expected IsEmpty to return true, but got false")
	}

	item := NewFace(0, 0, 0, 0)
	model.SetItem(0, item)

	result = model.IsEmpty()
	if result {
		t.Error("Expected IsEmpty to return false, but got true")
	}
}

func TestIndexModelCorrection_LastIndex(t *testing.T) {
	model := NewFaces()
	item := NewFace(5, 0, 0, 0)
	model.SetItem(5, item)

	result := model.LastIndex()
	if result != 5 {
		t.Errorf("Expected LastIndex to return 5, but got %d", result)
	}
}
