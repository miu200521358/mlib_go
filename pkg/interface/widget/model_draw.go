//go:build windows
// +build windows

package widget

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
)

type ModelSet struct {
	Model                        *pmx.PmxModel    // 現在描画中のモデル
	Motion                       *vmd.VmdMotion   // 現在描画中のモーション
	Meshes                       *renderer.Meshes // 現在描画中のメッシュ
	InvisibleMaterialIndexes     []int            // 非表示材質インデックス
	SelectedVertexIndexes        []int            // 選択頂点インデックス
	NextModel                    *pmx.PmxModel    // UIから渡された次のモデル
	NextMotion                   *vmd.VmdMotion   // UIから渡された次のモーション
	NextInvisibleMaterialIndexes []int            // UIから渡された次の非表示材質インデックス
	NextSelectedVertexIndexes    []int            // UIから渡された次の選択頂点インデックス
	PrevDeltas                   *delta.VmdDeltas // 前回のデフォーム情報
	BoneGlDeltas                 []mgl32.Mat4
	MeshGlDeltas                 []*renderer.MeshDelta
	VertexMorphIndexes           []int
	VertexMorphGlDeltas          [][]float32
	SelectedVertexIndexesDeltas  []int
	SelectedVertexGlDeltasDeltas [][]float32
	Vertices                     *renderer.VertexDeltas
	SelectedVertexDeltas         *renderer.SelectedVertexMorphDeltas
}

func NewModelSet() *ModelSet {
	return &ModelSet{
		InvisibleMaterialIndexes:     make([]int, 0),
		SelectedVertexIndexes:        make([]int, 0),
		NextInvisibleMaterialIndexes: make([]int, 0),
		NextSelectedVertexIndexes:    make([]int, 0),
		SelectedVertexDeltas:         renderer.NewSelectedVertexMorphDeltas(),
	}
}

func DeformsAll(
	modelPhysics *mbt.MPhysics,
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
					modelSets[ii].PrevDeltas = nil
				}

				modelSets[ii].PrevDeltas, modelSets[ii].MeshGlDeltas, modelSets[ii].VertexMorphIndexes, modelSets[ii].VertexMorphGlDeltas = deformBeforePhysics(
					modelPhysics, modelSets[ii].Model, modelSets[ii].Motion, modelSets[ii].PrevDeltas,
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

				modelSets[ii].PrevDeltas, modelSets[ii].BoneGlDeltas, modelSets[ii].SelectedVertexIndexesDeltas, modelSets[ii].SelectedVertexGlDeltasDeltas = deformAfterPhysics(
					modelPhysics, modelSets[ii].Model, modelSets[ii].Motion, modelSets[ii].PrevDeltas,
					modelSets[ii].SelectedVertexIndexes, modelSets[ii].NextSelectedVertexIndexes, modelSets[ii].SelectedVertexDeltas,
					int(frame), enablePhysics, resetPhysics,
				)
			}(i)
		}

		wg.Wait()
	}

	return modelSets
}

func deformBeforePhysics(
	modelPhysics *mbt.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	prevDeltas *delta.VmdDeltas,
	frame int,
	timeStep float32,
	enablePhysics, resetPhysics bool,
) (*delta.VmdDeltas, []*renderer.MeshDelta, []int, [][]float32) {
	if motion == nil {
		motion = vmd.NewVmdMotion("")
	}

	vds := delta.NewVmdDeltas(model.Vertices)

	// IKのON/OFF
	ikFrame := motion.IkFrames.Get(frame)

	if prevDeltas == nil {
		vds.Morphs = deform.DeformMorph(motion, motion.MorphFrames, frame, model, nil)
		vds.Bones = deform.DeformByPhysicsFlag(motion.BoneFrames, frame, model, nil, true,
			nil, vds.Morphs, ikFrame, false)
	} else {
		vds.Morphs = prevDeltas.Morphs
		vds.Bones = prevDeltas.Bones
	}

	MeshGlDeltas := make([]*renderer.MeshDelta, len(model.Materials.Data))
	for i, md := range vds.Morphs.Materials.Data {
		MeshGlDeltas[i] = renderer.MaterialMorphDeltaResult(md)
	}

	VertexMorphIndexes, VertexMorphGlDeltas := renderer.VertexMorphDeltasGL(vds.Morphs.Vertices)

	modelIndex := 0

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
		forceUpdate := mbt.UpdateFlags(modelIndex, modelPhysics, rigidBody, enablePhysics, resetPhysics)
		forceUpdate = timeStep == 0.0 || !enablePhysics || forceUpdate

		if rigidBody.PhysicsType != pmx.PHYSICS_TYPE_DYNAMIC || forceUpdate {
			// ボーン追従剛体・物理＋ボーン位置もしくは強制更新の場合のみ剛体位置更新
			boneTransform := bt.NewBtTransform()
			defer bt.DeleteBtTransform(boneTransform)
			mat := vds.Bones.Get(rigidBodyBone.Index).FilledGlobalMatrix().GL()
			boneTransform.SetFromOpenGLMatrix(&mat[0])

			mbt.UpdateTransform(modelIndex, modelPhysics, rigidBodyBone, boneTransform, rigidBody)
		}
	}

	return vds, MeshGlDeltas, VertexMorphIndexes, VertexMorphGlDeltas
}

