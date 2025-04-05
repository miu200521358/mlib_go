package pmx

import (
	"errors"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/merr"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// GetRoot すべての親
func (bones *Bones) GetRoot() (*Bone, error) {
	return bones.GetByName(ROOT.String())
}

// CreateRoot すべての親作成
func (bones *Bones) CreateRoot() (*Bone, error) {
	bone := NewBoneByName(ROOT.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE
	return bone, nil
}

// GetCenter センター取得
func (bones *Bones) GetCenter() (*Bone, error) {
	return bones.GetByName(CENTER.String())
}

func (bones *Bones) CreateCenter() (*Bone, error) {
	bone := NewBoneByName(CENTER.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if upper, err := bones.GetUpper(); err == nil {
		bone.Position.Y = upper.Position.Y * 0.6
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(ROOT, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetGroove グルーブ取得
func (bones *Bones) GetGroove() (*Bone, error) {
	return bones.GetByName(GROOVE.String())
}

// CreateGroove グルーブ作成
func (bones *Bones) CreateGroove() (*Bone, error) {
	bone := NewBoneByName(GROOVE.String())
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE

	// 位置
	if upper, err := bones.GetUpper(); err == nil {
		bone.Position.Y = upper.Position.Y * 0.7
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(GROOVE, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetWaist 腰取得
func (bones *Bones) GetWaist() (*Bone, error) {
	return bones.GetByName(WAIST.String())
}

// CreateWaist 腰作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(WAIST, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetTrunkRoot 体幹中心取得
func (bones *Bones) GetTrunkRoot() (*Bone, error) {
	return bones.GetByName(TRUNK_ROOT.String())
}

// CreateTrunkRoot 体幹中心作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(TRUNK_ROOT, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetLowerRoot 下半身根元取得
func (bones *Bones) GetLowerRoot() (*Bone, error) {
	return bones.GetByName(LOWER_ROOT.String())
}

// CreateLowerRoot 下半身根元作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(LOWER_ROOT, BONE_DIRECTION_TRUNK)

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

// CreateLegCenter 足中心作成
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

// CreateUpperRoot 上半身根元作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(UPPER_ROOT, BONE_DIRECTION_TRUNK)

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

// CreateUpper2 上半身2作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(UPPER2, BONE_DIRECTION_TRUNK)

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

// CreateNeckRoot 首根元作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(NECK_ROOT, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetNeck 首取得
func (bones *Bones) GetNeck() (*Bone, error) {
	return bones.GetByName(NECK.String())
}

// CreateNeck 首作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(NECK, BONE_DIRECTION_TRUNK)

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

// CreateHead 頭作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(HEAD, BONE_DIRECTION_TRUNK)

	// 表示先位置
	bone.TailPosition = &mmath.MVec3{X: 0, Y: 0.5, Z: 0}

	return bone, nil
}

// GetHeadTail 頭先取得
func (bones *Bones) GetHeadTail() (*Bone, error) {
	return bones.GetByName(HEAD_TAIL.String())
}

// CreateHeadTail 頭先作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(HEAD_TAIL, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetEyes 両目取得
func (bones *Bones) GetEyes() (*Bone, error) {
	return bones.GetByName(EYES.String())
}

// CreateEyes 両目作成
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
	bone.ParentIndex = bones.findParentIndexByConfig(EYES, BONE_DIRECTION_TRUNK)

	return bone, nil
}

// GetEye 目取得
func (bones *Bones) GetEye(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(EYE.StringFromDirection(direction))
}

// CreateEye 目作成
func (bones *Bones) CreateEye(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(EYE.StringFromDirection(direction))
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
	bone.ParentIndex = bones.findParentIndexByConfig(EYE, direction)

	return bone, nil
}

// GetShoulderRoot 肩根元取得
func (bones *Bones) GetShoulderRoot(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_ROOT.StringFromDirection(direction))
}

// CreateShoulderRoot 肩根元作成
func (bones *Bones) CreateShoulderRoot(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_ROOT.StringFromDirection(direction))

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
	return bones.GetByName(SHOULDER_P.StringFromDirection(direction))
}

// CreateShoulderP 肩P作成
func (bones *Bones) CreateShoulderP(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_P.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if shoulder, err := bones.GetShoulder(direction); err == nil {
		bone.Position = shoulder.Position.Copy()
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(SHOULDER_P, direction)

	return bone, nil
}

// GetShoulder 肩取得
func (bones *Bones) GetShoulder(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER.StringFromDirection(direction))
}

// GetShoulderC 肩C取得
func (bones *Bones) GetShoulderC(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(SHOULDER_C.StringFromDirection(direction))
}

// CreateShoulderC 肩C作成
func (bones *Bones) CreateShoulderC(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(SHOULDER_C.StringFromDirection(direction))

	// 位置
	if arm, err := bones.GetArm(direction); err == nil {
		bone.Position = arm.Position.Copy()
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
func (bones *Bones) GetArm(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ARM.StringFromDirection(direction))
}

// GetArmTwist 腕捩取得
func (bones *Bones) GetArmTwist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ARM_TWIST.StringFromDirection(direction))
}

// CreateArmTwist 腕捩作成
func (bones *Bones) CreateArmTwist(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ARM_TWIST.StringFromDirection(direction))
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
	bone.ParentIndex = bones.findParentIndexByConfig(ARM_TWIST, direction)

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
	return bones.GetByName(ARM_TWIST.StringFromDirectionAndIdx(direction, idx+1))
}

// CreateArmTwistChild 腕捩分割作成
func (bones *Bones) CreateArmTwistChild(direction BoneDirection, idx int) (*Bone, error) {
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
		bone.Position = &mmath.MVec3{
			X: arm.Position.X + ((elbow.Position.X - arm.Position.X) * ratio),
			Y: arm.Position.Y + ((elbow.Position.Y - arm.Position.Y) * ratio),
			Z: arm.Position.Z + ((elbow.Position.Z - arm.Position.Z) * ratio),
		}
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
func (bones *Bones) GetElbow(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ELBOW.StringFromDirection(direction))
}

// GetWristTwist 腕捩取得
func (bones *Bones) GetWristTwist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST_TWIST.StringFromDirection(direction))
}

// CreateWristTwist 腕捩作成
func (bones *Bones) CreateWristTwist(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WRIST_TWIST.StringFromDirection(direction))
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
	bone.ParentIndex = bones.findParentIndexByConfig(WRIST_TWIST, direction)

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
	return bones.GetByName(WRIST_TWIST.StringFromDirectionAndIdx(direction, idx+1))
}

