//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
)

// Vertex Buffer Object.
type VBO struct {
	id         uint32         // ID
	target     uint32         // gl.ARRAY_BUFFER
	size       int            // size in bytes
	ptr        unsafe.Pointer // verticesPtr
	stride     int32          // stride
	StrideSize int            // strideSize
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
	vbo.StrideSize = 3 + 3 + 2 + 2 + 1 + 4 + 4 + 1 + 3 + 3 + 3 + 3 + 4 + 4 + 3
	vbo.stride = int32(4 * vbo.StrideSize)
	vbo.size = count * 4

	return vbo
}

// Creates a new VBO with given faceDtype.
func NewVBOForCursorPosition() *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    nil,
	}
	// 剛体構造体のサイズ(全部floatとする)
	// position(3)
	vbo.stride = int32(4 * 3)

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindCursorPosition(cursorLines []float32) {
	// cursorLinesの要素数 * float32のサイズ
	v.size = len(cursorLines) * 4

	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, gl.Ptr(&cursorLines[0]), gl.STATIC_DRAW)

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

// Creates a new VBO with given faceDtype.
func NewVBOForOverride(ptr unsafe.Pointer, count int) *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    ptr,
	}
	// 頂点構造体のサイズ(全部floatとする)
	// position(3), texCoords(2)
	vbo.StrideSize = 3 + 2
	vbo.stride = int32(4 * vbo.StrideSize)
	vbo.size = count * 4

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindOverride() {
	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, v.ptr, gl.STATIC_DRAW)

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

	// 1: texCoords
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(
		1,
		2,
		gl.FLOAT,
		false,
		v.stride,
		3*4,
	)
}

func newVertexMorphDeltaGl(md *delta.VertexMorphDelta) []float32 {
	var p0, p1, p2 float32
	if md.Position != nil {
		p := NewGlVec3(md.Position)
		p0, p1, p2 = p[0], p[1], p[2]
	}
	var ap0, ap1, ap2 float32
	if md.AfterPosition != nil {
		ap := NewGlVec3(md.AfterPosition)
		ap0, ap1, ap2 = ap[0], ap[1], ap[2]
	}
	// UVは符号関係ないのでそのまま取得する
	var u0x, u0y, u1x, u1y float32
	if md.Uv != nil {
		u0x = float32(md.Uv.X)
		u0y = float32(md.Uv.Y)
	}
	if md.Uv1 != nil {
		u1x = float32(md.Uv1.X)
		u1y = float32(md.Uv1.Y)
	}
	return []float32{
		p0, p1, p2,
		u0x, u0y, 0, 0,
		u1x, u1y, 0, 0,
		ap0, ap1, ap2,
	}
}

