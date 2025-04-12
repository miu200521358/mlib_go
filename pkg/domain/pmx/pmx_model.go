package pmx

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/tiendc/go-deepcopy"
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
	model.DisplaySlots = NewInitialDisplaySlots()
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

func (model *PmxModel) EnglishName() string {
	return model.englishName
}

func (model *PmxModel) SetEnglishName(name string) {
	model.englishName = name
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
	h.Write([]byte(fmt.Sprintf("%d", model.Vertices.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.Faces.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.Textures.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.Materials.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.Bones.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.Morphs.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.DisplaySlots.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.RigidBodies.Length())))
	h.Write([]byte(fmt.Sprintf("%d", model.Joints.Length())))

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
	// セットアップ
	model.Materials.Setup(model.Vertices, model.Faces, model.Textures)
	model.Bones.Setup()

	// 剛体
	model.RigidBodies.ForEach(func(index int, rb *RigidBody) bool {
		if rb.BoneIndex >= 0 && model.Bones.Contains(rb.BoneIndex) {
			// 剛体に関連付けられたボーンが存在する場合、剛体とボーンを関連付ける
			if bone, err := model.Bones.Get(rb.BoneIndex); err == nil {
				bone.RigidBody = rb
				rb.Bone = bone
			}
		}
		return true
	})

	model.UpdateHash()
}

func (model *PmxModel) Copy() (*PmxModel, error) {
	copied := new(PmxModel)
	err := deepcopy.Copy(copied, model)
	return copied, err
}
