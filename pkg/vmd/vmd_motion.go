package vmd

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/mcore"
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

// AnimateBone 物理前ボーンのデフォーム計算する
func (m *VmdMotion) Animate(
	frame int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk, isClearIk, isCalcMorph bool,
) *VmdDeltas {
	vds := &VmdDeltas{}

	if isCalcMorph {
		vds.Morphs = m.AnimateMorph(frame, model, nil)
	}

	// 物理前のデフォーム情報
	beforeBoneDeltas := m.BoneFrames.Animate(frame, model, boneNames, isCalcIk, isClearIk, isCalcMorph, nil)

	// 物理後のデフォーム情報
	vds.Bones = m.BoneFrames.Animate(frame, model, boneNames, isCalcIk, isClearIk, isCalcMorph, beforeBoneDeltas)

	return vds
}
