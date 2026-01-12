// 指示: miu200521358
package motion

import "testing"

// TestSortedFrameIndexStoreBasics は索引の基本動作を確認する。
func TestSortedFrameIndexStoreBasics(t *testing.T) {
	store := NewSortedFrameIndexStore()
	if store.Len() != 0 {
		t.Fatalf("Len: got=%d", store.Len())
	}
	if store.Has(1) {
		t.Fatalf("Has should be false")
	}

	store.Upsert(5)
	store.Upsert(1)
	store.Upsert(10)
	if !store.Has(1) || !store.Has(5) || !store.Has(10) {
		t.Fatalf("Has missing")
	}
	if !store.IsDirty() {
		t.Fatalf("IsDirty should be true")
	}
	store.Finalize()
	if store.IsDirty() {
		t.Fatalf("IsDirty should be false")
	}
	if store.Len() != 3 {
		t.Fatalf("Len after upsert: got=%d", store.Len())
	}

	min, ok := store.Min()
	if !ok || min != 1 {
		t.Fatalf("Min: got=%v ok=%v", min, ok)
	}
	max, ok := store.Max()
	if !ok || max != 10 {
		t.Fatalf("Max: got=%v ok=%v", max, ok)
	}

	prev, found := store.Prev(1)
	if prev != 1 || found {
		t.Fatalf("Prev min: got=%v found=%v", prev, found)
	}
	prev, found = store.Prev(6)
	if prev != 5 || !found {
		t.Fatalf("Prev mid: got=%v found=%v", prev, found)
	}

	next, found := store.Next(5)
	if next != 10 || !found {
		t.Fatalf("Next mid: got=%v found=%v", next, found)
	}
	next, found = store.Next(10)
	if next != 10 || found {
		t.Fatalf("Next max: got=%v found=%v", next, found)
	}

	seen := make([]Frame, 0)
	store.ForEach(func(frame Frame) bool {
		seen = append(seen, frame)
		return true
	})
	if len(seen) != 3 || seen[0] != 1 || seen[1] != 5 || seen[2] != 10 {
		t.Fatalf("ForEach order: got=%v", seen)
	}

	store.Delete(5)
	if store.Has(5) {
		t.Fatalf("Delete failed")
	}
	store.Finalize()
	if store.Len() != 2 {
		t.Fatalf("Len after delete: got=%d", store.Len())
	}
}

// TestSortedFrameIndexStoreEmpty は空ストアの境界を確認する。
func TestSortedFrameIndexStoreEmpty(t *testing.T) {
	store := NewSortedFrameIndexStore()
	if _, ok := store.Max(); ok {
		t.Fatalf("Max should be false")
	}
	if _, ok := store.Min(); ok {
		t.Fatalf("Min should be false")
	}
	if prev, ok := store.Prev(3); prev != 0 || ok {
		t.Fatalf("Prev empty: got=%v ok=%v", prev, ok)
	}
	if next, ok := store.Next(3); next != 3 || ok {
		t.Fatalf("Next empty: got=%v ok=%v", next, ok)
	}
}
