package bone

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mmat4"
	"github.com/miu200521358/mlib_go/pkg/math/mrotation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type IkLink struct {
	// リンクボーンのボーンIndex
	BoneIndex int
	// 角度制限有無
	AngleLimit bool
	// 下限
	MinAngleLimit mrotation.T
	// 上限
	MaxAngleLimit mrotation.T
	// ローカル軸の角度制限有無
	LocalAngleLimit bool
	// ローカル軸制限の下限
	LocalMinAngleLimit mrotation.T
	// ローカル軸制限の上限
	LocalMaxAngleLimit mrotation.T
}

func NewIkLink() *IkLink {
	return &IkLink{
		BoneIndex:          -1,
		AngleLimit:         false,
		MinAngleLimit:      mrotation.T{},
		MaxAngleLimit:      mrotation.T{},
		LocalAngleLimit:    false,
		LocalMinAngleLimit: mrotation.T{},
		LocalMaxAngleLimit: mrotation.T{},
	}
}

func (t *IkLink) Copy() *IkLink {
	copied := *t
	return &copied
}

type Ik struct {
	// IKターゲットボーンのボーンIndex
	BoneIndex int
	// IKループ回数 (最大255)
	LoopCount int
	// IKループ計算時の1回あたりの制限角度(Xにのみ値が入っている)
	UnitRotation mrotation.T
	// IKリンクリスト
	Links []IkLink
}

func NewIk() *Ik {
	return &Ik{
		BoneIndex:    -1,
		LoopCount:    0,
		UnitRotation: mrotation.T{},
		Links:        []IkLink{},
	}
}

func (t *Ik) Copy() *Ik {
	copied := *t
	copied.Links = make([]IkLink, len(t.Links))
	for i, link := range t.Links {
		copied.Links[i] = *link.Copy()
	}
	return &copied
}

type Bone struct {
	*index_model.IndexModel
	// ボーン名
	Name string
	// ボーン名英
	EnglishName string
	// ボーン位置
	Position mvec3.T
	// 親ボーンのボーンIndex
	ParentIndex int
	// 変形階層
	Layer int
	// ボーンフラグ(16bit) 各bit 0:OFF 1:ON
	BoneFlag BoneFlag
	// 接続先:0 の場合 座標オフセット, ボーン位置からの相対分
	TailPosition mvec3.T
	// 接続先:1 の場合 接続先ボーンのボーンIndex
	TailIndex int
	// 回転付与:1 または 移動付与:1 の場合 付与親ボーンのボーンIndex
	EffectIndex int
	// 付与率
	EffectFactor float64
	// 軸固定:1 の場合 軸の方向ベクトル
	FixedAxis mvec3.T
	// ローカル軸:1 の場合 X軸の方向ベクトル
	LocalAxisX mvec3.T
	// ローカル軸:1 の場合 Z軸の方向ベクトル
	LocalAxisZ mvec3.T
	// 外部親変形:1 の場合 Key値
	ExternalKey int
	// IK:1 の場合 IKデータを格納
	Ik *Ik
	// 該当表示枠
	DisplaySlot int
	// システム計算用追加ボーン の場合 true
	IsSystem bool
	// 計算済みのX軸の方向ベクトル
	NormalizedLocalAxisX mvec3.T
	// 計算済みのY軸の方向ベクトル
	NormalizedLocalAxisY mvec3.T
	// 計算済みのZ軸の方向ベクトル
	NormalizedLocalAxisZ mvec3.T
	// 計算済みの軸制限ベクトル
	NormalizedFixedAxis mvec3.T
	// ローカル軸の方向ベクトル(CorrectedLocalXVectorの正規化ベクトル)
	LocalAxis mvec3.T
	// 親ボーンからの相対位置
	ParentRelativePosition mvec3.T
	// Tailボーンへの相対位置
	TailRelativePosition mvec3.T
	// 逆オフセット行列(親ボーンからの相対位置分を戻す)
	ParentRevertMatrix mmat4.T
	// オフセット行列 (自身の位置を原点に戻す行列)
	OffsetMatrix mmat4.T
	// 自分のボーンまでのボーンIndexのリスト
	TreeBoneIndexes []int
	// 関連ボーンINDEX一覧（付与親とかIKとか）
	RelativeBoneIndexes []int
	// 自分を親として登録しているボーンINDEX一覧
	ChildBoneIndexes []int
	// 自分を付与親として登録しているボーンINDEX一覧
	EffectiveBoneIndexes []int
	// IKリンクとして登録されているIKボーンのボーンIndex
	IkLinkBoneIndexes []int
	// IKターゲットとして登録されているIKボーンのボーンIndex
	IkTargetBoneIndexes []int
	// 自分がIKリンクボーンの角度制限がある場合、true
	AngleLimit bool
	// 自分がIKリンクボーンの角度制限の下限
	MinAngleLimit mrotation.T
	// 自分がIKリンクボーンの角度制限の上限
	MaxAngleLimit mrotation.T
	// 自分がIKリンクボーンのローカル軸角度制限がある場合、true
	LocalAngleLimit bool
	// 自分がIKリンクボーンのローカル軸角度制限の下限
	LocalMinAngleLimit mrotation.T
	// 自分がIKリンクボーンのローカル軸角度制限の上限
	LocalMaxAngleLimit mrotation.T
}

