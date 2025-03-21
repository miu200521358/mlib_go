package pmx

import (
	"errors"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// GetRoot すべての親
func (bones *Bones) GetRoot() (*Bone, error) {
	return bones.GetByName(ROOT.String())
}

// CreateRoot すべての親取得or作成
func (bones *Bones) CreateRoot() (*Bone, error) {
	bone := NewBoneByName(ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE
	return bone, nil
}

// GetCenter センター取得
func (bones *Bones) GetCenter() (*Bone, error) {
	return bones.GetByName(CENTER.String())
}

// GetGroove グルーブ取得
func (bones *Bones) GetGroove() (*Bone, error) {
	return bones.GetByName(GROOVE.String())
}

// CreateGroove グルーブ取得or作成
func (bones *Bones) CreateGroove() (*Bone, error) {
	bone := NewBoneByName(GROOVE.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if upper, err := bones.GetUpper(); err == nil {
		bone.Position.Y = upper.Position.Y * 0.7
	}

	// 親ボーン
	if center, err := bones.GetCenter(); err == nil {
		bone.ParentIndex = center.Index()
	}

	return bone, nil
}

// GetWaist 腰取得
func (bones *Bones) GetWaist() (*Bone, error) {
	return bones.GetByName(WAIST.String())
}

// CreateWaist 腰取得or作成
func (bones *Bones) CreateWaist() (*Bone, error) {
	bone := NewBoneByName(WAIST.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	upper, _ := bones.GetUpper()
	lower, _ := bones.GetLower()
	if upper != nil && lower != nil {
		bone.Position = &mmath.MVec3{
			X: (upper.Position.X + lower.Position.X) * 0.5,
			Y: (upper.Position.Y + lower.Position.Y) * 0.5,
			Z: (upper.Position.Z + lower.Position.Z) * 0.5,
		}
	} else if upper != nil {
		bone.Position = upper.Position.Copy()
	} else if lower != nil {
		bone.Position = lower.Position.Copy()
	}

	// 親ボーン
	if groove, err := bones.GetGroove(); err == nil {
		bone.ParentIndex = groove.Index()
	}

	return bone, nil
}

// GetTrunkRoot 体幹中心取得
func (bones *Bones) GetTrunkRoot() (*Bone, error) {
	return bones.GetByName(TRUNK_ROOT.String())
}

// CreateTrunkRoot 体幹中心取得or作成
func (bones *Bones) CreateTrunkRoot() (*Bone, error) {
	bone := NewBoneByName(TRUNK_ROOT.String())

	// 位置
	upper, _ := bones.GetUpper()
	lower, _ := bones.GetLower()
	if upper != nil && lower != nil {
		bone.Position = &mmath.MVec3{
			X: (upper.Position.X + lower.Position.X) * 0.5,
			Y: (upper.Position.Y + lower.Position.Y) * 0.5,
			Z: (upper.Position.Z + lower.Position.Z) * 0.5,
		}
	} else if upper != nil {
		bone.Position = upper.Position.Copy()
	} else if lower != nil {
		bone.Position = lower.Position.Copy()
	}

	// 親ボーン
	if waist, err := bones.GetWaist(); err == nil {
		bone.ParentIndex = waist.Index()
	}

	return bone, nil
}

// GetLowerRoot 下半身根元取得
func (bones *Bones) GetLowerRoot() (*Bone, error) {
	return bones.GetByName(LOWER_ROOT.String())
}

// CreateLowerRoot 下半身根元取得or作成
func (bones *Bones) CreateLowerRoot() (*Bone, error) {
	bone := NewBoneByName(LOWER_ROOT.String())

	// 位置
	upper, _ := bones.GetUpper()
	lower, _ := bones.GetLower()
	if upper != nil && lower != nil {
		bone.Position = &mmath.MVec3{
			X: (upper.Position.X + lower.Position.X) * 0.5,
			Y: (upper.Position.Y + lower.Position.Y) * 0.5,
			Z: (upper.Position.Z + lower.Position.Z) * 0.5,
		}
	} else if upper != nil {
		bone.Position = upper.Position.Copy()
	} else if lower != nil {
		bone.Position = lower.Position.Copy()
	}

	// 親ボーン
	if trunkRoot, err := bones.GetTrunkRoot(); err == nil {
		bone.ParentIndex = trunkRoot.Index()
	}

	return bone, nil
}

// GetLower 下半身
func (bones *Bones) GetLower() (*Bone, error) {
	return bones.GetByName(LOWER.String())
}

// GetLegCenter 足中心取得
func (bones *Bones) GetLegCenter() (*Bone, error) {
	return bones.GetByName(LEG_CENTER.String())
}

// CreateLegCenter 足中心取得or作成
func (bones *Bones) CreateLegCenter() (*Bone, error) {
	bone := NewBoneByName(LEG_CENTER.String())

	// 位置
	legLeft, _ := bones.GetLeg(BONE_DIRECTION_LEFT)
	legRight, _ := bones.GetLeg(BONE_DIRECTION_RIGHT)
	if legLeft != nil && legRight != nil {
		bone.Position = &mmath.MVec3{
			X: (legLeft.Position.X + legRight.Position.X) * 0.5,
			Y: (legLeft.Position.Y + legRight.Position.Y) * 0.5,
			Z: (legLeft.Position.Z + legRight.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if lower, err := bones.GetLower(); err == nil {
		bone.ParentIndex = lower.Index()
	}

	return bone, nil
}

// GetUpperRoot 上半身根元取得
func (bones *Bones) GetUpperRoot() (*Bone, error) {
	return bones.GetByName(UPPER_ROOT.String())
}

// CreateUpperRoot 上半身根元取得or作成
func (bones *Bones) CreateUpperRoot() (*Bone, error) {
	bone := NewBoneByName(UPPER_ROOT.String())

	// 位置
	upper, _ := bones.GetUpper()
	lower, _ := bones.GetLower()
	if upper != nil && lower != nil {
		bone.Position = &mmath.MVec3{
			X: (upper.Position.X + lower.Position.X) * 0.5,
			Y: (upper.Position.Y + lower.Position.Y) * 0.5,
			Z: (upper.Position.Z + lower.Position.Z) * 0.5,
		}
	} else if upper != nil {
		bone.Position = upper.Position.Copy()
	} else if lower != nil {
		bone.Position = lower.Position.Copy()
	}

	// 親ボーン
	if trunkRoot, err := bones.GetTrunkRoot(); err == nil {
		bone.ParentIndex = trunkRoot.Index()
	}

	return bone, nil
}

// GetUpper 上半身取得
func (bones *Bones) GetUpper() (*Bone, error) {
	return bones.GetByName(UPPER.String())
}

// GetUpper2 上半身2取得
func (bones *Bones) GetUpper2() (*Bone, error) {
	return bones.GetByName(UPPER2.String())
}

// CreateUpper2 上半身2取得or作成
func (bones *Bones) CreateUpper2() (*Bone, error) {
	bone := NewBoneByName(UPPER2.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	upper, _ := bones.GetUpper()
	neck, _ := bones.GetNeck()
	if upper != nil && neck != nil {
		bone.Position = &mmath.MVec3{
			X: (upper.Position.X + neck.Position.X) * 0.5,
			Y: (upper.Position.Y + neck.Position.Y) * 0.5,
			Z: (upper.Position.Z + neck.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if upper, err := bones.GetUpper(); err == nil {
		bone.ParentIndex = upper.Index()
	}

	// 表示先ボーン
	if neck, err := bones.GetNeck(); err == nil {
		bone.TailIndex = neck.Index()
		bone.BoneFlag |= BONE_FLAG_TAIL_IS_BONE
	}

	return bone, nil
}

// GetNeckRoot 首根元取得
func (bones *Bones) GetNeckRoot() (*Bone, error) {
	return bones.GetByName(NECK_ROOT.String())
}

// CreateNeckRoot 首根元取得or作成
func (bones *Bones) CreateNeckRoot() (*Bone, error) {
	bone := NewBoneByName(NECK_ROOT.String())

	// 位置
	armLeft, _ := bones.GetArm(BONE_DIRECTION_LEFT)
	armRight, _ := bones.GetArm(BONE_DIRECTION_RIGHT)
	if armLeft != nil && armRight != nil {
		bone.Position = &mmath.MVec3{
			X: (armLeft.Position.X + armRight.Position.X) * 0.5,
			Y: (armLeft.Position.Y + armRight.Position.Y) * 0.5,
			Z: (armLeft.Position.Z + armRight.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if upper2, err := bones.GetUpper2(); err == nil {
		bone.ParentIndex = upper2.Index()
	}

	return bone, nil
}

// GetNeck 首取得
func (bones *Bones) GetNeck() (*Bone, error) {
	return bones.GetByName(NECK.String())
}

// CreateNeck 首取得or作成
func (bones *Bones) CreateNeck() (*Bone, error) {
	bone := NewBoneByName(NECK.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	armLeft, _ := bones.GetArm(BONE_DIRECTION_LEFT)
	armRight, _ := bones.GetArm(BONE_DIRECTION_RIGHT)
	if armLeft != nil && armRight != nil {
		bone.Position = &mmath.MVec3{
			X: (armLeft.Position.X + armRight.Position.X) * 0.5,
			Y: (armLeft.Position.Y + armRight.Position.Y) * 0.5,
			Z: (armLeft.Position.Z + armRight.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if neckRoot, err := bones.GetNeckRoot(); err == nil {
		bone.ParentIndex = neckRoot.Index()
	}

	// 表示先ボーン
	if head, err := bones.GetHead(); err == nil {
		bone.TailIndex = head.Index()
		bone.BoneFlag |= BONE_FLAG_TAIL_IS_BONE
	}

	return bone, nil
}

// GetHead 頭取得
func (bones *Bones) GetHead() (*Bone, error) {
	return bones.GetByName(HEAD.String())
}

// CreateHead 頭取得or作成
func (bones *Bones) CreateHead() (*Bone, error) {
	bone := NewBoneByName(HEAD.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if neck, err := bones.GetNeck(); err == nil {
		bone.Position = &mmath.MVec3{
			X: neck.Position.X,
			Y: neck.Position.Y * 1.1,
			Z: neck.Position.Z,
		}
	}

	// 親ボーン
	if neck, err := bones.GetNeck(); err == nil {
		bone.ParentIndex = neck.Index()
	}

	// 表示先位置
	bone.TailPosition = &mmath.MVec3{X: 0, Y: 0.5, Z: 0}

	return bone, nil
}

// GetHeadTail 頭先取得
func (bones *Bones) GetHeadTail() (*Bone, error) {
	return bones.GetByName(HEAD_TAIL.String())
}

// CreateHeadTail 頭先取得or作成
func (bones *Bones) CreateHeadTail() (*Bone, error) {
	bone := NewBoneByName(HEAD_TAIL.String())

	// 位置
	neck, _ := bones.GetNeck()
	head, _ := bones.GetHead()
	if neck != nil && head != nil {
		bone.Position = &mmath.MVec3{
			X: head.Position.X,
			Y: head.Position.Y + (head.Position.Y - neck.Position.Y),
			Z: head.Position.Z,
		}
	}

	// 親ボーン
	if head, err := bones.GetHead(); err == nil {
		bone.ParentIndex = head.Index()
	}

	return bone, nil
}

// GetEyes 両目取得
func (bones *Bones) GetEyes() (*Bone, error) {
	return bones.GetByName(EYES.String())
}

// CreateEyes 両目取得or作成
func (bones *Bones) CreateEyes() (*Bone, error) {
	bone := NewBoneByName(EYES.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if head, err := bones.GetHead(); err == nil {
		bone.Position = &mmath.MVec3{
			X: head.Position.X,
			Y: head.Position.Y + 0.1,
			Z: head.Position.Z,
		}
	}

	// 親ボーン
	if head, err := bones.GetHead(); err == nil {
		bone.ParentIndex = head.Index()
	}

	return bone, nil
}

// GetEye 目取得
func (bones *Bones) GetEye(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(EYE.StringFromDirection(direction.String()))
}

// CreateEye 目取得or作成
func (bones *Bones) CreateEye(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(EYE.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if eyes, err := bones.GetEyes(); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = &mmath.MVec3{
				X: eyes.Position.X - 0.1,
				Y: eyes.Position.Y,
				Z: eyes.Position.Z,
			}
		case BONE_DIRECTION_RIGHT:
			bone.Position = &mmath.MVec3{
				X: eyes.Position.X + 0.1,
				Y: eyes.Position.Y,
				Z: eyes.Position.Z,
			}
		}

		// 付与親
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
		bone.EffectFactor = 0.3
		bone.EffectIndex = eyes.Index()
	}

	// 親ボーン
	if eyes, err := bones.GetEyes(); err == nil {
		bone.ParentIndex = eyes.Index()
	}

	return bone, nil
}

// GetShoulderRoot 肩根元取得
func (bones *Bones) GetShoulderRoot(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_ROOT.StringFromDirection(direction.String()))
}

// CreateShoulderRoot 肩根元取得or作成
func (bones *Bones) CreateShoulderRoot(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_ROOT.StringFromDirection(direction.String()))

	// 位置
	shoulderLeft, _ := bones.GetShoulder(BONE_DIRECTION_LEFT)
	shoulderRight, _ := bones.GetShoulder(BONE_DIRECTION_RIGHT)
	if shoulderLeft != nil && shoulderRight != nil {
		bone.Position = &mmath.MVec3{
			X: (shoulderLeft.Position.X + shoulderRight.Position.X) * 0.5,
			Y: (shoulderLeft.Position.Y + shoulderRight.Position.Y) * 0.5,
			Z: (shoulderLeft.Position.Z + shoulderRight.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if neckRoot, err := bones.GetNeckRoot(); err == nil {
		bone.ParentIndex = neckRoot.Index()
	}

	return bone, nil
}

// GetShoulderP 肩P取得
func (bones *Bones) GetShoulderP(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_P.StringFromDirection(direction.String()))
}

// CreateShoulderP 肩P取得or作成
func (bones *Bones) CreateShoulderP(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_P.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if shoulder, err := bones.GetShoulder(direction); err == nil {
		bone.Position = shoulder.Position.Copy()
	}

	// 親ボーン
	if shoulderRoot, err := bones.GetShoulderRoot(direction); err == nil {
		bone.ParentIndex = shoulderRoot.Index()
	}

	return bone, nil
}

// GetShoulder 肩取得
func (bones *Bones) GetShoulder(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER.StringFromDirection(direction.String()))
}

// GetShoulderC 肩C取得
func (bones *Bones) GetShoulderC(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_C.StringFromDirection(direction.String()))
}

// CreateShoulderC 肩C取得or作成
func (bones *Bones) CreateShoulderC(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_C.StringFromDirection(direction.String()))

	// 位置
	if arm, err := bones.GetArm(direction); err == nil {
		bone.Position = arm.Position.Copy()
	}

	// 親ボーン
	if shoulder, err := bones.GetShoulder(direction); err == nil {
		bone.ParentIndex = shoulder.Index()
	}

	// 付与親
	if shoulderP, err := bones.GetShoulderP(direction); err == nil {
		bone.EffectIndex = shoulderP.Index()
		bone.EffectFactor = -1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetArm 腕取得
func (bones *Bones) GetArm(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ARM.StringFromDirection(direction.String()))
}

// GetArmTwist 腕捩取得
func (bones *Bones) GetArmTwist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ARM_TWIST.StringFromDirection(direction.String()))
}

// CreateArmTwist 腕捩取得or作成
func (bones *Bones) CreateArmTwist(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ARM_TWIST.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	arm, _ := bones.GetArm(direction)
	elbow, _ := bones.GetElbow(direction)
	if arm != nil && elbow != nil {
		bone.Position = &mmath.MVec3{
			X: (arm.Position.X + elbow.Position.X) * 0.5,
			Y: (arm.Position.Y + elbow.Position.Y) * 0.5,
			Z: (arm.Position.Z + elbow.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if arm != nil {
		bone.ParentIndex = arm.Index()
	}

	if elbow != nil {
		// 固定軸
		bone.FixedAxis = elbow.Position.Subed(bone.Position).Normalize()
		bone.BoneFlag |= BONE_FLAG_HAS_FIXED_AXIS

		// ローカル軸
		bone.LocalAxisX = elbow.Position.Subed(bone.Position).Normalize()
		bone.LocalAxisZ = mmath.MVec3UnitYNeg.Cross(bone.LocalAxisX).Normalize()
		bone.BoneFlag |= BONE_FLAG_HAS_LOCAL_AXIS
	}

	return bone, nil
}

// GetArmTwistChild 腕捩分割取得
func (bones *Bones) GetArmTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	return bones.GetByName(ARM_TWIST.StringFromDirectionAndIdx(direction.String(), idx+1))
}

// CreateArmTwistChild 腕捩分割取得or作成
func (bones *Bones) CreateArmTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	bone := NewBoneByName(ARM_TWIST.StringFromDirectionAndIdx(direction.String(), idx+1))

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
		bone.Position = &mmath.MVec3{
			X: arm.Position.X + ((elbow.Position.X - arm.Position.X) * ratio),
			Y: arm.Position.Y + ((elbow.Position.Y - arm.Position.Y) * ratio),
			Z: arm.Position.Z + ((elbow.Position.Z - arm.Position.Z) * ratio),
		}
	}

	// 親ボーン
	if arm != nil {
		bone.ParentIndex = arm.Index()
	}

	// 付与親
	if armTwist, err := bones.GetArmTwist(direction); err == nil {
		bone.EffectIndex = armTwist.Index()
		bone.EffectFactor = ratio
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetElbowRoot ひじ取得
func (bones *Bones) GetElbow(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ELBOW.StringFromDirection(direction.String()))
}

// GetWristTwist 腕捩取得
func (bones *Bones) GetWristTwist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST_TWIST.StringFromDirection(direction.String()))
}

// CreateWristTwist 腕捩取得or作成
func (bones *Bones) CreateWristTwist(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WRIST_TWIST.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	elbow, _ := bones.GetElbow(direction)
	wrist, _ := bones.GetWrist(direction)
	if elbow != nil && wrist != nil {
		bone.Position = &mmath.MVec3{
			X: (elbow.Position.X + wrist.Position.X) * 0.5,
			Y: (elbow.Position.Y + wrist.Position.Y) * 0.5,
			Z: (elbow.Position.Z + wrist.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if elbow != nil {
		bone.ParentIndex = elbow.Index()
	}

	if wrist != nil {
		// 固定軸
		bone.FixedAxis = wrist.Position.Subed(bone.Position).Normalize()
		bone.BoneFlag |= BONE_FLAG_HAS_FIXED_AXIS

		// ローカル軸
		bone.LocalAxisX = wrist.Position.Subed(bone.Position).Normalize()
		bone.LocalAxisZ = mmath.MVec3UnitYNeg.Cross(bone.LocalAxisX).Normalize()
		bone.BoneFlag |= BONE_FLAG_HAS_LOCAL_AXIS
	}

	return bone, nil
}

// GetWristTwistChild 腕捩分割取得
func (bones *Bones) GetWristTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	return bones.GetByName(WRIST_TWIST.StringFromDirectionAndIdx(direction.String(), idx+1))
}

// CreateWristTwistChild 腕捩分割取得or作成
func (bones *Bones) CreateWristTwistChild(direction BoneDirection, idx int) (*Bone, error) {
	bone := NewBoneByName(WRIST_TWIST.StringFromDirectionAndIdx(direction.String(), idx+1))

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
		bone.Position = &mmath.MVec3{
			X: elbow.Position.X + ((wrist.Position.X - elbow.Position.X) * ratio),
			Y: elbow.Position.Y + ((wrist.Position.Y - elbow.Position.Y) * ratio),
			Z: elbow.Position.Z + ((wrist.Position.Z - elbow.Position.Z) * ratio),
		}
	}

	// 親ボーン
	if elbow != nil {
		bone.ParentIndex = elbow.Index()
	}

	// 付与親
	if wristTwist, err := bones.GetWristTwist(direction); err == nil {
		bone.EffectIndex = wristTwist.Index()
		bone.EffectFactor = ratio
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetWrist 手首取得
func (bones *Bones) GetWrist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST.StringFromDirection(direction.String()))
}

// GetWristTail 手首先先取得
func (bones *Bones) GetWristTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST_TAIL.StringFromDirection(direction.String()))
}

// CreateWristTail 手首先先取得or作成
func (bones *Bones) CreateWristTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WRIST_TAIL.StringFromDirection(direction.String()))

	if wrist, err := bones.GetWrist(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = &mmath.MVec3{
				X: wrist.Position.X + 0.2,
				Y: wrist.Position.Y - 0.5,
				Z: wrist.Position.Z - 0.2,
			}
		case BONE_DIRECTION_RIGHT:
			bone.Position = &mmath.MVec3{
				X: wrist.Position.X - 0.2,
				Y: wrist.Position.Y - 0.5,
				Z: wrist.Position.Z - 0.2,
			}
		}
	}

	// 親ボーン
	if wrist, err := bones.GetWrist(direction); err == nil {
		bone.ParentIndex = wrist.Index()
	}

	return bone, nil
}

// GetThumb 親指取得
func (bones *Bones) GetThumb(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(THUMB0.StringFromDirection(direction.String()))
	case 1:
		return bones.GetByName(THUMB1.StringFromDirection(direction.String()))
	case 2:
		return bones.GetByName(THUMB2.StringFromDirection(direction.String()))
	}

	return nil, errors.New("invalid idx")
}

// CreateThumb0 親指0取得or作成
func (bones *Bones) CreateThumb0(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(THUMB0.StringFromDirection(direction.String()))

	// 位置
	wrist, _ := bones.GetWrist(direction)
	thumb1, _ := bones.GetThumb(direction, 1)
	if wrist != nil && thumb1 != nil {
		bone.Position = &mmath.MVec3{
			X: wrist.Position.X + (thumb1.Position.X-wrist.Position.X)*0.5,
			Y: wrist.Position.Y + (thumb1.Position.Y-wrist.Position.Y)*0.5,
			Z: wrist.Position.Z + (thumb1.Position.Z-wrist.Position.Z)*0.5,
		}
	}

	// 親ボーン
	if wrist, err := bones.GetWrist(direction); err == nil {
		bone.ParentIndex = wrist.Index()
	}

	return bone, nil
}

// GetThumbTail 親指先先取得
func (bones *Bones) GetThumbTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(THUMB_TAIL.StringFromDirection(direction.String()))
}

// CreateThumbTail 親指先先取得or作成
func (bones *Bones) CreateThumbTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(THUMB_TAIL.StringFromDirection(direction.String()))

	// 位置
	thumb1, _ := bones.GetThumb(direction, 1)
	thumb2, _ := bones.GetThumb(direction, 2)
	if thumb1 != nil && thumb2 != nil {
		bone.Position = &mmath.MVec3{
			X: thumb2.Position.X + (thumb2.Position.X - thumb2.Position.X),
			Y: thumb2.Position.Y + (thumb2.Position.Y - thumb2.Position.Y),
			Z: thumb2.Position.Z + (thumb2.Position.Z - thumb2.Position.Z),
		}
	}

	// 親ボーン
	if thumb2, err := bones.GetThumb(direction, 2); err == nil {
		bone.ParentIndex = thumb2.Index()
	}

	return bone, nil
}

// GetIndex 人差し指取得
func (bones *Bones) GetIndex(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(INDEX1.StringFromDirection(direction.String()))
	case 1:
		return bones.GetByName(INDEX2.StringFromDirection(direction.String()))
	case 2:
		return bones.GetByName(INDEX3.StringFromDirection(direction.String()))
	}

	return nil, errors.New("invalid idx")
}

// GetIndexTail 人差し指先先取得
func (bones *Bones) GetIndexTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(INDEX_TAIL.StringFromDirection(direction.String()))
}

// CreateIndexTail 親指先先取得or作成
func (bones *Bones) CreateIndexTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(INDEX_TAIL.StringFromDirection(direction.String()))

	// 位置
	index1, _ := bones.GetIndex(direction, 1)
	index2, _ := bones.GetIndex(direction, 2)
	if index1 != nil && index2 != nil {
		bone.Position = &mmath.MVec3{
			X: index2.Position.X + (index2.Position.X - index2.Position.X),
			Y: index2.Position.Y + (index2.Position.Y - index2.Position.Y),
			Z: index2.Position.Z + (index2.Position.Z - index2.Position.Z),
		}
	}

	// 親ボーン
	if index2, err := bones.GetIndex(direction, 2); err == nil {
		bone.ParentIndex = index2.Index()
	}

	return bone, nil
}

// GetMiddle 中指取得
func (bones *Bones) GetMiddle(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(MIDDLE1.StringFromDirection(direction.String()))
	case 1:
		return bones.GetByName(MIDDLE2.StringFromDirection(direction.String()))
	case 2:
		return bones.GetByName(MIDDLE3.StringFromDirection(direction.String()))
	}

	return nil, errors.New("invalid idx")
}

// GetMiddleTail 中指先先取得
func (bones *Bones) GetMiddleTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(MIDDLE_TAIL.StringFromDirection(direction.String()))
}

// CreateMiddleTail 親指先先取得or作成
func (bones *Bones) CreateMiddleTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(MIDDLE_TAIL.StringFromDirection(direction.String()))

	// 位置
	middle1, _ := bones.GetMiddle(direction, 1)
	middle2, _ := bones.GetMiddle(direction, 2)
	if middle1 != nil && middle2 != nil {
		bone.Position = &mmath.MVec3{
			X: middle2.Position.X + (middle2.Position.X - middle2.Position.X),
			Y: middle2.Position.Y + (middle2.Position.Y - middle2.Position.Y),
			Z: middle2.Position.Z + (middle2.Position.Z - middle2.Position.Z),
		}
	}

	// 親ボーン
	if middle2, err := bones.GetMiddle(direction, 2); err == nil {
		bone.ParentIndex = middle2.Index()
	}

	return bone, nil
}

// GetRing 薬指取得
func (bones *Bones) GetRing(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(RING1.StringFromDirection(direction.String()))
	case 1:
		return bones.GetByName(RING2.StringFromDirection(direction.String()))
	case 2:
		return bones.GetByName(RING3.StringFromDirection(direction.String()))
	}

	return nil, errors.New("invalid idx")
}

