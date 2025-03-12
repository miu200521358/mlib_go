package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// BoneDelta は1つのボーンにおける変形（ポジション・回転・スケールなど）の差分を表す
type BoneDelta struct {
	Bone              *pmx.Bone
	Frame             float32
	GlobalIkOffMatrix *mmath.MMat4
	GlobalMatrix      *mmath.MMat4
	LocalMatrix       *mmath.MMat4
	GlobalPosition    *mmath.MVec3
	UnitMatrix        *mmath.MMat4

	FramePosition                *mmath.MVec3
	FrameMorphPosition           *mmath.MVec3
	FrameCancelablePosition      *mmath.MVec3
	FrameMorphCancelablePosition *mmath.MVec3

	FrameRotation                *mmath.MQuaternion
	FrameMorphRotation           *mmath.MQuaternion
	FrameCancelableRotation      *mmath.MQuaternion
	FrameMorphCancelableRotation *mmath.MQuaternion

	FrameScale                *mmath.MVec3
	FrameMorphScale           *mmath.MVec3
	FrameCancelableScale      *mmath.MVec3
	FrameMorphCancelableScale *mmath.MVec3

	FrameLocalMat      *mmath.MMat4
	FrameLocalMorphMat *mmath.MMat4
}

// NewBoneDelta は新規の BoneDelta を生成するコンストラクタ
func NewBoneDelta(bone *pmx.Bone, frame float32) *BoneDelta {
	return &BoneDelta{
		Bone:  bone,
		Frame: frame,
	}
}

// NewBoneDeltaByGlobalMatrix はグローバル行列をもとにボーンデルタを生成する
func NewBoneDeltaByGlobalMatrix(bone *pmx.Bone, frame float32, globalMatrix *mmath.MMat4, parent *BoneDelta) *BoneDelta {
	var parentGlobal *mmath.MMat4
	if parent != nil {
		parentGlobal = parent.FilledGlobalMatrix()
	} else {
		parentGlobal = mmath.NewMMat4()
	}
	localMat := globalMatrix.Muled(bone.OffsetMatrix)
	unitMat := parentGlobal.Inverted().Muled(globalMatrix)

	return &BoneDelta{
		Bone:         bone,
		Frame:        frame,
		GlobalMatrix: globalMatrix,
		LocalMatrix:  localMat,
		UnitMatrix:   unitMat,
		// 位置と回転を抜き出して保持
		// FramePosition: unitMat.Translation(),	// ちゃむちかの左袖白が落ちる
		FrameRotation: unitMat.Quaternion(),
	}
}

// 各種 Filled◯◯ メソッドは、 nil の場合に初期化して返す “getter + lazy init” 的な役割です。
func (bd *BoneDelta) FilledGlobalMatrix() *mmath.MMat4 {
	if bd.GlobalMatrix == nil {
		bd.GlobalMatrix = mmath.NewMMat4()
	}
	return bd.GlobalMatrix
}

func (bd *BoneDelta) FilledLocalMatrix() *mmath.MMat4 {
	if bd.LocalMatrix == nil {
		bd.LocalMatrix = bd.FilledGlobalMatrix().Muled(bd.Bone.OffsetMatrix)
	}
	return bd.LocalMatrix
}

func (bd *BoneDelta) FilledGlobalPosition() *mmath.MVec3 {
	if bd.GlobalPosition == nil {
		bd.GlobalPosition = bd.FilledGlobalMatrix().Translation()
	}
	return bd.GlobalPosition
}

// 回転は同じ処理をまとめ、骨特有の制御（FixedAxisなど）もここに含める
func (bd *BoneDelta) FilledFrameRotation() *mmath.MQuaternion {
	if bd.FrameRotation == nil {
		bd.FrameRotation = mmath.NewMQuaternion()
	}
	return bd.FrameRotation
}
func (bd *BoneDelta) FilledFrameMorphRotation() *mmath.MQuaternion {
	if bd.FrameMorphRotation == nil {
		bd.FrameMorphRotation = mmath.NewMQuaternion()
	}
	return bd.FrameMorphRotation
}
func (bd *BoneDelta) FilledFrameCancelableRotation() *mmath.MQuaternion {
	if bd.FrameCancelableRotation == nil {
		bd.FrameCancelableRotation = mmath.NewMQuaternion()
	}
	return bd.FrameCancelableRotation
}
func (bd *BoneDelta) FilledFrameMorphCancelableRotation() *mmath.MQuaternion {
	if bd.FrameMorphCancelableRotation == nil {
		bd.FrameMorphCancelableRotation = mmath.NewMQuaternion()
	}
	return bd.FrameMorphCancelableRotation
}

// トータルの回転(モーフ含む)を返す(なかった場合、Identを返す)
func (bd *BoneDelta) FilledTotalRotation() *mmath.MQuaternion {
	rot := bd.TotalRotation()
	if rot == nil {
		rot = mmath.MQuaternionIdent
	}
	return rot
}

