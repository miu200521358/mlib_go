//go:build windows
// +build windows

package pmx

import "github.com/miu200521358/mlib_go/pkg/mutils"

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
	sdefC, sdefR0, sdefR1 := v.Deform.GetSdefParams()
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
		0.0, 0.0, 0.0, // 頂点モーフ
		0.0, 0.0, 0.0, 0.0, // UVモーフ
		0.0, 0.0, 0.0, 0.0, // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}