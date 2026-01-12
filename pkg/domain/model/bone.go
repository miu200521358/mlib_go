// 指示: miu200521358
package model

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// BoneFlag はボーンフラグを表す。
type BoneFlag int

const (
	// BONE_FLAG_NONE はフラグなし。
	BONE_FLAG_NONE BoneFlag = 0x0000
	// BONE_FLAG_TAIL_IS_BONE は接続先ボーン指定。
	BONE_FLAG_TAIL_IS_BONE BoneFlag = 0x0001
	// BONE_FLAG_CAN_ROTATE は回転可能。
	BONE_FLAG_CAN_ROTATE BoneFlag = 0x0002
	// BONE_FLAG_CAN_TRANSLATE は移動可能。
	BONE_FLAG_CAN_TRANSLATE BoneFlag = 0x0004
	// BONE_FLAG_IS_VISIBLE は表示対象。
	BONE_FLAG_IS_VISIBLE BoneFlag = 0x0008
	// BONE_FLAG_CAN_MANIPULATE は操作可能。
	BONE_FLAG_CAN_MANIPULATE BoneFlag = 0x0010
	// BONE_FLAG_IS_IK はIK。
	BONE_FLAG_IS_IK BoneFlag = 0x0020
	// BONE_FLAG_IS_EXTERNAL_LOCAL はローカル付与。
	BONE_FLAG_IS_EXTERNAL_LOCAL BoneFlag = 0x0080
	// BONE_FLAG_IS_EXTERNAL_ROTATION は回転付与。
	BONE_FLAG_IS_EXTERNAL_ROTATION BoneFlag = 0x0100
	// BONE_FLAG_IS_EXTERNAL_TRANSLATION は移動付与。
	BONE_FLAG_IS_EXTERNAL_TRANSLATION BoneFlag = 0x0200
	// BONE_FLAG_HAS_FIXED_AXIS は軸固定。
	BONE_FLAG_HAS_FIXED_AXIS BoneFlag = 0x0400
	// BONE_FLAG_HAS_LOCAL_AXIS はローカル軸。
	BONE_FLAG_HAS_LOCAL_AXIS BoneFlag = 0x0800
	// BONE_FLAG_IS_AFTER_PHYSICS_DEFORM は物理後変形。
	BONE_FLAG_IS_AFTER_PHYSICS_DEFORM BoneFlag = 0x1000
	// BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM は外部親変形。
	BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM BoneFlag = 0x2000
)

// IkLink はIKリンクを表す。
type IkLink struct {
	BoneIndex          int
	AngleLimit         bool
	MinAngleLimit      mmath.Vec3
	MaxAngleLimit      mmath.Vec3
	LocalAngleLimit    bool
	LocalMinAngleLimit mmath.Vec3
	LocalMaxAngleLimit mmath.Vec3
}

// Ik はIK設定を表す。
type Ik struct {
	BoneIndex    int
	LoopCount    int
	UnitRotation mmath.Vec3
	Links        []IkLink
}

// Bone はボーン要素を表す。
type Bone struct {
	index            int
	name             string
	EnglishName      string
	Position         mmath.Vec3
	ParentIndex      int
	Layer            int
	BoneFlag         BoneFlag
	TailPosition     mmath.Vec3
	TailIndex        int
	EffectIndex      int
	EffectFactor     float64
	FixedAxis        mmath.Vec3
	LocalAxisX       mmath.Vec3
	LocalAxisZ       mmath.Vec3
	EffectorKey      int
	Ik               *Ik
	DisplaySlotIndex int
}

// Index はボーン index を返す。
func (b *Bone) Index() int {
	return b.index
}

// SetIndex はボーン index を設定する。
func (b *Bone) SetIndex(index int) {
	b.index = index
}

// Name はボーン名を返す。
func (b *Bone) Name() string {
	return b.name
}

// SetName はボーン名を設定する。
func (b *Bone) SetName(name string) {
	b.name = name
}

// IsValid はボーンが有効か判定する。
func (b *Bone) IsValid() bool {
	return b != nil && b.index >= 0
}
