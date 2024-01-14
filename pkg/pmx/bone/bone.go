package bone

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/core/index_model"
	"github.com/miu200521358/mlib_go/pkg/math/mmat4"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

type IkLink struct {
	BoneIndex          int
	AngleLimit         bool
	MinAngleLimit      mvec3.T
	MaxAngleLimit      mvec3.T
	LocalAngleLimit    bool
	LocalMinAngleLimit mvec3.T
	LocalMaxAngleLimit mvec3.T
}

func (ikLink *IkLink) IsValid() bool {
	return ikLink.BoneIndex >= 0
}

func (t *IkLink) Copy() *IkLink {
	copied := *t
	return &copied
}

type Ik struct {
	BoneIndex    int
	LoopCount    int
	UnitRotation float64
	Links        []IkLink
}

func (ik *Ik) IsValid() bool {
	return ik.BoneIndex >= 0
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
	index_model.IndexModel
	Index                  int
	Name                   string
	EnglishName            string
	Position               *mvec3.T
	ParentIndex            int
	Layer                  int
	BoneFlag               BoneFlag
	TailPosition           *mvec3.T
	TailIndex              int
	EffectIndex            int
	EffectFactor           float64
	FixedAxis              *mvec3.T
	LocalXVector           *mvec3.T
	LocalZVector           *mvec3.T
	ExternalKey            int
	Ik                     *Ik
	DisplaySlot            int
	IsSystem               bool
	CorrectedLocalYVector  *mvec3.T
	CorrectedLocalZVector  *mvec3.T
	CorrectedLocalXVector  *mvec3.T
	LocalAxis              *mvec3.T
	LocalMatrix            *mmat4.T
	IkLinkIndexes          []int
	IkTargetIndexes        []int
	ParentRelativePosition *mvec3.T
	TailRelativePosition   *mvec3.T
	CorrectedFixedAxis     *mvec3.T
	TreeIndexes            []int
	ParentRevertMatrix     *mmat4.T
	OffsetMatrix           *mmat4.T
	RelativeBoneIndexes    []int
	ChildBoneIndexes       []int
	EffectiveTargetIndexes []int
	AngleLimit             bool
	MinAngleLimit          *mvec3.T
	MaxAngleLimit          *mvec3.T
	LocalAngleLimit        bool
	LocalMinAngleLimit     *mvec3.T
	LocalMaxAngleLimit     *mvec3.T
}

func NewBone(
	name string,
	englishName string,
	position *mvec3.T,
	parentIndex int,
	layer int,
	boneFlag BoneFlag,
	tailPosition *mvec3.T,
	tailIndex int,
	effectIndex int,
	effectFactor float64,
	fixedAxis *mvec3.T,
	localXVector *mvec3.T,
	localZVector *mvec3.T,
	externalKey int,
	ik *Ik,
	displaySlot int,
	isSystem bool,
) *Bone {
	bone := &Bone{
		Index:                  0,
		Name:                   name,
		EnglishName:            englishName,
		Position:               position,
		ParentIndex:            parentIndex,
		Layer:                  layer,
		BoneFlag:               boneFlag,
		TailPosition:           tailPosition,
		TailIndex:              tailIndex,
		EffectIndex:            effectIndex,
		EffectFactor:           effectFactor,
		FixedAxis:              fixedAxis,
		LocalXVector:           localXVector,
		LocalZVector:           localZVector,
		ExternalKey:            externalKey,
		Ik:                     ik,
		DisplaySlot:            displaySlot,
		IsSystem:               isSystem,
		CorrectedLocalYVector:  &mvec3.T{},
		CorrectedLocalZVector:  &mvec3.T{},
		CorrectedLocalXVector:  &mvec3.T{},
		LocalAxis:              &mvec3.UnitX,
		LocalMatrix:            &mmat4.Ident,
		IkLinkIndexes:          []int{},
		IkTargetIndexes:        []int{},
		ParentRelativePosition: &mvec3.T{},
		TailRelativePosition:   &mvec3.T{},
		CorrectedFixedAxis:     &mvec3.T{},
		TreeIndexes:            []int{},
		ParentRevertMatrix:     &mmat4.Ident,
		OffsetMatrix:           &mmat4.Ident,
		RelativeBoneIndexes:    []int{},
		ChildBoneIndexes:       []int{},
		EffectiveTargetIndexes: []int{},
		AngleLimit:             false,
		MinAngleLimit:          &mvec3.T{},
		MaxAngleLimit:          &mvec3.T{},
		LocalAngleLimit:        false,
		LocalMinAngleLimit:     &mvec3.T{},
		LocalMaxAngleLimit:     &mvec3.T{},
	}
	bone.CorrectedLocalXVector = bone.LocalXVector.Copy()
	bone.CorrectedLocalYVector = bone.LocalZVector.Cross(bone.CorrectedLocalXVector)
	bone.CorrectedLocalZVector = bone.CorrectedLocalXVector.Cross(bone.CorrectedLocalYVector)
	bone.CorrectedFixedAxis = bone.FixedAxis.Copy()
	return bone
}

