package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type BoneFrame struct {
	*BaseFrame                            // キーフレ
	Position           *mmath.MVec3       // 位置
	CancelablePosition *mmath.MVec3       // 親キャンセル位置
	Rotation           *mmath.MQuaternion // 回転
	CancelableRotation *mmath.MQuaternion // 親キャンセル回転
	Scale              *mmath.MVec3       // スケール
	CancelableScale    *mmath.MVec3       // 親キャンセルスケール
	Curves             *BoneCurves        // 補間曲線
}

func NewBoneFrame(index float32) *BoneFrame {
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
	if bf.CancelablePosition != nil || v.CancelablePosition != nil {
		if bf.CancelablePosition == nil {
			bf.CancelablePosition = v.CancelablePosition.Copy()
		} else if v.CancelablePosition != nil {
			bf.CancelablePosition.Add(v.CancelablePosition)
		}
	}

	if bf.Rotation != nil || v.Rotation != nil {
		if bf.Rotation == nil {
			bf.Rotation = v.Rotation.Copy()
		} else if v.Rotation != nil {
			bf.Rotation.Mul(v.Rotation)
		}
	}

	if bf.CancelableRotation != nil || v.CancelableRotation != nil {
		if bf.CancelableRotation == nil {
			bf.CancelableRotation = v.CancelableRotation.Copy()
		} else if v.CancelableRotation != nil {
			bf.CancelableRotation.Mul(v.CancelableRotation)
		}
	}

	if bf.Scale != nil || v.Scale != nil {
		if bf.Scale == nil {
			bf.Scale = v.Scale.Copy()
		} else if v.Scale != nil {
			bf.Scale.Add(v.Scale)
		}
	}

	if bf.CancelableScale != nil || v.CancelableScale != nil {
		if bf.CancelableScale == nil {
			bf.CancelableScale = v.CancelableScale.Copy()
		} else if v.CancelableScale != nil {
			bf.CancelableScale.Add(v.CancelableScale)
		}
	}

	return bf
}

func (bf *BoneFrame) Added(v *BoneFrame) *BoneFrame {
	copied := bf.Copy().(*BoneFrame)
	return copied.Add(v)
}

func (bf *BoneFrame) Copy() IBaseFrame {
	copied := &BoneFrame{
		BaseFrame: NewFrame(bf.Index()).(*BaseFrame),
	}
	if bf.Position != nil {
		copied.Position = bf.Position.Copy()
	}
	if bf.Rotation != nil {
		copied.Rotation = bf.Rotation.Copy()
	}
	if bf.CancelablePosition != nil {
		copied.CancelablePosition = bf.CancelablePosition.Copy()
	}
	if bf.CancelableRotation != nil {
		copied.CancelableRotation = bf.CancelableRotation.Copy()
	}
	if bf.Scale != nil {
		copied.Scale = bf.Scale.Copy()
	}
	if bf.CancelableScale != nil {
		copied.CancelableScale = bf.CancelableScale.Copy()
	}
	if bf.Curves != nil {
		copied.Curves = bf.Curves.Copy()
	}

	return copied
}

