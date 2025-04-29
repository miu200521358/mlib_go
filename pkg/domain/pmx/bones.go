package pmx

import (
	"slices"
	"sort"
	"strings"

	"fmt"

	"github.com/miu200521358/mlib_go/pkg/config/merr"
	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

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

func (b *Bones) GetStandardBoneNames() []string {
	boneNames := make([]string, 0)
	b.ForEach(func(index int, bone *Bone) bool {
		if bone.Config() != nil {
			boneNames = append(boneNames, bone.Name())
		}
		return true
	})
	return boneNames
}

// 指定された範囲のボーンの範囲のINDEXを取得
func (bones *Bones) Range(fromBoneIndex, toBoneIndex int) []int {
	fromIndex := slices.Index(bones.LayerSortedIndexes, fromBoneIndex)
	toIndex := slices.Index(bones.LayerSortedIndexes, toBoneIndex)

	return bones.LayerSortedIndexes[fromIndex : toIndex+1]
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

func (bones *Bones) Insert(bone *Bone) error {
	afterBoneIndex := bones.getInsertAfterIndex(bone)

	// 挿入位置を探す
	insertPos := -1

	if afterBoneIndex < 0 {
		// afterIndexが-1の場合、ルートに挿入
		bone.Layer = 0
		insertPos = len(bones.LayerSortedIndexes)
	} else {
		for i, boneIndex := range bones.LayerSortedIndexes {
			if boneIndex == afterBoneIndex {
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
		return nil
	}

	// 新しい要素のLayerを決定
	var newLayer int
	if insertPos == len(bones.LayerSortedIndexes) {
		// 挿入位置が最後の場合
		if boneAtPrevPos, err := bones.Get(bones.LayerSortedIndexes[insertPos-1]); err == nil {
			if afterBoneIndex >= 0 {
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
			return err
		}
		currentLayer := boneAtPrevPos.Layer
		boneAtNextPos, err := bones.Get(bones.LayerSortedIndexes[insertPos])
		if err != nil {
			return err
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

	return nil
}

func (bones *Bones) SetParentFromConfig(bone *Bone) {
	// 設定上の子どもボーンで最初に取得できたボーンを子ルートボーンとする
	isBreak := false
	for _, childNames := range bone.Config().ChildBoneNames {
		// 兄弟（並列の子ども）は全部処理する
		for _, childName := range childNames {
			var childRootBone *Bone
			var err error
			if childRootBone, err = bones.GetByName(childName.StringFromDirection(bone.Direction())); err != nil && childRootBone == nil {
				continue
			}

			// 親ボーンの子ども
			var parentBone *Bone
			if parentBone, err = bones.Get(bone.ParentIndex); err != nil && parentBone == nil {
				continue
			}

			for _, childIndex := range parentBone.ChildBoneIndexes {
				if childIndex == bone.Index() {
					continue
				}

				var childBone *Bone
				if childBone, err = bones.Get(childIndex); err != nil && childBone == nil {
					continue
				}

				if slices.Contains(childBone.RelativeBoneIndexes, childRootBone.Index()) ||
					childBone.index == childRootBone.Index() {
					// 子どもの関連ボーンに子ルートボーンがいる場合のみ対象
					childBone.ParentIndex = bone.Index()
					isBreak = true
				}
			}

			if !isBreak {
				childRootBone.ParentIndex = bone.Index()
				isBreak = true
			}

			// isInParent := false
			// for _, boneIndex := range bones.LayerSortedIndexes {
			// 	if boneIndex == parentBone.Index() {
			// 		isInParent = true
			// 		continue
			// 	}
			// 	if boneIndex == childRootBone.Index() {
			// 		break
			// 	}

			// 	if isInParent {
			// 		if b, err := bones.Get(boneIndex); err == nil && b != nil {
			// 			if b.ParentIndex == parentBone.Index() {
			// 				b.ParentIndex = bone.Index()
			// 				continue
			// 			}
			// 		}
			// 	}
			// }
		}
		if isBreak {
			// 子ボーンの親が設定できたら、全部のループを抜ける
			break
		}
	}
}

func (bones *Bones) GetIkTargetByName(ikBoneName string) (*Bone, error) {
	if ikBoneName == "" || !bones.ContainsByName(ikBoneName) {
		return nil, merr.NewNameNotFoundError(ikBoneName, "ik target not found from name")
	}

	if ikBone, err := bones.GetByName(ikBoneName); err != nil ||
		!ikBone.IsIK() || !bones.Contains(ikBone.Ik.BoneIndex) {
		return nil, merr.NewNameNotFoundError(ikBoneName, "ik target not found from name")
	} else {
		if ikTargetBone, err := bones.Get(ikBone.Ik.BoneIndex); err != nil {
			return nil, merr.NewNameNotFoundError(ikBoneName, "ik target not found from name")
		} else {
			return ikTargetBone, nil
		}
	}
}

func (bones *Bones) GetIkTarget(ikBoneIndex int) (*Bone, error) {
	if ikBone, err := bones.Get(ikBoneIndex); err != nil ||
		!ikBone.IsIK() || !bones.Contains(ikBone.Ik.BoneIndex) {
		return nil, merr.NewNameNotFoundError(fmt.Sprint(ikBoneIndex), "ik target not found from index")
	} else {
		if ikTargetBone, err := bones.Get(ikBone.Ik.BoneIndex); err != nil {
			return nil, merr.NewNameNotFoundError(fmt.Sprint(ikBoneIndex), "ik target not found from index")
		} else {
			return ikTargetBone, nil
		}
	}
}

func (bones *Bones) Setup() {
	bones.LayerSortedIndexes = make([]int, 0)
	bones.LayerSortedBones = make(map[bool][]*Bone)
	bones.LayerSortedNames = make(map[bool]map[string]int)
	bones.LayerSortedBoneIndexes = make(map[bool][]int)
	bones.DeformBoneIndexes = make(map[int][]int)

	bones.ForEach(func(index int, bone *Bone) bool {
		// 関係ボーンリストを一旦クリア
		bone.IkLinkBoneIndexes = make([]int, 0)
		bone.IkTargetBoneIndexes = make([]int, 0)
		bone.EffectiveBoneIndexes = make([]int, 0)
		bone.ChildBoneIndexes = make([]int, 0)
		bone.RelativeBoneIndexes = make([]int, 0)
		bone.ParentBoneIndexes = make([]int, 0)
		bone.ParentBoneNames = make([]string, 0)
		bone.TreeBoneIndexes = make([]int, 0)

		return true
	})

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

		if parentBone, err := bones.Get(bone.ParentIndex); err == nil && parentBone != nil {
			// 親ボーンに子ボーンとして登録する
			parentBone.ChildBoneIndexes = append(parentBone.ChildBoneIndexes, bone.Index())
			// 親ボーンを登録
			bone.ParentBone = parentBone
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

	bones.ForEach(func(index int, bone *Bone) bool {
		// ボーンのデフォームINDEXリストを取得
		bones.createLayerSortedBones(bone)
		return true
	})
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
	bones.ForEach(func(index int, bone *Bone) bool {
		layerIndexes = append(layerIndexes,
			layerIndex{isAfterPhysics: bone.IsAfterPhysicsDeform(), layer: bone.Layer, index: index})
		return true
	})
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
