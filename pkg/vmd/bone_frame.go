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
	return &BoneFrame{
		BaseFrame:          NewFrame(index).(*BaseFrame),
		Position:           mmath.NewMVec3(),
		MorphPosition:      mmath.NewMVec3(),
		LocalPosition:      mmath.NewMVec3(),
		MorphLocalPosition: mmath.NewMVec3(),
		Rotation:           mmath.NewRotation(),
		MorphRotation:      mmath.NewRotation(),
		LocalRotation:      mmath.NewRotation(),
		MorphLocalRotation: mmath.NewRotation(),
		Scale:              mmath.NewMVec3(),
		MorphScale:         mmath.NewMVec3(),
		LocalScale:         mmath.NewMVec3(),
		MorphLocalScale:    mmath.NewMVec3(),
		IkRotation:         mmath.NewRotation(),
		Curves:             NewBoneCurves(),
	}
}

func (bf *BoneFrame) Add(v *BoneFrame) {
	if bf.Position != nil && v.Position != nil {
		bf.Position.Add(v.Position)
	}
	if bf.MorphPosition != nil && v.MorphPosition != nil {
		bf.MorphPosition.Add(v.MorphPosition)
	}
	if bf.LocalPosition != nil && v.LocalPosition != nil {
		bf.LocalPosition.Add(v.LocalPosition)
	}
	if bf.MorphLocalPosition != nil && v.MorphLocalPosition != nil {
		bf.MorphLocalPosition.Add(v.MorphLocalPosition)
	}
	bf.Rotation.Mul(v.Rotation)
	bf.MorphRotation.Mul(v.MorphRotation)
	bf.LocalRotation.Mul(v.LocalRotation)
	bf.MorphLocalRotation.Mul(v.MorphLocalRotation)
	if bf.Scale != nil && v.Scale != nil {
		bf.Scale.Add(v.Scale)
	}
	if bf.MorphScale != nil && v.MorphScale != nil {
		bf.MorphScale.Add(v.MorphScale)
	}
	if bf.LocalScale != nil && v.LocalScale != nil {
		bf.LocalScale.Add(v.LocalScale)
	}
	if bf.MorphLocalScale != nil && v.MorphLocalScale != nil {
		bf.MorphLocalScale.Add(v.MorphLocalScale)
	}
	bf.IkRotation.Mul(v.IkRotation)
}

func (bf *BoneFrame) Added(v *BoneFrame) *BoneFrame {
	copied := bf.Copy().(*BoneFrame)

	if bf.Position != nil && v.Position != nil {
		copied.Position.Add(v.Position)
	}
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
	copied := &BoneFrame{
		BaseFrame:          NewFrame(v.GetIndex()).(*BaseFrame),
		Position:           v.Position.Copy(),
		MorphPosition:      nil,
		LocalPosition:      nil,
		MorphLocalPosition: nil,
		Rotation:           v.Rotation.Copy(),
		MorphRotation:      nil,
		LocalRotation:      nil,
		MorphLocalRotation: nil,
		Scale:              nil,
		MorphScale:         nil,
		LocalScale:         nil,
		MorphLocalScale:    nil,
		IkRotation:         nil,
		Curves:             nil,
	}
	if v.MorphPosition != nil {
		copied.MorphPosition = v.MorphPosition.Copy()
	}
	if v.LocalPosition != nil {
		copied.LocalPosition = v.LocalPosition.Copy()
	}
	if v.MorphLocalPosition != nil {
		copied.MorphLocalPosition = v.MorphLocalPosition.Copy()
	}
	if v.MorphRotation != nil {
		copied.MorphRotation = v.MorphRotation.Copy()
	}
	if v.LocalRotation != nil {
		copied.LocalRotation = v.LocalRotation.Copy()
	}
	if v.MorphLocalRotation != nil {
		copied.MorphLocalRotation = v.MorphLocalRotation.Copy()
	}
	if v.Scale != nil {
		copied.Scale = v.Scale.Copy()
	}
	if v.MorphScale != nil {
		copied.MorphScale = v.MorphScale.Copy()
	}
	if v.LocalScale != nil {
		copied.LocalScale = v.LocalScale.Copy()
	}
	if v.MorphLocalScale != nil {
		copied.MorphLocalScale = v.MorphLocalScale.Copy()
	}
	if v.IkRotation != nil {
		copied.IkRotation = v.IkRotation.Copy()
	}
	if v.Curves != nil {
		copied.Curves = v.Curves.Copy()
	}
	return copied
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

	plpx := 0.0
	plpy := 0.0
	plpz := 0.0
	if prevBf.LocalPosition != nil {
		plpx = prevBf.LocalPosition.GetX()
		plpy = prevBf.LocalPosition.GetY()
		plpz = prevBf.LocalPosition.GetZ()
	}
	psx := 0.0
	psy := 0.0
	psz := 0.0
	if prevBf.Scale != nil {
		psx = prevBf.Scale.GetX()
		psy = prevBf.Scale.GetY()
		psz = prevBf.Scale.GetZ()
	}
	plsx := 0.0
	plsy := 0.0
	plsz := 0.0
	if prevBf.LocalScale != nil {
		plsx = prevBf.LocalScale.GetX()
		plsy = prevBf.LocalScale.GetY()
		plsz = prevBf.LocalScale.GetZ()
	}
	nlpx := 0.0
	nlpy := 0.0
	nlpz := 0.0
	if nextBf.LocalPosition != nil {
		nlpx = nextBf.LocalPosition.GetX()
		nlpy = nextBf.LocalPosition.GetY()
		nlpz = nextBf.LocalPosition.GetZ()
	}
	nsx := 0.0
	nsy := 0.0
	nsz := 0.0
	if nextBf.Scale != nil {
		nsx = nextBf.Scale.GetX()
		nsy = nextBf.Scale.GetY()
		nsz = nextBf.Scale.GetZ()
	}
	nlsx := 0.0
	nlsy := 0.0
	nlsz := 0.0
	if nextBf.LocalScale != nil {
		nlsx = nextBf.LocalScale.GetX()
		nlsy = nextBf.LocalScale.GetY()
		nlsz = nextBf.LocalScale.GetZ()
	}

	prevX := &mmath.MVec4{
		prevBf.Position.GetX(), plpx, psx, plsx}
	nextX := &mmath.MVec4{
		nextBf.Position.GetX(), nlpx, nsx, nlsx}
	nowX := mmath.LerpVec4(prevX, nextX, xy)
	bf.Position.SetX(nowX[0])

	prevY := &mmath.MVec4{
		prevBf.Position.GetY(), plpy, psy, plsy}
	nextY := &mmath.MVec4{
		nextBf.Position.GetY(), nlpy, nsy, nlsy}
	nowY := mmath.LerpVec4(prevY, nextY, yy)
	bf.Position.SetY(nowY[0])

	prevZ := &mmath.MVec4{
		prevBf.Position.GetZ(), plpz, psz, plsz}
	nextZ := &mmath.MVec4{
		nextBf.Position.GetZ(), nlpz, nsz, nlsz}
	nowZ := mmath.LerpVec4(prevZ, nextZ, zy)
	bf.Position.SetZ(nowZ[0])

	bf.LocalPosition = &mmath.MVec3{nowX[1], nowY[1], nowZ[1]}
	bf.Scale = &mmath.MVec3{nowX[2], nowY[2], nowZ[2]}
	bf.LocalScale = &mmath.MVec3{nowX[3], nowY[3], nowZ[3]}

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
