// 指示: miu200521358
package model

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model/merrors"
	"gonum.org/v1/gonum/spatial/r3"
)

// GetRoot すべての親
func (bones *BoneCollection) GetRoot() (*Bone, error) {
	return bones.GetByName(ROOT.String())
}

// CreateRoot すべての親作成
func (bones *BoneCollection) CreateRoot() (*Bone, error) {
	bone := NewBoneByName(ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE
	return bone, nil
}

// GetCenter センター取得
func (bones *BoneCollection) GetCenter() (*Bone, error) {
	return bones.GetByName(CENTER.String())
}

// GetGroove グルーブ取得
func (bones *BoneCollection) GetGroove() (*Bone, error) {
	return bones.GetByName(GROOVE.String())
}

// CreateGroove グルーブ作成
func (bones *BoneCollection) CreateGroove() (*Bone, error) {
	bone := NewBoneByName(GROOVE.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if upper, err := bones.GetUpper(); err == nil {
		bone.Position.Y = upper.Position.Y * 0.7
	} else {
		return nil, merrors.NewParentNotFoundError(
			GROOVE.String(),
			[]string{CENTER.String()},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(GROOVE, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetBodyAxis 体軸取得
func (bones *BoneCollection) GetBodyAxis() (*Bone, error) {
	return bones.GetByName(BODY_AXIS.String())
}

// CreateBodyAxis 体軸作成
func (bones *BoneCollection) CreateBodyAxis() (*Bone, error) {
	bone := NewBoneByName(BODY_AXIS.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	legLeft, _ := bones.GetLeg(BONE_DIRECTION_LEFT)
	legRight, _ := bones.GetLeg(BONE_DIRECTION_RIGHT)
	if legLeft != nil && legRight != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: 0.0,
			Y: (legLeft.Position.Y + legRight.Position.Y) * 0.5,
			Z: 0.0,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			BODY_AXIS.String(),
			[]string{LEG.Left(), LEG.Right()},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(GROOVE, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetWaist 腰取得
func (bones *BoneCollection) GetWaist() (*Bone, error) {
	return bones.GetByName(WAIST.String())
}

// CreateWaist 腰作成
func (bones *BoneCollection) CreateWaist() (*Bone, error) {
	bone := NewBoneByName(WAIST.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	upper, _ := bones.GetUpper()
	lower, _ := bones.GetLower()
	if upper != nil && lower != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (upper.Position.X + lower.Position.X) * 0.5,
			Y: (upper.Position.Y + lower.Position.Y) * 0.5,
			Z: (upper.Position.Z + lower.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			WAIST.String(),
			[]string{UPPER.String(), LOWER.String()},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WAIST, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetTrunkRoot 体幹中心取得
func (bones *BoneCollection) GetTrunkRoot() (*Bone, error) {
	return bones.GetByName(TRUNK_ROOT.String())
}

// CreateTrunkRoot 体幹中心作成
func (bones *BoneCollection) CreateTrunkRoot() (*Bone, error) {
	bone := NewBoneByName(TRUNK_ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	upper, _ := bones.GetUpper()
	lower, _ := bones.GetLower()
	if upper != nil && lower != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (upper.Position.X + lower.Position.X) * 0.5,
			Y: (upper.Position.Y + lower.Position.Y) * 0.5,
			Z: (upper.Position.Z + lower.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			TRUNK_ROOT.String(),
			[]string{UPPER.String(), LOWER.String()},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TRUNK_ROOT, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetLowerRoot 下半身根元取得
func (bones *BoneCollection) GetLowerRoot() (*Bone, error) {
	return bones.GetByName(LOWER_ROOT.String())
}

// CreateLowerRoot 下半身根元作成
func (bones *BoneCollection) CreateLowerRoot() (*Bone, error) {
	bone := NewBoneByName(LOWER_ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	lower, _ := bones.GetLower()
	if lower != nil {
		bone.Position = lower.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			LOWER_ROOT.String(),
			[]string{LOWER.String()},
		)
	}

	// 親ボーン
	for _, parentName := range bone.Config().ParentBoneNames {
		if parent, err := bones.GetByName(parentName.StringFromDirection(bone.Direction())); err == nil {
			bone.ParentIndex = parent.Index()
			break
		}
	}

	return bone, nil
}

// GetLower 下半身
func (bones *BoneCollection) GetLower() (*Bone, error) {
	return bones.GetByName(LOWER.String())
}

// GetLegCenter 足中心取得
func (bones *BoneCollection) GetLegCenter() (*Bone, error) {
	return bones.GetByName(LEG_CENTER.String())
}

// CreateLegCenter 足中心作成
func (bones *BoneCollection) CreateLegCenter() (*Bone, error) {
	bone := NewBoneByName(LEG_CENTER.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	legLeft, _ := bones.GetLeg(BONE_DIRECTION_LEFT)
	legRight, _ := bones.GetLeg(BONE_DIRECTION_RIGHT)
	if legLeft != nil && legRight != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (legLeft.Position.X + legRight.Position.X) * 0.5,
			Y: (legLeft.Position.Y + legRight.Position.Y) * 0.5,
			Z: (legLeft.Position.Z + legRight.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			LEG_CENTER.String(),
			[]string{LEG.Left(), LEG.Right()},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(LEG_CENTER, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetUpperRoot 上半身根元取得
func (bones *BoneCollection) GetUpperRoot() (*Bone, error) {
	return bones.GetByName(UPPER_ROOT.String())
}

// CreateUpperRoot 上半身根元作成
func (bones *BoneCollection) CreateUpperRoot() (*Bone, error) {
	bone := NewBoneByName(UPPER_ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	upper, _ := bones.GetUpper()
	if upper != nil {
		bone.Position = upper.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			UPPER_ROOT.String(),
			[]string{UPPER.String()},
		)
	}

	// 親ボーン
	for _, parentName := range bone.Config().ParentBoneNames {
		if parent, err := bones.GetByName(parentName.StringFromDirection(bone.Direction())); err == nil {
			bone.ParentIndex = parent.Index()
			break
		}
	}

	return bone, nil
}

// GetUpper 上半身取得
func (bones *BoneCollection) GetUpper() (*Bone, error) {
	return bones.GetByName(UPPER.String())
}

// GetUpper2 上半身2取得
func (bones *BoneCollection) GetUpper2() (*Bone, error) {
	return bones.GetByName(UPPER2.String())
}

// CreateUpper2 上半身2作成
func (bones *BoneCollection) CreateUpper2() (*Bone, error) {
	bone := NewBoneByName(UPPER2.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	upper, _ := bones.GetUpper()
	neck, _ := bones.GetNeck()
	if upper != nil && neck != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (upper.Position.X + neck.Position.X) * 0.5,
			Y: (upper.Position.Y + neck.Position.Y) * 0.5,
			Z: (upper.Position.Z + neck.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			UPPER2.String(),
			[]string{UPPER.String(), NECK.String()},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(UPPER2, BONE_DIRECTION_TRUNK)

	// 表示先ボーン
	if neck, err := bones.GetNeck(); err == nil {
		bone.TailIndex = neck.Index()
		bone.BoneFlag |= BONE_FLAG_TAIL_IS_BONE
	}

	return bone, nil
}

// GetNeckRoot 首根元取得
func (bones *BoneCollection) GetNeckRoot() (*Bone, error) {
	return bones.GetByName(NECK_ROOT.String())
}

// CreateNeckRoot 首根元作成
func (bones *BoneCollection) CreateNeckRoot() (*Bone, error) {
	bone := NewBoneByName(NECK_ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	armLeft, _ := bones.GetArm(BONE_DIRECTION_LEFT)
	armRight, _ := bones.GetArm(BONE_DIRECTION_RIGHT)
	if armLeft != nil && armRight != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (armLeft.Position.X + armRight.Position.X) * 0.5,
			Y: (armLeft.Position.Y + armRight.Position.Y) * 0.5,
			Z: (armLeft.Position.Z + armRight.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			NECK_ROOT.String(),
			[]string{ARM.Left(), ARM.Right()},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(NECK_ROOT, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetNeck 首取得
func (bones *BoneCollection) GetNeck() (*Bone, error) {
	return bones.GetByName(NECK.String())
}

// CreateNeck 首作成
func (bones *BoneCollection) CreateNeck() (*Bone, error) {
	bone := NewBoneByName(NECK.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	armLeft, _ := bones.GetArm(BONE_DIRECTION_LEFT)
	armRight, _ := bones.GetArm(BONE_DIRECTION_RIGHT)
	if armLeft != nil && armRight != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (armLeft.Position.X + armRight.Position.X) * 0.5,
			Y: (armLeft.Position.Y + armRight.Position.Y) * 0.5,
			Z: (armLeft.Position.Z + armRight.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			NECK.String(),
			[]string{ARM.Left(), ARM.Right()},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(NECK, BONE_DIRECTION_TRUNK)

	// 表示先ボーン
	if head, err := bones.GetHead(); err == nil {
		bone.TailIndex = head.Index()
		bone.BoneFlag |= BONE_FLAG_TAIL_IS_BONE
	}

	return bone, nil
}

// GetHead 頭取得
func (bones *BoneCollection) GetHead() (*Bone, error) {
	return bones.GetByName(HEAD.String())
}

// CreateHead 頭作成
func (bones *BoneCollection) CreateHead() (*Bone, error) {
	bone := NewBoneByName(HEAD.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if neck, err := bones.GetNeck(); err == nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: neck.Position.X,
			Y: neck.Position.Y * 1.1,
			Z: neck.Position.Z,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			HEAD.String(),
			[]string{NECK.String()},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(HEAD, BONE_DIRECTION_TRUNK)

	// 表示先位置
	bone.TailPosition = mmath.Vec3{Vec: r3.Vec{X: 0, Y: 0.5, Z: 0}}

	return bone, nil
}

// GetHeadTail 頭先取得
func (bones *BoneCollection) GetHeadTail() (*Bone, error) {
	return bones.GetByName(HEAD_TAIL.String())
}

// CreateHeadTail 頭先作成
func (bones *BoneCollection) CreateHeadTail() (*Bone, error) {
	bone := NewBoneByName(HEAD_TAIL.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	neck, _ := bones.GetNeck()
	head, _ := bones.GetHead()
	if neck != nil && head != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: head.Position.X,
			Y: head.Position.Y + (head.Position.Y - neck.Position.Y),
			Z: head.Position.Z,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			HEAD_TAIL.String(),
			[]string{NECK.String(), HEAD.String()},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(HEAD_TAIL, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetEyes 両目取得
func (bones *BoneCollection) GetEyes() (*Bone, error) {
	return bones.GetByName(EYES.String())
}

// GetEye 目取得
func (bones *BoneCollection) GetEye(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(EYE.StringFromDirection(direction))
}

// GetShoulderRoot 肩根元取得
func (bones *BoneCollection) GetShoulderRoot(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_ROOT.StringFromDirection(direction))
}

// CreateShoulderRoot 肩根元作成
func (bones *BoneCollection) CreateShoulderRoot(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_ROOT.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if shoulderBone, err := bones.GetShoulder(direction); err == nil && shoulderBone != nil {
		bone.Position = shoulderBone.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			SHOULDER_ROOT.StringFromDirection(direction),
			[]string{SHOULDER.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	if neckRoot, err := bones.GetNeckRoot(); err == nil {
		bone.ParentIndex = neckRoot.Index()
	}

	return bone, nil
}

// GetShoulderP 肩P取得
func (bones *BoneCollection) GetShoulderP(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_P.StringFromDirection(direction))
}

// CreateShoulderP 肩P作成
func (bones *BoneCollection) CreateShoulderP(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_P.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if shoulder, err := bones.GetShoulder(direction); err == nil {
		bone.Position = shoulder.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			SHOULDER_P.StringFromDirection(direction),
			[]string{SHOULDER.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(SHOULDER_P, direction)

	return bone, nil
}

// GetShoulder 肩取得
func (bones *BoneCollection) GetShoulder(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER.StringFromDirection(direction))
}

// GetShoulderC 肩C取得
func (bones *BoneCollection) GetShoulderC(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_C.StringFromDirection(direction))
}

// CreateShoulderC 肩C作成
func (bones *BoneCollection) CreateShoulderC(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_C.StringFromDirection(direction))

	// 位置
	if arm, err := bones.GetArm(direction); err == nil {
		bone.Position = arm.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			SHOULDER_C.StringFromDirection(direction),
			[]string{SHOULDER.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(SHOULDER_C, direction)

	// 付与親
	if shoulderP, err := bones.GetShoulderP(direction); err == nil {
		bone.EffectIndex = shoulderP.Index()
		bone.EffectFactor = -1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetArm 腕取得
func (bones *BoneCollection) GetArm(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ARM.StringFromDirection(direction))
}

// GetArmTwist 腕捩取得
func (bones *BoneCollection) GetArmTwist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ARM_TWIST.StringFromDirection(direction))
}

// CreateArmTwist 腕捩作成
func (bones *BoneCollection) CreateArmTwist(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ARM_TWIST.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	arm, _ := bones.GetArm(direction)
	elbow, _ := bones.GetElbow(direction)
	if arm != nil && elbow != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (arm.Position.X + elbow.Position.X) * 0.5,
			Y: (arm.Position.Y + elbow.Position.Y) * 0.5,
			Z: (arm.Position.Z + elbow.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			ARM_TWIST.StringFromDirection(direction),
			[]string{
				ARM.StringFromDirection(direction), ELBOW.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(ARM_TWIST, direction)

	// 固定軸
	bone.FixedAxis = elbow.Position.Subed(bone.Position).Normalized()
	bone.BoneFlag |= BONE_FLAG_HAS_FIXED_AXIS

	// ローカル軸
	bone.LocalAxisX = elbow.Position.Subed(bone.Position).Normalized()
	bone.LocalAxisZ = mmath.UNIT_Y_NEG_VEC3.Cross(bone.LocalAxisX).Normalized()
	bone.BoneFlag |= BONE_FLAG_HAS_LOCAL_AXIS

	return bone, nil
}

// GetArmTwistChild 腕捩分割取得
func (bones *BoneCollection) GetArmTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	return bones.GetByName(ARM_TWIST.StringFromDirectionAndIdx(direction, idx+1))
}

// CreateArmTwistChild 腕捩分割作成
func (bones *BoneCollection) CreateArmTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	bone := NewBoneByName(ARM_TWIST.StringFromDirectionAndIdx(direction, idx+1))

	var ratio float64
	switch idx {
	case 0:
		ratio = 0.25
	case 1:
		ratio = 0.5
	case 2:
		ratio = 0.75
	}

	// 位置
	arm, _ := bones.GetArm(direction)
	elbow, _ := bones.GetElbow(direction)
	if arm != nil && elbow != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: arm.Position.X + ((elbow.Position.X - arm.Position.X) * ratio),
			Y: arm.Position.Y + ((elbow.Position.Y - arm.Position.Y) * ratio),
			Z: arm.Position.Z + ((elbow.Position.Z - arm.Position.Z) * ratio),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			ARM_TWIST.StringFromDirectionAndIdx(direction, idx),
			[]string{
				ARM.StringFromDirection(direction), ELBOW.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(ARM_TWIST, direction)

	// 付与親
	if armTwist, err := bones.GetArmTwist(direction); err == nil {
		bone.EffectIndex = armTwist.Index()
		bone.EffectFactor = ratio
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetElbowRoot ひじ取得
func (bones *BoneCollection) GetElbow(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ELBOW.StringFromDirection(direction))
}

// GetWristTwist 腕捩取得
func (bones *BoneCollection) GetWristTwist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST_TWIST.StringFromDirection(direction))
}

// CreateWristTwist 腕捩作成
func (bones *BoneCollection) CreateWristTwist(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WRIST_TWIST.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	elbow, _ := bones.GetElbow(direction)
	wrist, _ := bones.GetWrist(direction)
	if elbow != nil && wrist != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (elbow.Position.X + wrist.Position.X) * 0.5,
			Y: (elbow.Position.Y + wrist.Position.Y) * 0.5,
			Z: (elbow.Position.Z + wrist.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			WRIST_TWIST.StringFromDirection(direction),
			[]string{
				ELBOW.StringFromDirection(direction), WRIST.StringFromDirection(direction)},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WRIST_TWIST, direction)

	// 固定軸
	bone.FixedAxis = wrist.Position.Subed(bone.Position).Normalized()
	bone.BoneFlag |= BONE_FLAG_HAS_FIXED_AXIS

	// ローカル軸
	bone.LocalAxisX = wrist.Position.Subed(bone.Position).Normalized()
	bone.LocalAxisZ = mmath.UNIT_Y_NEG_VEC3.Cross(bone.LocalAxisX).Normalized()
	bone.BoneFlag |= BONE_FLAG_HAS_LOCAL_AXIS

	return bone, nil
}

// GetWristTwistChild 腕捩分割取得
func (bones *BoneCollection) GetWristTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	return bones.GetByName(WRIST_TWIST.StringFromDirectionAndIdx(direction, idx+1))
}

// CreateWristTwistChild 腕捩分割作成
func (bones *BoneCollection) CreateWristTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	bone := NewBoneByName(WRIST_TWIST.StringFromDirectionAndIdx(direction, idx+1))

	var ratio float64
	switch idx {
	case 0:
		ratio = 0.25
	case 1:
		ratio = 0.5
	case 2:
		ratio = 0.75
	}

	// 位置
	elbow, _ := bones.GetElbow(direction)
	wrist, _ := bones.GetWrist(direction)
	if elbow != nil && wrist != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: elbow.Position.X + ((wrist.Position.X - elbow.Position.X) * ratio),
			Y: elbow.Position.Y + ((wrist.Position.Y - elbow.Position.Y) * ratio),
			Z: elbow.Position.Z + ((wrist.Position.Z - elbow.Position.Z) * ratio),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			WRIST_TWIST.StringFromDirectionAndIdx(direction, idx),
			[]string{
				ELBOW.StringFromDirection(direction), WRIST.StringFromDirection(direction)},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WRIST_TWIST, direction)

	// 付与親
	if wristTwist, err := bones.GetWristTwist(direction); err == nil {
		bone.EffectIndex = wristTwist.Index()
		bone.EffectFactor = ratio
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetWrist 手首取得
func (bones *BoneCollection) GetWrist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST.StringFromDirection(direction))
}

// GetWristTail 手首先先取得
func (bones *BoneCollection) GetWristTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST_TAIL.StringFromDirection(direction))
}

// CreateWristTail 手首先先作成
func (bones *BoneCollection) CreateWristTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WRIST_TAIL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	bonePositions := make([]mmath.Vec3, 0)
	for _, boneName := range []string{THUMB1.StringFromDirection(direction),
		INDEX1.StringFromDirection(direction), MIDDLE1.StringFromDirection(direction),
		RING1.StringFromDirection(direction), PINKY1.StringFromDirection(direction)} {
		if bone, err := bones.GetByName(boneName); err == nil {
			bonePositions = append(bonePositions, bone.Position)
		}
	}
	bone.Position = mmath.MeanVec3(bonePositions)
	if bone.Position.IsZero() {
		// 指がなくて位置が取れなかった場合、ひじからの相対位置を利用する
		if elbow, err := bones.GetElbow(direction); err == nil && elbow != nil {
			wrist, _ := bones.GetWrist(direction)
			if wrist != nil {
				bone.Position = mmath.Vec3{Vec: r3.Vec{
					X: wrist.Position.X + (wrist.Position.X-elbow.Position.X)*0.5,
					Y: wrist.Position.Y + (wrist.Position.Y-elbow.Position.Y)*0.5,
					Z: wrist.Position.Z + (wrist.Position.Z-elbow.Position.Z)*0.5,
				}}
			}
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WRIST_TAIL, direction)

	return bone, nil
}

// GetThumb 親指取得
func (bones *BoneCollection) GetThumb(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(THUMB0.StringFromDirection(direction))
	case 1:
		return bones.GetByName(THUMB1.StringFromDirection(direction))
	case 2:
		return bones.GetByName(THUMB2.StringFromDirection(direction))
	}

	return nil, merrors.NewInvalidIndexError(idx)
}

// CreateThumb0 親指0作成
func (bones *BoneCollection) CreateThumb0(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(THUMB0.StringFromDirection(direction))

	// 位置
	wrist, _ := bones.GetWrist(direction)
	thumb1, _ := bones.GetThumb(direction, 1)
	if wrist != nil && thumb1 != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: wrist.Position.X + (thumb1.Position.X-wrist.Position.X)*0.5,
			Y: wrist.Position.Y + (thumb1.Position.Y-wrist.Position.Y)*0.5,
			Z: wrist.Position.Z + (thumb1.Position.Z-wrist.Position.Z)*0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			THUMB0.StringFromDirection(direction),
			[]string{
				WRIST.StringFromDirection(direction), THUMB1.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(THUMB0, direction)

	return bone, nil
}

// GetThumbTail 親指先先取得
func (bones *BoneCollection) GetThumbTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(THUMB_TAIL.StringFromDirection(direction))
}

// CreateThumbTail 親指先先作成
func (bones *BoneCollection) CreateThumbTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(THUMB_TAIL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	thumb1, _ := bones.GetThumb(direction, 1)
	thumb2, _ := bones.GetThumb(direction, 2)
	if thumb1 != nil && thumb2 != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: thumb2.Position.X + (thumb2.Position.X - thumb1.Position.X),
			Y: thumb2.Position.Y + (thumb2.Position.Y - thumb1.Position.Y),
			Z: thumb2.Position.Z + (thumb2.Position.Z - thumb1.Position.Z),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			THUMB_TAIL.StringFromDirection(direction),
			[]string{
				THUMB1.StringFromDirection(direction), THUMB2.StringFromDirection(direction)},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(THUMB_TAIL, direction)

	return bone, nil
}

// GetIndex 人差し指取得
func (bones *BoneCollection) GetIndex(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(INDEX1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(INDEX2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(INDEX3.StringFromDirection(direction))
	}

	return nil, merrors.NewInvalidIndexError(idx)
}

// GetIndexTail 人差し指先先取得
func (bones *BoneCollection) GetIndexTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(INDEX_TAIL.StringFromDirection(direction))
}

// CreateIndexTail 親指先先作成
func (bones *BoneCollection) CreateIndexTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(INDEX_TAIL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	index1, _ := bones.GetIndex(direction, 1)
	index2, _ := bones.GetIndex(direction, 2)
	if index1 != nil && index2 != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: index2.Position.X + (index2.Position.X - index1.Position.X),
			Y: index2.Position.Y + (index2.Position.Y - index1.Position.Y),
			Z: index2.Position.Z + (index2.Position.Z - index1.Position.Z),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			INDEX_TAIL.StringFromDirection(direction),
			[]string{
				INDEX1.StringFromDirection(direction), INDEX2.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(INDEX_TAIL, direction)

	return bone, nil
}

// GetMiddle 中指取得
func (bones *BoneCollection) GetMiddle(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(MIDDLE1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(MIDDLE2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(MIDDLE3.StringFromDirection(direction))
	}

	return nil, merrors.NewInvalidIndexError(idx)
}

// GetMiddleTail 中指先先取得
func (bones *BoneCollection) GetMiddleTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(MIDDLE_TAIL.StringFromDirection(direction))
}

// CreateMiddleTail 親指先先作成
func (bones *BoneCollection) CreateMiddleTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(MIDDLE_TAIL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	middle1, _ := bones.GetMiddle(direction, 1)
	middle2, _ := bones.GetMiddle(direction, 2)
	if middle1 != nil && middle2 != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: middle2.Position.X + (middle2.Position.X - middle1.Position.X),
			Y: middle2.Position.Y + (middle2.Position.Y - middle1.Position.Y),
			Z: middle2.Position.Z + (middle2.Position.Z - middle1.Position.Z),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			MIDDLE_TAIL.StringFromDirection(direction),
			[]string{MIDDLE1.StringFromDirection(direction), MIDDLE2.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(MIDDLE_TAIL, direction)

	return bone, nil
}

// GetRing 薬指取得
func (bones *BoneCollection) GetRing(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(RING1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(RING2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(RING3.StringFromDirection(direction))
	}

	return nil, merrors.NewInvalidIndexError(idx)
}

// GetRingTail 薬指先先取得
func (bones *BoneCollection) GetRingTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(RING_TAIL.StringFromDirection(direction))
}

// CreateRingTail 親指先先作成
func (bones *BoneCollection) CreateRingTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(RING_TAIL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	ring1, _ := bones.GetRing(direction, 1)
	ring2, _ := bones.GetRing(direction, 2)
	if ring1 != nil && ring2 != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: ring2.Position.X + (ring2.Position.X - ring1.Position.X),
			Y: ring2.Position.Y + (ring2.Position.Y - ring1.Position.Y),
			Z: ring2.Position.Z + (ring2.Position.Z - ring1.Position.Z),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			RING_TAIL.StringFromDirection(direction),
			[]string{
				RING1.StringFromDirection(direction), RING2.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(RING_TAIL, direction)

	return bone, nil
}

// GetPinky 小指取得
func (bones *BoneCollection) GetPinky(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(PINKY1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(PINKY2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(PINKY3.StringFromDirection(direction))
	}

	return nil, merrors.NewInvalidIndexError(idx)
}

// GetPinkyTail 小指先先取得
func (bones *BoneCollection) GetPinkyTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(PINKY_TAIL.StringFromDirection(direction))
}

// CreatePinkyTail 小指先先作成
func (bones *BoneCollection) CreatePinkyTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(PINKY_TAIL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	pinky1, _ := bones.GetPinky(direction, 1)
	pinky2, _ := bones.GetPinky(direction, 2)
	if pinky1 != nil && pinky2 != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: pinky2.Position.X + (pinky2.Position.X - pinky1.Position.X),
			Y: pinky2.Position.Y + (pinky2.Position.Y - pinky1.Position.Y),
			Z: pinky2.Position.Z + (pinky2.Position.Z - pinky1.Position.Z),
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			PINKY_TAIL.StringFromDirection(direction),
			[]string{
				PINKY1.StringFromDirection(direction), PINKY2.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(PINKY_TAIL, direction)

	return bone, nil
}

// GetWaistCancel 腰キャンセル取得
func (bones *BoneCollection) GetWaistCancel(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WAIST_CANCEL.StringFromDirection(direction))
}

// CreateWaistCancel 腰キャンセル作成
func (bones *BoneCollection) CreateWaistCancel(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WAIST_CANCEL.StringFromDirection(direction))

	// 位置
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.Position = leg.Position
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WAIST_CANCEL, direction)

	// 付与親
	if waist, err := bones.GetWaist(); err == nil {
		bone.EffectIndex = waist.Index()
		bone.EffectFactor = -1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	} else {
		return nil, merrors.NewParentNotFoundError(
			WAIST_CANCEL.StringFromDirection(direction),
			[]string{WAIST.String()},
		)
	}

	return bone, nil
}

// GetLegRoot 足根元取得
func (bones *BoneCollection) GetLegRoot(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_ROOT.StringFromDirection(direction))
}

// CreateLegRoot 足根元作成
func (bones *BoneCollection) CreateLegRoot(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_ROOT.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if legBone, err := bones.GetLeg(direction); err == nil && legBone != nil {
		bone.Position = legBone.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			LEG_ROOT.StringFromDirection(direction),
			[]string{LEG.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(LEG_ROOT, direction)

	return bone, nil
}

// GetHip 腰骨取得
func (bones *BoneCollection) GetHip(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HIP.StringFromDirection(direction))
}

// CreateHip 腰骨作成
func (bones *BoneCollection) CreateHip(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HIP.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	// 親ボーン
	if lower, err := bones.GetLower(); err == nil && lower != nil {
		if leg, err := bones.GetLeg(direction); err == nil {
			bone.Position = mmath.Vec3{Vec: r3.Vec{
				X: leg.Position.X,
				Y: lower.Position.Y,
				Z: lower.Position.Z,
			}}
		}
	} else {
		return nil, merrors.NewParentNotFoundError(
			HIP.StringFromDirection(direction),
			[]string{LOWER.String(), LEG.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(HIP, direction)

	return bone, nil
}

// GetLeg 足取得
func (bones *BoneCollection) GetLeg(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG.StringFromDirection(direction))
}

// GetKnee ひざ取得
func (bones *BoneCollection) GetKnee(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(KNEE.StringFromDirection(direction))
}

// GetAnkle 足首取得
func (bones *BoneCollection) GetAnkle(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE.StringFromDirection(direction))
}

// GetHeel かかと取得
func (bones *BoneCollection) GetHeel(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HEEL.StringFromDirection(direction))
}

// CreateHeel かかと作成
func (bones *BoneCollection) CreateHeel(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HEEL.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	// 親ボーン
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: ankle.Position.X,
			Y: 0.0,
			Z: ankle.Position.Z + 0.2,
		}}

		bone.ParentIndex = ankle.Index()
	} else {
		return nil, merrors.NewParentNotFoundError(
			HEEL.StringFromDirection(direction),
			[]string{ANKLE.StringFromDirection(direction)},
		)
	}

	return bone, nil
}

// GetToeT つま先先取得
func (bones *BoneCollection) GetToeT(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_T.StringFromDirection(direction))
}

// CreateToeT つま先先作成
func (bones *BoneCollection) CreateToeT(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_T.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if toeIK, err := bones.GetToeIK(direction); err == nil {
		if toe, err := bones.Get(toeIK.Ik.BoneIndex); err == nil {
			// つま先IKのターゲットと同位置
			bone.Position = toe.Position
		}
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_T.StringFromDirection(direction),
			[]string{TOE_IK.StringFromDirection(direction)},
		)
	}

	bone.Position.Y = 0.0

	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.ParentIndex = ankle.Index()
	}

	return bone, nil
}

// GetToeP つま先親取得
func (bones *BoneCollection) GetToeP(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_P.StringFromDirection(direction))
}

// CreateToeP つま先親作成
func (bones *BoneCollection) CreateToeP(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_P.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	// 親ボーン
	if toeT, err := bones.GetToeT(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = mmath.Vec3{Vec: r3.Vec{
				X: toeT.Position.X - 1.0,
				Y: 0.0,
				Z: toeT.Position.Z,
			}}
		case BONE_DIRECTION_RIGHT:
			bone.Position = mmath.Vec3{Vec: r3.Vec{
				X: toeT.Position.X + 1.0,
				Y: 0.0,
				Z: toeT.Position.Z,
			}}
		}

		bone.ParentIndex = toeT.Index()
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_P.StringFromDirection(direction),
			[]string{TOE_T.StringFromDirection(direction)},
		)
	}
	return bone, nil
}

// GetToeC つま先子取得
func (bones *BoneCollection) GetToeC(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_C.StringFromDirection(direction))
}

// CreateToeC つま先子作成
func (bones *BoneCollection) CreateToeC(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_C.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if toeT, err := bones.GetToeT(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = mmath.Vec3{Vec: r3.Vec{
				X: toeT.Position.X + 1.0,
				Y: 0.0,
				Z: toeT.Position.Z,
			}}
		case BONE_DIRECTION_RIGHT:
			bone.Position = mmath.Vec3{Vec: r3.Vec{
				X: toeT.Position.X - 1.0,
				Y: 0.0,
				Z: toeT.Position.Z,
			}}
		}
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_C.StringFromDirection(direction),
			[]string{TOE_T.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_C, direction)

	return bone, nil
}

// GetLegD 足D取得
func (bones *BoneCollection) GetLegD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_D.StringFromDirection(direction))
}

// CreateLegD 足D作成
func (bones *BoneCollection) CreateLegD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.Position = leg.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			LEG_D.StringFromDirection(direction),
			[]string{LEG.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(LEG_D, direction)

	// 付与親
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.EffectIndex = leg.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetKneeD ひざD取得
func (bones *BoneCollection) GetKneeD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(KNEE_D.StringFromDirection(direction))
}

// CreateKneeD ひざD作成
func (bones *BoneCollection) CreateKneeD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(KNEE_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if knee, err := bones.GetKnee(direction); err == nil {
		bone.Position = knee.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			KNEE_D.StringFromDirection(direction),
			[]string{KNEE.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(KNEE_D, direction)

	// 付与親
	if knee, err := bones.GetKnee(direction); err == nil {
		bone.EffectIndex = knee.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetAnkleD 足首D取得
func (bones *BoneCollection) GetAnkleD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE_D.StringFromDirection(direction))
}

// CreateAnkleD 足首D作成
func (bones *BoneCollection) CreateAnkleD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ANKLE_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.Position = ankle.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			ANKLE_D.StringFromDirection(direction),
			[]string{ANKLE.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(ANKLE_D, direction)

	// 付与親
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.EffectIndex = ankle.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetAnkleDGround 足首地面取得
func (bones *BoneCollection) GetAnkleDGround(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE_D_GROUND.StringFromDirection(direction))
}

// CreateAnkleDGround 足首地面作成
func (bones *BoneCollection) CreateAnkleDGround(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ANKLE_D_GROUND.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if ankle, err := bones.GetAnkleD(direction); err == nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: ankle.Position.X,
			Y: 0.0,
			Z: ankle.Position.Z,
		}}

		bone.ParentIndex = ankle.Index()
	} else {
		return nil, merrors.NewParentNotFoundError(
			ANKLE_GROUND.StringFromDirection(direction),
			[]string{ANKLE.StringFromDirection(direction)},
		)
	}

	return bone, nil
}

// GetHeelD かかとD取得
func (bones *BoneCollection) GetHeelD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HEEL_D.StringFromDirection(direction))
}

// CreateHeelD かかとD作成
func (bones *BoneCollection) CreateHeelD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HEEL_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if heel, err := bones.GetHeel(direction); err == nil {
		bone.Position = heel.Position
	}

	// 親ボーン
	if ankleD, err := bones.GetAnkleD(direction); err == nil {
		bone.ParentIndex = ankleD.Index()
	}

	// 付与親
	if heel, err := bones.GetHeel(direction); err == nil {
		bone.EffectIndex = heel.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	} else {
		return nil, merrors.NewParentNotFoundError(
			HEEL_D.StringFromDirection(direction),
			[]string{HEEL.StringFromDirection(direction)},
		)
	}

	return bone, nil
}

// GetToeEx 足先EX取得
func (bones *BoneCollection) GetToeEx(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_EX.StringFromDirection(direction))
}

// CreateToeEx 足先EX作成
func (bones *BoneCollection) CreateToeEx(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_EX.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	ankle, _ := bones.GetAnkle(direction)
	toeT, _ := bones.GetToeT(direction)
	if ankle != nil && toeT != nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: (ankle.Position.X + toeT.Position.X) * 0.5,
			Y: (ankle.Position.Y + toeT.Position.Y) * 0.5,
			Z: (ankle.Position.Z + toeT.Position.Z) * 0.5,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_EX.StringFromDirection(direction),
			[]string{
				ANKLE.StringFromDirection(direction), TOE_T.StringFromDirection(direction)},
		)

	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_EX, direction)

	return bone, nil
}

// GetToeTD つま先先D取得
func (bones *BoneCollection) GetToeTD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_T_D.StringFromDirection(direction))
}

// CreateToeTD つま先先D作成
func (bones *BoneCollection) CreateToeTD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_T_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if toeT, err := bones.GetToeT(direction); err == nil {
		bone.Position = toeT.Position
	}
	bone.Position.Y = 0.0

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_T_D, direction)

	// 付与親
	if toeT, err := bones.GetToeT(direction); err == nil {
		bone.EffectIndex = toeT.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_T_D.StringFromDirection(direction),
			[]string{TOE_T.StringFromDirection(direction)},
		)
	}

	return bone, nil
}

// GetToePD つま先親D取得
func (bones *BoneCollection) GetToePD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_P_D.StringFromDirection(direction))
}

