//go:build windows
// +build windows

package render

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type materialGL struct {
	*pmx.Material
	texture           *textureGl // 通常テクスチャ
	sphereTexture     *textureGl // スフィアテクスチャ
	toonTexture       *textureGl // トゥーンテクスチャ
	prevVerticesCount int        // 前の材質までの頂点数
}

// newMaterialGL は、pmx.Material をもとに描画用の materialGL を生成します。
// 引数 m にはドメイン側の pmx.Material、prevVerticesCount には前の材質までの頂点数を指定します。
// 戻り値は、OpenGL描画に必要な拡張情報を保持する materialGL となります。
func newMaterialGL(m *pmx.Material, prevVerticesCount int, tm *TextureManager) *materialGL {
	var toonTexGl *textureGl
	if m.ToonSharingFlag == pmx.TOON_SHARING_INDIVIDUAL && m.ToonTextureIndex != -1 {
		// 個別Toon
		toonTexGl = tm.Texture(m.ToonTextureIndex)
	} else if m.ToonSharingFlag == pmx.TOON_SHARING_SHARING && m.ToonTextureIndex != -1 {
		// 共有Toon
		toonTexGl = tm.ToonTexture(m.ToonTextureIndex)
	}

	// スフィアは無効かテクスチャが無い場合は作らない
	var sphereTexGl *textureGl
	if m.SphereMode != pmx.SPHERE_MODE_INVALID && m.SphereTextureIndex != -1 {
		sphereTexGl = tm.Texture(m.SphereTextureIndex)
	}

	return &materialGL{
		Material:          m,
		prevVerticesCount: prevVerticesCount,
		texture:           tm.Texture(m.TextureIndex),
		sphereTexture:     sphereTexGl,
		toonTexture:       toonTexGl,
	}
}
