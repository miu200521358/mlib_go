package vmd

import (
	"slices"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"

)

type BoneNameFrames struct {
	*mcore.IndexModelCorrection[*BoneFrame]
	Name              string // ボーン名
	IkIndexes         []int  // IK計算済みキーフレリスト
	RegisteredIndexes []int  // 登録対象キーフレリスト
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		IndexModelCorrection: mcore.NewIndexModelCorrection[*BoneFrame](),
		Name:                 name,
		IkIndexes:            []int{},
		RegisteredIndexes:    []int{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (bnfs *BoneNameFrames) GetRangeIndexes(index int) (int, int) {
	if len(bnfs.Indexes) == 0 {
		return 0, 0
	}

	prevIndex := 0
	nextIndex := index

	if idx := sort.SearchInts(bnfs.Indexes, index); idx == 0 {
		prevIndex = 0
	} else {
		prevIndex = bnfs.Indexes[idx-1]
	}

	if idx := sort.SearchInts(bnfs.Indexes, index); idx == len(bnfs.Indexes) {
		nextIndex = slices.Max(bnfs.Indexes)
	} else {
		nextIndex = bnfs.Indexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (bnfs *BoneNameFrames) GetItem(index int) *BoneFrame {
	if bnfs == nil {
		return NewBoneFrame(index)
	}
	if index < 0 {
		// マイナス指定の場合、後ろからの順番に置き換える
		index = len(bnfs.Data) + index
		return bnfs.Data[bnfs.Indexes[index]]
	}
	if slices.Contains(bnfs.Indexes, index) {
		return bnfs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := bnfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && slices.Contains(bnfs.Indexes, nextIndex) {
		nextBf := bnfs.Data[nextIndex]
		copied := &BoneFrame{
			BaseFrame:     NewVmdBaseFrame(index),
			Position:      nextBf.Position.Copy(),
			LocalPosition: nextBf.LocalPosition.Copy(),
			Rotation:      nextBf.Rotation.Copy(),
			LocalRotation: nextBf.LocalRotation.Copy(),
			Scale:         nextBf.Scale.Copy(),
			LocalScale:    nextBf.LocalScale.Copy(),
			// IKとかの計算値はコピーしないで初期値
			IkRotation: mmath.NewRotationModelByDegrees(mmath.NewMVec3()),
		}
		return copied
	}

	var prevBf, nextBf *BoneFrame
	if slices.Contains(bnfs.Indexes, prevIndex) {
		prevBf = bnfs.Data[prevIndex]
	} else {
		prevBf = NewBoneFrame(index)
	}
	if slices.Contains(bnfs.Indexes, nextIndex) {
		nextBf = bnfs.Data[nextIndex]
	} else {
		nextBf = NewBoneFrame(index)
	}

	bf := NewBoneFrame(index)

	xy, yy, zy, ry := nextBf.Curves.Evaluate(prevIndex, index, nextIndex)

	qq := prevBf.Rotation.GetQuaternion().Slerp(nextBf.Rotation.GetQuaternion(), ry)
	bf.Rotation.SetQuaternion(&qq)

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

	return bf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (bnfs *BoneNameFrames) Append(value *BoneFrame) {
	if !slices.Contains(bnfs.Indexes, value.Index) {
		bnfs.Indexes = append(bnfs.Indexes, value.Index)
		sort.Ints(bnfs.Indexes)
	}
	if value.IkRegistered && !slices.Contains(bnfs.IkIndexes, value.Index) {
		bnfs.IkIndexes = append(bnfs.IkIndexes, value.Index)
		sort.Ints(bnfs.IkIndexes)
	}

	if value.Registered {
		if !slices.Contains(bnfs.RegisteredIndexes, value.Index) {
			bnfs.RegisteredIndexes = append(bnfs.RegisteredIndexes, value.Index)
			sort.Ints(bnfs.RegisteredIndexes)
		}
		// 補間曲線を分割する
		prevIndex, nextIndex := bnfs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextBf := bnfs.Data[nextIndex]
			// 自分の前後にフレームがある場合、分割する
			value.Curves.TranslateX, nextBf.Curves.TranslateX =
				mmath.SplitCurve(nextBf.Curves.TranslateX, prevIndex, value.Index, nextIndex)
			value.Curves.TranslateY, nextBf.Curves.TranslateY =
				mmath.SplitCurve(nextBf.Curves.TranslateY, prevIndex, value.Index, nextIndex)
			value.Curves.TranslateZ, nextBf.Curves.TranslateZ =
				mmath.SplitCurve(nextBf.Curves.TranslateZ, prevIndex, value.Index, nextIndex)
			value.Curves.Rotate, nextBf.Curves.Rotate =
				mmath.SplitCurve(nextBf.Curves.Rotate, prevIndex, value.Index, nextIndex)
		}
	}

	bnfs.Data[value.Index] = value
}
