package mmodel

import "testing"

func TestNewMaterial(t *testing.T) {
	m := NewMaterial()
	if m.Index() != -1 {
		t.Errorf("Index() = %v, want -1", m.Index())
	}
	if m.Name() != "" {
		t.Errorf("Name() = %v, want empty", m.Name())
	}
	if m.Diffuse == nil || m.Specular == nil || m.Ambient == nil || m.Edge == nil {
		t.Errorf("color vectors should not be nil")
	}
	if m.DrawFlag != DRAW_FLAG_NONE {
		t.Errorf("DrawFlag = %v, want DRAW_FLAG_NONE", m.DrawFlag)
	}
}

func TestMaterial_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		m := NewMaterial()
		if m.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		m := NewMaterial()
		m.SetIndex(0)
		if !m.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestMaterial_Copy(t *testing.T) {
	t.Run("基本コピー", func(t *testing.T) {
		m := NewMaterial()
		m.SetIndex(5)
		m.SetName("テスト材質")
		m.Diffuse.X = 1.0

		cp, err := m.Copy()
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}
		if cp.Index() != 5 {
			t.Errorf("Copy() Index = %v, want 5", cp.Index())
		}
		if cp.Name() != "テスト材質" {
			t.Errorf("Copy() Name = %v, want テスト材質", cp.Name())
		}
	})

	t.Run("別オブジェクト確認_Diffuse", func(t *testing.T) {
		m := NewMaterial()
		m.Diffuse.X = 1.0

		cp, _ := m.Copy()

		if m.Diffuse == cp.Diffuse {
			t.Errorf("Diffuse pointer should be different")
		}

		m.Diffuse.X = 100
		if cp.Diffuse.X == 100 {
			t.Errorf("Diffuse should be independent")
		}
	})

	t.Run("別オブジェクト確認_ポインタ比較", func(t *testing.T) {
		m := NewMaterial()

		cp, _ := m.Copy()

		if m == cp {
			t.Errorf("Material pointer should be different")
		}
	})
}

func TestDrawFlag(t *testing.T) {
	t.Run("フラグ設定と取得", func(t *testing.T) {
		f := DRAW_FLAG_NONE

		f = f.SetDoubleSidedDrawing(true)
		if !f.IsDoubleSidedDrawing() {
			t.Errorf("IsDoubleSidedDrawing() = false, want true")
		}

		f = f.SetDrawingEdge(true)
		if !f.IsDrawingEdge() {
			t.Errorf("IsDrawingEdge() = false, want true")
		}

		f = f.SetDoubleSidedDrawing(false)
		if f.IsDoubleSidedDrawing() {
			t.Errorf("IsDoubleSidedDrawing() = true, want false")
		}
		if !f.IsDrawingEdge() {
			t.Errorf("IsDrawingEdge() should still be true")
		}
	})
}

func TestSphereMode(t *testing.T) {
	m := NewMaterial()
	m.SphereMode = SPHERE_MODE_MULTIPLICATION
	if m.SphereMode != SPHERE_MODE_MULTIPLICATION {
		t.Errorf("SphereMode = %v, want SPHERE_MODE_MULTIPLICATION", m.SphereMode)
	}
}

func TestToonSharing(t *testing.T) {
	m := NewMaterial()
	m.ToonSharingFlag = TOON_SHARING_SHARED
	if m.ToonSharingFlag != TOON_SHARING_SHARED {
		t.Errorf("ToonSharingFlag = %v, want TOON_SHARING_SHARED", m.ToonSharingFlag)
	}
}
