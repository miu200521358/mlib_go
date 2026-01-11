package model

import (
	"testing"

	modelerrors "github.com/miu200521358/mlib_go/pkg/domain/model/errors"
)

func newBone(name string, layer int) *Bone {
	b := &Bone{Layer: layer}
	b.SetName(name)
	return b
}

func TestBoneCollectionAppendGetContains(t *testing.T) {
	bones := NewBoneCollection(0)
	b0 := newBone("a", 0)
	idx, res := bones.Append(b0)
	if idx != 0 || b0.Index() != 0 {
		t.Fatalf("Append index=%d bone.Index=%d", idx, b0.Index())
	}
	if res.Changed || len(res.Added) != 1 || res.Added[0] != 0 {
		t.Fatalf("Append result = %+v", res)
	}
	if !bones.Contains(0) {
		t.Fatalf("Contains should be true")
	}
	if _, err := bones.Get(0); err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if _, err := bones.Get(5); err == nil || !modelerrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Get out of range should return IndexOutOfRangeError")
	}
	if bones.Contains(5) {
		t.Fatalf("Contains out of range should be false")
	}

	delete(bones.indexToPos, 0)
	if _, err := bones.Get(0); err == nil || !modelerrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Get missing map should return IndexOutOfRangeError")
	}
	if bones.Contains(0) {
		t.Fatalf("Contains missing map should be false")
	}
	bones.rebuildNameIndex()
	bones.rebuildIndexToPos()
	if _, err := bones.Get(0); err != nil {
		t.Fatalf("Get after rebuild error: %v", err)
	}
}

func TestBoneCollectionInsertRoot(t *testing.T) {
	bones := NewBoneCollection(0)
	b0 := newBone("a", 0)
	b1 := newBone("b", 1)
	bones.Append(b0)
	bones.Append(b1)

	insert := newBone("root", 9)
	idx, res, err := bones.Insert(insert, -1)
	if err != nil {
		t.Fatalf("Insert error: %v", err)
	}
	if idx != 2 || insert.Index() != 2 {
		t.Fatalf("Insert index=%d insert.Index=%d", idx, insert.Index())
	}
	if insert.Layer != 0 {
		t.Fatalf("Insert layer=%d", insert.Layer)
	}
	if b0.Layer != 1 || b1.Layer != 2 {
		t.Fatalf("Existing layers = %d,%d", b0.Layer, b1.Layer)
	}
	if res.Changed {
		t.Fatalf("Insert result should not change indexes")
	}
}

func TestBoneCollectionInsertNotFound(t *testing.T) {
	bones := NewBoneCollection(0)
	b0 := newBone("a", 2)
	b1 := newBone("b", 5)
	bones.Append(b0)
	bones.Append(b1)

	insert := newBone("c", 0)
	_, _, err := bones.Insert(insert, 99)
	if err != nil {
		t.Fatalf("Insert error: %v", err)
	}
	if insert.Layer != b1.Layer {
		t.Fatalf("Insert layer=%d expected %d", insert.Layer, b1.Layer)
	}
	vals := bones.Values()
	if vals[len(vals)-1] != insert {
		t.Fatalf("Insert should append when not found")
	}
}

func TestBoneCollectionInsertFoundLayers(t *testing.T) {
	bones := NewBoneCollection(0)
	b0 := newBone("a", 0)
	b1 := newBone("b", 1)
	bones.Append(b0)
	bones.Append(b1)

	insertNoGap := newBone("c", 9)
	_, _, err := bones.Insert(insertNoGap, 0)
	if err != nil {
		t.Fatalf("Insert error: %v", err)
	}
	if insertNoGap.Layer != 0 {
		t.Fatalf("Insert no gap layer=%d", insertNoGap.Layer)
	}
	vals := bones.Values()
	if vals[1] != insertNoGap {
		t.Fatalf("Insert position mismatch")
	}

	bonesGap := NewBoneCollection(0)
	g0 := newBone("a", 0)
	g1 := newBone("b", 3)
	bonesGap.Append(g0)
	bonesGap.Append(g1)
	insertGap := newBone("c", 9)
	_, _, err = bonesGap.Insert(insertGap, 0)
	if err != nil {
		t.Fatalf("Insert gap error: %v", err)
	}
	if insertGap.Layer != 1 {
		t.Fatalf("Insert gap layer=%d", insertGap.Layer)
	}

	insertEnd := newBone("end", 0)
	_, _, err = bonesGap.Insert(insertEnd, g1.Index())
	if err != nil {
		t.Fatalf("Insert end error: %v", err)
	}
	if insertEnd.Layer != g1.Layer {
		t.Fatalf("Insert end layer=%d", insertEnd.Layer)
	}
}

func TestBoneCollectionRemoveUpdateRename(t *testing.T) {
	bones := NewBoneCollection(0)
	b0 := newBone("a", 0)
	b1 := newBone("a", 0)
	bones.Append(b0)
	bones.Append(b1)

	res, err := bones.Remove(0)
	if err != nil {
		t.Fatalf("Remove error: %v", err)
	}
	if !res.Changed || len(res.Removed) != 1 || res.Removed[0] != 0 {
		t.Fatalf("Remove result = %+v", res)
	}
	if b1.Index() != 0 {
		t.Fatalf("Remove should reindex remaining bones")
	}

	if _, err := bones.Update(0, newBone("x", 0)); err == nil || !modelerrors.IsNameMismatchError(err) {
		t.Fatalf("Update mismatch should return NameMismatchError")
	}

	upd := newBone("a", 1)
	if _, err := bones.Update(0, upd); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if upd.Index() != 0 {
		t.Fatalf("Update should keep index")
	}

	bones.Append(newBone("b", 0))
	if _, err := bones.Rename(0, "b"); err == nil || !modelerrors.IsNameConflictError(err) {
		t.Fatalf("Rename conflict should return NameConflictError")
	}

	changed, err := bones.Rename(0, "c")
	if err != nil || !changed {
		t.Fatalf("Rename should succeed")
	}

	if _, err := bones.GetByName("a"); err == nil || !modelerrors.IsNameNotFoundError(err) {
		t.Fatalf("GetByName should be updated after rename")
	}
}

func TestBoneCollectionRemoveInvalidMap(t *testing.T) {
	bones := NewBoneCollection(0)
	bones.Append(newBone("a", 0))
	delete(bones.indexToPos, 0)
	if _, err := bones.Remove(0); err == nil || !modelerrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Remove missing map should return IndexOutOfRangeError")
	}
}
