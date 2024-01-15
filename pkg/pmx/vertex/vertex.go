package vertex

import (
	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"

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
	ExtendedUVs *[]mvec4.T
	// ウェイト変形方式
	DeformType DeformType
	// デフォーム
	Deform DeformInterface
	// エッジ倍率
	EdgeFactor float64
}

func NewVertex(
	index int,
	position *mvec3.T,
	normal *mvec3.T,
	uv *mvec2.T,
	deformType DeformType,
	deform DeformInterface,
	edgeFactor float64,
) *Vertex {
	v := &Vertex{
		IndexModel:  &index_model.IndexModel{Index: index},
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
	*index_model.IndexModelCorrection[*Vertex]
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Vertex](),
	}
}
