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
	expectedPosition := mmath.NewVec3()
	if !rigidDelta.Position.NearEquals(expectedPosition, 1e-10) {
		t.Fatalf("position mismatch: got=%v", rigidDelta.Position)
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

// TestBuildPhysicsDeltasUseFramePositionWhenSpecified は位置指定時にフレーム値を使うことを確認する。
func TestBuildPhysicsDeltasUseFramePositionWhenSpecified(t *testing.T) {
	modelData := model.NewPmxModel()
	rigidBody := &model.RigidBody{Position: newTestVec3(1, 2, 3)}
	rigidBody.SetIndex(0)
	rigidBody.SetName("rb")
	modelData.RigidBodies.AppendRaw(rigidBody)

	motionData := motion.NewVmdMotion("")
	rf := motion.NewRigidBodyFrame(0)
	framePosition := newTestVec3(4, 5, 6)
	rf.Position = &framePosition
	motionData.AppendRigidBodyFrame("rb", rf)

	deltas := BuildPhysicsDeltas(modelData, motionData, 0)
	if deltas == nil || deltas.RigidBodies == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	rigidDelta := deltas.RigidBodies.Get(0)
	if rigidDelta == nil {
		t.Fatalf("rigid body delta should not be nil")
	}
	if !rigidDelta.Position.NearEquals(framePosition, 1e-10) {
		t.Fatalf("position mismatch: got=%v want=%v", rigidDelta.Position, framePosition)
	}
}

// TestBuildPhysicsDeltasUseFallbackPositionSizeWhenInterpolatedRigidBodyFrameIsUnset は補間中未指定の位置/サイズで既定値を使うことを確認する。
func TestBuildPhysicsDeltasUseFallbackPositionSizeWhenInterpolatedRigidBodyFrameIsUnset(t *testing.T) {
	modelData := model.NewPmxModel()
	rigidBody := &model.RigidBody{
		Position: newTestVec3(1, 2, 3),
		Size:     newTestVec3(4, 5, 6),
		Param:    model.RigidBodyParam{Mass: 0},
	}
	rigidBody.SetIndex(0)
	rigidBody.SetName("rb")
	modelData.RigidBodies.AppendRaw(rigidBody)

	motionData := motion.NewVmdMotion("")
	rf0 := motion.NewRigidBodyFrame(0)
	mass0 := 1.0
	rf0.Mass = &mass0
	motionData.AppendRigidBodyFrame("rb", rf0)
	rf10 := motion.NewRigidBodyFrame(10)
	mass10 := 3.0
	rf10.Mass = &mass10
	motionData.AppendRigidBodyFrame("rb", rf10)

	deltas := BuildPhysicsDeltas(modelData, motionData, 5)
	if deltas == nil || deltas.RigidBodies == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	rigidDelta := deltas.RigidBodies.Get(0)
	if rigidDelta == nil {
		t.Fatalf("rigid body delta should not be nil")
	}
	if !rigidDelta.Position.NearEquals(rigidBody.Position, 1e-10) {
		t.Fatalf("position fallback mismatch: got=%v want=%v", rigidDelta.Position, rigidBody.Position)
	}
	if !rigidDelta.Size.NearEquals(rigidBody.Size, 1e-10) {
		t.Fatalf("size fallback mismatch: got=%v want=%v", rigidDelta.Size, rigidBody.Size)
	}
	if !mmath.NearEquals(rigidDelta.Mass, 2.0, 1e-10) {
		t.Fatalf("mass interpolation mismatch: got=%v want=2.0", rigidDelta.Mass)
	}
}

// TestBuildPhysicsDeltasInterpolateRigidBodyPositionWhenFramesSpecified は剛体位置キーが両端にある場合に線形補間されることを確認する。
func TestBuildPhysicsDeltasInterpolateRigidBodyPositionWhenFramesSpecified(t *testing.T) {
	modelData := model.NewPmxModel()
	rigidBody := &model.RigidBody{Position: newTestVec3(9, 9, 9)}
	rigidBody.SetIndex(0)
	rigidBody.SetName("rb")
	modelData.RigidBodies.AppendRaw(rigidBody)

	motionData := motion.NewVmdMotion("")
	rf0 := motion.NewRigidBodyFrame(0)
	position0 := newTestVec3(0, 0, 0)
	rf0.Position = &position0
	motionData.AppendRigidBodyFrame("rb", rf0)
	rf10 := motion.NewRigidBodyFrame(10)
	position10 := newTestVec3(10, 0, 0)
	rf10.Position = &position10
	motionData.AppendRigidBodyFrame("rb", rf10)

	deltas := BuildPhysicsDeltas(modelData, motionData, 5)
	if deltas == nil || deltas.RigidBodies == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	rigidDelta := deltas.RigidBodies.Get(0)
	if rigidDelta == nil {
		t.Fatalf("rigid body delta should not be nil")
	}
	expectedPosition := newTestVec3(5, 0, 0)
	if !rigidDelta.Position.NearEquals(expectedPosition, 1e-10) {
		t.Fatalf("position interpolation mismatch: got=%v want=%v", rigidDelta.Position, expectedPosition)
	}
}

// TestBuildPhysicsDeltasConvertJointRotationLimitDegreesToRadians はジョイント回転制限を度からラジアンへ変換することを確認する。
func TestBuildPhysicsDeltasConvertJointRotationLimitDegreesToRadians(t *testing.T) {
	modelData := model.NewPmxModel()
	joint := &model.Joint{
		Param: model.JointParam{
			TranslationLimitMin: newTestVec3(-1, -2, -3),
			TranslationLimitMax: newTestVec3(1, 2, 3),
			RotationLimitMin:    newTestVec3(-0.1, -0.2, -0.3),
			RotationLimitMax:    newTestVec3(0.1, 0.2, 0.3),
		},
	}
	joint.SetIndex(0)
	joint.SetName("j")
	modelData.Joints.AppendRaw(joint)

	motionData := motion.NewVmdMotion("")
	jf := motion.NewJointFrame(0)
	translationMin := newTestVec3(-10, -20, -30)
	translationMax := newTestVec3(10, 20, 30)
	rotationMinDeg := newTestVec3(-90, -45, -30)
	rotationMaxDeg := newTestVec3(90, 45, 30)
	jf.TranslationLimitMin = &translationMin
	jf.TranslationLimitMax = &translationMax
	jf.RotationLimitMin = &rotationMinDeg
	jf.RotationLimitMax = &rotationMaxDeg
	motionData.AppendJointFrame("j", jf)

	deltas := BuildPhysicsDeltas(modelData, motionData, 0)
	if deltas == nil || deltas.Joints == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	jointDelta := deltas.Joints.Get(0)
	if jointDelta == nil {
		t.Fatalf("joint delta should not be nil")
	}
	if !jointDelta.TranslationLimitMin.NearEquals(translationMin, 1e-10) {
		t.Fatalf("translation min mismatch: got=%v want=%v", jointDelta.TranslationLimitMin, translationMin)
	}
	if !jointDelta.TranslationLimitMax.NearEquals(translationMax, 1e-10) {
		t.Fatalf("translation max mismatch: got=%v want=%v", jointDelta.TranslationLimitMax, translationMax)
	}
	expectedRotationMin := rotationMinDeg.DegToRad()
	expectedRotationMax := rotationMaxDeg.DegToRad()
	if !jointDelta.RotationLimitMin.NearEquals(expectedRotationMin, 1e-10) {
		t.Fatalf("rotation min mismatch: got=%v want=%v", jointDelta.RotationLimitMin, expectedRotationMin)
	}
	if !jointDelta.RotationLimitMax.NearEquals(expectedRotationMax, 1e-10) {
		t.Fatalf("rotation max mismatch: got=%v want=%v", jointDelta.RotationLimitMax, expectedRotationMax)
	}
}

// TestBuildPhysicsDeltasUseFallbackJointLimitsWhenInterpolatedJointFrameIsUnset は補間中未指定のジョイント制限値で既定値を使うことを確認する。
func TestBuildPhysicsDeltasUseFallbackJointLimitsWhenInterpolatedJointFrameIsUnset(t *testing.T) {
	modelData := model.NewPmxModel()
	joint := &model.Joint{
		Param: model.JointParam{
			TranslationLimitMin:       newTestVec3(-1, -2, -3),
			TranslationLimitMax:       newTestVec3(1, 2, 3),
			RotationLimitMin:          newTestVec3(-0.4, -0.5, -0.6),
			RotationLimitMax:          newTestVec3(0.4, 0.5, 0.6),
			SpringConstantTranslation: newTestVec3(7, 8, 9),
			SpringConstantRotation:    newTestVec3(10, 11, 12),
		},
	}
	joint.SetIndex(0)
	joint.SetName("j")
	modelData.Joints.AppendRaw(joint)

	motionData := motion.NewVmdMotion("")
	jf0 := motion.NewJointFrame(0)
	spring0 := newTestVec3(1, 2, 3)
	jf0.SpringConstantTranslation = &spring0
	motionData.AppendJointFrame("j", jf0)
	jf10 := motion.NewJointFrame(10)
	spring10 := newTestVec3(3, 4, 5)
	jf10.SpringConstantTranslation = &spring10
	motionData.AppendJointFrame("j", jf10)

	deltas := BuildPhysicsDeltas(modelData, motionData, 5)
	if deltas == nil || deltas.Joints == nil {
		t.Fatalf("physics deltas should not be nil")
	}
	jointDelta := deltas.Joints.Get(0)
	if jointDelta == nil {
		t.Fatalf("joint delta should not be nil")
	}
	if !jointDelta.TranslationLimitMin.NearEquals(joint.Param.TranslationLimitMin, 1e-10) {
		t.Fatalf("translation min fallback mismatch: got=%v want=%v", jointDelta.TranslationLimitMin, joint.Param.TranslationLimitMin)
	}
	if !jointDelta.TranslationLimitMax.NearEquals(joint.Param.TranslationLimitMax, 1e-10) {
		t.Fatalf("translation max fallback mismatch: got=%v want=%v", jointDelta.TranslationLimitMax, joint.Param.TranslationLimitMax)
	}
	if !jointDelta.RotationLimitMin.NearEquals(joint.Param.RotationLimitMin, 1e-10) {
		t.Fatalf("rotation min fallback mismatch: got=%v want=%v", jointDelta.RotationLimitMin, joint.Param.RotationLimitMin)
	}
	if !jointDelta.RotationLimitMax.NearEquals(joint.Param.RotationLimitMax, 1e-10) {
		t.Fatalf("rotation max fallback mismatch: got=%v want=%v", jointDelta.RotationLimitMax, joint.Param.RotationLimitMax)
	}
	expectedSpringTranslation := newTestVec3(2, 3, 4)
	if !jointDelta.SpringConstantTranslation.NearEquals(expectedSpringTranslation, 1e-10) {
		t.Fatalf("spring translation interpolation mismatch: got=%v want=%v", jointDelta.SpringConstantTranslation, expectedSpringTranslation)
	}
	if !jointDelta.SpringConstantRotation.NearEquals(joint.Param.SpringConstantRotation, 1e-10) {
		t.Fatalf("spring rotation fallback mismatch: got=%v want=%v", jointDelta.SpringConstantRotation, joint.Param.SpringConstantRotation)
	}
}

// newTestVec3 はテスト用3次元ベクトルを生成する。
func newTestVec3(x, y, z float64) mmath.Vec3 {
	v := mmath.NewVec3()
	v.X = x
	v.Y = y
	v.Z = z
	return v
}
