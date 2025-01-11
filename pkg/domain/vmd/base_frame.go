package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/petar/GoLLRB/llrb"
)

type IBaseFrame interface {
	Index() float32
	SetIndex(index float32)
	IsRegistered() bool
	IsRead() bool
	Less(than llrb.Item) bool
	lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame
	splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32)
	Copy() IBaseFrame
}

type BaseFrame struct {
	index      *mmath.LlrbItem[float32]
	Registered bool // 登録対象のキーフレであるか
	Read       bool // VMDファイルから読み込んだキーフレであるか
}

func NewFrame(index float32) IBaseFrame {
	return &BaseFrame{
		index:      mmath.NewLlrbItem(index),
		Registered: false,
		Read:       false,
	}
}

func (baseFrame *BaseFrame) Index() float32 {
	return baseFrame.index.Value()
}

func (baseFrame *BaseFrame) SetIndex(index float32) {
	baseFrame.index = mmath.NewLlrbItem(index)
}

func (baseFrame *BaseFrame) IsRegistered() bool {
	return baseFrame.Registered
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
		index:      baseFrame.index,
		Registered: baseFrame.Registered,
		Read:       baseFrame.Read,
	}
}

func (baseFrame *BaseFrame) lerpFrame(prevFrame IBaseFrame, index float32) IBaseFrame {
	return baseFrame.Copy()
}

func (baseFrame *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float32) {
}

type BaseFrames[T IBaseFrame] struct {
	values            map[float32]T               // キーフレリスト
	Indexes           *mmath.LlrbIndexes[float32] // 全キーフレリスト
	RegisteredIndexes *mmath.LlrbIndexes[float32] // 登録対象キーフレリスト
	lock              sync.RWMutex                // マップアクセス制御用
}

func NewBaseFrames[T IBaseFrame]() *BaseFrames[T] {
	return &BaseFrames[T]{
		values:            make(map[float32]T),
		Indexes:           &mmath.LlrbIndexes[float32]{LLRB: llrb.New()},
		RegisteredIndexes: &mmath.LlrbIndexes[float32]{LLRB: llrb.New()},
		lock:              sync.RWMutex{},
	}
}

func (baseFrames *BaseFrames[T]) NewFrame(index float32) T {
	return NewFrame(index).(T)
}

func (baseFrames *BaseFrames[T]) Get(index float32) T {
	baseFrames.lock.RLock()
	defer baseFrames.lock.RUnlock()

	if _, ok := baseFrames.values[index]; ok {
		return baseFrames.values[index]
	}

	if len(baseFrames.values) == 0 {
		// 指定INDEXで新フレームを作成
		return baseFrames.NewFrame(index)
	}

	prevFrame := baseFrames.PrevFrame(index)
	nextFrame := baseFrames.NextFrame(index)
	if nextFrame == prevFrame {
		// 次のキーフレが無い場合、最大キーフレのコピーを返す
		if baseFrames.Indexes.Len() == 0 {
			// 存在しない場合新フレームを作成
			return baseFrames.NewFrame(index)
		}
		copied := baseFrames.values[baseFrames.Indexes.Max()].Copy()
		copied.SetIndex(index)
		return copied.(T)
	}

	prevF := baseFrames.Get(prevFrame)
	nextF := baseFrames.Get(nextFrame)

	// 該当キーフレが無い場合、補間結果を返す
	return nextF.lerpFrame(prevF, index).(T)
}

func (baseFrames *BaseFrames[T]) PrevFrame(index float32) float32 {
	return baseFrames.Indexes.Prev(index)
}

func (baseFrames *BaseFrames[T]) NextFrame(index float32) float32 {
	return baseFrames.Indexes.Next(index)
}

func (baseFrames *BaseFrames[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		for v := range baseFrames.Indexes.Iterator() {
			ch <- baseFrames.Get(v)
		}
		close(ch)
	}()
	return ch
}

func (baseFrames *BaseFrames[T]) appendFrame(v T) {
	baseFrames.lock.Lock()
	defer baseFrames.lock.Unlock()

	if v.IsRegistered() {
		baseFrames.RegisteredIndexes.ReplaceOrInsert(mmath.NewLlrbItem(v.Index()))
	}

	baseFrames.values[v.Index()] = v
	baseFrames.Indexes.ReplaceOrInsert(mmath.NewLlrbItem(v.Index()))
}

func (baseFrames *BaseFrames[T]) MaxFrame() float32 {
	if baseFrames.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return baseFrames.RegisteredIndexes.Max()
}

func (baseFrames *BaseFrames[T]) MinFrame() float32 {
	if baseFrames.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return baseFrames.RegisteredIndexes.Min()
}

func (baseFrames *BaseFrames[T]) ContainsRegistered(index float32) bool {
	return baseFrames.RegisteredIndexes.Has(index)
}

func (baseFrames *BaseFrames[T]) Contains(index float32) bool {
	_, ok := baseFrames.values[index]
	return ok
}

func (baseFrames *BaseFrames[T]) Delete(index float32) {
	baseFrames.lock.Lock()
	defer baseFrames.lock.Unlock()

	if _, ok := baseFrames.values[index]; ok {
		delete(baseFrames.values, index)
		baseFrames.Indexes.Delete(mmath.NewLlrbItem(index))
	}

	if baseFrames.RegisteredIndexes.Has(index) {
		baseFrames.RegisteredIndexes.Delete(mmath.NewLlrbItem(index))
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
	baseFrames.values[f.Index()] = f
}

func (baseFrames *BaseFrames[T]) appendOrInsert(f T, isSplitCurve bool) {
	if f.IsRegistered() {
		if baseFrames.RegisteredIndexes.Len() == 0 {
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
	}

	baseFrames.appendFrame(f)
}

func (baseFrames *BaseFrames[T]) Length() int {
	return baseFrames.RegisteredIndexes.Length()
}
