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
		v1X      float64
		v1Y      float64
		v1Z      float64
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本的な減算",
			v1X:      5,
			v1Y:      3,
			v1Z:      2,
			v2:       NewVec3ByValues(2, 1, 1),
			expected: NewVec3ByValues(3, 2, 1),
		},
		{
			name:     "負の結果",
			v1X:      1,
			v1Y:      1,
			v1Z:      1,
			v2:       NewVec3ByValues(2, 3, 4),
			expected: NewVec3ByValues(-1, -2, -3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
			result := v1.Subed(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Subed() = %v, want %v", result, tt.expected)
			}
			// 元のベクトルが変更されていないことを確認
			if v1.X != tt.v1X || v1.Y != tt.v1Y || v1.Z != tt.v1Z {
				t.Errorf("Subed() should not modify original: got %v", v1)
			}
		})
	}
}

func TestVec3_Muled(t *testing.T) {
	tests := []struct {
		name     string
		v1X      float64
		v1Y      float64
		v1Z      float64
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本的な乗算",
			v1X:      2,
			v1Y:      3,
			v1Z:      4,
			v2:       NewVec3ByValues(5, 6, 7),
			expected: NewVec3ByValues(10, 18, 28),
		},
		{
			name:     "ゼロ乗算",
			v1X:      2,
			v1Y:      3,
			v1Z:      4,
			v2:       NewVec3(),
			expected: NewVec3(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
			result := v1.Muled(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Muled() = %v, want %v", result, tt.expected)
			}
			// 元のベクトルが変更されていないことを確認
			if v1.X != tt.v1X || v1.Y != tt.v1Y || v1.Z != tt.v1Z {
				t.Errorf("Muled() should not modify original: got %v", v1)
			}
		})
	}
}

func TestVec3_Dived(t *testing.T) {
	tests := []struct {
		name     string
		v1X      float64
		v1Y      float64
		v1Z      float64
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本的な除算",
			v1X:      10,
			v1Y:      18,
			v1Z:      28,
			v2:       NewVec3ByValues(5, 6, 7),
			expected: NewVec3ByValues(2, 3, 4),
		},
		{
			name:     "1で除算",
			v1X:      2,
			v1Y:      3,
			v1Z:      4,
			v2:       VEC3_ONE.Copy(),
			expected: NewVec3ByValues(2, 3, 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
			result := v1.Dived(tt.v2)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Dived() = %v, want %v", result, tt.expected)
			}
			// 元のベクトルが変更されていないことを確認
			if v1.X != tt.v1X || v1.Y != tt.v1Y || v1.Z != tt.v1Z {
				t.Errorf("Dived() should not modify original: got %v", v1)
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

// 破壊的メソッドのテスト

func TestVec3_Add(t *testing.T) {
	tests := []struct {
		name     string
		v1X      float64
		v1Y      float64
		v1Z      float64
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本加算",
			v1X:      1,
			v1Y:      2,
			v1Z:      3,
			v2:       NewVec3ByValues(4, 5, 6),
			expected: NewVec3ByValues(5, 7, 9),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
			result := v.Add(tt.v2)
			// 結果の確認
			if !v.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Add() = %v, want %v", v, tt.expected)
			}
			// 元のオブジェクトが変更されていることを確認
			if v.X != tt.expected.X || v.Y != tt.expected.Y || v.Z != tt.expected.Z {
				t.Errorf("Add() should modify original: v = %v", v)
			}
			// 戻り値がレシーバ自身であることを確認（チェーン呼び出し用）
			if result != v {
				t.Errorf("Add() should return receiver: result = %p, v = %p", result, v)
			}
		})
	}
}

func TestVec3_Sub(t *testing.T) {
	tests := []struct {
		name     string
		v1X      float64
		v1Y      float64
		v1Z      float64
		v2       *Vec3
		expected *Vec3
	}{
		{
			name:     "基本減算",
			v1X:      5,
			v1Y:      7,
			v1Z:      9,
			v2:       NewVec3ByValues(2, 3, 4),
			expected: NewVec3ByValues(3, 4, 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVec3ByValues(tt.v1X, tt.v1Y, tt.v1Z)
			result := v.Sub(tt.v2)
			// 結果の確認
			if !v.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Sub() = %v, want %v", v, tt.expected)
			}
			// 戻り値がレシーバ自身であることを確認
			if result != v {
				t.Errorf("Sub() should return receiver: result = %p, v = %p", result, v)
			}
		})
	}
}

