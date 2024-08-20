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
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
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
	windowIndex                           int                    // ウィンドウインデックス
	title                                 string                 // ウィンドウタイトル
	appConfig                             *mconfig.AppConfig     // アプリケーション設定
	appState                              state.IAppState        // アプリ状態
	channelState                          state.IChannelState    // ビューワー→コントローラーチャンネル
	physics                               *mbt.MPhysics          // 物理
	shader                                *mgl.MShader           // シェーダー
	leftButtonPressed                     bool                   // 左ボタン押下フラグ
	middleButtonPressed                   bool                   // 中ボタン押下フラグ
	rightButtonPressed                    bool                   // 右ボタン押下フラグ
	shiftPressed                          bool                   // Shiftキー押下フラグ
	ctrlPressed                           bool                   // Ctrlキー押下フラグ
	updatedPrevCursor                     bool                   // 前回のカーソル位置更新フラグ
	prevCursorPos                         *mmath.MVec2           // 前回のカーソル位置
	leftCursorWindowPositions             map[mgl32.Vec2]float32 // 左クリック位置ウィンドウ座標リスト
	leftCursorRemoveWindowPositions       map[mgl32.Vec2]float32 // 左クリック位置ウィンドウ座標リスト(削除用)
	leftCursorWorldHistoryPositions       []*mgl32.Vec3          // 左クリック位置ワールド座標リスト
	leftCursorRemoveWorldHistoryPositions []*mgl32.Vec3          // 左クリック位置ワールド座標リスト(削除用)
	yaw                                   float64                // カメラyaw
	pitch                                 float64                // カメラpitch
	size                                  *mmath.MVec2           // ウィンドウサイズ
	renderModels                          []*render.RenderModel  // レンダリングモデルリスト
}

func NewViewWindow(
	windowIndex int,
	appConfig *mconfig.AppConfig,
	appState state.IAppState,
	channelState state.IChannelState,
	title string,
	mainWindow *glfw.Window,
) *ViewWindow {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	glWindow, err := glfw.CreateWindow(
		appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height, title, nil, mainWindow)
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
		Window:                                glWindow,
		appState:                              appState,
		channelState:                          channelState,
		windowIndex:                           windowIndex,
		title:                                 title,
		appConfig:                             appConfig,
		shader:                                mgl.NewMShader(appConfig.ViewWindowSize.Width, appConfig.ViewWindowSize.Height),
		physics:                               mbt.NewMPhysics(),
		prevCursorPos:                         mmath.NewMVec2(),
		leftCursorWindowPositions:             make(map[mgl32.Vec2]float32),
		leftCursorRemoveWindowPositions:       make(map[mgl32.Vec2]float32),
		leftCursorWorldHistoryPositions:       make([]*mgl32.Vec3, 0),
		leftCursorRemoveWorldHistoryPositions: make([]*mgl32.Vec3, 0),
		size: &mmath.MVec2{X: float64(appConfig.ViewWindowSize.Width),
			Y: float64(appConfig.ViewWindowSize.Height)},
		renderModels: make([]*render.RenderModel, 0),
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

func (viewWindow *ViewWindow) Physics() *mbt.MPhysics {
	return viewWindow.physics
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
	} else if viewWindow.leftButtonPressed {
		// 左クリックはカーソル位置を取得
		if viewWindow.ctrlPressed {
			viewWindow.leftCursorRemoveWindowPositions[mgl32.Vec2{float32(xpos), float32(ypos)}] = 0.0
		} else {
			viewWindow.leftCursorWindowPositions[mgl32.Vec2{float32(xpos), float32(ypos)}] = 0.0
		}
	}

	viewWindow.prevCursorPos.X = xpos
	viewWindow.prevCursorPos.Y = ypos
}

func (viewWindow *ViewWindow) updateCameraAngleByCursor(xpos, ypos float64) {
	ratio := 0.1
	if viewWindow.shiftPressed {
		ratio *= 3
	} else if viewWindow.ctrlPressed {
		ratio *= 0.1
	}

	// 右クリックはカメラ中心をそのままにカメラ位置を変える
	xOffset := (xpos - viewWindow.prevCursorPos.X) * ratio
	yOffset := (ypos - viewWindow.prevCursorPos.Y) * ratio

	// 方位角と仰角を更新
	viewWindow.yaw += xOffset
	viewWindow.pitch += yOffset

	viewWindow.updateCameraAngle()
}