// GetRingTail 薬指先先取得
func (bones *Bones) GetRingTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(RING_TAIL.StringFromDirection(direction.String()))
}

// CreateRingTail 親指先先取得or作成
func (bones *Bones) CreateRingTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(RING_TAIL.StringFromDirection(direction.String()))

	// 位置
	ring1, _ := bones.GetRing(direction, 1)
	ring2, _ := bones.GetRing(direction, 2)
	if ring1 != nil && ring2 != nil {
		bone.Position = &mmath.MVec3{
			X: ring2.Position.X + (ring2.Position.X - ring2.Position.X),
			Y: ring2.Position.Y + (ring2.Position.Y - ring2.Position.Y),
			Z: ring2.Position.Z + (ring2.Position.Z - ring2.Position.Z),
		}
	}

	// 親ボーン
	if ring2, err := bones.GetRing(direction, 2); err == nil {
		bone.ParentIndex = ring2.Index()
	}

	return bone, nil
}

// GetPinky 小指取得
func (bones *Bones) GetPinky(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(PINKY1.StringFromDirection(direction.String()))
	case 1:
		return bones.GetByName(PINKY2.StringFromDirection(direction.String()))
	case 2:
		return bones.GetByName(PINKY3.StringFromDirection(direction.String()))
	}

	return nil, errors.New("invalid idx")
}