// Copy
func (t *Bone) Copy() index_model.IndexModelInterface {
	copied := *t
	copied.Ik = t.Ik.Copy()
	copied.Position = t.Position.Copy()
	copied.TailPosition = t.TailPosition.Copy()
	copied.FixedAxis = t.FixedAxis.Copy()
	copied.LocalXVector = t.LocalXVector.Copy()
	copied.LocalZVector = t.LocalZVector.Copy()
	copied.CorrectedLocalXVector = t.CorrectedLocalXVector.Copy()
	copied.CorrectedLocalYVector = t.CorrectedLocalYVector.Copy()
	copied.CorrectedLocalZVector = t.CorrectedLocalZVector.Copy()
	copied.LocalAxis = t.LocalAxis.Copy()
	copied.LocalMatrix = t.LocalMatrix.Copy()
	copied.ParentRelativePosition = t.ParentRelativePosition.Copy()
	copied.TailRelativePosition = t.TailRelativePosition.Copy()
	copied.CorrectedFixedAxis = t.CorrectedFixedAxis.Copy()
	copied.IkLinkIndexes = make([]int, len(t.IkLinkIndexes))
	copy(copied.IkLinkIndexes, t.IkLinkIndexes)
	copied.IkTargetIndexes = make([]int, len(t.IkTargetIndexes))
	copy(copied.IkTargetIndexes, t.IkTargetIndexes)
	copied.TreeIndexes = make([]int, len(t.TreeIndexes))
	copy(copied.TreeIndexes, t.TreeIndexes)
	copied.ParentRevertMatrix = t.ParentRevertMatrix.Copy()
	copied.OffsetMatrix = t.OffsetMatrix.Copy()
	copied.RelativeBoneIndexes = make([]int, len(t.RelativeBoneIndexes))
	copy(copied.RelativeBoneIndexes, t.RelativeBoneIndexes)
	copied.ChildBoneIndexes = make([]int, len(t.ChildBoneIndexes))
	copy(copied.ChildBoneIndexes, t.ChildBoneIndexes)
	copied.EffectiveTargetIndexes = make([]int, len(t.EffectiveTargetIndexes))
	copy(copied.EffectiveTargetIndexes, t.EffectiveTargetIndexes)
	copied.MinAngleLimit = t.MinAngleLimit.Copy()
	copied.MaxAngleLimit = t.MaxAngleLimit.Copy()
	copied.LocalMinAngleLimit = t.LocalMinAngleLimit.Copy()
	copied.LocalMaxAngleLimit = t.LocalMaxAngleLimit.Copy()
	return &copied
}

func (bone *Bone) CorrectFixedAxis(correctedFixedAxis mvec3.T) {
	bone.CorrectedFixedAxis = correctedFixedAxis.Normalize()
}

func (bone *Bone) CorrectLocalVector(correctedLocalXVector mvec3.T) {
	bone.CorrectedLocalXVector = correctedLocalXVector.Normalize()
	bone.CorrectedLocalYVector = bone.CorrectedLocalXVector.Cross(mvec3.UnitZ.Invert())
	bone.CorrectedLocalZVector = bone.CorrectedLocalXVector.Cross(bone.CorrectedLocalYVector)
}

// 表示先がボーンであるか
func (bone *Bone) IsTailBone() bool {
	return bone.BoneFlag&TAIL_IS_BONE != 0
}

// 回転可能であるか
func (bone *Bone) CanRotate() bool {
	return bone.BoneFlag&CAN_ROTATE != 0
}

// 移動可能であるか
func (bone *Bone) CanTranslate() bool {
	return bone.BoneFlag&CAN_TRANSLATE != 0
}

// 表示であるか
func (bone *Bone) IsVisible() bool {
	return bone.BoneFlag&IS_VISIBLE != 0
}

// 操作可であるか
func (bone *Bone) CanManipulate() bool {
	return bone.BoneFlag&CAN_MANIPULATE != 0
}

// IKであるか
func (bone *Bone) IsIK() bool {
	return bone.BoneFlag&IS_IK != 0
}

// ローカル付与であるか
func (bone *Bone) IsExternalLocal() bool {
	return bone.BoneFlag&IS_EXTERNAL_LOCAL != 0
}

// 回転付与であるか
func (bone *Bone) IsExternalRotation() bool {
	return bone.BoneFlag&IS_EXTERNAL_ROTATION != 0
}

// 移動付与であるか
func (bone *Bone) IsExternalTranslation() bool {
	return bone.BoneFlag&IS_EXTERNAL_TRANSLATION != 0
}

// 軸固定であるか
func (bone *Bone) HasFixedAxis() bool {
	return bone.BoneFlag&HAS_FIXED_AXIS != 0
}

// ローカル軸を持つか
func (bone *Bone) HasLocalCoordinate() bool {
	return bone.BoneFlag&HAS_LOCAL_COORDINATE != 0
}

// 物理後変形であるか
func (bone *Bone) IsAfterPhysicsDeform() bool {
	return bone.BoneFlag&IS_AFTER_PHYSICS_DEFORM != 0
}

// 外部親変形であるか
func (bone *Bone) IsExternalParentDeform() bool {
	return bone.BoneFlag&IS_EXTERNAL_PARENT_DEFORM != 0
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
	index_model.IndexModelCorrection[*Bone]
}

func NewBones(name string) *Bones {
	return &Bones{
		IndexModelCorrection: *index_model.NewIndexModelCorrection[*Bone](),
	}
}
