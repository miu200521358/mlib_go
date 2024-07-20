//go:build windows
// +build windows

package delta

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type RenderDeltas struct {
	InvisibleMaterialIndexes []int        // 非表示材質インデックス
	SelectedVertexIndexes    []int        // 選択頂点インデックス
	VertexMorphDeltaIndexes  []int        // 頂点モーフインデックス
	VertexMorphDeltas        [][]float32  // 頂点モーフデルタ
	MeshDeltas               []*MeshDelta // メッシュデルタ
}

func NewRenderDeltas() *RenderDeltas {
	return &RenderDeltas{
		InvisibleMaterialIndexes: make([]int, 0),
		SelectedVertexIndexes:    make([]int, 0),
		VertexMorphDeltaIndexes:  make([]int, 0),
		VertexMorphDeltas:        make([][]float32, 0),
		MeshDeltas:               make([]*MeshDelta, 0),
	}
}

type MeshDelta struct {
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
	Ambient  mgl32.Vec3
	Edge     mgl32.Vec4
	EdgeSize float32
}

func NewMeshDelta(md *MaterialMorphDelta) *MeshDelta {
	material := &MeshDelta{
		Diffuse:  diffuse(md),
		Specular: specular(md),
		Ambient:  ambient(md),
		Edge:     edge(md),
		EdgeSize: edgeSize(md),
	}

	return material
}

func diffuse(md *MaterialMorphDelta) mgl32.Vec4 {
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

func specular(md *MaterialMorphDelta) mgl32.Vec4 {
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

func ambient(md *MaterialMorphDelta) mgl32.Vec3 {
	a := md.Diffuse.GetXYZ().MuledScalar(float64(mgl.LIGHT_AMBIENT))
	am := md.MulMaterial.Ambient.Muled(md.MulRatios.Ambient)
	aa := md.AddMaterial.Ambient.Muled(md.AddRatios.Ambient)
	return mgl32.Vec3{
		float32(a.X*am.X + aa.X),
		float32(a.Y*am.Y + aa.Y),
		float32(a.Z*am.Z + aa.Z),
	}
}

func edge(md *MaterialMorphDelta) mgl32.Vec4 {
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

func edgeSize(md *MaterialMorphDelta) float32 {
	return float32(md.Material.EdgeSize*(md.MulMaterial.EdgeSize*md.MulRatios.EdgeSize) +
		(md.AddMaterial.EdgeSize * md.AddRatios.EdgeSize))
}
