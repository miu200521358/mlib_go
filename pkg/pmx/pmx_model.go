package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
)

type PmxModel struct {
	*mcore.HashModel
	Signature           string
	Version             float64
	ExtendedUVCountType int
	VertexCountType     int
	TextureCountType    int
	MaterialCountType   int
	BoneCountType       int
	MorphCountType      int
	RigidBodyCountType  int
	Name                string
	EnglishName         string
	Comment             string
	EnglishComment      string
	JsonData            map[string]interface{}
	Vertices            *Vertices
	Faces               *Faces
	Textures            *Textures
	ToonTextures        *ToonTextures
	Materials           *Materials
	Bones               *Bones
	Morphs              *Morphs
	DisplaySlots        *DisplaySlots
	RigidBodies         *RigidBodies
	Joints              *Joints
	Meshes              *Meshes
}

func NewPmxModel(path string) *PmxModel {
	model := &PmxModel{}
	model.HashModel = mcore.NewHashModel(path)

	model.Vertices = NewVertices()
	model.Faces = NewFaces()
	model.Textures = NewTextures()
	model.ToonTextures = NewToonTextures()
	model.Materials = NewMaterials()
	model.Bones = NewBones()
	model.Morphs = NewMorphs()
	model.DisplaySlots = NewDisplaySlots()
	model.RigidBodies = NewRigidBodies()
	model.Joints = NewJoints()

	return model
}

func (pm *PmxModel) InitializeDisplaySlots() {
	d01 := NewDisplaySlot()
	d01.Name = "Root"
	d01.EnglishName = "Root"
	d01.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d01)

	d02 := NewDisplaySlot()
	d02.Name = "表情"
	d02.EnglishName = "Exp"
	d02.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d02)
}

func (pm *PmxModel) setup() {
	// セットアップ
	pm.Bones.setup()

	// 位置マッピングのセットアップ
	pm.Vertices.SetupMapKeys()
	pm.Bones.SetupMapKeys()
	pm.RigidBodies.SetupMapKeys()
	pm.Joints.SetupMapKeys()

	// 剛体
	for i := range pm.RigidBodies.Len() {
		rb := pm.RigidBodies.Get(i)
		if rb.BoneIndex >= 0 && pm.Bones.Contains(rb.BoneIndex) {
			// 剛体に関連付けられたボーンが存在する場合、剛体とボーンを関連付ける
			pm.Bones.Get(rb.BoneIndex).RigidBody = rb
			rb.Bone = pm.Bones.Get(rb.BoneIndex)
		}
	}
	// ジョイント
	for i := range pm.Joints.Len() {
		joint := pm.Joints.Get(i)
		if joint.RigidbodyIndexA >= 0 && pm.RigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && pm.RigidBodies.Contains(joint.RigidbodyIndexB) &&
			pm.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex >= 0 &&
			pm.Bones.Contains(pm.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex) &&
			pm.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex >= 0 &&
			pm.Bones.Contains(pm.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex) {
			// 剛体AもBも存在する場合、剛体Aと剛体Bを関連付ける
			pm.RigidBodies.Get(joint.RigidbodyIndexA).JointedBone =
				pm.Bones.Get(pm.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex)
			pm.RigidBodies.Get(joint.RigidbodyIndexB).JointedBone =
				pm.Bones.Get(pm.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex)
		}
	}
}

func (m *PmxModel) Copy() mcore.IHashModel {
	copied := NewPmxModel("")
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}
