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

// FindInflectionFrames は、与えられた値の変曲点を探す
// frames は、フレーム番号の配列
// values は、値の配列 (framesと同じ長さ)
// tolerance は、許容値
func FindInflectionFrames(frames []float32, values []float64) []float32 {
	// framesの間が空いている所は停止区域として扱う
	inflectionFrames := make([]float32, 0, len(frames))
	rangeFrames := make([]float32, 0, len(frames))
	rangeValues := make([]float64, 0, len(frames))

	for i, f := range frames {
		if i > 0 && f-frames[i-1] > 1 {
			// 前回のフレームとの間に空白がある場合、停止区域もしくは単調増加区域として扱う
			inflectionFrames = append(inflectionFrames, frames[i-1], f)
			inflectionFrames = appendInflections(rangeFrames, rangeValues, inflectionFrames)

			// 空白がある場合、区間値リストを初期化
			rangeFrames = make([]float32, 0, len(frames))
			rangeValues = make([]float64, 0, len(frames))
		}

		rangeFrames = append(rangeFrames, f)
		rangeValues = append(rangeValues, values[i])
	}

	inflectionFrames = appendInflections(rangeFrames, rangeValues, inflectionFrames)

	if len(frames) > 0 {
		inflectionFrames = append(inflectionFrames, frames[0], frames[len(frames)-1])
	}

	// 重複を削除
	inflectionFrames = UniqueFloat32s(inflectionFrames)

	SortFloat32s(inflectionFrames)

	return inflectionFrames
}

func appendInflections(rangeFrames []float32, rangeValues []float64, inflectionFrames []float32) []float32 {

	if len(rangeValues) > 1 {
		// 2つ以上ある場合、区間値の変曲点を探す
		ysPrime := gradient(rangeValues)

		for j, v := range ysPrime {
			if j > 0 && ysPrime[j-1]*v < 0 {
				// 前回と符号が変わっている場合、変曲点と見なす
				inflectionFrames = append(inflectionFrames, rangeFrames[j])
			}
		}
	}

	// 区間の最初と最後もなければ追加
	if len(rangeFrames) > 0 {
		inflectionFrames = append(inflectionFrames, rangeFrames[0], rangeFrames[len(rangeFrames)-1])
	}

	return inflectionFrames
}
