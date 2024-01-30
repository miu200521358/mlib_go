package vmd

import (
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
		HashModel:  mcore.NewHashModel(path),
		ModelName:  "",
		BoneFrames: NewBoneFrames(),
	}
}

func (m *VmdMotion) GetName() string {
	return m.ModelName
}

func (m *VmdMotion) SetName(name string) {
	m.ModelName = name
}

func (m *VmdMotion) AppendBoneFrame(boneName string, bf *BoneFrame) {
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

func (m *VmdMotion) Animate(fno int, model *pmx.PmxModel) BoneTrees {
	return m.AnimateBone([]int{fno}, model, nil, false, false, "")
}

func (m *VmdMotion) AnimateBone(
	fnos []int,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// ボーン変形行列操作
	return m.BoneFrames.Animate(fnos, model, boneNames, isCalcIk, isOutLog, description)
}
