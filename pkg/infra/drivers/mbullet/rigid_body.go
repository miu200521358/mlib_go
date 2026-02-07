//go:build windows
// +build windows

// 指示: miu200521358
package mbullet

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet/bt"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// initRigidBodies はモデルの剛体を初期化します。
func (mp *PhysicsEngine) initRigidBodies(modelIndex int, pmxModel *model.PmxModel) {
	if pmxModel == nil || pmxModel.RigidBodies == nil {
		return
	}

	rigidBodies := pmxModel.RigidBodies.Values()
	mp.rigidBodies[modelIndex] = make([]*RigidBodyValue, len(rigidBodies))
	for _, rigidBody := range rigidBodies {
		if rigidBody == nil {
			continue
		}
		btRigidBodyTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBody.Position))
		mp.initRigidBody(modelIndex, pmxModel.Bones, rigidBody, btRigidBodyTransform, nil)
	}
}

// initRigidBodiesByBoneDeltas はボーンデルタ情報を使用して剛体を初期化します。
func (mp *PhysicsEngine) initRigidBodiesByBoneDeltas(
	modelIndex int,
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBodyDeltas *delta.RigidBodyDeltas,
) {
	if pmxModel == nil || pmxModel.RigidBodies == nil || pmxModel.Bones == nil || boneDeltas == nil {
		return
	}

	rigidBodies := pmxModel.RigidBodies.Values()
	mp.rigidBodies[modelIndex] = make([]*RigidBodyValue, len(rigidBodies))
	for _, rigidBody := range rigidBodies {
		if rigidBody == nil {
			continue
		}
		bone := mp.getRigidBodyBone(pmxModel.Bones, rigidBody)
		if bone == nil {
			if rigidBody.BoneIndex < 0 {
				var rigidBodyDelta *delta.RigidBodyDelta
				if rigidBodyDeltas != nil {
					rigidBodyDelta = rigidBodyDeltas.Get(rigidBody.Index())
				}
				btRigidBodyTransform, referenceRigidBody := mp.resolveBoneLessRigidBodyTransform(
					modelIndex,
					pmxModel,
					boneDeltas,
					rigidBody,
				)
				logger := logging.DefaultLogger()
				if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
					if referenceRigidBody != nil {
						logger.Verbose(
							logging.VERBOSE_INDEX_PHYSICS,
							"物理検証ボーン未紐付け剛体: model=%d rigid=%d(%s) reference=%d(%s)",
							modelIndex,
							rigidBody.Index(),
							rigidBody.Name(),
							referenceRigidBody.Index(),
							referenceRigidBody.Name(),
						)
					} else {
						logger.Verbose(
							logging.VERBOSE_INDEX_PHYSICS,
							"物理検証ボーン未紐付け剛体: model=%d rigid=%d(%s) reference=none(rest使用)",
							modelIndex,
							rigidBody.Index(),
							rigidBody.Name(),
						)
					}
				}
				if btRigidBodyTransform == nil {
					btRigidBodyTransform = bt.NewBtTransform(
						newBulletFromRad(rigidBody.Rotation),
						newBulletFromVec(rigidBody.Position),
					)
				}
				mp.initRigidBody(modelIndex, pmxModel.Bones, rigidBody, btRigidBodyTransform, rigidBodyDelta)
				continue
			}
			continue
		}
		if !boneDeltas.Contains(bone.Index()) {
			continue
		}

		btRigidBodyTransform := bt.NewBtTransform()
		boneTransform := bt.NewBtTransform()
		defer bt.DeleteBtTransform(boneTransform)

		mat := newMglMat4FromMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
		boneTransform.SetFromOpenGLMatrix(&mat[0])

		rigidBodyLocalPos := rigidBody.Position.Subed(bone.Position)
		btRigidBodyLocalTransform := bt.NewBtTransform(newBulletFromRad(rigidBody.Rotation), newBulletFromVec(rigidBodyLocalPos))
		defer bt.DeleteBtTransform(btRigidBodyLocalTransform)

		btRigidBodyTransform.Mult(boneTransform, btRigidBodyLocalTransform)

		var rigidBodyDelta *delta.RigidBodyDelta
		if rigidBodyDeltas != nil {
			rigidBodyDelta = rigidBodyDeltas.Get(rigidBody.Index())
		}

		mp.initRigidBody(modelIndex, pmxModel.Bones, rigidBody, btRigidBodyTransform, rigidBodyDelta)
	}
}

