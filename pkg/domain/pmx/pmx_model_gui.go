//go:build windows
// +build windows

package pmx

func (pm *PmxModel) DrawInitialize(windowIndex int) {
	if !pm.DrawInitialized {
		// モデルの初期化が終わっていない場合は初期化実行
		pm.InitDraw(windowIndex)
		pm.DrawInitialized = true
	}
}

func (pm *PmxModel) InitDraw(windowIndex int) {
	pm.ToonTextures = NewToonTextures()
	pm.ToonTextures.initGl(windowIndex)
	pm.Meshes = NewMeshes(pm, windowIndex)
}

func (pm *PmxModel) Delete() {
	pm.Meshes.delete()
	pm.DeletePhysics()
	pm.DrawInitialized = false
}

func (pm *PmxModel) DeletePhysics() {
	// if pm.physics != nil {
	// 	pm.physics.DeleteRigidBodies(pm.Index)
	// 	pm.physics.DeleteJoints(pm.Index)
	// }
}
