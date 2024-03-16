package pmx

import (
	"embed"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
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
	Physics             *mphysics.MPhysics
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

func (pm *PmxModel) InitializeDraw(physics *mphysics.MPhysics, windowIndex int, resourceFiles embed.FS) {
	pm.Physics = physics
	pm.ToonTextures.InitGl(windowIndex, resourceFiles)
	pm.Meshes = NewMeshes(pm, windowIndex, resourceFiles)
	pm.RigidBodies.InitPhysics(pm.Physics, pm.Bones)
	pm.Joints.InitPhysics(pm.Physics, pm.RigidBodies)
	pm.Bones.PrepareDraw()
}

func (pm *PmxModel) Draw(
	shader *mgl.MShader,
	boneMatrixes []*mgl32.Mat4,
	boneGlobalMatrixes []*mmath.MMat4,
	boneTransforms []*mbt.BtTransform,
	windowIndex int,
	frame float32,
	elapsed float32,
	isBoneDebug bool,
) {
	pm.UpdatePhysics(boneMatrixes, boneTransforms, frame, elapsed)
	pm.Meshes.Draw(shader, boneMatrixes, windowIndex)

	// 物理デバッグ表示
	pm.Physics.DrawWorld()

	// ボーンデバッグ表示
	if isBoneDebug {
		pm.Bones.Draw(shader, boneGlobalMatrixes, windowIndex)
	}
}

func (pm *PmxModel) UpdatePhysics(
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
	frame float32,
	elapsed float32,
) {
	if pm.Physics == nil {
		return
	}

	for _, rigidBody := range pm.RigidBodies.GetSortedData() {
		rigidBody.UpdateTransform(boneMatrixes, boneTransforms, elapsed == 0.0)
	}

	if frame > pm.Physics.Spf {
		pm.Physics.Update(elapsed)

		// 剛体位置を更新
		for _, rigidBody := range pm.RigidBodies.GetSortedData() {
			rigidBody.UpdateMatrix(boneMatrixes, boneTransforms)
		}
	}
}

func (pm *PmxModel) SetUp() {
	// ボーン情報のセットアップ
	pm.Bones.setup()

	// ジョイント
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
