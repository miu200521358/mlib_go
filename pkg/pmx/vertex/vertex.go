package vertex

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"

)

type T struct {
	index_model.T
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
) *T {
	v := &T{
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

func (m *T) Copy() *T {
	copied := *m
	return &copied
}

// 頂点リスト
type C struct {
	index_model.C
	data    map[int]*T
	Indexes []int
}

func NewVertices() *C {
	return &C{
		data:    make(map[int]*T),
		Indexes: make([]int, 0),
	}
}
