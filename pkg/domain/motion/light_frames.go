// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// LightFrame はライトフレームを表す。
type LightFrame struct {
	*BaseFrame
	Position mmath.Vec3
	Color    mmath.Vec3
}

// NewLightFrame はLightFrameを生成する。
func NewLightFrame(index Frame) *LightFrame {
	return &LightFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (f *LightFrame) Copy() (LightFrame, error) {
	if f == nil {
		return LightFrame{}, nil
	}
	copied := LightFrame{
		BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read},
		Position:  f.Position,
		Color:     f.Color,
	}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *LightFrame) lerpFrame(prev *LightFrame, index Frame) *LightFrame {
	if prev == nil && next == nil {
		return NewLightFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	if next == nil {
		return prev.copyWithIndex(index)
	}
	out := NewLightFrame(index)
	t := linearT(prev.Index(), index, next.Index())
	out.Position = prev.Position.Lerp(next.Position, t)
	out.Color = prev.Color.Lerp(next.Color, t)
	return out
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *LightFrame) copyWithIndex(index Frame) *LightFrame {
	if f == nil {
		return nil
	}
	return &LightFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Position:  f.Position,
		Color:     f.Color,
	}
}

// splitCurve は何もしない。
func (f *LightFrame) splitCurve(prev *LightFrame, next *LightFrame, index Frame) {
}

// IsDefault は既定値か判定する。
func (f *LightFrame) IsDefault() bool {
	if f == nil {
		return true
	}
	return f.Position.NearEquals(mmath.Vec3{}, 1e-8) && f.Color.NearEquals(mmath.Vec3{}, 1e-8)
}

// LightFrames はライトフレーム集合を表す。
type LightFrames struct {
	*BaseFrames[*LightFrame]
}

// NewLightFrames はLightFramesを生成する。
func NewLightFrames() *LightFrames {
	return &LightFrames{BaseFrames: NewBaseFrames(NewLightFrame, nilLightFrame)}
}

// Clean は既定値のみのフレームを削除する。
func (l *LightFrames) Clean() {
	if l == nil || l.Len() != 1 {
		return
	}
	var frameIndex Frame
	var frameValue *LightFrame
	l.ForEach(func(idx Frame, value *LightFrame) bool {
		frameIndex = idx
		frameValue = value
		return false
	})
	if frameValue == nil || frameValue.IsDefault() {
		l.Delete(frameIndex)
	}
}

// Copy はフレーム集合を複製する。
func (l *LightFrames) Copy() (LightFrames, error) {
	if l == nil {
		return LightFrames{}, nil
	}
	return deepCopy(*l)
}

// nilLightFrame は既定の空フレームを返す。
func nilLightFrame() *LightFrame {
	return nil
}