func (viewWindow *ViewWindow) updateCameraAngle() {
	// 球面座標系をデカルト座標系に変換
	radius := math.Abs(float64(mgl.INITIAL_CAMERA_POSITION_Z))

	// 四元数を使ってカメラの方向を計算
	yawRad := mgl64.DegToRad(viewWindow.yaw)
	pitchRad := mgl64.DegToRad(viewWindow.pitch)
	orientation := mmath.NewMQuaternionFromAxisAngles(mmath.MVec3UnitY, yawRad).Mul(
		mmath.NewMQuaternionFromAxisAngles(mmath.MVec3UnitX, pitchRad))
	forwardXYZ := orientation.Rotate(mmath.MVec3UnitZInv).MulScalar(radius)

	// カメラ位置を更新
	viewWindow.shader.CameraPosition.X = forwardXYZ.X
	viewWindow.shader.CameraPosition.Y = mgl.INITIAL_CAMERA_POSITION_Y + forwardXYZ.Y
	viewWindow.shader.CameraPosition.Z = forwardXYZ.Z

	// forward := viewWindow.shader.LookAtCenterPosition.Subed(viewWindow.shader.CameraPosition).Normalize()
	// right := forward.Cross(viewWindow.shader.CameraUp).Normalize()
	// viewWindow.shader.CameraUp = right.Cross(forward).Normalize()
}

func (viewWindow *ViewWindow) updateCameraPositionByCursor(xpos float64, ypos float64) {
	// 中ボタンが押された場合の処理
	ratio := 0.07
	if viewWindow.shiftPressed {
		ratio *= 3
	} else if viewWindow.ctrlPressed {
		ratio *= 0.1
	}

	xOffset := (viewWindow.prevCursorPos.X - xpos) * ratio
	yOffset := (viewWindow.prevCursorPos.Y - ypos) * ratio

	// カメラの向きに基づいて移動方向を計算
	forward := viewWindow.shader.LookAtCenterPosition.Subed(viewWindow.shader.CameraPosition)
	right := forward.Cross(viewWindow.shader.CameraUp).Normalize()
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
			if !viewWindow.leftButtonPressed {
				viewWindow.leftCursorWindowPositions = make(map[mgl32.Vec2]float32)
				viewWindow.leftCursorRemoveWindowPositions = make(map[mgl32.Vec2]float32)
				viewWindow.leftCursorWorldHistoryPositions = make([]*mgl32.Vec3, 0)
				viewWindow.leftCursorRemoveWorldHistoryPositions = make([]*mgl32.Vec3, 0)
			}
			viewWindow.leftButtonPressed = true
		case glfw.MouseButtonMiddle:
			viewWindow.middleButtonPressed = true
		case glfw.MouseButtonRight:
			viewWindow.rightButtonPressed = true
		}
	} else if action == glfw.Release {
		switch button {
		case glfw.MouseButtonLeft:
			viewWindow.leftButtonPressed = false
		case glfw.MouseButtonMiddle:
			viewWindow.middleButtonPressed = false
		case glfw.MouseButtonRight:
			viewWindow.rightButtonPressed = false
		}
	}
}

