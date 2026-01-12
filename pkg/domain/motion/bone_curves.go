// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// INITIAL_BONE_CURVES はボーン曲線の既定値を表す。
var INITIAL_BONE_CURVES = []byte{
	20,
	20,
	0,
	0,
	20,
	20,
	20,
	20,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	20,
	20,
	20,
	20,
	20,
	20,
	20,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	0,
	20,
	20,
	20,
	20,
	20,
	20,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	0,
	0,
	20,
	20,
	20,
	20,
	20,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	107,
	0,
	0,
	0,
}

// BoneCurves はボーン補間曲線を表す。
type BoneCurves struct {
	TranslateX *mmath.Curve
	TranslateY *mmath.Curve
	TranslateZ *mmath.Curve
	Rotate     *mmath.Curve
	Values     []byte
}

// NewBoneCurves はBoneCurvesを生成する。
func NewBoneCurves() *BoneCurves {
	values := append([]byte(nil), INITIAL_BONE_CURVES...)
	return NewBoneCurvesByValues(values)
}

// NewBoneCurvesByValues はVMD値からBoneCurvesを生成する。
func NewBoneCurvesByValues(values []byte) *BoneCurves {
	if len(values) < len(INITIAL_BONE_CURVES) {
		values = append([]byte(nil), INITIAL_BONE_CURVES...)
	}
	copied := append([]byte(nil), values...)
	return &BoneCurves{
		TranslateX: newCurveByValues(copied[0], copied[4], copied[8], copied[12]),
		TranslateY: newCurveByValues(copied[16], copied[20], copied[24], copied[28]),
		TranslateZ: newCurveByValues(copied[32], copied[36], copied[40], copied[44]),
		Rotate:     newCurveByValues(copied[48], copied[52], copied[56], copied[60]),
		Values:     copied,
	}
}

