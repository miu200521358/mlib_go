// 指示: miu200521358
package model

// Face は三角面を表す。
type Face struct {
	index         int
	VertexIndexes [3]int
}

// Index は面 index を返す。
func (f *Face) Index() int {
	return f.index
}

// SetIndex は面 index を設定する。
func (f *Face) SetIndex(index int) {
	f.index = index
}

// IsValid は面が有効か判定する。
func (f *Face) IsValid() bool {
	return f != nil && f.index >= 0
}
