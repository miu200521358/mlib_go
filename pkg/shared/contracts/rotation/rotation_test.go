// 指示: miu200521358
package rotation

import "testing"

// TestRotationPolicyDefault は既定の回転ポリシーを確認する。
func TestRotationPolicyDefault(t *testing.T) {
	if DEFAULT_ROTATION_POLICY.QuaternionOrder != QUATERNION_ORDER_XYZW {
		t.Errorf("QuaternionOrder: got=%v", DEFAULT_ROTATION_POLICY.QuaternionOrder)
	}
	if DEFAULT_ROTATION_POLICY.EulerOrder != EULER_ORDER_YXZ {
		t.Errorf("EulerOrder: got=%v", DEFAULT_ROTATION_POLICY.EulerOrder)
	}
	orders := DEFAULT_ROTATION_POLICY.IkLimitEulerOrders
	if len(orders) != 3 || orders[0] != IK_LIMIT_EULER_ORDER_ZXY || orders[1] != IK_LIMIT_EULER_ORDER_XYZ || orders[2] != IK_LIMIT_EULER_ORDER_YZX {
		t.Errorf("IkLimitEulerOrders: got=%v", orders)
	}
	if DEFAULT_ROTATION_POLICY.Interpolation != ROTATION_INTERPOLATION_SLERP {
		t.Errorf("Interpolation: got=%v", DEFAULT_ROTATION_POLICY.Interpolation)
	}
}

// TestVmdCameraRotationUnit はVMDカメラ回転の単位を確認する。
func TestVmdCameraRotationUnit(t *testing.T) {
	if VMD_CAMERA_ROTATION_UNIT != "degree" {
		t.Errorf("VMD_CAMERA_ROTATION_UNIT: got=%v", VMD_CAMERA_ROTATION_UNIT)
	}
}