func NewBone() *Bone {
	bone := &Bone{
		IndexModel:             &index_model.IndexModel{Index: -1},
		Name:                   "",
		EnglishName:            "",
		Position:               mvec3.T{},
		ParentIndex:            -1,
		Layer:                  -1,
		BoneFlag:               BONE_FLAG_NONE,
		TailPosition:           mvec3.T{},
		TailIndex:              -1,
		EffectIndex:            -1,
		EffectFactor:           0.0,
		FixedAxis:              mvec3.T{},
		LocalAxisX:             mvec3.T{},
		LocalAxisZ:             mvec3.T{},
		ExternalKey:            -1,
		Ik:                     NewIk(),
		DisplaySlot:            -1,
		IsSystem:               true,
		NormalizedLocalAxisX:   mvec3.T{},
		NormalizedLocalAxisY:   mvec3.T{},
		NormalizedLocalAxisZ:   mvec3.T{},
		LocalAxis:              mvec3.UnitX,
		IkLinkBoneIndexes:      []int{},
		IkTargetBoneIndexes:    []int{},
		ParentRelativePosition: mvec3.T{},
		TailRelativePosition:   mvec3.T{},
		NormalizedFixedAxis:    mvec3.T{},
		TreeBoneIndexes:        []int{},
		ParentRevertMatrix:     mmat4.Ident,
		OffsetMatrix:           mmat4.Ident,
		RelativeBoneIndexes:    []int{},
		ChildBoneIndexes:       []int{},
		EffectiveBoneIndexes:   []int{},
		AngleLimit:             false,
		MinAngleLimit:          mrotation.T{},
		MaxAngleLimit:          mrotation.T{},
		LocalAngleLimit:        false,
		LocalMinAngleLimit:     mrotation.T{},
		LocalMaxAngleLimit:     mrotation.T{},
	}
	bone.NormalizedLocalAxisX = *bone.LocalAxisX.Copy()
	bone.NormalizedLocalAxisZ = *bone.LocalAxisZ.Copy()
	bone.NormalizedLocalAxisY = *bone.NormalizedLocalAxisZ.Cross(&bone.NormalizedLocalAxisX)
	bone.NormalizedFixedAxis = *bone.FixedAxis.Copy()
	return bone
}

// Copy
func (t *Bone) Copy() index_model.IndexModelInterface {
	copied := *t
	copied.Ik = t.Ik.Copy()
	copied.Position = *t.Position.Copy()
	copied.TailPosition = *t.TailPosition.Copy()
	copied.FixedAxis = *t.FixedAxis.Copy()
	copied.LocalAxisX = *t.LocalAxisX.Copy()
	copied.LocalAxisZ = *t.LocalAxisZ.Copy()
	copied.NormalizedLocalAxisZ = *t.NormalizedLocalAxisZ.Copy()
	copied.NormalizedLocalAxisX = *t.NormalizedLocalAxisX.Copy()
	copied.NormalizedLocalAxisY = *t.NormalizedLocalAxisY.Copy()
	copied.LocalAxis = *t.LocalAxis.Copy()
	copied.ParentRelativePosition = *t.ParentRelativePosition.Copy()
	copied.TailRelativePosition = *t.TailRelativePosition.Copy()
	copied.NormalizedFixedAxis = *t.NormalizedFixedAxis.Copy()
	copied.IkLinkBoneIndexes = make([]int, len(t.IkLinkBoneIndexes))
	copy(copied.IkLinkBoneIndexes, t.IkLinkBoneIndexes)
	copied.IkTargetBoneIndexes = make([]int, len(t.IkTargetBoneIndexes))
	copy(copied.IkTargetBoneIndexes, t.IkTargetBoneIndexes)
	copied.TreeBoneIndexes = make([]int, len(t.TreeBoneIndexes))
	copy(copied.TreeBoneIndexes, t.TreeBoneIndexes)
	copied.ParentRevertMatrix = *t.ParentRevertMatrix.Copy()
	copied.OffsetMatrix = *t.OffsetMatrix.Copy()
	copied.RelativeBoneIndexes = make([]int, len(t.RelativeBoneIndexes))
	copy(copied.RelativeBoneIndexes, t.RelativeBoneIndexes)
	copied.ChildBoneIndexes = make([]int, len(t.ChildBoneIndexes))
	copy(copied.ChildBoneIndexes, t.ChildBoneIndexes)
	copied.EffectiveBoneIndexes = make([]int, len(t.EffectiveBoneIndexes))
	copy(copied.EffectiveBoneIndexes, t.EffectiveBoneIndexes)
	copied.MinAngleLimit = *t.MinAngleLimit.Copy()
	copied.MaxAngleLimit = *t.MaxAngleLimit.Copy()
	copied.LocalMinAngleLimit = *t.LocalMinAngleLimit.Copy()
	copied.LocalMaxAngleLimit = *t.LocalMaxAngleLimit.Copy()
	return &copied
}

