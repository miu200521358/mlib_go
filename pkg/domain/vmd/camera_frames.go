package vmd

type CameraFrames struct {
	*BaseFrames[*CameraFrame]
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		BaseFrames: NewBaseFrames[*CameraFrame](NewCameraFrame, NullCameraFrame),
	}
}

func (cameraFrames *CameraFrames) Copy() *CameraFrames {
	copied := NewCameraFrames()
	for _, frame := range cameraFrames.List() {
		copied.Append(frame.Copy().(*CameraFrame))
	}
	return copied
}
