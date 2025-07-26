package vmd

import "github.com/tiendc/go-deepcopy"

type CameraFrames struct {
	*BaseFrames[*CameraFrame]
}

func NewCameraFrames() *CameraFrames {
	return &CameraFrames{
		BaseFrames: NewBaseFrames[*CameraFrame](NewCameraFrame, nilCameraFrame),
	}
}

func nilCameraFrame() *CameraFrame {
	return nil
}

func (cameraFrames *CameraFrames) Clean() {
	if cameraFrames.Length() > 1 {
		return
	} else {
		cf := cameraFrames.Get(cameraFrames.values.Min())
		if !(cf.Position == nil || cf.Position.Length() == 0 ||
			cf.Degrees == nil || cf.Degrees.Length() == 0 ||
			cf.Distance == 0 || cf.ViewOfAngle == 0 || cf.IsPerspectiveOff) {
			return
		}
		cameraFrames.Delete(cf.Index())
	}
}

func (cameraFrames *CameraFrames) Copy() (*CameraFrames, error) {
	copied := new(CameraFrames)
	err := deepcopy.Copy(copied, cameraFrames)
	return copied, err
}
