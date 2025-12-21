package mmodel

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/tiendc/go-deepcopy"
)

// コレクション型エイリアス
type Vertices = mcore.IndexModels[*Vertex]
type Faces = mcore.IndexModels[*Face]
type Textures = mcore.IndexModels[*Texture]
type Materials = mcore.IndexModels[*Material]
type Bones = mcore.IndexNameModels[*Bone]
type Morphs = mcore.IndexNameModels[*Morph]
type DisplaySlots = mcore.IndexNameModels[*DisplaySlot]
type RigidBodies = mcore.IndexNameModels[*RigidBody]
type Joints = mcore.IndexNameModels[*Joint]

// NewVertices は新しいVerticesコレクションを生成します。
func NewVertices(capacity int) *Vertices {
	v, _ := mcore.NewIndexModels[*Vertex](capacity)
	return v
}

// NewFaces は新しいFacesコレクションを生成します。
func NewFaces(capacity int) *Faces {
	v, _ := mcore.NewIndexModels[*Face](capacity)
	return v
}

// NewTextures は新しいTexturesコレクションを生成します。
func NewTextures(capacity int) *Textures {
	v, _ := mcore.NewIndexModels[*Texture](capacity)
	return v
}

// NewMaterials は新しいMaterialsコレクションを生成します。
func NewMaterials(capacity int) *Materials {
	v, _ := mcore.NewIndexModels[*Material](capacity)
	return v
}

// NewBones は新しいBonesコレクションを生成します。
func NewBones(capacity int) *Bones {
	v, _ := mcore.NewIndexNameModels[*Bone](capacity)
	return v
}

// NewMorphs は新しいMorphsコレクションを生成します。
func NewMorphs(capacity int) *Morphs {
	v, _ := mcore.NewIndexNameModels[*Morph](capacity)
	return v
}

// NewDisplaySlots は新しいDisplaySlotsコレクションを生成します。
func NewDisplaySlots(capacity int) *DisplaySlots {
	v, _ := mcore.NewIndexNameModels[*DisplaySlot](capacity)
	return v
}

// NewInitialDisplaySlots は初期表示枠（RootとExp）を含むコレクションを生成します。
func NewInitialDisplaySlots() *DisplaySlots {
	slots := NewDisplaySlots(0)
	slots.Append(NewRootDisplaySlot())
	slots.Append(NewMorphDisplaySlot())
	return slots
}

// NewRigidBodies は新しいRigidBodiesコレクションを生成します。
func NewRigidBodies(capacity int) *RigidBodies {
	v, _ := mcore.NewIndexNameModels[*RigidBody](capacity)
	return v
}

// NewJoints は新しいJointsコレクションを生成します。
func NewJoints(capacity int) *Joints {
	v, _ := mcore.NewIndexNameModels[*Joint](capacity)
	return v
}

// PmxModel はPMXモデルを表します。
type PmxModel struct {
	index              int
	name               string
	englishName        string
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

// NewPmxModel は新しいPmxModelを生成します。
func NewPmxModel(path string) *PmxModel {
	return &PmxModel{
		index:        0,
		name:         "",
		englishName:  "",
		path:         path,
		hash:         "",
		Vertices:     NewVertices(0),
		Faces:        NewFaces(0),
		Textures:     NewTextures(0),
		Materials:    NewMaterials(0),
		Bones:        NewBones(0),
		Morphs:       NewMorphs(0),
		DisplaySlots: NewInitialDisplaySlots(),
		RigidBodies:  NewRigidBodies(0),
		Joints:       NewJoints(0),
	}
}

// Index はインデックスを返します。
func (m *PmxModel) Index() int {
	return m.index
}

// SetIndex はインデックスを設定します。
func (m *PmxModel) SetIndex(index int) {
	m.index = index
}

// Name は名前を返します。
func (m *PmxModel) Name() string {
	return m.name
}

// SetName は名前を設定します。
func (m *PmxModel) SetName(name string) {
	m.name = name
}

// EnglishName は英語名を返します。
func (m *PmxModel) EnglishName() string {
	return m.englishName
}

// SetEnglishName は英語名を設定します。
func (m *PmxModel) SetEnglishName(name string) {
	m.englishName = name
}

// Path はファイルパスを返します。
func (m *PmxModel) Path() string {
	return m.path
}

// SetPath はファイルパスを設定します。
func (m *PmxModel) SetPath(path string) {
	m.path = path
}

// Hash はハッシュを返します。
func (m *PmxModel) Hash() string {
	return m.hash
}

// SetHash はハッシュを設定します。
func (m *PmxModel) SetHash(hash string) {
	m.hash = hash
}

// SetRandHash はランダムなハッシュを設定します。
func (m *PmxModel) SetRandHash() {
	m.hash = fmt.Sprintf("%d", rand.Intn(10000))
}

// UpdateHash はモデルの内容からハッシュを更新します。
func (m *PmxModel) UpdateHash() {
	h := fnv.New32a()
	h.Write([]byte(m.Name()))
	h.Write([]byte(m.Path()))
	h.Write([]byte(fmt.Sprintf("%d", m.Vertices.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.Faces.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.Textures.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.Materials.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.Bones.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.Morphs.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.DisplaySlots.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.RigidBodies.Length())))
	h.Write([]byte(fmt.Sprintf("%d", m.Joints.Length())))
	m.SetHash(fmt.Sprintf("%x", h.Sum(nil)))
}

// Copy は深いコピーを作成します。
func (m *PmxModel) Copy() (*PmxModel, error) {
	cp := &PmxModel{}
	if err := deepcopy.Copy(cp, m); err != nil {
		return nil, err
	}
	return cp, nil
}
