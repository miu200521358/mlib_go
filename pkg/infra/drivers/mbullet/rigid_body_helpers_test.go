// 指示: miu200521358
package mbullet

import (
	"math"
	"testing"
)

func TestSelectReferenceRigidBodyCandidate(t *testing.T) {
	candidates := []referenceRigidBodyCandidate{
		{JointIndex: 5, RigidBodyIndex: 3},
		{JointIndex: 2, RigidBodyIndex: 7},
		{JointIndex: 2, RigidBodyIndex: 4},
		{JointIndex: 1, RigidBodyIndex: 7},
	}

	selected, normalized, ok := selectReferenceRigidBodyCandidate(candidates)
	if !ok {
		t.Fatalf("候補があるのに選択できませんでした")
	}
	if selected.RigidBodyIndex != 7 || selected.JointIndex != 1 {
		t.Fatalf("想定外の選択結果: selected=%+v", selected)
	}
	if len(normalized) != 3 {
		t.Fatalf("候補の正規化件数が不正です: got=%d want=3", len(normalized))
	}
	if normalized[0].RigidBodyIndex != 7 || normalized[1].RigidBodyIndex != 4 || normalized[2].RigidBodyIndex != 3 {
		t.Fatalf("候補の優先順位が不正です: %+v", normalized)
	}
}

func TestSelectReferenceRigidBodyCandidateEmpty(t *testing.T) {
	candidates := []referenceRigidBodyCandidate{
		{JointIndex: 1, RigidBodyIndex: -1},
	}

	_, normalized, ok := selectReferenceRigidBodyCandidate(candidates)
	if ok {
		t.Fatalf("無効候補のみなのに選択成功になりました")
	}
	if normalized != nil {
		t.Fatalf("無効候補のみでは正規化結果はnil想定です: %+v", normalized)
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
