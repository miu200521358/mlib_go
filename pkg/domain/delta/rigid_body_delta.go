package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

// RigidBodyDelta は1つのボーンにおける変形（ポジション・回転・スケールなど）の差分を表す
type RigidBodyDelta struct {
	RigidBody *pmx.RigidBody
	Frame     float32
	Size      *mmath.MVec3 // 剛体の大きさ
	Mass      float64      // 剛体の質量
}

// NewRigidBodyDelta は新規の RigidBodyDelta を生成するコンストラクタ
func NewRigidBodyDelta(rigidBody *pmx.RigidBody, frame float32) *RigidBodyDelta {
	return &RigidBodyDelta{
		RigidBody: rigidBody,
		Frame:     frame,
	}
}

// NewRigidBodyDelta は新規の RigidBodyDelta を生成するコンストラクタ
func NewRigidBodyDeltaByValue(
	rigidBody *pmx.RigidBody, frame float32, size *mmath.MVec3, mass float64,
) *RigidBodyDelta {
	return &RigidBodyDelta{
		RigidBody: rigidBody,
		Frame:     frame,
		Size:      size,
		Mass:      mass,
	}
}
