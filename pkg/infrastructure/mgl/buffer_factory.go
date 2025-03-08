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
func (f *BufferFactory) CreateVertexBuffer(ptr unsafe.Pointer, count int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddStandardVertexAttributes().
		SetData(ptr, count).
		Build()
}

// CreateBoneBuffer はボーン表示用バッファを作成
func (f *BufferFactory) CreateBoneBuffer(ptr unsafe.Pointer, count int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddBoneAttributes().
		SetData(ptr, count).
		Build()
}

// CreateDebugBuffer はデバッグ表示用バッファを作成
func (f *BufferFactory) CreateDebugBuffer() *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddPositionColorAttributes().
		Build()
}

// UpdateDebugBuffer はデバッグバッファのデータを更新
func (f *BufferFactory) UpdateDebugBuffer(handle *VertexBufferHandle, vertices []float32) {
	handle.VAO.Bind()
	handle.VBO.Bind()

	size := len(vertices) * 4 // float32のサイズ
	handle.VBO.BufferData(size, gl.Ptr(&vertices[0]), rendering.BufferUsageStatic)

	handle.VBO.Unbind()
	handle.VAO.Unbind()
}

// CreateOverrideBuffer はオーバーライド用バッファを作成
func (f *BufferFactory) CreateOverrideBuffer(ptr unsafe.Pointer, count int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddOverrideAttributes().
		SetData(ptr, count).
		Build()
}

// CreateElementBuffer はインデックスバッファを作成
func (f *BufferFactory) CreateElementBuffer(ptr unsafe.Pointer, count int) *ElementBuffer {
	ebo := newElementBuffer()

	ebo.Bind()

	size := count * 4 // uint32のサイズ
	ebo.BufferData(size, ptr, rendering.BufferUsageStatic)

	ebo.Unbind()

	return ebo
}

// CreateFloorBuffer は床表示用バッファを作成
func (f *BufferFactory) CreateFloorBuffer(floorVertices []float32, size int) *VertexBufferHandle {
	builder := NewVertexBufferBuilder()
	return builder.
		AddPositionColorAttributes().
		SetData(unsafe.Pointer(&floorVertices[0]), size).
		Build()
}
