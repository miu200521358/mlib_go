// 指示: miu200521358
package axis

// Axis は軸方向を表す。
type Axis int

const (
	// AXIS_X はX軸を表す。
	AXIS_X Axis = iota
	// AXIS_Y はY軸を表す。
	AXIS_Y
	// AXIS_Z はZ軸を表す。
	AXIS_Z
)

// AxisSign は軸の符号を表す。
type AxisSign int

const (
	// AXIS_SIGN_POS は正方向を表す。
	AXIS_SIGN_POS AxisSign = iota
	// AXIS_SIGN_NEG は負方向を表す。
	AXIS_SIGN_NEG
)

// AxisDir は軸と符号の組を表す。
type AxisDir struct {
	Axis Axis
	Sign AxisSign
}

// Handedness は座標系の手系を表す。
type Handedness int

const (
	// HANDEDNESS_RIGHT は右手系を表す。
	HANDEDNESS_RIGHT Handedness = iota
	// HANDEDNESS_LEFT は左手系を表す。
	HANDEDNESS_LEFT
)

// AxisPolicy は座標系の前提を表す。
type AxisPolicy struct {
	Handedness Handedness
	Up         AxisDir
	Forward    AxisDir
	Right      AxisDir
}

// MatrixLayout は行列のレイアウトを表す。
type MatrixLayout int

const (
	// MATRIX_LAYOUT_COLUMN_MAJOR は列優先を表す。
	MATRIX_LAYOUT_COLUMN_MAJOR MatrixLayout = iota
)

// VectorMultiply はベクトルの乗算方向を表す。
type VectorMultiply int

const (
	// VECTOR_MULTIPLY_MV はM*vを表す。
	VECTOR_MULTIPLY_MV VectorMultiply = iota
)

// ComposeOrder はグローバル合成順を表す。
type ComposeOrder int

const (
	// COMPOSE_ORDER_PARENT_LOCAL は親→ローカル合成を表す。
	COMPOSE_ORDER_PARENT_LOCAL ComposeOrder = iota
)

// LocalTransformOrder はローカル変換順を表す。
type LocalTransformOrder int

const (
	// LOCAL_TRANSFORM_ORDER_LOCAL_SCALE_POSITION_ROTATION はS→T→R順を表す。
	LOCAL_TRANSFORM_ORDER_LOCAL_SCALE_POSITION_ROTATION LocalTransformOrder = iota
)

// BoneOffsetOrder はボーンオフセットの扱い順を表す。
type BoneOffsetOrder int

const (
	// BONE_OFFSET_ORDER_REVERT_OFFSET_UNIT はオフセット単位の反転規約を表す。
	BONE_OFFSET_ORDER_REVERT_OFFSET_UNIT BoneOffsetOrder = iota
)

// NdcZRange はNDCのZレンジを表す。
type NdcZRange int

const (
	// NDC_Z_RANGE_NEG1_POS1 は-1〜+1を表す。
	NDC_Z_RANGE_NEG1_POS1 NdcZRange = iota
)

// MatrixPolicy は行列計算の前提を表す。
type MatrixPolicy struct {
	Layout              MatrixLayout
	VectorMultiply      VectorMultiply
	GlobalCompose       ComposeOrder
	LocalTransformOrder LocalTransformOrder
	BoneOffsetOrder     BoneOffsetOrder
}

// AXIS_POLICY_INTERNAL は内部座標系の既定ポリシー。
var AXIS_POLICY_INTERNAL = AxisPolicy{
	Handedness: HANDEDNESS_LEFT,
	Up:         AxisDir{Axis: AXIS_Y, Sign: AXIS_SIGN_POS},
	Forward:    AxisDir{Axis: AXIS_Z, Sign: AXIS_SIGN_NEG},
	Right:      AxisDir{Axis: AXIS_X, Sign: AXIS_SIGN_POS},
}

// AXIS_POLICY_IO_MMD はMMD I/O向けポリシー。
var AXIS_POLICY_IO_MMD = AxisPolicy{
	Handedness: HANDEDNESS_LEFT,
	Up:         AxisDir{Axis: AXIS_Y, Sign: AXIS_SIGN_POS},
	Forward:    AxisDir{Axis: AXIS_Z, Sign: AXIS_SIGN_NEG},
	Right:      AxisDir{Axis: AXIS_X, Sign: AXIS_SIGN_POS},
}

// MATRIX_POLICY_DEFAULT は行列計算の既定ポリシー。
var MATRIX_POLICY_DEFAULT = MatrixPolicy{
	Layout:              MATRIX_LAYOUT_COLUMN_MAJOR,
	VectorMultiply:      VECTOR_MULTIPLY_MV,
	GlobalCompose:       COMPOSE_ORDER_PARENT_LOCAL,
	LocalTransformOrder: LOCAL_TRANSFORM_ORDER_LOCAL_SCALE_POSITION_ROTATION,
	BoneOffsetOrder:     BONE_OFFSET_ORDER_REVERT_OFFSET_UNIT,
}

// CAMERA_NDC_Z_RANGE はカメラのNDCレンジを表す。
var CAMERA_NDC_Z_RANGE = NDC_Z_RANGE_NEG1_POS1
