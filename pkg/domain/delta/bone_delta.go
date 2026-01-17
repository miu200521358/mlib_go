// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	sharedtime "github.com/miu200521358/mlib_go/pkg/shared/contracts/time"
)

// BoneDelta は1ボーンの差分と派生行列を保持する。
type BoneDelta struct {
	Bone              *model.Bone
	Frame             sharedtime.Frame
	globalIkOffMatrix mmath.Mat4
	globalMatrix      mmath.Mat4
	localMatrix       mmath.Mat4
	unitMatrix        mmath.Mat4
	globalPosition    mmath.Vec3
	GlobalIkOffMatrix *mmath.Mat4
	GlobalMatrix      *mmath.Mat4
	LocalMatrix       *mmath.Mat4
	UnitMatrix        *mmath.Mat4
	GlobalPosition    *mmath.Vec3

	FramePosition                *mmath.Vec3
	FrameMorphPosition           *mmath.Vec3
	FrameCancelablePosition      *mmath.Vec3
	FrameMorphCancelablePosition *mmath.Vec3

	FrameRotation                *mmath.Quaternion
	FrameMorphRotation           *mmath.Quaternion
	FrameCancelableRotation      *mmath.Quaternion
	FrameMorphCancelableRotation *mmath.Quaternion

	FrameScale                *mmath.Vec3
	FrameMorphScale           *mmath.Vec3
	FrameCancelableScale      *mmath.Vec3
	FrameMorphCancelableScale *mmath.Vec3

	FrameLocalMat      *mmath.Mat4
	FrameLocalMorphMat *mmath.Mat4
}

// NewBoneDelta はBoneDeltaを生成する。
func NewBoneDelta(bone *model.Bone, frame sharedtime.Frame) *BoneDelta {
	return &BoneDelta{Bone: bone, Frame: frame}
}

// SetGlobalIkOffMatrix はIKオフ用の行列を設定する。
func (d *BoneDelta) SetGlobalIkOffMatrix(mat mmath.Mat4) {
	if d == nil {
		return
	}
	d.globalIkOffMatrix = mat
	d.GlobalIkOffMatrix = &d.globalIkOffMatrix
}

// SetGlobalMatrix はグローバル行列を設定する。
func (d *BoneDelta) SetGlobalMatrix(mat mmath.Mat4) {
	if d == nil {
		return
	}
	d.globalMatrix = mat
	d.GlobalMatrix = &d.globalMatrix
}

// SetLocalMatrix はローカル行列を設定する。
func (d *BoneDelta) SetLocalMatrix(mat mmath.Mat4) {
	if d == nil {
		return
	}
	d.localMatrix = mat
	d.LocalMatrix = &d.localMatrix
}

// SetUnitMatrix はユニット行列を設定する。
func (d *BoneDelta) SetUnitMatrix(mat mmath.Mat4) {
	if d == nil {
		return
	}
	d.unitMatrix = mat
	d.UnitMatrix = &d.unitMatrix
}

// SetGlobalPosition はグローバル位置を設定する。
func (d *BoneDelta) SetGlobalPosition(pos mmath.Vec3) {
	if d == nil {
		return
	}
	d.globalPosition = pos
	d.GlobalPosition = &d.globalPosition
}

// NewBoneDeltaByGlobalMatrix はグローバル行列をもとにボーン差分を生成する。
func NewBoneDeltaByGlobalMatrix(
	bone *model.Bone,
	frame sharedtime.Frame,
	globalMatrix mmath.Mat4,
	parent *BoneDelta,
) *BoneDelta {
	if bone == nil {
		return nil
	}
	parentGlobal := mmath.NewMat4()
	parentPos := mmath.NewVec3()
	if parent != nil {
		parentGlobal = parent.FilledGlobalMatrix()
		if parent.Bone != nil {
			parentPos = parent.Bone.Position
		}
	}
	localMat := globalMatrix.Muled(bone.Position.Negated().ToMat4())
	unitMat := parentGlobal.Inverted().Muled(globalMatrix)
	parentRelative := bone.Position.Subed(parentPos)
	framePos := unitMat.Translation().Subed(parentRelative)
	frameRot := unitMat.Quaternion()
	global := globalMatrix
	local := localMat
	unit := unitMat
	d := &BoneDelta{
		Bone:          bone,
		Frame:         frame,
		FramePosition: &framePos,
		FrameRotation: &frameRot,
	}
	d.globalMatrix = global
	d.GlobalMatrix = &d.globalMatrix
	d.localMatrix = local
	d.LocalMatrix = &d.localMatrix
	d.unitMatrix = unit
	d.UnitMatrix = &d.unitMatrix
	return d
}

// FilledGlobalMatrix はグローバル行列を返す。
func (d *BoneDelta) FilledGlobalMatrix() mmath.Mat4 {
	if d == nil || d.GlobalMatrix == nil {
		return mmath.NewMat4()
	}
	return *d.GlobalMatrix
}

