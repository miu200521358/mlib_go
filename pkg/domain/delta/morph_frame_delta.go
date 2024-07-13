package delta

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

type MorphFrameDelta struct {
	FramePosition *mmath.MVec3       // キーフレ位置の変動量
	FrameRotation *mmath.MQuaternion // キーフレ回転の変動量
	FrameScale    *mmath.MVec3       // キーフレスケールの変動量
}

func (md *MorphFrameDelta) GetFramePosition() *mmath.MVec3 {
	if md.FramePosition == nil {
		md.FramePosition = mmath.NewMVec3()
	}
	return md.FramePosition
}

func (md *MorphFrameDelta) GetFrameRotation() *mmath.MQuaternion {
	if md.FrameRotation == nil {
		md.FrameRotation = mmath.NewMQuaternion()
	}
	return md.FrameRotation
}

func (md *MorphFrameDelta) GetFrameScale() *mmath.MVec3 {
	if md.FrameScale == nil {
		md.FrameScale = mmath.NewMVec3()
	}
	return md.FrameScale
}

func NewMorphFrameDelta() *MorphFrameDelta {
	return &MorphFrameDelta{}
}

func (md *MorphFrameDelta) Copy() *MorphFrameDelta {
	return &MorphFrameDelta{
		FramePosition: md.GetFramePosition().Copy(),
		FrameRotation: md.GetFrameRotation().Copy(),
		FrameScale:    md.GetFrameScale().Copy(),
	}
}
