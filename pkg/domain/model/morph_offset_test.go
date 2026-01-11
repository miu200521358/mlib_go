package model

import "testing"

func TestMorphOffsetTypes(t *testing.T) {
	if (&VertexMorphOffset{}).MorphType() != MORPH_TYPE_VERTEX {
		t.Fatalf("VertexMorphOffset type mismatch")
	}
	if (&UvMorphOffset{UvType: MORPH_TYPE_EXTENDED_UV1}).MorphType() != MORPH_TYPE_EXTENDED_UV1 {
		t.Fatalf("UvMorphOffset type mismatch")
	}
	if (&BoneMorphOffset{}).MorphType() != MORPH_TYPE_BONE {
		t.Fatalf("BoneMorphOffset type mismatch")
	}
	if (&GroupMorphOffset{}).MorphType() != MORPH_TYPE_GROUP {
		t.Fatalf("GroupMorphOffset type mismatch")
	}
	if (&MaterialMorphOffset{}).MorphType() != MORPH_TYPE_MATERIAL {
		t.Fatalf("MaterialMorphOffset type mismatch")
	}
}
