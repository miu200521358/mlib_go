package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/petar/GoLLRB/llrb"
)

type IBaseFrame interface {
	Index() int
	getIntIndex() core.Int
	SetIndex(index int)
	IsRegistered() bool
	IsRead() bool
	Less(than llrb.Item) bool
	lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame
	splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int)
	Copy() IBaseFrame
	New(index int) IBaseFrame
}

type BaseFrame struct {
	index      core.Int
	Registered bool // 登録対象のキーフレであるか
	Read       bool // VMDファイルから読み込んだキーフレであるか
}

func (baseFrame *BaseFrame) Less(than llrb.Item) bool {
	return baseFrame.Index() < int(than.(core.Int))
}

func (baseFrame *BaseFrame) Copy() IBaseFrame {
	return &BaseFrame{
		index:      baseFrame.index,
		Registered: baseFrame.Registered,
		Read:       baseFrame.Read,
	}
}

func (baseFrame *BaseFrame) New(index int) IBaseFrame {
	return &BaseFrame{
		index: core.Int(index),
	}
}

func NewFrame(index int) IBaseFrame {
	return &BaseFrame{
		index:      core.NewInt(index),
		Registered: false,
		Read:       false,
	}
}

func (baseFrame *BaseFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	return baseFrame.Copy()
}

func (baseFrame *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}

func (baseFrame *BaseFrame) getIntIndex() core.Int {
	return baseFrame.index
}

func (baseFrame *BaseFrame) Index() int {
	return int(baseFrame.index)
}

func (baseFrame *BaseFrame) SetIndex(index int) {
	baseFrame.index = core.NewInt(index)
}

func (baseFrame *BaseFrame) IsRegistered() bool {
	return baseFrame.Registered
}

func (baseFrame *BaseFrame) IsRead() bool {
	return baseFrame.Read
}

type BaseFrames[T IBaseFrame] struct {
	data              map[int]T         // キーフレリスト
	Indexes           *core.IntIndexes  // 全キーフレリスト
	RegisteredIndexes *core.IntIndexes  // 登録対象キーフレリスト
	newFunc           func(index int) T // キーフレ生成関数
	nullFunc          func() T          // 空キーフレ生成関数
	lock              sync.RWMutex      // マップアクセス制御用
}

func NewBaseFrames[T IBaseFrame](newFunc func(index int) T, nullFunc func() T) *BaseFrames[T] {
	return &BaseFrames[T]{
		data:              make(map[int]T),
		Indexes:           core.NewIntIndexes(),
		RegisteredIndexes: core.NewIntIndexes(),
		newFunc:           newFunc,
		nullFunc:          nullFunc,
		lock:              sync.RWMutex{},
	}
}

func (baseFrames *BaseFrames[T]) NewFrame(index int) T {
	return NewFrame(index).(T)
}

func (baseFrames *BaseFrames[T]) Get(index int) T {
	baseFrames.lock.RLock()
	defer baseFrames.lock.RUnlock()

	if _, ok := baseFrames.data[index]; ok {
		return baseFrames.data[index]
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

func (baseFrames *BaseFrames[T]) prevFrame(index int) int {
	prevFrame := baseFrames.MinFrame()

	baseFrames.RegisteredIndexes.DescendLessOrEqual(core.Int(index), func(i llrb.Item) bool {
		prevFrame = int(i.(core.Int))
		return false
	})

	return prevFrame
}

func (baseFrames *BaseFrames[T]) nextFrame(index int) int {
	nextFrame := baseFrames.MaxFrame()

	baseFrames.RegisteredIndexes.AscendGreaterOrEqual(core.Int(index), func(i llrb.Item) bool {
		nextFrame = int(i.(core.Int))
		return false
	})

	return nextFrame
}

func (baseFrames *BaseFrames[T]) List() []T {
	list := make([]T, 0, baseFrames.RegisteredIndexes.Len())

	baseFrames.RegisteredIndexes.AscendRange(core.Int(0), core.Int(baseFrames.RegisteredIndexes.Max()), func(i llrb.Item) bool {
		if _, ok := baseFrames.data[int(i.(core.Int))]; !ok {
			list = append(list, baseFrames.data[int(i.(core.Int))])
		}
		return true
	})

	return list
}

func (baseFrames *BaseFrames[T]) appendFrame(v T) {
	baseFrames.lock.Lock()
	defer baseFrames.lock.Unlock()

	if v.IsRegistered() {
		baseFrames.RegisteredIndexes.ReplaceOrInsert(core.Int(v.Index()))
	}

	baseFrames.data[v.Index()] = v
	baseFrames.Indexes.ReplaceOrInsert(core.Int(v.Index()))
}

func (baseFrames *BaseFrames[T]) MaxFrame() int {
	if baseFrames.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return baseFrames.RegisteredIndexes.Max()
}

func (baseFrames *BaseFrames[T]) MinFrame() int {
	if baseFrames.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return baseFrames.RegisteredIndexes.Min()
}

func (baseFrames *BaseFrames[T]) ContainsRegistered(index int) bool {
	return baseFrames.RegisteredIndexes.Has(index)
}

func (baseFrames *BaseFrames[T]) Contains(index int) bool {
	if _, ok := baseFrames.data[index]; ok {
		return true
	}
	return false
}

func (baseFrames *BaseFrames[T]) Delete(index int) {
	baseFrames.lock.Lock()
	defer baseFrames.lock.Unlock()

	if _, ok := baseFrames.data[index]; ok {
		delete(baseFrames.data, index)
		baseFrames.Indexes.Delete(core.Int(index))
	}

	if baseFrames.RegisteredIndexes.Has(index) {
		baseFrames.RegisteredIndexes.Delete(core.Int(index))
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
