package window

import (
	"math"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

func (w *GlWindow) resizeBuffer(window *glfw.Window, width int, height int) {
	w.width = width
	w.height = height
	if width > 0 && height > 0 {
		gl.Viewport(0, 0, int32(width), int32(height))
	}
}

func (w *GlWindow) resize(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	w.shader.Resize(width, height)
}

func (w *GlWindow) handleKeyEvent(
	window *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if !(action == glfw.Press || action == glfw.Repeat) ||
		!(key == glfw.KeyKP0 ||
			key == glfw.KeyKP2 ||
			key == glfw.KeyKP4 ||
			key == glfw.KeyKP5 ||
			key == glfw.KeyKP6 ||
			key == glfw.KeyKP8 ||
			key == glfw.KeyLeftShift ||
			key == glfw.KeyRightShift ||
			key == glfw.KeyLeftControl ||
			key == glfw.KeyRightControl ||
			key == glfw.KeyLeft ||
			key == glfw.KeyRight ||
			key == glfw.KeyUp ||
			key == glfw.KeyDown) {
		return
	}

	if key == glfw.KeyRight || key == glfw.KeyUp {
		w.UiState.AddFrame(1.0)
	} else if key == glfw.KeyLeft || key == glfw.KeyDown {
		w.UiState.AddFrame(-1.0)
	}

	if key == glfw.KeyLeftShift || key == glfw.KeyRightShift {
		if action == glfw.Press {
			w.UiState.ShiftPressed = true
		} else if action == glfw.Release {
			w.UiState.ShiftPressed = false
		}
		return
	} else if key == glfw.KeyLeftControl || key == glfw.KeyRightControl {
		if action == glfw.Press {
			w.UiState.CtrlPressed = true
		} else if action == glfw.Release {
			w.UiState.CtrlPressed = false
		}
		return
	}

	w.TriggerViewReset()

	switch key {
	case glfw.KeyKP0: // 下面から
		w.yaw = RIGHT_ANGLE
		w.pitch = RIGHT_ANGLE
	case glfw.KeyKP2: // 正面から
		w.yaw = RIGHT_ANGLE
		w.pitch = 0
	case glfw.KeyKP4: // 左面から
		w.yaw = 180
		w.pitch = 0
	case glfw.KeyKP5: // 上面から
		w.yaw = RIGHT_ANGLE
		w.pitch = -RIGHT_ANGLE
	case glfw.KeyKP6: // 右面から
		w.yaw = 0
		w.pitch = 0
	case glfw.KeyKP8: // 背面から
		w.yaw = -RIGHT_ANGLE
		w.pitch = 0
	default:
		return
	}

	// カメラの新しい位置を計算
	radius := mgl.INITIAL_CAMERA_POSITION_Z

	// 球面座標系をデカルト座標系に変換
	cameraX := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw))
	cameraY := radius * math.Sin(mgl64.DegToRad(w.pitch))
	cameraZ := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw))

	// カメラ位置を更新
	w.shader.CameraPosition.X = cameraX
	w.shader.CameraPosition.Y = mgl.INITIAL_CAMERA_POSITION_Y + cameraY
	w.shader.CameraPosition.Z = cameraZ
}

func (w *GlWindow) handleScrollEvent(window *glfw.Window, xoff float64, yoff float64) {
	ratio := float32(1.0)
	if w.UiState.ShiftPressed {
		ratio *= 3
	} else if w.UiState.CtrlPressed {
		ratio *= 0.1
	}

	if yoff > 0 {
		w.shader.FieldOfViewAngle -= ratio
		if w.shader.FieldOfViewAngle < 2.0 {
			w.shader.FieldOfViewAngle = 2.0
		}
	} else if yoff < 0 {
		w.shader.FieldOfViewAngle += ratio
	}
}

func (w *GlWindow) handleMouseButtonEvent(
	window *glfw.Window,
	button glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey,
) {
	if button == glfw.MouseButtonMiddle {
		if action == glfw.Press {
			w.UiState.MiddleButtonPressed = true
			w.UiState.UpdatedPrev = false
		} else if action == glfw.Release {
			w.UiState.MiddleButtonPressed = false
		}
	} else if button == glfw.MouseButtonRight {
		if action == glfw.Press {
			w.UiState.RightButtonPressed = true
			w.UiState.UpdatedPrev = false
		} else if action == glfw.Release {
			w.UiState.RightButtonPressed = false
		}
	} else if button == glfw.MouseButtonLeft && w.worldPosFunc != nil && w.UiState.IsShowSelectedVertex {
		if action == glfw.Press {
			w.UiState.LeftButtonPressed = true
			w.UiState.UpdatedPrev = false
		} else if action == glfw.Release {
			w.UiState.LeftButtonPressed = false
			w.execWorldPos()
			w.nowCursorPos = &mmath.MVec2{X: 0, Y: 0}
		}
	}
}

