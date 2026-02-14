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
	boneLessPoseMoved := mp.isBoneLessPoseMoved(pmxModel, boneDeltas)
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
				btRigidBodyTransform, _ := mp.resolveBoneLessRigidBodyTransform(
					pmxModel,
					boneDeltas,
					rigidBody,
				)
				var referenceRigidBody *model.RigidBody
				rawScore, hasRawScore := mp.scoreBoneLessRigidBodyPositionByJoint(
					pmxModel,
					rigidBody.Index(),
					rigidBody.Position,
				)
				shouldResolve, _ := shouldResolveBoneLessByScore(
					boneLessPoseMoved,
					hasRawScore,
					rawScore,
					boneLessReferenceResolveScoreThreshold,
				)
				if shouldResolve {
					resolvedTransform, _ := mp.resolveBoneLessRigidBodyTransform(
						modelIndex,
						pmxModel,
						boneDeltas,
						rigidBody,
					)
					if resolvedTransform != nil {
						btRigidBodyTransform = resolvedTransform
					} else if boneLessPoseMoved {
						centerTransform := mp.resolveBoneLessRigidBodyTransformByCenter(pmxModel, boneDeltas, rigidBody)
						if centerTransform != nil {
							btRigidBodyTransform = centerTransform
						}
					}
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
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBody *model.RigidBody,
) (bt.BtTransform, *model.RigidBody) {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.Bones == nil || boneDeltas == nil || rigidBody == nil {
		return nil, nil
	}
	refRigidBody := mp.findReferenceRigidBody(pmxModel, boneDeltas, rigidBody.Index())
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

	targetRestPosition := mp.resolveBoneLessRigidBodyRestPosition(pmxModel, rigidBody, refBone)
	targetRestTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation),
		newBulletFromVec(targetRestPosition),
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

// resolveBoneLessRigidBodyRestPosition はボーン未紐付け剛体のレスト位置候補を評価して返す。
func (mp *PhysicsEngine) resolveBoneLessRigidBodyRestPosition(
	pmxModel *model.PmxModel,
	rigidBody *model.RigidBody,
	refBone *model.Bone,
) mmath.Vec3 {
	if rigidBody == nil {
		return mmath.ZERO_VEC3
	}
	if pmxModel == nil || pmxModel.Joints == nil || refBone == nil {
		return rigidBody.Position
	}
	rawPosition := rigidBody.Position
	relativeToRefBonePosition := rawPosition.Added(refBone.Position)

	rawScore, hasRawScore := mp.scoreBoneLessRigidBodyPositionByJoint(
		pmxModel,
		rigidBody.Index(),
		rawPosition,
	)
	relativeScore, hasRelativeScore := mp.scoreBoneLessRigidBodyPositionByJoint(
		pmxModel,
		rigidBody.Index(),
		relativeToRefBonePosition,
	)
	if !hasRawScore && !hasRelativeScore {
		return rawPosition
	}
	if !hasRawScore {
		return relativeToRefBonePosition
	}
	if !hasRelativeScore {
		return rawPosition
	}
	scoreDelta := rawScore - relativeScore
	scoreRatio := 1.0
	if rawScore > 1e-6 {
		scoreRatio = relativeScore / rawScore
	}
	if scoreDelta > boneLessRelativePositionScoreDeltaThreshold &&
		scoreRatio < boneLessRelativePositionScoreRatioThreshold {
		return relativeToRefBonePosition
	}
	return rawPosition
}

// scoreBoneLessRigidBodyPositionByJoint は候補位置と接続ジョイント位置の一致度を評価する。
func (mp *PhysicsEngine) scoreBoneLessRigidBodyPositionByJoint(
	pmxModel *model.PmxModel,
	rigidBodyIndex int,
	candidatePosition mmath.Vec3,
) (float64, bool) {
	if pmxModel == nil || pmxModel.Joints == nil {
		return 0, false
	}
	totalDistance := 0.0
	count := 0
	for _, joint := range pmxModel.Joints.Values() {
		if joint == nil {
			continue
		}
		if joint.RigidBodyIndexA != rigidBodyIndex && joint.RigidBodyIndexB != rigidBodyIndex {
			continue
		}
		totalDistance += candidatePosition.Distance(joint.Param.Position)
		count++
	}
	if count == 0 {
		return 0, false
	}
	return totalDistance / float64(count), true
}

// resolveBoneLessRigidBodyTransformByCenter はセンターボーン姿勢を使って剛体初期変換を推定する。
func (mp *PhysicsEngine) resolveBoneLessRigidBodyTransformByCenter(
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBody *model.RigidBody,
) bt.BtTransform {
	if pmxModel == nil || pmxModel.Bones == nil || boneDeltas == nil || rigidBody == nil {
		return nil
	}
	centerBone, err := pmxModel.Bones.GetCenter()
	if err != nil || centerBone == nil {
		return nil
	}
	if !boneDeltas.Contains(centerBone.Index()) {
		return nil
	}
	centerDelta := boneDeltas.Get(centerBone.Index())
	if centerDelta == nil {
		return nil
	}

	centerTransform := bt.NewBtTransform()
	mat := newMglMat4FromMat4(centerDelta.FilledGlobalMatrix())
	centerTransform.SetFromOpenGLMatrix(&mat[0])
	defer bt.DeleteBtTransform(centerTransform)

	localPos := rigidBody.Position.Subed(centerBone.Position)
	localTransform := bt.NewBtTransform(
		newBulletFromRad(rigidBody.Rotation),
		newBulletFromVec(localPos),
	)
	defer bt.DeleteBtTransform(localTransform)

	targetTransform := bt.NewBtTransform()
	targetTransform.Mult(centerTransform, localTransform)
	return targetTransform
}

