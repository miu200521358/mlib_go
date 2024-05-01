//go:build !for_linux
// +build !for_linux

package mwidget

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/deform"
	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

func Draw(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	shader *mgl.MShader,
	deltas *deform.VmdDeltas,
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

	updatePhysics(modelPhysics, model, boneMatrixes, boneTransforms, deltas, frame, elapsed, enablePhysics)
	model.Meshes.Draw(shader, boneMatrixes, vertexDeltas, materialDeltas, windowIndex)

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld()

	// ボーンデバッグ表示
	if isBoneDebug {
		model.Bones.Draw(shader, globalMatrixes, windowIndex)
	}
}

func updatePhysics(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
	deltas *deform.VmdDeltas,
	frame float32,
	elapsed float32,
	enablePhysics bool,
) {
	if modelPhysics == nil {
		return
	}

	for _, r := range model.RigidBodies.GetSortedData() {
		// 物理フラグが落ちている場合があるので、強制的に起こす
		forceUpdate := r.UpdateFlags(modelPhysics, enablePhysics)
		r.UpdateTransform(modelPhysics, boneTransforms, elapsed == 0.0 || !enablePhysics || forceUpdate)
	}

	if frame > modelPhysics.Spf {
		modelPhysics.Update(elapsed)

		// 剛体位置を更新
		for _, rigidBody := range model.RigidBodies.GetSortedData() {
			rigidBody.UpdateMatrix(modelPhysics, boneMatrixes, boneTransforms)
		}

		// // 物理後ボーン位置を更新
		// for boneIndex := range model.Bones.LayerSortedIndexes {
		// 	bone := model.Bones.GetItem(boneIndex)
		// 	if bone.IsAfterPhysicsDeform() && bone.ParentIndex == -1 && model.Bones.Contains(bone.ParentIndex) {
		// 		// 物理後ボーンで親が存在している場合、親の行列を取得する
		// 		parentMat := boneMatrixes[bone.ParentIndex]
		// 		pos := deltas.Bones.GetItem(bone.Name, frame).FramePosition.GL()
		// 		rot := deltas.Bones.GetItem(bone.Name, frame).FrameRotation.GL()
		// 		scl := deltas.Bones.GetItem(bone.Name, frame).FrameScale

		// 		// 自身の行列を作成
		// 		mat := parentMat.Mul4(mgl32.Translate3D(pos[0], pos[1], pos[2]))
		// 		mat = mat.Mul4(mgl32.HomogRotate3D(rot[3], mgl32.Vec3{rot[0], rot[1], rot[2]}))
		// 		mat = mat.Mul4(mgl32.Scale3D(float32(scl[0]), float32(scl[1]), float32(scl[2])))
		// 		boneMatrixes[boneIndex] = &mat
		// 	}
		// }
	}
}
