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

func NewBoneFrame(index int) *BoneFrame {
	position := mmath.NewMVec3()
	morphPosition := mmath.NewMVec3()
	localPosition := mmath.NewMVec3()
	morphLocalPosition := mmath.NewMVec3()
	scale := mmath.NewMVec3()
	morphScale := mmath.NewMVec3()
	localScale := mmath.NewMVec3()
	morphLocalScale := mmath.NewMVec3()

	return &BoneFrame{
		BaseFrame:          NewFrame(index).(*BaseFrame),
		Position:           &position,
		MorphPosition:      &morphPosition,
		LocalPosition:      &localPosition,
		MorphLocalPosition: &morphLocalPosition,
		Rotation:           mmath.NewRotation(),
		MorphRotation:      mmath.NewRotation(),
		LocalRotation:      mmath.NewRotation(),
		MorphLocalRotation: mmath.NewRotation(),
		Scale:              &scale,
		MorphScale:         &morphScale,
		LocalScale:         &localScale,
		MorphLocalScale:    &morphLocalScale,
		IkRotation:         mmath.NewRotation(),
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

func (v *BoneFrame) Copy() IBaseFrame {
	position := v.Position.Copy()
	morphPosition := v.MorphPosition.Copy()
	localPosition := v.LocalPosition.Copy()
	morphLocalPosition := v.MorphLocalPosition.Copy()
	scale := v.Scale.Copy()
	morphScale := v.MorphScale.Copy()
	localScale := v.LocalScale.Copy()
	morphLocalScale := v.MorphLocalScale.Copy()

	copied := BoneFrame{
		BaseFrame:          NewFrame(v.GetIndex()).(*BaseFrame),
		Position:           &position,
		MorphPosition:      &morphPosition,
		LocalPosition:      &localPosition,
		MorphLocalPosition: &morphLocalPosition,
		Rotation:           v.Rotation.Copy(),
		MorphRotation:      v.MorphRotation.Copy(),
		LocalRotation:      v.LocalRotation.Copy(),
		MorphLocalRotation: v.MorphLocalRotation.Copy(),
		Scale:              &scale,
		MorphScale:         &morphScale,
		LocalScale:         &localScale,
		MorphLocalScale:    &morphLocalScale,
		IkRotation:         v.IkRotation.Copy(),
		Curves:             v.Curves.Copy(),
	}
	return &copied
}

func (nextBf *BoneFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevBf := prevFrame.(*BoneFrame)

	if prevBf == nil || nextBf.GetIndex() <= index {
		// 前がないか、最後より後の場合、次のキーフレをコピーして返す
		frame := nextBf.Copy().(*BoneFrame)
		// 計算情報はクリア
		frame.IkRotation = mmath.NewRotation()
		return frame
	}

	bf := NewBoneFrame(index)

	xy, yy, zy, ry := nextBf.Curves.Evaluate(prevBf.GetIndex(), index, nextBf.GetIndex())

	qq := prevBf.Rotation.GetQuaternion().Slerp(nextBf.Rotation.GetQuaternion(), ry)
	bf.Rotation.SetQuaternion(qq)

	prevX := mmath.MVec4{
		prevBf.Position.GetX(), prevBf.LocalPosition.GetX(), prevBf.Scale.GetX(), prevBf.LocalScale.GetX()}
	nextX := mmath.MVec4{
		nextBf.Position.GetX(), nextBf.LocalPosition.GetX(), nextBf.Scale.GetX(), nextBf.LocalScale.GetX()}
	nowX := mmath.LerpVec4(&prevX, &nextX, xy)
	bf.Position.SetX(nowX[0])
	bf.LocalPosition.SetX(nowX[1])
	bf.Scale.SetX(nowX[2])
	bf.LocalScale.SetX(nowX[3])

	prevY := mmath.MVec4{
		prevBf.Position.GetY(), prevBf.LocalPosition.GetY(), prevBf.Scale.GetY(), prevBf.LocalScale.GetY()}
	nextY := mmath.MVec4{
		nextBf.Position.GetY(), nextBf.LocalPosition.GetY(), nextBf.Scale.GetY(), nextBf.LocalScale.GetY()}
	nowY := mmath.LerpVec4(&prevY, &nextY, yy)
	bf.Position.SetY(nowY[0])
	bf.LocalPosition.SetY(nowY[1])
	bf.Scale.SetY(nowY[2])
	bf.LocalScale.SetY(nowY[3])

	prevZ := mmath.MVec4{
		prevBf.Position.GetZ(), prevBf.LocalPosition.GetZ(), prevBf.Scale.GetZ(), prevBf.LocalScale.GetZ()}
	nextZ := mmath.MVec4{
		nextBf.Position.GetZ(), nextBf.LocalPosition.GetZ(), nextBf.Scale.GetZ(), nextBf.LocalScale.GetZ()}
	nowZ := mmath.LerpVec4(&prevZ, &nextZ, zy)
	bf.Position.SetZ(nowZ[0])
	bf.LocalPosition.SetZ(nowZ[1])
	bf.Scale.SetZ(nowZ[2])
	bf.LocalScale.SetZ(nowZ[3])

	return bf
}

func (bf *BoneFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
	nextBf := nextFrame.(*BoneFrame)
	prevBf := prevFrame.(*BoneFrame)

	bf.Curves.TranslateX, nextBf.Curves.TranslateX =
		mmath.SplitCurve(nextBf.Curves.TranslateX, prevBf.GetIndex(), bf.GetIndex(), nextBf.GetIndex())
	bf.Curves.TranslateY, nextBf.Curves.TranslateY =
		mmath.SplitCurve(nextBf.Curves.TranslateY, prevBf.GetIndex(), bf.GetIndex(), nextBf.GetIndex())
	bf.Curves.TranslateZ, nextBf.Curves.TranslateZ =
		mmath.SplitCurve(nextBf.Curves.TranslateZ, prevBf.GetIndex(), bf.GetIndex(), nextBf.GetIndex())
	bf.Curves.Rotate, nextBf.Curves.Rotate =
		mmath.SplitCurve(nextBf.Curves.Rotate, prevBf.GetIndex(), bf.GetIndex(), nextBf.GetIndex())
}
