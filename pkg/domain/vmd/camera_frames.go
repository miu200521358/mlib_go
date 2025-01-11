package vmd

type CameraFrames struct {
	*BaseFrames[*CameraFrame]
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		BaseFrames: NewBaseFrames[*CameraFrame](),
	}
}

func (cameraFrames *CameraFrames) Clean() {
	if cameraFrames.Length() > 1 {
		return
	} else {
		cf := cameraFrames.Get(cameraFrames.Indexes.Min())
		if !(cf.Position == nil || cf.Position.Length() == 0 ||
			cf.Radians == nil || cf.Radians.Length() == 0 ||
			cf.Distance == 0 || cf.ViewOfAngle == 0 || cf.IsPerspectiveOff) {
			return
		}
		cameraFrames.Delete(cf.Index())
	}
}
