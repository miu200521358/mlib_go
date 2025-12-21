package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
)

// Face は三角形面を表します。
type Face struct {
	mcore.IndexModel        // インデックス
	VertexIndexes    [3]int // 頂点インデックス（3つ）
}

// NewFace は新しいFaceを生成します。
func NewFace() *Face {
	return &Face{
		IndexModel:    *mcore.NewIndexModel(-1),
		VertexIndexes: [3]int{-1, -1, -1},
	}
}

// NewFaceByIndexes は頂点インデックスを指定してFaceを生成します。
func NewFaceByIndexes(v0, v1, v2 int) *Face {
	return &Face{
		IndexModel:    *mcore.NewIndexModel(-1),
		VertexIndexes: [3]int{v0, v1, v2},
	}
}

// IsValid はFaceが有効かどうかを返します。
func (f *Face) IsValid() bool {
	if f == nil || !f.IndexModel.IsValid() {
		return false
	}
	// 全頂点インデックスが有効かチェック
	for _, idx := range f.VertexIndexes {
		if idx < 0 {
			return false
		}
	}
	return true
}

// Copy は深いコピーを作成します。
func (f *Face) Copy() (*Face, error) {
	return &Face{
		IndexModel:    *mcore.NewIndexModel(f.Index()),
		VertexIndexes: f.VertexIndexes, // 配列はコピー
	}, nil
}