// resolveBoneLessRigidBodyTransform はボーン未紐付け剛体の初期変換を参照剛体から推定する。
func (mp *PhysicsEngine) resolveBoneLessRigidBodyTransform(
	modelIndex int,
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBody *model.RigidBody,
) (bt.BtTransform, *model.RigidBody) {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.Bones == nil || boneDeltas == nil || rigidBody == nil {
		return nil, nil
	}
	refRigidBody := mp.findReferenceRigidBody(modelIndex, pmxModel, boneDeltas, rigidBody.Index())
	if refRigidBody == nil {
		return nil, nil
	}
	refBone := mp.getRigidBodyBone(pmxModel.Bones, refRigidBody)
	if refBone == nil || !boneDeltas.Contains(refBone.Index()) {
		return nil, nil
	}
	refBoneDelta := boneDeltas.Get(refBone.Index())
	if refBoneDelta == nil {
		return nil, nil
	}

	// 参照剛体の現在ワールド変換を算出する。
	refBoneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(refBoneTransform)
	mat := newMglMat4FromMat4(refBoneDelta.FilledGlobalMatrix())
	refBoneTransform.SetFromOpenGLMatrix(&mat[0])

	refLocalPos := refRigidBody.Position.Subed(refBone.Position)
	refLocalTransform := bt.NewBtTransform(
		newBulletFromRad(refRigidBody.Rotation),
		newBulletFromVec(refLocalPos),
	)
	defer bt.DeleteBtTransform(refLocalTransform)

	refCurrentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(refCurrentTransform)
	refCurrentTransform.Mult(refBoneTransform, refLocalTransform)

	// 参照剛体のレスト変換と対象剛体のレスト変換から相対変換を算出する。
	refRestTransform := bt.NewBtTransform(
		newBulletFromRad(refRigidBody.Rotation),
		newBulletFromVec(refRigidBody.Position),
	)
	defer bt.DeleteBtTransform(refRestTransform)

	targetRestTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation),
		newBulletFromVec(rigidBody.Position),
	)
	defer bt.DeleteBtTransform(targetRestTransform)

	invRefRestTransform := refRestTransform.Inverse()
	defer bt.DeleteBtTransform(invRefRestTransform)

	relativeTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(relativeTransform)
	relativeTransform.Mult(invRefRestTransform, targetRestTransform)

	targetCurrentTransform := bt.NewBtTransform()
	targetCurrentTransform.Mult(refCurrentTransform, relativeTransform)
	return targetCurrentTransform, refRigidBody
}