// GetPinkyTail 小指先先取得
func (bones *Bones) GetPinkyTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(PINKY_TAIL.StringFromDirection(direction.String()))
}

// CreateLittleTail 親指先先取得or作成
func (bones *Bones) CreatePinkyTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(PINKY_TAIL.StringFromDirection(direction.String()))

	// 位置
	pinky1, _ := bones.GetPinky(direction, 1)
	pinky2, _ := bones.GetPinky(direction, 2)
	if pinky1 != nil && pinky2 != nil {
		bone.Position = &mmath.MVec3{
			X: pinky2.Position.X + (pinky2.Position.X - pinky2.Position.X),
			Y: pinky2.Position.Y + (pinky2.Position.Y - pinky2.Position.Y),
			Z: pinky2.Position.Z + (pinky2.Position.Z - pinky2.Position.Z),
		}
	}

	// 親ボーン
	if pinky2, err := bones.GetPinky(direction, 2); err == nil {
		bone.ParentIndex = pinky2.Index()
	}

	return bone, nil
}

// GetLegRoot 足根元取得
func (bones *Bones) GetLegRoot(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_ROOT.StringFromDirection(direction.String()))
}

// CreateLegRoot 足根元取得or作成
func (bones *Bones) CreateLegRoot(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_ROOT.StringFromDirection(direction.String()))

	// 位置
	legLeft, _ := bones.GetLeg(BONE_DIRECTION_LEFT)
	legRight, _ := bones.GetLeg(BONE_DIRECTION_RIGHT)
	if legLeft != nil && legRight != nil {
		bone.Position = &mmath.MVec3{
			X: (legLeft.Position.X + legRight.Position.X) * 0.5,
			Y: (legLeft.Position.Y + legRight.Position.Y) * 0.5,
			Z: (legLeft.Position.Z + legRight.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if legCenter, err := bones.GetLegCenter(); err == nil {
		bone.ParentIndex = legCenter.Index()
	} else if lower, err := bones.GetLower(); err == nil {
		bone.ParentIndex = lower.Index()
	}

	return bone, nil
}

// GetWaistCancel 腰キャンセル取得
func (bones *Bones) GetWaistCancel(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WAIST_CANCEL.StringFromDirection(direction.String()))
}

// CreateWaistCancel 腰キャンセル取得or作成
func (bones *Bones) CreateWaistCancel(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WAIST_CANCEL.StringFromDirection(direction.String()))

	// 位置
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.Position = leg.Position.Copy()
	}

	// 親ボーン
	if legRoot, err := bones.GetLegRoot(direction); err == nil {
		bone.ParentIndex = legRoot.Index()
	} else if legCenter, err := bones.GetLegCenter(); err == nil {
		bone.ParentIndex = legCenter.Index()
	} else if lower, err := bones.GetLower(); err == nil {
		bone.ParentIndex = lower.Index()
	}

	// 付与親
	if waist, err := bones.GetWaist(); err == nil {
		bone.EffectIndex = waist.Index()
		bone.EffectFactor = -1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetLeg 足取得
func (bones *Bones) GetLeg(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG.StringFromDirection(direction.String()))
}

