package mmath

import (
	"testing"
)

func TestNewCurve(t *testing.T) {
	tests := []struct {
		name           string
		expectedLinear bool
		expectedStartX float64
		expectedStartY float64
		expectedEndX   float64
		expectedEndY   float64
	}{
		{
			name:           "デフォルト曲線",
			expectedLinear: true,
			expectedStartX: 20,
			expectedStartY: 20,
			expectedEndX:   107,
			expectedEndY:   107,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCurve()
			if !c.IsLinear() {
				t.Errorf("NewCurve() should return linear curve")
			}
			if c.Start.X != tt.expectedStartX || c.Start.Y != tt.expectedStartY {
				t.Errorf("NewCurve() Start = %v, want (%v, %v)", c.Start, tt.expectedStartX, tt.expectedStartY)
			}
			if c.End.X != tt.expectedEndX || c.End.Y != tt.expectedEndY {
				t.Errorf("NewCurve() End = %v, want (%v, %v)", c.End, tt.expectedEndX, tt.expectedEndY)
			}
		})
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
	tests := []struct {
		name   string
		startX byte
		startY byte
		endX   byte
		endY   byte
	}{
		{"基本", 30, 40, 100, 110},
		{"線形", 20, 20, 107, 107},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCurveByValues(tt.startX, tt.startY, tt.endX, tt.endY)
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
		})
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
	tests := []struct {
		name         string
		curve        *Curve
		start        float32
		now          float32
		end          float32
		yGreaterThan float64
	}{
		{
			name: "イージングアウト",
			curve: &Curve{
				Start: Vec2{X: 0, Y: 64},
				End:   Vec2{X: 64, Y: 127},
			},
			start:        0,
			now:          50,
			end:          100,
			yGreaterThan: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, y, _ := Evaluate(tt.curve, tt.start, tt.now, tt.end)
			if y <= tt.yGreaterThan {
				t.Errorf("Evaluate() y = %v, should be > %v", y, tt.yGreaterThan)
			}
		})
	}
}

func TestSplitCurve(t *testing.T) {
	tests := []struct {
		name               string
		curve              *Curve
		start              float32
		now                float32
		end                float32
		startCurveIsLinear bool
		endCurveIsLinear   bool
	}{
		{
			name:               "線形補間曲線の分割",
			curve:              NewCurve(),
			start:              0,
			now:                50,
			end:                100,
			startCurveIsLinear: true,
			endCurveIsLinear:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startCurve, endCurve := SplitCurve(tt.curve, tt.start, tt.now, tt.end)
			if startCurve.IsLinear() != tt.startCurveIsLinear {
				t.Errorf("SplitCurve() startCurve.IsLinear() = %v, want %v", startCurve.IsLinear(), tt.startCurveIsLinear)
			}
			if endCurve.IsLinear() != tt.endCurveIsLinear {
				t.Errorf("SplitCurve() endCurve.IsLinear() = %v, want %v", endCurve.IsLinear(), tt.endCurveIsLinear)
			}
		})
	}
}

func TestSplitCurve_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		curve *Curve
		start float32
		now   float32
		end   float32
	}{
		{
			name:  "start == now",
			curve: NewCurve(),
			start: 0,
			now:   0,
			end:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startCurve, endCurve := SplitCurve(tt.curve, tt.start, tt.now, tt.end)
			if startCurve == nil || endCurve == nil {
				t.Errorf("SplitCurve() should not return nil for edge cases")
			}
		})
	}
}

func TestCurve_Normalize(t *testing.T) {
	tests := []struct {
		name   string
		curve  *Curve
		begin  *Vec2
		finish *Vec2
	}{
		{
			name: "基本正規化",
			curve: &Curve{
				Start: Vec2{X: 0.2, Y: 0.3},
				End:   Vec2{X: 0.8, Y: 0.9},
			},
			begin:  &Vec2{X: 0, Y: 0},
			finish: &Vec2{X: 1, Y: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.curve.Normalize(tt.begin, tt.finish)
			if tt.curve.Start.X < 0 || tt.curve.Start.X > CURVE_MAX {
				t.Errorf("Normalize() Start.X = %v, should be in [0, 127]", tt.curve.Start.X)
			}
			if tt.curve.End.X < 0 || tt.curve.End.X > CURVE_MAX {
				t.Errorf("Normalize() End.X = %v, should be in [0, 127]", tt.curve.End.X)
			}
		})
	}
}

// 既存mlib_goからの移行テスト

func TestEvaluate_NonLinear2(t *testing.T) {
	tests := []struct {
		name      string
		curve     *Curve
		start     float32
		now       float32
		end       float32
		expectedX float64
		expectedY float64
		expectedT float64
	}{
		{
			name: "既存テスト",
			curve: &Curve{
				Start: Vec2{X: 10.0, Y: 30.0},
				End:   Vec2{X: 100.0, Y: 80.0},
			},
			start:     0,
			now:       2,
			end:       10,
			expectedX: 0.2,
			expectedY: 0.24085271757748078,
			expectedT: 0.2900272452240925,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y, tVal := Evaluate(tt.curve, tt.start, tt.now, tt.end)
			if x != tt.expectedX {
				t.Errorf("Evaluate() x = %v, want %v", x, tt.expectedX)
			}
			if !NearEquals(y, tt.expectedY, 1e-10) {
				t.Errorf("Evaluate() y = %v, want %v", y, tt.expectedY)
			}
			if !NearEquals(tVal, tt.expectedT, 1e-10) {
				t.Errorf("Evaluate() t = %v, want %v", tVal, tt.expectedT)
			}
		})
	}
}

