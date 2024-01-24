package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IBO struct {
	id      uint32         // ID
	target  uint32         // gl.ELEMENT_ARRAY_BUFFER
	Dtype   uint32         // 面の型
	Dsize   int            // 面のbyte数
	size    int            // 一面のbyte数
	facePtr unsafe.Pointer // facePtr
}

func NewIBO(facePtr unsafe.Pointer, count int, dtype uint32) *IBO {
	var iboId uint32

	gl.GenBuffers(1, &iboId)

	ibo := &IBO{
		id:      iboId,
		target:  gl.ELEMENT_ARRAY_BUFFER,
		facePtr: facePtr,
	}

	ibo.Dtype = gl.UNSIGNED_INT
	ibo.Dsize = 4

	ibo.size = count * ibo.Dsize * 3 // ひとつの面につき、dtype(任意byte) * 3(三角形)

	return ibo
}

func (ibo *IBO) Bind() {
	gl.BindBuffer(ibo.target, ibo.id)
	gl.BufferData(ibo.target, ibo.size, ibo.facePtr, gl.STATIC_DRAW)
	mutils.CheckGLError()
}

func (ibo *IBO) Unbind() {
	gl.BindBuffer(ibo.target, 0)
}

func (ibo *IBO) Delete() {
	if ibo.id != 0 {
		gl.DeleteBuffers(1, &ibo.id)
	}
}
