// 指示: miu200521358
package collection

import (
	"sort"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model/merrors"
)

type testItem struct {
	index int
	name  string
	valid bool
}

func (t *testItem) Index() int {
	return t.index
}

func (t *testItem) SetIndex(index int) {
	t.index = index
}

func (t *testItem) IsValid() bool {
	return t != nil && t.index >= 0 && t.valid
}

func (t *testItem) Name() string {
	return t.name
}

func (t *testItem) SetName(name string) {
	t.name = name
}

func newItem(name string, index int, valid bool) *testItem {
	return &testItem{name: name, index: index, valid: valid}
}

func TestNameIndexRebuildAndNames(t *testing.T) {
	idx := NewNameIndex[*testItem]()
	items := []*testItem{
		newItem("a", 0, true),
		newItem("a", 1, true),
		newItem("b", 2, false),
		newItem("c", 3, true),
	}
	idx.Rebuild(items)

	if got, ok := idx.GetByName("a"); !ok || got != 0 {
		t.Fatalf("GetByName(a) = %v, %v", got, ok)
	}
	if _, ok := idx.GetByName("b"); ok {
		t.Fatalf("GetByName(b) should be missing")
	}
	if got, ok := idx.GetByName("c"); !ok || got != 3 {
		t.Fatalf("GetByName(c) = %v, %v", got, ok)
	}

	names := idx.Names()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "a" || names[1] != "c" {
		t.Fatalf("Names = %v", names)
	}

	if idx.SetIfAbsent("a", 99) {
		t.Fatalf("SetIfAbsent should not overwrite existing name")
	}
}

func TestNameIndexNil(t *testing.T) {
	var idx *NameIndex[*testItem]
	idx.Rebuild([]*testItem{newItem("a", 0, true)})
	if idx.SetIfAbsent("a", 0) {
		t.Fatalf("SetIfAbsent on nil should be false")
	}
	if _, ok := idx.GetByName("a"); ok {
		t.Fatalf("GetByName on nil should be false")
	}
	if names := idx.Names(); names != nil {
		t.Fatalf("Names on nil should be nil")
	}
}

func TestIndexedCollectionAppendInsertRemoveUpdateContainsGet(t *testing.T) {
	c := NewIndexedCollection[*testItem](0)

	item := newItem("x", 5, true)
	idx, res := c.Append(item)
	if idx != 0 || item.Index() != 0 {
		t.Fatalf("Append index = %d item.Index=%d", idx, item.Index())
	}
	if res.Changed || len(res.Added) != 1 || res.Added[0] != 0 {
		t.Fatalf("Append result = %+v", res)
	}
	if res.OldToNew == nil || res.NewToOld == nil {
		t.Fatalf("Append mappings should not be nil")
	}

	c.Append(newItem("y", 10, true))
	c.Append(newItem("z", 11, true))

	inserted := newItem("ins", 99, true)
	newIdx, insRes, err := c.Insert(inserted, 1)
	if err != nil {
		t.Fatalf("Insert error: %v", err)
	}
	if newIdx != 1 || inserted.Index() != 1 {
		t.Fatalf("Insert index = %d inserted.Index=%d", newIdx, inserted.Index())
	}
	if !insRes.Changed || len(insRes.Added) != 1 || insRes.Added[0] != 1 {
		t.Fatalf("Insert result = %+v", insRes)
	}
	if got := insRes.OldToNew; len(got) != 3 || got[0] != 0 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("Insert OldToNew = %v", got)
	}
	if got := insRes.NewToOld; len(got) != 3 || got[0] != 0 || got[1] != -1 || got[2] != 1 {
		t.Fatalf("Insert NewToOld = %v", got)
	}

	appendItem := newItem("append", 0, true)
	appendIdx, appendRes, err := c.Insert(appendItem, c.Len())
	if err != nil {
		t.Fatalf("Insert append error: %v", err)
	}
	if appendIdx != c.Len()-1 || appendRes.Changed {
		t.Fatalf("Insert append result = %+v", appendRes)
	}

	_, _, err = c.Insert(newItem("bad", 0, true), -1)
	if err == nil || !merrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Insert negative should return IndexOutOfRangeError")
	}
	_, _, err = c.Insert(newItem("bad", 0, true), c.Len()+1)
	if err == nil || !merrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Insert overflow should return IndexOutOfRangeError")
	}

	remRes, err := c.Remove(1)
	if err != nil {
		t.Fatalf("Remove error: %v", err)
	}
	if !remRes.Changed || len(remRes.Removed) != 1 || remRes.Removed[0] != 1 {
		t.Fatalf("Remove result = %+v", remRes)
	}
	if got := remRes.OldToNew; len(got) != 5 || got[0] != 0 || got[1] != -1 || got[2] != 1 || got[3] != 2 || got[4] != 3 {
		t.Fatalf("Remove OldToNew = %v", got)
	}
	if got := remRes.NewToOld; len(got) != 5 || got[0] != 0 || got[1] != 2 || got[2] != 3 || got[3] != 4 || got[4] != -1 {
		t.Fatalf("Remove NewToOld = %v", got)
	}

	if _, err := c.Remove(999); err == nil || !merrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Remove out of range should return IndexOutOfRangeError")
	}

	upd := newItem("upd", 10, true)
	updRes, err := c.Update(0, upd)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if updRes.Changed || upd.Index() != 0 {
		t.Fatalf("Update result = %+v index=%d", updRes, upd.Index())
	}
	if _, err := c.Update(999, upd); err == nil || !merrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Update out of range should return IndexOutOfRangeError")
	}

	invalid := newItem("inv", 0, false)
	c.Append(invalid)
	if c.Contains(c.Len() - 1) {
		t.Fatalf("Contains should be false for invalid item")
	}
	if c.Contains(999) {
		t.Fatalf("Contains should be false for out of range")
	}

	if _, err := c.Get(999); err == nil || !merrors.IsIndexOutOfRangeError(err) {
		t.Fatalf("Get out of range should return IndexOutOfRangeError")
	}
}

