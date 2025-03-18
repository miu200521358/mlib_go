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
func (b *VertexBufferBuilder) SetData(data unsafe.Pointer, arrayCount int) *VertexBufferBuilder {
	b.vao.Bind()
	b.vbo.Bind()

	// ストライドの設定
	stride := int32(b.strideSize * b.floatSize)
	fieldCount := 0
	for i := range b.attributes {
		b.attributes[i].Stride = stride
		fieldCount += b.attributes[i].Size
	}

	// データのバッファリング
	b.vbo.BufferData(arrayCount, fieldCount, b.floatSize, data, rendering.BufferUsageStatic)

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

func (h *VertexBufferHandle) BufferData() {
	gl.BufferData(h.VBO.target, h.VBO.size, h.VBO.ptr, gl.STATIC_DRAW)
}

// バッチ処理に必要な情報を収集
type updateInfo struct {
	offset int
	data   []float32
}

// UpdateVertexDeltas はバッファに頂点デルタを適用
func (h *VertexBufferHandle) UpdateVertexDeltas(vertices *delta.VertexMorphDeltas) {
	// 頂点モーフ用の頂点数が 0 なら、元の頂点データを再設定するなどの処理
	if vertices == nil || vertices.Length() == 0 {
		h.VBO.Bind()
		// 元のデータに戻すなら、BufferDataで再アップロード
		gl.BufferData(h.VBO.target, h.VBO.size, h.VBO.ptr, gl.STATIC_DRAW)
		return
	}

	// VBO をバインド
	h.VBO.Bind()

	// バッファ全体を一度にマップ (書き込み用)
	// MAP_INVALIDATE_RANGE_BIT は、現在のデータを破棄して書き込む際に使用
	mappedPtr := gl.MapBufferRange(h.VBO.target, 0, h.VBO.size,
		gl.MAP_WRITE_BIT|gl.MAP_INVALIDATE_RANGE_BIT)
	if mappedPtr == nil {
		// マッピング失敗時のエラーハンドリング
		return
	}
	defer gl.UnmapBuffer(h.VBO.target) // 関数終了時にアンマップ

	// mappedPtr を []float32 に変換する
	mappedSlice := unsafe.Slice((*float32)(mappedPtr), h.VBO.size/h.FloatSize)

	// もし「元データをベースにして、その上にモーフの差分を適用する」場合は、
	// まず mappedSlice に元の頂点配列 h.VBO.ptr をコピーしてから差分適用します。
	// (h.VBO.ptr が []float32 であることを仮定した例)
	baseSlice := unsafe.Slice((*float32)(h.VBO.ptr), h.VBO.size/h.FloatSize)
	copy(mappedSlice, baseSlice)

	// 以下、頂点モーフの差分を書き込む
	vboVertexSize := 3 + 3 + 2 + 2 + 1 + 4 + 4 + 1 + 3 + 3 + 3

	vertices.ForEach(func(vidx int, vertexDelta *delta.VertexMorphDelta) {
		if vertexDelta.IsZero() {
			return
		}

		vd := newVertexMorphDeltaGl(vertexDelta)
		if vd == nil {
			return
		}

		offsetStride := (vidx*h.StrideSize + vboVertexSize)
		copy(mappedSlice[offsetStride:offsetStride+len(vd)], vd)
	})
}

// UpdateBoneDeltas は、ボーン変形のデルタを一括更新します。
func (h *VertexBufferHandle) UpdateBoneDeltas(boneIndexes []int, boneDeltas [][]float32) {
	if boneIndexes == nil || boneDeltas == nil {
		return
	}

	// 対象の VBO をバインド
	h.VBO.Bind()

	// バッファ全体を一括で書き込み用にマップする
	mappedPtr := gl.MapBufferRange(h.VBO.target, 0, h.VBO.size,
		gl.MAP_WRITE_BIT|gl.MAP_INVALIDATE_RANGE_BIT)
	if mappedPtr == nil {
		// マッピング失敗時はエラー処理
		return
	}
	// 関数終了時にアンマップ
	defer gl.UnmapBuffer(h.VBO.target)

	// バッファサイズはバイト単位なので、float32 の要素数に換算
	numFloats := int(h.VBO.size / h.FloatSize)

	// マップしたメモリを []float32 に変換する
	mappedSlice := unsafe.Slice((*float32)(mappedPtr), numFloats)

	baseSlice := unsafe.Slice((*float32)(h.VBO.ptr), h.VBO.size/h.FloatSize)
	copy(mappedSlice, baseSlice)

	// ボーン変形の領域開始位置 (float 要素数)
	vboVertexSize := (3 + 4 + 4)

	// 各ボーンごとに、対象の領域に boneDelta をコピーする
	// ここで、各頂点の先頭から (vidx * StrideSize) の後に、さらにボーン変形用領域 vboVertexSize があると想定
	for i, vidx := range boneIndexes {
		vd := boneDeltas[i]
		// マッピングしたスライスは float32 の配列になっているので、バイト単位のオフセットではなく float 単位のオフセットで指定
		offset := vidx*h.StrideSize + vboVertexSize
		// 範囲外アクセスにならないよう、コピーする長さに注意
		copy(mappedSlice[offset:offset+len(vd)], vd)
	}
}