// getWorldPosition は、マウスクリック位置 (スクリーン座標) と深度値を基にワールド座標を計算します
func (viewWindow *ViewWindow) getWorldPosition(mouseX, mouseY, depth float32) *mgl32.Vec3 {
	width := float32(viewWindow.size.X)
	height := float32(viewWindow.size.Y)

	// クリップ座標に変換
	ndcX := max(min(float32((2.0*mouseX)/float32(width)-1.0), 1.0), -1.0)
	ndcY := max(min(float32(1.0-(2.0*mouseY)/float32(height)), 1.0), -1.0)
	ndcZ := max(min(depth*2.0-1.0, 1.0), -1.0)

	clipCoords := mgl32.Vec4{ndcX, ndcY, ndcZ, 1.0}

	projectionMatrix, _, viewMatrix := viewWindow.getCameraParameter()

	// 視点座標系に変換
	viewCoords := projectionMatrix.Inv().Mul4x1(clipCoords)
	var viewCoordPos mgl32.Vec4
	if viewCoords.W() == 0.0 {
		viewCoordPos = mgl32.Vec4{viewCoords.X(), viewCoords.Y(), viewCoords.Z(), 1.0}
	} else {
		viewCoordPos = mgl32.Vec4{viewCoords.X() / viewCoords.W(),
			viewCoords.Y() / viewCoords.W(), viewCoords.Z() / viewCoords.W(), 1.0}
	}

	// ワールド座標系に変換
	worldCoords := viewMatrix.Inv().Mul4x1(viewCoordPos)
	var globalPos *mgl32.Vec3
	if worldCoords.W() == 0.0 {
		globalPos = &mgl32.Vec3{worldCoords.X(), worldCoords.Y(), worldCoords.Z()}
	} else {
		globalPos = &mgl32.Vec3{worldCoords.X() / worldCoords.W(),
			worldCoords.Y() / worldCoords.W(), worldCoords.Z() / worldCoords.W()}
	}

	return globalPos
}

func (viewWindow *ViewWindow) keyCallback(
	w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey,
) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			viewWindow.shiftPressed = true
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			viewWindow.ctrlPressed = true
			return
		}
	} else if action == glfw.Release {
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			viewWindow.shiftPressed = false
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			viewWindow.ctrlPressed = false
			return
		}
	}

	switch key {
	case glfw.KeyKP1: // 下面から
		viewWindow.yaw = 0
		viewWindow.pitch = -rightAngle
		viewWindow.shader.Reset(false)
		viewWindow.updateCameraAngle()
	case glfw.KeyKP2: // 正面から
		viewWindow.yaw = 0
		viewWindow.pitch = 0
		viewWindow.shader.Reset(false)
		viewWindow.updateCameraAngle()
	case glfw.KeyKP4: // 左面から
		viewWindow.yaw = -rightAngle
		viewWindow.pitch = 0
		viewWindow.shader.Reset(false)
		viewWindow.updateCameraAngle()
	case glfw.KeyKP5: // 上面から
		viewWindow.yaw = 0
		viewWindow.pitch = rightAngle
		viewWindow.shader.Reset(false)
		viewWindow.updateCameraAngle()
	case glfw.KeyKP6: // 右面から
		viewWindow.yaw = rightAngle
		viewWindow.pitch = 0
		viewWindow.shader.Reset(false)
		viewWindow.updateCameraAngle()
	case glfw.KeyKP8: // 背面から
		viewWindow.yaw = 180
		viewWindow.pitch = 0
		viewWindow.shader.Reset(false)
		viewWindow.updateCameraAngle()
	default:
		return
	}
}

