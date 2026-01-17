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
	deformBoneIndexes := collectBoneIndexes(modelData, boneNames, includeIk, afterPhysics)
	morphDeltas := ComputeMorphDeltas(modelData, motionData, frame, nil)

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
				pos := *bf.Position
				d.FramePosition = &pos
			}
			if bf.Rotation != nil {
				rot := *bf.Rotation
				d.FrameRotation = &rot
			}
			if bf.CancelablePosition != nil {
				pos := *bf.CancelablePosition
				d.FrameCancelablePosition = &pos
			}
			if bf.CancelableRotation != nil {
				rot := *bf.CancelableRotation
				d.FrameCancelableRotation = &rot
			}
			if bf.Scale != nil {
				scale := *bf.Scale
				d.FrameScale = &scale
			}
			if bf.CancelableScale != nil {
				scale := *bf.CancelableScale
				d.FrameCancelableScale = &scale
			}
		}

		if morphDeltas != nil && morphDeltas.Bones() != nil {
			morphDelta := morphDeltas.Bones().Get(boneIndex)
			if morphDelta != nil {
				if morphDelta.FramePosition != nil {
					pos := *morphDelta.FramePosition
					d.FrameMorphPosition = &pos
				}
				if morphDelta.FrameRotation != nil {
					rot := *morphDelta.FrameRotation
					d.FrameMorphRotation = &rot
				}
				if morphDelta.FrameCancelablePosition != nil {
					pos := *morphDelta.FrameCancelablePosition
					d.FrameMorphCancelablePosition = &pos
				}
				if morphDelta.FrameCancelableRotation != nil {
					rot := *morphDelta.FrameCancelableRotation
					d.FrameMorphCancelableRotation = &rot
				}
				if morphDelta.FrameScale != nil {
					scale := *morphDelta.FrameScale
					d.FrameMorphScale = &scale
				}
				if morphDelta.FrameCancelableScale != nil {
					scale := *morphDelta.FrameCancelableScale
					d.FrameMorphCancelableScale = &scale
				}
				if morphDelta.FrameLocalMat != nil {
					mat := *morphDelta.FrameLocalMat
					d.FrameLocalMorphMat = &mat
				}
			}
		}

		boneDeltas.Update(d)
	}

	if includeIk {
		ApplyBoneMatrices(modelData, boneDeltas)
		applyIkDeltas(modelData, motionData, boneDeltas, frame, deformBoneIndexes, removeTwist)
		ApplyBoneMatrices(modelData, boneDeltas)
	}

	return boneDeltas, deformBoneIndexes
}

// ApplyBoneMatrices はボーン行列を合成して差分へ反映する。
func ApplyBoneMatrices(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas) {
	if modelData == nil || boneDeltas == nil {
		return
	}
	indexes := sortedBoneIndexes(modelData, boneDeltas)
	for _, boneIndex := range indexes {
		d := boneDeltas.Get(boneIndex)
		if d == nil || d.Bone == nil {
			continue
		}
		updateBoneDelta(modelData, boneDeltas, d)
		applyGlobalMatrix(boneDeltas, d)
		boneDeltas.Update(d)
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
	if !slices.Contains(motionData.BoneFrames.Names(), name) {
		return nil
	}
	return motionData.BoneFrames.Get(name).Get(frame)
}

// updateBoneDelta はユニット行列を更新する。
func updateBoneDelta(modelData *model.PmxModel, boneDeltas *delta.BoneDeltas, d *delta.BoneDelta) {
	if modelData == nil || boneDeltas == nil || d == nil || d.Bone == nil {
		return
	}
	unit := mmath.NewMat4()
	localMat := calculateTotalLocalMat(boneDeltas, d.Bone.Index())
	if !localMat.IsIdent() {
		unit = unit.Muled(localMat)
	}
	scaleMat := calculateTotalScaleMat(boneDeltas, d.Bone.Index())
	if !scaleMat.IsIdent() {
		unit = unit.Muled(scaleMat)
	}
	posMat := calculateTotalPositionMat(boneDeltas, d.Bone.Index())
	if !posMat.IsIdent() {
		unit = unit.Muled(posMat)
	}
	rotMat := calculateTotalRotationMat(boneDeltas, d.Bone.Index())
	if !rotMat.IsIdent() {
		unit = unit.Muled(rotMat)
	}
	revert := boneRevertOffsetMat(modelData, d.Bone)
	unit = revert.Muled(unit)
	d.UnitMatrix = &unit
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
		unit := mmath.NewMat4()
		d.UnitMatrix = &unit
	}
	parent := boneDeltas.Get(d.Bone.ParentIndex)
	var global mmath.Mat4
	switch {
	case parent != nil && parent.GlobalIkOffMatrix != nil && boneIsIk(parent.Bone):
		global = parent.GlobalIkOffMatrix.Muled(*d.UnitMatrix)
	case parent != nil && parent.GlobalMatrix != nil:
		global = parent.GlobalMatrix.Muled(*d.UnitMatrix)
	default:
		global = *d.UnitMatrix
	}
	d.GlobalMatrix = &global
	local := global.Muled(boneOffsetMat(d.Bone))
	d.LocalMatrix = &local
	pos := global.Translation()
	d.GlobalPosition = &pos
}

