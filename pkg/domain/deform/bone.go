// 指示: miu200521358
package deform

import (
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

const maxEffectorRecursion = 10

// ComputeBoneDeltas はボーン差分を算出して返す。
func ComputeBoneDeltas(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	frame motion.Frame,
	boneNames []string,
	includeIk bool,
	afterPhysics bool,
	removeTwist bool,
) (*delta.BoneDeltas, []int) {
	if modelData == nil {
		return delta.NewBoneDeltas(nil), nil
	}
	boneDeltas := delta.NewBoneDeltas(modelData.Bones)
	logger := logging.DefaultLogger()
	logDetail := logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) && afterPhysics && frame == 0
	if logDetail {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
			"物理後変形差分計算: frame=%v boneNames=%d includeIk=%t",
			frame,
			len(boneNames),
			includeIk,
		)
		if len(boneNames) > 0 {
			logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
				"物理後変形差分計算の対象名: %s",
				strings.Join(boneNames, ","),
			)
		}
	}
	var deformBoneIndexes []int
	if afterPhysics {
		deformBoneIndexes = collectAfterPhysicsBoneIndexes(modelData, boneNames)
	} else {
		deformBoneIndexes = collectBoneIndexes(modelData, boneNames, includeIk, afterPhysics)
	}
	if logDetail {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
			"物理後変形差分計算の対象index: count=%d indexes=%v",
			len(deformBoneIndexes),
			deformBoneIndexes,
		)
	}
	boneMorphDeltas := computeBoneMorphDeltas(modelData, motionData, frame, nil)

	for _, boneIndex := range deformBoneIndexes {
		bone, err := modelData.Bones.Get(boneIndex)
		if err != nil || bone == nil {
			continue
		}
		d := boneDeltas.Get(boneIndex)
		if d == nil {
			d = delta.NewBoneDelta(bone, frame)
		}
		bf := getBoneFrame(motionData, bone.Name(), frame)
		if bf != nil {
			if bf.Position != nil {
				d.FramePosition = bf.Position
			}
			if bf.Rotation != nil {
				d.FrameRotation = bf.Rotation
			}
			if bf.CancelablePosition != nil {
				d.FrameCancelablePosition = bf.CancelablePosition
			}
			if bf.CancelableRotation != nil {
				d.FrameCancelableRotation = bf.CancelableRotation
			}
			if bf.Scale != nil {
				d.FrameScale = bf.Scale
			}
			if bf.CancelableScale != nil {
				d.FrameCancelableScale = bf.CancelableScale
			}
		}

		if boneMorphDeltas != nil {
			morphDelta := boneMorphDeltas.Get(boneIndex)
			if morphDelta != nil {
				if morphDelta.FramePosition != nil {
					d.FrameMorphPosition = morphDelta.FramePosition
				}
				if morphDelta.FrameRotation != nil {
					d.FrameMorphRotation = morphDelta.FrameRotation
				}
				if morphDelta.FrameCancelablePosition != nil {
					d.FrameMorphCancelablePosition = morphDelta.FrameCancelablePosition
				}
				if morphDelta.FrameCancelableRotation != nil {
					d.FrameMorphCancelableRotation = morphDelta.FrameCancelableRotation
				}
				if morphDelta.FrameScale != nil {
					d.FrameMorphScale = morphDelta.FrameScale
				}
				if morphDelta.FrameCancelableScale != nil {
					d.FrameMorphCancelableScale = morphDelta.FrameCancelableScale
				}
				if morphDelta.FrameLocalMat != nil {
					d.FrameLocalMorphMat = morphDelta.FrameLocalMat
				}
			}
		}

		d.InvalidateTotals()
		boneDeltas.Update(d)
	}

	if includeIk {
		applyBoneMatricesWithIndexes(modelData, boneDeltas, deformBoneIndexes)
		applyIkDeltas(modelData, motionData, boneDeltas, frame, deformBoneIndexes, removeTwist)
		applyBoneMatricesWithIndexes(modelData, boneDeltas, deformBoneIndexes)
	}

	return boneDeltas, deformBoneIndexes
}

// ApplyBoneMatrices はボーン行列を合成して差分へ反映する。
func ApplyBoneMatrices(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	indexes := sortedBoneIndexes(modelData, boneDeltas)
	applyBoneMatricesWithIndexes(modelData, boneDeltas, indexes)
}

// ApplyBoneMatricesWithIndexes は指定インデックス順でボーン行列を合成する。
func ApplyBoneMatricesWithIndexes(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, indexes []int) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	applyBoneMatricesWithIndexes(modelData, boneDeltas, indexes)
}

// ApplyGlobalMatrices はグローバル行列のみ更新する。
func ApplyGlobalMatrices(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	indexes := sortedBoneIndexes(modelData, boneDeltas)
	applyGlobalMatricesWithIndexes(modelData, boneDeltas, indexes)
}

// ApplyGlobalMatricesWithIndexes は指定インデックス順でグローバル行列のみ更新する。
func ApplyGlobalMatricesWithIndexes(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, indexes []int) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	applyGlobalMatricesWithIndexes(modelData, boneDeltas, indexes)
}

// applyBoneMatricesWithIndexes は指定インデックス順でボーン行列を合成する。
func applyBoneMatricesWithIndexes(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, indexes []int) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	for _, boneIndex := range indexes {
		d := boneDeltas.Get(boneIndex)
		if d == nil || d.Bone == nil {
			continue
		}
		updateBoneDelta(modelData, boneDeltas, d)
		applyGlobalMatrix(boneDeltas, d)
	}
}

// applyGlobalMatricesWithIndexes は指定インデックス順でグローバル行列のみ更新する。
func applyGlobalMatricesWithIndexes(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, indexes []int) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	for _, boneIndex := range indexes {
		d := boneDeltas.Get(boneIndex)
		if d == nil || d.Bone == nil {
			continue
		}
		if d.UnitMatrix == nil {
			updateBoneDelta(modelData, boneDeltas, d)
		}
		applyGlobalMatrixNoLocal(boneDeltas, d)
	}
}

// collectBoneIndexes は変形対象ボーンのindexを収集する。
func collectBoneIndexes(
	modelData *model.PmxModel,
	boneNames []string,
	includeIk bool,
	afterPhysics bool,
) []int {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	if len(boneNames) == 0 {
		return collectAllBones(modelData, afterPhysics)
	}
	effectorChildren, ikTargets, ikLinks := buildReverseBoneRelations(modelData)
	indexes := make(map[int]struct{})
	for _, name := range boneNames {
		bone, err := modelData.Bones.GetByName(name)
		if err != nil || bone == nil {
			continue
		}
		addRelatedBoneIndexes(modelData, bone, includeIk, indexes, effectorChildren, ikTargets, ikLinks)
	}
	out := make([]int, 0, len(indexes))
	for idx := range indexes {
		out = append(out, idx)
	}
	sortBoneIndexes(modelData, out)
	return out
}

