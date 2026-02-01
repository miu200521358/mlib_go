// 指示: miu200521358
package mmath

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

// Curve は補間曲線を表す。
type Curve struct {
	Start Vec2
	End   Vec2
}

const CurveMax = 127.0

var (
	CURVE_MIN_VEC = Vec2{0.0, 0.0}
	CURVE_MAX_VEC = Vec2{CurveMax, CurveMax}
)

const (
	mathCurveFitFailedErrorID = "92106"
)

// newMathCurveFitFailed は曲線フィット失敗エラーを生成する。
func newMathCurveFitFailed(cause error) error {
	return merr.NewCommonError(mathCurveFitFailedErrorID, merr.ErrorKindInternal, "曲線フィットに失敗しました", cause)
}

// NewCurve は曲線を生成する。
func NewCurve() *Curve {
	return &Curve{
		Start: Vec2{20.0, 20.0},
		End:   Vec2{107.0, 107.0},
	}
}

// Normalize は正規化する。
func (curve *Curve) Normalize(begin, finish Vec2) {
	diff := finish.Subed(begin)

	curve.Start = curve.Start.Subed(begin).Dived(diff)
	curve.Start = curve.Start.Clamped(ZERO_VEC2, UNIT_XY_VEC2)

	curve.End = curve.End.Subed(begin).Dived(diff)
	curve.End = curve.End.Clamped(ZERO_VEC2, UNIT_XY_VEC2)

	if math.Abs(curve.Start.X-curve.Start.Y) <= 1e-6 && math.Abs(curve.End.X-curve.End.Y) <= 1e-6 {
		curve.Start = Vec2{20.0 / 127.0, 20.0 / 127.0}
		curve.End = Vec2{107.0 / 127.0, 107.0 / 127.0}
	}

	curve.Start = curve.Start.MuledScalar(CurveMax).Round()
	curve.End = curve.End.MuledScalar(CurveMax).Round()
}

// tryCurveNormalize は曲線の正規化を試みる。
func tryCurveNormalize(c0, c1, c2, c3 Vec2, decreasing bool) *Curve {
	p0 := Vec2{X: c0.X, Y: c0.Y}
	p3 := Vec2{X: c3.X, Y: c3.Y}

	diff := p3.Subed(p0)
	if diff.X == 0 {
		diff.X = 1
	}
	if diff.Y == 0 {
		diff.Y = 1
	}

	p1 := c1.Subed(p0).Dived(diff)
	p2 := c2.Subed(p0).Dived(diff)

	if math.Abs(p1.X-p1.Y) <= 1e-6 && math.Abs(p2.X-p2.Y) <= 1e-6 {
		return NewCurve()
	}

	if decreasing {
		p1, p2 = p2, p1
	}

	curve := &Curve{
		Start: p1.MuledScalar(CurveMax).Round(),
		End:   p2.MuledScalar(CurveMax).Round(),
	}

	if curve.Start.X < 0 || curve.Start.X > CurveMax || curve.Start.Y < 0 || curve.Start.Y > CurveMax ||
		curve.End.X < 0 || curve.End.X > CurveMax || curve.End.Y < 0 || curve.End.Y > CurveMax {
		return nil
	}

	return curve
}

// Evaluate は補間曲線を評価して値を返す。
func Evaluate(curve *Curve, start, now, end float32) (x, y, t float64) {
	if (now-start) == 0 || (end-start) == 0 {
		return 0, 0, 0
	}

	x = float64(now-start) / float64(end-start)
	if x >= 1 {
		return 1, 1, 1
	}

	if curve.Start.X == curve.Start.Y && curve.End.X == curve.End.Y {
		return x, x, x
	}

	x1 := curve.Start.X / CurveMax
	y1 := curve.Start.Y / CurveMax
	x2 := curve.End.X / CurveMax
	y2 := curve.End.Y / CurveMax

	// x に対応する t をニュートン法で求める。
	t = newton(x1, x2, x, x, 1e-15, 1e-20)
	s := 1.0 - t

	y = (3.0 * (math.Pow(s, 2.0)) * t * y1) + (3.0 * s * (math.Pow(t, 2.0)) * y2) + math.Pow(t, 3.0)

	return x, y, t
}

// bezierCoeffs はベジェ係数を計算する。
func bezierCoeffs(p1, p2 float64) (a, b, c float64) {
	a = 3*p1 - 3*p2 + 1
	b = -6*p1 + 3*p2
	c = 3 * p1
	return a, b, c
}

// bezierValue はベジェ曲線の値を計算する。
func bezierValue(a, b, c, t float64) float64 {
	return ((a*t+b)*t + c) * t
}

