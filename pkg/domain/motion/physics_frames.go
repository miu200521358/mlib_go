// 指示: miu200521358
package motion

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/performance"
)

// PhysicsResetType は物理リセット種別を表す。
type PhysicsResetType int

const (
	// PHYSICS_RESET_TYPE_NONE はリセットなし。
	PHYSICS_RESET_TYPE_NONE PhysicsResetType = iota
	// PHYSICS_RESET_TYPE_CONTINUE_FRAME は継続フレーム用リセット。
	PHYSICS_RESET_TYPE_CONTINUE_FRAME
	// PHYSICS_RESET_TYPE_START_FRAME は開始フレーム用リセット。
	PHYSICS_RESET_TYPE_START_FRAME
	// PHYSICS_RESET_TYPE_START_FIT_FRAME は開始フレーム用リセット（Yスタンス）。
	PHYSICS_RESET_TYPE_START_FIT_FRAME
)

// MaxSubStepsFrame は最大演算回数フレームを表す。
type MaxSubStepsFrame struct {
	*BaseFrame
	MaxSubSteps int
}

// NewMaxSubStepsFrame はMaxSubStepsFrameを生成する。
func NewMaxSubStepsFrame(index Frame) *MaxSubStepsFrame {
	return &MaxSubStepsFrame{BaseFrame: NewBaseFrame(index), MaxSubSteps: performance.DefaultMaxSubSteps}
}

