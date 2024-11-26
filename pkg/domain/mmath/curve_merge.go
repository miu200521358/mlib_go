package mmath

import (
	"fmt"
	"math"
)

type bezierCurve struct {
	P0, P1, P2, P3 MVec2
	Length         float64
}

func (curve *bezierCurve) String() string {
	return fmt.Sprintf("[P0=%v, P1=%v, P2=%v, P3=%v, Len=%.7f]", curve.P0, curve.P1, curve.P2, curve.P3, curve.Length)
}

// CombineCurves は、与えられたベジェ曲線のリストを結合し、可能であれば結合されたベジェ曲線のリストを返します。
func CombineCurves(frames []float32, values []float64, curves []*Curve) (combinedFrames []float32, combinedCurves []*Curve) {
	// ベジェ曲線の数を取得
	n := len(curves)
	if n < 2 || len(frames) != n+1 || len(values) != n+1 {
		// 入力データの長さが不正な場合、元のデータを返す
		combinedFrames = frames
		combinedCurves = curves
		return
	}

	// ベジェ曲線のリストを実際の値に変換し、bezierCurve構造体に格納
	bezierCurves := make([]*bezierCurve, n)
	for i := 0; i < n; i++ {
		startFrame := frames[i]
		endFrame := frames[i+1]
		startValue := values[i]
		endValue := values[i+1]
		curve := curves[i]

		// 制御点を実際の値に変換
		bezier := convertCurveToBezierCurve(curve, startFrame, endFrame, startValue, endValue)
		// 曲線長を計算
		bezier.Length = BezierCurveLength(bezier)
		bezierCurves[i] = bezier
	}

	// 結合処理
	i := 0
	for i < n {
		// 現在の曲線を取得
		bezier1 := bezierCurves[i]
		frame1Start := frames[i]
		// frame1End := frames[i+1]
		value1Start := values[i]
		// value1End := values[i+1]
		curve1 := curves[i]

		if i+1 < n {
			// 次の曲線を取得
			bezier2 := bezierCurves[i+1]
			// frame2Start := frames[i+1]
			frame2End := frames[i+2]
			// value2Start := values[i+1]
			value2End := values[i+2]
			// curve2 := curves[i+1]

			// 結合可能か判定
			canCombine := checkCanCombine(bezier1, bezier2)
			if canCombine {
				// 結合
				combinedBezier := combineBezierCurves(bezier1, bezier2)
				// 新しい曲線長を計算
				combinedBezier.Length = BezierCurveLength(combinedBezier)

				// 結合した結果の長さが元の長さの合計と一致するか確認
				totalLength := bezier1.Length + bezier2.Length
				if math.Abs(combinedBezier.Length-totalLength) < 1e-6 {
					// 制御点の微調整: G3連続性を考慮
					adjustControlPointsForG3Continuity(combinedBezier)

					// 新しいCurveを作成し、正規化
					newCurve := createCurveFromBezierCurve(combinedBezier)
					newCurve.Normalize(
						&MVec2{X: float64(frame1Start), Y: value1Start},
						&MVec2{X: float64(frame2End), Y: value2End})

					// 結合結果を追加
					combinedFrames = append(combinedFrames, frame1Start)
					combinedCurves = append(combinedCurves, newCurve)

					// インデックスを2つ進める
					i += 2
					continue
				}
			}
		}

		// 結合できなかった場合、元の曲線を追加
		combinedFrames = append(combinedFrames, frame1Start)
		combinedCurves = append(combinedCurves, curve1)
		i++
	}

	// 最後のフレームを追加
	combinedFrames = append(combinedFrames, frames[n])

	return
}

// Curveを実際のbezierCurveに変換する関数
func convertCurveToBezierCurve(curve *Curve, startFrame, endFrame float32, startValue, endValue float64) *bezierCurve {
	// 制御点を計算
	p0 := MVec2{X: float64(startFrame), Y: startValue}
	p3 := MVec2{X: float64(endFrame), Y: endValue}

	diffFrame := float64(endFrame - startFrame)
	diffValue := endValue - startValue

	// 制御点P1, P2を実際の値にスケーリング
	p1 := MVec2{
		X: float64(startFrame) + (curve.Start.X/127.0)*diffFrame,
		Y: startValue + (curve.Start.Y/127.0)*diffValue,
	}
	p2 := MVec2{
		X: float64(startFrame) + (curve.End.X/127.0)*diffFrame,
		Y: startValue + (curve.End.Y/127.0)*diffValue,
	}

	return &bezierCurve{
		P0: p0,
		P1: p1,
		P2: p2,
		P3: p3,
	}
}

// ベジェ曲線の長さを適応的シンプソン法で計算する関数
func BezierCurveLength(curve *bezierCurve) float64 {
	// 目標とする精度
	epsilon := 1e-6
	return adaptiveSimpson(curve, 0.0, 1.0, epsilon)
}

// ベジェ曲線の微分を計算する関数
func bezierDerivative(curve *bezierCurve, t float64) *MVec2 {
	// 一次導関数の計算
	mt := 1 - t
	a := curve.P0.MuledScalar(-3 * mt * mt)
	b := curve.P1.MuledScalar(3*mt*mt - 6*mt*t)
	c := curve.P2.MuledScalar(6*mt*t - 3*t*t)
	d := curve.P3.MuledScalar(3 * t * t)

	return a.Added(b).Added(c).Added(d)
}

