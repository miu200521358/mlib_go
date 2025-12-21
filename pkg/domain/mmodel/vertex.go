package mmodel

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mcore"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

// Vertex は頂点を表します。
type Vertex struct {
	mcore.IndexModel               // インデックス
	Position         *mmath.Vec3   // 頂点位置
	Normal           *mmath.Vec3   // 頂点法線
	Uv               *mmath.Vec2   // UV座標
	ExtendedUvs      []*mmath.Vec4 // 追加UV（最大4つ）
	DeformType       DeformType    // ウェイト変形方式
	Deform           IDeform       // デフォーム
	EdgeFactor       float64       // エッジ倍率
	MaterialIndexes  []int         // 割り当て材質インデックス
}

// NewVertex は新しいVertexを生成します。
func NewVertex() *Vertex {
	return &Vertex{
		IndexModel:      *mcore.NewIndexModel(-1),
		Position:        mmath.NewVec3(),
		Normal:          mmath.NewVec3(),
		Uv:              mmath.NewVec2(),
		ExtendedUvs:     make([]*mmath.Vec4, 0),
		DeformType:      DEFORM_BDEF1,
		Deform:          nil,
		EdgeFactor:      0.0,
		MaterialIndexes: make([]int, 0),
	}
}

// IsValid はVertexが有効かどうかを返します。
func (v *Vertex) IsValid() bool {
	return v != nil && v.IndexModel.IsValid() && v.Position != nil
}

// Copy は深いコピーを作成します。
func (v *Vertex) Copy() (*Vertex, error) {
	cp := &Vertex{}
	if err := deepcopy.Copy(cp, v); err != nil {
		return nil, err
	}
	return cp, nil
}
