package delta

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

// RigidBodyDeltas は RigidBodyDelta の集合を管理する
type RigidBodyDeltas struct {
	data        []*RigidBodyDelta
	rigidBodies *pmx.RigidBodies
}

// NewRigidBodyDeltas は剛体数に合わせて RigidBodyDelta のスライスを用意する
func NewRigidBodyDeltas(rigidBodies *pmx.RigidBodies) *RigidBodyDeltas {
	return &RigidBodyDeltas{
		data:        make([]*RigidBodyDelta, rigidBodies.Length()),
		rigidBodies: rigidBodies,
	}
}

func (bds *RigidBodyDeltas) Length() int {
	return len(bds.data)
}

// Get は boneIndex に対応する RigidBodyDelta を返す
func (bds *RigidBodyDeltas) Get(boneIndex int) *RigidBodyDelta {
	if boneIndex < 0 || boneIndex >= len(bds.data) {
		return nil
	}
	return bds.data[boneIndex]
}

// GetByName はボーン名に対応する RigidBodyDelta を返す
func (bds *RigidBodyDeltas) GetByName(boneName string) *RigidBodyDelta {
	if bone, err := bds.rigidBodies.GetByName(boneName); err == nil {
		return bds.Get(bone.Index())
	}

	return nil
}

// Delete は指定のインデックスの RigidBodyDelta を削除する
func (bds *RigidBodyDeltas) Delete(boneIndex int) {
	if boneIndex < 0 || boneIndex >= len(bds.data) {
		return
	}
	bds.data[boneIndex] = nil
}

// Update は RigidBodyDelta をデータにセットする
func (bds *RigidBodyDeltas) Update(bd *RigidBodyDelta) {
	if bd == nil || bd.RigidBody == nil {
		return
	}
	idx := bd.RigidBody.Index()
	if idx >= 0 && idx < len(bds.data) {
		bds.data[idx] = bd
	}
}

// Contains は指定のインデックスに RigidBodyDelta が存在するかを返す
func (bds *RigidBodyDeltas) Contains(boneIndex int) bool {
	return boneIndex >= 0 && boneIndex < len(bds.data) && bds.data[boneIndex] != nil
}

func (bds *RigidBodyDeltas) ContainsByName(boneName string) bool {
	if bone, err := bds.rigidBodies.GetByName(boneName); err == nil {
		return bds.Contains(bone.Index())
	}
	return false
}

// ForEach は全ての値をコールバック関数に渡します
func (bds *RigidBodyDeltas) ForEach(callback func(index int, value *RigidBodyDelta) bool) {
	for i, v := range bds.data {
		if !callback(i, v) {
			break
		}
	}
}
