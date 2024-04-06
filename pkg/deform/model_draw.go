package deform

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

func Draw(
	model *pmx.PmxModel,
	shader *mgl.MShader,
	deltas *VmdDeltas,
	windowIndex int,
	frame float32,
	elapsed float32,
	isBoneDebug bool,
	enablePhysics bool,
) {
	boneMatrixes := make([]*mgl32.Mat4, len(model.Bones.NameIndexes))
	globalMatrixes := make([]*mmath.MMat4, len(model.Bones.NameIndexes))
	boneTransforms := make([]*mbt.BtTransform, len(model.Bones.NameIndexes))
	vertexDeltas := make([][]float32, len(model.Vertices.Data))
	materialDeltas := make([]*pmx.Material, len(model.Materials.Data))
	for i, bone := range model.Bones.GetSortedData() {
		mat := deltas.Bones.GetItem(bone.Name, frame).LocalMatrix.GL()
		boneMatrixes[i] = mat
		globalMatrixes[i] = deltas.Bones.GetItem(bone.Name, frame).GlobalMatrix
		t := mbt.NewBtTransform()
		t.SetFromOpenGLMatrix(&mat[0])
		boneTransforms[i] = &t
	}
	// TODO: 並列化
	for i, vd := range deltas.Morphs.Vertices.Data {
		vertexDeltas[i] = vd.GL()
	}
	for i, md := range deltas.Morphs.Materials.Data {
		materialDeltas[i] = md.Material
	}

	updatePhysics(model, boneMatrixes, boneTransforms, frame, elapsed, enablePhysics)
	model.Meshes.Draw(shader, boneMatrixes, vertexDeltas, materialDeltas, windowIndex)

	// 物理デバッグ表示
	model.Physics.DebugDrawWorld()

	// ボーンデバッグ表示
	if isBoneDebug {
		model.Bones.Draw(shader, globalMatrixes, windowIndex)
	}
}

func updatePhysics(
	model *pmx.PmxModel,
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
	frame float32,
	elapsed float32,
	enablePhysics bool,
) {
	if model.Physics == nil {
		return
	}

	for _, r := range model.RigidBodies.GetSortedData() {
		// 物理フラグが落ちている場合があるので、強制的に起こす
		forceUpdate := r.UpdateFlags(enablePhysics)
		r.UpdateTransform(boneTransforms, elapsed == 0.0 || !enablePhysics || forceUpdate)
	}

	if frame > model.Physics.Spf {
		model.Physics.Update(elapsed)

		// 剛体位置を更新
		for _, rigidBody := range model.RigidBodies.GetSortedData() {
			rigidBody.UpdateMatrix(boneMatrixes, boneTransforms)
		}
	}
}
