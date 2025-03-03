//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

// VertexBufferConfig は頂点バッファの設定
type VertexBufferConfig struct {
	Attributes []rendering.VertexAttribute
	Data       unsafe.Pointer
	Count      int
	FloatSize  int
}

// AttributeType は標準頂点属性タイプ
type AttributeType int

const (
	AttributePosition         AttributeType = iota // 位置
	AttributeNormal                                // 法線
	AttributeUV                                    // UV
	AttributeExtendedUV                            // 拡張UV
	AttributeEdgeFactor                            // エッジの太さ
	AttributeDeformBoneIndex                       // ボーンインデックス
	AttributeDeformBoneWeight                      // ボーンウェイト
	AttributeIsSdef                                // SDEFかどうか
	AttributeSdefC                                 // SDEFのC
	AttributeSdefR0                                // SDEFのR0
	AttributeSdefR1                                // SDEFのR1
	AttributeVertexDelta                           // 頂点変形情報
	AttributeUVDelta                               // UV変形情報
	AttributeUV1Delta                              // UV1変形情報
	AttributeAfterVertexDelta                      // 頂点変形後情報
	AttributeColor                                 // 色
	AttributeTexCoords                             // テクスチャ座標
)

// VertexBufferBuilder は頂点バッファビルダー
type VertexBufferBuilder struct {
	vao         *VertexArray
	vbo         *VertexBuffer
	attributes  []rendering.VertexAttribute
	strideSize  int
	floatSize   int
	elementSize int
}

// NewVertexBufferBuilder は新しい頂点バッファビルダーを作成
func NewVertexBufferBuilder() *VertexBufferBuilder {
	return &VertexBufferBuilder{
		vao:         NewVertexArray(),
		vbo:         NewVertexBuffer(),
		attributes:  []rendering.VertexAttribute{},
		strideSize:  0,
		floatSize:   4, // float32のサイズ
		elementSize: 0,
	}
}

// AddAttribute は頂点属性を追加
func (b *VertexBufferBuilder) AddAttribute(attrType AttributeType, size int) *VertexBufferBuilder {
	offset := b.strideSize * b.floatSize

	attr := rendering.VertexAttribute{
		Index:     uint32(len(b.attributes)),
		Size:      size,
		Type:      gl.FLOAT,
		Normalize: false,
		Offset:    offset,
	}

	b.attributes = append(b.attributes, attr)
	b.strideSize += size
	b.elementSize++

	return b
}

// AddStandardVertexAttributes は標準的な頂点属性を追加
func (b *VertexBufferBuilder) AddStandardVertexAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition, 3).
		AddAttribute(AttributeNormal, 3).
		AddAttribute(AttributeUV, 2).
		AddAttribute(AttributeExtendedUV, 2).
		AddAttribute(AttributeEdgeFactor, 1).
		AddAttribute(AttributeDeformBoneIndex, 4).
		AddAttribute(AttributeDeformBoneWeight, 4).
		AddAttribute(AttributeIsSdef, 1).
		AddAttribute(AttributeSdefC, 3).
		AddAttribute(AttributeSdefR0, 3).
		AddAttribute(AttributeSdefR1, 3).
		AddAttribute(AttributeVertexDelta, 3).
		AddAttribute(AttributeUVDelta, 4).
		AddAttribute(AttributeUV1Delta, 4).
		AddAttribute(AttributeAfterVertexDelta, 3)
}

// AddDebugAttributes は位置と色の属性を追加
func (b *VertexBufferBuilder) AddDebugAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition, 3).
		AddAttribute(AttributeColor, 4)
}

// AddBoneAttributes はボーンの属性を追加
func (b *VertexBufferBuilder) AddBoneAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition, 3).
		AddAttribute(AttributeDeformBoneIndex, 4).
		AddAttribute(AttributeDeformBoneWeight, 4).
		AddAttribute(AttributeColor, 4)
}

// AddOverrideAttributes はオーバーライドの属性を追加
func (b *VertexBufferBuilder) AddOverrideAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition, 3).
		AddAttribute(AttributeTexCoords, 2)
}

// SetData はデータを設定
func (b *VertexBufferBuilder) SetData(data unsafe.Pointer, count int) *VertexBufferBuilder {
	b.vao.Bind()
	b.vbo.Bind()

	// ストライドの設定
	stride := int32(b.strideSize * b.floatSize)
	for i := range b.attributes {
		b.attributes[i].Stride = stride
	}

	// データのバッファリング
	size := count * b.strideSize * b.floatSize
	b.vbo.BufferData(size, data, rendering.BufferUsageStatic)

	// 属性の設定
	for _, attr := range b.attributes {
		b.vbo.SetAttribute(attr)
	}

	b.vbo.Unbind()
	b.vao.Unbind()

	return b
}

// Build はビルダーから頂点バッファ構成を構築
func (b *VertexBufferBuilder) Build() *VertexBufferHandle {
	return &VertexBufferHandle{
		VAO:        b.vao,
		VBO:        b.vbo,
		StrideSize: b.strideSize,
		FloatSize:  b.floatSize,
		Attributes: b.attributes,
	}
}

// VertexBufferHandle は頂点バッファへのハンドル
type VertexBufferHandle struct {
	VAO        *VertexArray
	VBO        *VertexBuffer
	StrideSize int
	FloatSize  int
	Attributes []rendering.VertexAttribute
}

// Bind はVAOをバインド
func (h *VertexBufferHandle) Bind() {
	h.VAO.Bind()
}

// Unbind はVAOをアンバインド
func (h *VertexBufferHandle) Unbind() {
	h.VAO.Unbind()
}

// Delete はリソースを解放
func (h *VertexBufferHandle) Delete() {
	h.VBO.Delete()
	h.VAO.Delete()
}

// UpdateVertexDeltas はバッファに頂点デルタを適用
func (h *VertexBufferHandle) UpdateVertexDeltas(vertices *delta.VertexMorphDeltas) {
	if vertices == nil || vertices.Length() == 0 {
		return
	}

	h.VAO.Bind()
	h.VBO.Bind()

	// 頂点モーフのVBOサイズ
	vboVertexSize := (3 + 3 + 2 + 2 + 1 + 4 + 4 + 1 + 3 + 3 + 3)

	// モーフ分の変動量を設定
	for v := range vertices.Iterator() {
		vidx := v.Index
		vertexDelta := v.Value

		vd := newVertexMorphDeltaGl(vertexDelta)
		if vd != nil {
			offsetStride := (vidx*h.StrideSize + vboVertexSize) * h.FloatSize
			h.VBO.BufferSubData(offsetStride, len(vd)*h.FloatSize, gl.Ptr(vd))
		}
	}

	h.VBO.Unbind()
	h.VAO.Unbind()
}
