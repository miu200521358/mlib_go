package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneNameFrames struct {
	*mcore.IndexFloatModelCorrection[*BoneFrame]
	Name              string                            // ボーン名
	IkIndexes         *mcore.TreeIndexes[mcore.Float32] // IK計算済みキーフレリスト
	RegisteredIndexes *mcore.TreeIndexes[mcore.Float32] // 登録対象キーフレリスト
	lock              sync.RWMutex                      // マップアクセス制御用
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*BoneFrame](),
		Name:                      name,
		IkIndexes:                 mcore.NewTreeIndexes[mcore.Float32](),
		RegisteredIndexes:         mcore.NewTreeIndexes[mcore.Float32](),
		lock:                      sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (bnfs *BoneNameFrames) GetRangeIndexes(index mcore.Float32) (mcore.Float32, mcore.Float32) {
	if bnfs.RegisteredIndexes.IsEmpty() {
		return 0.0, 0.0
	}

	nowIndexes := bnfs.RegisteredIndexes.Search(index)
	if nowIndexes != nil {
		return index, index
	}

	var prevIndex, nextIndex mcore.Float32

	prevIndexes := bnfs.RegisteredIndexes.SearchLeft(index)
	nextIndexes := bnfs.RegisteredIndexes.SearchRight(index)

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
func (bnfs *BoneNameFrames) GetItem(index mcore.Float32) *BoneFrame {
	if bnfs == nil {
		return NewBoneFrame(index)
	}

	bnfs.lock.RLock()
	defer bnfs.lock.RUnlock()

	if bnfs.Indexes.Contains(index) {
		return bnfs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := bnfs.GetRangeIndexes(index)

	if prevIndex == nextIndex {
		if bnfs.Indexes.Contains(nextIndex) {
			nextBf := bnfs.Data[nextIndex]
			copied := &BoneFrame{
				BaseFrame:          NewVmdBaseFrame(index),
				Position:           nextBf.Position.Copy(),
				MorphPosition:      nextBf.MorphPosition.Copy(),
				LocalPosition:      nextBf.LocalPosition.Copy(),
				MorphLocalPosition: nextBf.MorphLocalPosition.Copy(),
				Rotation:           nextBf.Rotation.Copy(),
				MorphRotation:      nextBf.MorphRotation.Copy(),
				LocalRotation:      nextBf.LocalRotation.Copy(),
				MorphLocalRotation: nextBf.MorphLocalRotation.Copy(),
				Scale:              nextBf.Scale.Copy(),
				MorphScale:         nextBf.MorphScale.Copy(),
				LocalScale:         nextBf.LocalScale.Copy(),
				MorphLocalScale:    nextBf.MorphLocalScale.Copy(),
				// IKとかの計算値はコピーしないで初期値
				IkRotation: mmath.NewRotationModel(),
			}
			return copied
		} else {
			return NewBoneFrame(index)
		}
	}

	var prevBf, nextBf *BoneFrame
	if bnfs.Indexes.Contains(prevIndex) {
		prevBf = bnfs.Data[prevIndex]
	} else {
		prevBf = NewBoneFrame(index)
	}
	if bnfs.Indexes.Contains(nextIndex) {
		nextBf = bnfs.Data[nextIndex]
	} else {
		nextBf = NewBoneFrame(index)
	}

	bf := NewBoneFrame(index)

	xy, yy, zy, ry := nextBf.Curves.Evaluate(float32(prevIndex), float32(index), float32(nextIndex))

	qq := prevBf.Rotation.GetQuaternion().Slerp(nextBf.Rotation.GetQuaternion(), ry)
	bf.Rotation.SetQuaternion(qq)

	prevX := mmath.MVec4{
		prevBf.Position.GetX(), prevBf.LocalPosition.GetX(), prevBf.Scale.GetX(), prevBf.LocalScale.GetX()}
	nextX := mmath.MVec4{
		nextBf.Position.GetX(), nextBf.LocalPosition.GetX(), nextBf.Scale.GetX(), nextBf.LocalScale.GetX()}
	nowX := mmath.LerpVec4(&prevX, &nextX, xy)
	bf.Position.SetX(nowX[0])
	bf.LocalPosition.SetX(nowX[1])
	bf.Scale.SetX(nowX[2])
	bf.LocalScale.SetX(nowX[3])

	prevY := mmath.MVec4{
		prevBf.Position.GetY(), prevBf.LocalPosition.GetY(), prevBf.Scale.GetY(), prevBf.LocalScale.GetY()}
	nextY := mmath.MVec4{
		nextBf.Position.GetY(), nextBf.LocalPosition.GetY(), nextBf.Scale.GetY(), nextBf.LocalScale.GetY()}
	nowY := mmath.LerpVec4(&prevY, &nextY, yy)
	bf.Position.SetY(nowY[0])
	bf.LocalPosition.SetY(nowY[1])
	bf.Scale.SetY(nowY[2])
	bf.LocalScale.SetY(nowY[3])

	prevZ := mmath.MVec4{
		prevBf.Position.GetZ(), prevBf.LocalPosition.GetZ(), prevBf.Scale.GetZ(), prevBf.LocalScale.GetZ()}
	nextZ := mmath.MVec4{
		nextBf.Position.GetZ(), nextBf.LocalPosition.GetZ(), nextBf.Scale.GetZ(), nextBf.LocalScale.GetZ()}
	nowZ := mmath.LerpVec4(&prevZ, &nextZ, zy)
	bf.Position.SetZ(nowZ[0])
	bf.LocalPosition.SetZ(nowZ[1])
	bf.Scale.SetZ(nowZ[2])
	bf.LocalScale.SetZ(nowZ[3])

	// IKとかの計算値はコピーしないで初期値
	bf.IkRotation = mmath.NewRotationModel()

	return bf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (bnfs *BoneNameFrames) Append(value *BoneFrame, isSort bool) {
	bnfs.lock.Lock()
	defer bnfs.lock.Unlock()

	if !bnfs.Indexes.Contains(value.Index) {
		bnfs.Indexes.Insert(value.Index)
	}
	if value.IkRegistered && !bnfs.IkIndexes.Contains(value.Index) {
		bnfs.IkIndexes.Insert(value.Index)
	}

	if value.Registered {
		if !bnfs.RegisteredIndexes.Contains(value.Index) {
			bnfs.RegisteredIndexes.Insert(value.Index)
		}
	}

	if value.Registered && isSort {
		// 補間曲線を分割する
		prevIndex, nextIndex := bnfs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextBf := bnfs.Data[nextIndex]
			pi := float32(prevIndex)
			vi := float32(value.Index)
			ni := float32(nextIndex)
			// 自分の前後にフレームがある場合、分割する
			value.Curves.TranslateX, nextBf.Curves.TranslateX =
				mmath.SplitCurve(nextBf.Curves.TranslateX, pi, vi, ni)
			value.Curves.TranslateY, nextBf.Curves.TranslateY =
				mmath.SplitCurve(nextBf.Curves.TranslateY, pi, vi, ni)
			value.Curves.TranslateZ, nextBf.Curves.TranslateZ =
				mmath.SplitCurve(nextBf.Curves.TranslateZ, pi, vi, ni)
			value.Curves.Rotate, nextBf.Curves.Rotate =
				mmath.SplitCurve(nextBf.Curves.Rotate, pi, vi, ni)
		}
	}

	bnfs.Data[value.Index] = value
}

func (bnfs *BoneNameFrames) GetMaxFrame() mcore.Float32 {
	if bnfs.RegisteredIndexes.IsEmpty() {
		return 0
	}

	return bnfs.RegisteredIndexes.GetMax()
}
