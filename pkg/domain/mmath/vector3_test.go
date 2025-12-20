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
		{
			name:     "既存テスト1",
			v:        NewVec3ByValues(1, 2, 3),
			expected: 14.0,
		},
		{
			name:     "既存テスト2",
			v:        NewVec3ByValues(2.3, 0.2, 9),
			expected: 86.33000000000001,
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
		{
			name:     "既存テスト",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(4, 5, 6),
			t:        0.5,
			expected: NewVec3ByValues(2.5, 3.5, 4.5),
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

func TestVec3_NearEquals(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		epsilon  float64
		expected bool
	}{
		{
			name:     "ほぼ等しい",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(1.000001, 2.000001, 3.000001),
			epsilon:  0.00001,
			expected: true,
		},
		{
			name:     "等しくない",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(1.0001, 2.0001, 3.0001),
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

func TestVec3_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected bool
	}{
		{
			name:     "小さい",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(4, 5, 6),
			expected: true,
		},
		{
			name:     "大きい",
			v1:       NewVec3ByValues(3, 4, 5),
			v2:       NewVec3ByValues(1, 2, 3),
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

func TestVec3_GreaterThan(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected bool
	}{
		{
			name:     "大きい",
			v1:       NewVec3ByValues(3, 4, 5),
			v2:       NewVec3ByValues(1, 2, 3),
			expected: true,
		},
		{
			name:     "小さい",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(3, 4, 5),
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

func TestVec3_Negated(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "正の値",
			v:        NewVec3ByValues(1, 2, 3),
			expected: NewVec3ByValues(-1, -2, -3),
		},
		{
			name:     "負の値",
			v:        NewVec3ByValues(-3, -4, -5),
			expected: NewVec3ByValues(3, 4, 5),
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

func TestVec3_Length(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected float64
	}{
		{
			name:     "既存テスト1",
			v:        NewVec3ByValues(1, 2, 3),
			expected: 3.7416573867739413,
		},
		{
			name:     "既存テスト2",
			v:        NewVec3ByValues(2.3, 0.2, 9),
			expected: 9.291393867445294,
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

func TestVec3_Normalized(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "既存テスト1",
			v:        NewVec3ByValues(1, 2, 3),
			expected: NewVec3ByValues(0.2672612419124244, 0.5345224838248488, 0.8017837257372732),
		},
		{
			name:     "既存テスト2",
			v:        NewVec3ByValues(2.3, 0.2, 9),
			expected: NewVec3ByValues(0.24754089997827142, 0.021525295650284472, 0.9686383042628013),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Normalized()
			if !result.NearEquals(tt.expected, 1e-8) {
				t.Errorf("Normalized() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Distance(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected float64
	}{
		{
			name:     "既存テスト1",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(2.3, 0.2, 9),
			expected: 6.397655820689325,
		},
		{
			name:     "既存テスト2",
			v1:       NewVec3ByValues(-1, -0.3, 3),
			v2:       NewVec3ByValues(-2.3, 0.2, 9.33333333),
			expected: 6.484682804030502,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Distance(tt.v2)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Distance() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_One(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "非ゼロ",
			v:        NewVec3ByValues(1, 2, 3.2),
			expected: NewVec3ByValues(1, 2, 3.2),
		},
		{
			name:     "X=0",
			v:        NewVec3ByValues(0, 2, 3.2),
			expected: NewVec3ByValues(1, 2, 3.2),
		},
		{
			name:     "Y=0",
			v:        NewVec3ByValues(1, 0, 3.2),
			expected: NewVec3ByValues(1, 1, 3.2),
		},
		{
			name:     "Y=0,Z=0",
			v:        NewVec3ByValues(2, 0, 0),
			expected: NewVec3ByValues(2, 1, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.One()
			if !result.NearEquals(tt.expected, 1e-8) {
				t.Errorf("One() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Degree(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected float64
	}{
		{
			name:     "90度",
			v1:       NewVec3ByValues(1, 0, 0),
			v2:       NewVec3ByValues(0, 1, 0),
			expected: 90.0,
		},
		{
			name:     "既存テスト",
			v1:       NewVec3ByValues(1, 0, 0),
			v2:       NewVec3ByValues(1, 1, 1),
			expected: 54.735610317245346,
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

func TestVec3_Hash(t *testing.T) {
	v := NewVec3ByValues(1, 2, 3)
	expected := uint64(17648364615301650315)
	result := v.Hash()
	if result != expected {
		t.Errorf("Hash() = %v, want %v", result, expected)
	}
}
