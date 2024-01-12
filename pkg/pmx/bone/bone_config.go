package bone

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type BoneFlag uint16

const (
	// 初期値
	NONE BoneFlag = 0x0000
	// 接続先(PMD子ボーン指定)表示方法 -> 0:座標オフセットで指定 1:ボーンで指定
	TAIL_IS_BONE BoneFlag = 0x0001
	// 回転可能
	CAN_ROTATE BoneFlag = 0x0002
	// 移動可能
	CAN_TRANSLATE BoneFlag = 0x0004
	// 表示
	IS_VISIBLE BoneFlag = 0x0008
	// 操作可
	CAN_MANIPULATE BoneFlag = 0x0010
	// IK
	IS_IK BoneFlag = 0x0020
	// ローカル付与 | 付与対象 0:ユーザー変形値／IKリンク／多重付与 1:親のローカル変形量
	IS_EXTERNAL_LOCAL BoneFlag = 0x0080
	// 回転付与
	IS_EXTERNAL_ROTATION BoneFlag = 0x0100
	// 移動付与
	IS_EXTERNAL_TRANSLATION BoneFlag = 0x0200
	// 軸固定
	HAS_FIXED_AXIS BoneFlag = 0x0400
	// ローカル軸
	HAS_LOCAL_COORDINATE BoneFlag = 0x0800
	// 物理後変形
	IS_AFTER_PHYSICS_DEFORM BoneFlag = 0x1000
	// 外部親変形
	IS_EXTERNAL_PARENT_DEFORM BoneFlag = 0x2000
)

type BoneCategory int

const (
	// ルート
	CATEGORY_ROOT BoneCategory = iota
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
	// 足首
	CATEGORY_ANKLE BoneCategory = iota
	// 捩
	CATEGORY_TWIST BoneCategory = iota
	// 頭
	CATEGORY_HEAD BoneCategory = iota
	// ローカル軸行列計算で親のキャンセルをさせないボーン
	CATEGORY_NO_LOCAL_CANCEL BoneCategory = iota
)

type BoneConfig struct {
	// 準標準ボーン名
	Name StandardBoneNames
	// 親ボーン名候補リスト
	ParentBoneNames []StandardBoneNames
	// 表示先位置(該当ボーンの位置との相対位置)
	TailPosition mvec3.T
	// 末端ボーン名候補リスト
	TailBoneNames []StandardBoneNames
	// ボーンの特性
	Flag BoneFlag
	// ボーンカテゴリ
	Categories []BoneCategory
}

func NewBoneConfig(
	name StandardBoneNames,
	parentBoneNames []StandardBoneNames,
	tailPosition mvec3.T,
	tailBoneNames []StandardBoneNames,
	flag BoneFlag,
	categories []BoneCategory,
) *BoneConfig {
	return &BoneConfig{
		Name:            name,
		ParentBoneNames: parentBoneNames,
		TailPosition:    tailPosition,
		TailBoneNames:   tailBoneNames,
		Flag:            flag,
		Categories:      categories,
	}
}

type StandardBoneNames string

const (
	ROOT          StandardBoneNames = "全ての親"
	CENTER        StandardBoneNames = "センター"
	GROOVE        StandardBoneNames = "グルーブ"
	WAIST         StandardBoneNames = "腰"
	LOWER         StandardBoneNames = "下半身"
	LEG_CENTER    StandardBoneNames = "足中心"
	UPPER         StandardBoneNames = "上半身"
	UPPER2        StandardBoneNames = "上半身2"
	UPPER3        StandardBoneNames = "上半身3"
	NECK_ROOT     StandardBoneNames = "首根元"
	NECK          StandardBoneNames = "首"
	HEAD          StandardBoneNames = "頭"
	EYES          StandardBoneNames = "両目"
	EYE           StandardBoneNames = "{d}目"
	BUST          StandardBoneNames = "{d}胸"
	SHOULDER_ROOT StandardBoneNames = "{d}肩根元"
	SHOULDER_P    StandardBoneNames = "{d}肩P"
	SHOULDER      StandardBoneNames = "{d}肩"
	SHOULDER_C    StandardBoneNames = "{d}肩C"
	ARM           StandardBoneNames = "{d}腕"
	ARM_TWIST     StandardBoneNames = "{d}腕捩"
	ELBOW         StandardBoneNames = "{d}ひじ"
	WRIST_TWIST   StandardBoneNames = "{d}手捩"
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
	LEG           StandardBoneNames = "{d}足"
	KNEE          StandardBoneNames = "{d}ひざ"
	ANKLE         StandardBoneNames = "{d}足首"
	TOE           StandardBoneNames = "{d}つま先"
	LEG_D         StandardBoneNames = "{d}足D"
	KNEE_D        StandardBoneNames = "{d}ひざD"
	ANKLE_D       StandardBoneNames = "{d}足首D"
	TOE_EX        StandardBoneNames = "{d}足先EX"
	LEG_IK_PARENT StandardBoneNames = "{d}足IK親"
	LEG_IK        StandardBoneNames = "{d}足IK"
	TOE_IK        StandardBoneNames = "{d}つま先IK"
)

