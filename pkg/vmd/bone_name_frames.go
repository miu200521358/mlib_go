package vmd

import (
	"sync"

	"github.com/petar/GoLLRB/llrb"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

type BoneNameFrames struct {
	*mcore.IndexFloatModels[*BoneFrame]
	Name              string              // ボーン名
	IkIndexes         *mcore.FloatIndexes // IK計算済みキーフレリスト
	RegisteredIndexes *mcore.FloatIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex        // マップアクセス制御用
}

func NewBoneNameFrames(name string) *BoneNameFrames {
	return &BoneNameFrames{
		IndexFloatModels:  mcore.NewIndexFloatModelCorrection[*BoneFrame](),
		Name:              name,
		IkIndexes:         mcore.NewFloatIndexes(),
		RegisteredIndexes: mcore.NewFloatIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *BoneNameFrames) GetRangeIndexes(index float32) (float32, float32) {
	if fs.RegisteredIndexes.Len() == 0 {
		return 0.0, 0.0
	}

	if fs.RegisteredIndexes.Max() < index {
		return float32(fs.RegisteredIndexes.Max()), float32(fs.RegisteredIndexes.Max())
	}

	prevIndex := mcore.Float32(0.0)
	nextIndex := mcore.Float32(index)

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
func (fs *BoneNameFrames) GetItem(index float32) *BoneFrame {
	if fs == nil {
		return NewBoneFrame(index)
	}

	fs.lock.RLock()
	defer fs.lock.RUnlock()

	if fs.Indexes.Has(index) {
		return fs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex {
		if fs.Indexes.Has(nextIndex) {
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
				IkRotation: mmath.NewRotationModel(),
			}
			return copied
		} else {
			return NewBoneFrame(index)
		}
	}

	var prevBf, nextBf *BoneFrame
	if fs.Indexes.Has(prevIndex) {
		prevBf = fs.Data[prevIndex]
	} else {
		prevBf = NewBoneFrame(index)
	}
	if fs.Indexes.Has(nextIndex) {
		nextBf = fs.Data[nextIndex]
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
func (fs *BoneNameFrames) Append(value *BoneFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Indexes.InsertNoReplace(mcore.Float32(value.Index))

	if value.IkRegistered {
		fs.IkIndexes.InsertNoReplace(mcore.Float32(value.Index))
	}

	if value.Registered {
		fs.RegisteredIndexes.InsertNoReplace(mcore.Float32(value.Index))
	}

	fs.Data[value.Index] = value
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (fs *BoneNameFrames) Insert(value *BoneFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	fs.Indexes.InsertNoReplace(mcore.Float32(value.Index))

	if value.IkRegistered {
		fs.IkIndexes.InsertNoReplace(mcore.Float32(value.Index))
	}

	if value.Registered && !fs.RegisteredIndexes.Has(value.Index) {
		// 補間曲線を分割する
		prevIndex, nextIndex := fs.GetRangeIndexes(value.Index)
		if nextIndex > value.Index && prevIndex < value.Index {
			nextBf := fs.Data[nextIndex]
			// 自分の前後にフレームがある場合、分割する
			value.Curves.TranslateX, nextBf.Curves.TranslateX =
				mmath.SplitCurve(nextBf.Curves.TranslateX, float32(prevIndex), float32(value.Index), float32(nextIndex))
			value.Curves.TranslateY, nextBf.Curves.TranslateY =
				mmath.SplitCurve(nextBf.Curves.TranslateY, float32(prevIndex), float32(value.Index), float32(nextIndex))
			value.Curves.TranslateZ, nextBf.Curves.TranslateZ =
				mmath.SplitCurve(nextBf.Curves.TranslateZ, float32(prevIndex), float32(value.Index), float32(nextIndex))
			value.Curves.Rotate, nextBf.Curves.Rotate =
				mmath.SplitCurve(nextBf.Curves.Rotate, float32(prevIndex), float32(value.Index), float32(nextIndex))
		}
	}

	if value.Registered {
		fs.RegisteredIndexes.InsertNoReplace(mcore.Float32(value.Index))
	}

	fs.Data[value.Index] = value
}