// GetKnee ひざ取得
func (bones *Bones) GetKnee(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(KNEE.StringFromDirection(direction.String()))
}

// GetAnkle 足首取得
func (bones *Bones) GetAnkle(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE.StringFromDirection(direction.String()))
}

// GetHeel かかと取得
func (bones *Bones) GetHeel(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HEEL.StringFromDirection(direction.String()))
}

// CreateHeel かかと取得or作成
func (bones *Bones) CreateHeel(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HEEL.StringFromDirection(direction.String()))

	// 位置
	// 親ボーン
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.Position = &mmath.MVec3{
			X: ankle.Position.X,
			Y: 0.0,
			Z: ankle.Position.Z + 0.2,
		}

		bone.ParentIndex = ankle.Index()
	}

	return bone, nil
}

// GetToeT つま先先取得
func (bones *Bones) GetToeT(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_T.StringFromDirection(direction.String()))
}

// CreateToeT つま先先取得or作成
func (bones *Bones) CreateToeT(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_T.StringFromDirection(direction.String()))

	// 位置
	if toeIK, err := bones.GetToeIK(direction); err == nil {
		if toe, err := bones.Get(toeIK.Ik.BoneIndex); err == nil {
			// つま先IKのターゲットと同位置
			bone.Position = toe.Position.Copy()

			// 親はつま先IKのターゲット
			bone.ParentIndex = toe.Index()
		}
	}

	return bone, nil
}

