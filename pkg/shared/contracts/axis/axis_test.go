// 指示: miu200521358
package axis

import "testing"

// TestAxisPolicies は座標系ポリシーの既定値を確認する。
func TestAxisPolicies(t *testing.T) {
	if AXIS_POLICY_INTERNAL.Handedness != HANDEDNESS_LEFT {
		t.Errorf("AXIS_POLICY_INTERNAL.Handedness: got=%v", AXIS_POLICY_INTERNAL.Handedness)
	}
	if AXIS_POLICY_INTERNAL.Up != (AxisDir{Axis: AXIS_Y, Sign: AXIS_SIGN_POS}) {
		t.Errorf("AXIS_POLICY_INTERNAL.Up: got=%v", AXIS_POLICY_INTERNAL.Up)
	}
	if AXIS_POLICY_INTERNAL.Forward != (AxisDir{Axis: AXIS_Z, Sign: AXIS_SIGN_NEG}) {
		t.Errorf("AXIS_POLICY_INTERNAL.Forward: got=%v", AXIS_POLICY_INTERNAL.Forward)
	}
	if AXIS_POLICY_INTERNAL.Right != (AxisDir{Axis: AXIS_X, Sign: AXIS_SIGN_POS}) {
		t.Errorf("AXIS_POLICY_INTERNAL.Right: got=%v", AXIS_POLICY_INTERNAL.Right)
	}

	if AXIS_POLICY_IO_MMD != AXIS_POLICY_INTERNAL {
		t.Errorf("AXIS_POLICY_IO_MMD: got=%v want=%v", AXIS_POLICY_IO_MMD, AXIS_POLICY_INTERNAL)
	}
}

// TestMatrixPolicyDefault は行列規約の既定値を確認する。
func TestMatrixPolicyDefault(t *testing.T) {
	if MATRIX_POLICY_DEFAULT.Layout != MATRIX_LAYOUT_COLUMN_MAJOR {
		t.Errorf("MATRIX_POLICY_DEFAULT.Layout: got=%v", MATRIX_POLICY_DEFAULT.Layout)
	}
	if MATRIX_POLICY_DEFAULT.VectorMultiply != VECTOR_MULTIPLY_MV {
		t.Errorf("MATRIX_POLICY_DEFAULT.VectorMultiply: got=%v", MATRIX_POLICY_DEFAULT.VectorMultiply)
	}
	if MATRIX_POLICY_DEFAULT.GlobalCompose != COMPOSE_ORDER_PARENT_LOCAL {
		t.Errorf("MATRIX_POLICY_DEFAULT.GlobalCompose: got=%v", MATRIX_POLICY_DEFAULT.GlobalCompose)
	}
	if MATRIX_POLICY_DEFAULT.LocalTransformOrder != LOCAL_TRANSFORM_ORDER_LOCAL_SCALE_POSITION_ROTATION {
		t.Errorf("MATRIX_POLICY_DEFAULT.LocalTransformOrder: got=%v", MATRIX_POLICY_DEFAULT.LocalTransformOrder)
	}
	if MATRIX_POLICY_DEFAULT.BoneOffsetOrder != BONE_OFFSET_ORDER_REVERT_OFFSET_UNIT {
		t.Errorf("MATRIX_POLICY_DEFAULT.BoneOffsetOrder: got=%v", MATRIX_POLICY_DEFAULT.BoneOffsetOrder)
	}
}

// TestCameraNdcZRange はカメラのNDCレンジ定数を確認する。
func TestCameraNdcZRange(t *testing.T) {
	if CAMERA_NDC_Z_RANGE != NDC_Z_RANGE_NEG1_POS1 {
		t.Errorf("CAMERA_NDC_Z_RANGE: got=%v", CAMERA_NDC_Z_RANGE)
	}
}
