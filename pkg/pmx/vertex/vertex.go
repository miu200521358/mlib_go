package vertex

import (
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"
	"github.com/miu200521358/mlib_go/pkg/pmx/deform"

)

type Vertex struct {
	Index       int
	Position    *mvec3.T
	Normal      *mvec3.T
	UV          *mvec2.T
	ExtendedUVs *[]mvec4.T
	DeformType  deform.DeformType
	Deform      deform.T
	EdgeFactor  float64
}

func NewVertex(
	index int,
	position *mvec3.T,
	normal *mvec3.T,
	uv *mvec2.T,
	deformType deform.DeformType,
	deform deform.T,
	edgeFactor float64,
) *Vertex {
	v := &Vertex{
		Index:       index,
		Position:    position,
		Normal:      normal,
		UV:          uv,
		ExtendedUVs: &[]mvec4.T{},
		DeformType:  deformType,
		Deform:      deform,
		EdgeFactor:  edgeFactor,
	}
	return v
}
