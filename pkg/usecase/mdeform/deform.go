// 指示: miu200521358
package mdeform

import (
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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

// dynamicBoneSyncMode は DYNAMIC_BONE 剛体の同期方式を表す。
type dynamicBoneSyncMode int

const (
	// DYNAMIC_BONE_SYNC_MODE_BULLET は Bullet 結果を優先し、物理前同期を行わない。
	DYNAMIC_BONE_SYNC_MODE_BULLET dynamicBoneSyncMode = iota
	// DYNAMIC_BONE_SYNC_MODE_FOLLOW_DELTA はボーン差分で剛体姿勢を追従更新する。
	DYNAMIC_BONE_SYNC_MODE_FOLLOW_DELTA
)

// defaultDynamicBoneSyncMode は通常時の DYNAMIC_BONE 同期方式。
// 見た目の基準復帰のため、既定は Bullet 結果優先（物理前同期なし）を使う。
const defaultDynamicBoneSyncMode = DYNAMIC_BONE_SYNC_MODE_BULLET

// playingDynamicBoneSyncMode は再生中の DYNAMIC_BONE 同期方式。
const playingDynamicBoneSyncMode = DYNAMIC_BONE_SYNC_MODE_FOLLOW_DELTA

// afterPhysicsVerboseTargetFrame は物理冗長ログの対象フレーム。
const afterPhysicsVerboseTargetFrame motion.Frame = 680

// DeformOptions は変形オプションを表す。
type DeformOptions struct {
	TargetBoneNames []string
	EnableIK        bool
	EnablePhysics   bool
	IkDebugFactory  deform.IIkDebugFactory
}

// rigidBodyBoneMatrixResult は剛体結果から得たボーン行列反映情報を保持する。
type rigidBodyBoneMatrixResult struct {
	rigidBody    *model.RigidBody
	bone         *model.Bone
	globalMatrix mmath.Mat4
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
	playing bool,
	resetType state.PhysicsResetType,
) *delta.VmdDeltas {
	if core == nil || modelData == nil || deltas == nil || deltas.Bones == nil {
		return deltas
	}
	dynamicBoneMode := resolveDynamicBoneSyncMode(playing)
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
		global := boneDelta.FilledGlobalMatrix()
		if resetType != state.PHYSICS_RESET_TYPE_NONE {
			core.UpdateTransform(modelIndex, bone, &global, rigidBody)
			continue
		}
		if !enabled {
			continue
		}
		switch rigidBody.PhysicsType {
		case model.PHYSICS_TYPE_STATIC:
			core.UpdateTransform(modelIndex, bone, &global, rigidBody)
		case model.PHYSICS_TYPE_DYNAMIC_BONE:
			syncDynamicBoneForPhysics(core, dynamicBoneMode, modelIndex, bone, &global, rigidBody)
		case model.PHYSICS_TYPE_DYNAMIC:
		}
	}
	return deltas
}

// resolveDynamicBoneSyncMode は再生状態に応じた DYNAMIC_BONE 同期方式を返す。
func resolveDynamicBoneSyncMode(playing bool) dynamicBoneSyncMode {
	if playing {
		return playingDynamicBoneSyncMode
	}
	return defaultDynamicBoneSyncMode
}

// syncDynamicBoneForPhysics は DYNAMIC_BONE の通常時同期を方式ごとに実行する。
func syncDynamicBoneForPhysics(
	core physics.IPhysicsCore,
	mode dynamicBoneSyncMode,
	modelIndex int,
	bone *model.Bone,
	global *mmath.Mat4,
	rigidBody *model.RigidBody,
) {
	if core == nil || bone == nil || global == nil || rigidBody == nil {
		return
	}
	switch mode {
	case DYNAMIC_BONE_SYNC_MODE_FOLLOW_DELTA:
		core.FollowDeltaTransform(modelIndex, bone, global, rigidBody)
	default:
		// 現状の既定は Bullet 結果優先。同期は行わず、物理後反映のみで追従させる。
	}
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
	logger := afterPhysicsVerboseLogger()
	results := map[int]rigidBodyBoneMatrixResult{}
	order := make([]int, 0)
	if core != nil && physicsEnabled {
		// 動的剛体の結果をボーンへ反映する。
		results = collectRigidBodyBoneMatrices(core, modelIndex, modelData)
		order = resolveRigidBodyBoneUpdateOrder(results)
		applyRigidBodyBoneMatrices(deltas, frame, results, order)
	}
	logAfterPhysicsStage(logger, "物理直後", modelIndex, frame, modelData, deltas.Bones, order)

	// 物理後変形対象ボーンを再計算して反映する。
	afterNames := afterPhysicsBoneNames(modelData)
	var afterIndexes []int
	if len(afterNames) > 0 {
		afterIndexes = deform.UpdateAfterPhysicsBoneDeltas(modelData, motionData, deltas.Bones, frame, afterNames)
	}
	rotationGrantIndexes := collectExternalRotationBoneIndexes(modelData, afterIndexes)
	logAfterPhysicsStage(logger, "回転付与", modelIndex, frame, modelData, deltas.Bones, rotationGrantIndexes)

	// 既存のユニット行列を使い物理後変形対象のグローバル行列のみ更新する。
	if len(afterIndexes) > 0 {
		deform.ApplyGlobalMatricesWithIndexes(modelData, deltas.Bones, afterIndexes)
	}
	logAfterPhysicsStage(logger, "物理後変形", modelIndex, frame, modelData, deltas.Bones, afterIndexes)
	finalIndexes := mergeUniqueBoneIndexes(afterIndexes, order)
	logAfterPhysicsStage(logger, "最終結果", modelIndex, frame, modelData, deltas.Bones, finalIndexes)
	return deltas
}

