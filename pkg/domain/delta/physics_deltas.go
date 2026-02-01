// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
)

// PhysicsDeltas は物理差分をまとめる。
type PhysicsDeltas struct {
	frame       mtime.Frame
	modelHash   string
	motionHash  string
	RigidBodies *RigidBodyDeltas
	Joints      *JointDeltas
}

// NewPhysicsDeltas はPhysicsDeltasを生成する。
func NewPhysicsDeltas(
	frame mtime.Frame,
	rigidBodies *collection.NamedCollection[*model.RigidBody],
	joints *collection.NamedCollection[*model.Joint],
	modelHash, motionHash string,
) *PhysicsDeltas {
	return &PhysicsDeltas{
		frame:       frame,
		modelHash:   modelHash,
		motionHash:  motionHash,
		RigidBodies: NewRigidBodyDeltas(rigidBodies),
		Joints:      NewJointDeltas(joints),
	}
}

// Frame はフレーム番号を返す。
func (p *PhysicsDeltas) Frame() mtime.Frame {
	if p == nil {
		return 0
	}
	return p.frame
}

// SetFrame はフレーム番号を設定する。
func (p *PhysicsDeltas) SetFrame(frame mtime.Frame) {
	if p == nil {
		return
	}
	p.frame = frame
}

// ModelHash はモデルハッシュを返す。
func (p *PhysicsDeltas) ModelHash() string {
	if p == nil {
		return ""
	}
	return p.modelHash
}

// SetModelHash はモデルハッシュを設定する。
func (p *PhysicsDeltas) SetModelHash(hash string) {
	if p == nil {
		return
	}
	p.modelHash = hash
}

// MotionHash はモーションハッシュを返す。
func (p *PhysicsDeltas) MotionHash() string {
	if p == nil {
		return ""
	}
	return p.motionHash
}

// SetMotionHash はモーションハッシュを設定する。
func (p *PhysicsDeltas) SetMotionHash(hash string) {
	if p == nil {
		return
	}
	p.motionHash = hash
}
