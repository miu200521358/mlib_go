//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// 使用方法のマッピング
var usageMap = map[rendering.BufferUsage]uint32{
	rendering.BufferUsageStatic:  gl.STATIC_DRAW,
	rendering.BufferUsageDynamic: gl.DYNAMIC_DRAW,
	rendering.BufferUsageStream:  gl.STREAM_DRAW,
}

// VertexBuffer はOpenGLの頂点バッファ実装
type VertexBuffer struct {
	id         uint32
	target     uint32
	itemCount  int            // 項目数（頂点の個数など）
	fieldCount int            // 1項目ごとの要素数（位置(vec3)、法線(vec3)、UV(vec2)の場合、8）
	size       int            // バッファサイズ(項目数*1項目ごとの要素数*要素サイズ(Float:4))
	ptr        unsafe.Pointer // verticesPtr
}

// newVertexBuffer は新しい頂点バッファを作成
func newVertexBuffer() *VertexBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	return &VertexBuffer{
		id:     id,
		target: gl.ARRAY_BUFFER,
	}
}

func (vb *VertexBuffer) Bind() {
	gl.BindBuffer(vb.target, vb.id)
}

func (vb *VertexBuffer) Unbind() {
	gl.BindBuffer(vb.target, 0)
}

func (vb *VertexBuffer) BufferData(arrayCount, fieldCount, floatSize int, data unsafe.Pointer, usage rendering.BufferUsage) {
	vb.itemCount = arrayCount / fieldCount
	vb.fieldCount = fieldCount
	vb.size = arrayCount * floatSize
	vb.ptr = data
	gl.BufferData(vb.target, vb.size, data, usageMap[usage])
}

func (vb *VertexBuffer) BufferSubData(offset int, size int, data unsafe.Pointer) {
	gl.BufferSubData(vb.target, offset, size, data)
}

func (vb *VertexBuffer) Delete() {
	if vb == nil || vb.id == 0 {
		return
	}

	gl.DeleteBuffers(1, &vb.id)
}

func (vb *VertexBuffer) GetID() uint32 {
	return vb.id
}

func (vb *VertexBuffer) SetAttribute(attr rendering.VertexAttribute) {
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

// ElementBuffer はOpenGLのインデックスバッファ実装
type ElementBuffer struct {
	id     uint32
	target uint32
}

// newElementBuffer は新しいインデックスバッファを作成
func newElementBuffer() *ElementBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	return &ElementBuffer{
		id:     id,
		target: gl.ELEMENT_ARRAY_BUFFER,
	}
}

func (eb *ElementBuffer) Bind() {
	gl.BindBuffer(eb.target, eb.id)
}

func (eb *ElementBuffer) Unbind() {
	gl.BindBuffer(eb.target, 0)
}

func (eb *ElementBuffer) BufferData(size int, data unsafe.Pointer, usage rendering.BufferUsage) {
	gl.BufferData(eb.target, size, data, usageMap[usage])
}

func (eb *ElementBuffer) Delete() {
	if eb == nil || eb.id == 0 {
		return
	}
	gl.DeleteBuffers(1, &eb.id)
}

func (eb *ElementBuffer) GetID() uint32 {
	return eb.id
}

// VertexArray はOpenGLの頂点配列実装
type VertexArray struct {
	id uint32
}

// newVertexArray は新しい頂点配列を作成
func newVertexArray() *VertexArray {
	var id uint32
	gl.GenVertexArrays(1, &id)
	return &VertexArray{id: id}
}

func (va *VertexArray) Bind() {
	gl.BindVertexArray(va.id)
}

func (va *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

func (va *VertexArray) Delete() {
	if va == nil || va.id == 0 {
		return
	}

	gl.DeleteVertexArrays(1, &va.id)
}

func (va *VertexArray) GetID() uint32 {
	return va.id
}
