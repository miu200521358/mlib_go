package vmd

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

var InitialBoneCurves = []byte{
	20,
	20,
	0,
	0,
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
	0,
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
	0,
	0,
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
	0,
	0,
	0,
}

type BoneCurves struct {
	TranslateX *mmath.Curve // 移動X
	TranslateY *mmath.Curve // 移動Y
	TranslateZ *mmath.Curve // 移動Z
	Rotate     *mmath.Curve // 回転
	Values     []byte       // 補間曲線の値
}

func NewBoneCurves() *BoneCurves {
	return &BoneCurves{
		TranslateX: mmath.NewCurve(),
		TranslateY: mmath.NewCurve(),
		TranslateZ: mmath.NewCurve(),
		Rotate:     mmath.NewCurve(),
		Values: []byte{
			20,
			20,
			0,
			0,
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
			0,
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
			0,
			0,
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
			0,
			0,
			0,
		},
	}
}

func NewBoneCurvesByValues(values []byte) *BoneCurves {
	curves := &BoneCurves{
		TranslateX: mmath.NewCurveByValues(values[0], values[4], values[8], values[12]),
		TranslateY: mmath.NewCurveByValues(values[16], values[20], values[24], values[28]),
		TranslateZ: mmath.NewCurveByValues(values[32], values[36], values[40], values[44]),
		Rotate:     mmath.NewCurveByValues(values[48], values[52], values[56], values[60]),
		Values:     values,
	}

	return curves
}

// 補間曲線の計算 (xy, yy, zy, ry)
func (boneCurves *BoneCurves) Evaluate(prevIndex, nowIndex, nextIndex float32) (float64, float64, float64, float64) {
	var xy, yy, zy, ry float64
	_, xy, _ = mmath.Evaluate(boneCurves.TranslateX, prevIndex, nowIndex, nextIndex)
	_, yy, _ = mmath.Evaluate(boneCurves.TranslateY, prevIndex, nowIndex, nextIndex)
	_, zy, _ = mmath.Evaluate(boneCurves.TranslateZ, prevIndex, nowIndex, nextIndex)
	_, ry, _ = mmath.Evaluate(boneCurves.Rotate, prevIndex, nowIndex, nextIndex)

	return xy, yy, zy, ry
}

func (boneCurves *BoneCurves) Merge(enablePhysics bool) []byte {
	c02 := byte(1)
	c03 := byte(1)
	c31 := byte(1)
	c46 := byte(1)
	c47 := byte(0)
	c61 := byte(1)
	c62 := byte(0)
	c63 := byte(0)
	if boneCurves.Values != nil {
		c31 = boneCurves.Values[31]
		c46 = boneCurves.Values[46]
		c47 = boneCurves.Values[47]
		c61 = boneCurves.Values[61]
		c62 = boneCurves.Values[62]
		c63 = boneCurves.Values[63]
	}

	if enablePhysics {
		c02 = byte(boneCurves.TranslateZ.Start.X)
		c03 = byte(boneCurves.Rotate.Start.X)
	} else {
		c02 = byte(99)
		c03 = byte(15)
	}

	return []byte{
		byte(boneCurves.TranslateX.Start.X),
		byte(boneCurves.TranslateY.Start.X),
		c02,
		c03,
		byte(boneCurves.TranslateX.Start.Y),
		byte(boneCurves.TranslateY.Start.Y),
		byte(boneCurves.TranslateZ.Start.Y),
		byte(boneCurves.Rotate.Start.Y),
		byte(boneCurves.TranslateX.End.X),
		byte(boneCurves.TranslateY.End.X),
		byte(boneCurves.TranslateZ.End.X),
		byte(boneCurves.Rotate.End.X),
		byte(boneCurves.TranslateX.End.Y),
		byte(boneCurves.TranslateY.End.Y),
		byte(boneCurves.TranslateZ.End.Y),
		byte(boneCurves.Rotate.End.Y),
		byte(boneCurves.TranslateY.Start.X),
		byte(boneCurves.TranslateZ.Start.X),
		byte(boneCurves.Rotate.Start.X),
		byte(boneCurves.TranslateX.Start.Y),
		byte(boneCurves.TranslateY.Start.Y),
		byte(boneCurves.TranslateZ.Start.Y),
		byte(boneCurves.Rotate.Start.Y),
		byte(boneCurves.TranslateX.End.X),
		byte(boneCurves.TranslateY.End.X),
		byte(boneCurves.TranslateZ.End.X),
		byte(boneCurves.Rotate.End.X),
		byte(boneCurves.TranslateX.End.Y),
		byte(boneCurves.TranslateY.End.Y),
		byte(boneCurves.TranslateZ.End.Y),
		byte(boneCurves.Rotate.End.Y),
		c31,
		byte(boneCurves.TranslateZ.Start.X),
		byte(boneCurves.Rotate.Start.X),
		byte(boneCurves.TranslateX.Start.Y),
		byte(boneCurves.TranslateY.Start.Y),
		byte(boneCurves.TranslateZ.Start.Y),
		byte(boneCurves.Rotate.Start.Y),
		byte(boneCurves.TranslateX.End.X),
		byte(boneCurves.TranslateY.End.X),
		byte(boneCurves.TranslateZ.End.X),
		byte(boneCurves.Rotate.End.X),
		byte(boneCurves.TranslateX.End.Y),
		byte(boneCurves.TranslateY.End.Y),
		byte(boneCurves.TranslateZ.End.Y),
		byte(boneCurves.Rotate.End.Y),
		c46,
		c47,
		byte(boneCurves.Rotate.Start.X),
		byte(boneCurves.TranslateX.Start.Y),
		byte(boneCurves.TranslateY.Start.Y),
		byte(boneCurves.TranslateZ.Start.Y),
		byte(boneCurves.Rotate.Start.Y),
		byte(boneCurves.TranslateX.End.X),
		byte(boneCurves.TranslateY.End.X),
		byte(boneCurves.TranslateZ.End.X),
		byte(boneCurves.Rotate.End.X),
		byte(boneCurves.TranslateX.End.Y),
		byte(boneCurves.TranslateY.End.Y),
		byte(boneCurves.TranslateZ.End.Y),
		byte(boneCurves.Rotate.End.Y),
		c61,
		c62,
		c63,
	}
}

func (boneCurves *BoneCurves) Copy() *BoneCurves {
	return &BoneCurves{
		TranslateX: boneCurves.TranslateX.Copy(),
		TranslateY: boneCurves.TranslateY.Copy(),
		TranslateZ: boneCurves.TranslateZ.Copy(),
		Rotate:     boneCurves.Rotate.Copy(),
		Values:     boneCurves.Values,
	}
}
