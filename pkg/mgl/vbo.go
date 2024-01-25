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
	// position(3), normal(3), uv(2), extendedUV(2), edgeFactor(1), deformBoneIndex(4), deformBoneWeight(4)
	vbo.stride = int32(4 * (3 + 3 + 2 + 2 + 1 + 4 + 4))
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

	// 2: uv
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(
		2,        // 属性のインデックス
		2,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		6*4,      // オフセット（構造体内のオフセット）
	)

	// 3: extendedUV
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(
		3,        // 属性のインデックス
		2,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		8*4,      // オフセット（構造体内のオフセット）
	)

	// 4: edgeFactor
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointerWithOffset(
		4,        // 属性のインデックス
		1,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		10*4,     // オフセット（構造体内のオフセット）
	)

	// 5: deformBoneIndex
	gl.EnableVertexAttribArray(5)
	gl.VertexAttribPointerWithOffset(
		5,        // 属性のインデックス
		4,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		11*4,     // オフセット（構造体内のオフセット）
	)

	// 6: deformBoneWeight
	gl.EnableVertexAttribArray(6)
	gl.VertexAttribPointerWithOffset(
		6,        // 属性のインデックス
		4,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		15*4,     // オフセット（構造体内のオフセット）
	)

}

// Unbinds.
func (v *VBO) Unbind() {
	gl.BindBuffer(v.target, 0)
}
