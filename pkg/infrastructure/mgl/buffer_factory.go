//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// BufferFactory はバッファ作成のためのファクトリー
type BufferFactory struct{}

// NewBufferFactory は新しいバッファファクトリーを作成
func NewBufferFactory() *BufferFactory {
	return &BufferFactory{}
}

// CreateVertexBuffer は標準的なモデル用頂点バッファを作成
func (f *BufferFactory) CreateVertexBuffer(ptr unsafe.Pointer, arrayCount int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddStandardVertexAttributes().
		SetData(ptr, arrayCount).
		Build()
}

// CreateBoneBuffer はボーン表示用バッファを作成
func (f *BufferFactory) CreateBoneBuffer(ptr unsafe.Pointer, arrayCount int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddBoneAttributes().
		SetData(ptr, arrayCount).
		Build()
}

// CreateDebugBuffer はデバッグ表示用バッファを作成
func (f *BufferFactory) CreateDebugBuffer(ptr unsafe.Pointer, arrayCount int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddPositionColorAttributes().
		SetData(ptr, arrayCount).
		Build()
}

// UpdateDebugBuffer はデバッグバッファのデータを更新
func (f *BufferFactory) UpdateDebugBuffer(handle *VertexBufferHandle, vertices []float32, fieldCount int) {
	handle.VAO.Bind()
	handle.VBO.Bind()

	handle.VBO.BufferData(len(vertices), fieldCount, handle.FloatSize, gl.Ptr(&vertices[0]), rendering.BufferUsageStatic)

	handle.VBO.Unbind()
	handle.VAO.Unbind()
}

// CreateOverrideBuffer はオーバーライド用バッファを作成
func (f *BufferFactory) CreateOverrideBuffer(overrideVertices []float32) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddOverrideAttributes().
		SetData(unsafe.Pointer(&overrideVertices[0]), len(overrideVertices)).
		Build()
}

// CreateElementBuffer はインデックスバッファを作成
func (f *BufferFactory) CreateElementBuffer(ptr unsafe.Pointer, arrayCount int) *ElementBuffer {
	ebo := newElementBuffer()

	ebo.Bind()

	size := arrayCount * 4 // uint32のサイズ
	ebo.BufferData(size, ptr, rendering.BufferUsageStatic)

	ebo.Unbind()

	return ebo
}

// CreateFloorBuffer は床表示用バッファを作成
func (f *BufferFactory) CreateFloorBuffer(floorVertices []float32) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddPositionColorAttributes().
		SetData(unsafe.Pointer(&floorVertices[0]), len(floorVertices)).
		Build()
}
