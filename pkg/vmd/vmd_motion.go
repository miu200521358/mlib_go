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
	m.BoneFrames.Get(boneName).Append(bf)
}

func (m *VmdMotion) AppendRegisteredBoneFrame(boneName string, bf *BoneFrame) {
	bf.Registered = true
	m.BoneFrames.Get(boneName).Append(bf)
}

func (m *VmdMotion) AppendMorphFrame(morphName string, mf *MorphFrame) {
	m.MorphFrames.Get(morphName).Append(mf)
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

func (m *VmdMotion) InsertBoneFrame(boneName string, bf *BoneFrame) {
	m.BoneFrames.Get(boneName).Insert(bf)
}

func (m *VmdMotion) InsertRegisteredBoneFrame(boneName string, bf *BoneFrame) {
	bf.Registered = true
	m.BoneFrames.Get(boneName).Insert(bf)
}

func (m *VmdMotion) InsertMorphFrame(morphName string, mf *MorphFrame) {
	m.MorphFrames.Get(morphName).Insert(mf)
}

func (m *VmdMotion) InsertCameraFrame(cf *CameraFrame) {
	m.CameraFrames.Insert(cf)
}

func (m *VmdMotion) InsertLightFrame(lf *LightFrame) {
	m.LightFrames.Insert(lf)
}

func (m *VmdMotion) InsertShadowFrame(sf *ShadowFrame) {
	m.ShadowFrames.Insert(sf)
}

func (m *VmdMotion) InsertIkFrame(ikf *IkFrame) {
	m.IkFrames.Insert(ikf)
}

func (m *VmdMotion) Animate(fno int, model *pmx.PmxModel) *VmdDeltas {
	vds := &VmdDeltas{}

	vds.Morphs = m.AnimateMorph(fno, model, nil)

	for i, bd := range vds.Morphs.Bones.Data {
		if bd == nil {
			continue
		}
		bone := model.Bones.Get(i)
		if !m.BoneFrames.Contains(bone.Name) {
			m.BoneFrames.Append(NewBoneNameFrames(bone.Name))
		}
		bf := m.BoneFrames.Get(bone.Name).Get(fno)

		// 一旦モーフの値をクリア
		bf.MorphPosition = nil
		bf.MorphLocalPosition = nil
		bf.MorphRotation = mmath.NewRotation()
		bf.MorphLocalRotation = mmath.NewRotation()
		bf.MorphScale = nil
		bf.MorphLocalScale = nil

		// 該当ボーンキーフレにモーフの値を加算
		bf.Add(bd.BoneFrame)
		m.AppendBoneFrame(bone.Name, bf)
	}

	// モーフ付きで変形を計算
	vds.Bones = m.animateBoneWithMorphs(fno, model, nil, true, true, true)

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

// AnimateBone IK結果をクリアして計算する
func (m *VmdMotion) AnimateBone(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
) *BoneDeltas {
	return m.animateBoneWithMorphs(frame, model, boneNames, isCalcIk, true, false)
}

// AnimateBoneContinueIk IK結果をクリアせずに計算する
func (m *VmdMotion) AnimateBoneContinueIk(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
) *BoneDeltas {
	return m.animateBoneWithMorphs(frame, model, boneNames, isCalcIk, false, false)
}

func (m *VmdMotion) animateBoneWithMorphs(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk, isClearIk, isCalcMorph bool,
) *BoneDeltas {
	return m.BoneFrames.Animate(frame, model, boneNames, isCalcIk, isClearIk, isCalcMorph)
}
