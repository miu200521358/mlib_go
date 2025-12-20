package mmath

import (
	"math"
	"testing"
)

func TestVec3_LengthSquared(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected float64
	}{
		{
			name:     "ゼロベクトル",
			v:        NewVec3(),
			expected: 0,
		},
		{
			name:     "単位ベクトルX",
			v:        VEC3_UNIT_X.Copy(),
			expected: 1,
		},
		{
			name:     "3-4-5三角形を3Dに拡張",
			v:        NewVec3ByValues(1, 2, 2),
			expected: 9, // 1^2 + 2^2 + 2^2 = 1 + 4 + 4 = 9
		},
		{
			name:     "負の値",
			v:        NewVec3ByValues(-1, -2, -2),
			expected: 9,
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

func TestVec3_Cos(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected float64
	}{
		{
			name:     "同じ方向",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_X.Copy(),
			expected: 1.0,
		},
		{
			name:     "反対方向",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_X_NEG.Copy(),
			expected: -1.0,
		},
		{
			name:     "直交ベクトル",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_Y.Copy(),
			expected: 0.0,
		},
		{
			name:     "45度",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       NewVec3ByValues(1, 1, 0).Normalized(),
			expected: math.Cos(math.Pi / 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Cos(tt.v2)
			if !NearEquals(result, tt.expected, 1e-6) {
				t.Errorf("Cos() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Subed(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本的な減算",
			v1:       NewVec3ByValues(5, 3, 2),
			v2:       NewVec3ByValues(2, 1, 1),
			expected: NewVec3ByValues(3, 2, 1),
		},
		{
			name:     "負の結果",
			v1:       NewVec3ByValues(1, 1, 1),
			v2:       NewVec3ByValues(2, 3, 4),
			expected: NewVec3ByValues(-1, -2, -3),
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

func TestVec3_Muled(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本的な乗算",
			v1:       NewVec3ByValues(2, 3, 4),
			v2:       NewVec3ByValues(5, 6, 7),
			expected: NewVec3ByValues(10, 18, 28),
		},
		{
			name:     "ゼロ乗算",
			v1:       NewVec3ByValues(2, 3, 4),
			v2:       NewVec3(),
			expected: NewVec3(),
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

func TestVec3_Dived(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本的な除算",
			v1:       NewVec3ByValues(10, 18, 28),
			v2:       NewVec3ByValues(5, 6, 7),
			expected: NewVec3ByValues(2, 3, 4),
		},
		{
			name:     "1で除算",
			v1:       NewVec3ByValues(2, 3, 4),
			v2:       VEC3_ONE.Copy(),
			expected: NewVec3ByValues(2, 3, 4),
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

func TestVec3_Absed(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "負の値",
			v:        NewVec3ByValues(-3, -4, -5),
			expected: NewVec3ByValues(3, 4, 5),
		},
		{
			name:     "混合",
			v:        NewVec3ByValues(-1, 2, -3),
			expected: NewVec3ByValues(1, 2, 3),
		},
		{
			name:     "正の値",
			v:        NewVec3ByValues(3, 4, 5),
			expected: NewVec3ByValues(3, 4, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Absed()
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Absed() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Cross(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "X x Y = Z",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_Y.Copy(),
			expected: VEC3_UNIT_Z.Copy(),
		},
		{
			name:     "Y x Z = X",
			v1:       VEC3_UNIT_Y.Copy(),
			v2:       VEC3_UNIT_Z.Copy(),
			expected: VEC3_UNIT_X.Copy(),
		},
		{
			name:     "Z x X = Y",
			v1:       VEC3_UNIT_Z.Copy(),
			v2:       VEC3_UNIT_X.Copy(),
			expected: VEC3_UNIT_Y.Copy(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Cross(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Cross() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Angle(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected float64 // ラジアン
	}{
		{
			name:     "直交ベクトル",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_Y.Copy(),
			expected: math.Pi / 2,
		},
		{
			name:     "同じ方向",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_X.Copy(),
			expected: 0,
		},
		{
			name:     "反対方向",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_X_NEG.Copy(),
			expected: math.Pi,
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

func TestVec3_Lerp(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		t        float64
		expected *Vec3
	}{
		{
			name:     "t=0",
			v1:       NewVec3ByValues(0, 0, 0),
			v2:       NewVec3ByValues(10, 10, 10),
			t:        0,
			expected: NewVec3ByValues(0, 0, 0),
		},
		{
			name:     "t=1",
			v1:       NewVec3ByValues(0, 0, 0),
			v2:       NewVec3ByValues(10, 10, 10),
			t:        1,
			expected: NewVec3ByValues(10, 10, 10),
		},
		{
			name:     "t=0.5",
			v1:       NewVec3ByValues(0, 0, 0),
			v2:       NewVec3ByValues(10, 10, 10),
			t:        0.5,
			expected: NewVec3ByValues(5, 5, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Lerp(tt.v2, tt.t)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Lerp() = %v, want %v", result, tt.expected)
			}
		})
	}
}
