// 指示: miu200521358
package mphysics

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// TestBuildPhysicsDeltasUseRigidBodyMassWhenFrameMassIsNil は質量未指定時に剛体既定値を使うことを確認する。
func TestBuildPhysicsDeltasUseRigidBodyMassWhenFrameMassIsNil(t *testing.T) {
	modelData := model.NewPmxModel()
	rigidBodySize := mmath.NewVec3()
	rigidBodySize.X = 1
	rigidBodySize.Y = 2
	rigidBodySize.Z = 3
	rigidBody := &model.RigidBody{
		Size:  rigidBodySize,
		Param: model.RigidBodyParam{Mass: 4.2},
	}
	rigidBody.SetIndex(0)
	rigidBody.SetName("rb")
	modelData.RigidBodies.AppendRaw(rigidBody)

	motionData := motion.NewVmdMotion("")
	rf := motion.NewRigidBodyFrame(0)
	frameSize := mmath.NewVec3()
	frameSize.X = 5
	frameSize.Y = 6
	frameSize.Z = 7
	rf.Size = &frameSize
	motionData.AppendRigidBodyFrame("rb", rf)

	deltas := BuildPhysicsDeltas(modelData, motionData, 0)
	if deltas == nil || deltas.RigidBodies == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	rigidDelta := deltas.RigidBodies.Get(0)
	if rigidDelta == nil {
		t.Fatalf("rigid body delta should not be nil")
	}
	expectedSize := mmath.NewVec3()
	expectedSize.X = 5
	expectedSize.Y = 6
	expectedSize.Z = 7
	if !rigidDelta.Size.NearEquals(expectedSize, 1e-10) {
		t.Fatalf("size mismatch: got=%v", rigidDelta.Size)
	}
	if !mmath.NearEquals(rigidDelta.Mass, 4.2, 1e-10) {
		t.Fatalf("mass mismatch: got=%v", rigidDelta.Mass)
	}
}

// TestBuildPhysicsDeltasUseFrameMassWhenSpecified は質量指定時にフレーム値を使うことを確認する。
func TestBuildPhysicsDeltasUseFrameMassWhenSpecified(t *testing.T) {
	modelData := model.NewPmxModel()
	rigidBody := &model.RigidBody{Param: model.RigidBodyParam{Mass: 9.9}}
	rigidBody.SetIndex(0)
	rigidBody.SetName("rb")
	modelData.RigidBodies.AppendRaw(rigidBody)

	motionData := motion.NewVmdMotion("")
	rf := motion.NewRigidBodyFrame(0)
	mass := 0.0
	rf.Mass = &mass
	motionData.AppendRigidBodyFrame("rb", rf)

	deltas := BuildPhysicsDeltas(modelData, motionData, 0)
	if deltas == nil || deltas.RigidBodies == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	rigidDelta := deltas.RigidBodies.Get(0)
	if rigidDelta == nil {
		t.Fatalf("rigid body delta should not be nil")
	}
	if !mmath.NearEquals(rigidDelta.Mass, 0.0, 1e-10) {
		t.Fatalf("mass mismatch: got=%v", rigidDelta.Mass)
	}
}