// isBoneLessPoseMoved はボーン未紐付け剛体を姿勢追従させるべきか判定する。
func (mp *PhysicsEngine) isBoneLessPoseMoved(
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
) bool {
	if pmxModel == nil || pmxModel.Bones == nil || boneDeltas == nil {
		return false
	}
	centerBone, err := pmxModel.Bones.GetCenter()
	if err != nil || centerBone == nil {
		return false
	}
	if !boneDeltas.Contains(centerBone.Index()) {
		return false
	}
	centerDelta := boneDeltas.Get(centerBone.Index())
	if centerDelta == nil {
		return false
	}
	currentPos := centerDelta.FilledGlobalPosition()
	if currentPos.Distance(centerBone.Position) > boneLessPoseMovedTranslationThreshold {
		return true
	}
	currentRot := centerDelta.FilledGlobalMatrix().Quaternion()
	return !currentRot.NearEquals(mmath.NewQuaternion(), boneLessPoseMovedRotationEpsilon)
}

// shouldResolveBoneLessRigidBodyTransform は参照剛体による補正を適用すべきか判定する。
func (mp *PhysicsEngine) shouldResolveBoneLessRigidBodyTransform(
	pmxModel *model.PmxModel,
	rigidBody *model.RigidBody,
	poseMoved bool,
) bool {
	if rigidBody == nil {
		return false
	}
	rawScore, hasRawScore := mp.scoreBoneLessRigidBodyPositionByJoint(
		pmxModel,
		rigidBody.Index(),
		rigidBody.Position,
	)
	resolve, _ := shouldResolveBoneLessByScore(
		poseMoved,
		hasRawScore,
		rawScore,
		boneLessReferenceResolveScoreThreshold,
	)
	return resolve
}

