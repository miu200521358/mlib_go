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

func (m *VmdMotion) Copy() mcore.HashModelInterface {
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

func (m *VmdMotion) GetMaxFrame() float32 {
	// TODO: モーフが入ったらモーフも考慮する
	return max(m.BoneFrames.GetMaxFrame(), m.MorphFrames.GetMaxFrame())
}

func (m *VmdMotion) AppendBoneFrame(boneName string, bf *BoneFrame, isSort bool) {
	m.BoneFrames.GetItem(boneName).Append(bf, isSort)
}

func (m *VmdMotion) SortBoneFrames() {
	for _, bnfs := range m.BoneFrames.Data {
		bnfs.Sort()
	}
}

func (m *VmdMotion) AppendMorphFrame(morphName string, mf *MorphFrame, isSort bool) {
	m.MorphFrames.GetItem(morphName).Append(mf, isSort)
}

func (m *VmdMotion) SortMorphFrames() {
	for _, mnfs := range m.MorphFrames.Data {
		mnfs.Sort()
	}
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

func (m *VmdMotion) Animate(fno float32, model *pmx.PmxModel) *VmdDeltas {
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
		m.AppendBoneFrame(bone.Name, bf, false)
	}

	// モーフ付きで変形を計算
	vds.Bones = m.AnimateBoneWithMorphs(fno, model, nil, true, true)

	return vds
}

func (m *VmdMotion) AnimateMorph(
	frame float32,
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
	frame float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
) *BoneDeltas {
	return m.AnimateBoneWithMorphs(frame, model, boneNames, isCalcIk, false)
}

func (m *VmdMotion) AnimateBoneWithMorphs(
	frame float32,
	model *pmx.PmxModel,
	boneNames []string,
	isCalcIk bool,
	isCalcMorph bool,
) *BoneDeltas {
	// ボーン変形行列操作
	// IKリンクボーンの回転量を初期化
	for _, bnfs := range m.BoneFrames.Data {
		for _, bf := range bnfs.Data {
			bf.IkRotation = mmath.NewRotationModel()
		}
	}

	for _, bone := range model.Bones.Data {
		// ボーンフレームを生成
		if !m.BoneFrames.Contains(bone.Name) {
			m.BoneFrames.Append(NewBoneNameFrames(bone.Name))
		}
	}

	return m.BoneFrames.Animate(frame, model, boneNames, isCalcIk, isCalcMorph)
}
