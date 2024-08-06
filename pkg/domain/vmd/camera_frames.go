package vmd

type CameraFrames struct {
	*BaseFrames[*CameraFrame]
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		BaseFrames: NewBaseFrames[*CameraFrame](NewCameraFrame, NullCameraFrame),
	}
}
