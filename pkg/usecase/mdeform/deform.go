// 指示: miu200521358
package mdeform

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/physics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// DeformStage は変形ステージを表す。
type DeformStage int

const (
	// DEFORM_STAGE_BEFORE_PHYSICS は物理前変形。
	DEFORM_STAGE_BEFORE_PHYSICS DeformStage = iota
	// DEFORM_STAGE_FOR_PHYSICS は物理準備。
	DEFORM_STAGE_FOR_PHYSICS
	// DEFORM_STAGE_AFTER_PHYSICS は物理後変形。
	DEFORM_STAGE_AFTER_PHYSICS
)

// DeformOptions は変形オプションを表す。
type DeformOptions struct {
	TargetBoneNames []string
	EnableIK        bool
	EnablePhysics   bool
}

// BuildBeforePhysics は物理前の変形差分を構築する。
func BuildBeforePhysics(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	base *delta.VmdDeltas,
	frame motion.Frame,
	opts *DeformOptions,
) *delta.VmdDeltas {
	if modelData == nil {
		return base
	}
	motionHash := motionHash(motionData)
	if base != nil && base.Frame() == frame && base.ModelHash() == modelData.Hash() && base.MotionHash() == motionHash {
		return base
	}
	deltas := ensureVmdDeltas(modelData, motionHash, base, frame)
	deformNames := targetBoneNames(opts)
	includeIk := enableIk(opts)

	deltas.Morphs = deform.ComputeMorphDeltas(modelData, motionData, frame, nil)
	boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, frame, deformNames, includeIk, false, false)
	deltas.Bones = boneDeltas
	return deltas
}

// RebuildBeforePhysics は物理リセット用に物理前変形を再構築する。
func RebuildBeforePhysics(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	base *delta.VmdDeltas,
	frame motion.Frame,
	opts *DeformOptions,
) *delta.VmdDeltas {
	if modelData == nil {
		return base
	}
	motionHash := motionHash(motionData)
	deltas := ensureVmdDeltas(modelData, motionHash, base, frame)
	deltas.Bones = delta.NewBoneDeltas(modelData.Bones)
	deformNames := targetBoneNames(opts)
	includeIk := enableIk(opts)

	deltas.Morphs = deform.ComputeMorphDeltas(modelData, motionData, frame, nil)
	boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, frame, deformNames, includeIk, false, false)
	deltas.Bones = boneDeltas
	return deltas
}

// BuildForPhysics は物理用の剛体更新を行う。
func BuildForPhysics(
	core physics_api.IPhysicsCore,
	modelIndex int,
	modelData *model.PmxModel,
	deltas *delta.VmdDeltas,
	physicsDeltas *delta.PhysicsDeltas,
	enabled bool,
	resetType state.PhysicsResetType,
) *delta.VmdDeltas {
	if core == nil || modelData == nil || deltas == nil || deltas.Bones == nil {
		return deltas
	}
	if physicsDeltas != nil && physicsDeltas.RigidBodies != nil {
		updateRigidBodyShapeMass(core, modelIndex, modelData, physicsDeltas)
	}
	if modelData.RigidBodies == nil {
		return deltas
	}
	for _, rigidBody := range modelData.RigidBodies.Values() {
		if rigidBody == nil {
			continue
		}
		bone := boneByRigidBody(modelData, rigidBody)
		if bone == nil {
			continue
		}
		boneDelta := deltas.Bones.Get(bone.Index())
		if boneDelta == nil {
			continue
		}
		if (enabled && rigidBody.PhysicsType != model.PHYSICS_TYPE_DYNAMIC) || resetType != state.PHYSICS_RESET_TYPE_NONE {
			global := boneDelta.FilledGlobalMatrix()
			core.UpdateTransform(modelIndex, bone, &global, rigidBody)
		}
	}
	return deltas
}

