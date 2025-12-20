package mmath

import (
	"math"
	"testing"
)

func TestMat4_NewMat4(t *testing.T) {
	m := NewMat4()

	// 単位行列であることを確認
	if !m.IsIdent() {
		t.Errorf("NewMat4() should return identity matrix")
	}
}

func TestMat4_At(t *testing.T) {
	m := NewMat4()

	// 対角要素が1であることを確認
	for i := 0; i < 4; i++ {
		if m.At(i, i) != 1 {
			t.Errorf("At(%d, %d) = %v, want 1", i, i, m.At(i, i))
		}
	}

	// 非対角要素が0であることを確認
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i != j && m.At(i, j) != 0 {
				t.Errorf("At(%d, %d) = %v, want 0", i, j, m.At(i, j))
			}
		}
	}
}

func TestMat4_Set(t *testing.T) {
	m := NewMat4()
	m.Set(0, 3, 5)

	if m.At(0, 3) != 5 {
		t.Errorf("Set() did not set value correctly")
	}
}

func TestMat4_Mul(t *testing.T) {
	// 単位行列との乗算
	m1 := NewMat4()
	m2 := NewMat4()
	result := m1.Muled(m2)

	if !result.IsIdent() {
		t.Errorf("Identity * Identity should be Identity")
	}
}

func TestMat4_Inverted(t *testing.T) {
	// 単位行列の逆行列は単位行列
	m := NewMat4()
	inv := m.Inverted()

	if !inv.IsIdent() {
		t.Errorf("Inverted identity should be identity")
	}

	// 逆行列と元の行列の積は単位行列
	m2 := NewMat4()
	m2.Translate(NewVec3ByValues(10, 20, 30))
	inv2 := m2.Inverted()
	product := m2.Muled(inv2)

	if !product.NearEquals(MAT4_IDENTITY, 1e-6) {
		t.Errorf("M * M^-1 should be identity, got %v", product)
	}
}

func TestMat4_Translation(t *testing.T) {
	m := NewMat4()
	m.Translate(NewVec3ByValues(10, 20, 30))

	trans := m.Translation()
	expected := NewVec3ByValues(10, 20, 30)

	if !trans.NearEquals(expected, 1e-10) {
		t.Errorf("Translation() = %v, want %v", trans, expected)
	}
}

func TestMat4_SetTranslation(t *testing.T) {
	m := NewMat4()
	m.SetTranslation(NewVec3ByValues(5, 10, 15))

	trans := m.Translation()
	expected := NewVec3ByValues(5, 10, 15)

	if !trans.NearEquals(expected, 1e-10) {
		t.Errorf("SetTranslation() result = %v, want %v", trans, expected)
	}
}

func TestMat4_Translate(t *testing.T) {
	m := NewMat4()
	m.Translate(NewVec3ByValues(1, 2, 3))

	// 原点を変換すると移動量になるはず
	v := NewVec3()
	result := m.MulVec3(v)
	expected := NewVec3ByValues(1, 2, 3)

	if !result.NearEquals(expected, 1e-10) {
		t.Errorf("Translate() * origin = %v, want %v", result, expected)
	}
}

func TestMat4_Scale(t *testing.T) {
	m := NewMat4()
	m.Scale(NewVec3ByValues(2, 3, 4))

	v := NewVec3ByValues(1, 1, 1)
	result := m.MulVec3(v)
	expected := NewVec3ByValues(2, 3, 4)

	if !result.NearEquals(expected, 1e-10) {
		t.Errorf("Scale() * (1,1,1) = %v, want %v", result, expected)
	}
}

func TestMat4_RotateX(t *testing.T) {
	m := NewMat4()
	m.RotateX(math.Pi / 2) // 90度

	v := NewVec3ByValues(0, 1, 0)
	result := m.MulVec3(v)
	expected := NewVec3ByValues(0, 0, 1)

	if !result.NearEquals(expected, 1e-6) {
		t.Errorf("RotateX(90°) * (0,1,0) = %v, want %v", result, expected)
	}
}

func TestMat4_RotateY(t *testing.T) {
	m := NewMat4()
	m.RotateY(math.Pi / 2) // 90度

	v := NewVec3ByValues(1, 0, 0)
	result := m.MulVec3(v)
	expected := NewVec3ByValues(0, 0, -1)

	if !result.NearEquals(expected, 1e-6) {
		t.Errorf("RotateY(90°) * (1,0,0) = %v, want %v", result, expected)
	}
}

func TestMat4_RotateZ(t *testing.T) {
	m := NewMat4()
	m.RotateZ(math.Pi / 2) // 90度

	v := NewVec3ByValues(1, 0, 0)
	result := m.MulVec3(v)
	expected := NewVec3ByValues(0, 1, 0)

	if !result.NearEquals(expected, 1e-6) {
		t.Errorf("RotateZ(90°) * (1,0,0) = %v, want %v", result, expected)
	}
}

func TestMat4_MulVec3(t *testing.T) {
	tests := []struct {
		name     string
		m        *Mat4
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "単位行列",
			m:        NewMat4(),
			v:        NewVec3ByValues(1, 2, 3),
			expected: NewVec3ByValues(1, 2, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.m.MulVec3(tt.v)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("MulVec3() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMat4_MulVec4(t *testing.T) {
	tests := []struct {
		name     string
		m        *Mat4
		v        *Vec4
		expected *Vec4
	}{
		{
			name:     "単位行列",
			m:        NewMat4(),
			v:        NewVec4ByValues(1, 2, 3, 1),
			expected: NewVec4ByValues(1, 2, 3, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.m.MulVec4(tt.v)
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("MulVec4() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMat4_Copy(t *testing.T) {
	m := NewMat4()
	m.Translate(NewVec3ByValues(1, 2, 3))

	copied := m.Copy()

	if !m.NearEquals(copied, 1e-10) {
		t.Errorf("Copy() returned different values")
	}

	// コピーを変更しても元の行列に影響しないことを確認
	copied.Set(0, 0, 100)
	if m.At(0, 0) == 100 {
		t.Errorf("Copy() did not create a deep copy")
	}
}

func TestMat4_NearEquals(t *testing.T) {
	m1 := NewMat4()
	m2 := NewMat4()

	if !m1.NearEquals(m2, 1e-10) {
		t.Errorf("Two identity matrices should be equal")
	}

	m2.Set(0, 0, 1.0000001)
	if !m1.NearEquals(m2, 1e-6) {
		t.Errorf("Matrices should be near equal with epsilon 1e-6")
	}

	if m1.NearEquals(m2, 1e-10) {
		t.Errorf("Matrices should not be equal with epsilon 1e-10")
	}
}

func TestMat4_GL(t *testing.T) {
	m := NewMat4()
	gl := m.GL()

	if len(gl) != 16 {
		t.Errorf("GL() length = %v, want 16", len(gl))
	}

	// 対角要素が1であることを確認
	if gl[0] != 1 || gl[5] != 1 || gl[10] != 1 || gl[15] != 1 {
		t.Errorf("GL() diagonal should be 1")
	}
}
