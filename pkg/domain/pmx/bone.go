package pmx

import (
	"slices"
	"sort"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IkLink struct {
	BoneIndex          int              // リンクボーンのボーンIndex
	AngleLimit         bool             // 角度制限有無
	MinAngleLimit      *mmath.MRotation // 下限
	MaxAngleLimit      *mmath.MRotation // 上限
	LocalAngleLimit    bool             // ローカル軸の角度制限有無
	LocalMinAngleLimit *mmath.MRotation // ローカル軸制限の下限
	LocalMaxAngleLimit *mmath.MRotation // ローカル軸制限の上限
}

func NewIkLink() *IkLink {
	return &IkLink{
		BoneIndex:          -1,
		AngleLimit:         false,
		MinAngleLimit:      mmath.NewMRotation(),
		MaxAngleLimit:      mmath.NewMRotation(),
		LocalAngleLimit:    false,
		LocalMinAngleLimit: mmath.NewMRotation(),
		LocalMaxAngleLimit: mmath.NewMRotation(),
	}
}

func (ikLink *IkLink) Copy() *IkLink {
	copied := &IkLink{
		BoneIndex:          ikLink.BoneIndex,
		AngleLimit:         ikLink.AngleLimit,
		MinAngleLimit:      ikLink.MinAngleLimit.Copy(),
		MaxAngleLimit:      ikLink.MaxAngleLimit.Copy(),
		LocalAngleLimit:    ikLink.LocalAngleLimit,
		LocalMinAngleLimit: ikLink.LocalMinAngleLimit.Copy(),
		LocalMaxAngleLimit: ikLink.LocalMaxAngleLimit.Copy(),
	}
	return copied
}

type Ik struct {
	BoneIndex    int              // IKターゲットボーンのボーンIndex
	LoopCount    int              // IKループ回数 (最大255)
	UnitRotation *mmath.MRotation // IKループ計算時の1回あたりの制限角度
	Links        []*IkLink        // IKリンクリスト
}

func NewIk() *Ik {
	return &Ik{
		BoneIndex:    -1,
		LoopCount:    0,
		UnitRotation: mmath.NewMRotation(),
		Links:        []*IkLink{},
	}
}

func (ik *Ik) Copy() *Ik {
	copied := &Ik{}
	copied.BoneIndex = ik.BoneIndex
	copied.LoopCount = ik.LoopCount
	copied.UnitRotation = ik.UnitRotation.Copy()
	copied.Links = make([]*IkLink, len(ik.Links))
	for i, link := range ik.Links {
		copied.Links[i] = link.Copy()
	}
	return copied
}

type Bone struct {
	index        int          // ボーンINDEX
	name         string       // ボーン名
	englishName  string       // ボーン英名
	Position     *mmath.MVec3 // ボーン位置
	ParentIndex  int          // 親ボーンのボーンIndex
	Layer        int          // 変形階層
	BoneFlag     BoneFlag     // ボーンフラグ(16bit) 各bit 0:OFF 1:ON
	TailPosition *mmath.MVec3 // 接続先:0 の場合 座標オフセット, ボーン位置からの相対分
	TailIndex    int          // 接続先:1 の場合 接続先ボーンのボーンIndex
	EffectIndex  int          // 回転付与:1 または 移動付与:1 の場合 付与親ボーンのボーンIndex
	EffectFactor float64      // 付与率
	FixedAxis    *mmath.MVec3 // 軸固定:1 の場合 軸の方向ベクトル
	LocalAxisX   *mmath.MVec3 // ローカル軸:1 の場合 X軸の方向ベクトル
	LocalAxisZ   *mmath.MVec3 // ローカル軸:1 の場合 Z軸の方向ベクトル
	EffectorKey  int          // 外部親変形:1 の場合 Key値
	Ik           *Ik          // IK:1 の場合 IKデータを格納
	DisplaySlot  int          // 該当表示枠
	IsSystem     bool         // システム計算用追加ボーン の場合 true
	Extend       *BoneExtend  // 拡張情報
}

type BoneExtend struct {
	NormalizedLocalAxisX   *mmath.MVec3     // 計算済みのX軸の方向ベクトル
	NormalizedLocalAxisY   *mmath.MVec3     // 計算済みのY軸の方向ベクトル
	NormalizedLocalAxisZ   *mmath.MVec3     // 計算済みのZ軸の方向ベクトル
	NormalizedFixedAxis    *mmath.MVec3     // 計算済みの軸制限ベクトル
	LocalAxis              *mmath.MVec3     // ローカル軸の方向ベクトル(CorrectedLocalXVectorの正規化ベクトル)
	ParentRelativePosition *mmath.MVec3     // 親ボーンからの相対位置
	ChildRelativePosition  *mmath.MVec3     // Tailボーンへの相対位置
	RevertOffsetMatrix     *mmath.MMat4     // 逆オフセット行列(親ボーンからの相対位置分を戻す)
	OffsetMatrix           *mmath.MMat4     // オフセット行列 (自身の位置を原点に戻す行列)
	TreeBoneIndexes        []int            // 自分のボーンまでのボーンIndexのリスト
	ParentBoneIndexes      []int            // 自分の親ボーンからルートまでのボーンIndexのリスト
	ParentBoneNames        []string         // 自分の親ボーンからルートまでのボーン名のリスト
	RelativeBoneIndexes    []int            // 関連ボーンINDEX一覧（付与親とかIKとか）
	ChildBoneIndexes       []int            // 自分を親として登録しているボーンINDEX一覧
	EffectiveBoneIndexes   []int            // 自分を付与親として登録しているボーンINDEX一覧
	IkLinkBoneIndexes      []int            // 自分をIKリンクとして登録されているIKボーンのボーンIndex
	IkTargetBoneIndexes    []int            // 自分をIKターゲットとして登録されているIKボーンのボーンIndex
	AngleLimit             bool             // 自分がIKリンクボーンの角度制限がある場合、true
	MinAngleLimit          *mmath.MRotation // 自分がIKリンクボーンの角度制限の下限
	MaxAngleLimit          *mmath.MRotation // 自分がIKリンクボーンの角度制限の上限
	LocalAngleLimit        bool             // 自分がIKリンクボーンのローカル軸角度制限がある場合、true
	LocalMinAngleLimit     *mmath.MRotation // 自分がIKリンクボーンのローカル軸角度制限の下限
	LocalMaxAngleLimit     *mmath.MRotation // 自分がIKリンクボーンのローカル軸角度制限の上限
	AxisSign               int              // ボーンの軸ベクトル(左は-1, 右は1)
	RigidBody              *RigidBody       // 物理演算用剛体
	OriginalLayer          int              // 元の変形階層
}

func (boneExtend *BoneExtend) Copy() *BoneExtend {
	var copiedNormalizedFixedAxis *mmath.MVec3
	if boneExtend.NormalizedFixedAxis != nil {
		copiedNormalizedFixedAxis = boneExtend.NormalizedFixedAxis.Copy()
	}
	var copiedMinAngleLimit *mmath.MRotation
	if boneExtend.MinAngleLimit != nil {
		copiedMinAngleLimit = boneExtend.MinAngleLimit.Copy()
	}
	var copiedMaxAngleLimit *mmath.MRotation
	if boneExtend.MaxAngleLimit != nil {
		copiedMaxAngleLimit = boneExtend.MaxAngleLimit.Copy()
	}
	var copiedLocalMinAngleLimit *mmath.MRotation
	if boneExtend.LocalMinAngleLimit != nil {
		copiedLocalMinAngleLimit = boneExtend.LocalMinAngleLimit.Copy()
	}
	var copiedLocalMaxAngleLimit *mmath.MRotation
	if boneExtend.LocalMaxAngleLimit != nil {
		copiedLocalMaxAngleLimit = boneExtend.LocalMaxAngleLimit.Copy()
	}

	return &BoneExtend{
		NormalizedLocalAxisX:   boneExtend.NormalizedLocalAxisX.Copy(),
		NormalizedLocalAxisY:   boneExtend.NormalizedLocalAxisY.Copy(),
		NormalizedLocalAxisZ:   boneExtend.NormalizedLocalAxisZ.Copy(),
		NormalizedFixedAxis:    copiedNormalizedFixedAxis,
		LocalAxis:              boneExtend.LocalAxis.Copy(),
		ParentRelativePosition: boneExtend.ParentRelativePosition.Copy(),
		ChildRelativePosition:  boneExtend.ChildRelativePosition.Copy(),
		RevertOffsetMatrix:     boneExtend.RevertOffsetMatrix.Copy(),
		OffsetMatrix:           boneExtend.OffsetMatrix.Copy(),
		TreeBoneIndexes:        mutils.DeepCopyIntSlice(boneExtend.TreeBoneIndexes),
		ParentBoneIndexes:      mutils.DeepCopyIntSlice(boneExtend.ParentBoneIndexes),
		ParentBoneNames:        mutils.DeepCopyStringSlice(boneExtend.ParentBoneNames),
		RelativeBoneIndexes:    mutils.DeepCopyIntSlice(boneExtend.RelativeBoneIndexes),
		ChildBoneIndexes:       mutils.DeepCopyIntSlice(boneExtend.ChildBoneIndexes),
		EffectiveBoneIndexes:   mutils.DeepCopyIntSlice(boneExtend.EffectiveBoneIndexes),
		IkLinkBoneIndexes:      mutils.DeepCopyIntSlice(boneExtend.IkLinkBoneIndexes),
		IkTargetBoneIndexes:    mutils.DeepCopyIntSlice(boneExtend.IkTargetBoneIndexes),
		AngleLimit:             boneExtend.AngleLimit,
		MinAngleLimit:          copiedMinAngleLimit,
		MaxAngleLimit:          copiedMaxAngleLimit,
		LocalAngleLimit:        boneExtend.LocalAngleLimit,
		LocalMinAngleLimit:     copiedLocalMinAngleLimit,
		LocalMaxAngleLimit:     copiedLocalMaxAngleLimit,
		AxisSign:               boneExtend.AxisSign,
		RigidBody:              nil,
	}
}

func NewBone() *Bone {
	bone := &Bone{
		index:        -1,
		name:         "",
		englishName:  "",
		Position:     mmath.NewMVec3(),
		ParentIndex:  -1,
		Layer:        -1,
		BoneFlag:     BONE_FLAG_NONE,
		TailPosition: mmath.NewMVec3(),
		TailIndex:    -1,
		EffectIndex:  -1,
		EffectFactor: 0.0,
		FixedAxis:    mmath.NewMVec3(),
		LocalAxisX:   mmath.NewMVec3(),
		LocalAxisZ:   mmath.NewMVec3(),
		EffectorKey:  -1,
		Ik:           NewIk(),
		DisplaySlot:  -1,
		IsSystem:     true,
		Extend: &BoneExtend{
			NormalizedLocalAxisX:   mmath.NewMVec3(),
			NormalizedLocalAxisY:   mmath.NewMVec3(),
			NormalizedLocalAxisZ:   mmath.NewMVec3(),
			LocalAxis:              &mmath.MVec3{X: 1, Y: 0, Z: 0},
			IkLinkBoneIndexes:      make([]int, 0),
			IkTargetBoneIndexes:    make([]int, 0),
			ParentRelativePosition: mmath.NewMVec3(),
			ChildRelativePosition:  mmath.NewMVec3(),
			NormalizedFixedAxis:    mmath.NewMVec3(),
			TreeBoneIndexes:        make([]int, 0),
			RevertOffsetMatrix:     mmath.NewMMat4(),
			OffsetMatrix:           mmath.NewMMat4(),
			ParentBoneIndexes:      make([]int, 0),
			ParentBoneNames:        make([]string, 0),
			RelativeBoneIndexes:    make([]int, 0),
			ChildBoneIndexes:       make([]int, 0),
			EffectiveBoneIndexes:   make([]int, 0),
			AngleLimit:             false,
			MinAngleLimit:          mmath.NewMRotation(),
			MaxAngleLimit:          mmath.NewMRotation(),
			LocalAngleLimit:        false,
			LocalMinAngleLimit:     mmath.NewMRotation(),
			LocalMaxAngleLimit:     mmath.NewMRotation(),
			AxisSign:               1,
			RigidBody:              nil,
		},
	}
	bone.Extend.NormalizedLocalAxisX = bone.LocalAxisX.Copy()
	bone.Extend.NormalizedLocalAxisZ = bone.LocalAxisZ.Copy()
	bone.Extend.NormalizedLocalAxisY = bone.Extend.NormalizedLocalAxisZ.Cross(bone.Extend.NormalizedLocalAxisX)
	bone.Extend.NormalizedFixedAxis = bone.FixedAxis.Copy()
	return bone
}

func NewBoneByName(name string) *Bone {
	bone := NewBone()
	bone.SetName(name)
	return bone
}

func (bone *Bone) Index() int {
	return bone.index
}

func (bone *Bone) SetIndex(index int) {
	bone.index = index
}

func (bone *Bone) Name() string {
	return bone.name
}

func (bone *Bone) SetName(name string) {
	bone.name = name
}

func (bone *Bone) EnglishName() string {
	return bone.englishName
}

func (bone *Bone) SetEnglishName(englishName string) {
	bone.englishName = englishName
}

func (bone *Bone) Direction() string {
	if strings.Contains(bone.name, "左") {
		return "左"
	} else if strings.Contains(bone.name, "右") {
		return "右"
	}
	return ""
}

func (bone *Bone) IsValid() bool {
	return bone != nil && bone.index >= 0
}

func (bone *Bone) Copy() core.IIndexNameModel {
	var copiedIk *Ik
	if bone.Ik != nil {
		copiedIk = bone.Ik.Copy()
	}

	copied := &Bone{
		index:        bone.index,
		name:         bone.name,
		englishName:  bone.englishName,
		Position:     bone.Position.Copy(),
		ParentIndex:  bone.ParentIndex,
		Layer:        bone.Layer,
		BoneFlag:     bone.BoneFlag,
		TailPosition: bone.TailPosition.Copy(),
		TailIndex:    bone.TailIndex,
		EffectIndex:  bone.EffectIndex,
		EffectFactor: bone.EffectFactor,
		FixedAxis:    bone.FixedAxis.Copy(),
		LocalAxisX:   bone.LocalAxisX.Copy(),
		LocalAxisZ:   bone.LocalAxisZ.Copy(),
		EffectorKey:  bone.EffectorKey,
		Ik:           copiedIk,
		DisplaySlot:  bone.DisplaySlot,
		IsSystem:     bone.IsSystem,
		Extend:       bone.Extend.Copy(),
	}
	return copied
}

func (bone *Bone) NormalizeFixedAxis(fixedAxis *mmath.MVec3) {
	bone.Extend.NormalizedFixedAxis = fixedAxis.Normalized()
}

func (bone *Bone) NormalizeLocalAxis(localAxisX *mmath.MVec3) {
	bone.Extend.NormalizedLocalAxisX = localAxisX.Normalized()
	bone.Extend.NormalizedLocalAxisY = bone.Extend.NormalizedLocalAxisX.Cross(mmath.MVec3UnitZInv)
	bone.Extend.NormalizedLocalAxisZ = bone.Extend.NormalizedLocalAxisX.Cross(bone.Extend.NormalizedLocalAxisY)
}

// 表示先がボーンであるか
func (bone *Bone) IsTailBone() bool {
	return bone.BoneFlag&BONE_FLAG_TAIL_IS_BONE == BONE_FLAG_TAIL_IS_BONE
}

// 回転可能であるか
func (bone *Bone) CanRotate() bool {
	return bone.BoneFlag&BONE_FLAG_CAN_ROTATE == BONE_FLAG_CAN_ROTATE
}

// 移動可能であるか
func (bone *Bone) CanTranslate() bool {
	return bone.BoneFlag&BONE_FLAG_CAN_TRANSLATE == BONE_FLAG_CAN_TRANSLATE
}

// 表示であるか
func (bone *Bone) IsVisible() bool {
	return bone.BoneFlag&BONE_FLAG_IS_VISIBLE == BONE_FLAG_IS_VISIBLE
}

// 操作可であるか
func (bone *Bone) CanManipulate() bool {
	return bone.BoneFlag&BONE_FLAG_CAN_MANIPULATE == BONE_FLAG_CAN_MANIPULATE
}

// IKであるか
func (bone *Bone) IsIK() bool {
	return bone.BoneFlag&BONE_FLAG_IS_IK == BONE_FLAG_IS_IK
}

// ローカル付与であるか
func (bone *Bone) IsEffectorLocal() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_LOCAL == BONE_FLAG_IS_EXTERNAL_LOCAL
}

