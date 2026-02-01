// 指示: miu200521358
package graphics_api

import (
	"unsafe"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// BufferType はバッファ種別を表す。
type BufferType int

const (
	// BufferTypeVertex は頂点バッファ。
	BufferTypeVertex BufferType = iota
	// BufferTypeIndex はインデックスバッファ。
	BufferTypeIndex
)

// BufferUsage はバッファ用途を表す。
type BufferUsage int

const (
	// BufferUsageStatic は静的用途。
	BufferUsageStatic BufferUsage = iota
	// BufferUsageDynamic は動的用途。
	BufferUsageDynamic
	// BufferUsageStream はストリーム用途。
	BufferUsageStream
)

// VertexAttribute は頂点属性を表す。
type VertexAttribute struct {
	Index     uint32
	Size      int
	Type      uint32
	Normalize bool
	Stride    int32
	Offset    int
}

// IVertexBuffer は頂点バッファの抽象。
type IVertexBuffer interface {
	// Bind はバッファをバインドする。
	Bind()
	// Unbind はバッファをアンバインドする。
	Unbind()
	// BufferData はバッファ全体を転送する。
	BufferData(size int, data unsafe.Pointer, usage BufferUsage)
	// BufferSubData は部分更新を行う。
	BufferSubData(offset int, size int, data unsafe.Pointer)
	// Delete はバッファを削除する。
	Delete()
	// GetID はバッファIDを返す。
	GetID() uint32
	// SetAttribute は頂点属性を設定する。
	SetAttribute(attr VertexAttribute)
}

// IIndexBuffer はインデックスバッファの抽象。
type IIndexBuffer interface {
	// Bind はバッファをバインドする。
	Bind()
	// Unbind はバッファをアンバインドする。
	Unbind()
	// BufferData はバッファ全体を転送する。
	BufferData(size int, data unsafe.Pointer, usage BufferUsage)
	// Delete はバッファを削除する。
	Delete()
	// GetID はバッファIDを返す。
	GetID() uint32
}

// IVertexArray は頂点配列の抽象。
type IVertexArray interface {
	// Bind は配列をバインドする。
	Bind()
	// Unbind は配列をアンバインドする。
	Unbind()
	// Delete は配列を削除する。
	Delete()
	// GetID は配列IDを返す。
	GetID() uint32
}

// DebugBoneHover はデバッグカーソル下のボーン情報を保持する。
type DebugBoneHover struct {
	ModelIndex int         // モデルインデックス
	Bone       *model.Bone // 検出されたボーン
	Distance   float64     // カーソルからボーンラインまでの最短距離
}
