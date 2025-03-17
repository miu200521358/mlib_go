package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type Vertex struct {
	index           int            // 頂点INDEX
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
		index:           -1,
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

func (vertex *Vertex) Index() int {
	return vertex.index
}

func (vertex *Vertex) SetIndex(index int) {
	vertex.index = index
}

func (vertex *Vertex) IsValid() bool {
	return vertex != nil && vertex.Index() >= 0
}

func (vertex *Vertex) Copy() core.IIndexModel {
	var copiedExtendedUvs []*mmath.MVec4
	for _, uv := range vertex.ExtendedUvs {
		copiedExtendedUvs = append(copiedExtendedUvs, uv.Copy())
	}

	return &Vertex{
		index:           vertex.index,
		Position:        vertex.Position.Copy(),
		Normal:          vertex.Normal.Copy(),
		Uv:              vertex.Uv.Copy(),
		ExtendedUvs:     copiedExtendedUvs,
		DeformType:      vertex.DeformType,
		Deform:          vertex.Deform,
		EdgeFactor:      vertex.EdgeFactor,
		MaterialIndexes: mmath.DeepCopy(vertex.MaterialIndexes),
	}
}
