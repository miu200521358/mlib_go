package mmodel

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestNewBone(t *testing.T) {
	b := NewBone()
	if b.Index() != -1 {
		t.Errorf("Index() = %v, want -1", b.Index())
	}
	if b.Name() != "" {
		t.Errorf("Name() = %v, want empty", b.Name())
	}
	if b.ParentIndex != -1 {
		t.Errorf("ParentIndex = %v, want -1", b.ParentIndex)
	}
	if b.Flag != BoneFlagNone {
		t.Errorf("Flag = %v, want BoneFlagNone", b.Flag)
	}
	if b.Position == nil {
		t.Errorf("Position should not be nil")
	}
}

func TestNewBoneByName(t *testing.T) {
	b := NewBoneByName("センター")
	if b.Name() != "センター" {
		t.Errorf("Name() = %v, want センター", b.Name())
	}
}

func TestBone_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		b := NewBone()
		if b.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		b := NewBone()
		b.SetIndex(0)
		if !b.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestBone_Direction(t *testing.T) {
	tests := []struct {
		name     string
		boneName string
		expected BoneDirection
	}{
		{"左腕は左", "左腕", BONE_DIRECTION_LEFT},
		{"右腕は右", "右腕", BONE_DIRECTION_RIGHT},
		{"センターは体幹", "センター", BONE_DIRECTION_TRUNK},
		{"上半身は体幹", "上半身", BONE_DIRECTION_TRUNK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoneByName(tt.boneName)
			if b.Direction() != tt.expected {
				t.Errorf("Direction() = %v, want %v", b.Direction(), tt.expected)
			}
		})
	}
}

func TestBone_Copy(t *testing.T) {
	t.Run("基本コピー", func(t *testing.T) {
		b := NewBoneByName("左腕")
		b.SetIndex(5)
		b.Position = mmath.NewVec3ByValues(1, 2, 3)
		b.ParentIndex = 4
		b.Layer = 1
		b.Flag = BoneFlagCanRotate | BoneFlagIsVisible

		cp, err := b.Copy()
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}
		if cp.Index() != 5 {
			t.Errorf("Copy() Index = %v, want 5", cp.Index())
		}
		if cp.Name() != "左腕" {
			t.Errorf("Copy() Name = %v, want 左腕", cp.Name())
		}
		if cp.ParentIndex != 4 {
			t.Errorf("Copy() ParentIndex = %v, want 4", cp.ParentIndex)
		}
	})

	t.Run("別オブジェクト確認_Position", func(t *testing.T) {
		b := NewBone()
		b.Position = mmath.NewVec3ByValues(1, 2, 3)

		cp, _ := b.Copy()

		b.Position.X = 100
		if cp.Position.X == 100 {
			t.Errorf("Position should be independent")
		}
	})

	t.Run("ParentBoneはnilになる", func(t *testing.T) {
		b := NewBone()
		b.ParentBone = NewBone()

		cp, _ := b.Copy()
		if cp.ParentBone != nil {
			t.Errorf("ParentBone should be nil after copy")
		}
	})
}

func TestBone_NormalizeLocalAxis(t *testing.T) {
	b := NewBone()
	b.NormalizeLocalAxis(mmath.NewVec3ByValues(1, 0, 0))

	if b.NormalizedLocalAxisX == nil {
		t.Fatalf("NormalizedLocalAxisX should not be nil")
	}
	if b.NormalizedLocalAxisY == nil {
		t.Fatalf("NormalizedLocalAxisY should not be nil")
	}
	if b.NormalizedLocalAxisZ == nil {
		t.Fatalf("NormalizedLocalAxisZ should not be nil")
	}

	// X軸正規化されている
	if !b.NormalizedLocalAxisX.NearEquals(mmath.VEC3_UNIT_X, 1e-10) {
		t.Errorf("NormalizedLocalAxisX = %v, want unit X", b.NormalizedLocalAxisX)
	}
}

func TestBone_NormalizeFixedAxis(t *testing.T) {
	b := NewBone()
	b.NormalizeFixedAxis(mmath.NewVec3ByValues(0, 2, 0))

	if b.NormalizedFixedAxis == nil {
		t.Fatalf("NormalizedFixedAxis should not be nil")
	}
	if !b.NormalizedFixedAxis.NearEquals(mmath.VEC3_UNIT_Y, 1e-10) {
		t.Errorf("NormalizedFixedAxis = %v, want unit Y", b.NormalizedFixedAxis)
	}
}
