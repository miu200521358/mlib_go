package pmx

import (
	"strings"
	"sync"
)

// MLIB_PREFIX システム用接頭辞
const MLIB_PREFIX string = "[mlib]"

type BoneFlag uint16

const (
	// 初期値
	BONE_FLAG_NONE BoneFlag = 0x0000
	// 接続先(PMD子ボーン指定)表示方法 -> 0:座標オフセットで指定 1:ボーンで指定
	BONE_FLAG_TAIL_IS_BONE BoneFlag = 0x0001
	// 回転可能
	BONE_FLAG_CAN_ROTATE BoneFlag = 0x0002
	// 移動可能
	BONE_FLAG_CAN_TRANSLATE BoneFlag = 0x0004
	// 表示
	BONE_FLAG_IS_VISIBLE BoneFlag = 0x0008
	// 操作可
	BONE_FLAG_CAN_MANIPULATE BoneFlag = 0x0010
	// IK
	BONE_FLAG_IS_IK BoneFlag = 0x0020
	// ローカル付与 | 付与対象 0:ユーザー変形値／IKリンク／多重付与 1:親のローカル変形量
	BONE_FLAG_IS_EXTERNAL_LOCAL BoneFlag = 0x0080
	// 回転付与
	BONE_FLAG_IS_EXTERNAL_ROTATION BoneFlag = 0x0100
	// 移動付与
	BONE_FLAG_IS_EXTERNAL_TRANSLATION BoneFlag = 0x0200
	// 軸固定
	BONE_FLAG_HAS_FIXED_AXIS BoneFlag = 0x0400
	// ローカル軸
	BONE_FLAG_HAS_LOCAL_AXIS BoneFlag = 0x0800
	// 物理後変形
	BONE_FLAG_IS_AFTER_PHYSICS_DEFORM BoneFlag = 0x1000
	// 外部親変形
	BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM BoneFlag = 0x2000
)

type BoneCategory int

const (
	// ルート
	CATEGORY_ROOT BoneCategory = iota
	// 体幹
	CATEGORY_TRUNK BoneCategory = iota
	// 上半身
	CATEGORY_UPPER BoneCategory = iota
	// 下半身
	CATEGORY_LOWER BoneCategory = iota
	// 肩
	CATEGORY_SHOULDER BoneCategory = iota
	// 腕
	CATEGORY_ARM BoneCategory = iota
	// ひじ
	CATEGORY_ELBOW BoneCategory = iota
	// 足
	CATEGORY_LEG BoneCategory = iota
	// 指
	CATEGORY_FINGER BoneCategory = iota
	// 先
	CATEGORY_TAIL BoneCategory = iota
	// 足D
	CATEGORY_LEG_D BoneCategory = iota
	// 肩P
	CATEGORY_SHOULDER_P BoneCategory = iota
	// 足FK
	CATEGORY_LEG_FK BoneCategory = iota
	// 足IK
	CATEGORY_LEG_IK BoneCategory = iota
	// 足首
	CATEGORY_ANKLE BoneCategory = iota
	// 靴底
	CATEGORY_SOLE BoneCategory = iota
	// 捩
	CATEGORY_TWIST BoneCategory = iota
	// 頭
	CATEGORY_HEAD BoneCategory = iota
	// フィッティングの時に移動だけさせるか
	CATEGORY_FITTING_ONLY_MOVE BoneCategory = iota
)

type BoneConfig struct {
	// 親ボーン名候補リスト
	ParentBoneNames []StandardBoneNames
	// 末端ボーン名候補リスト
	ChildBoneNames []StandardBoneNames
	// ボーンカテゴリ
	Categories []BoneCategory
	// バウンディングボックスの形
	BoundingBoxShape Shape
	// 準標準までのボーンか
	IsStandard bool
	// 重心
	CenterOfGravity float64
	// 重心先のボーン名
	CenterOfGravityBoneNames []StandardBoneNames
}

type StandardBoneNames string

