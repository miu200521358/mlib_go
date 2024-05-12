//go:build windows
// +build windows

package pmx

import (
	"embed"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/mview"
)

type MaterialGL struct {
	*Material
	Texture           *TextureGL // 通常テクスチャ
	SphereTexture     *TextureGL // スフィアテクスチャ
	ToonTexture       *TextureGL // トゥーンテクスチャ
	PrevVerticesCount int        // 前の材質までの頂点数
}

func (m *Material) GL(
	modelPath string,
	texture *Texture,
	toonTexture *Texture,
	sphereTexture *Texture,
	windowIndex int,
	prevVerticesCount int,
	resourceFiles embed.FS,
) *MaterialGL {
	var textureGL *TextureGL
	if texture != nil {
		textureGL = texture.GL(modelPath, TEXTURE_TYPE_TEXTURE, windowIndex, resourceFiles)
	}

	var sphereTextureGL *TextureGL
	if sphereTexture != nil {
		sphereTextureGL = sphereTexture.GL(modelPath, TEXTURE_TYPE_SPHERE, windowIndex, resourceFiles)
	}

	var tooTextureGL *TextureGL
	if toonTexture != nil {
		tooTextureGL = toonTexture.GL(modelPath, TEXTURE_TYPE_TOON, windowIndex, resourceFiles)
	}

	return &MaterialGL{
		Material:          m,
		Texture:           textureGL,
		SphereTexture:     sphereTextureGL,
		ToonTexture:       tooTextureGL,
		PrevVerticesCount: prevVerticesCount * 4,
	}
}

func (m *Material) DiffuseGL() mgl32.Vec4 {
	d1 := m.Diffuse.GetXYZ().Copy()
	d2 := d1.MulScalar(float64(mview.LIGHT_AMBIENT)).Add(&m.Ambient)
	diffuse := mgl32.Vec4{float32(d2.GetX()), float32(d2.GetY()), float32(d2.GetZ()), float32(m.Diffuse.GetW())}
	return diffuse
}

func (m *Material) AmbientGL() mgl32.Vec3 {
	a := m.Diffuse.GetXYZ().MuledScalar(float64(mview.LIGHT_AMBIENT))
	ambient := mgl32.Vec3{float32(a.GetX()), float32(a.GetY()), float32(a.GetZ())}
	return ambient
}

func (m *Material) SpecularGL() mgl32.Vec4 {
	s := m.Specular.GetXYZ().MuledScalar(float64(mview.LIGHT_AMBIENT))
	specular := mgl32.Vec4{float32(s.GetX()), float32(s.GetY()), float32(s.GetZ()), float32(m.Specular.GetW())}
	return specular
}

func (m *Material) EdgeGL() [4]float32 {
	e := m.Edge.GetXYZ().MuledScalar(float64(m.Diffuse.GetW()))
	edge := [4]float32{float32(e.GetX()), float32(e.GetY()), float32(e.GetZ()),
		float32(m.Edge.GetW()) * float32(m.Diffuse.GetW())}
	return edge
}
