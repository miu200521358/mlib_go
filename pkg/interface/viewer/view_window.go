//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/core"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/walk"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex int                // ウィンドウインデックス
	title       string             // ウィンドウタイトル
	appConfig   *mconfig.AppConfig // アプリケーション設定
	appState    core.IAppState     // アプリ状態
	physics     *mbt.MPhysics      // 物理
	shader      *mgl.MShader       // シェーダ
	// animationStates []*renderer.AnimationState // アニメーションステート
	// doResetPhysicsStart    bool               // 物理リセット開始フラグ
	// doResetPhysicsProgress bool               // 物理リセット中フラグ
	// doResetPhysicsCount    int                // 物理リセット処理回数
}

func NewViewWindow(
	windowIndex int,
	appConfig *mconfig.AppConfig,
	appState core.IAppState,
	title string,
) *ViewWindow {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	glWindow, err := glfw.CreateWindow(
		appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height, title, nil, nil)
	if err != nil {
		widget.RaiseError(err)
	}

	glWindow.MakeContextCurrent()
	glWindow.SetInputMode(glfw.StickyKeysMode, glfw.True)
	glWindow.SetIcon([]image.Image{appConfig.IconImage})

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		widget.RaiseError(err)
	}

	viewWindow := &ViewWindow{
		Window:      glWindow,
		appState:    appState,
		windowIndex: windowIndex,
		title:       title,
		appConfig:   appConfig,
		shader:      mgl.NewMShader(appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height),
		physics:     mbt.NewMPhysics(),
	}

	glWindow.SetCloseCallback(viewWindow.closeCallback)

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                        // 同期的なデバッグ出力有効
	gl.DebugMessageCallback(viewWindow.debugMessageCallback, nil) // デバッグコールバック

	return viewWindow
}

func (viewWindow *ViewWindow) closeCallback(w *glfw.Window) {
	if !viewWindow.appState.IsClosed() {
		if result := walk.MsgBox(nil, mi18n.T("終了確認"), mi18n.T("終了確認メッセージ"),
			walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel); result == walk.DlgCmdOK {
			viewWindow.appState.SetClosed(true)
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
		widget.RaiseError(fmt.Errorf("[HIGH] GL CRITICAL ERROR: %v type = 0x%x, severity = 0x%x, message = %s",
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

func (w *ViewWindow) AppState() core.IAppState {
	return w.appState
}

func (w *ViewWindow) TriggerClose(window *glfw.Window) {
	w.appState.SetClosed(true)
}

func (w *ViewWindow) GetWindow() *glfw.Window {
	return w.Window
}

func (w *ViewWindow) ResetPhysicsStart() {
	// // 物理ON・まだリセット中ではないの時だけリセット処理を行う
	// if w.physics.Enabled && !w.doResetPhysicsProgress {
	// 	// 一旦物理OFFにする
	// 	w.physics.Enabled = false
	// 	// 物理ワールドを作り直す
	// 	w.physics.ResetWorld()
	// 	w.doResetPhysicsStart = false
	// 	w.doResetPhysicsProgress = true
	// 	w.doResetPhysicsCount = 0
	// }
}

// func (w *ViewWindow) load(states []core.IAnimationState) {
// 	for _, s := range states {
// 		if len(w.animationStates) <= s.ModelIndex() {
// 			w.animationStates = append(w.animationStates, renderer.NewAnimationState(w.windowIndex, s.ModelIndex()))
// 		}

// 		if s.Model() != nil {
// 			w.animationStates[s.ModelIndex()].SetModel(s.Model())
// 			w.animationStates[s.ModelIndex()].Load()
// 			w.physics.AddModel(s.ModelIndex(), s.Model())
// 		}
// 		if s.Motion() != nil {
// 			w.animationStates[s.ModelIndex()].SetMotion(s.Motion())
// 		}
// 		if s.VmdDeltas() != nil {
// 			w.animationStates[s.ModelIndex()].SetVmdDeltas(s.VmdDeltas())
// 		}
// 		if s.RenderDeltas() != nil {
// 			w.animationStates[s.ModelIndex()].SetRenderDeltas(s.RenderDeltas())
// 		}
// 		w.animationStates[s.ModelIndex()].SetFrame(s.Frame())
// 	}
// }

func (w *ViewWindow) Render(animationStates []core.IAnimationState, timeStep float32) {
	glfw.PollEvents()

	w.MakeContextCurrent()

	// MSAAフレームバッファをバインド
	w.shader.Msaa.Bind()

	// 深度バッファのクリア
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 隠面消去
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// カメラの再計算
	w.updateCamera()

	// 床描画
	w.shader.DrawFloor()

	if animationStates != nil {
		// 何かしら描画対象情報がある場合の処理

		// アニメーション
		renderer.Animate(w.physics, animationStates, w.appState, timeStep)

		// モデル描画
		for _, animationState := range animationStates {
			if animationState != nil {
				animationState.Render(w.shader, w.appState)
			}
		}

		// 物理描画
		w.physics.DrawDebugLines(w.shader, w.appState.IsShowRigidBodyFront() || w.appState.IsShowRigidBodyBack(),
			w.appState.IsShowJoint(), w.appState.IsShowRigidBodyFront())

		// 深度解決
		w.shader.Msaa.Resolve()
	}

	w.shader.Msaa.Unbind()

	w.SwapBuffers()
}

func (w *ViewWindow) updateCamera() {
	// カメラの再計算
	projection := mgl32.Perspective(
		mgl32.DegToRad(w.shader.FieldOfViewAngle),
		float32(w.shader.Width)/float32(w.shader.Height),
		w.shader.NearPlane,
		w.shader.FarPlane,
	)

	// カメラの位置
	cameraPosition := mgl.NewGlVec3(w.shader.CameraPosition)

	// カメラの中心
	lookAtCenter := mgl.NewGlVec3(w.shader.LookAtCenterPosition)
	camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})

	for _, program := range w.shader.GetPrograms() {
		// プログラムの切り替え
		gl.UseProgram(program)

		// カメラの再計算
		projectionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_PROJECTION_MATRIX))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

		// カメラの位置
		cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CAMERA_POSITION))
		gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

		// カメラの中心
		cameraUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_MATRIX))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

		gl.UseProgram(0)
	}
}