// Clamp系テスト

func TestVec3_Clamped(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		min      *Vec3
		max      *Vec3
		expected *Vec3
	}{
		{
			name:     "範囲内",
			v:        NewVec3ByValues(0.5, 0.5, 0.5),
			min:      NewVec3(),
			max:      VEC3_ONE.Copy(),
			expected: NewVec3ByValues(0.5, 0.5, 0.5),
		},
		{
			name:     "下限クランプ",
			v:        NewVec3ByValues(-1, -1, -1),
			min:      NewVec3(),
			max:      VEC3_ONE.Copy(),
			expected: NewVec3(),
		},
		{
			name:     "上限クランプ",
			v:        NewVec3ByValues(2, 2, 2),
			min:      NewVec3(),
			max:      VEC3_ONE.Copy(),
			expected: VEC3_ONE.Copy(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Clamped(tt.min, tt.max)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Clamped() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Clamped01(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "範囲内",
			v:        NewVec3ByValues(0.5, 0.5, 0.5),
			expected: NewVec3ByValues(0.5, 0.5, 0.5),
		},
		{
			name:     "下限クランプ",
			v:        NewVec3ByValues(-0.5, -0.5, -0.5),
			expected: NewVec3(),
		},
		{
			name:     "上限クランプ",
			v:        NewVec3ByValues(1.5, 1.5, 1.5),
			expected: VEC3_ONE.Copy(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Clamped01()
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("Clamped01() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ユーティリティ関数テスト

func TestVec3_Copy(t *testing.T) {
	tests := []struct {
		name string
		v    *Vec3
	}{
		{
			name: "基本コピー",
			v:    NewVec3ByValues(1, 2, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := tt.v.Copy()
			if !copied.NearEquals(tt.v, 1e-10) {
				t.Errorf("Copy() = %v, want %v", copied, tt.v)
			}
			// ディープコピー確認
			copied.X = 999
			if tt.v.X == 999 {
				t.Errorf("Copy() did not create deep copy")
			}
		})
	}
}

func TestVec3_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected bool
	}{
		{
			name:     "ゼロベクトル",
			v:        NewVec3(),
			expected: true,
		},
		{
			name:     "非ゼロ",
			v:        NewVec3ByValues(1, 0, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.IsZero()
			if result != tt.expected {
				t.Errorf("IsZero() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Dot(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		expected float64
	}{
		{
			name:     "直交ベクトル",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_Y.Copy(),
			expected: 0,
		},
		{
			name:     "同じ方向",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_X.Copy(),
			expected: 1,
		},
		{
			name:     "一般的な場合",
			v1:       NewVec3ByValues(1, 2, 3),
			v2:       NewVec3ByValues(4, 5, 6),
			expected: 32, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18
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

func TestVec3_Project(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		onto     *Vec3
		expected *Vec3
	}{
		{
			name:     "X軸への射影",
			v:        NewVec3ByValues(1, 1, 0),
			onto:     VEC3_UNIT_X.Copy(),
			expected: VEC3_UNIT_X.Copy(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Project(tt.onto)
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("Project() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_Slerp(t *testing.T) {
	tests := []struct {
		name     string
		v1       *Vec3
		v2       *Vec3
		t        float64
		expected *Vec3
	}{
		{
			name:     "t=0",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_Y.Copy(),
			t:        0,
			expected: VEC3_UNIT_X.Copy(),
		},
		{
			name:     "t=1",
			v1:       VEC3_UNIT_X.Copy(),
			v2:       VEC3_UNIT_Y.Copy(),
			t:        1,
			expected: VEC3_UNIT_Y.Copy(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Slerp(tt.v2, tt.t)
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("Slerp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_RadToDeg(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "基本変換",
			v:        NewVec3ByValues(math.Pi, math.Pi/2, 0),
			expected: NewVec3ByValues(180, 90, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.RadToDeg()
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("RadToDeg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVec3_DegToRad(t *testing.T) {
	tests := []struct {
		name     string
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "基本変換",
			v:        NewVec3ByValues(180, 90, 0),
			expected: NewVec3ByValues(math.Pi, math.Pi/2, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.DegToRad()
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("DegToRad() = %v, want %v", result, tt.expected)
			}
		})
	}
}
