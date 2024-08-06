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
func (cameraCurves *CameraCurves) Evaluate(prevIndex, nowIndex, nextIndex float64) (float64, float64, float64, float64, float64, float64) {
	var xy, yy, zy, ry, dy, vy float64
	_, xy, _ = mmath.Evaluate(cameraCurves.TranslateX, prevIndex, nowIndex, nextIndex)
	_, yy, _ = mmath.Evaluate(cameraCurves.TranslateY, prevIndex, nowIndex, nextIndex)
	_, zy, _ = mmath.Evaluate(cameraCurves.TranslateZ, prevIndex, nowIndex, nextIndex)
	_, ry, _ = mmath.Evaluate(cameraCurves.Rotate, prevIndex, nowIndex, nextIndex)
	_, dy, _ = mmath.Evaluate(cameraCurves.Distance, prevIndex, nowIndex, nextIndex)
	_, vy, _ = mmath.Evaluate(cameraCurves.ViewOfAngle, prevIndex, nowIndex, nextIndex)

	return xy, yy, zy, ry, dy, vy
}

func (cameraCurves *CameraCurves) Merge() []byte {
	return []byte{
		byte(cameraCurves.TranslateX.Start.X),
		byte(cameraCurves.TranslateY.Start.X),
		byte(cameraCurves.TranslateZ.Start.X),
		byte(cameraCurves.Rotate.Start.X),
		byte(cameraCurves.Distance.Start.X),
		byte(cameraCurves.ViewOfAngle.Start.X),
		byte(cameraCurves.TranslateX.Start.Y),
		byte(cameraCurves.TranslateY.Start.Y),
		byte(cameraCurves.TranslateZ.Start.Y),
		byte(cameraCurves.Rotate.Start.Y),
		byte(cameraCurves.Distance.Start.Y),
		byte(cameraCurves.ViewOfAngle.Start.Y),
		byte(cameraCurves.TranslateX.End.X),
		byte(cameraCurves.TranslateY.End.X),
		byte(cameraCurves.TranslateZ.End.X),
		byte(cameraCurves.Rotate.End.X),
		byte(cameraCurves.Distance.End.X),
		byte(cameraCurves.ViewOfAngle.End.X),
		byte(cameraCurves.TranslateX.End.Y),
		byte(cameraCurves.TranslateY.End.Y),
		byte(cameraCurves.TranslateZ.End.Y),
		byte(cameraCurves.Rotate.End.Y),
		byte(cameraCurves.Distance.End.Y),
		byte(cameraCurves.ViewOfAngle.End.Y),
	}
}

func (cameraCurves *CameraCurves) Copy() *CameraCurves {
	return &CameraCurves{
		TranslateX:  cameraCurves.TranslateX.Copy(),
		TranslateY:  cameraCurves.TranslateY.Copy(),
		TranslateZ:  cameraCurves.TranslateZ.Copy(),
		Rotate:      cameraCurves.Rotate.Copy(),
		Distance:    cameraCurves.Distance.Copy(),
		ViewOfAngle: cameraCurves.ViewOfAngle.Copy(),
	}
}
