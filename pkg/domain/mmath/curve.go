// 指示: miu200521358
package mmath

import (
	"math"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

type Curve struct {
	Start MVec2
	End   MVec2
}

const (
	CURVE_MAX = 127.0
)

var CurveMin = &MVec2{0.0, 0.0}
var CurveMax = &MVec2{CURVE_MAX, CURVE_MAX}

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

func (curve *Curve) Copy() *Curve {
	copied := NewCurve()
	copied.Start.X = curve.Start.X
	copied.Start.Y = curve.Start.Y
	copied.End.X = curve.End.X
	copied.End.Y = curve.End.Y
	return copied
}

func (curve *Curve) Normalize(begin, finish *MVec2) {
	diff := finish.Subed(begin)

	curve.Start = *curve.Start.Sub(begin).Div(diff)

	if curve.Start.X < 0 {
		curve.Start.X = 0
	} else if curve.Start.X > 1 {
		curve.Start.X = 1
	}

	if curve.Start.Y < 0 {
		curve.Start.Y = 0
	} else if curve.Start.Y > 1 {
		curve.Start.Y = 1
	}

	curve.End = *curve.End.Sub(begin).Div(diff)

	if curve.End.X < 0 {
		curve.End.X = 0
	} else if curve.End.X > 1 {
		curve.End.X = 1
	}

	if curve.End.Y < 0 {
		curve.End.Y = 0
	} else if curve.End.Y > 1 {
		curve.End.Y = 1
	}

	if NearEquals(curve.Start.X, curve.Start.Y, 1e-6) && NearEquals(curve.End.X, curve.End.Y, 1e-6) {
		curve.Start = MVec2{20.0 / 127.0, 20.0 / 127.0}
		curve.End = MVec2{107.0 / 127.0, 107.0 / 127.0}
	}

	curve.Start = *curve.Start.MulScalar(CURVE_MAX).Round()
	curve.End = *curve.End.MulScalar(CURVE_MAX).Round()
}

func tryCurveNormalize(c0, c1, c2, c3 *MVec2, decreasing bool) *Curve {
	p0 := &MVec2{X: c0.X, Y: c0.Y}
	p3 := &MVec2{X: c3.X, Y: c3.Y}

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

func newtonFuncF(x1, x2, x, t float64) float64 {
	t1 := 1.0 - t
	return 3.0*(math.Pow(t1, 2.0))*t*x1 + 3.0*t1*(math.Pow(t, 2.0))*x2 + math.Pow(t, 3.0) - x
}

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

func SplitCurve(curve *Curve, start, now, end float32) (*Curve, *Curve) {
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
		startCurve = NewCurve()
		endCurve = NewCurve()
	}

	return startCurve, endCurve
}

type controlPoints struct {
	P0, P1, P2, P3 MVec2
}

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

	P0 := MVec2{X: xCoords[0], Y: yCoords[0]}
	P3 := MVec2{X: xCoords[len(xCoords)-1], Y: yCoords[len(yCoords)-1]}

	P1 := MVec2{X: xCoords[len(xCoords)/3], Y: yCoords[len(yCoords)/3]}
	P2 := MVec2{X: xCoords[2*len(xCoords)/3], Y: yCoords[2*len(yCoords)/3]}

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

func optimizePoints(xCoords, yCoords []float64, P0, P1, P2, P3 MVec2) (controlPoints, error) {
	initial := []float64{P1.X, P1.Y, P2.X, P2.Y}

	problem := optimize.Problem{
		Func: func(p []float64) float64 {
			P1 := MVec2{X: p[0], Y: p[1]}
			P2 := MVec2{X: p[2], Y: p[3]}
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

	P1 = MVec2{X: result.X[0], Y: result.X[1]}
	P2 = MVec2{X: result.X[2], Y: result.X[3]}

	return controlPoints{P0, P1, P2, P3}, nil
}

func calculateError(xCoords, yCoords []float64, P1, P2 MVec2) float64 {
	totalError := 0.0
	for i, x := range xCoords {
		t := newton(P1.X, P2.X, x, 0.5, 1e-15, 1e-20)
		s := 1.0 - t
		y := (3.0 * (math.Pow(s, 2.0)) * t * P1.Y) + (3.0 * s * (math.Pow(t, 2.0)) * P2.Y) + math.Pow(t, 3.0)
		totalError += math.Pow(yCoords[i]-y, 2)
	}
	return totalError
}

