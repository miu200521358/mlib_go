// 指示: miu200521358
package mphysics

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// BuildPhysicsDeltas は剛体/ジョイント差分を生成する。
func BuildPhysicsDeltas(modelData *model.PmxModel, motionData *motion.VmdMotion, frame motion.Frame) *delta.PhysicsDeltas {
	if modelData == nil {
		return nil
	}
	logger := logging.DefaultLogger()
	logSummary := logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) && shouldEmitPhysicsDeltaSummary(frame)
	var rigidDeltaCount int
	var jointDeltaCount int
	var rigidMassExplicitCount int
	var rigidMassFallbackCount int

	motionHash := motionHash(motionData)
	deltas := delta.NewPhysicsDeltas(frame, modelData.RigidBodies, modelData.Joints, modelData.Hash(), motionHash)
	if motionData == nil {
		if logSummary {
			logger.Verbose(
				logging.VERBOSE_INDEX_PHYSICS,
				"物理検証デルタ要約: frame=%v rigidDelta=%d jointDelta=%d mass(explicit=%d fallback=%d)",
				frame,
				rigidDeltaCount,
				jointDeltaCount,
				rigidMassExplicitCount,
				rigidMassFallbackCount,
			)
		}
		return deltas
	}
	if modelData.RigidBodies != nil && motionData.RigidBodyFrames != nil {
		for _, rigidBody := range modelData.RigidBodies.Values() {
			if rigidBody == nil {
				continue
			}
			frames := motionData.RigidBodyFrames.Get(rigidBody.Name())
			if frames == nil || frames.Len() == 0 {
				continue
			}
			rf := frames.Get(frame)
			if rf == nil {
				continue
			}
			size := resolveVec3(rf.Size, rigidBody.Size)
			mass := rf.Mass
			deltas.RigidBodies.Update(delta.NewRigidBodyDeltaByValue(rigidBody, frame, size, mass))
			rigidDeltaCount++
		}
	}
	if modelData.Joints != nil && motionData.JointFrames != nil {
		for _, joint := range modelData.Joints.Values() {
			if joint == nil {
				continue
			}
			frames := motionData.JointFrames.Get(joint.Name())
			if frames == nil || frames.Len() == 0 {
				continue
			}
			jf := frames.Get(frame)
			if jf == nil {
				continue
			}
			deltas.Joints.Update(delta.NewJointDeltaByValue(
				joint,
				frame,
				resolveVec3(jf.TranslationLimitMin, joint.Param.TranslationLimitMin),
				resolveVec3(jf.TranslationLimitMax, joint.Param.TranslationLimitMax),
				resolveVec3(jf.RotationLimitMin, joint.Param.RotationLimitMin),
				resolveVec3(jf.RotationLimitMax, joint.Param.RotationLimitMax),
				resolveVec3(jf.SpringConstantTranslation, joint.Param.SpringConstantTranslation),
				resolveVec3(jf.SpringConstantRotation, joint.Param.SpringConstantRotation),
			))
			jointDeltaCount++
		}
	}
	if logSummary {
		logger.Verbose(
			logging.VERBOSE_INDEX_PHYSICS,
			"物理検証デルタ要約: frame=%v rigidDelta=%d jointDelta=%d mass(explicit=%d fallback=%d)",
			frame,
			rigidDeltaCount,
			jointDeltaCount,
			rigidMassExplicitCount,
			rigidMassFallbackCount,
		)
	}
	return deltas
}

// motionHash はモーションのハッシュ文字列を返す。
func motionHash(motionData *motion.VmdMotion) string {
	if motionData == nil {
		return ""
	}
	return motionData.Hash()
}

// resolveVec3 はnilを既定値で補う。
func resolveVec3(value *mmath.Vec3, fallback mmath.Vec3) mmath.Vec3 {
	if value == nil {
		return fallback
	}
	return *value
}

// shouldEmitPhysicsDeltaSummary は物理デルタ要約を出力すべきフレームか判定する。
func shouldEmitPhysicsDeltaSummary(frame motion.Frame) bool {
	rounded := motion.Frame(math.Round(float64(frame)))
	if math.Abs(float64(frame-rounded)) > 1e-3 {
		return false
	}
	frameNumber := int(rounded)
	if frameNumber < 0 {
		return false
	}
	return frameNumber == 0 || frameNumber%30 == 0
}
