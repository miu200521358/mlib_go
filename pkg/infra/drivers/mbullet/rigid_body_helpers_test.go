// 指示: miu200521358
package mbullet

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestSelectReferenceRigidBodyCandidate(t *testing.T) {
	candidates := []referenceRigidBodyCandidate{
		{Depth: 1, JointIndex: 5, RigidBodyIndex: 3, SidePenalty: 0, JointScore: 2.0, Distance: 4.0},
		{Depth: 1, JointIndex: 2, RigidBodyIndex: 7, SidePenalty: 0, JointScore: 1.0, Distance: 3.0},
		{Depth: 1, JointIndex: 1, RigidBodyIndex: 4, SidePenalty: 0, JointScore: 1.0, Distance: 2.0},
		{Depth: 1, JointIndex: 9, RigidBodyIndex: 7, SidePenalty: 0, JointScore: 9.0, Distance: 9.0},
	}

	selected, normalized, ok := selectReferenceRigidBodyCandidate(candidates)
	if !ok {
		t.Fatalf("候補があるのに選択できませんでした")
	}
	if selected.RigidBodyIndex != 4 || selected.JointIndex != 1 || selected.Depth != 1 {
		t.Fatalf("想定外の選択結果: selected=%+v", selected)
	}
	if len(normalized) != 3 {
		t.Fatalf("候補の正規化件数が不正です: got=%d want=3", len(normalized))
	}
	if normalized[0].RigidBodyIndex != 4 || normalized[1].RigidBodyIndex != 7 || normalized[2].RigidBodyIndex != 3 {
		t.Fatalf("候補の優先順位が不正です: %+v", normalized)
	}
}

func TestSelectReferenceRigidBodyCandidateEmpty(t *testing.T) {
	candidates := []referenceRigidBodyCandidate{
		{Depth: 1, JointIndex: 1, RigidBodyIndex: -1},
	}

	_, normalized, ok := selectReferenceRigidBodyCandidate(candidates)
	if ok {
		t.Fatalf("無効候補のみなのに選択成功になりました")
	}
	if normalized != nil {
		t.Fatalf("無効候補のみでは正規化結果はnil想定です: %+v", normalized)
	}
}

func TestSelectReferenceRigidBodyCandidateDepthPriority(t *testing.T) {
	candidates := []referenceRigidBodyCandidate{
		{Depth: 3, JointIndex: 0, RigidBodyIndex: 1, SidePenalty: 0, JointScore: 0, Distance: 0},
		{Depth: 1, JointIndex: 9, RigidBodyIndex: 2, SidePenalty: 0, JointScore: 10, Distance: 10},
		{Depth: 1, JointIndex: 3, RigidBodyIndex: 3, SidePenalty: 0, JointScore: 9, Distance: 9},
	}

	selected, normalized, ok := selectReferenceRigidBodyCandidate(candidates)
	if !ok {
		t.Fatalf("候補があるのに選択できませんでした")
	}
	if selected.RigidBodyIndex != 3 || selected.Depth != 1 || selected.JointIndex != 3 {
		t.Fatalf("深さ優先の選択結果が不正です: selected=%+v", selected)
	}
	if len(normalized) != 3 {
		t.Fatalf("候補の正規化件数が不正です: got=%d want=3", len(normalized))
	}
}

func TestSelectReferenceRigidBodyCandidateSidePriority(t *testing.T) {
	candidates := []referenceRigidBodyCandidate{
		{Depth: 1, JointIndex: 1, RigidBodyIndex: 1, SidePenalty: 1, JointScore: 0.1, Distance: 0.1},
		{Depth: 1, JointIndex: 2, RigidBodyIndex: 2, SidePenalty: 0, JointScore: 9.0, Distance: 9.0},
	}

	selected, _, ok := selectReferenceRigidBodyCandidate(candidates)
	if !ok {
		t.Fatalf("候補があるのに選択できませんでした")
	}
	if selected.RigidBodyIndex != 2 {
		t.Fatalf("左右整合の優先順位が不正です: selected=%+v", selected)
	}
}

func TestCalculateBoneLessReferenceSidePenalty(t *testing.T) {
	if penalty := calculateBoneLessReferenceSidePenalty(1.0, -1.0); penalty != 1 {
		t.Fatalf("左右反転時のペナルティが不正です: got=%d want=1", penalty)
	}
	if penalty := calculateBoneLessReferenceSidePenalty(1.0, 0.0); penalty != 0 {
		t.Fatalf("中央参照時のペナルティが不正です: got=%d want=0", penalty)
	}
}

