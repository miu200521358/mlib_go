//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/walk/pkg/walk"
)

// 直角の定数値
const rightAngle = 89.9

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
