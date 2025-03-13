package viewer

import (
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
	"github.com/miu200521358/mlib_go/pkg/usecase/deform"
)

// loadModelRenderersメソッドを拡張
func (vw *ViewWindow) loadModelRenderers(shared *state.SharedState) {
	for i := range shared.ModelCount(vw.windowIndex) {
		vw.extendModelRenderers(shared)

		model := shared.LoadModel(vw.windowIndex, i)
		if vw.modelRenderers[i] != nil && vw.modelRenderers[i].Hash() != model.Hash() {
			vw.physics.DeleteModel(i)
			vw.modelRenderers[i].Delete()
			vw.modelRenderers[i] = nil
			vw.vmdDeltas[i] = nil
			// モデルが変わったらキャッシュもリセット
			if i < len(vw.speculativeCaches) && vw.speculativeCaches[i] != nil {
				vw.speculativeCaches[i].Reset()
			}
		}
		if vw.modelRenderers[i] == nil {
			vw.modelRenderers[i] = render.NewModelRenderer(vw.windowIndex, model)
			vw.physics.AddModel(i, model)
		}
	}
}

// モーションが変わった場合もキャッシュをリセット
func (vw *ViewWindow) loadMotions(shared *state.SharedState) {
	for i := range shared.MotionCount(vw.windowIndex) {
		vw.extendModelRenderers(shared)

		motion := shared.LoadMotion(vw.windowIndex, i)
		if vw.motions[i] != nil && vw.motions[i].Hash() != motion.Hash() {
			vw.motions[i] = nil
			// モーションが変わったらキャッシュもリセット
			if i < len(vw.speculativeCaches) && vw.speculativeCaches[i] != nil {
				vw.speculativeCaches[i].Reset()
			}
		}
		if vw.motions[i] == nil {
			vw.motions[i] = motion
		}
	}
}

func (vw *ViewWindow) extendModelRenderers(shared *state.SharedState) {
	for i := range shared.ModelCount(vw.windowIndex) {
		for i >= len(vw.modelRenderers) {
			vw.modelRenderers = append(vw.modelRenderers, nil)
		}
		for i >= len(vw.motions) {
			vw.motions = append(vw.motions, nil)
		}
		for i >= len(vw.vmdDeltas) {
			vw.vmdDeltas = append(vw.vmdDeltas, nil)
		}
		for i >= len(vw.speculativeCaches) {
			vw.speculativeCaches = append(vw.speculativeCaches, deform.NewSpeculativeDeformCache())
		}
	}
}
