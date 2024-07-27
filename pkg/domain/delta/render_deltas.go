//go:build windows
// +build windows

package delta

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	LIGHT_AMBIENT float64 = 154.0 / 255.0
)

type RenderDeltas struct {
	VertexMorphDeltaIndexes []int        // 選択頂点インデックス
	VertexMorphDeltas       [][]float32  // 選択頂点デルタ
	MeshDeltas              []*MeshDelta // メッシュデルタ
}

func NewRenderDeltas() *RenderDeltas {
	return &RenderDeltas{
		VertexMorphDeltaIndexes: make([]int, 0),
		VertexMorphDeltas:       make([][]float32, 0),
		MeshDeltas:              make([]*MeshDelta, 0),
	}
}

type MeshDelta struct {
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
	Ambient  mgl32.Vec3
	Edge     mgl32.Vec4
	EdgeSize float32
}

func NewMeshDelta(materialMorphDelta *MaterialMorphDelta) *MeshDelta {
	material := &MeshDelta{
		Diffuse:  diffuse(materialMorphDelta),
		Specular: specular(materialMorphDelta),
		Ambient:  ambient(materialMorphDelta),
		Edge:     edge(materialMorphDelta),
		EdgeSize: edgeSize(materialMorphDelta),
	}

	return material
}

func diffuse(materialMorphDelta *MaterialMorphDelta) mgl32.Vec4 {
	d1 := materialMorphDelta.Diffuse.XYZ().Copy()
	d2 := d1.MulScalar(float64(LIGHT_AMBIENT)).Add(materialMorphDelta.Ambient)
	dm := materialMorphDelta.MulMaterial.Diffuse.Muled(materialMorphDelta.MulRatios.Diffuse)
	da := materialMorphDelta.AddMaterial.Diffuse.Muled(materialMorphDelta.AddRatios.Diffuse)

	return mgl32.Vec4{
		float32(d2.X*dm.X + da.X),
		float32(d2.Y*dm.Y + da.Y),
		float32(d2.Z*dm.Z + da.Z),
		float32(materialMorphDelta.Diffuse.W*dm.W + da.W),
	}
}

func specular(materialMorphDelta *MaterialMorphDelta) mgl32.Vec4 {
	s1 := materialMorphDelta.Specular.XYZ().MuledScalar(float64(LIGHT_AMBIENT))
	sm := materialMorphDelta.MulMaterial.Specular.Muled(materialMorphDelta.MulRatios.Specular)
	sa := materialMorphDelta.AddMaterial.Specular.Muled(materialMorphDelta.AddRatios.Specular)

	return mgl32.Vec4{
		float32(s1.X*sm.X + sa.X),
		float32(s1.Y*sm.Y + sa.Y),
		float32(s1.Z*sm.Z + sa.Z),
		float32(materialMorphDelta.Specular.W*sm.W + sa.W),
	}
}

func ambient(materialMorphDelta *MaterialMorphDelta) mgl32.Vec3 {
	a := materialMorphDelta.Diffuse.XYZ().MuledScalar(float64(LIGHT_AMBIENT))
	am := materialMorphDelta.MulMaterial.Ambient.Muled(materialMorphDelta.MulRatios.Ambient)
	aa := materialMorphDelta.AddMaterial.Ambient.Muled(materialMorphDelta.AddRatios.Ambient)
	return mgl32.Vec3{
		float32(a.X*am.X + aa.X),
		float32(a.Y*am.Y + aa.Y),
		float32(a.Z*am.Z + aa.Z),
	}
}

func edge(materialMorphDelta *MaterialMorphDelta) mgl32.Vec4 {
	e := materialMorphDelta.Edge.XYZ().MuledScalar(float64(materialMorphDelta.Diffuse.W))
	em := materialMorphDelta.MulMaterial.Edge.Muled(materialMorphDelta.MulRatios.Edge)
	ea := materialMorphDelta.AddMaterial.Edge.Muled(materialMorphDelta.AddRatios.Edge)

	return mgl32.Vec4{
		float32(e.X*em.X + ea.X),
		float32(e.Y*em.Y + ea.Y),
		float32(e.Z*em.Z + ea.Z),
		float32(materialMorphDelta.Edge.W) * float32(materialMorphDelta.Diffuse.W*em.W+ea.W),
	}
}

func edgeSize(materialMorphDelta *MaterialMorphDelta) float32 {
	return float32(materialMorphDelta.Material.EdgeSize*(materialMorphDelta.MulMaterial.EdgeSize*materialMorphDelta.MulRatios.EdgeSize) +
		(materialMorphDelta.AddMaterial.EdgeSize * materialMorphDelta.AddRatios.EdgeSize))
}