func TestScoreBoneLessReferenceByJointPositions(t *testing.T) {
	jointPos1 := mmath.NewVec3()
	jointPos1.X = 1.0
	jointPos2 := mmath.NewVec3()
	jointPos2.X = 3.0
	targetJointPositions := []mmath.Vec3{
		jointPos1,
		jointPos2,
	}
	candidatePos := mmath.NewVec3()
	candidatePos.X = 2.0
	score := scoreBoneLessReferenceByJointPositions(targetJointPositions, candidatePos)
	if math.Abs(score-1.0) > 1e-12 {
		t.Fatalf("ジョイント整合度が不正です: got=%f want=1.0", score)
	}
}

func TestResolveQuaternionRotationAngleFromW(t *testing.T) {
	if angle := resolveQuaternionRotationAngleFromW(1.0); math.Abs(angle-0) > 1e-12 {
		t.Fatalf("w=1 の角度が不正です: got=%f want=0", angle)
	}
	if angle := resolveQuaternionRotationAngleFromW(-1.0); math.Abs(angle-0) > 1e-12 {
		t.Fatalf("w=-1 の角度が不正です: got=%f want=0", angle)
	}
	if angle := resolveQuaternionRotationAngleFromW(0.0); math.Abs(angle-math.Pi) > 1e-12 {
		t.Fatalf("w=0 の角度が不正です: got=%f want=%f", angle, math.Pi)
	}
}

func TestClampFollowDeltaVelocityRotationMaxRadians(t *testing.T) {
	if clamped := clampFollowDeltaVelocityRotationMaxRadians(math.NaN()); math.Abs(clamped-defaultFollowDeltaVelocityRotationMaxRadians) > 1e-12 {
		t.Fatalf("NaN クランプ結果が不正です: got=%f want=%f", clamped, defaultFollowDeltaVelocityRotationMaxRadians)
	}
	if clamped := clampFollowDeltaVelocityRotationMaxRadians(math.Pi * 2); math.Abs(clamped-math.Pi) > 1e-12 {
		t.Fatalf("上限クランプ結果が不正です: got=%f want=%f", clamped, math.Pi)
	}
	if clamped := clampFollowDeltaVelocityRotationMaxRadians(0.2); math.Abs(clamped-0.2) > 1e-12 {
		t.Fatalf("正常値クランプ結果が不正です: got=%f want=0.2", clamped)
	}
}

func TestShouldRotateVelocityByDeltaRotation(t *testing.T) {
	if !shouldRotateVelocityByDeltaRotation(0.15, 0.2) {
		t.Fatalf("しきい値未満で false は不正です")
	}
	if shouldRotateVelocityByDeltaRotation(0.25, 0.2) {
		t.Fatalf("しきい値超過で true は不正です")
	}
}

func TestShouldResolveBoneLessByScorePoseMoved(t *testing.T) {
	resolve, reason := shouldResolveBoneLessByScore(true, false, 0.0, 5.0)
	if !resolve {
		t.Fatalf("poseMoved=true なのに resolve=false は不正です")
	}
	if reason != boneLessResolveReasonPoseMoved {
		t.Fatalf("理由が不正です: got=%s want=%s", reason, boneLessResolveReasonPoseMoved)
	}
}

func TestShouldResolveBoneLessByScoreNoJoint(t *testing.T) {
	resolve, reason := shouldResolveBoneLessByScore(false, false, 0.0, 5.0)
	if resolve {
		t.Fatalf("ジョイントスコア無しで resolve=true は不正です")
	}
	if reason != boneLessResolveReasonNoJointScore {
		t.Fatalf("理由が不正です: got=%s want=%s", reason, boneLessResolveReasonNoJointScore)
	}
}

func TestShouldResolveBoneLessByScoreThreshold(t *testing.T) {
	resolve, reason := shouldResolveBoneLessByScore(false, true, 6.0, 5.0)
	if !resolve {
		t.Fatalf("閾値超過で resolve=false は不正です")
	}
	if reason != boneLessResolveReasonScoreThreshold {
		t.Fatalf("理由が不正です: got=%s want=%s", reason, boneLessResolveReasonScoreThreshold)
	}
}

func TestShouldResolveBoneLessByScoreJointConnected(t *testing.T) {
	resolve, reason := shouldResolveBoneLessByScore(false, true, 1.0, 5.0)
	if !resolve {
		t.Fatalf("閾値未満でもジョイント接続ありで resolve=false は不正です")
	}
	if reason != boneLessResolveReasonJointConnected {
		t.Fatalf("理由が不正です: got=%s want=%s", reason, boneLessResolveReasonJointConnected)
	}
}
