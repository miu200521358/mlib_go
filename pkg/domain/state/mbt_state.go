//go:build windows
// +build windows

package state

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type IPhysics interface {
	ResetWorld()
	AddModel(modelIndex int, model *pmx.PmxModel)
	DeleteModel(modelIndex int)
	StepSimulation(timeStep float32)
	UpdateTransform(modelIndex int, rigidBodyBone *pmx.Bone, boneGlobalMatrix *mmath.MMat4, r *pmx.RigidBody)
	GetRigidBodyBoneMatrix(modelIndex int, rigidBody *pmx.RigidBody) *mmath.MMat4
}
