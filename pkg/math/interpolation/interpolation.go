package interpolation

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/math/mvec2"

)

type T struct {
	Start mvec2.T
	End   mvec2.T
}

var (
	LinearStart = mvec2.T{20.0, 20.0}
	LinearEnd   = mvec2.T{107.0, 107.0}

	// MMDでの補間曲線の最大値
	Max = 127.0
)

// https://pomax.github.io/bezierinfo
// https://shspage.hatenadiary.org/entry/20140625/1403702735
// https://bezier.readthedocs.io/en/stable/python/reference/bezier.curve.html#bezier.curve.Curve.evaluate
// https://edvakf.hatenadiary.org/entry/20111016/1318716097
// Evaluate 補間曲線を求めます。
// return x（計算キーフレ時点のX値）, y（計算キーフレ時点のY値）, t（計算キーフレまでの変化量）
func Evaluate(interpolation *T, start, now, end int) (float64, float64, float64) {
	if (now-start) == 0.0 || (end-start) == 0.0 {
		return 0.0, 0.0, 0.0
	}

	x := float64(now-start) / float64(end-start)

	if x >= 1 {
		return 1.0, 1.0, 1.0
	}

	if interpolation.Start.GetX() == interpolation.Start.GetY() && interpolation.End.GetX() == interpolation.End.GetY() {
		// 前後が同じ場合、必ず線形補間になる
		return x, x, x
	}

	x1 := interpolation.Start.GetX() / Max
	y1 := interpolation.Start.GetY() / Max
	x2 := interpolation.End.GetX() / Max
	y2 := interpolation.End.GetY() / Max

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
