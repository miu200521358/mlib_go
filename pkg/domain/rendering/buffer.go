//go:build windows
// +build windows

package rendering

import "unsafe"

// BufferType はバッファの種類を表す
type BufferType int

const (
	BufferTypeVertex  BufferType = iota // 頂点バッファ
	BufferTypeElement                   // 要素（インデックス）バッファ
)

// BufferUsage はバッファの使用方法を表す
type BufferUsage int

const (
	BufferUsageStatic  BufferUsage = iota // 静的データ
	BufferUsageDynamic                    // 動的データ
	BufferUsageStream                     // ストリーミングデータ
)

// VertexAttribute は頂点属性を表す
type VertexAttribute struct {
	Index     uint32 // 属性インデックス
	Size      int    // 要素数（例: vec3なら3）
	Type      uint32 // データ型（OpenGL定数に対応）
	Normalize bool   // 正規化するか
	Stride    int32  // ストライド
	Offset    int    // オフセット
}

// IVertexBuffer は頂点バッファのインターフェース
type IVertexBuffer interface {
	Bind()
	Unbind()
	BufferData(size int, data unsafe.Pointer, usage BufferUsage)
	BufferSubData(offset int, size int, data unsafe.Pointer)
	Delete()
	GetID() uint32
	SetAttribute(attr VertexAttribute)
}

// IElementBuffer はインデックスバッファのインターフェース
type IElementBuffer interface {
	Bind()
	Unbind()
	BufferData(size int, data unsafe.Pointer, usage BufferUsage)
	Delete()
	GetID() uint32
}

// IVertexArray は頂点配列のインターフェース
type IVertexArray interface {
	Bind()
	Unbind()
	Delete()
	GetID() uint32
}
