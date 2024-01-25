package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mutils"

)

// Vertex Buffer Object.
type VBO struct {
	id          uint32         // ID
	target      uint32         // gl.ARRAY_BUFFER
	size        int            // size in bytes
	verticesPtr unsafe.Pointer // verticesPtr
	stride      int32          // stride
}

// Creates a new VBO with given faceDtype.
func NewVBO(verticesPtr unsafe.Pointer, count int) *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:          vboId,
		target:      gl.ARRAY_BUFFER,
		verticesPtr: verticesPtr,
	}
	// 頂点構造体のサイズ
	// position(3), normal(3)
	vbo.stride = int32(4 * (3 + 3))
	vbo.size = count * 4

	return vbo
}

// Delete this VBO.
func (v *VBO) Delete() {
	gl.DeleteBuffers(1, &v.id)
}

// Binds VBO for rendering.
func (v *VBO) Bind() {
	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, v.verticesPtr, gl.STATIC_DRAW)
	mutils.CheckGLError()

	// 0: position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(
		0,        // 属性のインデックス
		3,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		0,        // オフセット
	)
	mutils.CheckGLError()

	// 1: normal
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(
		1,        // 属性のインデックス
		3,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		3*4,      // オフセット（構造体内のオフセット）
	)
	mutils.CheckGLError()
}

// Unbinds.
func (v *VBO) Unbind() {
	gl.BindBuffer(v.target, 0)
}