// GetToeP つま先親取得
func (bones *Bones) GetToeP(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_P.StringFromDirection(direction.String()))
}

// CreateToeP つま先親取得or作成
func (bones *Bones) CreateToeP(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_P.StringFromDirection(direction.String()))

	// 位置
	// 親ボーン
	if toeT, err := bones.GetToeT(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X - 0.3,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		case BONE_DIRECTION_RIGHT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X + 0.3,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		}

		bone.ParentIndex = toeT.Index()
	}

	return bone, nil
}

// GetToeC つま先子取得
func (bones *Bones) GetToeC(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_C.StringFromDirection(direction.String()))
}

// CreateToeC つま先子取得or作成
func (bones *Bones) CreateToeC(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_C.StringFromDirection(direction.String()))

	// 位置
	if toeT, err := bones.GetToeT(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X + 0.3,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		case BONE_DIRECTION_RIGHT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X - 0.3,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		}
	}

	// 親ボーン
	if toeP, err := bones.GetToeP(direction); err == nil {
		bone.ParentIndex = toeP.Index()
	}

	return bone, nil
}

// GetLegD 足D取得
func (bones *Bones) GetLegD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_D.StringFromDirection(direction.String()))
}

// CreateLegD 足D取得or作成
func (bones *Bones) CreateLegD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_D.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.Position = leg.Position.Copy()
	}

	// 親ボーン
	if waistCancel, err := bones.GetWaistCancel(direction); err == nil {
		bone.ParentIndex = waistCancel.Index()
	} else if legRoot, err := bones.GetLegRoot(direction); err == nil {
		bone.ParentIndex = legRoot.Index()
	} else if legCenter, err := bones.GetLegCenter(); err == nil {
		bone.ParentIndex = legCenter.Index()
	} else if lower, err := bones.GetLower(); err == nil {
		bone.ParentIndex = lower.Index()
	}

	// 付与親
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.EffectIndex = leg.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetKneeD ひざD取得
func (bones *Bones) GetKneeD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(KNEE_D.StringFromDirection(direction.String()))
}

