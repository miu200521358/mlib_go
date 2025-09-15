package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type RigidBodyFrame struct {
	*BaseFrame              // キーフレ
	Position   *mmath.MVec3 // 位置
	Size       *mmath.MVec3 // サイズ
	Mass       float64      // 質量
}

func NewRigidBodyFrameByValues(index float32, position, size *mmath.MVec3, mass float64) *RigidBodyFrame {
	return &RigidBodyFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
		Position:  position,
		Size:      size, // サイズ
		Mass:      mass, // 質量
	}
}

func (mf *RigidBodyFrame) Copy() IBaseFrame {
	return NewRigidBodyFrameByValues(mf.Index(), mf.Position.Copy(), mf.Size.Copy(), mf.Mass)
}

func (nextMf *RigidBodyFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevMf := prevFrame.(*RigidBodyFrame)

	prevIndex := prevMf.Index()
	nextIndex := nextMf.Index()

	ry := float64(index-prevIndex) / float64(nextIndex-prevIndex)
	position := prevMf.Position.Lerp(nextMf.Position, ry)
	size := prevMf.Size.Lerp(nextMf.Size, ry)
	mass := mmath.Lerp(prevMf.Mass, nextMf.Mass, ry)

	// mlog.I("RigidBodyFrame Lerp: prevIndex=%.2f, nextIndex=%.2f, index=%.2f, ry=%.2f => pos=(%.2f, %.2f, %.2f), size=(%.2f, %.2f, %.2f), mass=%.2f",
	// 	prevIndex, nextIndex, index, ry,
	// 	position.X, position.Y, position.Z,
	// 	size.X, size.Y, size.Z,
	// 	mass,
	// )

	return NewRigidBodyFrameByValues(index, position, size, mass)
}

func (mf *RigidBodyFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}
