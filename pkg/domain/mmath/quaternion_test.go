package mmath

import (
	"math"
	"testing"
)

func TestQuaternion_Conjugate(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected *Quaternion
	}{
		{
			name:     "単位クォータニオン",
			q:        NewQuaternion(),
			expected: NewQuaternion(),
		},
		{
			name:     "一般的なクォータニオン",
			q:        NewQuaternionByValues(1, 2, 3, 4),
			expected: NewQuaternionByValues(-1, -2, -3, 4), // 虚部の符号が反転
		},
		{
			name:     "90度X軸回転",
			q:        NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/2),
			expected: NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), -math.Pi/2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.Conjugated()
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("Conjugated() = %v, want %v", result, tt.expected)
			}
			// 元のクォータニオンが変更されていないことを確認
			if tt.q.X() == result.X() && tt.q.Y() == result.Y() && tt.q.Z() == result.Z() && tt.q.X() != 0 {
				t.Errorf("Conjugated() should not modify original quaternion")
			}
		})
	}
}

func TestQuaternion_Conjugate_Destructive(t *testing.T) {
	q := NewQuaternionByValues(1, 2, 3, 4)
	q.Conjugate()

	expected := NewQuaternionByValues(-1, -2, -3, 4)
	if !q.NearEquals(expected, 1e-6) {
		t.Errorf("Conjugate() = %v, want %v", q, expected)
	}
}

func TestQuaternion_Log(t *testing.T) {
	tests := []struct {
		name string
		q    *Quaternion
	}{
		{
			name: "単位クォータニオン",
			q:    NewQuaternion(),
		},
		{
			name: "90度X軸回転",
			q:    NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/2),
		},
		{
			name: "45度Y軸回転",
			q:    NewQuaternionFromAxisAngle(VEC3_UNIT_Y.Copy(), math.Pi/4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logQ := tt.q.Log()
			// Log の結果は実数部が0（または非常に小さい）になるはず
			// Log(Exp(q)) = q を確認
			expLogQ := logQ.Exp()
			if !expLogQ.NearEquals(tt.q, 1e-6) {
				t.Errorf("Exp(Log(q)) = %v, want %v", expLogQ, tt.q)
			}
		})
	}
}

func TestQuaternion_Exp(t *testing.T) {
	tests := []struct {
		name string
		q    *Quaternion
	}{
		{
			name: "単位クォータニオン",
			q:    NewQuaternion(),
		},
		{
			name: "45度X軸回転",
			q:    NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/4),
		},
		{
			name: "180度Z軸回転",
			q:    NewQuaternionFromAxisAngle(VEC3_UNIT_Z.Copy(), math.Pi),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logQ := tt.q.Log()
			expLogQ := logQ.Exp()
			// Exp(Log(q)) = q を確認
			if !expLogQ.NearEquals(tt.q, 1e-6) {
				t.Errorf("Exp(Log(q)) = %v, want %v", expLogQ, tt.q)
			}
		})
	}
}

func TestQuaternion_Slerp(t *testing.T) {
	tests := []struct {
		name     string
		q1       *Quaternion
		q2       *Quaternion
		t        float64
		expected *Quaternion
	}{
		{
			name:     "t=0",
			q1:       NewQuaternion(),
			q2:       NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/2),
			t:        0,
			expected: NewQuaternion(),
		},
		{
			name:     "t=1",
			q1:       NewQuaternion(),
			q2:       NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/2),
			t:        1,
			expected: NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/2),
		},
		{
			name:     "t=0.5",
			q1:       NewQuaternion(),
			q2:       NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/2),
			t:        0.5,
			expected: NewQuaternionFromAxisAngle(VEC3_UNIT_X.Copy(), math.Pi/4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q1.Slerp(tt.q2, tt.t)
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("Slerp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_MulVec3(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		v        *Vec3
		expected *Vec3
	}{
		{
			name:     "単位クォータニオンで回転",
			q:        NewQuaternion(),
			v:        NewVec3ByValues(1, 0, 0),
			expected: NewVec3ByValues(1, 0, 0),
		},
		{
			name:     "90度Z軸回転",
			q:        NewQuaternionFromAxisAngle(VEC3_UNIT_Z.Copy(), math.Pi/2),
			v:        NewVec3ByValues(1, 0, 0),
			expected: NewVec3ByValues(0, 1, 0),
		},
		{
			name:     "180度Y軸回転",
			q:        NewQuaternionFromAxisAngle(VEC3_UNIT_Y.Copy(), math.Pi),
			v:        NewVec3ByValues(1, 0, 0),
			expected: NewVec3ByValues(-1, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.MulVec3(tt.v)
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("MulVec3() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_ToRadians(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected *Vec3
	}{
		{
			name:     "単位クォータニオン",
			q:        NewQuaternion(),
			expected: NewVec3(),
		},
		{
			name:     "90度X軸回転",
			q:        NewQuaternionFromRadians(math.Pi/2, 0, 0),
			expected: NewVec3ByValues(math.Pi/2, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.ToRadians()
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("ToRadians() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_FromAxisAngle(t *testing.T) {
	tests := []struct {
		name     string
		axis     *Vec3
		angle    float64
		expected *Quaternion
	}{
		{
			name:     "X軸90度",
			axis:     VEC3_UNIT_X.Copy(),
			angle:    math.Pi / 2,
			expected: NewQuaternionByValues(math.Sin(math.Pi/4), 0, 0, math.Cos(math.Pi/4)),
		},
		{
			name:     "Y軸180度",
			axis:     VEC3_UNIT_Y.Copy(),
			angle:    math.Pi,
			expected: NewQuaternionByValues(0, 1, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewQuaternionFromAxisAngle(tt.axis, tt.angle)
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("NewQuaternionFromAxisAngle() = %v, want %v", result, tt.expected)
			}
		})
	}
}
