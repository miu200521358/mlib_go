package pmx

import (
	"slices"
	"sort"
	"strings"

	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
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
		MinAngleLimit:      mmath.NewRotation(),
		MaxAngleLimit:      mmath.NewRotation(),
		LocalAngleLimit:    false,
		LocalMinAngleLimit: mmath.NewRotation(),
		LocalMaxAngleLimit: mmath.NewRotation(),
	}
}

func (t *IkLink) Copy() *IkLink {
	copied := *t
	return &copied
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
		UnitRotation: mmath.NewRotation(),
		Links:        []*IkLink{},
	}
}

func (t *Ik) Copy() *Ik {
	copied := &Ik{}
	copied.BoneIndex = t.BoneIndex
	copied.LoopCount = t.LoopCount
	copied.UnitRotation = t.UnitRotation.Copy()
	copied.Links = make([]*IkLink, len(t.Links))
	for i, link := range t.Links {
		copied.Links[i] = link.Copy()
	}
	return copied
}

type Bone struct {
	*mcore.IndexNameModel
	Position               *mmath.MVec3     // ボーン位置
	ParentIndex            int              // 親ボーンのボーンIndex
	Layer                  int              // 変形階層
	BoneFlag               BoneFlag         // ボーンフラグ(16bit) 各bit 0:OFF 1:ON
	TailPosition           *mmath.MVec3     // 接続先:0 の場合 座標オフセット, ボーン位置からの相対分
	TailIndex              int              // 接続先:1 の場合 接続先ボーンのボーンIndex
	EffectIndex            int              // 回転付与:1 または 移動付与:1 の場合 付与親ボーンのボーンIndex
	EffectFactor           float64          // 付与率
	FixedAxis              *mmath.MVec3     // 軸固定:1 の場合 軸の方向ベクトル
	LocalAxisX             *mmath.MVec3     // ローカル軸:1 の場合 X軸の方向ベクトル
	LocalAxisZ             *mmath.MVec3     // ローカル軸:1 の場合 Z軸の方向ベクトル
	EffectorKey            int              // 外部親変形:1 の場合 Key値
	Ik                     *Ik              // IK:1 の場合 IKデータを格納
	DisplaySlot            int              // 該当表示枠
	IsSystem               bool             // システム計算用追加ボーン の場合 true
	NormalizedLocalAxisX   *mmath.MVec3     // 計算済みのX軸の方向ベクトル
	NormalizedLocalAxisY   *mmath.MVec3     // 計算済みのY軸の方向ベクトル
	NormalizedLocalAxisZ   *mmath.MVec3     // 計算済みのZ軸の方向ベクトル
	NormalizedFixedAxis    *mmath.MVec3     // 計算済みの軸制限ベクトル
	LocalAxis              *mmath.MVec3     // ローカル軸の方向ベクトル(CorrectedLocalXVectorの正規化ベクトル)
	ParentRelativePosition *mmath.MVec3     // 親ボーンからの相対位置
	ChildRelativePosition  *mmath.MVec3     // Tailボーンへの相対位置
	LocalMatrix            *mmath.MMat4     // ローカル軸行列
	RevertOffsetMatrix     *mmath.MMat4     // 逆オフセット行列(親ボーンからの相対位置分を戻す)
	OffsetMatrix           *mmath.MMat4     // オフセット行列 (自身の位置を原点に戻す行列)
	TreeBoneIndexes        []int            // 自分のボーンまでのボーンIndexのリスト
	ParentBoneIndexes      []int            // 自分の親ボーンからルートまでのボーンIndexのリスト
	RelativeBoneIndexes    []int            // 関連ボーンINDEX一覧（付与親とかIKとか）
	ChildBoneIndexes       []int            // 自分を親として登録しているボーンINDEX一覧
	EffectiveBoneIndexes   []int            // 自分を付与親として登録しているボーンINDEX一覧
	IkLinkBoneIndexes      []int            // IKリンクとして登録されているIKボーンのボーンIndex
	IkTargetBoneIndexes    []int            // IKターゲットとして登録されているIKボーンのボーンIndex
	AngleLimit             bool             // 自分がIKリンクボーンの角度制限がある場合、true
	MinAngleLimit          *mmath.MRotation // 自分がIKリンクボーンの角度制限の下限
	MaxAngleLimit          *mmath.MRotation // 自分がIKリンクボーンの角度制限の上限
	LocalAngleLimit        bool             // 自分がIKリンクボーンのローカル軸角度制限がある場合、true
	LocalMinAngleLimit     *mmath.MRotation // 自分がIKリンクボーンのローカル軸角度制限の下限
	LocalMaxAngleLimit     *mmath.MRotation // 自分がIKリンクボーンのローカル軸角度制限の上限
	AxisSign               int              // ボーンの軸ベクトル(左は-1, 右は1)
}

