package vmd

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/petar/GoLLRB/llrb"
)

type MorphNameFrames struct {
	*mcore.IndexModels[*MorphFrame]
	Name              string            // ボーン名
	RegisteredIndexes *mcore.IntIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex      // マップアクセス制御用
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		IndexModels:       mcore.NewIndexModels[*MorphFrame](),
		Name:              name,
		RegisteredIndexes: mcore.NewIntIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *MorphNameFrames) GetRangeIndexes(index int) (int, int) {
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
func (fs *MorphNameFrames) GetItem(index int) *MorphFrame {
	if fs == nil {
		return NewMorphFrame(index)
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
			nextMf := fs.Data[nextIndex]
			copied := &MorphFrame{
				BaseFrame: NewVmdBaseFrame(index),
				Ratio:     nextMf.Ratio,
			}
			return copied
		} else {
			return NewMorphFrame(index)
		}
	}

	var prevMf, nextMf *MorphFrame
	if _, ok := fs.Data[prevIndex]; ok {
		prevMf = fs.Data[prevIndex]
	} else {
		prevMf = NewMorphFrame(index)
	}
	if _, ok := fs.Data[nextIndex]; ok {
		nextMf = fs.Data[nextIndex]
	} else {
		nextMf = NewMorphFrame(index)
	}

	mf := NewMorphFrame(index)

	ry := (index - prevIndex) / (nextIndex - prevIndex)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*float64(ry)

	return mf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (fs *MorphNameFrames) Append(value *MorphFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if _, ok := fs.Data[value.Index]; !ok {
		fs.Indexes[value.Index] = value.Index
	}

	if value.Registered {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Int(value.Index))
	}

	fs.Data[value.Index] = value
}

func (fs *MorphNameFrames) GetMaxFrame() int {
	return fs.RegisteredIndexes.Max()
}

func (fs *MorphNameFrames) GetMinFrame() int {
	return fs.RegisteredIndexes.Min()
}

func (fs *MorphNameFrames) AnimateVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			p := offset.Position.MuledScalar(mf.Ratio)
			delta.Position.Add(&p)
		}
	}
}

func (fs *MorphNameFrames) AnimateAfterVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			ap := offset.Position.MuledScalar(mf.Ratio)
			delta.AfterPosition.Add(&ap)
		}
	}
}

func (fs *MorphNameFrames) AnimateUv(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			uv := offset.Uv.MuledScalar(mf.Ratio).GetXY()
			delta.Uv.Add(uv)
		}
	}
}

func (fs *MorphNameFrames) AnimateUv1(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := fs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			uv := offset.Uv.MuledScalar(mf.Ratio)
			delta.Uv1.Add(uv.GetXY())
		}
	}
}

func (fs *MorphNameFrames) AnimateBone(
	frame int,
	model *pmx.PmxModel,
	deltas *BoneMorphDeltas,
) {
	mf := fs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex && offset.BoneIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.BoneIndex]
			mp := offset.Position.MuledScalar(mf.Ratio)
			delta.MorphPosition.Add(&mp)
			mlp := offset.LocalPosition.MuledScalar(mf.Ratio)
			delta.MorphLocalPosition.Add(&mlp)
			deltaRad := offset.Rotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphRotation.SetQuaternion(delta.MorphRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromRadians(deltaRad.GetX(), deltaRad.GetY(), deltaRad.GetZ())))
			deltaLocalRad := offset.LocalRotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphLocalRotation.SetQuaternion(delta.MorphLocalRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromRadians(deltaLocalRad.GetX(), deltaLocalRad.GetY(), deltaLocalRad.GetZ())))
			ms := offset.Scale.MuledScalar(mf.Ratio)
			delta.MorphScale.Add(&ms)
			mls := offset.LocalScale.MuledScalar(mf.Ratio)
			delta.MorphLocalScale.Add(&mls)
		}
	}
}

// AnimateMaterial 材質モーフの適用
func (fs *MorphNameFrames) AnimateMaterial(
	frame int,
	model *pmx.PmxModel,
	deltas *MaterialMorphDeltas,
) {
	mf := fs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(fs.Name)
	// 乗算→加算の順で処理
	for _, calcMode := range []pmx.MaterialMorphCalcMode{pmx.CALC_MODE_MULTIPLICATION, pmx.CALC_MODE_ADDITION} {
		for _, o := range morph.Offsets {
			offset := o.(*pmx.MaterialMorphOffset)
			if offset.CalcMode != calcMode {
				continue
			}
			if offset.MaterialIndex < 0 {
				// 全材質対象の場合
				for _, delta := range deltas.Data {
					if calcMode == pmx.CALC_MODE_MULTIPLICATION {
						delta.Mul(offset, mf.Ratio)
					} else {
						delta.Add(offset, mf.Ratio)
					}
				}
			} else if 0 < offset.MaterialIndex && offset.MaterialIndex <= len(deltas.Data) {
				// 特定材質のみの場合
				if calcMode == pmx.CALC_MODE_MULTIPLICATION {
					deltas.Data[offset.MaterialIndex].Mul(offset, mf.Ratio)
				} else {
					deltas.Data[offset.MaterialIndex].Add(offset, mf.Ratio)
				}
			}
		}
	}
}
