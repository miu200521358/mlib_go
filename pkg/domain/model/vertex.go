// 指示: miu200521358
package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// Vertex はモデル頂点を表す。
type Vertex struct {
	index           int
	Position        mmath.Vec3
	Normal          mmath.Vec3
	Uv              mmath.Vec2
	ExtendedUvs     []mmath.Vec4
	DeformType      DeformType
	Deform          IDeform
	EdgeFactor      float64
	MaterialIndexes []int
}

// Index は頂点 index を返す。
func (v *Vertex) Index() int {
	return v.index
}

// SetIndex は頂点 index を設定する。
func (v *Vertex) SetIndex(index int) {
	v.index = index
}

// IsValid は頂点が有効か判定する。
func (v *Vertex) IsValid() bool {
	return v != nil && v.index >= 0
}
