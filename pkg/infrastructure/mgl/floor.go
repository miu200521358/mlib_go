//go:build windows
// +build windows

package mgl

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
)

// FloorRenderer は床描画を担当するレンダラー
type FloorRenderer struct {
	bufferHandle *VertexBufferHandle
	vertexCount  int32
}

// NewFloorRenderer は新しい床レンダラーを作成
func NewFloorRenderer() *FloorRenderer {
	factory := NewBufferFactory()
	bufferHandle, vertexCount := factory.CreateFloorBuffer()

	return &FloorRenderer{
		bufferHandle: bufferHandle,
		vertexCount:  vertexCount,
	}
}

// Render は床を描画
func (f *FloorRenderer) Render(programID uint32) {
	gl.UseProgram(programID)

	// OpenGLエラーチェック
	if err := gl.GetError(); err != gl.NO_ERROR {
		mlog.E("OpenGL error before floor rendering: %v", err)
	}

	f.bufferHandle.Bind()

	// レンダリング状態の確認
	mlog.D("Floor rendering: vertices=%d", f.vertexCount)

	// LINESモードで描画
	gl.DrawArrays(gl.LINES, 0, f.vertexCount)

	// 描画後のエラーチェック
	if err := gl.GetError(); err != gl.NO_ERROR {
		mlog.E("OpenGL error after floor rendering: %v", err)
	}

	f.bufferHandle.Unbind()
	gl.UseProgram(0)
}

// Delete はリソースを解放
func (f *FloorRenderer) Delete() {
	if f.bufferHandle != nil {
		f.bufferHandle.Delete()
	}
}

// generateFloorVertices は床描画用の頂点データを生成
func generateFloorVertices() []float32 {
	// 結果格納用の配列
	vertices := make([]float32, 0)

	// グリッド設定
	gridSize := 50
	step := 5

	// 関数：1本の線を追加
	addLine := func(x1, y1, z1, x2, y2, z2 float32, r, g, b, a float32) {
		// 始点
		vertices = append(vertices, x1, y1, z1) // 位置
		vertices = append(vertices, r, g, b, a) // 色

		// 終点
		vertices = append(vertices, x2, y2, z2) // 位置
		vertices = append(vertices, r, g, b, a) // 色
	}

	normalColor := [4]float32{0.9, 0.9, 0.9, 0.7} // グリッド線の通常色（グレー）
	xColor := [4]float32{1.0, 0.0, 0.0, 1.0}      // X軸: 赤
	yColor := [4]float32{0.0, 1.0, 0.0, 1.0}      // Y軸: 緑
	zColor := [4]float32{0.0, 0.0, 1.0, 1.0}      // Z軸: 青

	// X座標固定のグリッド線（Z方向に伸びる線）
	for x := -gridSize; x <= gridSize; x += step {
		if x == 0 {
			// Z軸を2つに分け、負方向のみ青色に

			// Z軸負方向（青色）
			addLine(0, 0, float32(-gridSize), 0, 0, 0,
				zColor[0], zColor[1], zColor[2], zColor[3])

			// Z軸正方向（グレー）
			addLine(0, 0, 0, 0, 0, float32(gridSize),
				normalColor[0], normalColor[1], normalColor[2], normalColor[3])
		} else {
			// 通常の垂直グリッド線
			addLine(float32(x), 0, float32(-gridSize), float32(x), 0, float32(gridSize),
				normalColor[0], normalColor[1], normalColor[2], normalColor[3])
		}
	}

	// Z座標固定のグリッド線（X方向に伸びる線）
	for z := -gridSize; z <= gridSize; z += step {
		if z == 0 {
			// X軸を2つに分け、正方向のみ赤色に

			// X軸負方向（グレー）
			addLine(0, 0, 0, float32(gridSize), 0, 0,
				normalColor[0], normalColor[1], normalColor[2], normalColor[3])

			// X軸正方向（赤色）
			addLine(float32(-gridSize), 0, 0, 0, 0, 0,
				xColor[0], xColor[1], xColor[2], xColor[3])
		} else {
			// 通常の水平グリッド線
			addLine(float32(-gridSize), 0, float32(z), float32(gridSize), 0, float32(z),
				normalColor[0], normalColor[1], normalColor[2], normalColor[3])
		}
	}

	// Y軸線（垂直正方向）- 常に緑色
	addLine(0, 0, 0, 0, float32(gridSize), 0, yColor[0], yColor[1], yColor[2], yColor[3])

	return vertices
}
