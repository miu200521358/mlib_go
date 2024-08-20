package state

import "github.com/miu200521358/mlib_go/pkg/domain/mmath"

type ViewerParameter struct {
	Yaw              float64
	Pitch            float64
	FieldOfViewAngle float32
	Size             *mmath.MVec2
	CameraPos        *mmath.MVec3
	CameraUp         *mmath.MVec3
	LookAtCenter     *mmath.MVec3
}

func (vp *ViewerParameter) Equals(other *ViewerParameter) bool {
	if vp.Yaw != other.Yaw {
		return false
	}
	if vp.Pitch != other.Pitch {
		return false
	}
	if vp.FieldOfViewAngle != other.FieldOfViewAngle {
		return false
	}
	if !vp.Size.Equals(other.Size) {
		return false
	}
	if !vp.CameraPos.Equals(other.CameraPos) {
		return false
	}
	if !vp.CameraUp.Equals(other.CameraUp) {
		return false
	}
	if !vp.LookAtCenter.Equals(other.LookAtCenter) {
		return false
	}
	return true
}
