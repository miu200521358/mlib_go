package mmath

import (
	"math"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

// ----- 定数 -----

const (
	// CURVE_MAX はMMDでの補間曲線の最大値です
	CURVE_MAX = 127.0
)

var (
	// CURVE_MIN は補間曲線の最小値ベクトルです
	CURVE_MIN = &Vec2{X: 0.0, Y: 0.0}

	// CURVE_MAX_VEC は補間曲線の最大値ベクトルです
	CURVE_MAX_VEC = &Vec2{X: CURVE_MAX, Y: CURVE_MAX}

	// CURVE_LINEAR は線形補間曲線です
	CURVE_LINEAR = &Curve{
		Start: Vec2{X: 20.0, Y: 20.0},
		End:   Vec2{X: 107.0, Y: 107.0},
	}
)

// ----- 型定義 -----

// Curve はMMDの補間曲線（三次ベジェ曲線）を表します
type Curve struct {
	Start Vec2 // 三次ベジェ曲線のP1に相当
	End   Vec2 // 三次ベジェ曲線のP2に相当
}

// ----- コンストラクタ -----

// NewCurve は線形補間曲線を作成します
func NewCurve() *Curve {
	return &Curve{
		Start: Vec2{X: 20.0, Y: 20.0},
		End:   Vec2{X: 107.0, Y: 107.0},
	}
}

// NewCurveByValues は指定した値で補間曲線を作成します
func NewCurveByValues(startX, startY, endX, endY byte) *Curve {
	if startX == 20 && startY == 20 && endX == 107 && endY == 107 {
		return CURVE_LINEAR
	}

	return &Curve{
		Start: Vec2{X: float64(startX), Y: float64(startY)},
		End:   Vec2{X: float64(endX), Y: float64(endY)},
	}
}

// ----- メソッド -----

// Copy はコピーを返します
func (c *Curve) Copy() *Curve {
	return &Curve{
		Start: Vec2{X: c.Start.X, Y: c.Start.Y},
		End:   Vec2{X: c.End.X, Y: c.End.Y},
	}
}

// IsLinear は線形補間かどうかを返します
func (c *Curve) IsLinear() bool {
	return c.Start.X == c.Start.Y && c.End.X == c.End.Y
}

// Normalize は補間曲線を正規化します（破壊的）
func (c *Curve) Normalize(begin, finish *Vec2) {
	diff := finish.Subed(begin)

	c.Start = *c.Start.Subed(begin).Dived(diff)
	c.Start.X = Clamped(c.Start.X, 0.0, 1.0)
	c.Start.Y = Clamped(c.Start.Y, 0.0, 1.0)

	c.End = *c.End.Subed(begin).Dived(diff)
	c.End.X = Clamped(c.End.X, 0.0, 1.0)
	c.End.Y = Clamped(c.End.Y, 0.0, 1.0)

	if NearEquals(c.Start.X, c.Start.Y, 1e-6) && NearEquals(c.End.X, c.End.Y, 1e-6) {
		c.Start = Vec2{X: 20.0 / 127.0, Y: 20.0 / 127.0}
		c.End = Vec2{X: 107.0 / 127.0, Y: 107.0 / 127.0}
	}

	c.Start = *c.Start.MuledScalar(CURVE_MAX).Round()
	c.End = *c.End.MuledScalar(CURVE_MAX).Round()
}

// ----- 評価関数 -----

// Evaluate は補間曲線を評価します
// return x（計算キーフレ時点のX値）, y（計算キーフレ時点のY値）, t（計算キーフレまでの変化量）
func Evaluate(curve *Curve, start, now, end float32) (x, y, t float64) {
	if (now-start) == 0.0 || (end-start) == 0.0 {
		return 0.0, 0.0, 0.0
	}

	x = float64(now-start) / float64(end-start)

	if x >= 1 {
		return 1.0, 1.0, 1.0
	}

	if curve.Start.X == curve.Start.Y && curve.End.X == curve.End.Y {
		return x, x, x
	}

	x1 := curve.Start.X / CURVE_MAX
	y1 := curve.Start.Y / CURVE_MAX
	x2 := curve.End.X / CURVE_MAX
	y2 := curve.End.Y / CURVE_MAX

	t = newton(x1, x2, x, 0.5, 1e-15, 1e-20)
	s := 1.0 - t

	y = (3.0 * (math.Pow(s, 2.0)) * t * y1) + (3.0 * s * (math.Pow(t, 2.0)) * y2) + math.Pow(t, 3.0)

	return x, y, t
}

// newtonFuncF は解を求める関数です
func newtonFuncF(x1, x2, x, t float64) float64 {
	t1 := 1.0 - t
	return 3.0*(math.Pow(t1, 2.0))*t*x1 + 3.0*t1*(math.Pow(t, 2.0))*x2 + math.Pow(t, 3.0) - x
}

// newton はニュートン法で解を求めます
func newton(x1, x2, x, t0, eps, err float64) float64 {
	derivative := 2.0 * eps

	for i := 0; i < 20; i++ {
		funcFValue := newtonFuncF(x1, x2, x, t0)
		funcDF := (newtonFuncF(x1, x2, x, t0+eps) - newtonFuncF(x1, x2, x, t0-eps)) / derivative

		if math.Abs(funcDF) < eps {
			funcDF = 1
		}

		t1 := t0 - funcFValue/funcDF

		if err >= math.Abs(t1-t0) {
			break
		}

		t0 = t1
	}

	return t0
}

// ----- 分割 -----

// SplitCurve は補間曲線を指定キーフレで前後に分割します
func SplitCurve(curve *Curve, start, now, end float32) (*Curve, *Curve) {
	if (now-start) == 0 || (end-start) == 0 {
		return NewCurve(), NewCurve()
	}

	_, _, t := Evaluate(curve, start, now, end)

	iA := &Vec2{X: 0.0, Y: 0.0}
	iB := curve.Start.DivedScalar(CURVE_MAX)
	iC := curve.End.DivedScalar(CURVE_MAX)
	iD := &Vec2{X: 1.0, Y: 1.0}

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

	startCurve := &Curve{Start: *iE, End: *iH}
	startCurve.Normalize(iA, iJ)

	endCurve := &Curve{Start: *iI, End: *iG}
	endCurve.Normalize(iJ, iD)

	if startCurve.IsLinear() && endCurve.IsLinear() {
		startCurve = NewCurve()
		endCurve = NewCurve()
	}

	return startCurve, endCurve
}

// ----- 曲線フィッティング -----

// NewCurveFromValues は指定された値のリストに基づいてベジェ曲線を近似します
func NewCurveFromValues(values []float64, threshold float64) *Curve {
	if len(values) <= 2 {
		return NewCurve()
	}

	decreasing := values[0] > values[len(values)-1]

	xCoords := make([]float64, len(values))
	for i := range values {
		if decreasing {
			xCoords[i] = 1.0 - float64(i)/float64(len(values)-1)
		} else {
			xCoords[i] = float64(i) / float64(len(values)-1)
		}
	}

	yMin := Min(values)
	yMax := Max(values)
	if yMin == yMax {
		return NewCurve()
	}

	yCoords := make([]float64, len(values))
	for i, v := range values {
		yCoords[i] = (v - yMin) / (yMax - yMin)
	}

	if isLinearInterpolation(yCoords, threshold) {
		return NewCurve()
	}

	P0 := Vec2{X: xCoords[0], Y: yCoords[0]}
	P3 := Vec2{X: xCoords[len(xCoords)-1], Y: yCoords[len(yCoords)-1]}
	P1 := Vec2{X: xCoords[len(xCoords)/3], Y: yCoords[len(yCoords)/3]}
	P2 := Vec2{X: xCoords[2*len(xCoords)/3], Y: yCoords[2*len(yCoords)/3]}

	result, err := optimizePoints(xCoords, yCoords, P0, P1, P2, P3)
	if err != nil {
		return nil
	}

	return tryCurveNormalize(&result.P0, &result.P1, &result.P2, &result.P3, decreasing)
}

func isLinearInterpolation(yCoords []float64, threshold float64) bool {
	if IsAlmostAllSameValues(yCoords, threshold) {
		return true
	}

	diffs := make([]float64, len(yCoords)-1)
	for i := 1; i < len(yCoords); i++ {
		diffs[i-1] = yCoords[i] - yCoords[i-1]
	}

	return IsAlmostAllSameValues(diffs, threshold)
}

type controlPoints struct {
	P0, P1, P2, P3 Vec2
}

func optimizePoints(xCoords, yCoords []float64, P0, P1, P2, P3 Vec2) (controlPoints, error) {
	initial := []float64{P1.X, P1.Y, P2.X, P2.Y}

	problem := optimize.Problem{
		Func: func(p []float64) float64 {
			P1 := Vec2{X: p[0], Y: p[1]}
			P2 := Vec2{X: p[2], Y: p[3]}
			return calculateError(xCoords, yCoords, P1, P2)
		},
	}

	problem.Grad = func(grad, p []float64) {
		fd.Gradient(grad, problem.Func, p, nil)
	}

	gradientThreshold := 1e-6
	settings := &optimize.Settings{GradientThreshold: gradientThreshold, FuncEvaluations: 10000, MajorIterations: 1000}
	method := &optimize.BFGS{}

	result, err := optimize.Minimize(problem, initial, settings, method)
	if err != nil {
		return controlPoints{}, err
	}

	P1 = Vec2{X: result.X[0], Y: result.X[1]}
	P2 = Vec2{X: result.X[2], Y: result.X[3]}

	return controlPoints{P0, P1, P2, P3}, nil
}

func calculateError(xCoords, yCoords []float64, P1, P2 Vec2) float64 {
	totalError := 0.0
	for i, x := range xCoords {
		t := newton(P1.X, P2.X, x, 0.5, 1e-15, 1e-20)
		s := 1.0 - t
		y := (3.0 * (math.Pow(s, 2.0)) * t * P1.Y) + (3.0 * s * (math.Pow(t, 2.0)) * P2.Y) + math.Pow(t, 3.0)
		totalError += math.Pow(yCoords[i]-y, 2)
	}
	return totalError
}

func tryCurveNormalize(c0, c1, c2, c3 *Vec2, decreasing bool) *Curve {
	p0 := &Vec2{X: c0.X, Y: c0.Y}
	p3 := &Vec2{X: c3.X, Y: c3.Y}

	diff := p3.Subed(p0)
	if diff.X == 0 {
		diff.X = 1
	}
	if diff.Y == 0 {
		diff.Y = 1
	}

	p1 := c1.Subed(p0).Dived(diff)
	p2 := c2.Subed(p0).Dived(diff)

	if NearEquals(p1.X, p1.Y, 1e-6) && NearEquals(p2.X, p2.Y, 1e-6) {
		return NewCurve()
	}

	if decreasing {
		p1, p2 = p2, p1
	}

	curve := &Curve{
		Start: *p1.MuledScalar(CURVE_MAX).Round(),
		End:   *p2.MuledScalar(CURVE_MAX).Round(),
	}

	if curve.Start.X < 0 || curve.Start.X > CURVE_MAX || curve.Start.Y < 0 || curve.Start.Y > CURVE_MAX ||
		curve.End.X < 0 || curve.End.X > CURVE_MAX || curve.End.Y < 0 || curve.End.Y > CURVE_MAX {
		return nil
	}

	return curve
}