// findReferenceRigidBody はボーン未紐付け剛体の参照剛体をジョイントから探索する。
func (mp *PhysicsEngine) findReferenceRigidBody(
	modelIndex int,
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBodyIndex int,
) *model.RigidBody {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.RigidBodies == nil || pmxModel.Bones == nil || boneDeltas == nil {
		return nil
	}

	candidates := make([]referenceRigidBodyCandidate, 0)
	candidateBodies := make(map[int]*model.RigidBody)
	for _, joint := range pmxModel.Joints.Values() {
		if joint == nil {
			continue
		}
		otherIndex := -1
		if joint.RigidBodyIndexA == rigidBodyIndex {
			otherIndex = joint.RigidBodyIndexB
		} else if joint.RigidBodyIndexB == rigidBodyIndex {
			otherIndex = joint.RigidBodyIndexA
		}
		if otherIndex < 0 {
			continue
		}
		refRigidBody, err := pmxModel.RigidBodies.Get(otherIndex)
		if err != nil || refRigidBody == nil {
			continue
		}
		if refRigidBody.BoneIndex < 0 {
			continue
		}
		refBone, err := pmxModel.Bones.Get(refRigidBody.BoneIndex)
		if err != nil || refBone == nil {
			continue
		}
		if !boneDeltas.Contains(refBone.Index()) {
			continue
		}

		candidates = append(candidates, referenceRigidBodyCandidate{
			JointIndex:     joint.Index(),
			RigidBodyIndex: otherIndex,
		})
		candidateBodies[otherIndex] = refRigidBody
	}
	// 優先順位は「ジョイントindex昇順」->「剛体index昇順」で固定し、候補が複数でも決定的に選ぶ。
	selected, sortedCandidates, ok := selectReferenceRigidBodyCandidate(candidates)
	if !ok {
		return nil
	}
	if len(sortedCandidates) > 1 {
		candidateIndexes := make([]int, 0, len(sortedCandidates))
		for _, candidate := range sortedCandidates {
			candidateIndexes = append(candidateIndexes, candidate.RigidBodyIndex)
		}
		logger := logging.DefaultLogger()
		if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
			logger.Verbose(
				logging.VERBOSE_INDEX_PHYSICS,
				"物理検証参照候補複数: model=%d target=%d selected=%d candidates=%v",
				modelIndex,
				rigidBodyIndex,
				selected.RigidBodyIndex,
				candidateIndexes,
			)
		}
	}

	return candidateBodies[selected.RigidBodyIndex]
}

// initRigidBody は個別の剛体を初期化します。
func (mp *PhysicsEngine) initRigidBody(
	modelIndex int,
	bones *model.BoneCollection,
	rigidBody *model.RigidBody,
	btRigidBodyTransform bt.BtTransform,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	btCollisionShape := mp.createCollisionShape(rigidBody, rigidBodyDelta)
	mass, localInertia := mp.calculateMassAndInertia(rigidBody, btCollisionShape, rigidBodyDelta)

	bonePos := mp.getBonePosition(bones, rigidBody)
	btRigidBodyLocalTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation),
		newBulletFromVec(rigidBody.Position.Subed(bonePos)),
	)

	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)
	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	mp.configureRigidBody(btRigidBody, modelIndex, rigidBody)

	group := 1 << int(rigidBody.CollisionGroup.Group)
	mask := int(rigidBody.CollisionGroup.Mask)
	appliedSize, appliedMass := mp.resolveAppliedShapeMass(rigidBody, rigidBodyDelta)
	mp.world.AddRigidBody(btRigidBody, group, mask)
	mp.rigidBodies[modelIndex][rigidBody.Index()] = &RigidBodyValue{
		RigidBody:        rigidBody,
		BtRigidBody:      btRigidBody,
		BtLocalTransform: &btRigidBodyLocalTransform,
		Mask:             mask,
		Group:            group,
		AppliedSize:      appliedSize,
		AppliedMass:      appliedMass,
		HasAppliedParams: true,
	}

	mp.updateFlag(modelIndex, rigidBody)
}

// createCollisionShape は剛体の形状に基づいた衝突形状を生成します。
func (mp *PhysicsEngine) createCollisionShape(
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) bt.BtCollisionShape {
	size := rigidBody.Size
	if rigidBodyDelta != nil {
		size = rigidBodyDelta.Size
	}
	size.Clamp(mmath.ZERO_VEC3, mmath.VEC3_MAX_VAL)

	switch rigidBody.Shape {
	case model.SHAPE_SPHERE:
		return bt.NewBtSphereShape(float32(size.X))
	case model.SHAPE_BOX:
		return bt.NewBtBoxShape(bt.NewBtVector3(float32(size.X), float32(size.Y), float32(size.Z)))
	case model.SHAPE_CAPSULE:
		return bt.NewBtCapsuleShape(float32(size.X), float32(size.Y))
	default:
		return bt.NewBtSphereShape(float32(size.X))
	}
}

