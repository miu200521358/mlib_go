// 指示: miu200521358
package motion

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

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
	return &MaxSubStepsFrame{BaseFrame: NewBaseFrame(index), MaxSubSteps: 2}
}

// Copy はフレームを複製する。
func (f *MaxSubStepsFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*MaxSubStepsFrame)(nil), nil
	}
	copied := &MaxSubStepsFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, MaxSubSteps: f.MaxSubSteps}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *MaxSubStepsFrame) lerpFrame(prev *MaxSubStepsFrame, index Frame) *MaxSubStepsFrame {
	if prev == nil && next == nil {
		return NewMaxSubStepsFrame(index)
	}
	if prev == nil {
		copied, _ := next.Copy()
		out := copied.(*MaxSubStepsFrame)
		out.SetIndex(index)
		return out
	}
	copied, _ := prev.Copy()
	out := copied.(*MaxSubStepsFrame)
	out.SetIndex(index)
	return out
}

// splitCurve は何もしない。
func (f *MaxSubStepsFrame) splitCurve(prev *MaxSubStepsFrame, next *MaxSubStepsFrame, index Frame) {
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
func (m *MaxSubStepsFrames) Copy() (*MaxSubStepsFrames, error) {
	return deepCopy(m)
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
func (f *FixedTimeStepFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*FixedTimeStepFrame)(nil), nil
	}
	copied := &FixedTimeStepFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, FixedTimeStepNum: f.FixedTimeStepNum}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *FixedTimeStepFrame) lerpFrame(prev *FixedTimeStepFrame, index Frame) *FixedTimeStepFrame {
	if prev == nil && next == nil {
		return NewFixedTimeStepFrame(index)
	}
	if prev == nil {
		copied, _ := next.Copy()
		out := copied.(*FixedTimeStepFrame)
		out.SetIndex(index)
		return out
	}
	copied, _ := prev.Copy()
	out := copied.(*FixedTimeStepFrame)
	out.SetIndex(index)
	return out
}

// splitCurve は何もしない。
func (f *FixedTimeStepFrame) splitCurve(prev *FixedTimeStepFrame, next *FixedTimeStepFrame, index Frame) {
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
func (m *FixedTimeStepFrames) Copy() (*FixedTimeStepFrames, error) {
	return deepCopy(m)
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
func (f *GravityFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*GravityFrame)(nil), nil
	}
	copied := &GravityFrame{
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
		copied, _ := next.Copy()
		out := copied.(*GravityFrame)
		out.SetIndex(index)
		return out
	}
	copied, _ := prev.Copy()
	out := copied.(*GravityFrame)
	out.SetIndex(index)
	return out
}

// splitCurve は何もしない。
func (f *GravityFrame) splitCurve(prev *GravityFrame, next *GravityFrame, index Frame) {
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
func (m *GravityFrames) Copy() (*GravityFrames, error) {
	return deepCopy(m)
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
func (f *PhysicsResetFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*PhysicsResetFrame)(nil), nil
	}
	copied := &PhysicsResetFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, PhysicsResetType: f.PhysicsResetType}
	return copied, nil
}

// lerpFrame は補間せず前フレームを複製する。
func (next *PhysicsResetFrame) lerpFrame(prev *PhysicsResetFrame, index Frame) *PhysicsResetFrame {
	if prev == nil && next == nil {
		return NewPhysicsResetFrame(index)
	}
	if prev == nil {
		copied, _ := next.Copy()
		out := copied.(*PhysicsResetFrame)
		out.SetIndex(index)
		return out
	}
	copied, _ := prev.Copy()
	out := copied.(*PhysicsResetFrame)
	out.SetIndex(index)
	return out
}

// splitCurve は何もしない。
func (f *PhysicsResetFrame) splitCurve(prev *PhysicsResetFrame, next *PhysicsResetFrame, index Frame) {
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
func (m *PhysicsResetFrames) Copy() (*PhysicsResetFrames, error) {
	return deepCopy(m)
}

func nilMaxSubStepsFrame() *MaxSubStepsFrame {
	return nil
}

func nilFixedTimeStepFrame() *FixedTimeStepFrame {
	return nil
}

func nilGravityFrame() *GravityFrame {
	return nil
}

func nilPhysicsResetFrame() *PhysicsResetFrame {
	return nil
}