// buildReverseBoneRelations は逆引き関係を構築する。
func buildReverseBoneRelations(modelData *model.PmxModel) (map[int][]int, map[int][]int, map[int][]int) {
	effectorChildren := map[int][]int{}
	ikTargets := map[int][]int{}
	ikLinks := map[int][]int{}
	if modelData == nil || modelData.Bones == nil {
		return effectorChildren, ikTargets, ikLinks
	}
	for _, bone := range modelData.Bones.Values() {
		if bone == nil {
			continue
		}
		if boneHasEffector(bone) && bone.EffectIndex >= 0 {
			effectorChildren[bone.EffectIndex] = append(effectorChildren[bone.EffectIndex], bone.Index())
		}
		if boneIsIk(bone) && bone.Ik != nil {
			if bone.Ik.BoneIndex >= 0 {
				ikTargets[bone.Ik.BoneIndex] = append(ikTargets[bone.Ik.BoneIndex], bone.Index())
			}
			for _, link := range bone.Ik.Links {
				if link.BoneIndex < 0 {
					continue
				}
				ikLinks[link.BoneIndex] = append(ikLinks[link.BoneIndex], bone.Index())
			}
		}
	}
	return effectorChildren, ikTargets, ikLinks
}

// buildChildRelations は親子関係の逆引きを構築する。
func buildChildRelations(modelData *model.PmxModel) [][]int {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	children := make([][]int, modelData.Bones.Len())
	for _, bone := range modelData.Bones.Values() {
		if bone == nil || bone.ParentIndex < 0 {
			continue
		}
		if bone.ParentIndex >= len(children) {
			continue
		}
		children[bone.ParentIndex] = append(children[bone.ParentIndex], bone.Index())
	}
	return children
}

// buildEffectorChildrenSlice は付与関係の逆引きを構築する。
func buildEffectorChildrenSlice(modelData *model.PmxModel) [][]int {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	children := make([][]int, modelData.Bones.Len())
	for _, bone := range modelData.Bones.Values() {
		if bone == nil || !boneHasEffector(bone) || bone.EffectIndex < 0 {
			continue
		}
		if bone.EffectIndex >= len(children) {
			continue
		}
		children[bone.EffectIndex] = append(children[bone.EffectIndex], bone.Index())
	}
	return children
}

// addRelatedBoneIndexes は関連ボーンを再帰的に追加する。
func addRelatedBoneIndexes(
	modelData *model.PmxModel,
	bone *model.Bone,
	includeIk bool,
	indexes map[int]struct{},
	effectorChildren map[int][]int,
	ikTargets map[int][]int,
	ikLinks map[int][]int,
) {
	if modelData == nil || bone == nil {
		return
	}
	queue := []int{bone.Index()}
	visited := map[int]struct{}{}
	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]
		if _, ok := visited[idx]; ok {
			continue
		}
		visited[idx] = struct{}{}
		current, err := modelData.Bones.Get(idx)
		if err != nil || current == nil {
			continue
		}
		if current.ParentIndex >= 0 && modelData.Bones.Contains(current.ParentIndex) {
			queue = append(queue, current.ParentIndex)
		}
		if boneHasEffector(current) && current.EffectIndex >= 0 && modelData.Bones.Contains(current.EffectIndex) {
			queue = append(queue, current.EffectIndex)
		}
		if includeIk && boneIsIk(current) && current.Ik != nil {
			if current.Ik.BoneIndex >= 0 && modelData.Bones.Contains(current.Ik.BoneIndex) {
				queue = append(queue, current.Ik.BoneIndex)
			}
			for _, link := range current.Ik.Links {
				if link.BoneIndex >= 0 && modelData.Bones.Contains(link.BoneIndex) {
					queue = append(queue, link.BoneIndex)
				}
			}
		}
		if children := effectorChildren[idx]; len(children) > 0 {
			queue = append(queue, children...)
		}
		if includeIk {
			if targets := ikTargets[idx]; len(targets) > 0 {
				queue = append(queue, targets...)
			}
			if links := ikLinks[idx]; len(links) > 0 {
				queue = append(queue, links...)
			}
		}
	}
	for idx := range visited {
		indexes[idx] = struct{}{}
	}
}

// collectIkChainIndexes はIKボーンに関連するインデックスを収集する。
func collectIkChainIndexes(modelData *model.PmxModel, ikBone *model.Bone) []int {
	if modelData == nil || ikBone == nil {
		return nil
	}
	indexes := make(map[int]struct{})
	addBoneWithDependencies(modelData, ikBone, true, indexes)
	out := make([]int, 0, len(indexes))
	for idx := range indexes {
		out = append(out, idx)
	}
	sortBoneIndexes(modelData, out)
	return out
}

// collectDescendantFlags は開始ボーン以下をフラグ化する。
func collectDescendantFlags(children [][]int, start int, flags []bool, updated *[]int, queue *[]int) {
	if len(children) == 0 || len(flags) == 0 || start < 0 || start >= len(children) {
		return
	}
	q := queue
	if q == nil {
		return
	}
	*q = append((*q)[:0], start)
	for len(*q) > 0 {
		idx := (*q)[0]
		*q = (*q)[1:]
		if idx < 0 || idx >= len(flags) {
			continue
		}
		if flags[idx] {
			continue
		}
		flags[idx] = true
		if updated != nil {
			*updated = append(*updated, idx)
		}
		if next := children[idx]; len(next) > 0 {
			*q = append(*q, next...)
		}
	}
}

// buildRecalcIndexesByOrderFlags は変形順で再計算対象のindexを整列する。
func buildRecalcIndexesByOrderFlags(orderByRank []int, updatedFlags []bool, out []int) []int {
	if len(orderByRank) == 0 || len(updatedFlags) == 0 {
		return out[:0]
	}
	out = out[:0]
	for _, boneIndex := range orderByRank {
		if boneIndex < 0 || boneIndex >= len(updatedFlags) {
			continue
		}
		if updatedFlags[boneIndex] {
			out = append(out, boneIndex)
		}
	}
	return out
}

