// 指示: miu200521358
package motion

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// MorphFrame はモーフフレームを表す。
type MorphFrame struct {
	*BaseFrame
	Ratio float64
}

// NewMorphFrame はMorphFrameを生成する。
func NewMorphFrame(index Frame) *MorphFrame {
	return &MorphFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (f *MorphFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*MorphFrame)(nil), nil
	}
	copied := &MorphFrame{BaseFrame: &BaseFrame{index: f.Index(), Read: f.Read}, Ratio: f.Ratio}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *MorphFrame) lerpFrame(prev *MorphFrame, index Frame) *MorphFrame {
	if prev == nil && next == nil {
		return NewMorphFrame(index)
	}
	if prev == nil {
		return next.copyWithIndex(index)
	}
	if next == nil {
		return prev.copyWithIndex(index)
	}
	out := NewMorphFrame(index)
	t := linearT(prev.Index(), index, next.Index())
	out.Ratio = mmath.Effective(mmath.Lerp(prev.Ratio, next.Ratio, t))
	return out
}

// copyWithIndex は指定フレーム番号で複製する。
func (f *MorphFrame) copyWithIndex(index Frame) *MorphFrame {
	if f == nil {
		return nil
	}
	return &MorphFrame{BaseFrame: &BaseFrame{index: index, Read: f.Read}, Ratio: f.Ratio}
}

// splitCurve は何もしない。
func (f *MorphFrame) splitCurve(prev *MorphFrame, next *MorphFrame, index Frame) {
}

// MorphNameFrames はモーフ名ごとのフレーム集合を表す。
type MorphNameFrames struct {
	*BaseFrames[*MorphFrame]
	Name string
}

// NewMorphNameFrames はMorphNameFramesを生成する。
func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		BaseFrames: NewBaseFrames(NewMorphFrame, nilMorphFrame),
		Name:       name,
	}
}

// ContainsActive は有効なキーフレが存在するか判定する。
func (m *MorphNameFrames) ContainsActive() bool {
	if m == nil || m.Len() == 0 {
		return false
	}
	active := false
	m.ForEach(func(_ Frame, mf *MorphFrame) bool {
		if mf == nil {
			return true
		}
		if math.Abs(mf.Ratio) > 1e-2 {
			active = true
			return false
		}
		return true
	})
	return active
}

// Copy はフレーム集合を複製する。
func (m *MorphNameFrames) Copy() (*MorphNameFrames, error) {
	return deepCopy(m)
}

// MorphFrames はモーフ名ごとの集合を表す。
type MorphFrames struct {
	names     []string
	nameIndex map[string]int
	values    []*MorphNameFrames
}

// NewMorphFrames はMorphFramesを生成する。
func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		names:     make([]string, 0),
		nameIndex: make(map[string]int),
		values:    make([]*MorphNameFrames, 0),
	}
}

// Names は登録順の名前一覧を返す。
func (m *MorphFrames) Names() []string {
	if m == nil {
		return nil
	}
	return append([]string(nil), m.names...)
}

// Get は名前に対応するフレーム集合を返す。
func (m *MorphFrames) Get(name string) *MorphNameFrames {
	if m == nil {
		return nil
	}
	if idx, ok := m.nameIndex[name]; ok {
		return m.values[idx]
	}
	frames := NewMorphNameFrames(name)
	m.nameIndex[name] = len(m.values)
	m.names = append(m.names, name)
	m.values = append(m.values, frames)
	return frames
}

// Update はフレーム集合を更新する。
func (m *MorphFrames) Update(frames *MorphNameFrames) {
	if m == nil || frames == nil {
		return
	}
	if idx, ok := m.nameIndex[frames.Name]; ok {
		m.values[idx] = frames
		return
	}
	m.nameIndex[frames.Name] = len(m.values)
	m.names = append(m.names, frames.Name)
	m.values = append(m.values, frames)
}

// Delete は名前を削除する。
func (m *MorphFrames) Delete(name string) {
	if m == nil {
		return
	}
	idx, ok := m.nameIndex[name]
	if !ok {
		return
	}
	m.names = append(m.names[:idx], m.names[idx+1:]...)
	m.values = append(m.values[:idx], m.values[idx+1:]...)
	delete(m.nameIndex, name)
	for i := idx; i < len(m.names); i++ {
		m.nameIndex[m.names[i]] = i
	}
}

// Indexes は全トラックのフレーム番号を返す。
func (m *MorphFrames) Indexes() []int {
	if m == nil {
		return nil
	}
	indexes := make([]int, 0)
	for _, frames := range m.values {
		frames.ForEach(func(frame Frame, _ *MorphFrame) bool {
			indexes = append(indexes, int(frame))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

// Len は全フレーム数を返す。
func (m *MorphFrames) Len() int {
	if m == nil {
		return 0
	}
	count := 0
	for _, frames := range m.values {
		count += frames.Len()
	}
	return count
}

// MaxFrame は最大フレーム番号を返す。
func (m *MorphFrames) MaxFrame() Frame {
	if m == nil {
		return 0
	}
	maxFrame := Frame(0)
	for _, frames := range m.values {
		frame := frames.MaxFrame()
		if frame > maxFrame {
			maxFrame = frame
		}
	}
	return maxFrame
}

// MinFrame は最小フレーム番号を返す。
func (m *MorphFrames) MinFrame() Frame {
	if m == nil {
		return 0
	}
	minFrame := Frame(math.MaxFloat32)
	for _, frames := range m.values {
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

// Clean は無効なトラックを削除する。
func (m *MorphFrames) Clean() {
	if m == nil {
		return
	}
	keptNames := make([]string, 0, len(m.names))
	keptValues := make([]*MorphNameFrames, 0, len(m.values))
	keptIndex := make(map[string]int)
	for _, name := range m.names {
		frames := m.Get(name)
		if frames == nil || !frames.ContainsActive() {
			continue
		}
		keptIndex[name] = len(keptValues)
		keptNames = append(keptNames, name)
		keptValues = append(keptValues, frames)
	}
	m.names = keptNames
	m.values = keptValues
	m.nameIndex = keptIndex
}

// Copy はモーフ集合を複製する。
func (m *MorphFrames) Copy() (*MorphFrames, error) {
	copied, err := deepCopy(m)
	if err != nil {
		return nil, err
	}
	copied.rebuildNameIndex()
	return copied, nil
}

func (m *MorphFrames) rebuildNameIndex() {
	m.nameIndex = make(map[string]int, len(m.names))
	for i, name := range m.names {
		m.nameIndex[name] = i
	}
}

func nilMorphFrame() *MorphFrame {
	return nil
}
