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

func (m *VmdMotion) GetMaxFrame() float32 {
	// TODO: モーフが入ったらモーフも考慮する
	return m.BoneFrames.GetMaxFrame()
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
	return m.AnimateBone(fno, model, nil, true, false, "")
}

func (m *VmdMotion) AnimateBone(
	frame float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isOutLog bool,
	description string,
) BoneTrees {
	// ボーン変形行列操作
	// IKリンクボーンの回転量を初期化
	for _, bnfs := range m.BoneFrames.Data {
		for _, bf := range bnfs.Data {
			bf.IkRotation = nil
		}
	}

	for _, bone := range model.Bones.Data {
		// ボーンフレームを生成
		if !m.BoneFrames.Contains(bone.Name) {
			m.BoneFrames.Append(NewBoneNameFrames(bone.Name))
		}
	}

	return m.BoneFrames.Animate(frame, model, boneNames, isCalcIk, isOutLog, description)
}
