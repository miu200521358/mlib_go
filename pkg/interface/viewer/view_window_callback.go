//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/walk/pkg/walk"
)

// 直角の定数値
const rightAngle = 89.9

// closeCallback はウィンドウのクローズイベントを処理する
func (vw *ViewWindow) closeCallback(w *glfw.Window) {
	// controllerStateを読み取り
	if !vw.list.appConfig.IsCloseConfirm() {
		vw.list.shared.SetClosed(true)
		return
	}
	if !vw.list.shared.IsClosed() {
		// ビューワーがまだ閉じていない場合のみ、確認ダイアログを表示
		if result := walk.MsgBox(
			nil,
			mi18n.T("終了確認"),
			mi18n.T("終了確認メッセージ"),
			walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel,
		); result == walk.DlgCmdOK {
			vw.list.shared.SetClosed(true)
		}
	}
}

// CameraPreset はカメラの視点プリセットを定義
type CameraPreset struct {
	Name  string  // プリセット名（デバッグ用）
	Yaw   float64 // 水平方向の角度
	Pitch float64 // 垂直方向の角度
}

// カメラの視点プリセット定義
var cameraPresets = map[glfw.Key]CameraPreset{
	glfw.KeyKP1: {"Bottom", 0, -rightAngle}, // 下面から
	glfw.KeyKP2: {"Front", 0, 0},            // 正面から
	glfw.KeyKP4: {"Left", -rightAngle, 0},   // 左面から
	glfw.KeyKP5: {"Top", 0, rightAngle},     // 上面から
	glfw.KeyKP6: {"Right", rightAngle, 0},   // 右面から
	glfw.KeyKP8: {"Back", 180, 0},           // 背面から
}

// keyCallback はキーボードのイベントを処理する
func (vw *ViewWindow) keyCallback(
	w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey,
) {
	// 修飾キーの処理
	if action == glfw.Press {
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			vw.shiftPressed = true
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			vw.ctrlPressed = true
			return
		}
	} else if action == glfw.Release {
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			vw.shiftPressed = false
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			vw.ctrlPressed = false
			return
		}
	}

	// カメラプリセットの適用
	if preset, exists := cameraPresets[key]; exists {
		vw.resetCameraPosition(preset.Yaw, preset.Pitch)
	}
}

// mouseCallback はマウスボタンのイベントを処理する
func (vw *ViewWindow) mouseCallback(
	w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey,
) {
	if action == glfw.Press {
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = true
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = true
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = true
		}
	} else if action == glfw.Release {
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = false
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = false
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = false
		}
	}
}

// cursorPosCallback はカーソル位置のイベントを処理する
func (vw *ViewWindow) cursorPosCallback(w *glfw.Window, xpos, ypos float64) {

	if !vw.updatedPrevCursor {
		vw.prevCursorPos.X = xpos
		vw.prevCursorPos.Y = ypos
		vw.updatedPrevCursor = true
		return
	}

	if vw.rightButtonPressed {
		// 右クリックはカメラの角度を更新
		vw.updateCameraAngleByCursor(xpos, ypos)
	} else if vw.middleButtonPressed {
		// 中クリックはカメラ位置と中心を移動
		vw.updateCameraPositionByCursor(xpos, ypos)
	} else if vw.leftButtonPressed {
		// 左クリックはカーソル位置を取得
		// if vw.ctrlPressed {
		// 	vw.leftCursorRemoveWindowPositions[mgl32.Vec2{float32(xpos), float32(ypos)}] = 0.0
		// } else {
		// 	vw.leftCursorWindowPositions[mgl32.Vec2{float32(xpos), float32(ypos)}] = 0.0
		// }
	}

	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos
}

// updateCameraAngleByCursor はカメラの角度をカーソル位置に基づいて更新する
func (vw *ViewWindow) updateCameraAngleByCursor(xpos, ypos float64) {
	ratio := 0.1
	if vw.shiftPressed {
		ratio *= 3
	} else if vw.ctrlPressed {
		ratio *= 0.1
	}

	// 右クリックはカメラ中心をそのままにカメラ位置を変える
	xOffset := (xpos - vw.prevCursorPos.X) * ratio
	yOffset := (ypos - vw.prevCursorPos.Y) * ratio

	// 方位角と仰角を更新
	vw.resetCameraPosition(vw.shader.Camera().Yaw+xOffset, vw.shader.Camera().Pitch+yOffset)
}

// updateCameraPositionByCursor はカメラ位置と中心をカーソル位置に基づいて更新する
func (vw *ViewWindow) updateCameraPositionByCursor(xpos float64, ypos float64) {
	// 中ボタンが押された場合の処理
	ratio := 0.07
	if vw.shiftPressed {
		ratio *= 3
	} else if vw.ctrlPressed {
		ratio *= 0.1
	}

	xOffset := (vw.prevCursorPos.X - xpos) * ratio
	yOffset := (vw.prevCursorPos.Y - ypos) * ratio

	cam := vw.shader.Camera()

	// カメラの向きに基づいて移動方向を計算
	forward := cam.LookAtCenter.Subed(cam.Position)
	right := forward.Cross(cam.Up).Normalize()
	up := right.Cross(forward.Normalize()).Normalize()

	// 上下移動のベクトルを計算
	upMovement := up.MulScalar(-yOffset)
	// 左右移動のベクトルを計算
	rightMovement := right.MulScalar(-xOffset)

	// 移動ベクトルを合成してカメラ位置と中心を更新
	movement := upMovement.Add(rightMovement)
	cam.Position.Add(movement)
	cam.LookAtCenter.Add(movement)

	vw.shader.SetCamera(cam)
	// カメラ同期が有効なら、他のウィンドウへも同じカメラ設定を反映
	vw.syncCameraToOthers()
}

