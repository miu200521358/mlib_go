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

func (m *VmdMotion) Animate(fno float32, model *pmx.PmxModel) BoneTrees {
	return m.AnimateBone([]float32{fno}, model, nil, true, false, "")
}

func (m *VmdMotion) AnimateBone(
	frames []float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// ボーン変形行列操作
	return m.BoneFrames.Animate(frames, model, boneNames, isCalcIk, isOutLog, description)
}
