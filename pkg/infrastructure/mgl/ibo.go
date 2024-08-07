//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
)

type IBO struct {
	id      uint32         // ID
	target  uint32         // gl.ELEMENT_ARRAY_BUFFER
	size    int            // 一面のbyte数
	facePtr unsafe.Pointer // facePtr
}

func NewIBO(facePtr unsafe.Pointer, count int) *IBO {
	var iboId uint32

	gl.GenBuffers(1, &iboId)

	ibo := &IBO{
		id:      iboId,
		target:  gl.ELEMENT_ARRAY_BUFFER,
		facePtr: facePtr,
		size:    count * 4, // ひとつの面につき、dtype(UNSIGNED_INT)
	}

	return ibo
}

func (ibo *IBO) Bind() {
	gl.BindBuffer(ibo.target, ibo.id)
	gl.BufferData(ibo.target, ibo.size, ibo.facePtr, gl.STATIC_DRAW)
}

func (ibo *IBO) Unbind() {
	gl.BindBuffer(ibo.target, 0)
}

func (ibo *IBO) Delete() {
	gl.DeleteBuffers(1, &ibo.id)
}
