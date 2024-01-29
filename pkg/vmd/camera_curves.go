package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type CameraCurves struct {
	TranslateX  *mmath.Curve // 移動X
	TranslateY  *mmath.Curve // 移動Y
	TranslateZ  *mmath.Curve // 移動Z
	Rotate      *mmath.Curve // 回転
	Distance    *mmath.Curve // 距離
	ViewOfAngle *mmath.Curve // 視野角
}

func NewCameraCurves() *CameraCurves {
	return &CameraCurves{
		TranslateX:  mmath.NewCurve(),
		TranslateY:  mmath.NewCurve(),
		TranslateZ:  mmath.NewCurve(),
		Rotate:      mmath.NewCurve(),
		Distance:    mmath.NewCurve(),
		ViewOfAngle: mmath.NewCurve(),
	}
}

// 補間曲線の計算
func (v *CameraCurves) Evaluate(prevIndex int, nowIndex int, nextIndex int) (float64, float64, float64, float64, float64, float64) {
	var xy, yy, zy, ry, dy, vy float64
	_, xy, _ = mmath.Evaluate(v.TranslateX, prevIndex, nowIndex, nextIndex)
	_, yy, _ = mmath.Evaluate(v.TranslateY, prevIndex, nowIndex, nextIndex)
	_, zy, _ = mmath.Evaluate(v.TranslateZ, prevIndex, nowIndex, nextIndex)
	_, ry, _ = mmath.Evaluate(v.Rotate, prevIndex, nowIndex, nextIndex)
	_, dy, _ = mmath.Evaluate(v.Distance, prevIndex, nowIndex, nextIndex)
	_, vy, _ = mmath.Evaluate(v.ViewOfAngle, prevIndex, nowIndex, nextIndex)

	return xy, yy, zy, ry, dy, vy
}

func (c *CameraCurves) Merge() []int {
	return []int{
		int(c.TranslateX.Start.GetX()),
		int(c.TranslateY.Start.GetX()),
		int(c.TranslateZ.Start.GetX()),
		int(c.Rotate.Start.GetX()),
		int(c.Distance.Start.GetX()),
		int(c.ViewOfAngle.Start.GetX()),
		int(c.TranslateX.Start.GetY()),
		int(c.TranslateY.Start.GetY()),
		int(c.TranslateZ.Start.GetY()),
		int(c.Rotate.Start.GetY()),
		int(c.Distance.Start.GetY()),
		int(c.ViewOfAngle.Start.GetY()),
		int(c.TranslateX.End.GetX()),
		int(c.TranslateY.End.GetX()),
		int(c.TranslateZ.End.GetX()),
		int(c.Rotate.End.GetX()),
		int(c.Distance.End.GetX()),
		int(c.ViewOfAngle.End.GetX()),
		int(c.TranslateX.End.GetY()),
		int(c.TranslateY.End.GetY()),
		int(c.TranslateZ.End.GetY()),
		int(c.Rotate.End.GetY()),
		int(c.Distance.End.GetY()),
		int(c.ViewOfAngle.End.GetY()),
	}
}

func (c *CameraCurves) Copy() *CameraCurves {
	return &CameraCurves{
		TranslateX:  c.TranslateX.Copy(),
		TranslateY:  c.TranslateY.Copy(),
		TranslateZ:  c.TranslateZ.Copy(),
		Rotate:      c.Rotate.Copy(),
		Distance:    c.Distance.Copy(),
		ViewOfAngle: c.ViewOfAngle.Copy(),
	}
}
