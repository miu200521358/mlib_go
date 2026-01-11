package model

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"testing"
)

func TestPmxModelHashAndCopy(t *testing.T) {
	m := NewPmxModel()
	m.SetName("model")
	m.SetPath("path")
	m.SetFileModTime(123)

	m.Vertices.Append(&Vertex{})
	m.Faces.Append(&Face{})
	m.Textures.Append(NewTexture())
	m.Materials.Append(NewMaterial())
	m.Bones.Append(&Bone{})
	m.Morphs.Append(&Morph{})
	m.DisplaySlots.Append(&DisplaySlot{})
	m.RigidBodies.Append(&RigidBody{})
	m.Joints.Append(&Joint{})

	parts := m.GetHashParts()
	expectedParts := strings.Repeat("00000001", 9)
	if parts != expectedParts {
		t.Fatalf("GetHashParts = %s", parts)
	}

	h := fnv.New32a()
	_, _ = h.Write([]byte(m.Name()))
	_, _ = h.Write([]byte(m.Path()))
	_, _ = h.Write([]byte(strconv.FormatInt(m.FileModTime(), 10)))
	_, _ = h.Write([]byte(parts))
	expectedHash := fmt.Sprintf("%x", h.Sum(nil))
	m.UpdateHash()
	if m.Hash() != expectedHash {
		t.Fatalf("UpdateHash = %s expected %s", m.Hash(), expectedHash)
	}

	m.SetHash("fixed")
	copied, err := m.Copy()
	if err != nil {
		t.Fatalf("Copy error: %v", err)
	}
	if copied.Hash() == "fixed" {
		t.Fatalf("Copy should update random hash")
	}
	if copied.Vertices == m.Vertices {
		t.Fatalf("Copy should deep copy collections")
	}

	m.Vertices.Append(&Vertex{})
	if copied.Vertices.Len() != 1 {
		t.Fatalf("Copy should not share collection")
	}
}
