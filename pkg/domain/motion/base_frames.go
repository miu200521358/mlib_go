// 指示: miu200521358
package motion

// frameOps は補間/曲線分割を行うフレームの内部契約。
type frameOps[T any] interface {
	IBaseFrame
	lerpFrame(prev T, index Frame) T
	splitCurve(prev, next T, index Frame)
}

// BaseFrames はフレーム共通の格納を表す。
type BaseFrames[T frameOps[T]] struct {
	frames   map[Frame]T
	store    IFrameIndexStore
	newFunc  func(frame Frame) T
	nullFunc func() T
}

// NewBaseFrames はBaseFramesを生成する。
func NewBaseFrames[T frameOps[T]](newFunc func(frame Frame) T, nullFunc func() T) *BaseFrames[T] {
	return &BaseFrames[T]{
		frames:   make(map[Frame]T),
		store:    NewSortedFrameIndexStore(),
		newFunc:  newFunc,
		nullFunc: nullFunc,
	}
}

// Get は指定フレームの値を返す。
func (b *BaseFrames[T]) Get(frame Frame) T {
	if b == nil {
		return b.nullValue()
	}
	if v, ok := b.frames[frame]; ok {
		return v
	}
	if len(b.frames) == 0 {
		return b.newValue(frame)
	}

	prevFrame, _ := b.PrevFrame(frame)
	nextFrame, _ := b.NextFrame(frame)
	if nextFrame == frame {
		maxFrame := b.MaxFrame()
		maxValue, ok := b.frames[maxFrame]
		if !ok {
			return b.nullValue()
		}
		copied, err := maxValue.Copy()
		if err != nil {
			return b.nullValue()
		}
		casted, ok := copied.(T)
		if !ok {
			return b.nullValue()
		}
		casted.SetIndex(frame)
		return casted
	}

	prevValue, ok := b.frames[prevFrame]
	if !ok {
		return b.nullValue()
	}
	nextValue, ok := b.frames[nextFrame]
	if !ok {
		return b.nullValue()
	}
	return nextValue.lerpFrame(prevValue, frame)
}

// Append は末尾に追加する。
func (b *BaseFrames[T]) Append(frame T) {
	b.appendOrInsert(frame, false)
}

// Insert は補間曲線を分割して追加する。
func (b *BaseFrames[T]) Insert(frame T) {
	b.appendOrInsert(frame, true)
}

// Update は既存フレームを更新する。
func (b *BaseFrames[T]) Update(frame T) {
	if b == nil {
		return
	}
	if isNilValue(frame) {
		return
	}
	idx := frame.Index()
	b.frames[idx] = frame
	b.store.Upsert(idx)
}

// PrevFrame は直前のフレーム番号を返す。
func (b *BaseFrames[T]) PrevFrame(frame Frame) (Frame, bool) {
	if b == nil {
		return 0, false
	}
	return b.store.Prev(frame)
}

// NextFrame は直後のフレーム番号を返す。
func (b *BaseFrames[T]) NextFrame(frame Frame) (Frame, bool) {
	if b == nil {
		return frame, false
	}
	return b.store.Next(frame)
}

// Has はフレームの有無を判定する。
func (b *BaseFrames[T]) Has(frame Frame) bool {
	if b == nil {
		return false
	}
	_, ok := b.frames[frame]
	return ok
}

// Delete はフレームを削除する。
func (b *BaseFrames[T]) Delete(frame Frame) {
	if b == nil {
		return
	}
	delete(b.frames, frame)
	b.store.Delete(frame)
}

// ForEach は昇順で走査する。
func (b *BaseFrames[T]) ForEach(fn func(frame Frame, value T) bool) {
	if b == nil || fn == nil {
		return
	}
	b.store.ForEach(func(frame Frame) bool {
		v, ok := b.frames[frame]
		if !ok {
			return true
		}
		return fn(frame, v)
	})
}

// MaxFrame は最大フレーム番号を返す。
func (b *BaseFrames[T]) MaxFrame() Frame {
	if b == nil {
		return 0
	}
	max, ok := b.store.Max()
	if !ok {
		return 0
	}
	return max
}

// MinFrame は最小フレーム番号を返す。
func (b *BaseFrames[T]) MinFrame() Frame {
	if b == nil {
		return 0
	}
	min, ok := b.store.Min()
	if !ok {
		return 0
	}
	return min
}

// Len はフレーム数を返す。
func (b *BaseFrames[T]) Len() int {
	if b == nil {
		return 0
	}
	return len(b.frames)
}

// Finalize は索引の確定を行う。
func (b *BaseFrames[T]) Finalize() {
	if b == nil {
		return
	}
	b.store.Finalize()
}

func (b *BaseFrames[T]) appendOrInsert(frame T, split bool) {
	if b == nil || isNilValue(frame) {
		return
	}
	if split && len(b.frames) != 0 {
		prevFrame, _ := b.PrevFrame(frame.Index())
		nextFrame, _ := b.NextFrame(frame.Index())
		if prevFrame < frame.Index() && frame.Index() < nextFrame {
			prevValue, okPrev := b.frames[prevFrame]
			nextValue, okNext := b.frames[nextFrame]
			if okPrev && okNext {
				frame.splitCurve(prevValue, nextValue, frame.Index())
			}
		}
	}
	idx := frame.Index()
	b.frames[idx] = frame
	b.store.Upsert(idx)
}

func (b *BaseFrames[T]) newValue(frame Frame) T {
	if b == nil || b.newFunc == nil {
		var zero T
		return zero
	}
	return b.newFunc(frame)
}

func (b *BaseFrames[T]) nullValue() T {
	if b == nil || b.nullFunc == nil {
		var zero T
		return zero
	}
	return b.nullFunc()
}
