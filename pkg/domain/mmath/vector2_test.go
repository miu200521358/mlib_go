package mmath

import (
	"math"
	"testing"
)

func TestVec2_LengthSquared(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec2
		expected float64
	}{
		{
			name:     "ゼロベクトル",
			v:        NewVec2(),
			expected: 0,
		},
		{
			name:     "単位ベクトル",
			v:        &Vec2{X: 1, Y: 0},
			expected: 1,
		},
		{
			name:     "3-4-5三角形",
			v:        &Vec2{X: 3, Y: 4},
			expected: 25, // 3^2 + 4^2 = 9 + 16 = 25
		},
		{
			name:     "負の値",
			v:        &Vec2{X: -3, Y: -4},
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.LengthSquared()
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("LengthSquared() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec2
		expected float64 // 正規化後の長さ
	}{
		{
			name:     "単位ベクトル",
			v:        &Vec2{X: 1, Y: 0},
			expected: 1,
		},
		{
			name:     "通常ベクトル",
			v:        &Vec2{X: 3, Y: 4},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.v.Copy()
			v.Normalize()
			length := v.Length()
			if !NearEquals(length, tt.expected, 1e-10) {
				t.Errorf("Normalize() length = %v, want %v", length, tt.expected)
			}
		})
	}
}

func TestVec2_Subed(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected *Vec2
	}{
		{
			name:     "基本的な減算",
			v1:       &Vec2{X: 5, Y: 3},
			v2:       &Vec2{X: 2, Y: 1},
			expected: &Vec2{X: 3, Y: 2},
		},
		{
			name:     "負の結果",
			v1:       &Vec2{X: 1, Y: 1},
			v2:       &Vec2{X: 2, Y: 3},
			expected: &Vec2{X: -1, Y: -2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Subed(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Subed() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Muled(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected *Vec2
	}{
		{
			name:     "基本的な乗算",
			v1:       &Vec2{X: 2, Y: 3},
			v2:       &Vec2{X: 4, Y: 5},
			expected: &Vec2{X: 8, Y: 15},
		},
		{
			name:     "ゼロ乗算",
			v1:       &Vec2{X: 2, Y: 3},
			v2:       NewVec2(),
			expected: NewVec2(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Muled(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Muled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Dived(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected *Vec2
	}{
		{
			name:     "基本的な除算",
			v1:       &Vec2{X: 8, Y: 15},
			v2:       &Vec2{X: 4, Y: 5},
			expected: &Vec2{X: 2, Y: 3},
		},
		{
			name:     "1で除算",
			v1:       &Vec2{X: 2, Y: 3},
			v2:       &Vec2{X: 1, Y: 1},
			expected: &Vec2{X: 2, Y: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Dived(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Dived() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Absed(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec2
		expected *Vec2
	}{
		{
			name:     "負の値",
			v:        &Vec2{X: -3, Y: -4},
			expected: &Vec2{X: 3, Y: 4},
		},
		{
			name:     "混合",
			v:        &Vec2{X: -1, Y: 2},
			expected: &Vec2{X: 1, Y: 2},
		},
		{
			name:     "正の値",
			v:        &Vec2{X: 3, Y: 4},
			expected: &Vec2{X: 3, Y: 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Absed()
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Absed() = %v, want %v", result, tt.expected)
			}
			// 元のベクトルが変更されていないことを確認
			if tt.v.X > 0 && result.X != tt.v.X {
				t.Errorf("Absed() should not modify original vector")
			}
		})
	}
}

func TestVec2_Dot(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected float64
	}{
		{
			name:     "直交ベクトル",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 0, Y: 1},
			expected: 0,
		},
		{
			name:     "同じ方向",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 1, Y: 0},
			expected: 1,
		},
		{
			name:     "一般的な場合",
			v1:       &Vec2{X: 2, Y: 3},
			v2:       &Vec2{X: 4, Y: 5},
			expected: 23, // 2*4 + 3*5 = 8 + 15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Dot(tt.v2)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Dot() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Angle(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected float64 // ラジアン
	}{
		{
			name:     "直交ベクトル",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 0, Y: 1},
			expected: math.Pi / 2,
		},
		{
			name:     "同じ方向",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 1, Y: 0},
			expected: 0,
		},
		{
			name:     "反対方向",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: -1, Y: 0},
			expected: math.Pi,
		},
		{
			name:     "45度",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 1, Y: 1},
			expected: math.Pi / 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Angle(tt.v2)
			if !NearEquals(result, tt.expected, 1e-6) {
				t.Errorf("Angle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_NearEquals(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		epsilon  float64
		expected bool
	}{
		{
			name:     "ほぼ等しい",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 1.000001, Y: 2.000001},
			epsilon:  0.00001,
			expected: true,
		},
		{
			name:     "等しくない",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 1.0001, Y: 2.0001},
			epsilon:  0.00001,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.NearEquals(tt.v2, tt.epsilon)
			if result != tt.expected {
				t.Errorf("NearEquals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected bool
	}{
		{
			name:     "小さい",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 3, Y: 4},
			expected: true,
		},
		{
			name:     "大きい",
			v1:       &Vec2{X: 3, Y: 4},
			v2:       &Vec2{X: 1, Y: 2},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.LessThan(tt.v2)
			if result != tt.expected {
				t.Errorf("LessThan() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_GreaterThan(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected bool
	}{
		{
			name:     "大きい",
			v1:       &Vec2{X: 3, Y: 4},
			v2:       &Vec2{X: 1, Y: 2},
			expected: true,
		},
		{
			name:     "小さい",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 3, Y: 4},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.GreaterThan(tt.v2)
			if result != tt.expected {
				t.Errorf("GreaterThan() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Negated(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec2
		expected *Vec2
	}{
		{
			name:     "正の値",
			v:        &Vec2{X: 1, Y: 2},
			expected: &Vec2{X: -1, Y: -2},
		},
		{
			name:     "負の値",
			v:        &Vec2{X: -3, Y: -4},
			expected: &Vec2{X: 3, Y: 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Negated()
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Negated() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Lerp(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		t        float64
		expected *Vec2
	}{
		{
			name:     "t=0.5",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 3, Y: 4},
			t:        0.5,
			expected: &Vec2{X: 2, Y: 3},
		},
		{
			name:     "t=0",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 3, Y: 4},
			t:        0,
			expected: &Vec2{X: 1, Y: 2},
		},
		{
			name:     "t=1",
			v1:       &Vec2{X: 1, Y: 2},
			v2:       &Vec2{X: 3, Y: 4},
			t:        1,
			expected: &Vec2{X: 3, Y: 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Lerp(tt.v2, tt.t)
			if !result.NearEquals(tt.expected, 1e-8) {
				t.Errorf("Lerp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Degree(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec2
		v2       *Vec2
		expected float64
	}{
		{
			name:     "90度",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 0, Y: 1},
			expected: 90.0,
		},
		{
			name:     "45度",
			v1:       &Vec2{X: 1, Y: 0},
			v2:       &Vec2{X: 1, Y: 1},
			expected: 45.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Degree(tt.v2)
			if !NearEquals(result, tt.expected, 0.00001) {
				t.Errorf("Degree() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec2_Hash(t *testing.T) {
	v := &Vec2{X: 1, Y: 2}
	expected := uint64(4921663092573786862)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash() = %v, want %v", result, expected)
	}
}
