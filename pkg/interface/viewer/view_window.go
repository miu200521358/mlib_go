//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"image"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/animation"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/walk"
)

// 直角の定数値
const rightAngle = 89.9

type ViewWindow struct {
	*glfw.Window
	windowIndex         int                // ウィンドウインデックス
	title               string             // ウィンドウタイトル
	appConfig           *mconfig.AppConfig // アプリケーション設定
	appState            state.IAppState    // アプリ状態
	physics             *mbt.MPhysics      // 物理
	shader              *mgl.MShader       // シェーダ
	leftButtonPressed   bool               // 左ボタン押下フラグ
	middleButtonPressed bool               // 中ボタン押下フラグ
	rightButtonPressed  bool               // 右ボタン押下フラグ
	shiftPressed        bool               // Shiftキー押下フラグ
	ctrlPressed         bool               // Ctrlキー押下フラグ
	updatedPrevCursor   bool               // 前回のカーソル位置更新フラグ
	prevCursorPos       *mmath.MVec2       // 前回のカーソル位置
	prevLeftCursorPos   *mmath.MRect       // 前回の左クリック位置
	yaw                 float64            // カメラyaw
	pitch               float64            // カメラpitch
	size                *mmath.MVec2       // ウィンドウサイズ
}

func NewViewWindow(
	windowIndex int,
	appConfig *mconfig.AppConfig,
	appState state.IAppState,
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
		Window:            glWindow,
		appState:          appState,
		windowIndex:       windowIndex,
		title:             title,
		appConfig:         appConfig,
		shader:            mgl.NewMShader(appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height),
		physics:           mbt.NewMPhysics(),
		prevCursorPos:     mmath.NewMVec2(),
		prevLeftCursorPos: mmath.NewMRect(),
		size: &mmath.MVec2{X: float64(appConfig.ViewWindowSize.Width),
			Y: float64(appConfig.ViewWindowSize.Height)},
	}

	glWindow.SetCloseCallback(viewWindow.closeCallback)
	glWindow.SetScrollCallback(viewWindow.scrollCallback)
	glWindow.SetKeyCallback(viewWindow.keyCallback)
	glWindow.SetMouseButtonCallback(viewWindow.mouseCallback)
	glWindow.SetCursorPosCallback(viewWindow.cursorPosCallback)
	glWindow.SetSizeCallback(viewWindow.resizeCallback)
	glWindow.SetFramebufferSizeCallback(viewWindow.resizeCallback)

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                        // 同期的なデバッグ出力有効
	gl.DebugMessageCallback(viewWindow.debugMessageCallback, nil) // デバッグコールバック

	return viewWindow
}

func (viewWindow *ViewWindow) Title() string {
	return viewWindow.title
}

func (viewWindow *ViewWindow) resizeCallback(w *glfw.Window, width int, height int) {
	viewWindow.size.X = float64(width)
	viewWindow.size.Y = float64(height)
}

func (viewWindow *ViewWindow) cursorPosCallback(w *glfw.Window, xpos, ypos float64) {

	if !viewWindow.updatedPrevCursor {
		viewWindow.prevCursorPos.X = xpos
		viewWindow.prevCursorPos.Y = ypos
		viewWindow.updatedPrevCursor = true
		return
	}

	if viewWindow.rightButtonPressed {
		// 右クリックはカメラの角度を更新
		viewWindow.updateCameraAngleByCursor(xpos, ypos)
	} else if viewWindow.middleButtonPressed {
		// 中クリックはカメラ位置と中心を移動
		viewWindow.updateCameraPositionByCursor(xpos, ypos)
	}

	viewWindow.prevCursorPos.X = xpos
	viewWindow.prevCursorPos.Y = ypos
}

func (viewWindow *ViewWindow) updateCameraAngleByCursor(xpos, ypos float64) {
	ratio := 0.1
	if viewWindow.shiftPressed {
		ratio *= 10
	} else if viewWindow.ctrlPressed {
		ratio *= 0.1
	}

	// 右クリックはカメラ中心をそのままにカメラ位置を変える
	xOffset := (viewWindow.prevCursorPos.X - xpos) * ratio
	yOffset := (viewWindow.prevCursorPos.Y - ypos) * ratio

	// 方位角と仰角を更新
	viewWindow.yaw += xOffset
	viewWindow.pitch += yOffset

	// 仰角の制限（水平面より上下に行き過ぎないようにする）
	viewWindow.pitch = mmath.ClampedFloat(viewWindow.pitch, -rightAngle, rightAngle)

	// 方位角の制限（360度を超えないようにする）
	if viewWindow.yaw > 360.0 {
		viewWindow.yaw -= 360.0
	} else if viewWindow.yaw < -360.0 {
		viewWindow.yaw += 360.0
	}

	viewWindow.updateCameraAngle()
}

