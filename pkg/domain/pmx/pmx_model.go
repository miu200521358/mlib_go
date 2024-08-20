package pmx

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

type PmxModel struct {
	index              int
	name               string
	path               string
	hash               string
	Signature          string
	Version            float64
	ExtendedUVCount    int
	VertexCountType    int
	TextureCountType   int
	MaterialCountType  int
	BoneCountType      int
	MorphCountType     int
	RigidBodyCountType int
	englishName        string
	Comment            string
	EnglishComment     string
	Vertices           *Vertices
	Faces              *Faces
	Textures           *Textures
	Materials          *Materials
	Bones              *Bones
	Morphs             *Morphs
	DisplaySlots       *DisplaySlots
	RigidBodies        *RigidBodies
	Joints             *Joints
}

func NewPmxModel(path string) *PmxModel {
	model := &PmxModel{}
	model.index = 0
	model.name = ""
	model.path = path
	model.hash = ""

	model.Vertices = NewVertices(0)
	model.Faces = NewFaces(0)
	model.Textures = NewTextures(0)
	model.Materials = NewMaterials(0)
	model.Bones = NewBones(0)
	model.Morphs = NewMorphs(0)
	model.DisplaySlots = NewDisplaySlots(0)
	model.RigidBodies = NewRigidBodies(0)
	model.Joints = NewJoints(0)

	return model
}

func (model *PmxModel) Index() int {
	return model.index
}

func (model *PmxModel) SetIndex(index int) {
	model.index = index
}

func (model *PmxModel) Path() string {
	return model.path
}

func (model *PmxModel) SetPath(path string) {
	model.path = path
}

func (model *PmxModel) Name() string {
	return model.name
}

func (model *PmxModel) SetName(name string) {
	model.name = name
}

func (model *PmxModel) Hash() string {
	return model.hash
}

func (model *PmxModel) SetHash(hash string) {
	model.hash = hash
}

func (model *PmxModel) SetRandHash() {
	model.hash = fmt.Sprintf("%d", rand.Intn(10000))
}

func (model *PmxModel) UpdateHash() {

	h := fnv.New32a()
	// 名前をハッシュに含める
	h.Write([]byte(model.Name()))
	// ファイルパスをハッシュに含める
	h.Write([]byte(model.Path()))
	// 各要素の数をハッシュに含める
	h.Write([]byte(fmt.Sprintf("%d", model.Vertices.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.Faces.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.Textures.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.Materials.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.Bones.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.Morphs.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.DisplaySlots.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.RigidBodies.Len())))
	h.Write([]byte(fmt.Sprintf("%d", model.Joints.Len())))

	// ハッシュ値を16進数文字列に変換
	model.SetHash(fmt.Sprintf("%x", h.Sum(nil)))
}

func (model *PmxModel) InitializeDisplaySlots() {
	d01 := NewDisplaySlot()
	d01.SetName("Root")
	d01.SetEnglishName("Root")
	d01.SpecialFlag = SPECIAL_FLAG_ON
	model.DisplaySlots.Update(d01)

	d02 := NewDisplaySlot()
	d02.SetName("表情")
	d02.SetEnglishName("Exp")
	d02.SpecialFlag = SPECIAL_FLAG_ON
	model.DisplaySlots.Update(d02)
}

func (model *PmxModel) Setup() {
	if !model.Materials.IsDirty() && !model.Bones.IsDirty() && !model.RigidBodies.IsDirty() && !model.Joints.IsDirty() {
		return
	}

	// セットアップ
	model.Materials.setup(model.Vertices, model.Faces, model.Textures)
	model.Bones.setup()

	// 剛体
	for i, rb := range model.RigidBodies.Data {
		if rb.BoneIndex >= 0 && model.Bones.Contains(rb.BoneIndex) {
			// 剛体に関連付けられたボーンが存在する場合、剛体とボーンを関連付ける
			model.Bones.Data[rb.BoneIndex].Extend.RigidBody = rb
			model.RigidBodies.Data[i].Bone = model.Bones.Get(rb.BoneIndex)
		}
	}

	// ジョイント
	for _, joint := range model.Joints.Data {
		if joint.RigidbodyIndexA >= 0 && model.RigidBodies.Contains(joint.RigidbodyIndexA) &&
			joint.RigidbodyIndexB >= 0 && model.RigidBodies.Contains(joint.RigidbodyIndexB) {

			// 剛体AもBも存在する場合、剛体Aと剛体Bを関連付ける
			if model.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex >= 0 &&
				model.Bones.Contains(model.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex) {
				model.RigidBodies.Data[joint.RigidbodyIndexA].JointedBone =
					model.Bones.Get(model.RigidBodies.Get(joint.RigidbodyIndexB).BoneIndex)
			}
			if model.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex >= 0 &&
				model.Bones.Contains(model.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex) {
				model.RigidBodies.Data[joint.RigidbodyIndexB].JointedBone =
					model.Bones.Get(model.RigidBodies.Get(joint.RigidbodyIndexA).BoneIndex)
			}
		}
	}

	model.Bones.SetDirty(false)
	model.RigidBodies.SetDirty(false)
	model.Joints.SetDirty(false)

	model.UpdateHash()
}

func (model *PmxModel) Copy() core.IHashModel {
	return &PmxModel{
		index:              model.index,
		name:               model.name,
		englishName:        model.englishName,
		path:               model.path,
		hash:               model.hash,
		Signature:          model.Signature,
		Version:            model.Version,
		ExtendedUVCount:    model.ExtendedUVCount,
		VertexCountType:    model.VertexCountType,
		TextureCountType:   model.TextureCountType,
		MaterialCountType:  model.MaterialCountType,
		BoneCountType:      model.BoneCountType,
		MorphCountType:     model.MorphCountType,
		RigidBodyCountType: model.RigidBodyCountType,
		Comment:            model.Comment,
		EnglishComment:     model.EnglishComment,
		Vertices:           model.Vertices.Copy(),
		Faces:              model.Faces.Copy(),
		Textures:           model.Textures.Copy(),
		Materials:          model.Materials.Copy(),
		Bones:              model.Bones.Copy(),
		Morphs:             model.Morphs.Copy(),
		DisplaySlots:       model.DisplaySlots.Copy(),
		RigidBodies:        model.RigidBodies.Copy(),
		Joints:             model.Joints.Copy(),
	}
}

func (model *PmxModel) EnglishName() string {
	return model.englishName
}

func (model *PmxModel) SetEnglishName(name string) {
	model.englishName = name
}
