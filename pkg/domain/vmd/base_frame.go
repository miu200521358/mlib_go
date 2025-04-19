package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/petar/GoLLRB/llrb"
)

type IBaseFrame interface {
	Index() float32
	SetIndex(index float32)
	IsRead() bool
	Less(than llrb.Item) bool
	lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame
	splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32)
	Copy() IBaseFrame
}

type BaseFrame struct {
	index *mmath.LlrbItem[float32]
	Read  bool // VMDファイルから読み込んだキーフレであるか
}

func NewFrame(index float32) IBaseFrame {
	return &BaseFrame{
		index: mmath.NewLlrbItem(index),
		Read:  false,
	}
}

func (baseFrame *BaseFrame) Index() float32 {
	return baseFrame.index.Value()
}

func (baseFrame *BaseFrame) SetIndex(index float32) {
	baseFrame.index = mmath.NewLlrbItem(index)
}

func (baseFrame *BaseFrame) IsRead() bool {
	return baseFrame.Read
}

func (baseFrame *BaseFrame) Less(than llrb.Item) bool {
	other, ok := than.(mmath.LlrbItem[float32])
	if !ok {
		return false
	}
	return baseFrame.index.Less(other)
}

func (baseFrame *BaseFrame) Copy() IBaseFrame {
	return &BaseFrame{
		index: baseFrame.index,
		Read:  baseFrame.Read,
	}
}

func (baseFrame *BaseFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	return baseFrame.Copy()
}

func (baseFrame *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type BaseFrames[T IBaseFrame] struct {
	values       []T                         // キーフレームの値
	valueIndexes []float32                   // キーフレームのフレーム番号
	Indexes      *mmath.LlrbIndexes[float32] // 全キーフレリスト
	newFunc      func(index float32) T       // キーフレ生成関数
	nullFunc     func() T                    // 空キーフレ生成関数
}

func NewBaseFrames[T IBaseFrame](newFunc func(index float32) T, nullFunc func() T) *BaseFrames[T] {
	return &BaseFrames[T]{
		values:       make([]T, 0),
		valueIndexes: make([]float32, 0),
		Indexes:      &mmath.LlrbIndexes[float32]{LLRB: llrb.New()},
		newFunc:      newFunc,
		nullFunc:     nullFunc,
	}
}

func (baseFrames *BaseFrames[T]) NewFrame(index float32) T {
	return NewFrame(index).(T)
}

func (baseFrames *BaseFrames[T]) Get(frame float32) T {
	index := slices.Index(baseFrames.valueIndexes, frame)
	if index >= 0 {
		return baseFrames.values[index]
	}

	if len(baseFrames.values) <= 1 {
		// 指定INDEXで新フレームを作成
		return baseFrames.newFunc(frame)
	}

	prevFrame := baseFrames.PrevFrame(frame)
	nextFrame := baseFrames.NextFrame(frame)
	if nextFrame == frame {
		// 次のキーフレが無い場合、最大キーフレのコピーを返す
		if baseFrames.Indexes.Len() == 0 {
			// 存在しない場合nilを返す
			return baseFrames.nullFunc()
		}

		index := slices.Index(baseFrames.valueIndexes, baseFrames.Indexes.Max())
		copied := baseFrames.values[index].Copy()
		copied.SetIndex(frame)
		return copied.(T)
	}

	prevF := baseFrames.Get(prevFrame)
	nextF := baseFrames.Get(nextFrame)

	// 該当キーフレが無い場合、補間結果を返す
	return nextF.lerpFrame(prevF, frame).(T)
}

func (baseFrames *BaseFrames[T]) PrevFrame(index float32) float32 {
	return baseFrames.Indexes.Prev(index)
}

func (baseFrames *BaseFrames[T]) NextFrame(index float32) float32 {
	return baseFrames.Indexes.Next(index)
}

func (baseFrames *BaseFrames[T]) ForEach(callback func(index float32, value T) bool) {
	baseFrames.Indexes.ForEach(func(index float32) bool {
		return callback(index, baseFrames.Get(index))
	})

	// for _, v := range  {
	// 	if !callback(v.Index(), v) {
	// 		return
	// 	}
	// }

	// baseFrames.Indexes.ForEach(func(index float32) bool {
	// 	return callback(index, baseFrames.Get(index))
	// })
}

func (baseFrames *BaseFrames[T]) appendFrame(v T) {
	baseFrames.Indexes.ReplaceOrInsert(mmath.NewLlrbItem(v.Index()))
	baseFrames.valueIndexes = append(baseFrames.valueIndexes, v.Index())
	baseFrames.values = append(baseFrames.values, v)
}

func (baseFrames *BaseFrames[T]) MaxFrame() float32 {
	if baseFrames.Indexes.Len() == 0 {
		return 0
	}
	return baseFrames.Indexes.Max()
}

func (baseFrames *BaseFrames[T]) MinFrame() float32 {
	if baseFrames.Indexes.Len() == 0 {
		return 0
	}
	return baseFrames.Indexes.Min()
}

func (baseFrames *BaseFrames[T]) ContainsRegistered(index float32) bool {
	return baseFrames.Indexes.Has(index)
}

func (baseFrames *BaseFrames[T]) Contains(frame float32) bool {
	index := slices.Index(baseFrames.valueIndexes, frame)
	return index >= 0
}

func (baseFrames *BaseFrames[T]) Delete(frame float32) {
	index := slices.Index(baseFrames.valueIndexes, frame)
	if index < 0 {
		return
	}

	baseFrames.valueIndexes = append(baseFrames.valueIndexes[:index], baseFrames.valueIndexes[index+1:]...)
	baseFrames.values = append(baseFrames.values[:index], baseFrames.values[index+1:]...)

	if baseFrames.Indexes.Has(frame) {
		baseFrames.Indexes.Delete(mmath.NewLlrbItem(frame))
	}
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
	index := slices.Index(baseFrames.valueIndexes, f.Index())
	baseFrames.values[index] = f
}

func (baseFrames *BaseFrames[T]) appendOrInsert(f T, isSplitCurve bool) {
	if baseFrames.Indexes.Len() == 0 {
		// フレームがない場合、何もしない
		baseFrames.appendFrame(f)
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
	return baseFrames.Indexes.Length()
}
