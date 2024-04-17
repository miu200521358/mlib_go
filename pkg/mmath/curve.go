package mmath

import (
	"math"
	"sort"

	"github.com/jinzhu/copier"
)

type Curve struct {
	Start *MVec2
	End   *MVec2
}

const (
	// MMDでの補間曲線の最大値
	CURVE_MAX = 127.0
)

func NewCurve() *Curve {
	return &Curve{
		Start: &MVec2{20.0, 20.0},
		End:   &MVec2{107.0, 107.0},
	}
}

// Copy
func (v *Curve) Copy() *Curve {
	copied := NewCurve()
	copier.CopyWithOption(copied, v, copier.Option{DeepCopy: true})
	return copied
}

func (v *Curve) Normalize(begin, finish *MVec2) {
	diff := finish.Subed(begin)

	v.Start = v.Start.Sub(begin)
	v.Start = v.Start.Div(diff)
	v.Start = v.Start.MulScalar(CURVE_MAX)
	v.Start = v.Start.Round()

	if v.Start.GetX() < 0 {
		v.Start.SetX(0)
	} else if v.Start.GetX() > CURVE_MAX {
		v.Start.SetX(CURVE_MAX)
	}

	if v.Start.GetY() < 0 {
		v.Start.SetY(0)
	} else if v.Start.GetY() > CURVE_MAX {
		v.Start.SetY(CURVE_MAX)
	}

	v.End = v.End.Sub(begin)
	v.End = v.End.Div(diff)
	v.End = v.End.MulScalar(CURVE_MAX)
	v.End = v.End.Round()

	if v.End.GetX() < 0 {
		v.End.SetX(0)
	} else if v.End.GetX() > CURVE_MAX {
		v.End.SetX(CURVE_MAX)
	}

	if v.End.GetY() < 0 {
		v.End.SetY(0)
	} else if v.End.GetY() > CURVE_MAX {
		v.End.SetY(CURVE_MAX)
	}
}

// https://pomax.github.io/bezierinfo
// https://shspage.hatenadiary.org/entry/20140625/1403702735
// https://bezier.readthedocs.io/en/stable/python/reference/bezier.curve.html#bezier.curve.Curve.evaluate
// https://edvakf.hatenadiary.org/entry/20111016/1318716097
// Evaluate 補間曲線を求めます。
// return x（計算キーフレ時点のX値）, y（計算キーフレ時点のY値）, t（計算キーフレまでの変化量）
func Evaluate(curve *Curve, start, now, end float64) (float64, float64, float64) {
	if (now-start) == 0.0 || (end-start) == 0.0 {
		return 0.0, 0.0, 0.0
	}

	x := float64(now-start) / float64(end-start)

	if x >= 1 {
		return 1.0, 1.0, 1.0
	}

	if curve.Start.GetX() == curve.Start.GetY() && curve.End.GetX() == curve.End.GetY() {
		// 前後が同じ場合、必ず線形補間になる
		return x, x, x
	}

	x1 := curve.Start.GetX() / CURVE_MAX
	y1 := curve.Start.GetY() / CURVE_MAX
	x2 := curve.End.GetX() / CURVE_MAX
	y2 := curve.End.GetY() / CURVE_MAX

	t := newton(x1, x2, x, 0.5, 1e-15, 1e-20)
	s := 1.0 - t

	y := (3.0 * (math.Pow(s, 2.0)) * t * y1) + (3.0 * s * (math.Pow(t, 2.0)) * y2) + math.Pow(t, 3.0)

	return x, y, t
}

// 解を求める関数
func newtonFuncF(x1, x2, x, t float64) float64 {
	t1 := 1.0 - t
	return 3.0*(math.Pow(t1, 2.0))*t*x1 + 3.0*t1*(math.Pow(t, 2.0))*x2 + math.Pow(t, 3.0) - x
}

// Newton法（方程式の関数項、探索の開始点、微小量、誤差範囲、最大反復回数）
func newton(x1, x2, x, t0, eps, err float64) float64 {
	derivative := 2.0 * eps

	for i := 0; i < 20; i++ {
		funcFValue := newtonFuncF(x1, x2, x, t0)
		// 中心差分による微分値
		funcDF := (newtonFuncF(x1, x2, x, t0+eps) - newtonFuncF(x1, x2, x, t0-eps)) / derivative

		// 次の解を計算
		t1 := t0 - funcFValue/funcDF

		if err >= math.Abs(t1-t0) {
			// 「誤差範囲が一定値以下」ならば終了
			break
		}

		// 解を更新
		t0 = t1
	}

	return t0
}

