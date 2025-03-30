package controller

import "github.com/miu200521358/walk/pkg/walk"

type MProgressBar struct {
	*walk.ProgressBar
}

func NewMProgressBar(parent walk.Container) (*MProgressBar, error) {
	pb := new(MProgressBar)

	var err error
	pb.ProgressBar, err = walk.NewProgressBar(parent)
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func (pb *MProgressBar) SetMax(max int) {
	pb.SetRange(0, max)
}

func (pb *MProgressBar) Increment() {
	pb.SetValue(pb.Value() + 1)
	if pb.Value() >= pb.MaxValue() {
		pb.SetValue(0)
		pb.SetMax(0)
	}
}

func (pb *MProgressBar) Add(value int) {
	pb.SetValue(pb.Value() + value)
	if pb.Value() >= pb.MaxValue() {
		pb.SetValue(0)
		pb.SetMax(0)
	}
}