// collectEffectorRelatedFlags は付与関係で影響するボーンをフラグ化する。
func collectEffectorRelatedFlags(
	children [][]int,
	start int,
	flags []bool,
	updated *[]int,
	queue *[]int,
) {
	if len(children) == 0 || len(flags) == 0 || start < 0 || start >= len(children) {
		return
	}
	q := queue
	if q == nil {
		return
	}
	*q = append((*q)[:0], start)
	for len(*q) > 0 {
		idx := (*q)[0]
		*q = (*q)[1:]
		if idx < 0 || idx >= len(flags) {
			continue
		}
		if flags[idx] {
			continue
		}
		flags[idx] = true
		if updated != nil {
			*updated = append(*updated, idx)
		}
		if next := children[idx]; len(next) > 0 {
			*q = append(*q, next...)
		}
	}
}

// resetFlagsByList はlistで指定されたフラグを解除する。
func resetFlagsByList(flags []bool, list []int) {
	for _, idx := range list {
		flags[idx] = false
	}
}

// collectAllBones は全ボーンindexを収集する。
func collectAllBones(modelData *model.PmxModel, afterPhysics bool) []int {
	values := modelData.Bones.Values()
	indexes := make(map[int]struct{}, len(values))
	for _, bone := range values {
		if bone == nil {
			continue
		}
		if afterPhysics {
			if boneIsAfterPhysics(bone) {
				addBoneAndParents(modelData, bone, indexes)
			}
			continue
		}
		indexes[bone.Index()] = struct{}{}
	}
	out := make([]int, 0, len(indexes))
	for idx := range indexes {
		out = append(out, idx)
	}
	sortBoneIndexes(modelData, out)
	return out
}

// collectAfterPhysicsBoneIndexes は物理後変形対象のボーンindexを収集する。
func collectAfterPhysicsBoneIndexes(modelData *model.PmxModel, boneNames []string) []int {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	indexes := make([]int, 0)
	if len(boneNames) == 0 {
		for _, bone := range modelData.Bones.Values() {
			if bone == nil || !boneIsAfterPhysics(bone) {
				continue
			}
			indexes = append(indexes, bone.Index())
		}
		sortBoneIndexes(modelData, indexes)
		return indexes
	}
	for _, name := range boneNames {
		bone, err := modelData.Bones.GetByName(name)
		if err != nil || bone == nil || !boneIsAfterPhysics(bone) {
			continue
		}
		indexes = append(indexes, bone.Index())
	}
	sortBoneIndexes(modelData, indexes)
	return indexes
}

// addBoneWithDependencies は親や付与元を含めて追加する。
func addBoneWithDependencies(
	modelData *model.PmxModel,
	bone *model.Bone,
	includeIk bool,
	indexes map[int]struct{},
) {
	if modelData == nil || bone == nil {
		return
	}
	addBoneAndParents(modelData, bone, indexes)
	if boneHasEffector(bone) {
		if effector, err := modelData.Bones.Get(bone.EffectIndex); err == nil && effector != nil {
			addBoneAndParents(modelData, effector, indexes)
		}
	}
	if includeIk && boneIsIk(bone) {
		addIkDependencies(modelData, bone, indexes)
	}
}

// addRelatedBonesByTargets は選択済みボーンに紐づく関連ボーンを追加する。
func addRelatedBonesByTargets(modelData *model.PmxModel, includeIk bool, indexes map[int]struct{}) {
	if modelData == nil || modelData.Bones == nil {
		return
	}
	for {
		added := false
		for _, bone := range modelData.Bones.Values() {
			if bone == nil {
				continue
			}
			if _, ok := indexes[bone.Index()]; ok {
				continue
			}
			if includeIk && boneIsIk(bone) && ikHasTargetInIndexes(bone.Ik, indexes) {
				prevLen := len(indexes)
				addBoneWithDependencies(modelData, bone, includeIk, indexes)
				if len(indexes) != prevLen {
					added = true
				}
				continue
			}
			if boneHasEffector(bone) {
				if _, ok := indexes[bone.EffectIndex]; ok {
					prevLen := len(indexes)
					addBoneWithDependencies(modelData, bone, includeIk, indexes)
					if len(indexes) != prevLen {
						added = true
					}
				}
			}
		}
		if !added {
			break
		}
	}
}

// ikHasTargetInIndexes はIKターゲット/リンクが指定集合にあるか判定する。
func ikHasTargetInIndexes(ik *model.Ik, indexes map[int]struct{}) bool {
	if ik == nil {
		return false
	}
	if ik.BoneIndex >= 0 {
		if _, ok := indexes[ik.BoneIndex]; ok {
			return true
		}
	}
	for _, link := range ik.Links {
		if link.BoneIndex < 0 {
			continue
		}
		if _, ok := indexes[link.BoneIndex]; ok {
			return true
		}
	}
	return false
}

// addBoneAndParents はボーンと親階層を追加する。
func addBoneAndParents(modelData *model.PmxModel, bone *model.Bone, indexes map[int]struct{}) {
	if modelData == nil || bone == nil {
		return
	}
	current := bone
	for current != nil {
		indexes[current.Index()] = struct{}{}
		if current.ParentIndex < 0 {
			break
		}
		parent, err := modelData.Bones.Get(current.ParentIndex)
		if err != nil || parent == nil {
			break
		}
		current = parent
	}
}

// addIkDependencies はIK関連のボーンを追加する。
func addIkDependencies(modelData *model.PmxModel, bone *model.Bone, indexes map[int]struct{}) {
	if modelData == nil || bone == nil || bone.Ik == nil {
		return
	}
	if bone.Ik.BoneIndex >= 0 {
		if target, err := modelData.Bones.Get(bone.Ik.BoneIndex); err == nil && target != nil {
			addBoneAndParents(modelData, target, indexes)
		}
	}
	for _, link := range bone.Ik.Links {
		if link.BoneIndex < 0 {
			continue
		}
		linkBone, err := modelData.Bones.Get(link.BoneIndex)
		if err != nil || linkBone == nil {
			continue
		}
		addBoneAndParents(modelData, linkBone, indexes)
	}
}

// sortedBoneIndexes はBoneDeltasに含まれるボーンindexを返す。
func sortedBoneIndexes(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas) []int {
	if modelData == nil || boneDeltas == nil {
		return nil
	}
	indexes := make([]int, 0)
	boneDeltas.ForEach(func(index int, delta *delta.BoneDelta) bool {
		if delta == nil {
			return true
		}
		indexes = append(indexes, index)
		return true
	})
	sortBoneIndexes(modelData, indexes)
	return indexes
}

