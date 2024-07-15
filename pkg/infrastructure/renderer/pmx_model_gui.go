//go:build windows
// +build windows

package renderer

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

func DrawInitialize(windowIndex int, pm *pmx.PmxModel) *Meshes {
	if !pm.DrawInitialized {
		// モデルの初期化が終わっていない場合は初期化実行
		pm.ToonTextures = pmx.NewToonTextures()
		initToonTexturesGl(windowIndex, pm.ToonTextures)
		Meshes := NewMeshes(pm, windowIndex)
		pm.DrawInitialized = true
		return Meshes
	}
	return nil
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