// calculateMassAndInertia は剛体の質量と慣性を計算します。
func (mp *PhysicsEngine) calculateMassAndInertia(
	rigidBody *model.RigidBody,
	btCollisionShape bt.BtCollisionShape,
	rigidBodyDelta *delta.RigidBodyDelta,
) (float32, bt.BtVector3) {
	mass := float32(0.0)
	localInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))

	if rigidBody.PhysicsType != model.PHYSICS_TYPE_STATIC {
		baseMass := rigidBody.Param.Mass
		if rigidBodyDelta != nil {
			baseMass = rigidBodyDelta.Mass
		}
		mass = float32(mmath.Clamped(baseMass, 0, math.MaxFloat64))
	}

	if mass != 0 {
		btCollisionShape.CalculateLocalInertia(mass, localInertia)
	}

	return mass, localInertia
}

// resolveAppliedShapeMass は剛体に適用すべきサイズと質量を返します。
func (mp *PhysicsEngine) resolveAppliedShapeMass(
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) (mmath.Vec3, float64) {
	size := rigidBody.Size
	mass := rigidBody.Param.Mass
	if rigidBodyDelta != nil {
		size = rigidBodyDelta.Size
		mass = rigidBodyDelta.Mass
	}
	if rigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
		mass = 0
	}
	return size, mass
}

// getRigidBodyBone は剛体に紐づくボーンを取得します。
func (mp *PhysicsEngine) getRigidBodyBone(bones *model.BoneCollection, rigidBody *model.RigidBody) *model.Bone {
	if bones == nil || rigidBody == nil || rigidBody.BoneIndex < 0 {
		return nil
	}
	bone, err := bones.Get(rigidBody.BoneIndex)
	if err != nil {
		return nil
	}
	return bone
}

// getBonePosition は剛体に紐づくボーン位置を取得します。
func (mp *PhysicsEngine) getBonePosition(bones *model.BoneCollection, rigidBody *model.RigidBody) mmath.Vec3 {
	bone := mp.getRigidBodyBone(bones, rigidBody)
	if bone == nil {
		return mmath.NewVec3()
	}
	return bone.Position
}

// configureRigidBody は剛体の物理パラメータを設定します。
func (mp *PhysicsEngine) configureRigidBody(btRigidBody bt.BtRigidBody, modelIndex int, rigidBody *model.RigidBody) {
	btRigidBody.SetDamping(float32(rigidBody.Param.LinearDamping), float32(rigidBody.Param.AngularDamping))
	btRigidBody.SetRestitution(float32(rigidBody.Param.Restitution))
	btRigidBody.SetFriction(float32(rigidBody.Param.Friction))
	btRigidBody.SetUserIndex(modelIndex)
	btRigidBody.SetUserIndex2(rigidBody.Index())
}

// deleteRigidBodies はモデルの全剛体を削除します。
func (mp *PhysicsEngine) deleteRigidBodies(modelIndex int) {
	for _, r := range mp.rigidBodies[modelIndex] {
		if r == nil || r.BtRigidBody == nil {
			continue
		}
		mp.world.RemoveRigidBody(r.BtRigidBody)
		bt.DeleteBtRigidBody(r.BtRigidBody)
	}
	mp.rigidBodies[modelIndex] = nil
}

// updateFlag は剛体の物理フラグを更新します。
func (mp *PhysicsEngine) updateFlag(modelIndex int, rigidBody *model.RigidBody) {
	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody

	if rigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
		btRigidBody.SetCollisionFlags(
			btRigidBody.GetCollisionFlags() | int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT),
		)
		btRigidBody.SetActivationState(bt.DISABLE_SIMULATION)
		return
	}

	btRigidBody.SetCollisionFlags(
		btRigidBody.GetCollisionFlags() &
			^int(bt.BtCollisionObjectCF_NO_CONTACT_RESPONSE) &
			^int(bt.BtCollisionObjectCF_KINEMATIC_OBJECT),
	)
	btRigidBody.SetActivationState(bt.DISABLE_DEACTIVATION)
}

// getRigidBodyValue は物理エンジン内の剛体情報を取得する。
func (mp *PhysicsEngine) getRigidBodyValue(modelIndex int, rigidBody *model.RigidBody) *RigidBodyValue {
	if rigidBody == nil {
		return nil
	}
	bodies, ok := mp.rigidBodies[modelIndex]
	if !ok || bodies == nil {
		return nil
	}
	rigidBodyIndex := rigidBody.Index()
	if rigidBodyIndex < 0 || rigidBodyIndex >= len(bodies) {
		return nil
	}
	body := bodies[rigidBodyIndex]
	if body == nil || body.BtRigidBody == nil || body.BtLocalTransform == nil {
		return nil
	}
	return body
}