func (nextBf *BoneFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	prevBf := prevFrame.(*BoneFrame)

	if prevBf == nil || nextBf.Index() <= index {
		// 前がないか、最後より後の場合、次のキーフレをコピーして返す
		return nextBf.Copy().(*BoneFrame)
	}

	bf := NewBoneFrame(index)
	var xy, yy, zy, ry float64
	if nextBf.Curves == nil {
		t := float64(index-prevBf.Index()) / float64(nextBf.Index()-prevBf.Index())
		xy = t
		yy = t
		zy = t
		ry = t
	} else {
		xy, yy, zy, ry = nextBf.Curves.Evaluate(prevBf.Index(), index, nextBf.Index())
	}

	bf.Rotation = prevBf.Rotation.Slerp(nextBf.Rotation, ry)

	ppx := 0.0
	ppy := 0.0
	ppz := 0.0
	if prevBf.Position != nil {
		ppx = prevBf.Position.X
		ppy = prevBf.Position.Y
		ppz = prevBf.Position.Z
	}
	npx := 0.0
	npy := 0.0
	npz := 0.0
	if nextBf.Position != nil {
		npx = nextBf.Position.X
		npy = nextBf.Position.Y
		npz = nextBf.Position.Z
	}

	plpx := 0.0
	plpy := 0.0
	plpz := 0.0
	if prevBf.CancelablePosition != nil {
		plpx = prevBf.CancelablePosition.X
		plpy = prevBf.CancelablePosition.Y
		plpz = prevBf.CancelablePosition.Z
	}
	psx := 1.0
	psy := 1.0
	psz := 1.0
	if prevBf.Scale != nil {
		psx = prevBf.Scale.X
		psy = prevBf.Scale.Y
		psz = prevBf.Scale.Z
	}
	plsx := 1.0
	plsy := 1.0
	plsz := 1.0
	if prevBf.CancelableScale != nil {
		plsx = prevBf.CancelableScale.X
		plsy = prevBf.CancelableScale.Y
		plsz = prevBf.CancelableScale.Z
	}
	nlpx := 0.0
	nlpy := 0.0
	nlpz := 0.0
	if nextBf.CancelablePosition != nil {
		nlpx = nextBf.CancelablePosition.X
		nlpy = nextBf.CancelablePosition.Y
		nlpz = nextBf.CancelablePosition.Z
	}
	nsx := 1.0
	nsy := 1.0
	nsz := 1.0
	if nextBf.Scale != nil {
		nsx = nextBf.Scale.X
		nsy = nextBf.Scale.Y
		nsz = nextBf.Scale.Z
	}
	nlsx := 1.0
	nlsy := 1.0
	nlsz := 1.0
	if nextBf.CancelableScale != nil {
		nlsx = nextBf.CancelableScale.X
		nlsy = nextBf.CancelableScale.Y
		nlsz = nextBf.CancelableScale.Z
	}

	prevX := &mmath.MVec4{X: ppx, Y: plpx, Z: psx, W: plsx}
	nextX := &mmath.MVec4{X: npx, Y: nlpx, Z: nsx, W: nlsx}
	nowX := prevX.Lerp(nextX, xy)

	prevY := &mmath.MVec4{X: ppy, Y: plpy, Z: psy, W: plsy}
	nextY := &mmath.MVec4{X: npy, Y: nlpy, Z: nsy, W: nlsy}
	nowY := prevY.Lerp(nextY, yy)

	prevZ := &mmath.MVec4{X: ppz, Y: plpz, Z: psz, W: plsz}
	nextZ := &mmath.MVec4{X: npz, Y: nlpz, Z: nsz, W: nlsz}
	nowZ := prevZ.Lerp(nextZ, zy)

	bf.Position = &mmath.MVec3{X: nowX.X, Y: nowY.X, Z: nowZ.X}
	bf.CancelablePosition = &mmath.MVec3{X: nowX.Y, Y: nowY.Y, Z: nowZ.Y}
	bf.Scale = &mmath.MVec3{X: nowX.Z, Y: nowY.Z, Z: nowZ.Z}
	bf.CancelableScale = &mmath.MVec3{X: nowX.W, Y: nowY.W, Z: nowZ.W}

	return bf
}

func (bf *BoneFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
	nextBf := nextFrame.(*BoneFrame)
	prevBf := prevFrame.(*BoneFrame)
	if nextBf.Curves == nil {
		// 次の補間曲線がない場合は線形補間
		bf.Curves = &BoneCurves{
			TranslateX: mmath.LINER_CURVE.Copy(),
			TranslateY: mmath.LINER_CURVE.Copy(),
			TranslateZ: mmath.LINER_CURVE.Copy(),
			Rotate:     mmath.LINER_CURVE.Copy(),
		}
		return
	}
	if bf.Curves == nil {
		bf.Curves = NewBoneCurves()
	}

	bf.Curves.TranslateX, nextBf.Curves.TranslateX =
		mmath.SplitCurve(nextBf.Curves.TranslateX, prevBf.Index(), bf.Index(), nextBf.Index())
	bf.Curves.TranslateY, nextBf.Curves.TranslateY =
		mmath.SplitCurve(nextBf.Curves.TranslateY, prevBf.Index(), bf.Index(), nextBf.Index())
	bf.Curves.TranslateZ, nextBf.Curves.TranslateZ =
		mmath.SplitCurve(nextBf.Curves.TranslateZ, prevBf.Index(), bf.Index(), nextBf.Index())
	bf.Curves.Rotate, nextBf.Curves.Rotate =
		mmath.SplitCurve(nextBf.Curves.Rotate, prevBf.Index(), bf.Index(), nextBf.Index())
}
