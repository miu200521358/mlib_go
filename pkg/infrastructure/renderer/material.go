//go:build windows
// +build windows

package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

type materialGL struct {
	*pmx.Material
	texture           *textureGl // 通常テクスチャ
	sphereTexture     *textureGl // スフィアテクスチャ
	toonTexture       *textureGl // トゥーンテクスチャ
	prevVerticesCount int        // 前の材質までの頂点数
}

type MeshDelta struct {
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
	Ambient  mgl32.Vec3
	Edge     mgl32.Vec4
	EdgeSize float32
}
