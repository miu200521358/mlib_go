package deform

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type BoneDelta struct {
	BoneName                   string             // ボーン名
	Frame                      float32            // キーフレーム
	GlobalMatrix               *mmath.MMat4       // グローバル行列
	LocalMatrix                *mmath.MMat4       // ローカル行列
	Position                   *mmath.MVec3       // グローバル位置
	FramePosition              *mmath.MVec3       // キーフレ位置の変動量
	FrameRotation              *mmath.MQuaternion // キーフレ回転の変動量
	FrameRotationWithoutEffect *mmath.MQuaternion // キーフレ回転の変動量(付与親無視)
	FrameScale                 *mmath.MVec3       // キーフレスケールの変動量
	Matrix                     *mmath.MMat4       // ボーンの変動行列
}

func NewBoneDelta(
	boneName string,
	frame float32,
	globalMatrix, localMatrix *mmath.MMat4,
	framePosition *mmath.MVec3,
	frameRotation *mmath.MQuaternion,
	frameRotationWithoutEffect *mmath.MQuaternion,
	frameScale *mmath.MVec3,
	matrix *mmath.MMat4,
) *BoneDelta {
	p := globalMatrix.Translation()
	return &BoneDelta{
		BoneName:                   boneName,
		Frame:                      frame,
		GlobalMatrix:               globalMatrix,
		LocalMatrix:                localMatrix,
		Position:                   p,
		FramePosition:              framePosition,
		FrameRotation:              frameRotation,
		FrameRotationWithoutEffect: frameRotationWithoutEffect,
		FrameScale:                 frameScale,
		Matrix:                     matrix,
	}
}

type BoneNameFrameNo struct {
	BoneName string
	Frame    float32
}

type BoneDeltas struct {
	Data map[BoneNameFrameNo]*BoneDelta
}

func NewBoneDeltas() *BoneDeltas {
	return &BoneDeltas{
		Data: make(map[BoneNameFrameNo]*BoneDelta, 0),
	}
}

func (bts *BoneDeltas) GetItem(boneName string, frame float32) *BoneDelta {
	return bts.Data[BoneNameFrameNo{boneName, frame}]
}

func (bts *BoneDeltas) SetItem(boneName string, frame float32, boneDelta *BoneDelta) {
	bts.Data[BoneNameFrameNo{boneName, frame}] = boneDelta
}

func (bts *BoneDeltas) GetBoneNames() []string {
	boneNames := make([]string, 0)
	for key := range bts.Data {
		if !slices.Contains(boneNames, key.BoneName) {
			boneNames = append(boneNames, key.BoneName)
		}
	}
	return boneNames
}

func (bts *BoneDeltas) GetFrameNos() []float32 {
	frames := make([]float32, 0)
	for key := range bts.Data {
		if slices.Contains(frames, key.Frame) {
			frames = append(frames, key.Frame)
		}
	}
	return frames
}

func (bts *BoneDeltas) Contains(boneName string, frame float32) bool {
	_, ok := bts.Data[BoneNameFrameNo{boneName, frame}]
	return ok
}