// CreateKneeD ひざD取得or作成
func (bones *Bones) CreateKneeD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(KNEE_D.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if knee, err := bones.GetKnee(direction); err == nil {
		bone.Position = knee.Position.Copy()
	}

	// 親ボーン
	if legD, err := bones.GetLegD(direction); err == nil {
		bone.ParentIndex = legD.Index()
	}

	// 付与親
	if knee, err := bones.GetKnee(direction); err == nil {
		bone.EffectIndex = knee.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetAnkleD 足首D取得
func (bones *Bones) GetAnkleD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE_D.StringFromDirection(direction.String()))
}

// CreateAnkleD 足首D取得or作成
func (bones *Bones) CreateAnkleD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ANKLE_D.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.Position = ankle.Position.Copy()
	}

	// 親ボーン
	if kneeD, err := bones.GetKneeD(direction); err == nil {
		bone.ParentIndex = kneeD.Index()
	}

	// 付与親
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.EffectIndex = ankle.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetHeelD かかとD取得
func (bones *Bones) GetHeelD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HEEL_D.StringFromDirection(direction.String()))
}

// CreateHeelD かかとD取得or作成
func (bones *Bones) CreateHeelD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HEEL_D.StringFromDirection(direction.String()))

	// 位置
	if heel, err := bones.GetHeel(direction); err == nil {
		bone.Position = heel.Position.Copy()
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
	}

	return bone, nil
}