// CreateWristTwistChild 腕捩分割作成
func (bones *Bones) CreateWristTwistChild(direction BoneDirection, idx int) (*Bone, error) {
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
		bone.Position = &mmath.MVec3{
			X: elbow.Position.X + ((wrist.Position.X - elbow.Position.X) * ratio),
			Y: elbow.Position.Y + ((wrist.Position.Y - elbow.Position.Y) * ratio),
			Z: elbow.Position.Z + ((wrist.Position.Z - elbow.Position.Z) * ratio),
		}
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
func (bones *Bones) GetWrist(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST.StringFromDirection(direction))
}

// GetWristTail 手首先先取得
func (bones *Bones) GetWristTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WRIST_TAIL.StringFromDirection(direction))
}

// CreateWristTail 手首先先作成
func (bones *Bones) CreateWristTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WRIST_TAIL.StringFromDirection(direction))

	switch direction {
	case BONE_DIRECTION_LEFT:
		bonePositions := make([]*mmath.MVec3, 0)
		for _, boneName := range []string{THUMB1.Left(), INDEX1.Left(), MIDDLE1.Left(), RING1.Left(), PINKY1.Left()} {
			if bone, err := bones.GetByName(boneName); err == nil {
				bonePositions = append(bonePositions, bone.Position)
			}
		}
		bone.Position = mmath.MeanVec3(bonePositions)
	case BONE_DIRECTION_RIGHT:
		bonePositions := make([]*mmath.MVec3, 0)
		for _, boneName := range []string{THUMB1.Right(), INDEX1.Right(), MIDDLE1.Right(), RING1.Right(), PINKY1.Right()} {
			if bone, err := bones.GetByName(boneName); err == nil {
				bonePositions = append(bonePositions, bone.Position)
			}
		}
		bone.Position = mmath.MeanVec3(bonePositions)
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WRIST_TAIL, direction)

	return bone, nil
}

