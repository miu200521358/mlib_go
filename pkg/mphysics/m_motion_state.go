package mphysics

import "github.com/miu200521358/mlib_go/pkg/mbt"

type MMotionState struct {
	mbt.BtDefaultMotionState
}

func NewMMotionState() *MMotionState {
	return &MMotionState{
		BtDefaultMotionState: mbt.NewBtDefaultMotionState(),
	}
}

// ボーン追従剛体用MotionState
type StaticMotionState struct {
	MMotionState
}

func NewStaticMotionState() *StaticMotionState {
	return &StaticMotionState{
		MMotionState: *NewMMotionState(),
	}
}

// 物理剛体用MotionState
type DynamicMotionState struct {
	MMotionState
}

func NewDynamicMotionState() *DynamicMotionState {
	return &DynamicMotionState{
		MMotionState: *NewMMotionState(),
	}
}

// 物理剛体+ボーン位置合わせ用MotionState
type DynamicBoneMotionState struct {
	MMotionState
}

func NewDynamicBoneMotionState() *DynamicBoneMotionState {
	return &DynamicBoneMotionState{
		MMotionState: *NewMMotionState(),
	}
}