// CreateToePD つま先親D作成
func (bones *BoneCollection) CreateToePD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_P_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if toeP, err := bones.GetToeP(direction); err == nil {
		bone.Position = toeP.Position
	}
	bone.Position.Y = 0.0

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_P_D, direction)

	// 付与親
	if toeP, err := bones.GetToeP(direction); err == nil {
		bone.EffectIndex = toeP.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_P_D.StringFromDirection(direction),
			[]string{TOE_P.StringFromDirection(direction)},
		)
	}

	return bone, nil
}

// GetToeCD つま先子D取得
func (bones *BoneCollection) GetToeCD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_C_D.StringFromDirection(direction))
}

// CreateToeCD つま先子D作成
func (bones *BoneCollection) CreateToeCD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_C_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if toeC, err := bones.GetToeC(direction); err == nil {
		bone.Position = toeC.Position
	} else {
		return nil, merrors.NewParentNotFoundError(
			TOE_C_D.StringFromDirection(direction),
			[]string{TOE_C.StringFromDirection(direction)},
		)
	}
	bone.Position.Y = 0.0

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_C_D, direction)

	// 付与親
	if toeC, err := bones.GetToeC(direction); err == nil {
		bone.EffectIndex = toeC.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetLegIkParent 足IK親取得
func (bones *BoneCollection) GetLegIkParent(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_IK_PARENT.StringFromDirection(direction))
}