func TestSplitCurve_Variants(t *testing.T) {
	tests := []struct {
		name               string
		curve              *Curve
		start              float32
		now                float32
		end                float32
		expectedStartStart Vec2
		expectedStartEnd   Vec2
		expectedEndStart   Vec2
		expectedEndEnd     Vec2
		epsilon            float64
	}{
		{
			name: "非線形",
			curve: &Curve{
				Start: Vec2{X: 89.0, Y: 2.0},
				End:   Vec2{X: 52.0, Y: 106.0},
			},
			start:              0,
			now:                2,
			end:                10,
			expectedStartStart: Vec2{X: 50, Y: 7},
			expectedStartEnd:   Vec2{X: 91, Y: 52},
			expectedEndStart:   Vec2{X: 71, Y: 21},
			expectedEndEnd:     Vec2{X: 44, Y: 108},
			epsilon:            1e-1,
		},
		{
			name: "線形",
			curve: &Curve{
				Start: Vec2{X: 20.0, Y: 20.0},
				End:   Vec2{X: 107.0, Y: 107.0},
			},
			start:              0,
			now:                50,
			end:                100,
			expectedStartStart: Vec2{X: 20, Y: 20},
			expectedStartEnd:   Vec2{X: 107, Y: 107},
			expectedEndStart:   Vec2{X: 20, Y: 20},
			expectedEndEnd:     Vec2{X: 107, Y: 107},
			epsilon:            0,
		},
		{
			name: "同一点",
			curve: &Curve{
				Start: Vec2{X: 10.0, Y: 10.0},
				End:   Vec2{X: 10.0, Y: 10.0},
			},
			start:              0,
			now:                2,
			end:                10,
			expectedStartStart: Vec2{X: 20, Y: 20},
			expectedStartEnd:   Vec2{X: 107, Y: 107},
			expectedEndStart:   Vec2{X: 20, Y: 20},
			expectedEndEnd:     Vec2{X: 107, Y: 107},
			epsilon:            0,
		},
		{
			name: "範囲外",
			curve: &Curve{
				Start: Vec2{X: 25.0, Y: 101.0},
				End:   Vec2{X: 127.0, Y: 12.0},
			},
			start:              0,
			now:                2,
			end:                10,
			expectedStartStart: Vec2{X: 27, Y: 65},
			expectedStartEnd:   Vec2{X: 73, Y: 103},
			expectedEndStart:   Vec2{X: 49, Y: 44},
			expectedEndEnd:     Vec2{X: 127, Y: 0},
			epsilon:            0,
		},
		{
			name: "NaN対策",
			curve: &Curve{
				Start: Vec2{X: 127.0, Y: 0.0},
				End:   Vec2{X: 0.0, Y: 127.0},
			},
			start:              0,
			now:                2,
			end:                10,
			expectedStartStart: Vec2{X: 50, Y: 0},
			expectedStartEnd:   Vec2{X: 92, Y: 45},
			expectedEndStart:   Vec2{X: 104, Y: 17},
			expectedEndEnd:     Vec2{X: 0, Y: 127},
			epsilon:            0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startCurve, endCurve := SplitCurve(tt.curve, tt.start, tt.now, tt.end)
			if tt.epsilon > 0 {
				if !startCurve.Start.NearEquals(&tt.expectedStartStart, tt.epsilon) {
					t.Errorf("startCurve.Start = %v, want %v", startCurve.Start, tt.expectedStartStart)
				}
				if !startCurve.End.NearEquals(&tt.expectedStartEnd, tt.epsilon) {
					t.Errorf("startCurve.End = %v, want %v", startCurve.End, tt.expectedStartEnd)
				}
				if !endCurve.Start.NearEquals(&tt.expectedEndStart, tt.epsilon) {
					t.Errorf("endCurve.Start = %v, want %v", endCurve.Start, tt.expectedEndStart)
				}
				if !endCurve.End.NearEquals(&tt.expectedEndEnd, tt.epsilon) {
					t.Errorf("endCurve.End = %v, want %v", endCurve.End, tt.expectedEndEnd)
				}
			} else {
				if !startCurve.Start.Equals(&tt.expectedStartStart) {
					t.Errorf("startCurve.Start = %v, want %v", startCurve.Start, tt.expectedStartStart)
				}
				if !startCurve.End.Equals(&tt.expectedStartEnd) {
					t.Errorf("startCurve.End = %v, want %v", startCurve.End, tt.expectedStartEnd)
				}
				if !endCurve.Start.Equals(&tt.expectedEndStart) {
					t.Errorf("endCurve.Start = %v, want %v", endCurve.Start, tt.expectedEndStart)
				}
				if !endCurve.End.Equals(&tt.expectedEndEnd) {
					t.Errorf("endCurve.End = %v, want %v", endCurve.End, tt.expectedEndEnd)
				}
			}
		})
	}
}
