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
	t.Run("基本コピー", func(t *testing.T) {
		tex := NewTextureByPath("test.png")
		tex.SetIndex(5)

		cp, err := tex.Copy()
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}
		if cp.Index() != 5 {
			t.Errorf("Copy() Index = %v, want 5", cp.Index())
		}
		if cp.Path != "test.png" {
			t.Errorf("Copy() Path = %v, want test.png", cp.Path)
		}
	})

	t.Run("別オブジェクト確認_ポインタ比較", func(t *testing.T) {
		tex := NewTextureByPath("test.png")

		cp, _ := tex.Copy()

		// コピー元とコピー先が異なるオブジェクトであることを確認
		if tex == cp {
			t.Errorf("Texture pointer should be different")
		}
	})

	t.Run("別オブジェクト確認_Path変更", func(t *testing.T) {
		tex := NewTextureByPath("original.png")

		cp, _ := tex.Copy()

		// Pathを変更してもコピー先に影響しないことを確認
		tex.Path = "changed.png"
		if cp.Path == "changed.png" {
			t.Errorf("Path should be independent")
		}
	})
}
