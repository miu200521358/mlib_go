package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type CameraFrames struct {
	*mcore.IndexModels[*CameraFrame]
	RegisteredIndexes []int // 登録対象キーフレリスト
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		IndexModels:       mcore.NewIndexModels[*CameraFrame](),
		RegisteredIndexes: []int{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (cfs *CameraFrames) GetRangeIndexes(index int) (int, int) {

	prevIndex := int(0)
	nextIndex := index

	if idx := mutils.SearchInts(cfs.RegisteredIndexes, index); idx == 0 {
		prevIndex = 0.0
	} else {
		prevIndex = cfs.RegisteredIndexes[idx-1]
	}

	if idx := mutils.SearchInts(cfs.RegisteredIndexes, index); idx == len(cfs.RegisteredIndexes) {
		nextIndex = slices.Max(cfs.RegisteredIndexes)
	} else {
		nextIndex = cfs.RegisteredIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (cfs *CameraFrames) GetItem(index int) *CameraFrame {
	if val, ok := cfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := cfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && slices.Contains(cfs.RegisteredIndexes, nextIndex) {
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
	if slices.Contains(cfs.RegisteredIndexes, prevIndex) {
		prevCf = cfs.Data[prevIndex]
	} else {
		prevCf = NewCameraFrame(index)
	}
	if slices.Contains(cfs.RegisteredIndexes, nextIndex) {
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
	if !slices.Contains(cfs.RegisteredIndexes, value.Index) {
		cfs.RegisteredIndexes = append(cfs.RegisteredIndexes, value.Index)
		mutils.SortInts(cfs.RegisteredIndexes)
	}

	if value.Registered {
		if !slices.Contains(cfs.RegisteredIndexes, value.Index) {
			cfs.RegisteredIndexes = append(cfs.RegisteredIndexes, value.Index)
			mutils.SortInts(cfs.RegisteredIndexes)
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