// sortBoneIndexes はlayer->index順に整列する。
func sortBoneIndexes(modelData *model.PmxModel, indexes []int) {
	if modelData == nil {
		return
	}
	sort.Slice(indexes, func(i, j int) bool {
		boneI, _ := modelData.Bones.Get(indexes[i])
		boneJ, _ := modelData.Bones.Get(indexes[j])
		layerI := 0
		layerJ := 0
		if boneI != nil {
			layerI = boneI.Layer
		}
		if boneJ != nil {
			layerJ = boneJ.Layer
		}
		if layerI == layerJ {
			return indexes[i] < indexes[j]
		}
		return layerI < layerJ
	})
}

// getBoneFrame は対象ボーンのフレームを返す。
func getBoneFrame(motionData *motion.VmdMotion, name string, frame motion.Frame) *motion.BoneFrame {
	if motionData == nil || motionData.BoneFrames == nil {
		return nil
	}
	if !motionData.BoneFrames.Has(name) {
		return nil
	}
	return motionData.BoneFrames.Get(name).Get(frame)
}

// updateBoneDelta はユニット行列を更新する。
func updateBoneDelta(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, d *delta.BoneDelta) {
	if modelData == nil || boneDeltas == nil || d == nil || d.Bone == nil {
		return
	}
	unit := d.ResetUnitMatrix()
	if unit == nil {
		return
	}
	localMat := calculateTotalLocalMat(boneDeltas, d.Bone.Index())
	if localMat != nil {
		unit.MulToPtr(localMat, unit)
	}
	var scaleMat mmath.Mat4
	if calculateTotalScaleMat(boneDeltas, d.Bone.Index(), &scaleMat) {
		unit.MulToPtr(&scaleMat, unit)
	}
	var pos mmath.Vec3
	if calculateTotalPosition(boneDeltas, d.Bone.Index(), &pos) {
		unit.MulTranslateTo(&pos, unit)
	}
	var rotMat mmath.Mat4
	if calculateTotalRotationMat(boneDeltas, d.Bone.Index(), &rotMat) {
		unit.MulToPtr(&rotMat, unit)
	}
	var revert mmath.Mat4
	boneRevertOffsetMatTo(modelData, d.Bone, &revert)
	revert.MulToPtr(unit, unit)
	d.GlobalMatrix = nil
	d.LocalMatrix = nil
	d.GlobalPosition = nil
	boneDeltas.Update(d)
}

// applyGlobalMatrix は親行列を用いてグローバル行列を更新する。
func applyGlobalMatrix(boneDeltas *delta.BoneDeltas, d *delta.BoneDelta) {
	if boneDeltas == nil || d == nil || d.Bone == nil {
		return
	}
	if d.UnitMatrix == nil {
		d.ResetUnitMatrix()
	}
	parent := boneDeltas.Get(d.Bone.ParentIndex)
	var global *mmath.Mat4
	unitIsIdent := d.UnitMatrix != nil && *d.UnitMatrix == mmath.IDENT_MAT4
	switch {
	case parent != nil && parent.GlobalIkOffMatrix != nil && boneIsIk(parent.Bone):
		global = d.GlobalMatrixPtr()
		if unitIsIdent {
			copyMat4(global, parent.GlobalIkOffMatrix)
		} else {
			parent.GlobalIkOffMatrix.MulToPtr(d.UnitMatrix, global)
		}
	case parent != nil && parent.GlobalMatrix != nil:
		global = d.GlobalMatrixPtr()
		if unitIsIdent {
			copyMat4(global, parent.GlobalMatrix)
		} else {
			parent.GlobalMatrix.MulToPtr(d.UnitMatrix, global)
		}
	default:
		d.GlobalMatrix = d.UnitMatrix
		global = d.UnitMatrix
	}
	if global == nil {
		return
	}
	local := d.LocalMatrixPtr()
	if local == nil {
		return
	}
	neg := d.Bone.Position.Negated()
	global.MulTranslateTo(&neg, local)
	d.SetGlobalPosition(global.Translation())
}

// applyGlobalMatrixNoLocal はグローバル行列のみ更新する。
func applyGlobalMatrixNoLocal(boneDeltas *delta.BoneDeltas, d *delta.BoneDelta) {
	if boneDeltas == nil || d == nil || d.Bone == nil {
		return
	}
	if d.UnitMatrix == nil {
		d.ResetUnitMatrix()
	}
	parent := boneDeltas.Get(d.Bone.ParentIndex)
	var global *mmath.Mat4
	unitIsIdent := d.UnitMatrix != nil && *d.UnitMatrix == mmath.IDENT_MAT4
	switch {
	case parent != nil && parent.GlobalIkOffMatrix != nil && boneIsIk(parent.Bone):
		global = d.GlobalMatrixPtr()
		if unitIsIdent {
			copyMat4(global, parent.GlobalIkOffMatrix)
		} else {
			parent.GlobalIkOffMatrix.MulToPtr(d.UnitMatrix, global)
		}
	case parent != nil && parent.GlobalMatrix != nil:
		global = d.GlobalMatrixPtr()
		if unitIsIdent {
			copyMat4(global, parent.GlobalMatrix)
		} else {
			parent.GlobalMatrix.MulToPtr(d.UnitMatrix, global)
		}
	default:
		d.GlobalMatrix = d.UnitMatrix
		global = d.UnitMatrix
	}
	if global == nil {
		return
	}
	d.LocalMatrix = nil
	d.GlobalPosition = nil
}

// copyMat4 はMat4の値をdstへコピーする。
func copyMat4(dst, src *mmath.Mat4) {
	if dst == nil || src == nil {
		return
	}
	dst[0] = src[0]
	dst[1] = src[1]
	dst[2] = src[2]
	dst[3] = src[3]
	dst[4] = src[4]
	dst[5] = src[5]
	dst[6] = src[6]
	dst[7] = src[7]
	dst[8] = src[8]
	dst[9] = src[9]
	dst[10] = src[10]
	dst[11] = src[11]
	dst[12] = src[12]
	dst[13] = src[13]
	dst[14] = src[14]
	dst[15] = src[15]
}

// calculateTotalRotationMat は総回転行列をoutへ書き込み、反映有無を返す。
func calculateTotalRotationMat(boneDeltas *delta.BoneDeltas, boneIndex int, out *mmath.Mat4) bool {
	if out == nil {
		return false
	}
	rot := accumulateTotalRotation(boneDeltas, boneIndex, 0, 1.0)
	hasRot := rot != nil && !rot.IsIdent()
	if rot != nil {
		rot.ToMat4To(out)
	} else {
		*out = mmath.NewMat4()
	}
	return applyCancelableRotation(boneDeltas, boneIndex, out, hasRot)
}

