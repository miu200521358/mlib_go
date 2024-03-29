package vmd

import (
	"sync"

	"github.com/petar/GoLLRB/llrb"

	"github.com/miu200521358/mlib_go/pkg/mcore"

)

type IkFrames struct {
	*mcore.IndexFloatModels[*IkFrame]
	RegisteredIndexes *mcore.FloatIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex        // マップアクセス制御用
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		IndexFloatModels:  mcore.NewIndexFloatModels[*IkFrame](),
		RegisteredIndexes: mcore.NewFloatIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *IkFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := mcore.Float32(0)
	nextIndex := mcore.Float32(index)

	if fs.RegisteredIndexes.Max() < index {
		return float32(fs.RegisteredIndexes.Max()), float32(fs.RegisteredIndexes.Max())
	}

	fs.RegisteredIndexes.DescendLessOrEqual(mcore.Float32(index), func(i llrb.Item) bool {
		prevIndex = i.(mcore.Float32)
		return false
	})

	fs.RegisteredIndexes.AscendGreaterOrEqual(mcore.Float32(index), func(i llrb.Item) bool {
		nextIndex = i.(mcore.Float32)
		return false
	})

	return float32(prevIndex), float32(nextIndex)
}

// キーフレ計算結果を返す
func (fs *IkFrames) GetItem(index float32) *IkFrame {
	if fs == nil {
		return NewIkFrame(index)
	}

	fs.lock.RLock()
	defer fs.lock.RUnlock()

	if fs.Has(index) {
		return fs.Get(index)
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex && fs.Has(nextIndex) {
		nextIf := fs.Get(nextIndex)
		copied := nextIf.Copy()
		return copied.(*IkFrame)
	}

	var prevIf *IkFrame
	if fs.Has(prevIndex) {
		prevIf = fs.Get(prevIndex)
	} else {
		prevIf = NewIkFrame(index)
	}

	nif := prevIf.Copy()
	return nif.(*IkFrame)
}

func (fs *IkFrames) Append(value *IkFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if value.Registered {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Float32(value.Index))
	}

	fs.ReplaceOrInsert(value)
}
