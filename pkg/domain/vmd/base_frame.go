package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/petar/GoLLRB/llrb"
)

type IBaseFrame interface {
	Index() float64
	getIntIndex() core.Float
	SetIndex(index float64)
	IsRegistered() bool
	IsRead() bool
	Less(than llrb.Item) bool
	lerpFrame(prevFrame IBaseFrame, index float64) IBaseFrame
	splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float64)
	Copy() IBaseFrame
	New(index float64) IBaseFrame
}

type BaseFrame struct {
	index      core.Float
	Registered bool // 登録対象のキーフレであるか
	Read       bool // VMDファイルから読み込んだキーフレであるか
}

func (baseFrame *BaseFrame) Less(than llrb.Item) bool {
	return baseFrame.Index() < float64(than.(core.Float))
}

func (baseFrame *BaseFrame) Copy() IBaseFrame {
	return &BaseFrame{
		index:      baseFrame.index,
		Registered: baseFrame.Registered,
		Read:       baseFrame.Read,
	}
}

func (baseFrame *BaseFrame) New(index float64) IBaseFrame {
	return &BaseFrame{
		index: core.Float(index),
	}
}

func NewFrame(index float64) IBaseFrame {
	return &BaseFrame{
		index:      core.NewFloat(index),
		Registered: false,
		Read:       false,
	}
}

func (baseFrame *BaseFrame) lerpFrame(prevFrame IBaseFrame, index float64) IBaseFrame {
	return baseFrame.Copy()
}

func (baseFrame *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index float64) {
}

func (baseFrame *BaseFrame) getIntIndex() core.Float {
	return baseFrame.index
}

func (baseFrame *BaseFrame) Index() float64 {
	return float64(baseFrame.index)
}

func (baseFrame *BaseFrame) SetIndex(index float64) {
	baseFrame.index = core.NewFloat(index)
}

func (baseFrame *BaseFrame) IsRegistered() bool {
	return baseFrame.Registered
}

func (baseFrame *BaseFrame) IsRead() bool {
	return baseFrame.Read
}

type BaseFrames[T IBaseFrame] struct {
	data              map[float64]T         // キーフレリスト
	Indexes           *core.FloatIndexes    // 全キーフレリスト
	RegisteredIndexes *core.FloatIndexes    // 登録対象キーフレリスト
	newFunc           func(index float64) T // キーフレ生成関数
	nullFunc          func() T              // 空キーフレ生成関数
	lock              sync.RWMutex          // マップアクセス制御用
}

func NewBaseFrames[T IBaseFrame](newFunc func(index float64) T, nullFunc func() T) *BaseFrames[T] {
	return &BaseFrames[T]{
		data:              make(map[float64]T),
		Indexes:           core.NewFloatIndexes(),
		RegisteredIndexes: core.NewFloatIndexes(),
		newFunc:           newFunc,
		nullFunc:          nullFunc,
		lock:              sync.RWMutex{},
	}
}

func (baseFrames *BaseFrames[T]) NewFrame(index float64) T {
	return NewFrame(index).(T)
}

func (baseFrames *BaseFrames[T]) Get(index float64) T {
	baseFrames.lock.RLock()
	defer baseFrames.lock.RUnlock()

	if _, ok := baseFrames.data[float64(index)]; ok {
		return baseFrames.data[float64(index)]
	}

	if len(baseFrames.data) == 0 {
		return baseFrames.newFunc(index)
	}

	prevFrame := baseFrames.prevFrame(index)
	nextFrame := baseFrames.nextFrame(index)
	if nextFrame == prevFrame {
		// 次のキーフレが無い場合、最大キーフレのコピーを返す
		if baseFrames.Indexes.Len() == 0 {
			// 存在しない場合nullを返す
			return baseFrames.nullFunc()
		}
		copied := baseFrames.data[baseFrames.Indexes.Max()].Copy()
		copied.SetIndex(index)
		return copied.(T)
	}

	prevF := baseFrames.Get(prevFrame)
	nextF := baseFrames.Get(nextFrame)

	// 該当キーフレが無い場合、補間結果を返す
	return nextF.lerpFrame(prevF, index).(T)
}