// calculateTotalPosition は総移動ベクトルをoutへ書き込み、反映有無を返す。
func calculateTotalPosition(boneDeltas *delta.BoneDeltas, boneIndex int, out *mmath.Vec3) bool {
	if out == nil {
		return false
	}
	pos := accumulateTotalPosition(boneDeltas, boneIndex, 0)
	if pos != nil {
		*out = *pos
	} else {
		*out = mmath.NewVec3()
	}
	return applyCancelablePosition(boneDeltas, boneIndex, out)
}

// calculateTotalScaleMat は総スケール行列をoutへ書き込み、反映有無を返す。
func calculateTotalScaleMat(boneDeltas *delta.BoneDeltas, boneIndex int, out *mmath.Mat4) bool {
	if out == nil {
		return false
	}
	scale := accumulateTotalScale(boneDeltas, boneIndex, 0)
	hasScale := scale != nil && !scale.IsOne()
	if scale != nil {
		scale.ToScaleMat4To(out)
	} else {
		*out = mmath.NewMat4()
	}
	return applyCancelableScale(boneDeltas, boneIndex, out, hasScale)
}

// calculateTotalLocalMat は総ローカル行列を返す。
func calculateTotalLocalMat(boneDeltas *delta.BoneDeltas, boneIndex int) *mmath.Mat4 {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return nil
	}
	return bd.TotalLocalMat()
}

// accumulateTotalRotation は回転付与を再帰合成する。
func accumulateTotalRotation(
	boneDeltas *delta.BoneDeltas,
	boneIndex int,
	recursion int,
	factor float64,
) *mmath.Quaternion {
	if recursion > maxEffectorRecursion {
		return nil
	}
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return nil
	}
	rot := bd.TotalRotation()
	if boneIsEffectorRotation(bd.Bone) {
		effectorRot := accumulateTotalRotation(boneDeltas, bd.Bone.EffectIndex, recursion+1, bd.Bone.EffectFactor)
		if effectorRot != nil {
			if rot == nil {
				rot = effectorRot
			} else {
				tmp := rot.Muled(*effectorRot)
				rot = &tmp
			}
		}
	}
	if rot == nil {
		return nil
	}
	out := rot.MuledScalar(factor)
	return &out
}

// accumulateTotalPosition は移動付与を再帰合成する。
func accumulateTotalPosition(boneDeltas *delta.BoneDeltas, boneIndex int, recursion int) *mmath.Vec3 {
	if recursion > maxEffectorRecursion {
		return nil
	}
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return nil
	}
	pos := bd.TotalPosition()
	if boneIsEffectorTranslation(bd.Bone) {
		effectorPos := accumulateTotalPosition(boneDeltas, bd.Bone.EffectIndex, recursion+1)
		if effectorPos != nil {
			factorPos := effectorPos.MuledScalar(bd.Bone.EffectFactor)
			if pos == nil {
				out := factorPos
				return &out
			}
			out := pos.Added(factorPos)
			return &out
		}
	}
	return pos
}

// accumulateTotalScale はスケールを合成する。
func accumulateTotalScale(boneDeltas *delta.BoneDeltas, boneIndex int, recursion int) *mmath.Vec3 {
	if recursion > maxEffectorRecursion {
		return nil
	}
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return nil
	}
	return bd.TotalScale()
}

// applyCancelableRotation はキャンセル回転を適用し、反映有無を返す。
func applyCancelableRotation(boneDeltas *delta.BoneDeltas, boneIndex int, rotMat *mmath.Mat4, hasRot bool) bool {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil || rotMat == nil {
		return false
	}
	parentMat := getParentCancelableRotationMat(boneDeltas, bd.Bone.ParentIndex)
	hasSelfCancel := (bd.FrameCancelableRotation != nil && !bd.FrameCancelableRotation.IsIdent()) ||
		(bd.FrameMorphCancelableRotation != nil && !bd.FrameMorphCancelableRotation.IsIdent())
	if !hasSelfCancel {
		if parentMat == nil {
			return hasRot
		}
		var inv mmath.Mat4
		parentMat.InvertedTo(&inv)
		rotMat.MulToPtr(&inv, rotMat)
		return true
	}
	if bd.FrameCancelableRotation != nil && !bd.FrameCancelableRotation.IsIdent() {
		var cancelMat mmath.Mat4
		bd.FrameCancelableRotation.ToMat4To(&cancelMat)
		rotMat.MulToPtr(&cancelMat, rotMat)
	}
	if bd.FrameMorphCancelableRotation != nil && !bd.FrameMorphCancelableRotation.IsIdent() {
		var morphMat mmath.Mat4
		bd.FrameMorphCancelableRotation.ToMat4To(&morphMat)
		rotMat.MulToPtr(&morphMat, rotMat)
	}
	if parentMat == nil {
		return true
	}
	var inv mmath.Mat4
	parentMat.InvertedTo(&inv)
	rotMat.MulToPtr(&inv, rotMat)
	return true
}

// applyCancelablePosition はキャンセル移動を適用し、反映有無を返す。
func applyCancelablePosition(boneDeltas *delta.BoneDeltas, boneIndex int, pos *mmath.Vec3) bool {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil || pos == nil {
		return false
	}
	var parentPos mmath.Vec3
	hasParent := getParentCancelablePosition(boneDeltas, bd.Bone.ParentIndex, &parentPos)
	hasSelfCancel := (bd.FrameCancelablePosition != nil && !bd.FrameCancelablePosition.IsZero()) ||
		(bd.FrameMorphCancelablePosition != nil && !bd.FrameMorphCancelablePosition.IsZero())

	if hasSelfCancel {
		if bd.FrameCancelablePosition != nil && !bd.FrameCancelablePosition.IsZero() {
			*pos = pos.Added(*bd.FrameCancelablePosition)
		}
		if bd.FrameMorphCancelablePosition != nil && !bd.FrameMorphCancelablePosition.IsZero() {
			*pos = pos.Added(*bd.FrameMorphCancelablePosition)
		}
	}
	if hasParent {
		*pos = pos.Subed(parentPos)
	}
	return !pos.IsZero()
}

