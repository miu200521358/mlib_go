//go:build windows
// +build windows

package mwidget

import (
	"math"
	"runtime"
	"sync"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
	"github.com/miu200521358/mlib_go/pkg/mview"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

func Draw(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	shader *mview.MShader,
	deltas *vmd.VmdDeltas,
	windowIndex int,
	fno int,
	elapsed float32,
	isBoneDebug bool,
	enablePhysics bool,
) {
	boneMatrixes := make([]*mgl32.Mat4, len(model.Bones.NameIndexes))
	globalMatrixes := make([]*mmath.MMat4, len(model.Bones.NameIndexes))
	boneTransforms := make([]*mbt.BtTransform, len(model.Bones.NameIndexes))
	materialDeltas := make([]*pmx.Material, len(model.Materials.Data))
	for i, bone := range model.Bones.GetSortedData() {
		mat := deltas.Bones.GetItem(bone.Name, fno).LocalMatrix.GL()
		boneMatrixes[i] = mat
		globalMatrixes[i] = &deltas.Bones.GetItem(bone.Name, fno).GlobalMatrix
		t := mbt.NewBtTransform()
		t.SetFromOpenGLMatrix(&mat[0])
		boneTransforms[i] = &t
	}

	for i, md := range deltas.Morphs.Materials.Data {
		materialDeltas[i] = md.Material
	}

	vertexDeltas := fetchVertexDeltas(model, deltas)

	updatePhysics(modelPhysics, model, boneMatrixes, boneTransforms, deltas, fno, elapsed, enablePhysics)
	model.Meshes.Draw(shader, boneMatrixes, vertexDeltas, materialDeltas, windowIndex)

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld()

	// ボーンデバッグ表示
	if isBoneDebug {
		model.Bones.Draw(shader, globalMatrixes, windowIndex)
	}
}

func fetchVertexDeltas(model *pmx.PmxModel, deltas *vmd.VmdDeltas) [][]float32 {
	vertexDeltas := make([][]float32, len(model.Vertices.Data))

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	chunkSize := int(math.Ceil(float64(len(model.Vertices.Data)) / float64(runtime.NumCPU())))
	for chunkStart := 0; chunkStart < len(model.Vertices.Data); chunkStart += chunkSize {
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > len(model.Vertices.Data) {
			chunkEnd = len(model.Vertices.Data)
		}

		go func(chunkStart, chunkEnd int) {
			defer wg.Done()
			for i := chunkStart; i < chunkEnd; i++ {
				if !deltas.Morphs.Vertices.Data[i].Position.IsZero() ||
					!deltas.Morphs.Vertices.Data[i].Uv.IsZero() ||
					!deltas.Morphs.Vertices.Data[i].Uv1.IsZero() ||
					!deltas.Morphs.Vertices.Data[i].AfterPosition.IsZero() {
					// 必要な場合にのみ部分更新するよう設定
					vertexDeltas[i] = deltas.Morphs.Vertices.Data[i].GL()
				}
			}
		}(chunkStart, chunkEnd)
	}
	wg.Wait()

	return vertexDeltas
}

func updatePhysics(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	boneMatrixes []*mgl32.Mat4,
	boneTransforms []*mbt.BtTransform,
	deltas *vmd.VmdDeltas,
	fno int,
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

	if float32(fno) > modelPhysics.Spf {
		modelPhysics.Update(elapsed)

		// 剛体位置を更新
		for _, rigidBody := range model.RigidBodies.GetSortedData() {
			rigidBody.UpdateMatrix(modelPhysics, boneMatrixes, boneTransforms)
		}

		// // 物理後ボーン位置を更新
		// for boneIndex := range model.Bones.LayerSortedIndexes {
		// 	bone := model.Bones.GetItem(boneIndex)
		// 	if bone.IsAfterPhysicsvmd.) && bone.ParentIndex == -1 && model.Bones.Contains(bone.ParentIndex) {
		// 		// 物理後ボーンで親が存在している場合、親の行列を取得する
		// 		parentMat := boneMatrixes[bone.ParentIndex]
		// 		pos := deltas.Bones.GetItem(bone.Name, fno).FramePosition.GL()
		// 		rot := deltas.Bones.GetItem(bone.Name, fno).FrameRotation.GL()
		// 		scl := deltas.Bones.GetItem(bone.Name, fno).FrameScale

		// 		// 自身の行列を作成
		// 		mat := parentMat.Mul4(mgl32.Translate3D(pos[0], pos[1], pos[2]))
		// 		mat = mat.Mul4(mgl32.HomogRotate3D(rot[3], mgl32.Vec3{rot[0], rot[1], rot[2]}))
		// 		mat = mat.Mul4(mgl32.Scale3D(float32(scl[0]), float32(scl[1]), float32(scl[2])))
		// 		boneMatrixes[boneIndex] = &mat
		// 	}
		// }
	}
}
