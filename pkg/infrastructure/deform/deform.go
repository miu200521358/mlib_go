package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

func DeformBeforePhysics(
	appState state.IAppState, model *pmx.PmxModel, motion *vmd.VmdMotion,
) (*delta.VmdDeltas, *delta.RenderDeltas) {
	frame := int(appState.Frame())

	vmdDeltas := delta.NewVmdDeltas(model.Materials, model.Bones)
	vmdDeltas.Morphs = DeformMorph(model, motion.MorphFrames, frame, nil)
	vmdDeltas = DeformBoneByPhysicsFlag(model, motion, vmdDeltas, true, float64(frame), nil, false)

	renderDeltas := delta.NewRenderDeltas()
	renderDeltas.MeshDeltas = make([]*delta.MeshDelta, len(model.Materials.Data))
	for i, md := range vmdDeltas.Morphs.Materials.Data {
		renderDeltas.MeshDeltas[i] = delta.NewMeshDelta(md)
	}

	return vmdDeltas, renderDeltas
}
