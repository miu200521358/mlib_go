// 指示: miu200521358
package mdeform

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/mlib_go/pkg/usecase/port/physics"
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
	IkDebugFactory  deform.IIkDebugFactory
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
	debugFactory := ikDebugFactory(opts)

	deltas.Morphs = deform.ComputeMorphDeltas(modelData, motionData, frame, nil)
	boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, frame, deformNames, includeIk, false, false, debugFactory)
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
	debugFactory := ikDebugFactory(opts)

	deltas.Morphs = deform.ComputeMorphDeltas(modelData, motionData, frame, nil)
	boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, frame, deformNames, includeIk, false, false, debugFactory)
	deltas.Bones = boneDeltas
	return deltas
}

// BuildForPhysics は物理用の剛体更新を行う。
func BuildForPhysics(
	core physics.IPhysicsCore,
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
	frame := deltas.Frame()
	logSummary := shouldEmitPhysicsVerificationSummary(frame)
	logger := logging.DefaultLogger()
	counters := physicsSyncCounters{}
	if physicsDeltas != nil && physicsDeltas.RigidBodies != nil {
		updateRigidBodyShapeMass(core, modelIndex, modelData, physicsDeltas)
	}
	if modelData.RigidBodies == nil {
		if logSummary && logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
			logger.Verbose(
				logging.VERBOSE_INDEX_PHYSICS,
				"物理検証同期要約: model=%d frame=%v rigidBodies=0",
				modelIndex,
				frame,
			)
		}
		return deltas
	}
	for _, rigidBody := range modelData.RigidBodies.Values() {
		if rigidBody == nil {
			continue
		}
		counters.Total++
		switch rigidBody.PhysicsType {
		case model.PHYSICS_TYPE_STATIC:
			counters.Static++
		case model.PHYSICS_TYPE_DYNAMIC:
			counters.Dynamic++
		case model.PHYSICS_TYPE_DYNAMIC_BONE:
			counters.DynamicBone++
		}
		bone := boneByRigidBody(modelData, rigidBody)
		if bone == nil {
			counters.MissingBone++
			continue
		}
		boneDelta := deltas.Bones.Get(bone.Index())
		if boneDelta == nil {
			counters.MissingBoneDelta++
			continue
		}
		global := boneDelta.FilledGlobalMatrix()
		if resetType != state.PHYSICS_RESET_TYPE_NONE {
			counters.HardSyncByReset++
			core.UpdateTransform(modelIndex, bone, &global, rigidBody)
			continue
		}
		if !enabled {
			counters.SkippedByPhysicsDisabled++
			continue
		}
		switch rigidBody.PhysicsType {
		case model.PHYSICS_TYPE_STATIC:
			counters.HardSyncStatic++
			core.UpdateTransform(modelIndex, bone, &global, rigidBody)
		case model.PHYSICS_TYPE_DYNAMIC_BONE:
			// 旧 mlib と同等に、動的ボーン追従剛体も毎フレームの姿勢同期はハード更新で扱う。
			counters.HardSyncDynamicBone++
			core.UpdateTransform(modelIndex, bone, &global, rigidBody)
		case model.PHYSICS_TYPE_DYNAMIC:
			counters.SkippedDynamic++
		}
	}
	if logSummary && logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		logger.Verbose(
			logging.VERBOSE_INDEX_PHYSICS,
			"物理検証同期要約: model=%d frame=%v enabled=%t resetType=%d total=%d static=%d dynamic=%d dynamicBone=%d hardStatic=%d hardDynamicBone=%d hardReset=%d skipDynamic=%d skipDisabled=%d missingBone=%d missingBoneDelta=%d",
			modelIndex,
			frame,
			enabled,
			resetType,
			counters.Total,
			counters.Static,
			counters.Dynamic,
			counters.DynamicBone,
			counters.HardSyncStatic,
			counters.HardSyncDynamicBone,
			counters.HardSyncByReset,
			counters.SkippedDynamic,
			counters.SkippedByPhysicsDisabled,
			counters.MissingBone,
			counters.MissingBoneDelta,
		)
	}
	return deltas
}

// BuildAfterPhysics は物理結果を反映し、物理後変形を行う。
func BuildAfterPhysics(
	core physics.IPhysicsCore,
	physicsEnabled bool,
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
	if core != nil && physicsEnabled {
		// 動的剛体の結果をボーンへ反映する。
		if modelData.RigidBodies != nil {
			rigidBodies := modelData.RigidBodies.Values()
			for _, rigidBody := range rigidBodies {
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
			}
		}
	}

	// 物理後変形対象ボーンを再計算して反映する。
	afterNames := afterPhysicsBoneNames(modelData)
	var afterIndexes []int
	if len(afterNames) > 0 {
		afterIndexes = deform.UpdateAfterPhysicsBoneDeltas(modelData, motionData, deltas.Bones, frame, afterNames)
	}

	// 既存のユニット行列を使い物理後変形対象のグローバル行列のみ更新する。
	if len(afterIndexes) > 0 {
		deform.ApplyGlobalMatricesWithIndexes(modelData, deltas.Bones, afterIndexes)
	}
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

// ikDebugFactory はIKデバッグ出力I/Fを返す。
func ikDebugFactory(opts *DeformOptions) deform.IIkDebugFactory {
	if opts == nil {
		return nil
	}
	return opts.IkDebugFactory
}

// updateRigidBodyShapeMass は剛体形状と質量を差分に応じて更新する。
func updateRigidBodyShapeMass(core physics.IPhysicsCore, modelIndex int, modelData *model.PmxModel, physicsDeltas *delta.PhysicsDeltas) {
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
		core.UpdateRigidBodyShapeMass(modelIndex, rigidBody, rigidDelta)
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

// physicsSyncCounters は物理同期経路の集計値を保持する。
type physicsSyncCounters struct {
	Total                    int
	Static                   int
	Dynamic                  int
	DynamicBone              int
	HardSyncStatic           int
	HardSyncDynamicBone      int
	HardSyncByReset          int
	SkippedDynamic           int
	SkippedByPhysicsDisabled int
	MissingBone              int
	MissingBoneDelta         int
}

// shouldEmitPhysicsVerificationSummary は物理検証要約ログを出力すべきフレームか判定する。
func shouldEmitPhysicsVerificationSummary(frame motion.Frame) bool {
	rounded := motion.Frame(math.Round(float64(frame)))
	if math.Abs(float64(frame-rounded)) > 1e-3 {
		return false
	}
	frameNumber := int(rounded)
	if frameNumber < 0 {
		return false
	}
	return frameNumber == 0 || frameNumber%30 == 0
}
