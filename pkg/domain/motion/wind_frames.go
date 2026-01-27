// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

// WindEnabledFrame は風有効フレームを表す。
type WindEnabledFrame struct {
	*BaseFrame
	Enabled bool
}

// NewWindEnabledFrame はWindEnabledFrameを生成する。
func NewWindEnabledFrame(index Frame) *WindEnabledFrame {
	return &WindEnabledFrame{BaseFrame: NewBaseFrame(index), Enabled: false}
}

// Copy はフレームを複製する。
func (f *WindEnabledFrame) Copy() (WindEnabledFrame, error) {
	if f == nil {
		return WindEnabledFrame{}, nil
	}
	copied := WindEnabledFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, Enabled: f.Enabled}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindEnabledFrame) lerpFrame(prev *WindEnabledFrame, index Frame) *WindEnabledFrame {
	if prev == nil && next == nil {
		return NewWindEnabledFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindEnabledFrame) splitCurve(prev *WindEnabledFrame, next *WindEnabledFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindEnabledFrame) copyWithIndex(index Frame) *WindEnabledFrame {
	if f == nil {
		return nil
	}
	return &WindEnabledFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Enabled:   f.Enabled,
	}
}

// WindEnabledFrames は風有効フレーム集合を表す。
type WindEnabledFrames struct {
	*BaseFrames[*WindEnabledFrame]
}

// NewWindEnabledFrames はWindEnabledFramesを生成する。
func NewWindEnabledFrames() *WindEnabledFrames {
	return &WindEnabledFrames{BaseFrames: NewBaseFrames(NewWindEnabledFrame, nilWindEnabledFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindEnabledFrames) Get(frame Frame) *WindEnabledFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindEnabledFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindEnabledFrames) Copy() (WindEnabledFrames, error) {
	if w == nil {
		return WindEnabledFrames{}, nil
	}
	return deepCopy(*w)
}

// WindDirectionFrame は風向きフレームを表す。
type WindDirectionFrame struct {
	*BaseFrame
	Direction *mmath.Vec3
}

// NewWindDirectionFrame はWindDirectionFrameを生成する。
func NewWindDirectionFrame(index Frame) *WindDirectionFrame {
	return &WindDirectionFrame{BaseFrame: NewBaseFrame(index), Direction: &mmath.Vec3{}}
}

// Copy はフレームを複製する。
func (f *WindDirectionFrame) Copy() (WindDirectionFrame, error) {
	if f == nil {
		return WindDirectionFrame{}, nil
	}
	copied := WindDirectionFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, Direction: copyVec3(f.Direction)}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindDirectionFrame) lerpFrame(prev *WindDirectionFrame, index Frame) *WindDirectionFrame {
	if prev == nil && next == nil {
		return NewWindDirectionFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindDirectionFrame) splitCurve(prev *WindDirectionFrame, next *WindDirectionFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindDirectionFrame) copyWithIndex(index Frame) *WindDirectionFrame {
	if f == nil {
		return nil
	}
	return &WindDirectionFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Direction: copyVec3(f.Direction),
	}
}

// WindDirectionFrames は風向きフレーム集合を表す。
type WindDirectionFrames struct {
	*BaseFrames[*WindDirectionFrame]
}

// NewWindDirectionFrames はWindDirectionFramesを生成する。
func NewWindDirectionFrames() *WindDirectionFrames {
	return &WindDirectionFrames{BaseFrames: NewBaseFrames(NewWindDirectionFrame, nilWindDirectionFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindDirectionFrames) Get(frame Frame) *WindDirectionFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindDirectionFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindDirectionFrames) Copy() (WindDirectionFrames, error) {
	if w == nil {
		return WindDirectionFrames{}, nil
	}
	return deepCopy(*w)
}

// WindLiftCoeffFrame は風揚力係数フレームを表す。
type WindLiftCoeffFrame struct {
	*BaseFrame
	LiftCoeff float64
}

// NewWindLiftCoeffFrame はWindLiftCoeffFrameを生成する。
func NewWindLiftCoeffFrame(index Frame) *WindLiftCoeffFrame {
	return &WindLiftCoeffFrame{BaseFrame: NewBaseFrame(index), LiftCoeff: 60}
}

// WindLiftCoeff は1/係数を返す。
func (f *WindLiftCoeffFrame) WindLiftCoeff() float64 {
	if f == nil || f.LiftCoeff <= 0 {
		return 1
	}
	return 1.0 / f.LiftCoeff
}

// Copy はフレームを複製する。
func (f *WindLiftCoeffFrame) Copy() (WindLiftCoeffFrame, error) {
	if f == nil {
		return WindLiftCoeffFrame{}, nil
	}
	copied := WindLiftCoeffFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, LiftCoeff: f.LiftCoeff}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindLiftCoeffFrame) lerpFrame(prev *WindLiftCoeffFrame, index Frame) *WindLiftCoeffFrame {
	if prev == nil && next == nil {
		return NewWindLiftCoeffFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindLiftCoeffFrame) splitCurve(prev *WindLiftCoeffFrame, next *WindLiftCoeffFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindLiftCoeffFrame) copyWithIndex(index Frame) *WindLiftCoeffFrame {
	if f == nil {
		return nil
	}
	return &WindLiftCoeffFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		LiftCoeff: f.LiftCoeff,
	}
}

