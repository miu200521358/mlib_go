package mphysics

import (
	"github.com/miu200521358/mlib_go/pkg/mbt"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type MMotionStateInterface interface {
	mbt.BtDefaultMotionState
	Reset()
}

type MMotionState struct {
	mbt.BtDefaultMotionState
	BoneTransform         mbt.BtTransform
	OffsetRigidBodyMatrix *mmath.MMat4
}

func NewMMotionState(boneTransform mbt.BtTransform, offsetRigidBodyMatrix *mmath.MMat4) MMotionState {
	return MMotionState{
		BtDefaultMotionState:  mbt.NewBtDefaultMotionState(boneTransform),
		BoneTransform:         boneTransform,
		OffsetRigidBodyMatrix: offsetRigidBodyMatrix,
	}
}

func (m *MMotionState) Swigcptr() uintptr {
	return m.BtDefaultMotionState.Swigcptr()
}

func (m *MMotionState) SwigIsBtDefaultMotionState() {
	m.BtDefaultMotionState.SwigIsBtDefaultMotionState()
}

func (m *MMotionState) SetM_graphicsWorldTrans(arg2 mbt.BtTransform) {
	m.BtDefaultMotionState.SetM_graphicsWorldTrans(arg2)
}

func (m *MMotionState) GetM_graphicsWorldTrans() mbt.BtTransform {
	return m.BtDefaultMotionState.GetM_graphicsWorldTrans()
}

func (m *MMotionState) SetM_centerOfMassOffset(arg2 mbt.BtTransform) {
	m.BtDefaultMotionState.SetM_centerOfMassOffset(arg2)
}

func (m *MMotionState) GetM_centerOfMassOffset() mbt.BtTransform {
	return m.BtDefaultMotionState.GetM_centerOfMassOffset()
}

func (m *MMotionState) SetM_startWorldTrans(arg2 mbt.BtTransform) {
	m.BtDefaultMotionState.SetM_startWorldTrans(arg2)
}

func (m *MMotionState) GetM_startWorldTrans() mbt.BtTransform {
	return m.BtDefaultMotionState.GetM_startWorldTrans()
}

func (m *MMotionState) SetM_userPointer(arg2 uintptr) {
	m.BtDefaultMotionState.SetM_userPointer(arg2)
}

func (m *MMotionState) GetM_userPointer() uintptr {
	return m.BtDefaultMotionState.GetM_userPointer()
}

func (m *MMotionState) GetWorldTransform(arg2 mbt.BtTransform) {
	m.BtDefaultMotionState.GetWorldTransform(arg2)
}

func (m *MMotionState) SetWorldTransform(arg2 mbt.BtTransform) {
	m.BtDefaultMotionState.SetWorldTransform(arg2)
}

func (m *MMotionState) SwigIsBtMotionState() {
	m.BtDefaultMotionState.SwigIsBtMotionState()
}

func (m *MMotionState) SwigGetBtMotionState() mbt.BtMotionState {
	return m.BtDefaultMotionState.SwigGetBtMotionState()
}

func (m *MMotionState) Reset() {
}

// ボーン追従剛体用MotionState
type StaticMotionState struct {
	MMotionState
}

func NewStaticMotionState(
	boneTransform mbt.BtTransform,
	offsetRigidBodyMatrix *mmath.MMat4,
) *StaticMotionState {
	return &StaticMotionState{MMotionState: NewMMotionState(boneTransform, offsetRigidBodyMatrix)}
}

// 物理剛体用MotionState
type DynamicMotionState struct {
	MMotionState
}

func NewDynamicMotionState(
	boneTransform mbt.BtTransform,
	offsetRigidBodyMatrix *mmath.MMat4,
) *DynamicMotionState {
	return &DynamicMotionState{MMotionState: NewMMotionState(boneTransform, offsetRigidBodyMatrix)}
}

func (m *DynamicMotionState) Reset() {
	m.BoneTransform.SetFromOpenGLMatrix(&m.OffsetRigidBodyMatrix.GL()[0])
}

// 物理剛体+ボーン位置合わせ用MotionState
type DynamicBoneMotionState struct {
	MMotionState
}

func NewDynamicBoneMotionState(
	boneTransform mbt.BtTransform,
	offsetRigidBodyMatrix *mmath.MMat4,
) *DynamicBoneMotionState {
	return &DynamicBoneMotionState{MMotionState: NewMMotionState(boneTransform, offsetRigidBodyMatrix)}
}

func (m *DynamicBoneMotionState) Reset() {
	m.BoneTransform.SetFromOpenGLMatrix(&m.OffsetRigidBodyMatrix.GL()[0])
}