// collectRigidBodyBoneMatrices は剛体シミュレーション結果のボーン行列を収集する。
func collectRigidBodyBoneMatrices(
	core physics.IPhysicsCore,
	modelIndex int,
	modelData *model.PmxModel,
) map[int]rigidBodyBoneMatrixResult {
	results := map[int]rigidBodyBoneMatrixResult{}
	if core == nil || modelData == nil || modelData.RigidBodies == nil {
		return results
	}
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
		// 同一ボーンへ複数剛体が紐づく場合は従来互換で後勝ちにする。
		results[bone.Index()] = rigidBodyBoneMatrixResult{
			rigidBody:    rigidBody,
			bone:         bone,
			globalMatrix: *mat,
		}
	}
	return results
}

// resolveRigidBodyBoneUpdateOrder は親剛体ボーンを先行させる反映順を返す。
func resolveRigidBodyBoneUpdateOrder(results map[int]rigidBodyBoneMatrixResult) []int {
	if len(results) == 0 {
		return nil
	}
	keys := make([]int, 0, len(results))
	for index := range results {
		keys = append(keys, index)
	}
	sort.Slice(keys, func(i, j int) bool {
		left := results[keys[i]].bone
		right := results[keys[j]].bone
		leftLayer := 0
		rightLayer := 0
		if left != nil {
			leftLayer = left.Layer
		}
		if right != nil {
			rightLayer = right.Layer
		}
		if leftLayer == rightLayer {
			return keys[i] < keys[j]
		}
		return leftLayer < rightLayer
	})

	order := make([]int, 0, len(keys))
	visited := map[int]bool{}
	visiting := map[int]bool{}
	var visit func(index int)
	visit = func(index int) {
		if visited[index] {
			return
		}
		if visiting[index] {
			return
		}
		visiting[index] = true
		result, ok := results[index]
		if ok && result.bone != nil {
			if _, exists := results[result.bone.ParentIndex]; exists {
				visit(result.bone.ParentIndex)
			}
		}
		visiting[index] = false
		visited[index] = true
		order = append(order, index)
	}
	for _, index := range keys {
		visit(index)
	}
	return order
}

// applyRigidBodyBoneMatrices は反映順に従って剛体行列をボーン差分へ適用する。
func applyRigidBodyBoneMatrices(
	deltas *delta.VmdDeltas,
	frame motion.Frame,
	results map[int]rigidBodyBoneMatrixResult,
	order []int,
) int {
	if deltas == nil || deltas.Bones == nil || len(results) == 0 || len(order) == 0 {
		return 0
	}
	reflectedCount := 0
	for _, boneIndex := range order {
		result, ok := results[boneIndex]
		if !ok || result.bone == nil {
			continue
		}
		parent := deltas.Bones.Get(result.bone.ParentIndex)
		bd := delta.NewBoneDeltaByGlobalMatrix(result.bone, frame, result.globalMatrix, parent)
		if bd == nil {
			continue
		}
		deltas.Bones.Update(bd)
		reflectedCount++
	}
	return reflectedCount
}

// afterPhysicsVerboseLogger は物理冗長ログ出力用の logger を返す。
func afterPhysicsVerboseLogger() logging.ILogger {
	logger := logging.DefaultLogger()
	if logger == nil || !logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		return nil
	}
	return logger
}

