package mmath

import (
	"math"
	"testing"
)

func TestDegToRad(t *testing.T) {
	tests := []struct {
		name     string
		deg      float64
		expected float64
	}{
		{"0度", 0, 0},
		{"90度", 90, math.Pi / 2},
		{"180度", 180, math.Pi},
		{"360度", 360, 2 * math.Pi},
		{"-90度", -90, -math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DegToRad(tt.deg)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("DegToRad(%v) = %v, want %v", tt.deg, result, tt.expected)
			}
		})
	}
}

func TestRadToDeg(t *testing.T) {
	tests := []struct {
		name     string
		rad      float64
		expected float64
	}{
		{"0ラジアン", 0, 0},
		{"π/2ラジアン", math.Pi / 2, 90},
		{"πラジアン", math.Pi, 180},
		{"2πラジアン", 2 * math.Pi, 360},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RadToDeg(tt.rad)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("RadToDeg(%v) = %v, want %v", tt.rad, result, tt.expected)
			}
		})
	}
}

func TestNearEquals(t *testing.T) {
	tests := []struct {
		name     string
		v1       float64
		v2       float64
		epsilon  float64
		expected bool
	}{
		{"同じ値", 1.0, 1.0, 1e-10, true},
		{"わずかに異なる", 1.0, 1.0000000001, 1e-9, true},
		{"異なる値", 1.0, 1.1, 1e-10, false},
		{"epsilon内", 1.0, 1.05, 0.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NearEquals(tt.v1, tt.v2, tt.epsilon)
			if result != tt.expected {
				t.Errorf("NearEquals(%v, %v, %v) = %v, want %v", tt.v1, tt.v2, tt.epsilon, result, tt.expected)
			}
		})
	}
}

func TestClamped(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		min      float64
		max      float64
		expected float64
	}{
		{"範囲内", 5, 0, 10, 5},
		{"下限以下", -5, 0, 10, 0},
		{"上限以上", 15, 0, 10, 10},
		{"下限と同じ", 0, 0, 10, 0},
		{"上限と同じ", 10, 0, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamped(tt.v, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("Clamped(%v, %v, %v) = %v, want %v", tt.v, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

func TestClamped01(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		expected float64
	}{
		{"範囲内", 0.5, 0.5},
		{"下限以下", -0.5, 0},
		{"上限以上", 1.5, 1},
		{"0", 0, 0},
		{"1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamped01(tt.v)
			if result != tt.expected {
				t.Errorf("Clamped01(%v) = %v, want %v", tt.v, result, tt.expected)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		name     string
		v1       float64
		v2       float64
		t        float64
		expected float64
	}{
		{"t=0", 0, 10, 0, 0},
		{"t=1", 0, 10, 1, 10},
		{"t=0.5", 0, 10, 0.5, 5},
		{"t=0.25", 0, 100, 0.25, 25},
		{"負のt", 0, 10, -0.5, 0},
		{"1を超えるt", 0, 10, 1.5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Lerp(tt.v1, tt.v2, tt.t)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Lerp(%v, %v, %v) = %v, want %v", tt.v1, tt.v2, tt.t, result, tt.expected)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"空配列", []float64{}, 0},
		{"単一要素", []float64{5}, 5},
		{"複数要素", []float64{1, 2, 3, 4, 5}, 15},
		{"負の値含む", []float64{-1, 2, -3, 4}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.values)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Sum(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"空配列", []float64{}, 0},
		{"単一要素", []float64{5}, 5},
		{"複数要素", []float64{1, 2, 3, 4, 5}, 3},
		{"小数", []float64{1.5, 2.5, 3.0}, 7.0 / 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mean(tt.values)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Mean(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"空配列", []float64{}, 0},
		{"単一要素", []float64{5}, 5},
		{"複数要素", []float64{1, 5, 3, 2, 4}, 5},
		{"負の値", []float64{-5, -1, -3}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.values)
			if result != tt.expected {
				t.Errorf("Max(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"空配列", []float64{}, 0},
		{"単一要素", []float64{5}, 5},
		{"複数要素", []float64{5, 1, 3, 2, 4}, 1},
		{"負の値", []float64{-5, -1, -3}, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.values)
			if result != tt.expected {
				t.Errorf("Min(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		epsilon  float64
		expected float64
	}{
		{"ゼロ", 0, 1e-10, 0},
		{"epsilon未満", 1e-11, 1e-10, 0},
		{"epsilon以上", 1e-9, 1e-10, 1e-9},
		{"負のepsilon未満", -1e-11, 1e-10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.v, tt.epsilon)
			if result != tt.expected {
				t.Errorf("Truncate(%v, %v) = %v, want %v", tt.v, tt.epsilon, result, tt.expected)
			}
		})
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		expected float64
	}{
		{"正", 5, 1},
		{"負", -5, -1},
		{"ゼロ", 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sign(tt.v)
			if result != tt.expected {
				t.Errorf("Sign(%v) = %v, want %v", tt.v, result, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected []int
	}{
		{"空配列", []int{}, []int{}},
		{"重複なし", []int{1, 2, 3}, []int{1, 2, 3}},
		{"重複あり", []int{1, 2, 2, 3, 3, 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.values)
			if len(result) != len(tt.expected) {
				t.Errorf("Unique(%v) length = %v, want %v", tt.values, len(result), len(tt.expected))
			}
		})
	}
}

func TestBoolToInt(t *testing.T) {
	if BoolToInt(true) != 1 {
		t.Errorf("BoolToInt(true) = %v, want 1", BoolToInt(true))
	}
	if BoolToInt(false) != 0 {
		t.Errorf("BoolToInt(false) = %v, want 0", BoolToInt(false))
	}
}

func TestBoolToFlag(t *testing.T) {
	if BoolToFlag(true) != 1.0 {
		t.Errorf("BoolToFlag(true) = %v, want 1.0", BoolToFlag(true))
	}
	if BoolToFlag(false) != -1.0 {
		t.Errorf("BoolToFlag(false) = %v, want -1.0", BoolToFlag(false))
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{"1", 1, true},
		{"2", 2, true},
		{"4", 4, true},
		{"8", 8, true},
		{"16", 16, true},
		{"3", 3, false},
		{"5", 5, false},
		{"0", 0, false},
		{"-1", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPowerOfTwo(tt.n)
			if result != tt.expected {
				t.Errorf("IsPowerOfTwo(%v) = %v, want %v", tt.n, result, tt.expected)
			}
		})
	}
}
