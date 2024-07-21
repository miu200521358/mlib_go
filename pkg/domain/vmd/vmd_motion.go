package vmd

import (
	"github.com/jinzhu/copier"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
)

type VmdMotion struct {
	*core.HashModel
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
		HashModel:    core.NewHashModel(path),
		ModelName:    "",
		BoneFrames:   NewBoneFrames(),
		MorphFrames:  NewMorphFrames(),
		CameraFrames: NewCameraFrames(),
		LightFrames:  NewLightFrames(),
		ShadowFrames: NewShadowFrames(),
		IkFrames:     NewIkFrames(),
	}
}

func (motion *VmdMotion) Copy() core.IHashModel {
	copied := NewVmdMotion("")
	copier.CopyWithOption(copied, motion, copier.Option{DeepCopy: true})
	return copied
}

func (motion *VmdMotion) Name() string {
	return motion.ModelName
}

func (motion *VmdMotion) SetName(name string) {
	motion.ModelName = name
}

func (motion *VmdMotion) MaxFrame() int {
	return max(motion.BoneFrames.MaxFrame(), motion.MorphFrames.MaxFrame())
}

func (motion *VmdMotion) MinFrame() int {
	return min(motion.BoneFrames.MinFrame(), motion.MorphFrames.MinFrame())
}

func (motion *VmdMotion) AppendBoneFrame(boneName string, bf *BoneFrame) {
	motion.BoneFrames.Get(boneName).Append(bf)
}

func (motion *VmdMotion) AppendRegisteredBoneFrame(boneName string, bf *BoneFrame) {
	bf.Registered = true
	motion.BoneFrames.Get(boneName).Append(bf)
}

func (motion *VmdMotion) AppendMorphFrame(morphName string, mf *MorphFrame) {
	motion.MorphFrames.Get(morphName).Append(mf)
}

func (motion *VmdMotion) AppendRegisteredMorphFrame(morphName string, mf *MorphFrame) {
	mf.Registered = true
	motion.MorphFrames.Get(morphName).Append(mf)
}

func (motion *VmdMotion) AppendCameraFrame(cf *CameraFrame) {
	motion.CameraFrames.Append(cf)
}

func (motion *VmdMotion) AppendLightFrame(lf *LightFrame) {
	motion.LightFrames.Append(lf)
}

func (motion *VmdMotion) AppendShadowFrame(sf *ShadowFrame) {
	motion.ShadowFrames.Append(sf)
}

func (motion *VmdMotion) AppendIkFrame(ikf *IkFrame) {
	motion.IkFrames.Append(ikf)
}

func (motion *VmdMotion) InsertBoneFrame(boneName string, bf *BoneFrame) {
	motion.BoneFrames.Get(boneName).Insert(bf)
}

func (motion *VmdMotion) InsertRegisteredBoneFrame(boneName string, bf *BoneFrame) {
	bf.Registered = true
	motion.BoneFrames.Get(boneName).Insert(bf)
}

func (motion *VmdMotion) InsertMorphFrame(morphName string, mf *MorphFrame) {
	motion.MorphFrames.Get(morphName).Insert(mf)
}

func (motion *VmdMotion) InsertCameraFrame(cf *CameraFrame) {
	motion.CameraFrames.Insert(cf)
}

func (motion *VmdMotion) InsertLightFrame(lf *LightFrame) {
	motion.LightFrames.Insert(lf)
}

func (motion *VmdMotion) InsertShadowFrame(sf *ShadowFrame) {
	motion.ShadowFrames.Insert(sf)
}

func (motion *VmdMotion) InsertIkFrame(ikf *IkFrame) {
	motion.IkFrames.Insert(ikf)
}
