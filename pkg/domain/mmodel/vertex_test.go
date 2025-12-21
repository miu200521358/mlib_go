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
	t.Run("基本コピー", func(t *testing.T) {
		v := NewVertex()
		v.SetIndex(5)
		v.Position = mmath.NewVec3ByValues(1, 2, 3)
		v.Normal = mmath.NewVec3ByValues(0, 1, 0)
		v.Uv = mmath.NewVec2ByValues(0.5, 0.5)
		v.ExtendedUvs = append(v.ExtendedUvs, mmath.NewVec4ByValues(1, 2, 3, 4))
		v.EdgeFactor = 1.5
		v.MaterialIndexes = []int{0, 1}
		v.Deform = NewBdef1(10)

		cp, err := v.Copy()
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}
		if cp.Index() != 5 {
			t.Errorf("Copy() Index = %v, want 5", cp.Index())
		}
		if !cp.Position.NearEquals(v.Position, 1e-10) {
			t.Errorf("Copy() Position mismatch")
		}
	})

	t.Run("別オブジェクト確認_Position", func(t *testing.T) {
		v := NewVertex()
		v.Position = mmath.NewVec3ByValues(1, 2, 3)

		cp, _ := v.Copy()

		// ポインタが異なることを確認
		if v.Position == cp.Position {
			t.Errorf("Position pointer should be different")
		}

		// 値変更が影響しないことを確認
		v.Position.X = 100
		if cp.Position.X == 100 {
			t.Errorf("Position should be independent")
		}
	})

	t.Run("別オブジェクト確認_Normal", func(t *testing.T) {
		v := NewVertex()
		v.Normal = mmath.NewVec3ByValues(0, 1, 0)

		cp, _ := v.Copy()

		if v.Normal == cp.Normal {
			t.Errorf("Normal pointer should be different")
		}
	})

	t.Run("別オブジェクト確認_Uv", func(t *testing.T) {
		v := NewVertex()
		v.Uv = mmath.NewVec2ByValues(0.5, 0.5)

		cp, _ := v.Copy()

		if v.Uv == cp.Uv {
			t.Errorf("Uv pointer should be different")
		}
	})

	t.Run("別オブジェクト確認_ExtendedUvs", func(t *testing.T) {
		v := NewVertex()
		v.ExtendedUvs = append(v.ExtendedUvs, mmath.NewVec4ByValues(1, 2, 3, 4))

		cp, _ := v.Copy()

		if len(cp.ExtendedUvs) != 1 {
			t.Fatalf("ExtendedUvs length = %v, want 1", len(cp.ExtendedUvs))
		}
		if v.ExtendedUvs[0] == cp.ExtendedUvs[0] {
			t.Errorf("ExtendedUvs[0] pointer should be different")
		}
	})

	t.Run("別オブジェクト確認_MaterialIndexes", func(t *testing.T) {
		v := NewVertex()
		v.MaterialIndexes = []int{0, 1, 2}

		cp, _ := v.Copy()

		v.MaterialIndexes[0] = 999
		if cp.MaterialIndexes[0] == 999 {
			t.Errorf("MaterialIndexes should be independent")
		}
	})
}

func TestVertex_SetIndex(t *testing.T) {
	v := NewVertex()
	v.SetIndex(10)
	if v.Index() != 10 {
		t.Errorf("SetIndex() Index = %v, want 10", v.Index())
	}
}