// Copy はフレームを複製する。
func (f *MaxSubStepsFrame) Copy() (MaxSubStepsFrame, error) {
	if f == nil {
		return MaxSubStepsFrame{}, nil
	}
	copied := MaxSubStepsFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, MaxSubSteps: f.MaxSubSteps}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *MaxSubStepsFrame) lerpFrame(prev *MaxSubStepsFrame, index Frame) *MaxSubStepsFrame {
	if prev == nil && next == nil {
		return NewMaxSubStepsFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *MaxSubStepsFrame) splitCurve(prev *MaxSubStepsFrame, next *MaxSubStepsFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *MaxSubStepsFrame) copyWithIndex(index Frame) *MaxSubStepsFrame {
	if f == nil {
		return nil
	}
	return &MaxSubStepsFrame{
		BaseFrame:   &BaseFrame{index: index, Read: f.Read},
		MaxSubSteps: f.MaxSubSteps,
	}
}

// MaxSubStepsFrames は最大演算回数フレーム集合を表す。
type MaxSubStepsFrames struct {
	*BaseFrames[*MaxSubStepsFrame]
}

// NewMaxSubStepsFrames はMaxSubStepsFramesを生成する。
func NewMaxSubStepsFrames() *MaxSubStepsFrames {
	return &MaxSubStepsFrames{BaseFrames: NewBaseFrames(NewMaxSubStepsFrame, nilMaxSubStepsFrame)}
}

// Get は次フレーム優先で値を返す。
func (m *MaxSubStepsFrames) Get(frame Frame) *MaxSubStepsFrame {
	if m == nil {
		return nil
	}
	if m.Len() == 0 {
		return NewMaxSubStepsFrame(frame)
	}
	if m.Has(frame) {
		return m.frames[frame]
	}
	if next, ok := m.NextFrame(frame); ok && next > frame {
		return m.frames[next]
	}
	prev, _ := m.PrevFrame(frame)
	return m.frames[prev]
}

// Copy はフレーム集合を複製する。
func (m *MaxSubStepsFrames) Copy() (MaxSubStepsFrames, error) {
	if m == nil {
		return MaxSubStepsFrames{}, nil
	}
	return deepCopy(*m)
}

// FixedTimeStepFrame は演算頻度フレームを表す。
type FixedTimeStepFrame struct {
	*BaseFrame
	FixedTimeStepNum float64
}

// NewFixedTimeStepFrame はFixedTimeStepFrameを生成する。
func NewFixedTimeStepFrame(index Frame) *FixedTimeStepFrame {
	return &FixedTimeStepFrame{BaseFrame: NewBaseFrame(index), FixedTimeStepNum: 60}
}

// FixedTimeStep は1/FixedTimeStepNumを返す。
func (f *FixedTimeStepFrame) FixedTimeStep() float64 {
	if f == nil || f.FixedTimeStepNum <= 0 {
		return 1.0 / 60.0
	}
	return 1.0 / f.FixedTimeStepNum
}

// Copy はフレームを複製する。
func (f *FixedTimeStepFrame) Copy() (FixedTimeStepFrame, error) {
	if f == nil {
		return FixedTimeStepFrame{}, nil
	}
	copied := FixedTimeStepFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, FixedTimeStepNum: f.FixedTimeStepNum}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *FixedTimeStepFrame) lerpFrame(prev *FixedTimeStepFrame, index Frame) *FixedTimeStepFrame {
	if prev == nil && next == nil {
		return NewFixedTimeStepFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *FixedTimeStepFrame) splitCurve(prev *FixedTimeStepFrame, next *FixedTimeStepFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *FixedTimeStepFrame) copyWithIndex(index Frame) *FixedTimeStepFrame {
	if f == nil {
		return nil
	}
	return &FixedTimeStepFrame{
		BaseFrame:        &BaseFrame{index: index, Read: f.Read},
		FixedTimeStepNum: f.FixedTimeStepNum,
	}
}

// FixedTimeStepFrames は演算頻度フレーム集合を表す。
type FixedTimeStepFrames struct {
	*BaseFrames[*FixedTimeStepFrame]
}

// NewFixedTimeStepFrames はFixedTimeStepFramesを生成する。
func NewFixedTimeStepFrames() *FixedTimeStepFrames {
	return &FixedTimeStepFrames{BaseFrames: NewBaseFrames(NewFixedTimeStepFrame, nilFixedTimeStepFrame)}
}

// Get は次フレーム優先で値を返す。
func (m *FixedTimeStepFrames) Get(frame Frame) *FixedTimeStepFrame {
	if m == nil {
		return nil
	}
	if m.Len() == 0 {
		return NewFixedTimeStepFrame(frame)
	}
	if m.Has(frame) {
		return m.frames[frame]
	}
	if next, ok := m.NextFrame(frame); ok && next > frame {
		return m.frames[next]
	}
	prev, _ := m.PrevFrame(frame)
	return m.frames[prev]
}

// Copy はフレーム集合を複製する。
func (m *FixedTimeStepFrames) Copy() (FixedTimeStepFrames, error) {
	if m == nil {
		return FixedTimeStepFrames{}, nil
	}
	return deepCopy(*m)
}

// GravityFrame は重力フレームを表す。
type GravityFrame struct {
	*BaseFrame
	Gravity *mmath.Vec3
}

// NewGravityFrame はGravityFrameを生成する。
func NewGravityFrame(index Frame) *GravityFrame {
	return &GravityFrame{BaseFrame: NewBaseFrame(index), Gravity: vec3Ptr(0, -9.8, 0)}
}

// Copy はフレームを複製する。
func (f *GravityFrame) Copy() (GravityFrame, error) {
	if f == nil {
		return GravityFrame{}, nil
	}
	copied := GravityFrame{
		BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read},
		Gravity:   copyVec3(f.Gravity),
	}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *GravityFrame) lerpFrame(prev *GravityFrame, index Frame) *GravityFrame {
	if prev == nil && next == nil {
		return NewGravityFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *GravityFrame) splitCurve(prev *GravityFrame, next *GravityFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *GravityFrame) copyWithIndex(index Frame) *GravityFrame {
	if f == nil {
		return nil
	}
	return &GravityFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Gravity:   copyVec3(f.Gravity),
	}
}

// GravityFrames は重力フレーム集合を表す。
type GravityFrames struct {
	*BaseFrames[*GravityFrame]
}

// NewGravityFrames はGravityFramesを生成する。
func NewGravityFrames() *GravityFrames {
	return &GravityFrames{BaseFrames: NewBaseFrames(NewGravityFrame, nilGravityFrame)}
}

