package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type BoneCurves struct {
	TranslateX *mmath.Curve // 移動X
	TranslateY *mmath.Curve // 移動Y
	TranslateZ *mmath.Curve // 移動Z
	Rotate     *mmath.Curve // 回転
	values     []byte       // 補間曲線の値
}

func NewBoneCurves() *BoneCurves {
	return &BoneCurves{
		TranslateX: mmath.NewCurve(),
		TranslateY: mmath.NewCurve(),
		TranslateZ: mmath.NewCurve(),
		Rotate:     mmath.NewCurve(),
		values: []byte{
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
		values:     values,
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
func (v *BoneCurves) Evaluate(prevIndex int, nowIndex int, nextIndex int) (float64, float64, float64, float64) {
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
		c.values[1],
		c.values[2],
		c.values[3],
		byte(c.TranslateX.Start.GetY()),
		c.values[5],
		c.values[6],
		c.values[7],
		byte(c.TranslateX.End.GetX()),
		c.values[9],
		c.values[10],
		c.values[11],
		byte(c.TranslateX.End.GetY()),
		c.values[13],
		c.values[14],
		c.values[15],
		byte(c.TranslateY.Start.GetX()),
		c.values[17],
		c.values[18],
		c.values[19],
		byte(c.TranslateY.Start.GetY()),
		c.values[21],
		c.values[22],
		c.values[23],
		byte(c.TranslateY.End.GetX()),
		c.values[25],
		c.values[26],
		c.values[27],
		byte(c.TranslateY.End.GetY()),
		c.values[29],
		c.values[30],
		c.values[31],
		byte(c.TranslateZ.Start.GetX()),
		c.values[33],
		c.values[34],
		c.values[35],
		byte(c.TranslateZ.Start.GetY()),
		c.values[37],
		c.values[38],
		c.values[39],
		byte(c.TranslateZ.End.GetX()),
		c.values[41],
		c.values[42],
		c.values[43],
		byte(c.TranslateZ.End.GetY()),
		c.values[45],
		c.values[46],
		c.values[47],
		byte(c.Rotate.Start.GetX()),
		c.values[49],
		c.values[50],
		c.values[51],
		byte(c.Rotate.Start.GetY()),
		c.values[53],
		c.values[54],
		c.values[55],
		byte(c.Rotate.End.GetX()),
		c.values[57],
		c.values[58],
		c.values[59],
		byte(c.Rotate.End.GetY()),
		c.values[61],
		c.values[62],
		c.values[63],
	}
}

func (c *BoneCurves) Copy() *BoneCurves {
	return &BoneCurves{
		TranslateX: c.TranslateX.Copy(),
		TranslateY: c.TranslateY.Copy(),
		TranslateZ: c.TranslateZ.Copy(),
		Rotate:     c.Rotate.Copy(),
		values:     c.values,
	}
}