// FilledLocalMatrix はローカル行列を返す。
func (d *BoneDelta) FilledLocalMatrix() mmath.Mat4 {
	if d == nil {
		return mmath.NewMat4()
	}
	if d.LocalMatrix == nil {
		global := d.FilledGlobalMatrix()
		offset := mmath.NewMat4()
		if d.Bone != nil {
			offset = d.Bone.Position.Negated().ToMat4()
		}
		d.localMatrix = global.Muled(offset)
		d.LocalMatrix = &d.localMatrix
	}
	return *d.LocalMatrix
}

// FilledUnitMatrix はユニット行列を返す。
func (d *BoneDelta) FilledUnitMatrix() mmath.Mat4 {
	if d == nil || d.UnitMatrix == nil {
		return mmath.NewMat4()
	}
	return *d.UnitMatrix
}

// FilledGlobalPosition はグローバル位置を返す。
func (d *BoneDelta) FilledGlobalPosition() mmath.Vec3 {
	if d == nil {
		return mmath.NewVec3()
	}
	if d.GlobalPosition == nil {
		d.globalPosition = d.FilledGlobalMatrix().Translation()
		d.GlobalPosition = &d.globalPosition
	}
	return *d.GlobalPosition
}

// FilledFrameRotation はフレーム回転を返す。
func (d *BoneDelta) FilledFrameRotation() mmath.Quaternion {
	if d == nil || d.FrameRotation == nil {
		return mmath.NewQuaternion()
	}
	return *d.FrameRotation
}

// FilledFrameMorphRotation はフレームモーフ回転を返す。
func (d *BoneDelta) FilledFrameMorphRotation() mmath.Quaternion {
	if d == nil || d.FrameMorphRotation == nil {
		return mmath.NewQuaternion()
	}
	return *d.FrameMorphRotation
}

// FilledFrameCancelableRotation はフレームキャンセル回転を返す。
func (d *BoneDelta) FilledFrameCancelableRotation() mmath.Quaternion {
	if d == nil || d.FrameCancelableRotation == nil {
		return mmath.NewQuaternion()
	}
	return *d.FrameCancelableRotation
}

// FilledFrameMorphCancelableRotation はフレームモーフキャンセル回転を返す。
func (d *BoneDelta) FilledFrameMorphCancelableRotation() mmath.Quaternion {
	if d == nil || d.FrameMorphCancelableRotation == nil {
		return mmath.NewQuaternion()
	}
	return *d.FrameMorphCancelableRotation
}

// TotalRotation はモーフ込みの回転を返す。
func (d *BoneDelta) TotalRotation() *mmath.Quaternion {
	if d == nil {
		return nil
	}
	var rot mmath.Quaternion
	hasRot := false
	if d.FrameRotation != nil {
		rot = *d.FrameRotation
		hasRot = true
	}
	if d.FrameMorphRotation != nil && !d.FrameMorphRotation.IsIdent() {
		if !hasRot {
			rot = *d.FrameMorphRotation
			hasRot = true
		} else {
			rot = rot.Muled(*d.FrameMorphRotation)
		}
	}
	if !hasRot {
		return nil
	}
	if boneHasFixedAxis(d.Bone) {
		rot = rot.ToFixedAxisRotation(d.Bone.FixedAxis.Normalized())
	}
	return &rot
}

// FilledTotalRotation はモーフ込みの回転を返す。
func (d *BoneDelta) FilledTotalRotation() mmath.Quaternion {
	rot := d.TotalRotation()
	if rot == nil {
		return mmath.NewQuaternion()
	}
	return *rot
}

// FilledFramePosition はフレーム位置を返す。
func (d *BoneDelta) FilledFramePosition() mmath.Vec3 {
	if d == nil || d.FramePosition == nil {
		return mmath.NewVec3()
	}
	return *d.FramePosition
}

// FilledFrameMorphPosition はフレームモーフ位置を返す。
func (d *BoneDelta) FilledFrameMorphPosition() mmath.Vec3 {
	if d == nil || d.FrameMorphPosition == nil {
		return mmath.NewVec3()
	}
	return *d.FrameMorphPosition
}

// FilledFrameCancelablePosition はフレームキャンセル位置を返す。
func (d *BoneDelta) FilledFrameCancelablePosition() mmath.Vec3 {
	if d == nil || d.FrameCancelablePosition == nil {
		return mmath.NewVec3()
	}
	return *d.FrameCancelablePosition
}

// FilledFrameMorphCancelablePosition はフレームモーフキャンセル位置を返す。
func (d *BoneDelta) FilledFrameMorphCancelablePosition() mmath.Vec3 {
	if d == nil || d.FrameMorphCancelablePosition == nil {
		return mmath.NewVec3()
	}
	return *d.FrameMorphCancelablePosition
}

