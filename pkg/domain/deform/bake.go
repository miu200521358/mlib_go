package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// BakeBoneFrameByGlobalMatrix はグローバル行列を現モデルの親子構造に再分解してVMDフレームへ反映する。
func BakeBoneFrameByGlobalMatrix(
	modelData *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	motionData *motion.VmdMotion,
	boneName string,
	frame motion.Frame,
	globalMatrix mmath.Mat4,
) *delta.BoneDelta {
	if modelData == nil || boneDeltas == nil || motionData == nil {
		return nil
	}
	bone, err := modelData.Bones.GetByName(boneName)
	if err != nil || bone == nil {
		return nil
	}
	parent := boneDeltas.Get(bone.ParentIndex)
	bakedDelta := delta.NewBoneDeltaByGlobalMatrix(bone, frame, globalMatrix, parent)
	if bakedDelta == nil {
		return nil
	}
	boneDeltas.Update(bakedDelta)
	bf := motionData.BoneFrames.Get(boneName).Get(frame)
	pos := bakedDelta.FilledFramePosition()
	rot := bakedDelta.FilledFrameRotation()
	bf.Position = &pos
	bf.Rotation = &rot
	return bakedDelta
}
