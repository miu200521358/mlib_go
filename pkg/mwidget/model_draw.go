//go:build windows
// +build windows

package mwidget

import (
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
	prevDeltas *vmd.VmdDeltas,
	frame int,
	elapsed float32,
	isDeform, enablePhysics bool,
) *vmd.VmdDeltas {
	vds := &vmd.VmdDeltas{}
	var beforeBoneDeltas *vmd.BoneDeltas

	// IKのON/OFF
	ikFrame := motion.IkFrames.Get(frame)

	if isDeform || prevDeltas == nil {
		vds.Morphs = motion.DeformMorph(frame, model, nil)
		beforeBoneDeltas = motion.BoneFrames.DeformByPhysicsFlag(frame, model, nil, true,
			nil, vds.Morphs, ikFrame, false)
	} else {
		vds.Morphs = prevDeltas.Morphs
		beforeBoneDeltas = prevDeltas.Bones
	}

	// 物理更新
	updatePhysics(modelPhysics, model, beforeBoneDeltas, frame, elapsed, enablePhysics)

	// 物理後のデフォーム情報
	vds.Bones = motion.BoneFrames.DeformByPhysicsFlag(frame, model, nil, true,
		beforeBoneDeltas, vds.Morphs, ikFrame, true)

	return vds
}

func draw(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	shader *mview.MShader,
	prevDeltas *vmd.VmdDeltas,
	windowIndex int,
	frame int,
	elapsed float32,
	isDeform, enablePhysics bool,
	isDrawNormal bool,
	isDrawWire bool,
	isDrawBones map[pmx.BoneFlag]bool,
) *vmd.VmdDeltas {
	deltas := deform(modelPhysics, model, motion, prevDeltas, frame, elapsed, isDeform, enablePhysics)

	boneDeltas := make([]mgl32.Mat4, len(model.Bones.Data))
	for i, bone := range model.Bones.Data {
		boneDeltas[i] = deltas.Bones.Get(bone.Index).LocalMatrix().GL()
	}

	meshDeltas := make([]*pmx.MeshDelta, len(model.Materials.Data))
	for i, md := range deltas.Morphs.Materials.Data {
		meshDeltas[i] = md.Result()
	}

	vertexDeltas := fetchVertexDeltas(model, deltas)

	model.Meshes.Draw(shader, boneDeltas, vertexDeltas, meshDeltas, windowIndex,
		isDrawNormal, isDrawWire, isDrawBones, model.Bones)

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld()

	return deltas
}

func fetchVertexDeltas(model *pmx.PmxModel, deltas *vmd.VmdDeltas) [][]float32 {
	vertexDeltas := make([][]float32, len(model.Vertices.Data))

	for _, v := range deltas.Morphs.Vertices.Data {
		if v != nil && ((v.Position != nil && !v.Position.IsZero()) ||
			(v.Uv != nil && !v.Uv.IsZero()) ||
			(v.Uv1 != nil && !v.Uv1.IsZero()) ||
			(v.AfterPosition != nil && !v.AfterPosition.IsZero())) {
			// 必要な場合にのみ部分更新するよう設定
			vertexDeltas[v.Index] = v.GL()
		}
	}

	return vertexDeltas
}

// func fetchVertexDeltasParallel(model *pmx.PmxModel, deltas *vmd.VmdDeltas) [][]float32 {
// 	vertexDeltas := make([][]float32, len(model.Vertices.Data))

// 	var wg sync.WaitGroup

// 	chunkSize := int(math.Ceil(float64(len(model.Vertices.Data)) / float64(runtime.NumCPU())))
// 	wg.Add(runtime.NumCPU())
// 	for chunkStart := 0; chunkStart < len(model.Vertices.Data); chunkStart += chunkSize {
// 		chunkEnd := chunkStart + chunkSize
// 		if chunkEnd > len(model.Vertices.Data) {
// 			chunkEnd = len(model.Vertices.Data)
// 		}

