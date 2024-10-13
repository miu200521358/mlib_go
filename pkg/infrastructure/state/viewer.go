package state

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

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
	if !mmath.NearEquals(vp.Yaw, other.Yaw, 1e-2) {
		return false
	}
	if !mmath.NearEquals(vp.Pitch, other.Pitch, 1e-2) {
		return false
	}
	if !mmath.NearEquals(float64(vp.FieldOfViewAngle), float64(other.FieldOfViewAngle), 1e-2) {
		return false
	}
	if !vp.Size.NearEquals(other.Size, 1e-2) {
		return false
	}
	if !vp.CameraPos.NearEquals(other.CameraPos, 1e-2) {
		return false
	}
	if !vp.CameraUp.NearEquals(other.CameraUp, 1e-2) {
		return false
	}
	if !vp.LookAtCenter.NearEquals(other.LookAtCenter, 1e-2) {
		return false
	}
	return true
}

func (vp *ViewerParameter) Copy() *ViewerParameter {
	return &ViewerParameter{
		Yaw:              vp.Yaw,
		Pitch:            vp.Pitch,
		FieldOfViewAngle: vp.FieldOfViewAngle,
		Size:             vp.Size.Copy(),
		CameraPos:        vp.CameraPos.Copy(),
		CameraUp:         vp.CameraUp.Copy(),
		LookAtCenter:     vp.LookAtCenter.Copy(),
	}
}

func (vp *ViewerParameter) String() string {
	return fmt.Sprintf("Yaw: %.8f, Pitch: %.8f, FieldOfViewAngle: %.8f, Size: %s, CameraPos: %s, CameraUp: %s, LookAtCenter: %s", vp.Yaw, vp.Pitch, vp.FieldOfViewAngle, vp.Size, vp.CameraPos, vp.CameraUp, vp.LookAtCenter)
}
