package pmx

import (
	"slices"
	"sort"
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

// ------------------------------------------------------------

// ボーンリスト
type Bones struct {
	*core.IndexNameModels[*Bone]
	LayerSortedBones       map[bool][]*Bone
	LayerSortedNames       map[bool]map[string]int
	LayerSortedBoneIndexes map[bool][]int
	DeformBoneIndexes      map[int][]int
	LayerSortedIndexes     []int
}

func NewBones(capacity int) *Bones {
	return &Bones{
		IndexNameModels:        core.NewIndexNameModels[*Bone](capacity),
		LayerSortedBones:       make(map[bool][]*Bone),
		LayerSortedNames:       make(map[bool]map[string]int),
		LayerSortedBoneIndexes: make(map[bool][]int),
		LayerSortedIndexes:     make([]int, 0),
	}
}

func (bones *Bones) getInsertAfterIndex(bone *Bone) int {
	parentLayerIndex := slices.Index(bones.LayerSortedIndexes, bone.ParentIndex)
	ikBoneIndex := -1
	if bone.IsIK() && bone.Ik != nil {
		ikBoneIndex = slices.Index(bones.LayerSortedIndexes, bone.Ik.BoneIndex)
	}
	effectIndex := -1
	effectIkIndex := -1
	effectLayerIndex := -1
	effectIkLayerIndex := -1
	if bone.EffectIndex != -1 {
		effectBone, err := bones.Get(bone.EffectIndex)
		if err != nil {
			return -1
		}
		effectIndex = effectBone.Index()
		effectLayerIndex = slices.Index(bones.LayerSortedIndexes, effectIndex)
		if len(effectBone.IkLinkBoneIndexes) > 0 {
			if ikBone, err := bones.Get(effectBone.IkLinkBoneIndexes[0]); err == nil {
				effectIkIndex = ikBone.Index()
				effectIkLayerIndex = slices.Index(bones.LayerSortedIndexes, ikBoneIndex)
			}
		}
	}

	switch mmath.ArgMax([]int{parentLayerIndex, ikBoneIndex, effectLayerIndex, effectIkLayerIndex}) {
	case 0:
		return bone.ParentIndex
	case 1:
		return bone.Ik.BoneIndex
	case 2:
		return effectIndex
	case 3:
		return effectIkIndex
	}

	return -1
}

func (bones *Bones) Insert(bone *Bone) {
	afterIndex := bones.getInsertAfterIndex(bone)

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
		if lastBone, err := bones.Get(bones.LayerSortedIndexes[len(bones.LayerSortedIndexes)-1]); err == nil {
			bone.Layer = lastBone.Layer
			bones.Append(bone)
		}
		return
	}

	// 新しい要素のLayerを決定
	var newLayer int
	if insertPos == len(bones.LayerSortedIndexes) {
		// 挿入位置が最後の場合
		if boneAtPrevPos, err := bones.Get(bones.LayerSortedIndexes[insertPos-1]); err == nil {
			if afterIndex >= 0 {
				newLayer = boneAtPrevPos.Layer
			} else {
				// ルートに挿入の場合、全てのLayerをインクリメント
				newLayer = 0
				for _, boneIndex := range bones.LayerSortedIndexes {
					if boneToAdjust, err := bones.Get(boneIndex); err == nil {
						boneToAdjust.Layer++
					}
				}
			}
		}
	} else {
		// 挿入位置が途中の場合
		boneAtPrevPos, err := bones.Get(bones.LayerSortedIndexes[insertPos-1])
		if err != nil {
			return
		}
		currentLayer := boneAtPrevPos.Layer
		boneAtNextPos, err := bones.Get(bones.LayerSortedIndexes[insertPos])
		if err != nil {
			return
		}
		nextLayer := boneAtNextPos.Layer

		if currentLayer == nextLayer {
			// 新しい要素のLayerをcurrentLayerに設定
			newLayer = currentLayer
			// 挿入位置以降の要素のLayerをインクリメント
			for i := insertPos; i < len(bones.LayerSortedIndexes); i++ {
				if boneToAdjust, err := bones.Get(bones.LayerSortedIndexes[i]); err == nil {
					boneToAdjust.Layer++
				}
			}
		} else if currentLayer+1 < nextLayer {
			// Layerの隙間がある場合
			newLayer = currentLayer + 1
		} else {
			// 新しい要素のLayerをcurrentLayerに設定
			newLayer = currentLayer
			// 挿入位置以降の要素のLayerをインクリメント
			for i := insertPos; i < len(bones.LayerSortedIndexes); i++ {
				if boneToAdjust, err := bones.Get(bones.LayerSortedIndexes[i]); err == nil {
					boneToAdjust.Layer++
				}
			}
		}
	}

	bone.Layer = newLayer
	bones.Append(bone)
}