// Binds VBO for rendering.
func (v *VBO) BindVertex(vertices *delta.VertexMorphDeltas) {
	gl.BindBuffer(v.target, v.id)

	if vertices != nil && vertices.Length() > 0 {
		vboVertexSize := (3 + 3 + 2 + 2 + 1 + 4 + 4 + 1 + 3 + 3 + 3)

		// モーフ分の変動量を設定
		for vertexDelta := range vertices.Iterator() {
			vd := newVertexMorphDeltaGl(vertexDelta.Value)
			if vd != nil {
				offsetStride := (vertexDelta.Index*v.StrideSize + vboVertexSize) * 4
				// 必要な場合にのみ部分更新
				gl.BufferSubData(v.target, offsetStride, len(vd)*4, gl.Ptr(vd))
			}
		}
	} else {
		gl.BufferData(v.target, v.size, v.ptr, gl.STATIC_DRAW)
	}

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

func NewVBOForBone(ptr unsafe.Pointer, count int) *VBO {
	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    ptr,
	}
	// 頂点構造体のサイズ(全部floatとする)
	// position(3), deformBoneIndex(4), deformBoneWeight(4), color(4)
	vbo.StrideSize = 3 + 4 + 4 + 4
	vbo.stride = int32(4 * vbo.StrideSize)
	vbo.size = count * 4

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindBone(boneIndexes []int, boneDeltas [][]float32) {
	gl.BindBuffer(v.target, v.id)

	if boneIndexes != nil {
		vboVertexSize := (3 + 4 + 4)

		// モーフ分の変動量を設定
		for i, vidx := range boneIndexes {
			vd := boneDeltas[i]
			offsetStride := (vidx*v.StrideSize + vboVertexSize) * 4
			// 必要な場合にのみ部分更新
			gl.BufferSubData(v.target, offsetStride, len(vd)*4, gl.Ptr(vd))
		}
	} else {
		gl.BufferData(v.target, v.size, v.ptr, gl.STATIC_DRAW)
	}

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

	// 1: deformBoneIndex
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(
		1,
		4,
		gl.FLOAT,
		false,
		v.stride,
		3*4,
	)

	// 2: deformBoneWeight
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(
		2,
		4,
		gl.FLOAT,
		false,
		v.stride,
		7*4,
	)

	// 3: color
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(
		3,
		4,
		gl.FLOAT,
		false,
		v.stride,
		11*4,
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
	// position(3), color(4)
	vbo.stride = int32(4 * (3 + 4))

	return vbo
}

// Binds VBO for rendering.
func (v *VBO) BindDebug(vertices []float32) {
	// verticesの要素数 * float32のサイズ
	v.size = len(vertices) * 4

	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, gl.Ptr(&vertices[0]), gl.STATIC_DRAW)

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

// 床のVBOを作成
func NewVBOForFloor() (*VBO, int32) {

	// 床のラインの頂点データ
	floorVertices := make([]float32, 0)
	for x := -50; x < 0; x += 5 {
		floorVertices = append(floorVertices, float32(x))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(-50))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)

		floorVertices = append(floorVertices, float32(x))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(50))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)
	}
	for x := 0; x <= 50; x += 5 {
		floorVertices = append(floorVertices, float32(x))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(-50))

		if x == 0 {
			// 原点Z軸ライン
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 1.0)
			floorVertices = append(floorVertices, 1.0)
		} else {
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.7)
		}

		floorVertices = append(floorVertices, float32(x))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(0))

		if x == 0 {
			// 原点Z軸ライン
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 1.0)
			floorVertices = append(floorVertices, 1.0)
		} else {
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.7)
		}

		floorVertices = append(floorVertices, float32(x))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(0))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)

		floorVertices = append(floorVertices, float32(x))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(50))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)
	}
	for z := -50; z < 0; z += 5 {
		floorVertices = append(floorVertices, float32(-50))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(z))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)

		floorVertices = append(floorVertices, float32(50))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(z))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)
	}
	for z := 0; z <= 50; z += 5 {
		floorVertices = append(floorVertices, float32(-50))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(z))

		if z == 0 {
			// 原点X軸ライン
			floorVertices = append(floorVertices, 1.0)
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 1.0)
		} else {
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.7)
		}

		floorVertices = append(floorVertices, float32(0))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(z))

		if z == 0 {
			// 原点X軸ライン
			floorVertices = append(floorVertices, 1.0)
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 0.0)
			floorVertices = append(floorVertices, 1.0)
		} else {
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.9)
			floorVertices = append(floorVertices, 0.7)
		}

		floorVertices = append(floorVertices, float32(0))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(z))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)

		floorVertices = append(floorVertices, float32(50))
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, float32(z))

		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.9)
		floorVertices = append(floorVertices, 0.7)
	}

	// 原点Y軸ライン
	{
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 0.0)

		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 1.0)
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 1.0)

		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 50.0)
		floorVertices = append(floorVertices, 0.0)

		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 1.0)
		floorVertices = append(floorVertices, 0.0)
		floorVertices = append(floorVertices, 1.0)
	}

	var vboId uint32
	gl.GenBuffers(1, &vboId)

	vbo := &VBO{
		id:     vboId,
		target: gl.ARRAY_BUFFER,
		ptr:    gl.Ptr(&floorVertices[0]),
	}
	// 剛体構造体のサイズ(全部floatとする)
	// position(3)
	vbo.stride = int32(4 * (3 + 4))
	// floorVerticesの要素数 * float32のサイズ
	vbo.size = int(len(floorVertices) * 4)

	return vbo, int32(len(floorVertices))
}

// 床描画
func (v *VBO) BindFloor() {
	gl.BindBuffer(v.target, v.id)
	gl.BufferData(v.target, v.size, v.ptr, gl.STATIC_DRAW)

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
