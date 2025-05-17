package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/petar/GoLLRB/llrb"
)

type IBaseFrame interface {
	Index() float32
	SetIndex(index float32)
	IsRead() bool
	Llrb() *mmath.LlrbItem[float32]
	Less(than llrb.Item) bool
	lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame
	splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32)
	Copy() IBaseFrame
}

// ----------------------------------------------------------

type BaseFrame struct {
	*mmath.LlrbItem[float32]
	Read bool // VMDファイルから読み込んだキーフレであるか
}

func NewFrame(index float32) IBaseFrame {
	return &BaseFrame{
		LlrbItem: mmath.NewLlrbItem(index),
		Read:     false,
	}
}

func (baseFrame *BaseFrame) Llrb() *mmath.LlrbItem[float32] {
	return baseFrame.LlrbItem
}

func (baseFrame *BaseFrame) Index() float32 {
	return baseFrame.LlrbItem.Value()
}

func (baseFrame *BaseFrame) SetIndex(index float32) {
	baseFrame.LlrbItem = mmath.NewLlrbItem(index)
}

func (baseFrame *BaseFrame) IsRead() bool {
	return baseFrame.Read
}

func (baseFrame *BaseFrame) Less(than llrb.Item) bool {
	other, ok := than.(IBaseFrame)
	if ok {
		return baseFrame.LlrbItem.Less(other.Llrb())
	}

	other2, ok := than.(*mmath.LlrbItem[float32])
	if ok {
		return baseFrame.LlrbItem.Less(other2)
	}

	return false
}

func (baseFrame *BaseFrame) Copy() IBaseFrame {
	return &BaseFrame{
		LlrbItem: mmath.NewLlrbItem(baseFrame.Index()),
		Read:     baseFrame.Read,
	}
}

func (baseFrame *BaseFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	return baseFrame.Copy()
}

func (baseFrame *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type BaseFrames[T IBaseFrame] struct {
	values   *LlrbFrameIndexes     // 全キーフレリスト
	newFunc  func(index float32) T // キーフレ生成関数
	nullFunc func() T              // 空キーフレ生成関数
}

func NewBaseFrames[T IBaseFrame](newFunc func(index float32) T, nullFunc func() T) *BaseFrames[T] {
	return &BaseFrames[T]{
		values:   NewLlrbFrameIndexes(),
		newFunc:  newFunc,
		nullFunc: nullFunc,
	}
}

func (baseFrames *BaseFrames[T]) Get(frame float32) T {
	if baseFrames.values.Has(frame) {
		// キーフレが登録済みの場合、値を返す
		return baseFrames.values.Get(NewFrame(frame)).(T)
	}

	if baseFrames.values.Len() == 0 {
		// 指定INDEXで新フレームを作成
		return baseFrames.newFunc(frame)
	}

	prevFrame := baseFrames.PrevFrame(frame)
	nextFrame := baseFrames.NextFrame(frame)
	if nextFrame == frame {
		// 次のキーフレが無い場合、最大キーフレのコピーを返す
		if baseFrames.values.Len() == 0 {
			// 存在しない場合新規を返す
			return baseFrames.newFunc(frame)
		}

		copied := baseFrames.Get(baseFrames.MaxFrame()).Copy()
		copied.SetIndex(frame)
		return copied.(T)
	}

	prevF := baseFrames.Get(prevFrame)
	nextF := baseFrames.Get(nextFrame)

	// 該当キーフレが無い場合、補間結果を返す
	return nextF.lerpFrame(prevF, frame).(T)
}

func (baseFrames *BaseFrames[T]) PrevFrame(index float32) float32 {
	return baseFrames.values.Prev(index)
}

func (baseFrames *BaseFrames[T]) NextFrame(index float32) float32 {
	return baseFrames.values.Next(index)
}

func (baseFrames *BaseFrames[T]) ForEach(callback func(index float32, value T) bool) {
	baseFrames.values.ForEach(func(index float32) bool {
		return callback(index, baseFrames.Get(index))
	})
}

func (baseFrames *BaseFrames[T]) appendFrame(v T) {
	baseFrames.values.ReplaceOrInsert(v)
}

func (baseFrames *BaseFrames[T]) MaxFrame() float32 {
	if baseFrames.values.Len() == 0 {
		return 0
	}
	return baseFrames.values.Max()
}

func (baseFrames *BaseFrames[T]) MinFrame() float32 {
	if baseFrames.values.Len() == 0 {
		return 0
	}
	return baseFrames.values.Min()
}

func (baseFrames *BaseFrames[T]) Contains(frame float32) bool {
	return baseFrames.values.Has(frame)
}

func (baseFrames *BaseFrames[T]) Delete(frame float32) {
	baseFrames.values.Delete(mmath.NewLlrbItem(frame))
}

// Append 補間曲線は分割しない
func (baseFrames *BaseFrames[T]) Append(f T) {
	baseFrames.appendOrInsert(f, false)
}

// Insert Registered が true の場合、補間曲線を分割して登録する
func (baseFrames *BaseFrames[T]) Insert(f T) {
	baseFrames.appendOrInsert(f, true)
}

// Update 登録済みのキーフレームを更新する
func (baseFrames *BaseFrames[T]) Update(f T) {
	baseFrames.appendFrame(f)
}

func (baseFrames *BaseFrames[T]) appendOrInsert(f T, isSplitCurve bool) {
	if baseFrames.values.Len() == 0 {
		// フレームがない場合、何もしない
		baseFrames.appendFrame(f)
		return
	}

	if baseFrames.values.Has(f.Index()) {
		// 既に登録済みの場合、更新
		baseFrames.Update(f)
		return
	}

	if isSplitCurve {
		// 補間曲線を分割する
		prevF := baseFrames.Get(baseFrames.PrevFrame(f.Index()))
		nextF := baseFrames.Get(baseFrames.NextFrame(f.Index()))

		// 補間曲線を分割する
		if nextF.Index() > f.Index() && prevF.Index() < f.Index() {
			index := f.Index()
			f.splitCurve(prevF, nextF, index)
		}
	}

	baseFrames.appendFrame(f)
}

func (baseFrames *BaseFrames[T]) Length() int {
	return baseFrames.values.Length()
}
