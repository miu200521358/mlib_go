// 指示: miu200521358
package mbullet

import (
	"math"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

const defaultFollowDeltaVelocityRotationMaxRadians = math.Pi / 6.0
const boneLessReferencePreferredDepth = 2
const boneLessReferenceSideThreshold = 0.2
const boneLessReferenceScoreEpsilon = 1e-6
const enableStaticRigidBodyWorldTransformSync = false
const minCcdSweptSphereRadius = 0.005
const minCcdMotionThreshold = 0.02
const ccdSweptSphereRadiusScale = 0.2
const ccdMotionThresholdScale = 0.5

const boneLessResolveReasonPoseMoved = "pose_moved"
const boneLessResolveReasonNoJointScore = "no_joint_score"
const boneLessResolveReasonScoreThreshold = "score_over_threshold"
const boneLessResolveReasonJointConnected = "joint_connected"

// resolveBulletCollisionGroup はドメインの衝突グループを Bullet のグループビットへ変換する。
func resolveBulletCollisionGroup(collisionGroup byte) int {
	return 1 << int(collisionGroup)
}

// resolveBulletCollisionMask は PMX 由来の衝突マスク値を Bullet へそのまま渡す。
func resolveBulletCollisionMask(collisionMask uint16) int {
	return int(collisionMask)
}

// normalizeRigidBodySize は剛体サイズを物理計算で扱える安全範囲へ正規化する。
func normalizeRigidBodySize(size mmath.Vec3) mmath.Vec3 {
	return size.Clamped(mmath.ZERO_VEC3, mmath.VEC3_MAX_VAL)
}

// resolveCcdParameters は剛体属性からCCDのしきい値と半径を解決する。
func resolveCcdParameters(
	shape model.Shape,
	size mmath.Vec3,
	appliedMass float64,
) (float32, float32, bool) {
	if appliedMass <= 0 {
		return 0, 0, false
	}
	characteristicLength := resolveCcdCharacteristicLength(shape, size)
	if characteristicLength <= 0 || math.IsNaN(characteristicLength) || math.IsInf(characteristicLength, 0) {
		return 0, 0, false
	}

	ccdRadius := math.Max(characteristicLength*ccdSweptSphereRadiusScale, minCcdSweptSphereRadius)
	ccdMotionThreshold := math.Max(characteristicLength*ccdMotionThresholdScale, minCcdMotionThreshold)
	return float32(ccdMotionThreshold), float32(ccdRadius), true
}

// resolveCcdCharacteristicLength はCCD判定に使う代表長を返す。
func resolveCcdCharacteristicLength(shape model.Shape, size mmath.Vec3) float64 {
	positiveMin := resolvePositiveMinSizeComponent(size)
	switch shape {
	case model.SHAPE_SPHERE:
		return max(size.X, 0)
	case model.SHAPE_CAPSULE:
		return max(size.X, 0)
	case model.SHAPE_BOX:
		return positiveMin
	default:
		return positiveMin
	}
}

// resolvePositiveMinSizeComponent は正のサイズ成分の最小値を返す。
func resolvePositiveMinSizeComponent(size mmath.Vec3) float64 {
	values := []float64{size.X, size.Y, size.Z}
	minimum := math.MaxFloat64
	for _, value := range values {
		if value > 0 && value < minimum {
			minimum = value
		}
	}
	if minimum == math.MaxFloat64 {
		return 0
	}
	return minimum
}

// referenceRigidBodyCandidate はボーン未紐付け剛体の参照候補を表す。
type referenceRigidBodyCandidate struct {
	Depth          int
	JointIndex     int
	RigidBodyIndex int
	SidePenalty    int
	JointScore     float64
	Distance       float64
}

// normalizeReferenceRigidBodyCandidates は候補を剛体ごとに正規化して優先順位順へ並べる。
func normalizeReferenceRigidBodyCandidates(candidates []referenceRigidBodyCandidate) []referenceRigidBodyCandidate {
	if len(candidates) == 0 {
		return nil
	}

	deduped := make(map[int]referenceRigidBodyCandidate)
	for _, candidate := range candidates {
		if candidate.RigidBodyIndex < 0 {
			continue
		}
		current, exists := deduped[candidate.RigidBodyIndex]
		if !exists || isHigherPriorityReferenceRigidBodyCandidate(candidate, current) {
			deduped[candidate.RigidBodyIndex] = candidate
		}
	}
	if len(deduped) == 0 {
		return nil
	}

	normalized := make([]referenceRigidBodyCandidate, 0, len(deduped))
	for _, candidate := range deduped {
		normalized = append(normalized, candidate)
	}
	sort.Slice(normalized, func(i, j int) bool {
		left := normalized[i]
		right := normalized[j]
		return isHigherPriorityReferenceRigidBodyCandidate(left, right)
	})
	return normalized
}

// selectReferenceRigidBodyCandidate は候補から優先順位に基づいて参照剛体を決定する。
func selectReferenceRigidBodyCandidate(
	candidates []referenceRigidBodyCandidate,
) (referenceRigidBodyCandidate, []referenceRigidBodyCandidate, bool) {
	normalized := normalizeReferenceRigidBodyCandidates(candidates)
	if len(normalized) == 0 {
		return referenceRigidBodyCandidate{}, nil, false
	}
	return normalized[0], normalized, true
}

// isHigherPriorityReferenceRigidBodyCandidate は候補優先度比較を行う。
func isHigherPriorityReferenceRigidBodyCandidate(
	left referenceRigidBodyCandidate,
	right referenceRigidBodyCandidate,
) bool {
	if left.SidePenalty != right.SidePenalty {
		return left.SidePenalty < right.SidePenalty
	}
	leftDepthPreferred := left.Depth <= boneLessReferencePreferredDepth
	rightDepthPreferred := right.Depth <= boneLessReferencePreferredDepth
	if leftDepthPreferred != rightDepthPreferred {
		return leftDepthPreferred
	}
	if left.JointScore < right.JointScore-boneLessReferenceScoreEpsilon {
		return true
	}
	if right.JointScore < left.JointScore-boneLessReferenceScoreEpsilon {
		return false
	}
	if left.Distance < right.Distance-boneLessReferenceScoreEpsilon {
		return true
	}
	if right.Distance < left.Distance-boneLessReferenceScoreEpsilon {
		return false
	}
	if left.Depth != right.Depth {
		return left.Depth < right.Depth
	}
	if left.JointIndex != right.JointIndex {
		return left.JointIndex < right.JointIndex
	}
	return left.RigidBodyIndex < right.RigidBodyIndex
}

// calculateBoneLessReferenceSidePenalty は左右反転参照を避けるためのペナルティを返す。
func calculateBoneLessReferenceSidePenalty(targetPositionX, referencePositionX float64) int {
	targetSide := rigidBodySideByPositionX(targetPositionX)
	referenceSide := rigidBodySideByPositionX(referencePositionX)
	if targetSide == 0 || referenceSide == 0 || targetSide == referenceSide {
		return 0
	}
	return 1
}

// scoreBoneLessReferenceByJointPositions は対象ジョイント近傍との整合度を評価する。
func scoreBoneLessReferenceByJointPositions(
	targetJointPositions []mmath.Vec3,
	candidatePosition mmath.Vec3,
) float64 {
	if len(targetJointPositions) == 0 {
		return 0
	}
	totalDistance := 0.0
	for _, jointPosition := range targetJointPositions {
		totalDistance += candidatePosition.Distance(jointPosition)
	}
	return totalDistance / float64(len(targetJointPositions))
}

// rigidBodySideByPositionX はX座標から左右判定値（左:+1/中央:0/右:-1）を返す。
func rigidBodySideByPositionX(positionX float64) int {
	if positionX > boneLessReferenceSideThreshold {
		return 1
	}
	if positionX < -boneLessReferenceSideThreshold {
		return -1
	}
	return 0
}

// resolveQuaternionRotationAngleFromW はクォータニオンのW成分から回転角[rad]を求める。
func resolveQuaternionRotationAngleFromW(w float64) float64 {
	clampedW := math.Max(-1, math.Min(1, w))
	return 2.0 * math.Acos(math.Abs(clampedW))
}

// shouldRotateVelocityByDeltaRotation は差分回転に速度ベクトル追従を適用するか判定する。
func shouldRotateVelocityByDeltaRotation(rotationAngleRad, maxAngleRad float64) bool {
	return rotationAngleRad <= maxAngleRad+1e-6
}

// clampFollowDeltaVelocityRotationMaxRadians は速度回転許容角度を安全な範囲へ丸める。
func clampFollowDeltaVelocityRotationMaxRadians(maxAngleRad float64) float64 {
	if math.IsNaN(maxAngleRad) || math.IsInf(maxAngleRad, 0) || maxAngleRad <= 0 {
		return defaultFollowDeltaVelocityRotationMaxRadians
	}
	if maxAngleRad > math.Pi {
		return math.Pi
	}
	return maxAngleRad
}

// shouldResolveBoneLessByScore はボーン未紐付け剛体の参照推定を試行すべきか判定する。
func shouldResolveBoneLessByScore(
	poseMoved bool,
	hasRawScore bool,
	rawScore float64,
	scoreThreshold float64,
) (bool, string) {
	if poseMoved {
		return true, boneLessResolveReasonPoseMoved
	}
	if !hasRawScore {
		return false, boneLessResolveReasonNoJointScore
	}
	if rawScore >= scoreThreshold {
		return true, boneLessResolveReasonScoreThreshold
	}
	// 閾値未満でもジョイント接続がある剛体は参照探索を試行する。
	// 採用姿勢は resolveBoneLessRigidBodyRestPosition が raw/relative を比較して決める。
	return true, boneLessResolveReasonJointConnected
}
