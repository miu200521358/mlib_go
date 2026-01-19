//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"unsafe"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
)

// BufferFactory はバッファ作成のためのファクトリー。
type BufferFactory struct{}

// NewBufferFactory はバッファファクトリーを生成する。
func NewBufferFactory() *BufferFactory {
	return &BufferFactory{}
}

// NewVertexBuffer は標準的なモデル用頂点バッファを生成する。
func (f *BufferFactory) NewVertexBuffer(ptr unsafe.Pointer, arrayCount int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddStandardVertexAttributes().
		SetData(ptr, arrayCount).
		Build()
}

// NewBoneBuffer はボーン表示用バッファを生成する。
func (f *BufferFactory) NewBoneBuffer(ptr unsafe.Pointer, arrayCount int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddBoneAttributes().
		SetData(ptr, arrayCount).
		Build()
}

// NewDebugBuffer はデバッグ表示用バッファを生成する。
func (f *BufferFactory) NewDebugBuffer(ptr unsafe.Pointer, arrayCount int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddPositionColorAttributes().
		SetData(ptr, arrayCount).
		Build()
}

// NewOverrideBuffer はオーバーライド用バッファを生成する。
func (f *BufferFactory) NewOverrideBuffer(overrideVertices []float32) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	var ptr unsafe.Pointer
	if len(overrideVertices) > 0 {
		ptr = unsafe.Pointer(&overrideVertices[0])
	}
	return builder.
		AddOverrideAttributes().
		SetData(ptr, len(overrideVertices)).
		Build()
}

// NewIndexBuffer はインデックスバッファを生成する。
func (f *BufferFactory) NewIndexBuffer(ptr unsafe.Pointer, arrayCount int) *IndexBuffer {
	ebo := newIndexBuffer()

	ebo.Bind()

	size := arrayCount * 4 // uint32のサイズ
	ebo.BufferData(size, ptr, graphics_api.BufferUsageStatic)

	ebo.Unbind()

	return ebo
}

// NewFloorBuffer は床表示用バッファを生成する。
func (f *BufferFactory) NewFloorBuffer(floorVertices []float32) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	var ptr unsafe.Pointer
	if len(floorVertices) > 0 {
		ptr = unsafe.Pointer(&floorVertices[0])
	}
	return builder.
		AddPositionColorAttributes().
		SetData(ptr, len(floorVertices)).
		Build()
}
