package mmodel

import "testing"

func TestNewFace(t *testing.T) {
	f := NewFace()
	if f.Index() != -1 {
		t.Errorf("Index() = %v, want -1", f.Index())
	}
	for i, idx := range f.VertexIndexes {
		if idx != -1 {
			t.Errorf("VertexIndexes[%d] = %v, want -1", i, idx)
		}
	}
}

func TestNewFaceByIndexes(t *testing.T) {
	f := NewFaceByIndexes(0, 1, 2)
	if f.VertexIndexes[0] != 0 || f.VertexIndexes[1] != 1 || f.VertexIndexes[2] != 2 {
		t.Errorf("VertexIndexes = %v, want [0, 1, 2]", f.VertexIndexes)
	}
}

func TestFace_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		f := NewFace()
		if f.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後も頂点未設定で無効", func(t *testing.T) {
		f := NewFace()
		f.SetIndex(0)
		if f.IsValid() {
			t.Errorf("IsValid() = true, want false (vertex not set)")
		}
	})

	t.Run("全て設定後は有効", func(t *testing.T) {
		f := NewFaceByIndexes(0, 1, 2)
		f.SetIndex(0)
		if !f.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestFace_Copy(t *testing.T) {
	f := NewFaceByIndexes(10, 20, 30)
	f.SetIndex(5)

	copied, err := f.Copy()
	if err != nil {
		t.Errorf("Copy() error = %v", err)
	}
	if copied.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", copied.Index())
	}
	if copied.VertexIndexes != f.VertexIndexes {
		t.Errorf("Copy() VertexIndexes mismatch")
	}

	// 独立性確認
	f.VertexIndexes[0] = 999
	if copied.VertexIndexes[0] == 999 {
		t.Errorf("Copy() is not independent")
	}
}