// GetToeEx 足先EX取得
func (bones *Bones) GetToeEx(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_EX.StringFromDirection(direction.String()))
}

// CreateToeEx 足先EX取得or作成
func (bones *Bones) CreateToeEx(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_EX.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	ankle, _ := bones.GetAnkle(direction)
	toeT, _ := bones.GetToeT(direction)
	if ankle != nil && toeT != nil {
		bone.Position = &mmath.MVec3{
			X: (ankle.Position.X + toeT.Position.X) * 0.5,
			Y: (ankle.Position.Y + toeT.Position.Y) * 0.5,
			Z: (ankle.Position.Z + toeT.Position.Z) * 0.5,
		}
	}

	// 親ボーン
	if ankleD, err := bones.GetAnkleD(direction); err == nil {
		bone.ParentIndex = ankleD.Index()
	}

	return bone, nil
}

// GetToeTD つま先先D取得
func (bones *Bones) GetToeTD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_T_D.StringFromDirection(direction.String()))
}

// CreateToeTD つま先先D取得or作成
func (bones *Bones) CreateToeTD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_T_D.StringFromDirection(direction.String()))

	// 位置
	if toeT, err := bones.GetToeT(direction); err == nil {
		bone.Position = toeT.Position.Copy()
	}

	// 親ボーン
	if toeEx, err := bones.GetToeEx(direction); err == nil {
		bone.ParentIndex = toeEx.Index()
	}

	// 付与親
	if toeT, err := bones.GetToeT(direction); err == nil {
		bone.EffectIndex = toeT.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetToePD つま先親D取得
func (bones *Bones) GetToePD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_P_D.StringFromDirection(direction.String()))
}

