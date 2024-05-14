package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type Vertex struct {
	*mcore.IndexModel
	Position    mmath.MVec3     // 頂点位置
	Normal      mmath.MVec3     // 頂点法線
	UV          mmath.MVec2     // UV
	ExtendedUVs []*mmath.MVec4  // 追加UV
	DeformType  DeformType      // ウェイト変形方式
	Deform      DeformInterface // デフォーム
	EdgeFactor  float64         // エッジ倍率
}

func NewVertex() *Vertex {
	v := &Vertex{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Position:    mmath.NewMVec3(),
		Normal:      mmath.NewMVec3(),
		UV:          mmath.NewMVec2(),
		ExtendedUVs: make([]*mmath.MVec4, 0),
		DeformType:  BDEF1,
		Deform:      NewBdef1(0),
		EdgeFactor:  0.0,
	}
	return v
}

func (v *Vertex) Copy() mcore.IIndexModel {
	copied := NewTexture()
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return copied
}

// 頂点リスト
type Vertices struct {
	*mcore.IndexModels[*Vertex]
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModels: mcore.NewIndexModels[*Vertex](),
	}
}