// GetThumb 親指取得
func (bones *Bones) GetThumb(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(THUMB0.StringFromDirection(direction))
	case 1:
		return bones.GetByName(THUMB1.StringFromDirection(direction))
	case 2:
		return bones.GetByName(THUMB2.StringFromDirection(direction))
	}

	return nil, errors.New("invalid idx")
}

// CreateThumb0 親指0作成
func (bones *Bones) CreateThumb0(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(THUMB0.StringFromDirection(direction))

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
	bone.ParentIndex = bones.findParentIndexByConfig(THUMB0, direction)

	return bone, nil
}

// GetThumbTail 親指先先取得
func (bones *Bones) GetThumbTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(THUMB_TAIL.StringFromDirection(direction))
}

// CreateThumbTail 親指先先作成
func (bones *Bones) CreateThumbTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(THUMB_TAIL.StringFromDirection(direction))

	// 位置
	thumb1, _ := bones.GetThumb(direction, 1)
	thumb2, _ := bones.GetThumb(direction, 2)
	if thumb1 != nil && thumb2 != nil {
		bone.Position = &mmath.MVec3{
			X: thumb2.Position.X + (thumb2.Position.X - thumb1.Position.X),
			Y: thumb2.Position.Y + (thumb2.Position.Y - thumb1.Position.Y),
			Z: thumb2.Position.Z + (thumb2.Position.Z - thumb1.Position.Z),
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(THUMB_TAIL, direction)

	return bone, nil
}

// GetIndex 人差し指取得
func (bones *Bones) GetIndex(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(INDEX1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(INDEX2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(INDEX3.StringFromDirection(direction))
	}

	return nil, errors.New("invalid idx")
}

// GetIndexTail 人差し指先先取得
func (bones *Bones) GetIndexTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(INDEX_TAIL.StringFromDirection(direction))
}

// CreateIndexTail 親指先先作成
func (bones *Bones) CreateIndexTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(INDEX_TAIL.StringFromDirection(direction))

	// 位置
	index1, _ := bones.GetIndex(direction, 1)
	index2, _ := bones.GetIndex(direction, 2)
	if index1 != nil && index2 != nil {
		bone.Position = &mmath.MVec3{
			X: index2.Position.X + (index2.Position.X - index1.Position.X),
			Y: index2.Position.Y + (index2.Position.Y - index1.Position.Y),
			Z: index2.Position.Z + (index2.Position.Z - index1.Position.Z),
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(INDEX_TAIL, direction)

	return bone, nil
}

// GetMiddle 中指取得
func (bones *Bones) GetMiddle(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(MIDDLE1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(MIDDLE2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(MIDDLE3.StringFromDirection(direction))
	}

	return nil, errors.New("invalid idx")
}

// GetMiddleTail 中指先先取得
func (bones *Bones) GetMiddleTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(MIDDLE_TAIL.StringFromDirection(direction))
}

// CreateMiddleTail 親指先先作成
func (bones *Bones) CreateMiddleTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(MIDDLE_TAIL.StringFromDirection(direction))

	// 位置
	middle1, _ := bones.GetMiddle(direction, 1)
	middle2, _ := bones.GetMiddle(direction, 2)
	if middle1 != nil && middle2 != nil {
		bone.Position = &mmath.MVec3{
			X: middle2.Position.X + (middle2.Position.X - middle1.Position.X),
			Y: middle2.Position.Y + (middle2.Position.Y - middle1.Position.Y),
			Z: middle2.Position.Z + (middle2.Position.Z - middle1.Position.Z),
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(MIDDLE_TAIL, direction)

	return bone, nil
}

// GetRing 薬指取得
func (bones *Bones) GetRing(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(RING1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(RING2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(RING3.StringFromDirection(direction))
	}

	return nil, errors.New("invalid idx")
}

// GetRingTail 薬指先先取得
func (bones *Bones) GetRingTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(RING_TAIL.StringFromDirection(direction))
}

// CreateRingTail 親指先先作成
func (bones *Bones) CreateRingTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(RING_TAIL.StringFromDirection(direction))

	// 位置
	ring1, _ := bones.GetRing(direction, 1)
	ring2, _ := bones.GetRing(direction, 2)
	if ring1 != nil && ring2 != nil {
		bone.Position = &mmath.MVec3{
			X: ring2.Position.X + (ring2.Position.X - ring1.Position.X),
			Y: ring2.Position.Y + (ring2.Position.Y - ring1.Position.Y),
			Z: ring2.Position.Z + (ring2.Position.Z - ring1.Position.Z),
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(RING_TAIL, direction)

	return bone, nil
}

// GetPinky 小指取得
func (bones *Bones) GetPinky(direction BoneDirection, idx int) (*Bone, error) {
	switch idx {
	case 0:
		return bones.GetByName(PINKY1.StringFromDirection(direction))
	case 1:
		return bones.GetByName(PINKY2.StringFromDirection(direction))
	case 2:
		return bones.GetByName(PINKY3.StringFromDirection(direction))
	}

	return nil, errors.New("invalid idx")
}

// GetPinkyTail 小指先先取得
func (bones *Bones) GetPinkyTail(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(PINKY_TAIL.StringFromDirection(direction))
}

// CreateLittleTail 親指先先作成
func (bones *Bones) CreatePinkyTail(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(PINKY_TAIL.StringFromDirection(direction))

	// 位置
	pinky1, _ := bones.GetPinky(direction, 1)
	pinky2, _ := bones.GetPinky(direction, 2)
	if pinky1 != nil && pinky2 != nil {
		bone.Position = &mmath.MVec3{
			X: pinky2.Position.X + (pinky2.Position.X - pinky1.Position.X),
			Y: pinky2.Position.Y + (pinky2.Position.Y - pinky1.Position.Y),
			Z: pinky2.Position.Z + (pinky2.Position.Z - pinky1.Position.Z),
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(PINKY_TAIL, direction)

	return bone, nil
}

// GetLegRoot 足根元取得
func (bones *Bones) GetLegRoot(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_ROOT.StringFromDirection(direction))
}

// CreateLegRoot 足根元作成
func (bones *Bones) CreateLegRoot(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_ROOT.StringFromDirection(direction))

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
	bone.ParentIndex = bones.findParentIndexByConfig(LEG_ROOT, direction)

	return bone, nil
}

// GetWaistCancel 腰キャンセル取得
func (bones *Bones) GetWaistCancel(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(WAIST_CANCEL.StringFromDirection(direction))
}

// CreateWaistCancel 腰キャンセル作成
func (bones *Bones) CreateWaistCancel(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(WAIST_CANCEL.StringFromDirection(direction))

	// 位置
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.Position = leg.Position.Copy()
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(WAIST_CANCEL, direction)

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
	return bones.GetByName(LEG.StringFromDirection(direction))
}

// GetKnee ひざ取得
func (bones *Bones) GetKnee(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(KNEE.StringFromDirection(direction))
}

// GetAnkle 足首取得
func (bones *Bones) GetAnkle(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE.StringFromDirection(direction))
}

func (bones *Bones) GetToeIk(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_IK.StringFromDirection(direction))
}

// GetHeel かかと取得
func (bones *Bones) GetHeel(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HEEL.StringFromDirection(direction))
}

// CreateHeel かかと作成
func (bones *Bones) CreateHeel(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HEEL.StringFromDirection(direction))

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
	return bones.GetByName(TOE_T.StringFromDirection(direction))
}

