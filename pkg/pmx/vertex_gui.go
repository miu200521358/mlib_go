//go:build windows
// +build windows

package pmx

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

func (v *Vertex) GL() []float32 {
	p := v.Position.GL()
	n := v.Normal.GL()
	eu := [2]float32{0.0, 0.0}
	if len(v.ExtendedUvs) > 0 {
		eu[0] = float32(v.ExtendedUvs[0].GetX())
		eu[1] = float32(v.ExtendedUvs[0].GetY())
	}
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == SDEF))
	sdefC, sdefR0, sdefR1 := v.Deform.GetSdefParams()
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(v.Uv.GetX()), float32(v.Uv.GetY()), // UV
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

func (v *Vertex) NormalGL() []float32 {
	p := v.Position.GL()
	n := v.Normal.MuledScalar(0.5).GL()
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == SDEF))
	sdefC, sdefR0, sdefR1 := v.Deform.GetSdefParams()
	return []float32{
		p[0] + n[0], p[1] + n[1], p[2] + n[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(0), float32(0), // UV
		float32(0), float32(0), // 追加UV
		float32(0),             // エッジ倍率
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

func (v *Vertex) WireGL() []float32 {
	p := v.Position.GL()
	n := v.Normal.GL()
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == SDEF))
	sdefC, sdefR0, sdefR1 := v.Deform.GetSdefParams()
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(1), float32(0), // UV(Xは明示的に1)
		float32(0), float32(0), // 追加UV
		float32(0),             // エッジ倍率
		d[0], d[1], d[2], d[3], // デフォームボーンINDEX
		d[4], d[5], d[6], d[7], // デフォームボーンウェイト
		s,                            // SDEFであるか否か
		sdefC[0], sdefC[1], sdefC[2], // SDEF-C
		sdefR0[0], sdefR0[1], sdefR0[2], // SDEF-R0
		sdefR1[0], sdefR1[1], sdefR1[2], // SDEF-R1
		0.0, 0.0, 0.0, // 頂点モーフ
		1.0, 0.0, 0.0, 0.0, // UVモーフ(Xは明示的に1にしてフラグを立てる（表示状態）)
		0.0, 0.0, 0.0, 0.0, // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}

func (v *Vertex) SelectedGL() []float32 {
	p := v.Position.GL()
	n := v.Normal.GL()
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == SDEF))
	sdefC, sdefR0, sdefR1 := v.Deform.GetSdefParams()
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(-1), float32(0), // UV(Xは明示的に-1)
		float32(0), float32(0), // 追加UV
		float32(0),             // エッジ倍率
		d[0], d[1], d[2], d[3], // デフォームボーンINDEX
		d[4], d[5], d[6], d[7], // デフォームボーンウェイト
		s,                            // SDEFであるか否か
		sdefC[0], sdefC[1], sdefC[2], // SDEF-C
		sdefR0[0], sdefR0[1], sdefR0[2], // SDEF-R0
		sdefR1[0], sdefR1[1], sdefR1[2], // SDEF-R1
		0.0, 0.0, 0.0, // 頂点モーフ
		1.0, 0.0, 0.0, 0.0, // UVモーフ(Xは明示的に1にしてフラグを立てる（選択状態）)
		0.0, 0.0, 0.0, 0.0, // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}
