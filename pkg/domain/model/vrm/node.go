// 指示: miu200521358
package vrm

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// Node はVRMノード情報を表す。
type Node struct {
	Index       int
	Name        string
	ParentIndex int
	Children    []int
	Translation mmath.Vec3
}

// NewNode はNodeを既定値で生成する。
func NewNode(index int) *Node {
	return &Node{
		Index:       index,
		ParentIndex: -1,
		Children:    make([]int, 0),
	}
}
