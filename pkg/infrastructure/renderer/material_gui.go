//go:build windows
// +build windows

package renderer

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

type MaterialGL struct {
	*pmx.Material
	Texture           *TextureGL // 通常テクスチャ
	SphereTexture     *TextureGL // スフィアテクスチャ
	ToonTexture       *TextureGL // トゥーンテクスチャ
	PrevVerticesCount int        // 前の材質までの頂点数
}

func materialGL(
	m *pmx.Material,
	modelPath string,
	texture *pmx.Texture,
	toonTexture *pmx.Texture,
	sphereTexture *pmx.Texture,
	windowIndex int,
	prevVerticesCount int,
) *MaterialGL {
	var textureGL *TextureGL
	if texture != nil {
		textureGL = textureGLInit(texture, modelPath, pmx.TEXTURE_TYPE_TEXTURE, windowIndex)
	}

	var sphereTextureGL *TextureGL
	if sphereTexture != nil {
		sphereTextureGL = textureGLInit(sphereTexture, modelPath, pmx.TEXTURE_TYPE_SPHERE, windowIndex)
	}

	var tooTextureGL *TextureGL
	if toonTexture != nil {
		tooTextureGL = textureGLInit(toonTexture, modelPath, pmx.TEXTURE_TYPE_TOON, windowIndex)
	}

	return &MaterialGL{
		Material:          m,
		Texture:           textureGL,
		SphereTexture:     sphereTextureGL,
		ToonTexture:       tooTextureGL,
		PrevVerticesCount: prevVerticesCount * 4,
	}
}