// WindLiftCoeffFrames は風揚力係数フレーム集合を表す。
type WindLiftCoeffFrames struct {
	*BaseFrames[*WindLiftCoeffFrame]
}

// NewWindLiftCoeffFrames はWindLiftCoeffFramesを生成する。
func NewWindLiftCoeffFrames() *WindLiftCoeffFrames {
	return &WindLiftCoeffFrames{BaseFrames: NewBaseFrames(NewWindLiftCoeffFrame, nilWindLiftCoeffFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindLiftCoeffFrames) Get(frame Frame) *WindLiftCoeffFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindLiftCoeffFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindLiftCoeffFrames) Copy() (WindLiftCoeffFrames, error) {
	if w == nil {
		return WindLiftCoeffFrames{}, nil
	}
	return deepCopy(*w)
}

// WindDragCoeffFrame は風抗力係数フレームを表す。
type WindDragCoeffFrame struct {
	*BaseFrame
	DragCoeff float64
}

// NewWindDragCoeffFrame はWindDragCoeffFrameを生成する。
func NewWindDragCoeffFrame(index Frame) *WindDragCoeffFrame {
	return &WindDragCoeffFrame{BaseFrame: NewBaseFrame(index), DragCoeff: 60}
}

// WindDragCoeff は1/係数を返す。
func (f *WindDragCoeffFrame) WindDragCoeff() float64 {
	if f == nil || f.DragCoeff <= 0 {
		return 1
	}
	return 1.0 / f.DragCoeff
}

// Copy はフレームを複製する。
func (f *WindDragCoeffFrame) Copy() (WindDragCoeffFrame, error) {
	if f == nil {
		return WindDragCoeffFrame{}, nil
	}
	copied := WindDragCoeffFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, DragCoeff: f.DragCoeff}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindDragCoeffFrame) lerpFrame(prev *WindDragCoeffFrame, index Frame) *WindDragCoeffFrame {
	if prev == nil && next == nil {
		return NewWindDragCoeffFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindDragCoeffFrame) splitCurve(prev *WindDragCoeffFrame, next *WindDragCoeffFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindDragCoeffFrame) copyWithIndex(index Frame) *WindDragCoeffFrame {
	if f == nil {
		return nil
	}
	return &WindDragCoeffFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		DragCoeff: f.DragCoeff,
	}
}

// WindDragCoeffFrames は風抗力係数フレーム集合を表す。
type WindDragCoeffFrames struct {
	*BaseFrames[*WindDragCoeffFrame]
}

