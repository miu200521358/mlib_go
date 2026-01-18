//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
)

// usageMap はBufferUsageをOpenGL定数へ変換する。
var usageMap = map[graphics_api.BufferUsage]uint32{
	graphics_api.BufferUsageStatic:  gl.STATIC_DRAW,
	graphics_api.BufferUsageDynamic: gl.DYNAMIC_DRAW,
	graphics_api.BufferUsageStream:  gl.STREAM_DRAW,
}

// VertexBuffer はOpenGLの頂点バッファ実装。
type VertexBuffer struct {
	id         uint32
	target     uint32
	itemCount  int
	fieldCount int
	size       int
	ptr        unsafe.Pointer
}

// newVertexBuffer は頂点バッファを生成する。
func newVertexBuffer() *VertexBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	return &VertexBuffer{
		id:     id,
		target: gl.ARRAY_BUFFER,
	}
}

// Bind はバッファをバインドする。
func (vb *VertexBuffer) Bind() {
	gl.BindBuffer(vb.target, vb.id)
}

// Unbind はバッファをアンバインドする。
func (vb *VertexBuffer) Unbind() {
	gl.BindBuffer(vb.target, 0)
}

// BufferData は頂点データを転送する。
func (vb *VertexBuffer) BufferData(size int, data unsafe.Pointer, usage graphics_api.BufferUsage) {
	vb.size = size
	vb.ptr = data
	gl.BufferData(vb.target, vb.size, data, usageMap[usage])
}

// BufferDataWithLayout は頂点レイアウトを含めて転送する。
func (vb *VertexBuffer) BufferDataWithLayout(arrayCount, fieldCount, floatSize int, data unsafe.Pointer, usage graphics_api.BufferUsage) {
	vb.itemCount = arrayCount / fieldCount
	vb.fieldCount = fieldCount
	vb.size = arrayCount * floatSize
	vb.ptr = data
	gl.BufferData(vb.target, vb.size, data, usageMap[usage])
}

// BufferSubData は部分更新を行う。
func (vb *VertexBuffer) BufferSubData(offset int, size int, data unsafe.Pointer) {
	gl.BufferSubData(vb.target, offset, size, data)
}

// Delete はバッファを削除する。
func (vb *VertexBuffer) Delete() {
	if vb == nil || vb.id == 0 {
		return
	}

	gl.DeleteBuffers(1, &vb.id)
}

// GetID はバッファIDを返す。
func (vb *VertexBuffer) GetID() uint32 {
	return vb.id
}

// SetAttribute は頂点属性を設定する。
func (vb *VertexBuffer) SetAttribute(attr graphics_api.VertexAttribute) {
	gl.EnableVertexAttribArray(attr.Index)
	gl.VertexAttribPointerWithOffset(
		attr.Index,
		int32(attr.Size),
		attr.Type,
		attr.Normalize,
		attr.Stride,
		uintptr(attr.Offset),
	)
}

// IndexBuffer はOpenGLのインデックスバッファ実装。
type IndexBuffer struct {
	id     uint32
	target uint32
}

// newIndexBuffer はインデックスバッファを生成する。
func newIndexBuffer() *IndexBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	return &IndexBuffer{
		id:     id,
		target: gl.ELEMENT_ARRAY_BUFFER,
	}
}

// Bind はバッファをバインドする。
func (eb *IndexBuffer) Bind() {
	gl.BindBuffer(eb.target, eb.id)
}

// Unbind はバッファをアンバインドする。
func (eb *IndexBuffer) Unbind() {
	gl.BindBuffer(eb.target, 0)
}

// BufferData はインデックスデータを転送する。
func (eb *IndexBuffer) BufferData(size int, data unsafe.Pointer, usage graphics_api.BufferUsage) {
	gl.BufferData(eb.target, size, data, usageMap[usage])
}

// Delete はバッファを削除する。
func (eb *IndexBuffer) Delete() {
	if eb == nil || eb.id == 0 {
		return
	}
	gl.DeleteBuffers(1, &eb.id)
}

// GetID はバッファIDを返す。
func (eb *IndexBuffer) GetID() uint32 {
	return eb.id
}

// VertexArray はOpenGLの頂点配列実装。
type VertexArray struct {
	id uint32
}

// newVertexArray は頂点配列を生成する。
func newVertexArray() *VertexArray {
	var id uint32
	gl.GenVertexArrays(1, &id)
	return &VertexArray{id: id}
}

// Bind は配列をバインドする。
func (va *VertexArray) Bind() {
	gl.BindVertexArray(va.id)
}

// Unbind は配列をアンバインドする。
func (va *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

// Delete は配列を削除する。
func (va *VertexArray) Delete() {
	if va == nil || va.id == 0 {
		return
	}

	gl.DeleteVertexArrays(1, &va.id)
}

// GetID は配列IDを返す。
func (va *VertexArray) GetID() uint32 {
	return va.id
}
