// 指示: miu200521358
package deform

import (
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

// IkLocalAxes はIK計算用のローカル軸を表す。
type IkLocalAxes struct {
	X mmath.Vec3
	Y mmath.Vec3
	Z mmath.Vec3
}

// IkSolveInput は角度制限計算の入力を表す。
type IkSolveInput struct {
	TotalIkQuat   mmath.Quaternion
	MinAngleLimit mmath.Vec3
	MaxAngleLimit mmath.Vec3
	AxisX         mmath.Vec3
	AxisY         mmath.Vec3
	AxisZ         mmath.Vec3
	Loop          int
	LoopCount     int
}

// IkSolveStepInput はIK計算1ステップの入力を表す。
type IkSolveStepInput struct {
	LinkRotation    mmath.Quaternion
	LimitedAxis     mmath.Vec3
	LinkAngle       float64
	MinAngleLimit   mmath.Vec3
	MaxAngleLimit   mmath.Vec3
	LocalMinLimit   mmath.Vec3
	LocalMaxLimit   mmath.Vec3
	AngleLimit      bool
	LocalAngleLimit bool
	Loop            int
	LoopCount       int
	RemoveTwist     bool
	FixedAxis       mmath.Vec3
	ChildAxis       mmath.Vec3
	LocalAxes       IkLocalAxes
}

// SolveIk は角度制限を適用した回転を返す。
func SolveIk(input IkSolveInput) mmath.Quaternion {
	ikMat := input.TotalIkQuat.ToMat4()
	minLimit := input.MinAngleLimit
	maxLimit := input.MaxAngleLimit
	switch {
	case minLimit.X > -math.Pi/2 && maxLimit.X < math.Pi/2:
		return solveIkAxisX(ikMat, input)
	case minLimit.Y > -math.Pi/2 && maxLimit.Y < math.Pi/2:
		return solveIkAxisY(ikMat, input)
	default:
		return solveIkAxisZ(ikMat, input)
	}
}

// SolveIkStep はIK計算1ステップの回転を返す。
func SolveIkStep(input IkSolveStepInput) mmath.Quaternion {
	linkAngle := input.LinkAngle
	axis := input.LimitedAxis
	// 回転軸が不正または角度がゼロなら、回転を適用せず現状維持する。
	if isInvalidFloat(linkAngle) || isInvalidVec3(axis) || axis.IsZero() || linkAngle == 0 {
		return input.LinkRotation
	}

	var ikQuat mmath.Quaternion
	if !input.FixedAxis.IsZero() {
		fixed := input.FixedAxis.Normalized()
		if !input.AngleLimit && !input.LocalAngleLimit {
			tmp := mmath.NewQuaternionFromAxisAngles(axis, linkAngle)
			ikQuat, _ = tmp.SeparateTwistByAxis(fixed)
		} else {
			if axis.Dot(fixed) < 0 {
				linkAngle = -linkAngle
			}
			ikQuat = mmath.NewQuaternionFromAxisAngles(fixed, linkAngle)
		}
	} else {
		ikQuat = mmath.NewQuaternionFromAxisAngles(axis, linkAngle)
		if input.RemoveTwist {
			_, ikQuat = ikQuat.SeparateTwistByAxis(input.ChildAxis)
		}
	}

	totalIkQuat := input.LinkRotation.Muled(ikQuat)
	if input.AngleLimit {
		return SolveIk(IkSolveInput{
			TotalIkQuat:   totalIkQuat,
			MinAngleLimit: input.MinAngleLimit,
			MaxAngleLimit: input.MaxAngleLimit,
			AxisX:         mmath.UNIT_X_VEC3,
			AxisY:         mmath.UNIT_Y_VEC3,
			AxisZ:         mmath.UNIT_Z_VEC3,
			Loop:          input.Loop,
			LoopCount:     input.LoopCount,
		})
	}
	if input.LocalAngleLimit {
		return SolveIk(IkSolveInput{
			TotalIkQuat:   totalIkQuat,
			MinAngleLimit: input.LocalMinLimit,
			MaxAngleLimit: input.LocalMaxLimit,
			AxisX:         input.LocalAxes.X,
			AxisY:         input.LocalAxes.Y,
			AxisZ:         input.LocalAxes.Z,
			Loop:          input.Loop,
			LoopCount:     input.LoopCount,
		})
	}
	return totalIkQuat
}

// getLinkAxis はリンク回転軸を返す。
func getLinkAxis(link model.IkLink, ikTargetLocalPos, ikLocalPos mmath.Vec3) mmath.Vec3 {
	axis := ikTargetLocalPos.Cross(ikLocalPos).Normalized()
	minLimit := link.MinAngleLimit
	maxLimit := link.MaxAngleLimit
	if link.LocalAngleLimit {
		minLimit = link.LocalMinAngleLimit
		maxLimit = link.LocalMaxAngleLimit
	}
	switch {
	case minLimit.IsOnlyX() || maxLimit.IsOnlyX():
		if axis.X < 0 {
			return mmath.UNIT_X_NEG_VEC3
		}
		return mmath.UNIT_X_VEC3
	case minLimit.IsOnlyY() || maxLimit.IsOnlyY():
		if axis.Y < 0 {
			return mmath.UNIT_Y_NEG_VEC3
		}
		return mmath.UNIT_Y_VEC3
	case minLimit.IsOnlyZ() || maxLimit.IsOnlyZ():
		if axis.Z < 0 {
			return mmath.UNIT_Z_NEG_VEC3
		}
		return mmath.UNIT_Z_VEC3
	}
	return axis
}

// fixedAxisOrZero は固定軸を返す。
func fixedAxisOrZero(bone *model.Bone) mmath.Vec3 {
	if boneHasFixedAxis(bone) {
		return bone.FixedAxis
	}
	return mmath.NewVec3()
}

// localAxes はローカル軸を返す。
func localAxes(modelData *model.PmxModel, bone *model.Bone) IkLocalAxes {
	x, y, z := boneLocalAxes(modelData, bone)
	return IkLocalAxes{X: x, Y: y, Z: z}
}

// solveIkAxisX はX軸制限の回転を返す。
func solveIkAxisX(ikMat mmath.Mat4, input IkSolveInput) mmath.Quaternion {
	fSX := -ikMat.AxisZ().Y
	fX := math.Asin(fSX)
	fCX := math.Cos(fX)
	if math.Abs(fX) > mmath.Gimbal1Rad {
		if fX < 0 {
			fX = -mmath.Gimbal1Rad
		} else {
			fX = mmath.Gimbal1Rad
		}
		fCX = math.Cos(fX)
	}
	fCXInv := 1.0 / fCX
	fSY := ikMat.AxisZ().X * fCXInv
	fCY := ikMat.AxisZ().Z * fCXInv
	fY := math.Atan2(fSY, fCY)
	fSZ := ikMat.AxisX().Y * fCXInv
	fCZ := ikMat.AxisY().Y * fCXInv
	fZ := math.Atan2(fSZ, fCZ)

	fX = getIkAxisValue(fX, input.MinAngleLimit.X, input.MaxAngleLimit.X, input.Loop, input.LoopCount)
	fY = getIkAxisValue(fY, input.MinAngleLimit.Y, input.MaxAngleLimit.Y, input.Loop, input.LoopCount)
	fZ = getIkAxisValue(fZ, input.MinAngleLimit.Z, input.MaxAngleLimit.Z, input.Loop, input.LoopCount)

	xQuat := mmath.NewQuaternionFromAxisAngles(input.AxisX, fX)
	yQuat := mmath.NewQuaternionFromAxisAngles(input.AxisY, fY)
	zQuat := mmath.NewQuaternionFromAxisAngles(input.AxisZ, fZ)
	return yQuat.Muled(xQuat).Muled(zQuat)
}

// solveIkAxisY はY軸制限の回転を返す。
func solveIkAxisY(ikMat mmath.Mat4, input IkSolveInput) mmath.Quaternion {
	fSY := -ikMat.AxisX().Z
	fY := math.Asin(fSY)
	fCY := math.Cos(fY)
	if math.Abs(fY) > mmath.Gimbal1Rad {
		if fY < 0 {
			fY = -mmath.Gimbal1Rad
		} else {
			fY = mmath.Gimbal1Rad
		}
		fCY = math.Cos(fY)
	}
	fCYInv := 1.0 / fCY
	fSX := ikMat.AxisY().Z * fCYInv
	fCX := ikMat.AxisZ().Z * fCYInv
	fX := math.Atan2(fSX, fCX)
	fSZ := ikMat.AxisX().Y * fCYInv
	fCZ := ikMat.AxisX().X * fCYInv
	fZ := math.Atan2(fSZ, fCZ)

	fX = getIkAxisValue(fX, input.MinAngleLimit.X, input.MaxAngleLimit.X, input.Loop, input.LoopCount)
	fY = getIkAxisValue(fY, input.MinAngleLimit.Y, input.MaxAngleLimit.Y, input.Loop, input.LoopCount)
	fZ = getIkAxisValue(fZ, input.MinAngleLimit.Z, input.MaxAngleLimit.Z, input.Loop, input.LoopCount)

	xQuat := mmath.NewQuaternionFromAxisAngles(input.AxisX, fX)
	yQuat := mmath.NewQuaternionFromAxisAngles(input.AxisY, fY)
	zQuat := mmath.NewQuaternionFromAxisAngles(input.AxisZ, fZ)
	return zQuat.Muled(yQuat).Muled(xQuat)
}

// solveIkAxisZ はZ軸制限の回転を返す。
func solveIkAxisZ(ikMat mmath.Mat4, input IkSolveInput) mmath.Quaternion {
	fSZ := ikMat.AxisY().X
	fZ := math.Asin(fSZ)
	fCZ := math.Cos(fZ)
	if math.Abs(fZ) > mmath.Gimbal1Rad {
		if fZ < 0 {
			fZ = -mmath.Gimbal1Rad
		} else {
			fZ = mmath.Gimbal1Rad
		}
		fCZ = math.Cos(fZ)
	}
	fCZInv := 1.0 / fCZ
	fSX := ikMat.AxisY().Z * fCZInv
	fCX := ikMat.AxisY().Y * fCZInv
	fX := math.Atan2(fSX, fCX)
	fSY := ikMat.AxisX().X * fCZInv
	fCY := ikMat.AxisZ().X * fCZInv
	fY := math.Atan2(fSY, fCY)

	fX = getIkAxisValue(fX, input.MinAngleLimit.X, input.MaxAngleLimit.X, input.Loop, input.LoopCount)
	fY = getIkAxisValue(fY, input.MinAngleLimit.Y, input.MaxAngleLimit.Y, input.Loop, input.LoopCount)
	fZ = getIkAxisValue(fZ, input.MinAngleLimit.Z, input.MaxAngleLimit.Z, input.Loop, input.LoopCount)

	xQuat := mmath.NewQuaternionFromAxisAngles(input.AxisX, fX)
	yQuat := mmath.NewQuaternionFromAxisAngles(input.AxisY, fY)
	zQuat := mmath.NewQuaternionFromAxisAngles(input.AxisZ, fZ)
	return xQuat.Muled(zQuat).Muled(yQuat)
}

// getIkAxisValue は角度制限を反映する。
func getIkAxisValue(fV, minAngle, maxAngle float64, loop, loopCount int) float64 {
	isInLoop := float64(loop) < float64(loopCount)/2.0
	if fV < minAngle {
		tf := 2*minAngle - fV
		if tf <= maxAngle && isInLoop {
			fV = tf
		} else {
			fV = minAngle
		}
	}
	if fV > maxAngle {
		tf := 2*maxAngle - fV
		if tf >= minAngle && isInLoop {
			fV = tf
		} else {
			fV = maxAngle
		}
	}
	return fV
}
