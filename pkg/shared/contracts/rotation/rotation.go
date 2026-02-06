// 指示: miu200521358
package rotation

// QuaternionOrder はクォータニオンの成分順を表す。
type QuaternionOrder int

const (
	// QUATERNION_ORDER_XYZW はXYZW順を表す。
	QUATERNION_ORDER_XYZW QuaternionOrder = iota
)

// EulerOrder はオイラー角の回転順を表す。
type EulerOrder int

const (
	// EULER_ORDER_XYZ はXYZ順を表す。
	EULER_ORDER_XYZ EulerOrder = iota
	// EULER_ORDER_YXZ はYXZ順を表す。
	EULER_ORDER_YXZ
	// EULER_ORDER_ZXY はZXY順を表す。
	EULER_ORDER_ZXY
	// EULER_ORDER_YZX はYZX順を表す。
	EULER_ORDER_YZX
)

// IkLimitEulerOrder はIK制限用のオイラー順を表す。
type IkLimitEulerOrder int

const (
	// IK_LIMIT_EULER_ORDER_ZXY はZXY順を表す。
	IK_LIMIT_EULER_ORDER_ZXY IkLimitEulerOrder = iota
	// IK_LIMIT_EULER_ORDER_XYZ はXYZ順を表す。
	IK_LIMIT_EULER_ORDER_XYZ
	// IK_LIMIT_EULER_ORDER_YZX はYZX順を表す。
	IK_LIMIT_EULER_ORDER_YZX
)

// RotationInterpolation は回転補間方式を表す。
type RotationInterpolation int

const (
	// ROTATION_INTERPOLATION_SLERP は球面線形補間を表す。
	ROTATION_INTERPOLATION_SLERP RotationInterpolation = iota
)

// RotationPolicy は回転表現の前提を表す。
type RotationPolicy struct {
	QuaternionOrder    QuaternionOrder
	EulerOrder         EulerOrder
	IkLimitEulerOrders []IkLimitEulerOrder
	Interpolation      RotationInterpolation
}

// VMD_CAMERA_ROTATION_UNIT はVMDカメラ回転の単位を表す。
const VMD_CAMERA_ROTATION_UNIT = "degree"

// DEFAULT_ROTATION_POLICY は既定の回転ポリシー。
var DEFAULT_ROTATION_POLICY = RotationPolicy{
	QuaternionOrder:    QUATERNION_ORDER_XYZW,
	EulerOrder:         EULER_ORDER_YXZ,
	IkLimitEulerOrders: []IkLimitEulerOrder{IK_LIMIT_EULER_ORDER_ZXY, IK_LIMIT_EULER_ORDER_XYZ, IK_LIMIT_EULER_ORDER_YZX},
	Interpolation:      ROTATION_INTERPOLATION_SLERP,
}
