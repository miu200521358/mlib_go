package mmath

import (
	"testing"
)

func TestNewCurve(t *testing.T) {
	c := NewCurve()

	// 線形補間曲線であることを確認
	if !c.IsLinear() {
		t.Errorf("NewCurve() should return linear curve")
	}

	if c.Start.X != 20 || c.Start.Y != 20 {
		t.Errorf("NewCurve() Start = %v, want (20, 20)", c.Start)
	}

	if c.End.X != 107 || c.End.Y != 107 {
		t.Errorf("NewCurve() End = %v, want (107, 107)", c.End)
	}
}

func TestNewCurveByValues(t *testing.T) {
	tests := []struct {
		name           string
		startX, startY byte
		endX, endY     byte
		expectedLinear bool
	}{
		{
			name:           "線形補間",
			startX:         20,
			startY:         20,
			endX:           107,
			endY:           107,
			expectedLinear: true,
		},
		{
			name:           "非線形補間",
			startX:         64,
			startY:         0,
			endX:           64,
			endY:           127,
			expectedLinear: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCurveByValues(tt.startX, tt.startY, tt.endX, tt.endY)
			if c.IsLinear() != tt.expectedLinear {
				t.Errorf("NewCurveByValues() IsLinear() = %v, want %v", c.IsLinear(), tt.expectedLinear)
			}
		})
	}
}

func TestCurve_Copy(t *testing.T) {
	c := NewCurveByValues(30, 40, 100, 110)
	copied := c.Copy()

	if c.Start.X != copied.Start.X || c.Start.Y != copied.Start.Y {
		t.Errorf("Copy() Start mismatch")
	}
	if c.End.X != copied.End.X || c.End.Y != copied.End.Y {
		t.Errorf("Copy() End mismatch")
	}

	// コピーを変更しても元に影響しないことを確認
	copied.Start.X = 999
	if c.Start.X == 999 {
		t.Errorf("Copy() did not create a deep copy")
	}
}

func TestCurve_IsLinear(t *testing.T) {
	tests := []struct {
		name     string
		curve    *Curve
		expected bool
	}{
		{
			name:     "線形補間",
			curve:    NewCurve(),
			expected: true,
		},
		{
			name:     "Start.X != Start.Y",
			curve:    &Curve{Start: Vec2{X: 20, Y: 30}, End: Vec2{X: 107, Y: 107}},
			expected: false,
		},
		{
			name:     "End.X != End.Y",
			curve:    &Curve{Start: Vec2{X: 20, Y: 20}, End: Vec2{X: 100, Y: 107}},
			expected: false,
		},
		{
			name:     "両方異なる",
			curve:    &Curve{Start: Vec2{X: 64, Y: 0}, End: Vec2{X: 64, Y: 127}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.curve.IsLinear()
			if result != tt.expected {
				t.Errorf("IsLinear() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name      string
		curve     *Curve
		start     float32
		now       float32
		end       float32
		expectedX float64
		expectedY float64
	}{
		{
			name:      "線形補間 t=0",
			curve:     NewCurve(),
			start:     0,
			now:       0,
			end:       100,
			expectedX: 0,
			expectedY: 0,
		},
		{
			name:      "線形補間 t=0.5",
			curve:     NewCurve(),
			start:     0,
			now:       50,
			end:       100,
			expectedX: 0.5,
			expectedY: 0.5,
		},
		{
			name:      "線形補間 t=1",
			curve:     NewCurve(),
			start:     0,
			now:       100,
			end:       100,
			expectedX: 1.0,
			expectedY: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y, _ := Evaluate(tt.curve, tt.start, tt.now, tt.end)
			if !NearEquals(x, tt.expectedX, 1e-6) {
				t.Errorf("Evaluate() x = %v, want %v", x, tt.expectedX)
			}
			if !NearEquals(y, tt.expectedY, 1e-6) {
				t.Errorf("Evaluate() y = %v, want %v", y, tt.expectedY)
			}
		})
	}
}

func TestEvaluate_NonLinear(t *testing.T) {
	// 非線形補間曲線のテスト
	// イージングアウト（最初が速く後が遅い）
	curve := &Curve{
		Start: Vec2{X: 0, Y: 64},
		End:   Vec2{X: 64, Y: 127},
	}

	// 中間点での評価
	_, y, _ := Evaluate(curve, 0, 50, 100)

	// イージングアウトなので、中間点でのY値は0.5より大きいはず
	if y <= 0.5 {
		t.Errorf("Ease-out curve at t=0.5 should have y > 0.5, got %v", y)
	}
}

func TestSplitCurve(t *testing.T) {
	// 線形補間曲線の分割
	curve := NewCurve()
	startCurve, endCurve := SplitCurve(curve, 0, 50, 100)

	// 分割後も線形補間であるべき
	if !startCurve.IsLinear() {
		t.Errorf("SplitCurve() startCurve should be linear")
	}
	if !endCurve.IsLinear() {
		t.Errorf("SplitCurve() endCurve should be linear")
	}
}

func TestSplitCurve_EdgeCases(t *testing.T) {
	curve := NewCurve()

	// start == now の場合
	startCurve, endCurve := SplitCurve(curve, 0, 0, 100)
	if startCurve == nil || endCurve == nil {
		t.Errorf("SplitCurve() should not return nil for edge cases")
	}
}

func TestCurve_Normalize(t *testing.T) {
	c := &Curve{
		Start: Vec2{X: 0.2, Y: 0.3},
		End:   Vec2{X: 0.8, Y: 0.9},
	}

	begin := &Vec2{X: 0, Y: 0}
	finish := &Vec2{X: 1, Y: 1}

	c.Normalize(begin, finish)

	// 正規化後の値が0-127の範囲内であることを確認
	if c.Start.X < 0 || c.Start.X > CURVE_MAX {
		t.Errorf("Normalize() Start.X = %v, should be in [0, 127]", c.Start.X)
	}
	if c.End.X < 0 || c.End.X > CURVE_MAX {
		t.Errorf("Normalize() End.X = %v, should be in [0, 127]", c.End.X)
	}
}
