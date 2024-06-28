package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type BoneFrame struct {
	*BaseFrame                       // キーフレ
	Position      *mmath.MVec3       // 位置
	LocalPosition *mmath.MVec3       // ローカル位置
	Rotation      *mmath.MQuaternion // 回転
	LocalRotation *mmath.MQuaternion // ローカル回転
	Scale         *mmath.MVec3       // スケール
	LocalScale    *mmath.MVec3       // ローカルスケール
	Curves        *BoneCurves        // 補間曲線
}

func NewBoneFrame(index int) *BoneFrame {
	return &BoneFrame{
		BaseFrame: NewFrame(index).(*BaseFrame),
	}
}

func NullBoneFrame() *BoneFrame {
	return nil
}

func (bf *BoneFrame) Add(v *BoneFrame) *BoneFrame {
	if bf.Position != nil || v.Position != nil {
		if bf.Position == nil {
			bf.Position = v.Position.Copy()
		} else if v.Position != nil {
			bf.Position.Add(v.Position)
		}
	}
	if bf.LocalPosition != nil || v.LocalPosition != nil {
		if bf.LocalPosition == nil {
			bf.LocalPosition = v.LocalPosition.Copy()
		} else if v.LocalPosition != nil {
			bf.LocalPosition.Add(v.LocalPosition)
		}
	}

	if bf.Rotation != nil || v.Rotation != nil {
		if bf.Rotation == nil {
			bf.Rotation = v.Rotation.Copy()
		} else if v.Rotation != nil {
			bf.Rotation.Mul(v.Rotation)
		}
	}

	if bf.LocalRotation != nil || v.LocalRotation != nil {
		if bf.LocalRotation == nil {
			bf.LocalRotation = v.LocalRotation.Copy()
		} else if v.LocalRotation != nil {
			bf.LocalRotation.Mul(v.LocalRotation)
		}
	}

	if bf.Scale != nil || v.Scale != nil {
		if bf.Scale == nil {
			bf.Scale = v.Scale.Copy()
		} else if v.Scale != nil {
			bf.Scale.Add(v.Scale)
		}
	}

	if bf.LocalScale != nil || v.LocalScale != nil {
		if bf.LocalScale == nil {
			bf.LocalScale = v.LocalScale.Copy()
		} else if v.LocalScale != nil {
			bf.LocalScale.Add(v.LocalScale)
		}
	}

	return bf
}

func (bf *BoneFrame) Added(v *BoneFrame) *BoneFrame {
	copied := bf.Copy().(*BoneFrame)
	return copied.Add(v)
}

func (v *BoneFrame) Copy() IBaseFrame {
	copied := &BoneFrame{
		BaseFrame:     NewFrame(v.GetIndex()).(*BaseFrame),
		Position:      v.Position.Copy(),
		LocalPosition: nil,
		Rotation:      v.Rotation.Copy(),
		LocalRotation: nil,
		Scale:         nil,
		LocalScale:    nil,
		Curves:        NewBoneCurves(),
	}
	if v.LocalPosition != nil {
		copied.LocalPosition = v.LocalPosition.Copy()
	}
	if v.LocalRotation != nil {
		copied.LocalRotation = v.LocalRotation.Copy()
	}
	if v.Scale != nil {
		copied.Scale = v.Scale.Copy()
	}
	if v.LocalScale != nil {
		copied.LocalScale = v.LocalScale.Copy()
	}
	copied.Curves = v.Curves.Copy()

	return copied
}

func (nextBf *BoneFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevBf := prevFrame.(*BoneFrame)

	if prevBf == nil || nextBf.GetIndex() <= index {
		// 前がないか、最後より後の場合、次のキーフレをコピーして返す
		return nextBf.Copy().(*BoneFrame)
	}

	bf := NewBoneFrame(index)
	var xy, yy, zy, ry float64
	if nextBf.Curves == nil {
		t := float64(index-prevBf.GetIndex()) / float64(nextBf.GetIndex()-prevBf.GetIndex())
		xy = t
		yy = t
		zy = t
		ry = t
	} else {
		xy, yy, zy, ry = nextBf.Curves.Evaluate(prevBf.GetIndex(), index, nextBf.GetIndex())
	}

	bf.Rotation = prevBf.Rotation.Slerp(nextBf.Rotation, ry)

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
	nowX := prevX.Lerp(nextX, xy)

	prevY := &mmath.MVec4{
		prevBf.Position.GetY(), plpy, psy, plsy}
	nextY := &mmath.MVec4{
		nextBf.Position.GetY(), nlpy, nsy, nlsy}
	nowY := prevY.Lerp(nextY, yy)

	prevZ := &mmath.MVec4{
		prevBf.Position.GetZ(), plpz, psz, plsz}
	nextZ := &mmath.MVec4{
		nextBf.Position.GetZ(), nlpz, nsz, nlsz}
	nowZ := prevZ.Lerp(nextZ, zy)

	bf.Position = &mmath.MVec3{nowX[0], nowY[0], nowZ[0]}
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
