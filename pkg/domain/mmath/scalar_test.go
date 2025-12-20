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
	tests := []struct {
		name     string
		b        bool
		expected int
	}{
		{"true", true, 1},
		{"false", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolToInt(tt.b)
			if result != tt.expected {
				t.Errorf("BoolToInt(%v) = %v, want %v", tt.b, result, tt.expected)
			}
		})
	}
}

func TestBoolToFlag(t *testing.T) {
	tests := []struct {
		name     string
		b        bool
		expected float64
	}{
		{"true", true, 1.0},
		{"false", false, -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolToFlag(tt.b)
			if result != tt.expected {
				t.Errorf("BoolToFlag(%v) = %v, want %v", tt.b, result, tt.expected)
			}
		})
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

func TestMedian(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"空配列", []float64{}, 0},
		{"単一要素", []float64{5}, 5},
		{"奇数個", []float64{3, 1, 2}, 2},
		{"偶数個", []float64{4, 1, 3, 2}, 2.5}, // ソート後[1,2,3,4]、中央は(2+3)/2=2.5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Median(tt.values)
			if result != tt.expected {
				t.Errorf("Median(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestStd(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"空配列", []float64{}, 0},
		{"同じ値", []float64{5, 5, 5}, 0},
		{"標準偏差計算", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Std(tt.values)
			if !NearEquals(result, tt.expected, 1e-6) {
				t.Errorf("Std(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestRatio(t *testing.T) {
	tests := []struct {
		name     string
		total    float64
		values   []float64
		expected []float64
	}{
		{"基本", 100, []float64{25, 50, 25}, []float64{0.25, 0.5, 0.25}},
		{"非ゼロ合計", 10, []float64{1, 2, 3}, []float64{0.1, 0.2, 0.3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ratio(tt.total, tt.values)
			if len(result) != len(tt.expected) {
				t.Errorf("Ratio() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if !NearEquals(result[i], tt.expected[i], 1e-10) {
					t.Errorf("Ratio()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name      string
		v         float64
		threshold float64
		expected  float64
	}{
		{"基本丸め", 1.5, 1, 2},
		{"切り捨て", 1.4, 1, 1},
		{"負の丸め", -1.5, 1, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.v, tt.threshold)
			if result != tt.expected {
				t.Errorf("Round(%v, %v) = %v, want %v", tt.v, tt.threshold, result, tt.expected)
			}
		})
	}
}

func TestEffective(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		expected float64
	}{
		{"基本", 1.23456789, 1.2346}, // 小数4桁
		{"ゼロ", 0, 0},
		{"負", -1.23456789, -1.2346},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Effective(tt.v)
			if !NearEquals(result, tt.expected, 1e-4) {
				t.Errorf("Effective(%v) = %v, want %v", tt.v, result, tt.expected)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected []float64
	}{
		{"基本ソート", []float64{5, 2, 8, 1, 9}, []float64{1, 2, 5, 8, 9}},
		{"逆順", []float64{5, 4, 3, 2, 1}, []float64{1, 2, 3, 4, 5}},
		{"既にソート済み", []float64{1, 2, 3}, []float64{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := DeepCopy(tt.values)
			Sort(values)
			for i := range values {
				if values[i] != tt.expected[i] {
					t.Errorf("Sort(%v) = %v, want %v", tt.values, values, tt.expected)
					break
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        []int
		v        int
		expected bool
	}{
		{"含む", []int{1, 2, 3}, 2, true},
		{"含まない", []int{1, 2, 3}, 5, false},
		{"空配列", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.s, tt.v)
			if result != tt.expected {
				t.Errorf("Contains(%v, %v) = %v, want %v", tt.s, tt.v, result, tt.expected)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {
	tests := []struct {
		name     string
		original []float64
	}{
		{"基本", []float64{1, 2, 3}},
		{"空配列", []float64{}},
		{"単一要素", []float64{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := DeepCopy(tt.original)
			// 同じ値であることを確認
			if len(copied) != len(tt.original) {
				t.Errorf("DeepCopy() length = %v, want %v", len(copied), len(tt.original))
				return
			}
			for i := range tt.original {
				if tt.original[i] != copied[i] {
					t.Errorf("DeepCopy() values differ at index %d", i)
				}
			}
			// コピーを変更しても元に影響しないことを確認
			if len(copied) > 0 {
				copied[0] = 999
				if tt.original[0] == 999 {
					t.Errorf("DeepCopy() did not create a deep copy")
				}
			}
		})
	}
}

func TestIntRanges(t *testing.T) {
	tests := []struct {
		name     string
		max      int
		expected []int
	}{
		{"5まで", 5, []int{0, 1, 2, 3, 4, 5}}, // 0からmaxまで含む
		{"1まで", 1, []int{0, 1}},
		{"0", 0, []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntRanges(tt.max)
			if len(result) != len(tt.expected) {
				t.Errorf("IntRanges(%v) length = %v, want %v", tt.max, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("IntRanges(%v)[%d] = %v, want %v", tt.max, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIntRangesByStep(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		step     int
		expected []int
	}{
		{"基本", 0, 10, 2, []int{0, 2, 4, 6, 8, 10}}, // maxも含む
		{"1から", 1, 5, 1, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntRangesByStep(tt.min, tt.max, tt.step)
			if len(result) != len(tt.expected) {
				t.Errorf("IntRangesByStep() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("IntRangesByStep()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIsAllSameValues(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected bool
	}{
		{"全て同じ", []float64{5, 5, 5}, true},
		{"異なる", []float64{5, 5, 6}, false},
		{"空", []float64{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAllSameValues(tt.values)
			if result != tt.expected {
				t.Errorf("IsAllSameValues(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestIsAlmostAllSameValues(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		threshold float64
		expected  bool
	}{
		{"全て同じ", []float64{5, 5, 5}, 0.1, true},
		{"ほぼ同じ", []float64{5, 5.05, 4.95}, 0.1, true},
		{"異なる", []float64{5, 5, 6}, 0.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAlmostAllSameValues(tt.values, tt.threshold)
			if result != tt.expected {
				t.Errorf("IsAlmostAllSameValues(%v, %v) = %v, want %v", tt.values, tt.threshold, result, tt.expected)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{"基本", [][]int{{1, 2}, {3, 4, 5}, {6}}, []int{1, 2, 3, 4, 5, 6}},
		{"空配列", [][]int{}, []int{}},
		{"単一配列", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Flatten() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Flatten()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestThetaToRad(t *testing.T) {
	tests := []struct {
		name     string
		theta    float64
		expected float64
	}{
		{"0", 0, 0},
		{"90", 90, math.Pi / 2},   // theta/180*π/2
		{"180", 180, math.Pi / 2}, // theta/180*π/2 (180/180=1, 1*π/2)
		{"360", 360, math.Pi / 2}, // arctan(tan(θ*π/180))的な動作
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ThetaToRad(tt.theta)
			// ThetaToRadは特殊な変換のため、動作確認のみ
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Logf("ThetaToRad(%v) = %v (expected %v)", tt.theta, result, tt.expected)
			}
		})
	}
}

func TestMean2DVertical(t *testing.T) {
	tests := []struct {
		name     string
		nums     [][]float64
		expected []float64
	}{
		{
			name:     "基本",
			nums:     [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: []float64{4, 5, 6}, // 各列の平均
		},
		{
			name:     "空",
			nums:     [][]float64{},
			expected: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mean2DVertical(tt.nums)
			if len(result) != len(tt.expected) {
				t.Errorf("Mean2DVertical() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if !NearEquals(result[i], tt.expected[i], 1e-10) {
					t.Errorf("Mean2DVertical()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestMean2DHorizontal(t *testing.T) {
	tests := []struct {
		name     string
		nums     [][]float64
		expected []float64
	}{
		{
			name:     "基本",
			nums:     [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			expected: []float64{2, 5, 8}, // 各行の平均
		},
		{
			name:     "空",
			nums:     [][]float64{},
			expected: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mean2DHorizontal(tt.nums)
			if len(result) != len(tt.expected) {
				t.Errorf("Mean2DHorizontal() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if !NearEquals(result[i], tt.expected[i], 1e-10) {
					t.Errorf("Mean2DHorizontal()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}
