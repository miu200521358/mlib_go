//go:build windows
// +build windows

package mgl

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
)

// ModelRenderer はモデル描画を担当
type ModelRenderer struct {
	bufferFactory *BufferFactory
	vertexBuffer  *VertexBufferHandle
	indexBuffer   *ElementBuffer
	indexCount    int
}

// NewModelRenderer は新しいモデルレンダラーを作成
func NewModelRenderer(
	vertexPtr unsafe.Pointer, vertexCount int, indexPtr unsafe.Pointer, indexCount int,
) *ModelRenderer {
	factory := NewBufferFactory()
	vertexBuffer := factory.CreateVertexBuffer(vertexPtr, vertexCount)

	indexBuffer := factory.CreateElementBuffer(indexPtr, indexCount)

	return &ModelRenderer{
		bufferFactory: factory,
		vertexBuffer:  vertexBuffer,
		indexBuffer:   indexBuffer,
		indexCount:    indexCount,
	}
}

// UpdateVertexMorphs は頂点モーフの更新
func (r *ModelRenderer) UpdateVertexMorphs(morphDeltas *delta.VertexMorphDeltas) {
	r.vertexBuffer.UpdateVertexDeltas(morphDeltas)
}

// Render はモデルの描画
func (r *ModelRenderer) Render() {
	r.vertexBuffer.Bind()
	r.indexBuffer.Bind()

	gl.DrawElements(gl.TRIANGLES, int32(r.indexCount), gl.UNSIGNED_INT, nil)

	r.indexBuffer.Unbind()
	r.vertexBuffer.Unbind()
}

// Delete はリソース解放
func (r *ModelRenderer) Delete() {
	r.vertexBuffer.Delete()
	r.indexBuffer.Delete()
}
