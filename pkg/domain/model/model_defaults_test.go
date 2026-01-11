package model

import "testing"

func TestMaterialDefaults(t *testing.T) {
	mat := NewMaterial()
	if mat.Index() != -1 {
		t.Fatalf("Index=%d", mat.Index())
	}
	if mat.Name() != "" || mat.EnglishName != "" || mat.Memo != "" {
		t.Fatalf("Name/EnglishName/Memo defaults mismatch")
	}
	if mat.DrawFlag != DRAW_FLAG_NONE {
		t.Fatalf("DrawFlag=%v", mat.DrawFlag)
	}
	if mat.EdgeSize != 0.0 {
		t.Fatalf("EdgeSize=%v", mat.EdgeSize)
	}
	if mat.TextureIndex != -1 || mat.SphereTextureIndex != -1 || mat.ToonTextureIndex != -1 {
		t.Fatalf("Texture indexes defaults mismatch")
	}
	if mat.SphereMode != SPHERE_MODE_INVALID {
		t.Fatalf("SphereMode=%v", mat.SphereMode)
	}
	if mat.ToonSharingFlag != TOON_SHARING_INDIVIDUAL {
		t.Fatalf("ToonSharingFlag=%v", mat.ToonSharingFlag)
	}
	if mat.VerticesCount != 0 {
		t.Fatalf("VerticesCount=%d", mat.VerticesCount)
	}
}

func TestTextureDefaults(t *testing.T) {
	tex := NewTexture()
	if tex.Index() != -1 {
		t.Fatalf("Index=%d", tex.Index())
	}
	if tex.Name() != "" || tex.EnglishName != "" {
		t.Fatalf("Name/EnglishName defaults mismatch")
	}
	if tex.TextureType != TEXTURE_TYPE_TEXTURE {
		t.Fatalf("TextureType=%v", tex.TextureType)
	}
	if tex.IsValid() {
		t.Fatalf("IsValid should be false by default")
	}
}
