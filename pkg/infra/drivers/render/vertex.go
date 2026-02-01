//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
)

// 頂点データのサイズ定数
//
//	位置 (position): 3要素
//	法線 (normal): 3要素
//	UV: 2要素
//	追加UV: 2要素
//	エッジ倍率: 1要素
//	デフォームボーンINDEX: 4要素
//	デフォームボーンウェイト: 4要素
//	SDEFフラグ: 1要素
//	SDEF-C: 3要素
//	SDEF-R0: 3要素
//	SDEF-R1: 3要素
//	頂点モーフ: 3要素
//	UVモーフ: 4要素
//	追加UV1モーフ: 4要素
//	変形後頂点モーフ: 3要素
const vertexDataSize = 43

// newVertexGl はOpenGL向け頂点データを生成する。
func newVertexGl(v *model.Vertex) []float32 {
	p := mgl.NewGlVec3(&v.Position)
	n := mgl.NewGlVec3(&v.Normal)
	eu := [2]float32{0.0, 0.0}
	if len(v.ExtendedUvs) > 0 {
		eu[0] = float32(v.ExtendedUvs[0].X)
		eu[1] = float32(v.ExtendedUvs[0].Y)
	}
	d := packDeform(v.Deform)
	s := float32(mmath.BoolToInt(v.DeformType == model.SDEF))
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

// newVertexNormalGl は法線描画用の頂点データを生成する。
func newVertexNormalGl(v *model.Vertex) []float32 {
	p := mgl.NewGlVec3(&v.Position)
	normal := v.Normal.MuledScalar(0.5)
	n := mgl.NewGlVec3(&normal)
	d := packDeform(v.Deform)
	s := float32(mmath.BoolToInt(v.DeformType == model.SDEF))
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

// newSelectedVertexGl は選択頂点表示用データを生成する。
func newSelectedVertexGl(v *model.Vertex) []float32 {
	p := mgl.NewGlVec3(&v.Position)
	n := mgl.NewGlVec3(&v.Normal)
	d := packDeform(v.Deform)
	s := float32(mmath.BoolToInt(v.DeformType == model.SDEF))
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

// getSdefParams はSDEF用パラメータを取得する。
func getSdefParams(d model.IDeform) (mgl32.Vec3, mgl32.Vec3, mgl32.Vec3) {
	s, ok := d.(*model.Sdef)
	if ok {
		weights := s.Weights()
		weight0 := 0.0
		if len(weights) > 0 {
			weight0 = weights[0]
		}

		weight := s.SdefR0.MuledScalar(weight0).Added(s.SdefR1.MuledScalar(1 - weight0))
		sdefR0 := s.SdefC.Added(s.SdefR0).Subed(weight)
		sdefR1 := s.SdefC.Added(s.SdefR1).Subed(weight)

		return mgl.NewGlVec3(&s.SdefC), mgl.NewGlVec3(&sdefR0), mgl.NewGlVec3(&sdefR1)
	}
	return mgl32.Vec3{}, mgl32.Vec3{}, mgl32.Vec3{}
}

// packDeform はデフォーム情報を4ボーン+4ウェイトに詰める。
func packDeform(d model.IDeform) [8]float32 {
	var packed [8]float32
	if d == nil {
		return packed
	}
	indexes := d.Indexes()
	for i := 0; i < len(indexes) && i < 4; i++ {
		packed[i] = float32(indexes[i])
	}
	weights := d.Weights()
	for i := 0; i < len(weights) && i < 4; i++ {
		packed[4+i] = float32(weights[i])
	}
	return packed
}