// CreateToePD つま先親D取得or作成
func (bones *Bones) CreateToePD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_P_D.StringFromDirection(direction.String()))

	// 位置
	if toeP, err := bones.GetToeP(direction); err == nil {
		bone.Position = toeP.Position.Copy()
	}

	// 親ボーン
	if toeTD, err := bones.GetToeTD(direction); err == nil {
		bone.ParentIndex = toeTD.Index()
	}

	// 付与親
	if toeP, err := bones.GetToeP(direction); err == nil {
		bone.EffectIndex = toeP.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetToeCD つま先子D取得
func (bones *Bones) GetToeCD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_C_D.StringFromDirection(direction.String()))
}

// CreateToeCD つま先子D取得or作成
func (bones *Bones) CreateToeCD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_C_D.StringFromDirection(direction.String()))

	// 位置
	if toeC, err := bones.GetToeC(direction); err == nil {
		bone.Position = toeC.Position.Copy()
	}

	// 親ボーン
	if toePD, err := bones.GetToePD(direction); err == nil {
		bone.ParentIndex = toePD.Index()
	}

	// 付与親
	if toeC, err := bones.GetToeC(direction); err == nil {
		bone.EffectIndex = toeC.Index()
		bone.EffectFactor = 1.0
		bone.BoneFlag |= BONE_FLAG_IS_EXTERNAL_ROTATION
	}

	return bone, nil
}

// GetLegIkParent 足IK親取得
func (bones *Bones) GetLegIkParent(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_IK_PARENT.StringFromDirection(direction.String()))
}

// CreateLegIkParent 足IK親取得or作成
func (bones *Bones) CreateLegIkParent(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_IK_PARENT.StringFromDirection(direction.String()))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if legIk, err := bones.GetLegIk(direction); err == nil {
		bone.Position = &mmath.MVec3{
			X: legIk.Position.X,
			Y: 0.0,
			Z: legIk.Position.Z,
		}
	}

	// 親ボーン
	if root, err := bones.GetRoot(); err == nil {
		bone.ParentIndex = root.Index()
	}

	return bone, nil
}

// GetLegIk 足IK取得
func (bones *Bones) GetLegIk(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_IK.StringFromDirection(direction.String()))
}

// GetToeIK つま先IK取得
func (bones *Bones) GetToeIK(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_IK.StringFromDirection(direction.String()))
}