// トータルの回転(モーフ含む)を返す(なかった場合、nilを返す)
func (bd *BoneDelta) TotalRotation() *mmath.MQuaternion {
	rot := bd.FrameRotation
	morphRot := bd.FrameMorphRotation
	if morphRot != nil && !morphRot.IsIdent() {
		if rot == nil {
			rot = morphRot.Copy()
		} else {
			rot = rot.Muled(morphRot)
		}
	}

	// ボーンが固定軸を持っている場合の最終調整
	if rot != nil && bd.Bone.HasFixedAxis() {
		rot = rot.ToFixedAxisRotation(bd.Bone.NormalizedFixedAxis)
	}

	return rot
}

func (bd *BoneDelta) FilledFramePosition() *mmath.MVec3 {
	if bd.FramePosition == nil {
		bd.FramePosition = mmath.NewMVec3()
	}
	return bd.FramePosition
}
func (bd *BoneDelta) FilledFrameMorphPosition() *mmath.MVec3 {
	if bd.FrameMorphPosition == nil {
		bd.FrameMorphPosition = mmath.NewMVec3()
	}
	return bd.FrameMorphPosition
}
func (bd *BoneDelta) FilledFrameCancelablePosition() *mmath.MVec3 {
	if bd.FrameCancelablePosition == nil {
		bd.FrameCancelablePosition = mmath.NewMVec3()
	}
	return bd.FrameCancelablePosition
}
func (bd *BoneDelta) FilledFrameMorphCancelablePosition() *mmath.MVec3 {
	if bd.FrameMorphCancelablePosition == nil {
		bd.FrameMorphCancelablePosition = mmath.NewMVec3()
	}
	return bd.FrameMorphCancelablePosition
}

// トータルの位置(モーフ含む)を返す(なかった場合、Zeroを返す)
func (bd *BoneDelta) FilledTotalPosition() *mmath.MVec3 {
	pos := bd.TotalPosition()
	if pos == nil {
		return mmath.MVec3Zero
	}
	return pos
}

// トータルの位置(モーフ含む)を返す(なかった場合、nilを返す)
func (bd *BoneDelta) TotalPosition() *mmath.MVec3 {
	pos := bd.FramePosition
	morphPos := bd.FrameMorphPosition
	if morphPos != nil && !morphPos.IsZero() {
		if pos == nil {
			pos = morphPos
		} else {
			pos = pos.Added(morphPos)
		}
	}
	return pos
}

// スケール系
func (bd *BoneDelta) FilledFrameScale() *mmath.MVec3 {
	if bd.FrameScale == nil {
		bd.FrameScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return bd.FrameScale
}
func (bd *BoneDelta) FilledFrameMorphScale() *mmath.MVec3 {
	if bd.FrameMorphScale == nil {
		bd.FrameMorphScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return bd.FrameMorphScale
}
func (bd *BoneDelta) FilledFrameCancelableScale() *mmath.MVec3 {
	if bd.FrameCancelableScale == nil {
		bd.FrameCancelableScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return bd.FrameCancelableScale
}
func (bd *BoneDelta) FilledFrameMorphCancelableScale() *mmath.MVec3 {
	if bd.FrameMorphCancelableScale == nil {
		bd.FrameMorphCancelableScale = &mmath.MVec3{X: 1, Y: 1, Z: 1}
	}
	return bd.FrameMorphCancelableScale
}

// トータルのスケール(モーフ含む)を返す(なかった場合、Oneを返す)
func (bd *BoneDelta) FilledTotalScale() *mmath.MVec3 {
	scale := bd.TotalScale()
	if scale == nil {
		return mmath.MVec3One
	}
	return scale
}

// トータルのスケール(モーフ含む)を返す(なかった場合、nilを返す)
func (bd *BoneDelta) TotalScale() *mmath.MVec3 {
	scale := bd.FrameScale
	morphScale := bd.FrameMorphScale
	if morphScale != nil && !morphScale.IsOne() {
		if scale == nil {
			scale = morphScale
		} else {
			scale = scale.Muled(morphScale)
		}
	}
	return scale
}

// ローカル行列
func (bd *BoneDelta) FilledFrameLocalMat() *mmath.MMat4 {
	if bd.FrameLocalMat == nil {
		bd.FrameLocalMat = mmath.NewMMat4()
	}
	return bd.FrameLocalMat
}
func (bd *BoneDelta) FilledFrameLocalMorphMat() *mmath.MMat4 {
	if bd.FrameLocalMorphMat == nil {
		bd.FrameLocalMorphMat = mmath.NewMMat4()
	}
	return bd.FrameLocalMorphMat
}

// トータルのローカル行列(モーフ含む)(なかった場合、Identを返す)
func (bd *BoneDelta) FilledTotalLocalMat() *mmath.MMat4 {
	mat := bd.TotalLocalMat()
	if mat == nil {
		return mmath.MMat4Ident
	}
	return mat
}

// トータルのローカル行列(モーフ含む)(なかった場合、nilを返す)
func (bd *BoneDelta) TotalLocalMat() *mmath.MMat4 {
	mat := bd.FrameLocalMat
	morphMat := bd.FrameLocalMorphMat
	if morphMat != nil && !morphMat.IsIdent() {
		if mat == nil {
			mat = morphMat
		} else {
			mat = mat.Muled(morphMat)
		}
	}
	return mat
}
