//go:build windows
// +build windows

package mwidget

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/mmath"
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
	timeStep float32,
	enablePhysics, resetPhysics bool,
) *vmd.VmdDeltas {
	vds := vmd.NewVmdDeltas()
	var beforeBoneDeltas *vmd.BoneDeltas

	// IKのON/OFF
	ikFrame := motion.IkFrames.Get(frame)

	if prevDeltas == nil {
		vds.Morphs = motion.DeformMorph(frame, model, nil)
		beforeBoneDeltas = motion.BoneFrames.DeformByPhysicsFlag(frame, model, nil, true,
			nil, vds.Morphs, ikFrame, false)
	} else {
		vds.Morphs = prevDeltas.Morphs
		beforeBoneDeltas = prevDeltas.Bones
	}

	// 物理更新
	updatePhysics(modelPhysics, model, beforeBoneDeltas, frame, timeStep, enablePhysics, resetPhysics)

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
	invisibleMaterialIndexes, nextInvisibleMaterialIndexes, selectedVertexIndexes, nextSelectedVertexIndexes []int,
	windowIndex int,
	frame int,
	timeStep float32,
	enablePhysics, resetPhysics, isDrawNormal, isDrawWire, isDrawSelectedVertex bool,
	isDrawBones map[pmx.BoneFlag]bool,
) *vmd.VmdDeltas {
	deltas := deform(modelPhysics, model, motion, prevDeltas, frame, timeStep, enablePhysics, resetPhysics)

	boneDeltas := make([]mgl32.Mat4, len(model.Bones.Data))
	for i, bone := range model.Bones.Data {
		boneDeltas[i] = deltas.Bones.Get(bone.Index).LocalMatrix().GL()
	}

	meshDeltas := make([]*pmx.MeshDelta, len(model.Materials.Data))
	for i, md := range deltas.Morphs.Materials.Data {
		meshDeltas[i] = md.Result()
	}

	vertexMorphIndexes, vertexMorphDeltas, selectedVertexMorphIndexes, selectedVertexDeltas :=
		fetchVertexDeltas(model, deltas, selectedVertexIndexes, nextSelectedVertexIndexes)

	vertexPositions := model.Meshes.Draw(shader, boneDeltas, vertexMorphIndexes, vertexMorphDeltas,
		selectedVertexMorphIndexes, selectedVertexDeltas, meshDeltas,
		invisibleMaterialIndexes, nextInvisibleMaterialIndexes, windowIndex,
		isDrawNormal, isDrawWire, isDrawSelectedVertex, isDrawBones, model.Bones)

	for i, pos := range vertexPositions {
		deltas.Vertices.Data[i] = vmd.NewVertexDelta(&mmath.MVec3{float64(-pos[0]), float64(pos[1]), float64(pos[2])})
	}

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld()

	return deltas
}

func fetchVertexDeltas(
	model *pmx.PmxModel, deltas *vmd.VmdDeltas,
	selectedVertexIndexes, nextSelectedVertexIndexes []int,
) ([]int, [][]float32, []int, [][]float32) {
	vertexMorphIndexes, vertexMorphDeltas := deltas.Morphs.Vertices.GL()
	selectedVertexMorphIndexes, selectedVertexDeltas := deltas.SelectedVertexDeltas.GL(
		model, selectedVertexIndexes, nextSelectedVertexIndexes)

	return vertexMorphIndexes, vertexMorphDeltas, selectedVertexMorphIndexes, selectedVertexDeltas
}

func updatePhysics(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	boneDeltas *vmd.BoneDeltas,
	frame int,
	timeStep float32,
	enablePhysics bool,
	resetPhysics bool,
) {
	if modelPhysics == nil {
		return
	}

	// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[1]", frame))

	for i := range model.RigidBodies.Len() {
		rigidBody := model.RigidBodies.Get(i)
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
		forceUpdate := rigidBody.UpdateFlags(model.Index, modelPhysics, enablePhysics, resetPhysics)
		rigidBody.UpdateTransform(model.Index, modelPhysics, rigidBodyBone, boneTransform,
			timeStep == 0.0 || !enablePhysics || forceUpdate)

		// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[2][%d]", frame, rigidBody.Index))
	}

	if (enablePhysics || resetPhysics) && timeStep >= 1e-5 {
		modelPhysics.Update(timeStep)

		// 剛体位置を更新
		physicsBoneIndexes := make([]int, 0)
		for i := range model.RigidBodies.Len() {
			rigidBody := model.RigidBodies.Get(i)
			bonePhysicsGlobalMatrix := rigidBody.GetRigidBodyBoneMatrix(model.Index, modelPhysics)
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
					boneDeltas.Update(&vmd.BoneDelta{Bone: rigidBody.Bone, Frame: frame})
				}
				boneDeltas.SetGlobalMatrix(rigidBody.Bone, bonePhysicsGlobalMatrix)
				physicsBoneIndexes = append(physicsBoneIndexes, rigidBody.Bone.Index)
			}

			// mlog.Memory(fmt.Sprintf("[%d] updatePhysics[4][%d]", frame, rigidBody.Index))
		}

		// グローバル行列を埋め終わったらローカル行列の計算
		boneDeltas.FillLocalMatrix(physicsBoneIndexes)
	}
}
