///go:build windows
//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
)

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

// AttributeInfo は頂点属性の情報
type AttributeInfo struct {
	Size      int    // 要素数（例：Position は3要素）
	Type      uint32 // データ型
	Normalize bool   // 正規化フラグ
	Name      string // シェーダーでの変数名（ドキュメント目的）
}

// 頂点属性タイプに対応する情報のマッピング
var attributeInfoMap = map[AttributeType]AttributeInfo{
	AttributePosition:         {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "position"},
	AttributeNormal:           {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "normal"},
	AttributeUV:               {Size: 2, Type: gl.FLOAT, Normalize: false, Name: "texCoord"},
	AttributeExtendedUV:       {Size: 2, Type: gl.FLOAT, Normalize: false, Name: "extendedUV"},
	AttributeEdgeFactor:       {Size: 1, Type: gl.FLOAT, Normalize: false, Name: "edgeFactor"},
	AttributeDeformBoneIndex:  {Size: 4, Type: gl.FLOAT, Normalize: false, Name: "deformBoneIndex"},
	AttributeDeformBoneWeight: {Size: 4, Type: gl.FLOAT, Normalize: false, Name: "deformBoneWeight"},
	AttributeIsSdef:           {Size: 1, Type: gl.FLOAT, Normalize: false, Name: "isSdef"},
	AttributeSdefC:            {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "sdefC"},
	AttributeSdefR0:           {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "sdefR0"},
	AttributeSdefR1:           {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "sdefR1"},
	AttributeVertexDelta:      {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "vertexDelta"},
	AttributeUVDelta:          {Size: 4, Type: gl.FLOAT, Normalize: false, Name: "uvDelta"},
	AttributeUV1Delta:         {Size: 4, Type: gl.FLOAT, Normalize: false, Name: "uv1Delta"},
	AttributeAfterVertexDelta: {Size: 3, Type: gl.FLOAT, Normalize: false, Name: "afterVertexDelta"},
	AttributeColor:            {Size: 4, Type: gl.FLOAT, Normalize: false, Name: "color"},
	AttributeTexCoords:        {Size: 2, Type: gl.FLOAT, Normalize: false, Name: "texCoords"},
}

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
		vao:         newVertexArray(),
		vbo:         newVertexBuffer(),
		attributes:  []rendering.VertexAttribute{},
		strideSize:  0,
		floatSize:   4, // float32のサイズ
		elementSize: 0,
	}
}

// AddAttribute は頂点属性を追加（属性タイプから自動的にサイズを取得）
func (b *VertexBufferBuilder) AddAttribute(attrType AttributeType) *VertexBufferBuilder {
	info, exists := attributeInfoMap[attrType]
	if !exists {
		// 未定義の属性タイプの場合のエラー処理
		// ログ出力するか、パニックを起こすかは要件による
		return b
	}

	offset := b.strideSize * b.floatSize

	attr := rendering.VertexAttribute{
		Index:     uint32(len(b.attributes)),
		Size:      info.Size,
		Type:      info.Type,
		Normalize: info.Normalize,
		Offset:    offset,
	}

	b.attributes = append(b.attributes, attr)
	b.strideSize += info.Size
	b.elementSize++

	return b
}

// AddAttributeWithSize は頂点属性をカスタムサイズで追加（後方互換性用）
func (b *VertexBufferBuilder) AddAttributeWithSize(attrType AttributeType, size int) *VertexBufferBuilder {
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
		AddAttribute(AttributePosition).
		AddAttribute(AttributeNormal).
		AddAttribute(AttributeUV).
		AddAttribute(AttributeExtendedUV).
		AddAttribute(AttributeEdgeFactor).
		AddAttribute(AttributeDeformBoneIndex).
		AddAttribute(AttributeDeformBoneWeight).
		AddAttribute(AttributeIsSdef).
		AddAttribute(AttributeSdefC).
		AddAttribute(AttributeSdefR0).
		AddAttribute(AttributeSdefR1).
		AddAttribute(AttributeVertexDelta).
		AddAttribute(AttributeUVDelta).
		AddAttribute(AttributeUV1Delta).
		AddAttribute(AttributeAfterVertexDelta)
}

// AddPositionColorAttributes は位置と色の属性を追加
func (b *VertexBufferBuilder) AddPositionColorAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition).
		AddAttribute(AttributeColor)
}

// AddBoneAttributes はボーンの属性を追加
func (b *VertexBufferBuilder) AddBoneAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition).
		AddAttribute(AttributeDeformBoneIndex).
		AddAttribute(AttributeDeformBoneWeight).
		AddAttribute(AttributeColor)
}

// AddOverrideAttributes はオーバーライドの属性を追加
func (b *VertexBufferBuilder) AddOverrideAttributes() *VertexBufferBuilder {
	return b.
		AddAttribute(AttributePosition).
		AddAttribute(AttributeTexCoords)
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
	size := count * b.floatSize
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

// Bind はVAO, VBOをバインド
func (h *VertexBufferHandle) Bind() {
	h.VAO.Bind()
	h.VBO.Bind()
}

// Unbind はVAO, VBOをアンバインド
func (h *VertexBufferHandle) Unbind() {
	h.VBO.Unbind()
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

	gl.BindBuffer(h.VBO.target, h.VBO.id)

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
}
