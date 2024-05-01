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

func (pm *PmxModel) SetUp() {
	// ボーン情報のセットアップ
	pm.Bones.setup()

	// 剛体・ジョイント
	for _, joint := range pm.Joints.GetSortedData() {
		if joint.RigidbodyIndexA >= 0 && pm.RigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && pm.RigidBodies.Contains(joint.RigidbodyIndexB) {
			// 剛体AもBも存在する場合、剛体Aと剛体Bを関連付ける
			pm.RigidBodies.GetItem(joint.RigidbodyIndexA).JointedBoneIndex =
				pm.RigidBodies.GetItem(joint.RigidbodyIndexB).BoneIndex
			pm.RigidBodies.GetItem(joint.RigidbodyIndexB).JointedBoneIndex =
				pm.RigidBodies.GetItem(joint.RigidbodyIndexA).BoneIndex
		}
	}
}
func (m *PmxModel) Copy() mcore.HashModelInterface {
	copied := NewPmxModel("")
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}