// applyCancelableScale はキャンセルスケールを適用し、反映有無を返す。
func applyCancelableScale(boneDeltas *delta.BoneDeltas, boneIndex int, scaleMat *mmath.Mat4, hasScale bool) bool {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil || scaleMat == nil {
		return false
	}
	parentMat := getParentCancelableScaleMat(boneDeltas, bd.Bone.ParentIndex)
	hasSelfCancel := (bd.FrameCancelableScale != nil && !bd.FrameCancelableScale.IsZero()) ||
		(bd.FrameMorphCancelableScale != nil && !bd.FrameMorphCancelableScale.IsZero())
	if !hasSelfCancel {
		if parentMat == nil {
			return hasScale
		}
		var inv mmath.Mat4
		parentMat.InvertedTo(&inv)
		scaleMat.MulToPtr(&inv, scaleMat)
		return true
	}
	if bd.FrameCancelableScale != nil && !bd.FrameCancelableScale.IsZero() {
		var cancelMat mmath.Mat4
		bd.FrameCancelableScale.ToScaleMat4To(&cancelMat)
		scaleMat.MulToPtr(&cancelMat, scaleMat)
	}
	if bd.FrameMorphCancelableScale != nil && !bd.FrameMorphCancelableScale.IsZero() {
		var morphMat mmath.Mat4
		bd.FrameMorphCancelableScale.ToScaleMat4To(&morphMat)
		scaleMat.MulToPtr(&morphMat, scaleMat)
	}
	if parentMat == nil {
		return true
	}
	var inv mmath.Mat4
	parentMat.InvertedTo(&inv)
	scaleMat.MulToPtr(&inv, scaleMat)
	return true
}

// getParentCancelableRotationMat は親キャンセル回転行列を返す。
func getParentCancelableRotationMat(boneDeltas *delta.BoneDeltas, parentIndex int) *mmath.Mat4 {
	if boneDeltas == nil || !boneDeltas.Contains(parentIndex) {
		return nil
	}
	parent := boneDeltas.Get(parentIndex)
	if parent == nil {
		return nil
	}
	var mat *mmath.Mat4
	if parent.FrameCancelableRotation != nil && !parent.FrameCancelableRotation.IsIdent() {
		var tmp mmath.Mat4
		parent.FrameCancelableRotation.ToMat4To(&tmp)
		mat = &tmp
	}
	if parent.FrameMorphCancelableRotation != nil && !parent.FrameMorphCancelableRotation.IsIdent() {
		var tmp mmath.Mat4
		parent.FrameMorphCancelableRotation.ToMat4To(&tmp)
		if mat == nil {
			mat = &tmp
		} else {
			mat.MulToPtr(&tmp, mat)
		}
	}
	return mat
}

// getParentCancelablePosition は親キャンセル移動ベクトルを返す。
func getParentCancelablePosition(boneDeltas *delta.BoneDeltas, parentIndex int, out *mmath.Vec3) bool {
	if boneDeltas == nil || out == nil || !boneDeltas.Contains(parentIndex) {
		return false
	}
	parent := boneDeltas.Get(parentIndex)
	if parent == nil {
		return false
	}
	hasPos := false
	var pos mmath.Vec3
	if parent.FrameCancelablePosition != nil && !parent.FrameCancelablePosition.IsZero() {
		pos = *parent.FrameCancelablePosition
		hasPos = true
	}
	if parent.FrameMorphCancelablePosition != nil && !parent.FrameMorphCancelablePosition.IsZero() {
		if !hasPos {
			pos = *parent.FrameMorphCancelablePosition
			hasPos = true
		} else {
			pos = pos.Added(*parent.FrameMorphCancelablePosition)
		}
	}
	if !hasPos {
		return false
	}
	*out = pos
	return true
}

// getParentCancelableScaleMat は親キャンセルスケール行列を返す。
func getParentCancelableScaleMat(boneDeltas *delta.BoneDeltas, parentIndex int) *mmath.Mat4 {
	if boneDeltas == nil || !boneDeltas.Contains(parentIndex) {
		return nil
	}
	parent := boneDeltas.Get(parentIndex)
	if parent == nil {
		return nil
	}
	var mat *mmath.Mat4
	if parent.FrameCancelableScale != nil && !parent.FrameCancelableScale.IsZero() {
		var tmp mmath.Mat4
		parent.FrameCancelableScale.ToScaleMat4To(&tmp)
		mat = &tmp
	}
	if parent.FrameMorphCancelableScale != nil && !parent.FrameMorphCancelableScale.IsZero() {
		var tmp mmath.Mat4
		parent.FrameMorphCancelableScale.ToScaleMat4To(&tmp)
		if mat == nil {
			mat = &tmp
		} else {
			mat.MulToPtr(&tmp, mat)
		}
	}
	return mat
}

// boneRevertOffsetMatTo は逆オフセット行列をoutへ書き込む。
func boneRevertOffsetMatTo(modelData *model.PmxModel, bone *model.Bone, out *mmath.Mat4) {
	if out == nil {
		return
	}
	if bone == nil {
		*out = mmath.NewMat4()
		return
	}
	parentPos := mmath.NewVec3()
	if modelData != nil && bone.ParentIndex >= 0 {
		parent, err := modelData.Bones.Get(bone.ParentIndex)
		if err == nil && parent != nil {
			parentPos = parent.Position
		}
	}
	relative := bone.Position.Subed(parentPos)
	relative.ToMat4To(out)
}

// boneIsAfterPhysics は物理後変形ボーンか判定する。
func boneIsAfterPhysics(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&model.BONE_FLAG_IS_AFTER_PHYSICS_DEFORM != 0
}

// boneIsIk はIKボーンか判定する。
func boneIsIk(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&model.BONE_FLAG_IS_IK != 0 && bone.Ik != nil
}

// boneHasEffector は付与ボーンか判定する。
func boneHasEffector(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return boneIsEffectorRotation(bone) || boneIsEffectorTranslation(bone)
}

// boneIsEffectorRotation は回転付与か判定する。
func boneIsEffectorRotation(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_ROTATION != 0 && bone.EffectIndex >= 0
}

// boneIsEffectorTranslation は移動付与か判定する。
func boneIsEffectorTranslation(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&model.BONE_FLAG_IS_EXTERNAL_TRANSLATION != 0 && bone.EffectIndex >= 0
}

// boneHasFixedAxis は固定軸を持つか判定する。
func boneHasFixedAxis(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&model.BONE_FLAG_HAS_FIXED_AXIS != 0 && !bone.FixedAxis.IsZero()
}

