package delta

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

// JointDeltas は JointDelta の集合を管理する
type JointDeltas struct {
	data   []*JointDelta
	joints *pmx.Joints
}

// NewJointDeltas はジョイント数に合わせて JointDelta のスライスを用意する
func NewJointDeltas(joints *pmx.Joints) *JointDeltas {
	return &JointDeltas{
		data:   make([]*JointDelta, joints.Length()),
		joints: joints,
	}
}

func (bds *JointDeltas) Length() int {
	return len(bds.data)
}

// Get は jointIndex に対応する JointDelta を返す
func (bds *JointDeltas) Get(jointIndex int) *JointDelta {
	if jointIndex < 0 || jointIndex >= len(bds.data) {
		return nil
	}
	return bds.data[jointIndex]
}

// GetByName はボーン名に対応する JointDelta を返す
func (bds *JointDeltas) GetByName(jointName string) *JointDelta {
	if joint, err := bds.joints.GetByName(jointName); err == nil {
		return bds.Get(joint.Index())
	}

	return nil
}

// Delete は指定のインデックスの JointDelta を削除する
func (bds *JointDeltas) Delete(jointIndex int) {
	if jointIndex < 0 || jointIndex >= len(bds.data) {
		return
	}
	bds.data[jointIndex] = nil
}

// Update は JointDelta をデータにセットする
func (bds *JointDeltas) Update(bd *JointDelta) {
	if bd == nil || bd.Joint == nil {
		return
	}
	idx := bd.Joint.Index()
	if idx >= 0 && idx < len(bds.data) {
		bds.data[idx] = bd
	}
}

// Contains は指定のインデックスに JointDelta が存在するかを返す
func (bds *JointDeltas) Contains(jointIndex int) bool {
	return jointIndex >= 0 && jointIndex < len(bds.data) && bds.data[jointIndex] != nil
}

func (bds *JointDeltas) ContainsByName(jointName string) bool {
	if joint, err := bds.joints.GetByName(jointName); err == nil {
		return bds.Contains(joint.Index())
	}
	return false
}

// ForEach は全ての値をコールバック関数に渡します
func (bds *JointDeltas) ForEach(callback func(index int, value *JointDelta) bool) {
	for i, v := range bds.data {
		if !callback(i, v) {
			break
		}
	}
}
