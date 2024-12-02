package mmath

import "math"

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

// FindInflectionFrames は、与えられた値の変曲点を探す(重複あり、順不同)
// frames は、フレーム番号の配列
// values は、値の配列 (framesと同じ長さ)
func FindInflectionFrames(frames []float32, values []float64) []float32 {
	if len(frames) < 2 || len(values) < 2 {
		return frames
	}

	inflectionFrames := make([]float32, 0, len(frames))
	inflectionFrames = append(inflectionFrames, frames[0])

	// 2つ以上ある場合、区間値の変曲点を探す
	grad := Gradient(values, 1)

	// 変曲点を見つける
	deltaThreshold := 10.0
	for i := range len(grad) {
		if i > 0 && (grad[i-1]*grad[i] < 0 ||
			(values[i-1] == 0 && math.Abs(values[i]) > deltaThreshold) ||
			(values[i-1] != 0 && math.Abs(values[i]/values[i-1]) > deltaThreshold)) {
			inflectionFrames = append(inflectionFrames, frames[i])
		}
	}

	if !Contains(inflectionFrames, frames[len(frames)-1]) {
		inflectionFrames = append(inflectionFrames, frames[len(frames)-1])
	}

	return inflectionFrames
}
