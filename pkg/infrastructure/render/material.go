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
func newMaterialGL(m *pmx.Material, prevVerticesCount int) *materialGL {
	return &materialGL{
		Material:          m,
		prevVerticesCount: prevVerticesCount,
		// 後からテクスチャは texture manager などで初期化されるため、初期状態は nil
		texture:       nil,
		sphereTexture: nil,
		toonTexture:   nil,
	}
}
