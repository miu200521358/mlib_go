//go:build windows
// +build windows

package delta

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	LIGHT_DIFFUSE  float32 = 0.0
	LIGHT_SPECULAR float32 = 154.0 / 255.0
	LIGHT_AMBIENT  float32 = 154.0 / 255.0
)

type MeshDelta struct {
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
	Ambient  mgl32.Vec3
	Emissive mgl32.Vec3
	Edge     mgl32.Vec4
	EdgeSize float32
}

func NewMeshDelta(materialMorphDelta *MaterialMorphDelta) *MeshDelta {
	material := &MeshDelta{
		Diffuse:  diffuse(materialMorphDelta),
		Specular: specular(materialMorphDelta),
		Ambient:  ambient(materialMorphDelta),
		Emissive: emissive(materialMorphDelta),
		Edge:     edge(materialMorphDelta),
		EdgeSize: edgeSize(materialMorphDelta),
	}

	return material
}

func diffuse(materialMorphDelta *MaterialMorphDelta) mgl32.Vec4 {
	return mgl32.Vec4{
		float32(materialMorphDelta.Diffuse.X),
		float32(materialMorphDelta.Diffuse.Y),
		float32(materialMorphDelta.Diffuse.Z),
		float32(materialMorphDelta.Diffuse.W*materialMorphDelta.MulMaterial.Diffuse.W +
			materialMorphDelta.AddMaterial.Diffuse.W),
	}
}

func specular(materialMorphDelta *MaterialMorphDelta) mgl32.Vec4 {
	s := materialMorphDelta.Specular.XYZ()
	sm := materialMorphDelta.MulMaterial.Specular
	sa := materialMorphDelta.AddMaterial.Specular

	return mgl32.Vec4{
		float32(s.X*sm.X + sa.X),
		float32(s.Y*sm.Y + sa.Y),
		float32(s.Z*sm.Z + sa.Z),
		float32(materialMorphDelta.Specular.W*sm.W + sa.W),
	}
}

func ambient(materialMorphDelta *MaterialMorphDelta) mgl32.Vec3 {
	a := materialMorphDelta.Diffuse
	am := materialMorphDelta.MulMaterial.Diffuse
	aa := materialMorphDelta.AddMaterial.Diffuse

	return mgl32.Vec3{
		float32(a.X*am.X + aa.X),
		float32(a.Y*am.Y + aa.Y),
		float32(a.Z*am.Z + aa.Z),
	}
}

func emissive(materialMorphDelta *MaterialMorphDelta) mgl32.Vec3 {
	d := materialMorphDelta.Ambient
	dm := materialMorphDelta.MulMaterial.Ambient
	da := materialMorphDelta.AddMaterial.Ambient

	return mgl32.Vec3{
		float32(d.X*dm.X + da.X),
		float32(d.Y*dm.Y + da.Y),
		float32(d.Z*dm.Z + da.Z),
	}
}

func edge(materialMorphDelta *MaterialMorphDelta) mgl32.Vec4 {
	e := materialMorphDelta.Edge.XYZ()
	em := materialMorphDelta.MulMaterial.Edge
	ea := materialMorphDelta.AddMaterial.Edge

	return mgl32.Vec4{
		float32(e.X*em.X + ea.X),
		float32(e.Y*em.Y + ea.Y),
		float32(e.Z*em.Z + ea.Z),
		float32(materialMorphDelta.Edge.W*em.W + ea.W),
	}
}

func edgeSize(materialMorphDelta *MaterialMorphDelta) float32 {
	return float32(materialMorphDelta.Material.EdgeSize*
		materialMorphDelta.MulMaterial.EdgeSize + materialMorphDelta.AddMaterial.EdgeSize)
}