// UpdateTransform はボーン行列に基づいて剛体の位置を更新します。
func (mp *PhysicsEngine) UpdateTransform(
	modelIndex int,
	rigidBodyBone *model.Bone,
	boneGlobalMatrix *mmath.Mat4,
	rigidBody *model.RigidBody,
) {
	if rigidBodyBone == nil || boneGlobalMatrix == nil || rigidBody == nil {
		return
	}
	body := mp.getRigidBodyValue(modelIndex, rigidBody)
	if body == nil {
		return
	}

	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	mat := newMglMat4FromMat4(*boneGlobalMatrix)
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	btRigidBody := body.BtRigidBody
	btRigidBodyLocalTransform := *body.BtLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	worldTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(worldTransform)

	worldTransform.Mult(boneTransform, btRigidBodyLocalTransform)
	// 旧 mlib と同等に、剛体追従は MotionState 更新を基準とする。
	// 動的剛体へ直接 SetWorldTransform を繰り返すとソルバの warm start を崩して不安定化しやすいため、
	// ここでは MotionState の更新に限定する。
	motionState.SetWorldTransform(worldTransform)
	body.PrevBoneMatrix = *boneGlobalMatrix
	body.HasPrevBone = true
}

// FollowDeltaTransform は前回ボーン姿勢との差分で剛体姿勢を追従更新する。
func (mp *PhysicsEngine) FollowDeltaTransform(
	modelIndex int,
	rigidBodyBone *model.Bone,
	boneGlobalMatrix *mmath.Mat4,
	rigidBody *model.RigidBody,
) {
	if rigidBodyBone == nil || boneGlobalMatrix == nil || rigidBody == nil {
		return
	}
	body := mp.getRigidBodyValue(modelIndex, rigidBody)
	if body == nil {
		return
	}
	if !body.HasPrevBone {
		// 初回は差分を計算できないため、ハード同期を行って基準姿勢を確定する。
		mp.UpdateTransform(modelIndex, rigidBodyBone, boneGlobalMatrix, rigidBody)
		return
	}
	// 同一フレーム停止中の微小誤差を追従すると速度にノイズが蓄積するため、差分が極小なら更新しない。
	if boneGlobalMatrix.NearEquals(body.PrevBoneMatrix, 1e-7) {
		body.PrevBoneMatrix = *boneGlobalMatrix
		return
	}
	btRigidBody := body.BtRigidBody
	btRigidBodyLocalTransform := *body.BtLocalTransform
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)
	motionState.GetWorldTransform(currentTransform)

	// 物理後反映済みの同一姿勢を再追従すると、停止中に差分が二重適用されて発散する。
	// 入力ボーン姿勢が現在剛体姿勢由来と一致する場合は追従処理を行わない。
	currentRigidBodyBoneMatrix := mp.getBoneMatrixByTransforms(currentTransform, btRigidBodyLocalTransform)
	if boneGlobalMatrix.NearEquals(currentRigidBodyBoneMatrix, 1e-6) {
		body.PrevBoneMatrix = *boneGlobalMatrix
		body.HasPrevBone = true
		return
	}

	prevBoneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(prevBoneTransform)
	prevMat := newMglMat4FromMat4(body.PrevBoneMatrix)
	prevBoneTransform.SetFromOpenGLMatrix(&prevMat[0])

	currBoneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currBoneTransform)
	currMat := newMglMat4FromMat4(*boneGlobalMatrix)
	currBoneTransform.SetFromOpenGLMatrix(&currMat[0])

	invPrevBoneTransform := prevBoneTransform.Inverse()
	defer bt.DeleteBtTransform(invPrevBoneTransform)

	deltaBoneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(deltaBoneTransform)
	deltaBoneTransform.Mult(currBoneTransform, invPrevBoneTransform)

	targetTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(targetTransform)
	targetTransform.Mult(deltaBoneTransform, currentTransform)

	// 剛体姿勢とモーションステートを同時更新して、ボーン差分に追従させる。
	btRigidBody.SetWorldTransform(targetTransform)
	btRigidBody.SetInterpolationWorldTransform(targetTransform)
	motionState.SetWorldTransform(targetTransform)

	// 速度ベクトルも差分回転に追従させ、見た目の連続性を維持する。
	deltaRotation := deltaBoneTransform.GetRotation()
	defer bt.DeleteBtQuaternion(deltaRotation)
	normalizedDeltaRotation := deltaRotation.Normalized()
	defer bt.DeleteBtQuaternion(normalizedDeltaRotation)

	deltaRotationAngle := resolveQuaternionRotationAngleFromW(float64(normalizedDeltaRotation.GetW()))
	if shouldRotateVelocityByDeltaRotation(deltaRotationAngle, mp.followDeltaVelocityRotationMaxRad) {
		linearVelocity := btRigidBody.GetLinearVelocity()
		angularVelocity := btRigidBody.GetAngularVelocity()
		rotatedLinearVelocity := bt.QuatRotate(normalizedDeltaRotation, linearVelocity)
		rotatedAngularVelocity := bt.QuatRotate(normalizedDeltaRotation, angularVelocity)
		defer bt.DeleteBtVector3(rotatedLinearVelocity)
		defer bt.DeleteBtVector3(rotatedAngularVelocity)

		btRigidBody.SetLinearVelocity(rotatedLinearVelocity)
		btRigidBody.SetAngularVelocity(rotatedAngularVelocity)
	}
	btRigidBody.Activate(true)

	body.PrevBoneMatrix = *boneGlobalMatrix
	body.HasPrevBone = true
}

