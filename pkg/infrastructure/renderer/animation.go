package renderer

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type AnimationState struct {
	RenderModel              *RenderModel     // 描画モデル
	Model                    *pmx.PmxModel    // モデル
	Motion                   *vmd.VmdMotion   // モーション
	VmdDeltas                *delta.VmdDeltas // モーション変化量
	InvisibleMaterialIndexes []int            // 非表示材質インデックス
	SelectedVertexIndexes    []int            // 選択頂点インデックス
	vertexMorphDeltaIndexes  []int            // 頂点モーフインデックス
	vertexMorphDeltas        [][]float32      // 頂点モーフデルタ
	meshDeltas               []*MeshDelta     // メッシュデルタ
}

type AnimationStates struct {
	Now  *AnimationState
	Next *AnimationState
}

func NewAnimationState() *AnimationState {
	return &AnimationState{
		InvisibleMaterialIndexes: make([]int, 0),
		SelectedVertexIndexes:    make([]int, 0),
	}
}

func NewAnimationStates() *AnimationStates {
	return &AnimationStates{
		Now:  NewAnimationState(),
		Next: NewAnimationState(),
	}
}

func Animate(
	physics *mbt.MPhysics, animationStates []*AnimationStates,
	frame int, timeStep float32, enabledPhysics, resetPhysics bool,
) []*AnimationStates {
	// 物理前デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i].Now.Model == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii].Now =
					animateBeforePhysics(physics, animationStates[ii].Now, int(frame), enabledPhysics)
			}(i)
		}

		wg.Wait()
	}

	if enabledPhysics || resetPhysics {
		// 物理更新
		physics.StepSimulation(timeStep)
	}

	// 物理後デフォーム
	{
		var wg sync.WaitGroup

		for i := range animationStates {
			if animationStates[i].Now.RenderModel == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()
				animationStates[ii].Now =
					animateAfterPhysics(physics, animationStates[ii].Now, int(frame), enabledPhysics, resetPhysics)
			}(i)
		}

		wg.Wait()
	}

	return animationStates
}

func animateBeforePhysics(
	physics *mbt.MPhysics, animationState *AnimationState,
	frame int, enabledPhysics bool,
) *AnimationState {
	if animationState.Motion == nil {
		animationState.Motion = vmd.NewVmdMotion("")
	}

	deltas := delta.NewVmdDeltas(animationState.Model.Materials, animationState.Model.Bones)

	if animationState.VmdDeltas == nil {
		deltas.Morphs = deform.DeformMorph(animationState.Model, animationState.Motion.MorphFrames, frame, nil)
		deltas = deform.DeformBoneByPhysicsFlag(animationState.Model,
			animationState.Motion, deltas, true, frame, nil, false)

		animationState.vertexMorphDeltaIndexes, animationState.vertexMorphDeltas =
			newVertexMorphDeltasGl(deltas.Morphs.Vertices)

		animationState.meshDeltas = make([]*MeshDelta, len(animationState.Model.Materials.Data))
		for i, md := range deltas.Morphs.Materials.Data {
			animationState.meshDeltas[i] = newMeshDelta(md)
		}
	} else {
		deltas.Morphs = animationState.VmdDeltas.Morphs
		deltas.Bones = animationState.VmdDeltas.Bones
	}

	modelIndex := 0

	for _, rigidBody := range animationState.Model.RigidBodies.Data {
		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}
		if rigidBodyBone == nil || deltas.Bones.Get(rigidBodyBone.Index) == nil {
			continue
		}

		if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC || !enabledPhysics {
			// ボーン追従剛体・物理＋ボーン位置もしくは強制更新の場合のみ剛体位置更新
			boneTransform := bt.NewBtTransform()
			defer bt.DeleteBtTransform(boneTransform)
			mat := mgl.NewGlMat4(deltas.Bones.Get(rigidBodyBone.Index).FilledGlobalMatrix())
			boneTransform.SetFromOpenGLMatrix(&mat[0])

			physics.UpdateTransform(modelIndex, rigidBodyBone, boneTransform, rigidBody)
		}
	}

	animationState.VmdDeltas = deltas

	return animationState
}

func animateAfterPhysics(
	physics *mbt.MPhysics, animationState *AnimationState,
	frame int, enabledPhysics, resetPhysics bool,
) *AnimationState {
	modelIndex := 0

	// 物理剛体位置を更新
	if enabledPhysics || resetPhysics {
		for _, isAfterPhysics := range []bool{false, true} {
			for _, bone := range animationState.Model.Bones.LayerSortedBones[isAfterPhysics] {
				if bone.Extend.RigidBody == nil || bone.Extend.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
					continue
				}
				bonePhysicsGlobalMatrix := physics.GetRigidBodyBoneMatrix(modelIndex, bone.Extend.RigidBody)
				if animationState.VmdDeltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
					bd := delta.NewBoneDeltaByGlobalMatrix(bone, frame,
						bonePhysicsGlobalMatrix, animationState.VmdDeltas.Bones.Get(bone.ParentIndex))
					animationState.VmdDeltas.Bones.Update(bd)
				}
			}
		}
	}

	// 物理後のデフォーム情報
	animationState.VmdDeltas = deform.DeformBoneByPhysicsFlag(animationState.Model,
		animationState.Motion, animationState.VmdDeltas, true, frame, nil, true)

	// // 選択頂点モーフの設定は常に更新する
	// SelectedVertexIndexesDeltas, SelectedVertexGlDeltas := renderer.SelectedVertexMorphDeltasGL(
	// 	SelectedVertexDeltas, model, selectedVertexIndexes, nextSelectedVertexIndexes)

	return animationState
}
