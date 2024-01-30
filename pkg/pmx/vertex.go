package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type Vertex struct {
	*mcore.IndexModel
	Position    *mmath.MVec3    // 頂点位置
	Normal      *mmath.MVec3    // 頂点法線
	UV          *mmath.MVec2    // UV
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
		UV:          &mmath.MVec2{},
		ExtendedUVs: []*mmath.MVec4{},
		DeformType:  BDEF1,
		Deform:      NewBdef1(0),
		EdgeFactor:  0.0,
	}
	return v
}

func (v *Vertex) GL() []float32 {
	p := v.Position.GL()
	n := v.Normal.GL()
	eu := [2]float32{0.0, 0.0}
	if len(v.ExtendedUVs) > 0 {
		eu[0] = float32(v.ExtendedUVs[0].GetX())
		eu[1] = float32(v.ExtendedUVs[0].GetY())
	}
	d := v.Deform.NormalizedDeform()
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(v.UV.GetX()), float32(v.UV.GetY()), // UV
		eu[0], eu[1], // 追加UV
		float32(v.EdgeFactor),  // エッジ倍率
		d[0], d[1], d[2], d[3], // デフォームボーンINDEX
		d[4], d[5], d[6], d[7], // デフォームボーンウェイト
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
