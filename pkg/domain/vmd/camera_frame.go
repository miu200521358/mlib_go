package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

type CameraFrame struct {
	*BaseFrame                          // キーフレ
	Position         *mmath.MVec3       // 位置
	Degrees          *mmath.MVec3       // 回転(オイラー角)
	Quaternion       *mmath.MQuaternion // 回転(クォータニオン)
	Distance         float64            // 距離
	ViewOfAngle      int                // 視野角
	IsPerspectiveOff bool               // パースOFF
	Curves           *CameraCurves      // 補間曲線
}

func NewCameraFrame(index float32) *CameraFrame {
	return &CameraFrame{
		BaseFrame:        NewFrame(index).(*BaseFrame),
		Position:         mmath.NewMVec3(),
		Degrees:          mmath.NewMVec3(),
		Distance:         0.0,
		ViewOfAngle:      0,
		IsPerspectiveOff: true,
		Curves:           NewCameraCurves(),
	}
}

func (cf *CameraFrame) Add(v *CameraFrame) {
	cf.Position.Add(v.Position)
	cf.Degrees.Mul(v.Degrees)
	cf.Distance += v.Distance
	cf.ViewOfAngle += v.ViewOfAngle
}

func (cf *CameraFrame) Added(v *CameraFrame) *CameraFrame {
	copied := cf.Copy().(*CameraFrame)

	copied.Position.Add(v.Position)
	copied.Degrees.Mul(v.Degrees)
	copied.Distance += v.Distance
	copied.ViewOfAngle += v.ViewOfAngle

	return copied
}

func (cf *CameraFrame) Copy() IBaseFrame {
	copied := NewCameraFrame(cf.Index())
	copied.Position = cf.Position
	copied.Degrees = cf.Degrees.Copy()
	copied.Distance = cf.Distance
	copied.ViewOfAngle = cf.ViewOfAngle
	copied.IsPerspectiveOff = cf.IsPerspectiveOff
	copied.Curves = cf.Curves.Copy()

	return copied
}

func (nextCf *CameraFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
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

	q1 := mmath.NewMQuaternionFromRadians(prevCf.Degrees.X, prevCf.Degrees.Y, prevCf.Degrees.Z)
	q2 := mmath.NewMQuaternionFromRadians(nextCf.Degrees.X, nextCf.Degrees.Y, nextCf.Degrees.Z)

	cf.Quaternion = q1.Slerp(q2, ry)

	cf.Position.X = mmath.Lerp(prevCf.Position.X, nextCf.Position.X, xy)
	cf.Position.Y = mmath.Lerp(prevCf.Position.Y, nextCf.Position.Y, yy)
	cf.Position.Z = mmath.Lerp(prevCf.Position.Z, nextCf.Position.Z, zy)

	cf.Distance = mmath.Lerp(prevCf.Distance, nextCf.Distance, dy)
	cf.ViewOfAngle = int(mmath.Lerp(float64(prevCf.ViewOfAngle), float64(nextCf.ViewOfAngle), vy))
	cf.IsPerspectiveOff = nextCf.IsPerspectiveOff

	return cf
}

func (cf *CameraFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
	prevCf := prevFrame.(*CameraFrame)
	nextCf := nextFrame.(*CameraFrame)
	if cf.Curves == nil {
		cf.Curves = NewCameraCurves()
	}

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