// logAfterPhysicsStage は物理後変形の段階別ボーン情報を冗長出力する。
func logAfterPhysicsStage(
	logger logging.ILogger,
	stage string,
	modelIndex int,
	frame motion.Frame,
	modelData *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	indexes []int,
) {
	if logger == nil || modelData == nil || boneDeltas == nil || !isAfterPhysicsVerboseTargetFrame(frame) {
		return
	}
	targetIndexes := filterAfterPhysicsVerboseIndexes(boneDeltas, indexes)
	if len(targetIndexes) == 0 {
		return
	}
	logger.Verbose(
		logging.VERBOSE_INDEX_PHYSICS,
		"物理段階: stage=%s model=%d frame=%v bones=%d",
		stage,
		modelIndex,
		frame,
		len(targetIndexes),
	)
	for _, boneIndex := range targetIndexes {
		boneDelta := boneDeltas.Get(boneIndex)
		if boneDelta == nil || boneDelta.Bone == nil {
			continue
		}
		bone := boneDelta.Bone
		parentName := boneNameByIndex(modelData, bone.ParentIndex)
		effectName := boneNameByIndex(modelData, bone.EffectIndex)
		isAfterPhysics := bone.BoneFlag&model.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM != 0
		isExternalRotation := bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_ROTATION != 0 && bone.EffectIndex >= 0
		globalPos := boneDelta.FilledGlobalPosition().String()
		globalRot := boneDelta.FilledGlobalMatrix().Quaternion().ToMMDDegrees().String()
		totalRot := boneDelta.FilledTotalRotation().ToMMDDegrees().String()
		framePos := boneDelta.FilledFramePosition().String()
		frameRot := boneDelta.FilledFrameRotation().ToMMDDegrees().String()
		effectGlobalRot, effectTotalRot := afterPhysicsEffectRotationInfo(boneDeltas, bone.EffectIndex)
		logger.Verbose(
			logging.VERBOSE_INDEX_PHYSICS,
			"物理段階詳細: stage=%s model=%d frame=%v bone=%s(%d) parent=%s(%d) effect=%s(%d,%.3f) flags=0x%04x afterPhysics=%t extRot=%t framePos=%s frameRotDeg=%s globalPos=%s globalRotDeg=%s totalRotDeg=%s effectGlobalRotDeg=%s effectTotalRotDeg=%s",
			stage,
			modelIndex,
			frame,
			bone.Name(),
			bone.Index(),
			parentName,
			bone.ParentIndex,
			effectName,
			bone.EffectIndex,
			bone.EffectFactor,
			int(bone.BoneFlag),
			isAfterPhysics,
			isExternalRotation,
			framePos,
			frameRot,
			globalPos,
			globalRot,
			totalRot,
			effectGlobalRot,
			effectTotalRot,
		)
	}
}

// isAfterPhysicsVerboseTargetFrame は物理冗長ログの対象フレームか判定する。
func isAfterPhysicsVerboseTargetFrame(frame motion.Frame) bool {
	return frame == afterPhysicsVerboseTargetFrame
}

// afterPhysicsEffectRotationInfo は付与元ボーンの回転情報を返す。
func afterPhysicsEffectRotationInfo(boneDeltas *delta.BoneDeltas, effectIndex int) (string, string) {
	if boneDeltas == nil || effectIndex < 0 {
		return "-", "-"
	}
	effectDelta := boneDeltas.Get(effectIndex)
	if effectDelta == nil {
		return "-", "-"
	}
	globalRot := effectDelta.FilledGlobalMatrix().Quaternion().ToMMDDegrees().String()
	totalRot := effectDelta.FilledTotalRotation().ToMMDDegrees().String()
	return globalRot, totalRot
}

// filterAfterPhysicsVerboseIndexes は物理冗長ログ対象のボーンindexへ絞り込む。
func filterAfterPhysicsVerboseIndexes(boneDeltas *delta.BoneDeltas, indexes []int) []int {
	if boneDeltas == nil || len(indexes) == 0 {
		return nil
	}
	out := make([]int, 0, len(indexes))
	for _, boneIndex := range indexes {
		boneDelta := boneDeltas.Get(boneIndex)
		if boneDelta == nil || boneDelta.Bone == nil {
			continue
		}
		if !isAfterPhysicsVerboseTargetBone(boneDelta.Bone) {
			continue
		}
		out = append(out, boneIndex)
	}
	return out
}

// isAfterPhysicsVerboseTargetBone は物理冗長ログに出力する対象ボーンか判定する。
func isAfterPhysicsVerboseTargetBone(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return strings.Contains(bone.Name(), "右袖")
}

// collectExternalRotationBoneIndexes は回転付与ボーンの index を返す。
func collectExternalRotationBoneIndexes(modelData *model.PmxModel, indexes []int) []int {
	if modelData == nil || len(indexes) == 0 {
		return nil
	}
	out := make([]int, 0, len(indexes))
	for _, boneIndex := range indexes {
		bone, err := modelData.Bones.Get(boneIndex)
		if err != nil || bone == nil {
			continue
		}
		if bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_ROTATION == 0 || bone.EffectIndex < 0 {
			continue
		}
		out = append(out, boneIndex)
	}
	return out
}

// mergeUniqueBoneIndexes は index 群を重複排除しつつ連結する。
func mergeUniqueBoneIndexes(primary []int, secondary []int) []int {
	out := make([]int, 0, len(primary)+len(secondary))
	used := map[int]struct{}{}
	for _, index := range primary {
		if _, exists := used[index]; exists {
			continue
		}
		used[index] = struct{}{}
		out = append(out, index)
	}
	for _, index := range secondary {
		if _, exists := used[index]; exists {
			continue
		}
		used[index] = struct{}{}
		out = append(out, index)
	}
	return out
}

// boneNameByIndex は index に対応するボーン名を返す。
func boneNameByIndex(modelData *model.PmxModel, index int) string {
	if modelData == nil || modelData.Bones == nil || index < 0 {
		return "-"
	}
	bone, err := modelData.Bones.Get(index)
	if err != nil || bone == nil {
		return "-"
	}
	return bone.Name()
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
