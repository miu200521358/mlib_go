package pmx

import (
	"slices"
	"strconv"
	"strings"
	"sync"
)

// MLIB_PREFIX システム用接頭辞
const MLIB_PREFIX string = "[mlib]"
const BONE_DIRECTION_PREFIX string = "{d}"

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

type BoneDirection string

const (
	// 右
	BONE_DIRECTION_RIGHT BoneDirection = "右"
	// 左
	BONE_DIRECTION_LEFT BoneDirection = "左"
	// 体幹
	BONE_DIRECTION_TRUNK BoneDirection = ""
)

func (d BoneDirection) String() string {
	return string(d)
}

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
	ParentBoneNames []StandardBoneName
	// 末端ボーン名候補リスト
	ChildBoneNames [][]StandardBoneName
	// ボーンカテゴリ
	Categories []BoneCategory
	// 表示枠
	DisplaySlot StandardBoneName
	// バウンディングボックスの形
	BoundingBoxShape Shape
	// 準標準までのボーンか
	IsStandard bool
	// 重心先ボーン名
	GravityTargetBoneName StandardBoneName
	// 略称
	Abbreviation StandardBoneName
}

type StandardBoneName string

const (
	ROOT           StandardBoneName = "全ての親"
	CENTER         StandardBoneName = "センター"
	GROOVE         StandardBoneName = "グルーブ"
	BODY_AXIS      StandardBoneName = "体軸"
	TRUNK_ROOT     StandardBoneName = "体幹中心"
	WAIST          StandardBoneName = "腰"
	LOWER_ROOT     StandardBoneName = "下半身根元"
	LOWER          StandardBoneName = "下半身"
	HIP            StandardBoneName = "{d}腰骨"
	LEG_CENTER     StandardBoneName = "足中心"
	UPPER_ROOT     StandardBoneName = "上半身根元"
	UPPER          StandardBoneName = "上半身"
	UPPER2         StandardBoneName = "上半身2"
	NECK_ROOT      StandardBoneName = "首根元"
	NECK           StandardBoneName = "首"
	HEAD           StandardBoneName = "頭"
	HEAD_TAIL      StandardBoneName = "頭先先"
	EYES           StandardBoneName = "両目"
	EYE            StandardBoneName = "{d}目"
	SHOULDER_ROOT  StandardBoneName = "{d}肩根元"
	SHOULDER_P     StandardBoneName = "{d}肩P"
	SHOULDER       StandardBoneName = "{d}肩"
	SHOULDER_C     StandardBoneName = "{d}肩C"
	ARM            StandardBoneName = "{d}腕"
	ARM_TWIST      StandardBoneName = "{d}腕捩"
	ARM_TWIST1     StandardBoneName = "{d}腕捩1"
	ARM_TWIST2     StandardBoneName = "{d}腕捩2"
	ARM_TWIST3     StandardBoneName = "{d}腕捩3"
	ELBOW          StandardBoneName = "{d}ひじ"
	WRIST_TWIST    StandardBoneName = "{d}手捩"
	WRIST_TWIST1   StandardBoneName = "{d}手捩1"
	WRIST_TWIST2   StandardBoneName = "{d}手捩2"
	WRIST_TWIST3   StandardBoneName = "{d}手捩3"
	WRIST          StandardBoneName = "{d}手首"
	WRIST_TAIL     StandardBoneName = "{d}手首先先"
	THUMB0         StandardBoneName = "{d}親指０"
	THUMB1         StandardBoneName = "{d}親指１"
	THUMB2         StandardBoneName = "{d}親指２"
	THUMB_TAIL     StandardBoneName = "{d}親指先先"
	INDEX1         StandardBoneName = "{d}人指１"
	INDEX2         StandardBoneName = "{d}人指２"
	INDEX3         StandardBoneName = "{d}人指３"
	INDEX_TAIL     StandardBoneName = "{d}人指先先"
	MIDDLE1        StandardBoneName = "{d}中指１"
	MIDDLE2        StandardBoneName = "{d}中指２"
	MIDDLE3        StandardBoneName = "{d}中指３"
	MIDDLE_TAIL    StandardBoneName = "{d}中指先先"
	RING1          StandardBoneName = "{d}薬指１"
	RING2          StandardBoneName = "{d}薬指２"
	RING3          StandardBoneName = "{d}薬指３"
	RING_TAIL      StandardBoneName = "{d}薬指先先"
	PINKY1         StandardBoneName = "{d}小指１"
	PINKY2         StandardBoneName = "{d}小指２"
	PINKY3         StandardBoneName = "{d}小指３"
	PINKY_TAIL     StandardBoneName = "{d}小指先先"
	WAIST_CANCEL   StandardBoneName = "腰キャンセル{d}"
	LEG_ROOT       StandardBoneName = "{d}足根元"
	LEG            StandardBoneName = "{d}足"
	KNEE           StandardBoneName = "{d}ひざ"
	ANKLE          StandardBoneName = "{d}足首"
	ANKLE_GROUND   StandardBoneName = "{d}足首地面"
	HEEL           StandardBoneName = "{d}かかと"
	TOE_T          StandardBoneName = "{d}つま先先"
	TOE_P          StandardBoneName = "{d}つま先親"
	TOE_C          StandardBoneName = "{d}つま先子"
	LEG_D          StandardBoneName = "{d}足D"
	KNEE_D         StandardBoneName = "{d}ひざD"
	HEEL_D         StandardBoneName = "{d}かかとD"
	ANKLE_D        StandardBoneName = "{d}足首D"
	ANKLE_D_GROUND StandardBoneName = "{d}足首地面D"
	TOE_T_D        StandardBoneName = "{d}つま先先D"
	TOE_P_D        StandardBoneName = "{d}つま先親D"
	TOE_C_D        StandardBoneName = "{d}つま先子D"
	TOE_EX         StandardBoneName = "{d}足先EX"
	LEG_IK_PARENT  StandardBoneName = "{d}足IK親"
	LEG_IK         StandardBoneName = "{d}足ＩＫ"
	TOE_IK         StandardBoneName = "{d}つま先ＩＫ"
)

