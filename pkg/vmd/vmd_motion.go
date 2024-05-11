package vmd

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type VmdMotion struct {
	*mcore.HashModel
	Signature    string // vmdバージョン
	ModelName    string // モデル名
	BoneFrames   *BoneFrames
	MorphFrames  *MorphFrames
	CameraFrames *CameraFrames
	LightFrames  *LightFrames
	ShadowFrames *ShadowFrames
	IkFrames     *IkFrames
}

func NewVmdMotion(path string) *VmdMotion {
	return &VmdMotion{
		HashModel:    mcore.NewHashModel(path),
		ModelName:    "",
		BoneFrames:   NewBoneFrames(),
		MorphFrames:  NewMorphFrames(),
		CameraFrames: NewCameraFrames(),
		LightFrames:  NewLightFrames(),
		ShadowFrames: NewShadowFrames(),
		IkFrames:     NewIkFrames(),
	}
}

func (m *VmdMotion) Copy() mcore.IHashModel {
	copied := NewVmdMotion("")
	copier.CopyWithOption(copied, m, copier.Option{DeepCopy: true})
	return copied
}

func (m *VmdMotion) GetName() string {
	return m.ModelName
}

func (m *VmdMotion) SetName(name string) {
	m.ModelName = name
}

func (m *VmdMotion) GetMaxFrame() int {
	return max(m.BoneFrames.GetMaxFrame(), m.MorphFrames.GetMaxFrame())
}

func (m *VmdMotion) GetMinFrame() int {
	return min(m.BoneFrames.GetMinFrame(), m.MorphFrames.GetMinFrame())
}

func (m *VmdMotion) AppendBoneFrame(boneName string, bf *BoneFrame) {
	m.BoneFrames.GetItem(boneName).Append(bf)
}

func (m *VmdMotion) AppendRegisteredBoneFrame(boneName string, bf *BoneFrame) {
	bf.Registered = true
	m.BoneFrames.GetItem(boneName).Append(bf)
}

func (m *VmdMotion) AppendMorphFrame(morphName string, mf *MorphFrame) {
	m.MorphFrames.GetItem(morphName).Append(mf)
}

func (m *VmdMotion) AppendCameraFrame(cf *CameraFrame) {
	m.CameraFrames.Append(cf)
}

func (m *VmdMotion) AppendLightFrame(lf *LightFrame) {
	m.LightFrames.Append(lf)
}

func (m *VmdMotion) AppendShadowFrame(sf *ShadowFrame) {
	m.ShadowFrames.Append(sf)
}

func (m *VmdMotion) AppendIkFrame(ikf *IkFrame) {
	m.IkFrames.Append(ikf)
}

func (m *VmdMotion) Animate(fno int, model *pmx.PmxModel) *VmdDeltas {
	vds := &VmdDeltas{}

	vds.Morphs = m.AnimateMorph(fno, model, nil)

	for i, bd := range vds.Morphs.Bones.Data {
		bone := model.Bones.GetItem(i)
		if !m.BoneFrames.Contains(bone.Name) {
			m.BoneFrames.Append(NewBoneNameFrames(bone.Name))
		}
		bf := m.BoneFrames.GetItem(bone.Name).GetItem(fno)

		// 一旦モーフの値をクリア
		bf.MorphPosition = mmath.NewMVec3()
		bf.MorphLocalPosition = mmath.NewMVec3()
		bf.MorphRotation.SetQuaternion(mmath.NewMQuaternion())
		bf.MorphLocalRotation.SetQuaternion(mmath.NewMQuaternion())
		bf.MorphScale = mmath.NewMVec3()
		bf.MorphLocalScale = mmath.NewMVec3()

		// 該当ボーンキーフレにモーフの値を加算
		bf.Add(bd.BoneFrame)
		m.AppendBoneFrame(bone.Name, bf)
	}

	// モーフ付きで変形を計算
	vds.Bones = m.AnimateBoneWithMorphs(fno, model, nil, true, true)

	return vds
}

func (m *VmdMotion) AnimateMorph(
	frame int,
	model *pmx.PmxModel,
	morphNames []string,
) *MorphDeltas {
	if morphNames == nil {
		morphNames = make([]string, 0)
	}

	for _, morph := range model.Morphs.Data {
		// モーフフレームを生成
		if !m.MorphFrames.Contains(morph.Name) {
			m.MorphFrames.Append(NewMorphNameFrames(morph.Name))
		}
		morphNames = append(morphNames, morph.Name)
	}

	return m.MorphFrames.Animate(frame, model, morphNames)
}

func (m *VmdMotion) AnimateBone(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
) *BoneDeltas {
	return m.AnimateBoneWithMorphs(frame, model, boneNames, isCalcIk, false)
}

func (m *VmdMotion) AnimateBoneWithMorphs(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isCalcMorph bool,
) *BoneDeltas {
	return m.BoneFrames.Animate(frame, model, boneNames, isCalcIk, isCalcMorph)
}