func (bones *Bones) GetIkTarget(ikBoneName string) *Bone {
	if ikBoneName == "" || !bones.ContainsByName(ikBoneName) {
		return nil
	}

	if ikBone, err := bones.GetByName(ikBoneName); err != nil ||
		!ikBone.IsIK() || !bones.Contains(ikBone.Ik.BoneIndex) {
		return nil
	} else {
		return ikBone
	}
}

func (bones *Bones) Setup() {
	bones.LayerSortedIndexes = make([]int, 0)
	bones.LayerSortedBones = make(map[bool][]*Bone)
	bones.LayerSortedNames = make(map[bool]map[string]int)
	bones.LayerSortedBoneIndexes = make(map[bool][]int)
	bones.DeformBoneIndexes = make(map[int][]int)

	for bone := range bones.Iterator() {
		// 関係ボーンリストを一旦クリア
		bone.IkLinkBoneIndexes = make([]int, 0)
		bone.IkTargetBoneIndexes = make([]int, 0)
		bone.EffectiveBoneIndexes = make([]int, 0)
		bone.ChildBoneIndexes = make([]int, 0)
		bone.RelativeBoneIndexes = make([]int, 0)
		bone.ParentBoneIndexes = make([]int, 0)
		bone.ParentBoneNames = make([]string, 0)
		bone.TreeBoneIndexes = make([]int, 0)
	}

	// 関連ボーンINDEX情報を設定
	for i := range bones.Length() {
		bone, err := bones.Get(i)
		if err != nil {
			continue
		}
		if strings.HasPrefix(bone.Name(), "左") {
			bone.AxisSign = -1
		}
		if bone.IsIK() && bone.Ik != nil {
			// IKのリンクとターゲット
			for _, link := range bone.Ik.Links {
				if linkBone, err := bones.Get(link.BoneIndex); err == nil && bones.Contains(link.BoneIndex) &&
					!slices.Contains(linkBone.IkLinkBoneIndexes, bone.Index()) {
					// リンクボーンにフラグを立てる
					linkBone := linkBone
					linkBone.IkLinkBoneIndexes = append(linkBone.IkLinkBoneIndexes, bone.Index())
					// リンクの制限をコピーしておく
					linkBone.AngleLimit = link.AngleLimit
					linkBone.MinAngleLimit = link.MinAngleLimit
					linkBone.MaxAngleLimit = link.MaxAngleLimit
					linkBone.LocalAngleLimit = link.LocalAngleLimit
					linkBone.LocalMinAngleLimit = link.LocalMinAngleLimit
					linkBone.LocalMaxAngleLimit = link.LocalMaxAngleLimit
				}
			}
			if ikBone, err := bones.Get(bone.Ik.BoneIndex); err == nil && !slices.Contains(ikBone.IkTargetBoneIndexes, bone.Index()) {
				// ターゲットボーンにもフラグを立てる
				ikBone.IkTargetBoneIndexes = append(ikBone.IkTargetBoneIndexes, bone.Index())
			}
		}
		if effectBone, err := bones.Get(bone.EffectIndex); err == nil && bone.EffectIndex >= 0 &&
			bones.Contains(bone.EffectIndex) && !slices.Contains(effectBone.EffectiveBoneIndexes, bone.Index()) {
			// 付与親の方に付与子情報を保持
			effectBone.EffectiveBoneIndexes = append(effectBone.EffectiveBoneIndexes, bone.Index())
		}
	}

	for i := range bones.Length() {
		bone, err := bones.Get(i)
		if err != nil {
			continue
		}
		// 影響があるボーンINDEXリスト
		bone.ParentBoneIndexes, bone.RelativeBoneIndexes = bones.getRelativeBoneIndexes(bone.Index(), make([]int, 0), make([]int, 0))

		// ボーンINDEXリストからボーン名リストを作成
		bone.ParentBoneNames = make([]string, len(bone.ParentBoneIndexes))
		for i, parentBoneIndex := range bone.ParentBoneIndexes {
			if parentBone, err := bones.Get(parentBoneIndex); err == nil {
				bone.ParentBoneNames[i] = parentBone.Name()
			}
		}

		// 親ボーンに子ボーンとして登録する
		if parentBone, err := bones.Get(bone.ParentIndex); err == nil {
			parentBone.ChildBoneIndexes = append(parentBone.ChildBoneIndexes, bone.Index())
		}
		// 親からの相対位置
		bone.ParentRelativePosition = bones.getParentRelativePosition(bone.Index())
		// 子への相対位置
		bone.ChildRelativePosition = bones.getChildRelativePosition(bone.Index())
		// ボーン単体のセットアップ
		bone.setup()
	}

	// 変形階層・ボーンINDEXでソート

	// 変形前と変形後に分けてINDEXリストを生成
	bones.createLayerIndexes()

	for bone := range bones.Iterator() {
		// ボーンのデフォームINDEXリストを取得
		bones.createLayerSortedBones(bone)
	}
}

