// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
)

// RigidBodyDeltas は剛体差分の集合を表す。
type RigidBodyDeltas struct {
	data        []*RigidBodyDelta
	rigidBodies *collection.NamedCollection[*model.RigidBody]
}

// NewRigidBodyDeltas はRigidBodyDeltasを生成する。
func NewRigidBodyDeltas(rigidBodies *collection.NamedCollection[*model.RigidBody]) *RigidBodyDeltas {
	length := 0
	if rigidBodies != nil {
		length = rigidBodies.Len()
	}
	return &RigidBodyDeltas{
		data:        make([]*RigidBodyDelta, length),
		rigidBodies: rigidBodies,
	}
}

// Len は要素数を返す。
func (d *RigidBodyDeltas) Len() int {
	if d == nil {
		return 0
	}
	return len(d.data)
}

// Get はindexの差分を返す。
func (d *RigidBodyDeltas) Get(index int) *RigidBodyDelta {
	if d == nil || index < 0 || index >= len(d.data) {
		return nil
	}
	return d.data[index]
}

// GetByName は名前に対応する差分を返す。
func (d *RigidBodyDeltas) GetByName(name string) *RigidBodyDelta {
	if d == nil || d.rigidBodies == nil {
		return nil
	}
	body, err := d.rigidBodies.GetByName(name)
	if err != nil {
		return nil
	}
	return d.Get(body.Index())
}

// Update は差分を更新する。
func (d *RigidBodyDeltas) Update(delta *RigidBodyDelta) {
	if d == nil || delta == nil || delta.RigidBody == nil {
		return
	}
	idx := delta.RigidBody.Index()
	if idx < 0 || idx >= len(d.data) {
		return
	}
	d.data[idx] = delta
}

// Delete はindexの差分を削除する。
func (d *RigidBodyDeltas) Delete(index int) {
	if d == nil || index < 0 || index >= len(d.data) {
		return
	}
	d.data[index] = nil
}

// Contains はindexの差分が存在するか判定する。
func (d *RigidBodyDeltas) Contains(index int) bool {
	if d == nil || index < 0 || index >= len(d.data) {
		return false
	}
	return d.data[index] != nil
}

// ContainsByName は名前に対応する差分が存在するか判定する。
func (d *RigidBodyDeltas) ContainsByName(name string) bool {
	if d == nil || d.rigidBodies == nil {
		return false
	}
	body, err := d.rigidBodies.GetByName(name)
	if err != nil {
		return false
	}
	return d.Contains(body.Index())
}

// ForEach は全要素を走査する。
func (d *RigidBodyDeltas) ForEach(fn func(index int, delta *RigidBodyDelta) bool) {
	if d == nil || fn == nil {
		return
	}
	for i, v := range d.data {
		if !fn(i, v) {
			return
		}
	}
}
