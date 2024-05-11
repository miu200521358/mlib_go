package vmd

import (
	"slices"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/mcore"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type MorphNameFrames struct {
	*mcore.IndexModels[*MorphFrame]
	Name              string       // ボーン名
	RegisteredIndexes []int        // 登録対象キーフレリスト
	lock              sync.RWMutex // マップアクセス制御用
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		IndexModels:       mcore.NewIndexModels[*MorphFrame](),
		Name:              name,
		RegisteredIndexes: []int{},
		lock:              sync.RWMutex{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (mnfs *MorphNameFrames) GetRangeIndexes(index int) (int, int) {
	if len(mnfs.RegisteredIndexes) == 0 {
		return 0.0, 0.0
	}

	prevIndex := 0
	nextIndex := index

	if idx := mutils.SearchInts(mnfs.RegisteredIndexes, index); idx == 0 {
		prevIndex = 0.0
	} else {
		prevIndex = mnfs.RegisteredIndexes[idx-1]
	}

	if idx := mutils.SearchInts(mnfs.RegisteredIndexes, index); idx >= len(mnfs.RegisteredIndexes) {
		nextIndex = slices.Max(mnfs.RegisteredIndexes)
	} else {
		nextIndex = mnfs.RegisteredIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (mnfs *MorphNameFrames) GetItem(index int) *MorphFrame {
	if mnfs == nil {
		return NewMorphFrame(index)
	}

	mnfs.lock.RLock()
	defer mnfs.lock.RUnlock()

	if slices.Contains(mnfs.RegisteredIndexes, index) {
		return mnfs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := mnfs.GetRangeIndexes(index)

	if prevIndex == nextIndex {
		if slices.Contains(mnfs.RegisteredIndexes, nextIndex) {
			nextMf := mnfs.Data[nextIndex]
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
	if slices.Contains(mnfs.RegisteredIndexes, prevIndex) {
		prevMf = mnfs.Data[prevIndex]
	} else {
		prevMf = NewMorphFrame(index)
	}
	if slices.Contains(mnfs.RegisteredIndexes, nextIndex) {
		nextMf = mnfs.Data[nextIndex]
	} else {
		nextMf = NewMorphFrame(index)
	}

	mf := NewMorphFrame(index)

	ry := (index - prevIndex) / (nextIndex - prevIndex)
	mf.Ratio = prevMf.Ratio + (nextMf.Ratio-prevMf.Ratio)*float64(ry)

	return mf
}

// bf.Registered が true の場合、補間曲線を分割して登録する
func (mnfs *MorphNameFrames) Append(value *MorphFrame) {
	mnfs.lock.Lock()
	defer mnfs.lock.Unlock()

	if !slices.Contains(mnfs.RegisteredIndexes, value.Index) {
		mnfs.RegisteredIndexes = append(mnfs.RegisteredIndexes, value.Index)
		mutils.SortInts(mnfs.RegisteredIndexes)
	}

	if value.Registered && !slices.Contains(mnfs.RegisteredIndexes, value.Index) {
		mnfs.RegisteredIndexes = append(mnfs.RegisteredIndexes, value.Index)
		mutils.SortInts(mnfs.RegisteredIndexes)
	}

	mnfs.Data[value.Index] = value
}

func (mnfs *MorphNameFrames) GetMaxFrame() int {
	if len(mnfs.RegisteredIndexes) == 0 {
		return 0
	}

	return slices.Max(mnfs.RegisteredIndexes)
}

func (mnfs *MorphNameFrames) GetMinFrame() int {
	if len(mnfs.RegisteredIndexes) == 0 {
		return 0
	}

	return slices.Min(mnfs.RegisteredIndexes)
}

func (mnfs *MorphNameFrames) AnimateVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := mnfs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(mnfs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			delta.Position.Add(offset.Position.MuledScalar(mf.Ratio))
		}
	}
}

func (mnfs *MorphNameFrames) AnimateAfterVertex(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := mnfs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(mnfs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.VertexMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			delta.AfterPosition.Add(offset.Position.MuledScalar(mf.Ratio))
		}
	}
}

func (mnfs *MorphNameFrames) AnimateUv(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := mnfs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(mnfs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			uv := offset.Uv.MuledScalar(mf.Ratio).GetXY()
			delta.Uv.Add(uv)
		}
	}
}

func (mnfs *MorphNameFrames) AnimateUv1(
	frame int,
	model *pmx.PmxModel,
	deltas *VertexMorphDeltas,
) {
	mf := mnfs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(mnfs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.UvMorphOffset)
		if 0 < offset.VertexIndex && offset.VertexIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.VertexIndex]
			uv := offset.Uv.MuledScalar(mf.Ratio)
			delta.Uv1.Add(uv.GetXY())
		}
	}
}

func (mnfs *MorphNameFrames) AnimateBone(
	frame int,
	model *pmx.PmxModel,
	deltas *BoneMorphDeltas,
) {
	mf := mnfs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(mnfs.Name)
	for _, o := range morph.Offsets {
		offset := o.(*pmx.BoneMorphOffset)
		if 0 < offset.BoneIndex && offset.BoneIndex <= len(deltas.Data) {
			delta := deltas.Data[offset.BoneIndex]
			delta.MorphPosition.Add(offset.Position.MuledScalar(mf.Ratio))
			delta.MorphLocalPosition.Add(offset.LocalPosition.MuledScalar(mf.Ratio))
			deltaRad := offset.Rotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphRotation.SetQuaternion(delta.MorphRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromRadians(deltaRad.GetX(), deltaRad.GetY(), deltaRad.GetZ())))
			deltaLocalRad := offset.LocalRotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphLocalRotation.SetQuaternion(delta.MorphLocalRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromRadians(deltaLocalRad.GetX(), deltaLocalRad.GetY(), deltaLocalRad.GetZ())))
			delta.MorphScale.Add(offset.Scale.MuledScalar(mf.Ratio))
			delta.MorphLocalScale.Add(offset.LocalScale.MuledScalar(mf.Ratio))
		}
	}
}

// AnimateMaterial 材質モーフの適用
func (mnfs *MorphNameFrames) AnimateMaterial(
	frame int,
	model *pmx.PmxModel,
	deltas *MaterialMorphDeltas,
) {
	mf := mnfs.GetItem(frame)
	if mf.Ratio == 0.0 {
		return
	}

	morph := model.Morphs.GetItemByName(mnfs.Name)
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
