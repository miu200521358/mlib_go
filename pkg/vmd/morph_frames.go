package vmd

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
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

func (mfs *MorphFrames) Update(fs *MorphNameFrames) {
	mfs.Data[fs.Name] = fs
}

func (mfs *MorphFrames) GetNames() []string {
	names := make([]string, 0, len(mfs.Data))
	for name := range mfs.Data {
		names = append(names, name)
	}
	return names
}

func (mfs *MorphFrames) Get(morphName string) *MorphNameFrames {
	if !mfs.Contains(morphName) {
		mfs.Update(NewMorphNameFrames(morphName))
	}
	return mfs.Data[morphName]
}

func (mfs *MorphFrames) Deform(
	frame int,
	model *pmx.PmxModel,
	morphNames []string,
) *MorphDeltas {
	mds := NewMorphDeltas(model.Materials, model.Bones)
	for _, morphName := range morphNames {
		if !mfs.Contains(morphName) || !model.Morphs.ContainsByName(morphName) {
			continue
		}

		mf := mfs.Get(morphName).Get(frame)
		if mf == nil {
			continue
		}

		morph := model.Morphs.GetByName(morphName)
		switch morph.MorphType {
		case pmx.MORPH_TYPE_VERTEX:
			mds.Vertices = mf.DeformVertex(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_AFTER_VERTEX:
			mds.Vertices = mf.DeformAfterVertex(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_UV:
			mds.Vertices = mf.DeformUv(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_EXTENDED_UV1:
			mds.Vertices = mf.DeformUv1(morphName, model, mds.Vertices, mf.Ratio)
		case pmx.MORPH_TYPE_BONE:
			mds.Bones = mf.DeformBone(morphName, model, mds.Bones, mf.Ratio)
		case pmx.MORPH_TYPE_MATERIAL:
			mds.Materials = mf.DeformMaterial(morphName, model, mds.Materials, mf.Ratio)
		case pmx.MORPH_TYPE_GROUP:
			// グループモーフは細分化
			for _, offset := range morph.Offsets {
				groupOffset := offset.(*pmx.GroupMorphOffset)
				groupMorph := model.Morphs.Get(groupOffset.MorphIndex)
				if groupMorph == nil {
					continue
				}
				gmf := mfs.Get(groupMorph.Name).Get(frame)
				switch groupMorph.MorphType {
				case pmx.MORPH_TYPE_VERTEX:
					mds.Vertices = gmf.DeformVertex(
						groupMorph.Name, model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_AFTER_VERTEX:
					mds.Vertices = gmf.DeformAfterVertex(
						groupMorph.Name, model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_UV:
					mds.Vertices = gmf.DeformUv(
						groupMorph.Name, model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_EXTENDED_UV1:
					mds.Vertices = gmf.DeformUv1(
						groupMorph.Name, model, mds.Vertices, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_BONE:
					mds.Bones = gmf.DeformBone(
						groupMorph.Name, model, mds.Bones, mf.Ratio*groupOffset.MorphFactor)
				case pmx.MORPH_TYPE_MATERIAL:
					mds.Materials = gmf.DeformMaterial(
						groupMorph.Name, model, mds.Materials, mf.Ratio*groupOffset.MorphFactor)
				}
			}
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

func (fs *MorphFrames) Len() int {
	count := 0
	for _, fs := range fs.Data {
		count += fs.RegisteredIndexes.Len()
	}
	return count
}
