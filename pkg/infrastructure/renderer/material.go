//go:build windows
// +build windows

package renderer

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

type materialGL struct {
	*pmx.Material
	Texture           *TextureGL // 通常テクスチャ
	SphereTexture     *TextureGL // スフィアテクスチャ
	ToonTexture       *TextureGL // トゥーンテクスチャ
	PrevVerticesCount int        // 前の材質までの頂点数
}

func newMaterialGl(
	m *pmx.Material,
	modelPath string,
	texture *pmx.Texture,
	toonTexture *pmx.Texture,
	sphereTexture *pmx.Texture,
	windowIndex int,
	prevVerticesCount int,
) *materialGL {
	var textureGL *TextureGL
	if texture != nil {
		textureGL = newTextureGl(texture, modelPath, pmx.TEXTURE_TYPE_TEXTURE, windowIndex)
	}

	var sphereTextureGL *TextureGL
	if sphereTexture != nil {
		sphereTextureGL = newTextureGl(sphereTexture, modelPath, pmx.TEXTURE_TYPE_SPHERE, windowIndex)
	}

	var tooTextureGL *TextureGL
	if toonTexture != nil {
		tooTextureGL = newTextureGl(toonTexture, modelPath, pmx.TEXTURE_TYPE_TOON, windowIndex)
	}

	return &materialGL{
		Material:          m,
		Texture:           textureGL,
		SphereTexture:     sphereTextureGL,
		ToonTexture:       tooTextureGL,
		PrevVerticesCount: prevVerticesCount * 4,
	}
}
