//go:build windows
// +build windows

package rendering

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// カメラ周りの各種初期値
const (
	InitialCameraPositionY float64 = 11.0
	InitialCameraPositionZ float64 = -40.0
	InitialLookAtCenterY   float64 = 11.0
	FieldOfViewAngle       float32 = 40.0
)

// Camera はカメラの位置と設定を保持する構造体
type Camera struct {
	Position     *mmath.MVec3
	LookAtCenter *mmath.MVec3
	Up           *mmath.MVec3
	FieldOfView  float32
	NearPlane    float32
	FarPlane     float32
}

// NewDefaultCamera はデフォルト設定のカメラを作成する
func NewDefaultCamera() *Camera {
	return &Camera{
		Position:     &mmath.MVec3{X: 0.0, Y: InitialCameraPositionY, Z: InitialCameraPositionZ},
		LookAtCenter: &mmath.MVec3{X: 0.0, Y: InitialCameraPositionY, Z: 0.0},
		Up:           &mmath.MVec3{X: 0.0, Y: 1.0, Z: 0.0},
		FieldOfView:  FieldOfViewAngle,
		NearPlane:    0.1,
		FarPlane:     1000.0,
	}
}

// Reset はカメラをデフォルト設定にリセットする
func (c *Camera) Reset() {
	defaultCam := NewDefaultCamera()
	c.Position = defaultCam.Position
	c.LookAtCenter = defaultCam.LookAtCenter
	c.Up = defaultCam.Up
	c.FieldOfView = defaultCam.FieldOfView
}
