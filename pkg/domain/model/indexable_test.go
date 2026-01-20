// 指示: miu200521358
package model

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
)

func TestIsValidNil(t *testing.T) {
	var v *Vertex
	if v.IsValid() {
		t.Fatalf("Vertex nil should be invalid")
	}
	var f *Face
	if f.IsValid() {
		t.Fatalf("Face nil should be invalid")
	}
	var tex *Texture
	if tex.IsValid() {
		t.Fatalf("Texture nil should be invalid")
	}
	var mat *Material
	if mat.IsValid() {
		t.Fatalf("Material nil should be invalid")
	}
	var bone *Bone
	if bone.IsValid() {
		t.Fatalf("Bone nil should be invalid")
	}
	var morph *Morph
	if morph.IsValid() {
		t.Fatalf("Morph nil should be invalid")
	}
	var slot *DisplaySlot
	if slot.IsValid() {
		t.Fatalf("DisplaySlot nil should be invalid")
	}
	var rb *RigidBody
	if rb.IsValid() {
		t.Fatalf("RigidBody nil should be invalid")
	}
	var joint *Joint
	if joint.IsValid() {
		t.Fatalf("Joint nil should be invalid")
	}
}

func TestIndexableMethods(t *testing.T) {
	cases := []struct {
		name string
		new  func() collection.IIndexable
	}{
		{"Vertex", func() collection.IIndexable { return &Vertex{} }},
		{"Face", func() collection.IIndexable { return &Face{} }},
		{"Material", func() collection.IIndexable { return &Material{} }},
		{"Bone", func() collection.IIndexable { return &Bone{} }},
		{"Morph", func() collection.IIndexable { return &Morph{} }},
		{"DisplaySlot", func() collection.IIndexable { return &DisplaySlot{} }},
		{"RigidBody", func() collection.IIndexable { return &RigidBody{} }},
		{"Joint", func() collection.IIndexable { return &Joint{} }},
	}

	for _, tc := range cases {
		item := tc.new()
		item.SetIndex(-1)
		if item.IsValid() {
			t.Fatalf("%s should be invalid", tc.name)
		}
		item.SetIndex(0)
		if !item.IsValid() {
			t.Fatalf("%s should be valid", tc.name)
		}
		if item.Index() != 0 {
			t.Fatalf("%s Index=%d", tc.name, item.Index())
		}
	}
}

func TestNameableMethods(t *testing.T) {
	cases := []struct {
		name string
		new  func() collection.INameable
	}{
		{"Texture", func() collection.INameable { return NewTexture() }},
		{"Material", func() collection.INameable { return &Material{} }},
		{"Bone", func() collection.INameable { return &Bone{} }},
		{"Morph", func() collection.INameable { return &Morph{} }},
		{"DisplaySlot", func() collection.INameable { return &DisplaySlot{} }},
		{"RigidBody", func() collection.INameable { return &RigidBody{} }},
		{"Joint", func() collection.INameable { return &Joint{} }},
	}

	for _, tc := range cases {
		item := tc.new()
		item.SetName("name")
		if item.Name() != "name" {
			t.Fatalf("%s Name=%s", tc.name, item.Name())
		}
	}
}

func TestTextureValidity(t *testing.T) {
	tex := NewTexture()
	tex.SetIndex(0)
	if tex.IsValid() {
		t.Fatalf("texture should be invalid before SetValid")
	}
	tex.SetValid(true)
	if !tex.IsValid() {
		t.Fatalf("texture should be valid after SetValid")
	}
}
