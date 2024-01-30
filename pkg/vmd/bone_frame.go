package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneFrame struct {
	*BaseFrame                     // キーフレ
	Position      *mmath.MVec3     // 位置
	LocalPosition *mmath.MVec3     // ローカル位置
	Rotation      *mmath.MRotation // 回転
	LocalRotation *mmath.MRotation // ローカル回転
	Scale         *mmath.MVec3     // スケール
	LocalScale    *mmath.MVec3     // ローカルスケール
	IkRotation    *mmath.MRotation // IK回転
	Curves        *BoneCurves      // 補間曲線
	IkRegistered  bool             // IK計算済み
}

func NewBoneFrame(index int) *BoneFrame {
	return &BoneFrame{
		BaseFrame:     NewVmdBaseFrame(index),
		Position:      mmath.NewMVec3(),
		LocalPosition: mmath.NewMVec3(),
		Rotation:      mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		LocalRotation: mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		Scale:         mmath.NewMVec3(),
		LocalScale:    mmath.NewMVec3(),
		IkRotation:    mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		Curves:        NewBoneCurves(),
	}
}

func (bf *BoneFrame) Add(v *BoneFrame) {
	bf.Position.Add(v.Position)
	bf.LocalPosition.Add(v.LocalPosition)
	bf.Rotation.Mul(v.Rotation)
	bf.LocalRotation.Mul(v.LocalRotation)
	bf.Scale.Add(v.Scale)
	bf.LocalScale.Add(v.LocalScale)
	bf.IkRotation.Mul(v.IkRotation)
}

func (bf *BoneFrame) Added(v *BoneFrame) *BoneFrame {
	copied := bf.Copy().(*BoneFrame)

	copied.Position.Add(v.Position)
	copied.LocalPosition.Add(v.LocalPosition)
	copied.Rotation.Mul(v.Rotation)
	copied.LocalRotation.Mul(v.LocalRotation)
	copied.Scale.Add(v.Scale)
	copied.LocalScale.Add(v.LocalScale)
	copied.IkRotation.Mul(v.IkRotation)

	return copied
}
