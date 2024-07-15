package vmd

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

var InitialCameraCurves = []byte{
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
}

type CameraCurves struct {
	TranslateX  *mmath.Curve // 移動X
	TranslateY  *mmath.Curve // 移動Y
	TranslateZ  *mmath.Curve // 移動Z
	Rotate      *mmath.Curve // 回転
	Distance    *mmath.Curve // 距離
	ViewOfAngle *mmath.Curve // 視野角
	Values      []byte       // 補間曲線の値
}

func NewCameraCurves() *CameraCurves {
	return &CameraCurves{
		TranslateX:  mmath.NewCurve(),
		TranslateY:  mmath.NewCurve(),
		TranslateZ:  mmath.NewCurve(),
		Rotate:      mmath.NewCurve(),
		Distance:    mmath.NewCurve(),
		ViewOfAngle: mmath.NewCurve(),
		Values: []byte{
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
	curves := &CameraCurves{
		TranslateX:  mmath.NewCurveByValues(values[0], values[6], values[12], values[18]),  // 移動X
		TranslateY:  mmath.NewCurveByValues(values[1], values[7], values[13], values[19]),  // 移動Y
		TranslateZ:  mmath.NewCurveByValues(values[2], values[8], values[14], values[20]),  // 移動Z
		Rotate:      mmath.NewCurveByValues(values[3], values[9], values[15], values[21]),  // 回転
		Distance:    mmath.NewCurveByValues(values[4], values[10], values[16], values[22]), // 距離
		ViewOfAngle: mmath.NewCurveByValues(values[5], values[11], values[17], values[23]), // 視野角
		Values:      values,
	}
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
		byte(c.TranslateX.Start.X),
		byte(c.TranslateY.Start.X),
		byte(c.TranslateZ.Start.X),
		byte(c.Rotate.Start.X),
		byte(c.Distance.Start.X),
		byte(c.ViewOfAngle.Start.X),
		byte(c.TranslateX.Start.Y),
		byte(c.TranslateY.Start.Y),
		byte(c.TranslateZ.Start.Y),
		byte(c.Rotate.Start.Y),
		byte(c.Distance.Start.Y),
		byte(c.ViewOfAngle.Start.Y),
		byte(c.TranslateX.End.X),
		byte(c.TranslateY.End.X),
		byte(c.TranslateZ.End.X),
		byte(c.Rotate.End.X),
		byte(c.Distance.End.X),
		byte(c.ViewOfAngle.End.X),
		byte(c.TranslateX.End.Y),
		byte(c.TranslateY.End.Y),
		byte(c.TranslateZ.End.Y),
		byte(c.Rotate.End.Y),
		byte(c.Distance.End.Y),
		byte(c.ViewOfAngle.End.Y),
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