func TestNamedCollectionBehavior(t *testing.T) {
	c := NewNamedCollection[*testItem](0)

	item0 := newItem("a", 0, true)
	c.Append(item0)
	item1 := newItem("a", 1, true)
	c.Append(item1)

	got, err := c.GetByName("a")
	if err != nil || got != item0 {
		t.Fatalf("GetByName should return first win")
	}

	inserted := newItem("a", 0, true)
	_, _, err = c.Insert(inserted, 0)
	if err != nil {
		t.Fatalf("Insert error: %v", err)
	}
	got, err = c.GetByName("a")
	if err != nil || got != inserted {
		t.Fatalf("GetByName after insert should return inserted")
	}

	_, err = c.Remove(0)
	if err != nil {
		t.Fatalf("Remove error: %v", err)
	}
	got, err = c.GetByName("a")
	if err != nil || got != item0 {
		t.Fatalf("GetByName after remove should return next")
	}

	if _, err := c.Update(0, newItem("x", 0, true)); err == nil || !merrors.IsNameMismatchError(err) {
		t.Fatalf("Update mismatch should return NameMismatchError")
	}

	upd := newItem("a", 0, true)
	if _, err := c.Update(0, upd); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	changed, err := c.Rename(0, "a")
	if err != nil || changed {
		t.Fatalf("Rename same name should be no-op")
	}

	if _, err := c.Rename(0, "a2"); err != nil {
		t.Fatalf("Rename error: %v", err)
	}
	if _, err := c.Rename(0, "a2"); err != nil {
		t.Fatalf("Rename idempotent error: %v", err)
	}

	if _, err := c.Rename(0, "a"); err == nil || !merrors.IsNameConflictError(err) {
		t.Fatalf("Rename back should return NameConflictError")
	}

	dup := newItem("dup", 0, true)
	c.Append(dup)
	if _, err := c.Rename(0, "dup"); err == nil || !merrors.IsNameConflictError(err) {
		t.Fatalf("Rename conflict should return NameConflictError")
	}

	if _, err := c.GetByName("missing"); err == nil || !merrors.IsNameNotFoundError(err) {
		t.Fatalf("GetByName missing should return NameNotFoundError")
	}
}
