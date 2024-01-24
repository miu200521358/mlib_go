package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type Vertex struct {
	*mcore.IndexModel
	// 頂点位置
	Position *mmath.MVec3
	// 頂点法線
	Normal *mmath.MVec3
	// UV
	UV *mmath.MVec2
	// 追加UV
	ExtendedUVs []mmath.MVec4
	// ウェイト変形方式
	DeformType DeformType
	// デフォーム
	Deform DeformInterface
	// エッジ倍率
	EdgeFactor float64
}

type VertexGL struct {
	Position [3]float32
	Normal   [3]float32
	UV       [2]float32
}

func NewVertex() *Vertex {
	v := &Vertex{
		IndexModel:  &mcore.IndexModel{Index: -1},
		Position:    &mmath.MVec3{},
		Normal:      &mmath.MVec3{},
		UV:          &mmath.MVec2{},
		ExtendedUVs: []mmath.MVec4{},
		DeformType:  BDEF1,
		Deform:      NewBdef1(0),
		EdgeFactor:  0.0,
	}
	return v
}

func (v *Vertex) GL() *VertexGL {
	return &VertexGL{
		Position: v.Position.GL(),
		Normal:   v.Normal.GL(),
		UV:       v.UV.GL(),
	}
}

// 頂点リスト
type Vertices struct {
	*mcore.IndexModelCorrection[*Vertex]
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Vertex](),
	}
}
