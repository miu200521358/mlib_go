//go:build windows
// +build windows

package renderer

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

func DrawInitialize(windowIndex int, pm *pmx.PmxModel) *Meshes {
	// モデルの初期化が終わっていない場合は初期化実行
	pm.ToonTextures = pmx.NewToonTextures()
	initToonTexturesGl(windowIndex, pm.ToonTextures)
	Meshes := NewMeshes(pm, windowIndex)
	return Meshes
}

// func (pm *PmxModel) Delete() {
// 	pm.Meshes.delete()
// 	pm.DeletePhysics()
// 	pm.DrawInitialized = false
// }

// func (pm *PmxModel) DeletePhysics() {
// 	// if pm.physics != nil {
// 	// 	pm.physics.DeleteRigidBodies(pm.Index)
// 	// 	pm.physics.DeleteJoints(pm.Index)
// 	// }
// }