const (
	ROOT          StandardBoneNames = "全ての親"
	CENTER        StandardBoneNames = "センター"
	GROOVE        StandardBoneNames = "グルーブ"
	WAIST         StandardBoneNames = "腰"
	TRUNK_ROOT    StandardBoneNames = "体幹中心"
	LOWER_ROOT    StandardBoneNames = "下半身根元"
	LOWER         StandardBoneNames = "下半身"
	LEG_CENTER    StandardBoneNames = "足中心"
	UPPER_ROOT    StandardBoneNames = "上半身根元"
	UPPER         StandardBoneNames = "上半身"
	UPPER2        StandardBoneNames = "上半身2"
	NECK_ROOT     StandardBoneNames = "首根元"
	NECK          StandardBoneNames = "首"
	HEAD          StandardBoneNames = "頭"
	HEAD_TAIL     StandardBoneNames = "頭先先"
	EYES          StandardBoneNames = "両目"
	EYE           StandardBoneNames = "{d}目"
	SHOULDER_ROOT StandardBoneNames = "{d}肩根元"
	SHOULDER_P    StandardBoneNames = "{d}肩P"
	SHOULDER      StandardBoneNames = "{d}肩"
	SHOULDER_C    StandardBoneNames = "{d}肩C"
	ARM           StandardBoneNames = "{d}腕"
	ARM_TWIST     StandardBoneNames = "{d}腕捩"
	ARM_TWIST1    StandardBoneNames = "{d}腕捩1"
	ARM_TWIST2    StandardBoneNames = "{d}腕捩2"
	ARM_TWIST3    StandardBoneNames = "{d}腕捩3"
	ELBOW         StandardBoneNames = "{d}ひじ"
	WRIST_TWIST   StandardBoneNames = "{d}手捩"
	WRIST_TWIST1  StandardBoneNames = "{d}手捩1"
	WRIST_TWIST2  StandardBoneNames = "{d}手捩2"
	WRIST_TWIST3  StandardBoneNames = "{d}手捩3"
	WRIST         StandardBoneNames = "{d}手首"
	WRIST_TAIL    StandardBoneNames = "{d}手首先先"
	THUMB0        StandardBoneNames = "{d}親指０"
	THUMB1        StandardBoneNames = "{d}親指１"
	THUMB2        StandardBoneNames = "{d}親指２"
	THUMB_TAIL    StandardBoneNames = "{d}親指先"
	INDEX1        StandardBoneNames = "{d}人指１"
	INDEX2        StandardBoneNames = "{d}人指２"
	INDEX3        StandardBoneNames = "{d}人指３"
	INDEX_TAIL    StandardBoneNames = "{d}人指先"
	MIDDLE1       StandardBoneNames = "{d}中指１"
	MIDDLE2       StandardBoneNames = "{d}中指２"
	MIDDLE3       StandardBoneNames = "{d}中指３"
	MIDDLE_TAIL   StandardBoneNames = "{d}中指先"
	RING1         StandardBoneNames = "{d}薬指１"
	RING2         StandardBoneNames = "{d}薬指２"
	RING3         StandardBoneNames = "{d}薬指３"
	RING_TAIL     StandardBoneNames = "{d}薬指先"
	PINKY1        StandardBoneNames = "{d}小指１"
	PINKY2        StandardBoneNames = "{d}小指２"
	PINKY3        StandardBoneNames = "{d}小指３"
	PINKY_TAIL    StandardBoneNames = "{d}小指先"
	WAIST_CANCEL  StandardBoneNames = "腰キャンセル{d}"
	LEG_ROOT      StandardBoneNames = "{d}足根元"
	LEG           StandardBoneNames = "{d}足"
	KNEE          StandardBoneNames = "{d}ひざ"
	HEEL          StandardBoneNames = "{d}かかと"
	ANKLE         StandardBoneNames = "{d}足首"
	TOE_T         StandardBoneNames = "{d}つま先先"
	TOE_P         StandardBoneNames = "{d}つま先親"
	TOE_C         StandardBoneNames = "{d}つま先子"
	LEG_D         StandardBoneNames = "{d}足D"
	KNEE_D        StandardBoneNames = "{d}ひざD"
	HEEL_D        StandardBoneNames = "{d}かかとD"
	ANKLE_D       StandardBoneNames = "{d}足首D"
	TOE_T_D       StandardBoneNames = "{d}つま先先D"
	TOE_P_D       StandardBoneNames = "{d}つま先親D"
	TOE_C_D       StandardBoneNames = "{d}つま先子D"
	TOE_EX        StandardBoneNames = "{d}足先EX"
	LEG_IK_PARENT StandardBoneNames = "{d}足IK親"
	LEG_IK        StandardBoneNames = "{d}足ＩＫ"
	TOE_IK        StandardBoneNames = "{d}つま先ＩＫ"
)

