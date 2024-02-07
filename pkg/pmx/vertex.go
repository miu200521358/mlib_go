package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"

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
	s := float32(mutils.BoolToInt(v.DeformType == SDEF))
	// s = 0.0
	sdefC := [3]float32{0.0, 0.0, 0.0}
	sdefR0 := [3]float32{0.0, 0.0, 0.0}
	sdefR1 := [3]float32{0.0, 0.0, 0.0}
	if v.DeformType == SDEF {
		sdefC[0] = float32(v.Deform.(*Sdef).SdefC.GetX())
		sdefC[1] = float32(v.Deform.(*Sdef).SdefC.GetY())
		sdefC[2] = float32(v.Deform.(*Sdef).SdefC.GetZ())

		// c := v.Deform.(*Sdef).SdefC
		r0 := v.Deform.(*Sdef).SdefR0
		r1 := v.Deform.(*Sdef).SdefR1
		// rMuled0 := r0.MuledScalar(v.Deform.GetAllWeights()[0])
		// rMuled1 := r1.MuledScalar(1 - v.Deform.GetAllWeights()[0])
		// rOffset := rMuled0.Added(&rMuled1)
		// r00 := c.Added(r0)
		// r01 := c.Added(r1)
		// normalizedR0 := r00.Subed(&rOffset)
		// normalizedR1 := r01.Subed(&rOffset)

		// sdefR0[0] = float32(normalizedR0.GetX())
		// sdefR0[1] = float32(normalizedR0.GetY())
		// sdefR0[2] = float32(normalizedR0.GetZ())

		// sdefR1[0] = float32(normalizedR1.GetX())
		// sdefR1[1] = float32(normalizedR1.GetY())
		// sdefR1[2] = float32(normalizedR1.GetZ())

		sdefR0[0] = float32(r0.GetX())
		sdefR0[1] = float32(r0.GetY())
		sdefR0[2] = float32(r0.GetZ())

		sdefR1[0] = float32(r1.GetX())
		sdefR1[1] = float32(r1.GetY())
		sdefR1[2] = float32(r1.GetZ())
	}
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(v.UV.GetX()), float32(v.UV.GetY()), // UV
		eu[0], eu[1], // 追加UV
		float32(v.EdgeFactor),  // エッジ倍率
		d[0], d[1], d[2], d[3], // デフォームボーンINDEX
		d[4], d[5], d[6], d[7], // デフォームボーンウェイト
		s,                            // SDEFであるか否か
		sdefC[0], sdefC[1], sdefC[2], // SDEF-C
		sdefR0[0], sdefR0[1], sdefR0[2], // SDEF-R0
		sdefR1[0], sdefR1[1], sdefR1[2], // SDEF-R1
	}
}

// 頂点リスト
type Vertices struct {
	*mcore.IndexModelCorrection[*Vertex]
	SDEFs []int // SDEF対象頂点INDEXリスト
}

func NewVertices() *Vertices {
	return &Vertices{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*Vertex](),
		SDEFs:                []int{},
	}
}