// Evaluate は補間係数を返す。
func (c *BoneCurves) Evaluate(prevIndex, nowIndex, nextIndex Frame) (float64, float64, float64, float64) {
	if c == nil {
		return 0, 0, 0, 0
	}
	_, xy, _ := mmath.Evaluate(c.TranslateX, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, yy, _ := mmath.Evaluate(c.TranslateY, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, zy, _ := mmath.Evaluate(c.TranslateZ, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, ry, _ := mmath.Evaluate(c.Rotate, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	return xy, yy, zy, ry
}

// Merge はVMD形式の配列へ変換する。
func (c *BoneCurves) Merge(enablePhysics bool) []byte {
	if c == nil {
		return append([]byte(nil), INITIAL_BONE_CURVES...)
	}
	var c02 byte
	var c03 byte
	c31 := byte(1)
	c46 := byte(1)
	c47 := byte(0)
	c61 := byte(1)
	c62 := byte(0)
	c63 := byte(0)
	if c.Values != nil {
		if len(c.Values) > 31 {
			c31 = c.Values[31]
		}
		if len(c.Values) > 46 {
			c46 = c.Values[46]
		}
		if len(c.Values) > 47 {
			c47 = c.Values[47]
		}
		if len(c.Values) > 61 {
			c61 = c.Values[61]
		}
		if len(c.Values) > 62 {
			c62 = c.Values[62]
		}
		if len(c.Values) > 63 {
			c63 = c.Values[63]
		}
	}
	if enablePhysics {
		c02 = byte(curveStartX(c.TranslateZ))
		c03 = byte(curveStartX(c.Rotate))
	} else {
		c02 = byte(99)
		c03 = byte(15)
	}

	return []byte{
		byte(curveStartX(c.TranslateX)),
		byte(curveStartX(c.TranslateY)),
		c02,
		c03,
		byte(curveStartY(c.TranslateX)),
		byte(curveStartY(c.TranslateY)),
		byte(curveStartY(c.TranslateZ)),
		byte(curveStartY(c.Rotate)),
		byte(curveEndX(c.TranslateX)),
		byte(curveEndX(c.TranslateY)),
		byte(curveEndX(c.TranslateZ)),
		byte(curveEndX(c.Rotate)),
		byte(curveEndY(c.TranslateX)),
		byte(curveEndY(c.TranslateY)),
		byte(curveEndY(c.TranslateZ)),
		byte(curveEndY(c.Rotate)),
		byte(curveStartX(c.TranslateY)),
		byte(curveStartX(c.TranslateZ)),
		byte(curveStartX(c.Rotate)),
		byte(curveStartY(c.TranslateX)),
		byte(curveStartY(c.TranslateY)),
		byte(curveStartY(c.TranslateZ)),
		byte(curveStartY(c.Rotate)),
		byte(curveEndX(c.TranslateX)),
		byte(curveEndX(c.TranslateY)),
		byte(curveEndX(c.TranslateZ)),
		byte(curveEndX(c.Rotate)),
		byte(curveEndY(c.TranslateX)),
		byte(curveEndY(c.TranslateY)),
		byte(curveEndY(c.TranslateZ)),
		byte(curveEndY(c.Rotate)),
		c31,
		byte(curveStartX(c.TranslateZ)),
		byte(curveStartX(c.Rotate)),
		byte(curveStartY(c.TranslateX)),
		byte(curveStartY(c.TranslateY)),
		byte(curveStartY(c.TranslateZ)),
		byte(curveStartY(c.Rotate)),
		byte(curveEndX(c.TranslateX)),
		byte(curveEndX(c.TranslateY)),
		byte(curveEndX(c.TranslateZ)),
		byte(curveEndX(c.Rotate)),
		byte(curveEndY(c.TranslateX)),
		byte(curveEndY(c.TranslateY)),
		byte(curveEndY(c.TranslateZ)),
		byte(curveEndY(c.Rotate)),
		c46,
		c47,
		byte(curveStartX(c.Rotate)),
		byte(curveStartY(c.TranslateX)),
		byte(curveStartY(c.TranslateY)),
		byte(curveStartY(c.TranslateZ)),
		byte(curveStartY(c.Rotate)),
		byte(curveEndX(c.TranslateX)),
		byte(curveEndX(c.TranslateY)),
		byte(curveEndX(c.TranslateZ)),
		byte(curveEndX(c.Rotate)),
		byte(curveEndY(c.TranslateX)),
		byte(curveEndY(c.TranslateY)),
		byte(curveEndY(c.TranslateZ)),
		byte(curveEndY(c.Rotate)),
		c61,
		c62,
		c63,
	}
}

// Copy は曲線を複製する。
func (c *BoneCurves) Copy() *BoneCurves {
	if c == nil {
		return nil
	}
	values := append([]byte(nil), c.Values...)
	return &BoneCurves{
		TranslateX: copyCurve(c.TranslateX),
		TranslateY: copyCurve(c.TranslateY),
		TranslateZ: copyCurve(c.TranslateZ),
		Rotate:     copyCurve(c.Rotate),
		Values:     values,
	}
}

func newCurveByValues(x1, y1, x2, y2 byte) *mmath.Curve {
	return &mmath.Curve{
		Start: mmath.Vec2{X: float64(x1), Y: float64(y1)},
		End:   mmath.Vec2{X: float64(x2), Y: float64(y2)},
	}
}

func copyCurve(src *mmath.Curve) *mmath.Curve {
	if src == nil {
		return nil
	}
	return &mmath.Curve{Start: src.Start, End: src.End}
}

func curveStartX(curve *mmath.Curve) float64 {
	if curve == nil {
		return mmath.NewCurve().Start.X
	}
	return curve.Start.X
}

func curveStartY(curve *mmath.Curve) float64 {
	if curve == nil {
		return mmath.NewCurve().Start.Y
	}
	return curve.Start.Y
}

func curveEndX(curve *mmath.Curve) float64 {
	if curve == nil {
		return mmath.NewCurve().End.X
	}
	return curve.End.X
}

func curveEndY(curve *mmath.Curve) float64 {
	if curve == nil {
		return mmath.NewCurve().End.Y
	}
	return curve.End.Y
}
