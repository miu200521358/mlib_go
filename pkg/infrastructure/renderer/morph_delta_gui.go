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

func VertexMorphDeltasGL(mds *delta.VertexMorphDeltas) ([]int, [][]float32) {
	vertices := make([][]float32, 0)
	indices := make([]int, 0)
	for i, md := range mds.Data {
		vertices = append(vertices, VertexMorphDeltaGL(md))
		indices = append(indices, i)
	}
	return indices, vertices
}

func VertexMorphDeltaGL(md *delta.VertexMorphDelta) []float32 {
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
		u0x = float32(md.Uv.GetX())
		u0y = float32(md.Uv.GetY())
	}
	if md.Uv1 != nil {
		u1x = float32(md.Uv1.GetX())
		u1y = float32(md.Uv1.GetY())
	}
	return []float32{
		p0, p1, p2,
		u0x, u0y, 0, 0,
		u1x, u1y, 0, 0,
		ap0, ap1, ap2,
	}
}

func SelectedVertexMorphDeltasGL(
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

func MaterialMorphDeltaResult(md *delta.MaterialMorphDelta) *MeshDelta {
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
		float32(d2.GetX()*dm.GetX() + da.GetX()),
		float32(d2.GetY()*dm.GetY() + da.GetY()),
		float32(d2.GetZ()*dm.GetZ() + da.GetZ()),
		float32(md.Diffuse.GetW()*dm.GetW() + da.GetW()),
	}
}

func specular(md *delta.MaterialMorphDelta) mgl32.Vec4 {
	s1 := md.Specular.GetXYZ().MuledScalar(float64(mgl.LIGHT_AMBIENT))
	sm := md.MulMaterial.Specular.Muled(md.MulRatios.Specular)
	sa := md.AddMaterial.Specular.Muled(md.AddRatios.Specular)

	return mgl32.Vec4{
		float32(s1.GetX()*sm.GetX() + sa.GetX()),
		float32(s1.GetY()*sm.GetY() + sa.GetY()),
		float32(s1.GetZ()*sm.GetZ() + sa.GetZ()),
		float32(md.Specular.GetW()*sm.GetW() + sa.GetW()),
	}
}

func ambient(md *delta.MaterialMorphDelta) mgl32.Vec3 {
	a := md.Diffuse.GetXYZ().MuledScalar(float64(mgl.LIGHT_AMBIENT))
	am := md.MulMaterial.Ambient.Muled(md.MulRatios.Ambient)
	aa := md.AddMaterial.Ambient.Muled(md.AddRatios.Ambient)
	return mgl32.Vec3{
		float32(a.GetX()*am.GetX() + aa.GetX()),
		float32(a.GetY()*am.GetY() + aa.GetY()),
		float32(a.GetZ()*am.GetZ() + aa.GetZ()),
	}
}

func edge(md *delta.MaterialMorphDelta) mgl32.Vec4 {
	e := md.Edge.GetXYZ().MuledScalar(float64(md.Diffuse.GetW()))
	em := md.MulMaterial.Edge.Muled(md.MulRatios.Edge)
	ea := md.AddMaterial.Edge.Muled(md.AddRatios.Edge)

	return mgl32.Vec4{
		float32(e.GetX()*em.GetX() + ea.GetX()),
		float32(e.GetY()*em.GetY() + ea.GetY()),
		float32(e.GetZ()*em.GetZ() + ea.GetZ()),
		float32(md.Edge.GetW()) * float32(md.Diffuse.GetW()*em.GetW()+ea.GetW()),
	}
}

func edgeSize(md *delta.MaterialMorphDelta) float32 {
	return float32(md.Material.EdgeSize*(md.MulMaterial.EdgeSize*md.MulRatios.EdgeSize) +
		(md.AddMaterial.EdgeSize * md.AddRatios.EdgeSize))
}
