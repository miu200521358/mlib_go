package vmd

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type BoneNameFrames struct {
	*mcore.IndexFloatModelCorrection[*BoneFrame]
	Name              string              // ボーン名
	IkIndexes         map[float32]float32 // IK計算済みキーフレリスト
	RegisteredIndexes map[float32]float32 // 登録対象キーフレリスト
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*BoneFrame](),
		Name:                      name,
		IkIndexes:                 make(map[float32]float32, 0),
		RegisteredIndexes:         make(map[float32]float32, 0),
	}
}

func (c *BoneNameFrames) ContainsIk(key float32) bool {
	_, ok := c.IkIndexes[key]
	return ok
}

func (c *BoneNameFrames) ContainsRegistered(key float32) bool {
	_, ok := c.RegisteredIndexes[key]
	return ok
}

func (c *BoneNameFrames) GetSortedIkIndexes() []float32 {
	keys := make([]float32, 0, len(c.IkIndexes))
	for key := range c.IkIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func (c *BoneNameFrames) GetSortedRegisteredIndexes() []float32 {
	keys := make([]float32, 0, len(c.RegisteredIndexes))
	for key := range c.RegisteredIndexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// 指定したキーフレの前後のキーフレ番号を返す
func (bnfs *BoneNameFrames) GetRangeIndexes(index float32) (float32, float32) {
	if len(bnfs.RegisteredIndexes) == 0 {
		return 0.0, 0.0
	}

	prevIndex := float32(0.0)
	nextIndex := index

	registeredIndexes := bnfs.GetSortedRegisteredIndexes()

	if idx := mutils.SearchFloat32s(registeredIndexes, index); idx == 0 {
		prevIndex = 0.0
	} else {
		prevIndex = registeredIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(registeredIndexes, index); idx == len(registeredIndexes) {
		nextIndex = slices.Max(registeredIndexes)
	} else {
		nextIndex = registeredIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (bnfs *BoneNameFrames) GetItem(index float32) *BoneFrame {
	if bnfs == nil {
		return NewBoneFrame(index)
	}
	if bnfs.Contains(index) {
		return bnfs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := bnfs.GetRangeIndexes(index)

	if prevIndex == nextIndex && bnfs.Contains(nextIndex) {
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
	if bnfs.Contains(prevIndex) {
		prevBf = bnfs.Data[prevIndex]
	} else {
		prevBf = NewBoneFrame(index)
	}
	if bnfs.Contains(nextIndex) {
		nextBf = bnfs.Data[nextIndex]
	} else {
		nextBf = NewBoneFrame(index)
	}

	bf := NewBoneFrame(index)

	xy, yy, zy, ry := nextBf.Curves.Evaluate(prevIndex, index, nextIndex)

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

	return bf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (bnfs *BoneNameFrames) Append(value *BoneFrame) {
	if !bnfs.Contains(value.Index) {
		bnfs.Indexes[value.Index] = value.Index
	}
	if value.IkRegistered && !bnfs.ContainsIk(value.Index) {
		bnfs.IkIndexes[value.Index] = value.Index
	}

	if value.Registered {
		if !bnfs.ContainsRegistered(value.Index) {
			bnfs.RegisteredIndexes[value.Index] = value.Index
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