// boneLocalAxes はローカル軸ベクトルを返す。
func boneLocalAxes(modelData *model.PmxModel, bone *model.Bone) (mmath.Vec3, mmath.Vec3, mmath.Vec3) {
	if bone == nil {
		return mmath.UNIT_X_VEC3, mmath.UNIT_Y_VEC3, mmath.UNIT_Z_NEG_VEC3
	}
	xAxis := boneChildDirection(modelData, bone)
	if boneHasFixedAxis(bone) {
		xAxis = bone.FixedAxis
	}
	if xAxis.IsZero() {
		xAxis = mmath.UNIT_X_VEC3
	}
	xAxis = xAxis.Normalized()
	yAxis := xAxis.Cross(mmath.UNIT_Z_NEG_VEC3)
	if yAxis.Length() == 0 {
		yAxis = mmath.UNIT_Y_VEC3
	}
	zAxis := xAxis.Cross(yAxis)
	if zAxis.Length() == 0 {
		zAxis = mmath.UNIT_Z_NEG_VEC3
	}
	return xAxis, yAxis, zAxis
}

// boneChildDirection は子方向の軸ベクトルを返す。
func boneChildDirection(modelData *model.PmxModel, bone *model.Bone) mmath.Vec3 {
	if bone == nil {
		return mmath.UNIT_X_VEC3
	}
	if modelData != nil && bone.TailIndex >= 0 {
		if tail, err := modelData.Bones.Get(bone.TailIndex); err == nil && tail != nil {
			dir := tail.Position.Subed(bone.Position)
			if !dir.IsZero() {
				return dir.Normalized()
			}
		}
	}
	if !bone.TailPosition.IsZero() {
		return bone.TailPosition.Normalized()
	}
	return mmath.UNIT_X_VEC3
}

// applyIkDeltas はIKボーンの差分を適用する。
func applyIkDeltas(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	boneDeltas *delta.BoneDeltas,
	frame motion.Frame,
	deformBoneIndexes []int,
	removeTwist bool,
) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	effectorChildren := buildEffectorChildrenSlice(modelData)
	children := buildChildRelations(modelData)
	scratch := newIkScratch(modelData, deformBoneIndexes)
	ikFrame := getIkFrame(motionData, frame)
	for _, boneIndex := range deformBoneIndexes {
		bone, err := modelData.Bones.Get(boneIndex)
		if err != nil || bone == nil || !boneIsIk(bone) {
			continue
		}
		if ikFrame != nil && !ikFrame.IsEnable(bone.Name()) {
			continue
		}
		ikIndexes := collectIkChainIndexes(modelData, bone)
		for _, idx := range ikIndexes {
			d := boneDeltas.Get(idx)
			if d == nil {
				continue
			}
			off := d.FilledGlobalMatrix()
			d.SetGlobalIkOffMatrix(off)
			boneDeltas.Update(d)
		}
		applyIkForBone(modelData, motionData, boneDeltas, bone, frame, deformBoneIndexes, effectorChildren, children, scratch, removeTwist)
	}
}

// getIkFrame はIKフレームを返す。
func getIkFrame(motionData *motion.VmdMotion, frame motion.Frame) *motion.IkFrame {
	if motionData == nil || motionData.IkFrames == nil {
		return nil
	}
	return motionData.IkFrames.Get(frame)
}

// ikScratch はIK計算の作業用バッファを保持する。
type ikScratch struct {
	unitUpdatedFlags   []bool
	unitUpdatedList    []int
	unitQueue          []int
	globalUpdatedFlags []bool
	globalUpdatedList  []int
	globalQueue        []int
	globalRecalc       []int
	orderByRank        []int
}

// newIkScratch はIK計算の作業バッファを初期化する。
func newIkScratch(modelData *model.PmxModel, deformBoneIndexes []int) *ikScratch {
	size := 0
	if modelData != nil && modelData.Bones != nil {
		size = modelData.Bones.Len()
	}
	orderByRank := make([]int, size)
	for i := range orderByRank {
		orderByRank[i] = i
	}
	sortBoneIndexes(modelData, orderByRank)
	return &ikScratch{
		unitUpdatedFlags:   make([]bool, size),
		unitUpdatedList:    make([]int, 0, size),
		unitQueue:          make([]int, 0, size),
		globalUpdatedFlags: make([]bool, size),
		globalUpdatedList:  make([]int, 0, size),
		globalQueue:        make([]int, 0, size),
		globalRecalc:       make([]int, 0, size),
		orderByRank:        orderByRank,
	}
}