// TotalPosition はモーフ込みの位置を返す。
func (d *BoneDelta) TotalPosition() *mmath.Vec3 {
	if d == nil {
		return nil
	}
	var pos mmath.Vec3
	hasPos := false
	if d.FramePosition != nil {
		pos = *d.FramePosition
		hasPos = true
	}
	if d.FrameMorphPosition != nil && !d.FrameMorphPosition.IsZero() {
		if !hasPos {
			pos = *d.FrameMorphPosition
			hasPos = true
		} else {
			pos = pos.Added(*d.FrameMorphPosition)
		}
	}
	if !hasPos {
		return nil
	}
	return &pos
}

// FilledTotalPosition はモーフ込みの位置を返す。
func (d *BoneDelta) FilledTotalPosition() mmath.Vec3 {
	pos := d.TotalPosition()
	if pos == nil {
		return mmath.NewVec3()
	}
	return *pos
}

// FilledFrameScale はフレームスケールを返す。
func (d *BoneDelta) FilledFrameScale() mmath.Vec3 {
	if d == nil || d.FrameScale == nil {
		return mmath.ONE_VEC3
	}
	return *d.FrameScale
}

// FilledFrameMorphScale はフレームモーフスケールを返す。
func (d *BoneDelta) FilledFrameMorphScale() mmath.Vec3 {
	if d == nil || d.FrameMorphScale == nil {
		return mmath.ONE_VEC3
	}
	return *d.FrameMorphScale
}

// FilledFrameCancelableScale はフレームキャンセルスケールを返す。
func (d *BoneDelta) FilledFrameCancelableScale() mmath.Vec3 {
	if d == nil || d.FrameCancelableScale == nil {
		return mmath.ONE_VEC3
	}
	return *d.FrameCancelableScale
}

// FilledFrameMorphCancelableScale はフレームモーフキャンセルスケールを返す。
func (d *BoneDelta) FilledFrameMorphCancelableScale() mmath.Vec3 {
	if d == nil || d.FrameMorphCancelableScale == nil {
		return mmath.ONE_VEC3
	}
	return *d.FrameMorphCancelableScale
}

// TotalScale はモーフ込みのスケールを返す。
func (d *BoneDelta) TotalScale() *mmath.Vec3 {
	if d == nil {
		return nil
	}
	var scale mmath.Vec3
	hasScale := false
	if d.FrameScale != nil {
		scale = *d.FrameScale
		hasScale = true
	}
	if d.FrameMorphScale != nil && !d.FrameMorphScale.IsOne() {
		if !hasScale {
			scale = *d.FrameMorphScale
			hasScale = true
		} else {
			scale = scale.Muled(*d.FrameMorphScale)
		}
	}
	if !hasScale {
		return nil
	}
	return &scale
}

// FilledTotalScale はモーフ込みのスケールを返す。
func (d *BoneDelta) FilledTotalScale() mmath.Vec3 {
	scale := d.TotalScale()
	if scale == nil {
		return mmath.ONE_VEC3
	}
	return *scale
}

// FilledFrameLocalMat はフレームローカル行列を返す。
func (d *BoneDelta) FilledFrameLocalMat() mmath.Mat4 {
	if d == nil || d.FrameLocalMat == nil {
		return mmath.NewMat4()
	}
	return *d.FrameLocalMat
}

// FilledFrameLocalMorphMat はフレームモーフローカル行列を返す。
func (d *BoneDelta) FilledFrameLocalMorphMat() mmath.Mat4 {
	if d == nil || d.FrameLocalMorphMat == nil {
		return mmath.NewMat4()
	}
	return *d.FrameLocalMorphMat
}

// TotalLocalMat はモーフ込みのローカル行列を返す。
func (d *BoneDelta) TotalLocalMat() *mmath.Mat4 {
	if d == nil {
		return nil
	}
	var mat mmath.Mat4
	hasMat := false
	if d.FrameLocalMat != nil && !d.FrameLocalMat.IsIdent() {
		mat = *d.FrameLocalMat
		hasMat = true
	}
	if d.FrameLocalMorphMat != nil && !d.FrameLocalMorphMat.IsIdent() {
		if !hasMat {
			mat = *d.FrameLocalMorphMat
			hasMat = true
		} else {
			mat = mat.Muled(*d.FrameLocalMorphMat)
		}
	}
	if !hasMat {
		return nil
	}
	return &mat
}

// FilledTotalLocalMat はモーフ込みのローカル行列を返す。
func (d *BoneDelta) FilledTotalLocalMat() mmath.Mat4 {
	mat := d.TotalLocalMat()
	if mat == nil {
		return mmath.NewMat4()
	}
	return *mat
}

// boneHasFixedAxis は固定軸を持つか判定する。
func boneHasFixedAxis(bone *model.Bone) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&model.BONE_FLAG_HAS_FIXED_AXIS != 0 && !bone.FixedAxis.IsZero()
}
