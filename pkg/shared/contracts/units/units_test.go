// 指示: miu200521358
package units

import "testing"

// TestUnitPolicyDefault は既定の単位ポリシーを確認する。
func TestUnitPolicyDefault(t *testing.T) {
	if DEFAULT_UNIT_POLICY.Length != LENGTH_UNIT_MMD {
		t.Errorf("Length: got=%v", DEFAULT_UNIT_POLICY.Length)
	}
	if DEFAULT_UNIT_POLICY.Angle != ANGLE_UNIT_RADIAN {
		t.Errorf("Angle: got=%v", DEFAULT_UNIT_POLICY.Angle)
	}
	if DEFAULT_UNIT_POLICY.CameraFov != CAMERA_FOV_UNIT_DEGREE {
		t.Errorf("CameraFov: got=%v", DEFAULT_UNIT_POLICY.CameraFov)
	}
	if DEFAULT_UNIT_POLICY.ScaleNormalization != SCALE_NORMALIZATION_NONE {
		t.Errorf("ScaleNormalization: got=%v", DEFAULT_UNIT_POLICY.ScaleNormalization)
	}
	if DEFAULT_UNIT_POLICY.Uv != UV_POLICY_NO_V_FLIP {
		t.Errorf("Uv: got=%v", DEFAULT_UNIT_POLICY.Uv)
	}
	if DEFAULT_UNIT_POLICY.Physics != PHYSICS_UNIT_ABSTRACT {
		t.Errorf("Physics: got=%v", DEFAULT_UNIT_POLICY.Physics)
	}
	if DEFAULT_UNIT_POLICY.Wind != WIND_UNIT_REAL_WORLD {
		t.Errorf("Wind: got=%v", DEFAULT_UNIT_POLICY.Wind)
	}
}
