//go:build windows
// +build windows

package renderer

import (
	"slices"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

func newVertexMorphDeltasGl(mds *delta.VertexMorphDeltas) ([]int, [][]float32) {
	vertices := make([][]float32, 0)
	indices := make([]int, 0)
	for i, md := range mds.Data {
		vertices = append(vertices, newVertexMorphDeltaGl(md))
		indices = append(indices, i)
	}
	return indices, vertices
}

func newVertexMorphDeltaGl(md *delta.VertexMorphDelta) []float32 {
	var p0, p1, p2 float32
	if md.Position != nil {
		p := mgl.NewGlVec3(md.Position)
		p0, p1, p2 = p[0], p[1], p[2]
	}
	var ap0, ap1, ap2 float32
	if md.AfterPosition != nil {
		ap := mgl.NewGlVec3(md.AfterPosition)
		ap0, ap1, ap2 = ap[0], ap[1], ap[2]
	}
	// UVは符号関係ないのでそのまま取得する
	var u0x, u0y, u1x, u1y float32
	if md.Uv != nil {
		u0x = float32(md.Uv.X)
		u0y = float32(md.Uv.Y)
	}
	if md.Uv1 != nil {
		u1x = float32(md.Uv1.X)
		u1y = float32(md.Uv1.Y)
	}
	return []float32{
		p0, p1, p2,
		u0x, u0y, 0, 0,
		u1x, u1y, 0, 0,
		ap0, ap1, ap2,
	}
}

func newSelectedVertexMorphDeltasGL(
	mds *SelectedVertexMorphDeltas,
	model *pmx.PmxModel, selectedVertexIndexes, nextSelectedVertexIndexes []int,
) ([]int, [][]float32) {
	indices := make([]int, 0)
	vertices := make([][]float32, 0)
	for i := range len(model.Vertices.Data) {
		// 選択頂点
		var selectedVertexDelta []float32
		// 選択頂点
		if selectedVertexIndexes != nil && nextSelectedVertexIndexes != nil {
			if slices.Contains(selectedVertexIndexes, i) {
				// 選択されている頂点のUVXを＋にして（フラグをたてて）非表示にする
				selectedVertexDelta = []float32{
					0, 0, 0,
					1, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0,
				}
			}
			if slices.Contains(nextSelectedVertexIndexes, i) {
				// 選択されている頂点のUVXを0にして（フラグを落として）表示する
				selectedVertexDelta = []float32{
					0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0,
				}
			}
		} else if selectedVertexIndexes != nil && slices.Contains(selectedVertexIndexes, i) {
			// 選択されている頂点のUVXを0にして（フラグを落として）表示する
			selectedVertexDelta = []float32{
				0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
				0, 0, 0,
			}
		}

		if selectedVertexDelta != nil {
			vertices = append(vertices, selectedVertexDelta)
			indices = append(indices, i)
		}
	}
	return indices, vertices
}

func newMeshDelta(md *delta.MaterialMorphDelta) *MeshDelta {
	material := &MeshDelta{
		Diffuse:  diffuse(md),
		Specular: specular(md),
		Ambient:  ambient(md),
		Edge:     edge(md),
		EdgeSize: edgeSize(md),
	}

	return material
}

func diffuse(md *delta.MaterialMorphDelta) mgl32.Vec4 {
	d1 := md.Diffuse.GetXYZ().Copy()
	d2 := d1.MulScalar(float64(mgl.LIGHT_AMBIENT)).Add(md.Ambient)
	dm := md.MulMaterial.Diffuse.Muled(md.MulRatios.Diffuse)
	da := md.AddMaterial.Diffuse.Muled(md.AddRatios.Diffuse)

	return mgl32.Vec4{
		float32(d2.X*dm.X + da.X),
		float32(d2.Y*dm.Y + da.Y),
		float32(d2.Z*dm.Z + da.Z),
		float32(md.Diffuse.W*dm.W + da.W),
	}
}

func specular(md *delta.MaterialMorphDelta) mgl32.Vec4 {
	s1 := md.Specular.GetXYZ().MuledScalar(float64(mgl.LIGHT_AMBIENT))
	sm := md.MulMaterial.Specular.Muled(md.MulRatios.Specular)
	sa := md.AddMaterial.Specular.Muled(md.AddRatios.Specular)

	return mgl32.Vec4{
		float32(s1.X*sm.X + sa.X),
		float32(s1.Y*sm.Y + sa.Y),
		float32(s1.Z*sm.Z + sa.Z),
		float32(md.Specular.W*sm.W + sa.W),
	}
}

func ambient(md *delta.MaterialMorphDelta) mgl32.Vec3 {
	a := md.Diffuse.GetXYZ().MuledScalar(float64(mgl.LIGHT_AMBIENT))
	am := md.MulMaterial.Ambient.Muled(md.MulRatios.Ambient)
	aa := md.AddMaterial.Ambient.Muled(md.AddRatios.Ambient)
	return mgl32.Vec3{
		float32(a.X*am.X + aa.X),
		float32(a.Y*am.Y + aa.Y),
		float32(a.Z*am.Z + aa.Z),
	}
}

func edge(md *delta.MaterialMorphDelta) mgl32.Vec4 {
	e := md.Edge.GetXYZ().MuledScalar(float64(md.Diffuse.W))
	em := md.MulMaterial.Edge.Muled(md.MulRatios.Edge)
	ea := md.AddMaterial.Edge.Muled(md.AddRatios.Edge)

	return mgl32.Vec4{
		float32(e.X*em.X + ea.X),
		float32(e.Y*em.Y + ea.Y),
		float32(e.Z*em.Z + ea.Z),
		float32(md.Edge.W) * float32(md.Diffuse.W*em.W+ea.W),
	}
}

func edgeSize(md *delta.MaterialMorphDelta) float32 {
	return float32(md.Material.EdgeSize*(md.MulMaterial.EdgeSize*md.MulRatios.EdgeSize) +
		(md.AddMaterial.EdgeSize * md.AddRatios.EdgeSize))
}
