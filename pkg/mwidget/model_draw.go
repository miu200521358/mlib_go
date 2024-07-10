//go:build windows
// +build windows

package mwidget

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mphysics/mbt"
	"github.com/miu200521358/mlib_go/pkg/mview"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

type ModelSet struct {
	Model                        *pmx.PmxModel  // 現在描画中のモデル
	Motion                       *vmd.VmdMotion // 現在描画中のモーション
	InvisibleMaterialIndexes     []int          // 非表示材質インデックス
	SelectedVertexIndexes        []int          // 選択頂点インデックス
	NextModel                    *pmx.PmxModel  // UIから渡された次のモデル
	NextMotion                   *vmd.VmdMotion // UIから渡された次のモーション
	NextInvisibleMaterialIndexes []int          // UIから渡された次の非表示材質インデックス
	NextSelectedVertexIndexes    []int          // UIから渡された次の選択頂点インデックス
	prevDeltas                   *vmd.VmdDeltas // 前回のデフォーム情報
}

func NewModelSet() *ModelSet {
	return &ModelSet{
		InvisibleMaterialIndexes:     make([]int, 0),
		SelectedVertexIndexes:        make([]int, 0),
		NextInvisibleMaterialIndexes: make([]int, 0),
		NextSelectedVertexIndexes:    make([]int, 0),
	}
}

func deformsAll(
	modelPhysics *mphysics.MPhysics,
	modelSets []*ModelSet,
	frame, prevFrame int,
	timeStep float32,
	enablePhysics, resetPhysics bool,
) []*ModelSet {

	// 物理前デフォーム
	{
		var wg sync.WaitGroup

		for i := range modelSets {
			if modelSets[i].Model == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()

				if int(frame) != prevFrame {
					// フレーム番号が変わっている場合は前回デフォームを破棄
					modelSets[ii].prevDeltas = nil
				}

				modelSets[ii].prevDeltas = deformBeforePhysics(
					modelPhysics, modelSets[ii].Model, modelSets[ii].Motion, modelSets[ii].prevDeltas,
					int(frame), timeStep, enablePhysics, resetPhysics,
				)
			}(i)
		}

		wg.Wait()
	}

	if enablePhysics || resetPhysics {
		// 物理更新
		modelPhysics.Update(timeStep)
	}

	// 物理後デフォーム
	{
		var wg sync.WaitGroup

		for i := range modelSets {
			if modelSets[i].Model == nil {
				continue
			}

			wg.Add(1)
			go func(ii int) {
				defer wg.Done()

				modelSets[ii].prevDeltas = deformAfterPhysics(
					modelPhysics, modelSets[ii].Model, modelSets[ii].Motion, modelSets[ii].prevDeltas,
					modelSets[ii].SelectedVertexIndexes, modelSets[ii].NextSelectedVertexIndexes,
					int(frame), enablePhysics, resetPhysics,
				)
			}(i)
		}

		wg.Wait()
	}

	return modelSets
}

func deformBeforePhysics(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	prevDeltas *vmd.VmdDeltas,
	frame int,
	timeStep float32,
	enablePhysics, resetPhysics bool,
) *vmd.VmdDeltas {
	if motion == nil {
		motion = vmd.NewVmdMotion("")
	}

	vds := vmd.NewVmdDeltas(model.Vertices)

	// IKのON/OFF
	ikFrame := motion.IkFrames.Get(frame)

	if prevDeltas == nil {
		vds.Morphs = motion.DeformMorph(frame, model, nil)
		vds.Bones = motion.BoneFrames.DeformByPhysicsFlag(frame, model, nil, true,
			nil, vds.Morphs, ikFrame, false)

		vds.MeshGlDeltas = make([]*pmx.MeshDelta, len(model.Materials.Data))
		for i, md := range vds.Morphs.Materials.Data {
			vds.MeshGlDeltas[i] = md.Result()
		}

		vds.VertexMorphIndexes, vds.VertexMorphGlDeltas = vds.Morphs.Vertices.GL()
	} else {
		vds.Morphs = prevDeltas.Morphs
		vds.Bones = prevDeltas.Bones

		vds.MeshGlDeltas = prevDeltas.MeshGlDeltas

		vds.VertexMorphIndexes = prevDeltas.VertexMorphIndexes
		vds.VertexMorphGlDeltas = prevDeltas.VertexMorphGlDeltas
	}

	for _, rigidBody := range model.RigidBodies.Data {
		// 現在のボーン変形情報を保持
		rigidBodyBone := rigidBody.Bone
		if rigidBodyBone == nil {
			rigidBodyBone = rigidBody.JointedBone
		}
		if rigidBodyBone == nil || vds.Bones.Get(rigidBodyBone.Index) == nil {
			continue
		}

		// 物理フラグが落ちている場合があるので、強制的に起こす
		forceUpdate := rigidBody.UpdateFlags(model.Index, modelPhysics, enablePhysics, resetPhysics)
		forceUpdate = timeStep == 0.0 || !enablePhysics || forceUpdate

		if rigidBody.CorrectPhysicsType != pmx.PHYSICS_TYPE_DYNAMIC || forceUpdate {
			// ボーン追従剛体・物理＋ボーン位置もしくは強制更新の場合のみ剛体位置更新
			boneTransform := mbt.NewBtTransform()
			defer mbt.DeleteBtTransform(boneTransform)
			mat := vds.Bones.Get(rigidBodyBone.Index).GlobalMatrix().GL()
			boneTransform.SetFromOpenGLMatrix(&mat[0])

			rigidBody.UpdateTransform(model.Index, modelPhysics, rigidBodyBone, boneTransform)
		}
	}

	return vds
}

