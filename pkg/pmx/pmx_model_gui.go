//go:build windows
// +build windows

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

func (pm *PmxModel) DrawInitialize(windowIndex int, physics *mphysics.MPhysics) {
	if !pm.DrawInitialized {
		// モデルの初期化が終わっていない場合は初期化実行
		pm.InitDraw(windowIndex)
		pm.InitPhysics(physics)
		pm.DrawInitialized = true
	}
}

func (pm *PmxModel) InitDraw(windowIndex int) {
	pm.ToonTextures.initGl(windowIndex)
	pm.Meshes = NewMeshes(pm, windowIndex)
}

func (pm *PmxModel) InitPhysics(physics *mphysics.MPhysics) {
	pm.physics = physics
	pm.RigidBodies.initPhysics(pm.Index, physics)
	pm.Joints.initPhysics(pm.Index, physics, pm.RigidBodies)
}

func (pm *PmxModel) Delete() {
	pm.Meshes.delete()
	pm.DeletePhysics()
	pm.DrawInitialized = false
}

func (pm *PmxModel) DeletePhysics() {
	if pm.physics != nil {
		pm.physics.DeleteRigidBodies(pm.Index)
		pm.physics.DeleteJoints(pm.Index)
	}
}