// BuildAfterPhysics は物理結果を反映し、物理後変形を行う。
func BuildAfterPhysics(
	shared state.ISharedState,
	core physics_api.IPhysicsCore,
	modelIndex int,
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	deltas *delta.VmdDeltas,
	frame motion.Frame,
) *delta.VmdDeltas {
	if modelData == nil || deltas == nil || deltas.Bones == nil {
		return deltas
	}
	deltas.SetFrame(frame)
	dynamicBones := map[int]struct{}{}

	if core != nil && shared != nil && shared.HasFlag(state.STATE_FLAG_PHYSICS_ENABLED) {
		// 動的剛体の結果をボーンへ反映する。
		if modelData.RigidBodies != nil {
			for _, rigidBody := range modelData.RigidBodies.Values() {
				if rigidBody == nil || rigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
					continue
				}
				bone := boneByRigidBody(modelData, rigidBody)
				if bone == nil {
					continue
				}
				mat := core.GetRigidBodyBoneMatrix(modelIndex, rigidBody)
				if mat == nil {
					continue
				}
				parent := deltas.Bones.Get(bone.ParentIndex)
				bd := delta.NewBoneDeltaByGlobalMatrix(bone, frame, *mat, parent)
				deltas.Bones.Update(bd)
				dynamicBones[bone.Index()] = struct{}{}
			}
		}
	}

	// 物理後変形対象ボーンを再計算して反映する。
	afterNames := afterPhysicsBoneNames(modelData)
	if len(afterNames) > 0 {
		afterDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, frame, afterNames, true, true, false)
		mergeBoneDeltas(deltas.Bones, afterDeltas, dynamicBones)
	}

	// 既存のユニット行列を使いグローバル行列のみ更新する。
	deform.ApplyGlobalMatrices(modelData, deltas.Bones)
	return deltas
}

// ApplySkinning はスキニングを適用して頂点/法線を更新する。
func ApplySkinning(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, morphDeltas *delta.MorphDeltas) {
	if modelData == nil || modelData.Vertices == nil {
		return
	}
	deform.ApplySkinning(modelData.Vertices, boneDeltas, morphDeltas)
}

// ensureVmdDeltas は差分の再利用/再生成を行う。
func ensureVmdDeltas(modelData *model.PmxModel, motionHash string, base *delta.VmdDeltas, frame motion.Frame) *delta.VmdDeltas {
	if modelData == nil {
		return base
	}
	if base == nil {
		return delta.NewVmdDeltas(frame, modelData.Bones, modelData.Hash(), motionHash)
	}
	base.SetFrame(frame)
	base.SetModelHash(modelData.Hash())
	base.SetMotionHash(motionHash)
	base.Bones = delta.NewBoneDeltas(modelData.Bones)
	return base
}

// motionHash はモーションのハッシュ文字列を返す。
func motionHash(motionData *motion.VmdMotion) string {
	if motionData == nil {
		return ""
	}
	return motionData.Hash()
}

// targetBoneNames は変形対象ボーン名を返す。
func targetBoneNames(opts *DeformOptions) []string {
	if opts == nil {
		return nil
	}
	return opts.TargetBoneNames
}

// enableIk はIK有効フラグを返す。
func enableIk(opts *DeformOptions) bool {
	if opts == nil {
		return true
	}
	return opts.EnableIK
}

// updateRigidBodyShapeMass は剛体形状と質量を差分に応じて更新する。
func updateRigidBodyShapeMass(core physics_api.IPhysicsCore, modelIndex int, modelData *model.PmxModel, physicsDeltas *delta.PhysicsDeltas) {
	if core == nil || modelData == nil || physicsDeltas == nil || physicsDeltas.RigidBodies == nil {
		return
	}
	if modelData.RigidBodies == nil {
		return
	}
	for index, rigidBody := range modelData.RigidBodies.Values() {
		if rigidBody == nil {
			continue
		}
		rigidDelta := physicsDeltas.RigidBodies.Get(index)
		if rigidDelta == nil {
			continue
		}
		if !rigidDelta.Size.IsZero() || rigidDelta.Mass != 0.0 {
			core.UpdateRigidBodyShapeMass(modelIndex, rigidBody, rigidDelta)
		}
	}
}

// boneByRigidBody は剛体に対応するボーンを返す。
func boneByRigidBody(modelData *model.PmxModel, rigidBody *model.RigidBody) *model.Bone {
	if modelData == nil || rigidBody == nil || modelData.Bones == nil {
		return nil
	}
	if rigidBody.BoneIndex < 0 {
		return nil
	}
	bone, err := modelData.Bones.Get(rigidBody.BoneIndex)
	if err != nil {
		return nil
	}
	return bone
}

// afterPhysicsBoneNames は物理後変形対象ボーン名を収集する。
func afterPhysicsBoneNames(modelData *model.PmxModel) []string {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	names := make([]string, 0)
	for _, bone := range modelData.Bones.Values() {
		if bone == nil {
			continue
		}
		if bone.BoneFlag&model.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM == 0 {
			continue
		}
		names = append(names, bone.Name())
	}
	return names
}

// mergeBoneDeltas は差分を既存のボーン差分へ統合する。
func mergeBoneDeltas(dst *delta.BoneDeltas, src *delta.BoneDeltas, skip map[int]struct{}) {
	if dst == nil || src == nil {
		return
	}
	src.ForEach(func(index int, bd *delta.BoneDelta) bool {
		if bd == nil {
			return true
		}
		if skip != nil {
			if _, ok := skip[index]; ok {
				return true
			}
		}
		dst.Update(bd)
		return true
	})
}