func (s StandardBoneNames) String() string {
	return string(s)
}

func (s StandardBoneNames) StringFromDirection(direction string) string {
	return strings.ReplaceAll(string(s), "{d}", direction)
}

func (s StandardBoneNames) Right() string {
	return strings.ReplaceAll(string(s), "{d}", "右")
}

func (s StandardBoneNames) Left() string {
	return strings.ReplaceAll(string(s), "{d}", "左")
}

var configOnce sync.Once
var standardBoneConfigs map[StandardBoneNames]*BoneConfig

func GetStandardBoneConfigs() map[StandardBoneNames]*BoneConfig {
	configOnce.Do(func() {
		standardBoneConfigs = map[StandardBoneNames]*BoneConfig{
			ROOT: {
				ParentBoneNames:  []StandardBoneNames{},
				ChildBoneNames:   []StandardBoneNames{CENTER},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			CENTER: {
				ParentBoneNames:  []StandardBoneNames{ROOT},
				ChildBoneNames:   []StandardBoneNames{GROOVE, WAIST, TRUNK_ROOT, UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			GROOVE: {
				ParentBoneNames:  []StandardBoneNames{CENTER},
				ChildBoneNames:   []StandardBoneNames{WAIST, TRUNK_ROOT, UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			WAIST: {
				ParentBoneNames:  []StandardBoneNames{GROOVE, CENTER},
				ChildBoneNames:   []StandardBoneNames{TRUNK_ROOT, UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			TRUNK_ROOT: {
				ParentBoneNames:  []StandardBoneNames{WAIST, GROOVE, CENTER},
				ChildBoneNames:   []StandardBoneNames{UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				Categories:       []BoneCategory{CATEGORY_ROOT},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			LOWER_ROOT: {
				ParentBoneNames:          []StandardBoneNames{TRUNK_ROOT, WAIST, GROOVE, CENTER},
				ChildBoneNames:           []StandardBoneNames{LOWER},
				Categories:               []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER},
				BoundingBoxShape:         SHAPE_NONE,
				CenterOfGravity:          0.12,
				CenterOfGravityBoneNames: []StandardBoneNames{LEG_CENTER},
				IsStandard:               false},
			LOWER: {
				ParentBoneNames:  []StandardBoneNames{LOWER_ROOT, TRUNK_ROOT, WAIST, GROOVE, CENTER},
				ChildBoneNames:   []StandardBoneNames{LEG_CENTER, LEG_ROOT, WAIST_CANCEL, LEG},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			UPPER_ROOT: {
				ParentBoneNames:          []StandardBoneNames{TRUNK_ROOT, WAIST, GROOVE, CENTER},
				ChildBoneNames:           []StandardBoneNames{UPPER},
				Categories:               []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER},
				BoundingBoxShape:         SHAPE_NONE,
				CenterOfGravity:          0.24,
				CenterOfGravityBoneNames: []StandardBoneNames{NECK_ROOT},
				IsStandard:               false},
			UPPER: {
				ParentBoneNames:  []StandardBoneNames{UPPER_ROOT, TRUNK_ROOT, WAIST, GROOVE, CENTER},
				ChildBoneNames:   []StandardBoneNames{UPPER2},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			UPPER2: {
				ParentBoneNames:  []StandardBoneNames{UPPER},
				ChildBoneNames:   []StandardBoneNames{NECK_ROOT},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			NECK_ROOT: {
				ParentBoneNames:          []StandardBoneNames{UPPER2, UPPER},
				ChildBoneNames:           []StandardBoneNames{NECK},
				Categories:               []BoneCategory{CATEGORY_UPPER},
				BoundingBoxShape:         SHAPE_NONE,
				CenterOfGravity:          0.12,
				CenterOfGravityBoneNames: []StandardBoneNames{HEAD},
				IsStandard:               false},
			NECK: {
				ParentBoneNames:  []StandardBoneNames{NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:   []StandardBoneNames{HEAD},
				Categories:       []BoneCategory{CATEGORY_UPPER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			HEAD: {
				ParentBoneNames:  []StandardBoneNames{NECK},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD},
				BoundingBoxShape: SHAPE_SPHERE,
				IsStandard:       true},
			HEAD_TAIL: {
				ParentBoneNames:  []StandardBoneNames{HEAD},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			EYES: {
				ParentBoneNames:  []StandardBoneNames{HEAD},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			EYE: {
				ParentBoneNames:  []StandardBoneNames{HEAD},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD},
				BoundingBoxShape: SHAPE_SPHERE,
				IsStandard:       true},
			SHOULDER_ROOT: {
				ParentBoneNames:  []StandardBoneNames{NECK_ROOT},
				ChildBoneNames:   []StandardBoneNames{SHOULDER_P, SHOULDER},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			SHOULDER_P: {
				ParentBoneNames:  []StandardBoneNames{SHOULDER_ROOT, UPPER2, UPPER},
				ChildBoneNames:   []StandardBoneNames{SHOULDER, SHOULDER_C, ARM},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			SHOULDER: {
				ParentBoneNames:  []StandardBoneNames{SHOULDER_P, SHOULDER_ROOT, UPPER2, UPPER},
				ChildBoneNames:   []StandardBoneNames{SHOULDER_C, ARM},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			SHOULDER_C: {
				ParentBoneNames:  []StandardBoneNames{SHOULDER, SHOULDER_P, SHOULDER_ROOT, UPPER2, UPPER},
				ChildBoneNames:   []StandardBoneNames{ARM, ARM_TWIST, ELBOW},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			ARM: {
				ParentBoneNames:          []StandardBoneNames{SHOULDER_C, SHOULDER},
				ChildBoneNames:           []StandardBoneNames{ARM_TWIST, ELBOW},
				Categories:               []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM},
				BoundingBoxShape:         SHAPE_CAPSULE,
				CenterOfGravity:          0.03,
				CenterOfGravityBoneNames: []StandardBoneNames{ELBOW},
				IsStandard:               true},
			ARM_TWIST: {
				ParentBoneNames:  []StandardBoneNames{ARM},
				ChildBoneNames:   []StandardBoneNames{ELBOW},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			ARM_TWIST1: {
				ParentBoneNames:  []StandardBoneNames{ARM},
				ChildBoneNames:   []StandardBoneNames{ELBOW},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			ARM_TWIST2: {
				ParentBoneNames:  []StandardBoneNames{ARM},
				ChildBoneNames:   []StandardBoneNames{ELBOW},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			ARM_TWIST3: {
				ParentBoneNames:  []StandardBoneNames{ARM},
				ChildBoneNames:   []StandardBoneNames{ELBOW},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			ELBOW: {
				ParentBoneNames:          []StandardBoneNames{ARM_TWIST, ARM},
				ChildBoneNames:           []StandardBoneNames{WRIST_TWIST},
				Categories:               []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_ARM},
				BoundingBoxShape:         SHAPE_CAPSULE,
				CenterOfGravity:          0.03,
				CenterOfGravityBoneNames: []StandardBoneNames{WRIST},
				IsStandard:               true},
			WRIST_TWIST: {
				ParentBoneNames:  []StandardBoneNames{ELBOW},
				ChildBoneNames:   []StandardBoneNames{WRIST},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST, CATEGORY_ARM},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			WRIST_TWIST1: {
				ParentBoneNames:  []StandardBoneNames{ELBOW},
				ChildBoneNames:   []StandardBoneNames{WRIST},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			WRIST_TWIST2: {
				ParentBoneNames:  []StandardBoneNames{ELBOW},
				ChildBoneNames:   []StandardBoneNames{WRIST},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			WRIST_TWIST3: {
				ParentBoneNames:  []StandardBoneNames{ELBOW},
				ChildBoneNames:   []StandardBoneNames{WRIST},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			WRIST: {
				ParentBoneNames:  []StandardBoneNames{WRIST_TWIST, ELBOW},
				ChildBoneNames:   []StandardBoneNames{WRIST_TAIL},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_ARM},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			WRIST_TAIL: {
				ParentBoneNames:  []StandardBoneNames{WRIST},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			THUMB0: {
				ParentBoneNames:  []StandardBoneNames{WRIST},
				ChildBoneNames:   []StandardBoneNames{THUMB1},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			THUMB1: {
				ParentBoneNames:  []StandardBoneNames{THUMB0},
				ChildBoneNames:   []StandardBoneNames{THUMB2},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			THUMB2: {
				ParentBoneNames:  []StandardBoneNames{THUMB1},
				ChildBoneNames:   []StandardBoneNames{THUMB_TAIL},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			THUMB_TAIL: {
				ParentBoneNames:  []StandardBoneNames{THUMB2},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			INDEX1: {
				ParentBoneNames:  []StandardBoneNames{WRIST},
				ChildBoneNames:   []StandardBoneNames{INDEX2},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			INDEX2: {
				ParentBoneNames:  []StandardBoneNames{INDEX1},
				ChildBoneNames:   []StandardBoneNames{INDEX3},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			INDEX3: {
				ParentBoneNames:  []StandardBoneNames{INDEX2},
				ChildBoneNames:   []StandardBoneNames{INDEX_TAIL},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			INDEX_TAIL: {
				ParentBoneNames:  []StandardBoneNames{INDEX3},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			MIDDLE1: {
				ParentBoneNames:  []StandardBoneNames{WRIST},
				ChildBoneNames:   []StandardBoneNames{MIDDLE2},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			MIDDLE2: {
				ParentBoneNames:  []StandardBoneNames{MIDDLE1},
				ChildBoneNames:   []StandardBoneNames{MIDDLE3},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			MIDDLE3: {
				ParentBoneNames:  []StandardBoneNames{MIDDLE2},
				ChildBoneNames:   []StandardBoneNames{MIDDLE_TAIL},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			MIDDLE_TAIL: {
				ParentBoneNames:  []StandardBoneNames{MIDDLE3},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			RING1: {
				ParentBoneNames:  []StandardBoneNames{WRIST},
				ChildBoneNames:   []StandardBoneNames{RING2},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			RING2: {
				ParentBoneNames:  []StandardBoneNames{RING1},
				ChildBoneNames:   []StandardBoneNames{RING3},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			RING3: {
				ParentBoneNames:  []StandardBoneNames{RING2},
				ChildBoneNames:   []StandardBoneNames{RING_TAIL},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			RING_TAIL: {
				ParentBoneNames:  []StandardBoneNames{RING3},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			PINKY1: {
				ParentBoneNames:  []StandardBoneNames{WRIST},
				ChildBoneNames:   []StandardBoneNames{PINKY2},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			PINKY2: {
				ParentBoneNames:  []StandardBoneNames{PINKY1},
				ChildBoneNames:   []StandardBoneNames{PINKY3},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			PINKY3: {
				ParentBoneNames:  []StandardBoneNames{PINKY2},
				ChildBoneNames:   []StandardBoneNames{PINKY_TAIL},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			PINKY_TAIL: {
				ParentBoneNames:  []StandardBoneNames{PINKY3},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			LEG_CENTER: {
				ParentBoneNames:  []StandardBoneNames{LOWER, LOWER_ROOT},
				ChildBoneNames:   []StandardBoneNames{LEG_ROOT, WAIST_CANCEL, LEG},
				Categories:       []BoneCategory{CATEGORY_LOWER},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			LEG_ROOT: {
				ParentBoneNames:  []StandardBoneNames{LEG_CENTER, LOWER, LOWER_ROOT},
				ChildBoneNames:   []StandardBoneNames{WAIST_CANCEL, LEG},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			WAIST_CANCEL: {
				ParentBoneNames:  []StandardBoneNames{LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:   []StandardBoneNames{LEG, KNEE},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			LEG: {
				ParentBoneNames:  []StandardBoneNames{WAIST_CANCEL, LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:   []StandardBoneNames{KNEE},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			KNEE: {
				ParentBoneNames:  []StandardBoneNames{LEG},
				ChildBoneNames:   []StandardBoneNames{ANKLE},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			ANKLE: {
				ParentBoneNames:  []StandardBoneNames{KNEE},
				ChildBoneNames:   []StandardBoneNames{TOE_T},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			HEEL: {
				ParentBoneNames:  []StandardBoneNames{ANKLE},
				ChildBoneNames:   []StandardBoneNames{TOE_T},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			TOE_T: {
				ParentBoneNames:  []StandardBoneNames{ANKLE},
				ChildBoneNames:   []StandardBoneNames{TOE_P},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE, CATEGORY_TAIL},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			TOE_P: {
				ParentBoneNames:  []StandardBoneNames{TOE_T},
				ChildBoneNames:   []StandardBoneNames{TOE_C},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			TOE_C: {
				ParentBoneNames:  []StandardBoneNames{TOE_P},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			LEG_D: {
				ParentBoneNames:          []StandardBoneNames{WAIST_CANCEL, LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:           []StandardBoneNames{KNEE_D},
				Categories:               []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D},
				BoundingBoxShape:         SHAPE_CAPSULE,
				CenterOfGravity:          0.12,
				CenterOfGravityBoneNames: []StandardBoneNames{KNEE_D},
				IsStandard:               true},
			KNEE_D: {
				ParentBoneNames:          []StandardBoneNames{LEG},
				ChildBoneNames:           []StandardBoneNames{ANKLE_D},
				Categories:               []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D},
				BoundingBoxShape:         SHAPE_CAPSULE,
				CenterOfGravity:          0.08,
				CenterOfGravityBoneNames: []StandardBoneNames{HEEL_D},
				IsStandard:               true},
			ANKLE_D: {
				ParentBoneNames:  []StandardBoneNames{KNEE_D},
				ChildBoneNames:   []StandardBoneNames{TOE_EX},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			HEEL_D: {
				ParentBoneNames:  []StandardBoneNames{ANKLE_D},
				ChildBoneNames:   []StandardBoneNames{TOE_T_D},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			TOE_EX: {
				ParentBoneNames:  []StandardBoneNames{ANKLE_D},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true},
			TOE_T_D: {
				ParentBoneNames:  []StandardBoneNames{TOE_EX},
				ChildBoneNames:   []StandardBoneNames{TOE_P_D},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			TOE_P_D: {
				ParentBoneNames:  []StandardBoneNames{TOE_T_D},
				ChildBoneNames:   []StandardBoneNames{TOE_C_D},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			TOE_C_D: {
				ParentBoneNames:  []StandardBoneNames{TOE_P_D},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false},
			LEG_IK_PARENT: {
				ParentBoneNames:  []StandardBoneNames{ROOT},
				ChildBoneNames:   []StandardBoneNames{LEG_IK},
				Categories:       []BoneCategory{CATEGORY_LEG_IK, CATEGORY_SOLE, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			LEG_IK: {
				ParentBoneNames:  []StandardBoneNames{LEG_IK_PARENT, ROOT},
				ChildBoneNames:   []StandardBoneNames{TOE_IK},
				Categories:       []BoneCategory{CATEGORY_LEG_IK, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
			TOE_IK: {
				ParentBoneNames:  []StandardBoneNames{LEG_IK},
				ChildBoneNames:   []StandardBoneNames{},
				Categories:       []BoneCategory{CATEGORY_LEG_IK, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true},
		}
	})
	return standardBoneConfigs
}
