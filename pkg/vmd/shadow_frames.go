package vmd

import (
	"sync"

	"github.com/petar/GoLLRB/llrb"

	"github.com/miu200521358/mlib_go/pkg/mcore"

)

type ShadowFrames struct {
	*mcore.IndexFloatModels[*ShadowFrame]
	RegisteredIndexes *mcore.FloatIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex        // マップアクセス制御用
}

func NewShadowFrames() *ShadowFrames {
	return &ShadowFrames{
		IndexFloatModels:  mcore.NewIndexFloatModelCorrection[*ShadowFrame](),
		RegisteredIndexes: mcore.NewFloatIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *ShadowFrames) GetRangeIndexes(index float32) (float32, float32) {

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
func (fs *ShadowFrames) GetItem(index float32) *ShadowFrame {
	if val, ok := fs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex && fs.Indexes.Has(nextIndex) {
		nextSf := fs.Data[nextIndex]
		copied := nextSf.Copy()
		return copied.(*ShadowFrame)
	}

	var prevSf *ShadowFrame
	if fs.Indexes.Has(prevIndex) {
		prevSf = fs.Data[prevIndex]
	} else {
		prevSf = NewShadowFrame(index)
	}

	nif := prevSf.Copy()
	return nif.(*ShadowFrame)
}

func (fs *ShadowFrames) Append(value *ShadowFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Indexes.InsertNoReplace(mcore.Float32(value.Index))

	if value.Registered {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Float32(value.Index))
	}

	fs.Data[value.Index] = value
}
