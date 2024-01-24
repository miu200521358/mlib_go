package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
)

type IBO struct {
	id        uint32         // ID
	target    uint32         // gl.ELEMENT_ARRAY_BUFFER
	Dtype     uint32         // 面の型
	FaceSize  int            // 一面のbyte数
	FaceCount int32          // 面の数
	facePtr   unsafe.Pointer // facePtr
}

func NewIBO(facePtr unsafe.Pointer, faceCount int, faceDtype uint32) *IBO {
	var iboId uint32

	gl.GenBuffers(1, &iboId)

	ibo := &IBO{
		id:        iboId,
		target:    gl.ELEMENT_ARRAY_BUFFER,
		FaceCount: int32(faceCount),
		facePtr:   facePtr,
	}

	if faceDtype == uint32(8) {
		ibo.Dtype = gl.UNSIGNED_BYTE
	} else if faceDtype == uint32(16) {
		ibo.Dtype = gl.UNSIGNED_SHORT
	} else if faceDtype == uint32(32) {
		ibo.Dtype = gl.UNSIGNED_INT
	}

	ibo.FaceSize = faceCount * int(ibo.Dtype)

	ibo.Bind()
	ibo.Unbind()

	return ibo
}

func (ibo *IBO) Bind() {
	gl.BindBuffer(ibo.target, ibo.id)
	gl.BufferData(ibo.target, ibo.FaceSize, ibo.facePtr, gl.STATIC_DRAW)
}

func (ibo *IBO) Unbind() {
	gl.BindBuffer(ibo.target, 0)
}

func (ibo *IBO) Delete() {
	if ibo.id != 0 {
		gl.DeleteBuffers(1, &ibo.id)
	}
}
