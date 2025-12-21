package mmodel

// BoneFlag はボーンフラグを表します（16bit）。
type BoneFlag uint16

const (
	BONE_FLAG_NONE                      BoneFlag = 0x0000 // 初期値
	BONE_FLAG_TAIL_IS_BONE              BoneFlag = 0x0001 // 接続先:1=ボーンで指定
	BONE_FLAG_CAN_ROTATE                BoneFlag = 0x0002 // 回転可能
	BONE_FLAG_CAN_TRANSLATE             BoneFlag = 0x0004 // 移動可能
	BONE_FLAG_IS_VISIBLE                BoneFlag = 0x0008 // 表示
	BONE_FLAG_CAN_MANIPULATE            BoneFlag = 0x0010 // 操作可
	BONE_FLAG_IS_IK                     BoneFlag = 0x0020 // IK
	BONE_FLAG_IS_EXTERNAL_LOCAL         BoneFlag = 0x0080 // ローカル付与
	BONE_FLAG_IS_EXTERNAL_ROTATION      BoneFlag = 0x0100 // 回転付与
	BONE_FLAG_IS_EXTERNAL_TRANSLATION   BoneFlag = 0x0200 // 移動付与
	BONE_FLAG_HAS_FIXED_AXIS            BoneFlag = 0x0400 // 軸固定
	BONE_FLAG_HAS_LOCAL_AXIS            BoneFlag = 0x0800 // ローカル軸
	BONE_FLAG_IS_AFTER_PHYSICS_DEFORM   BoneFlag = 0x1000 // 物理後変形
	BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM BoneFlag = 0x2000 // 外部親変形
)

// IsTailBone は接続先がボーンかどうかを返します。
func (f BoneFlag) IsTailBone() bool {
	return f&BONE_FLAG_TAIL_IS_BONE != 0
}

// CanRotate は回転可能かどうかを返します。
func (f BoneFlag) CanRotate() bool {
	return f&BONE_FLAG_CAN_ROTATE != 0
}

// CanTranslate は移動可能かどうかを返します。
func (f BoneFlag) CanTranslate() bool {
	return f&BONE_FLAG_CAN_TRANSLATE != 0
}

// IsVisible は表示かどうかを返します。
func (f BoneFlag) IsVisible() bool {
	return f&BONE_FLAG_IS_VISIBLE != 0
}

// CanManipulate は操作可かどうかを返します。
func (f BoneFlag) CanManipulate() bool {
	return f&BONE_FLAG_CAN_MANIPULATE != 0
}

// IsIK はIKかどうかを返します。
func (f BoneFlag) IsIK() bool {
	return f&BONE_FLAG_IS_IK != 0
}

// IsExternalLocal はローカル付与かどうかを返します。
func (f BoneFlag) IsExternalLocal() bool {
	return f&BONE_FLAG_IS_EXTERNAL_LOCAL != 0
}

// IsExternalRotation は回転付与かどうかを返します。
func (f BoneFlag) IsExternalRotation() bool {
	return f&BONE_FLAG_IS_EXTERNAL_ROTATION != 0
}

// IsExternalTranslation は移動付与かどうかを返します。
func (f BoneFlag) IsExternalTranslation() bool {
	return f&BONE_FLAG_IS_EXTERNAL_TRANSLATION != 0
}

// HasFixedAxis は軸固定かどうかを返します。
func (f BoneFlag) HasFixedAxis() bool {
	return f&BONE_FLAG_HAS_FIXED_AXIS != 0
}

// HasLocalAxis はローカル軸を持つかどうかを返します。
func (f BoneFlag) HasLocalAxis() bool {
	return f&BONE_FLAG_HAS_LOCAL_AXIS != 0
}

// IsAfterPhysicsDeform は物理後変形かどうかを返します。
func (f BoneFlag) IsAfterPhysicsDeform() bool {
	return f&BONE_FLAG_IS_AFTER_PHYSICS_DEFORM != 0
}

// IsExternalParentDeform は外部親変形かどうかを返します。
func (f BoneFlag) IsExternalParentDeform() bool {
	return f&BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM != 0
}

// SetTailIsBone は接続先ボーンフラグを設定します。
func (f BoneFlag) SetTailIsBone(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_TAIL_IS_BONE
	}
	return f &^ BONE_FLAG_TAIL_IS_BONE
}

// SetCanRotate は回転可能フラグを設定します。
func (f BoneFlag) SetCanRotate(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_CAN_ROTATE
	}
	return f &^ BONE_FLAG_CAN_ROTATE
}

// SetCanTranslate は移動可能フラグを設定します。
func (f BoneFlag) SetCanTranslate(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_CAN_TRANSLATE
	}
	return f &^ BONE_FLAG_CAN_TRANSLATE
}

// SetIsVisible は表示フラグを設定します。
func (f BoneFlag) SetIsVisible(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_IS_VISIBLE
	}
	return f &^ BONE_FLAG_IS_VISIBLE
}

// SetCanManipulate は操作可フラグを設定します。
func (f BoneFlag) SetCanManipulate(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_CAN_MANIPULATE
	}
	return f &^ BONE_FLAG_CAN_MANIPULATE
}

// SetIsIK はIKフラグを設定します。
func (f BoneFlag) SetIsIK(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_IS_IK
	}
	return f &^ BONE_FLAG_IS_IK
}

// SetExternalRotation は回転付与フラグを設定します。
func (f BoneFlag) SetExternalRotation(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_IS_EXTERNAL_ROTATION
	}
	return f &^ BONE_FLAG_IS_EXTERNAL_ROTATION
}

// SetExternalTranslation は移動付与フラグを設定します。
func (f BoneFlag) SetExternalTranslation(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_IS_EXTERNAL_TRANSLATION
	}
	return f &^ BONE_FLAG_IS_EXTERNAL_TRANSLATION
}

// SetHasFixedAxis は軸固定フラグを設定します。
func (f BoneFlag) SetHasFixedAxis(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_HAS_FIXED_AXIS
	}
	return f &^ BONE_FLAG_HAS_FIXED_AXIS
}

// SetHasLocalAxis はローカル軸フラグを設定します。
func (f BoneFlag) SetHasLocalAxis(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_HAS_LOCAL_AXIS
	}
	return f &^ BONE_FLAG_HAS_LOCAL_AXIS
}

// SetAfterPhysicsDeform は物理後変形フラグを設定します。
func (f BoneFlag) SetAfterPhysicsDeform(on bool) BoneFlag {
	if on {
		return f | BONE_FLAG_IS_AFTER_PHYSICS_DEFORM
	}
	return f &^ BONE_FLAG_IS_AFTER_PHYSICS_DEFORM
}
