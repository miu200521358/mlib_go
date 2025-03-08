package viewer

import (
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
)

func (vw *ViewWindow) loadModelRenderers(shared *state.SharedState) {
	vw.MakeContextCurrent()

	for i := range shared.ModelCount(vw.windowIndex) {
		for i >= len(vw.modelRenderers) {
			vw.modelRenderers = append(vw.modelRenderers, nil)
		}
		vw.modelRenderers[i] = render.NewModelRenderer(vw.windowIndex, shared.LoadModel(vw.windowIndex, i))
	}
}

func (vw *ViewWindow) loadMotions(shared *state.SharedState) {
	vw.MakeContextCurrent()

	for i := range shared.MotionCount(vw.windowIndex) {
		for i >= len(vw.motions) {
			vw.motions = append(vw.motions, nil)
			vw.vmdDeltas = append(vw.vmdDeltas, nil)
		}
		vw.motions[i] = shared.LoadMotion(vw.windowIndex, i)
	}
}
