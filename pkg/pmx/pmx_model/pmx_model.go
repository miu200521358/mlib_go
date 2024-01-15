package pmx_model

import (
	"github.com/miu200521358/mlib_go/pkg/core/hash_model"
	"github.com/miu200521358/mlib_go/pkg/pmx/bone"
	"github.com/miu200521358/mlib_go/pkg/pmx/display_slot"
	"github.com/miu200521358/mlib_go/pkg/pmx/face"
	"github.com/miu200521358/mlib_go/pkg/pmx/joint"
	"github.com/miu200521358/mlib_go/pkg/pmx/material"
	"github.com/miu200521358/mlib_go/pkg/pmx/morph"
	"github.com/miu200521358/mlib_go/pkg/pmx/rigidbody"
	"github.com/miu200521358/mlib_go/pkg/pmx/texture"
	"github.com/miu200521358/mlib_go/pkg/pmx/vertex"
)

type PmxModel struct {
	*hash_model.HashModel
	Signature           string
	Version             float64
	ExtendedUVCount     int
	VertexCount         int
	TextureCount        int
	MaterialCount       int
	BoneCount           int
	MorphCount          int
	RigidbodyCount      int
	Name                string
	EnglishName         string
	Comment             string
	EnglishComment      string
	JsonData            map[string]interface{}
	Vertices            vertex.Vertices
	Faces               face.Faces
	Textures            texture.Textures
	ToonTextures        texture.ToonTextures
	Materials           material.Materials
	Bones               bone.Bones
	Morphs              morph.Morphs
	DisplaySlots        display_slot.DisplaySlots
	RigidBodies         rigidbody.RigidBodies
	Joints              joint.Joints
	VerticesByBones     map[int][]int
	VerticesByMaterials map[int][]int
	FacesByMaterials    map[int][]int
}

func NewPmxModel(path string) *PmxModel {
	model := &PmxModel{}
	model.HashModel = hash_model.NewBaseHashModel(path)
	model.InitializeDisplaySlots()
	return model
}

func (pm *PmxModel) InitializeDisplaySlots() {
	d01 := display_slot.NewDisplaySlot("Root", "Root", display_slot.SPECIAL_FLAG_ON)
	pm.DisplaySlots.Append(d01, false)

	d02 := display_slot.NewDisplaySlot("表情", "Exp", display_slot.SPECIAL_FLAG_ON)
	pm.DisplaySlots.Append(d02, false)
}
