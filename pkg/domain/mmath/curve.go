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
	n := len(values) - 1

	// 特殊ケース処理
	switch len(values) {
	case 1, 2:
		// 値が1、2つだけの場合、線形補間を生成
		return NewCurve()
	case 3:
		// 値が3つだけの場合、始点と終点の中間に制御点を配置
		P0 := MVec2{X: 0, Y: values[0]}
		P1 := MVec2{X: 0.33, Y: values[1]} // 中間制御点
		P2 := MVec2{X: 0.66, Y: values[1]} // 同じ位置の制御点
		P3 := MVec2{X: 1, Y: values[2]}
		return tryCurveNormalize(&P0, &P1, &P2, &P3)
	}

	if IsAllSameValues(values) {
		return NewCurve()
	}

	// 値を正規化（スケーリング）
	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// 正規化した値を作成
	scale := maxVal - minVal
	normalizedValues := make([]float64, len(values))
	for i, v := range values {
		normalizedValues[i] = (v - minVal) / scale
	}

	// ステップ1: t パラメータの計算（0から1に正規化）
	t := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		t[i] = float64(i) / float64(n)
	}

	// ステップ2: 基底関数の計算
	b0 := make([]float64, n+1)
	b1 := make([]float64, n+1)
	b2 := make([]float64, n+1)
	b3 := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		ti := t[i]
		oneMinusT := 1 - ti
		b0[i] = oneMinusT * oneMinusT * oneMinusT
		b1[i] = 3 * oneMinusT * oneMinusT * ti
		b2[i] = 3 * oneMinusT * ti * ti
		b3[i] = ti * ti * ti
	}

	// ステップ3: 始点 P0 と 終点 P3 の設定
	P0 := MVec2{X: t[0], Y: normalizedValues[0]}
	P3 := MVec2{X: t[n], Y: normalizedValues[n]}

	// ステップ4: 行列 A とベクトル Y の構築（Y は Y 成分のみ）
	AData := make([]float64, 2*(n+1))
	YData := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		// 行列 A の要素（b1 と b2）
		AData[i*2] = b1[i]
		AData[i*2+1] = b2[i]
		// ベクトル Y の要素（Y 成分のみ）
		YData[i] = normalizedValues[i] - (b0[i]*P0.Y + b3[i]*P3.Y)
	}

	A := mat.NewDense(n+1, 2, AData)
	Y := mat.NewVecDense(n+1, YData)

	// ステップ5: 正規方程式の構築と解法
	// AT = A^T * A
	var AT mat.Dense
	AT.Mul(A.T(), A)
	// ATY = A^T * Y
	var ATY mat.VecDense
	ATY.MulVec(A.T(), Y)

	// 制御点 P1 と P2 の Y 値の計算
	PY := mat.NewVecDense(2, nil)
	err := PY.SolveVec(&AT, &ATY)
	if err != nil {
		return nil
	}

	// 制御点の設定（X 値は t パラメータに基づく）
	P1 := MVec2{X: t[1], Y: PY.AtVec(0)}
	P2 := MVec2{X: t[n-1], Y: PY.AtVec(1)}

	// Yの正規化
	yMin := MinFloat([]float64{P0.Y, P1.Y, P2.Y, P3.Y})
	yMax := MaxFloat([]float64{P0.Y, P1.Y, P2.Y, P3.Y})
	yDiff := yMax - yMin
	if yDiff == 0 {
		return NewCurve()
	}

	P0.Y = (P0.Y - yMin) / yDiff
	P1.Y = (P1.Y - yMin) / yDiff
	P2.Y = (P2.Y - yMin) / yDiff
	P3.Y = (P3.Y - yMin) / yDiff

	// 単調減少している場合、反転
	if P0.Y > P3.Y {
		P0.Y = 1 - P0.Y
		P1.Y = 1 - P1.Y
		P2.Y = 1 - P2.Y
		P3.Y = 1 - P3.Y
	}

	// 最適化された制御点
	return tryCurveNormalize(&P0, &P1, &P2, &P3)
}
