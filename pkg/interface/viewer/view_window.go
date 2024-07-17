package viewer

import (
	"image"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/window"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex            int                // ウィンドウインデックス
	appConfig              *mconfig.AppConfig // アプリケーション設定
	appState               window.IAppState   // UI状態
	physics                *mbt.MPhysics      // 物理
	shader                 *mgl.MShader       // シェーダ
	doResetPhysicsStart    bool               // 物理リセット開始フラグ
	doResetPhysicsProgress bool               // 物理リセット中フラグ
	doResetPhysicsCount    int                // 物理リセット処理回数
}

func NewGlWindow(
	windowIndex int,
	appConfig *mconfig.AppConfig,
	uiState window.IAppState,
) *ViewWindow {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	w, err := glfw.CreateWindow(appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height,
		mi18n.T("ビューワー"), nil, nil)
	if err != nil {
		mlog.E("Failed to create window: %v", err)
		return nil
	}

	w.MakeContextCurrent()
	w.SetInputMode(glfw.StickyKeysMode, glfw.True)
	w.SetIcon([]image.Image{*appConfig.IconImage})

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		mlog.E("Failed to initialize OpenGL: %v", err)
		return nil
	}

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)             // 同期的なデバッグ出力有効
	gl.DebugMessageCallback(debugMessageCallback, nil) // デバッグコールバック

	viewWindow := &ViewWindow{
		Window:      w,
		windowIndex: windowIndex,
		appConfig:   appConfig,
		appState:    uiState,
		shader:      mgl.NewMShader(appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height),
		physics:     mbt.NewMPhysics(),
	}

	return viewWindow
}

func debugMessageCallback(
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
		mlog.E("[HIGH] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
			source, glType, severity, message)
		panic("critical OpenGL error")
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

func (w *ViewWindow) Dispose() {
	w.Window.Destroy()
}

func (w *ViewWindow) Close() {
	w.Window.Destroy()
}

func (w *ViewWindow) Size() (int, int) {
	return w.appConfig.ViewWindowSize.Width, w.appConfig.ViewWindowSize.Height
}

func (w *ViewWindow) SetPosition(x, y int) {
	w.SetPos(x, y)
}

func (w *ViewWindow) TriggerClose(window *glfw.Window) {
	w.appState.SetClosed(true)
}

func (w *ViewWindow) GetWindow() *glfw.Window {
	return w.Window
}

func (w *ViewWindow) ResetPhysicsStart() {
	// 物理ON・まだリセット中ではないの時だけリセット処理を行う
	if w.physics.Enabled && !w.doResetPhysicsProgress {
		// 一旦物理OFFにする
		w.physics.Enabled = false
		// 物理ワールドを作り直す
		w.physics.ResetWorld()
		w.doResetPhysicsStart = false
		w.doResetPhysicsProgress = true
		w.doResetPhysicsCount = 0
	}
}

func (w *ViewWindow) Render() {
	w.SwapBuffers()
	glfw.PollEvents()
}
