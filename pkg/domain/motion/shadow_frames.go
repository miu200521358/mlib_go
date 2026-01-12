// 指示: miu200521358
package motion

// ShadowFrame はシャドウフレームを表す。
type ShadowFrame struct {
	*BaseFrame
	ShadowMode int
	Distance   float64
}

// NewShadowFrame はShadowFrameを生成する。
func NewShadowFrame(index Frame) *ShadowFrame {
	return &ShadowFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (f *ShadowFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*ShadowFrame)(nil), nil
	}
	copied := &ShadowFrame{
		BaseFrame:  &BaseFrame{index: f.Index(), Read: f.Read},
		ShadowMode: f.ShadowMode,
		Distance:   f.Distance,
	}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *ShadowFrame) lerpFrame(prev *ShadowFrame, index Frame) *ShadowFrame {
	if prev == nil && next == nil {
		return NewShadowFrame(index)
	}
	if prev == nil {
		copied, _ := next.Copy()
		out := copied.(*ShadowFrame)
		out.SetIndex(index)
		return out
	}
	if next == nil {
		copied, _ := prev.Copy()
		out := copied.(*ShadowFrame)
		out.SetIndex(index)
		return out
	}
	out := NewShadowFrame(index)
	t := linearT(prev.Index(), index, next.Index())
	out.ShadowMode = prev.ShadowMode
	out.Distance = prev.Distance + (next.Distance-prev.Distance)*t
	return out
}

// splitCurve は何もしない。
func (f *ShadowFrame) splitCurve(prev *ShadowFrame, next *ShadowFrame, index Frame) {
}

// IsDefault は既定値か判定する。
func (f *ShadowFrame) IsDefault() bool {
	if f == nil {
		return true
	}
	return f.ShadowMode == 0 && f.Distance == 0
}

// ShadowFrames はシャドウフレーム集合を表す。
type ShadowFrames struct {
	*BaseFrames[*ShadowFrame]
}

// NewShadowFrames はShadowFramesを生成する。
func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{BaseFrames: NewBaseFrames(NewShadowFrame, nilShadowFrame)}
}

// Clean は既定値のみのフレームを削除する。
func (s *ShadowFrames) Clean() {
	if s == nil || s.Len() != 1 {
		return
	}
	var frameIndex Frame
	var frameValue *ShadowFrame
	s.ForEach(func(idx Frame, value *ShadowFrame) bool {
		frameIndex = idx
		frameValue = value
		return false
	})
	if frameValue == nil || frameValue.IsDefault() {
		s.Delete(frameIndex)
	}
}

// Copy はフレーム集合を複製する。
func (s *ShadowFrames) Copy() (*ShadowFrames, error) {
	return deepCopy(s)
}

func nilShadowFrame() *ShadowFrame {
	return nil
}
