package vmd

type CameraFrames struct {
	*BaseFrames[*CameraFrame]
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		BaseFrames: NewBaseFrames[*CameraFrame](NewCameraFrame, NullCameraFrame),
	}
}

func (i *CameraFrames) NewFrame(index int) *CameraFrame {
	return NewCameraFrame(index)
}