// CreateLegIkParent 足IK親作成
func (bones *BoneCollection) CreateLegIkParent(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_IK_PARENT.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if legIk, err := bones.GetLegIk(direction); err == nil {
		bone.Position = mmath.Vec3{Vec: r3.Vec{
			X: legIk.Position.X,
			Y: 0.0,
			Z: legIk.Position.Z,
		}}
	} else {
		return nil, merrors.NewParentNotFoundError(
			LEG_IK_PARENT.StringFromDirection(direction),
			[]string{LEG_IK.StringFromDirection(direction)},
		)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(LEG_IK_PARENT, direction)

	return bone, nil
}

// GetLegIk 足IK取得
func (bones *BoneCollection) GetLegIk(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_IK.StringFromDirection(direction))
}

// GetToeIK つま先IK取得
func (bones *BoneCollection) GetToeIK(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_IK.StringFromDirection(direction))
}

// findParentIndexByConfig は標準ボーン設定から親ボーンindexを推定する。
func (bones *BoneCollection) findParentIndexByConfig(boneName StandardBoneName, direction BoneDirection) int {
	boneConfig := GetStandardBoneConfigs()[boneName]
	// 子ボーンが定義されているなら、その子ボーンの親を探す
	for _, tailBoneNames := range boneConfig.ChildBoneNames {
		for _, tailBoneName := range tailBoneNames {
			if bones.ContainsByName(tailBoneName.StringFromDirection(direction)) {
				if bone, err := bones.GetByName(tailBoneName.StringFromDirection(direction)); bone != nil && err == nil {
					return bone.ParentIndex
				}
			} else if bones.ContainsByName(tailBoneName.String()) {
				if bone, err := bones.GetByName(tailBoneName.String()); bone != nil && err == nil {
					return bone.ParentIndex
				}
			} else if bones.ContainsByName(tailBoneName.Left()) {
				// 左右ボーンの親が左右の場合、その親までは辿る
				if bone, err := bones.GetByName(tailBoneName.Left()); bone != nil && err == nil {
					parentIndex := bone.ParentIndex
					for i := 0; i < 5; i++ {
						if parentIndex < 0 {
							break
						}
						parentBone, err := bones.Get(parentIndex)
						if err != nil || parentBone == nil {
							break
						}
						if parentBone.Direction() == direction {
							return parentBone.Index()
						}
						parentIndex = parentBone.ParentIndex
					}
				}
			} else if bones.ContainsByName(tailBoneName.Right()) {
				// 左右ボーンの親が左右の場合、その親までは辿る
				if bone, err := bones.GetByName(tailBoneName.Right()); bone != nil && err == nil {
					parentIndex := bone.ParentIndex
					for i := 0; i < 5; i++ {
						if parentIndex < 0 {
							break
						}
						parentBone, err := bones.Get(parentIndex)
						if err != nil || parentBone == nil {
							break
						}
						if parentBone.Direction() == direction {
							return parentBone.Index()
						}
						parentIndex = parentBone.ParentIndex
					}
				}
			}
		}
	}

	// 親ボーンが定義されているなら、そのボーンを探す
	for _, parentBoneName := range boneConfig.ParentBoneNames {
		if bones.ContainsByName(parentBoneName.StringFromDirection(direction)) {
			if bone, err := bones.GetByName(parentBoneName.StringFromDirection(direction)); bone != nil && err == nil {
				return bone.Index()
			}
		} else if bones.ContainsByName(parentBoneName.String()) {
			if bone, err := bones.GetByName(parentBoneName.String()); bone != nil && err == nil {
				return bone.Index()
			}
		}
	}

	return -1
}

