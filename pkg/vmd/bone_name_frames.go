package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/petar/GoLLRB/llrb"
)

type BoneNameFrames struct {
	*mcore.IndexModels[*BoneFrame]
	Name              string            // ボーン名
	IkIndexes         *mcore.IntIndexes // IK計算済みキーフレリスト
	RegisteredIndexes *mcore.IntIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex      // マップアクセス制御用
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		IndexModels:       mcore.NewIndexModels[*BoneFrame](),
		Name:              name,
		IkIndexes:         mcore.NewIntIndexes(),
		RegisteredIndexes: mcore.NewIntIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *BoneNameFrames) GetRangeIndexes(index int) (int, int) {
	if fs.RegisteredIndexes.Len() == 0 {
		return 0, 0
	}

	if fs.RegisteredIndexes.Max() < index {
		return fs.RegisteredIndexes.Max(), fs.RegisteredIndexes.Max()
	}

	prevIndex := mcore.Int(0)
	nextIndex := mcore.Int(index)

	fs.RegisteredIndexes.DescendLessOrEqual(mcore.Int(index), func(i llrb.Item) bool {
		prevIndex = i.(mcore.Int)
		return false
	})

	fs.RegisteredIndexes.AscendGreaterOrEqual(mcore.Int(index), func(i llrb.Item) bool {
		nextIndex = i.(mcore.Int)
		return false
	})

	return int(prevIndex), int(nextIndex)
}

// キーフレ計算結果を返す
func (fs *BoneNameFrames) GetItem(index int) *BoneFrame {
	if fs == nil {
		return NewBoneFrame(index)
	}

	fs.lock.RLock()
	defer fs.lock.RUnlock()

	if _, ok := fs.Data[index]; ok {
		return fs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex {
		if _, ok := fs.Data[nextIndex]; ok {
			nextBf := fs.Data[nextIndex]
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
				IkRotation: mmath.NewRotation(),
				// 補間曲線はとりあえず初期値
				Curves: NewBoneCurves(),
			}
			return copied
		} else {
			return NewBoneFrame(index)
		}
	}

	var prevBf, nextBf *BoneFrame
	if _, ok := fs.Data[prevIndex]; ok {
		prevBf = fs.Data[prevIndex]
	} else {
		prevBf = NewBoneFrame(index)
	}
	if _, ok := fs.Data[nextIndex]; ok {
		nextBf = fs.Data[nextIndex]
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

	// IKとかの計算値はコピーしないで初期値
	bf.IkRotation = mmath.NewRotation()

	return bf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (fs *BoneNameFrames) Append(value *BoneFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if _, ok := fs.Data[value.Index]; !ok {
		fs.Indexes[value.Index] = value.Index
	}
	if value.IkRegistered {
		fs.IkIndexes.ReplaceOrInsert(mcore.Int(value.Index))
	}

	if value.Registered {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Int(value.Index))

		// 補間曲線を分割する
		prevIndex, nextIndex := fs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextBf := fs.Data[nextIndex]
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

	fs.Data[value.Index] = value
}

func (fs *BoneNameFrames) GetMaxFrame() int {
	return fs.RegisteredIndexes.Max()
}

func (fs *BoneNameFrames) GetMinFrame() int {
	return fs.RegisteredIndexes.Min()
}
