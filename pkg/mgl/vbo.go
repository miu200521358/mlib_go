package mgl

import (
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/mbt"
)

// Vertex Buffer Object.
type VBO struct {
	id         uint32         // ID
	target     uint32         // gl.ARRAY_BUFFER
	size       int            // size in bytes
	ptr        unsafe.Pointer // verticesPtr
	stride     int32          // stride
	strideSize int            // strideSize
}

// Delete this VBO.
func (v *VBO) Delete() {
	gl.DeleteBuffers(1, &v.id)
}

// Unbinds.
func (v *VBO) Unbind() {
	gl.BindBuffer(v.target, 0)
}

// Creates a new VBO with given faceDtype.
func NewVBOForVertex(ptr unsafe.Pointer, count int) *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    ptr,
	}
	// 頂点構造体のサイズ(全部floatとする)
	// position(3), normal(3), uv(2), extendedUV(2), edgeFactor(1), deformBoneIndex(4), deformBoneWeight(4),
	// isSdef(1), sdefC(3), sdefR0(3), sdefR1(3), vertexDelta(3), uvDelta(4), uv1Delta(4), afterVertexDelta(3)
	vbo.strideSize = 3 + 3 + 2 + 2 + 1 + 4 + 4 + 1 + 3 + 3 + 3 + 3 + 4 + 4 + 3
	vbo.stride = int32(4 * vbo.strideSize)
	vbo.size = count * 4

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindVertex(vertices []float32, vertexDeltas [][]float32, vertexChunkSize int) {
	gl.BindBuffer(v.target, v.id)

	if vertices != nil && vertexDeltas != nil {
		// モーフ分の変動量を設定
		numChunks := len(vertexDeltas)/vertexChunkSize + 1

		var wg sync.WaitGroup
		wg.Add(numChunks)

		for i := 0; i < numChunks; i++ {
			go func(chunkIndex int) {
				defer wg.Done()

				// Calculate the start and end indices for the current chunk
				startIndex := chunkIndex * vertexChunkSize
				endIndex := (chunkIndex + 1) * vertexChunkSize

				// Process the vertex deltas for the current chunk
				for i := startIndex; i < endIndex; i++ {
					if i >= len(vertexDeltas) {
						break
					}
					vd := vertexDeltas[i]

					// Calculate the offset for the current vertex delta
					offset := i*v.strideSize + (3 + 3 + 2 + 2 + 1 + 4 + 4 + 1 + 3 + 3 + 3)
					nextOffset := (i + 1) * v.strideSize

					// Copy the vertex delta to the vertices slice
					copy(vertices[offset:nextOffset], vd)
				}
			}(i)
		}

		// Wait for all goroutines to finish
		wg.Wait()
		gl.BufferData(v.target, v.size, gl.Ptr(vertices), gl.STATIC_DRAW)
	} else {
		gl.BufferData(v.target, v.size, v.ptr, gl.STATIC_DRAW)
	}

	CheckGLError()

	// 0: position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(
		0,        // 属性のインデックス
		3,        // 属性のサイズ
		gl.FLOAT, // データの型
		false,    // 正規化するかどうか
		v.stride, // ストライド
		0,        // オフセット（構造体内の位置）
	)

	// 1: normal
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(
		1,
		3,
		gl.FLOAT,
		false,
		v.stride,
		3*4,
	)

	// 2: uv
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(
		2,
		2,
		gl.FLOAT,
		false,
		v.stride,
		6*4,
	)

	// 3: extendedUV
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(
		3,
		2,
		gl.FLOAT,
		false,
		v.stride,
		8*4,
	)

	// 4: edgeFactor
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointerWithOffset(
		4,
		1,
		gl.FLOAT,
		false,
		v.stride,
		10*4,
	)

	// 5: deformBoneIndex
	gl.EnableVertexAttribArray(5)
	gl.VertexAttribPointerWithOffset(
		5,
		4,
		gl.FLOAT,
		false,
		v.stride,
		11*4,
	)

	// 6: deformBoneWeight
	gl.EnableVertexAttribArray(6)
	gl.VertexAttribPointerWithOffset(
		6,
		4,
		gl.FLOAT,
		false,
		v.stride,
		15*4,
	)

	// 7: isSdef
	gl.EnableVertexAttribArray(7)
	gl.VertexAttribPointerWithOffset(
		7,
		1,
		gl.FLOAT,
		false,
		v.stride,
		19*4,
	)

	// 8: SDEF-C
	gl.EnableVertexAttribArray(8)
	gl.VertexAttribPointerWithOffset(
		8,
		3,
		gl.FLOAT,
		false,
		v.stride,
		20*4,
	)

	// 9: SDEF-R0
	gl.EnableVertexAttribArray(9)
	gl.VertexAttribPointerWithOffset(
		9,
		3,
		gl.FLOAT,
		false,
		v.stride,
		23*4,
	)

	// 10: SDEF-R1
	gl.EnableVertexAttribArray(10)
	gl.VertexAttribPointerWithOffset(
		10,
		3,
		gl.FLOAT,
		false,
		v.stride,
		26*4,
	)

	// 11: vertexDelta
	gl.EnableVertexAttribArray(11)
	gl.VertexAttribPointerWithOffset(
		11,
		3,
		gl.FLOAT,
		false,
		v.stride,
		29*4,
	)

	// 12: uvDelta
	gl.EnableVertexAttribArray(12)
	gl.VertexAttribPointerWithOffset(
		12,
		4,
		gl.FLOAT,
		false,
		v.stride,
		32*4,
	)

	// 13: uv1Delta
	gl.EnableVertexAttribArray(13)
	gl.VertexAttribPointerWithOffset(
		13,
		4,
		gl.FLOAT,
		false,
		v.stride,
		36*4,
	)

	// 14: vertexDelta
	gl.EnableVertexAttribArray(14)
	gl.VertexAttribPointerWithOffset(
		14,
		3,
		gl.FLOAT,
		false,
		v.stride,
		40*4,
	)

}

// Creates a new VBO with given faceDtype.
func NewVBOForBone(ptr unsafe.Pointer, count int) *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    ptr,
	}
	// 剛体構造体のサイズ(全部floatとする)
	// position(3), color(4)
	vbo.stride = int32(4 * (3 + 4))
	vbo.size = count * 4

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindBone() {
	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, v.ptr, gl.STATIC_DRAW)
	CheckGLError()

	// 0: position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(
		0,
		3,
		gl.FLOAT,
		false,
		v.stride,
		0*4,
	)

	// 1: color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(
		1,
		4,
		gl.FLOAT,
		false,
		v.stride,
		3*4,
	)

}

// Creates a new VBO with given faceDtype.
func NewVBOForDebug() *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    nil,
	}
	// 剛体構造体のサイズ(全部floatとする)
	// position(3)
	vbo.stride = int32(4 * (3))
	vbo.size = 6 * 4

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindDebug(from mbt.BtVector3, to mbt.BtVector3) {
	vertices := []float32{
		from.GetX(), from.GetY(), from.GetZ(),
		to.GetX(), to.GetY(), to.GetZ(),
	}

	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, gl.Ptr(&vertices[0]), gl.STATIC_DRAW)
	CheckGLError()

	// 0: position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(
		0,
		3,
		gl.FLOAT,
		false,
		v.stride,
		0*4,
	)

}