// SplitCurve 補間曲線を指定キーフレで前後に分割する
func SplitCurve(curve *Curve, start, now, end float64) (*Curve, *Curve) {
	if (now-start) == 0 || (end-start) == 0 {
		return NewCurve(), NewCurve()
	}

	_, _, t := Evaluate(curve, start, now, end)

	iA := &MVec2{0.0, 0.0}
	iB := curve.Start.DivedScalar(CURVE_MAX)
	iC := curve.End.DivedScalar(CURVE_MAX)
	iD := &MVec2{1.0, 1.0}

	iAt1 := iA.MuledScalar(1 - t)
	iBt1 := iB.MuledScalar(1 - t)
	iBt2 := iB.MuledScalar(t)
	iCt1 := iC.MuledScalar(1 - t)
	iCt2 := iC.MuledScalar(t)
	iDt2 := iD.MuledScalar(t)

	iE := iAt1.Added(iBt2)
	iF := iBt1.Added(iCt2)
	iG := iCt1.Added(iDt2)

	iEt1 := iE.MuledScalar(1 - t)
	iFt1 := iF.MuledScalar(1 - t)
	iFt2 := iF.MuledScalar(t)
	iGt2 := iG.MuledScalar(t)

	iH := iEt1.Added(iFt2)
	iI := iFt1.Added(iGt2)

	iHt1 := iH.MuledScalar(1 - t)
	iIt2 := iI.MuledScalar(t)

	iJ := iHt1.Added(iIt2)

	// 新たな4つのベジェ曲線の制御点は、A側がAEHJ、C側がJIGDとなる。
	startCurve := &Curve{
		Start: iE,
		End:   iH,
	}
	startCurve.Normalize(iA, iJ)

	endCurve := &Curve{
		Start: iI,
		End:   iG,
	}
	endCurve.Normalize(iJ, iD)

	if startCurve.Start.GetX() == startCurve.Start.GetY() &&
		startCurve.End.GetX() == startCurve.End.GetY() &&
		endCurve.Start.GetX() == endCurve.Start.GetY() &&
		endCurve.End.GetX() == endCurve.End.GetY() {
		// 線形の場合初期化
		startCurve = NewCurve()
		endCurve = NewCurve()
	}

	return startCurve, endCurve
}

// GetInfections 補間曲線を分割するキーフレを求めます
func GetInfections(values []float64, threshold float64, tolerance float64) []int {
	// Calculate the first derivative
	ysPrime := make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		if i == 0 {
			ysPrime[i] = values[i+1] - values[i]
		} else if i == len(values)-1 {
			ysPrime[i] = values[i] - values[i-1]
		} else {
			ysPrime[i] = (values[i+1] - values[i-1]) / 2
		}
	}

	// Calculate the second derivative
	ysDoublePrime := make([]float64, len(ysPrime))
	for i := 0; i < len(ysPrime); i++ {
		if i == 0 {
			ysDoublePrime[i] = ysPrime[i+1] - ysPrime[i]
		} else if i == len(ysPrime)-1 {
			ysDoublePrime[i] = ysPrime[i] - ysPrime[i-1]
		} else {
			ysDoublePrime[i] = (ysPrime[i+1] - ysPrime[i-1]) / 2
		}
	}

	// Find the inflection points
	inflectionPoints := make([]int, 0)
	for i := 0; i < len(ysDoublePrime); i++ {
		if math.Abs(ysDoublePrime[i]) >= threshold {
			inflectionPoints = append(inflectionPoints, i)
		}
	}

	// Sign the second derivative with tolerance
	ysDoublePrimeSign := make([]float64, len(ysDoublePrime))
	for i := 0; i < len(ysDoublePrime); i++ {
		if ysDoublePrime[i] >= tolerance {
			ysDoublePrimeSign[i] = 1
		} else if ysDoublePrime[i] <= -tolerance {
			ysDoublePrimeSign[i] = -1
		} else {
			ysDoublePrimeSign[i] = 0
		}
	}

	// Calculate the difference between each element
	ysDoublePrimeDiff := make([]float64, len(ysDoublePrimeSign)-1)
	for i := 0; i < len(ysDoublePrimeSign)-1; i++ {
		ysDoublePrimeDiff[i] = ysDoublePrimeSign[i+1] - ysDoublePrimeSign[i]
	}

	// Find the indices where the difference is not 0 (i.e., the sign has changed)
	signChangeIndices := make([]int, 0)
	for i := 0; i < len(ysDoublePrimeDiff); i++ {
		if ysDoublePrimeDiff[i] != 0 {
			signChangeIndices = append(signChangeIndices, i)
		}
	}

	inflections := make([]int, 0)
	inflections = append(inflections, inflectionPoints...)
	inflections = append(inflections, signChangeIndices...)
	inflections = append(inflections, 0, len(values)-1)
	sort.Ints(inflections)

	return inflections
}
