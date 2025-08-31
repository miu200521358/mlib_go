package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type RigidBodyFrame struct {
	*BaseFrame              // キーフレ
	Size       *mmath.MVec3 // サイズ
	Mass       float64      // 質量
}

func NewRigidBodyFrameByValues(index float32, size *mmath.MVec3, mass float64) *RigidBodyFrame {
	return &RigidBodyFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Size:      size, // サイズ
		Mass:      mass, // 質量
	}
}

func (mf *RigidBodyFrame) Copy() IBaseFrame {
	return NewRigidBodyFrameByValues(mf.Index(), mf.Size.Copy(), mf.Mass)
}

func (nextMf *RigidBodyFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*RigidBodyFrame)

	prevIndex := prevMf.Index()
	nextIndex := nextMf.Index()

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	size := prevMf.Size.Lerp(nextMf.Size, ry)
	mass := mmath.Lerp(prevMf.Mass, nextMf.Mass, ry)

	return NewRigidBodyFrameByValues(index, size, mass)
}

func (mf *RigidBodyFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}
