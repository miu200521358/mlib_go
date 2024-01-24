package pmx

import "github.com/miu200521358/mlib_go/pkg/mcore"

type PmxModel struct {
	*mcore.HashModel
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
	Vertices            Vertices
	Faces               Faces
	Textures            Textures
	ToonTextures        ToonTextures
	Materials           Materials
	Bones               Bones
	Morphs              Morphs
	DisplaySlots        DisplaySlots
	RigidBodies         RigidBodies
	Joints              Joints
	VerticesByBones     map[int][]int
	VerticesByMaterials map[int][]int
	FacesByMaterials    map[int][]int
	Meshes              *Meshes
}

func NewPmxModel(path string) *PmxModel {
	model := &PmxModel{}
	model.HashModel = mcore.NewHashModel(path)

	model.Vertices = *NewVertices()
	model.Faces = *NewFaces()
	model.Textures = *NewTextures()
	model.ToonTextures = *NewToonTextures()
	model.Materials = *NewMaterials()
	model.Bones = *NewBones()
	model.Morphs = *NewMorphs()
	model.DisplaySlots = *NewDisplaySlots()
	model.RigidBodies = *NewRigidBodies()
	model.Joints = *NewJoints()
	model.VerticesByBones = make(map[int][]int)
	model.VerticesByMaterials = make(map[int][]int)
	model.FacesByMaterials = make(map[int][]int)

	return model
}

func (pm *PmxModel) InitializeDisplaySlots() {
	d01 := NewDisplaySlot()
	d01.Name = "Root"
	d01.EnglishName = "Root"
	d01.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d01, false)

	d02 := NewDisplaySlot()
	d02.Name = "表情"
	d02.EnglishName = "Exp"
	d02.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d02, false)
}

func (pm *PmxModel) InitializeDraw() {
}
