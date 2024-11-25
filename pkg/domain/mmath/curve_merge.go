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
func CombineCurves(
	frames []float32, values []float64, curves []*Curve,
) (combinedFrames []float32, combinedCurves []*Curve) {
	combinedFrames = []float32{}
	combinedCurves = []*Curve{}

	if len(curves) < 2 {
		combinedFrames = frames
		combinedCurves = curves
		return
	}

	// 最初のフレームと曲線を追加
	combinedFrames = append(combinedFrames, frames[0])
	combinedCurves = append(combinedCurves, curves[0])

	i := 1
	for i < len(curves) {
		curve1 := curves[i-1]
		frameStart := frames[i-1]
		valueStart := values[i-1]

		curve2 := curves[i]
		frameMid := frames[i]
		valueMid := values[i]

		// 次の曲線が存在するか確認
		if i+1 < len(frames) {
			frameNext := frames[i+1]
			valueNext := values[i+1]

			// ベジェ曲線を実際の値に合わせて調整
			bezier1 := adjustCurveToValues(curve1, frameStart, frameMid, valueStart, valueMid)
			bezier2 := adjustCurveToValues(curve2, frameMid, frameNext, valueMid, valueNext)

			// 曲線の長さを計算
			length1 := bezier1.Length
			length2 := bezier2.Length

			// 結合可能かどうかを判定
			canCombine := checkIfCurvesCanBeCombined(bezier1, bezier2, length1, length2)

			if canCombine {
				// 曲線を結合
				combinedBezier := combineTwoCurvesWithoutOptimization(bezier1, bezier2)

				// 正規化してCurve構造体に格納
				combinedCurve := createCurveFromBezier(combinedBezier)

				// 結合された曲線をリストに追加
				combinedFrames = append(combinedFrames, frameNext)
				combinedCurves = append(combinedCurves, combinedCurve)

				// インデックスを2つ進める
				i += 2
				continue
			}
		}

		// 結合できない場合、元の曲線を追加
		combinedFrames = append(combinedFrames, frameMid)
		combinedCurves = append(combinedCurves, curve2)

		// インデックスを1つ進める
		i++
	}

	// 最後のフレームが追加されていない場合は追加
	if combinedFrames[len(combinedFrames)-1] != frames[len(frames)-1] {
		combinedFrames = append(combinedFrames, frames[len(frames)-1])
	}

	return
}

// adjustCurveToValues は、ベジェ曲線を実際の値に合わせて制御点を調整します。
func adjustCurveToValues(curve *Curve, frameStart, frameEnd float32, valueStart, valueEnd float64) *bezierCurve {
	// P0とP3を設定
	P0 := MVec2{X: float64(frameStart), Y: valueStart}
	P3 := MVec2{X: float64(frameEnd), Y: valueEnd}

	// P1とP2を計算
	P1 := MVec2{
		X: P0.X + (curve.Start.X/CURVE_MAX)*(P3.X-P0.X),
		Y: P0.Y + (curve.Start.Y/CURVE_MAX)*(P3.Y-P0.Y),
	}
	P2 := MVec2{
		X: P0.X + (curve.End.X/CURVE_MAX)*(P3.X-P0.X),
		Y: P0.Y + (curve.End.Y/CURVE_MAX)*(P3.Y-P0.Y),
	}

	// ベジェ曲線を作成
	bezier := &bezierCurve{
		P0: P0,
		P1: P1,
		P2: P2,
		P3: P3,
	}

	// 曲線の長さを計算
	bezier.Length = computeCurveLength(bezier)

	return bezier
}

// computeCurveLength は、適応的シンプソン法を使用してベジェ曲線の長さを計算します。
func computeCurveLength(bezier *bezierCurve) float64 {
	// 適応的シンプソン法の初期パラメータ
	tol := 1e-5
	maxDepth := 20

	length := adaptiveSimpson(bezierLengthFunc(bezier), 0.0, 1.0, tol, maxDepth)

	return length
}

// bezierLengthFunc は、ベジェ曲線の微分の長さを返す関数を生成します。
func bezierLengthFunc(bezier *bezierCurve) func(t float64) float64 {
	return func(t float64) float64 {
		dx, dy := bezierDerivative(bezier, t)
		return math.Sqrt(dx*dx + dy*dy)
	}
}

