package vmd

import (
	"sync"

	"github.com/petar/GoLLRB/llrb"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type CameraFrames struct {
	*mcore.IndexFloatModels[*CameraFrame]
	RegisteredIndexes *mcore.FloatIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex        // マップアクセス制御用
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		IndexFloatModels:  mcore.NewIndexFloatModelCorrection[*CameraFrame](),
		RegisteredIndexes: mcore.NewFloatIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *CameraFrames) GetRangeIndexes(index float32) (float32, float32) {

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
func (fs *CameraFrames) GetItem(index float32) *CameraFrame {
	if val, ok := fs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex && fs.Indexes.Has(nextIndex) {
		nextCf := fs.Data[nextIndex]
		copied := &CameraFrame{
			BaseFrame:        NewVmdBaseFrame(index),
			Position:         nextCf.Position.Copy(),
			Rotation:         nextCf.Rotation.Copy(),
			Distance:         nextCf.Distance,
			ViewOfAngle:      nextCf.ViewOfAngle,
			IsPerspectiveOff: nextCf.IsPerspectiveOff,
			Curves:           nextCf.Curves.Copy(),
			// # IKとかの計算値はコピーしない
		}
		return copied
	}

	var prevCf, nextCf *CameraFrame
	if fs.Indexes.Has(prevIndex) {
		prevCf = fs.Data[prevIndex]
	} else {
		prevCf = NewCameraFrame(index)
	}
	if fs.Indexes.Has(nextIndex) {
		nextCf = fs.Data[nextIndex]
	} else {
		nextCf = NewCameraFrame(index)
	}

	cf := NewCameraFrame(index)

	xy, yy, zy, ry, dy, vy := nextCf.Curves.Evaluate(float32(prevIndex), float32(index), float32(nextIndex))

	qq := prevCf.Rotation.GetQuaternion().Slerp(nextCf.Rotation.GetQuaternion(), ry)
	cf.Rotation.SetQuaternion(qq)

	cf.Position.SetX(mmath.LerpFloat(prevCf.Position.GetX(), nextCf.Position.GetX(), xy))
	cf.Position.SetY(mmath.LerpFloat(prevCf.Position.GetY(), nextCf.Position.GetY(), yy))
	cf.Position.SetZ(mmath.LerpFloat(prevCf.Position.GetZ(), nextCf.Position.GetZ(), zy))

	cf.Distance = mmath.LerpFloat(prevCf.Distance, nextCf.Distance, dy)
	cf.ViewOfAngle = int(mmath.LerpFloat(float64(prevCf.ViewOfAngle), float64(nextCf.ViewOfAngle), vy))
	cf.IsPerspectiveOff = nextCf.IsPerspectiveOff

	return cf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (fs *CameraFrames) Append(value *CameraFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Indexes.InsertNoReplace(mcore.Float32(value.Index))

	if value.Registered {
		fs.RegisteredIndexes.InsertNoReplace(mcore.Float32(value.Index))
	}

	fs.Data[value.Index] = value
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (fs *CameraFrames) Insert(value *CameraFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Indexes.InsertNoReplace(mcore.Float32(value.Index))

	if value.Registered && !fs.RegisteredIndexes.Has(value.Index) {
		// 補間曲線を分割する
		prevIndex, nextIndex := fs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextCf := fs.Data[nextIndex]
			// 自分の前後にフレームがある場合、分割する
			value.Curves.TranslateX, nextCf.Curves.TranslateX =
				mmath.SplitCurve(nextCf.Curves.TranslateX, prevIndex, value.Index, nextIndex)
			value.Curves.TranslateY, nextCf.Curves.TranslateY =
				mmath.SplitCurve(nextCf.Curves.TranslateY, prevIndex, value.Index, nextIndex)
			value.Curves.TranslateZ, nextCf.Curves.TranslateZ =
				mmath.SplitCurve(nextCf.Curves.TranslateZ, prevIndex, value.Index, nextIndex)
			value.Curves.Rotate, nextCf.Curves.Rotate =
				mmath.SplitCurve(nextCf.Curves.Rotate, prevIndex, value.Index, nextIndex)
			value.Curves.Distance, nextCf.Curves.Distance =
				mmath.SplitCurve(nextCf.Curves.Distance, prevIndex, value.Index, nextIndex)
			value.Curves.ViewOfAngle, nextCf.Curves.ViewOfAngle =
				mmath.SplitCurve(nextCf.Curves.ViewOfAngle, prevIndex, value.Index, nextIndex)
		}
	}

	if value.Registered {
		fs.RegisteredIndexes.InsertNoReplace(mcore.Float32(value.Index))
	}

	fs.Data[value.Index] = value
}
