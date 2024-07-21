package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type Vertex struct {
	*core.IndexModel
	Position        *mmath.MVec3   // 頂点位置
	Normal          *mmath.MVec3   // 頂点法線
	Uv              *mmath.MVec2   // UV
	ExtendedUvs     []*mmath.MVec4 // 追加UV
	DeformType      DeformType     // ウェイト変形方式
	Deform          IDeform        // デフォーム
	EdgeFactor      float64        // エッジ倍率
	MaterialIndexes []int          // 割り当て材質インデックス
}

func NewVertex() *Vertex {
	v := &Vertex{
		IndexModel:      core.NewIndexModel(-1),
		Position:        mmath.NewMVec3(),
		Normal:          mmath.NewMVec3(),
		Uv:              mmath.NewMVec2(),
		ExtendedUvs:     make([]*mmath.MVec4, 0),
		DeformType:      BDEF1,
		Deform:          nil,
		EdgeFactor:      0.0,
		MaterialIndexes: make([]int, 0),
	}
	return v
}

func (vertex *Vertex) Copy() core.IIndexModel {
	copied := NewVertex()
	copier.CopyWithOption(copied, vertex, copier.Option{DeepCopy: true})
	return copied
}

// 頂点リスト
type Vertices struct {
	*core.IndexModels[*Vertex]
}

func NewVertices(count int) *Vertices {
	return &Vertices{
		IndexModels: core.NewIndexModels[*Vertex](count, func() *Vertex { return nil }),
	}
}
