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

func (b *BaseFrame) Less(than llrb.Item) bool {
	return b.Index() < int(than.(core.Int))
}

func (b *BaseFrame) Copy() IBaseFrame {
	return &BaseFrame{
		index:      b.index,
		Registered: b.Registered,
		Read:       b.Read,
	}
}

func (b *BaseFrame) New(index int) IBaseFrame {
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

func (b *BaseFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	return b.Copy()
}

func (b *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}

func (bf *BaseFrame) getIntIndex() core.Int {
	return bf.index
}

func (bf *BaseFrame) Index() int {
	return int(bf.index)
}

func (bf *BaseFrame) SetIndex(index int) {
	bf.index = core.NewInt(index)
}

func (bf *BaseFrame) IsRegistered() bool {
	return bf.Registered
}

func (bf *BaseFrame) IsRead() bool {
	return bf.Read
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

func (fs *BaseFrames[T]) NewFrame(index int) T {
	return NewFrame(index).(T)
}

func (fs *BaseFrames[T]) Get(index int) T {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	if _, ok := fs.data[index]; ok {
		return fs.data[index]
	}

	if len(fs.data) == 0 {
		return fs.newFunc(index)
	}

	prevFrame := fs.prevFrame(index)
	nextFrame := fs.nextFrame(index)
	if nextFrame == prevFrame {
		// 次のキーフレが無い場合、最大キーフレのコピーを返す
		if fs.Indexes.Len() == 0 {
			// 存在しない場合nullを返す
			return fs.nullFunc()
		}
		copied := fs.data[fs.Indexes.Max()].Copy()
		copied.SetIndex(index)
		return copied.(T)
	}

	prevF := fs.Get(prevFrame)
	nextF := fs.Get(nextFrame)

	// 該当キーフレが無い場合、補間結果を返す
	return nextF.lerpFrame(prevF, index).(T)
}

func (fs *BaseFrames[T]) prevFrame(index int) int {
	prevFrame := fs.GetMinFrame()

	fs.RegisteredIndexes.DescendLessOrEqual(core.Int(index), func(i llrb.Item) bool {
		prevFrame = int(i.(core.Int))
		return false
	})

	return prevFrame
}

func (fs *BaseFrames[T]) nextFrame(index int) int {
	nextFrame := fs.GetMaxFrame()

	fs.RegisteredIndexes.AscendGreaterOrEqual(core.Int(index), func(i llrb.Item) bool {
		nextFrame = int(i.(core.Int))
		return false
	})

	return nextFrame
}

func (fs *BaseFrames[T]) List() []T {
	list := make([]T, 0, fs.RegisteredIndexes.Len())

	fs.RegisteredIndexes.AscendRange(core.Int(0), core.Int(fs.RegisteredIndexes.Max()), func(i llrb.Item) bool {
		if _, ok := fs.data[int(i.(core.Int))]; !ok {
			list = append(list, fs.data[int(i.(core.Int))])
		}
		return true
	})

	return list
}

func (fs *BaseFrames[T]) appendFrame(v T) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if v.IsRegistered() {
		fs.RegisteredIndexes.ReplaceOrInsert(core.Int(v.Index()))
	}

	fs.data[v.Index()] = v
	fs.Indexes.ReplaceOrInsert(core.Int(v.Index()))
}

func (fs *BaseFrames[T]) GetMaxFrame() int {
	if fs.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return fs.RegisteredIndexes.Max()
}

func (fs *BaseFrames[T]) GetMinFrame() int {
	if fs.RegisteredIndexes.Len() == 0 {
		return 0
	}
	return fs.RegisteredIndexes.Min()
}

func (fs *BaseFrames[T]) ContainsRegistered(index int) bool {
	return fs.RegisteredIndexes.Has(index)
}

func (fs *BaseFrames[T]) Contains(index int) bool {
	if _, ok := fs.data[index]; ok {
		return true
	}
	return false
}

func (fs *BaseFrames[T]) Delete(index int) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if _, ok := fs.data[index]; ok {
		delete(fs.data, index)
		fs.Indexes.Delete(core.Int(index))
	}

	if fs.RegisteredIndexes.Has(index) {
		fs.RegisteredIndexes.Delete(core.Int(index))
	}
}

// Append 補間曲線は分割しない
func (fs *BaseFrames[T]) Append(f T) {
	fs.appendOrInsert(f, false)
}

// Insert Registered が true の場合、補間曲線を分割して登録する
func (fs *BaseFrames[T]) Insert(f T) {
	fs.appendOrInsert(f, true)
}

func (fs *BaseFrames[T]) appendOrInsert(f T, isSplitCurve bool) {
	if f.IsRegistered() {
		if fs.RegisteredIndexes.Len() == 0 {
			// フレームがない場合、何もしない
			fs.appendFrame(f)
			return
		}

		if isSplitCurve {
			// 補間曲線を分割する
			prevF := fs.Get(fs.prevFrame(f.Index()))
			nextF := fs.Get(fs.nextFrame(f.Index()))

			// 補間曲線を分割する
			if nextF.Index() > f.Index() && prevF.Index() < f.Index() {
				index := f.Index()
				f.splitCurve(prevF, nextF, index)
			}
		}
	}

	fs.appendFrame(f)
}

func (fs *BaseFrames[T]) Len() int {
	return fs.RegisteredIndexes.Len()
}