// getBoneMatrixByTransforms は剛体ワールド変換とローカル変換からボーングローバル行列を生成する。
func (mp *PhysicsEngine) getBoneMatrixByTransforms(
	rigidBodyGlobalTransform bt.BtTransform,
	rigidBodyLocalTransform bt.BtTransform,
) mmath.Mat4 {
	boneGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneGlobalTransform)

	invRigidBodyLocalTransform := rigidBodyLocalTransform.Inverse()
	defer bt.DeleteBtTransform(invRigidBodyLocalTransform)

	boneGlobalTransform.Mult(rigidBodyGlobalTransform, invRigidBodyLocalTransform)

	boneGlobalMatrixGL := mgl32.Mat4{}
	boneGlobalTransform.GetOpenGLMatrix(&boneGlobalMatrixGL[0])

	return newMat4FromMgl(&boneGlobalMatrixGL)
}

// GetRigidBodyBoneMatrix は剛体に基づいてボーン行列を取得します。
func (mp *PhysicsEngine) GetRigidBodyBoneMatrix(
	modelIndex int,
	rigidBody *model.RigidBody,
) *mmath.Mat4 {
	body := mp.getRigidBodyValue(modelIndex, rigidBody)
	if body == nil {
		return nil
	}
	btRigidBody := body.BtRigidBody
	btRigidBodyLocalTransform := *body.BtLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	rigidBodyGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(rigidBodyGlobalTransform)

	motionState.GetWorldTransform(rigidBodyGlobalTransform)

	boneGlobalMatrix := mp.getBoneMatrixByTransforms(rigidBodyGlobalTransform, btRigidBodyLocalTransform)
	return &boneGlobalMatrix
}

// UpdateRigidBodiesSelectively は変更が必要な剛体のみを選択的に更新します。
func (mp *PhysicsEngine) UpdateRigidBodiesSelectively(
	modelIndex int,
	pmxModel *model.PmxModel,
	rigidBodyDeltas *delta.RigidBodyDeltas,
) {
	if pmxModel == nil || rigidBodyDeltas == nil {
		return
	}

	rigidBodyDeltas.ForEach(func(index int, rigidBodyDelta *delta.RigidBodyDelta) bool {
		if rigidBodyDelta == nil || rigidBodyDelta.RigidBody == nil {
			return true
		}

		mp.UpdateRigidBodyShapeMass(modelIndex, rigidBodyDelta.RigidBody, rigidBodyDelta)

		return true
	})
}

