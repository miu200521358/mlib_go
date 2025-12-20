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
		{
			name:     "既存テスト1",
			axis:     NewVec3ByValues(1, 2, 3),
			angle:    DegToRad(30),
			expected: NewQuaternionByValues(0.0691722994246875, 0.138344598849375, 0.207516898274062, 0.965925826289068),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewQuaternionFromAxisAngle(tt.axis, tt.angle)
			if !result.NearEquals(tt.expected, 1e-5) {
				t.Errorf("NewQuaternionFromAxisAngle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_ToDegrees(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected *Vec3
	}{
		{
			name:     "ゼロ回転",
			q:        NewQuaternionByValues(0, 0, 0, 1),
			expected: NewVec3ByValues(0, 0, 0),
		},
		{
			name:     "X軸10度",
			q:        NewQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455),
			expected: NewVec3ByValues(10, 0, 0),
		},
		{
			name:     "10,20,30度",
			q:        NewQuaternionByValues(0.12767944, 0.14487813, 0.23929834, 0.95154852),
			expected: NewVec3ByValues(10, 20, 30),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.ToDegrees()
			if !result.NearEquals(tt.expected, 1e-5) {
				t.Errorf("ToDegrees() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_FromDegrees(t *testing.T) {
	tests := []struct {
		name                 string
		xPitch, yHead, zRoll float64
		expected             *Quaternion
	}{
		{
			name:   "ゼロ",
			xPitch: 0, yHead: 0, zRoll: 0,
			expected: NewQuaternionByValues(0, 0, 0, 1),
		},
		{
			name:   "X軸10度",
			xPitch: 10, yHead: 0, zRoll: 0,
			expected: NewQuaternionByValues(0.08715574, 0.0, 0.0, 0.9961947),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewQuaternionFromDegrees(tt.xPitch, tt.yHead, tt.zRoll)
			if !result.NearEquals(tt.expected, 1e-6) {
				t.Errorf("NewQuaternionFromDegrees() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_Dot(t *testing.T) {
	tests := []struct {
		name     string
		q1       *Quaternion
		q2       *Quaternion
		expected float64
	}{
		{
			name:     "既存テスト1",
			q1:       NewQuaternionByValues(0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844),
			q2:       NewQuaternionByValues(0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885),
			expected: 0.6491836986795888,
		},
		{
			name:     "既存テスト2",
			q1:       NewQuaternionByValues(0.1549093965157679, 0.15080756177478563, 0.3575205710320892, 0.908536845412201),
			q2:       NewQuaternionByValues(0.15799222008931638, 0.1243359045760714, 0.33404459937562386, 0.9208654879256133),
			expected: 0.9992933154462645,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q1.Dot(tt.q2)
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("Dot() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_Normalized(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected *Quaternion
	}{
		{
			name:     "既存テスト1",
			q:        NewQuaternionByValues(2, 3, 4, 1),
			expected: NewQuaternionByValues(0.36514837, 0.54772256, 0.73029674, 0.18257419),
		},
		{
			name:     "単位クォータニオン",
			q:        NewQuaternionByValues(0, 0, 0, 1),
			expected: NewQuaternionByValues(0, 0, 0, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.Normalized()
			if !result.NearEquals(tt.expected, 1e-7) {
				t.Errorf("Normalized() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_Muled(t *testing.T) {
	tests := []struct {
		name     string
		q1       *Quaternion
		q2       *Quaternion
		expected *Quaternion
	}{
		{
			name:     "既存テスト1",
			q1:       NewQuaternionByValues(0.4738680537545347, 0.20131048764138487, -0.48170221425083437, 0.7091446481376844),
			q2:       NewQuaternionByValues(0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885),
			expected: NewQuaternionByValues(0.6594130183457979, 0.11939693791117263, -0.24571599091322077, 0.7003873887093154),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q1.Muled(tt.q2)
			if !result.NearEquals(tt.expected, 1e-8) {
				t.Errorf("Muled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_ToDegree(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected float64
	}{
		{
			name:     "10度",
			q:        NewQuaternionByValues(0.08715574274765817, 0.0, 0.0, 0.9961946980917455),
			expected: 10.0,
		},
		{
			name:     "既存テスト",
			q:        NewQuaternionByValues(0.12767944069578063, 0.14487812541736916, 0.2392983377447303, 0.9515485246437885),
			expected: 35.81710117358426,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.ToDegree()
			if !NearEquals(result, tt.expected, 1e-10) {
				t.Errorf("ToDegree() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_Rotate(t *testing.T) {
	tests := []struct {
		name     string
		from     *Vec3
		to       *Vec3
		expected *Quaternion
	}{
		{
			name:     "既存テスト1",
			from:     NewVec3ByValues(1, 2, 3),
			to:       NewVec3ByValues(4, 5, 6),
			expected: NewQuaternionByValues(-0.04597839511020707, 0.0919567902204141, -0.04597839511020706, 0.9936377222602503),
		},
		{
			name:     "既存テスト2",
			from:     NewVec3ByValues(-10, 20, -15),
			to:       NewVec3ByValues(40, -5, 6),
			expected: NewQuaternionByValues(0.042643949239185255, -0.511727390870223, -0.7107324873197542, 0.48080755245182594),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewQuaternionRotate(tt.from, tt.to)
			if !result.NearEquals(tt.expected, 1e-5) {
				t.Errorf("NewQuaternionRotate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// 既存mlib_goからの移行テスト

func TestQuaternion_ToMat4(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected *Mat4
	}{
		{
			name: "単位クォータニオン",
			q:    NewQuaternion(),
			expected: &Mat4{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		{
			name: "90度回転",
			q:    NewQuaternionByValues(0.5, 0.5, 0.5, 0.5),
			expected: &Mat4{
				0.0, 1.0, 0.0, 0.0,
				0.0, 0.0, 1.0, 0.0,
				1.0, 0.0, 0.0, 0.0,
				0.0, 0.0, 0.0, 1.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.ToMat4()
			if !result.NearEquals(tt.expected, 1e-10) {
				t.Errorf("ToMat4() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_Length(t *testing.T) {
	q := NewQuaternionByValues(1, 2, 3, 4)
	// 正規化されていないクォータニオンの長さ
	expected := 5.477225575051661 // sqrt(1+4+9+16) = sqrt(30)

	result := q.Length()
	if !NearEquals(result, expected, 1e-10) {
		t.Errorf("Length() = %v, want %v", result, expected)
	}

	// 正規化後のクォータニオンの長さは1
	qNorm := q.Normalized()
	if !NearEquals(qNorm.Length(), 1.0, 1e-10) {
		t.Errorf("Normalized().Length() = %v, want 1.0", qNorm.Length())
	}
}

func TestQuaternion_Lerp(t *testing.T) {
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
			q2:       NewQuaternionFromDegrees(90, 0, 0),
			t:        0,
			expected: NewQuaternion(),
		},
		{
			name:     "t=1",
			q1:       NewQuaternion(),
			q2:       NewQuaternionFromDegrees(90, 0, 0),
			t:        1,
			expected: NewQuaternionFromDegrees(90, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q1.Lerp(tt.q2, tt.t)
			if !result.NearEquals(tt.expected, 1e-5) {
				t.Errorf("Lerp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestQuaternion_Inverted(t *testing.T) {
	q := NewQuaternionFromDegrees(30, 45, 60)
	inv := q.Inverted()

	// q * inv = 単位クォータニオン
	result := q.Muled(inv)
	expected := NewQuaternion()

	if !result.NearEquals(expected, 1e-6) {
		t.Errorf("q * q.Inverted() = %v, want identity", result)
	}
}

func TestQuaternion_Copy(t *testing.T) {
	q := NewQuaternionByValues(1, 2, 3, 4)
	copied := q.Copy()

	if !q.NearEquals(copied, 1e-10) {
		t.Errorf("Copy() returned different values")
	}

	// コピーを変更しても元のクォータニオンに影響しないことを確認
	copied.SetX(999)
	if q.X() == 999 {
		t.Errorf("Copy() did not create a deep copy")
	}
}

func TestQuaternion_IsIdent(t *testing.T) {
	tests := []struct {
		name     string
		q        *Quaternion
		expected bool
	}{
		{
			name:     "単位クォータニオン",
			q:        NewQuaternion(),
			expected: true,
		},
		{
			name:     "非単位",
			q:        NewQuaternionFromDegrees(10, 0, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.q.IsIdent()
			if result != tt.expected {
				t.Errorf("IsIdent() = %v, want %v", result, tt.expected)
			}
		})
	}
}
