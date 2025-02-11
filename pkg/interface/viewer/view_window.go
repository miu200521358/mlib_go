package viewer

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
	"github.com/miu200521358/walk/pkg/walk"
)

// 直角の定数値
const rightAngle = 89.9

type ViewWindow struct {
	*glfw.Window
	windowIndex int                // ウィンドウインデックス
	title       string             // ウィンドウタイトル
	shared      *state.SharedState // SharedState への参照
	appConfig   *mconfig.AppConfig // アプリケーション設定
}

func NewViewWindow(
	windowIndex int,
	title string,
	shared *state.SharedState,
	appConfig *mconfig.AppConfig,
	mainWindow *glfw.Window,
) (*ViewWindow, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	glWindow, err := glfw.CreateWindow(
		appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height, title, nil, mainWindow)
	if err != nil {
		return nil, err
	}

	glWindow.MakeContextCurrent()
	glWindow.SetInputMode(glfw.StickyKeysMode, glfw.True)
	glWindow.SetIcon([]image.Image{appConfig.IconImage})

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		return nil, err
	}

	viewWindow := &ViewWindow{
		Window:      glWindow,
		windowIndex: windowIndex,
		title:       title,
		shared:      shared,
		appConfig:   appConfig,
	}

	glWindow.SetCloseCallback(viewWindow.closeCallback)
	// glWindow.SetScrollCallback(viewWindow.scrollCallback)
	// glWindow.SetKeyCallback(viewWindow.keyCallback)
	// glWindow.SetMouseButtonCallback(viewWindow.mouseCallback)
	// glWindow.SetCursorPosCallback(viewWindow.cursorPosCallback)
	// glWindow.SetSizeCallback(viewWindow.resizeCallback)
	// glWindow.SetFramebufferSizeCallback(viewWindow.resizeCallback)

	if !appConfig.IsEnvProd() {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                        // 同期的なデバッグ出力有効
		gl.DebugMessageCallback(viewWindow.debugMessageCallback, nil) // デバッグコールバック
	}

	return viewWindow, nil
}

func (viewWindow *ViewWindow) closeCallback(w *glfw.Window) {
	// controllerStateを読み取り
	if !viewWindow.shared.IsClosed() {
		// ビューワーがまだ閉じていない場合のみ、確認ダイアログを表示
		if result := walk.MsgBox(
			nil,
			mi18n.T("終了確認"),
			mi18n.T("終了確認メッセージ"),
			walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel,
		); result == walk.DlgCmdOK {
			viewWindow.shared.SetClosed(true)
		}
	}
}

func (viewWindow *ViewWindow) debugMessageCallback(
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

func (viewWindow *ViewWindow) Run() {
	for {
		if viewWindow.shared.IsClosed() {
			break
		}

		glfw.PollEvents()
		viewWindow.MakeContextCurrent()

		// 深度バッファのクリア
		gl.ClearColor(0.7, 0.7, 0.7, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 隠面消去
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)

		// ブレンディングを有効にする
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}

	viewWindow.Destroy()
}