// CreateToeT つま先先作成
func (bones *Bones) CreateToeT(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_T.StringFromDirection(direction))

	// 位置
	if toeIK, err := bones.GetToeIK(direction); err == nil {
		if toe, err := bones.Get(toeIK.Ik.BoneIndex); err == nil {
			// つま先IKのターゲットと同位置
			bone.Position = toe.Position.Copy()
		}
	}

	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.ParentIndex = ankle.Index()
	}

	return bone, nil
}

// GetToeP つま先親取得
func (bones *Bones) GetToeP(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_P.StringFromDirection(direction))
}

// CreateToeP つま先親作成
func (bones *Bones) CreateToeP(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_P.StringFromDirection(direction))

	// 位置
	// 親ボーン
	if toeT, err := bones.GetToeT(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X - 0.8,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		case BONE_DIRECTION_RIGHT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X + 0.8,
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
	return bones.GetByName(TOE_C.StringFromDirection(direction))
}

// CreateToeC つま先子作成
func (bones *Bones) CreateToeC(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_C.StringFromDirection(direction))

	// 位置
	if toeT, err := bones.GetToeT(direction); err == nil {
		switch direction {
		case BONE_DIRECTION_LEFT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X + 0.8,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		case BONE_DIRECTION_RIGHT:
			bone.Position = &mmath.MVec3{
				X: toeT.Position.X - 0.8,
				Y: toeT.Position.Y,
				Z: toeT.Position.Z,
			}
		}
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_C, direction)

	return bone, nil
}

