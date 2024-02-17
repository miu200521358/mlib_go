package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

type CameraFrames struct {
	*mcore.IndexFloatModelCorrection[*CameraFrame]
	RegisteredIndexes map[float32]float32 // 登録対象キーフレリスト
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*CameraFrame](),
		RegisteredIndexes:         make(map[float32]float32, 0),
	}
}

func (c *CameraFrames) ContainsRegistered(key float32) bool {
	_, ok := c.RegisteredIndexes[key]
	return ok
}

func (c *CameraFrames) GetSortedRegisteredIndexes() []float32 {
	keys := make([]float32, 0, len(c.RegisteredIndexes))
	for key := range c.RegisteredIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// 指定したキーフレの前後のキーフレ番号を返す
func (cfs *CameraFrames) GetRangeIndexes(index float32) (float32, float32) {

	prevIndex := float32(0)
	nextIndex := index

	sortedIndexes := cfs.GetSortedIndexes()

	if idx := mutils.SearchFloat32s(sortedIndexes, index); idx == 0 {
		prevIndex = 0.0
	} else {
		prevIndex = sortedIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(sortedIndexes, index); idx == len(sortedIndexes) {
		nextIndex = slices.Max(sortedIndexes)
	} else {
		nextIndex = sortedIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (cfs *CameraFrames) GetItem(index float32) *CameraFrame {
	if val, ok := cfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := cfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && cfs.Contains(nextIndex) {
		nextCf := cfs.Data[nextIndex]
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
	if cfs.Contains(prevIndex) {
		prevCf = cfs.Data[prevIndex]
	} else {
		prevCf = NewCameraFrame(index)
	}
	if cfs.Contains(nextIndex) {
		nextCf = cfs.Data[nextIndex]
	} else {
		nextCf = NewCameraFrame(index)
	}

	cf := NewCameraFrame(index)

	xy, yy, zy, ry, dy, vy := nextCf.Curves.Evaluate(prevIndex, index, nextIndex)

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
func (cfs *CameraFrames) Append(value *CameraFrame) {
	if !cfs.Contains(value.Index) {
		cfs.Indexes[value.Index] = value.Index
	}

	if value.Registered {
		if !cfs.ContainsRegistered(value.Index) {
			cfs.RegisteredIndexes[value.Index] = value.Index
		}
		// 補間曲線を分割する
		prevIndex, nextIndex := cfs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextCf := cfs.Data[nextIndex]
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

	cfs.Data[value.Index] = value
}