func (bones *Bones) createLayerSortedBones(bone *Bone) {
	deformBoneIndexes := make([]int, 0)
	for _, boneIndex := range bones.LayerSortedIndexes {
		if slices.Contains(bone.RelativeBoneIndexes, boneIndex) || boneIndex == bone.Index() {
			deformBoneIndexes = append(deformBoneIndexes, boneIndex)
		}
	}

	bones.DeformBoneIndexes[bone.Index()] = deformBoneIndexes
}

func (bones *Bones) createLayerIndexes() {
	bones.LayerSortedBones[false] = make([]*Bone, 0)
	bones.LayerSortedNames[false] = make(map[string]int)
	bones.LayerSortedBoneIndexes[false] = make([]int, 0)

	bones.LayerSortedBones[true] = make([]*Bone, 0)
	bones.LayerSortedNames[true] = make(map[string]int)
	bones.LayerSortedBoneIndexes[true] = make([]int, 0)

	layerIndexes := make(layerIndexes, 0, bones.Length())
	for bone := range bones.Iterator() {
		layerIndexes = append(layerIndexes, layerIndex{isAfterPhysics: bone.IsAfterPhysicsDeform(), layer: bone.Layer, index: bone.Index()})
	}
	sort.Sort(layerIndexes)

	for i, layerBone := range layerIndexes {
		bone, err := bones.Get(layerBone.index)
		if err != nil {
			continue
		}
		bones.LayerSortedNames[layerBone.isAfterPhysics][bone.Name()] = i
		bones.LayerSortedBones[layerBone.isAfterPhysics] =
			append(bones.LayerSortedBones[layerBone.isAfterPhysics], bone)
		bones.LayerSortedBoneIndexes[layerBone.isAfterPhysics] =
			append(bones.LayerSortedBoneIndexes[layerBone.isAfterPhysics], layerBone.index)
		bones.LayerSortedIndexes = append(bones.LayerSortedIndexes, bone.Index())
		i++
	}
}

// 指定されたボーンのうち、もっとも変形階層が小さいINDEXを取得
func (bones *Bones) MinBoneIndex(boneIndexes []int) int {
	layerIndexes := make(layerIndexes, len(boneIndexes))
	for i, boneIndex := range boneIndexes {
		if bone, err := bones.Get(boneIndex); err == nil {
			layerIndexes[i] = layerIndex{isAfterPhysics: bone.IsAfterPhysicsDeform(), layer: bone.Layer, index: boneIndex}
		}
	}
	sort.Sort(layerIndexes)

	return layerIndexes[0].index
}