// applyIkForBone はIKボーンの回転を更新する。
func applyIkForBone(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	boneDeltas *delta.BoneDeltas,
	ikBone *model.Bone,
	frame motion.Frame,
	deformBoneIndexes []int,
	effectorChildren [][]int,
	children [][]int,
	scratch *ikScratch,
	removeTwist bool,
) {
	if modelData == nil || boneDeltas == nil || ikBone == nil || ikBone.Ik == nil {
		return
	}
	if len(ikBone.Ik.Links) == 0 {
		return
	}
	ikTargetIndex := ikBone.Ik.BoneIndex
	if ikTargetIndex < 0 {
		return
	}
	loopCount := max(ikBone.Ik.LoopCount, 1)
	targetBeforeIk := isTargetBeforeIk(deformBoneIndexes, ikTargetIndex, ikBone.Index())
	if targetBeforeIk {
		loopCount++
	}

	isSingleIk := len(ikBone.Ik.Links) == 1
	ikDelta := boneDeltas.Get(ikBone.Index())
	if ikDelta == nil {
		ikDelta = delta.NewBoneDelta(ikBone, frame)
	}
	ikPos := ikDelta.FilledGlobalPosition()
	ikOnPos := ikPos
	useToeIk := false
	if targetBeforeIk && len(ikBone.Ik.Links) == 1 && isToeIkBone(ikBone) && motionData != nil {
		if targetBone, err := modelData.Bones.Get(ikTargetIndex); err == nil && targetBone != nil {
			ikOffDeltas, ikIndexes := ComputeBoneDeltas(modelData, motionData, frame, []string{targetBone.Name()}, false, false, false)
			applyBoneMatricesWithIndexes(modelData, ikOffDeltas, ikIndexes)
			if ikOffDeltas != nil {
				if targetDelta := ikOffDeltas.Get(ikTargetIndex); targetDelta != nil {
					ikPos = targetDelta.FilledGlobalPosition()
					useToeIk = true
				}
			}
		}
	}

	if scratch == nil {
		scratch = newIkScratch(modelData, deformBoneIndexes)
	}
	unitUpdatedFlags := scratch.unitUpdatedFlags
	unitUpdatedList := scratch.unitUpdatedList[:0]
	unitQueue := scratch.unitQueue[:0]
	globalUpdatedFlags := scratch.globalUpdatedFlags
	globalUpdatedList := scratch.globalUpdatedList[:0]
	globalQueue := scratch.globalQueue[:0]
	globalRecalc := scratch.globalRecalc[:0]
	orderByRank := scratch.orderByRank
	for loop := 0; loop < loopCount; loop++ {
		for linkIndex, link := range ikBone.Ik.Links {
			linkBone, err := modelData.Bones.Get(link.BoneIndex)
			if err != nil || linkBone == nil {
				continue
			}
			linkDelta := boneDeltas.Get(linkBone.Index())
			if linkDelta == nil {
				linkDelta = delta.NewBoneDelta(linkBone, frame)
			}
			linkQuat := linkDelta.FilledTotalRotation()

			if useToeIk && loop == 1 && linkIndex == 0 {
				ikPos = ikOnPos
			}

			ikTargetDelta := boneDeltas.Get(ikTargetIndex)
			if ikTargetDelta == nil {
				continue
			}
			ikTargetPos := ikTargetDelta.FilledGlobalPosition()
			linkGlobal := linkDelta.FilledGlobalMatrix()
			linkInv := linkGlobal.Inverted()
			ikTargetLocalPos := linkInv.MulVec3(ikTargetPos).Normalized()
			ikLocalPos := linkInv.MulVec3(ikPos).Normalized()
			if ikTargetLocalPos.Length() == 0 || ikLocalPos.Length() == 0 {
				continue
			}

			unitRad := ikBone.Ik.UnitRotation.X * float64(linkIndex+1)
			linkAngle := mmath.VectorToRadian(ikTargetLocalPos, ikLocalPos)
			if linkAngle > unitRad {
				linkAngle = unitRad
			}
			limitedAxis := ikTargetLocalPos.Cross(ikLocalPos).Normalized()
			if (!isSingleIk || linkAngle > mmath.Gimbal1Rad) && (link.AngleLimit || link.LocalAngleLimit) {
				limitedAxis = getLinkAxis(link, ikTargetLocalPos, ikLocalPos)
			}

			resultQuat := SolveIkStep(IkSolveStepInput{
				LinkRotation:    linkQuat,
				LimitedAxis:     limitedAxis,
				LinkAngle:       linkAngle,
				MinAngleLimit:   link.MinAngleLimit,
				MaxAngleLimit:   link.MaxAngleLimit,
				LocalMinLimit:   link.LocalMinAngleLimit,
				LocalMaxLimit:   link.LocalMaxAngleLimit,
				AngleLimit:      link.AngleLimit,
				LocalAngleLimit: link.LocalAngleLimit,
				Loop:            loop,
				LoopCount:       loopCount,
				RemoveTwist:     removeTwist,
				FixedAxis:       fixedAxisOrZero(linkBone),
				ChildAxis:       boneChildDirection(modelData, linkBone),
				LocalAxes:       localAxes(modelData, linkBone),
			})

			if linkDelta.FrameMorphRotation != nil && !linkDelta.FrameMorphRotation.IsIdent() {
				resultQuat = resultQuat.Muled(linkDelta.FrameMorphRotation.Inverted())
			}
			rot := resultQuat
			linkDelta.FrameRotation = &rot
			linkDelta.InvalidateTotals()
			updateBoneDelta(modelData, boneDeltas, linkDelta)
			collectEffectorRelatedFlags(effectorChildren, linkBone.Index(), unitUpdatedFlags, &unitUpdatedList, &unitQueue)
			for _, idx := range unitUpdatedList {
				if idx == linkBone.Index() {
					continue
				}
				d := boneDeltas.Get(idx)
				if d == nil || d.Bone == nil {
					continue
				}
				updateBoneDelta(modelData, boneDeltas, d)
			}
			for _, idx := range unitUpdatedList {
				collectDescendantFlags(children, idx, globalUpdatedFlags, &globalUpdatedList, &globalQueue)
			}
			globalRecalc = buildRecalcIndexesByOrderFlags(orderByRank, globalUpdatedFlags, globalRecalc)
			applyGlobalMatricesWithIndexes(modelData, boneDeltas, globalRecalc)
			resetFlagsByList(unitUpdatedFlags, unitUpdatedList)
			unitUpdatedList = unitUpdatedList[:0]
			resetFlagsByList(globalUpdatedFlags, globalUpdatedList)
			globalUpdatedList = globalUpdatedList[:0]
		}

		threshold := ikTargetDistance(boneDeltas, ikBone.Index(), ikTargetIndex)
		if threshold <= 1e-5 {
			break
		}
	}
	scratch.unitUpdatedList = unitUpdatedList
	scratch.unitQueue = unitQueue
	scratch.globalUpdatedList = globalUpdatedList
	scratch.globalQueue = globalQueue
	scratch.globalRecalc = globalRecalc
}

// ikTargetDistance はIKターゲット距離を返す。
func ikTargetDistance(boneDeltas *delta.BoneDeltas, ikBoneIndex, targetIndex int) float64 {
	if boneDeltas == nil {
		return math.MaxFloat64
	}
	ik := boneDeltas.Get(ikBoneIndex)
	target := boneDeltas.Get(targetIndex)
	if ik == nil || target == nil {
		return math.MaxFloat64
	}
	return ik.FilledGlobalPosition().Distance(target.FilledGlobalPosition())
}

// snapshotLinkRotations はリンク回転を退避する。
func snapshotLinkRotations(
	boneDeltas *delta.BoneDeltas,
	links []model.IkLink,
) map[int]mmath.Quaternion {
	rotations := map[int]mmath.Quaternion{}
	for _, link := range links {
		if link.BoneIndex < 0 {
			continue
		}
		linkDelta := boneDeltas.Get(link.BoneIndex)
		if linkDelta == nil {
			continue
		}
		rotations[link.BoneIndex] = linkDelta.FilledFrameRotation()
	}
	return rotations
}

// isTargetBeforeIk はターゲットが先に変形されるか判定する。
func isTargetBeforeIk(indexes []int, targetIndex, ikIndex int) bool {
	if len(indexes) == 0 {
		return false
	}
	targetPos := slices.Index(indexes, targetIndex)
	ikPos := slices.Index(indexes, ikIndex)
	if targetPos < 0 || ikPos < 0 {
		return false
	}
	return targetPos < ikPos
}

// isToeIkBone はつま先IKボーンか判定する。
func isToeIkBone(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	name := bone.Name()
	return strings.Contains(name, "つま先ＩＫ") || strings.Contains(name, "つま先IK")
}

// max は最大値を返す。
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min は最小値を返す。
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
