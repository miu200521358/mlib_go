package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type MorphFrame struct {
	*BaseFrame         // キーフレ
	Ratio      float64 // モーフの割合
}

func NewMorphFrame(index int) *MorphFrame {
	return &MorphFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Ratio:     0.0,
	}
}

func NullMorphFrame() *MorphFrame {
	return nil
}

func (mf *MorphFrame) Add(v *MorphFrame) {
	mf.Ratio += v.Ratio
}

func (mf *MorphFrame) Added(v *MorphFrame) *MorphFrame {
	copied := mf.Copy().(*MorphFrame)
	copied.Ratio += v.Ratio
	return copied
}

func (mf *MorphFrame) Copy() IBaseFrame {
	copied := NewMorphFrame(mf.GetIndex())
	copied.Ratio = mf.Ratio
	return copied
}

func (nextMf *MorphFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevMf := prevFrame.(*MorphFrame)

	prevIndex := prevMf.GetIndex()
	nextIndex := nextMf.GetIndex()

	mf := NewMorphFrame(index)

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*ry

	return mf
}

func (mf *MorphFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}

type MorphFrameDelta struct {
	framePosition *mmath.MVec3       // キーフレ位置の変動量
	frameRotation *mmath.MQuaternion // キーフレ回転の変動量
	frameScale    *mmath.MVec3       // キーフレスケールの変動量
}

func (md *MorphFrameDelta) FramePosition() *mmath.MVec3 {
	if md.framePosition == nil {
		md.framePosition = mmath.NewMVec3()
	}
	return md.framePosition
}

func (md *MorphFrameDelta) FrameRotation() *mmath.MQuaternion {
	if md.frameRotation == nil {
		md.frameRotation = mmath.NewMQuaternion()
	}
	return md.frameRotation
}

func (md *MorphFrameDelta) FrameScale() *mmath.MVec3 {
	if md.frameScale == nil {
		md.frameScale = mmath.NewMVec3()
	}
	return md.frameScale
}

func NewMorphFrameDelta() *MorphFrameDelta {
	return &MorphFrameDelta{}
}

func (md *MorphFrameDelta) Copy() *MorphFrameDelta {
	return &MorphFrameDelta{
		framePosition: md.FramePosition().Copy(),
		frameRotation: md.FrameRotation().Copy(),
		frameScale:    md.FrameScale().Copy(),
	}
}
