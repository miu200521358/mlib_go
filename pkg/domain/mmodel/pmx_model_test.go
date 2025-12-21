package mmodel

import "testing"

func TestNewPmxModel(t *testing.T) {
	m := NewPmxModel("test.pmx")
	if m.Path() != "test.pmx" {
		t.Errorf("Path() = %v, want test.pmx", m.Path())
	}
	if m.Vertices == nil {
		t.Errorf("Vertices should not be nil")
	}
	if m.Bones == nil {
		t.Errorf("Bones should not be nil")
	}
	if m.DisplaySlots == nil {
		t.Errorf("DisplaySlots should not be nil")
	}
	// 初期DisplaySlotsはRoot + 表情の2つ
	if m.DisplaySlots.Length() != 2 {
		t.Errorf("DisplaySlots.Length() = %v, want 2", m.DisplaySlots.Length())
	}
}

func TestPmxModel_SettersGetters(t *testing.T) {
	m := NewPmxModel("")
	m.SetName("テストモデル")
	m.SetEnglishName("TestModel")
	m.SetPath("path/to/model.pmx")

	if m.Name() != "テストモデル" {
		t.Errorf("Name() = %v, want テストモデル", m.Name())
	}
	if m.EnglishName() != "TestModel" {
		t.Errorf("EnglishName() = %v, want TestModel", m.EnglishName())
	}
	if m.Path() != "path/to/model.pmx" {
		t.Errorf("Path() = %v, want path/to/model.pmx", m.Path())
	}
}

func TestPmxModel_UpdateHash(t *testing.T) {
	m := NewPmxModel("test.pmx")
	m.SetName("TestModel")
	m.UpdateHash()

	if m.Hash() == "" {
		t.Errorf("Hash() should not be empty after UpdateHash()")
	}

	// 同じ内容なら同じハッシュ
	m2 := NewPmxModel("test.pmx")
	m2.SetName("TestModel")
	m2.UpdateHash()

	if m.Hash() != m2.Hash() {
		t.Errorf("Same content should produce same hash")
	}
}

func TestPmxModel_SetRandHash(t *testing.T) {
	m := NewPmxModel("test.pmx")
	m.SetRandHash()

	if m.Hash() == "" {
		t.Errorf("Hash() should not be empty after SetRandHash()")
	}
}

func TestPmxModel_Copy(t *testing.T) {
	m := NewPmxModel("test.pmx")
	m.SetName("オリジナル")
	m.SetIndex(5)

	cp, err := m.Copy()
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	if cp.Name() != "オリジナル" {
		t.Errorf("Copy() Name = %v, want オリジナル", cp.Name())
	}
	if cp.Path() != "test.pmx" {
		t.Errorf("Copy() Path = %v, want test.pmx", cp.Path())
	}
}

func TestNewVertices(t *testing.T) {
	v := NewVertices(0)
	if v == nil {
		t.Errorf("NewVertices() should not return nil")
	}
	if v.Length() != 0 {
		t.Errorf("NewVertices(0).Length() = %v, want 0", v.Length())
	}
}

func TestNewBones(t *testing.T) {
	b := NewBones(0)
	if b == nil {
		t.Errorf("NewBones() should not return nil")
	}
	if b.Length() != 0 {
		t.Errorf("NewBones(0).Length() = %v, want 0", b.Length())
	}
}

func TestNewInitialDisplaySlots(t *testing.T) {
	slots := NewInitialDisplaySlots()
	if slots.Length() != 2 {
		t.Errorf("NewInitialDisplaySlots().Length() = %v, want 2", slots.Length())
	}

	root, _ := slots.Get(0)
	if root.Name() != "Root" {
		t.Errorf("First slot Name = %v, want Root", root.Name())
	}

	exp, _ := slots.Get(1)
	if exp.Name() != "表情" {
		t.Errorf("Second slot Name = %v, want 表情", exp.Name())
	}
}