// syncCameraToOthers は、現在のウィンドウのカメラ設定を他のウィンドウに反映する
func (vw *ViewWindow) syncCameraToOthers() {
	if !vw.list.shared.IsCameraSync() {
		return
	}

	currentCam := vw.shader.Camera()
	for _, otherVW := range vw.list.windowList {
		if otherVW.windowIndex != vw.windowIndex {
			otherCam := otherVW.shader.Camera()
			otherCam.Position.X = currentCam.Position.X
			otherCam.Position.Y = currentCam.Position.Y
			otherCam.Position.Z = currentCam.Position.Z
			otherCam.LookAtCenter.X = currentCam.LookAtCenter.X
			otherCam.LookAtCenter.Y = currentCam.LookAtCenter.Y
			otherCam.LookAtCenter.Z = currentCam.LookAtCenter.Z
			otherCam.Up.X = currentCam.Up.X
			otherCam.Up.Y = currentCam.Up.Y
			otherCam.Up.Z = currentCam.Up.Z
			otherCam.FieldOfView = currentCam.FieldOfView
			otherCam.AspectRatio = currentCam.AspectRatio
			otherCam.NearPlane = currentCam.NearPlane
			otherCam.FarPlane = currentCam.FarPlane
			otherCam.Yaw = currentCam.Yaw
			otherCam.Pitch = currentCam.Pitch
			otherVW.shader.SetCamera(otherCam)
		}
	}
}

// scrollCallback はマウスホイールのスクロールイベントを処理する
func (vw *ViewWindow) scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	step := float32(1.0)
	if vw.shiftPressed {
		step *= 5
	} else if vw.ctrlPressed {
		step *= 0.1
	}

	cam := vw.shader.Camera()

	if yoff > 0 {
		cam.FieldOfView -= step
		if cam.FieldOfView < 1.0 {
			cam.FieldOfView = 1.0
		}
	} else if yoff < 0 {
		cam.FieldOfView += step
	}

	vw.shader.SetCamera(cam)
	// カメラ同期が有効なら、他のウィンドウへも同じカメラ設定を反映
	vw.syncCameraToOthers()
}

func (vw *ViewWindow) focusCallback(w *glfw.Window, focused bool) {
	if !vw.list.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	if focused && !vw.list.shared.IsLinkingFocus() {
		// ユーザー操作等でウィンドウが前面になった場合に連動フォーカスを発火
		vw.list.shared.TriggerLinkedFocus(vw.windowIndex)
	}
}

func (vw *ViewWindow) sizeCallback(w *glfw.Window, width int, height int) {
	if !vw.list.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	vw.syncSizeToOthers(width, height)
}

// syncCameraToOthers は、現在のウィンドウのカメラ設定を他のウィンドウに反映する
func (vw *ViewWindow) syncSizeToOthers(width, height int) {
	if !vw.list.shared.IsShowOverride() {
		// オーバーライドが無効ならスルー
		return
	}

	for _, otherVW := range vw.list.windowList {
		if otherVW.windowIndex != vw.windowIndex {
			otherVW.SetSize(width, height)
		}
	}
}

func (vw *ViewWindow) iconifyCallback(w *glfw.Window, iconified bool) {
	if !vw.list.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	if iconified {
		vw.list.shared.SyncMinimize(vw.windowIndex)
	} else {
		vw.list.shared.SyncRestore(vw.windowIndex)
	}
}

// debugMessageCallback はOpenGLのデバッグメッセージを処理する
func (vw *ViewWindow) debugMessageCallback(
	source uint32,
	glType uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer,
) {
	switch severity {
	case gl.DEBUG_SEVERITY_HIGH:
		// 対象のエラーメッセージを検出
		if strings.Contains(message, "glBlitFramebuffer failed because the framebuffer configurations require that the source and destination sample counts match") {
			// 現在のウィンドウサイズを取得
			width, height := vw.GetSize()
			// middleMsaa を生成
			newMsaa := mgl.NewMiddleMsaa(width, height)
			// shader.SetMsaa() で差し替え
			vw.shader.SetMsaa(newMsaa)
			mlog.W(mi18n.T("中間MSAA差し替え"), message)
			return
		}

		panic(fmt.Errorf("[HIGH] GL CRITICAL ERROR: %v type = 0x%x, severity = 0x%x, message = %s",
			source, glType, severity, message))
	case gl.DEBUG_SEVERITY_MEDIUM:
		mlog.V("[MEDIUM] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
			source, glType, severity, message)
	case gl.DEBUG_SEVERITY_LOW:
		mlog.V("[LOW] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
			source, glType, severity, message)
		// case gl.DEBUG_SEVERITY_NOTIFICATION:
		// 	mlog.D("[NOTIFICATION] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
		// 		source, glType, severity, message)
	}
}
