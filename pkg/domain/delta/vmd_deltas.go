// 指示: miu200521358
package delta

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
)

// VmdDeltas はボーン/モーフの差分をまとめる。
type VmdDeltas struct {
	frame      mtime.Frame
	modelHash  string
	motionHash string
	Bones      *BoneDeltas
	Morphs     *MorphDeltas
}

// NewVmdDeltas はVmdDeltasを生成する。
func NewVmdDeltas(frame mtime.Frame, bones *model.BoneCollection, modelHash, motionHash string) *VmdDeltas {
	return &VmdDeltas{
		frame:      frame,
		modelHash:  modelHash,
		motionHash: motionHash,
		Bones:      NewBoneDeltas(bones),
		Morphs:     nil,
	}
}

// Frame はフレーム番号を返す。
func (v *VmdDeltas) Frame() mtime.Frame {
	if v == nil {
		return 0
	}
	return v.frame
}

// SetFrame はフレーム番号を設定する。
func (v *VmdDeltas) SetFrame(frame mtime.Frame) {
	if v == nil {
		return
	}
	v.frame = frame
}

// ModelHash はモデルハッシュを返す。
func (v *VmdDeltas) ModelHash() string {
	if v == nil {
		return ""
	}
	return v.modelHash
}

// SetModelHash はモデルハッシュを設定する。
func (v *VmdDeltas) SetModelHash(hash string) {
	if v == nil {
		return
	}
	v.modelHash = hash
}

// MotionHash はモーションハッシュを返す。
func (v *VmdDeltas) MotionHash() string {
	if v == nil {
		return ""
	}
	return v.motionHash
}

// SetMotionHash はモーションハッシュを設定する。
func (v *VmdDeltas) SetMotionHash(hash string) {
	if v == nil {
		return
	}
	v.motionHash = hash
}