// NewWindDragCoeffFrames はWindDragCoeffFramesを生成する。
func NewWindDragCoeffFrames() *WindDragCoeffFrames {
	return &WindDragCoeffFrames{BaseFrames: NewBaseFrames(NewWindDragCoeffFrame, nilWindDragCoeffFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindDragCoeffFrames) Get(frame Frame) *WindDragCoeffFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindDragCoeffFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindDragCoeffFrames) Copy() (WindDragCoeffFrames, error) {
	if w == nil {
		return WindDragCoeffFrames{}, nil
	}
	return deepCopy(*w)
}

// WindRandomnessFrame は風乱流係数フレームを表す。
type WindRandomnessFrame struct {
	*BaseFrame
	Randomness float64
}

// NewWindRandomnessFrame はWindRandomnessFrameを生成する。
func NewWindRandomnessFrame(index Frame) *WindRandomnessFrame {
	return &WindRandomnessFrame{BaseFrame: NewBaseFrame(index), Randomness: 0}
}

// WindRandomness は1/係数を返す。
func (f *WindRandomnessFrame) WindRandomness() float64 {
	if f == nil || f.Randomness <= 0 {
		return 0
	}
	return 1.0 / f.Randomness
}

// Copy はフレームを複製する。
func (f *WindRandomnessFrame) Copy() (WindRandomnessFrame, error) {
	if f == nil {
		return WindRandomnessFrame{}, nil
	}
	copied := WindRandomnessFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, Randomness: f.Randomness}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindRandomnessFrame) lerpFrame(prev *WindRandomnessFrame, index Frame) *WindRandomnessFrame {
	if prev == nil && next == nil {
		return NewWindRandomnessFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindRandomnessFrame) splitCurve(prev *WindRandomnessFrame, next *WindRandomnessFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindRandomnessFrame) copyWithIndex(index Frame) *WindRandomnessFrame {
	if f == nil {
		return nil
	}
	return &WindRandomnessFrame{
		BaseFrame:  &BaseFrame{index: index, Read: f.Read},
		Randomness: f.Randomness,
	}
}

// WindRandomnessFrames は風乱流係数フレーム集合を表す。
type WindRandomnessFrames struct {
	*BaseFrames[*WindRandomnessFrame]
}

// NewWindRandomnessFrames はWindRandomnessFramesを生成する。
func NewWindRandomnessFrames() *WindRandomnessFrames {
	return &WindRandomnessFrames{BaseFrames: NewBaseFrames(NewWindRandomnessFrame, nilWindRandomnessFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindRandomnessFrames) Get(frame Frame) *WindRandomnessFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindRandomnessFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindRandomnessFrames) Copy() (WindRandomnessFrames, error) {
	if w == nil {
		return WindRandomnessFrames{}, nil
	}
	return deepCopy(*w)
}

// WindSpeedFrame は風速フレームを表す。
type WindSpeedFrame struct {
	*BaseFrame
	Speed float64
}

// NewWindSpeedFrame はWindSpeedFrameを生成する。
func NewWindSpeedFrame(index Frame) *WindSpeedFrame {
	return &WindSpeedFrame{BaseFrame: NewBaseFrame(index), Speed: 60}
}

// WindSpeed は1/係数を返す。
func (f *WindSpeedFrame) WindSpeed() float64 {
	if f == nil || f.Speed <= 0 {
		return 1
	}
	return 1.0 / f.Speed
}

// Copy はフレームを複製する。
func (f *WindSpeedFrame) Copy() (WindSpeedFrame, error) {
	if f == nil {
		return WindSpeedFrame{}, nil
	}
	copied := WindSpeedFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, Speed: f.Speed}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindSpeedFrame) lerpFrame(prev *WindSpeedFrame, index Frame) *WindSpeedFrame {
	if prev == nil && next == nil {
		return NewWindSpeedFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindSpeedFrame) splitCurve(prev *WindSpeedFrame, next *WindSpeedFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindSpeedFrame) copyWithIndex(index Frame) *WindSpeedFrame {
	if f == nil {
		return nil
	}
	return &WindSpeedFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Speed:     f.Speed,
	}
}

// WindSpeedFrames は風速フレーム集合を表す。
type WindSpeedFrames struct {
	*BaseFrames[*WindSpeedFrame]
}

// NewWindSpeedFrames はWindSpeedFramesを生成する。
func NewWindSpeedFrames() *WindSpeedFrames {
	return &WindSpeedFrames{BaseFrames: NewBaseFrames(NewWindSpeedFrame, nilWindSpeedFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindSpeedFrames) Get(frame Frame) *WindSpeedFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindSpeedFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindSpeedFrames) Copy() (WindSpeedFrames, error) {
	if w == nil {
		return WindSpeedFrames{}, nil
	}
	return deepCopy(*w)
}

// WindTurbulenceFreqHzFrame は風乱流周波数フレームを表す。
type WindTurbulenceFreqHzFrame struct {
	*BaseFrame
	TurbulenceFreqHz float64
}

// NewWindTurbulenceFreqHzFrame はWindTurbulenceFreqHzFrameを生成する。
func NewWindTurbulenceFreqHzFrame(index Frame) *WindTurbulenceFreqHzFrame {
	return &WindTurbulenceFreqHzFrame{BaseFrame: NewBaseFrame(index), TurbulenceFreqHz: 0}
}

// WindTurbulenceFreqHz は1/係数を返す。
func (f *WindTurbulenceFreqHzFrame) WindTurbulenceFreqHz() float64 {
	if f == nil || f.TurbulenceFreqHz <= 0 {
		return 0
	}
	return 1.0 / f.TurbulenceFreqHz
}

// Copy はフレームを複製する。
func (f *WindTurbulenceFreqHzFrame) Copy() (WindTurbulenceFreqHzFrame, error) {
	if f == nil {
		return WindTurbulenceFreqHzFrame{}, nil
	}
	copied := WindTurbulenceFreqHzFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, TurbulenceFreqHz: f.TurbulenceFreqHz}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *WindTurbulenceFreqHzFrame) lerpFrame(prev *WindTurbulenceFreqHzFrame, index Frame) *WindTurbulenceFreqHzFrame {
	if prev == nil && next == nil {
		return NewWindTurbulenceFreqHzFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *WindTurbulenceFreqHzFrame) splitCurve(prev *WindTurbulenceFreqHzFrame, next *WindTurbulenceFreqHzFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *WindTurbulenceFreqHzFrame) copyWithIndex(index Frame) *WindTurbulenceFreqHzFrame {
	if f == nil {
		return nil
	}
	return &WindTurbulenceFreqHzFrame{
		BaseFrame:        &BaseFrame{index: index, Read: f.Read},
		TurbulenceFreqHz: f.TurbulenceFreqHz,
	}
}

// WindTurbulenceFreqHzFrames は風乱流周波数フレーム集合を表す。
type WindTurbulenceFreqHzFrames struct {
	*BaseFrames[*WindTurbulenceFreqHzFrame]
}

// NewWindTurbulenceFreqHzFrames はWindTurbulenceFreqHzFramesを生成する。
func NewWindTurbulenceFreqHzFrames() *WindTurbulenceFreqHzFrames {
	return &WindTurbulenceFreqHzFrames{BaseFrames: NewBaseFrames(NewWindTurbulenceFreqHzFrame, nilWindTurbulenceFreqHzFrame)}
}

// Get は次フレーム優先で値を返す。
func (w *WindTurbulenceFreqHzFrames) Get(frame Frame) *WindTurbulenceFreqHzFrame {
	if w == nil {
		return nil
	}
	if w.Len() == 0 {
		return NewWindTurbulenceFreqHzFrame(frame)
	}
	if w.Has(frame) {
		return w.frames[frame]
	}
	if next, ok := w.NextFrame(frame); ok && next > frame {
		return w.frames[next]
	}
	prev, _ := w.PrevFrame(frame)
	return w.frames[prev]
}

// Copy はフレーム集合を複製する。
func (w *WindTurbulenceFreqHzFrames) Copy() (WindTurbulenceFreqHzFrames, error) {
	if w == nil {
		return WindTurbulenceFreqHzFrames{}, nil
	}
	return deepCopy(*w)
}

// nilWindEnabledFrame は空の風有効化フレームを返す。
func nilWindEnabledFrame() *WindEnabledFrame {
	return nil
}

// nilWindDirectionFrame は空の風向きフレームを返す。
func nilWindDirectionFrame() *WindDirectionFrame {
	return nil
}

// nilWindLiftCoeffFrame は空の揚力係数フレームを返す。
func nilWindLiftCoeffFrame() *WindLiftCoeffFrame {
	return nil
}

// nilWindDragCoeffFrame は空の抗力係数フレームを返す。
func nilWindDragCoeffFrame() *WindDragCoeffFrame {
	return nil
}

// nilWindRandomnessFrame は空の乱流係数フレームを返す。
func nilWindRandomnessFrame() *WindRandomnessFrame {
	return nil
}

// nilWindSpeedFrame は空の風速フレームを返す。
func nilWindSpeedFrame() *WindSpeedFrame {
	return nil
}

// nilWindTurbulenceFreqHzFrame は空の乱流周波数フレームを返す。
func nilWindTurbulenceFreqHzFrame() *WindTurbulenceFreqHzFrame {
	return nil
}