func deformAfterPhysics(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *vmd.VmdDeltas,
	selectedVertexIndexes, nextSelectedVertexIndexes []int,
	frame int,
	enablePhysics, resetPhysics bool,
) *vmd.VmdDeltas {
	if motion == nil {
		motion = vmd.NewVmdMotion("")
	}

	// IKのON/OFF
	ikFrame := motion.IkFrames.Get(frame)

	if enablePhysics || resetPhysics {
		// 物理剛体位置を更新
		for _, isAfterPhysics := range []bool{false, true} {
			for _, bone := range model.Bones.LayerSortedBones[isAfterPhysics] {
				if bone.RigidBody == nil || bone.RigidBody.CorrectPhysicsType == pmx.PHYSICS_TYPE_STATIC {
					continue
				}
				bonePhysicsGlobalMatrix := bone.RigidBody.GetRigidBodyBoneMatrix(model.Index, modelPhysics)
				if deltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
					bd := vmd.NewBoneDeltaByGlobalMatrix(bone, frame,
						bonePhysicsGlobalMatrix, deltas.Bones.Get(bone.ParentIndex))
					deltas.Bones.Update(bd)
				}
			}
		}
	}

	// 物理後のデフォーム情報
	deltas.Bones = motion.BoneFrames.DeformByPhysicsFlag(frame, model, nil, true,
		deltas.Bones, deltas.Morphs, ikFrame, true)

	// GL描画用データの作成
	deltas.BoneGlDeltas = make([]mgl32.Mat4, len(model.Bones.Data))
	for i, bone := range model.Bones.Data {
		delta := deltas.Bones.Get(bone.Index)
		if delta != nil {
			deltas.BoneGlDeltas[i] = delta.LocalMatrix().GL()
		}
	}

	// 選択頂点モーフの設定は常に更新する
	deltas.SelectedVertexIndexes, deltas.SelectedVertexGlDeltas = deltas.SelectedVertexDeltas.GL(
		model, selectedVertexIndexes, nextSelectedVertexIndexes)

	return deltas
}

func draw(
	modelPhysics *mphysics.MPhysics,
	model *pmx.PmxModel,
	shader *mview.MShader,
	deltas *vmd.VmdDeltas,
	invisibleMaterialIndexes, nextInvisibleMaterialIndexes []int,
	windowIndex int,
	isDrawNormal, isDrawWire, isDrawSelectedVertex bool,
	isDrawBones map[pmx.BoneFlag]bool,
) *vmd.VmdDeltas {
	vertexPositions := model.Meshes.Draw(
		shader, deltas.BoneGlDeltas, deltas.VertexMorphIndexes, deltas.VertexMorphGlDeltas,
		deltas.SelectedVertexIndexes, deltas.SelectedVertexGlDeltas, deltas.MeshGlDeltas,
		invisibleMaterialIndexes, nextInvisibleMaterialIndexes, windowIndex,
		isDrawNormal, isDrawWire, isDrawSelectedVertex, isDrawBones, model.Bones)

	for i, pos := range vertexPositions {
		deltas.Vertices.Data[i] = vmd.NewVertexDelta(&mmath.MVec3{float64(-pos[0]), float64(pos[1]), float64(pos[2])})
	}

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld()

	return deltas
}