// calculateTotalRotationMat は総回転行列を返す。
func calculateTotalRotationMat(boneDeltas *delta.BoneDeltas, boneIndex int) mmath.Mat4 {
	rot := accumulateTotalRotation(boneDeltas, boneIndex, 0, 1.0)
	rotMat := mmath.NewMat4()
	if rot != nil {
		rotMat = rot.ToMat4()
	}
	return applyCancelableRotation(boneDeltas, boneIndex, rotMat)
}

// calculateTotalPositionMat は総移動行列を返す。
func calculateTotalPositionMat(boneDeltas *delta.BoneDeltas, boneIndex int) mmath.Mat4 {
	pos := accumulateTotalPosition(boneDeltas, boneIndex, 0)
	posMat := mmath.NewMat4()
	if pos != nil {
		posMat = pos.ToMat4()
	}
	return applyCancelablePosition(boneDeltas, boneIndex, posMat)
}

// calculateTotalScaleMat は総スケール行列を返す。
func calculateTotalScaleMat(boneDeltas *delta.BoneDeltas, boneIndex int) mmath.Mat4 {
	scale := accumulateTotalScale(boneDeltas, boneIndex, 0)
	scaleMat := mmath.NewMat4()
	if scale != nil {
		scaleMat = scale.ToScaleMat4()
	}
	return applyCancelableScale(boneDeltas, boneIndex, scaleMat)
}

// calculateTotalLocalMat は総ローカル行列を返す。
func calculateTotalLocalMat(boneDeltas *delta.BoneDeltas, boneIndex int) mmath.Mat4 {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return mmath.NewMat4()
	}
	return bd.FilledTotalLocalMat()
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
				out := factorPos.ToMat4().Translation()
				return &out
			}
			out := pos.ToMat4().Muled(factorPos.ToMat4()).Translation()
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

// applyCancelableRotation はキャンセル回転を適用する。
func applyCancelableRotation(boneDeltas *delta.BoneDeltas, boneIndex int, rotMat mmath.Mat4) mmath.Mat4 {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return rotMat
	}
	parentMat := getParentCancelableRotationMat(boneDeltas, bd.Bone.ParentIndex)
	hasSelfCancel := (bd.FrameCancelableRotation != nil && !bd.FrameCancelableRotation.IsIdent()) ||
		(bd.FrameMorphCancelableRotation != nil && !bd.FrameMorphCancelableRotation.IsIdent())
	if !hasSelfCancel {
		if parentMat == nil {
			return rotMat
		}
		return rotMat.Muled(parentMat.Inverted())
	}
	if bd.FrameCancelableRotation != nil && !bd.FrameCancelableRotation.IsIdent() {
		rotMat = rotMat.Muled(bd.FrameCancelableRotation.ToMat4())
	}
	if bd.FrameMorphCancelableRotation != nil && !bd.FrameMorphCancelableRotation.IsIdent() {
		rotMat = rotMat.Muled(bd.FrameMorphCancelableRotation.ToMat4())
	}
	if parentMat == nil {
		return rotMat
	}
	return rotMat.Muled(parentMat.Inverted())
}

