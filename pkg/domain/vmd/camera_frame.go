package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type CameraFrame struct {
	*BaseFrame                        // キーフレ
	Position         *mmath.MVec3     // 位置
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
		Rotation:         mmath.NewMRotation(),
		Distance:         0.0,
		ViewOfAngle:      0,
		IsPerspectiveOff: true,
		Curves:           NewCameraCurves(),
	}
}

func NullCameraFrame() *CameraFrame {
	return nil
}

func (cf *CameraFrame) Add(v *CameraFrame) {
	cf.Position.Add(v.Position)
	cf.Rotation.Mul(v.Rotation)
	cf.Distance += v.Distance
	cf.ViewOfAngle += v.ViewOfAngle
}

func (cf *CameraFrame) Added(v *CameraFrame) *CameraFrame {
	copied := cf.Copy().(*CameraFrame)

	copied.Position.Add(v.Position)
	copied.Rotation.Mul(v.Rotation)
	copied.Distance += v.Distance
	copied.ViewOfAngle += v.ViewOfAngle

	return copied
}

func (cf *CameraFrame) Copy() IBaseFrame {
	copied := NewCameraFrame(cf.Index())
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

	if prevCf == nil || nextCf.Index() <= index {
		// 前がないか、最後より後の場合、次のキーフレをコピーして返す
		frame := nextCf.Copy()
		return frame
	}

	if nextCf == nil {
		frame := prevCf.Copy()
		return frame
	}

	cf := NewCameraFrame(index)

	xy, yy, zy, ry, dy, vy := nextCf.Curves.Evaluate(prevCf.Index(), index, nextCf.Index())

	qq := prevCf.Rotation.GetQuaternion().Slerp(nextCf.Rotation.GetQuaternion(), ry)
	cf.Rotation.SetQuaternion(qq)

	cf.Position.X = mmath.LerpFloat(prevCf.Position.X, nextCf.Position.X, xy)
	cf.Position.Y = mmath.LerpFloat(prevCf.Position.Y, nextCf.Position.Y, yy)
	cf.Position.Z = mmath.LerpFloat(prevCf.Position.Z, nextCf.Position.Z, zy)

	cf.Distance = mmath.LerpFloat(prevCf.Distance, nextCf.Distance, dy)
	cf.ViewOfAngle = int(mmath.LerpFloat(float64(prevCf.ViewOfAngle), float64(nextCf.ViewOfAngle), vy))
	cf.IsPerspectiveOff = nextCf.IsPerspectiveOff

	return cf
}

func (cf *CameraFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
	prevCf := prevFrame.(*CameraFrame)
	nextCf := nextFrame.(*CameraFrame)

	cf.Curves.TranslateX, nextCf.Curves.TranslateX =
		mmath.SplitCurve(nextCf.Curves.TranslateX, prevCf.Index(), index, nextCf.Index())
	cf.Curves.TranslateY, nextCf.Curves.TranslateY =
		mmath.SplitCurve(nextCf.Curves.TranslateY, prevCf.Index(), index, nextCf.Index())
	cf.Curves.TranslateZ, nextCf.Curves.TranslateZ =
		mmath.SplitCurve(nextCf.Curves.TranslateZ, prevCf.Index(), index, nextCf.Index())
	cf.Curves.Rotate, nextCf.Curves.Rotate =
		mmath.SplitCurve(nextCf.Curves.Rotate, prevCf.Index(), index, nextCf.Index())
	cf.Curves.Distance, nextCf.Curves.Distance =
		mmath.SplitCurve(nextCf.Curves.Distance, prevCf.Index(), index, nextCf.Index())
	cf.Curves.ViewOfAngle, nextCf.Curves.ViewOfAngle =
		mmath.SplitCurve(nextCf.Curves.ViewOfAngle, prevCf.Index(), index, nextCf.Index())
}
