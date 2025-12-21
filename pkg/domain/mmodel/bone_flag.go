package mmodel

// BoneFlag はボーンフラグを表します（16bit）。
type BoneFlag uint16

const (
	BoneFlagNone                   BoneFlag = 0x0000 // 初期値
	BoneFlagTailIsBone             BoneFlag = 0x0001 // 接続先:1=ボーンで指定
	BoneFlagCanRotate              BoneFlag = 0x0002 // 回転可能
	BoneFlagCanTranslate           BoneFlag = 0x0004 // 移動可能
	BoneFlagIsVisible              BoneFlag = 0x0008 // 表示
	BoneFlagCanManipulate          BoneFlag = 0x0010 // 操作可
	BoneFlagIsIK                   BoneFlag = 0x0020 // IK
	BoneFlagIsExternalLocal        BoneFlag = 0x0080 // ローカル付与
	BoneFlagIsExternalRotation     BoneFlag = 0x0100 // 回転付与
	BoneFlagIsExternalTranslation  BoneFlag = 0x0200 // 移動付与
	BoneFlagHasFixedAxis           BoneFlag = 0x0400 // 軸固定
	BoneFlagHasLocalAxis           BoneFlag = 0x0800 // ローカル軸
	BoneFlagIsAfterPhysicsDeform   BoneFlag = 0x1000 // 物理後変形
	BoneFlagIsExternalParentDeform BoneFlag = 0x2000 // 外部親変形
)

// IsTailBone は接続先がボーンかどうかを返します。
func (f BoneFlag) IsTailBone() bool {
	return f&BoneFlagTailIsBone != 0
}

// CanRotate は回転可能かどうかを返します。
func (f BoneFlag) CanRotate() bool {
	return f&BoneFlagCanRotate != 0
}

// CanTranslate は移動可能かどうかを返します。
func (f BoneFlag) CanTranslate() bool {
	return f&BoneFlagCanTranslate != 0
}

// IsVisible は表示かどうかを返します。
func (f BoneFlag) IsVisible() bool {
	return f&BoneFlagIsVisible != 0
}

// CanManipulate は操作可かどうかを返します。
func (f BoneFlag) CanManipulate() bool {
	return f&BoneFlagCanManipulate != 0
}

// IsIK はIKかどうかを返します。
func (f BoneFlag) IsIK() bool {
	return f&BoneFlagIsIK != 0
}

// IsExternalLocal はローカル付与かどうかを返します。
func (f BoneFlag) IsExternalLocal() bool {
	return f&BoneFlagIsExternalLocal != 0
}

// IsExternalRotation は回転付与かどうかを返します。
func (f BoneFlag) IsExternalRotation() bool {
	return f&BoneFlagIsExternalRotation != 0
}

// IsExternalTranslation は移動付与かどうかを返します。
func (f BoneFlag) IsExternalTranslation() bool {
	return f&BoneFlagIsExternalTranslation != 0
}

// HasFixedAxis は軸固定かどうかを返します。
func (f BoneFlag) HasFixedAxis() bool {
	return f&BoneFlagHasFixedAxis != 0
}

// HasLocalAxis はローカル軸を持つかどうかを返します。
func (f BoneFlag) HasLocalAxis() bool {
	return f&BoneFlagHasLocalAxis != 0
}

// IsAfterPhysicsDeform は物理後変形かどうかを返します。
func (f BoneFlag) IsAfterPhysicsDeform() bool {
	return f&BoneFlagIsAfterPhysicsDeform != 0
}

// IsExternalParentDeform は外部親変形かどうかを返します。
func (f BoneFlag) IsExternalParentDeform() bool {
	return f&BoneFlagIsExternalParentDeform != 0
}

// SetTailIsBone は接続先ボーンフラグを設定します。
func (f BoneFlag) SetTailIsBone(on bool) BoneFlag {
	if on {
		return f | BoneFlagTailIsBone
	}
	return f &^ BoneFlagTailIsBone
}

// SetCanRotate は回転可能フラグを設定します。
func (f BoneFlag) SetCanRotate(on bool) BoneFlag {
	if on {
		return f | BoneFlagCanRotate
	}
	return f &^ BoneFlagCanRotate
}

// SetCanTranslate は移動可能フラグを設定します。
func (f BoneFlag) SetCanTranslate(on bool) BoneFlag {
	if on {
		return f | BoneFlagCanTranslate
	}
	return f &^ BoneFlagCanTranslate
}

// SetIsVisible は表示フラグを設定します。
func (f BoneFlag) SetIsVisible(on bool) BoneFlag {
	if on {
		return f | BoneFlagIsVisible
	}
	return f &^ BoneFlagIsVisible
}

// SetCanManipulate は操作可フラグを設定します。
func (f BoneFlag) SetCanManipulate(on bool) BoneFlag {
	if on {
		return f | BoneFlagCanManipulate
	}
	return f &^ BoneFlagCanManipulate
}

// SetIsIK はIKフラグを設定します。
func (f BoneFlag) SetIsIK(on bool) BoneFlag {
	if on {
		return f | BoneFlagIsIK
	}
	return f &^ BoneFlagIsIK
}

// SetExternalRotation は回転付与フラグを設定します。
func (f BoneFlag) SetExternalRotation(on bool) BoneFlag {
	if on {
		return f | BoneFlagIsExternalRotation
	}
	return f &^ BoneFlagIsExternalRotation
}

// SetExternalTranslation は移動付与フラグを設定します。
func (f BoneFlag) SetExternalTranslation(on bool) BoneFlag {
	if on {
		return f | BoneFlagIsExternalTranslation
	}
	return f &^ BoneFlagIsExternalTranslation
}

// SetHasFixedAxis は軸固定フラグを設定します。
func (f BoneFlag) SetHasFixedAxis(on bool) BoneFlag {
	if on {
		return f | BoneFlagHasFixedAxis
	}
	return f &^ BoneFlagHasFixedAxis
}

// SetHasLocalAxis はローカル軸フラグを設定します。
func (f BoneFlag) SetHasLocalAxis(on bool) BoneFlag {
	if on {
		return f | BoneFlagHasLocalAxis
	}
	return f &^ BoneFlagHasLocalAxis
}

// SetAfterPhysicsDeform は物理後変形フラグを設定します。
func (f BoneFlag) SetAfterPhysicsDeform(on bool) BoneFlag {
	if on {
		return f | BoneFlagIsAfterPhysicsDeform
	}
	return f &^ BoneFlagIsAfterPhysicsDeform
}
