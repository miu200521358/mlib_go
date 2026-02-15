// 指示: miu200521358
package mphysics

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// BuildPhysicsDeltas は剛体/ジョイント差分を生成する。
func BuildPhysicsDeltas(modelData *model.PmxModel, motionData *motion.VmdMotion, frame motion.Frame) *delta.PhysicsDeltas {
	if modelData == nil {
		return nil
	}

	motionHash := motionHash(motionData)
	deltas := delta.NewPhysicsDeltas(frame, modelData.RigidBodies, modelData.Joints, modelData.Hash(), motionHash)
	if motionData == nil {
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
			position := resolveRigidBodyFrameVec3(
				frames,
				frame,
				rigidBody.Position,
				func(frameValue *motion.RigidBodyFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.Position
				},
				nil,
			)
			size := resolveRigidBodyFrameVec3(
				frames,
				frame,
				rigidBody.Size,
				func(frameValue *motion.RigidBodyFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.Size
				},
				nil,
			)
			mass := rigidBody.Param.Mass
			if rf.Mass != nil {
				mass = *rf.Mass
			}
			deltas.RigidBodies.Update(delta.NewRigidBodyDeltaByValue(rigidBody, frame, position, size, mass))
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
			translationLimitMin := resolveJointFrameVec3(
				frames,
				frame,
				joint.Param.TranslationLimitMin,
				func(frameValue *motion.JointFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.TranslationLimitMin
				},
				nil,
			)
			translationLimitMax := resolveJointFrameVec3(
				frames,
				frame,
				joint.Param.TranslationLimitMax,
				func(frameValue *motion.JointFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.TranslationLimitMax
				},
				nil,
			)
			rotationLimitMin := resolveJointFrameVec3(
				frames,
				frame,
				joint.Param.RotationLimitMin,
				func(frameValue *motion.JointFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.RotationLimitMin
				},
				func(value mmath.Vec3) mmath.Vec3 {
					return value.DegToRad()
				},
			)
			rotationLimitMax := resolveJointFrameVec3(
				frames,
				frame,
				joint.Param.RotationLimitMax,
				func(frameValue *motion.JointFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.RotationLimitMax
				},
				func(value mmath.Vec3) mmath.Vec3 {
					return value.DegToRad()
				},
			)
			springConstantTranslation := resolveJointFrameVec3(
				frames,
				frame,
				joint.Param.SpringConstantTranslation,
				func(frameValue *motion.JointFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.SpringConstantTranslation
				},
				nil,
			)
			springConstantRotation := resolveJointFrameVec3(
				frames,
				frame,
				joint.Param.SpringConstantRotation,
				func(frameValue *motion.JointFrame) *mmath.Vec3 {
					if frameValue == nil {
						return nil
					}
					return frameValue.SpringConstantRotation
				},
				nil,
			)
			deltas.Joints.Update(delta.NewJointDeltaByValue(
				joint,
				frame,
				translationLimitMin,
				translationLimitMax,
				rotationLimitMin,
				rotationLimitMax,
				springConstantTranslation,
				springConstantRotation,
			))
		}
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

// resolveRigidBodyFrameVec3 は剛体フレームのVec3値を補間を考慮して解決する。
func resolveRigidBodyFrameVec3(
	frames *motion.RigidBodyNameFrames,
	frame motion.Frame,
	fallback mmath.Vec3,
	extractor func(frameValue *motion.RigidBodyFrame) *mmath.Vec3,
	transform func(value mmath.Vec3) mmath.Vec3,
) mmath.Vec3 {
	if frames == nil || extractor == nil {
		return fallback
	}
	if frames.Has(frame) {
		frameValue := frames.Get(frame)
		value := extractor(frameValue)
		if value == nil {
			return fallback
		}
		return applyVec3Transform(*value, transform)
	}
	prevFrame, hasPrev := frames.PrevFrame(frame)
	nextFrame, hasNext := frames.NextFrame(frame)
	var prevValue *mmath.Vec3
	if hasPrev {
		prevValue = extractor(frames.Get(prevFrame))
	}
	var nextValue *mmath.Vec3
	if hasNext {
		nextValue = extractor(frames.Get(nextFrame))
	}
	return resolveFrameRangeVec3(
		frame,
		prevFrame,
		nextFrame,
		hasPrev,
		hasNext,
		prevValue,
		nextValue,
		fallback,
		transform,
	)
}

// resolveJointFrameVec3 はジョイントフレームのVec3値を補間を考慮して解決する。
func resolveJointFrameVec3(
	frames *motion.JointNameFrames,
	frame motion.Frame,
	fallback mmath.Vec3,
	extractor func(frameValue *motion.JointFrame) *mmath.Vec3,
	transform func(value mmath.Vec3) mmath.Vec3,
) mmath.Vec3 {
	if frames == nil || extractor == nil {
		return fallback
	}
	if frames.Has(frame) {
		frameValue := frames.Get(frame)
		value := extractor(frameValue)
		if value == nil {
			return fallback
		}
		return applyVec3Transform(*value, transform)
	}
	prevFrame, hasPrev := frames.PrevFrame(frame)
	nextFrame, hasNext := frames.NextFrame(frame)
	var prevValue *mmath.Vec3
	if hasPrev {
		prevValue = extractor(frames.Get(prevFrame))
	}
	var nextValue *mmath.Vec3
	if hasNext {
		nextValue = extractor(frames.Get(nextFrame))
	}
	return resolveFrameRangeVec3(
		frame,
		prevFrame,
		nextFrame,
		hasPrev,
		hasNext,
		prevValue,
		nextValue,
		fallback,
		transform,
	)
}

// resolveFrameRangeVec3 は前後フレーム値の有無に応じてVec3を補間または補完する。
func resolveFrameRangeVec3(
	frame motion.Frame,
	prevFrame motion.Frame,
	nextFrame motion.Frame,
	hasPrev bool,
	hasNext bool,
	prevValue *mmath.Vec3,
	nextValue *mmath.Vec3,
	fallback mmath.Vec3,
	transform func(value mmath.Vec3) mmath.Vec3,
) mmath.Vec3 {
	if !hasPrev {
		return fallback
	}
	if !hasNext {
		if prevValue == nil {
			return fallback
		}
		return applyVec3Transform(*prevValue, transform)
	}
	if prevValue != nil && nextValue != nil {
		if nextFrame == prevFrame {
			return applyVec3Transform(*prevValue, transform)
		}
		t := resolveFrameLerpT(prevFrame, frame, nextFrame)
		value := prevValue.Lerp(*nextValue, t)
		return applyVec3Transform(value, transform)
	}
	if prevValue != nil {
		return applyVec3Transform(*prevValue, transform)
	}
	return fallback
}

// resolveFrameLerpT はフレーム番号から線形補間係数を計算する。
func resolveFrameLerpT(prevFrame motion.Frame, frame motion.Frame, nextFrame motion.Frame) float64 {
	denom := float64(nextFrame - prevFrame)
	if denom == 0 {
		return 0
	}
	return float64(frame-prevFrame) / denom
}

// applyVec3Transform は必要時のみVec3変換を適用する。
func applyVec3Transform(value mmath.Vec3, transform func(value mmath.Vec3) mmath.Vec3) mmath.Vec3 {
	if transform == nil {
		return value
	}
	return transform(value)
}