func deformAfterPhysics(
	modelPhysics *mbt.MPhysics,
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	deltas *delta.VmdDeltas,
	selectedVertexIndexes, nextSelectedVertexIndexes []int,
	SelectedVertexDeltas *renderer.SelectedVertexMorphDeltas,
	frame int,
	enablePhysics, resetPhysics bool,
) (*delta.VmdDeltas, []mgl32.Mat4, []int, [][]float32) {
	if motion == nil {
		motion = vmd.NewVmdMotion("")
	}

	modelIndex := 0

	// 物理剛体位置を更新
	if enablePhysics || resetPhysics {
		// IKのON/OFF
		ikFrame := motion.IkFrames.Get(frame)

		for _, isAfterPhysics := range []bool{false, true} {
			for _, bone := range model.Bones.LayerSortedBones[isAfterPhysics] {
				if bone.RigidBody == nil || bone.RigidBody.PhysicsType == pmx.PHYSICS_TYPE_STATIC {
					continue
				}
				bonePhysicsGlobalMatrix := mbt.GetRigidBodyBoneMatrix(modelIndex, modelPhysics, bone.RigidBody)
				if deltas.Bones != nil && bonePhysicsGlobalMatrix != nil {
					bd := delta.NewBoneDeltaByGlobalMatrix(bone, frame,
						bonePhysicsGlobalMatrix, deltas.Bones.Get(bone.ParentIndex))
					deltas.Bones.Update(bd)
				}
			}
		}

		// 物理後のデフォーム情報
		deltas.Bones = deform.DeformByPhysicsFlag(motion.BoneFrames, frame, model, nil, true,
			deltas.Bones, deltas.Morphs, ikFrame, true)
	}

	// GL描画用データの作成
	BoneGlDeltas := make([]mgl32.Mat4, len(model.Bones.Data))
	for i, bone := range model.Bones.Data {
		delta := deltas.Bones.Get(bone.Index)
		if delta != nil {
			BoneGlDeltas[i] = delta.FilledLocalMatrix().GL()
		}
	}

	// 選択頂点モーフの設定は常に更新する
	SelectedVertexIndexesDeltas, SelectedVertexGlDeltas := renderer.SelectedVertexMorphDeltasGL(
		SelectedVertexDeltas, model, selectedVertexIndexes, nextSelectedVertexIndexes)

	return deltas, BoneGlDeltas, SelectedVertexIndexesDeltas, SelectedVertexGlDeltas
}

func Draw(
	modelPhysics *mbt.MPhysics,
	model *pmx.PmxModel,
	meshes *renderer.Meshes,
	shader *mgl.MShader,
	deltas *delta.VmdDeltas,
	invisibleMaterialIndexes, nextInvisibleMaterialIndexes []int,
	BoneGlDeltas []mgl32.Mat4,
	MeshGlDeltas []*renderer.MeshDelta,
	VertexMorphIndexes []int,
	VertexMorphGlDeltas [][]float32,
	SelectedVertexIndexesDeltas []int,
	SelectedVertexGlDeltasDeltas [][]float32,
	windowIndex int,
	isDrawNormal, isDrawWire, isDrawSelectedVertex bool,
	isDrawBones map[pmx.BoneFlag]bool,
	isDrawRigidBodyFront, visibleRigidBody, visibleJoint bool,
) *delta.VmdDeltas {
	vertexPositions := meshes.Draw(
		shader, BoneGlDeltas, VertexMorphIndexes, VertexMorphGlDeltas,
		SelectedVertexIndexesDeltas, SelectedVertexGlDeltasDeltas, MeshGlDeltas,
		invisibleMaterialIndexes, nextInvisibleMaterialIndexes, windowIndex,
		isDrawNormal, isDrawWire, isDrawSelectedVertex, isDrawBones, model.Bones)

	Vertices := renderer.NewVertexDeltas(model.Vertices)

	for i, pos := range vertexPositions {
		Vertices.Data[i] = renderer.NewVertexDelta(&mmath.MVec3{float64(-pos[0]), float64(pos[1]), float64(pos[2])})
	}

	// 物理デバッグ表示
	modelPhysics.DebugDrawWorld(visibleRigidBody, visibleJoint)

	return deltas
}
