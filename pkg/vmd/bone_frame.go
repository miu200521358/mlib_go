package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type BoneFrame struct {
	*BaseFrame                          // キーフレ
	Position           *mmath.MVec3     // 位置
	MorphPosition      *mmath.MVec3     // モーフ位置
	LocalPosition      *mmath.MVec3     // ローカル位置
	MorphLocalPosition *mmath.MVec3     // モーフローカル位置
	Rotation           *mmath.MRotation // 回転
	MorphRotation      *mmath.MRotation // モーフ回転
	LocalRotation      *mmath.MRotation // ローカル回転
	MorphLocalRotation *mmath.MRotation // モーフローカル回転
	Scale              *mmath.MVec3     // スケール
	MorphScale         *mmath.MVec3     // モーフスケール
	LocalScale         *mmath.MVec3     // ローカルスケール
	MorphLocalScale    *mmath.MVec3     // モーフローカルスケール
	IkRotation         *mmath.MRotation // IK回転
	Curves             *BoneCurves      // 補間曲線
	IkRegistered       bool             // IK計算済み
}

func NewBoneFrame(index float32) *BoneFrame {
	return &BoneFrame{
		BaseFrame:          NewVmdBaseFrame(index),
		Position:           mmath.NewMVec3(),
		MorphPosition:      mmath.NewMVec3(),
		LocalPosition:      mmath.NewMVec3(),
		MorphLocalPosition: mmath.NewMVec3(),
		Rotation:           mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		MorphRotation:      mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		LocalRotation:      mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		MorphLocalRotation: mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		Scale:              mmath.NewMVec3(),
		MorphScale:         mmath.NewMVec3(),
		LocalScale:         mmath.NewMVec3(),
		MorphLocalScale:    mmath.NewMVec3(),
		IkRotation:         mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		Curves:             NewBoneCurves(),
	}
}

func (bf *BoneFrame) Add(v *BoneFrame) {
	bf.Position.Add(v.Position)
	bf.MorphPosition.Add(v.MorphPosition)
	bf.LocalPosition.Add(v.LocalPosition)
	bf.MorphLocalPosition.Add(v.MorphLocalPosition)
	bf.Rotation.Mul(v.Rotation)
	bf.MorphRotation.Mul(v.MorphRotation)
	bf.LocalRotation.Mul(v.LocalRotation)
	bf.MorphLocalRotation.Mul(v.MorphLocalRotation)
	bf.Scale.Add(v.Scale)
	bf.MorphScale.Add(v.MorphScale)
	bf.LocalScale.Add(v.LocalScale)
	bf.MorphLocalScale.Add(v.MorphLocalScale)
	bf.IkRotation.Mul(v.IkRotation)
}

func (bf *BoneFrame) Added(v *BoneFrame) *BoneFrame {
	copied := bf.Copy().(*BoneFrame)

	copied.Position.Add(v.Position)
	copied.MorphPosition.Add(v.MorphPosition)
	copied.LocalPosition.Add(v.LocalPosition)
	copied.MorphLocalPosition.Add(v.MorphLocalPosition)
	copied.Rotation.Mul(v.Rotation)
	copied.MorphRotation.Mul(v.MorphRotation)
	copied.LocalRotation.Mul(v.LocalRotation)
	copied.MorphLocalRotation.Mul(v.MorphLocalRotation)
	copied.Scale.Add(v.Scale)
	copied.MorphScale.Add(v.MorphScale)
	copied.LocalScale.Add(v.LocalScale)
	copied.MorphLocalScale.Add(v.MorphLocalScale)
	copied.IkRotation.Mul(v.IkRotation)

	return copied
}
