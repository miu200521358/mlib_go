package vmd

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/tiendc/go-deepcopy"
)

type VmdMotion struct {
	name         string
	path         string
	hash         string
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
		BoneFrames:   NewBoneFrames(),
		MorphFrames:  NewMorphFrames(),
		CameraFrames: NewCameraFrames(),
		LightFrames:  NewLightFrames(),
		ShadowFrames: NewShadowFrames(),
		IkFrames:     NewIkFrames(),
	}
}

func (motion *VmdMotion) IsVpd() bool {
	return strings.Contains(strings.ToLower(motion.path), ".vpd")
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
	h.Write([]byte(fmt.Sprintf("%d", motion.BoneFrames.Length())))
	h.Write([]byte(fmt.Sprintf("%d", motion.MorphFrames.Length())))
	h.Write([]byte(fmt.Sprintf("%d", motion.CameraFrames.Length())))
	h.Write([]byte(fmt.Sprintf("%d", motion.LightFrames.Length())))
	h.Write([]byte(fmt.Sprintf("%d", motion.ShadowFrames.Length())))
	h.Write([]byte(fmt.Sprintf("%d", motion.IkFrames.Length())))

	// ハッシュ値を16進数文字列に変換
	motion.SetHash(fmt.Sprintf("%x", h.Sum(nil)))
}

func (motion *VmdMotion) MaxFrame() float32 {
	return max(motion.BoneFrames.MaxFrame(), motion.MorphFrames.MaxFrame())
}

func (motion *VmdMotion) MinFrame() float32 {
	return min(motion.BoneFrames.MinFrame(), motion.MorphFrames.MinFrame())
}

func (motion *VmdMotion) RegisteredIndexes() []int {
	boneFrames := motion.BoneFrames.RegisteredIndexes()
	morphFrames := motion.MorphFrames.RegisteredIndexes()

	frames := make([]int, 0, len(boneFrames)+len(morphFrames))
	for f := range boneFrames {
		frames = append(frames, int(f))
	}
	for f := range morphFrames {
		frames = append(frames, int(f))
	}

	mmath.Unique(frames)
	mmath.Sort(frames)

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

func (motion *VmdMotion) Copy() (*VmdMotion, error) {
	copied := new(VmdMotion)
	err := deepcopy.Copy(copied, motion)
	return copied, err
}
