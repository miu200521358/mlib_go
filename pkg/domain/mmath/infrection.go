// 指示: miu200521358
package mmath

func Gradient(data []float64, dx float64) []float64 {
	n := len(data)
	grad := make([]float64, n)

	grad[0] = (data[1] - data[0]) / dx

	for i := 1; i < n-1; i++ {
		grad[i] = (data[i+1] - data[i-1]) / (2 * dx)
	}

	grad[n-1] = (data[n-1] - data[n-2]) / dx

	return grad
}

func FindInflectionFrames(frames []float32, values []float64, threshold float64) (inflectionFrames []float32) {
	if len(frames) <= 2 || len(values) <= 2 {
		return frames
	}

	inflectionFrames = []float32{frames[0]}

	for i := 2; i < len(values); i++ {
		delta := values[i] - values[i-1]

		if (delta > threshold && values[i-1] < values[i-2]) || (delta < -threshold && values[i-1] > values[i-2]) {
			inflectionFrames = append(inflectionFrames, frames[i-1])
		}
	}

	firstDerivative := Gradient(values, 1)

	secondDerivative := Gradient(firstDerivative, 1)

	for i := 1; i < len(secondDerivative); i++ {
		d1 := Round(secondDerivative[i-1], threshold)
		d2 := Round(secondDerivative[i], threshold)
		if d1*d2 < 0 || (d1 == 0 && d2 < 0) || (d1 < 0 && d2 == 0) {
			inflectionFrames = append(inflectionFrames, frames[i])
		}
	}

	if len(frames) > 1 && !Contains(inflectionFrames, frames[len(frames)-1]) {
		inflectionFrames = append(inflectionFrames, frames[len(frames)-1])
	}

	inflectionFrames = Unique(inflectionFrames)
	Sort(inflectionFrames)

	return inflectionFrames
}