func (viewWindow *ViewWindow) updateCameraAngle() {
	// 球面座標系をデカルト座標系に変換
	radius := float64(mgl.INITIAL_CAMERA_POSITION_Z)
	cameraX := radius * math.Cos(mgl64.DegToRad(viewWindow.pitch)) * math.Cos(mgl64.DegToRad(viewWindow.yaw))
	cameraY := radius * math.Sin(mgl64.DegToRad(viewWindow.pitch))
	cameraZ := radius * math.Cos(mgl64.DegToRad(viewWindow.pitch)) * math.Sin(mgl64.DegToRad(viewWindow.yaw))

	// カメラ位置を更新
	viewWindow.shader.CameraPosition.X = cameraX
	viewWindow.shader.CameraPosition.Y = mgl.INITIAL_CAMERA_POSITION_Y + cameraY
	viewWindow.shader.CameraPosition.Z = cameraZ
}

func (viewWindow *ViewWindow) updateCameraPositionByCursor(xpos float64, ypos float64) {
	// 中ボタンが押された場合の処理
	ratio := 0.07
	if viewWindow.shiftPressed {
		ratio *= 10
	} else if viewWindow.ctrlPressed {
		ratio *= 0.1
	}

	xOffset := (viewWindow.prevCursorPos.X - xpos) * ratio
	yOffset := (viewWindow.prevCursorPos.Y - ypos) * ratio

	// カメラの向きに基づいて移動方向を計算
	forward := viewWindow.shader.LookAtCenterPosition.Subed(viewWindow.shader.CameraPosition)
	right := forward.Cross(mmath.MVec3UnitY).Normalize()
	up := right.Cross(forward.Normalize()).Normalize()

	// 上下移動のベクトルを計算
	upMovement := up.MulScalar(-yOffset)
	// 左右移動のベクトルを計算
	rightMovement := right.MulScalar(-xOffset)

	// 移動ベクトルを合成してカメラ位置と中心を更新
	movement := upMovement.Add(rightMovement)
	viewWindow.shader.CameraPosition.Add(movement)
	viewWindow.shader.LookAtCenterPosition.Add(movement)
}

func (viewWindow *ViewWindow) mouseCallback(
	w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey,
) {
	if action == glfw.Press {
		switch button {
		case glfw.MouseButtonLeft:
			viewWindow.leftButtonPressed = true
			viewWindow.prevLeftCursorPos.Min.X, viewWindow.prevLeftCursorPos.Min.Y = viewWindow.GetCursorPos()
		case glfw.MouseButtonMiddle:
			viewWindow.middleButtonPressed = true
		case glfw.MouseButtonRight:
			viewWindow.rightButtonPressed = true
		}
	} else if action == glfw.Release {
		switch button {
		case glfw.MouseButtonLeft:
			viewWindow.leftButtonPressed = false
			viewWindow.prevLeftCursorPos.Max.X, viewWindow.prevLeftCursorPos.Max.Y = viewWindow.GetCursorPos()
		case glfw.MouseButtonMiddle:
			viewWindow.middleButtonPressed = false
		case glfw.MouseButtonRight:
			viewWindow.rightButtonPressed = false
		}
	}
}

func (viewWindow *ViewWindow) keyCallback(
	w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey,
) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			viewWindow.shiftPressed = true
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			viewWindow.ctrlPressed = true
		}
	} else if action == glfw.Release {
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			viewWindow.shiftPressed = false
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			viewWindow.ctrlPressed = false
		}
	}

	switch key {
	case glfw.KeyKP0: // 下面から
		viewWindow.yaw = rightAngle
		viewWindow.pitch = rightAngle
	case glfw.KeyKP2: // 正面から
		viewWindow.yaw = rightAngle
		viewWindow.pitch = 0
	case glfw.KeyKP4: // 左面から
		viewWindow.yaw = 180
		viewWindow.pitch = 0
	case glfw.KeyKP5: // 上面から
		viewWindow.yaw = rightAngle
		viewWindow.pitch = -rightAngle
	case glfw.KeyKP6: // 右面から
		viewWindow.yaw = 0
		viewWindow.pitch = 0
	case glfw.KeyKP8: // 背面から
		viewWindow.yaw = -rightAngle
		viewWindow.pitch = 0
	default:
		return
	}

	viewWindow.updateCameraAngle()
}

func (viewWindow *ViewWindow) scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	ratio := float32(1.0)
	if viewWindow.shiftPressed {
		ratio *= 10
	} else if viewWindow.ctrlPressed {
		ratio *= 0.1
	}

	if yoff > 0 {
		viewWindow.shader.FieldOfViewAngle -= ratio
		if viewWindow.shader.FieldOfViewAngle < 1.0 {
			viewWindow.shader.FieldOfViewAngle = 1.0
		}
	} else if yoff < 0 {
		viewWindow.shader.FieldOfViewAngle += ratio
	}
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

func (viewWindow *ViewWindow) Dispose() {
	viewWindow.Window.Destroy()
}

func (viewWindow *ViewWindow) Close() {
	viewWindow.Window.Destroy()
}

func (viewWindow *ViewWindow) Size() (int, int) {
	return viewWindow.appConfig.ViewWindowSize.Width, viewWindow.appConfig.ViewWindowSize.Height
}

func (viewWindow *ViewWindow) SetPosition(x, y int) {
	viewWindow.SetPos(x, y)
}

