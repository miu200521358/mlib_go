package vmd

import (
	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type CameraFrames struct {
	*mcore.IndexFloatModelCorrection[*CameraFrame]
	RegisteredIndexes *mcore.TreeIndexes[mcore.Float32] // 登録対象キーフレリスト
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*CameraFrame](),
		RegisteredIndexes:         mcore.NewTreeIndexes[mcore.Float32](),
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (cfs *CameraFrames) GetRangeIndexes(index mcore.Float32) (mcore.Float32, mcore.Float32) {
	if cfs.RegisteredIndexes.IsEmpty() {
		return 0.0, 0.0
	}

	nowIndexes := cfs.RegisteredIndexes.Search(index)
	if nowIndexes != nil {
		return index, index
	}

	var prevIndex, nextIndex mcore.Float32

	prevIndexes := cfs.RegisteredIndexes.SearchLeft(index)
	nextIndexes := cfs.RegisteredIndexes.SearchRight(index)

	if prevIndexes == nil {
		prevIndex = mcore.NewFloat32(0)
	} else {
		prevIndex = prevIndexes.Value
	}

	if nextIndexes == nil {
		nextIndex = index
	} else {
		nextIndex = nextIndexes.Value
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (cfs *CameraFrames) GetItem(index mcore.Float32) *CameraFrame {
	if val, ok := cfs.Data[index]; ok {
		return val
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := cfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && cfs.Indexes.Contains(nextIndex) {
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
	if cfs.Indexes.Contains(prevIndex) {
		prevCf = cfs.Data[prevIndex]
	} else {
		prevCf = NewCameraFrame(index)
	}
	if cfs.Indexes.Contains(nextIndex) {
		nextCf = cfs.Data[nextIndex]
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
func (cfs *CameraFrames) Append(value *CameraFrame) {
	if !cfs.Indexes.Contains(value.Index) {
		cfs.Indexes.Insert(value.Index)
	}

	if value.Registered {
		if !cfs.RegisteredIndexes.Contains(value.Index) {
			cfs.RegisteredIndexes.Insert(value.Index)
		}
		// 補間曲線を分割する
		prevIndex, nextIndex := cfs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextCf := cfs.Data[nextIndex]
			pi := float32(prevIndex)
			vi := float32(value.Index)
			ni := float32(nextIndex)
			// 自分の前後にフレームがある場合、分割する
			value.Curves.TranslateX, nextCf.Curves.TranslateX =
				mmath.SplitCurve(nextCf.Curves.TranslateX, pi, vi, ni)
			value.Curves.TranslateY, nextCf.Curves.TranslateY =
				mmath.SplitCurve(nextCf.Curves.TranslateY, pi, vi, ni)
			value.Curves.TranslateZ, nextCf.Curves.TranslateZ =
				mmath.SplitCurve(nextCf.Curves.TranslateZ, pi, vi, ni)
			value.Curves.Rotate, nextCf.Curves.Rotate =
				mmath.SplitCurve(nextCf.Curves.Rotate, pi, vi, ni)
			value.Curves.Distance, nextCf.Curves.Distance =
				mmath.SplitCurve(nextCf.Curves.Distance, pi, vi, ni)
			value.Curves.ViewOfAngle, nextCf.Curves.ViewOfAngle =
				mmath.SplitCurve(nextCf.Curves.ViewOfAngle, pi, vi, ni)
		}
	}

	cfs.Data[value.Index] = value
}
