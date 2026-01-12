// 指示: miu200521358
package motion

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// RigidBodyFrame は剛体フレームを表す。
type RigidBodyFrame struct {
	*BaseFrame
	Position *mmath.Vec3
	Size     *mmath.Vec3
	Mass     float64
}

// NewRigidBodyFrame はRigidBodyFrameを生成する。
func NewRigidBodyFrame(index Frame) *RigidBodyFrame {
	return &RigidBodyFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (f *RigidBodyFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*RigidBodyFrame)(nil), nil
	}
	copied := &RigidBodyFrame{
		BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read},
		Position:  copyVec3(f.Position),
		Size:      copyVec3(f.Size),
		Mass:      f.Mass,
	}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *RigidBodyFrame) lerpFrame(prev *RigidBodyFrame, index Frame) *RigidBodyFrame {
	if prev == nil && next == nil {
		return NewRigidBodyFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	if next == nil {
		return prev.copyWithIndex(index)
	}
	out := NewRigidBodyFrame(index)
	t := linearT(prev.Index(), index, next.Index())
	prevPos := vec3OrZero(prev.Position)
	nextPos := vec3OrZero(next.Position)
	prevSize := vec3OrZero(prev.Size)
	nextSize := vec3OrZero(next.Size)
	pos := prevPos.Lerp(nextPos, t)
	size := prevSize.Lerp(nextSize, t)
	out.Position = &pos
	out.Size = &size
	out.Mass = mmath.Lerp(prev.Mass, next.Mass, t)
	return out
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *RigidBodyFrame) copyWithIndex(index Frame) *RigidBodyFrame {
	if f == nil {
		return nil
	}
	return &RigidBodyFrame{
		BaseFrame: &BaseFrame{index: index, Read: f.Read},
		Position:  copyVec3(f.Position),
		Size:      copyVec3(f.Size),
		Mass:      f.Mass,
	}
}

// splitCurve は何もしない。
func (f *RigidBodyFrame) splitCurve(prev *RigidBodyFrame, next *RigidBodyFrame, index Frame) {
}

// RigidBodyNameFrames は剛体名ごとのフレーム集合を表す。
type RigidBodyNameFrames struct {
	*BaseFrames[*RigidBodyFrame]
	Name string
}

// NewRigidBodyNameFrames はRigidBodyNameFramesを生成する。
func NewRigidBodyNameFrames(name string) *RigidBodyNameFrames {
	return &RigidBodyNameFrames{
		BaseFrames: NewBaseFrames(newRigidBodyFrame, nilRigidBodyFrame),
		Name:       name,
	}
}

// Copy はフレーム集合を複製する。
func (r *RigidBodyNameFrames) Copy() (*RigidBodyNameFrames, error) {
	return deepCopy(r)
}

// RigidBodyFrames は剛体名ごとの集合を表す。
type RigidBodyFrames struct {
	names     []string
	nameIndex map[string]int
	values    []*RigidBodyNameFrames
}

// NewRigidBodyFrames はRigidBodyFramesを生成する。
func NewRigidBodyFrames() *RigidBodyFrames {
	return &RigidBodyFrames{
		names:     make([]string, 0),
		nameIndex: make(map[string]int),
		values:    make([]*RigidBodyNameFrames, 0),
	}
}

// Names は登録順の名前一覧を返す。
func (r *RigidBodyFrames) Names() []string {
	if r == nil {
		return nil
	}
	return append([]string(nil), r.names...)
}

// Get は名前に対応するフレーム集合を返す。
func (r *RigidBodyFrames) Get(name string) *RigidBodyNameFrames {
	if r == nil {
		return nil
	}
	if idx, ok := r.nameIndex[name]; ok {
		return r.values[idx]
	}
	return nil
}

// Update はフレーム集合を更新する。
func (r *RigidBodyFrames) Update(frames *RigidBodyNameFrames) {
	if r == nil || frames == nil {
		return
	}
	if idx, ok := r.nameIndex[frames.Name]; ok {
		r.values[idx] = frames
		return
	}
	r.nameIndex[frames.Name] = len(r.values)
	r.names = append(r.names, frames.Name)
	r.values = append(r.values, frames)
}

// Delete は名前を削除する。
func (r *RigidBodyFrames) Delete(name string) {
	if r == nil {
		return
	}
	idx, ok := r.nameIndex[name]
	if !ok {
		return
	}
	r.names = append(r.names[:idx], r.names[idx+1:]...)
	r.values = append(r.values[:idx], r.values[idx+1:]...)
	delete(r.nameIndex, name)
	for i := idx; i < len(r.names); i++ {
		r.nameIndex[r.names[i]] = i
	}
}

// Indexes は全トラックのフレーム番号を返す。
func (r *RigidBodyFrames) Indexes() []int {
	if r == nil {
		return nil
	}
	indexes := make([]int, 0)
	for _, frames := range r.values {
		frames.ForEach(func(frame Frame, _ *RigidBodyFrame) bool {
			indexes = append(indexes, int(frame))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

// Len は全フレーム数を返す。
func (r *RigidBodyFrames) Len() int {
	if r == nil {
		return 0
	}
	count := 0
	for _, frames := range r.values {
		count += frames.Len()
	}
	return count
}

// MaxFrame は最大フレーム番号を返す。
func (r *RigidBodyFrames) MaxFrame() Frame {
	if r == nil {
		return 0
	}
	maxFrame := Frame(0)
	for _, frames := range r.values {
		frame := frames.MaxFrame()
		if frame > maxFrame {
			maxFrame = frame
		}
	}
	return maxFrame
}

// MinFrame は最小フレーム番号を返す。
func (r *RigidBodyFrames) MinFrame() Frame {
	if r == nil {
		return 0
	}
	minFrame := Frame(math.MaxFloat32)
	for _, frames := range r.values {
		frame := frames.MinFrame()
		if frame < minFrame {
			minFrame = frame
		}
	}
	if minFrame == Frame(math.MaxFloat32) {
		return 0
	}
	return minFrame
}

// Copy は剛体集合を複製する。
func (r *RigidBodyFrames) Copy() (*RigidBodyFrames, error) {
	copied, err := deepCopy(r)
	if err != nil {
		return nil, err
	}
	copied.rebuildNameIndex()
	return copied, nil
}

func (r *RigidBodyFrames) rebuildNameIndex() {
	r.nameIndex = make(map[string]int, len(r.names))
	for i, name := range r.names {
		r.nameIndex[name] = i
	}
}

func newRigidBodyFrame(frame Frame) *RigidBodyFrame {
	return nil
}

func nilRigidBodyFrame() *RigidBodyFrame {
	return nil
}
