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

func (v *Vertex) GL() []float32 {
	p := v.Position.GL()
	// n := v.Normal.GL()
	// d := v.Deform.NormalizedDeform()
	return []float32{
		p[0], p[1], p[2],
		// n[0], n[1], n[2],
		// d[0], d[1], d[2], d[3], d[4], d[5], d[6], d[7],
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