func (baseFrames *BaseFrames[T]) prevFrame(index float64) float64 {
	prevFrame := baseFrames.MinFrame()

	baseFrames.RegisteredIndexes.DescendLessOrEqual(core.Float(index), func(i llrb.Item) bool {
		prevFrame = float64(i.(core.Float))
		return false
	})

	return prevFrame
}

func (baseFrames *BaseFrames[T]) nextFrame(index float64) float64 {
	nextFrame := baseFrames.MaxFrame()

	baseFrames.RegisteredIndexes.AscendGreaterOrEqual(core.Float(index), func(i llrb.Item) bool {
		nextFrame = float64(i.(core.Float))
		return false
	})

	return nextFrame
}

func (baseFrames *BaseFrames[T]) List() []T {
	list := make([]T, 0, baseFrames.RegisteredIndexes.Len())

	baseFrames.RegisteredIndexes.AscendRange(core.Float(0), core.Float(baseFrames.RegisteredIndexes.Max()), func(i llrb.Item) bool {
		if _, ok := baseFrames.data[float64(i.(core.Float))]; !ok {
			list = append(list, baseFrames.data[float64(i.(core.Float))])
		}
		return true
	})

	return list
}

func (baseFrames *BaseFrames[T]) appendFrame(v T) {
	baseFrames.lock.Lock()
	defer baseFrames.lock.Unlock()

	if v.IsRegistered() {
		baseFrames.RegisteredIndexes.ReplaceOrInsert(core.Float(v.Index()))
	}

	baseFrames.data[v.Index()] = v
	baseFrames.Indexes.ReplaceOrInsert(core.Float(v.Index()))
}

func (baseFrames *BaseFrames[T]) MaxFrame() float64 {
	if baseFrames.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return baseFrames.RegisteredIndexes.Max()
}

func (baseFrames *BaseFrames[T]) MinFrame() float64 {
	if baseFrames.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return baseFrames.RegisteredIndexes.Min()
}

func (baseFrames *BaseFrames[T]) ContainsRegistered(index float64) bool {
	return baseFrames.RegisteredIndexes.Has(index)
}

func (baseFrames *BaseFrames[T]) Contains(index float64) bool {
	if _, ok := baseFrames.data[index]; ok {
		return true
	}
	return false
}

func (baseFrames *BaseFrames[T]) Delete(index float64) {
	baseFrames.lock.Lock()
	defer baseFrames.lock.Unlock()

	if _, ok := baseFrames.data[index]; ok {
		delete(baseFrames.data, index)
		baseFrames.Indexes.Delete(core.Float(index))
	}

	if baseFrames.RegisteredIndexes.Has(index) {
		baseFrames.RegisteredIndexes.Delete(core.Float(index))
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

func (baseFrames *BaseFrames[T]) appendOrInsert(f T, isSplitCurve bool) {
	if f.IsRegistered() {
		if baseFrames.RegisteredIndexes.Len() == 0 {
			// フレームがない場合、何もしない
			baseFrames.appendFrame(f)
			return
		}

		if isSplitCurve {
			// 補間曲線を分割する
			prevF := baseFrames.Get(baseFrames.prevFrame(f.Index()))
			nextF := baseFrames.Get(baseFrames.nextFrame(f.Index()))

			// 補間曲線を分割する
			if nextF.Index() > f.Index() && prevF.Index() < f.Index() {
				index := f.Index()
				f.splitCurve(prevF, nextF, index)
			}
		}
	}

	baseFrames.appendFrame(f)
}

func (baseFrames *BaseFrames[T]) Len() int {
	return baseFrames.RegisteredIndexes.Len()
}