// GetLegD 足D取得
func (bones *Bones) GetLegD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_D.StringFromDirection(direction))
}

// CreateLegD 足D作成
func (bones *Bones) CreateLegD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if leg, err := bones.GetLeg(direction); err == nil {
		bone.Position = leg.Position.Copy()
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
func (bones *Bones) GetKneeD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(KNEE_D.StringFromDirection(direction))
}

// CreateKneeD ひざD作成
func (bones *Bones) CreateKneeD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(KNEE_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if knee, err := bones.GetKnee(direction); err == nil {
		bone.Position = knee.Position.Copy()
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
func (bones *Bones) GetAnkleD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(ANKLE_D.StringFromDirection(direction))
}

// CreateAnkleD 足首D作成
func (bones *Bones) CreateAnkleD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(ANKLE_D.StringFromDirection(direction))
	bone.BoneFlag = BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE | BONE_FLAG_CAN_ROTATE

	// 位置
	if ankle, err := bones.GetAnkle(direction); err == nil {
		bone.Position = ankle.Position.Copy()
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

// GetHeelD かかとD取得
func (bones *Bones) GetHeelD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(HEEL_D.StringFromDirection(direction))
}

// CreateHeelD かかとD作成
func (bones *Bones) CreateHeelD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(HEEL_D.StringFromDirection(direction))

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
	return bones.GetByName(TOE_EX.StringFromDirection(direction))
}

// CreateToeEx 足先EX作成
func (bones *Bones) CreateToeEx(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_EX.StringFromDirection(direction))
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
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_EX, direction)

	return bone, nil
}

// GetToeTD つま先先D取得
func (bones *Bones) GetToeTD(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_T_D.StringFromDirection(direction))
}

// CreateToeTD つま先先D作成
func (bones *Bones) CreateToeTD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_T_D.StringFromDirection(direction))

	// 位置
	if toeT, err := bones.GetToeT(direction); err == nil {
		bone.Position = toeT.Position.Copy()
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_T_D, direction)

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
	return bones.GetByName(TOE_P_D.StringFromDirection(direction))
}

