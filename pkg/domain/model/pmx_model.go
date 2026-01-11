package model

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
	modelerrors "github.com/miu200521358/mlib_go/pkg/domain/model/errors"
	"github.com/miu200521358/mlib_go/pkg/shared/hashable"
	"github.com/tiendc/go-deepcopy"
)

// PmxModel は domain 層の PMX/X モデルを表す。
type PmxModel struct {
	*hashable.HashableBase
	EnglishName    string
	Comment        string
	EnglishComment string
	Vertices       *collection.IndexedCollection[*Vertex]
	Faces          *collection.IndexedCollection[*Face]
	Textures       *collection.NamedCollection[*Texture]
	Materials      *collection.NamedCollection[*Material]
	Bones          *BoneCollection
	Morphs         *collection.NamedCollection[*Morph]
	DisplaySlots   *collection.NamedCollection[*DisplaySlot]
	RigidBodies    *collection.NamedCollection[*RigidBody]
	Joints         *collection.NamedCollection[*Joint]
}

// NewPmxModel は PmxModel を生成する。
func NewPmxModel() *PmxModel {
	return &PmxModel{
		HashableBase: hashable.NewHashableBase("", ""),
		Vertices:     collection.NewIndexedCollection[*Vertex](0),
		Faces:        collection.NewIndexedCollection[*Face](0),
		Textures:     collection.NewNamedCollection[*Texture](0),
		Materials:    collection.NewNamedCollection[*Material](0),
		Bones:        NewBoneCollection(0),
		Morphs:       collection.NewNamedCollection[*Morph](0),
		DisplaySlots: collection.NewNamedCollection[*DisplaySlot](0),
		RigidBodies:  collection.NewNamedCollection[*RigidBody](0),
		Joints:       collection.NewNamedCollection[*Joint](0),
	}
}

// GetHashParts はハッシュ用の追加要素文字列を返す。
func (m *PmxModel) GetHashParts() string {
	vertices := 0
	faces := 0
	textures := 0
	materials := 0
	bones := 0
	morphs := 0
	displaySlots := 0
	rigidBodies := 0
	joints := 0
	if m.Vertices != nil {
		vertices = m.Vertices.Len()
	}
	if m.Faces != nil {
		faces = m.Faces.Len()
	}
	if m.Textures != nil {
		textures = m.Textures.Len()
	}
	if m.Materials != nil {
		materials = m.Materials.Len()
	}
	if m.Bones != nil {
		bones = m.Bones.Len()
	}
	if m.Morphs != nil {
		morphs = m.Morphs.Len()
	}
	if m.DisplaySlots != nil {
		displaySlots = m.DisplaySlots.Len()
	}
	if m.RigidBodies != nil {
		rigidBodies = m.RigidBodies.Len()
	}
	if m.Joints != nil {
		joints = m.Joints.Len()
	}
	return fmt.Sprintf(
		"%08d%08d%08d%08d%08d%08d%08d%08d%08d",
		vertices,
		faces,
		textures,
		materials,
		bones,
		morphs,
		displaySlots,
		rigidBodies,
		joints,
	)
}

// UpdateHash は name/path/modtime と追加要素でハッシュを更新する。
func (m *PmxModel) UpdateHash() {
	h := fnv.New32a()
	_, _ = h.Write([]byte(m.Name()))
	_, _ = h.Write([]byte(m.Path()))
	_, _ = h.Write([]byte(strconv.FormatInt(m.FileModTime(), 10)))
	_, _ = h.Write([]byte(m.GetHashParts()))
	m.SetHash(fmt.Sprintf("%x", h.Sum(nil)))
}

// Copy はモデルを複製し、ランダムハッシュを更新する。
func (m *PmxModel) Copy() (PmxModel, error) {
	var copied PmxModel
	if err := deepcopy.Copy(&copied, m); err != nil {
		return PmxModel{}, modelerrors.NewModelCopyFailed(err)
	}
	copied.UpdateRandomHash()
	return copied, nil
}