func (bone *Bone) NormalizeFixedAxis(fixedAxis mvec3.T) {
	bone.NormalizedFixedAxis = fixedAxis.Normalized()
}

func (bone *Bone) NormalizeLocalAxis(localAxisX mvec3.T) {
	bone.NormalizedLocalAxisX = localAxisX.Normalized()
	bone.NormalizedLocalAxisY = *bone.NormalizedLocalAxisX.Cross(mvec3.UnitZ.Invert())
	bone.NormalizedLocalAxisZ = *bone.NormalizedLocalAxisX.Cross(&bone.NormalizedLocalAxisY)
}

// 表示先がボーンであるか
func (bone *Bone) IsTailBone() bool {
	return bone.BoneFlag&BONE_FLAG_TAIL_IS_BONE != 0
}

// 回転可能であるか
func (bone *Bone) CanRotate() bool {
	return bone.BoneFlag&BONE_FLAG_CAN_ROTATE != 0
}

// 移動可能であるか
func (bone *Bone) CanTranslate() bool {
	return bone.BoneFlag&BONE_FLAG_CAN_TRANSLATE != 0
}

// 表示であるか
func (bone *Bone) IsVisible() bool {
	return bone.BoneFlag&BONE_FLAG_IS_VISIBLE != 0
}

// 操作可であるか
func (bone *Bone) CanManipulate() bool {
	return bone.BoneFlag&BONE_FLAG_CAN_MANIPULATE != 0
}

// IKであるか
func (bone *Bone) IsIK() bool {
	return bone.BoneFlag&BONE_FLAG_IS_IK != 0
}

// ローカル付与であるか
func (bone *Bone) IsExternalLocal() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_LOCAL != 0
}

// 回転付与であるか
func (bone *Bone) IsExternalRotation() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_ROTATION != 0
}

// 移動付与であるか
func (bone *Bone) IsExternalTranslation() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_TRANSLATION != 0
}

// 軸固定であるか
func (bone *Bone) HasFixedAxis() bool {
	return bone.BoneFlag&BONE_FLAG_HAS_FIXED_AXIS != 0
}

// ローカル軸を持つか
func (bone *Bone) HasLocalAxis() bool {
	return bone.BoneFlag&BONE_FLAG_HAS_LOCAL_AXIS != 0
}

// 物理後変形であるか
func (bone *Bone) IsAfterPhysicsDeform() bool {
	return bone.BoneFlag&BONE_FLAG_IS_AFTER_PHYSICS_DEFORM != 0
}

// 外部親変形であるか
func (bone *Bone) IsExternalParentDeform() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM != 0
}

// 足D系列であるか
func (bone *Bone) IsLegD() bool {
	return bone.containsCategory(CATEGORY_LEG_D)
}

// 肩P系列であるか
func (bone *Bone) IsShoulderP() bool {
	return bone.containsCategory(CATEGORY_SHOULDER_P)
}

// 足FK系列であるか
func (bone *Bone) IsLegFK() bool {
	return bone.containsCategory(CATEGORY_LEG_FK)
}

// 足首から先であるか
func (bone *Bone) IsAnkle() bool {
	return bone.containsCategory(CATEGORY_ANKLE)
}

// 捩りボーンであるか
func (bone *Bone) IsTwist() bool {
	return bone.containsCategory(CATEGORY_TWIST)
}

// 腕系ボーンであるか(指は含まない)
func (bone *Bone) IsArm() bool {
	return bone.containsCategory(CATEGORY_ARM)
}

// 指系ボーンであるか
func (bone *Bone) IsFinger() bool {
	return bone.containsCategory(CATEGORY_FINGER)
}

// 頭系であるか
func (bone *Bone) IsHead() bool {
	return bone.containsCategory(CATEGORY_HEAD)
}

// 下半身系であるか
func (bone *Bone) IsLower() bool {
	return bone.containsCategory(CATEGORY_LOWER)
}

// 上半身系であるか
func (bone *Bone) IsUpper() bool {
	return bone.containsCategory(CATEGORY_UPPER)
}

// ローカル軸行列計算で親のキャンセルをさせないボーンであるか
func (bone *Bone) IsNoLocalCancel() bool {
	// 捩り分散ボーンも含む
	if strings.Contains(bone.Name, "捩") {
		return true
	}

	return bone.containsCategory(CATEGORY_NO_LOCAL_CANCEL)
}

// 指定したカテゴリーに属するか
func (bone *Bone) containsCategory(category BoneCategory) bool {
	for _, boneConfig := range StandardBoneConfigs {
		for _, c := range boneConfig.Categories {
			if c == category && (boneConfig.Name.String() == bone.Name ||
				boneConfig.Name.Right() == bone.Name ||
				boneConfig.Name.Left() == bone.Name) {
				return true
			}
		}
	}
	return false
}

// ボーンリスト
type Bones struct {
	*index_model.IndexModelCorrection[*Bone]
}

func NewBones() *Bones {
	return &Bones{
		IndexModelCorrection: index_model.NewIndexModelCorrection[*Bone](),
	}
}