// 		go func(chunkStart, chunkEnd int) {
// 			defer wg.Done()
// 			for i := chunkStart; i < chunkEnd; i++ {
// 				if deltas.Morphs.Vertices.Data[i] != nil &&
// 					(!deltas.Morphs.Vertices.Data[i].Position.IsZero() ||
// 						!deltas.Morphs.Vertices.Data[i].Uv.IsZero() ||
// 						!deltas.Morphs.Vertices.Data[i].Uv1.IsZero() ||
// 						!deltas.Morphs.Vertices.Data[i].AfterPosition.IsZero()) {
// 					// 必要な場合にのみ部分更新するよう設定
// 					vertexDeltas[i] = deltas.Morphs.Vertices.Data[i].GL()
// 				}
// 			}
// 		}(chunkStart, chunkEnd)
// 	}
// 	wg.Wait()

// 	return vertexDeltas
// }

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

	// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[1]", frame))

	for _, rigidBody := range model.RigidBodies.Data {
		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}
		if rigidBodyBone == nil || boneDeltas.Get(rigidBodyBone.Index) == nil {
			continue
		}

		boneTransform := mbt.NewBtTransform()
		defer mbt.DeleteBtTransform(boneTransform)
		// if r.CorrectPhysicsType == pmx.PHYSICS_TYPE_DYNAMIC_BONE {
		// 	mat := boneDeltas.Get(rigidBodyBone.Index).GlobalMatrix()
		// 	bonePhysicsGlobalMatrix := r.GetRigidBodyBoneMatrix(modelPhysics)

		// 	bonePhysicsGlobalMatrix[3] = mat[3]
		// 	bonePhysicsGlobalMatrix[7] = mat[7]
		// 	bonePhysicsGlobalMatrix[11] = mat[11]
		// } else {
		mat := boneDeltas.Get(rigidBodyBone.Index).GlobalMatrix().GL()
		boneTransform.SetFromOpenGLMatrix(&mat[0])
		// }

		// 物理フラグが落ちている場合があるので、強制的に起こす
		forceUpdate := rigidBody.UpdateFlags(modelPhysics, enablePhysics)
		rigidBody.UpdateTransform(modelPhysics, rigidBodyBone, boneTransform,
			elapsed == 0.0 || !enablePhysics || forceUpdate)

		// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[2][%d]", frame, rigidBody.Index))
	}

	if elapsed >= 1e-5 {
		modelPhysics.Update(elapsed)

		// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[3]", frame))

		// 剛体位置を更新
		for _, rigidBody := range model.RigidBodies.Data {
			bonePhysicsGlobalMatrix := rigidBody.GetRigidBodyBoneMatrix(modelPhysics)
			if boneDeltas != nil && bonePhysicsGlobalMatrix != nil && rigidBody.Bone != nil {
				// if rigidBody.CorrectPhysicsType == pmx.PHYSICS_TYPE_DYNAMIC_BONE &&
				// 	boneDeltas.Get(rigidBody.Bone.Index) != nil {
				// 	// ボーン位置合わせの場合、位置情報は計算したのを使う
				// 	globalMatrix := boneDeltas.Get(rigidBody.Bone.Index).GlobalMatrix()
				// 	bonePhysicsGlobalMatrix[3] = globalMatrix[3]
				// 	bonePhysicsGlobalMatrix[7] = globalMatrix[7]
				// 	bonePhysicsGlobalMatrix[11] = globalMatrix[11]
				// }
				if boneDeltas.Get(rigidBody.Bone.Index) == nil {
					boneDeltas.Append(&vmd.BoneDelta{Bone: rigidBody.Bone, Frame: frame})
				}
				boneDeltas.SetGlobalMatrix(rigidBody.Bone, bonePhysicsGlobalMatrix)
			}

			// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[4][%d]", frame, rigidBody.Index))
		}

		// グローバル行列を埋め終わったらローカル行列の計算
		boneDeltas.FillLocalMatrix()
	}
}
