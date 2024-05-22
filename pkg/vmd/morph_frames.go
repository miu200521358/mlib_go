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

func (mfs *MorphFrames) Get(morphName string) *MorphNameFrames {
	if !mfs.Contains(morphName) {
		mfs.Append(NewMorphNameFrames(morphName))
	}
	return mfs.Data[morphName]
}

func (mfs *MorphFrames) Deform(
	frame int,
	model *pmx.PmxModel,
	morphNames []string,
) *MorphDeltas {
	mds := NewMorphDeltas(len(model.Vertices.Data), len(model.Bones.Data), model.Materials)
	for _, morphName := range morphNames {
		if !mfs.Contains(morphName) || !model.Morphs.ContainsByName(morphName) {
			continue
		}

		morph := model.Morphs.GetByName(morphName)
		if morph.MorphType == pmx.MORPH_TYPE_VERTEX {
			mfs.Get(morphName).DeformVertex(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_AFTER_VERTEX {
			mfs.Get(morphName).DeformAfterVertex(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_UV {
			mfs.Get(morphName).DeformUv(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_EXTENDED_UV1 {
			mfs.Get(morphName).DeformUv1(frame, model, mds.Vertices)
		} else if morph.MorphType == pmx.MORPH_TYPE_BONE {
			mfs.Get(morphName).DeformBone(frame, model, mds.Bones)
		} else if morph.MorphType == pmx.MORPH_TYPE_MATERIAL {
			mfs.Get(morphName).DeformMaterial(frame, model, mds.Materials)
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
		count += fs.RegisteredIndexes.Len()
	}
	return count
}
