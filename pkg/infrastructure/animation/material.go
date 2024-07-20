//go:build windows
// +build windows

package animation

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
