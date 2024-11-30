package mmath

func gradient(values []float64) []float64 {
	result := make([]float64, len(values))
	for i := 1; i < len(values)-1; i++ {
		result[i] = (values[i+1] - values[i-1]) / 2.0
	}
	result[0] = values[1] - values[0]
	result[len(values)-1] = values[len(values)-1] - values[len(values)-2]
	return result
}

// FindInflectionFrames は、与えられた値の変曲点を探す(重複あり、順不同)
// frames は、フレーム番号の配列
// values は、値の配列 (framesと同じ長さ)
func FindInflectionFrames(frames []float32, values []float64) []float32 {
	if len(frames) < 2 || len(values) < 2 {
		return frames
	}

	inflectionFrames := make([]float32, 0, len(frames))

	// 2つ以上ある場合、区間値の変曲点を探す
	ysPrime := gradient(values)

	for j, v := range ysPrime {
		if j > 0 && ysPrime[j-1]*v < 0 {
			// 前回と符号が変わっている場合、変曲点と見なす
			inflectionFrames = append(inflectionFrames, frames[j])
		}
	}

	// 区間の最初と最後もなければ追加
	inflectionFrames = append(inflectionFrames, frames[0], frames[len(frames)-1])

	return inflectionFrames
}
