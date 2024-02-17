package pmx

import (
	"embed"
	"slices"
	"sort"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

type PmxModel struct {
	*mcore.HashModel
	Signature           string
	Version             float64
	ExtendedUVCountType int
	VertexCountType     int
	TextureCountType    int
	MaterialCountType   int
	BoneCountType       int
	MorphCountType      int
	RigidBodyCountType  int
	Name                string
	EnglishName         string
	Comment             string
	EnglishComment      string
	JsonData            map[string]interface{}
	Vertices            *Vertices
	Faces               *Faces
	Textures            *Textures
	ToonTextures        *ToonTextures
	Materials           *Materials
	Bones               *Bones
	Morphs              *Morphs
	DisplaySlots        *DisplaySlots
	RigidBodies         *RigidBodies
	Joints              *Joints
	VerticesByBones     map[int][]int
	VerticesByMaterials map[int][]int
	FacesByMaterials    map[int][]int
	Meshes              *Meshes
}

func NewPmxModel(path string) *PmxModel {
	model := &PmxModel{}
	model.HashModel = mcore.NewHashModel(path)

	model.Vertices = NewVertices()
	model.Faces = NewFaces()
	model.Textures = NewTextures()
	model.ToonTextures = NewToonTextures()
	model.Materials = NewMaterials()
	model.Bones = NewBones()
	model.Morphs = NewMorphs()
	model.DisplaySlots = NewDisplaySlots()
	model.RigidBodies = NewRigidBodies()
	model.Joints = NewJoints()
	model.VerticesByBones = make(map[int][]int)
	model.VerticesByMaterials = make(map[int][]int)
	model.FacesByMaterials = make(map[int][]int)

	return model
}

func (pm *PmxModel) InitializeDisplaySlots() {
	d01 := NewDisplaySlot()
	d01.Name = "Root"
	d01.EnglishName = "Root"
	d01.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d01)

	d02 := NewDisplaySlot()
	d02.Name = "表情"
	d02.EnglishName = "Exp"
	d02.SpecialFlag = SPECIAL_FLAG_ON
	pm.DisplaySlots.Append(d02)
}

func (pm *PmxModel) InitializeDraw(windowIndex int, resourceFiles embed.FS) {
	pm.ToonTextures.InitGl(windowIndex, resourceFiles)
	pm.Meshes = NewMeshes(pm, windowIndex, resourceFiles)
}

func (pm *PmxModel) Draw(
	shader *mgl.MShader,
	boneMatrixes []*mgl32.Mat4,
	windowIndex int,
) {
	pm.Meshes.Draw(shader, boneMatrixes, windowIndex)
}

