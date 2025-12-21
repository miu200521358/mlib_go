package mmodel

import "testing"

func TestNewTexture(t *testing.T) {
	tex := NewTexture()
	if tex.Index() != -1 {
		t.Errorf("Index() = %v, want -1", tex.Index())
	}
	if tex.Path != "" {
		t.Errorf("Path = %v, want empty", tex.Path)
	}
}

func TestNewTextureByPath(t *testing.T) {
	tex := NewTextureByPath("textures/body.png")
	if tex.Path != "textures/body.png" {
		t.Errorf("Path = %v, want textures/body.png", tex.Path)
	}
}

func TestTexture_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		tex := NewTexture()
		if tex.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		tex := NewTexture()
		tex.SetIndex(0)
		if !tex.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestTexture_Copy(t *testing.T) {
	tex := NewTextureByPath("test.png")
	tex.SetIndex(5)

	copied, err := tex.Copy()
	if err != nil {
		t.Errorf("Copy() error = %v", err)
	}
	if copied.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", copied.Index())
	}
	if copied.Path != "test.png" {
		t.Errorf("Copy() Path = %v, want test.png", copied.Path)
	}
}
