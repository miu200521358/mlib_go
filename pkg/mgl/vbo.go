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
}

// Creates a new VBO with given faceDtype.
func NewVBO(verticesPtr unsafe.Pointer, count int, vertexSize int) *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:          vboId,
		target:      gl.ARRAY_BUFFER,
		size:        count * 4 * vertexSize, // ひとつの頂点につき、float(4byte) * vertexSize(1つの頂点の要素数)
		verticesPtr: verticesPtr,
	}

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
		3,        // 属性のサイズ（例: vec3 の場合は3）
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		0,        // ストライド（バイト単位）
		0,        // オフセット（ポインタまたは整数値）
	)
	mutils.CheckGLError()
}

// Unbinds.
func (v *VBO) Unbind() {
	gl.BindBuffer(v.target, 0)
}
