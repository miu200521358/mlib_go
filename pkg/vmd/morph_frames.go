package vmd

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type MorphFrames struct {
	Data map[string]*MorphNameFrames
}

func NewMorphFrames() *MorphFrames {
	return &MorphFrames{
		Data: make(map[string]*MorphNameFrames, 0),
	}
}

func (mfs *MorphFrames) Contains(morphName string) bool {
	_, ok := mfs.Data[morphName]
	return ok
}

func (mfs *MorphFrames) Append(fs *MorphNameFrames) {
	mfs.Data[fs.Name] = fs
}

func (mfs *MorphFrames) GetItem(morphName string) *MorphNameFrames {
	if !mfs.Contains(morphName) {
		mfs.Append(NewMorphNameFrames(morphName))
	}
	return mfs.Data[morphName]
}

func (mfs *MorphFrames) Animate(
	frame int,
	model *pmx.PmxModel,
	morphNames []string,
) *MorphDeltas {
	mds := NewMorphDeltas(len(model.Vertices.Data), len(model.Bones.Data), model.Materials)
	for _, morphName := range morphNames {
		if !mfs.Contains(morphName) || !model.Morphs.ContainsByName(morphName) {
			continue
		}

		morph := model.Morphs.GetItemByName(morphName)
		if morph.MorphType == pmx.MORPH_TYPE_VERTEX {
			mfs.GetItem(morphName).AnimateVertex(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_AFTER_VERTEX {
			mfs.GetItem(morphName).AnimateAfterVertex(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_UV {
			mfs.GetItem(morphName).AnimateUv(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_EXTENDED_UV1 {
			mfs.GetItem(morphName).AnimateUv1(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_BONE {
			mfs.GetItem(morphName).AnimateBone(frame, model, mds.Bones)
		} else if morph.MorphType == pmx.MORPH_TYPE_MATERIAL {
			mfs.GetItem(morphName).AnimateMaterial(frame, model, mds.Materials)
		}
	}

	return mds
}

func (mfs *MorphFrames) GetMaxFrame() int {
	maxFno := int(0)
	for _, mnfs := range mfs.Data {
		fno := mnfs.GetMaxFrame()
		if fno > maxFno {
			maxFno = fno
		}
	}
	return maxFno
}

func (mfs *MorphFrames) GetMinFrame() int {
	minFno := math.MaxInt
	for _, mnfs := range mfs.Data {
		fno := mnfs.GetMinFrame()
		if fno < minFno {
			minFno = fno
		}
	}
	return minFno
}

func (fs *MorphFrames) GetCount() int {
	count := 0
	for _, fs := range fs.Data {
		count += len(fs.RegisteredIndexes)
	}
	return count
}
