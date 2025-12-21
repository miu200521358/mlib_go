package mmodel

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestMorphPanel_PanelName(t *testing.T) {
	tests := []struct {
		panel    MorphPanel
		expected string
	}{
		{MORPH_PANEL_SYSTEM, "システム"},
		{MORPH_PANEL_EYEBROW_LOWER_LEFT, "眉"},
		{MORPH_PANEL_EYE_UPPER_LEFT, "目"},
		{MORPH_PANEL_LIP_UPPER_RIGHT, "口"},
		{MORPH_PANEL_OTHER_LOWER_RIGHT, "他"},
	}

	for _, tt := range tests {
		if tt.panel.PanelName() != tt.expected {
			t.Errorf("PanelName() = %v, want %v", tt.panel.PanelName(), tt.expected)
		}
	}
}

func TestNewMorph(t *testing.T) {
	m := NewMorph()
	if m.Index() != -1 {
		t.Errorf("Index() = %v, want -1", m.Index())
	}
	if m.Panel != MORPH_PANEL_SYSTEM {
		t.Errorf("Panel = %v, want MORPH_PANEL_SYSTEM", m.Panel)
	}
	if m.MorphType != MORPH_TYPE_VERTEX {
		t.Errorf("MorphType = %v, want MORPH_TYPE_VERTEX", m.MorphType)
	}
	if len(m.Offsets) != 0 {
		t.Errorf("Offsets length = %v, want 0", len(m.Offsets))
	}
}

func TestMorph_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		m := NewMorph()
		if m.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		m := NewMorph()
		m.SetIndex(0)
		if !m.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestVertexMorphOffset(t *testing.T) {
	offset := NewVertexMorphOffset(10, mmath.NewVec3ByValues(1, 2, 3))
	if offset.VertexIndex != 10 {
		t.Errorf("VertexIndex = %v, want 10", offset.VertexIndex)
	}
	if offset.Type() != MORPH_TYPE_VERTEX {
		t.Errorf("Type() = %v, want MORPH_TYPE_VERTEX", offset.Type())
	}

	cp, err := offset.Copy()
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}
	if cp.(*VertexMorphOffset).VertexIndex != 10 {
		t.Errorf("Copy() VertexIndex = %v, want 10", cp.(*VertexMorphOffset).VertexIndex)
	}
}

func TestUvMorphOffset(t *testing.T) {
	offset := NewUvMorphOffset(5, mmath.NewVec4ByValues(0.1, 0.2, 0.3, 0.4))
	if offset.VertexIndex != 5 {
		t.Errorf("VertexIndex = %v, want 5", offset.VertexIndex)
	}
	if offset.Type() != MORPH_TYPE_UV {
		t.Errorf("Type() = %v, want MORPH_TYPE_UV", offset.Type())
	}
}

func TestBoneMorphOffset(t *testing.T) {
	offset := NewBoneMorphOffset(3)
	if offset.BoneIndex != 3 {
		t.Errorf("BoneIndex = %v, want 3", offset.BoneIndex)
	}
	if offset.Type() != MORPH_TYPE_BONE {
		t.Errorf("Type() = %v, want MORPH_TYPE_BONE", offset.Type())
	}
	if offset.Position == nil {
		t.Errorf("Position should not be nil")
	}
	if offset.Rotation == nil {
		t.Errorf("Rotation should not be nil")
	}
}

func TestGroupMorphOffset(t *testing.T) {
	offset := NewGroupMorphOffset(2, 0.5)
	if offset.MorphIndex != 2 {
		t.Errorf("MorphIndex = %v, want 2", offset.MorphIndex)
	}
	if offset.MorphFactor != 0.5 {
		t.Errorf("MorphFactor = %v, want 0.5", offset.MorphFactor)
	}
	if offset.Type() != MORPH_TYPE_GROUP {
		t.Errorf("Type() = %v, want MORPH_TYPE_GROUP", offset.Type())
	}
}

func TestMaterialMorphOffset(t *testing.T) {
	offset := NewMaterialMorphOffset(
		0,
		CALC_MODE_MULTIPLICATION,
		mmath.NewVec4ByValues(1, 1, 1, 1),
		mmath.NewVec4ByValues(1, 1, 1, 1),
		mmath.NewVec3ByValues(1, 1, 1),
		mmath.NewVec4ByValues(0, 0, 0, 1),
		1.0,
		mmath.NewVec4ByValues(1, 1, 1, 1),
		mmath.NewVec4ByValues(1, 1, 1, 1),
		mmath.NewVec4ByValues(1, 1, 1, 1),
	)
	if offset.MaterialIndex != 0 {
		t.Errorf("MaterialIndex = %v, want 0", offset.MaterialIndex)
	}
	if offset.CalcMode != CALC_MODE_MULTIPLICATION {
		t.Errorf("CalcMode = %v, want CALC_MODE_MULTIPLICATION", offset.CalcMode)
	}
	if offset.Type() != MORPH_TYPE_MATERIAL {
		t.Errorf("Type() = %v, want MORPH_TYPE_MATERIAL", offset.Type())
	}
}

func TestMorph_Copy(t *testing.T) {
	m := NewMorph()
	m.SetIndex(5)
	m.SetName("あ")
	m.Panel = MORPH_PANEL_LIP_UPPER_RIGHT
	m.MorphType = MORPH_TYPE_VERTEX
	m.Offsets = append(m.Offsets, NewVertexMorphOffset(0, mmath.NewVec3ByValues(1, 2, 3)))
	m.Offsets = append(m.Offsets, NewVertexMorphOffset(1, mmath.NewVec3ByValues(4, 5, 6)))

	cp, err := m.Copy()
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	if cp.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", cp.Index())
	}
	if cp.Name() != "あ" {
		t.Errorf("Copy() Name = %v, want あ", cp.Name())
	}
	if len(cp.Offsets) != 2 {
		t.Errorf("Copy() Offsets length = %v, want 2", len(cp.Offsets))
	}

	// 独立性確認
	m.Offsets[0].(*VertexMorphOffset).Position.X = 100
	if cp.Offsets[0].(*VertexMorphOffset).Position.X == 100 {
		t.Errorf("Offsets should be independent")
	}
}
