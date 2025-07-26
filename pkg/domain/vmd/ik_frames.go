package vmd

import "github.com/tiendc/go-deepcopy"

type IkFrames struct {
	*BaseFrames[*IkFrame]
}

func NewIkFrames() *IkFrames {
	return &IkFrames{
		BaseFrames: NewBaseFrames[*IkFrame](NewIkFrame, nilIkFrame),
	}
}

func nilIkFrame() *IkFrame {
	return nil
}

func (ikFrames *IkFrames) Clean() {
	if ikFrames.Length() > 1 {
		return
	} else {
		cf := ikFrames.Get(ikFrames.values.Min())
		if !(cf.IkList == nil || len(cf.IkList) == 0) {
			return
		}
		ikFrames.Delete(cf.Index())
	}
}

func (ikFrames *IkFrames) Copy() (*IkFrames, error) {
	copied := new(IkFrames)
	err := deepcopy.Copy(copied, ikFrames)
	return copied, err
}
