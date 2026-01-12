// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
)

// JointDeltas はジョイント差分の集合を表す。
type JointDeltas struct {
	data   []*JointDelta
	joints *collection.NamedCollection[*model.Joint]
}

// NewJointDeltas はJointDeltasを生成する。
func NewJointDeltas(joints *collection.NamedCollection[*model.Joint]) *JointDeltas {
	length := 0
	if joints != nil {
		length = joints.Len()
	}
	return &JointDeltas{
		data:   make([]*JointDelta, length),
		joints: joints,
	}
}

// Len は要素数を返す。
func (d *JointDeltas) Len() int {
	if d == nil {
		return 0
	}
	return len(d.data)
}

// Get はindexの差分を返す。
func (d *JointDeltas) Get(index int) *JointDelta {
	if d == nil || index < 0 || index >= len(d.data) {
		return nil
	}
	return d.data[index]
}

// GetByName は名前に対応する差分を返す。
func (d *JointDeltas) GetByName(name string) *JointDelta {
	if d == nil || d.joints == nil {
		return nil
	}
	joint, err := d.joints.GetByName(name)
	if err != nil {
		return nil
	}
	return d.Get(joint.Index())
}

// Update は差分を更新する。
func (d *JointDeltas) Update(delta *JointDelta) {
	if d == nil || delta == nil || delta.Joint == nil {
		return
	}
	idx := delta.Joint.Index()
	if idx < 0 || idx >= len(d.data) {
		return
	}
	d.data[idx] = delta
}

// Delete はindexの差分を削除する。
func (d *JointDeltas) Delete(index int) {
	if d == nil || index < 0 || index >= len(d.data) {
		return
	}
	d.data[index] = nil
}

// Contains はindexの差分が存在するか判定する。
func (d *JointDeltas) Contains(index int) bool {
	if d == nil || index < 0 || index >= len(d.data) {
		return false
	}
	return d.data[index] != nil
}

// ContainsByName は名前に対応する差分が存在するか判定する。
func (d *JointDeltas) ContainsByName(name string) bool {
	if d == nil || d.joints == nil {
		return false
	}
	joint, err := d.joints.GetByName(name)
	if err != nil {
		return false
	}
	return d.Contains(joint.Index())
}

// ForEach は全要素を走査する。
func (d *JointDeltas) ForEach(fn func(index int, delta *JointDelta) bool) {
	if d == nil || fn == nil {
		return
	}
	for i, v := range d.data {
		if !fn(i, v) {
			return
		}
	}
}
