// 指示: miu200521358
package mbullet

import (
	"math"
	"sort"
)

const defaultFollowDeltaVelocityRotationMaxRadians = math.Pi / 6.0

// referenceRigidBodyCandidate はボーン未紐付け剛体の参照候補を表す。
type referenceRigidBodyCandidate struct {
	JointIndex     int
	RigidBodyIndex int
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
		if !exists || candidate.JointIndex < current.JointIndex {
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
		if left.JointIndex != right.JointIndex {
			return left.JointIndex < right.JointIndex
		}
		return left.RigidBodyIndex < right.RigidBodyIndex
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