// 関連ボーンリストの取得
func (pm *PmxModel) GetRelativeBoneIndexes(boneIndex int, parentBoneIndexes, relativeBoneIndexes []int) ([]int, []int) {

	if boneIndex <= 0 || !pm.Bones.Contains(boneIndex) {
		return parentBoneIndexes, relativeBoneIndexes
	}

	bone := pm.Bones.GetItem(boneIndex)
	if pm.Bones.Contains(bone.ParentIndex) && !slices.Contains(relativeBoneIndexes, bone.ParentIndex) {
		// 親ボーンを辿る(親から子の順番)
		parentBoneIndexes = append([]int{bone.ParentIndex}, parentBoneIndexes...)
		relativeBoneIndexes = append(relativeBoneIndexes, bone.ParentIndex)
		parentBoneIndexes, relativeBoneIndexes =
			pm.GetRelativeBoneIndexes(bone.ParentIndex, parentBoneIndexes, relativeBoneIndexes)
	}
	if (bone.IsEffectorRotation() || bone.IsEffectorTranslation()) &&
		pm.Bones.Contains(bone.EffectIndex) && !slices.Contains(relativeBoneIndexes, bone.EffectIndex) {
		// 付与親ボーンを辿る
		relativeBoneIndexes = append(relativeBoneIndexes, bone.EffectIndex)
		_, relativeBoneIndexes =
			pm.GetRelativeBoneIndexes(bone.EffectIndex, parentBoneIndexes, relativeBoneIndexes)
	}
	if bone.IsIK() {
		if pm.Bones.Contains(bone.Ik.BoneIndex) && !slices.Contains(relativeBoneIndexes, bone.Ik.BoneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, bone.Ik.BoneIndex)
			_, relativeBoneIndexes =
				pm.GetRelativeBoneIndexes(bone.Ik.BoneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
		for _, link := range bone.Ik.Links {
			if pm.Bones.Contains(link.BoneIndex) && !slices.Contains(relativeBoneIndexes, link.BoneIndex) {
				// IKリンクボーンを辿る
				relativeBoneIndexes = append(relativeBoneIndexes, link.BoneIndex)
				_, relativeBoneIndexes =
					pm.GetRelativeBoneIndexes(link.BoneIndex, parentBoneIndexes, relativeBoneIndexes)
			}
		}
	}
	for _, boneIndex := range bone.EffectiveBoneIndexes {
		if pm.Bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// 外部子ボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				pm.GetRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.IkTargetBoneIndexes {
		if pm.Bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKターゲットボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				pm.GetRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}
	for _, boneIndex := range bone.IkLinkBoneIndexes {
		if pm.Bones.Contains(boneIndex) && !slices.Contains(relativeBoneIndexes, boneIndex) {
			// IKリンクボーンを辿る
			relativeBoneIndexes = append(relativeBoneIndexes, boneIndex)
			_, relativeBoneIndexes =
				pm.GetRelativeBoneIndexes(boneIndex, parentBoneIndexes, relativeBoneIndexes)
		}
	}

	return parentBoneIndexes, relativeBoneIndexes
}

func (pm *PmxModel) SetUp() {
	for _, bone := range pm.Bones.Data {
		// 関係ボーンリストを一旦クリア
		bone.IkLinkBoneIndexes = make([]int, 0)
		bone.IkTargetBoneIndexes = make([]int, 0)
		bone.EffectiveBoneIndexes = make([]int, 0)
		bone.ChildBoneIndexes = make([]int, 0)
		bone.ChildIkBoneIndexes = make([]int, 0)
	}

	// IK親（子として登録されていない）ボーンINDEXリスト
	parentIkBoneIndexes := make([]int, 0)
	// 関連ボーンINDEX情報を設定
	for _, bone := range pm.Bones.GetSortedData() {
		if bone.IsIK() && bone.Ik != nil {
			parentIkBoneIndexes = append(parentIkBoneIndexes, bone.Index)
			// IKのリンクとターゲット
			for _, link := range bone.Ik.Links {
				if pm.Bones.Contains(link.BoneIndex) &&
					!slices.Contains(pm.Bones.GetItem(link.BoneIndex).IkLinkBoneIndexes, bone.Index) {
					// リンクボーンにフラグを立てる
					linkBone := pm.Bones.GetItem(link.BoneIndex)
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
			if pm.Bones.Contains(bone.Ik.BoneIndex) &&
				!slices.Contains(pm.Bones.GetItem(bone.Ik.BoneIndex).IkTargetBoneIndexes, bone.Index) {
				// ターゲットボーンにもフラグを立てる
				pm.Bones.GetItem(bone.Ik.BoneIndex).IkTargetBoneIndexes = append(pm.Bones.GetItem(bone.Ik.BoneIndex).IkTargetBoneIndexes, bone.Index)
			}
		}
		if bone.EffectIndex >= 0 && pm.Bones.Contains(bone.EffectIndex) &&
			!slices.Contains(pm.Bones.GetItem(bone.EffectIndex).EffectiveBoneIndexes, bone.Index) {
			// 付与親の方に付与子情報を保持
			pm.Bones.GetItem(bone.EffectIndex).EffectiveBoneIndexes = append(pm.Bones.GetItem(bone.EffectIndex).EffectiveBoneIndexes, bone.Index)
		}
	}

	for _, bone := range pm.Bones.GetSortedData() {
		// 影響があるボーンINDEXリスト
		bone.ParentBoneIndexes, bone.RelativeBoneIndexes = pm.GetRelativeBoneIndexes(bone.Index, []int{}, []int{})

		// 親ボーンに子ボーンとして登録する
		if bone.ParentIndex >= 0 && pm.Bones.Contains(bone.ParentIndex) {
			pm.Bones.GetItem(bone.ParentIndex).ChildBoneIndexes = append(pm.Bones.GetItem(bone.ParentIndex).ChildBoneIndexes, bone.Index)
		}
		// 親からの相対位置
		bone.ParentRelativePosition = pm.Bones.getParentRelativePosition(bone.Index)
		// 子への相対位置
		bone.ChildRelativePosition = pm.Bones.getChildRelativePosition(bone.Index)
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

	// IK子ボーンINDEXリスト
	ikBoneLayerIndexes := make(map[int]LayerIndexes, 0)
	for _, bone := range pm.Bones.GetSortedData() {
		if bone.IsIK() {
			ikBoneLayerIndexes[bone.Index] =
				append(ikBoneLayerIndexes[bone.Index], LayerIndex{Index: bone.Index, Layer: bone.Layer})
			for _, parentIndex := range bone.ParentBoneIndexes {
				parentBone := pm.Bones.GetItem(parentIndex)
				if parentBone.IsIK() {
					// IK子ボーンとして追加
					if _, ok := ikBoneLayerIndexes[parentBone.Index]; !ok {
						ikBoneLayerIndexes[parentBone.Index] = make(LayerIndexes, 0)
					}
					ikBoneLayerIndexes[parentBone.Index] =
						append(ikBoneLayerIndexes[parentBone.Index], LayerIndex{Index: bone.Index, Layer: bone.Layer})
					// IK親ボーンINDEXリストからbone.Indexを削除
					parentIkBoneIndexes = mutils.RemoveFromSlice(parentIkBoneIndexes, bone.Index)
				}
				if len(parentBone.IkLinkBoneIndexes) > 0 {
					// IKリンクボーンとして追加
					for _, linkIkIndex := range parentBone.IkLinkBoneIndexes {
						linkIkBone := pm.Bones.GetItem(linkIkIndex)
						if linkIkBone.IsIK() {
							// IK子ボーンとして追加
							if _, ok := ikBoneLayerIndexes[linkIkBone.Index]; !ok {
								ikBoneLayerIndexes[linkIkBone.Index] = make(LayerIndexes, 0)
							}
							ikBoneLayerIndexes[linkIkBone.Index] =
								append(ikBoneLayerIndexes[linkIkBone.Index],
									LayerIndex{Index: bone.Index, Layer: bone.Layer})
							// IK親ボーンINDEXリストからbone.Indexを削除
							parentIkBoneIndexes = mutils.RemoveFromSlice(parentIkBoneIndexes, bone.Index)
						}
					}
				}
			}
		}
	}

	// 並列計算可能な親IKボーンの子IKボーンを変形階層でソートして設定
	for boneIndex, childLayerIndexes := range ikBoneLayerIndexes {
		if slices.Contains(parentIkBoneIndexes, boneIndex) {
			sort.Sort(childLayerIndexes)
			for _, childLayerIndex := range childLayerIndexes {
				pm.Bones.GetItem(boneIndex).ChildIkBoneIndexes = append(pm.Bones.GetItem(boneIndex).ChildIkBoneIndexes, childLayerIndex.Index)
			}
		}
	}

	// 変形階層・ボーンINDEXでソート
	pm.Bones.LayerSortedIndexes = make(map[int]string, len(pm.Bones.Data))
	pm.Bones.LayerSortedNames = make(map[string]int, len(pm.Bones.Data))

	i := 0
	for _, boneIndex := range pm.Bones.GetLayerIndexes() {
		bone := pm.Bones.GetItem(boneIndex)
		pm.Bones.LayerSortedNames[bone.Name] = i
		pm.Bones.LayerSortedIndexes[i] = bone.Name
		i++
	}

}
