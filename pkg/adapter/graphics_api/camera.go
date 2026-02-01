// 指示: miu200521358
package graphics_api

import (
	"fmt"
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"gonum.org/v1/gonum/spatial/r3"
)

const (
	// InitialCameraPositionY は初期カメラ位置のY成分。
	InitialCameraPositionY float64 = 11.0
	// InitialCameraPositionZ は初期カメラ位置のZ成分。
	InitialCameraPositionZ float64 = -40.0
	// InitialLookAtCenterY は注視点の初期Y成分。
	InitialLookAtCenterY float64 = 11.0
	// FieldOfViewAngle は視野角（度）。
	FieldOfViewAngle float32 = 40.0
)

// Camera はカメラの位置と設定を保持する。
type Camera struct {
	Position     *mmath.Vec3
	LookAtCenter *mmath.Vec3
	Up           *mmath.Vec3
	FieldOfView  float32
	AspectRatio  float32
	NearPlane    float32
	FarPlane     float32
	Yaw          float64
	Pitch        float64
}

// NewDefaultCamera は既定値のカメラを生成する。
func NewDefaultCamera(width, height int) *Camera {
	return &Camera{
		Position:     &mmath.Vec3{Vec: r3.Vec{X: 0.0, Y: InitialCameraPositionY, Z: InitialCameraPositionZ}},
		LookAtCenter: &mmath.Vec3{Vec: r3.Vec{X: 0.0, Y: InitialLookAtCenterY, Z: 0.0}},
		Up:           &mmath.Vec3{Vec: r3.Vec{X: 0.0, Y: 1.0, Z: 0.0}},
		FieldOfView:  FieldOfViewAngle,
		AspectRatio:  float32(width) / float32(height),
		NearPlane:    0.1,
		FarPlane:     1000.0,
	}
}

// String はカメラ情報を文字列化する。
func (c *Camera) String() string {
	return fmt.Sprintf("Camera: Position: %v, LookAtCenter: %v, Up: %v, FieldOfView: %.5f, AspectRatio: %.5f, NearPlane: %.5f, FarPlane: %.5f, Yaw: %.5f, Pitch: %.5f",
		c.Position, c.LookAtCenter, c.Up, c.FieldOfView, c.AspectRatio, c.NearPlane, c.FarPlane, c.Yaw, c.Pitch)
}

// UpdateAspectRatio はアスペクト比を更新する。
func (c *Camera) UpdateAspectRatio(width, height int) {
	if c == nil || height == 0 {
		return
	}
	c.AspectRatio = float32(width) / float32(height)
}

// Reset はカメラを既定値に戻す。
func (c *Camera) Reset(width, height int) {
	if c == nil {
		return
	}
	defaultCam := NewDefaultCamera(width, height)
	c.Position = defaultCam.Position
	c.LookAtCenter = defaultCam.LookAtCenter
	c.Up = defaultCam.Up
	c.FieldOfView = defaultCam.FieldOfView
	c.AspectRatio = defaultCam.AspectRatio
	c.NearPlane = defaultCam.NearPlane
	c.FarPlane = defaultCam.FarPlane
}

// GetProjectionMatrix は射影行列を返す。
func (c *Camera) GetProjectionMatrix(width, height int) mmath.Mat4 {
	if c == nil {
		return mmath.NewMat4()
	}
	if height == 0 {
		return mmath.NewMat4()
	}
	aspect := float64(width) / float64(height)
	if aspect == 0 {
		return mmath.NewMat4()
	}
	fovRad := mmath.DegToRad(float64(c.FieldOfView))
	f := 1.0 / math.Tan(fovRad*0.5)
	near := float64(c.NearPlane)
	far := float64(c.FarPlane)
	denom := near - far
	if denom == 0 {
		return mmath.NewMat4()
	}
	return mmath.NewMat4ByValues(
		f/aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far+near)/denom, (2*far*near)/denom,
		0, 0, -1, 0,
	)
}

// GetViewMatrix はビュー行列を返す。
func (c *Camera) GetViewMatrix() mmath.Mat4 {
	if c == nil || c.Position == nil || c.LookAtCenter == nil || c.Up == nil {
		return mmath.NewMat4()
	}
	return mmath.NewMat4FromLookAt(*c.Position, *c.LookAtCenter, *c.Up)
}

// ResetPosition はカメラの位置を角度指定で更新する。
func (c *Camera) ResetPosition(yaw, pitch float64) {
	if c == nil {
		return
	}
	c.Yaw = yaw
	c.Pitch = pitch

	radius := math.Abs(InitialCameraPositionZ)
	yawRad := mmath.DegToRad(yaw)
	pitchRad := mmath.DegToRad(pitch)
	orientation := mmath.NewQuaternionFromAxisAngles(mmath.UNIT_Y_VEC3, yawRad).
		Muled(mmath.NewQuaternionFromAxisAngles(mmath.UNIT_X_VEC3, pitchRad))
	forwardXYZ := orientation.MulVec3(mmath.UNIT_Z_NEG_VEC3).MuledScalar(radius)

	if c.Position == nil {
		c.Position = &mmath.Vec3{}
	}
	c.Position.X = forwardXYZ.X
	c.Position.Y = InitialCameraPositionY + forwardXYZ.Y
	c.Position.Z = forwardXYZ.Z
}
