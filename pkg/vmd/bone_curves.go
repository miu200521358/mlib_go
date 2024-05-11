package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type BoneCurves struct {
	TranslateX mmath.Curve // 移動X
	TranslateY mmath.Curve // 移動Y
	TranslateZ mmath.Curve // 移動Z
	Rotate     mmath.Curve // 回転
	Values     []byte      // 補間曲線の値
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
		TranslateX: mmath.NewCurve(),
		TranslateY: mmath.NewCurve(),
		TranslateZ: mmath.NewCurve(),
		Rotate:     mmath.NewCurve(),
		Values:     values,
	}
	curves.TranslateX.Start.SetX(float64(values[0]))
	curves.TranslateX.Start.SetY(float64(values[4]))
	curves.TranslateX.End.SetX(float64(values[8]))
	curves.TranslateX.End.SetY(float64(values[12]))
	curves.TranslateY.Start.SetX(float64(values[16]))
	curves.TranslateY.Start.SetY(float64(values[20]))
	curves.TranslateY.End.SetX(float64(values[24]))
	curves.TranslateY.End.SetY(float64(values[28]))
	curves.TranslateZ.Start.SetX(float64(values[32]))
	curves.TranslateZ.Start.SetY(float64(values[36]))
	curves.TranslateZ.End.SetX(float64(values[40]))
	curves.TranslateZ.End.SetY(float64(values[44]))
	curves.Rotate.Start.SetX(float64(values[48]))
	curves.Rotate.Start.SetY(float64(values[52]))
	curves.Rotate.End.SetX(float64(values[56]))
	curves.Rotate.End.SetY(float64(values[60]))

	return curves
}

// 補間曲線の計算 (xy, yy, zy, ry)
func (v *BoneCurves) Evaluate(prevIndex, nowIndex, nextIndex int) (float64, float64, float64, float64) {
	var xy, yy, zy, ry float64
	_, xy, _ = mmath.Evaluate(v.TranslateX, prevIndex, nowIndex, nextIndex)
	_, yy, _ = mmath.Evaluate(v.TranslateY, prevIndex, nowIndex, nextIndex)
	_, zy, _ = mmath.Evaluate(v.TranslateZ, prevIndex, nowIndex, nextIndex)
	_, ry, _ = mmath.Evaluate(v.Rotate, prevIndex, nowIndex, nextIndex)

	return xy, yy, zy, ry
}

func (c *BoneCurves) Merge() []byte {
	return []byte{
		byte(c.TranslateX.Start.GetX()),
		c.Values[1],
		c.Values[2],
		c.Values[3],
		byte(c.TranslateX.Start.GetY()),
		c.Values[5],
		c.Values[6],
		c.Values[7],
		byte(c.TranslateX.End.GetX()),
		c.Values[9],
		c.Values[10],
		c.Values[11],
		byte(c.TranslateX.End.GetY()),
		c.Values[13],
		c.Values[14],
		c.Values[15],
		byte(c.TranslateY.Start.GetX()),
		c.Values[17],
		c.Values[18],
		c.Values[19],
		byte(c.TranslateY.Start.GetY()),
		c.Values[21],
		c.Values[22],
		c.Values[23],
		byte(c.TranslateY.End.GetX()),
		c.Values[25],
		c.Values[26],
		c.Values[27],
		byte(c.TranslateY.End.GetY()),
		c.Values[29],
		c.Values[30],
		c.Values[31],
		byte(c.TranslateZ.Start.GetX()),
		c.Values[33],
		c.Values[34],
		c.Values[35],
		byte(c.TranslateZ.Start.GetY()),
		c.Values[37],
		c.Values[38],
		c.Values[39],
		byte(c.TranslateZ.End.GetX()),
		c.Values[41],
		c.Values[42],
		c.Values[43],
		byte(c.TranslateZ.End.GetY()),
		c.Values[45],
		c.Values[46],
		c.Values[47],
		byte(c.Rotate.Start.GetX()),
		c.Values[49],
		c.Values[50],
		c.Values[51],
		byte(c.Rotate.Start.GetY()),
		c.Values[53],
		c.Values[54],
		c.Values[55],
		byte(c.Rotate.End.GetX()),
		c.Values[57],
		c.Values[58],
		c.Values[59],
		byte(c.Rotate.End.GetY()),
		c.Values[61],
		c.Values[62],
		c.Values[63],
	}
}

func (c *BoneCurves) Copy() *BoneCurves {
	return &BoneCurves{
		TranslateX: c.TranslateX.Copy(),
		TranslateY: c.TranslateY.Copy(),
		TranslateZ: c.TranslateZ.Copy(),
		Rotate:     c.Rotate.Copy(),
		Values:     c.Values,
	}
}
