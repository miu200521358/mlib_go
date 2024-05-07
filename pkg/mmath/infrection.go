package mmath

import (
	"math"
	"sort"
)

func gradient(values []float64) []float64 {
	result := make([]float64, len(values))
	for i := 1; i < len(values)-1; i++ {
		result[i] = (values[i+1] - values[i-1]) / 2.0
	}
	result[0] = values[1] - values[0]
	result[len(values)-1] = values[len(values)-1] - values[len(values)-2]
	return result
}

func clip(values []float64, min, max float64) []float64 {
	result := make([]float64, len(values))
	for i, v := range values {
		if v < min {
			result[i] = min
		} else if v > max {
			result[i] = max
		} else {
			result[i] = v
		}
	}
	return result
}

func diff(values []float64) []float64 {
	result := make([]float64, len(values)-1)
	for i := range result {
		result[i] = values[i+1] - values[i]
	}
	return result
}

func CreateInfections(values []float64, threshold, tolerance float64) []int {
	ysPrime := gradient(values)
	ysDoublePrime := gradient(ysPrime)

	var inflectionPoints []int
	for i, v := range ysDoublePrime {
		if math.Abs(v) >= threshold {
			inflectionPoints = append(inflectionPoints, i)
		}
	}

	ysDoublePrimeSign := make([]float64, len(ysDoublePrime))
	for i, v := range clip(ysDoublePrime, -tolerance, tolerance) {
		ysDoublePrimeSign[i] = math.Copysign(1, v)
	}

	ysDoublePrimeDiff := diff(ysDoublePrimeSign)

	var signChangeIndices []int
	for i, v := range ysDoublePrimeDiff {
		if v != 0 {
			signChangeIndices = append(signChangeIndices, i)
		}
	}

	inflections := append(inflectionPoints, signChangeIndices...)
	inflections = append(inflections, 0, len(values)-1)

	sort.Ints(inflections)

	return inflections
}
