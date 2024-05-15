package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type CameraCurves struct {
	TranslateX  mmath.Curve // 移動X
	TranslateY  mmath.Curve // 移動Y
	TranslateZ  mmath.Curve // 移動Z
	Rotate      mmath.Curve // 回転
	Distance    mmath.Curve // 距離
	ViewOfAngle mmath.Curve // 視野角
	values      []byte      // 補間曲線の値
}

func NewCameraCurves() *CameraCurves {
	return &CameraCurves{
		TranslateX:  mmath.NewCurve(),
		TranslateY:  mmath.NewCurve(),
		TranslateZ:  mmath.NewCurve(),
		Rotate:      mmath.NewCurve(),
		Distance:    mmath.NewCurve(),
		ViewOfAngle: mmath.NewCurve(),
		values: []byte{
			20,
			20,
			20,
			20,
			20,
			20,
			20,
			20,
			20,
			20,
			20,
			20,
			107,
			107,
			107,
			107,
			107,
			107,
			107,
			107,
			107,
			107,
			107,
			107,
		},
	}
}

func NewCameraCurvesByValues(values []byte) *CameraCurves {
	curves := NewCameraCurves()
	curves.values = values
	curves.TranslateX.Start.SetX(float64(values[0]))
	curves.TranslateY.Start.SetX(float64(values[1]))
	curves.TranslateZ.Start.SetX(float64(values[2]))
	curves.Rotate.Start.SetX(float64(values[3]))
	curves.Distance.Start.SetX(float64(values[4]))
	curves.ViewOfAngle.Start.SetX(float64(values[5]))
	curves.TranslateX.Start.SetY(float64(values[6]))
	curves.TranslateY.Start.SetY(float64(values[7]))
	curves.TranslateZ.Start.SetY(float64(values[8]))
	curves.Rotate.Start.SetY(float64(values[9]))
	curves.Distance.Start.SetY(float64(values[10]))
	curves.ViewOfAngle.Start.SetY(float64(values[11]))
	curves.TranslateX.End.SetX(float64(values[12]))
	curves.TranslateY.End.SetX(float64(values[13]))
	curves.TranslateZ.End.SetX(float64(values[14]))
	curves.Rotate.End.SetX(float64(values[15]))
	curves.Distance.End.SetX(float64(values[16]))
	curves.ViewOfAngle.End.SetX(float64(values[17]))
	curves.TranslateX.End.SetY(float64(values[18]))
	curves.TranslateY.End.SetY(float64(values[19]))
	curves.TranslateZ.End.SetY(float64(values[20]))
	curves.Rotate.End.SetY(float64(values[21]))
	curves.Distance.End.SetY(float64(values[22]))
	curves.ViewOfAngle.End.SetY(float64(values[23]))
	return curves
}

// 補間曲線の計算
func (v *CameraCurves) Evaluate(prevIndex, nowIndex, nextIndex int) (float64, float64, float64, float64, float64, float64) {
	var xy, yy, zy, ry, dy, vy float64
	_, xy, _ = mmath.Evaluate(v.TranslateX, prevIndex, nowIndex, nextIndex)
	_, yy, _ = mmath.Evaluate(v.TranslateY, prevIndex, nowIndex, nextIndex)
	_, zy, _ = mmath.Evaluate(v.TranslateZ, prevIndex, nowIndex, nextIndex)
	_, ry, _ = mmath.Evaluate(v.Rotate, prevIndex, nowIndex, nextIndex)
	_, dy, _ = mmath.Evaluate(v.Distance, prevIndex, nowIndex, nextIndex)
	_, vy, _ = mmath.Evaluate(v.ViewOfAngle, prevIndex, nowIndex, nextIndex)

	return xy, yy, zy, ry, dy, vy
}

func (c *CameraCurves) Merge() []byte {
	return []byte{
		byte(c.TranslateX.Start.GetX()),
		byte(c.TranslateY.Start.GetX()),
		byte(c.TranslateZ.Start.GetX()),
		byte(c.Rotate.Start.GetX()),
		byte(c.Distance.Start.GetX()),
		byte(c.ViewOfAngle.Start.GetX()),
		byte(c.TranslateX.Start.GetY()),
		byte(c.TranslateY.Start.GetY()),
		byte(c.TranslateZ.Start.GetY()),
		byte(c.Rotate.Start.GetY()),
		byte(c.Distance.Start.GetY()),
		byte(c.ViewOfAngle.Start.GetY()),
		byte(c.TranslateX.End.GetX()),
		byte(c.TranslateY.End.GetX()),
		byte(c.TranslateZ.End.GetX()),
		byte(c.Rotate.End.GetX()),
		byte(c.Distance.End.GetX()),
		byte(c.ViewOfAngle.End.GetX()),
		byte(c.TranslateX.End.GetY()),
		byte(c.TranslateY.End.GetY()),
		byte(c.TranslateZ.End.GetY()),
		byte(c.Rotate.End.GetY()),
		byte(c.Distance.End.GetY()),
		byte(c.ViewOfAngle.End.GetY()),
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
