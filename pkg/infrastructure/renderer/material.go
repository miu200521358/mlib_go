//go:build windows
// +build windows

package renderer

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

type materialGL struct {
	*pmx.Material
	Texture           *textureGl // 通常テクスチャ
	SphereTexture     *textureGl // スフィアテクスチャ
	ToonTexture       *textureGl // トゥーンテクスチャ
	PrevVerticesCount int        // 前の材質までの頂点数
}
