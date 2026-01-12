// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// INITIAL_CAMERA_CURVES はカメラ曲線の既定値を表す。
var INITIAL_CAMERA_CURVES = []byte{
	20,
	20,
	20,
	20,
	20,
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
	107,
	107,
	107,
	107,
}

// CameraCurves はカメラ補間曲線を表す。
type CameraCurves struct {
	TranslateX  *mmath.Curve
	TranslateY  *mmath.Curve
	TranslateZ  *mmath.Curve
	Rotate      *mmath.Curve
	Distance    *mmath.Curve
	ViewOfAngle *mmath.Curve
	Values      []byte
}

// NewCameraCurves はCameraCurvesを生成する。
func NewCameraCurves() *CameraCurves {
	values := append([]byte(nil), INITIAL_CAMERA_CURVES...)
	return NewCameraCurvesByValues(values)
}

// NewCameraCurvesByValues はVMD値からCameraCurvesを生成する。
func NewCameraCurvesByValues(values []byte) *CameraCurves {
	if len(values) < len(INITIAL_CAMERA_CURVES) {
		values = append([]byte(nil), INITIAL_CAMERA_CURVES...)
	}
	copied := append([]byte(nil), values...)
	return &CameraCurves{
		TranslateX:  newCurveByValues(copied[0], copied[6], copied[12], copied[18]),
		TranslateY:  newCurveByValues(copied[1], copied[7], copied[13], copied[19]),
		TranslateZ:  newCurveByValues(copied[2], copied[8], copied[14], copied[20]),
		Rotate:      newCurveByValues(copied[3], copied[9], copied[15], copied[21]),
		Distance:    newCurveByValues(copied[4], copied[10], copied[16], copied[22]),
		ViewOfAngle: newCurveByValues(copied[5], copied[11], copied[17], copied[23]),
		Values:      copied,
	}
}

// Evaluate は補間係数を返す。
func (c *CameraCurves) Evaluate(prevIndex, nowIndex, nextIndex Frame) (float64, float64, float64, float64, float64, float64) {
	if c == nil {
		return 0, 0, 0, 0, 0, 0
	}
	_, xy, _ := mmath.Evaluate(c.TranslateX, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, yy, _ := mmath.Evaluate(c.TranslateY, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, zy, _ := mmath.Evaluate(c.TranslateZ, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, ry, _ := mmath.Evaluate(c.Rotate, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, dy, _ := mmath.Evaluate(c.Distance, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	_, vy, _ := mmath.Evaluate(c.ViewOfAngle, float32(prevIndex), float32(nowIndex), float32(nextIndex))
	return xy, yy, zy, ry, dy, vy
}

// Merge はVMD形式の配列へ変換する。
func (c *CameraCurves) Merge() []byte {
	if c == nil {
		return append([]byte(nil), INITIAL_CAMERA_CURVES...)
	}
	return []byte{
		byte(curveStartX(c.TranslateX)),
		byte(curveStartX(c.TranslateY)),
		byte(curveStartX(c.TranslateZ)),
		byte(curveStartX(c.Rotate)),
		byte(curveStartX(c.Distance)),
		byte(curveStartX(c.ViewOfAngle)),
		byte(curveStartY(c.TranslateX)),
		byte(curveStartY(c.TranslateY)),
		byte(curveStartY(c.TranslateZ)),
		byte(curveStartY(c.Rotate)),
		byte(curveStartY(c.Distance)),
		byte(curveStartY(c.ViewOfAngle)),
		byte(curveEndX(c.TranslateX)),
		byte(curveEndX(c.TranslateY)),
		byte(curveEndX(c.TranslateZ)),
		byte(curveEndX(c.Rotate)),
		byte(curveEndX(c.Distance)),
		byte(curveEndX(c.ViewOfAngle)),
		byte(curveEndY(c.TranslateX)),
		byte(curveEndY(c.TranslateY)),
		byte(curveEndY(c.TranslateZ)),
		byte(curveEndY(c.Rotate)),
		byte(curveEndY(c.Distance)),
		byte(curveEndY(c.ViewOfAngle)),
	}
}

// Copy は曲線を複製する。
func (c *CameraCurves) Copy() *CameraCurves {
	if c == nil {
		return nil
	}
	values := append([]byte(nil), c.Values...)
	return &CameraCurves{
		TranslateX:  copyCurve(c.TranslateX),
		TranslateY:  copyCurve(c.TranslateY),
		TranslateZ:  copyCurve(c.TranslateZ),
		Rotate:      copyCurve(c.Rotate),
		Distance:    copyCurve(c.Distance),
		ViewOfAngle: copyCurve(c.ViewOfAngle),
		Values:      values,
	}
}
