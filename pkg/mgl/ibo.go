package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
)

type IBO struct {
	id     uint32
	target uint32
	Dtype  uint32
	dsize  int
}

func NewIBO(data []byte, isSub bool) (*IBO, error) {
	var iboId uint32

	gl.GenBuffers(1, &iboId)

	ibo := &IBO{target: gl.ELEMENT_ARRAY_BUFFER, id: iboId}
	ibo.Dtype = gl.UNSIGNED_BYTE
	ibo.dsize = int(unsafe.Sizeof(byte(0)))

	ibo.Bind(isSub)
	ibo.SetIndices(data)
	ibo.Unbind()

	return ibo, nil
}

func (ibo *IBO) Bind(isSub bool) {
	gl.BindBuffer(ibo.target, ibo.id)
}

func (ibo *IBO) Unbind() {
	gl.BindBuffer(ibo.target, 0)
}

func (ibo *IBO) SetIndices(data []byte) {
	gl.BufferData(ibo.target, len(data), gl.Ptr(data), gl.STATIC_DRAW)
}

func (ibo *IBO) Delete() {
	if ibo.id != 0 {
		gl.DeleteBuffers(1, &ibo.id)
	}
}
