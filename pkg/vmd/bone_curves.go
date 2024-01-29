package vmd

import "github.com/miu200521358/mlib_go/pkg/mmath"

type BoneCurves struct {
	TranslateX *mmath.Curve // 移動X
	TranslateY *mmath.Curve // 移動Y
	TranslateZ *mmath.Curve // 移動Z
	Rotate     *mmath.Curve // 回転
	values     [64]int      // 補間曲線の値
}

func NewBoneCurves() *BoneCurves {
	return &BoneCurves{
		TranslateX: mmath.NewCurve(),
		TranslateY: mmath.NewCurve(),
		TranslateZ: mmath.NewCurve(),
		Rotate:     mmath.NewCurve(),
		values: [64]int{
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

// 補間曲線の計算 (xy, yy, zy, ry)
func (v *BoneCurves) Evaluate(prevIndex int, nowIndex int, nextIndex int) (float64, float64, float64, float64) {
	var xy, yy, zy, ry float64
	_, xy, _ = mmath.Evaluate(v.TranslateX, prevIndex, nowIndex, nextIndex)
	_, yy, _ = mmath.Evaluate(v.TranslateY, prevIndex, nowIndex, nextIndex)
	_, zy, _ = mmath.Evaluate(v.TranslateZ, prevIndex, nowIndex, nextIndex)
	_, ry, _ = mmath.Evaluate(v.Rotate, prevIndex, nowIndex, nextIndex)

	return xy, yy, zy, ry
}

func (c *BoneCurves) Merge() []int {
	return []int{
		int(c.TranslateX.Start.GetX()),
		c.values[1],
		c.values[2],
		c.values[3],
		int(c.TranslateX.Start.GetY()),
		c.values[5],
		c.values[6],
		c.values[7],
		int(c.TranslateX.End.GetX()),
		c.values[9],
		c.values[10],
		c.values[11],
		int(c.TranslateX.End.GetY()),
		c.values[13],
		c.values[14],
		c.values[15],
		int(c.TranslateY.Start.GetX()),
		c.values[17],
		c.values[18],
		c.values[19],
		int(c.TranslateY.Start.GetY()),
		c.values[21],
		c.values[22],
		c.values[23],
		int(c.TranslateY.End.GetX()),
		c.values[25],
		c.values[26],
		c.values[27],
		int(c.TranslateY.End.GetY()),
		c.values[29],
		c.values[30],
		c.values[31],
		int(c.TranslateZ.Start.GetX()),
		c.values[33],
		c.values[34],
		c.values[35],
		int(c.TranslateZ.Start.GetY()),
		c.values[37],
		c.values[38],
		c.values[39],
		int(c.TranslateZ.End.GetX()),
		c.values[41],
		c.values[42],
		c.values[43],
		int(c.TranslateZ.End.GetY()),
		c.values[45],
		c.values[46],
		c.values[47],
		int(c.Rotate.Start.GetX()),
		c.values[49],
		c.values[50],
		c.values[51],
		int(c.Rotate.Start.GetY()),
		c.values[53],
		c.values[54],
		c.values[55],
		int(c.Rotate.End.GetX()),
		c.values[57],
		c.values[58],
		c.values[59],
		int(c.Rotate.End.GetY()),
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
