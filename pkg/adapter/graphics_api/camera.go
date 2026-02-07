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
	// InitialLookAtCenterX は注視点の初期X成分。
	InitialLookAtCenterX float64 = 0.0
	// InitialLookAtCenterY は注視点の初期Y成分。
	InitialLookAtCenterY float64 = 11.0
	// InitialLookAtCenterZ は注視点の初期Z成分。
	InitialLookAtCenterZ float64 = 0.0
	// FieldOfViewAngle は視野角（度）。
	FieldOfViewAngle float32 = 40.0
)

// Camera はカメラの位置と設定を保持する。
type Camera struct {
	Position     *mmath.Vec3
	LookAtCenter *mmath.Vec3
	Up           *mmath.Vec3
	Orientation  mmath.Quaternion
	FieldOfView  float32
	AspectRatio  float32
	NearPlane    float32
	FarPlane     float32
}

// NewDefaultCamera は既定値のカメラを生成する。
func NewDefaultCamera(width, height int) *Camera {
	cam := &Camera{
		Position:     &mmath.Vec3{},
		LookAtCenter: &mmath.Vec3{Vec: r3.Vec{X: InitialLookAtCenterX, Y: InitialLookAtCenterY, Z: InitialLookAtCenterZ}},
		Up:           &mmath.Vec3{Vec: r3.Vec{X: 0.0, Y: 1.0, Z: 0.0}},
		Orientation:  mmath.NewQuaternion(),
		FieldOfView:  FieldOfViewAngle,
		AspectRatio:  float32(width) / float32(height),
		NearPlane:    0.1,
		FarPlane:     1000.0,
	}
	cam.applyOrientation(math.Abs(InitialCameraPositionZ))
	return cam
}

// String はカメラ情報を文字列化する。
func (c *Camera) String() string {
	return fmt.Sprintf(
		"Camera: Position: %v, LookAtCenter: %v, Up: %v, Orientation: %v, FieldOfView: %.5f, AspectRatio: %.5f, NearPlane: %.5f, FarPlane: %.5f",
		c.Position, c.LookAtCenter, c.Up, c.Orientation, c.FieldOfView, c.AspectRatio, c.NearPlane, c.FarPlane,
	)
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
	c.Orientation = defaultCam.Orientation
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
	c.ensureState()
	radius := c.OrbitDistance()
	yawRad := mmath.DegToRad(yaw)
	pitchRad := mmath.DegToRad(pitch)
	c.Orientation = mmath.NewQuaternionFromAxisAngles(mmath.UNIT_Y_VEC3, yawRad).
		Muled(mmath.NewQuaternionFromAxisAngles(mmath.UNIT_X_VEC3, pitchRad)).
		Normalized()
	c.applyOrientation(radius)
}

// ResetPresetPosition はプリセット視点用に水平中心を初期値へ戻してからカメラ位置を更新する。
func (c *Camera) ResetPresetPosition(yaw, pitch float64) {
	if c == nil {
		return
	}
	c.ensureState()
	radius := c.OrbitDistance()
	// プリセット視点では水平軸のズレを持ち越さず、真正面/真上/真下を原点基準で揃える。
	c.LookAtCenter.X = InitialLookAtCenterX
	c.LookAtCenter.Z = InitialLookAtCenterZ
	yawRad := mmath.DegToRad(yaw)
	pitchRad := mmath.DegToRad(pitch)
	c.Orientation = mmath.NewQuaternionFromAxisAngles(mmath.UNIT_Y_VEC3, yawRad).
		Muled(mmath.NewQuaternionFromAxisAngles(mmath.UNIT_X_VEC3, pitchRad)).
		Normalized()
	c.applyOrientation(radius)
}

// SetByMotionValues はVMDカメラ値からカメラ状態を設定する。
func (c *Camera) SetByMotionValues(center, degrees mmath.Vec3, distance float64, viewOfAngle int) {
	if c == nil {
		return
	}
	c.ensureState()

	// 参照実装に合わせ、MMDカメラ角度から姿勢クォータニオンを生成する。
	orientationFromMotion := NewMmdCameraQuaternionFromDegrees(degrees)
	eye := center.Added(orientationFromMotion.MulVec3(mmath.UNIT_Z_VEC3).MuledScalar(distance))
	up := orientationFromMotion.MulVec3(mmath.UNIT_Y_VEC3).Normalized()
	if up.LengthSqr() == 0 {
		up = mmath.UNIT_Y_VEC3
	}

	c.LookAtCenter.X = center.X
	c.LookAtCenter.Y = center.Y
	c.LookAtCenter.Z = center.Z
	c.Position.X = eye.X
	c.Position.Y = eye.Y
	c.Position.Z = eye.Z
	c.Up.X = up.X
	c.Up.Y = up.Y
	c.Up.Z = up.Z
	if viewOfAngle > 0 {
		c.FieldOfView = float32(viewOfAngle)
	}

	forward := center.Subed(eye).Normalized()
	if forward.LengthSqr() == 0 {
		c.Orientation = orientationFromMotion
		return
	}
	c.Orientation = mmath.NewQuaternionFromDirection(forward, up).Normalized()
}

