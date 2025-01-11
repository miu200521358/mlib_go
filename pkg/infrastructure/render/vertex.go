//go:build windows
// +build windows

package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

func newVertexGl(v *pmx.Vertex) []float32 {
	p := mgl.NewGlVec3(v.Position)
	n := mgl.NewGlVec3(v.Normal)
	eu := [2]float32{0.0, 0.0}
	if len(v.ExtendedUvs) > 0 {
		eu[0] = float32(v.ExtendedUvs[0].X)
		eu[1] = float32(v.ExtendedUvs[0].Y)
	}
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == pmx.SDEF))
	sdefC, sdefR0, sdefR1 := getSdefParams(v.Deform)
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(v.Uv.X), float32(v.Uv.Y), // UV
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

func newVertexNormalGl(v *pmx.Vertex) []float32 {
	p := mgl.NewGlVec3(v.Position)
	n := mgl.NewGlVec3(v.Normal.MuledScalar(0.5))
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == pmx.SDEF))
	sdefC, sdefR0, sdefR1 := getSdefParams(v.Deform)
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

func newSelectedVertexGl(v *pmx.Vertex) []float32 {
	p := mgl.NewGlVec3(v.Position)
	n := mgl.NewGlVec3(v.Normal)
	d := v.Deform.NormalizedDeform()
	s := float32(mmath.BoolToInt(v.DeformType == pmx.SDEF))
	sdefC, sdefR0, sdefR1 := getSdefParams(v.Deform)
	return []float32{
		p[0], p[1], p[2], // 位置
		n[0], n[1], n[2], // 法線
		float32(-0.1), float32(0), // UV (Xは明示的にマイナスにしておく)
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

// SDEF用パラメーターを返す
func getSdefParams(d pmx.IDeform) (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	// pmx.Sdefにキャストできる場合
	if _, ok := d.(*pmx.Sdef); ok {
		s := d.(*pmx.Sdef)
		// CがR0とR1より先にいかないよう、重みに基づいて補正
		copiedSdefR0 := s.SdefR0.Copy()
		copiedSdefR1 := s.SdefR1.Copy()
		copiedSdefCR0 := s.SdefC.Copy()
		copiedSdefCR1 := s.SdefC.Copy()

		weight := copiedSdefR0.MulScalar(s.AllWeights()[0]).Add(copiedSdefR1.MulScalar(1 - s.AllWeights()[0]))
		sdefR0 := copiedSdefCR0.Add(s.SdefR0).Sub(weight)
		sdefR1 := copiedSdefCR1.Add(s.SdefR1).Sub(weight)

		return mgl.NewGlVec3(s.SdefC), mgl.NewGlVec3(sdefR0), mgl.NewGlVec3(sdefR1)
	}
	return mgl32.Vec3{}, mgl32.Vec3{}, mgl32.Vec3{}
}