// UpdateRigidBodyShapeMass はサイズ・質量変更時に剛体の形状を更新します。
func (mp *PhysicsEngine) UpdateRigidBodyShapeMass(
	modelIndex int,
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	if rigidBody == nil || rigidBodyDelta == nil {
		return
	}

	r := mp.rigidBodies[modelIndex][rigidBody.Index()]
	if r == nil || r.BtRigidBody == nil {
		return
	}
	nextSize, nextMass := mp.resolveAppliedShapeMass(rigidBody, rigidBodyDelta)
	if r.HasAppliedParams &&
		nextSize.NearEquals(r.AppliedSize, 1e-10) &&
		mmath.NearEquals(nextMass, r.AppliedMass, 1e-10) {
		return
	}
	sizeChanged := !r.HasAppliedParams || !nextSize.NearEquals(r.AppliedSize, 1e-10)

	btRigidBody := r.BtRigidBody
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)
	motionState.GetWorldTransform(currentTransform)

	currentLinearVel := btRigidBody.GetLinearVelocity()
	currentAngularVel := btRigidBody.GetAngularVelocity()
	newMass := float32(mmath.Clamped(nextMass, 0, math.MaxFloat64))

	needShapeUpdate := sizeChanged
	if needShapeUpdate {
		mp.world.RemoveRigidBody(btRigidBody)

		oldShape := btRigidBody.GetCollisionShape()
		if oldShape != nil {
			if collisionShape, ok := oldShape.(bt.BtCollisionShape); ok {
				bt.DeleteBtCollisionShape(collisionShape)
			}
		}

		newShape := mp.createCollisionShape(rigidBody, rigidBodyDelta)
		btRigidBody.SetCollisionShape(newShape)
	}

	newInertia := bt.NewBtVector3(float32(0.0), float32(0.0), float32(0.0))
	if newMass > 0 && rigidBody.PhysicsType != model.PHYSICS_TYPE_STATIC {
		currentShape := btRigidBody.GetCollisionShape()
		if currentShape != nil {
			if collisionShape, ok := currentShape.(bt.BtCollisionShape); ok {
				collisionShape.CalculateLocalInertia(newMass, newInertia)
			}
		}
	}

	btRigidBody.SetMassProps(newMass, newInertia)
	btRigidBody.UpdateInertiaTensor()

	motionState.SetWorldTransform(currentTransform)
	btRigidBody.SetLinearVelocity(currentLinearVel)
	btRigidBody.SetAngularVelocity(currentAngularVel)

	if needShapeUpdate {
		mp.world.AddRigidBody(btRigidBody, r.Group, r.Mask)
	}

	mp.updateFlag(modelIndex, rigidBody)
	btRigidBody.Activate(true)
	r.AppliedSize = nextSize
	r.AppliedMass = nextMass
	r.HasAppliedParams = true
}

// findRigidBodyByCollisionObject は衝突オブジェクトから剛体参照を取得します。
func (mp *PhysicsEngine) findRigidBodyByCollisionObject(
	hitObj bt.BtCollisionObject,
	hasHit bool,
) (modelIndex int, rigidBodyIndex int, ok bool) {
	if hitObj == nil || !hasHit {
		return -1, -1, false
	}
	defer func() {
		if r := recover(); r != nil {
			modelIndex = -1
			rigidBodyIndex = -1
			ok = false
		}
	}()

	modelIndex = hitObj.GetUserIndex()
	rigidBodyIndex = hitObj.GetUserIndex2()

	bodies, exists := mp.rigidBodies[modelIndex]
	if !exists {
		return -1, -1, false
	}
	if rigidBodyIndex < 0 || rigidBodyIndex >= len(bodies) {
		return -1, -1, false
	}
	if bodies[rigidBodyIndex] == nil || bodies[rigidBodyIndex].BtRigidBody == nil {
		return -1, -1, false
	}

	return modelIndex, rigidBodyIndex, true
}