// bezierDerivative は、ベジェ曲線の微分を計算します。
func bezierDerivative(bezier *bezierCurve, t float64) (dx, dy float64) {
	// ベジェ曲線の微分
	mt := 1 - t
	dx = 3*mt*mt*(bezier.P1.X-bezier.P0.X) + 6*mt*t*(bezier.P2.X-bezier.P1.X) + 3*t*t*(bezier.P3.X-bezier.P2.X)
	dy = 3*mt*mt*(bezier.P1.Y-bezier.P0.Y) + 6*mt*t*(bezier.P2.Y-bezier.P1.Y) + 3*t*t*(bezier.P3.Y-bezier.P2.Y)
	return
}

// adaptiveSimpson は、適応的シンプソン法で数値積分を行います。
func adaptiveSimpson(f func(t float64) float64, a, b, tol float64, maxDepth int) float64 {
	c := (a + b) / 2
	fa := f(a)
	fb := f(b)
	fc := f(c)
	return adaptiveSimpsonAux(f, a, b, tol, fa, fb, fc, maxDepth)
}

func adaptiveSimpsonAux(f func(t float64) float64, a, b, tol, fa, fb, fc float64, depth int) float64 {
	c := (a + b) / 2
	h := b - a
	d := (a + c) / 2
	e := (c + b) / 2
	fd := f(d)
	fe := f(e)
	Sleft := (h / 12) * (fa + 4*fd + fc)
	Sright := (h / 12) * (fc + 4*fe + fb)
	S := Sleft + Sright
	Sapprox := (h / 6) * (fa + 4*fc + fb)
	if depth <= 0 || math.Abs(S-Sapprox) <= 15*tol {
		return S + (S-Sapprox)/15
	}
	return adaptiveSimpsonAux(f, a, c, tol/2, fa, fc, fd, depth-1) + adaptiveSimpsonAux(f, c, b, tol/2, fc, fb, fe, depth-1)
}

// checkIfCurvesCanBeCombined は、2つのベジェ曲線が結合可能かどうかを判定します。
func checkIfCurvesCanBeCombined(bezier1, bezier2 *bezierCurve, length1, length2 float64) bool {
	// 長さと値の変化量の確認
	totalLength := length1 + length2
	combinedBezier := combineTwoCurvesWithoutOptimization(bezier1, bezier2)
	combinedLength := computeCurveLength(combinedBezier)

	// 相対誤差を計算
	relativeError := math.Abs(totalLength-combinedLength) / totalLength

	// 許容する相対誤差の閾値内であることを確認
	return relativeError <= 0.1
}

// combineTwoCurvesWithoutOptimization は、制御点の重み付き平均を使用してベジェ曲線を結合します（最適化なし）。
func combineTwoCurvesWithoutOptimization(bezier1, bezier2 *bezierCurve) *bezierCurve {
	// 重みは曲線の長さとする
	w1 := bezier1.Length
	w2 := bezier2.Length

	// 新しい制御点を計算
	P1 := MVec2{
		X: (w1*bezier1.P1.X + w2*bezier2.P1.X) / (w1 + w2),
		Y: (w1*bezier1.P1.Y + w2*bezier2.P1.Y) / (w1 + w2),
	}
	P2 := MVec2{
		X: (w1*bezier1.P2.X + w2*bezier2.P2.X) / (w1 + w2),
		Y: (w1*bezier1.P2.Y + w2*bezier2.P2.Y) / (w1 + w2),
	}

	combinedBezier := &bezierCurve{
		P0: bezier1.P0,
		P1: P1,
		P2: P2,
		P3: bezier2.P3,
	}

	return combinedBezier
}

func createCurveFromBezier(bezier *bezierCurve) *Curve {
	frameRange := bezier.P3.X - bezier.P0.X
	valueRange := bezier.P3.Y - bezier.P0.Y

	if frameRange == 0 {
		frameRange = 1e-6
	}
	if valueRange == 0 {
		valueRange = 1e-6
	}

	curve := &Curve{}

	curve.Start.X = math.Round(((bezier.P1.X - bezier.P0.X) / frameRange) * CURVE_MAX)
	curve.End.X = math.Round(((bezier.P2.X - bezier.P0.X) / frameRange) * CURVE_MAX)

	curve.Start.Y = math.Round(((bezier.P1.Y - bezier.P0.Y) / valueRange) * CURVE_MAX)
	curve.End.Y = math.Round(((bezier.P2.Y - bezier.P0.Y) / valueRange) * CURVE_MAX)

	if NearEquals(curve.Start.X, curve.Start.Y, 1e-6) && NearEquals(curve.End.X, curve.End.Y, 1e-6) {
		// 開始点と終了点が同じ場合、線形補間に変更
		curve.Start.X = 20.0
		curve.Start.Y = 20.0
		curve.End.X = 107.0
		curve.End.Y = 107.0
	}

	return curve
}