// CreateToePD つま先親D作成
func (bones *Bones) CreateToePD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_P_D.StringFromDirection(direction))

	// 位置
	if toeP, err := bones.GetToeP(direction); err == nil {
		bone.Position = toeP.Position.Copy()
	}

	// 親ボーン
	bone.ParentIndex = bones.findParentIndexByConfig(TOE_P_D, direction)

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
	return bones.GetByName(TOE_C_D.StringFromDirection(direction))
}

// CreateToeCD つま先子D作成
func (bones *Bones) CreateToeCD(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(TOE_C_D.StringFromDirection(direction))

	// 位置
	if toeC, err := bones.GetToeC(direction); err == nil {
		bone.Position = toeC.Position.Copy()
	}

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
func (bones *Bones) GetLegIkParent(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_IK_PARENT.StringFromDirection(direction))
}

// CreateLegIkParent 足IK親作成
func (bones *Bones) CreateLegIkParent(direction BoneDirection) (*Bone, error) {
	bone := NewBoneByName(LEG_IK_PARENT.StringFromDirection(direction))
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
	bone.ParentIndex = bones.findParentIndexByConfig(LEG_IK_PARENT, direction)

	return bone, nil
}

// GetLegIk 足IK取得
func (bones *Bones) GetLegIk(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(LEG_IK.StringFromDirection(direction))
}

// GetToeIK つま先IK取得
func (bones *Bones) GetToeIK(direction BoneDirection) (*Bone, error) {
	return bones.GetByName(TOE_IK.StringFromDirection(direction))
}

func (bones *Bones) findParentIndexByConfig(boneName StandardBoneName, direction BoneDirection) int {
	boneConfig := GetStandardBoneConfigs()[boneName]
	// 子ボーンが定義されているなら、その子ボーンの親を探す
	for _, tailBoneName := range boneConfig.ChildBoneNames {
		if bones.ContainsByName(tailBoneName.StringFromDirection(direction)) {
			if bone, err := bones.GetByName(tailBoneName.StringFromDirection(direction)); bone != nil && err == nil {
				return bone.ParentIndex
			}
		} else if bones.ContainsByName(tailBoneName.String()) {
			if bone, err := bones.GetByName(tailBoneName.String()); bone != nil && err == nil {
				return bone.ParentIndex
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
func (bones *Bones) InsertShortageOverrideBones() error {

	// 体幹系
	for _, funcs := range [][]func() (*Bone, error){
		{bones.GetTrunkRoot, bones.CreateTrunkRoot},
		{bones.GetLegCenter, bones.CreateLegCenter},
		{bones.GetNeckRoot, bones.CreateNeckRoot},
	} {
		getFunc := funcs[0]
		createFunc := funcs[1]

		if bone, err := getFunc(); err != nil && err == merr.NameNotFoundError && bone == nil {
			if bone, err := createFunc(); err == nil && bone != nil {
				if err := bones.Insert(bone); err != nil {
					return err
				} else {
					// 追加したボーンの親ボーンを、同じく親ボーンに設定しているボーンの親ボーンを追加ボーンに置き換える
					bones.ForEach(func(i int, b *Bone) bool {
						if b.ParentIndex == bone.ParentIndex && b.Index() != bone.Index() &&
							b.EffectIndex != bone.Index() && bone.EffectIndex != b.Index() &&
							((strings.Contains(bone.Name(), "上") && !strings.Contains(b.Name(), "下") &&
								!strings.Contains(b.Name(), "左") && !strings.Contains(b.Name(), "右")) ||
								(strings.Contains(bone.Name(), "下") && !strings.Contains(b.Name(), "上") &&
									!strings.Contains(b.Name(), "左") && !strings.Contains(b.Name(), "右"))) {
							b.ParentIndex = bone.Index()
						}
						return true
					})
					// セットアップしなおし
					bones.Setup()
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
