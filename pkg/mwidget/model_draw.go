//go:build windows
// +build windows

package mwidget

import (
	"math"
	"runtime"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
	"github.com/miu200521358/mlib_go/pkg/mview"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

func deform(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	frame int,
	elapsed float32,
	enablePhysics bool,
) *vmd.VmdDeltas {
	vds := &vmd.VmdDeltas{}

	vds.Morphs = motion.DeformMorph(frame, model, nil)

	for i, bd := range vds.Morphs.Bones.Data {
		if bd == nil {
			continue
		}
		bone := model.Bones.Get(i)
		if !motion.BoneFrames.Contains(bone.Name) {
			motion.BoneFrames.Append(vmd.NewBoneNameFrames(bone.Name))
		}
		bf := motion.BoneFrames.Get(bone.Name).Get(frame)

		// 一旦モーフの値をクリア
		bf.MorphPosition = nil
		bf.MorphLocalPosition = nil
		bf.MorphRotation = nil
		bf.MorphLocalRotation = nil
		bf.MorphScale = nil
		bf.MorphLocalScale = nil

		// 該当ボーンキーフレにモーフの値を加算
		bf.Add(bd.BoneFrame)
		motion.AppendBoneFrame(bone.Name, bf)
	}

	// 物理前のデフォーム情報
	beforeBoneDeltas := motion.BoneFrames.Deform(frame, model, nil, true, nil)

	// 物理更新
	updatePhysics(modelPhysics, model, beforeBoneDeltas, frame, elapsed, enablePhysics)

	// 物理後のデフォーム情報
	vds.Bones = motion.BoneFrames.Deform(frame, model, nil, true, beforeBoneDeltas)

	return vds
}

func draw(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	shader *mview.MShader,
	windowIndex int,
	frame int,
	elapsed float32,
	enablePhysics bool,
	isDrawNormal bool,
	isDrawBone bool,
) {
	deltas := deform(modelPhysics, model, motion, frame, elapsed, enablePhysics)

	boneDeltas := make([]mgl32.Mat4, len(model.Bones.Data))
	for i, bone := range model.Bones.Data {
		boneDeltas[i] = deltas.Bones.Get(bone.Index).LocalMatrix().GL()
	}

	materialDeltas := make([]*pmx.Material, len(model.Materials.Data))
	for i, md := range deltas.Morphs.Materials.Data {
		materialDeltas[i] = md.Material
	}

	vertexDeltas := fetchVertexDeltas(model, deltas)

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
	boneDeltas *vmd.BoneDeltas,
	frame int,
	elapsed float32,
	enablePhysics bool,
) {
	if modelPhysics == nil {
		return
	}

	boneTransforms := make([]*mbt.BtTransform, len(model.Bones.Data))
	for _, delta := range boneDeltas.Data {
		mat := delta.GlobalMatrix().GL()
		t := mbt.NewBtTransform()
		t.SetFromOpenGLMatrix(&mat[0])
		boneTransforms[delta.Bone.Index] = &t
	}

	for _, r := range model.RigidBodies.GetSortedData() {
		// 物理フラグが落ちている場合があるので、強制的に起こす
		forceUpdate := r.UpdateFlags(modelPhysics, enablePhysics)
		r.UpdateTransform(modelPhysics, boneTransforms, elapsed == 0.0 || !enablePhysics || forceUpdate)
	}

	if float32(frame) > modelPhysics.Spf {
		modelPhysics.Update(elapsed)

		// 剛体位置を更新
		for _, rigidBody := range model.RigidBodies.GetSortedData() {
			bonePhysicsGlobalMatrix := rigidBody.UpdateMatrix(modelPhysics)
			if boneDeltas != nil && bonePhysicsGlobalMatrix != nil && rigidBody.Bone != nil {
				if rigidBody.CorrectPhysicsType == pmx.PHYSICS_TYPE_DYNAMIC_BONE &&
					boneDeltas.Get(rigidBody.Bone.Index) != nil {
					// ボーン位置合わせの場合、位置情報は計算したのを使う
					globalMatrix := boneDeltas.Get(rigidBody.Bone.Index).GlobalMatrix()
					bonePhysicsGlobalMatrix[3] = globalMatrix[3]
					bonePhysicsGlobalMatrix[7] = globalMatrix[7]
					bonePhysicsGlobalMatrix[11] = globalMatrix[11]
				}
				if boneDeltas.Get(rigidBody.Bone.Index) == nil {
					boneDeltas.Append(&vmd.BoneDelta{Bone: rigidBody.Bone, Frame: frame})
				}
				boneDeltas.SetGlobalMatrix(rigidBody.Bone, bonePhysicsGlobalMatrix)
			}
		}
	}
}
