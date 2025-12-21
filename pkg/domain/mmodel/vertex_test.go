package mmodel

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestNewVertex(t *testing.T) {
	v := NewVertex()
	if v.Index() != -1 {
		t.Errorf("Index() = %v, want -1", v.Index())
	}
	if v.Position == nil || v.Normal == nil || v.Uv == nil {
		t.Errorf("vectors should not be nil")
	}
	if v.DeformType != DEFORM_BDEF1 {
		t.Errorf("DeformType = %v, want DEFORM_BDEF1", v.DeformType)
	}
}

func TestVertex_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		v := NewVertex()
		if v.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		v := NewVertex()
		v.SetIndex(0)
		if !v.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestVertex_Copy(t *testing.T) {
	v := NewVertex()
	v.SetIndex(5)
	v.Position = mmath.NewVec3ByValues(1, 2, 3)
	v.Normal = mmath.NewVec3ByValues(0, 1, 0)
	v.Uv = mmath.NewVec2ByValues(0.5, 0.5)
	v.ExtendedUvs = append(v.ExtendedUvs, mmath.NewVec4ByValues(1, 2, 3, 4))
	v.EdgeFactor = 1.5
	v.MaterialIndexes = []int{0, 1}
	v.Deform = NewBdef1(10)

	copied, err := v.Copy()
	if err != nil {
		t.Errorf("Copy() error = %v", err)
	}
	if copied.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", copied.Index())
	}
	if !copied.Position.NearEquals(v.Position, 1e-10) {
		t.Errorf("Copy() Position mismatch")
	}

	// 独立性確認
	v.Position.X = 100
	if copied.Position.X == 100 {
		t.Errorf("Copy() is not independent")
	}
}

func TestVertex_SetIndex(t *testing.T) {
	v := NewVertex()
	v.SetIndex(10)
	if v.Index() != 10 {
		t.Errorf("SetIndex() Index = %v, want 10", v.Index())
	}
}
