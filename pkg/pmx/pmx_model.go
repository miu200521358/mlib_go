package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

type PmxModel struct {
	*mcore.HashModel
	physics            *mphysics.MPhysics
	Signature          string
	Version            float64
	ExtendedUVCount    int
	VertexCountType    int
	TextureCountType   int
	MaterialCountType  int
	BoneCountType      int
	MorphCountType     int
	RigidBodyCountType int
	Name               string
	EnglishName        string
	Comment            string
	EnglishComment     string
	JsonData           map[string]interface{}
	Vertices           *Vertices
	Faces              *Faces
	Textures           *Textures
	ToonTextures       *ToonTextures
	Materials          *Materials
	Bones              *Bones
	Morphs             *Morphs
	DisplaySlots       *DisplaySlots
	RigidBodies        *RigidBodies
	Joints             *Joints
	Meshes             *Meshes
	DrawInitialized    bool
	Index              int // モデルインデックス(複数人モデル対応)
}

func NewPmxModel(path string) *PmxModel {
	model := &PmxModel{}
	model.HashModel = mcore.NewHashModel(path)

	model.Vertices = NewVertices(0)
	model.Faces = NewFaces(0)
	model.Textures = NewTextures(0)
	model.ToonTextures = NewToonTextures()
	model.Materials = NewMaterials(0)
	model.Bones = NewBones(0)
	model.Morphs = NewMorphs(0)
	model.DisplaySlots = NewDisplaySlots(0)
	model.RigidBodies = NewRigidBodies(0)
	model.Joints = NewJoints(0)
	model.DrawInitialized = false

	return model
}

func (pm *PmxModel) InitializeDisplaySlots() {
	d01 := NewDisplaySlot()
	d01.Name = "Root"
	d01.EnglishName = "Root"
	d01.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Update(d01)

	d02 := NewDisplaySlot()
	d02.Name = "表情"
	d02.EnglishName = "Exp"
	d02.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Update(d02)
}

func (pm *PmxModel) Setup() {
	if !pm.Materials.IsDirty() && !pm.Bones.IsDirty() && !pm.RigidBodies.IsDirty() && !pm.Joints.IsDirty() {
		return
	}

	// セットアップ
	pm.Materials.setup(pm.Vertices, pm.Faces)
	pm.Bones.setup()

	// 剛体
	for i, rb := range pm.RigidBodies.Data {
		if rb.BoneIndex >= 0 && pm.Bones.Contains(rb.BoneIndex) {
			// 剛体に関連付けられたボーンが存在する場合、剛体とボーンを関連付ける
			pm.Bones.Data[rb.BoneIndex].RigidBody = rb
			pm.RigidBodies.Data[i].Bone = pm.Bones.Get(rb.BoneIndex)
		}
	}

	// ジョイント
	for _, joint := range pm.Joints.Data {
		if joint.RigidbodyIndexA >= 0 && pm.RigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && pm.RigidBodies.Contains(joint.RigidbodyIndexB) &&
			pm.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex >= 0 &&
			pm.Bones.Contains(pm.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex) &&
			pm.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex >= 0 &&
			pm.Bones.Contains(pm.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex) {

			// 剛体AもBも存在する場合、剛体Aと剛体Bを関連付ける
			pm.RigidBodies.Data[joint.RigidbodyIndexA].JointedBone =
				pm.Bones.Get(pm.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex)
			pm.RigidBodies.Data[joint.RigidbodyIndexB].JointedBone =
				pm.Bones.Get(pm.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex)
		}
	}

	pm.Bones.SetDirty(false)
	pm.RigidBodies.SetDirty(false)
	pm.Joints.SetDirty(false)
}

func (m *PmxModel) Copy() mcore.IHashModel {
	copied := NewPmxModel("")
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}
