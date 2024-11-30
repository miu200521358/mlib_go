package mmath

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type Curve struct {
	Start MVec2 // 三次ベジェ曲線のP1に相当
	End   MVec2 // 三次ベジェ曲線のP2に相当
}

const (
	// MMDでの補間曲線の最大値
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

// Copy
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

	curve.Start.MulScalar(CURVE_MAX).Round()
	curve.End.MulScalar(CURVE_MAX).Round()
}

func tryCurveNormalize(c0, c1, c2, c3 *MVec2) *Curve {
	p0 := c0
	p3 := c3

	diff := p3.Subed(p0)
	if diff.X == 0 {
		// 割算用なので1にしておく
		diff.X = 1
	}
	if diff.Y == 0 {
		// 割算用なので1にしておく
		diff.Y = 1
	}

	p1 := *c1.Subed(p0).Dived(diff)

	if p1.X < 0 || p1.X > 1 || p1.Y < 0 || p1.Y > 1 {
		return nil
	}

	p2 := *c2.Subed(p0).Dived(diff)

	if p2.X < 0 || p2.X > 1 || p2.Y < 0 || p2.Y > 1 {
		return nil
	}

	if NearEquals(p1.X, p1.Y, 1e-6) && NearEquals(p2.X, p2.Y, 1e-6) {
		return NewCurve()
	}

	curve := &Curve{
		Start: *p1.MuledScalar(CURVE_MAX).Round(),
		End:   *p2.MuledScalar(CURVE_MAX).Round(),
	}

	return curve
}

// https://pomax.github.io/bezierinfo
// https://shspage.hatenadiary.org/entry/20140625/1403702735
// https://bezier.readthedocs.io/en/stable/python/reference/bezier.curve.html#bezier.curve.Curve.evaluate
// https://edvakf.hatenadiary.org/entry/20111016/1318716097
// Evaluate 補間曲線を求めます。
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
		// 前後が同じ場合、必ず線形補間になる
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

		if math.Abs(funcDF) < eps {
			// 微分値が小さすぎる場合、微分値を1に設定
			funcDF = 1
		}

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

// NewCurveFromValues 関数は、与えられた値から補間曲線を生成します。
func NewCurveFromValues(values []float64) *Curve {
	n := len(values)
	if n < 3 {
		return NewCurve()
	}

	if IsAllSameValues(values) {
		return NewCurve()
	}

	// P0とP3の定義
	P0 := MVec2{X: 0, Y: values[0]}
	P3 := MVec2{X: float64(n - 1), Y: values[n-1]}

	// 最初と最後のポイントを除外
	m := n - 2
	A := mat.NewDense(m, 2, nil)
	bx := make([]float64, m)
	by := make([]float64, m)

	for i := 1; i < n-1; i++ {
		t := float64(i) / float64(n-1)
		B0 := (1 - t) * (1 - t) * (1 - t)
		B1 := 3 * (1 - t) * (1 - t) * t
		B2 := 3 * (1 - t) * t * t
		B3 := t * t * t

		// x座標用
		x := float64(i)
		Cx := B0*P0.X + B3*P3.X
		dx := x - Cx
		A.Set(i-1, 0, B1)
		A.Set(i-1, 1, B2)
		bx[i-1] = dx

		// y座標用
		y := values[i]
		Cy := B0*P0.Y + B3*P3.Y
		dy := y - Cy
		by[i-1] = dy
	}

	// 正則化パラメータの設定
	lambda := 1e-2

	// x座標用の正則化
	AtA := mat.NewDense(2, 2, nil)
	AtA.Mul(A.T(), A)
	lambdaI := mat.NewDiagDense(2, []float64{lambda, lambda})
	AtA.Add(AtA, lambdaI)
	AtASym := mat.NewSymDense(2, AtA.RawMatrix().Data)
	AtbX := mat.NewVecDense(2, nil)
	AtbX.MulVec(A.T(), mat.NewVecDense(m, bx))

	// x座標の解を求める
	var cholX mat.Cholesky
	if ok := cholX.Factorize(AtASym); !ok {
		return nil
	}
	pxVec := mat.NewVecDense(2, nil)
	err := cholX.SolveVecTo(pxVec, AtbX)
	if err != nil {
		return nil
	}
	px := pxVec.RawVector().Data

	// y座標用の正則化
	AtbY := mat.NewVecDense(2, nil)
	AtbY.MulVec(A.T(), mat.NewVecDense(m, by))

	// y座標の解を求める
	var cholY mat.Cholesky
	if ok := cholY.Factorize(AtASym); !ok {
		panic("Cholesky分解に失敗しました (y座標)")
	}
	pyVec := mat.NewVecDense(2, nil)
	err = cholY.SolveVecTo(pyVec, AtbY)
	if err != nil {
		panic(err)
	}
	py := pyVec.RawVector().Data

	// 計算された制御点を割り当て
	P1 := MVec2{X: px[0], Y: py[0]}
	P2 := MVec2{X: px[1], Y: py[1]}

	// 最適化された制御点
	return tryCurveNormalize(&P0, &P1, &P2, &P3)
}
