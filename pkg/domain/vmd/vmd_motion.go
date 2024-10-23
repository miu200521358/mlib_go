package vmd

import (
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type VmdMotion struct {
	name         string
	path         string
	hash         string
	Processing   bool   // 処理中フラグ
	Signature    string // vmdバージョン
	BoneFrames   *BoneFrames
	MorphFrames  *MorphFrames
	CameraFrames *CameraFrames
	LightFrames  *LightFrames
	ShadowFrames *ShadowFrames
	IkFrames     *IkFrames
}

func NewVmdMotion(path string) *VmdMotion {
	return &VmdMotion{
		name:         "",
		path:         path,
		hash:         "",
		Processing:   false,
		BoneFrames:   NewBoneFrames(),
		MorphFrames:  NewMorphFrames(),
		CameraFrames: NewCameraFrames(),
		LightFrames:  NewLightFrames(),
		ShadowFrames: NewShadowFrames(),
		IkFrames:     NewIkFrames(),
	}
}

func (motion *VmdMotion) Path() string {
	return motion.path
}

func (motion *VmdMotion) SetPath(path string) {
	motion.path = path
}

func (motion *VmdMotion) Name() string {
	return motion.name
}

func (motion *VmdMotion) SetName(name string) {
	motion.name = name
}

func (motion *VmdMotion) Hash() string {
	return motion.hash
}

func (motion *VmdMotion) SetHash(hash string) {
	motion.hash = hash
}

func (motion *VmdMotion) SetRandHash() {
	motion.hash = fmt.Sprintf("%d", rand.Intn(10000))
}

func (motion *VmdMotion) UpdateHash() {

	h := fnv.New32a()
	// 名前をハッシュに含める
	h.Write([]byte(motion.Name()))
	// ファイルパスをハッシュに含める
	h.Write([]byte(motion.Path()))
	// 各要素の数をハッシュに含める
	h.Write([]byte(fmt.Sprintf("%d", motion.BoneFrames.Len())))
	h.Write([]byte(fmt.Sprintf("%d", motion.MorphFrames.Len())))
	h.Write([]byte(fmt.Sprintf("%d", motion.CameraFrames.Len())))
	h.Write([]byte(fmt.Sprintf("%d", motion.LightFrames.Len())))
	h.Write([]byte(fmt.Sprintf("%d", motion.ShadowFrames.Len())))
	h.Write([]byte(fmt.Sprintf("%d", motion.IkFrames.Len())))

	// ハッシュ値を16進数文字列に変換
	motion.SetHash(fmt.Sprintf("%x", h.Sum(nil)))
}

func (motion *VmdMotion) Copy() core.IHashModel {
	return &VmdMotion{
		name:         motion.name,
		path:         motion.path,
		hash:         motion.hash,
		Signature:    motion.Signature,
		BoneFrames:   motion.BoneFrames.Copy(),
		MorphFrames:  motion.MorphFrames.Copy(),
		CameraFrames: motion.CameraFrames.Copy(),
		LightFrames:  motion.LightFrames.Copy(),
		ShadowFrames: motion.ShadowFrames.Copy(),
		IkFrames:     motion.IkFrames.Copy(),
	}
}

func (motion *VmdMotion) MaxFrame() float32 {
	return max(motion.BoneFrames.MaxFrame(), motion.MorphFrames.MaxFrame())
}

func (motion *VmdMotion) MinFrame() float32 {
	return min(motion.BoneFrames.MinFrame(), motion.MorphFrames.MinFrame())
}

func (motion *VmdMotion) RegisteredFrames(boneNames []string) []int {
	boneFrames := motion.BoneFrames.registeredFramesMap(boneNames)
	morphFrames := motion.MorphFrames.registeredFramesMap()

	frames := make([]int, 0, len(boneFrames)+len(morphFrames))
	for f := range boneFrames {
		frames = append(frames, int(f))
	}
	for f := range morphFrames {
		frames = append(frames, int(f))
	}
	mmath.SortInts(frames)

	return frames
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

func (motion *VmdMotion) Clean() {
	motion.BoneFrames.Clean()
	motion.MorphFrames.Clean()
	motion.CameraFrames.Clean()
	motion.LightFrames.Clean()
	motion.ShadowFrames.Clean()
	motion.IkFrames.Clean()
}
