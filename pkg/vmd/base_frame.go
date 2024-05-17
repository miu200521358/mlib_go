package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/petar/GoLLRB/llrb"
)

type IBaseFrame interface {
	GetIndex() int
	getIntIndex() mcore.Int
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
	Index      mcore.Int
	Registered bool // 登録対象のキーフレであるか
	Read       bool // VMDファイルから読み込んだキーフレであるか
}

func (b *BaseFrame) Less(than llrb.Item) bool {
	return b.GetIndex() < int(than.(mcore.Int))
}

func (b *BaseFrame) Copy() IBaseFrame {
	return &BaseFrame{
		Index:      b.Index,
		Registered: b.Registered,
		Read:       b.Read,
	}
}

func (b *BaseFrame) New(index int) IBaseFrame {
	return &BaseFrame{
		Index: mcore.Int(index),
	}
}

func NewFrame(index int) IBaseFrame {
	return &BaseFrame{
		Index:      mcore.NewInt(index),
		Registered: false,
		Read:       false,
	}
}

func (b *BaseFrame) lerpFrame(prevFrame IBaseFrame, index int) IBaseFrame {
	return b.Copy()
}

func (b *BaseFrame) splitCurve(prevFrame IBaseFrame, nextFrame IBaseFrame, index int) {
}

func (bf *BaseFrame) getIntIndex() mcore.Int {
	return bf.Index
}

func (bf *BaseFrame) GetIndex() int {
	return int(bf.Index)
}

func (bf *BaseFrame) SetIndex(index int) {
	bf.Index = mcore.NewInt(index)
}

func (bf *BaseFrame) IsRegistered() bool {
	return bf.Registered
}

func (bf *BaseFrame) IsRead() bool {
	return bf.Read
}

type BaseFrames[T IBaseFrame] struct {
	data              map[int]T         // キーフレリスト
	RegisteredIndexes *mcore.IntIndexes // 登録対象キーフレリスト
	newFunc           func(index int) T // キーフレ生成関数
	lock              sync.RWMutex      // マップアクセス制御用
}

func NewBaseFrames[T IBaseFrame](newFunc func(index int) T) *BaseFrames[T] {
	return &BaseFrames[T]{
		data:              make(map[int]T),
		RegisteredIndexes: mcore.NewIntIndexes(),
		newFunc:           newFunc,
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
		if fs.RegisteredIndexes.Len() == 0 {
			// 登録キーが無い場合、現在登録されているすべてのキーフレの中から最大のものをコピーして返す
			// 上でデータが無いことのチェックは済んでいるので何かしらのキーはあるはず
			indexes := mcore.NewIntIndexes()
			for i := range fs.data {
				indexes.ReplaceOrInsert(mcore.Int(i))
			}
			copied := fs.data[indexes.Max()].Copy().(T)
			copied.SetIndex(index)
			return copied
		}
		copied := fs.data[fs.RegisteredIndexes.Max()].Copy().(T)
		copied.SetIndex(index)
		return copied
	}

	prevF := fs.Get(prevFrame)
	nextF := fs.Get(nextFrame)

	// 該当キーフレが無い場合、補間結果を返す
	return nextF.lerpFrame(prevF, index).(T)
}

func (fs *BaseFrames[T]) prevFrame(index int) int {
	prevFrame := fs.GetMinFrame()

	fs.RegisteredIndexes.DescendLessOrEqual(mcore.Int(index), func(i llrb.Item) bool {
		prevFrame = int(i.(mcore.Int))
		return false
	})

	return prevFrame
}

func (fs *BaseFrames[T]) nextFrame(index int) int {
	nextFrame := fs.GetMaxFrame()

	fs.RegisteredIndexes.AscendGreaterOrEqual(mcore.Int(index), func(i llrb.Item) bool {
		nextFrame = int(i.(mcore.Int))
		return false
	})

	return nextFrame
}

func (fs *BaseFrames[T]) List() []T {
	list := make([]T, 0, fs.RegisteredIndexes.Len())

	fs.RegisteredIndexes.AscendRange(mcore.Int(0), mcore.Int(fs.RegisteredIndexes.Max()), func(i llrb.Item) bool {
		if _, ok := fs.data[int(i.(mcore.Int))]; !ok {
			list = append(list, fs.data[int(i.(mcore.Int))])
		}
		return true
	})

	return list
}

func (fs *BaseFrames[T]) appendFrame(v T) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if v.IsRegistered() {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Int(v.GetIndex()))
	}
	fs.data[v.GetIndex()] = v
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
			prevF := fs.Get(fs.prevFrame(f.GetIndex()))
			nextF := fs.Get(fs.nextFrame(f.GetIndex()))

			// 補間曲線を分割する
			if nextF.GetIndex() > f.GetIndex() && prevF.GetIndex() < f.GetIndex() {
				index := f.GetIndex()
				f.splitCurve(prevF, nextF, index)
			}
		}
	}

	fs.appendFrame(f)
}

func (fs *BaseFrames[T]) Len() int {
	return fs.RegisteredIndexes.Len()
}