func (s StandardBoneName) String() string {
	return string(s)
}

func (s StandardBoneName) StringFromDirection(direction BoneDirection) string {
	return strings.ReplaceAll(string(s), BONE_DIRECTION_PREFIX, string(direction))
}

func (s StandardBoneName) StringFromDirectionAndIdx(direction BoneDirection, idx int) string {
	return strings.ReplaceAll(string(s), BONE_DIRECTION_PREFIX, string(direction)) + strconv.Itoa(idx)
}

func (s StandardBoneName) Right() string {
	return strings.ReplaceAll(string(s), BONE_DIRECTION_PREFIX, "右")
}

func (s StandardBoneName) Left() string {
	return strings.ReplaceAll(string(s), BONE_DIRECTION_PREFIX, "左")
}

func BoneConfigFromName(name string) *BoneConfig {
	for _, direction := range []BoneDirection{BONE_DIRECTION_TRUNK, BONE_DIRECTION_RIGHT, BONE_DIRECTION_LEFT} {
		for boneName, v := range standardBoneConfigs {
			if boneName.StringFromDirection(direction) == name {
				return v
			}
		}
	}
	return nil
}

var configOnce sync.Once
var standardBoneConfigs map[StandardBoneName]*BoneConfig

func GetStandardBoneConfigs() map[StandardBoneName]*BoneConfig {
	configOnce.Do(func() {
		standardBoneConfigs = map[StandardBoneName]*BoneConfig{
			ROOT: {
				ParentBoneNames:  []StandardBoneName{},
				ChildBoneNames:   [][]StandardBoneName{{CENTER}, {LEG_IK_PARENT}, {LEG_IK}},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("全")},
			CENTER: {
				ParentBoneNames:  []StandardBoneName{ROOT},
				ChildBoneNames:   [][]StandardBoneName{{GROOVE}, {BODY_AXIS}, {TRUNK_ROOT}, {WAIST}, {UPPER_ROOT, LOWER_ROOT}, {UPPER, LOWER}},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      CENTER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("C")},
			GROOVE: {
				ParentBoneNames:  []StandardBoneName{CENTER},
				ChildBoneNames:   [][]StandardBoneName{{BODY_AXIS}, {TRUNK_ROOT}, {WAIST}, {UPPER_ROOT, LOWER_ROOT}, {UPPER, LOWER}},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      CENTER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("G")},
			BODY_AXIS: {
				ParentBoneNames:  []StandardBoneName{CENTER, GROOVE},
				ChildBoneNames:   [][]StandardBoneName{{TRUNK_ROOT}, {WAIST}, {UPPER_ROOT, LOWER_ROOT}, {UPPER, LOWER}},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      CENTER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("体軸")},
			TRUNK_ROOT: {
				ParentBoneNames:  []StandardBoneName{BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:   [][]StandardBoneName{{WAIST}, {UPPER_ROOT, LOWER_ROOT}, {UPPER, LOWER}},
				Categories:       []BoneCategory{CATEGORY_ROOT},
				DisplaySlot:      CENTER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("体幹")},
			WAIST: {
				ParentBoneNames:  []StandardBoneName{TRUNK_ROOT, BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:   [][]StandardBoneName{{UPPER_ROOT, LOWER_ROOT}, {UPPER, LOWER}},
				Categories:       []BoneCategory{CATEGORY_ROOT, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      CENTER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("腰")},
			LOWER_ROOT: {
				ParentBoneNames:       []StandardBoneName{WAIST, TRUNK_ROOT, BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:        [][]StandardBoneName{{LOWER}},
				Categories:            []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER},
				DisplaySlot:           LOWER,
				BoundingBoxShape:      SHAPE_NONE,
				IsStandard:            false,
				GravityTargetBoneName: LEG_CENTER,
				Abbreviation:          StandardBoneName("下根")},
			LOWER: {
				ParentBoneNames:  []StandardBoneName{LOWER_ROOT, WAIST, TRUNK_ROOT, BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:   [][]StandardBoneName{{LEG_CENTER, HIP}, {LEG_ROOT}, {WAIST_CANCEL}, {LEG, LEG_D}},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER},
				DisplaySlot:      LOWER,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("下")},
			HIP: {
				ParentBoneNames:  []StandardBoneName{LOWER, LOWER_ROOT, WAIST, TRUNK_ROOT, BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:   [][]StandardBoneName{{LEG_CENTER}, {LEG_ROOT}, {WAIST_CANCEL}, {LEG, LEG_D}},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_LOWER},
				DisplaySlot:      LEG,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}尻")},
			UPPER_ROOT: {
				ParentBoneNames:       []StandardBoneName{WAIST, TRUNK_ROOT, BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:        [][]StandardBoneName{{UPPER}},
				Categories:            []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER},
				DisplaySlot:           UPPER,
				BoundingBoxShape:      SHAPE_NONE,
				IsStandard:            false,
				GravityTargetBoneName: NECK_ROOT,
				Abbreviation:          StandardBoneName("上根")},
			UPPER: {
				ParentBoneNames:  []StandardBoneName{UPPER_ROOT, WAIST, TRUNK_ROOT, BODY_AXIS, GROOVE, CENTER},
				ChildBoneNames:   [][]StandardBoneName{{UPPER2}},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER},
				DisplaySlot:      UPPER,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("上")},
			UPPER2: {
				ParentBoneNames:  []StandardBoneName{UPPER},
				ChildBoneNames:   [][]StandardBoneName{{NECK_ROOT}},
				Categories:       []BoneCategory{CATEGORY_TRUNK, CATEGORY_UPPER},
				DisplaySlot:      UPPER,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("上2")},
			NECK_ROOT: {
				ParentBoneNames:       []StandardBoneName{UPPER2, UPPER},
				ChildBoneNames:        [][]StandardBoneName{{NECK, SHOULDER_ROOT}, {SHOULDER_P}, {SHOULDER}},
				Categories:            []BoneCategory{CATEGORY_UPPER},
				DisplaySlot:           UPPER,
				IsStandard:            false,
				GravityTargetBoneName: HEAD,
				Abbreviation:          StandardBoneName("首根")},
			NECK: {
				ParentBoneNames:  []StandardBoneName{NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:   [][]StandardBoneName{{HEAD}},
				Categories:       []BoneCategory{CATEGORY_UPPER},
				DisplaySlot:      UPPER,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("首")},
			HEAD: {
				ParentBoneNames:  []StandardBoneName{NECK},
				ChildBoneNames:   [][]StandardBoneName{{HEAD_TAIL}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD},
				DisplaySlot:      HEAD,
				BoundingBoxShape: SHAPE_SPHERE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("頭")},
			HEAD_TAIL: {
				ParentBoneNames:  []StandardBoneName{HEAD},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD, CATEGORY_TAIL},
				DisplaySlot:      HEAD,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("頭先")},
			EYES: {
				ParentBoneNames:  []StandardBoneName{HEAD},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD},
				DisplaySlot:      HEAD,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("両目")},
			EYE: {
				ParentBoneNames:  []StandardBoneName{HEAD},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD},
				DisplaySlot:      HEAD,
				BoundingBoxShape: SHAPE_SPHERE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}目")},
			SHOULDER_ROOT: {
				ParentBoneNames:  []StandardBoneName{NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:   [][]StandardBoneName{{SHOULDER_P}, {SHOULDER}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				DisplaySlot:      SHOULDER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}肩根")},
			SHOULDER_P: {
				ParentBoneNames:  []StandardBoneName{SHOULDER_ROOT, NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:   [][]StandardBoneName{{SHOULDER}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				DisplaySlot:      SHOULDER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}肩P")},
			SHOULDER: {
				ParentBoneNames:  []StandardBoneName{SHOULDER_P, SHOULDER_ROOT, NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:   [][]StandardBoneName{{SHOULDER_C}, {ARM}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				DisplaySlot:      SHOULDER,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}肩")},
			SHOULDER_C: {
				ParentBoneNames:  []StandardBoneName{SHOULDER, SHOULDER_P, SHOULDER_ROOT, NECK_ROOT, UPPER2, UPPER},
				ChildBoneNames:   [][]StandardBoneName{{ARM}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_SHOULDER},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}肩C")},
			ARM: {
				ParentBoneNames:       []StandardBoneName{SHOULDER},
				ChildBoneNames:        [][]StandardBoneName{{ARM_TWIST}, {ELBOW}},
				Categories:            []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM},
				DisplaySlot:           ARM,
				BoundingBoxShape:      SHAPE_CAPSULE,
				IsStandard:            true,
				GravityTargetBoneName: ELBOW,
				Abbreviation:          StandardBoneName("{d}腕")},
			ARM_TWIST: {
				ParentBoneNames:  []StandardBoneName{ARM},
				ChildBoneNames:   [][]StandardBoneName{{ELBOW}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST},
				DisplaySlot:      ARM,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}腕捩")},
			ARM_TWIST1: {
				ParentBoneNames:  []StandardBoneName{ARM},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}腕捩1")},
			ARM_TWIST2: {
				ParentBoneNames:  []StandardBoneName{ARM},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}腕捩2")},
			ARM_TWIST3: {
				ParentBoneNames:  []StandardBoneName{ARM},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}腕捩3")},
			ELBOW: {
				ParentBoneNames:       []StandardBoneName{ARM_TWIST, ARM},
				ChildBoneNames:        [][]StandardBoneName{{WRIST_TWIST}},
				Categories:            []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_ARM},
				DisplaySlot:           ELBOW,
				BoundingBoxShape:      SHAPE_CAPSULE,
				IsStandard:            true,
				GravityTargetBoneName: WRIST,
				Abbreviation:          StandardBoneName("{d}肘")},
			WRIST_TWIST: {
				ParentBoneNames:  []StandardBoneName{ELBOW},
				ChildBoneNames:   [][]StandardBoneName{{WRIST}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST, CATEGORY_ARM},
				DisplaySlot:      ELBOW,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}手捩")},
			WRIST_TWIST1: {
				ParentBoneNames:  []StandardBoneName{ELBOW},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}手捩1")},
			WRIST_TWIST2: {
				ParentBoneNames:  []StandardBoneName{ELBOW},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}手捩2")},
			WRIST_TWIST3: {
				ParentBoneNames:  []StandardBoneName{ELBOW},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TWIST},
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}手捩3")},
			WRIST: {
				ParentBoneNames:       []StandardBoneName{WRIST_TWIST, ELBOW},
				ChildBoneNames:        [][]StandardBoneName{{WRIST_TAIL}},
				Categories:            []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_ARM},
				DisplaySlot:           WRIST,
				BoundingBoxShape:      SHAPE_CAPSULE,
				IsStandard:            true,
				GravityTargetBoneName: WRIST_TAIL,
				Abbreviation:          StandardBoneName("{d}手首")},
			WRIST_TAIL: {
				ParentBoneNames:  []StandardBoneName{WRIST},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_ELBOW, CATEGORY_TAIL},
				DisplaySlot:      WRIST,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}手先")},
			THUMB0: {
				ParentBoneNames:  []StandardBoneName{WRIST},
				ChildBoneNames:   [][]StandardBoneName{{THUMB1}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      THUMB1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}親0")},
			THUMB1: {
				ParentBoneNames:  []StandardBoneName{THUMB0},
				ChildBoneNames:   [][]StandardBoneName{{THUMB2}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      THUMB1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}親1")},
			THUMB2: {
				ParentBoneNames:  []StandardBoneName{THUMB1},
				ChildBoneNames:   [][]StandardBoneName{{THUMB_TAIL}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      THUMB1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}親2")},
			THUMB_TAIL: {
				ParentBoneNames:  []StandardBoneName{THUMB2},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				DisplaySlot:      THUMB1,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}親先")},
			INDEX1: {
				ParentBoneNames:  []StandardBoneName{WRIST},
				ChildBoneNames:   [][]StandardBoneName{{INDEX2}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      INDEX1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}人1")},
			INDEX2: {
				ParentBoneNames:  []StandardBoneName{INDEX1},
				ChildBoneNames:   [][]StandardBoneName{{INDEX3}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      INDEX1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}人2")},
			INDEX3: {
				ParentBoneNames:  []StandardBoneName{INDEX2},
				ChildBoneNames:   [][]StandardBoneName{{INDEX_TAIL}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      INDEX1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}人3")},
			INDEX_TAIL: {
				ParentBoneNames:  []StandardBoneName{INDEX3},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				DisplaySlot:      INDEX1,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}人先")},
			MIDDLE1: {
				ParentBoneNames:  []StandardBoneName{WRIST},
				ChildBoneNames:   [][]StandardBoneName{{MIDDLE2}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      MIDDLE1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}中1")},
			MIDDLE2: {
				ParentBoneNames:  []StandardBoneName{MIDDLE1},
				ChildBoneNames:   [][]StandardBoneName{{MIDDLE3}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      MIDDLE1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}中2")},
			MIDDLE3: {
				ParentBoneNames:  []StandardBoneName{MIDDLE2},
				ChildBoneNames:   [][]StandardBoneName{{MIDDLE_TAIL}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      MIDDLE1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}中3")},
			MIDDLE_TAIL: {
				ParentBoneNames:  []StandardBoneName{MIDDLE3},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				DisplaySlot:      MIDDLE1,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}中先")},
			RING1: {
				ParentBoneNames:  []StandardBoneName{WRIST},
				ChildBoneNames:   [][]StandardBoneName{{RING2}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      RING1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}薬1")},
			RING2: {
				ParentBoneNames:  []StandardBoneName{RING1},
				ChildBoneNames:   [][]StandardBoneName{{RING3}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      RING1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}薬2")},
			RING3: {
				ParentBoneNames:  []StandardBoneName{RING2},
				ChildBoneNames:   [][]StandardBoneName{{RING_TAIL}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      RING1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}薬3")},
			RING_TAIL: {
				ParentBoneNames:  []StandardBoneName{RING3},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				DisplaySlot:      RING1,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}薬先")},
			PINKY1: {
				ParentBoneNames:  []StandardBoneName{WRIST},
				ChildBoneNames:   [][]StandardBoneName{{PINKY2}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      PINKY1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}小1")},
			PINKY2: {
				ParentBoneNames:  []StandardBoneName{PINKY1},
				ChildBoneNames:   [][]StandardBoneName{{PINKY3}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      PINKY1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}小2")},
			PINKY3: {
				ParentBoneNames:  []StandardBoneName{PINKY2},
				ChildBoneNames:   [][]StandardBoneName{{PINKY_TAIL}},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER},
				DisplaySlot:      PINKY1,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}小3")},
			PINKY_TAIL: {
				ParentBoneNames:  []StandardBoneName{PINKY3},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_UPPER, CATEGORY_FINGER, CATEGORY_TAIL},
				DisplaySlot:      PINKY1,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}小先")},
			LEG_CENTER: {
				ParentBoneNames:  []StandardBoneName{LOWER, LOWER_ROOT},
				ChildBoneNames:   [][]StandardBoneName{{LEG_ROOT}, {WAIST_CANCEL}, {LEG}},
				Categories:       []BoneCategory{CATEGORY_LOWER},
				DisplaySlot:      LOWER,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("足中")},
			WAIST_CANCEL: {
				ParentBoneNames:  []StandardBoneName{LEG_CENTER, LOWER, LOWER_ROOT},
				ChildBoneNames:   [][]StandardBoneName{{LEG_ROOT}, {LEG}, {KNEE}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG},
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}腰C")},
			LEG_ROOT: {
				ParentBoneNames:  []StandardBoneName{WAIST_CANCEL, LEG_CENTER, LOWER, LOWER_ROOT},
				ChildBoneNames:   [][]StandardBoneName{{LEG, LEG_D}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG},
				DisplaySlot:      LEG,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}足根")},
			LEG: {
				ParentBoneNames:  []StandardBoneName{WAIST_CANCEL, LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:   [][]StandardBoneName{{KNEE}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK},
				DisplaySlot:      LEG,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}足")},
			KNEE: {
				ParentBoneNames:  []StandardBoneName{LEG},
				ChildBoneNames:   [][]StandardBoneName{{ANKLE}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK},
				DisplaySlot:      KNEE,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}膝")},
			ANKLE: {
				ParentBoneNames:  []StandardBoneName{KNEE},
				ChildBoneNames:   [][]StandardBoneName{{ANKLE_GROUND, TOE_T, HEEL}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}足首")},
			ANKLE_GROUND: {
				ParentBoneNames:  []StandardBoneName{ANKLE},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}首地")},
			HEEL: {
				ParentBoneNames:  []StandardBoneName{ANKLE},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}踵")},
			TOE_T: {
				ParentBoneNames:  []StandardBoneName{ANKLE},
				ChildBoneNames:   [][]StandardBoneName{{TOE_P}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE, CATEGORY_TAIL},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}爪先")},
			TOE_P: {
				ParentBoneNames:  []StandardBoneName{TOE_T},
				ChildBoneNames:   [][]StandardBoneName{{TOE_C}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}爪先P")},
			TOE_C: {
				ParentBoneNames:  []StandardBoneName{TOE_P},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}爪先C")},
			LEG_D: {
				ParentBoneNames:       []StandardBoneName{WAIST_CANCEL, LEG_ROOT, LEG_CENTER, LOWER},
				ChildBoneNames:        [][]StandardBoneName{{KNEE_D}},
				Categories:            []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D},
				DisplaySlot:           LEG,
				BoundingBoxShape:      SHAPE_CAPSULE,
				IsStandard:            true,
				GravityTargetBoneName: KNEE_D,
				Abbreviation:          StandardBoneName("{d}足D")},
			KNEE_D: {
				ParentBoneNames:       []StandardBoneName{LEG_D},
				ChildBoneNames:        [][]StandardBoneName{{ANKLE_D}},
				Categories:            []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D},
				DisplaySlot:           KNEE,
				BoundingBoxShape:      SHAPE_CAPSULE,
				IsStandard:            true,
				GravityTargetBoneName: ANKLE_D,
				Abbreviation:          StandardBoneName("{d}膝D")},
			ANKLE_D: {
				ParentBoneNames:       []StandardBoneName{KNEE_D},
				ChildBoneNames:        [][]StandardBoneName{{ANKLE_D_GROUND, TOE_EX, HEEL_D}},
				Categories:            []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE},
				DisplaySlot:           ANKLE,
				BoundingBoxShape:      SHAPE_CAPSULE,
				IsStandard:            true,
				GravityTargetBoneName: TOE_T_D,
				Abbreviation:          StandardBoneName("{d}足首D")},
			ANKLE_D_GROUND: {
				ParentBoneNames:  []StandardBoneName{ANKLE_D},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}首地D")},
			HEEL_D: {
				ParentBoneNames:  []StandardBoneName{ANKLE_D},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}踵D")},
			TOE_EX: {
				ParentBoneNames:  []StandardBoneName{ANKLE_D},
				ChildBoneNames:   [][]StandardBoneName{{TOE_T_D}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_CAPSULE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}EX")},
			TOE_T_D: {
				ParentBoneNames:  []StandardBoneName{TOE_EX},
				ChildBoneNames:   [][]StandardBoneName{{TOE_P_D}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}爪先D")},
			TOE_P_D: {
				ParentBoneNames:  []StandardBoneName{TOE_T_D},
				ChildBoneNames:   [][]StandardBoneName{{TOE_C_D}},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}爪先PD")},
			TOE_C_D: {
				ParentBoneNames:  []StandardBoneName{TOE_P_D},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_SOLE},
				DisplaySlot:      ANKLE,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       false,
				Abbreviation:     StandardBoneName("{d}爪先CD")},
			LEG_IK_PARENT: {
				ParentBoneNames:  []StandardBoneName{ROOT},
				ChildBoneNames:   [][]StandardBoneName{{LEG_IK}},
				Categories:       []BoneCategory{CATEGORY_LEG_IK, CATEGORY_SOLE, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      LEG_IK,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}足K親")},
			LEG_IK: {
				ParentBoneNames:  []StandardBoneName{LEG_IK_PARENT, ROOT},
				ChildBoneNames:   [][]StandardBoneName{{TOE_IK}},
				Categories:       []BoneCategory{CATEGORY_LEG_IK, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      LEG_IK,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}足K")},
			TOE_IK: {
				ParentBoneNames:  []StandardBoneName{LEG_IK},
				ChildBoneNames:   [][]StandardBoneName{},
				Categories:       []BoneCategory{CATEGORY_LEG_IK, CATEGORY_FITTING_ONLY_MOVE},
				DisplaySlot:      LEG_IK,
				BoundingBoxShape: SHAPE_NONE,
				IsStandard:       true,
				Abbreviation:     StandardBoneName("{d}爪K")},
		}
	})
	return standardBoneConfigs
}

func (c *BoneConfig) IsFinger() bool {
	return slices.Contains(c.Categories, CATEGORY_FINGER)
}

func (c *BoneConfig) IsArm() bool {
	return slices.Contains(c.Categories, CATEGORY_ARM)
}