func (w *GlWindow) updateCameraAngle(xpos, ypos float64) {

	ratio := 0.1
	if w.UiState.ShiftPressed {
		ratio *= 10
	} else if w.UiState.CtrlPressed {
		ratio *= 0.1
	}

	// 右クリックはカメラ中心をそのままにカメラ位置を変える
	xOffset := (w.prevCursorPos.X - xpos) * ratio
	yOffset := (w.prevCursorPos.Y - ypos) * ratio

	// 方位角と仰角を更新
	w.yaw += xOffset
	w.pitch += yOffset

	// 仰角の制限（水平面より上下に行き過ぎないようにする）
	if w.pitch > RIGHT_ANGLE {
		w.pitch = RIGHT_ANGLE
	} else if w.pitch < -RIGHT_ANGLE {
		w.pitch = -RIGHT_ANGLE
	}

	// 方位角の制限（360度を超えないようにする）
	if w.yaw > 360.0 {
		w.yaw -= 360.0
	} else if w.yaw < -360.0 {
		w.yaw += 360.0
	}

	// 球面座標系をデカルト座標系に変換
	// radius := float64(-w.shader.CameraPosition.Sub(w.shader.LookAtCenterPosition).Length())
	radius := float64(mgl.INITIAL_CAMERA_POSITION_Z)
	cameraX := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw))
	cameraY := radius * math.Sin(mgl64.DegToRad(w.pitch))
	cameraZ := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw))

	// カメラ位置を更新
	w.shader.CameraPosition.X = cameraX
	w.shader.CameraPosition.Y = mgl.INITIAL_CAMERA_POSITION_Y + cameraY
	w.shader.CameraPosition.X = cameraZ
	// mlog.D("xOffset %.7f, yOffset %.7f, CameraPosition: %s, LookAtCenterPosition: %s\n",
	// 	xOffset, yOffset, w.shader.CameraPosition.String(), w.shader.LookAtCenterPosition.String())
}

func (w *GlWindow) updateCameraPosition(xpos, ypos float64) {

	ratio := 0.07
	if w.UiState.ShiftPressed {
		ratio *= 5
	} else if w.UiState.CtrlPressed {
		ratio *= 0.1
	}
	// 中ボタンが押された場合の処理
	if w.UiState.MiddleButtonPressed {
		ratio := 0.07
		if w.UiState.ShiftPressed {
			ratio *= 5
		} else if w.UiState.CtrlPressed {
			ratio *= 0.1
		}

		xOffset := (w.prevCursorPos.X - xpos) * ratio
		yOffset := (w.prevCursorPos.Y - ypos) * ratio

		// カメラの向きに基づいて移動方向を計算
		forward := w.shader.LookAtCenterPosition.Subed(w.shader.CameraPosition)
		right := forward.Cross(mmath.MVec3UnitY).Normalize()
		up := right.Cross(forward.Normalize()).Normalize()

		// 上下移動のベクトルを計算
		upMovement := up.MulScalar(-yOffset)
		// 左右移動のベクトルを計算
		rightMovement := right.MulScalar(-xOffset)

		// 移動ベクトルを合成してカメラ位置と中心を更新
		movement := upMovement.Add(rightMovement)
		w.shader.CameraPosition.Add(movement)
		w.shader.LookAtCenterPosition.Add(movement)
	}
}

func (w *GlWindow) handleCursorPosEvent(window *glfw.Window, xpos float64, ypos float64) {
	// mlog.D("[start] yaw %.7f, pitch %.7f, CameraPosition: %s, LookAtCenterPosition: %s\n",
	// 	w.yaw, w.pitch, w.shader.CameraPosition.String(), w.shader.LookAtCenterPosition.String())

	if !w.UiState.UpdatedPrev {
		w.prevCursorPos.X = xpos
		w.prevCursorPos.Y = ypos
		w.UiState.UpdatedPrev = true
		return
	}

	if w.UiState.LeftButtonPressed {
		w.nowCursorPos.X = xpos
		w.nowCursorPos.Y = ypos
		return
	} else if w.UiState.RightButtonPressed {
		// 右クリックはカメラの角度を更新
		w.updateCameraAngle(xpos, ypos)
	} else if w.UiState.MiddleButtonPressed {
		// 中クリックはカメラ位置と中心を移動
		w.updateCameraPosition(xpos, ypos)
	}

	w.prevCursorPos.X = xpos
	w.prevCursorPos.Y = ypos
}
