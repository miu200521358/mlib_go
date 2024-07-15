package mmath

import (
	"math"

	"gonum.org/v1/gonum/optimize"
)

type Curve struct {
	Start MVec2
	End   MVec2
}

const (
	// MMDでの補間曲線の最大値
	CURVE_MAX = 127.0
)

var LINER_CURVE = &Curve{
	Start: MVec2{20.0, 20.0},
	End:   MVec2{107.0, 107.0},
}

func NewCurve() *Curve {
	return &Curve{
		Start: MVec2{20.0, 20.0},
		End:   MVec2{107.0, 107.0},
	}
}

func NewCurveByValues(startX, startY, endX, endY byte) *Curve {
	if startX == 20 && startY == 20 && endX == 107 && endY == 107 {
		return LINER_CURVE
	}

	return &Curve{
		Start: MVec2{float64(startX), float64(startY)},
		End:   MVec2{float64(endX), float64(endY)},
	}
}

// Copy
func (v *Curve) Copy() *Curve {
	copied := NewCurve()
	copied.Start.X = v.Start.X
	copied.Start.Y = v.Start.Y
	copied.End.X = v.End.X
	copied.End.Y = v.End.Y
	return copied
}

func (v *Curve) Normalize(begin, finish *MVec2) {
	diff := finish.Subed(begin)

	v.Start = *v.Start.Sub(begin).Div(diff).MulScalar(CURVE_MAX).Round()

	if v.Start.X < 0 {
		v.Start.X = 0
	} else if v.Start.X > CURVE_MAX {
		v.Start.X = CURVE_MAX
	}

	if v.Start.Y < 0 {
		v.Start.Y = 0
	} else if v.Start.Y > CURVE_MAX {
		v.Start.Y = CURVE_MAX
	}

	v.End = *v.End.Sub(begin).Div(diff).MulScalar(CURVE_MAX).Round()

	if v.End.X < 0 {
		v.End.X = 0
	} else if v.End.X > CURVE_MAX {
		v.End.X = CURVE_MAX
	}

	if v.End.Y < 0 {
		v.End.Y = 0
	} else if v.End.Y > CURVE_MAX {
		v.End.Y = CURVE_MAX
	}
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

	if curve.Start.X == curve.Start.Y && curve.End.X == curve.End.Y {
		// 前後が同じ場合、必ず線形補間になる
		return x, x, x
	}

	x1 := curve.Start.X / CURVE_MAX
	y1 := curve.Start.Y / CURVE_MAX
	x2 := curve.End.X / CURVE_MAX
	y2 := curve.End.Y / CURVE_MAX

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
func SplitCurve(curve *Curve, start, now, end int) (*Curve, *Curve) {
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
		Start: *iE,
		End:   *iH,
	}
	startCurve.Normalize(iA, iJ)

	endCurve := &Curve{
		Start: *iI,
		End:   *iG,
	}
	endCurve.Normalize(iJ, iD)

	if startCurve.Start.X == startCurve.Start.Y &&
		startCurve.End.X == startCurve.End.Y &&
		endCurve.Start.X == endCurve.Start.Y &&
		endCurve.End.X == endCurve.End.Y {
		// 線形の場合初期化
		startCurve = NewCurve()
		endCurve = NewCurve()
	}

	return startCurve, endCurve
}

func Bezier(t float64, p0, p1, p2, p3 *MVec2) *MVec2 {
	t2 := t * t
	t3 := t2 * t
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt

	bx := mt3*p0.X + 3*mt2*t*p1.X + 3*mt*t2*p2.X + t3*p3.X
	by := mt3*p0.Y + 3*mt2*t*p1.Y + 3*mt*t2*p2.Y + t3*p3.Y

	return &MVec2{bx, by}
}

// NewCurveFromValues calculates the control points of a cubic Bezier curve
// that best fits a given set of y-values.
func NewCurveFromValues(values []float64) *Curve {
	n := len(values)
	if n <= 2 {
		return NewCurve()
	}

	// Set start and end points
	p0 := &MVec2{0, values[0]}
	p3 := &MVec2{float64(n - 1), values[n-1]}

	// Initial guesses for control points
	p1 := &MVec2{1, values[1]}
	p2 := &MVec2{float64(n - 2), values[n-2]}

	funcEval := func(x []float64) float64 {
		p1 = &MVec2{x[0], x[1]}
		p2 = &MVec2{x[2], x[3]}
		sumSq := 0.0
		for i, y := range values {
			t := float64(i) / float64(n-1)
			bp := Bezier(t, p0, p1, p2, p3)
			diff := bp.Y - y
			sumSq += diff * diff
		}
		return sumSq
	}

	// Define the optimization problem
	problem := optimize.Problem{
		Func: funcEval,
		Grad: func(grad, x []float64) {
			h := 1e-6
			fx := funcEval(x)
			for i := range x {
				orig := x[i]
				x[i] += h
				fxh := funcEval(x)
				x[i] = orig
				grad[i] = (fxh - fx) / h
			}
		},
	}

	// Optimization settings
	settings := &optimize.Settings{
		MajorIterations:   1000,  // 最大イテレーション数を増やす
		FuncEvaluations:   10000, // 関数評価の最大数を増やす
		GradientThreshold: 1e-6,  // 勾配の閾値を変更する
	}
	method := &optimize.LBFGS{}

	// Initial control points vector
	initial := []float64{p1.X, p1.Y, p2.X, p2.Y}

	// Perform optimization
	result, err := optimize.Minimize(problem, initial, settings, method)
	if err != nil {
		return NewCurve()
	}

	// Update p1 and p2 with optimized values
	c := &Curve{
		Start: MVec2{result.X[0], result.X[1]},
		End:   MVec2{result.X[2], result.X[3]},
	}
	c.Normalize(p0, p3)

	if c.Start.X/c.Start.Y == 1.0 && c.End.X/c.End.Y == 1.0 {
		// 線形の場合初期化
		c = NewCurve()
	}

	return c
}
