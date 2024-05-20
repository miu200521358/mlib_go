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

func draw(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	shader *mview.MShader,
	deltas *vmd.VmdDeltas,
	windowIndex int,
	fno int,
	elapsed float32,
	enablePhysics bool,
	isDrawNormal bool,
	isDrawBone bool,
) {
	boneDeltas := make([]*mgl32.Mat4, len(model.Bones.Data))
	globalMatrixes := make([]*mmath.MMat4, len(model.Bones.Data))
	boneTransforms := make([]*mbt.BtTransform, len(model.Bones.Data))
	materialDeltas := make([]*pmx.Material, len(model.Materials.Data))
	for i, bone := range model.Bones.Data {
		delta := deltas.Bones.Get(bone.Name)
		mat := delta.LocalMatrix.GL()
		boneDeltas[i] = &mat
		globalMatrixes[i] = delta.GlobalMatrix
		t := mbt.NewBtTransform()
		t.SetFromOpenGLMatrix(&mat[0])
		boneTransforms[i] = &t
	}

	for i, md := range deltas.Morphs.Materials.Data {
		materialDeltas[i] = md.Material
	}

	vertexDeltas := fetchVertexDeltas(model, deltas)

	updatePhysics(modelPhysics, model, boneDeltas, boneTransforms, deltas, fno, elapsed, enablePhysics)

	model.Meshes.Draw(shader, boneDeltas, vertexDeltas, materialDeltas, windowIndex, isDrawNormal, isDrawBone)

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld()
}

func fetchVertexDeltas(model *pmx.PmxModel, deltas *vmd.VmdDeltas) [][]float32 {
	vertexDeltas := make([][]float32, len(model.Vertices.Data))

	chunkSize := int(math.Ceil(float64(len(model.Vertices.Data)) / float64(runtime.NumCPU())))
	if chunkSize > 20000 {
		return fetchVertexDeltasParallel(model, deltas)
	}

	for i, v := range deltas.Morphs.Vertices.Data {
		if v != nil && (!v.Position.IsZero() || !v.Uv.IsZero() || !v.Uv1.IsZero() || !v.AfterPosition.IsZero()) {
			// 必要な場合にのみ部分更新するよう設定
			vertexDeltas[i] = v.GL()
		}
	}

	return vertexDeltas
}

func fetchVertexDeltasParallel(model *pmx.PmxModel, deltas *vmd.VmdDeltas) [][]float32 {
	vertexDeltas := make([][]float32, len(model.Vertices.Data))

	var wg sync.WaitGroup

	chunkSize := int(math.Ceil(float64(len(model.Vertices.Data)) / float64(runtime.NumCPU())))
	wg.Add(runtime.NumCPU())
	for chunkStart := 0; chunkStart < len(model.Vertices.Data); chunkStart += chunkSize {
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > len(model.Vertices.Data) {
			chunkEnd = len(model.Vertices.Data)
		}

		go func(chunkStart, chunkEnd int) {
			defer wg.Done()
			for i := chunkStart; i < chunkEnd; i++ {
				if deltas.Morphs.Vertices.Data[i] != nil &&
					(!deltas.Morphs.Vertices.Data[i].Position.IsZero() ||
						!deltas.Morphs.Vertices.Data[i].Uv.IsZero() ||
						!deltas.Morphs.Vertices.Data[i].Uv1.IsZero() ||
						!deltas.Morphs.Vertices.Data[i].AfterPosition.IsZero()) {
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
	boneDeltas []*mgl32.Mat4,
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
			physicsBoneMatrix := rigidBody.UpdateMatrix(modelPhysics, boneDeltas, boneTransforms)
			if len(model.Bones.AfterPhysicsBoneIndexes) > 0 && enablePhysics && physicsBoneMatrix != nil {
				// 物理後ボーン用に行列を更新
				bone := model.Bones.Get(rigidBody.BoneIndex)
				delta := deltas.Bones.Get(bone.Name)
				delta.LocalMatrix = mmath.NewMMat4ByMgl(physicsBoneMatrix)
				delta.FramePosition = delta.LocalMatrix.Translation()
				delta.FrameRotation = delta.LocalMatrix.Quaternion()
			}
		}

		if len(model.Bones.AfterPhysicsBoneIndexes) > 0 && enablePhysics {
			// 物理後ボーン位置更新の場合
			for _, boneIndex := range model.Bones.AfterPhysicsBoneIndexes {
				// 物理後ボーンで親が存在している場合、親の行列を取得する
				updateBoneMatrixAfterPhysics(model, boneDeltas, deltas, model.Bones.Get(boneIndex))
			}
		}
	}
}

func updateBoneMatrixAfterPhysics(
	model *pmx.PmxModel,
	boneDeltas []*mgl32.Mat4,
	deltas *vmd.VmdDeltas,
	bone *pmx.Bone,
) {
	delta := deltas.Bones.Get(bone.Name)

	pos := delta.FramePosition
	rot := delta.FrameRotationWithoutEffect
	scl := delta.FrameScale

	if bone.IsEffectorTranslation() && bone.EffectIndex >= 0 {
		effectBone := model.Bones.Get(bone.EffectIndex)
		pos.Add(deltas.Bones.Get(effectBone.Name).FramePosition)
	}

	if bone.IsEffectorRotation() && bone.EffectIndex >= 0 {
		effectBone := model.Bones.Get(bone.EffectIndex)
		rot = rot.Mul(deltas.Bones.Get(effectBone.Name).FrameRotation)
	}

	matrix := mmath.NewMMat4()
	matrix.Scale(scl)
	matrix.Rotate(rot)
	matrix.Translate(pos)
	matrix.Mul(bone.RevertOffsetMatrix)

	// 自身の行列を再作成
	parentBone := model.Bones.Get(bone.ParentIndex)
	mat := matrix.Muled(deltas.Bones.Get(parentBone.Name).LocalMatrix).GL()
	boneDeltas[bone.Index] = &mat

	// for _, childBoneIndex := range bone.ChildBoneIndexes {
	// 	// 子ボーンがいる場合、子ボーンの行列を再計算
	// 	updateBoneMatrixAfterPhysics(model, boneDeltas, deltas, fno, model.Bones.Get(childBoneIndex))
	// }
}