// Get は次フレーム優先で値を返す。
func (m *GravityFrames) Get(frame Frame) *GravityFrame {
	if m == nil {
		return nil
	}
	if m.Len() == 0 {
		return NewGravityFrame(frame)
	}
	if m.Has(frame) {
		return m.frames[frame]
	}
	if next, ok := m.NextFrame(frame); ok && next > frame {
		return m.frames[next]
	}
	prev, _ := m.PrevFrame(frame)
	return m.frames[prev]
}

// Copy はフレーム集合を複製する。
func (m *GravityFrames) Copy() (GravityFrames, error) {
	if m == nil {
		return GravityFrames{}, nil
	}
	return deepCopy(*m)
}

// PhysicsResetFrame は物理リセットフレームを表す。
type PhysicsResetFrame struct {
	*BaseFrame
	PhysicsResetType PhysicsResetType
}

// NewPhysicsResetFrame はPhysicsResetFrameを生成する。
func NewPhysicsResetFrame(index Frame) *PhysicsResetFrame {
	return &PhysicsResetFrame{BaseFrame: NewBaseFrame(index), PhysicsResetType: PHYSICS_RESET_TYPE_NONE}
}

// Copy はフレームを複製する。
func (f *PhysicsResetFrame) Copy() (PhysicsResetFrame, error) {
	if f == nil {
		return PhysicsResetFrame{}, nil
	}
	copied := PhysicsResetFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, PhysicsResetType: f.PhysicsResetType}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *PhysicsResetFrame) lerpFrame(prev *PhysicsResetFrame, index Frame) *PhysicsResetFrame {
	if prev == nil && next == nil {
		return NewPhysicsResetFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	return prev.copyWithIndex(index)
}

// splitCurve は何もしない。
func (f *PhysicsResetFrame) splitCurve(prev *PhysicsResetFrame, next *PhysicsResetFrame, index Frame) {
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *PhysicsResetFrame) copyWithIndex(index Frame) *PhysicsResetFrame {
	if f == nil {
		return nil
	}
	return &PhysicsResetFrame{
		BaseFrame:        &BaseFrame{index: index, Read: f.Read},
		PhysicsResetType: f.PhysicsResetType,
	}
}

// PhysicsResetFrames は物理リセットフレーム集合を表す。
type PhysicsResetFrames struct {
	*BaseFrames[*PhysicsResetFrame]
}

// NewPhysicsResetFrames はPhysicsResetFramesを生成する。
func NewPhysicsResetFrames() *PhysicsResetFrames {
	return &PhysicsResetFrames{BaseFrames: NewBaseFrames(NewPhysicsResetFrame, nilPhysicsResetFrame)}
}

// Get は次フレーム優先で値を返す。
func (m *PhysicsResetFrames) Get(frame Frame) *PhysicsResetFrame {
	if m == nil {
		return nil
	}
	if m.Len() == 0 {
		return NewPhysicsResetFrame(frame)
	}
	if m.Has(frame) {
		return m.frames[frame]
	}
	if next, ok := m.NextFrame(frame); ok && next > frame {
		return m.frames[next]
	}
	prev, _ := m.PrevFrame(frame)
	return m.frames[prev]
}

// Copy はフレーム集合を複製する。
func (m *PhysicsResetFrames) Copy() (PhysicsResetFrames, error) {
	if m == nil {
		return PhysicsResetFrames{}, nil
	}
	return deepCopy(*m)
}

// nilMaxSubStepsFrame は既定の空フレームを返す。
func nilMaxSubStepsFrame() *MaxSubStepsFrame {
	return nil
}

// nilFixedTimeStepFrame は既定の空フレームを返す。
func nilFixedTimeStepFrame() *FixedTimeStepFrame {
	return nil
}

// nilGravityFrame は既定の空フレームを返す。
func nilGravityFrame() *GravityFrame {
	return nil
}

// nilPhysicsResetFrame は既定の空フレームを返す。
func nilPhysicsResetFrame() *PhysicsResetFrame {
	return nil
}
