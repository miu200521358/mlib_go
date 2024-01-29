package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type CameraFrame struct {
	*BaseFrame                       // キーフレ
	Position      *mmath.MVec3       // 位置
	Rotation      *mmath.MQuaternion // 回転
	Distance      float64            // 距離
	ViewOfAngle   int                // 視野角
	IsPerspective bool               // パースON,OFF
	Curves        *CameraCurves      // 補間曲線
}

func NewCameraFrame(index int) *CameraFrame {
	return &CameraFrame{
		BaseFrame:     NewVmdBaseFrame(index),
		Position:      mmath.NewMVec3(),
		Rotation:      mmath.NewMQuaternion(),
		Distance:      0.0,
		ViewOfAngle:   0,
		IsPerspective: true,
		Curves:        NewCameraCurves(),
	}
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