func (s StandardBoneNames) String() string {
	return string(s)
}

func (s StandardBoneNames) Right() string {
	return strings.ReplaceAll(string(s), "{d}", "右")
}

func (s StandardBoneNames) Left() string {
	return strings.ReplaceAll(string(s), "{d}", "左")
}

var StandardBoneConfigs = map[StandardBoneNames]BoneConfig{
	ROOT: *NewBoneConfig(ROOT,
		[]StandardBoneNames{},
		mvec3.UnitY,
		[]StandardBoneNames{CENTER},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_ROOT}),
	CENTER: *NewBoneConfig(CENTER,
		[]StandardBoneNames{ROOT},
		mvec3.UnitY,
		[]StandardBoneNames{GROOVE, WAIST, UPPER, LOWER},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_ROOT}),
	GROOVE: *NewBoneConfig(GROOVE,
		[]StandardBoneNames{CENTER},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{WAIST, UPPER, LOWER},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_ROOT}),
	WAIST: *NewBoneConfig(WAIST,
		[]StandardBoneNames{GROOVE, CENTER},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{UPPER, LOWER},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_ROOT}),
	LOWER: *NewBoneConfig(LOWER,
		[]StandardBoneNames{WAIST, GROOVE, CENTER},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{LEG_CENTER},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_LOWER}),
	LEG_CENTER: *NewBoneConfig(LEG_CENTER,
		[]StandardBoneNames{LOWER},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{LEG},
		BoneFlag(CAN_ROTATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_LOWER}),
	UPPER: *NewBoneConfig(UPPER,
		[]StandardBoneNames{WAIST, GROOVE, CENTER},
		mvec3.UnitY,
		[]StandardBoneNames{UPPER2},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER}),
	UPPER2: *NewBoneConfig(UPPER2,
		[]StandardBoneNames{UPPER},
		mvec3.UnitY,
		[]StandardBoneNames{UPPER3, NECK_ROOT},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER}),
	UPPER3: *NewBoneConfig(UPPER3,
		[]StandardBoneNames{UPPER2},
		mvec3.UnitY,
		[]StandardBoneNames{NECK_ROOT},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER}),
	NECK_ROOT: *NewBoneConfig(NECK_ROOT,
		[]StandardBoneNames{UPPER3, UPPER2, UPPER},
		mvec3.UnitY,
		[]StandardBoneNames{NECK},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER}),
	NECK: *NewBoneConfig(NECK,
		[]StandardBoneNames{NECK_ROOT, UPPER3, UPPER2, UPPER},
		mvec3.UnitY,
		[]StandardBoneNames{HEAD},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}),
	HEAD: *NewBoneConfig(HEAD,
		[]StandardBoneNames{NECK},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}),
	EYES: *NewBoneConfig(EYES,
		[]StandardBoneNames{HEAD},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}),
	EYE: *NewBoneConfig(EYE,
		[]StandardBoneNames{HEAD},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|IS_EXTERNAL_ROTATION),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_HEAD}),
	BUST: *NewBoneConfig(BUST,
		[]StandardBoneNames{UPPER3, UPPER2, UPPER},
		mvec3.UnitZ.Inverted(),
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER}),
	SHOULDER_ROOT: *NewBoneConfig(SHOULDER_ROOT,
		[]StandardBoneNames{NECK_ROOT},
		mvec3.UnitY,
		[]StandardBoneNames{SHOULDER_P, SHOULDER},
		BoneFlag(CAN_ROTATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	SHOULDER_P: *NewBoneConfig(SHOULDER_P,
		[]StandardBoneNames{SHOULDER_ROOT},
		mvec3.UnitY,
		[]StandardBoneNames{SHOULDER},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	SHOULDER: *NewBoneConfig(SHOULDER,
		[]StandardBoneNames{SHOULDER_P, SHOULDER_ROOT},
		mvec3.UnitY,
		[]StandardBoneNames{SHOULDER_C, ARM},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	SHOULDER_C: *NewBoneConfig(SHOULDER_C,
		[]StandardBoneNames{SHOULDER},
		mvec3.UnitY,
		[]StandardBoneNames{ARM},
		BoneFlag(CAN_ROTATE|IS_EXTERNAL_ROTATION),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	ARM: *NewBoneConfig(ARM,
		[]StandardBoneNames{SHOULDER_C, SHOULDER},
		mvec3.UnitY,
		[]StandardBoneNames{ARM_TWIST, ELBOW},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	ARM_TWIST: *NewBoneConfig(ARM_TWIST,
		[]StandardBoneNames{ARM},
		mvec3.UnitY,
		[]StandardBoneNames{ELBOW},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|HAS_FIXED_AXIS),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST, CATEGORY_NO_LOCAL_CANCEL}),
	ELBOW: *NewBoneConfig(ELBOW,
		[]StandardBoneNames{ARM_TWIST, ARM},
		mvec3.UnitY,
		[]StandardBoneNames{WRIST_TWIST, WRIST},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	WRIST_TWIST: *NewBoneConfig(WRIST_TWIST,
		[]StandardBoneNames{ELBOW},
		mvec3.UnitY,
		[]StandardBoneNames{WRIST},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|HAS_FIXED_AXIS),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_TWIST, CATEGORY_NO_LOCAL_CANCEL}),
	WRIST: *NewBoneConfig(WRIST,
		[]StandardBoneNames{WRIST_TWIST, ELBOW},
		mvec3.UnitY,
		[]StandardBoneNames{INDEX1, MIDDLE1, RING1, PINKY1},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM}),
	THUMB0: *NewBoneConfig(THUMB0,
		[]StandardBoneNames{WRIST},
		mvec3.UnitY,
		[]StandardBoneNames{THUMB1},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	THUMB1: *NewBoneConfig(THUMB1,
		[]StandardBoneNames{THUMB0},
		mvec3.UnitY,
		[]StandardBoneNames{THUMB2},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	THUMB2: *NewBoneConfig(THUMB2,
		[]StandardBoneNames{THUMB1},
		mvec3.UnitY,
		[]StandardBoneNames{THUMB_TAIL},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	THUMB_TAIL: *NewBoneConfig(THUMB_TAIL,
		[]StandardBoneNames{THUMB2},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	INDEX1: *NewBoneConfig(INDEX1,
		[]StandardBoneNames{WRIST},
		mvec3.UnitY,
		[]StandardBoneNames{INDEX2},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	INDEX2: *NewBoneConfig(INDEX2,
		[]StandardBoneNames{INDEX1},
		mvec3.UnitY,
		[]StandardBoneNames{INDEX3},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	INDEX3: *NewBoneConfig(INDEX3,
		[]StandardBoneNames{INDEX2},
		mvec3.UnitY,
		[]StandardBoneNames{INDEX_TAIL},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	INDEX_TAIL: *NewBoneConfig(INDEX_TAIL,
		[]StandardBoneNames{INDEX3},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	MIDDLE1: *NewBoneConfig(MIDDLE1,
		[]StandardBoneNames{WRIST},
		mvec3.UnitY,
		[]StandardBoneNames{MIDDLE2},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	MIDDLE2: *NewBoneConfig(MIDDLE2,
		[]StandardBoneNames{MIDDLE1},
		mvec3.UnitY,
		[]StandardBoneNames{MIDDLE3},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	MIDDLE3: *NewBoneConfig(MIDDLE3,
		[]StandardBoneNames{MIDDLE2},
		mvec3.UnitY,
		[]StandardBoneNames{MIDDLE_TAIL},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	MIDDLE_TAIL: *NewBoneConfig(MIDDLE_TAIL,
		[]StandardBoneNames{MIDDLE3},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	RING1: *NewBoneConfig(RING1,
		[]StandardBoneNames{WRIST},
		mvec3.UnitY,
		[]StandardBoneNames{RING2},

		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	RING2: *NewBoneConfig(RING2,
		[]StandardBoneNames{RING1},
		mvec3.UnitY,
		[]StandardBoneNames{RING3},

		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	RING3: *NewBoneConfig(RING3,
		[]StandardBoneNames{RING2},
		mvec3.UnitY,

		[]StandardBoneNames{RING_TAIL},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	RING_TAIL: *NewBoneConfig(RING_TAIL,
		[]StandardBoneNames{RING3},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	PINKY1: *NewBoneConfig(PINKY1,
		[]StandardBoneNames{WRIST},
		mvec3.UnitY,
		[]StandardBoneNames{PINKY2},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	PINKY2: *NewBoneConfig(PINKY2,
		[]StandardBoneNames{PINKY1},
		mvec3.UnitY,
		[]StandardBoneNames{PINKY3},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	PINKY3: *NewBoneConfig(PINKY3,
		[]StandardBoneNames{PINKY2},
		mvec3.UnitY,
		[]StandardBoneNames{PINKY_TAIL},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|TAIL_IS_BONE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	PINKY_TAIL: *NewBoneConfig(PINKY_TAIL,
		[]StandardBoneNames{PINKY3},
		mvec3.UnitY,
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_UPPER, CATEGORY_ARM, CATEGORY_FINGER, CATEGORY_NO_LOCAL_CANCEL}),
	WAIST_CANCEL: *NewBoneConfig(WAIST_CANCEL,
		[]StandardBoneNames{WAIST},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|IS_EXTERNAL_ROTATION),
		[]BoneCategory{CATEGORY_LOWER}),
	LEG: *NewBoneConfig(LEG,
		[]StandardBoneNames{WAIST_CANCEL, LEG_CENTER, LOWER},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{KNEE},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK}),
	KNEE: *NewBoneConfig(KNEE,
		[]StandardBoneNames{LEG},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{ANKLE},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK}),
	ANKLE: *NewBoneConfig(ANKLE,
		[]StandardBoneNames{KNEE},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{TOE},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE}),
	TOE: *NewBoneConfig(TOE,
		[]StandardBoneNames{ANKLE},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_FK, CATEGORY_ANKLE}),
	LEG_D: *NewBoneConfig(LEG_D,
		[]StandardBoneNames{WAIST_CANCEL, LEG_CENTER, LOWER},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{KNEE_D},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D}),
	KNEE_D: *NewBoneConfig(KNEE_D,
		[]StandardBoneNames{LEG_D},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{ANKLE_D},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D}),
	ANKLE_D: *NewBoneConfig(ANKLE_D,
		[]StandardBoneNames{KNEE_D},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{TOE_EX},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE}),
	TOE_EX: *NewBoneConfig(TOE_EX,
		[]StandardBoneNames{ANKLE_D},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_LOWER, CATEGORY_LEG, CATEGORY_LEG_D, CATEGORY_ANKLE, CATEGORY_NO_LOCAL_CANCEL}),
	LEG_IK_PARENT: *NewBoneConfig(LEG_IK_PARENT,
		[]StandardBoneNames{ROOT},
		mvec3.UnitY,
		[]StandardBoneNames{LEG_IK},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE),
		[]BoneCategory{CATEGORY_LOWER}),
	LEG_IK: *NewBoneConfig(LEG_IK,
		[]StandardBoneNames{LEG_IK_PARENT, ROOT},
		mvec3.UnitY,
		[]StandardBoneNames{TOE_IK},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE|IS_IK|TAIL_IS_BONE),
		[]BoneCategory{CATEGORY_LOWER}),
	TOE_IK: *NewBoneConfig(TOE_IK,
		[]StandardBoneNames{LEG_IK},
		mvec3.UnitY.Inverted(),
		[]StandardBoneNames{},
		BoneFlag(CAN_ROTATE|CAN_TRANSLATE|CAN_MANIPULATE|IS_VISIBLE|IS_IK),
		[]BoneCategory{CATEGORY_LOWER}),
}
