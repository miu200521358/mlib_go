package mmath

import (
	"math"

)

type Curve struct {
	Start *MVec2
	End   *MVec2
}

var (
	LinearStart = MVec2{20.0, 20.0}
	LinearEnd   = MVec2{107.0, 107.0}

	// MMDでの補間曲線の最大値
	CurveMax = 127.0
)

func NewCurve() *Curve {
	return &Curve{
		Start: &LinearStart,
		End:   &LinearEnd,
	}
}

// Copy
func (v *Curve) Copy() *Curve {
	return &Curve{
		Start: v.Start.Copy(),
		End:   v.End.Copy(),
	}
}

func (v *Curve) Normalize(begin, finish *MVec2) {
	diff := finish.Sub(begin)
	v.Start = v.Start.Sub(begin).Div(diff)
	v.End = v.End.Sub(begin).Div(diff)
}

// https://pomax.github.io/bezierinfo
// https://shspage.hatenadiary.org/entry/20140625/1403702735
// https://bezier.readthedocs.io/en/stable/python/reference/bezier.curve.html#bezier.curve.Curve.evaluate
// https://edvakf.hatenadiary.org/entry/20111016/1318716097
// Evaluate 補間曲線を求めます。
// return x（計算キーフレ時点のX値）, y（計算キーフレ時点のY値）, t（計算キーフレまでの変化量）
func Evaluate(curve *Curve, start, now, end int) (float64, float64, float64) {
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

	x1 := curve.Start.GetX() / CurveMax
	y1 := curve.Start.GetY() / CurveMax
	x2 := curve.End.GetX() / CurveMax
	y2 := curve.End.GetY() / CurveMax

	t := newton(x1, x2, x, 0.5, 1e-10, 1e-15)
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

	for i := 0; i < 10; i++ {
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
func SplitCurve(interpolation *Curve, start, now, end int) (*Curve, *Curve) {
	if (now-start) == 0 || (end-start) == 0 {
		return NewCurve(), NewCurve()
	}

	_, _, t := Evaluate(interpolation, start, now, end)

	iA := MVec2{0.0, 0.0}
	iB := interpolation.Start.DivScalar(CurveMax)
	iC := interpolation.End.DivScalar(CurveMax)
	iD := MVec2{1.0, 1.0}

	iE := iA.MulScalar(1 - t).Added(iB.MulScalar(t))
	iF := iB.MulScalar(1 - t).Added(iC.MulScalar(t))
	iG := iC.MulScalar(1 - t).Added(iD.MulScalar(t))
	iH := iE.MulScalar(1 - t).Added(iF.MulScalar(t))
	iI := iF.MulScalar(1 - t).Added(iG.MulScalar(t))
	iJ := iH.MulScalar(1 - t).Added(iI.MulScalar(t))

	// 新たな4つのベジェ曲線の制御点は、A側がAEHJ、C側がJIGDとなる。
	startCurve := &Curve{
		Start: &iE,
		End:   &iH,
	}
	startCurve.Normalize(&iA, &iJ)

	endCurve := &Curve{
		Start: &iI,
		End:   &iG,
	}
	endCurve.Normalize(&iJ, &iD)

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