func (viewWindow *ViewWindow) scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	ratio := float32(1.0)
	if viewWindow.shiftPressed {
		ratio *= 5
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

func (viewWindow *ViewWindow) updateCursorPositions() ([]*mgl32.Vec3, []*mgl32.Vec3) {
	var leftCursorWorldPositions []*mgl32.Vec3       // 左クリック位置ワールド座標リスト
	var leftCursorRemoveWorldPositions []*mgl32.Vec3 // 左クリック位置ワールド座標リスト(削除用)

	if viewWindow.appState.IsShowSelectedVertex() {
		leftCursorWorldPositions = make([]*mgl32.Vec3, 0)
		leftCursorRemoveWorldPositions = make([]*mgl32.Vec3, 0)

		// 頂点選択ONの場合、深度が取れてないのを取得
		for screenPos, depth := range viewWindow.leftCursorWindowPositions {
			if depth == 0 {
				for range 20 {
					depth := viewWindow.shader.Msaa.ReadDepthAt(int(screenPos.X()), int(screenPos.Y()))
					if depth > 0.0 && depth < 1.0 {
						viewWindow.leftCursorWindowPositions[screenPos] = depth

						worldPos := viewWindow.getWorldPosition(screenPos.X(), screenPos.Y(), depth)
						leftCursorWorldPositions = append(leftCursorWorldPositions, worldPos)
						viewWindow.leftCursorWorldHistoryPositions = append(
							viewWindow.leftCursorWorldHistoryPositions, worldPos)

						break
					}
				}
			}
		}
		for screenPos, depth := range viewWindow.leftCursorRemoveWindowPositions {
			if depth == 0 {
				for range 20 {
					depth := viewWindow.shader.Msaa.ReadDepthAt(int(screenPos.X()), int(screenPos.Y()))
					if depth > 0.0 && depth < 1.0 {
						viewWindow.leftCursorRemoveWindowPositions[screenPos] = depth

						worldPos := viewWindow.getWorldPosition(screenPos.X(), screenPos.Y(), depth)
						leftCursorRemoveWorldPositions = append(leftCursorRemoveWorldPositions, worldPos)
						viewWindow.leftCursorRemoveWorldHistoryPositions = append(
							viewWindow.leftCursorRemoveWorldHistoryPositions, worldPos)

						break
					}
				}
			}
		}
	}

	return leftCursorWorldPositions, leftCursorRemoveWorldPositions
}

func (viewWindow *ViewWindow) Render(
	models []*pmx.PmxModel, vmdDeltas []*delta.VmdDeltas, invisibleMaterials [][]int,
	windowSelectedVertexes [][]int, windowNoSelectedVertexes [][]int,
) [][]int {
	glfw.PollEvents()

	if viewWindow.size.X == 0 || viewWindow.size.Y == 0 {
		return make([][]int, 0)
	}

	viewWindow.MakeContextCurrent()

	if viewWindow.size.X != float64(viewWindow.shader.Width) || viewWindow.size.Y != float64(viewWindow.shader.Height) {
		viewWindow.shader.Resize(int(viewWindow.size.X), int(viewWindow.size.Y))
	}

	// カーソル位置の更新
	leftCursorWorldPositions, leftCursorRemoveWorldPositions := viewWindow.updateCursorPositions()

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

	// モデル描画
	selectedVertexes := make([][]int, len(models))
	for i, renderModel := range viewWindow.renderModels {
		if renderModel != nil && vmdDeltas[i] != nil {
			if i < len(invisibleMaterials) {
				renderModel.SetInvisibleMaterials(invisibleMaterials[i])
			}
			if i < len(windowNoSelectedVertexes) {
				renderModel.UpdateNoSelectedVertexes(windowNoSelectedVertexes[i])
			}
			if i < len(windowSelectedVertexes) {
				renderModel.UpdateSelectedVertexes(windowSelectedVertexes[i])
			}
			renderModel.Render(viewWindow.shader, viewWindow.appState, vmdDeltas[i],
				leftCursorWorldPositions, leftCursorRemoveWorldPositions,
				viewWindow.leftCursorWorldHistoryPositions, viewWindow.leftCursorRemoveWorldHistoryPositions)
			selectedVertexes[i] = renderModel.SelectedVertexes()
		}
	}

	// 物理デバッグ描画
	viewWindow.physics.DrawDebugLines(viewWindow.shader,
		viewWindow.appState.IsShowRigidBodyFront() || viewWindow.appState.IsShowRigidBodyBack(),
		viewWindow.appState.IsShowJoint(), viewWindow.appState.IsShowRigidBodyFront())

	// 深度解決
	viewWindow.shader.Msaa.Resolve()
	if viewWindow.appState.IsShowOverride() && viewWindow.windowIndex == 0 {
		viewWindow.drawOverride()
	}
	viewWindow.shader.Msaa.Unbind()
	viewWindow.SwapBuffers()

	if !viewWindow.leftButtonPressed {
		viewWindow.leftCursorWindowPositions = make(map[mgl32.Vec2]float32)
		viewWindow.leftCursorRemoveWindowPositions = make(map[mgl32.Vec2]float32)
		viewWindow.leftCursorWorldHistoryPositions = make([]*mgl32.Vec3, 0)
		viewWindow.leftCursorRemoveWorldHistoryPositions = make([]*mgl32.Vec3, 0)
	}

	return selectedVertexes
}

func (viewWindow *ViewWindow) LoadModels(models []*pmx.PmxModel) {
	viewWindow.MakeContextCurrent()

	for i, model := range models {
		if model == nil {
			continue
		}

		for i >= len(viewWindow.renderModels) {
			viewWindow.renderModels = append(viewWindow.renderModels, nil)
		}
		if viewWindow.renderModels[i] != nil && viewWindow.renderModels[i].Hash() != model.Hash() {
			// 既存モデルがいる場合、削除
			viewWindow.physics.DeleteModel(i)
			viewWindow.renderModels[i].Delete()
			viewWindow.renderModels[i] = nil
		}
		if viewWindow.renderModels[i] == nil {
			viewWindow.renderModels[i] = render.NewRenderModel(viewWindow.windowIndex, model)
			viewWindow.physics.AddModel(i, model)
		}
	}
}

func (viewWindow *ViewWindow) GetViewerParameter() *state.ViewerParameter {
	return &state.ViewerParameter{
		viewWindow.yaw,
		viewWindow.pitch,
		viewWindow.shader.FieldOfViewAngle,
		viewWindow.size,
		viewWindow.shader.CameraPosition,
		viewWindow.shader.CameraUp,
		viewWindow.shader.LookAtCenterPosition,
	}
}

func (viewWindow *ViewWindow) UpdateViewerParameter(
	viewerParameter *state.ViewerParameter,
) {
	viewWindow.yaw = viewerParameter.Yaw
	viewWindow.pitch = viewerParameter.Pitch

	viewWindow.shader.FieldOfViewAngle = viewerParameter.FieldOfViewAngle
	viewWindow.shader.CameraPosition = viewerParameter.CameraPos.Copy()
	viewWindow.shader.CameraUp = viewerParameter.CameraUp.Copy()
	viewWindow.shader.LookAtCenterPosition = viewerParameter.LookAtCenter.Copy()

	viewWindow.updateCameraAngle()

	copiedSize := viewerParameter.Size.Copy()
	isResize := !viewWindow.size.Equals(copiedSize)
	if isResize {
		viewWindow.Window.SetSize(int(copiedSize.X), int(copiedSize.Y))
		viewWindow.size = copiedSize
	}
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
	up := mgl.NewGlVec3(viewWindow.shader.CameraUp)
	viewMatrix := mgl32.LookAtV(cameraPosition, lookAtCenter, up)

	return projectionMatrix, cameraPosition, viewMatrix
}

func (viewWindow *ViewWindow) updateCamera() {
	projectionMatrix, cameraPosition, viewMatrix := viewWindow.getCameraParameter()

	for _, program := range viewWindow.shader.Programs() {
		// プログラムの切り替え
		gl.UseProgram(program)

		// カメラの再計算
		projectionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_PROJECTION_MATRIX))
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projectionMatrix[0])

		// カメラの位置
		cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CAMERA_POSITION))
		gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

		// カメラの中心
		cameraUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_VIEW_MATRIX))
		gl.UniformMatrix4fv(cameraUniform, 1, false, &viewMatrix[0])

		gl.UseProgram(0)
	}
}

func (viewWindow *ViewWindow) SetOverrideTextureId(overrideTextureId uint32) {
	viewWindow.shader.Msaa.SetOverrideTargetTexture(overrideTextureId)
}

func (viewWindow *ViewWindow) OverrideTextureId() uint32 {
	return viewWindow.shader.Msaa.OverrideTextureId()
}

func (viewWindow *ViewWindow) drawOverride() {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := viewWindow.shader.Program(mgl.PROGRAM_TYPE_OVERRIDE)
	gl.UseProgram(program)

	viewWindow.shader.Msaa.BindOverrideTexture(viewWindow.windowIndex, program)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	viewWindow.shader.Msaa.UnbindOverrideTexture()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}
