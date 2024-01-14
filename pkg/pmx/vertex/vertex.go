package vertex

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"
)

type Vertex struct {
	index_model.IndexModel
	Index       int
	Position    *mvec3.T
	Normal      *mvec3.T
	UV          *mvec2.T
	ExtendedUVs *[]mvec4.T
	DeformType  DeformType
	Deform      Deform
	EdgeFactor  float64
}

func NewVertex(
	index int,
	position *mvec3.T,
	normal *mvec3.T,
	uv *mvec2.T,
	deformType DeformType,
	deform Deform,
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

// 頂点リスト
type Vertices struct {
	index_model.IndexModelCorrection[*Vertex]
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Vertex](),
	}
}