// bezierDerivative はベジェ曲線の導関数値を計算する。
func bezierDerivative(a, b, c, t float64) float64 {
	return (3*a*t+2*b)*t + c
}

// newton はニュートン法で解を求める。
func newton(x1, x2, x, t0, eps, err float64) float64 {
	a, b, c := bezierCoeffs(x1, x2)
	// ベジェ関数の逆関数をニュートン法で解く。
	for i := 0; i < 20; i++ {
		funcFValue := bezierValue(a, b, c, t0) - x
		funcDF := bezierDerivative(a, b, c, t0)
		if math.Abs(funcDF) < eps {
			funcDF = 1
		}
		t1 := t0 - funcFValue/funcDF
		if math.Abs(t1-t0) <= err {
			t0 = t1
			break
		}
		t0 = t1
	}
	return t0
}

// SplitCurve は曲線を分割する。
func SplitCurve(curve *Curve, start, now, end float32) (*Curve, *Curve) {
	if (now-start) == 0 || (end-start) == 0 {
		return NewCurve(), NewCurve()
	}

	_, _, t := Evaluate(curve, start, now, end)

	// de Casteljau 法で分割点を求める。
	iA := Vec2{0.0, 0.0}
	iB := curve.Start.DivedScalar(CurveMax)
	iC := curve.End.DivedScalar(CurveMax)
	iD := Vec2{1.0, 1.0}

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

	startCurve := &Curve{Start: iE, End: iH}
	startCurve.Normalize(iA, iJ)

	endCurve := &Curve{Start: iI, End: iG}
	endCurve.Normalize(iJ, iD)

	if startCurve.Start.X == startCurve.Start.Y &&
		startCurve.End.X == startCurve.End.Y &&
		endCurve.Start.X == endCurve.Start.Y &&
		endCurve.End.X == endCurve.End.Y {
		startCurve = NewCurve()
		endCurve = NewCurve()
	}

	return startCurve, endCurve
}

type controlPoints struct {
	P0, P1, P2, P3 Vec2
}

var optimizePointsFunc = optimizePoints

// NewCurveFromValues は曲線を生成する。
func NewCurveFromValues(values []float64, threshold float64) (*Curve, error) {
	if len(values) <= 2 {
		return NewCurve(), nil
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
		return NewCurve(), nil
	}

	yCoords := make([]float64, len(values))
	for i, v := range values {
		yCoords[i] = (v - yMin) / (yMax - yMin)
	}

	if isLinearInterpolation(yCoords, threshold) {
		return NewCurve(), nil
	}

	P0 := Vec2{X: xCoords[0], Y: yCoords[0]}
	P3 := Vec2{X: xCoords[len(xCoords)-1], Y: yCoords[len(yCoords)-1]}

	P1 := Vec2{X: xCoords[len(xCoords)/3], Y: yCoords[len(yCoords)/3]}
	P2 := Vec2{X: xCoords[2*len(xCoords)/3], Y: yCoords[2*len(yCoords)/3]}

	result, err := optimizePointsFunc(xCoords, yCoords, P0, P1, P2, P3)
	if err != nil {
		return nil, newMathCurveFitFailed(err)
	}

	curve := tryCurveNormalize(result.P0, result.P1, result.P2, result.P3, decreasing)
	if curve == nil {
		return nil, newMathCurveFitFailed(nil)
	}
	return curve, nil
}

// isLinearInterpolation は線形補間か判定する。
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

// optimizePoints は制御点を最適化する。
func optimizePoints(xCoords, yCoords []float64, P0, P1, P2, P3 Vec2) (controlPoints, error) {
	initial := []float64{P1.X, P1.Y, P2.X, P2.Y}

	// 誤差最小化問題をBFGSで解く。
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

	settings := &optimize.Settings{GradientThreshold: 1e-6, FuncEvaluations: 10000, MajorIterations: 1000}
	method := &optimize.BFGS{}

	result, err := optimize.Minimize(problem, initial, settings, method)
	if err != nil {
		return controlPoints{}, err
	}

	P1 = Vec2{X: result.X[0], Y: result.X[1]}
	P2 = Vec2{X: result.X[2], Y: result.X[3]}

	return controlPoints{P0, P1, P2, P3}, nil
}

// calculateError は誤差を計算する。
func calculateError(xCoords, yCoords []float64, P1, P2 Vec2) float64 {
	totalError := 0.0
	for i, x := range xCoords {
		t := newton(P1.X, P2.X, x, x, 1e-15, 1e-20)
		s := 1.0 - t
		y := (3.0 * (math.Pow(s, 2.0)) * t * P1.Y) + (3.0 * s * (math.Pow(t, 2.0)) * P2.Y) + math.Pow(t, 3.0)
		dy := yCoords[i] - y
		totalError += dy * dy
	}
	return totalError
}
