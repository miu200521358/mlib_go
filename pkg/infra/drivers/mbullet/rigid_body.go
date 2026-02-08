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

const (
	// boneLessReferenceResolveScoreThreshold は参照補正を試行するジョイント整合スコア閾値。
	boneLessReferenceResolveScoreThreshold = 5.0
	// boneLessPoseMovedTranslationThreshold は姿勢移動検知に用いる平行移動閾値。
	boneLessPoseMovedTranslationThreshold = 1e-3
	// boneLessPoseMovedRotationEpsilon は姿勢移動検知に用いる回転比較許容値。
	boneLessPoseMovedRotationEpsilon = 1e-3
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
		btRigidBodyTransform := newBulletTransform(rigidBody.Rotation, rigidBody.Position)
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
				appliedPosition := mp.resolveAppliedPosition(rigidBody, rigidBodyDelta)
				btRigidBodyTransform := newBulletTransform(rigidBody.Rotation, appliedPosition)
				rawScore, hasRawScore := mp.scoreBoneLessRigidBodyPositionByJoint(
					pmxModel,
					rigidBody.Index(),
					appliedPosition,
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
						} else {
							bone0Transform := mp.resolveBoneLessRigidBodyTransformByBone0(pmxModel, boneDeltas, rigidBody)
							if bone0Transform != nil {
								btRigidBodyTransform = bone0Transform
							}
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

		var rigidBodyDelta *delta.RigidBodyDelta
		if rigidBodyDeltas != nil {
			rigidBodyDelta = rigidBodyDeltas.Get(rigidBody.Index())
		}
		appliedPosition := mp.resolveAppliedPosition(rigidBody, rigidBodyDelta)

		btRigidBodyTransform := bt.NewBtTransform()
		boneTransform := bt.NewBtTransform()
		defer bt.DeleteBtTransform(boneTransform)

		mat := newMglMat4FromMat4(boneDeltas.Get(bone.Index()).FilledGlobalMatrix())
		boneTransform.SetFromOpenGLMatrix(&mat[0])

		rigidBodyLocalPos := appliedPosition.Subed(bone.Position)
		btRigidBodyLocalTransform := newBulletTransform(rigidBody.Rotation, rigidBodyLocalPos)
		defer bt.DeleteBtTransform(btRigidBodyLocalTransform)

		btRigidBodyTransform.Mult(boneTransform, btRigidBodyLocalTransform)

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
	refLocalTransform := newBulletTransform(refRigidBody.Rotation, refLocalPos)
	defer bt.DeleteBtTransform(refLocalTransform)

	refCurrentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(refCurrentTransform)
	refCurrentTransform.Mult(refBoneTransform, refLocalTransform)

	// 参照剛体のレスト変換と対象剛体のレスト変換から相対変換を算出する。
	refRestTransform := newBulletTransform(refRigidBody.Rotation, refRigidBody.Position)
	defer bt.DeleteBtTransform(refRestTransform)

	targetRestPosition := mp.resolveBoneLessRigidBodyRestPosition(pmxModel, rigidBody)
	targetRestTransform := newBulletTransform(rigidBody.Rotation, targetRestPosition)
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

// resolveBoneLessRigidBodyRestPosition はボーン未紐付け剛体のレスト位置を優先順で解釈して返す。
func (mp *PhysicsEngine) resolveBoneLessRigidBodyRestPosition(
	pmxModel *model.PmxModel,
	rigidBody *model.RigidBody,
) mmath.Vec3 {
	if rigidBody == nil {
		return mmath.ZERO_VEC3
	}
	rawPosition := rigidBody.Position
	if pmxModel == nil || pmxModel.Bones == nil {
		return rawPosition
	}
	rawScore, hasRawScore := mp.scoreBoneLessRigidBodyPositionByJoint(
		pmxModel,
		rigidBody.Index(),
		rawPosition,
	)

	// 第1候補: センターボーンからの相対位置を絶対位置に変換する。
	if centerBone, err := pmxModel.Bones.GetCenter(); err == nil && centerBone != nil {
		centerCandidate := rawPosition.Added(centerBone.Position)
		centerScore, hasCenterScore := mp.scoreBoneLessRigidBodyPositionByJoint(
			pmxModel,
			rigidBody.Index(),
			centerCandidate,
		)
		if shouldAdoptBoneLessRestCandidate(hasRawScore, rawScore, hasCenterScore, centerScore) {
			return centerCandidate
		}
	}
	// 第2候補: bone_index=0 からの相対位置を絶対位置に変換する。
	if bone0, err := pmxModel.Bones.Get(0); err == nil && bone0 != nil {
		bone0Candidate := rawPosition.Added(bone0.Position)
		bone0Score, hasBone0Score := mp.scoreBoneLessRigidBodyPositionByJoint(
			pmxModel,
			rigidBody.Index(),
			bone0Candidate,
		)
		if shouldAdoptBoneLessRestCandidate(hasRawScore, rawScore, hasBone0Score, bone0Score) {
			return bone0Candidate
		}
	}
	// 最終候補: 変換できない場合は raw のまま使う。
	return rawPosition
}

// shouldAdoptBoneLessRestCandidate は未紐付け剛体の位置候補を採用すべきか判定する。
func shouldAdoptBoneLessRestCandidate(
	hasRawScore bool,
	rawScore float64,
	hasCandidateScore bool,
	candidateScore float64,
) bool {
	if !hasCandidateScore {
		return false
	}
	if !hasRawScore {
		return true
	}
	return candidateScore <= rawScore+boneLessReferenceScoreEpsilon
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

	restAbsPos := mp.resolveBoneLessRigidBodyRestPosition(pmxModel, rigidBody)
	localPos := restAbsPos.Subed(centerBone.Position)
	localTransform := newBulletTransform(rigidBody.Rotation, localPos)
	defer bt.DeleteBtTransform(localTransform)

	targetTransform := bt.NewBtTransform()
	targetTransform.Mult(centerTransform, localTransform)
	return targetTransform
}

// resolveBoneLessRigidBodyTransformByBone0 は bone_index=0 姿勢を使って剛体初期変換を推定する。
func (mp *PhysicsEngine) resolveBoneLessRigidBodyTransformByBone0(
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBody *model.RigidBody,
) bt.BtTransform {
	if pmxModel == nil || pmxModel.Bones == nil || boneDeltas == nil || rigidBody == nil {
		return nil
	}
	bone0, err := pmxModel.Bones.Get(0)
	if err != nil || bone0 == nil {
		return nil
	}
	if !boneDeltas.Contains(bone0.Index()) {
		return nil
	}
	bone0Delta := boneDeltas.Get(bone0.Index())
	if bone0Delta == nil {
		return nil
	}

	bone0Transform := bt.NewBtTransform()
	mat := newMglMat4FromMat4(bone0Delta.FilledGlobalMatrix())
	bone0Transform.SetFromOpenGLMatrix(&mat[0])
	defer bt.DeleteBtTransform(bone0Transform)

	restAbsPos := mp.resolveBoneLessRigidBodyRestPosition(pmxModel, rigidBody)
	localPos := restAbsPos.Subed(bone0.Position)
	localTransform := newBulletTransform(rigidBody.Rotation, localPos)
	defer bt.DeleteBtTransform(localTransform)

	targetTransform := bt.NewBtTransform()
	targetTransform.Mult(bone0Transform, localTransform)
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
	modelIndex int,
	pmxModel *model.PmxModel,
	boneDeltas *delta.BoneDeltas,
	rigidBodyIndex int,
) *model.RigidBody {
	if pmxModel == nil || pmxModel.Joints == nil || pmxModel.RigidBodies == nil || pmxModel.Bones == nil || boneDeltas == nil {
		return nil
	}
	const boneLessReferenceSearchMaxDepth = 3

	type rigidBodyJointConnection struct {
		JointIndex     int
		RigidBodyIndex int
	}
	type referenceSearchState struct {
		Depth           int
		FirstJointIndex int
	}

	adjacency := make(map[int][]rigidBodyJointConnection)
	targetJointPositions := make([]mmath.Vec3, 0)
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
		candidates = append(candidates, referenceRigidBodyCandidate{
			Depth:          searchState.Depth,
			JointIndex:     searchState.FirstJointIndex,
			RigidBodyIndex: candidateRigidBodyIndex,
			SidePenalty: calculateBoneLessReferenceSidePenalty(
				targetRigidBody.Position.X,
				refRigidBody.Position.X,
			),
			JointScore: scoreBoneLessReferenceByJointPositions(targetJointPositions, refRigidBody.Position),
			Distance:   refRigidBody.Position.Distance(targetRigidBody.Position),
		})
		candidateBodies[candidateRigidBodyIndex] = refRigidBody
	}
	// 優先順位は「左右整合」->「優先深さ」->「ジョイント整合度」->「距離」->「深さ」の順で固定する。
	selected, _, ok := selectReferenceRigidBodyCandidate(candidates)
	if !ok {
		return nil
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
	if btRigidBodyTransform != nil {
		defer bt.DeleteBtTransform(btRigidBodyTransform)
	}
	appliedPosition := mp.resolveAppliedPosition(rigidBody, rigidBodyDelta)
	btCollisionShape := mp.createCollisionShape(rigidBody, rigidBodyDelta)
	mass, localInertia := mp.calculateMassAndInertia(rigidBody, btCollisionShape, rigidBodyDelta)
	defer bt.DeleteBtVector3(localInertia)

	bonePos := mp.getBonePosition(bones, rigidBody)
	btRigidBodyLocalTransform := newBulletTransform(
		rigidBody.Rotation,
		appliedPosition.Subed(bonePos),
	)

	motionState := bt.NewBtDefaultMotionState(btRigidBodyTransform)
	btRigidBody := bt.NewBtRigidBody(mass, motionState, btCollisionShape, localInertia)
	mp.configureRigidBody(btRigidBody, modelIndex, rigidBody)

	group := resolveBulletCollisionGroup(rigidBody.CollisionGroup.Group)
	mask := resolveBulletCollisionMask(rigidBody.CollisionGroup.Mask)
	appliedSize, appliedMass := mp.resolveAppliedShapeMass(rigidBody, rigidBodyDelta)
	mp.world.AddRigidBody(btRigidBody, group, mask)
	if rigidBody.Index() < 0 || rigidBody.Index() >= len(mp.rigidBodies[modelIndex]) {
		mp.world.RemoveRigidBody(btRigidBody)
		bt.DeleteBtMotionState(motionState)
		bt.DeleteBtCollisionShape(btCollisionShape)
		bt.DeleteBtRigidBody(btRigidBody)
		return
	}
	mp.rigidBodies[modelIndex][rigidBody.Index()] = &RigidBodyValue{
		RigidBody:        rigidBody,
		BtRigidBody:      btRigidBody,
		BtLocalTransform: &btRigidBodyLocalTransform,
		Mask:             mask,
		Group:            group,
		AppliedPosition:  appliedPosition,
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
	size := mp.resolveAppliedSize(rigidBody, rigidBodyDelta)

	switch rigidBody.Shape {
	case model.SHAPE_SPHERE:
		return bt.NewBtSphereShape(float32(size.X))
	case model.SHAPE_BOX:
		sizeVec := bt.NewBtVector3(float32(size.X), float32(size.Y), float32(size.Z))
		defer bt.DeleteBtVector3(sizeVec)
		return bt.NewBtBoxShape(sizeVec)
	case model.SHAPE_CAPSULE:
		return bt.NewBtCapsuleShape(float32(size.X), float32(size.Y))
	default:
		return bt.NewBtSphereShape(float32(size.X))
	}
}

// resolveAppliedPosition は剛体に適用すべき位置を返します。
func (mp *PhysicsEngine) resolveAppliedPosition(
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) mmath.Vec3 {
	if rigidBody == nil {
		return mmath.ZERO_VEC3
	}
	if rigidBodyDelta != nil {
		return rigidBodyDelta.Position
	}
	return rigidBody.Position
}

// resolveAppliedSize は剛体に適用すべきサイズを正規化して返します。
func (mp *PhysicsEngine) resolveAppliedSize(
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) mmath.Vec3 {
	if rigidBody == nil {
		return mmath.ZERO_VEC3
	}
	size := rigidBody.Size
	if rigidBodyDelta != nil {
		size = rigidBodyDelta.Size
	}
	return normalizeRigidBodySize(size)
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
	size := mp.resolveAppliedSize(rigidBody, rigidBodyDelta)
	mass := rigidBody.Param.Mass
	if rigidBodyDelta != nil {
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
		motionStateAny := r.BtRigidBody.GetMotionState()
		if motionStateAny != nil {
			if motionState, ok := motionStateAny.(bt.BtMotionState); ok && motionState != nil {
				bt.DeleteBtMotionState(motionState)
			}
		}
		shapeAny := r.BtRigidBody.GetCollisionShape()
		if shapeAny != nil {
			if shape, ok := shapeAny.(bt.BtCollisionShape); ok && shape != nil {
				bt.DeleteBtCollisionShape(shape)
			}
		}
		if r.BtLocalTransform != nil {
			bt.DeleteBtTransform(*r.BtLocalTransform)
			r.BtLocalTransform = nil
		}
		bt.DeleteBtRigidBody(r.BtRigidBody)
		r.BtRigidBody = nil
	}
	mp.rigidBodies[modelIndex] = nil
}

// updateFlag は剛体の物理フラグを更新します。
func (mp *PhysicsEngine) updateFlag(modelIndex int, rigidBody *model.RigidBody) {
	body := mp.getRigidBodyValue(modelIndex, rigidBody)
	if body == nil {
		return
	}
	btRigidBody := body.BtRigidBody

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
	return mp.getRigidBodyValueByIndex(modelIndex, rigidBody.Index())
}

// getRigidBodyValueByIndex は物理エンジン内の剛体情報をインデックスで取得する。
func (mp *PhysicsEngine) getRigidBodyValueByIndex(modelIndex int, rigidBodyIndex int) *RigidBodyValue {
	bodies, ok := mp.rigidBodies[modelIndex]
	if !ok || bodies == nil {
		return nil
	}
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
	if rigidBody.PhysicsType == model.PHYSICS_TYPE_STATIC {
		// ボーン未紐付け剛体が静的剛体を参照するケースでは、静的剛体のワールド姿勢が
		// Bullet 本体側にも反映されていないと拘束が崩れやすい。静的剛体のみ明示更新する。
		btRigidBody.SetWorldTransform(worldTransform)
		btRigidBody.SetInterpolationWorldTransform(worldTransform)
	}
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

// offsetTransformByPositionDelta は MMD 座標差分ぶん Bullet 変換の平行移動を更新します。
func (mp *PhysicsEngine) offsetTransformByPositionDelta(
	transform bt.BtTransform,
	positionDelta mmath.Vec3,
) {
	if transform == nil || positionDelta.IsZero() {
		return
	}

	originAny := transform.GetOrigin()
	origin, ok := originAny.(bt.BtVector3)
	if !ok || origin == nil {
		return
	}

	currentPosition := mmath.NewVec3()
	currentPosition.X = -float64(origin.GetX())
	currentPosition.Y = float64(origin.GetY())
	currentPosition.Z = float64(origin.GetZ())
	nextPosition := currentPosition.Added(positionDelta)
	nextOrigin := newBulletFromVec(nextPosition)
	defer bt.DeleteBtVector3(nextOrigin)
	transform.SetOrigin(nextOrigin)
}

// UpdateRigidBodyShapeMass は位置・サイズ・質量変更時に剛体状態を更新します。
func (mp *PhysicsEngine) UpdateRigidBodyShapeMass(
	modelIndex int,
	rigidBody *model.RigidBody,
	rigidBodyDelta *delta.RigidBodyDelta,
) {
	if rigidBody == nil || rigidBodyDelta == nil {
		return
	}

	r := mp.getRigidBodyValue(modelIndex, rigidBody)
	if r == nil || r.BtRigidBody == nil {
		return
	}
	nextPosition := mp.resolveAppliedPosition(rigidBody, rigidBodyDelta)
	nextSize, nextMass := mp.resolveAppliedShapeMass(rigidBody, rigidBodyDelta)
	if r.HasAppliedParams &&
		nextPosition.NearEquals(r.AppliedPosition, 1e-10) &&
		nextSize.NearEquals(r.AppliedSize, 1e-10) &&
		mmath.NearEquals(nextMass, r.AppliedMass, 1e-10) {
		return
	}
	positionChanged := !r.HasAppliedParams || !nextPosition.NearEquals(r.AppliedPosition, 1e-10)
	sizeChanged := !r.HasAppliedParams || !nextSize.NearEquals(r.AppliedSize, 1e-10)

	btRigidBody := r.BtRigidBody
	motionState := btRigidBody.GetMotionState().(bt.BtMotionState)
	currentTransform := bt.NewBtTransform()
	defer bt.DeleteBtTransform(currentTransform)
	motionState.GetWorldTransform(currentTransform)
	if positionChanged {
		positionDelta := nextPosition.Subed(r.AppliedPosition)
		mp.offsetTransformByPositionDelta(currentTransform, positionDelta)
		if r.BtLocalTransform != nil {
			mp.offsetTransformByPositionDelta(*r.BtLocalTransform, positionDelta)
		}
	}

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
	defer bt.DeleteBtVector3(newInertia)
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
	r.AppliedPosition = nextPosition
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

	modelIndex = hitObj.GetUserIndex()
	rigidBodyIndex = hitObj.GetUserIndex2()
	if mp.getRigidBodyValueByIndex(modelIndex, rigidBodyIndex) == nil {
		return -1, -1, false
	}

	return modelIndex, rigidBodyIndex, true
}