// 適応的シンプソン法による数値積分
func adaptiveSimpson(curve *bezierCurve, a, b, epsilon float64) float64 {
	// シンプソン近似を計算
	c := (a + b) / 2
	simpson := func(fa, fb, fc, h float64) float64 {
		return (h / 6) * (fa + 4*fc + fb)
	}

	h := b - a
	fa := bezierDerivative(curve, a).Length()
	fb := bezierDerivative(curve, b).Length()
	fc := bezierDerivative(curve, c).Length()
	s := simpson(fa, fb, fc, h)

	return adaptiveSimpsonAux(curve, a, b, epsilon, s, fa, fb, fc)
}

func adaptiveSimpsonAux(curve *bezierCurve, a, b, epsilon, s, fa, fb, fc float64) float64 {
	c := (a + b) / 2
	h := b - a
	d := (a + c) / 2
	e := (c + b) / 2
	fd := bezierDerivative(curve, d).Length()
	fe := bezierDerivative(curve, e).Length()

	simpson := func(fa, fb, fc, h float64) float64 {
		return (h / 6) * (fa + 4*fc + fb)
	}

	sLeft := simpson(fa, fc, fd, h/2)
	sRight := simpson(fc, fb, fe, h/2)
	s2 := sLeft + sRight

	if math.Abs(s2-s) < 15*epsilon {
		return s2 + (s2-s)/15
	} else {
		return adaptiveSimpsonAux(curve, a, c, epsilon/2, sLeft, fa, fc, fd) +
			adaptiveSimpsonAux(curve, c, b, epsilon/2, sRight, fc, fb, fe)
	}
}

// 2つのベジェ曲線が結合可能か判定する関数
func checkCanCombine(bezier1, bezier2 *bezierCurve) bool {
	// 単調増加または単調減少であるかを確認
	if !isMonotonic(bezier1) || !isMonotonic(bezier2) {
		return false
	}

	// 変曲点があるか確認（詳細な検出が必要なら実装）
	if hasInflectionPoint(bezier1) || hasInflectionPoint(bezier2) {
		return false
	}

	return true
}

// ベジェ曲線が単調増加または単調減少かを判定する関数
func isMonotonic(curve *bezierCurve) bool {
	// tを細かく分割して、曲線の増減を確認
	prevY := curve.P0.Y
	isIncreasing := true
	isDecreasing := true

	for t := 0.01; t <= 1.0; t += 0.01 {
		point := bezierPoint(curve, t)
		if point.Y < prevY {
			isIncreasing = false
		}
		if point.Y > prevY {
			isDecreasing = false
		}
		prevY = point.Y
	}

	return isIncreasing || isDecreasing
}

// ベジェ曲線上の点を計算する関数
func bezierPoint(curve *bezierCurve, t float64) *MVec2 {
	mt := 1 - t
	mt2 := mt * mt
	t2 := t * t

	a := curve.P0.MuledScalar(mt2 * mt)
	b := curve.P1.MuledScalar(3 * mt2 * t)
	c := curve.P2.MuledScalar(3 * mt * t2)
	d := curve.P3.MuledScalar(t2 * t)

	return a.Added(b).Added(c).Added(d)
}

// ベジェ曲線に変曲点があるかを判定する関数（簡易的な実装）
func hasInflectionPoint(curve *bezierCurve) bool {
	// 二次導関数を用いて変曲点を検出
	// 実装が複雑になるため、ここでは仮実装としてfalseを返す
	return false
}

// 2つのベジェ曲線を結合する関数
func combineBezierCurves(bezier1, bezier2 *bezierCurve) *bezierCurve {
	// 長さを取得
	L1 := bezier1.Length
	L2 := bezier2.Length

	// 重み付き平均で制御点を計算
	P1_new := bezier1.P1.MuledScalar(L1).Added(bezier2.P1.MuledScalar(L2)).DivedScalar(L1 + L2)
	P2_new := bezier1.P2.MuledScalar(L1).Added(bezier2.P2.MuledScalar(L2)).DivedScalar(L1 + L2)

	// 新しいベジェ曲線を作成
	combinedBezier := &bezierCurve{
		P0: bezier1.P0,
		P1: *P1_new,
		P2: *P2_new,
		P3: bezier2.P3,
	}

	return combinedBezier
}

// G3連続性を考慮して制御点を微調整する関数
func adjustControlPointsForG3Continuity(curve *bezierCurve) {
	// エネルギー関数の最小化など、実装が複雑なためここでは仮実装とする
	// 必要に応じて最適化アルゴリズムを実装
}

// bezierCurveからCurveを作成する関数
func createCurveFromBezierCurve(bezier *bezierCurve) *Curve {
	// CurveのStartとEndに制御点P1とP2を設定
	return &Curve{
		Start: MVec2{
			X: bezier.P1.X,
			Y: bezier.P1.Y,
		},
		End: MVec2{
			X: bezier.P2.X,
			Y: bezier.P2.Y,
		},
	}
}
