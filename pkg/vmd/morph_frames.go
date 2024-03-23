package vmd

import "github.com/miu200521358/mlib_go/pkg/pmx"

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

func (mfs *MorphFrames) Append(bnfs *MorphNameFrames) {
	mfs.Data[bnfs.Name] = bnfs
}

func (mfs *MorphFrames) GetItem(morphName string) *MorphNameFrames {
	if !mfs.Contains(morphName) {
		mfs.Append(NewMorphNameFrames(morphName))
	}
	return mfs.Data[morphName]
}

func (mfs *MorphFrames) Animate(
	frame float32,
	model *pmx.PmxModel,
	morphNames []string,
	isOutLog bool,
	description string,
) *MorphDeltas {
	mds := NewMorphDeltas(len(model.Vertices.Data))
	for _, morphName := range morphNames {
		if !mfs.Contains(morphName) || !model.Morphs.ContainsByName(morphName) {
			continue
		}

		morph := model.Morphs.GetItemByName(morphName)
		if morph.MorphType == pmx.MORPH_TYPE_VERTEX {
			mfs.GetItem(morphName).AnimateVertex(frame, model, mds.Vertices)
		}
	}

	return mds
}
