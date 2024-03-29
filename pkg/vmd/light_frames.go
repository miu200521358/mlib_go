package vmd

import (
	"sync"

	"github.com/petar/GoLLRB/llrb"

	"github.com/miu200521358/mlib_go/pkg/mcore"

)

type LightFrames struct {
	*mcore.IndexFloatModels[*LightFrame]
	RegisteredIndexes *mcore.FloatIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex        // マップアクセス制御用
}

func NewLightFrames() *LightFrames {
	return &LightFrames{
		IndexFloatModels:  mcore.NewIndexFloatModelCorrection[*LightFrame](),
		RegisteredIndexes: mcore.NewFloatIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *LightFrames) GetRangeIndexes(index float32) (float32, float32) {

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
func (fs *LightFrames) GetItem(index float32) *LightFrame {
	if val, ok := fs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex && fs.Indexes.Has(nextIndex) {
		nextIf := fs.Data[nextIndex]
		copied := nextIf.Copy()
		return copied.(*LightFrame)
	}

	var prevIf *LightFrame
	if fs.Indexes.Has(prevIndex) {
		prevIf = fs.Data[prevIndex]
	} else {
		prevIf = NewLightFrame(index)
	}

	nif := prevIf.Copy()
	return nif.(*LightFrame)
}

func (fs *LightFrames) Append(value *LightFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Indexes.InsertNoReplace(mcore.Float32(value.Index))

	if value.Registered {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Float32(value.Index))
	}

	fs.Data[value.Index] = value
}
