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

func (cameraFrames *CameraFrames) Clean() {
	if cameraFrames.Len() > 1 {
		return
	} else {
		cf := cameraFrames.Get(cameraFrames.Indexes.Min())
		if !(cf.Position == nil || cf.Position.Length() == 0 ||
			cf.Rotation == nil || cf.Rotation.Degrees().Length() == 0 ||
			cf.Distance == 0 || cf.ViewOfAngle == 0 || cf.IsPerspectiveOff) {
			return
		}
		cameraFrames.Delete(cf.Index())
	}
}
