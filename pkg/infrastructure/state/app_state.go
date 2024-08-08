//go:build windows
// +build windows

package state

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type IAppState interface {
	Frame() float32
	SetFrame(frame float32)
	MaxFrame() float32
	SetMaxFrame(maxFrame float32)
	UpdateMaxFrame(maxFrame float32)
	IsEnabledFrameDrop() bool
	SetEnabledFrameDrop(enabled bool)
	IsEnabledPhysics() bool
	SetEnabledPhysics(enabled bool)
	IsPhysicsReset() bool
	SetPhysicsReset(reset bool)
	IsShowNormal() bool
	SetShowNormal(show bool)
	IsShowWire() bool
	SetShowWire(show bool)
	IsShowOverride() bool
	SetShowOverride(show bool)
	IsShowSelectedVertex() bool
	SetShowSelectedVertex(show bool)
	IsShowBoneAll() bool
	SetShowBoneAll(show bool)
	IsShowBoneIk() bool
	SetShowBoneIk(show bool)
	IsShowBoneEffector() bool
	SetShowBoneEffector(show bool)
	IsShowBoneFixed() bool
	SetShowBoneFixed(show bool)
	IsShowBoneRotate() bool
	SetShowBoneRotate(show bool)
	IsShowBoneTranslate() bool
	SetShowBoneTranslate(show bool)
	IsShowBoneVisible() bool
	SetShowBoneVisible(show bool)
	IsShowRigidBodyFront() bool
	SetShowRigidBodyFront(show bool)
	IsShowRigidBodyBack() bool
	SetShowRigidBodyBack(show bool)
	IsShowJoint() bool
	SetShowJoint(show bool)
	IsShowInfo() bool
	SetShowInfo(show bool)
	IsLimitFps30() bool
	SetLimitFps30(limit bool)
	IsLimitFps60() bool
	SetLimitFps60(limit bool)
	IsUnLimitFps() bool
	SetUnLimitFps(limit bool)
	IsCameraSync() bool
	SetCameraSync(sync bool)
	IsClosed() bool
	SetClosed(closed bool)
	Playing() bool
	SetPlaying(p bool)
	SpfLimit() float64
	SetSpfLimit(spf float64)
	SetGetModels(f func() [][]*pmx.PmxModel)
	SetGetMotions(f func() [][]*vmd.VmdMotion)
	GetModels() [][]*pmx.PmxModel
	GetMotions() [][]*vmd.VmdMotion
	SetFrameChannel(v float32)
	SetMaxFrameChannel(v float32)
	SetEnabledFrameDropChannel(v bool)
	SetEnabledPhysicsChannel(v bool)
	SetPhysicsResetChannel(v bool)
	SetShowNormalChannel(v bool)
	SetShowWireChannel(v bool)
	SetShowOverrideChannel(v bool)
	SetShowSelectedVertexChannel(v bool)
	SetShowBoneAllChannel(v bool)
	SetShowBoneIkChannel(v bool)
	SetShowBoneEffectorChannel(v bool)
	SetShowBoneFixedChannel(v bool)
	SetShowBoneRotateChannel(v bool)
	SetShowBoneTranslateChannel(v bool)
	SetShowBoneVisibleChannel(v bool)
	SetShowRigidBodyFrontChannel(v bool)
	SetShowRigidBodyBackChannel(v bool)
	SetShowJointChannel(v bool)
	SetShowInfoChannel(v bool)
	SetLimitFps30Channel(v bool)
	SetLimitFps60Channel(v bool)
	SetUnLimitFpsChannel(v bool)
	SetUnLimitFpsDeformChannel(v bool)
	SetCameraSyncChannel(v bool)
	SetClosedChannel(v bool)
	SetPlayingChannel(v bool)
	SetSpfLimitChannel(v float64)
	SetSelectedVertexIndexesChannel(v [][][]int)
}
