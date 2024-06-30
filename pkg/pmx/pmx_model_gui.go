//go:build windows
// +build windows

package pmx

import (
	"embed"

	"github.com/miu200521358/mlib_go/pkg/mphysics"
)

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
	pm.RigidBodies.initPhysics(physics)
	pm.Joints.initPhysics(physics, pm.RigidBodies)
}

func (pm *PmxModel) Delete(physics *mphysics.MPhysics) {
	pm.Meshes.delete()
	pm.DeletePhysics(physics)
}

func (pm *PmxModel) DeletePhysics(physics *mphysics.MPhysics) {
	physics.DeleteRigidBodies()
	physics.DeleteJoints()
}
