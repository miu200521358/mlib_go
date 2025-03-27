package pmx

import (
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type IkLink struct {
	BoneIndex          int          // リンクボーンのボーンIndex
	AngleLimit         bool         // 角度制限有無
	MinAngleLimit      *mmath.MVec3 // 下限
	MaxAngleLimit      *mmath.MVec3 // 上限
	LocalAngleLimit    bool         // ローカル軸の角度制限有無
	LocalMinAngleLimit *mmath.MVec3 // ローカル軸制限の下限
	LocalMaxAngleLimit *mmath.MVec3 // ローカル軸制限の上限
}

func NewIkLink() *IkLink {
	return &IkLink{
		BoneIndex:          -1,
		AngleLimit:         false,
		MinAngleLimit:      mmath.NewMVec3(),
		MaxAngleLimit:      mmath.NewMVec3(),
		LocalAngleLimit:    false,
		LocalMinAngleLimit: mmath.NewMVec3(),
		LocalMaxAngleLimit: mmath.NewMVec3(),
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
	BoneIndex    int          // IKターゲットボーンのボーンIndex
	LoopCount    int          // IKループ回数 (最大255)
	UnitRotation *mmath.MVec3 // IKループ計算時の1回あたりの制限角度
	Links        []*IkLink    // IKリンクリスト
}

func NewIk() *Ik {
	return &Ik{
		BoneIndex:    -1,
		LoopCount:    0,
		UnitRotation: mmath.NewMVec3(),
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
	index                  int          // ボーンINDEX
	name                   string       // ボーン名
	englishName            string       // ボーン英名
	Position               *mmath.MVec3 // ボーン位置
	ParentIndex            int          // 親ボーンのボーンIndex
	Layer                  int          // 変形階層
	BoneFlag               BoneFlag     // ボーンフラグ(16bit) 各bit 0:OFF 1:ON
	TailPosition           *mmath.MVec3 // 接続先:0 の場合 座標オフセット, ボーン位置からの相対分
	TailIndex              int          // 接続先:1 の場合 接続先ボーンのボーンIndex
	EffectIndex            int          // 回転付与:1 または 移動付与:1 の場合 付与親ボーンのボーンIndex
	EffectFactor           float64      // 付与率
	FixedAxis              *mmath.MVec3 // 軸固定:1 の場合 軸の方向ベクトル
	LocalAxisX             *mmath.MVec3 // ローカル軸:1 の場合 X軸の方向ベクトル
	LocalAxisZ             *mmath.MVec3 // ローカル軸:1 の場合 Z軸の方向ベクトル
	EffectorKey            int          // 外部親変形:1 の場合 Key値
	Ik                     *Ik          // IK:1 の場合 IKデータを格納
	DisplaySlot            int          // 該当表示枠
	IsSystem               bool         // システム計算用追加ボーン の場合 true
	NormalizedLocalAxisX   *mmath.MVec3 // 計算済みのX軸の方向ベクトル
	NormalizedLocalAxisY   *mmath.MVec3 // 計算済みのY軸の方向ベクトル
	NormalizedLocalAxisZ   *mmath.MVec3 // 計算済みのZ軸の方向ベクトル
	NormalizedFixedAxis    *mmath.MVec3 // 計算済みの軸制限ベクトル
	LocalAxis              *mmath.MVec3 // ローカル軸の方向ベクトル(CorrectedLocalXVectorの正規化ベクトル)
	ParentRelativePosition *mmath.MVec3 // 親ボーンからの相対位置
	ChildRelativePosition  *mmath.MVec3 // Tailボーンへの相対位置
	RevertOffsetMatrix     *mmath.MMat4 // 逆オフセット行列(親ボーンからの相対位置分を戻す)
	OffsetMatrix           *mmath.MMat4 // オフセット行列 (自身の位置を原点に戻す行列)
	TreeBoneIndexes        []int        // 自分のボーンまでのボーンIndexのリスト
	ParentBoneIndexes      []int        // 自分の親ボーンからルートまでのボーンIndexのリスト
	ParentBoneNames        []string     // 自分の親ボーンからルートまでのボーン名のリスト
	RelativeBoneIndexes    []int        // 関連ボーンINDEX一覧（付与親とかIKとか）
	ChildBoneIndexes       []int        // 自分を親として登録しているボーンINDEX一覧
	EffectiveBoneIndexes   []int        // 自分を付与親として登録しているボーンINDEX一覧
	IkLinkBoneIndexes      []int        // 自分をIKリンクとして登録されているIKボーンのボーンIndex
	IkTargetBoneIndexes    []int        // 自分をIKターゲットとして登録されているIKボーンのボーンIndex
	AngleLimit             bool         // 自分がIKリンクボーンの角度制限がある場合、true
	MinAngleLimit          *mmath.MVec3 // 自分がIKリンクボーンの角度制限の下限
	MaxAngleLimit          *mmath.MVec3 // 自分がIKリンクボーンの角度制限の上限
	LocalAngleLimit        bool         // 自分がIKリンクボーンのローカル軸角度制限がある場合、true
	LocalMinAngleLimit     *mmath.MVec3 // 自分がIKリンクボーンのローカル軸角度制限の下限
	LocalMaxAngleLimit     *mmath.MVec3 // 自分がIKリンクボーンのローカル軸角度制限の上限
	AxisSign               int          // ボーンの軸ベクトル(左は-1, 右は1)
	RigidBody              *RigidBody   // 物理演算用剛体
	OriginalLayer          int          // 元の変形階層
}

func NewBone() *Bone {
	bone := &Bone{
		index:                  -1,
		name:                   "",
		englishName:            "",
		Position:               mmath.NewMVec3(),
		ParentIndex:            -1,
		Layer:                  -1,
		BoneFlag:               BONE_FLAG_NONE,
		TailPosition:           mmath.NewMVec3(),
		TailIndex:              -1,
		EffectIndex:            -1,
		EffectFactor:           0.0,
		FixedAxis:              nil,
		LocalAxisX:             nil,
		LocalAxisZ:             nil,
		EffectorKey:            -1,
		Ik:                     nil,
		DisplaySlot:            -1,
		IsSystem:               false,
		NormalizedLocalAxisX:   &mmath.MVec3{X: 1, Y: 0, Z: 0},
		NormalizedLocalAxisY:   &mmath.MVec3{X: 0, Y: 1, Z: 0},
		NormalizedLocalAxisZ:   &mmath.MVec3{X: 0, Y: 0, Z: -1},
		LocalAxis:              &mmath.MVec3{X: 1, Y: 0, Z: 0},
		IkLinkBoneIndexes:      make([]int, 0),
		IkTargetBoneIndexes:    make([]int, 0),
		ParentRelativePosition: mmath.NewMVec3(),
		ChildRelativePosition:  mmath.NewMVec3(),
		NormalizedFixedAxis:    nil,
		TreeBoneIndexes:        make([]int, 0),
		RevertOffsetMatrix:     mmath.NewMMat4(),
		OffsetMatrix:           mmath.NewMMat4(),
		ParentBoneIndexes:      make([]int, 0),
		ParentBoneNames:        make([]string, 0),
		RelativeBoneIndexes:    make([]int, 0),
		ChildBoneIndexes:       make([]int, 0),
		EffectiveBoneIndexes:   make([]int, 0),
		AngleLimit:             false,
		MinAngleLimit:          nil,
		MaxAngleLimit:          nil,
		LocalAngleLimit:        false,
		LocalMinAngleLimit:     nil,
		LocalMaxAngleLimit:     nil,
		AxisSign:               1,
		RigidBody:              nil,
		OriginalLayer:          -1,
	}
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

func (bone *Bone) Direction() BoneDirection {
	if strings.Contains(bone.name, "左") {
		return BONE_DIRECTION_LEFT
	} else if strings.Contains(bone.name, "右") {
		return BONE_DIRECTION_RIGHT
	}
	return BONE_DIRECTION_TRUNK
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
	}
	return copied
}

func (bone *Bone) NormalizeFixedAxis(fixedAxis *mmath.MVec3) {
	bone.NormalizedFixedAxis = fixedAxis.Normalized()
}

func (bone *Bone) NormalizeLocalAxis(localAxisX *mmath.MVec3) {
	bone.NormalizedLocalAxisX = localAxisX.Normalized()
	bone.NormalizedLocalAxisY = bone.NormalizedLocalAxisX.Cross(mmath.MVec3UnitZNeg)
	bone.NormalizedLocalAxisZ = bone.NormalizedLocalAxisX.Cross(bone.NormalizedLocalAxisY)
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

func (bone *Bone) Config() *BoneConfig {
	for boneConfigName, boneConfig := range getStandardBoneConfigs() {
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
	for boneConfigName, boneConfig := range getStandardBoneConfigs() {
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
	for boneConfigName, boneConfig := range getStandardBoneConfigs() {
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

func (bone *Bone) setup() {
	// 各ボーンのローカル軸
	bone.LocalAxis = bone.ChildRelativePosition.Normalized()

	if bone.HasFixedAxis() {
		bone.NormalizeFixedAxis(bone.FixedAxis)
		bone.NormalizeLocalAxis(bone.FixedAxis)
	} else {
		bone.NormalizeLocalAxis(bone.LocalAxis)
	}

	// オフセット行列は自身の位置を原点に戻す行列
	bone.OffsetMatrix = bone.Position.Negated().ToMat4()

	// 逆オフセット行列は親ボーンからの相対位置分
	bone.RevertOffsetMatrix = bone.ParentRelativePosition.ToMat4()
}
