// 指示: miu200521358
package units

// LengthUnit は長さ単位を表す。
type LengthUnit int

const (
	// LENGTH_UNIT_MMD はMMD単位を表す。
	LENGTH_UNIT_MMD LengthUnit = iota
)

// AngleUnit は角度単位を表す。
type AngleUnit int

const (
	// ANGLE_UNIT_RADIAN はラジアンを表す。
	ANGLE_UNIT_RADIAN AngleUnit = iota
)

// CameraFovUnit はカメラFOV単位を表す。
type CameraFovUnit int

const (
	// CAMERA_FOV_UNIT_DEGREE は度を表す。
	CAMERA_FOV_UNIT_DEGREE CameraFovUnit = iota
)

// ScaleNormalization はスケール正規化方針を表す。
type ScaleNormalization int

const (
	// SCALE_NORMALIZATION_NONE は正規化しない。
	SCALE_NORMALIZATION_NONE ScaleNormalization = iota
)

// UvPolicy はUV座標系の方針を表す。
type UvPolicy int

const (
	// UV_POLICY_NO_V_FLIP はV反転なしを表す。
	UV_POLICY_NO_V_FLIP UvPolicy = iota
)

// PhysicsUnit は物理単位系を表す。
type PhysicsUnit int

const (
	// PHYSICS_UNIT_ABSTRACT は抽象単位を表す。
	PHYSICS_UNIT_ABSTRACT PhysicsUnit = iota
)

// WindUnit は風単位系を表す。
type WindUnit int

const (
	// WIND_UNIT_REAL_WORLD は実世界相当を表す。
	WIND_UNIT_REAL_WORLD WindUnit = iota
)

// UnitPolicy は単位/スケールの前提を表す。
type UnitPolicy struct {
	Length            LengthUnit
	Angle             AngleUnit
	CameraFov         CameraFovUnit
	ScaleNormalization ScaleNormalization
	Uv                UvPolicy
	Physics           PhysicsUnit
	Wind              WindUnit
}

// DEFAULT_UNIT_POLICY は単位/スケールの既定ポリシー。
var DEFAULT_UNIT_POLICY = UnitPolicy{
	Length:            LENGTH_UNIT_MMD,
	Angle:             ANGLE_UNIT_RADIAN,
	CameraFov:         CAMERA_FOV_UNIT_DEGREE,
	ScaleNormalization: SCALE_NORMALIZATION_NONE,
	Uv:                UV_POLICY_NO_V_FLIP,
	Physics:           PHYSICS_UNIT_ABSTRACT,
	Wind:              WIND_UNIT_REAL_WORLD,
}