// 回転付与であるか
func (bone *Bone) IsEffectorRotation() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_ROTATION == BONE_FLAG_IS_EXTERNAL_ROTATION
}

// 移動付与であるか
func (bone *Bone) IsEffectorTranslation() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_TRANSLATION == BONE_FLAG_IS_EXTERNAL_TRANSLATION
}

// 軸固定であるか
func (bone *Bone) HasFixedAxis() bool {
	return bone.BoneFlag&BONE_FLAG_HAS_FIXED_AXIS == BONE_FLAG_HAS_FIXED_AXIS
}

// ローカル軸を持つか
func (bone *Bone) HasLocalAxis() bool {
	return bone.BoneFlag&BONE_FLAG_HAS_LOCAL_AXIS == BONE_FLAG_HAS_LOCAL_AXIS
}

// 物理後変形であるか
func (bone *Bone) IsAfterPhysicsDeform() bool {
	return bone.BoneFlag&BONE_FLAG_IS_AFTER_PHYSICS_DEFORM == BONE_FLAG_IS_AFTER_PHYSICS_DEFORM
}

// 外部親変形であるか
func (bone *Bone) IsEffectorParentDeform() bool {
	return bone.BoneFlag&BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM == BONE_FLAG_IS_EXTERNAL_PARENT_DEFORM
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

// 足IK系列であるか
func (bone *Bone) IsLegIK() bool {
	return bone.containsCategory(CATEGORY_LEG_IK)
}

// 足首から先であるか
func (bone *Bone) IsAnkle() bool {
	return bone.containsCategory(CATEGORY_ANKLE)
}

// 靴底であるか
func (bone *Bone) IsSole() bool {
	return bone.containsCategory(CATEGORY_SOLE)
}

// 捩りボーンであるか
func (bone *Bone) IsTwist() bool {
	return bone.containsCategory(CATEGORY_TWIST)
}

// 腕系ボーンであるか(指は含まない)
func (bone *Bone) IsArm() bool {
	return bone.containsCategory(CATEGORY_ARM)
}

// ひじ系ボーンであるか(指は含まない)
func (bone *Bone) IsElbow() bool {
	return bone.containsCategory(CATEGORY_ELBOW)
}

// 指系ボーンであるか
func (bone *Bone) IsFinger() bool {
	return bone.containsCategory(CATEGORY_FINGER)
}

// 指先系ボーンであるか
func (bone *Bone) IsFingerTail() bool {
	return bone.containsCategory(CATEGORY_FINGER_TAIL)
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

// 体幹系であるか
func (bone *Bone) IsTrunk() bool {
	return bone.containsCategory(CATEGORY_TRUNK)
}

func (bone *Bone) IsStandard() bool {
	if bone.Config() == nil {
		return false
	}

	return bone.Config().IsStandard
}

// フィッティングの時に移動だけ行うか
func (bone *Bone) CanFitOnlyMove() bool {
	return bone.containsCategory(CATEGORY_FITTING_ONLY_MOVE)
}

func (bone *Bone) Config() *BoneConfig {
	for boneConfigName, boneConfig := range GetStandardBoneConfigs() {
		if boneConfigName.String() == bone.Name() ||
			boneConfigName.Right() == bone.Name() ||
			boneConfigName.Left() == bone.Name() {
			return boneConfig
		}
	}
	return nil
}

// 定義上の親ボーン名
func (bone *Bone) ConfigParentBoneNames() []string {
	for boneConfigName, boneConfig := range GetStandardBoneConfigs() {
		if boneConfigName.String() == bone.Name() ||
			boneConfigName.Right() == bone.Name() ||
			boneConfigName.Left() == bone.Name() {

			boneNames := make([]string, 0)
			for _, parentBoneName := range boneConfig.ParentBoneNames {
				if boneConfigName.Right() == bone.Name() {
					boneNames = append(boneNames, parentBoneName.Right())
				} else if boneConfigName.Left() == bone.Name() {
					boneNames = append(boneNames, parentBoneName.Left())
				} else {
					boneNames = append(boneNames, parentBoneName.String())
				}
			}
			return boneNames
		}
	}
	return []string{}
}

// 定義上の子ボーン名
func (bone *Bone) ConfigChildBoneNames() []string {
	for boneConfigName, boneConfig := range GetStandardBoneConfigs() {
		if boneConfigName.String() == bone.Name() ||
			boneConfigName.Right() == bone.Name() ||
			boneConfigName.Left() == bone.Name() {

			boneNames := make([]string, 0)
			for _, tailBoneName := range boneConfig.ChildBoneNames {
				if boneConfigName.Right() == bone.Name() {
					boneNames = append(boneNames, tailBoneName.Right())
				} else if boneConfigName.Left() == bone.Name() {
					boneNames = append(boneNames, tailBoneName.Left())
				} else {
					boneNames = append(boneNames, tailBoneName.String())
				}
			}
			return boneNames
		}
	}
	return []string{}
}

// 定義上のUP Fromボーン名
func (bone *Bone) ConfigUpFromBoneNames() []string {
	for boneConfigName, boneConfig := range GetStandardBoneConfigs() {
		if boneConfigName.String() == bone.Name() ||
			boneConfigName.Right() == bone.Name() ||
			boneConfigName.Left() == bone.Name() {

			boneNames := make([]string, 0)
			for _, upFromBoneName := range boneConfig.UpFromBoneNames {
				if boneConfigName.Right() == bone.Name() {
					boneNames = append(boneNames, upFromBoneName.Right())
				} else if boneConfigName.Left() == bone.Name() {
					boneNames = append(boneNames, upFromBoneName.Left())
				} else {
					boneNames = append(boneNames, upFromBoneName.String())
				}
			}
			return boneNames
		}
	}
	return []string{}
}

// 定義上のUP Toボーン名
func (bone *Bone) ConfigUpToBoneNames() []string {
	for boneConfigName, boneConfig := range GetStandardBoneConfigs() {
		if boneConfigName.String() == bone.Name() ||
			boneConfigName.Right() == bone.Name() ||
			boneConfigName.Left() == bone.Name() {

			boneNames := make([]string, 0)
			for _, upToBoneName := range boneConfig.UpToBoneNames {
				if boneConfigName.Right() == bone.Name() {
					boneNames = append(boneNames, upToBoneName.Right())
				} else if boneConfigName.Left() == bone.Name() {
					boneNames = append(boneNames, upToBoneName.Left())
				} else {
					boneNames = append(boneNames, upToBoneName.String())
				}
			}
			return boneNames
		}
	}
	return []string{}
}

// 指定したカテゴリーに属するか
func (bone *Bone) containsCategory(category BoneCategory) bool {
	for boneConfigName, boneConfig := range GetStandardBoneConfigs() {
		for _, c := range boneConfig.Categories {
			if c == category && (boneConfigName.String() == bone.Name() ||
				boneConfigName.Right() == bone.Name() ||
				boneConfigName.Left() == bone.Name()) {
				return true
			}
		}
	}
	return false
}

func (bone *Bone) setup() {
	// 各ボーンのローカル軸
	bone.Extend.LocalAxis = bone.Extend.ChildRelativePosition.Normalized()

	if bone.HasFixedAxis() {
		bone.NormalizeFixedAxis(bone.FixedAxis)
		bone.NormalizeLocalAxis(bone.FixedAxis)
	} else {
		bone.NormalizeLocalAxis(bone.Extend.LocalAxis)
	}

	// オフセット行列は自身の位置を原点に戻す行列
	bone.Extend.OffsetMatrix = bone.Position.Inverted().ToMat4()

	// 逆オフセット行列は親ボーンからの相対位置分
	bone.Extend.RevertOffsetMatrix = bone.Extend.ParentRelativePosition.ToMat4()
}

func (bones *Bones) getParentRelativePosition(boneIndex int) *mmath.MVec3 {
	bone := bones.Get(boneIndex)
	if bone == nil {
		return mmath.NewMVec3()
	}

	if bone.ParentIndex >= 0 && bones.Contains(bone.ParentIndex) {
		return bone.Position.Subed(bones.Get(bone.ParentIndex).Position)
	}
	// 親が見つからない場合、自分の位置を原点からの相対位置として返す
	return bone.Position.Copy()
}

func (bones *Bones) getChildRelativePosition(boneIndex int) *mmath.MVec3 {
	bone := bones.Get(boneIndex)
	if bone == nil {
		return mmath.NewMVec3()
	}

	fromPosition := bone.Position
	var toPosition *mmath.MVec3

	configChildBoneNames := bone.ConfigChildBoneNames()
	if len(configChildBoneNames) > 0 {
		for _, childBoneName := range configChildBoneNames {
			childBone := bones.GetByName(childBoneName)
			if childBone != nil {
				toPosition = childBone.Position
				break
			}
		}
	}

	if toPosition == nil {
		if bone.IsTailBone() && bone.TailIndex >= 0 && slices.Contains(bones.Indexes(), bone.TailIndex) {
			toPosition = bones.Get(bone.TailIndex).Position
		} else if !bone.IsTailBone() && bone.TailPosition.Length() > 0 {
			toPosition = bone.TailPosition.Added(bone.Position)
		} else if bone.ParentIndex < 0 || !bones.Contains(bone.ParentIndex) {
			return mmath.NewMVec3()
		} else {
			fromPosition = bones.Get(bone.ParentIndex).Position
			toPosition = bone.Position
		}
	}

	v := toPosition.Subed(fromPosition)
	return v
}

func (bones *Bones) getLayerIndexes(isAfterPhysics bool) []int {
	layerIndexes := make(layerIndexes, 0, len(bones.NameIndexes))
	for _, bone := range bones.Data {
		if (isAfterPhysics && bone.IsAfterPhysicsDeform()) || (!isAfterPhysics && !bone.IsAfterPhysicsDeform()) {
			// 物理前後でフィルタリング
			layerIndexes = append(layerIndexes, layerIndex{layer: bone.Layer, index: bone.Index()})
		}
	}
	sort.Sort(layerIndexes)

	indexes := make([]int, len(layerIndexes))
	for i, layerIndex := range layerIndexes {
		indexes[i] = layerIndex.index
	}

	return indexes
}

// 関連ボーンリストの取得
func (bones *Bones) getRelativeBoneIndexes(boneIndex int, parentBoneIndexes, relativeBoneIndexes []int) ([]int, []int) {

	if boneIndex <= 0 || !bones.Contains(boneIndex) {
		return parentBoneIndexes, relativeBoneIndexes
	}

	bone := bones.Get(boneIndex)
	if bones.Contains(bone.ParentIndex) && !slices.Contains(relativeBoneIndexes, bone.ParentIndex) {
		// 親ボーンを辿る(子から親の順番)
		parentBoneIndexes = append(parentBoneIndexes, bone.ParentIndex)
		relativeBoneIndexes = append(relativeBoneIndexes, bone.ParentIndex)
		parentBoneIndexes, relativeBoneIndexes =
			bones.getRelativeBoneIndexes(bone.ParentIndex, parentBoneIndexes, relativeBoneIndexes)
	}
	if (bone.IsEffectorRotation() || bone.IsEffectorTranslation()) &&
		bones.Contains(bone.EffectIndex) && !slices.Contains(relativeBoneIndexes, bone.EffectIndex) {
		// 付与親ボーンを辿る
		relativeBoneIndexes = append(relativeBoneIndexes, bone.EffectIndex)
		_, relativeBoneIndexes =
			bones.getRelativeBoneIndexes(bone.EffectIndex, parentBoneIndexes, relativeBoneIndexes)
	}
	if bone.IsIK() {
		if bones.Contains(bone.Ik.BoneIndex) && !slices.Contains(relativeBoneIndexes, bone.Ik.BoneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, bone.Ik.BoneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(bone.Ik.BoneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
		for _, link := range bone.Ik.Links {
			if bones.Contains(link.BoneIndex) && !slices.Contains(relativeBoneIndexes, link.BoneIndex) {
				// IKリンクボーンを辿る
				relativeBoneIndexes = append(relativeBoneIndexes, link.BoneIndex)
				_, relativeBoneIndexes =
					bones.getRelativeBoneIndexes(link.BoneIndex, parentBoneIndexes, relativeBoneIndexes)
			}
		}
	}
	for _, boneIndex := range bone.Extend.EffectiveBoneIndexes {
		if bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// 外部子ボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.Extend.IkTargetBoneIndexes {
		if bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.Extend.IkLinkBoneIndexes {
		if bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKリンクボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}

	return parentBoneIndexes, relativeBoneIndexes
}

func (bones *Bones) Setup() {
	bones.LayerSortedIndexes = make([]int, 0)
	bones.LayerSortedBones = make(map[bool][]*Bone)
	bones.LayerSortedNames = make(map[bool]map[string]int)
	bones.DeformBoneIndexes = make(map[int][]int)

	for _, bone := range bones.Data {
		// 関係ボーンリストを一旦クリア
		if bone.Extend == nil {
			bone.Extend = &BoneExtend{}
		} else {
			bone.Extend.IkLinkBoneIndexes = make([]int, 0)
			bone.Extend.IkTargetBoneIndexes = make([]int, 0)
			bone.Extend.EffectiveBoneIndexes = make([]int, 0)
			bone.Extend.ChildBoneIndexes = make([]int, 0)
			bone.Extend.RelativeBoneIndexes = make([]int, 0)
			bone.Extend.ParentBoneIndexes = make([]int, 0)
			bone.Extend.ParentBoneNames = make([]string, 0)
			bone.Extend.TreeBoneIndexes = make([]int, 0)
		}
	}

	// 関連ボーンINDEX情報を設定
	for i := range len(bones.Data) {
		bone := bones.Get(i)
		if strings.HasPrefix(bone.Name(), "左") {
			bone.Extend.AxisSign = -1
		}
		if bone.IsIK() && bone.Ik != nil {
			// IKのリンクとターゲット
			for _, link := range bone.Ik.Links {
				if bones.Contains(link.BoneIndex) &&
					!slices.Contains(bones.Get(link.BoneIndex).Extend.IkLinkBoneIndexes, bone.Index()) {
					// リンクボーンにフラグを立てる
					linkBone := bones.Get(link.BoneIndex)
					linkBone.Extend.IkLinkBoneIndexes = append(linkBone.Extend.IkLinkBoneIndexes, bone.Index())
					// リンクの制限をコピーしておく
					linkBone.Extend.AngleLimit = link.AngleLimit
					linkBone.Extend.MinAngleLimit = link.MinAngleLimit
					linkBone.Extend.MaxAngleLimit = link.MaxAngleLimit
					linkBone.Extend.LocalAngleLimit = link.LocalAngleLimit
					linkBone.Extend.LocalMinAngleLimit = link.LocalMinAngleLimit
					linkBone.Extend.LocalMaxAngleLimit = link.LocalMaxAngleLimit
				}
			}
			if bones.Contains(bone.Ik.BoneIndex) &&
				!slices.Contains(bones.Get(bone.Ik.BoneIndex).Extend.IkTargetBoneIndexes, bone.Index()) {
				// ターゲットボーンにもフラグを立てる
				bones.Get(bone.Ik.BoneIndex).Extend.IkTargetBoneIndexes = append(bones.Get(bone.Ik.BoneIndex).Extend.IkTargetBoneIndexes, bone.Index())
			}
		}
		if bone.EffectIndex >= 0 && bones.Contains(bone.EffectIndex) &&
			!slices.Contains(bones.Get(bone.EffectIndex).Extend.EffectiveBoneIndexes, bone.Index()) {
			// 付与親の方に付与子情報を保持
			bones.Get(bone.EffectIndex).Extend.EffectiveBoneIndexes = append(bones.Get(bone.EffectIndex).Extend.EffectiveBoneIndexes, bone.Index())
		}
	}

	for i := range bones.Len() {
		bone := bones.Get(i)
		// 影響があるボーンINDEXリスト
		bone.Extend.ParentBoneIndexes, bone.Extend.RelativeBoneIndexes = bones.getRelativeBoneIndexes(bone.Index(), make([]int, 0), make([]int, 0))

		// ボーンINDEXリストからボーン名リストを作成
		bone.Extend.ParentBoneNames = make([]string, len(bone.Extend.ParentBoneIndexes))
		for i, parentBoneIndex := range bone.Extend.ParentBoneIndexes {
			bone.Extend.ParentBoneNames[i] = bones.Get(parentBoneIndex).Name()
		}

		// 親ボーンに子ボーンとして登録する
		if bone.ParentIndex >= 0 && bones.Contains(bone.ParentIndex) {
			bones.Get(bone.ParentIndex).Extend.ChildBoneIndexes = append(bones.Get(bone.ParentIndex).Extend.ChildBoneIndexes, bone.Index())
		}
		// 親からの相対位置
		bone.Extend.ParentRelativePosition = bones.getParentRelativePosition(bone.Index())
		// 子への相対位置
		bone.Extend.ChildRelativePosition = bones.getChildRelativePosition(bone.Index())
		// ボーン単体のセットアップ
		bone.setup()
	}

	// 変形階層・ボーンINDEXでソート

	// 変形前と変形後に分けてINDEXリストを生成
	bones.createLayerIndexes(false)
	bones.createLayerIndexes(true)

	for _, bone := range bones.Data {
		// ボーンのデフォームINDEXリストを取得
		bones.createLayerSortedBones(bone)
	}
}

func (bones *Bones) createLayerSortedBones(bone *Bone) {
	// 関連ボーンINDEXリスト（順不同）
	relativeBoneIndexes := make(map[int]struct{})

	// 対象のボーンは常に追加
	relativeBoneIndexes[bone.Index()] = struct{}{}

	// 関連するボーンの追加
	for _, index := range bone.Extend.RelativeBoneIndexes {
		if _, ok := relativeBoneIndexes[index]; !ok {
			relativeBoneIndexes[index] = struct{}{}
		}
	}

	deformBoneIndexes := make([]int, 0)
	for _, ap := range []bool{false, true} {
		for _, bone := range bones.LayerSortedBones[ap] {
			if _, ok := relativeBoneIndexes[bone.Index()]; ok {
				deformBoneIndexes = append(deformBoneIndexes, bone.Index())
			}
		}
	}

	bones.DeformBoneIndexes[bone.Index()] = deformBoneIndexes
}

func (bones *Bones) createLayerIndexes(isAfterPhysics bool) {
	bones.LayerSortedBones[isAfterPhysics] = make([]*Bone, 0)
	bones.LayerSortedNames[isAfterPhysics] = make(map[string]int)
	layerIndexes := bones.getLayerIndexes(isAfterPhysics)

	for i, boneIndex := range layerIndexes {
		bone := bones.Get(boneIndex)
		bones.LayerSortedNames[isAfterPhysics][bone.Name()] = i
		bones.LayerSortedBones[isAfterPhysics] = append(bones.LayerSortedBones[isAfterPhysics], bone)
		bones.LayerSortedIndexes = append(bones.LayerSortedIndexes, bone.Index())
		i++
	}
}

// 変形階層とINDEXのソート用構造体
type layerIndex struct {
	layer int
	index int
}

type layerIndexes []layerIndex

func (li layerIndexes) Len() int {
	return len(li)
}
func (li layerIndexes) Less(i, j int) bool {
	return li[i].layer < li[j].layer || (li[i].layer == li[j].layer && li[i].index < li[j].index)
}
func (li layerIndexes) Swap(i, j int) {
	li[i], li[j] = li[j], li[i]
}

func (li layerIndexes) Contains(index int) bool {
	for _, layerIndex := range li {
		if layerIndex.index == index {
			return true
		}
	}
	return false
}

// ボーンリスト
type Bones struct {
	*core.IndexNameModels[*Bone]
	LayerSortedBones   map[bool][]*Bone
	LayerSortedNames   map[bool]map[string]int
	DeformBoneIndexes  map[int][]int
	LayerSortedIndexes []int
}

func NewBones(count int) *Bones {
	return &Bones{
		IndexNameModels:    core.NewIndexNameModels[*Bone](count, func() *Bone { return nil }),
		LayerSortedBones:   make(map[bool][]*Bone),
		LayerSortedNames:   make(map[bool]map[string]int),
		LayerSortedIndexes: make([]int, 0),
	}
}

func (bones *Bones) Copy() *Bones {
	copied := NewBones(len(bones.Data))
	for i, bone := range bones.Data {
		copied.SetItem(i, bone.Copy().(*Bone))
	}
	return copied
}

func (bones *Bones) Insert(bone *Bone, afterIndex int) {
	// 挿入位置を探す
	insertPos := -1

	if afterIndex < 0 {
		// afterIndexが-1の場合、ルートに挿入
		bone.Layer = 0
		insertPos = len(bones.LayerSortedIndexes)
	} else {
		for i, boneIndex := range bones.LayerSortedIndexes {
			if boneIndex == afterIndex {
				insertPos = i + 1
				break
			}
		}
	}

	if insertPos < 0 {
		// 挿入場所が見つからない場合、最後に挿入
		lastBone := bones.Get(bones.LayerSortedIndexes[len(bones.LayerSortedIndexes)-1])
		bone.Layer = lastBone.Layer
		bones.Append(bone)
		bones.Setup()
		return
	}

	// 新しい要素のLayerを決定
	var newLayer int
	if insertPos == len(bones.LayerSortedIndexes) {
		// 挿入位置が最後の場合
		boneAtPrevPos := bones.Get(bones.LayerSortedIndexes[insertPos-1])
		if afterIndex >= 0 {
			newLayer = boneAtPrevPos.Layer
		} else {
			// ルートに挿入の場合、全てのLayerをインクリメント
			newLayer = 0
			for _, boneIndex := range bones.LayerSortedIndexes {
				boneToAdjust := bones.Get(boneIndex)
				boneToAdjust.Layer++
			}
		}
	} else {
		// 挿入位置が途中の場合
		boneAtPrevPos := bones.Get(bones.LayerSortedIndexes[insertPos-1])
		currentLayer := boneAtPrevPos.Layer
		boneAtNextPos := bones.Get(bones.LayerSortedIndexes[insertPos])
		nextLayer := boneAtNextPos.Layer

		if currentLayer == nextLayer {
			// 新しい要素のLayerをcurrentLayerに設定
			newLayer = currentLayer
			// 挿入位置以降の要素のLayerをインクリメント
			for i := insertPos; i < len(bones.LayerSortedIndexes); i++ {
				boneToAdjust := bones.Get(bones.LayerSortedIndexes[i])
				boneToAdjust.Layer++
			}
		} else if currentLayer+1 < nextLayer {
			// Layerの隙間がある場合
			newLayer = currentLayer + 1
		} else {
			// 新しい要素のLayerをcurrentLayerに設定
			newLayer = currentLayer
			// 挿入位置以降の要素のLayerをインクリメント
			for i := insertPos; i < len(bones.LayerSortedIndexes); i++ {
				boneToAdjust := bones.Get(bones.LayerSortedIndexes[i])
				boneToAdjust.Layer++
			}
		}
	}

	bone.Layer = newLayer
	bones.Append(bone)
	bones.Setup()
}

func (bones *Bones) GetIkTarget(ikBoneName string) *Bone {
	if ikBoneName == "" || !bones.ContainsByName(ikBoneName) {
		return nil
	}

	ikBone := bones.GetByName(ikBoneName)
	if ikBone == nil || !ikBone.IsIK() || !bones.Contains(ikBone.Ik.BoneIndex) {
		return nil
	}

	return bones.Get(ikBone.Ik.BoneIndex)
}
