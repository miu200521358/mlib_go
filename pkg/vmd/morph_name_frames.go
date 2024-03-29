package vmd

import (
	"sync"

	"github.com/petar/GoLLRB/llrb"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/pmx"

)

type MorphNameFrames struct {
	*mcore.IndexFloatModels[*MorphFrame]
	Name              string              // モーフ名
	RegisteredIndexes *mcore.FloatIndexes // 登録対象キーフレリスト
	lock              sync.RWMutex        // マップアクセス制御用
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		IndexFloatModels:  mcore.NewIndexFloatModels[*MorphFrame](),
		Name:              name,
		RegisteredIndexes: mcore.NewFloatIndexes(),
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (fs *MorphNameFrames) GetRangeIndexes(index float32) (float32, float32) {
	if fs.RegisteredIndexes.Len() == 0 {
		return 0.0, 0.0
	}

	prevIndex := mcore.Float32(0.0)
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
func (fs *MorphNameFrames) GetItem(index float32) *MorphFrame {
	if fs == nil {
		return NewMorphFrame(index)
	}

	fs.lock.RLock()
	defer fs.lock.RUnlock()

	if fs.Has(index) {
		return fs.Get(index)
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := fs.GetRangeIndexes(index)

	if prevIndex == nextIndex {
		if fs.Has(nextIndex) {
			nextMf := fs.Get(nextIndex)
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
	if fs.Has(prevIndex) {
		prevMf = fs.Get(prevIndex)
	} else {
		prevMf = NewMorphFrame(index)
	}
	if fs.Has(nextIndex) {
		nextMf = fs.Get(nextIndex)
	} else {
		nextMf = NewMorphFrame(index)
	}

	mf := NewMorphFrame(index)

	ry := (index - prevIndex) / (nextIndex - prevIndex)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*float64(ry)

	return mf
}

// Append モーフフレームを追加する
func (fs *MorphNameFrames) Append(value *MorphFrame) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if value.Registered {
		fs.RegisteredIndexes.ReplaceOrInsert(mcore.Float32(value.Index))
	}

	fs.ReplaceOrInsert(value)
}

func (fs *MorphNameFrames) AnimateVertex(
	frame float32,
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
			delta.Position.Add(offset.Position.MuledScalar(mf.Ratio))
		}
	}
}

func (fs *MorphNameFrames) AnimateAfterVertex(
	frame float32,
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
			delta.AfterPosition.Add(offset.Position.MuledScalar(mf.Ratio))
		}
	}
}

func (fs *MorphNameFrames) AnimateUv(
	frame float32,
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
	frame float32,
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
	frame float32,
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
			delta.MorphPosition.Add(offset.Position.MuledScalar(mf.Ratio))
			delta.MorphLocalPosition.Add(offset.LocalPosition.MuledScalar(mf.Ratio))
			deltaRad := offset.Rotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphRotation.SetQuaternion(delta.MorphRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromEulerAngles(deltaRad.GetX(), deltaRad.GetY(), deltaRad.GetZ())))
			deltaLocalRad := offset.LocalRotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphLocalRotation.SetQuaternion(delta.MorphLocalRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromEulerAngles(deltaLocalRad.GetX(), deltaLocalRad.GetY(), deltaLocalRad.GetZ())))
			delta.MorphScale.Add(offset.Scale.MuledScalar(mf.Ratio))
			delta.MorphLocalScale.Add(offset.LocalScale.MuledScalar(mf.Ratio))
		}
	}
}

// AnimateMaterial 材質モーフの適用
func (fs *MorphNameFrames) AnimateMaterial(
	frame float32,
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
