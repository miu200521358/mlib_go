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
func FindInflectionFrames(frames []float32, values []float64, threshold float64) (inflectionFrames []float32) {
	if len(frames) < 2 || len(values) < 2 {
		return frames
	}

	inflectionFrames = make([]float32, 0, len(frames))
	inflectionFrames = append(inflectionFrames, frames[0])

	// 2つ以上ある場合、区間値の変曲点を探す
	grad := Gradient(values, 1)

	// 変曲点を見つける
	for i := range len(grad) {
		if i < 2 {
			continue
		}

		if ((grad[i-1] < 0 && grad[i] >= 0) || (grad[i-1] >= 0 && grad[i] < 0)) && math.Abs(grad[i-1]-grad[i]) > threshold {
			inflectionFrames = append(inflectionFrames, frames[i])
		} else if ((grad[i-2]-grad[i-1] < 0 && grad[i-1]-grad[i] >= 0) ||
			(grad[i-2]-grad[i-1] >= 0 && grad[i-1]-grad[i] < 0)) &&
			math.Abs(grad[i-2]-grad[i-1]) > threshold && math.Abs(grad[i-1]-grad[i]) > threshold {
			inflectionFrames = append(inflectionFrames, frames[i-1])
			// } else if math.Abs(math.Min(values[i-1], values[i])) > 0 &&
			// 	math.Abs(math.Max(values[i-1], values[i])/math.Min(values[i-1], values[i])) > vTol {
			// 	// 直前との差が大きい場合も対象
			// 	inflectionFrames = append(inflectionFrames, frames[i])
			// } else if math.Abs(values[i-1]-values[i]) > aTol {
			// 	// 直前との差が大きい場合も対象(直前が0の場合)
			// 	inflectionFrames = append(inflectionFrames, frames[i])
		}
	}

	if !Contains(inflectionFrames, frames[len(frames)-1]) {
		inflectionFrames = append(inflectionFrames, frames[len(frames)-1])
	}

	return inflectionFrames
}
