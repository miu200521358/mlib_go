package mmodel

import "strings"

// BoneDirection はボーンの方向を表します。
type BoneDirection string

const (
	BONE_DIRECTION_RIGHT BoneDirection = "右" // 右
	BONE_DIRECTION_LEFT  BoneDirection = "左" // 左
	BONE_DIRECTION_TRUNK BoneDirection = ""  // 体幹
)

// BoneDirectionPrefix は方向を示すプレースホルダです。
const BoneDirectionPrefix = "{d}"

// String は方向を文字列で返します。
func (d BoneDirection) String() string {
	return string(d)
}

// Sign は方向に応じた符号を返します（左:-1, 右:1, 体幹:0）。
func (d BoneDirection) Sign() float64 {
	switch d {
	case BONE_DIRECTION_LEFT:
		return -1.0
	case BONE_DIRECTION_RIGHT:
		return 1.0
	}
	return 0.0
}

// BoneCategory はボーンのカテゴリを表します。
type BoneCategory int

const (
	CATEGORY_ROOT              BoneCategory = iota // ルート
	CATEGORY_TRUNK                                 // 体幹
	CATEGORY_UPPER                                 // 上半身
	CATEGORY_LOWER                                 // 下半身
	CATEGORY_SHOULDER                              // 肩
	CATEGORY_ARM                                   // 腕
	CATEGORY_ELBOW                                 // ひじ
	CATEGORY_LEG                                   // 足
	CATEGORY_FINGER                                // 指
	CATEGORY_TAIL                                  // 先
	CATEGORY_LEG_D                                 // 足D
	CATEGORY_SHOULDER_P                            // 肩P
	CATEGORY_LEG_FK                                // 足FK
	CATEGORY_LEG_IK                                // 足IK
	CATEGORY_ANKLE                                 // 足首
	CATEGORY_SOLE                                  // 靴底
	CATEGORY_TWIST                                 // 捩
	CATEGORY_HEAD                                  // 頭
	CATEGORY_EYE                                   // 目
	CATEGORY_FITTING_ONLY_MOVE                     // フィッティング時に移動のみ
)

// StandardBoneName は標準ボーン名を表します。
type StandardBoneName string

const (
	BONE_ROOT          StandardBoneName = "全ての親"
	BONE_CENTER        StandardBoneName = "センター"
	BONE_GROOVE        StandardBoneName = "グルーブ"
	BONE_BODY_AXIS     StandardBoneName = "体軸"
	BONE_TRUNK_ROOT    StandardBoneName = "体幹中心"
	BONE_WAIST         StandardBoneName = "腰"
	BONE_LOWER_ROOT    StandardBoneName = "下半身根元"
	BONE_LOWER         StandardBoneName = "下半身"
	BONE_HIP           StandardBoneName = "{d}腰骨"
	BONE_LEG_CENTER    StandardBoneName = "足中心"
	BONE_UPPER_ROOT    StandardBoneName = "上半身根元"
	BONE_UPPER         StandardBoneName = "上半身"
	BONE_UPPER2        StandardBoneName = "上半身2"
	BONE_NECK_ROOT     StandardBoneName = "首根元"
	BONE_NECK          StandardBoneName = "首"
	BONE_HEAD          StandardBoneName = "頭"
	BONE_HEAD_TAIL     StandardBoneName = "頭先先"
	BONE_EYES          StandardBoneName = "両目"
	BONE_EYE           StandardBoneName = "{d}目"
	BONE_SHOULDER_ROOT StandardBoneName = "{d}肩根元"
	BONE_SHOULDER_P    StandardBoneName = "{d}肩P"
	BONE_SHOULDER      StandardBoneName = "{d}肩"
	BONE_SHOULDER_C    StandardBoneName = "{d}肩C"
	BONE_ARM           StandardBoneName = "{d}腕"
	BONE_ARM_TWIST     StandardBoneName = "{d}腕捩"
	BONE_ELBOW         StandardBoneName = "{d}ひじ"
	BONE_WRIST_TWIST   StandardBoneName = "{d}手捩"
	BONE_WRIST         StandardBoneName = "{d}手首"
	BONE_WRIST_TAIL    StandardBoneName = "{d}手首先先"
	BONE_THUMB0        StandardBoneName = "{d}親指０"
	BONE_THUMB1        StandardBoneName = "{d}親指１"
	BONE_THUMB2        StandardBoneName = "{d}親指２"
	BONE_THUMB_TAIL    StandardBoneName = "{d}親指先先"
	BONE_INDEX1        StandardBoneName = "{d}人指１"
	BONE_INDEX2        StandardBoneName = "{d}人指２"
	BONE_INDEX3        StandardBoneName = "{d}人指３"
	BONE_INDEX_TAIL    StandardBoneName = "{d}人指先先"
	BONE_MIDDLE1       StandardBoneName = "{d}中指１"
	BONE_MIDDLE2       StandardBoneName = "{d}中指２"
	BONE_MIDDLE3       StandardBoneName = "{d}中指３"
	BONE_MIDDLE_TAIL   StandardBoneName = "{d}中指先先"
	BONE_RING1         StandardBoneName = "{d}薬指１"
	BONE_RING2         StandardBoneName = "{d}薬指２"
	BONE_RING3         StandardBoneName = "{d}薬指３"
	BONE_RING_TAIL     StandardBoneName = "{d}薬指先先"
	BONE_PINKY1        StandardBoneName = "{d}小指１"
	BONE_PINKY2        StandardBoneName = "{d}小指２"
	BONE_PINKY3        StandardBoneName = "{d}小指３"
	BONE_PINKY_TAIL    StandardBoneName = "{d}小指先先"
	BONE_WAIST_CANCEL  StandardBoneName = "腰キャンセル{d}"
	BONE_LEG_ROOT      StandardBoneName = "{d}足根元"
	BONE_LEG           StandardBoneName = "{d}足"
	BONE_KNEE          StandardBoneName = "{d}ひざ"
	BONE_ANKLE         StandardBoneName = "{d}足首"
	BONE_HEEL          StandardBoneName = "{d}かかと"
	BONE_TOE_T         StandardBoneName = "{d}つま先先"
	BONE_LEG_D         StandardBoneName = "{d}足D"
	BONE_KNEE_D        StandardBoneName = "{d}ひざD"
	BONE_ANKLE_D       StandardBoneName = "{d}足首D"
	BONE_TOE_EX        StandardBoneName = "{d}足先EX"
	BONE_LEG_IK_PARENT StandardBoneName = "{d}足IK親"
	BONE_LEG_IK        StandardBoneName = "{d}足ＩＫ"
	BONE_TOE_IK        StandardBoneName = "{d}つま先ＩＫ"
)

// String はボーン名を文字列で返します。
func (s StandardBoneName) String() string {
	return string(s)
}

// StringFromDirection は方向を適用したボーン名を返します。
func (s StandardBoneName) StringFromDirection(d BoneDirection) string {
	return strings.ReplaceAll(string(s), BoneDirectionPrefix, string(d))
}

// Right は右側のボーン名を返します。
func (s StandardBoneName) Right() string {
	return strings.ReplaceAll(string(s), BoneDirectionPrefix, "右")
}

// Left は左側のボーン名を返します。
func (s StandardBoneName) Left() string {
	return strings.ReplaceAll(string(s), BoneDirectionPrefix, "左")
}