// NewMmdCameraQuaternionFromDegrees はMMDカメラ角度（度）から姿勢クォータニオンを生成する。
func NewMmdCameraQuaternionFromDegrees(degrees mmath.Vec3) mmath.Quaternion {
	base := mmath.NewQuaternionFromDegrees(-degrees.X, degrees.Y, degrees.Z)
	return mmath.NewQuaternionByValues(-base.X(), base.Y(), base.Z(), -base.W()).Normalized()
}

// RotateOrbit は軌道回転を加算してカメラ姿勢を更新する。
func (c *Camera) RotateOrbit(yawDelta, pitchDelta float64) {
	if c == nil {
		return
	}
	c.ensureState()
	radius := c.OrbitDistance()

	// ワールドYのヨー回転と、回転後ローカルX軸のピッチ回転を順に適用する。
	yawQuat := mmath.NewQuaternionFromAxisAngles(mmath.UNIT_Y_VEC3, mmath.DegToRad(yawDelta))
	rotated := yawQuat.Muled(c.Orientation).Normalized()
	rightAxis := rotated.MulVec3(mmath.UNIT_X_VEC3).Normalized()
	if rightAxis.LengthSqr() == 0 {
		rightAxis = mmath.UNIT_X_VEC3
	}
	pitchQuat := mmath.NewQuaternionFromAxisAngles(rightAxis, mmath.DegToRad(pitchDelta))
	c.Orientation = pitchQuat.Muled(rotated).Normalized()
	c.applyOrientation(radius)
}

// OrbitDistance は注視点からカメラまでの距離を返す。
func (c *Camera) OrbitDistance() float64 {
	if c == nil || c.Position == nil || c.LookAtCenter == nil {
		return math.Abs(InitialCameraPositionZ)
	}
	distance := c.Position.Subed(*c.LookAtCenter).Length()
	if distance <= 1e-8 {
		return math.Abs(InitialCameraPositionZ)
	}
	return distance
}

// RightVector はカメラ右方向ベクトルを返す。
func (c *Camera) RightVector() mmath.Vec3 {
	if c == nil {
		return mmath.UNIT_X_VEC3
	}
	c.ensureState()
	right := c.Orientation.MulVec3(mmath.UNIT_X_VEC3).Normalized()
	if right.LengthSqr() == 0 {
		return mmath.UNIT_X_VEC3
	}
	return right
}

// UpVector はカメラ上方向ベクトルを返す。
func (c *Camera) UpVector() mmath.Vec3 {
	if c == nil {
		return mmath.UNIT_Y_VEC3
	}
	c.ensureState()
	up := c.Orientation.MulVec3(mmath.UNIT_Y_VEC3).Normalized()
	if up.LengthSqr() == 0 {
		return mmath.UNIT_Y_VEC3
	}
	return up
}

// ForwardVector はカメラ前方向ベクトルを返す。
func (c *Camera) ForwardVector() mmath.Vec3 {
	if c == nil {
		return mmath.UNIT_Z_VEC3
	}
	c.ensureState()
	forward := c.Orientation.MulVec3(mmath.UNIT_Z_VEC3).Normalized()
	if forward.LengthSqr() == 0 {
		return mmath.UNIT_Z_VEC3
	}
	return forward
}

// ensureState はカメラ内部状態を初期化・補正する。
func (c *Camera) ensureState() {
	if c.Position == nil {
		c.Position = &mmath.Vec3{Vec: r3.Vec{X: 0.0, Y: InitialCameraPositionY, Z: InitialCameraPositionZ}}
	}
	if c.LookAtCenter == nil {
		c.LookAtCenter = &mmath.Vec3{Vec: r3.Vec{X: InitialLookAtCenterX, Y: InitialLookAtCenterY, Z: InitialLookAtCenterZ}}
	}
	if c.Up == nil {
		c.Up = &mmath.Vec3{Vec: r3.Vec{X: 0.0, Y: 1.0, Z: 0.0}}
	}
	if c.Orientation.Length() <= 1e-8 {
		forward := c.LookAtCenter.Subed(*c.Position)
		if forward.LengthSqr() == 0 {
			c.Orientation = mmath.NewQuaternion()
		} else {
			c.Orientation = mmath.NewQuaternionFromDirection(forward, *c.Up)
		}
	}
	c.Orientation = c.Orientation.Normalized()
}

// applyOrientation は姿勢と距離からカメラ位置と上方向を再構成する。
func (c *Camera) applyOrientation(distance float64) {
	if c == nil {
		return
	}
	c.ensureState()
	if distance <= 1e-8 {
		distance = math.Abs(InitialCameraPositionZ)
	}

	offset := c.Orientation.MulVec3(mmath.UNIT_Z_NEG_VEC3).MuledScalar(distance)
	c.Position.X = c.LookAtCenter.X + offset.X
	c.Position.Y = c.LookAtCenter.Y + offset.Y
	c.Position.Z = c.LookAtCenter.Z + offset.Z

	up := c.Orientation.MulVec3(mmath.UNIT_Y_VEC3).Normalized()
	if up.LengthSqr() == 0 {
		up = mmath.UNIT_Y_VEC3
	}
	c.Up.X = up.X
	c.Up.Y = up.Y
	c.Up.Z = up.Z
}