// InsertShortageOverrideBones 不足ボーン作成
func (bones *BoneCollection) InsertShortageOverrideBones() error {

	// 体幹系
	for _, funcs := range [][]func() (*Bone, error){
		{bones.GetTrunkRoot, bones.CreateTrunkRoot},
		{bones.GetLegCenter, bones.CreateLegCenter},
		{bones.GetNeckRoot, bones.CreateNeckRoot},
	} {
		getFunc := funcs[0]
		createFunc := funcs[1]

		if bone, err := getFunc(); err != nil && merrors.IsNameNotFoundError(err) && bone == nil {
			if bone, err := createFunc(); err == nil && bone != nil {
				if _, _, err := bones.Insert(bone, bone.ParentIndex); err != nil {
					return err
				} else {
					// 追加したボーンの親ボーンを、同じく親ボーンに設定しているボーンの親ボーンを追加ボーンに置き換える
					for _, b := range bones.values {
						if b == nil {
							continue
						}
						if b.ParentIndex == bone.ParentIndex && b.Index() != bone.Index() &&
							b.EffectIndex != bone.Index() && bone.EffectIndex != b.Index() &&
							((strings.Contains(bone.Name(), "上") && !strings.Contains(b.Name(), "下") &&
								!strings.Contains(b.Name(), "左") && !strings.Contains(b.Name(), "右")) ||
								(strings.Contains(bone.Name(), "下") && !strings.Contains(b.Name(), "上") &&
									!strings.Contains(b.Name(), "左") && !strings.Contains(b.Name(), "右"))) {
							b.ParentIndex = bone.Index()
						}
					}
				}
			} else {
				return err
			}
		} else if err != nil {
			return err
		} else {
			switch bone.Name() {
			case NECK.String():
				if neckRoot, err := bones.GetNeckRoot(); err == nil {
					bone.ParentIndex = neckRoot.Index()
				}
			}
		}
	}

	return nil
}

