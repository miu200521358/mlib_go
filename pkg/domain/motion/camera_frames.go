// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// CameraFrame はカメラフレームを表す。
type CameraFrame struct {
	*BaseFrame
	Position         *mmath.Vec3
	Degrees          *mmath.Vec3
	Quaternion       *mmath.Quaternion
	Distance         float64
	ViewOfAngle      int
	IsPerspectiveOff bool
	Curves           *CameraCurves
}

// NewCameraFrame はCameraFrameを生成する。
func NewCameraFrame(index Frame) *CameraFrame {
	return &CameraFrame{
		BaseFrame:        NewBaseFrame(index),
		Position:         &mmath.Vec3{},
		Degrees:          &mmath.Vec3{},
		Quaternion:       nil,
		Distance:         0,
		ViewOfAngle:      0,
		IsPerspectiveOff: true,
		Curves:           NewCameraCurves(),
	}
}

// Copy はフレームを複製する。
func (f *CameraFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*CameraFrame)(nil), nil
	}
	var curves *CameraCurves
	if f.Curves != nil {
		curves = f.Curves.Copy()
	}
	copied := &CameraFrame{
		BaseFrame:        &BaseFrame{index: f.Index(), Read: f.Read},
		Position:         copyVec3(f.Position),
		Degrees:          copyVec3(f.Degrees),
		Quaternion:       copyQuaternion(f.Quaternion),
		Distance:         f.Distance,
		ViewOfAngle:      f.ViewOfAngle,
		IsPerspectiveOff: f.IsPerspectiveOff,
		Curves:           curves,
	}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *CameraFrame) lerpFrame(prev *CameraFrame, index Frame) *CameraFrame {
	if prev == nil && next == nil {
		return NewCameraFrame(index)
	}
	if prev == nil {
		copied, _ := next.Copy()
		out := copied.(*CameraFrame)
		out.SetIndex(index)
		return out
	}
	if next == nil {
		copied, _ := prev.Copy()
		out := copied.(*CameraFrame)
		out.SetIndex(index)
		return out
	}

	cf := NewCameraFrame(index)
	xy, yy, zy, ry, dy, vy := cameraCurveT(prev.Index(), index, next.Index(), next.Curves)

	prevDeg := vec3OrZero(prev.Degrees)
	nextDeg := vec3OrZero(next.Degrees)
	q1 := mmath.NewQuaternionFromDegrees(prevDeg.X, prevDeg.Y, prevDeg.Z)
	q2 := mmath.NewQuaternionFromDegrees(nextDeg.X, nextDeg.Y, nextDeg.Z)
	q := q1.Slerp(q2, ry)
	cf.Quaternion = &q
	deg := q.ToDegrees()
	cf.Degrees = &deg

	prevPos := vec3OrZero(prev.Position)
	nextPos := vec3OrZero(next.Position)
	pos := vec3(
		mmath.Lerp(prevPos.X, nextPos.X, xy),
		mmath.Lerp(prevPos.Y, nextPos.Y, yy),
		mmath.Lerp(prevPos.Z, nextPos.Z, zy),
	)
	cf.Position = &pos

	cf.Distance = mmath.Lerp(prev.Distance, next.Distance, dy)
	cf.ViewOfAngle = int(mmath.Lerp(float64(prev.ViewOfAngle), float64(next.ViewOfAngle), vy))
	cf.IsPerspectiveOff = next.IsPerspectiveOff

	return cf
}

// splitCurve は曲線を分割する。
func (f *CameraFrame) splitCurve(prev *CameraFrame, next *CameraFrame, index Frame) {
	if f == nil || prev == nil || next == nil || next.Curves == nil {
		return
	}
	if f.Curves == nil {
		f.Curves = NewCameraCurves()
	}
	f.Curves.TranslateX, next.Curves.TranslateX =
		mmath.SplitCurve(next.Curves.TranslateX, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.TranslateY, next.Curves.TranslateY =
		mmath.SplitCurve(next.Curves.TranslateY, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.TranslateZ, next.Curves.TranslateZ =
		mmath.SplitCurve(next.Curves.TranslateZ, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.Rotate, next.Curves.Rotate =
		mmath.SplitCurve(next.Curves.Rotate, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.Distance, next.Curves.Distance =
		mmath.SplitCurve(next.Curves.Distance, float32(prev.Index()), float32(index), float32(next.Index()))
	f.Curves.ViewOfAngle, next.Curves.ViewOfAngle =
		mmath.SplitCurve(next.Curves.ViewOfAngle, float32(prev.Index()), float32(index), float32(next.Index()))
}

// IsDefault は既定値か判定する。
func (f *CameraFrame) IsDefault() bool {
	if f == nil {
		return true
	}
	pos := vec3OrZero(f.Position)
	deg := vec3OrZero(f.Degrees)
	if !pos.NearEquals(mmath.Vec3{}, 1e-8) || !deg.NearEquals(mmath.Vec3{}, 1e-8) {
		return false
	}
	if f.Distance != 0 || f.ViewOfAngle != 0 || !f.IsPerspectiveOff {
		return false
	}
	if f.Curves == nil {
		return true
	}
	return isDefaultCameraCurves(f.Curves)
}

// CameraFrames はカメラフレーム集合を表す。
type CameraFrames struct {
	*BaseFrames[*CameraFrame]
}

// NewCameraFrames はCameraFramesを生成する。
func NewCameraFrames() *CameraFrames {
	return &CameraFrames{BaseFrames: NewBaseFrames(NewCameraFrame, nilCameraFrame)}
}

// Clean は既定値のみのフレームを削除する。
func (c *CameraFrames) Clean() {
	if c == nil || c.Len() != 1 {
		return
	}
	var frameIndex Frame
	var frameValue *CameraFrame
	c.ForEach(func(idx Frame, value *CameraFrame) bool {
		frameIndex = idx
		frameValue = value
		return false
	})
	if frameValue == nil || frameValue.IsDefault() {
		c.Delete(frameIndex)
	}
}

// Copy はフレーム集合を複製する。
func (c *CameraFrames) Copy() (*CameraFrames, error) {
	return deepCopy(c)
}

func nilCameraFrame() *CameraFrame {
	return nil
}

// cameraCurveT はカメラ補間係数を返す。
func cameraCurveT(prev, now, next Frame, curves *CameraCurves) (float64, float64, float64, float64, float64, float64) {
	if curves == nil {
		t := linearT(prev, now, next)
		return t, t, t, t, t, t
	}
	return curves.Evaluate(prev, now, next)
}

func isDefaultCameraCurves(curves *CameraCurves) bool {
	if curves == nil {
		return true
	}
	return isLinearCurve(curves.TranslateX) && isLinearCurve(curves.TranslateY) && isLinearCurve(curves.TranslateZ) &&
		isLinearCurve(curves.Rotate) && isLinearCurve(curves.Distance) && isLinearCurve(curves.ViewOfAngle)
}

func isLinearCurve(curve *mmath.Curve) bool {
	if curve == nil {
		return true
	}
	return curve.Start == (mmath.Vec2{X: 20, Y: 20}) && curve.End == (mmath.Vec2{X: 107, Y: 107})
}
