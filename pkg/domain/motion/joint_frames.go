// 指示: miu200521358
package motion

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// JointFrame はジョイントフレームを表す。
type JointFrame struct {
	*BaseFrame
	TranslationLimitMin       *mmath.Vec3
	TranslationLimitMax       *mmath.Vec3
	RotationLimitMin          *mmath.Vec3
	RotationLimitMax          *mmath.Vec3
	SpringConstantTranslation *mmath.Vec3
	SpringConstantRotation    *mmath.Vec3
}

// NewJointFrame はJointFrameを生成する。
func NewJointFrame(index Frame) *JointFrame {
	return &JointFrame{BaseFrame: NewBaseFrame(index)}
}

// Copy はフレームを複製する。
func (f *JointFrame) Copy() (IBaseFrame, error) {
	if f == nil {
		return (*JointFrame)(nil), nil
	}
	copied := &JointFrame{
		BaseFrame:                 &BaseFrame{index: f.Index(), Read: f.Read},
		TranslationLimitMin:       copyVec3(f.TranslationLimitMin),
		TranslationLimitMax:       copyVec3(f.TranslationLimitMax),
		RotationLimitMin:          copyVec3(f.RotationLimitMin),
		RotationLimitMax:          copyVec3(f.RotationLimitMax),
		SpringConstantTranslation: copyVec3(f.SpringConstantTranslation),
		SpringConstantRotation:    copyVec3(f.SpringConstantRotation),
	}
	return copied, nil
}

// lerpFrame は補間結果を返す。
func (next *JointFrame) lerpFrame(prev *JointFrame, index Frame) *JointFrame {
	if prev == nil && next == nil {
		return NewJointFrame(index)
	}
	if prev == nil {
		copied, _ := next.Copy()
		out := copied.(*JointFrame)
		out.SetIndex(index)
		return out
	}
	if next == nil {
		copied, _ := prev.Copy()
		out := copied.(*JointFrame)
		out.SetIndex(index)
		return out
	}
	out := NewJointFrame(index)
	t := linearT(prev.Index(), index, next.Index())
	minT := vec3OrZero(prev.TranslationLimitMin).Lerp(vec3OrZero(next.TranslationLimitMin), t)
	maxT := vec3OrZero(prev.TranslationLimitMax).Lerp(vec3OrZero(next.TranslationLimitMax), t)
	minR := vec3OrZero(prev.RotationLimitMin).Lerp(vec3OrZero(next.RotationLimitMin), t)
	maxR := vec3OrZero(prev.RotationLimitMax).Lerp(vec3OrZero(next.RotationLimitMax), t)
	springT := vec3OrZero(prev.SpringConstantTranslation).Lerp(vec3OrZero(next.SpringConstantTranslation), t)
	springR := vec3OrZero(prev.SpringConstantRotation).Lerp(vec3OrZero(next.SpringConstantRotation), t)
	out.TranslationLimitMin = &minT
	out.TranslationLimitMax = &maxT
	out.RotationLimitMin = &minR
	out.RotationLimitMax = &maxR
	out.SpringConstantTranslation = &springT
	out.SpringConstantRotation = &springR
	return out
}

// splitCurve は何もしない。
func (f *JointFrame) splitCurve(prev *JointFrame, next *JointFrame, index Frame) {
}

// JointNameFrames はジョイント名ごとのフレーム集合を表す。
type JointNameFrames struct {
	*BaseFrames[*JointFrame]
	Name string
}

// NewJointNameFrames はJointNameFramesを生成する。
func NewJointNameFrames(name string) *JointNameFrames {
	return &JointNameFrames{
		BaseFrames: NewBaseFrames(newJointFrame, nilJointFrame),
		Name:       name,
	}
}

// Copy はフレーム集合を複製する。
func (j *JointNameFrames) Copy() (*JointNameFrames, error) {
	return deepCopy(j)
}

// JointFrames はジョイント名ごとの集合を表す。
type JointFrames struct {
	names     []string
	nameIndex map[string]int
	values    []*JointNameFrames
}

// NewJointFrames はJointFramesを生成する。
func NewJointFrames() *JointFrames {
	return &JointFrames{
		names:     make([]string, 0),
		nameIndex: make(map[string]int),
		values:    make([]*JointNameFrames, 0),
	}
}

// Names は登録順の名前一覧を返す。
func (j *JointFrames) Names() []string {
	if j == nil {
		return nil
	}
	return append([]string(nil), j.names...)
}

// Get は名前に対応するフレーム集合を返す。
func (j *JointFrames) Get(name string) *JointNameFrames {
	if j == nil {
		return nil
	}
	if idx, ok := j.nameIndex[name]; ok {
		return j.values[idx]
	}
	return nil
}

// Update はフレーム集合を更新する。
func (j *JointFrames) Update(frames *JointNameFrames) {
	if j == nil || frames == nil {
		return
	}
	if idx, ok := j.nameIndex[frames.Name]; ok {
		j.values[idx] = frames
		return
	}
	j.nameIndex[frames.Name] = len(j.values)
	j.names = append(j.names, frames.Name)
	j.values = append(j.values, frames)
}

// Delete は名前を削除する。
func (j *JointFrames) Delete(name string) {
	if j == nil {
		return
	}
	idx, ok := j.nameIndex[name]
	if !ok {
		return
	}
	j.names = append(j.names[:idx], j.names[idx+1:]...)
	j.values = append(j.values[:idx], j.values[idx+1:]...)
	delete(j.nameIndex, name)
	for i := idx; i < len(j.names); i++ {
		j.nameIndex[j.names[i]] = i
	}
}

// Indexes は全トラックのフレーム番号を返す。
func (j *JointFrames) Indexes() []int {
	if j == nil {
		return nil
	}
	indexes := make([]int, 0)
	for _, frames := range j.values {
		frames.ForEach(func(frame Frame, _ *JointFrame) bool {
			indexes = append(indexes, int(frame))
			return true
		})
	}
	indexes = mmath.Unique(indexes)
	mmath.Sort(indexes)
	return indexes
}

// Len は全フレーム数を返す。
func (j *JointFrames) Len() int {
	if j == nil {
		return 0
	}
	count := 0
	for _, frames := range j.values {
		count += frames.Len()
	}
	return count
}

// MaxFrame は最大フレーム番号を返す。
func (j *JointFrames) MaxFrame() Frame {
	if j == nil {
		return 0
	}
	maxFrame := Frame(0)
	for _, frames := range j.values {
		frame := frames.MaxFrame()
		if frame > maxFrame {
			maxFrame = frame
		}
	}
	return maxFrame
}

// MinFrame は最小フレーム番号を返す。
func (j *JointFrames) MinFrame() Frame {
	if j == nil {
		return 0
	}
	minFrame := Frame(math.MaxFloat32)
	for _, frames := range j.values {
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

// Copy はジョイント集合を複製する。
func (j *JointFrames) Copy() (*JointFrames, error) {
	copied, err := deepCopy(j)
	if err != nil {
		return nil, err
	}
	copied.rebuildNameIndex()
	return copied, nil
}

func (j *JointFrames) rebuildNameIndex() {
	j.nameIndex = make(map[string]int, len(j.names))
	for i, name := range j.names {
		j.nameIndex[name] = i
	}
}

func newJointFrame(frame Frame) *JointFrame {
	return nil
}

func nilJointFrame() *JointFrame {
	return nil
}