func (viewWindow *ViewWindow) AppState() state.IAppState {
	return viewWindow.appState
}

func (viewWindow *ViewWindow) TriggerClose(window *glfw.Window) {
	viewWindow.appState.SetClosed(true)
}

func (viewWindow *ViewWindow) GetWindow() *glfw.Window {
	return viewWindow.Window
}

func (viewWindow *ViewWindow) ResetPhysics(animationStates []state.IAnimationState) {
	// 物理を削除
	for _, s := range animationStates {
		if s.Model() != nil && s.RenderModel() != nil {
			viewWindow.physics.DeleteModel(s.ModelIndex())
		}
	}
	// 物理ワールドを作り直す
	viewWindow.physics.ResetWorld()
	// 物理を登録
	for _, s := range animationStates {
		if s.Model() != nil && s.RenderModel() != nil {
			viewWindow.physics.AddModel(s.ModelIndex(), s.Model())
		}
	}
}

func (viewWindow *ViewWindow) Animate(
	animationStates []state.IAnimationState, nextStates []state.IAnimationState, timeStep float32,
) ([]state.IAnimationState, []state.IAnimationState) {
	glfw.PollEvents()

	if viewWindow.size.X == 0 || viewWindow.size.Y == 0 {
		return animationStates, nextStates
	}

	viewWindow.MakeContextCurrent()

	if viewWindow.size.X != float64(viewWindow.shader.Width) || viewWindow.size.Y != float64(viewWindow.shader.Height) {
		viewWindow.shader.Resize(int(viewWindow.size.X), int(viewWindow.size.Y))
	}

	// MSAAフレームバッファをバインド
	viewWindow.shader.Msaa.Bind()

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
	viewWindow.updateCamera()

	// 床描画
	viewWindow.shader.DrawFloor()

	if animationStates != nil {
		// 何かしら描画対象情報がある場合の処理
		for _, nextState := range nextStates {
			if nextState != nil {
				// モデルが指定されてたら初期化してセット
				if nextState.Model() != nil {
					modelIndex := nextState.ModelIndex()
					viewWindow.physics.DeleteModel(modelIndex)
					animationStates[nextState.ModelIndex()].Load(nextState.Model())
					viewWindow.physics.AddModel(modelIndex, nextState.Model())
					nextState.SetModel(nil)
				}
				// モーションが指定されてたらセット
				if nextState.Motion() != nil {
					animationStates[nextState.ModelIndex()].SetMotion(nextState.Motion())
					nextState.SetMotion(nil)
				}
				if nextState.VmdDeltas() != nil {
					animationStates[nextState.ModelIndex()].SetVmdDeltas(nextState.VmdDeltas())
					nextState.SetVmdDeltas(nil)
				}
				if nextState.RenderDeltas() != nil {
					animationStates[nextState.ModelIndex()].SetRenderDeltas(nextState.RenderDeltas())
					nextState.SetRenderDeltas(nil)
				}
			}
		}

		// デフォーム
		animation.Deform(viewWindow.physics, animationStates, viewWindow.appState, timeStep)

		// モデル描画
		for _, animationState := range animationStates {
			if animationState != nil {
				animationState.Render(viewWindow.shader, viewWindow.appState)
			}
		}

		// 物理描画
		viewWindow.physics.DrawDebugLines(viewWindow.shader,
			viewWindow.appState.IsShowRigidBodyFront() || viewWindow.appState.IsShowRigidBodyBack(),
			viewWindow.appState.IsShowJoint(), viewWindow.appState.IsShowRigidBodyFront())

		// 深度解決
		viewWindow.shader.Msaa.Resolve()
	}

	viewWindow.shader.Msaa.Unbind()

	viewWindow.SwapBuffers()

	return animationStates, nextStates
}

func (viewWindow *ViewWindow) getCameraParameter() (mgl32.Mat4, mgl32.Vec3, mgl32.Mat4) {
	// カメラの再計算
	projectionMatrix := mgl32.Perspective(
		mgl32.DegToRad(viewWindow.shader.FieldOfViewAngle),
		float32(viewWindow.shader.Width)/float32(viewWindow.shader.Height),
		viewWindow.shader.NearPlane,
		viewWindow.shader.FarPlane,
	)

	// カメラの位置
	cameraPosition := mgl.NewGlVec3(viewWindow.shader.CameraPosition)

	// カメラの中心
	lookAtCenter := mgl.NewGlVec3(viewWindow.shader.LookAtCenterPosition)
	viewMatrix := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})

	return projectionMatrix, cameraPosition, viewMatrix
}

func (viewWindow *ViewWindow) updateCamera() {
	projectionMatrix, cameraPosition, viewMatrix := viewWindow.getCameraParameter()

	for _, program := range viewWindow.shader.Programs() {
		// プログラムの切り替え
		gl.UseProgram(program)

		// カメラの再計算
		projectionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_PROJECTION_MATRIX))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projectionMatrix[0])

		// カメラの位置
		cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CAMERA_POSITION))
		gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

		// カメラの中心
		cameraUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_MATRIX))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

		gl.UseProgram(0)
	}
}
