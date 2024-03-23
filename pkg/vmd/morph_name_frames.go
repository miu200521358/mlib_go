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
	*mcore.IndexFloatModelCorrection[*MorphFrame]
	Name              string       // ボーン名
	RegisteredIndexes []float32    // 登録対象キーフレリスト
	lock              sync.RWMutex // マップアクセス制御用
}

func NewMorphNameFrames(name string) *MorphNameFrames {
	return &MorphNameFrames{
		IndexFloatModelCorrection: mcore.NewIndexFloatModelCorrection[*MorphFrame](),
		Name:                      name,
		RegisteredIndexes:         []float32{},
	}
}

// 指定したキーフレの前後のキーフレ番号を返す
func (mnfs *MorphNameFrames) GetRangeIndexes(index float32) (float32, float32) {
	if len(mnfs.RegisteredIndexes) == 0 {
		return 0.0, 0.0
	}

	prevIndex := float32(0.0)
	nextIndex := index

	if idx := mutils.SearchFloat32s(mnfs.RegisteredIndexes, index); idx == 0 {
		prevIndex = 0.0
	} else {
		prevIndex = mnfs.RegisteredIndexes[idx-1]
	}

	if idx := mutils.SearchFloat32s(mnfs.RegisteredIndexes, index); idx == len(mnfs.RegisteredIndexes) {
		nextIndex = slices.Max(mnfs.RegisteredIndexes)
	} else {
		nextIndex = mnfs.RegisteredIndexes[idx]
	}

	return prevIndex, nextIndex
}

// キーフレ計算結果を返す
func (mnfs *MorphNameFrames) GetItem(index float32) *MorphFrame {
	if mnfs == nil {
		return NewMorphFrame(index)
	}

	mnfs.lock.RLock()
	defer mnfs.lock.RUnlock()

	if slices.Contains(mnfs.Indexes, index) {
		return mnfs.Data[index]
	}

	// なかったら補間計算して返す
	prevIndex, nextIndex := mnfs.GetRangeIndexes(index)

	if prevIndex == nextIndex {
		if slices.Contains(mnfs.Indexes, nextIndex) {
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
	if slices.Contains(mnfs.Indexes, prevIndex) {
		prevMf = mnfs.Data[prevIndex]
	} else {
		prevMf = NewMorphFrame(index)
	}
	if slices.Contains(mnfs.Indexes, nextIndex) {
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

	if !slices.Contains(mnfs.Indexes, value.Index) {
		mnfs.Indexes = append(mnfs.Indexes, value.Index)
		mutils.SortFloat32s(mnfs.Indexes)
	}

	if value.Registered && !slices.Contains(mnfs.RegisteredIndexes, value.Index) {
		mnfs.RegisteredIndexes = append(mnfs.RegisteredIndexes, value.Index)
		mutils.SortFloat32s(mnfs.RegisteredIndexes)
	}

	mnfs.Data[value.Index] = value
}

func (mnfs *MorphNameFrames) GetMaxFrame() float32 {
	if len(mnfs.RegisteredIndexes) == 0 {
		return 0
	}

	return slices.Max(mnfs.RegisteredIndexes)
}

func (mnfs *MorphNameFrames) AnimateVertex(
	frame float32,
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
	frame float32,
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
	frame float32,
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
			// UVは左上原点なので、Yを反転させる
			delta.Uv.Add(&mmath.MVec2{uv[0], 1 - uv[1]})
		}
	}
}

func (mnfs *MorphNameFrames) AnimateUv1(
	frame float32,
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
			// UVは左上原点なので、Yを反転させる
			delta.Uv1.Add(&mmath.MVec2{uv[0], 1 - uv[1]})
		}
	}
}

func (mnfs *MorphNameFrames) AnimateBone(
	frame float32,
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
				mmath.NewMQuaternionFromEulerAngles(deltaRad.GetX(), deltaRad.GetY(), deltaRad.GetZ())))
			deltaLocalRad := offset.LocalRotation.GetRadians().MuledScalar(mf.Ratio)
			delta.MorphLocalRotation.SetQuaternion(delta.MorphLocalRotation.GetQuaternion().Muled(
				mmath.NewMQuaternionFromEulerAngles(deltaLocalRad.GetX(), deltaLocalRad.GetY(), deltaLocalRad.GetZ())))
			delta.MorphScale.Add(offset.Scale.MuledScalar(mf.Ratio))
			delta.MorphLocalScale.Add(offset.LocalScale.MuledScalar(mf.Ratio))
		}
	}
}