// 指定されたボーンのうち、もっとも変形階層が大きいINDEXを取得
func (bones *Bones) MaxBoneIndex(boneIndexes []int) int {
	layerIndexes := make(layerIndexes, len(boneIndexes))
	for i, boneIndex := range boneIndexes {
		if bone, err := bones.Get(boneIndex); err == nil {
			layerIndexes[i] = layerIndex{isAfterPhysics: bone.IsAfterPhysicsDeform(), layer: bone.Layer, index: boneIndex}
		}
	}
	sort.Sort(layerIndexes)

	return layerIndexes[len(boneIndexes)-1].index
}

func (bones *Bones) getParentRelativePosition(boneIndex int) *mmath.MVec3 {
	bone, err := bones.Get(boneIndex)
	if err != nil {
		return mmath.NewMVec3()
	}

	if bone.ParentIndex >= 0 && bones.Contains(bone.ParentIndex) {
		if parentBone, err := bones.Get(bone.ParentIndex); err == nil {
			return bone.Position.Subed(parentBone.Position)
		}
	}
	// 親が見つからない場合、自分の位置を原点からの相対位置として返す
	return bone.Position.Copy()
}

func (bones *Bones) getChildRelativePosition(boneIndex int) *mmath.MVec3 {
	bone, err := bones.Get(boneIndex)
	if err != nil {
		return mmath.NewMVec3()
	}

	fromPosition := bone.Position
	var toPosition *mmath.MVec3

	configChildBoneNames := bone.ConfigChildBoneNames()
	if len(configChildBoneNames) > 0 {
		for _, childBoneName := range configChildBoneNames {
			if childBone, err := bones.GetByName(childBoneName); err == nil {
				toPosition = childBone.Position
				break
			}
		}
	}

	if toPosition == nil {
		if bone.IsTailBone() && bone.TailIndex >= 0 && slices.Contains(bones.Indexes(), bone.TailIndex) {
			if toBone, err := bones.Get(bone.TailIndex); err == nil {
				toPosition = toBone.Position
			}
		} else if !bone.IsTailBone() && bone.TailPosition.Length() > 0 {
			toPosition = bone.TailPosition.Added(bone.Position)
		} else if bone.ParentIndex < 0 || !bones.Contains(bone.ParentIndex) {
			return mmath.NewMVec3()
		} else {
			if parentBone, err := bones.Get(bone.ParentIndex); err == nil {
				fromPosition = parentBone.Position
			}
			toPosition = bone.Position
		}
	}

	v := toPosition.Subed(fromPosition)
	return v
}

// 関連ボーンリストの取得
func (bones *Bones) getRelativeBoneIndexes(boneIndex int, parentBoneIndexes, relativeBoneIndexes []int) ([]int, []int) {

	if boneIndex <= 0 || !bones.Contains(boneIndex) {
		return parentBoneIndexes, relativeBoneIndexes
	}

	bone, err := bones.Get(boneIndex)
	if err != nil {
		return parentBoneIndexes, relativeBoneIndexes
	}
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
	for _, boneIndex := range bone.EffectiveBoneIndexes {
		if bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// 外部子ボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.IkTargetBoneIndexes {
		if bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.IkLinkBoneIndexes {
		if bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKリンクボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				bones.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}

	return parentBoneIndexes, relativeBoneIndexes
}

// ------------------------------------------------------------

// 変形階層とINDEXのソート用構造体
type layerIndex struct {
	isAfterPhysics bool
	layer          int
	index          int
}

type layerIndexes []layerIndex

func (li layerIndexes) Len() int {
	return len(li)
}
func (li layerIndexes) Less(i, j int) bool {
	ia := 0
	if li[i].isAfterPhysics {
		ia = 1
	}
	ib := 0
	if li[j].isAfterPhysics {
		ib = 1
	}

	return ia < ib || (ia == ib && li[i].layer < li[j].layer) ||
		(ia == ib && li[i].layer == li[j].layer && li[i].index < li[j].index)
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
