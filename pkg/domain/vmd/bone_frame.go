package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
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
	if bf.LocalPosition != nil {
		copied.LocalPosition = bf.LocalPosition.Copy()
	}
	if bf.LocalRotation != nil {
		copied.LocalRotation = bf.LocalRotation.Copy()
	}
	if bf.Scale != nil {
		copied.Scale = bf.Scale.Copy()
	}
	if bf.LocalScale != nil {
		copied.LocalScale = bf.LocalScale.Copy()
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
	if prevBf.LocalPosition != nil {
		plpx = prevBf.LocalPosition.X
		plpy = prevBf.LocalPosition.Y
		plpz = prevBf.LocalPosition.Z
	}
	psx := 0.0
	psy := 0.0
	psz := 0.0
	if prevBf.Scale != nil {
		psx = prevBf.Scale.X
		psy = prevBf.Scale.Y
		psz = prevBf.Scale.Z
	}
	plsx := 0.0
	plsy := 0.0
	plsz := 0.0
	if prevBf.LocalScale != nil {
		plsx = prevBf.LocalScale.X
		plsy = prevBf.LocalScale.Y
		plsz = prevBf.LocalScale.Z
	}
	nlpx := 0.0
	nlpy := 0.0
	nlpz := 0.0
	if nextBf.LocalPosition != nil {
		nlpx = nextBf.LocalPosition.X
		nlpy = nextBf.LocalPosition.Y
		nlpz = nextBf.LocalPosition.Z
	}
	nsx := 0.0
	nsy := 0.0
	nsz := 0.0
	if nextBf.Scale != nil {
		nsx = nextBf.Scale.X
		nsy = nextBf.Scale.Y
		nsz = nextBf.Scale.Z
	}
	nlsx := 0.0
	nlsy := 0.0
	nlsz := 0.0
	if nextBf.LocalScale != nil {
		nlsx = nextBf.LocalScale.X
		nlsy = nextBf.LocalScale.Y
		nlsz = nextBf.LocalScale.Z
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
	bf.LocalPosition = &mmath.MVec3{X: nowX.Y, Y: nowY.Y, Z: nowZ.Y}
	bf.Scale = &mmath.MVec3{X: nowX.Z, Y: nowY.Z, Z: nowZ.Z}
	bf.LocalScale = &mmath.MVec3{X: nowX.W, Y: nowY.W, Z: nowZ.W}

	return bf
}

func (bf *BoneFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
	nextBf := nextFrame.(*BoneFrame)
	prevBf := prevFrame.(*BoneFrame)

	bf.Curves.TranslateX, nextBf.Curves.TranslateX =
		mmath.SplitCurve(nextBf.Curves.TranslateX, prevBf.Index(), bf.Index(), nextBf.Index())
	bf.Curves.TranslateY, nextBf.Curves.TranslateY =
		mmath.SplitCurve(nextBf.Curves.TranslateY, prevBf.Index(), bf.Index(), nextBf.Index())
	bf.Curves.TranslateZ, nextBf.Curves.TranslateZ =
		mmath.SplitCurve(nextBf.Curves.TranslateZ, prevBf.Index(), bf.Index(), nextBf.Index())
	bf.Curves.Rotate, nextBf.Curves.Rotate =
		mmath.SplitCurve(nextBf.Curves.Rotate, prevBf.Index(), bf.Index(), nextBf.Index())
}