// applyCancelablePosition はキャンセル移動を適用する。
func applyCancelablePosition(boneDeltas *delta.BoneDeltas, boneIndex int, posMat mmath.Mat4) mmath.Mat4 {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return posMat
	}
	parentMat := getParentCancelablePositionMat(boneDeltas, bd.Bone.ParentIndex)
	hasSelfCancel := (bd.FrameCancelablePosition != nil && !bd.FrameCancelablePosition.IsZero()) ||
		(bd.FrameMorphCancelablePosition != nil && !bd.FrameMorphCancelablePosition.IsZero())
	if !hasSelfCancel {
		if parentMat == nil {
			return posMat
		}
		return posMat.Muled(parentMat.Inverted())
	}
	if bd.FrameCancelablePosition != nil && !bd.FrameCancelablePosition.IsZero() {
		posMat = posMat.Muled(bd.FrameCancelablePosition.ToMat4())
	}
	if bd.FrameMorphCancelablePosition != nil && !bd.FrameMorphCancelablePosition.IsZero() {
		posMat = posMat.Muled(bd.FrameMorphCancelablePosition.ToMat4())
	}
	if parentMat == nil {
		return posMat
	}
	return posMat.Muled(parentMat.Inverted())
}

// applyCancelableScale はキャンセルスケールを適用する。
func applyCancelableScale(boneDeltas *delta.BoneDeltas, boneIndex int, scaleMat mmath.Mat4) mmath.Mat4 {
	bd := boneDeltas.Get(boneIndex)
	if bd == nil {
		return scaleMat
	}
	parentMat := getParentCancelableScaleMat(boneDeltas, bd.Bone.ParentIndex)
	hasSelfCancel := (bd.FrameCancelableScale != nil && !bd.FrameCancelableScale.IsZero()) ||
		(bd.FrameMorphCancelableScale != nil && !bd.FrameMorphCancelableScale.IsZero())
	if !hasSelfCancel {
		if parentMat == nil {
			return scaleMat
		}
		return scaleMat.Muled(parentMat.Inverted())
	}
	if bd.FrameCancelableScale != nil && !bd.FrameCancelableScale.IsZero() {
		scaleMat = scaleMat.Muled(bd.FrameCancelableScale.ToScaleMat4())
	}
	if bd.FrameMorphCancelableScale != nil && !bd.FrameMorphCancelableScale.IsZero() {
		scaleMat = scaleMat.Muled(bd.FrameMorphCancelableScale.ToScaleMat4())
	}
	if parentMat == nil {
		return scaleMat
	}
	return scaleMat.Muled(parentMat.Inverted())
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
		tmp := parent.FrameCancelableRotation.ToMat4()
		mat = &tmp
	}
	if parent.FrameMorphCancelableRotation != nil && !parent.FrameMorphCancelableRotation.IsIdent() {
		tmp := parent.FrameMorphCancelableRotation.ToMat4()
		if mat == nil {
			mat = &tmp
		} else {
			m := mat.Muled(tmp)
			mat = &m
		}
	}
	return mat
}

