//go:build windows
// +build windows

package rendering

import (
	"github.com/go-gl/mathgl/mgl32"
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
	AspectRatio  float32
	NearPlane    float32
	FarPlane     float32
}

// NewDefaultCamera はデフォルト設定のカメラを作成する
func NewDefaultCamera(width, height int) *Camera {
	defaultCam := &Camera{
		Position:     &mmath.MVec3{X: 0.0, Y: InitialCameraPositionY, Z: InitialCameraPositionZ},
		LookAtCenter: &mmath.MVec3{X: 0.0, Y: InitialCameraPositionY, Z: 0.0},
		Up:           &mmath.MVec3{X: 0.0, Y: 1.0, Z: 0.0},
		FieldOfView:  FieldOfViewAngle,
		AspectRatio:  float32(width) / float32(height),
		NearPlane:    0.1,
		FarPlane:     1000.0,
	}
	return defaultCam
}

// UpdateAspectRatio はアスペクト比を更新する
func (c *Camera) UpdateAspectRatio(width, height int) {
	c.AspectRatio = float32(width) / float32(height)
}

// Reset はカメラをデフォルト設定にリセットする
func (c *Camera) Reset(width, height int) {
	defaultCam := NewDefaultCamera(width, height)
	c.Position = defaultCam.Position
	c.LookAtCenter = defaultCam.LookAtCenter
	c.Up = defaultCam.Up
	c.FieldOfView = defaultCam.FieldOfView
	c.AspectRatio = defaultCam.AspectRatio
}

// プロジェクション行列を取得
func (c *Camera) GetProjectionMatrix(width, height int) mgl32.Mat4 {
	return mgl32.Perspective(
		mgl32.DegToRad(c.FieldOfView),
		float32(width)/float32(height),
		c.NearPlane,
		c.FarPlane,
	)
}

// ビュー行列を取得
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	cameraPosition := mmath.NewGlVec3(c.Position)
	lookAtCenter := mmath.NewGlVec3(c.LookAtCenter)
	up := mmath.NewGlVec3(c.Up)
	return mgl32.LookAtV(cameraPosition, lookAtCenter, up)
}
