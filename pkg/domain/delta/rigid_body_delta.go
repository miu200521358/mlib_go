// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
)

// RigidBodyDelta は剛体差分を表す。
type RigidBodyDelta struct {
	RigidBody *model.RigidBody
	Frame     mtime.Frame
	Position  mmath.Vec3
	Size      mmath.Vec3
	Mass      float64
}

// NewRigidBodyDelta はRigidBodyDeltaを生成する。
func NewRigidBodyDelta(rigidBody *model.RigidBody, frame mtime.Frame) *RigidBodyDelta {
	if rigidBody == nil {
		return nil
	}
	return &RigidBodyDelta{
		RigidBody: rigidBody,
		Frame:     frame,
		Position:  rigidBody.Position,
		Size:      rigidBody.Size,
		Mass:      rigidBody.Param.Mass,
	}
}

// NewRigidBodyDeltaByValue は値を指定してRigidBodyDeltaを生成する。
func NewRigidBodyDeltaByValue(
	rigidBody *model.RigidBody,
	frame mtime.Frame,
	position mmath.Vec3,
	size mmath.Vec3,
	mass float64,
) *RigidBodyDelta {
	if rigidBody == nil {
		return nil
	}
	return &RigidBodyDelta{
		RigidBody: rigidBody,
		Frame:     frame,
		Position:  position,
		Size:      size,
		Mass:      mass,
	}
}
