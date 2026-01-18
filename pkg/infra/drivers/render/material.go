//go:build windows
// +build windows

package render

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// materialGL は描画用に拡張した材質情報を保持する。
type materialGL struct {
	*model.Material
	texture           *textureGl // 通常テクスチャ
	sphereTexture     *textureGl // スフィアテクスチャ
	toonTexture       *textureGl // トゥーンテクスチャ
	prevVerticesCount int        // 前の材質までの頂点数
}

// newMaterialGL は描画用の materialGL を生成する。
func newMaterialGL(m *model.Material, prevVerticesCount int, tm *TextureManager) *materialGL {
	var toonTexGl *textureGl
	if m.ToonSharingFlag == model.TOON_SHARING_INDIVIDUAL && m.ToonTextureIndex != -1 {
		// 個別Toon
		toonTexGl = tm.Texture(m.ToonTextureIndex)
	} else if m.ToonSharingFlag == model.TOON_SHARING_SHARING && m.ToonTextureIndex != -1 {
		// 共有Toon
		toonTexGl = tm.ToonTexture(m.ToonTextureIndex)
	}

	// スフィアは無効かテクスチャが無い場合は作らない
	var sphereTexGl *textureGl
	if m.SphereMode != model.SPHERE_MODE_INVALID &&
		m.SphereTextureIndex != -1 &&
		m.SphereTextureIndex != m.TextureIndex {
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