// findReferenceRigidBody はボーン未紐付け剛体の参照剛体をジョイント連結から探索する。
func (mp *PhysicsEngine) findReferenceRigidBody(
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBodyIndex int,
) *model.RigidBody {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.RigidBodies == nil || pmxModel.Bones == nil || boneDeltas == nil {
		return nil
	}
	for _, joint := range pmxModel.Joints.Values() {
		if joint == nil {
			continue
		}
		if joint.RigidBodyIndexA >= 0 && joint.RigidBodyIndexB >= 0 {
			adjacency[joint.RigidBodyIndexA] = append(adjacency[joint.RigidBodyIndexA], rigidBodyJointConnection{
				JointIndex:     joint.Index(),
				RigidBodyIndex: joint.RigidBodyIndexB,
			})
			adjacency[joint.RigidBodyIndexB] = append(adjacency[joint.RigidBodyIndexB], rigidBodyJointConnection{
				JointIndex:     joint.Index(),
				RigidBodyIndex: joint.RigidBodyIndexA,
			})
			if joint.RigidBodyIndexA == rigidBodyIndex || joint.RigidBodyIndexB == rigidBodyIndex {
				targetJointPositions = append(targetJointPositions, joint.Param.Position)
			}
		}
	}
	targetRigidBody, err := pmxModel.RigidBodies.Get(rigidBodyIndex)
	if err != nil || targetRigidBody == nil {
		return nil
	}

	searchStates := make(map[int]referenceSearchState)
	searchStates[rigidBodyIndex] = referenceSearchState{
		Depth:           0,
		FirstJointIndex: -1,
	}
	queue := []int{rigidBodyIndex}
	for len(queue) > 0 {
		currentRigidBodyIndex := queue[0]
		queue = queue[1:]
		currentState := searchStates[currentRigidBodyIndex]
		if currentState.Depth >= boneLessReferenceSearchMaxDepth {
			continue
		}

		for _, connection := range adjacency[currentRigidBodyIndex] {
			nextRigidBodyIndex := connection.RigidBodyIndex
			if nextRigidBodyIndex < 0 {
				continue
			}
			nextDepth := currentState.Depth + 1
			nextFirstJointIndex := currentState.FirstJointIndex
			if currentState.Depth == 0 {
				nextFirstJointIndex = connection.JointIndex
			}
			existingState, exists := searchStates[nextRigidBodyIndex]
			if exists {
				if existingState.Depth < nextDepth {
					continue
				}
				if existingState.Depth == nextDepth && existingState.FirstJointIndex <= nextFirstJointIndex {
					continue
				}
			}
			searchStates[nextRigidBodyIndex] = referenceSearchState{
				Depth:           nextDepth,
				FirstJointIndex: nextFirstJointIndex,
			}
			queue = append(queue, nextRigidBodyIndex)
		}
	}

	candidates := make([]referenceRigidBodyCandidate, 0)
	candidateBodies := make(map[int]*model.RigidBody)
	for candidateRigidBodyIndex, searchState := range searchStates {
		if candidateRigidBodyIndex == rigidBodyIndex || searchState.Depth <= 0 {
			continue
		}
		refRigidBody, err := pmxModel.RigidBodies.Get(candidateRigidBodyIndex)
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
		return refRigidBody
	}
	return nil
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
	appliedSize, appliedMass := mp.resolveAppliedShapeMass(rigidBody, rigidBodyDelta)
	mp.configureRigidBody(btRigidBody, modelIndex, rigidBody, appliedSize, appliedMass)

	group := 1 << int(rigidBody.CollisionGroup.Group)
	mask := int(rigidBody.CollisionGroup.Mask)
	mp.world.AddRigidBody(btRigidBody, group, mask)
	mp.rigidBodies[modelIndex][rigidBody.Index()] = &RigidBodyValue{
		RigidBody:        rigidBody,
		BtRigidBody:      btRigidBody,
		BtLocalTransform: &btRigidBodyLocalTransform,
		Mask:             mask,
		Group:            group,
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
func (mp *PhysicsEngine) configureRigidBody(
	btRigidBody bt.BtRigidBody,
	modelIndex int,
	rigidBody *model.RigidBody,
	appliedSize mmath.Vec3,
	appliedMass float64,
) {
	btRigidBody.SetDamping(float32(rigidBody.Param.LinearDamping), float32(rigidBody.Param.AngularDamping))
	btRigidBody.SetRestitution(float32(rigidBody.Param.Restitution))
	btRigidBody.SetFriction(float32(rigidBody.Param.Friction))
	ccdMotionThreshold, ccdSweptSphereRadius, ccdEnabled := resolveCcdParameters(
		rigidBody.Shape,
		appliedSize,
		appliedMass,
	)
	if ccdEnabled {
		btRigidBody.SetCcdMotionThreshold(ccdMotionThreshold)
		btRigidBody.SetCcdSweptSphereRadius(ccdSweptSphereRadius)
	}
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

	boneTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneTransform)

	mat := newMglMat4FromMat4(*boneGlobalMatrix)
	boneTransform.SetFromOpenGLMatrix(&mat[0])

	btRigidBody := mp.rigidBodies[modelIndex][rigidBody.Index()].BtRigidBody
	btRigidBodyLocalTransform := *mp.rigidBodies[modelIndex][rigidBody.Index()].BtLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)

	currentTransform.Mult(boneTransform, btRigidBodyLocalTransform)
	// 物理シミュレーションに反映させるため、剛体とモーションステートの両方を更新する。
	btRigidBody.SetWorldTransform(currentTransform)
	btRigidBody.SetInterpolationWorldTransform(currentTransform)
	motionState.SetWorldTransform(currentTransform)
	btRigidBody.Activate(true)
}

// GetRigidBodyBoneMatrix は剛体に基づいてボーン行列を取得します。
func (mp *PhysicsEngine) GetRigidBodyBoneMatrix(
	modelIndex int,
	rigidBody *model.RigidBody,
) *mmath.Mat4 {
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
	btRigidBody := body.BtRigidBody
	btRigidBodyLocalTransform := *body.BtLocalTransform

	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)

	rigidBodyGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(rigidBodyGlobalTransform)

	motionState.GetWorldTransform(rigidBodyGlobalTransform)

	boneGlobalTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(boneGlobalTransform)

	invRigidBodyLocalTransform := btRigidBodyLocalTransform.Inverse()
	defer bt.DeleteBtTransform(invRigidBodyLocalTransform)

	boneGlobalTransform.Mult(rigidBodyGlobalTransform, invRigidBodyLocalTransform)

	boneGlobalMatrixGL := mgl32.Mat4{}
	boneGlobalTransform.GetOpenGLMatrix(&boneGlobalMatrixGL[0])

	boneGlobalMatrix := newMat4FromMgl(&boneGlobalMatrixGL)
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

	massChanged := !mmath.NearEquals(rigidBodyDelta.Mass, rigidBody.Param.Mass, 1e-10)
	sizeChanged := !rigidBodyDelta.Size.NearEquals(rigidBody.Size, 1e-10)
	if !massChanged && !sizeChanged {
		return
	}

	r := mp.rigidBodies[modelIndex][rigidBody.Index()]
	if r == nil || r.BtRigidBody == nil {
		return
	}

	btRigidBody := r.BtRigidBody
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)
	motionState.GetWorldTransform(currentTransform)

	currentLinearVel := btRigidBody.GetLinearVelocity()
	currentAngularVel := btRigidBody.GetAngularVelocity()
	currentMass := float32(0.0)
	if btRigidBody.GetInvMass() != 0 {
		currentMass = 1.0 / btRigidBody.GetInvMass()
	}

	newMass := currentMass
	if rigidBody.PhysicsType != model.PHYSICS_TYPE_STATIC {
		newMass = float32(mmath.Clamped(rigidBodyDelta.Mass, 0, math.MaxFloat64))
	}

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
