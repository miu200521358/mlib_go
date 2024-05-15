package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type CameraFrame struct {
	*BaseFrame                        // キーフレ
	Position         mmath.MVec3      // 位置
	Rotation         *mmath.MRotation // 回転
	Distance         float64          // 距離
	ViewOfAngle      int              // 視野角
	IsPerspectiveOff bool             // パースOFF
	Curves           *CameraCurves    // 補間曲線
}

func NewCameraFrame(index int) *CameraFrame {
	return &CameraFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		Position:         mmath.NewMVec3(),
		Rotation:         mmath.NewRotationByDegrees(&mmath.MVec3{0, 0, 0}),
		Distance:         0.0,
		ViewOfAngle:      0,
		IsPerspectiveOff: true,
		Curves:           NewCameraCurves(),
	}
}

func (cf *CameraFrame) Add(v *CameraFrame) {
	cf.Position.Add(&v.Position)
	cf.Rotation.Mul(v.Rotation)
	cf.Distance += v.Distance
	cf.ViewOfAngle += v.ViewOfAngle
}

func (cf *CameraFrame) Added(v *CameraFrame) *CameraFrame {
	copied := cf.Copy().(*CameraFrame)

	copied.Position.Add(&v.Position)
	copied.Rotation.Mul(v.Rotation)
	copied.Distance += v.Distance
	copied.ViewOfAngle += v.ViewOfAngle

	return copied
}

func (cf *CameraFrame) Copy() IBaseFrame {
	copied := NewCameraFrame(cf.GetIndex())
	copied.Position = cf.Position
	copied.Rotation = cf.Rotation.Copy()
	copied.Distance = cf.Distance
	copied.ViewOfAngle = cf.ViewOfAngle
	copied.IsPerspectiveOff = cf.IsPerspectiveOff
	copied.Curves = cf.Curves.Copy()

	return copied
}

func (nextCf *CameraFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	prevCf := prevFrame.(*CameraFrame)

	if prevCf == nil || nextCf.GetIndex() <= index {
		// 前がないか、最後より後の場合、次のキーフレをコピーして返す
		frame := nextCf.Copy()
		return frame
	}

	if nextCf == nil {
		frame := prevCf.Copy()
		return frame
	}

	cf := NewCameraFrame(index)

	xy, yy, zy, ry, dy, vy := nextCf.Curves.Evaluate(prevCf.GetIndex(), index, nextCf.GetIndex())

	qq := prevCf.Rotation.GetQuaternion().Slerp(nextCf.Rotation.GetQuaternion(), ry)
	cf.Rotation.SetQuaternion(qq)

	cf.Position.SetX(mmath.LerpFloat(prevCf.Position.GetX(), nextCf.Position.GetX(), xy))
	cf.Position.SetY(mmath.LerpFloat(prevCf.Position.GetY(), nextCf.Position.GetY(), yy))
	cf.Position.SetZ(mmath.LerpFloat(prevCf.Position.GetZ(), nextCf.Position.GetZ(), zy))

	cf.Distance = mmath.LerpFloat(prevCf.Distance, nextCf.Distance, dy)
	cf.ViewOfAngle = int(mmath.LerpFloat(float64(prevCf.ViewOfAngle), float64(nextCf.ViewOfAngle), vy))
	cf.IsPerspectiveOff = nextCf.IsPerspectiveOff

	return cf
}

func (cf *CameraFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
	prevCf := prevFrame.(*CameraFrame)
	nextCf := nextFrame.(*CameraFrame)

	cf.Curves.TranslateX, nextCf.Curves.TranslateX =
		mmath.SplitCurve(nextCf.Curves.TranslateX, prevCf.GetIndex(), index, nextCf.GetIndex())
	cf.Curves.TranslateY, nextCf.Curves.TranslateY =
		mmath.SplitCurve(nextCf.Curves.TranslateY, prevCf.GetIndex(), index, nextCf.GetIndex())
	cf.Curves.TranslateZ, nextCf.Curves.TranslateZ =
		mmath.SplitCurve(nextCf.Curves.TranslateZ, prevCf.GetIndex(), index, nextCf.GetIndex())
	cf.Curves.Rotate, nextCf.Curves.Rotate =
		mmath.SplitCurve(nextCf.Curves.Rotate, prevCf.GetIndex(), index, nextCf.GetIndex())
	cf.Curves.Distance, nextCf.Curves.Distance =
		mmath.SplitCurve(nextCf.Curves.Distance, prevCf.GetIndex(), index, nextCf.GetIndex())
	cf.Curves.ViewOfAngle, nextCf.Curves.ViewOfAngle =
		mmath.SplitCurve(nextCf.Curves.ViewOfAngle, prevCf.GetIndex(), index, nextCf.GetIndex())
}