func NewBone() *Bone {
	bone := &Bone{
		IndexNameModel:         &mcore.IndexNameModel{Index: -1, Name: "", EnglishName: ""},
		Position:               mmath.NewMVec3(),
		ParentIndex:            -1,
		Layer:                  -1,
		BoneFlag:               BONE_FLAG_NONE,
		TailPosition:           mmath.NewMVec3(),
		TailIndex:              -1,
		EffectIndex:            -1,
		EffectFactor:           0.0,
		FixedAxis:              mmath.NewMVec3(),
		LocalAxisX:             mmath.NewMVec3(),
		LocalAxisZ:             mmath.NewMVec3(),
		EffectorKey:            -1,
		Ik:                     NewIk(),
		DisplaySlot:            -1,
		IsSystem:               true,
		NormalizedLocalAxisX:   mmath.NewMVec3(),
		NormalizedLocalAxisY:   mmath.NewMVec3(),
		NormalizedLocalAxisZ:   mmath.NewMVec3(),
		LocalAxis:              &mmath.MVec3{1, 0, 0},
		IkLinkBoneIndexes:      []int{},
		IkTargetBoneIndexes:    []int{},
		ParentRelativePosition: mmath.NewMVec3(),
		ChildRelativePosition:  mmath.NewMVec3(),
		NormalizedFixedAxis:    mmath.NewMVec3(),
		TreeBoneIndexes:        []int{},
		RevertOffsetMatrix:     mmath.NewMMat4(),
		OffsetMatrix:           mmath.NewMMat4(),
		ParentBoneIndexes:      []int{},
		RelativeBoneIndexes:    []int{},
		ChildBoneIndexes:       []int{},
		EffectiveBoneIndexes:   []int{},
		AngleLimit:             false,
		MinAngleLimit:          mmath.NewRotationByRadians(mmath.NewMVec3()),
		MaxAngleLimit:          mmath.NewRotationByRadians(mmath.NewMVec3()),
		LocalAngleLimit:        false,
		LocalMinAngleLimit:     mmath.NewRotationByRadians(mmath.NewMVec3()),
		LocalMaxAngleLimit:     mmath.NewRotationByRadians(mmath.NewMVec3()),
		AxisSign:               1,
	}
	bone.NormalizedLocalAxisX = bone.LocalAxisX.Copy()
	bone.NormalizedLocalAxisZ = bone.LocalAxisZ.Copy()
	bone.NormalizedLocalAxisY = bone.NormalizedLocalAxisZ.Cross(bone.NormalizedLocalAxisX)
	bone.NormalizedFixedAxis = bone.FixedAxis.Copy()
	return bone
}

func NewBoneByName(name string) *Bone {
	bone := NewBone()
	bone.Name = name
	return bone
}

func (v *Bone) Copy() mcore.IIndexNameModel {
	copied := NewBone()
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return copied
}

func (bone *Bone) NormalizeFixedAxis(fixedAxis *mmath.MVec3) {
	bone.NormalizedFixedAxis = fixedAxis.Normalize()
}

