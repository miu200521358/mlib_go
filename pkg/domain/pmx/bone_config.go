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
	// 腕
	CATEGORY_ARM BoneCategory = iota
	// 足
	CATEGORY_LEG BoneCategory = iota
	// 指
	CATEGORY_FINGER BoneCategory = iota
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
	// UP From ボーン名候補リスト
	UpFromBoneNames []StandardBoneNames
	// UP To ボーン名候補リスト
	UpToBoneNames []StandardBoneNames
	// ボーンカテゴリ
	Categories []BoneCategory
}

type StandardBoneNames string

const (
	ROOT          StandardBoneNames = "全ての親"
	CENTER        StandardBoneNames = "センター"
	GROOVE        StandardBoneNames = "グルーブ"
	WAIST         StandardBoneNames = "腰"
	WAIST_CENTER  StandardBoneNames = "体幹中心"
	LOWER_ROOT    StandardBoneNames = "下半身根元"
	LOWER         StandardBoneNames = "下半身"
	LEG_CENTER    StandardBoneNames = "足中心"
	UPPER_ROOT    StandardBoneNames = "上半身根元"
	UPPER         StandardBoneNames = "上半身"
	UPPER2        StandardBoneNames = "上半身2"
	NECK_ROOT     StandardBoneNames = "首根元"
	NECK          StandardBoneNames = "首"
	HEAD          StandardBoneNames = "頭"
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
	TOE           StandardBoneNames = "{d}つま先"
	TOE_P         StandardBoneNames = "{d}つま先親"
	TOE_C         StandardBoneNames = "{d}つま先子"
	LEG_D         StandardBoneNames = "{d}足D"
	KNEE_D        StandardBoneNames = "{d}ひざD"
	HEEL_D        StandardBoneNames = "{d}かかとD"
	ANKLE_D       StandardBoneNames = "{d}足首D"
	TOE_D         StandardBoneNames = "{d}つま先D"
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
				ParentBoneNames: []StandardBoneNames{},
				ChildBoneNames:  []StandardBoneNames{CENTER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE}},
			CENTER: {
				ParentBoneNames: []StandardBoneNames{ROOT},
				ChildBoneNames:  []StandardBoneNames{GROOVE, WAIST, WAIST_CENTER, UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE}},
			GROOVE: {
				ParentBoneNames: []StandardBoneNames{CENTER},
				ChildBoneNames:  []StandardBoneNames{WAIST, WAIST_CENTER, UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE}},
			WAIST: {
				ParentBoneNames: []StandardBoneNames{GROOVE, CENTER},
				ChildBoneNames:  []StandardBoneNames{WAIST_CENTER, UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE}},
			WAIST_CENTER: {
				ParentBoneNames: []StandardBoneNames{WAIST, GROOVE, CENTER},
				ChildBoneNames:  []StandardBoneNames{UPPER_ROOT, LOWER_ROOT, UPPER, LOWER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_ROOT}},
			LOWER_ROOT: {
				ParentBoneNames: []StandardBoneNames{WAIST_CENTER, WAIST, GROOVE, CENTER},
				ChildBoneNames:  []StandardBoneNames{LOWER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER}},
			LOWER: {
				ParentBoneNames: []StandardBoneNames{LOWER_ROOT, WAIST_CENTER, WAIST, GROOVE, CENTER},
				ChildBoneNames:  []StandardBoneNames{LEG_CENTER, LEG_ROOT, WAIST_CANCEL, LEG},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER}},
			UPPER_ROOT: {
				ParentBoneNames: []StandardBoneNames{WAIST_CENTER, WAIST, GROOVE, CENTER},
				ChildBoneNames:  []StandardBoneNames{UPPER},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER}},
			UPPER: {
				ParentBoneNames: []StandardBoneNames{UPPER_ROOT, WAIST_CENTER, WAIST, GROOVE, CENTER},
				ChildBoneNames:  []StandardBoneNames{UPPER2},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER}},
			UPPER2: {
				ParentBoneNames: []StandardBoneNames{UPPER},
				ChildBoneNames:  []StandardBoneNames{NECK_ROOT},
				UpFromBoneNames: []StandardBoneNames{UPPER},
				UpToBoneNames:   []StandardBoneNames{UPPER2},
				Categories:      []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER}},
			NECK_ROOT: {
				ParentBoneNames: []StandardBoneNames{UPPER2, UPPER},
				ChildBoneNames:  []StandardBoneNames{NECK},
				UpFromBoneNames: []StandardBoneNames{UPPER2},
				UpToBoneNames:   []StandardBoneNames{NECK_ROOT},
				Categories:      []BoneCategory{CATEGORY_UPPER}},
			NECK: {
				ParentBoneNames: []StandardBoneNames{NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:  []StandardBoneNames{HEAD},
				UpFromBoneNames: []StandardBoneNames{NECK_ROOT},
				UpToBoneNames:   []StandardBoneNames{NECK},
				Categories:      []BoneCategory{CATEGORY_UPPER}},
			HEAD: {
				ParentBoneNames: []StandardBoneNames{NECK},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{NECK},
				UpToBoneNames:   []StandardBoneNames{HEAD},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}},
			EYES: {
				ParentBoneNames: []StandardBoneNames{HEAD},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}},
			EYE: {
				ParentBoneNames: []StandardBoneNames{HEAD},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}},
			SHOULDER_ROOT: {
				ParentBoneNames: []StandardBoneNames{NECK_ROOT},
				ChildBoneNames:  []StandardBoneNames{SHOULDER_P, SHOULDER},
				UpFromBoneNames: []StandardBoneNames{NECK_ROOT, UPPER2},
				UpToBoneNames:   []StandardBoneNames{SHOULDER_ROOT},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			SHOULDER_P: {
				ParentBoneNames: []StandardBoneNames{SHOULDER_ROOT, UPPER2, UPPER},
				ChildBoneNames:  []StandardBoneNames{SHOULDER, SHOULDER_C, ARM},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			SHOULDER: {
				ParentBoneNames: []StandardBoneNames{SHOULDER_P, SHOULDER_ROOT, UPPER2, UPPER},
				ChildBoneNames:  []StandardBoneNames{SHOULDER_C, ARM},
				UpFromBoneNames: []StandardBoneNames{UPPER2},
				UpToBoneNames:   []StandardBoneNames{NECK_ROOT},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			SHOULDER_C: {
				ParentBoneNames: []StandardBoneNames{SHOULDER, SHOULDER_P, SHOULDER_ROOT, UPPER2, UPPER},
				ChildBoneNames:  []StandardBoneNames{ARM, ARM_TWIST, ELBOW},
				UpFromBoneNames: []StandardBoneNames{SHOULDER},
				UpToBoneNames:   []StandardBoneNames{ARM},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			ARM: {
				ParentBoneNames: []StandardBoneNames{SHOULDER_C, SHOULDER},
				ChildBoneNames:  []StandardBoneNames{ARM_TWIST, ELBOW},
				UpFromBoneNames: []StandardBoneNames{SHOULDER},
				UpToBoneNames:   []StandardBoneNames{ARM},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			ARM_TWIST: {
				ParentBoneNames: []StandardBoneNames{ARM},
				ChildBoneNames:  []StandardBoneNames{ELBOW},
				UpFromBoneNames: []StandardBoneNames{SHOULDER},
				UpToBoneNames:   []StandardBoneNames{ARM},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			ARM_TWIST1: {
				ParentBoneNames: []StandardBoneNames{ARM},
				ChildBoneNames:  []StandardBoneNames{ELBOW},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			ARM_TWIST2: {
				ParentBoneNames: []StandardBoneNames{ARM},
				ChildBoneNames:  []StandardBoneNames{ELBOW},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			ARM_TWIST3: {
				ParentBoneNames: []StandardBoneNames{ARM},
				ChildBoneNames:  []StandardBoneNames{ELBOW},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			ELBOW: {
				ParentBoneNames: []StandardBoneNames{ARM_TWIST, ARM},
				ChildBoneNames:  []StandardBoneNames{WRIST_TWIST},
				UpFromBoneNames: []StandardBoneNames{UPPER2},
				UpToBoneNames:   []StandardBoneNames{NECK_ROOT},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			WRIST_TWIST: {
				ParentBoneNames: []StandardBoneNames{ELBOW},
				ChildBoneNames:  []StandardBoneNames{WRIST},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			WRIST_TWIST1: {
				ParentBoneNames: []StandardBoneNames{ELBOW},
				ChildBoneNames:  []StandardBoneNames{WRIST},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			WRIST_TWIST2: {
				ParentBoneNames: []StandardBoneNames{ELBOW},
				ChildBoneNames:  []StandardBoneNames{WRIST},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			WRIST_TWIST3: {
				ParentBoneNames: []StandardBoneNames{ELBOW},
				ChildBoneNames:  []StandardBoneNames{WRIST},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST}},
			WRIST: {
				ParentBoneNames: []StandardBoneNames{WRIST_TWIST, ELBOW},
				ChildBoneNames:  []StandardBoneNames{MIDDLE1},
				UpFromBoneNames: []StandardBoneNames{WRIST_TWIST, ELBOW},
				UpToBoneNames:   []StandardBoneNames{WRIST},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}},
			THUMB0: {
				ParentBoneNames: []StandardBoneNames{WRIST},
				ChildBoneNames:  []StandardBoneNames{THUMB1},
				UpFromBoneNames: []StandardBoneNames{WRIST},
				UpToBoneNames:   []StandardBoneNames{THUMB0},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			THUMB1: {
				ParentBoneNames: []StandardBoneNames{THUMB0},
				ChildBoneNames:  []StandardBoneNames{THUMB2},
				UpFromBoneNames: []StandardBoneNames{THUMB0},
				UpToBoneNames:   []StandardBoneNames{THUMB1},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			THUMB2: {
				ParentBoneNames: []StandardBoneNames{THUMB1},
				ChildBoneNames:  []StandardBoneNames{THUMB_TAIL},
				UpFromBoneNames: []StandardBoneNames{THUMB1},
				UpToBoneNames:   []StandardBoneNames{THUMB2},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			THUMB_TAIL: {
				ParentBoneNames: []StandardBoneNames{THUMB2},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			INDEX1: {
				ParentBoneNames: []StandardBoneNames{WRIST},
				ChildBoneNames:  []StandardBoneNames{INDEX2},
				UpFromBoneNames: []StandardBoneNames{WRIST},
				UpToBoneNames:   []StandardBoneNames{INDEX1},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			INDEX2: {
				ParentBoneNames: []StandardBoneNames{INDEX1},
				ChildBoneNames:  []StandardBoneNames{INDEX3},
				UpFromBoneNames: []StandardBoneNames{INDEX1},
				UpToBoneNames:   []StandardBoneNames{INDEX2},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			INDEX3: {
				ParentBoneNames: []StandardBoneNames{INDEX2},
				ChildBoneNames:  []StandardBoneNames{INDEX_TAIL},
				UpFromBoneNames: []StandardBoneNames{INDEX2},
				UpToBoneNames:   []StandardBoneNames{INDEX3},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			INDEX_TAIL: {
				ParentBoneNames: []StandardBoneNames{INDEX3},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			MIDDLE1: {
				ParentBoneNames: []StandardBoneNames{WRIST},
				ChildBoneNames:  []StandardBoneNames{MIDDLE2},
				UpFromBoneNames: []StandardBoneNames{WRIST},
				UpToBoneNames:   []StandardBoneNames{MIDDLE1},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			MIDDLE2: {
				ParentBoneNames: []StandardBoneNames{MIDDLE1},
				ChildBoneNames:  []StandardBoneNames{MIDDLE3},
				UpFromBoneNames: []StandardBoneNames{MIDDLE1},
				UpToBoneNames:   []StandardBoneNames{MIDDLE2},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			MIDDLE3: {
				ParentBoneNames: []StandardBoneNames{MIDDLE2},
				ChildBoneNames:  []StandardBoneNames{MIDDLE_TAIL},
				UpFromBoneNames: []StandardBoneNames{MIDDLE2},
				UpToBoneNames:   []StandardBoneNames{MIDDLE3},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			MIDDLE_TAIL: {
				ParentBoneNames: []StandardBoneNames{MIDDLE3},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			RING1: {
				ParentBoneNames: []StandardBoneNames{WRIST},
				ChildBoneNames:  []StandardBoneNames{RING2},
				UpFromBoneNames: []StandardBoneNames{WRIST},
				UpToBoneNames:   []StandardBoneNames{RING1},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			RING2: {
				ParentBoneNames: []StandardBoneNames{RING1},
				ChildBoneNames:  []StandardBoneNames{RING3},
				UpFromBoneNames: []StandardBoneNames{RING1},
				UpToBoneNames:   []StandardBoneNames{RING2},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			RING3: {
				ParentBoneNames: []StandardBoneNames{RING2},
				ChildBoneNames:  []StandardBoneNames{RING_TAIL},
				UpFromBoneNames: []StandardBoneNames{RING2},
				UpToBoneNames:   []StandardBoneNames{RING3},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			RING_TAIL: {
				ParentBoneNames: []StandardBoneNames{RING3},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			PINKY1: {
				ParentBoneNames: []StandardBoneNames{WRIST},
				ChildBoneNames:  []StandardBoneNames{PINKY2},
				UpFromBoneNames: []StandardBoneNames{WRIST},
				UpToBoneNames:   []StandardBoneNames{PINKY1},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			PINKY2: {
				ParentBoneNames: []StandardBoneNames{PINKY1},
				ChildBoneNames:  []StandardBoneNames{PINKY3},
				UpFromBoneNames: []StandardBoneNames{PINKY1},
				UpToBoneNames:   []StandardBoneNames{PINKY2},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			PINKY3: {
				ParentBoneNames: []StandardBoneNames{PINKY2},
				ChildBoneNames:  []StandardBoneNames{PINKY_TAIL},
				UpFromBoneNames: []StandardBoneNames{PINKY2},
				UpToBoneNames:   []StandardBoneNames{PINKY3},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			PINKY_TAIL: {
				ParentBoneNames: []StandardBoneNames{PINKY3},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER}},
			LEG_CENTER: {
				ParentBoneNames: []StandardBoneNames{LOWER, LOWER_ROOT},
				ChildBoneNames:  []StandardBoneNames{LEG_ROOT, WAIST_CANCEL, LEG},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER}},
			LEG_ROOT: {
				ParentBoneNames: []StandardBoneNames{LEG_CENTER, LOWER, LOWER_ROOT},
				ChildBoneNames:  []StandardBoneNames{WAIST_CANCEL, LEG},
				UpFromBoneNames: []StandardBoneNames{LOWER},
				UpToBoneNames:   []StandardBoneNames{LEG_CENTER},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG}},
			WAIST_CANCEL: {
				ParentBoneNames: []StandardBoneNames{LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:  []StandardBoneNames{LEG, KNEE},
				UpFromBoneNames: []StandardBoneNames{LEG_ROOT},
				UpToBoneNames:   []StandardBoneNames{LEG},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG}},
			LEG: {
				ParentBoneNames: []StandardBoneNames{WAIST_CANCEL, LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:  []StandardBoneNames{KNEE},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK}},
			KNEE: {
				ParentBoneNames: []StandardBoneNames{LEG},
				ChildBoneNames:  []StandardBoneNames{ANKLE},
				UpFromBoneNames: []StandardBoneNames{LEG},
				UpToBoneNames:   []StandardBoneNames{KNEE},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK}},
			ANKLE: {
				ParentBoneNames: []StandardBoneNames{KNEE},
				ChildBoneNames:  []StandardBoneNames{TOE},
				UpFromBoneNames: []StandardBoneNames{KNEE},
				UpToBoneNames:   []StandardBoneNames{ANKLE},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE}},
			HEEL: {
				ParentBoneNames: []StandardBoneNames{ANKLE},
				ChildBoneNames:  []StandardBoneNames{TOE},
				UpFromBoneNames: []StandardBoneNames{TOE_P},
				UpToBoneNames:   []StandardBoneNames{TOE_C},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE: {
				ParentBoneNames: []StandardBoneNames{ANKLE},
				ChildBoneNames:  []StandardBoneNames{TOE_P},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE_P: {
				ParentBoneNames: []StandardBoneNames{TOE},
				ChildBoneNames:  []StandardBoneNames{TOE_C},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE_C: {
				ParentBoneNames: []StandardBoneNames{TOE_P},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE, CATEGORY_SOLE}},
			LEG_D: {
				ParentBoneNames: []StandardBoneNames{WAIST_CANCEL, LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:  []StandardBoneNames{KNEE_D},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D}},
			KNEE_D: {
				ParentBoneNames: []StandardBoneNames{LEG},
				ChildBoneNames:  []StandardBoneNames{ANKLE_D},
				UpFromBoneNames: []StandardBoneNames{LEG_D},
				UpToBoneNames:   []StandardBoneNames{KNEE_D},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D}},
			ANKLE_D: {
				ParentBoneNames: []StandardBoneNames{KNEE_D},
				ChildBoneNames:  []StandardBoneNames{TOE_EX},
				UpFromBoneNames: []StandardBoneNames{KNEE_D},
				UpToBoneNames:   []StandardBoneNames{ANKLE_D},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE}},
			HEEL_D: {
				ParentBoneNames: []StandardBoneNames{ANKLE_D},
				ChildBoneNames:  []StandardBoneNames{TOE_D},
				UpFromBoneNames: []StandardBoneNames{TOE_P_D},
				UpToBoneNames:   []StandardBoneNames{TOE_C_D},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE_EX: {
				ParentBoneNames: []StandardBoneNames{ANKLE_D},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE_D: {
				ParentBoneNames: []StandardBoneNames{TOE_EX},
				ChildBoneNames:  []StandardBoneNames{TOE_P_D},
				UpFromBoneNames: []StandardBoneNames{ANKLE_D},
				UpToBoneNames:   []StandardBoneNames{TOE_D},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE_P_D: {
				ParentBoneNames: []StandardBoneNames{TOE_D},
				ChildBoneNames:  []StandardBoneNames{TOE_C_D},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE, CATEGORY_SOLE}},
			TOE_C_D: {
				ParentBoneNames: []StandardBoneNames{TOE_P_D},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE, CATEGORY_SOLE}},
			LEG_IK_PARENT: {
				ParentBoneNames: []StandardBoneNames{ROOT},
				ChildBoneNames:  []StandardBoneNames{LEG_IK},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LEG_IK, CATEGORY_SOLE, CATEGORY_FITTING_ONLY_MOVE}},
			LEG_IK: {
				ParentBoneNames: []StandardBoneNames{LEG_IK_PARENT, ROOT},
				ChildBoneNames:  []StandardBoneNames{TOE_IK},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LEG_IK, CATEGORY_FITTING_ONLY_MOVE}},
			TOE_IK: {
				ParentBoneNames: []StandardBoneNames{LEG_IK},
				ChildBoneNames:  []StandardBoneNames{},
				UpFromBoneNames: []StandardBoneNames{},
				UpToBoneNames:   []StandardBoneNames{},
				Categories:      []BoneCategory{CATEGORY_LEG_IK, CATEGORY_FITTING_ONLY_MOVE}},
		}
	})
	return standardBoneConfigs
}
