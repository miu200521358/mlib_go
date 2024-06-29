package pmx

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type Vertex struct {
	*mcore.IndexModel
	Position    *mmath.MVec3   // 頂点位置
	Normal      *mmath.MVec3   // 頂点法線
	UV          *mmath.MVec2   // UV
	ExtendedUVs []*mmath.MVec4 // 追加UV
	DeformType  DeformType     // ウェイト変形方式
	Deform      IDeform        // デフォーム
	EdgeFactor  float64        // エッジ倍率
}

func NewVertex() *Vertex {
	v := &Vertex{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Position:    mmath.NewMVec3(),
		Normal:      mmath.NewMVec3(),
		UV:          mmath.NewMVec2(),
		ExtendedUVs: make([]*mmath.MVec4, 0),
		DeformType:  BDEF1,
		Deform:      nil,
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
	Positions []*mmath.MVec3
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModels: mcore.NewIndexModels[*Vertex](func() *Vertex { return nil }),
	}
}

func (vs *Vertices) setup() {
	vs.Positions = make([]*mmath.MVec3, vs.IndexModels.Len())
	for i := range vs.IndexModels.Len() {
		vs.Positions[i] = vs.IndexModels.Get(i).Position
	}
}
