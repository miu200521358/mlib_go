package mmath

import (
	"testing"
)

func TestVec4_Add(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec4
		v2       *Vec4
		expected *Vec4
	}{
		{
			name:     "基本的な加算",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(5, 6, 7, 8),
			expected: NewVec4ByValues(6, 8, 10, 12),
		},
		{
			name:     "ゼロ加算",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4(),
			expected: NewVec4ByValues(1, 2, 3, 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.v1.Copy()
			result := v.Add(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec4_Added(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec4
		v2       *Vec4
		expected *Vec4
	}{
		{
			name:     "基本的な加算",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(5, 6, 7, 8),
			expected: NewVec4ByValues(6, 8, 10, 12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Added(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Added() = %v, want %v", result, tt.expected)
			}
			// 元のベクトルが変更されていないことを確認
			if tt.v1.X != 1 || tt.v1.Y != 2 {
				t.Errorf("Added() modified original vector")
			}
		})
	}
}

func TestVec4_Subed(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec4
		v2       *Vec4
		expected *Vec4
	}{
		{
			name:     "基本的な減算",
			v1:       NewVec4ByValues(5, 6, 7, 8),
			v2:       NewVec4ByValues(1, 2, 3, 4),
			expected: NewVec4ByValues(4, 4, 4, 4),
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

func TestVec4_MulScalar(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec4
		s        float64
		expected *Vec4
	}{
		{
			name:     "2倍",
			v:        NewVec4ByValues(1, 2, 3, 4),
			s:        2,
			expected: NewVec4ByValues(2, 4, 6, 8),
		},
		{
			name:     "ゼロ倍",
			v:        NewVec4ByValues(1, 2, 3, 4),
			s:        0,
			expected: NewVec4(),
		},
		{
			name:     "負の倍率",
			v:        NewVec4ByValues(1, 2, 3, 4),
			s:        -1,
			expected: NewVec4ByValues(-1, -2, -3, -4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.v.Copy()
			result := v.MulScalar(tt.s)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("MulScalar(%v) = %v, want %v", tt.s, result, tt.expected)
			}
		})
	}
}

func TestVec4_MuledScalar(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec4
		s        float64
		expected *Vec4
	}{
		{
			name:     "2倍",
			v:        NewVec4ByValues(1, 2, 3, 4),
			s:        2,
			expected: NewVec4ByValues(2, 4, 6, 8),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.MuledScalar(tt.s)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("MuledScalar(%v) = %v, want %v", tt.s, result, tt.expected)
			}
			// 元のベクトルが変更されていないことを確認
			if tt.v.X != 1 || tt.v.Y != 2 {
				t.Errorf("MuledScalar() modified original vector")
			}
		})
	}
}

func TestVec4_Length(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec4
		expected float64
	}{
		{
			name:     "ゼロベクトル",
			v:        NewVec4(),
			expected: 0,
		},
		{
			name:     "単位ベクトル",
			v:        NewVec4ByValues(1, 0, 0, 0),
			expected: 1,
		},
		{
			name:     "2-2-2-2ベクトル",
			v:        NewVec4ByValues(2, 2, 2, 2),
			expected: 4, // sqrt(4+4+4+4) = sqrt(16) = 4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Length()
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Length() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec4_Dot(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec4
		v2       *Vec4
		expected float64
	}{
		{
			name:     "直交ベクトル",
			v1:       NewVec4ByValues(1, 0, 0, 0),
			v2:       NewVec4ByValues(0, 1, 0, 0),
			expected: 0,
		},
		{
			name:     "同じベクトル",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(1, 2, 3, 4),
			expected: 30, // 1+4+9+16
		},
		{
			name:     "一般的な場合",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(5, 6, 7, 8),
			expected: 70, // 5+12+21+32
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

func TestVec4_NearEquals(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec4
		v2       *Vec4
		epsilon  float64
		expected bool
	}{
		{
			name:     "同じベクトル",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(1, 2, 3, 4),
			epsilon:  1e-10,
			expected: true,
		},
		{
			name:     "わずかに異なる",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(1.0000001, 2, 3, 4),
			epsilon:  1e-6,
			expected: true,
		},
		{
			name:     "異なるベクトル",
			v1:       NewVec4ByValues(1, 2, 3, 4),
			v2:       NewVec4ByValues(5, 6, 7, 8),
			epsilon:  1e-10,
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

func TestVec4_Copy(t *testing.T) {
	v := NewVec4ByValues(1, 2, 3, 4)
	copied := v.Copy()

	if !v.NearEquals(copied, 1e-10) {
		t.Errorf("Copy() returned different values")
	}

	// コピーを変更しても元のベクトルに影響しないことを確認
	copied.X = 100
	if v.X == 100 {
		t.Errorf("Copy() did not create a deep copy")
	}
}

func TestVec4_Vector(t *testing.T) {
	v := NewVec4ByValues(1, 2, 3, 4)
	result := v.Vector()

	if len(result) != 4 {
		t.Errorf("Vector() length = %v, want 4", len(result))
	}
	if result[0] != 1 || result[1] != 2 || result[2] != 3 || result[3] != 4 {
		t.Errorf("Vector() = %v, want [1, 2, 3, 4]", result)
	}
}

func TestVec4_GetXYZ(t *testing.T) {
	v := NewVec4ByValues(1, 2, 3, 4)
	result := v.GetXYZ()

	expected := NewVec3ByValues(1, 2, 3)
	if !result.NearEquals(expected, 1e-10) {
		t.Errorf("GetXYZ() = %v, want %v", result, expected)
	}
}
