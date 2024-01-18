package vertex

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"
	"github.com/miu200521358/mlib_go/pkg/pmx/vertex/deform"
)

type Vertex struct {
	*index_model.IndexModel
	// 頂点位置
	Position *mvec3.T
	// 頂点法線
	Normal *mvec3.T
	// UV
	UV *mvec2.T
	// 追加UV
	ExtendedUVs []mvec4.T
	// ウェイト変形方式
	DeformType deform.DeformType
	// デフォーム
	Deform deform.DeformInterface
	// エッジ倍率
	EdgeFactor float64
}

func NewVertex() *Vertex {
	v := &Vertex{
		IndexModel:  &index_model.IndexModel{Index: -1},
		Position:    &mvec3.T{},
		Normal:      &mvec3.T{},
		UV:          &mvec2.T{},
		ExtendedUVs: []mvec4.T{},
		DeformType:  deform.BDEF1,
		Deform:      deform.NewBdef1(0),
		EdgeFactor:  0.0,
	}
	return v
}

// 頂点リスト
type Vertices struct {
	*index_model.IndexModelCorrection[*Vertex]
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Vertex](),
	}
}
