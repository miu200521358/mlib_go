package deform

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

// DeformPhysics 前回情報ありで物理デフォーム処理を実行する
func DeformPhysics(
	model *pmx.PmxModel,
	motion *vmd.VmdMotion,
	frame float32,
) *delta.PhysicsDeltas {
	deltas := delta.NewPhysicsDeltas(frame, model.RigidBodies, model.Joints, model.Hash(), motion.Hash())

	model.RigidBodies.ForEach(func(rigidBodyIndex int, rigidBody *pmx.RigidBody) bool {
		if rigidBody == nil {
			return true
		}

		rf := motion.RigidBodyFrames.Get(rigidBody.Name()).Get(frame)
		if rf == nil {
			return true
		}

		deltas.RigidBodies.Update(delta.NewRigidBodyDeltaByValue(rigidBody, frame, rf.Size, rf.Mass))

		return true
	})

	model.Joints.ForEach(func(jointIndex int, joint *pmx.Joint) bool {
		if joint == nil {
			return true
		}

		jf := motion.JointFrames.Get(joint.Name()).Get(frame)
		if jf == nil {
			return true
		}

		deltas.Joints.Update(delta.NewJointDeltaByValue(joint, frame, jf.TranslationLimitMin,
			jf.TranslationLimitMax, jf.RotationLimitMin, jf.RotationLimitMax,
			jf.SpringConstantTranslation, jf.SpringConstantRotation))
		return true
	})

	return deltas
}
