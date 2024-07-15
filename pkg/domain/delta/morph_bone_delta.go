package delta

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

type MorphBoneDelta struct {
	MorphPosition *mmath.MVec3       // キーフレ位置の変動量
	MorphRotation *mmath.MQuaternion // キーフレ回転の変動量
	MorphScale    *mmath.MVec3       // キーフレスケールの変動量
}

func NewMorphBoneDelta() *MorphBoneDelta {
	return &MorphBoneDelta{}
}

func (md *MorphBoneDelta) FilledMorphPosition() *mmath.MVec3 {
	if md.MorphPosition == nil {
		md.MorphPosition = mmath.NewMVec3()
	}
	return md.MorphPosition
}

func (md *MorphBoneDelta) FilledMorphRotation() *mmath.MQuaternion {
	if md.MorphRotation == nil {
		md.MorphRotation = mmath.NewMQuaternion()
	}
	return md.MorphRotation
}

func (md *MorphBoneDelta) FilledMorphScale() *mmath.MVec3 {
	if md.MorphScale == nil {
		md.MorphScale = mmath.NewMVec3()
	}
	return md.MorphScale
}

func (md *MorphBoneDelta) Copy() *MorphBoneDelta {
	return &MorphBoneDelta{
		MorphPosition: md.FilledMorphPosition().Copy(),
		MorphRotation: md.FilledMorphRotation().Copy(),
		MorphScale:    md.FilledMorphScale().Copy(),
	}
}