func (bone *Bone) NormalizeLocalAxis(localAxisX *mmath.MVec3) {
	bone.NormalizedLocalAxisX = localAxisX.Normalize()
	bone.NormalizedLocalAxisY = bone.NormalizedLocalAxisX.Cross(mmath.MVec3UnitZ.Invert())
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
	for _, boneConfig := range GetStandardBoneConfigs() {
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

func (bone *Bone) normalizeFixedAxis(fixedAxis *mmath.MVec3) {
	v := fixedAxis.Normalized()
	bone.NormalizedFixedAxis = v
}

func (bone *Bone) normalizeLocalAxis(localXVector *mmath.MVec3) {
	v := localXVector.Normalized()
	bone.NormalizedLocalAxisX = v
	bone.NormalizedLocalAxisY = v.Cross(&mmath.MVec3{0, 0, -1})
	bone.NormalizedLocalAxisZ = v.Cross(bone.NormalizedLocalAxisY)
}

func (bone *Bone) setup() {
	// 各ボーンのローカル軸
	localAxis := bone.ChildRelativePosition.Normalized()
	bone.LocalAxis = localAxis
	// ローカル軸行列
	bone.LocalMatrix = localAxis.ToLocalMatrix4x4()

	if bone.HasFixedAxis() {
		bone.normalizeFixedAxis(bone.FixedAxis)
		bone.normalizeLocalAxis(bone.FixedAxis)
	} else {
		bone.normalizeLocalAxis(bone.LocalAxis)
	}

	// オフセット行列は自身の位置を原点に戻す行列
	bone.OffsetMatrix.Translate(bone.Position.Inverted())

	// 逆オフセット行列は親ボーンからの相対位置分
	bone.RevertOffsetMatrix.Translate(bone.ParentRelativePosition.Copy())
}

func (b *Bones) getParentRelativePosition(boneIndex int) *mmath.MVec3 {
	bone := b.GetItem(boneIndex)
	if bone.ParentIndex >= 0 && b.Contains(bone.ParentIndex) {
		v := bone.Position.Subed(b.GetItem(bone.ParentIndex).Position)
		return v
	}
	// 親が見つからない場合、自分の位置を原点からの相対位置として返す
	return bone.Position
}

func (b *Bones) getChildRelativePosition(boneIndex int) *mmath.MVec3 {
	bone := b.GetItem(boneIndex)

	fromPosition := *bone.Position
	var toPosition *mmath.MVec3

	if bone.IsTailBone() && bone.TailIndex >= 0 && slices.Contains(b.GetIndexes(), bone.TailIndex) {
		toPosition = b.GetItem(bone.TailIndex).Position
	} else if !bone.IsTailBone() && bone.TailPosition.Length() > 0 {
		toPosition = bone.TailPosition.Add(bone.Position)
	} else if bone.ParentIndex < 0 || !b.Contains(bone.ParentIndex) {
		return mmath.NewMVec3()
	} else {
		fromPosition = *b.GetItem(bone.ParentIndex).Position
		toPosition = bone.Position
	}

	v := toPosition.Subed(&fromPosition)
	return v
}

func (b *Bones) GetInitializeLocalPosition(boneIndex int) *mmath.MVec3 {
	if boneIndex < 0 || !b.Contains(boneIndex) {
		return mmath.NewMVec3()
	}

	bone := b.GetItem(boneIndex)
	matrix := mmath.NewMMat4()
	for _, parentIndex := range bone.ParentBoneIndexes {
		if parentIndex >= 0 && b.Contains(parentIndex) {
			// 親ボーンの変形行列を掛ける(親->子の順で掛ける)
			matrix.Mul(b.GetItem(parentIndex).RevertOffsetMatrix)
		}
	}
	// 最後に対象ボーン自身の行列をかける
	matrix.Mul(bone.RevertOffsetMatrix)
	// BOf行列: 自身のボーンのボーンオフセット行列
	matrix.Mul(bone.OffsetMatrix)
	return matrix.Translation()
}

func (b *Bones) GetLayerIndexes() []int {
	layerIndexes := make(LayerIndexes, 0)
	for _, boneIndex := range b.GetIndexes() {
		bone := b.GetItem(boneIndex)
		layerIndexes = append(layerIndexes, LayerIndex{Layer: bone.Layer, Index: boneIndex})
	}
	sort.Sort(layerIndexes)

	indexes := make([]int, len(layerIndexes))
	for i, layerIndex := range layerIndexes {
		indexes[i] = layerIndex.Index
	}

	return indexes
}

// 関連ボーンリストの取得
func (b *Bones) getRelativeBoneIndexes(boneIndex int, parentBoneIndexes, relativeBoneIndexes []int) ([]int, []int) {

	if boneIndex <= 0 || !b.Contains(boneIndex) {
		return parentBoneIndexes, relativeBoneIndexes
	}

	bone := b.GetItem(boneIndex)
	if b.Contains(bone.ParentIndex) && !slices.Contains(relativeBoneIndexes, bone.ParentIndex) {
		// 親ボーンを辿る(親から子の順番)
		parentBoneIndexes = append([]int{bone.ParentIndex}, parentBoneIndexes...)
		relativeBoneIndexes = append(relativeBoneIndexes, bone.ParentIndex)
		parentBoneIndexes, relativeBoneIndexes =
			b.getRelativeBoneIndexes(bone.ParentIndex, parentBoneIndexes, relativeBoneIndexes)
	}
	if (bone.IsEffectorRotation() || bone.IsEffectorTranslation()) &&
		b.Contains(bone.EffectIndex) && !slices.Contains(relativeBoneIndexes, bone.EffectIndex) {
		// 付与親ボーンを辿る
		relativeBoneIndexes = append(relativeBoneIndexes, bone.EffectIndex)
		_, relativeBoneIndexes =
			b.getRelativeBoneIndexes(bone.EffectIndex, parentBoneIndexes, relativeBoneIndexes)
	}
	if bone.IsIK() {
		if b.Contains(bone.Ik.BoneIndex) && !slices.Contains(relativeBoneIndexes, bone.Ik.BoneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, bone.Ik.BoneIndex)
			_, relativeBoneIndexes =
				b.getRelativeBoneIndexes(bone.Ik.BoneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
		for _, link := range bone.Ik.Links {
			if b.Contains(link.BoneIndex) && !slices.Contains(relativeBoneIndexes, link.BoneIndex) {
				// IKリンクボーンを辿る
				relativeBoneIndexes = append(relativeBoneIndexes, link.BoneIndex)
				_, relativeBoneIndexes =
					b.getRelativeBoneIndexes(link.BoneIndex, parentBoneIndexes, relativeBoneIndexes)
			}
		}
	}
	for _, boneIndex := range bone.EffectiveBoneIndexes {
		if b.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// 外部子ボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				b.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.IkTargetBoneIndexes {
		if b.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				b.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.IkLinkBoneIndexes {
		if b.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKリンクボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				b.getRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}

	return parentBoneIndexes, relativeBoneIndexes
}

// IKツリーの親INDEXを取得
func (b *Bones) getIkTreeIndex(bone *Bone) *Bone {
	if bone == nil || bone.ParentIndex < 0 || !b.Contains(bone.ParentIndex) {
		return nil
	}

	parentBone := b.GetItem(bone.ParentIndex)
	if parentBone.Index < 0 {
		return nil
	}

	if _, ok := b.IkTreeIndexes[parentBone.Index]; ok {
		return parentBone
	} else {
		parentLayerBone := b.getIkTreeIndex(parentBone)
		if parentLayerBone != nil {
			return parentLayerBone
		}
	}

	if bone.IsEffectorRotation() || bone.IsEffectorTranslation() {
		effectBone := b.GetItem(bone.EffectIndex)
		if _, ok := b.IkTreeIndexes[effectBone.Index]; ok {
			return effectBone
		} else {
			effectorLayerBone := b.getIkTreeIndex(effectBone)
			if effectorLayerBone != nil {
				return effectorLayerBone
			}
		}
	}

	return nil
}

func (b *Bones) setup() {
	for _, bone := range b.Data {
		// 関係ボーンリストを一旦クリア
		bone.IkLinkBoneIndexes = make([]int, 0)
		bone.IkTargetBoneIndexes = make([]int, 0)
		bone.EffectiveBoneIndexes = make([]int, 0)
		bone.ChildBoneIndexes = make([]int, 0)
	}

	// 関連ボーンINDEX情報を設定
	for _, bone := range b.GetSortedData() {
		if strings.HasPrefix(bone.Name, "左") {
			bone.AxisSign = -1
		}
		if bone.IsIK() && bone.Ik != nil {
			// IKのリンクとターゲット
			for _, link := range bone.Ik.Links {
				if b.Contains(link.BoneIndex) &&
					!slices.Contains(b.GetItem(link.BoneIndex).IkLinkBoneIndexes, bone.Index) {
					// リンクボーンにフラグを立てる
					linkBone := b.GetItem(link.BoneIndex)
					linkBone.IkLinkBoneIndexes = append(linkBone.IkLinkBoneIndexes, bone.Index)
					// リンクの制限をコピーしておく
					linkBone.AngleLimit = link.AngleLimit
					linkBone.MinAngleLimit = link.MinAngleLimit
					linkBone.MaxAngleLimit = link.MaxAngleLimit
					linkBone.LocalAngleLimit = link.LocalAngleLimit
					linkBone.LocalMinAngleLimit = link.LocalMinAngleLimit
					linkBone.LocalMaxAngleLimit = link.LocalMaxAngleLimit
				}
			}
			if b.Contains(bone.Ik.BoneIndex) &&
				!slices.Contains(b.GetItem(bone.Ik.BoneIndex).IkTargetBoneIndexes, bone.Index) {
				// ターゲットボーンにもフラグを立てる
				b.GetItem(bone.Ik.BoneIndex).IkTargetBoneIndexes = append(b.GetItem(bone.Ik.BoneIndex).IkTargetBoneIndexes, bone.Index)
			}
		}
		if bone.EffectIndex >= 0 && b.Contains(bone.EffectIndex) &&
			!slices.Contains(b.GetItem(bone.EffectIndex).EffectiveBoneIndexes, bone.Index) {
			// 付与親の方に付与子情報を保持
			b.GetItem(bone.EffectIndex).EffectiveBoneIndexes = append(b.GetItem(bone.EffectIndex).EffectiveBoneIndexes, bone.Index)
		}
	}

	for _, bone := range b.GetSortedData() {
		// 影響があるボーンINDEXリスト
		bone.ParentBoneIndexes, bone.RelativeBoneIndexes = b.getRelativeBoneIndexes(bone.Index, []int{}, []int{})

		// 親ボーンに子ボーンとして登録する
		if bone.ParentIndex >= 0 && b.Contains(bone.ParentIndex) {
			b.GetItem(bone.ParentIndex).ChildBoneIndexes = append(b.GetItem(bone.ParentIndex).ChildBoneIndexes, bone.Index)
		}
		// 親からの相対位置
		bone.ParentRelativePosition = b.getParentRelativePosition(bone.Index)
		// 子への相対位置
		bone.ChildRelativePosition = b.getChildRelativePosition(bone.Index)
		// ボーン単体のセットアップ
		bone.setup()
	}

	// 変形階層・ボーンINDEXでソート
	b.LayerSortedIndexes = make(map[int]string, len(b.Data))
	b.LayerSortedNames = make(map[string]int, len(b.Data))
	layerIndexes := b.GetLayerIndexes()

	i := 0
	for _, boneIndex := range layerIndexes {
		bone := b.GetItem(boneIndex)
		b.LayerSortedNames[bone.Name] = i
		b.LayerSortedIndexes[i] = bone.Name
		i++
	}

	// 変形階層順に親子を繋げていく
ikLoop:
	for _, boneIndex := range layerIndexes {
		bone := b.GetItem(boneIndex)
		if bone.IsIK() {
			ikLayerBone := b.getIkTreeIndex(bone)
			if ikLayerBone != nil {
				// 合致するIKツリーがある場合、そのレイヤーに登録
				b.IkTreeIndexes[ikLayerBone.Index] =
					append(b.IkTreeIndexes[ikLayerBone.Index], bone.Index)
				continue ikLoop
			}
			for _, link := range bone.Ik.Links {
				linkBone := b.GetItem(link.BoneIndex)
				linkLayerBone := b.getIkTreeIndex(linkBone)
				if linkLayerBone != nil {
					// 合致するIKツリーがある場合、そのレイヤーに登録
					b.IkTreeIndexes[linkLayerBone.Index] =
						append(b.IkTreeIndexes[linkLayerBone.Index], bone.Index)
					continue ikLoop
				}
			}

			// 関連親がIKツリーに登録されていない場合、新規にIKツリーを作成
			linkBone := b.GetItem(bone.Ik.Links[len(bone.Ik.Links)-1].BoneIndex)
			// b.IkTreeIndexes[linkBone.Index] = []int{bone.Index}
			if linkBone.ParentIndex >= 0 && b.Contains(linkBone.ParentIndex) {
				parentBone := b.GetItem(linkBone.ParentIndex)
				b.IkTreeIndexes[parentBone.Index] = []int{bone.Index}
			} else {
				b.IkTreeIndexes[bone.Index] = []int{bone.Index}
			}
		}
	}
}

// 変形階層とINDEXのソート用構造体
type LayerIndex struct {
	Layer int
	Index int
}

type LayerIndexes []LayerIndex

func (p LayerIndexes) Len() int {
	return len(p)
}
func (p LayerIndexes) Less(i, j int) bool {
	return p[i].Layer < p[j].Layer || (p[i].Layer == p[j].Layer && p[i].Index < p[j].Index)
}
func (p LayerIndexes) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p LayerIndexes) Contains(index int) bool {
	for _, layerIndex := range p {
		if layerIndex.Index == index {
			return true
		}
	}
	return false
}