// getParentCancelablePositionMat は親キャンセル移動行列を返す。
func getParentCancelablePositionMat(boneDeltas *delta.BoneDeltas, parentIndex int) *mmath.Mat4 {
	if boneDeltas == nil || !boneDeltas.Contains(parentIndex) {
		return nil
	}
	parent := boneDeltas.Get(parentIndex)
	if parent == nil {
		return nil
	}
	var mat *mmath.Mat4
	if parent.FrameCancelablePosition != nil && !parent.FrameCancelablePosition.IsZero() {
		tmp := parent.FrameCancelablePosition.ToMat4()
		mat = &tmp
	}
	if parent.FrameMorphCancelablePosition != nil && !parent.FrameMorphCancelablePosition.IsZero() {
		tmp := parent.FrameMorphCancelablePosition.ToMat4()
		if mat == nil {
			mat = &tmp
		} else {
			m := mat.Muled(tmp)
			mat = &m
		}
	}
	return mat
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
		tmp := parent.FrameCancelableScale.ToScaleMat4()
		mat = &tmp
	}
	if parent.FrameMorphCancelableScale != nil && !parent.FrameMorphCancelableScale.IsZero() {
		tmp := parent.FrameMorphCancelableScale.ToScaleMat4()
		if mat == nil {
			mat = &tmp
		} else {
			m := mat.Muled(tmp)
			mat = &m
		}
	}
	return mat
}

// boneOffsetMat はオフセット行列を返す。
func boneOffsetMat(bone *model.Bone) mmath.Mat4 {
	if bone == nil {
		return mmath.NewMat4()
	}
	return bone.Position.Negated().ToMat4()
}

// boneRevertOffsetMat は逆オフセット行列を返す。
func boneRevertOffsetMat(modelData *model.PmxModel, bone *model.Bone) mmath.Mat4 {
	if bone == nil {
		return mmath.NewMat4()
	}
	parentPos := mmath.NewVec3()
	if modelData != nil && bone.ParentIndex >= 0 {
		parent, err := modelData.Bones.Get(bone.ParentIndex)
		if err == nil && parent != nil {
			parentPos = parent.Position
		}
	}
	relative := bone.Position.Subed(parentPos)
	return relative.ToMat4()
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
			d.GlobalIkOffMatrix = &off
			boneDeltas.Update(d)
		}
		applyIkForBone(modelData, motionData, boneDeltas, bone, frame, deformBoneIndexes, removeTwist)
	}
}

// getIkFrame はIKフレームを返す。
func getIkFrame(motionData *motion.VmdMotion, frame motion.Frame) *motion.IkFrame {
	if motionData == nil || motionData.IkFrames == nil {
		return nil
	}
	return motionData.IkFrames.Get(frame)
}

// applyIkForBone はIKボーンの回転を更新する。
func applyIkForBone(
	modelData *model.PmxModel,
	motionData *motion.VmdMotion,
	boneDeltas *delta.BoneDeltas,
	ikBone *model.Bone,
	frame motion.Frame,
	deformBoneIndexes []int,
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
			ikOffDeltas, _ := ComputeBoneDeltas(modelData, motionData, frame, []string{targetBone.Name()}, false, false, false)
			ApplyBoneMatrices(modelData, ikOffDeltas)
			if ikOffDeltas != nil {
				if targetDelta := ikOffDeltas.Get(ikTargetIndex); targetDelta != nil {
					ikPos = targetDelta.FilledGlobalPosition()
					useToeIk = true
				}
			}
		}
	}

	bestThreshold := math.MaxFloat64
	bestRotations := map[int]mmath.Quaternion{}

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
			updateBoneDelta(modelData, boneDeltas, linkDelta)
			ApplyBoneMatrices(modelData, boneDeltas)
		}

		threshold := ikTargetDistance(boneDeltas, ikBone.Index(), ikTargetIndex)
		if threshold < bestThreshold {
			bestThreshold = threshold
			bestRotations = snapshotLinkRotations(boneDeltas, ikBone.Ik.Links)
		}
		if threshold <= 1e-5 {
			break
		}
	}

	if len(bestRotations) > 0 {
		for linkIndex, rot := range bestRotations {
			linkDelta := boneDeltas.Get(linkIndex)
			if linkDelta == nil {
				continue
			}
			copyRot := rot
			linkDelta.FrameRotation = &copyRot
			updateBoneDelta(modelData, boneDeltas, linkDelta)
		}
		ApplyBoneMatrices(modelData, boneDeltas)
	}
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
