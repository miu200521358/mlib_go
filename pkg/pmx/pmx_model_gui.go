//go:build windows
// +build windows

package pmx

import (
	"embed"

	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

func (pm *PmxModel) DrawInitialize(windowIndex int, resourceFiles embed.FS, physics *mphysics.MPhysics) {
	if !pm.DrawInitialized {
		// モデルの初期化が終わっていない場合は初期化実行
		pm.InitDraw(windowIndex, resourceFiles)
		pm.InitPhysics(physics)
		pm.DrawInitialized = true
	}
}

func (pm *PmxModel) InitDraw(windowIndex int, resourceFiles embed.FS) {
	pm.ToonTextures.initGl(windowIndex, resourceFiles)
	pm.Meshes = NewMeshes(pm, windowIndex, resourceFiles)

	// 位置マッピングのセットアップ
	pm.Vertices.SetupMapKeys()
	pm.Bones.SetupMapKeys()
	pm.RigidBodies.SetupMapKeys()
	pm.Joints.SetupMapKeys()
}

func (pm *PmxModel) InitPhysics(physics *mphysics.MPhysics) {
	pm.physics = physics
	pm.RigidBodies.initPhysics(physics)
	pm.Joints.initPhysics(physics, pm.RigidBodies)
}

func (pm *PmxModel) Delete() {
	pm.Meshes.delete()
	pm.DeletePhysics()
	pm.DrawInitialized = false
}

func (pm *PmxModel) DeletePhysics() {
	if pm.physics != nil {
		pm.physics.DeleteRigidBodies()
		pm.physics.DeleteJoints()
	}
}