// InsertSystemTailBones システム用不足ボーン作成
func (bones *BoneCollection) InsertSystemTailBones() error {

	// 体幹系
	for _, funcs := range [][]func(direction BoneDirection) (*Bone, error){
		{bones.GetToeT, bones.CreateToeT},
		{bones.GetHeel, bones.CreateHeel},
		{bones.GetWristTail, bones.CreateWristTail},
	} {
		getFunc := funcs[0]
		createFunc := funcs[1]

		for _, direction := range []BoneDirection{BONE_DIRECTION_LEFT, BONE_DIRECTION_RIGHT} {
			if bone, err := getFunc(direction); err != nil && merrors.IsNameNotFoundError(err) && bone == nil {
				if bone, err := createFunc(direction); err == nil && bone != nil {
					if _, _, err := bones.Insert(bone, bone.ParentIndex); err != nil {
						return err
					} else {
						// 追加したボーンの親ボーンを、同じく親ボーンに設定しているボーンの親ボーンを追加ボーンに置き換える
						for _, b := range bones.values {
							if b == nil {
								continue
							}
							if b.ParentIndex == bone.ParentIndex && b.Index() != bone.Index() &&
								b.EffectIndex != bone.Index() && bone.EffectIndex != b.Index() &&
								((strings.Contains(bone.Name(), "上") && !strings.Contains(b.Name(), "下") &&
									!strings.Contains(b.Name(), "左") && !strings.Contains(b.Name(), "右")) ||
									(strings.Contains(bone.Name(), "下") && !strings.Contains(b.Name(), "上") &&
										!strings.Contains(b.Name(), "左") && !strings.Contains(b.Name(), "右"))) {
								b.ParentIndex = bone.Index()
							}
						}
					}
				} else {
					return err
				}
			} else if err != nil {
				return err
			}
		}
	}

	return nil
}
