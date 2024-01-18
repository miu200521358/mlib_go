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
	RigidBodyCount      int
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
	model.HashModel = hash_model.NewHashModel(path)

	model.Vertices = *vertex.NewVertices()
	model.Faces = *face.NewFaces()
	model.Textures = *texture.NewTextures()
	model.ToonTextures = *texture.NewToonTextures()
	model.Materials = *material.NewMaterials()
	model.Bones = *bone.NewBones()
	model.Morphs = *morph.NewMorphs()
	model.DisplaySlots = *display_slot.NewDisplaySlots()
	model.RigidBodies = *rigidbody.NewRigidBodies()
	model.Joints = *joint.NewJoints()
	model.VerticesByBones = make(map[int][]int)
	model.VerticesByMaterials = make(map[int][]int)
	model.FacesByMaterials = make(map[int][]int)

	return model
}

func (pm *PmxModel) InitializeDisplaySlots() {
	d01 := display_slot.NewDisplaySlot()
	d01.Name = "Root"
	d01.EnglishName = "Root"
	d01.SpecialFlag = display_slot.SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d01, false)

	d02 := display_slot.NewDisplaySlot()
	d02.Name = "表情"
	d02.EnglishName = "Exp"
	d02.SpecialFlag = display_slot.SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d02, false)
}
