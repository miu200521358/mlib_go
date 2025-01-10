package mmath

// Gradient computes the numerical gradient of a 1D array.
func Gradient(data []float64, dx float64) []float64 {
	n := len(data)
	grad := make([]float64, n)

	// Forward difference for the first element
	grad[0] = (data[1] - data[0]) / dx

	// Central difference for the middle elements
	for i := 1; i < n-1; i++ {
		grad[i] = (data[i+1] - data[i-1]) / (2 * dx)
	}

	// Backward difference for the last element
	grad[n-1] = (data[n-1] - data[n-2]) / dx

	return grad
}

// FindInflectionFrames は、与えられた値の増減の切り替わるフレーム番号を探す
// frames は、フレーム番号の配列
// values は、値の配列 (framesと同じ長さ)
func FindInflectionFrames(frames []float32, values []float64, threshold float64) (inflectionFrames []float32) {
	if len(frames) <= 2 || len(values) <= 2 {
		return frames
	}

	inflectionFrames = []float32{frames[0]}

	// 増減の確認
	for i := 2; i < len(values); i++ {
		// 現在の値と前の値の差分を計算
		delta := values[i] - values[i-1]

		// 増加から減少、または減少から増加に切り替わったか確認
		if (delta > threshold && values[i-1] < values[i-2]) || (delta < -threshold && values[i-1] > values[i-2]) {
			inflectionFrames = append(inflectionFrames, frames[i-1])
		}
	}

	// 1次導関数を計算
	firstDerivative := Gradient(values, 1)

	// 2次導関数を計算
	secondDerivative := Gradient(firstDerivative, 1)

	// 符号の変化を確認
	for i := 1; i < len(secondDerivative); i++ {
		d1 := Round(secondDerivative[i-1], threshold)
		d2 := Round(secondDerivative[i], threshold)
		if d1*d2 < 0 || (d1 == 0 && d2 < 0) || (d1 < 0 && d2 == 0) {
			inflectionFrames = append(inflectionFrames, frames[i])
		}
	}

	// 最後のフレームを追加
	if len(frames) > 1 && !Contains(inflectionFrames, frames[len(frames)-1]) {
		inflectionFrames = append(inflectionFrames, frames[len(frames)-1])
	}

	inflectionFrames = Unique(inflectionFrames)
	Sort(inflectionFrames)

	return inflectionFrames
}
