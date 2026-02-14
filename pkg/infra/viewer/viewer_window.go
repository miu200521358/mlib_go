//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mlib_go/pkg/adapter/physics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mbullet"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/render"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

const (
	rightAngle = 89.9
	// boneHoverMaxScreenDistance はボーンホバー判定の最大距離（ピクセル）。
	boneHoverMaxScreenDistance = 20.0
	// rigidBodyHoverRadiusScale は剛体ホバー判定の範囲を少し広げる補正係数。
	rigidBodyHoverRadiusScale = 1.2
	// rigidBodyHoverMinPixelRadius は剛体ホバー判定の最小ピクセル半径。
	rigidBodyHoverMinPixelRadius = 6.0
	// collisionAllFilterMask は全ての剛体を対象にするフィルタマスク。
	collisionAllFilterMask = int(^uint32(0))
	// cameraOperationBlockedWarnCooldown はカメラ操作ブロック警告の最小通知間隔。
	cameraOperationBlockedWarnCooldown = time.Second
)

// modelRendererLoadState はモデル描画の非同期読み込み状態を保持する。
type modelRendererLoadState struct {
	token       uint64
	hash        string
	sourceModel *model.PmxModel
	cancel      context.CancelFunc
	inProgress  bool
}

// modelRendererLoadResult は非同期読み込みの完了結果を保持する。
type modelRendererLoadResult struct {
	modelIndex  int
	token       uint64
	sourceModel *model.PmxModel
	renderModel *model.PmxModel
	bufferData  *render.ModelRendererBufferData
	err         error
}

// CameraPreset はカメラ視点プリセットを表す。
type CameraPreset struct {
	Name  string
	Yaw   float64
	Pitch float64
}

var cameraPresets = map[glfw.Key]CameraPreset{
	glfw.KeyKP1: {Name: "Bottom", Yaw: 0, Pitch: -rightAngle},
	glfw.KeyKP2: {Name: "Front", Yaw: 0, Pitch: 0},
	glfw.KeyKP4: {Name: "Left", Yaw: -rightAngle, Pitch: 0},
	glfw.KeyKP5: {Name: "Top", Yaw: 0, Pitch: rightAngle},
	glfw.KeyKP6: {Name: "Right", Yaw: rightAngle, Pitch: 0},
	glfw.KeyKP8: {Name: "Back", Yaw: 180, Pitch: 0},
}

// ViewerWindow はビューワーウィンドウを表す。
type ViewerWindow struct {
	*glfw.Window
	windowIndex int
	title       string
	list        *ViewerManager
	shader      graphics_api.IShader
	physics     *mbullet.PhysicsEngine

	tooltipRenderer               *mgl.TooltipRenderer
	boneHighlighter               *mgl.BoneHighlighter
	boneHoverActive               bool
	lastBoneHoverAt               time.Time
	rigidBodyHighlighter          *mgl.RigidBodyHighlighter
	rigidBodyHoverActive          bool
	lastRigidBodyHoverAt          time.Time
	jointHoverNames               []string
	jointHoverActive              bool
	lastJointHoverAt              time.Time
	selectedVertexHoverActive     bool
	selectedVertexHoverIndex      int
	selectedVertexHoverModelIndex int
	lastSelectedVertexHoverAt     time.Time
	lastCameraBlockedWarnAt       time.Time

	modelRenderers       []*render.ModelRenderer
	modelRendererLoads   []modelRendererLoadState
	loadResults          chan *modelRendererLoadResult
	motions              []*motion.VmdMotion
	cameraMotion         *motion.VmdMotion
	cameraManualOverride bool
	vmdDeltas            []*delta.VmdDeltas
	physicsModelHashes   []string

	leftButtonPressed               bool
	middleButtonPressed             bool
	rightButtonPressed              bool
	shiftPressed                    bool
	ctrlPressed                     bool
	updatedPrevCursor               bool
	prevCursorPos                   mmath.Vec2
	cursorX                         float64
	cursorY                         float64
	leftCursorWindowPositions       map[mmath.Vec2]float32
	leftCursorRemoveWindowPositions map[mmath.Vec2]float32
	// カーソル履歴の順序を保持するための記録リスト。
	leftCursorWindowOrder                 []mmath.Vec2
	leftCursorRemoveWindowOrder           []mmath.Vec2
	leftCursorWorldHistoryPositions       []*mmath.Vec3
	leftCursorRemoveWorldHistoryPositions []*mmath.Vec3
	boxSelectionDragging                  bool
	boxSelectionPending                   bool
	boxSelectionRemove                    bool
	boxSelectionStart                     mmath.Vec2
	boxSelectionEnd                       mmath.Vec2
}

// newViewerWindow はビューワーウィンドウを生成して初期化する。
func newViewerWindow(windowIndex int, title string, width, height, positionX, positionY int,
	appConfig *config.AppConfig, icon image.Image, mainWindow *glfw.Window, list *ViewerManager) (*ViewerWindow, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	glWindow, err := glfw.CreateWindow(width, height, title, nil, mainWindow)
	if err != nil {
		return nil, err
	}
	glWindow.MakeContextCurrent()
	glWindow.SetInputMode(glfw.StickyKeysMode, glfw.True)
	if icon != nil {
		glWindow.SetIcon([]image.Image{icon})
	}
	glfw.SwapInterval(0)

	if err := gl.Init(); err != nil {
		return nil, graphics_api.NewGraphicsContextInitFailed("OpenGLの初期化に失敗しました", err)
	}

	shaderFactory := mgl.NewShaderFactory(windowIndex)
	shader, err := shaderFactory.CreateShader(width, height)
	if err != nil {
		return nil, err
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	tooltipRenderer, err := mgl.NewTooltipRenderer()
	if err != nil {
		return nil, err
	}

	gravity := resolveInitialGravity(list, windowIndex)
	physics := mbullet.NewPhysicsEngine(&gravity)

	vw := &ViewerWindow{
		Window:                                glWindow,
		windowIndex:                           windowIndex,
		title:                                 title,
		list:                                  list,
		shader:                                shader,
		physics:                               physics,
		tooltipRenderer:                       tooltipRenderer,
		boneHighlighter:                       mgl.NewBoneHighlighter(),
		rigidBodyHighlighter:                  mgl.NewRigidBodyHighlighter(),
		modelRenderers:                        make([]*render.ModelRenderer, 0),
		modelRendererLoads:                    make([]modelRendererLoadState, 0),
		loadResults:                           make(chan *modelRendererLoadResult, 16),
		motions:                               make([]*motion.VmdMotion, 0),
		vmdDeltas:                             make([]*delta.VmdDeltas, 0),
		physicsModelHashes:                    make([]string, 0),
		selectedVertexHoverIndex:              -1,
		selectedVertexHoverModelIndex:         -1,
		leftCursorWindowPositions:             make(map[mmath.Vec2]float32),
		leftCursorRemoveWindowPositions:       make(map[mmath.Vec2]float32),
		leftCursorWindowOrder:                 make([]mmath.Vec2, 0),
		leftCursorRemoveWindowOrder:           make([]mmath.Vec2, 0),
		leftCursorWorldHistoryPositions:       make([]*mmath.Vec3, 0),
		leftCursorRemoveWorldHistoryPositions: make([]*mmath.Vec3, 0),
	}

	glWindow.SetCloseCallback(vw.closeCallback)
	glWindow.SetScrollCallback(vw.scrollCallback)
	glWindow.SetKeyCallback(vw.keyCallback)
	glWindow.SetMouseButtonCallback(vw.mouseCallback)
	glWindow.SetCursorPosCallback(vw.cursorPosCallback)
	glWindow.SetFocusCallback(vw.focusCallback)
	glWindow.SetIconifyCallback(vw.iconifyCallback)
	glWindow.SetSizeCallback(vw.sizeCallback)

	if appConfig == nil || !appConfig.IsProd() {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	}

	glWindow.SetPos(positionX, positionY)
	list.shared.SetViewerWindowReady(windowIndex, true)
	handle := int32(uintptr(unsafe.Pointer(glfw.GetCurrentContext().GetWin32Window())))
	list.shared.SetViewerWindowHandle(windowIndex, state.WindowHandle(handle))

	return vw, nil
}

// resolveInitialGravity は初期重力ベクトルを取得する。
func resolveInitialGravity(list *ViewerManager, windowIndex int) mmath.Vec3 {
	fallback := mmath.UNIT_Y_NEG_VEC3.MuledScalar(9.8)
	if list == nil || list.shared == nil {
		return fallback
	}
	raw := list.shared.PhysicsWorldMotion(windowIndex)
	motionData, ok := raw.(*motion.VmdMotion)
	if !ok || motionData == nil || motionData.GravityFrames == nil {
		return fallback
	}
	gravityFrame := motionData.GravityFrames.Get(0)
	if gravityFrame == nil || gravityFrame.Gravity == nil {
		return fallback
	}
	return *gravityFrame.Gravity
}

// ensureContextCurrent はOpenGLコンテキストを必要時だけカレントにする。
func (vw *ViewerWindow) ensureContextCurrent() {
	if vw == nil || vw.Window == nil {
		return
	}
	if glfw.GetCurrentContext() == vw.Window {
		return
	}
	vw.MakeContextCurrent()
}

// isSelectionEnabledInWindow はこのウィンドウでクリック選択/頂点選択を有効にするか返す。
func (vw *ViewerWindow) isSelectionEnabledInWindow() bool {
	return vw != nil && vw.windowIndex == 0
}

// isNonVertexHoverEnabledInWindow はこのウィンドウでボーン/剛体/ジョイントのホバーを有効にするか返す。
func (vw *ViewerWindow) isNonVertexHoverEnabledInWindow() bool {
	return vw != nil
}

// Title はタイトルを返す。
func (vw *ViewerWindow) Title() string {
	return vw.title
}

// SetTitle はタイトルを設定する。
func (vw *ViewerWindow) SetTitle(title string) {
	vw.title = title
	vw.Window.SetTitle(title)
}

// resetCameraPosition はカメラの視点をリセットして同期する。
func (vw *ViewerWindow) resetCameraPosition(yaw, pitch float64) {
	cam := vw.shader.Camera()
	if cam == nil {
		return
	}
	vw.markCameraManualOverride()
	cam.ResetPosition(yaw, pitch)
	if cam.FieldOfView == 0 {
		cam.FieldOfView = graphics_api.FieldOfViewAngle
	}
	vw.shader.SetCamera(cam)
	vw.syncCameraToOthers()
}

// resetCameraPositionForPreset はカメラの視点とFOVを既定値へ戻して同期する。
func (vw *ViewerWindow) resetCameraPositionForPreset(yaw, pitch float64) {
	cam := vw.shader.Camera()
	if cam == nil {
		return
	}
	vw.markCameraManualOverride()
	cam.ResetPosition(yaw, pitch)
	cam.FieldOfView = graphics_api.FieldOfViewAngle
	vw.shader.SetCamera(cam)
	vw.syncCameraToOthers()
}

// render は1フレーム分の描画を行う。
func (vw *ViewerWindow) render(frame motion.Frame) {
	winW, winH := vw.GetSize()
	fbW, fbH := vw.GetFramebufferSize()
	if fbW == 0 || fbH == 0 {
		fbW, fbH = winW, winH
	}
	if fbW == 0 || fbH == 0 {
		return
	}
	vw.ensureContextCurrent()
	vw.shader.Resize(fbW, fbH)

	showSelectedVertex := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX)
	selectionEnabled := vw.isSelectionEnabledInWindow()
	nonVertexHoverEnabled := vw.isNonVertexHoverEnabledInWindow()
	selectionMode := vw.selectedVertexMode()
	selectionDepthMode := state.SELECTED_VERTEX_DEPTH_MODE_ALL
	if vw.list != nil && vw.list.shared != nil {
		selectionDepthMode = vw.list.shared.SelectedVertexDepthMode()
	}
	applyBoxSelection := false
	boxSelectionRemove := false
	boxSelectionMin := mmath.Vec2{}
	boxSelectionMax := mmath.Vec2{}
	if showSelectedVertex && selectionEnabled && selectionMode == state.SELECTED_VERTEX_MODE_BOX {
		if minPos, maxPos, remove, ok := vw.consumeBoxSelectionRect(winW, winH, fbW, fbH); ok {
			applyBoxSelection = true
			boxSelectionRemove = remove
			boxSelectionMin = minPos
			boxSelectionMax = maxPos
			logger := logging.DefaultLogger()
			if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER) {
				logger.Verbose(
					logging.VERBOSE_INDEX_VIEWER,
					"ボックス選択: win=(%d,%d) fb=(%d,%d) start=(%.2f,%.2f) end=(%.2f,%.2f) rect=(%.2f,%.2f)-(%.2f,%.2f) remove=%t",
					winW,
					winH,
					fbW,
					fbH,
					vw.boxSelectionStart.X,
					vw.boxSelectionStart.Y,
					vw.boxSelectionEnd.X,
					vw.boxSelectionEnd.Y,
					boxSelectionMin.X,
					boxSelectionMin.Y,
					boxSelectionMax.X,
					boxSelectionMax.Y,
					boxSelectionRemove,
				)
			}
		}
	}

	if len(vw.list.windowList) > 1 && vw.list.shared.IsShowOverride() && vw.windowIndex != 0 {
		vw.shader.OverrideRenderer().Bind()
	} else {
		vw.shader.Msaa().Bind()
	}

	// 深度バッファのクリア
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 隠面消去
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// マルチサンプル有効化
	gl.Enable(gl.MULTISAMPLE)

	vw.applyCameraMotion(frame)
	vw.shader.UpdateCamera()

	vw.renderFloor()

	var debugBoneHover []*graphics_api.DebugBoneHover
	if selectionEnabled {
		vw.updateSelectedVertexHoverExpire()
	} else {
		vw.clearSelectedVertexHover()
	}
	if nonVertexHoverEnabled {
		vw.updateBoneHoverExpire()
		vw.updateRigidBodyHoverExpire()
		vw.updateJointHoverExpire()
		if vw.boneHighlighter != nil {
			debugBoneHover = vw.boneHighlighter.DebugBoneHoverInfo()
		}
	} else {
		vw.clearNonVertexHovers()
	}

	baseResults := make([]*render.ModelRenderBaseResult, len(vw.modelRenderers))
	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil {
			continue
		}
		limit := 0
		if vw.list.appConfig != nil {
			limit = vw.list.appConfig.CursorPositionLimit
		}
		renderer.SetCursorPositionLimit(limit)
		renderer.SetModelIndex(i)

		vmdDeltas := vw.vmdDeltas[i]
		baseResults[i] = renderer.RenderBase(
			vw.shader,
			vw.list.shared,
			vmdDeltas,
			debugBoneHover,
		)
	}

	var leftCursorWorldPositions []*mmath.Vec3
	var leftCursorDepths []float32
	var leftCursorRemoveWorldPositions []*mmath.Vec3
	var leftCursorRemoveDepths []float32
	hoverCursorPositions := []float32(nil)
	selectedCursorPositions := []float32(nil)
	selectedCursorDepths := []float32(nil)
	removeSelectedCursorPositions := []float32(nil)
	removeSelectedCursorDepths := []float32(nil)
	cursorLinePositions := []float32(nil)
	removeCursorLinePositions := []float32(nil)
	boxLinePositions := []float32(nil)
	applyPointSelection := false
	removePointSelection := false
	if showSelectedVertex && selectionEnabled {
		if resolver, ok := vw.shader.Msaa().(interface {
			ResolveDepth()
		}); ok {
			resolver.ResolveDepth()
		}
		if selectionMode == state.SELECTED_VERTEX_MODE_POINT {
			leftCursorWorldPositions, leftCursorDepths, leftCursorRemoveWorldPositions, leftCursorRemoveDepths = vw.updateCursorPositions()
		}
		limit := 0
		if vw.list.appConfig != nil {
			limit = vw.list.appConfig.CursorPositionLimit
		}
		hoverCursorPositions = vw.buildCursorPositions(vw.cursorX, vw.cursorY, limit)
		selectedCursorPositions = flattenCursorPositions(leftCursorWorldPositions, limit)
		selectedCursorDepths = flattenCursorDepths(leftCursorDepths, limit)
		removeSelectedCursorPositions = flattenCursorPositions(leftCursorRemoveWorldPositions, limit)
		removeSelectedCursorDepths = flattenCursorDepths(leftCursorRemoveDepths, limit)
		if selectionMode == state.SELECTED_VERTEX_MODE_POINT {
			// 軌跡表示は履歴全体を使うため、上限指定は行わない。
			cursorLinePositions = flattenCursorPositions(vw.leftCursorWorldHistoryPositions, 0)
			removeCursorLinePositions = flattenCursorPositions(vw.leftCursorRemoveWorldHistoryPositions, 0)
		}
		if selectionMode == state.SELECTED_VERTEX_MODE_BOX &&
			(vw.boxSelectionDragging || vw.boxSelectionPending) {
			boxLinePositions = vw.buildBoxSelectionLinePositions(winW, winH, fbW, fbH)
		}
		applyPointSelection = len(selectedCursorPositions) > 0 || len(removeSelectedCursorPositions) > 0
		removePointSelection = len(removeSelectedCursorPositions) > 0
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil {
			continue
		}
		base := baseResults[i]
		if base == nil {
			continue
		}
		selectedVertexIndexes := []int{}
		if showSelectedVertex {
			selectedVertexIndexes = vw.list.shared.SelectedVertexIndexes(vw.windowIndex, i)
		}
		cursorPositions := hoverCursorPositions
		cursorDepths := []float32(nil)
		removeCursorPositions := []float32(nil)
		removeCursorDepths := []float32(nil)
		applySelectionForModel := false
		selectionRequest := (*render.VertexSelectionRequest)(nil)
		if showSelectedVertex && selectionEnabled {
			selectionRequest = &render.VertexSelectionRequest{
				Mode:                      selectionMode,
				DepthMode:                 selectionDepthMode,
				Apply:                     false,
				Remove:                    false,
				CursorPositions:           cursorPositions,
				CursorDepths:              cursorDepths,
				RemoveCursorPositions:     removeCursorPositions,
				RemoveCursorDepths:        removeCursorDepths,
				CursorLinePositions:       nil,
				RemoveCursorLinePositions: nil,
				ScreenWidth:               fbW,
				ScreenHeight:              fbH,
				RectMin:                   boxSelectionMin,
				RectMax:                   boxSelectionMax,
				HasRect:                   false,
			}
		}
		if showSelectedVertex && selectionEnabled && i == 0 {
			if selectionMode == state.SELECTED_VERTEX_MODE_POINT && selectionRequest != nil {
				selectionRequest.CursorLinePositions = cursorLinePositions
				selectionRequest.RemoveCursorLinePositions = removeCursorLinePositions
			}
			if selectionMode == state.SELECTED_VERTEX_MODE_BOX && selectionRequest != nil {
				selectionRequest.CursorLinePositions = boxLinePositions
			}
			switch selectionMode {
			case state.SELECTED_VERTEX_MODE_BOX:
				if applyBoxSelection && selectionRequest != nil {
					selectionRequest.Apply = true
					selectionRequest.Remove = boxSelectionRemove
					selectionRequest.HasRect = true
					applySelectionForModel = true
				}
			default:
				if applyPointSelection && selectionRequest != nil {
					if removePointSelection {
						cursorPositions = removeSelectedCursorPositions
						cursorDepths = removeSelectedCursorDepths
						removeCursorPositions = cursorPositions
						removeCursorDepths = cursorDepths
					} else {
						cursorPositions = selectedCursorPositions
						cursorDepths = selectedCursorDepths
					}
					selectionRequest.CursorPositions = cursorPositions
					selectionRequest.CursorDepths = cursorDepths
					selectionRequest.RemoveCursorPositions = removeCursorPositions
					selectionRequest.RemoveCursorDepths = removeCursorDepths
					selectionRequest.Apply = true
					selectionRequest.Remove = removePointSelection
					applySelectionForModel = true
				}
			}
		}
		updatedSelected, hoverIndex := renderer.RenderSelection(
			vw.shader,
			vw.list.shared,
			selectedVertexIndexes,
			selectionRequest,
			base,
		)
		if applySelectionForModel {
			vw.list.shared.SetSelectedVertexIndexes(vw.windowIndex, i, updatedSelected)
		}
		if showSelectedVertex && selectionEnabled && i == 0 {
			vw.updateSelectedVertexHover(i, hoverIndex)
		}
	}
	if !showSelectedVertex {
		vw.clearSelectedVertexHover()
	}

	drawRigidBodyFront := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_RIGID_BODY_FRONT)
	drawRigidBodyBack := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_RIGID_BODY_BACK)
	drawJoint := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_JOINT)
	drawRigidBody := drawRigidBodyFront || drawRigidBodyBack

	if vw.physics != nil && (drawRigidBody || drawJoint) {
		vw.physics.DrawDebugLines(vw.shader, drawRigidBody, drawJoint, drawRigidBodyFront)
	}
	if vw.rigidBodyHighlighter != nil {
		if drawRigidBody {
			vw.rigidBodyHighlighter.DrawDebugHighlight(vw.shader, drawRigidBodyFront)
		}
		vw.rigidBodyHighlighter.CheckAndClearHighlightOnDebugChange(drawRigidBody)
	}

	vw.renderTooltip(drawRigidBody, drawJoint, winW, winH)

	if len(vw.list.windowList) > 1 && vw.list.shared.IsShowOverride() && vw.windowIndex != 0 {
		vw.shader.OverrideRenderer().Unbind()
	} else {
		vw.shader.Msaa().Resolve()
		vw.shader.Msaa().Unbind()
	}
	if len(vw.list.windowList) > 1 && vw.list.shared.IsShowOverride() && vw.windowIndex == 0 && vw.shader.OverrideRenderer().SharedTextureIDPtr() != nil {
		vw.shader.OverrideRenderer().Resolve()
	}

	vw.captureScreenshotIfRequested(fbW, fbH)
	vw.SwapBuffers()
	if !showSelectedVertex {
		vw.resetBoxSelection()
	}
	if showSelectedVertex && !vw.leftButtonPressed {
		vw.boxSelectionDragging = false
		if selectionMode == state.SELECTED_VERTEX_MODE_POINT {
			vw.leftCursorWindowPositions = make(map[mmath.Vec2]float32)
			vw.leftCursorRemoveWindowPositions = make(map[mmath.Vec2]float32)
			vw.leftCursorWindowOrder = make([]mmath.Vec2, 0)
			vw.leftCursorRemoveWindowOrder = make([]mmath.Vec2, 0)
			vw.leftCursorWorldHistoryPositions = make([]*mmath.Vec3, 0)
			vw.leftCursorRemoveWorldHistoryPositions = make([]*mmath.Vec3, 0)
		}
	}
}

// captureScreenshotIfRequested はスクリーンショット要求があれば保存する。
func (vw *ViewerWindow) captureScreenshotIfRequested(width, height int) {
	if vw == nil || vw.list == nil || vw.list.shared == nil {
		return
	}
	request, ok := vw.list.shared.DequeueScreenshot(vw.windowIndex)
	if !ok {
		return
	}
	err := vw.saveScreenshot(request.Path, width, height)
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
		logging.DefaultLogger().Warn("スクリーンショット保存に失敗しました: %s", errMessage)
	}
	vw.list.shared.CompleteScreenshot(state.ScreenshotResult{
		ID:         request.ID,
		Path:       request.Path,
		ErrMessage: errMessage,
	})
}

// saveScreenshot は現在のフレームバッファを画像として保存する。
func (vw *ViewerWindow) saveScreenshot(path string, width, height int) error {
	if path == "" {
		return fmt.Errorf("保存先パスが空です")
	}
	if width <= 0 || height <= 0 {
		return fmt.Errorf("スクリーンショットサイズが不正です: %dx%d", width, height)
	}

	pixels := make([]uint8, width*height*4)

	var prevReadBuffer int32
	gl.GetIntegerv(gl.READ_BUFFER, &prevReadBuffer)
	gl.ReadBuffer(gl.BACK)

	var prevPackAlign int32
	gl.GetIntegerv(gl.PACK_ALIGNMENT, &prevPackAlign)
	gl.PixelStorei(gl.PACK_ALIGNMENT, 1)

	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))

	gl.PixelStorei(gl.PACK_ALIGNMENT, prevPackAlign)
	gl.ReadBuffer(uint32(prevReadBuffer))

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		srcOffset := (height - 1 - y) * width * 4
		dstOffset := y * img.Stride
		copy(img.Pix[dstOffset:dstOffset+width*4], pixels[srcOffset:srcOffset+width*4])
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	}
	return nil
}

// cleanupResources はビューワーが保持するOpenGLリソースを解放する。
func (vw *ViewerWindow) cleanupResources() {
	if vw == nil {
		return
	}
	if vw.Window != nil {
		vw.ensureContextCurrent()
	}
	vw.cancelAllModelRendererLoads()

	for _, renderer := range vw.modelRenderers {
		if renderer != nil {
			renderer.Delete()
		}
	}
	vw.modelRenderers = nil

	if vw.tooltipRenderer != nil {
		vw.tooltipRenderer.Delete()
	}

	if vw.shader != nil {
		vw.shader.Cleanup()
	}
}

// renderFloor は床グリッドを描画する。
func (vw *ViewerWindow) renderFloor() {
	vw.shader.UseProgram(graphics_api.ProgramTypeFloor)
	vw.shader.FloorRenderer().Bind()
	vw.shader.FloorRenderer().Render()
	vw.shader.FloorRenderer().Unbind()
	vw.shader.ResetProgram()
}

// prepareFrame は描画前にモデル/モーションを同期し、描画停止が必要か返す。
func (vw *ViewerWindow) prepareFrame() bool {
	vw.ensureContextCurrent()
	loading := vw.loadModelRenderers()
	vw.loadMotions()
	vw.ensurePhysicsModelSlots()
	return loading
}

// ensurePhysicsModelSlots は物理同期用のスロット長をモデル数に合わせる。
func (vw *ViewerWindow) ensurePhysicsModelSlots() {
	for len(vw.physicsModelHashes) < len(vw.modelRenderers) {
		vw.physicsModelHashes = append(vw.physicsModelHashes, "")
	}
	if len(vw.physicsModelHashes) > len(vw.modelRenderers) {
		vw.physicsModelHashes = vw.physicsModelHashes[:len(vw.modelRenderers)]
	}
}

// syncPhysicsModel は物理エンジンにモデルを同期する。
func (vw *ViewerWindow) syncPhysicsModel(modelIndex int, modelData *model.PmxModel, vmdDeltas *delta.VmdDeltas, physicsDeltas *delta.PhysicsDeltas) {
	if vw.physics == nil || modelIndex < 0 {
		return
	}
	if modelIndex >= len(vw.physicsModelHashes) {
		return
	}
	if modelData == nil {
		if vw.physicsModelHashes[modelIndex] != "" {
			vw.physics.DeleteModel(modelIndex)
			vw.physicsModelHashes[modelIndex] = ""
		}
		return
	}
	if vmdDeltas == nil || vmdDeltas.Bones == nil {
		return
	}
	modelHash := modelData.Hash()
	if vw.physicsModelHashes[modelIndex] != modelHash {
		if vw.physicsModelHashes[modelIndex] != "" {
			vw.physics.DeleteModel(modelIndex)
		}
		vw.physics.AddModelByDeltas(modelIndex, modelData, vmdDeltas.Bones, physicsDeltas)
		vw.physicsModelHashes[modelIndex] = modelHash
	}
}

// loadModelRenderers はモデルレンダラを共有状態から同期し、描画停止が必要か返す。
func (vw *ViewerWindow) loadModelRenderers() bool {
	vw.applyModelRendererLoadResults()
	loading := false
	modelCount := vw.list.shared.ModelCount(vw.windowIndex)
	motionCount := vw.list.shared.MotionCount(vw.windowIndex)
	count := max(modelCount, motionCount)
	for i := 0; i < count; i++ {
		if i >= len(vw.modelRenderers) {
			vw.modelRenderers = append(vw.modelRenderers, nil)
		}
		if i >= len(vw.modelRendererLoads) {
			vw.modelRendererLoads = append(vw.modelRendererLoads, modelRendererLoadState{})
		}
		if i >= len(vw.motions) {
			vw.motions = append(vw.motions, nil)
		}
		if i >= len(vw.vmdDeltas) {
			vw.vmdDeltas = append(vw.vmdDeltas, nil)
		}

		var modelData *model.PmxModel
		if raw := vw.list.shared.Model(vw.windowIndex, i); raw != nil {
			if m, ok := raw.(*model.PmxModel); ok {
				modelData = m
			}
		}

		existing := vw.modelRenderers[i]
		if modelData == nil {
			if existing != nil {
				existing.Delete()
				vw.modelRenderers[i] = nil
				vw.vmdDeltas[i] = nil
			}
			vw.cancelModelRendererLoad(i)
			continue
		}
		shouldReload := existing == nil
		if !shouldReload {
			if existing.Hash() != modelData.Hash() {
				shouldReload = true
			} else if existing.SourceModel != modelData {
				shouldReload = true
			}
		}
		if vw.isModelRendererLoading(i, modelData) {
			loading = true
			continue
		}
		if shouldReload {
			if existing != nil {
				existing.Delete()
			}
			vw.modelRenderers[i] = nil
			vw.vmdDeltas[i] = nil
			vw.startModelRendererLoad(i, modelData)
			loading = true
		}
	}
	return loading
}

// applyModelRendererLoadResults は非同期読み込みの完了結果を反映する。
func (vw *ViewerWindow) applyModelRendererLoadResults() {
	if vw == nil || vw.loadResults == nil {
		return
	}
	for {
		select {
		case result := <-vw.loadResults:
			if result == nil {
				continue
			}
			if result.modelIndex < 0 || result.modelIndex >= len(vw.modelRendererLoads) {
				continue
			}
			state := &vw.modelRendererLoads[result.modelIndex]
			// 取消/再読み込み済みの結果は破棄する。
			if result.token != state.token {
				continue
			}
			state.inProgress = false
			state.cancel = nil
			if result.err != nil {
				logging.DefaultLogger().Warn("描画バッファの準備に失敗しました: %v", result.err)
				continue
			}
			if result.renderModel == nil || result.bufferData == nil {
				logging.DefaultLogger().Warn("描画バッファの準備結果が不正です")
				continue
			}
			if existing := vw.modelRenderers[result.modelIndex]; existing != nil {
				existing.Delete()
			}
			renderer := render.NewModelRendererWithPreparedData(vw.windowIndex, result.renderModel, result.bufferData)
			renderer.SourceModel = result.sourceModel
			vw.modelRenderers[result.modelIndex] = renderer
			vw.vmdDeltas[result.modelIndex] = nil
		default:
			return
		}
	}
}

// isModelRendererLoading は指定モデルが読み込み中か判定する。
func (vw *ViewerWindow) isModelRendererLoading(modelIndex int, modelData *model.PmxModel) bool {
	if vw == nil || modelData == nil {
		return false
	}
	if modelIndex < 0 || modelIndex >= len(vw.modelRendererLoads) {
		return false
	}
	state := &vw.modelRendererLoads[modelIndex]
	if !state.inProgress {
		return false
	}
	if state.sourceModel != modelData {
		return false
	}
	return state.hash == modelData.Hash()
}

// startModelRendererLoad はモデル描画の非同期読み込みを開始する。
func (vw *ViewerWindow) startModelRendererLoad(modelIndex int, modelData *model.PmxModel) {
	if vw == nil || modelData == nil {
		return
	}
	if modelIndex < 0 || modelIndex >= len(vw.modelRendererLoads) {
		return
	}
	vw.cancelModelRendererLoad(modelIndex)

	state := &vw.modelRendererLoads[modelIndex]
	state.token++
	state.hash = modelData.Hash()
	state.sourceModel = modelData
	state.inProgress = true

	ctx, cancel := context.WithCancel(context.Background())
	state.cancel = cancel
	token := state.token

	go func() {
		// モデル本体は不変とし、変形結果は別バッファへ出力するため複製しない。
		renderModel := modelData
		if ctx.Err() != nil {
			return
		}

		// UI応答性を確保するため1コア分は空ける。
		workerCount := max(1, runtime.GOMAXPROCS(0)-1)
		// 法線/選択/ボーンは必要時に生成するため、基本バッファのみ準備する。
		bufferData, err := render.PrepareModelRendererBufferDataWithOptions(
			ctx,
			renderModel,
			workerCount,
			render.ModelRendererBufferOptions{},
		)
		if ctx.Err() != nil {
			return
		}

		result := &modelRendererLoadResult{
			modelIndex:  modelIndex,
			token:       token,
			sourceModel: modelData,
			renderModel: renderModel,
			bufferData:  bufferData,
			err:         err,
		}
		select {
		case vw.loadResults <- result:
		default:
			logging.DefaultLogger().Warn("描画読み込み結果の通知に失敗しました: modelIndex=%d", modelIndex)
		}
	}()
}

// cancelModelRendererLoad は指定モデルの読み込みを中断する。
func (vw *ViewerWindow) cancelModelRendererLoad(modelIndex int) {
	if vw == nil {
		return
	}
	if modelIndex < 0 || modelIndex >= len(vw.modelRendererLoads) {
		return
	}
	state := &vw.modelRendererLoads[modelIndex]
	if state.cancel != nil {
		state.cancel()
		state.cancel = nil
	}
	state.token++
	state.hash = ""
	state.sourceModel = nil
	state.inProgress = false
}

// cancelAllModelRendererLoads は全モデルの読み込みを中断する。
func (vw *ViewerWindow) cancelAllModelRendererLoads() {
	if vw == nil {
		return
	}
	for i := range vw.modelRendererLoads {
		vw.cancelModelRendererLoad(i)
	}
}

// loadMotions はモーションを共有状態から同期する。
func (vw *ViewerWindow) loadMotions() {
	motionCount := vw.list.shared.MotionCount(vw.windowIndex)
	for i := 0; i < motionCount; i++ {
		if i >= len(vw.motions) {
			vw.motions = append(vw.motions, nil)
		}
		var motionData *motion.VmdMotion
		if raw := vw.list.shared.Motion(vw.windowIndex, i); raw != nil {
			if m, ok := raw.(*motion.VmdMotion); ok {
				motionData = m
			}
		}
		if motionData == nil {
			vw.motions[i] = nil
			continue
		}
		if vw.motions[i] == nil || vw.motions[i].Hash() != motionData.Hash() {
			vw.motions[i] = motionData
		}
	}
	vw.loadCameraMotion()
}

// loadCameraMotion はカメラ専用モーションを共有状態から同期する。
func (vw *ViewerWindow) loadCameraMotion() {
	if vw == nil || vw.list == nil || vw.list.shared == nil {
		return
	}
	if raw := vw.list.shared.CameraMotion(vw.windowIndex); raw != nil {
		if cameraMotion, ok := raw.(*motion.VmdMotion); ok && cameraMotion != nil {
			if vw.cameraMotion == nil || vw.cameraMotion.Hash() != cameraMotion.Hash() {
				vw.cameraMotion = cameraMotion
				vw.cameraManualOverride = false
			}
			return
		}
	}
	vw.cameraMotion = nil
	vw.cameraManualOverride = false
}

// applyCameraMotion は現在フレームのカメラモーションをカメラへ適用する。
func (vw *ViewerWindow) applyCameraMotion(frame motion.Frame) {
	if vw == nil || vw.shader == nil {
		return
	}
	cameraMotion := vw.resolveCameraMotion()
	if cameraMotion == nil || cameraMotion.CameraFrames == nil || cameraMotion.CameraFrames.Len() == 0 {
		return
	}
	if !vw.isPlaying() {
		vw.cameraManualOverride = false
	}
	if vw.cameraManualOverride {
		return
	}

	cf := cameraMotion.CameraFrames.Get(frame)
	if cf == nil || cf.Position == nil || cf.Degrees == nil {
		return
	}

	cam := vw.shader.Camera()
	if cam == nil {
		return
	}
	cam.SetByMotionValues(*cf.Position, *cf.Degrees, cf.Distance, cf.ViewOfAngle)
	vw.shader.SetCamera(cam)
}

// resolveCameraMotion はカメラフレームを持つモーションを返す。
func (vw *ViewerWindow) resolveCameraMotion() *motion.VmdMotion {
	if vw == nil {
		return nil
	}
	return vw.cameraMotion
}

// closeCallback はウィンドウ終了時の処理を行う。
func (vw *ViewerWindow) closeCallback(_ *glfw.Window) {
	vw.list.shared.SetClosed(true)
}

// keyCallback はキー入力を処理する。
func (vw *ViewerWindow) keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	logger := logging.DefaultLogger()
	verbose := logger != nil && logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER)
	if verbose && action != glfw.Release {
		if (key >= glfw.KeyKP0 && key <= glfw.KeyKP9) || key == glfw.KeyKPDecimal {
			actionLabel := "不明"
			switch action {
			case glfw.Press:
				actionLabel = "押下"
			case glfw.Release:
				actionLabel = "離し"
			case glfw.Repeat:
				actionLabel = "リピート"
			}
			logger.Verbose(
				logging.VERBOSE_INDEX_VIEWER,
				"操作: テンキー key=%d action=%s shift=%t ctrl=%t",
				key,
				actionLabel,
				vw.isShiftPressed(),
				vw.isCtrlPressed(),
			)
		}
	}
	switch action {
	case glfw.Press:
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			vw.shiftPressed = true
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			vw.ctrlPressed = true
			return
		}
	case glfw.Release:
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			vw.shiftPressed = false
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			vw.ctrlPressed = false
			return
		}
	}
	if action != glfw.Press {
		return
	}

	if preset, ok := cameraPresets[key]; ok {
		if vw.shouldSkipCameraManualOperation(messages.ViewerWindowKey001) {
			return
		}
		vw.resetCameraPositionForPreset(preset.Yaw, preset.Pitch)
		if verbose {
			cam := vw.shader.Camera()
			fov := float32(0)
			if cam != nil {
				fov = cam.FieldOfView
			}
			logger.Verbose(logging.VERBOSE_INDEX_VIEWER,
				"操作: カメラプリセット name=%s key=%d yaw=%.3f pitch=%.3f fov=%.3f",
				preset.Name,
				key,
				preset.Yaw,
				preset.Pitch,
				fov,
			)
		}
	}
}

// mouseCallback はマウスボタン入力を処理する。
func (vw *ViewerWindow) mouseCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	selectionEnabled := vw.isSelectionEnabledInWindow()
	switch action {
	case glfw.Press:
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = true
			if selectionEnabled && vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) &&
				vw.selectedVertexMode() == state.SELECTED_VERTEX_MODE_BOX {
				vw.boxSelectionDragging = true
				vw.boxSelectionStart = mmath.Vec2{X: vw.cursorX, Y: vw.cursorY}
				vw.boxSelectionEnd = vw.boxSelectionStart
			}
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = true
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = true
		}
	case glfw.Release:
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = false
			if !selectionEnabled {
				vw.boxSelectionDragging = false
				vw.boxSelectionPending = false
				return
			}
			if vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
				if vw.selectedVertexMode() == state.SELECTED_VERTEX_MODE_BOX {
					vw.boxSelectionDragging = false
					vw.boxSelectionPending = true
					vw.boxSelectionRemove = vw.isCtrlPressed()
					vw.boxSelectionEnd = mmath.Vec2{X: vw.cursorX, Y: vw.cursorY}
					return
				}
				vw.queueSelectedVertexSelection(vw.cursorX, vw.cursorY, vw.isCtrlPressed())
				return
			}
			drawRigidBody := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_RIGID_BODY_FRONT) ||
				vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_RIGID_BODY_BACK)
			drawJoint := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_JOINT)
			switch {
			case drawRigidBody:
				vw.selectRigidBodyByCursor(vw.cursorX, vw.cursorY)
			case drawJoint:
				vw.selectJointByCursor(vw.cursorX, vw.cursorY)
			case vw.list.shared.IsAnyBoneVisible():
				vw.selectBoneByCursor(vw.cursorX, vw.cursorY)
			default:
				vw.clearAllHovers()
			}
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = false
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = false
		}
	}
}

// cursorPosCallback はカーソル位置の更新を処理する。
func (vw *ViewerWindow) cursorPosCallback(_ *glfw.Window, xpos, ypos float64) {
	vw.cursorX = xpos
	vw.cursorY = ypos
	selectionEnabled := vw.isSelectionEnabledInWindow()
	nonVertexHoverEnabled := vw.isNonVertexHoverEnabledInWindow()

	if !vw.updatedPrevCursor {
		vw.prevCursorPos.X = xpos
		vw.prevCursorPos.Y = ypos
		vw.updatedPrevCursor = true
		return
	}
	if vw.rightButtonPressed {
		if !vw.shouldSkipCameraManualOperation(messages.ViewerWindowKey002) {
			vw.updateCameraAngleByCursor(xpos, ypos)
		}
	} else if vw.middleButtonPressed {
		if !vw.shouldSkipCameraManualOperation(messages.ViewerWindowKey003) {
			vw.updateCameraPositionByCursor(xpos, ypos)
		}
	}
	if selectionEnabled && vw.leftButtonPressed && vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
		if vw.selectedVertexMode() == state.SELECTED_VERTEX_MODE_BOX {
			if vw.boxSelectionDragging {
				vw.boxSelectionEnd = mmath.Vec2{X: xpos, Y: ypos}
			}
		} else {
			screenPos := mmath.Vec2{X: xpos, Y: ypos}
			if vw.isCtrlPressed() {
				if _, exists := vw.leftCursorRemoveWindowPositions[screenPos]; !exists {
					vw.leftCursorRemoveWindowOrder = append(vw.leftCursorRemoveWindowOrder, screenPos)
				}
				vw.leftCursorRemoveWindowPositions[screenPos] = 0
			} else {
				if _, exists := vw.leftCursorWindowPositions[screenPos]; !exists {
					vw.leftCursorWindowOrder = append(vw.leftCursorWindowOrder, screenPos)
				}
				vw.leftCursorWindowPositions[screenPos] = 0
			}
		}
	}
	if !vw.rightButtonPressed && !vw.middleButtonPressed {
		if !selectionEnabled && !nonVertexHoverEnabled {
			vw.prevCursorPos.X = xpos
			vw.prevCursorPos.Y = ypos
			return
		}
		if vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
			vw.clearNonVertexHovers()
			vw.prevCursorPos.X = xpos
			vw.prevCursorPos.Y = ypos
			return
		}
		if nonVertexHoverEnabled {
			drawRigidBody := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_RIGID_BODY_FRONT) ||
				vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_RIGID_BODY_BACK)
			drawJoint := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_JOINT)
			switch {
			case drawRigidBody:
				vw.selectRigidBodyByCursor(xpos, ypos)
				if drawJoint {
					if vw.rigidBodyHoverActive {
						vw.clearJointHover()
					} else {
						vw.selectJointByCursor(xpos, ypos)
					}
				} else {
					vw.clearJointHover()
				}
			case drawJoint:
				vw.selectJointByCursor(xpos, ypos)
			case vw.list.shared.IsAnyBoneVisible():
				vw.selectBoneByCursor(xpos, ypos)
			default:
				vw.clearAllHovers()
			}
		} else {
			vw.clearNonVertexHovers()
		}
	}
	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos
}

// updateCursorPositions は左クリックで記録したスクリーン座標をワールド座標に変換する。
func (vw *ViewerWindow) updateCursorPositions() ([]*mmath.Vec3, []float32, []*mmath.Vec3, []float32) {
	if vw == nil || vw.shader == nil || vw.list == nil || vw.list.shared == nil {
		return nil, nil, nil, nil
	}
	if !vw.isSelectionEnabledInWindow() {
		return nil, nil, nil, nil
	}
	if !vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
		return nil, nil, nil, nil
	}
	if vw.selectedVertexMode() != state.SELECTED_VERTEX_MODE_POINT {
		return nil, nil, nil, nil
	}
	leftCursorWorldPositions := make([]*mmath.Vec3, 0)
	leftCursorDepths := make([]float32, 0)
	leftCursorRemoveWorldPositions := make([]*mmath.Vec3, 0)
	leftCursorRemoveDepths := make([]float32, 0)

	winW, winH := vw.GetSize()
	fbW, fbH := vw.GetFramebufferSize()
	if fbW <= 0 || fbH <= 0 {
		return leftCursorWorldPositions, leftCursorDepths, leftCursorRemoveWorldPositions, leftCursorRemoveDepths
	}
	scaleX, scaleY := 1.0, 1.0
	if winW > 0 && winH > 0 {
		scaleX = float64(fbW) / float64(winW)
		scaleY = float64(fbH) / float64(winH)
	}

	for _, screenPos := range vw.leftCursorWindowOrder {
		depth, exists := vw.leftCursorWindowPositions[screenPos]
		if !exists || depth != 0 {
			continue
		}
		// 深度は最新の深度バッファを参照して確定する。
		x := float32(screenPos.X * scaleX)
		y := float32(screenPos.Y * scaleY)
		newDepth := vw.shader.Msaa().ReadDepthAt(int(x), int(y), fbW, fbH)
		if newDepth > 0.0 && newDepth < 1.0 {
			vw.leftCursorWindowPositions[screenPos] = newDepth
			worldPos := vw.getWorldPosition(x, y, newDepth, fbW, fbH)
			if worldPos != nil {
				leftCursorWorldPositions = append(leftCursorWorldPositions, worldPos)
				leftCursorDepths = append(leftCursorDepths, newDepth)
				vw.leftCursorWorldHistoryPositions = append(vw.leftCursorWorldHistoryPositions, worldPos)
			}
		}
	}
	for _, screenPos := range vw.leftCursorRemoveWindowOrder {
		depth, exists := vw.leftCursorRemoveWindowPositions[screenPos]
		if !exists || depth != 0 {
			continue
		}
		// 深度は最新の深度バッファを参照して確定する。
		x := float32(screenPos.X * scaleX)
		y := float32(screenPos.Y * scaleY)
		newDepth := vw.shader.Msaa().ReadDepthAt(int(x), int(y), fbW, fbH)
		if newDepth > 0.0 && newDepth < 1.0 {
			vw.leftCursorRemoveWindowPositions[screenPos] = newDepth
			worldPos := vw.getWorldPosition(x, y, newDepth, fbW, fbH)
			if worldPos != nil {
				leftCursorRemoveWorldPositions = append(leftCursorRemoveWorldPositions, worldPos)
				leftCursorRemoveDepths = append(leftCursorRemoveDepths, newDepth)
				vw.leftCursorRemoveWorldHistoryPositions = append(vw.leftCursorRemoveWorldHistoryPositions, worldPos)
			}
		}
	}

	return leftCursorWorldPositions, leftCursorDepths, leftCursorRemoveWorldPositions, leftCursorRemoveDepths
}

// flattenCursorPositions はワールド座標配列をGPU向け配列に変換する。
func flattenCursorPositions(positions []*mmath.Vec3, limit int) []float32 {
	if len(positions) == 0 {
		return nil
	}
	out := make([]float32, 0, len(positions)*3)
	for _, pos := range positions {
		if pos == nil {
			continue
		}
		out = append(out, float32(pos.X), float32(pos.Y), float32(pos.Z))
		if limit > 0 && len(out)/3 >= limit {
			break
		}
	}
	return out
}

// flattenCursorDepths は深度配列をGPU向けに上限付きで切り出す。
func flattenCursorDepths(depths []float32, limit int) []float32 {
	if len(depths) == 0 {
		return nil
	}
	count := len(depths)
	if limit > 0 && count > limit {
		count = limit
	}
	out := make([]float32, count)
	copy(out, depths[:count])
	return out
}

// queueSelectedVertexSelection は選択頂点の更新要求をキューに積む。
func (vw *ViewerWindow) queueSelectedVertexSelection(xpos, ypos float64, remove bool) {
	if !vw.isSelectionEnabledInWindow() {
		return
	}
	screenPos := mmath.Vec2{X: xpos, Y: ypos}
	if remove {
		if _, exists := vw.leftCursorRemoveWindowPositions[screenPos]; !exists {
			vw.leftCursorRemoveWindowOrder = append(vw.leftCursorRemoveWindowOrder, screenPos)
		}
		vw.leftCursorRemoveWindowPositions[screenPos] = 0
		return
	}
	if _, exists := vw.leftCursorWindowPositions[screenPos]; !exists {
		vw.leftCursorWindowOrder = append(vw.leftCursorWindowOrder, screenPos)
	}
	vw.leftCursorWindowPositions[screenPos] = 0
}

// updateBoneHoverExpire はボーンホバー終了後のツールチップ消去を遅延する。
func (vw *ViewerWindow) updateBoneHoverExpire() {
	if vw.boneHighlighter == nil {
		return
	}
	if vw.boneHoverActive {
		vw.lastBoneHoverAt = time.Now()
		return
	}
	if vw.lastBoneHoverAt.IsZero() {
		return
	}
	if time.Since(vw.lastBoneHoverAt) >= time.Second {
		vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
		vw.lastBoneHoverAt = time.Time{}
	}
}

// updateRigidBodyHoverExpire は剛体ホバー終了後のツールチップ消去を遅延する。
func (vw *ViewerWindow) updateRigidBodyHoverExpire() {
	if vw.rigidBodyHighlighter == nil {
		return
	}
	if vw.rigidBodyHoverActive {
		vw.lastRigidBodyHoverAt = time.Now()
		return
	}
	if vw.lastRigidBodyHoverAt.IsZero() {
		return
	}
	if time.Since(vw.lastRigidBodyHoverAt) >= time.Second {
		vw.rigidBodyHighlighter.UpdateDebugHoverByRigidBody(0, nil, false)
		vw.lastRigidBodyHoverAt = time.Time{}
	}
}

// updateJointHoverExpire はジョイントホバー終了後のツールチップ消去を遅延する。
func (vw *ViewerWindow) updateJointHoverExpire() {
	if vw.jointHoverActive {
		vw.lastJointHoverAt = time.Now()
		return
	}
	if vw.lastJointHoverAt.IsZero() {
		return
	}
	if time.Since(vw.lastJointHoverAt) >= time.Second {
		vw.clearJointHover()
	}
}

// updateSelectedVertexHoverExpire は選択頂点ホバー終了後のツールチップ消去を遅延する。
func (vw *ViewerWindow) updateSelectedVertexHoverExpire() {
	if vw.selectedVertexHoverActive {
		vw.lastSelectedVertexHoverAt = time.Now()
		return
	}
	if vw.lastSelectedVertexHoverAt.IsZero() {
		return
	}
	if time.Since(vw.lastSelectedVertexHoverAt) >= time.Second {
		vw.clearSelectedVertexHover()
	}
}

// updateSelectedVertexHover は選択頂点のホバー情報を更新する。
func (vw *ViewerWindow) updateSelectedVertexHover(modelIndex, hoverIndex int) {
	if hoverIndex < 0 {
		vw.selectedVertexHoverActive = false
		return
	}
	vw.selectedVertexHoverActive = true
	vw.selectedVertexHoverIndex = hoverIndex
	vw.selectedVertexHoverModelIndex = modelIndex
	vw.lastSelectedVertexHoverAt = time.Now()
}

// clearSelectedVertexHover は選択頂点のホバー情報をクリアする。
func (vw *ViewerWindow) clearSelectedVertexHover() {
	vw.selectedVertexHoverActive = false
	vw.selectedVertexHoverIndex = -1
	vw.selectedVertexHoverModelIndex = -1
	vw.lastSelectedVertexHoverAt = time.Time{}
}

// renderTooltip はホバー対象に応じてツールチップを描画する。
func (vw *ViewerWindow) renderTooltip(drawRigidBody, drawJoint bool, width, height int) {
	if vw.tooltipRenderer == nil {
		return
	}

	tooltipText := ""
	if vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) && vw.selectedVertexHoverActive && vw.selectedVertexHoverIndex >= 0 {
		tooltipText = fmt.Sprintf("頂点 %d", vw.selectedVertexHoverIndex)
	}
	if drawRigidBody && vw.rigidBodyHighlighter != nil {
		if hover := vw.rigidBodyHighlighter.DebugHoverInfo(); hover != nil && hover.RigidBody != nil {
			tooltipText = hover.RigidBody.Name()
		}
	}

	if tooltipText == "" && drawJoint && len(vw.jointHoverNames) > 0 {
		tooltipText = strings.Join(vw.jointHoverNames, ", ")
	}

	if tooltipText == "" && vw.list.shared.IsAnyBoneVisible() && vw.boneHighlighter != nil {
		if boneHover := vw.boneHighlighter.DebugBoneHoverInfo(); len(boneHover) > 0 {
			names := make([]string, 0, len(boneHover))
			for _, hover := range boneHover {
				if hover != nil && hover.Bone != nil {
					names = append(names, hover.Bone.Name())
				}
			}
			if len(names) > 0 {
				tooltipText = strings.Join(names, ", ")
			}
		}
	}

	if tooltipText == "" {
		return
	}

	vw.tooltipRenderer.Render(tooltipText, float32(vw.cursorX), float32(vw.cursorY), width, height)
}

// clearJointHover はジョイントのホバー情報をクリアする。
func (vw *ViewerWindow) clearJointHover() {
	vw.jointHoverNames = nil
	vw.jointHoverActive = false
	vw.lastJointHoverAt = time.Time{}
}

// clearAllHovers はボーン/剛体/ジョイントのホバー情報をクリアする。
func (vw *ViewerWindow) clearAllHovers() {
	vw.clearNonVertexHovers()
	vw.clearSelectedVertexHover()
}

// clearNonVertexHovers はボーン/剛体/ジョイントのホバー情報をクリアする。
func (vw *ViewerWindow) clearNonVertexHovers() {
	if vw.boneHighlighter != nil {
		vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
	}
	if vw.rigidBodyHighlighter != nil {
		vw.rigidBodyHighlighter.UpdateDebugHoverByRigidBody(0, nil, false)
	}
	vw.boneHoverActive = false
	vw.rigidBodyHoverActive = false
	vw.clearJointHover()
}

// calculateRayFromTo はカーソル位置からレイを生成する。
func (vw *ViewerWindow) calculateRayFromTo(xpos, ypos float64) (*mmath.Vec3, *mmath.Vec3) {
	windowWidth, windowHeight := vw.GetSize()
	framebufferWidth, framebufferHeight := vw.GetFramebufferSize()
	if framebufferWidth == 0 || framebufferHeight == 0 {
		return nil, nil
	}
	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return nil, nil
	}

	// GLFWのカーソル座標はウィンドウ座標なので、フレームバッファ座標に変換する。
	if windowWidth > 0 && windowHeight > 0 {
		scaleX := float64(framebufferWidth) / float64(windowWidth)
		scaleY := float64(framebufferHeight) / float64(windowHeight)
		xpos *= scaleX
		ypos *= scaleY
	}

	projection := toMgl32Mat4(cam.GetProjectionMatrix(framebufferWidth, framebufferHeight))
	view := toMgl32Mat4(cam.GetViewMatrix())

	nearWorld, errNear := mgl32.UnProject(
		mgl32.Vec3{float32(xpos), float32(framebufferHeight) - float32(ypos), 0.0},
		view, projection, 0, 0, framebufferWidth, framebufferHeight,
	)
	farWorld, errFar := mgl32.UnProject(
		mgl32.Vec3{float32(xpos), float32(framebufferHeight) - float32(ypos), 1.0},
		view, projection, 0, 0, framebufferWidth, framebufferHeight,
	)
	if errNear == nil && errFar == nil {
		rayFrom := mmath.Vec3{}
		rayFrom.X = float64(nearWorld.X())
		rayFrom.Y = float64(nearWorld.Y())
		rayFrom.Z = float64(nearWorld.Z())
		rayTo := mmath.Vec3{}
		rayTo.X = float64(farWorld.X())
		rayTo.Y = float64(farWorld.Y())
		rayTo.Z = float64(farWorld.Z())
		return &rayFrom, &rayTo
	}

	ndcX := (2.0*float64(xpos))/float64(framebufferWidth) - 1.0
	ndcY := 1.0 - (2.0*float64(ypos))/float64(framebufferHeight)
	aspect := float64(cam.AspectRatio)
	fovRad := mmath.DegToRad(float64(cam.FieldOfView))
	tanFov := math.Tan(fovRad * 0.5)

	dirCam := mmath.Vec3{}
	dirCam.X = ndcX * aspect * tanFov
	dirCam.Y = ndcY * tanFov
	dirCam.Z = -1.0
	dirCam = dirCam.Normalized()

	forward := cam.LookAtCenter.Subed(*cam.Position).Normalized()
	right := forward.Cross(*cam.Up).Normalized()
	up := right.Cross(forward).Normalized()

	dirWorld := mmath.Vec3{}
	dirWorld.X = dirCam.X*right.X + dirCam.Y*up.X + dirCam.Z*forward.X
	dirWorld.Y = dirCam.X*right.Y + dirCam.Y*up.Y + dirCam.Z*forward.Y
	dirWorld.Z = dirCam.X*right.Z + dirCam.Y*up.Z + dirCam.Z*forward.Z
	dirWorld = dirWorld.Normalized()

	from := cam.Position.Added(dirWorld.MuledScalar(float64(cam.NearPlane)))
	to := cam.Position.Added(dirWorld.MuledScalar(float64(cam.FarPlane)))
	return &from, &to
}

// cursorWorldPosition はスクリーン座標からワールド座標を求める。
func (vw *ViewerWindow) cursorWorldPosition(xpos, ypos float64) (*mmath.Vec3, bool) {
	if vw.shader == nil {
		return nil, false
	}
	windowWidth, windowHeight := vw.GetSize()
	framebufferWidth, framebufferHeight := vw.GetFramebufferSize()
	if framebufferWidth <= 0 || framebufferHeight <= 0 {
		return nil, false
	}
	if windowWidth > 0 && windowHeight > 0 {
		scaleX := float64(framebufferWidth) / float64(windowWidth)
		scaleY := float64(framebufferHeight) / float64(windowHeight)
		xpos *= scaleX
		ypos *= scaleY
	}
	// 深度は最新の深度バッファを参照して取得する。
	depth := vw.shader.Msaa().ReadDepthAt(int(xpos), int(ypos), framebufferWidth, framebufferHeight)
	if depth <= 0.0 || depth >= 1.0 {
		return nil, false
	}
	world := vw.getWorldPosition(float32(xpos), float32(ypos), depth, framebufferWidth, framebufferHeight)
	if world != nil {
		return world, true
	}
	return nil, false
}

// getWorldPosition はスクリーン座標と深度値からワールド座標を計算する。
func (vw *ViewerWindow) getWorldPosition(mouseX, mouseY, depth float32, width, height int) *mmath.Vec3 {
	projection, view, ok := vw.getCameraMatricesForGL(width, height)
	if !ok {
		return nil
	}

	ndcX := clampNormalized(2.0*mouseX/float32(width) - 1.0)
	ndcY := clampNormalized(1.0 - (2.0*mouseY)/float32(height))
	ndcZ := clampNormalized(depth*2.0 - 1.0)
	clip := mgl32.Vec4{ndcX, ndcY, ndcZ, 1.0}

	viewCoords := projection.Inv().Mul4x1(clip)
	viewPos := mgl32.Vec4{viewCoords.X(), viewCoords.Y(), viewCoords.Z(), 1.0}
	if viewCoords.W() != 0.0 {
		viewPos = mgl32.Vec4{
			viewCoords.X() / viewCoords.W(),
			viewCoords.Y() / viewCoords.W(),
			viewCoords.Z() / viewCoords.W(),
			1.0,
		}
	}

	worldCoords := view.Inv().Mul4x1(viewPos)
	worldPos := mgl32.Vec3{worldCoords.X(), worldCoords.Y(), worldCoords.Z()}
	if worldCoords.W() != 0.0 {
		worldPos = mgl32.Vec3{
			worldCoords.X() / worldCoords.W(),
			worldCoords.Y() / worldCoords.W(),
			worldCoords.Z() / worldCoords.W(),
		}
	}

	out := mmath.Vec3{}
	out.X = float64(worldPos.X())
	out.Y = float64(worldPos.Y())
	out.Z = float64(worldPos.Z())
	return &out
}

// getCameraMatricesForGL はOpenGL座標系の射影行列とビュー行列を返す。
func (vw *ViewerWindow) getCameraMatricesForGL(width, height int) (mgl32.Mat4, mgl32.Mat4, bool) {
	if vw == nil || vw.shader == nil || width <= 0 || height <= 0 {
		return mgl32.Mat4{}, mgl32.Mat4{}, false
	}
	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return mgl32.Mat4{}, mgl32.Mat4{}, false
	}

	projection := mgl32.Perspective(
		mgl32.DegToRad(cam.FieldOfView),
		float32(width)/float32(height),
		cam.NearPlane,
		cam.FarPlane,
	)
	view := mgl32.LookAtV(
		mgl.NewGlVec3(cam.Position),
		mgl.NewGlVec3(cam.LookAtCenter),
		mgl.NewGlVec3(cam.Up),
	)
	return projection, view, true
}

// clampNormalized は-1〜1にクランプする。
func clampNormalized(value float32) float32 {
	return min(max(value, -1.0), 1.0)
}

// buildCursorPositions はカーソル位置に沿ったサンプル点群を生成する。
func (vw *ViewerWindow) buildCursorPositions(xpos, ypos float64, limit int) []float32 {
	_ = limit
	worldPos, ok := vw.cursorWorldPosition(xpos, ypos)
	if !ok || worldPos == nil {
		return nil
	}
	// ワールド座標を点として渡す。
	positions := make([]float32, 0, 3)
	positions = append(positions, float32(worldPos.X), float32(worldPos.Y), float32(worldPos.Z))
	return positions
}

// selectedVertexMode は選択頂点のモードを返す。
func (vw *ViewerWindow) selectedVertexMode() state.SelectedVertexMode {
	if vw == nil || vw.list == nil || vw.list.shared == nil {
		return state.SELECTED_VERTEX_MODE_POINT
	}
	return vw.list.shared.SelectedVertexMode()
}

// resetBoxSelection はボックス選択状態をクリアする。
func (vw *ViewerWindow) resetBoxSelection() {
	vw.boxSelectionDragging = false
	vw.boxSelectionPending = false
	vw.boxSelectionRemove = false
}

// consumeBoxSelectionRect はボックス選択の矩形を取り出してクリアする。
func (vw *ViewerWindow) consumeBoxSelectionRect(winW, winH, fbW, fbH int) (mmath.Vec2, mmath.Vec2, bool, bool) {
	if vw == nil || !vw.boxSelectionPending || winW <= 0 || winH <= 0 || fbW <= 0 || fbH <= 0 {
		return mmath.Vec2{}, mmath.Vec2{}, false, false
	}
	vw.boxSelectionPending = false
	remove := vw.boxSelectionRemove
	vw.boxSelectionRemove = false

	scaleX := float64(fbW) / float64(winW)
	scaleY := float64(fbH) / float64(winH)
	startX := vw.boxSelectionStart.X * scaleX
	startY := vw.boxSelectionStart.Y * scaleY
	endX := vw.boxSelectionEnd.X * scaleX
	endY := vw.boxSelectionEnd.Y * scaleY

	minX := math.Min(startX, endX)
	maxX := math.Max(startX, endX)
	minY := math.Min(startY, endY)
	maxY := math.Max(startY, endY)

	minX = max(minX, 0)
	minY = max(minY, 0)
	maxX = min(maxX, float64(fbW))
	maxY = min(maxY, float64(fbH))

	return mmath.Vec2{X: minX, Y: minY}, mmath.Vec2{X: maxX, Y: maxY}, remove, true
}

// buildBoxSelectionLinePositions はボックス選択矩形のライン頂点をワールド座標で生成する。
func (vw *ViewerWindow) buildBoxSelectionLinePositions(winW, winH, fbW, fbH int) []float32 {
	if vw == nil || vw.shader == nil || winW <= 0 || winH <= 0 || fbW <= 0 || fbH <= 0 {
		return nil
	}
	scaleX := float64(fbW) / float64(winW)
	scaleY := float64(fbH) / float64(winH)
	startX := vw.boxSelectionStart.X * scaleX
	startY := vw.boxSelectionStart.Y * scaleY
	endX := vw.boxSelectionEnd.X * scaleX
	endY := vw.boxSelectionEnd.Y * scaleY

	minX := math.Min(startX, endX)
	maxX := math.Max(startX, endX)
	minY := math.Min(startY, endY)
	maxY := math.Max(startY, endY)

	minX = max(minX, 0)
	minY = max(minY, 0)
	maxX = min(maxX, float64(fbW))
	maxY = min(maxY, float64(fbH))
	if maxX-minX <= 0 || maxY-minY <= 0 {
		return nil
	}

	// 画面矩形が一致するように、近平面近傍のワールド座標へ変換する。
	// 近接しすぎると描画が消えることがあるため、少し奥に寄せる。
	depth := float32(0.01)
	p0 := vw.getWorldPosition(float32(minX), float32(minY), depth, fbW, fbH)
	p1 := vw.getWorldPosition(float32(maxX), float32(minY), depth, fbW, fbH)
	p2 := vw.getWorldPosition(float32(maxX), float32(maxY), depth, fbW, fbH)
	p3 := vw.getWorldPosition(float32(minX), float32(maxY), depth, fbW, fbH)
	if p0 == nil || p1 == nil || p2 == nil || p3 == nil {
		return nil
	}
	return []float32{
		float32(p0.X), float32(p0.Y), float32(p0.Z),
		float32(p1.X), float32(p1.Y), float32(p1.Z),
		float32(p2.X), float32(p2.Y), float32(p2.Z),
		float32(p3.X), float32(p3.Y), float32(p3.Z),
		float32(p0.X), float32(p0.Y), float32(p0.Z),
	}
}

// toMgl32Mat4 はmmath.Mat4をmgl32.Mat4に変換する。
func toMgl32Mat4(src mmath.Mat4) mgl32.Mat4 {
	var dst mgl32.Mat4
	for i := range src {
		dst[i] = float32(src[i])
	}
	return dst
}

// selectRigidBodyByCursor はカーソル位置に最も近い剛体を検出してハイライトを更新する。
func (vw *ViewerWindow) selectRigidBodyByCursor(xpos, ypos float64) {
	if vw.rigidBodyHighlighter == nil {
		vw.rigidBodyHoverActive = false
		return
	}
	// 物理レイキャストが有効な場合は最優先で使う。
	if vw.physics != nil {
		rayFrom, rayTo := vw.calculateRayFromTo(xpos, ypos)
		if rayFrom != nil && rayTo != nil {
			hit := vw.physics.RayTest(rayFrom, rayTo, &physics_api.RaycastFilter{
				Group: collisionAllFilterMask,
				Mask:  collisionAllFilterMask,
			})
			if hit != nil && hit.ModelIndex >= 0 && hit.ModelIndex < len(vw.modelRenderers) {
				renderer := vw.modelRenderers[hit.ModelIndex]
				if renderer != nil && renderer.Model != nil && renderer.Model.RigidBodies != nil {
					rigidBody, err := renderer.Model.RigidBodies.Get(hit.RigidBodyIndex)
					if err == nil && rigidBody != nil {
						var boneDeltas *delta.BoneDeltas
						if hit.ModelIndex < len(vw.vmdDeltas) && vw.vmdDeltas[hit.ModelIndex] != nil {
							boneDeltas = vw.vmdDeltas[hit.ModelIndex].Bones
						}
						worldMatrix := vw.rigidBodyWorldMatrix(rigidBody, renderer.Model.Bones, boneDeltas)
						vw.rigidBodyHighlighter.UpdateDebugHoverByRigidBodyWithMatrix(hit.ModelIndex, rigidBody, worldMatrix, true)
						vw.rigidBodyHoverActive = true
						vw.lastRigidBodyHoverAt = time.Now()
						return
					}
				}
			}
		}
	}

	// レイキャストで拾えない剛体はスクリーン距離で拾う。
	modelIndex, rigidBody, ok := vw.findRigidBodyByScreen(xpos, ypos)
	if !ok {
		vw.rigidBodyHoverActive = false
		return
	}
	var boneDeltas *delta.BoneDeltas
	if modelIndex < len(vw.vmdDeltas) && vw.vmdDeltas[modelIndex] != nil {
		boneDeltas = vw.vmdDeltas[modelIndex].Bones
	}
	var bones *model.BoneCollection
	if modelIndex < len(vw.modelRenderers) && vw.modelRenderers[modelIndex] != nil && vw.modelRenderers[modelIndex].Model != nil {
		bones = vw.modelRenderers[modelIndex].Model.Bones
	}
	worldMatrix := vw.rigidBodyWorldMatrix(rigidBody, bones, boneDeltas)
	vw.rigidBodyHighlighter.UpdateDebugHoverByRigidBodyWithMatrix(modelIndex, rigidBody, worldMatrix, true)
	vw.rigidBodyHoverActive = true
	vw.lastRigidBodyHoverAt = time.Now()
}

// findRigidBodyByScreen は画面投影による最近傍の剛体を探索する。
func (vw *ViewerWindow) findRigidBodyByScreen(xpos, ypos float64) (int, *model.RigidBody, bool) {
	if vw.shader == nil {
		return -1, nil, false
	}
	width, height := vw.GetSize()
	if width == 0 || height == 0 {
		return -1, nil, false
	}
	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return -1, nil, false
	}

	projection := mgl32.Perspective(
		mgl32.DegToRad(cam.FieldOfView),
		float32(width)/float32(height),
		cam.NearPlane,
		cam.FarPlane,
	)
	view := mgl32.LookAtV(
		mgl.NewGlVec3(cam.Position),
		mgl.NewGlVec3(cam.LookAtCenter),
		mgl.NewGlVec3(cam.Up),
	)
	forward := cam.LookAtCenter.Subed(*cam.Position).Normalized()
	right := forward.Cross(*cam.Up).Normalized()

	closestDistance := math.MaxFloat64
	closestModelIndex := -1
	var closestRigidBody *model.RigidBody

	for modelIndex, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil || renderer.Model.RigidBodies == nil {
			continue
		}
		var boneDeltas *delta.BoneDeltas
		if modelIndex < len(vw.vmdDeltas) && vw.vmdDeltas[modelIndex] != nil {
			boneDeltas = vw.vmdDeltas[modelIndex].Bones
		}
		for _, rigidBody := range renderer.Model.RigidBodies.Values() {
			if rigidBody == nil || !rigidBody.IsValid() {
				continue
			}
			pos, ok := vw.rigidBodyWorldPosition(rigidBody, renderer.Model.Bones, boneDeltas)
			if !ok {
				continue
			}
			screenX, screenY, ok := projectToScreen(pos, view, projection, width, height)
			if !ok {
				continue
			}
			pickRadius := rigidBodyScreenRadius(pos, rigidBody, right, view, projection, width, height, screenX, screenY)
			dist := math.Hypot(screenX-xpos, screenY-ypos)
			if dist > pickRadius {
				continue
			}
			if dist < closestDistance {
				closestDistance = dist
				closestModelIndex = modelIndex
				closestRigidBody = rigidBody
			}
		}
	}

	if closestRigidBody == nil {
		return -1, nil, false
	}
	return closestModelIndex, closestRigidBody, true
}

// rigidBodyWorldPosition は剛体のワールド位置を取得する。
func (vw *ViewerWindow) rigidBodyWorldPosition(
	rigidBody *model.RigidBody,
	bones *model.BoneCollection,
	boneDeltas *delta.BoneDeltas,
) (mmath.Vec3, bool) {
	if rigidBody == nil {
		return mmath.NewVec3(), false
	}
	if bones == nil || rigidBody.BoneIndex < 0 {
		if isInvalidViewerVec3(rigidBody.Position) {
			return mmath.NewVec3(), false
		}
		return rigidBody.Position, true
	}
	bone, err := bones.Get(rigidBody.BoneIndex)
	if err != nil || bone == nil {
		if isInvalidViewerVec3(rigidBody.Position) {
			return mmath.NewVec3(), false
		}
		return rigidBody.Position, true
	}
	if boneDeltas == nil || !boneDeltas.Contains(bone.Index()) {
		if isInvalidViewerVec3(rigidBody.Position) {
			return mmath.NewVec3(), false
		}
		return rigidBody.Position, true
	}
	boneDelta := boneDeltas.Get(bone.Index())
	if boneDelta == nil {
		if isInvalidViewerVec3(rigidBody.Position) {
			return mmath.NewVec3(), false
		}
		return rigidBody.Position, true
	}

	// ボーンのグローバル行列にローカルオフセットを掛けて剛体位置を推定する。
	localOffset := rigidBody.Position.Subed(bone.Position)
	pos := boneDelta.FilledGlobalMatrix().MulVec3(localOffset)
	if isInvalidViewerVec3(pos) {
		return mmath.NewVec3(), false
	}
	return pos, true
}

// rigidBodyWorldMatrix は剛体のワールド行列を推定する。
func (vw *ViewerWindow) rigidBodyWorldMatrix(
	rigidBody *model.RigidBody,
	bones *model.BoneCollection,
	boneDeltas *delta.BoneDeltas,
) *mmath.Mat4 {
	if rigidBody == nil {
		return nil
	}

	localPos := rigidBody.Position
	var boneDelta *delta.BoneDelta
	if bones != nil && rigidBody.BoneIndex >= 0 {
		bone, err := bones.Get(rigidBody.BoneIndex)
		if err == nil && bone != nil {
			localPos = rigidBody.Position.Subed(bone.Position)
			if boneDeltas != nil && boneDeltas.Contains(bone.Index()) {
				boneDelta = boneDeltas.Get(bone.Index())
			}
		}
	}

	localRot := rigidBody.Rotation.RadToQuaternion().ToMat4()
	localRot.Translate(localPos)

	if boneDelta == nil {
		local := localRot
		return &local
	}
	world := boneDelta.FilledGlobalMatrix().Muled(localRot)
	return &world
}

// rigidBodyScreenRadius は剛体の画面上半径を推定する。
func rigidBodyScreenRadius(
	center mmath.Vec3,
	rigidBody *model.RigidBody,
	right mmath.Vec3,
	view, projection mgl32.Mat4,
	width, height int,
	screenX, screenY float64,
) float64 {
	if rigidBody == nil || width == 0 || height == 0 {
		return boneHoverMaxScreenDistance
	}
	radius := rigidBodyBaseRadius(rigidBody) * rigidBodyHoverRadiusScale
	if radius <= 0 {
		return boneHoverMaxScreenDistance
	}
	edgePos := center.Added(right.MuledScalar(radius))
	edgeX, edgeY, ok := projectToScreen(edgePos, view, projection, width, height)
	if !ok {
		return boneHoverMaxScreenDistance
	}
	pixelRadius := math.Hypot(edgeX-screenX, edgeY-screenY)
	if pixelRadius < rigidBodyHoverMinPixelRadius {
		pixelRadius = rigidBodyHoverMinPixelRadius
	}
	return pixelRadius
}

// rigidBodyBaseRadius は剛体サイズから基準半径を計算する。
func rigidBodyBaseRadius(rigidBody *model.RigidBody) float64 {
	if rigidBody == nil {
		return 0
	}
	size := rigidBody.Size.Absed()
	switch rigidBody.Shape {
	case model.SHAPE_SPHERE:
		return size.X
	case model.SHAPE_CAPSULE:
		return size.X + size.Y*0.5
	case model.SHAPE_BOX:
		return math.Sqrt(size.X*size.X + size.Y*size.Y + size.Z*size.Z)
	default:
		return math.Sqrt(size.X*size.X + size.Y*size.Y + size.Z*size.Z)
	}
}

// isInvalidViewerVec3 はNaN/Infを含むベクトルか判定する。
func isInvalidViewerVec3(v mmath.Vec3) bool {
	return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z) ||
		math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) || math.IsInf(v.Z, 0)
}

// selectJointByCursor はカーソル位置に最も近いジョイントを検出してホバー情報を更新する。
func (vw *ViewerWindow) selectJointByCursor(xpos, ypos float64) {
	if vw.shader == nil {
		vw.jointHoverActive = false
		return
	}
	width, height := vw.GetSize()
	if width == 0 || height == 0 {
		return
	}
	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return
	}
	projection := mgl32.Perspective(
		mgl32.DegToRad(cam.FieldOfView),
		float32(width)/float32(height),
		cam.NearPlane,
		cam.FarPlane,
	)
	view := mgl32.LookAtV(
		mgl.NewGlVec3(cam.Position),
		mgl.NewGlVec3(cam.LookAtCenter),
		mgl.NewGlVec3(cam.Up),
	)

	closestDistance := math.MaxFloat64
	closestNames := make([]string, 0)
	for _, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil || renderer.Model.Joints == nil {
			continue
		}
		for _, joint := range renderer.Model.Joints.Values() {
			if joint == nil {
				continue
			}
			pos := joint.Param.Position
			screenX, screenY, ok := projectToScreen(pos, view, projection, width, height)
			if !ok {
				continue
			}
			dist := math.Hypot(screenX-xpos, screenY-ypos)
			if dist > boneHoverMaxScreenDistance {
				continue
			}
			if dist+0.01 < closestDistance {
				closestDistance = dist
				closestNames = closestNames[:0]
			}
			if math.Abs(dist-closestDistance) <= 0.01 {
				closestNames = append(closestNames, joint.Name())
			}
		}
	}
	if len(closestNames) > 0 {
		vw.jointHoverNames = closestNames
		vw.jointHoverActive = true
		vw.lastJointHoverAt = time.Now()
	} else {
		vw.jointHoverActive = false
	}
}

// selectBoneByCursor はカーソル位置に最も近いボーンを検出してハイライトを更新する。
func (vw *ViewerWindow) selectBoneByCursor(xpos, ypos float64) {
	if vw.boneHighlighter == nil || vw.shader == nil || vw.shader.Msaa() == nil {
		return
	}
	if len(vw.vmdDeltas) == 0 || vw.vmdDeltas[0] == nil {
		return
	}

	width, height := vw.GetSize()
	if width == 0 || height == 0 {
		return
	}
	if xpos < 0 || ypos < 0 || xpos > float64(width) || ypos > float64(height) {
		return
	}

	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return
	}

	// 描画と同じ行列でスクリーン座標→ワールド座標を復元する。
	projection := mgl32.Perspective(
		mgl32.DegToRad(cam.FieldOfView),
		float32(width)/float32(height),
		cam.NearPlane,
		cam.FarPlane,
	)
	view := mgl32.LookAtV(
		mgl.NewGlVec3(cam.Position),
		mgl.NewGlVec3(cam.LookAtCenter),
		mgl.NewGlVec3(cam.Up),
	)

	closestDistance := math.MaxFloat64
	closestBones := make([]*graphics_api.DebugBoneHover, 0)

	for modelIndex, vmdDeltas := range vw.vmdDeltas {
		if vmdDeltas == nil || vmdDeltas.Bones == nil {
			continue
		}
		info := buildBoneDebugInfo(vmdDeltas.Bones)
		vmdDeltas.Bones.ForEach(func(_ int, boneDelta *delta.BoneDelta) bool {
			if boneDelta == nil || boneDelta.Bone == nil {
				return true
			}
			if !isBoneVisibleForHover(boneDelta.Bone, vw.list.shared, info) {
				return true
			}
			bonePos := boneDelta.FilledGlobalPosition()
			screenX, screenY, ok := projectToScreen(bonePos, view, projection, width, height)
			if !ok {
				return true
			}
			dx := screenX - xpos
			dy := screenY - ypos
			dist := math.Hypot(dx, dy)
			if dist > boneHoverMaxScreenDistance {
				return true
			}
			if dist+0.01 < closestDistance {
				closestDistance = dist
				closestBones = closestBones[:0]
			}
			if math.Abs(dist-closestDistance) <= 0.01 {
				closestBones = append(closestBones, &graphics_api.DebugBoneHover{
					ModelIndex: modelIndex,
					Bone:       boneDelta.Bone,
					Distance:   dist,
				})
			}
			return true
		})
	}

	if len(closestBones) > 0 {
		vw.boneHighlighter.UpdateDebugHoverByBones(closestBones, true)
		vw.boneHoverActive = true
		vw.lastBoneHoverAt = time.Now()
	} else {
		vw.boneHoverActive = false
	}
}

// boneDebugInfo はボーン可視判定用の情報を保持する。
type boneDebugInfo struct {
	ikTargets       map[int]struct{}
	ikLinks         map[int]struct{}
	effectorParents map[int]struct{}
}

// isIkTarget はIKターゲットか判定する。
func (info boneDebugInfo) isIkTarget(index int) bool {
	_, ok := info.ikTargets[index]
	return ok
}

// isIkLink はIKリンクか判定する。
func (info boneDebugInfo) isIkLink(index int) bool {
	_, ok := info.ikLinks[index]
	return ok
}

// isEffectorParent は付与親ボーンか判定する。
func (info boneDebugInfo) isEffectorParent(index int) bool {
	_, ok := info.effectorParents[index]
	return ok
}

// buildBoneDebugInfo はIK/付与関係の判定情報を構築する。
func buildBoneDebugInfo(bones *delta.BoneDeltas) boneDebugInfo {
	info := boneDebugInfo{
		ikTargets:       map[int]struct{}{},
		ikLinks:         map[int]struct{}{},
		effectorParents: map[int]struct{}{},
	}
	if bones == nil {
		return info
	}
	bones.ForEach(func(_ int, delta *delta.BoneDelta) bool {
		if delta == nil || delta.Bone == nil {
			return true
		}
		bone := delta.Bone
		if bone.Ik != nil {
			if bone.Ik.BoneIndex >= 0 {
				info.ikTargets[bone.Ik.BoneIndex] = struct{}{}
			}
			for _, link := range bone.Ik.Links {
				if link.BoneIndex >= 0 {
					info.ikLinks[link.BoneIndex] = struct{}{}
				}
			}
		}
		if hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_ROTATION) || hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_TRANSLATION) {
			if bone.EffectIndex >= 0 {
				info.effectorParents[bone.EffectIndex] = struct{}{}
			}
		}
		return true
	})
	return info
}

// isBoneVisibleForHover は状態フラグに応じてホバー対象か判定する。
func isBoneVisibleForHover(bone *model.Bone, shared *state.SharedState, info boneDebugInfo) bool {
	if bone == nil || shared == nil {
		return false
	}
	showAll := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_ALL)
	showVisible := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_VISIBLE)
	showIk := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_IK)
	showEffector := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_EFFECTOR)
	showFixed := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_FIXED)
	showRotate := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_ROTATE)
	showTranslate := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_TRANSLATE)

	switch {
	case (showAll || showVisible || showIk) && hasBoneFlag(bone, model.BONE_FLAG_IS_IK):
		return true
	case (showAll || showVisible || showIk) && info.isIkLink(bone.Index()):
		return true
	case (showAll || showVisible || showIk) && info.isIkTarget(bone.Index()):
		return true
	case (showAll || showVisible || showEffector) &&
		(hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_ROTATION) || hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_TRANSLATION)):
		return true
	case (showAll || showVisible || showEffector) && info.isEffectorParent(bone.Index()):
		return true
	case (showAll || showVisible || showFixed) && hasBoneFlag(bone, model.BONE_FLAG_HAS_FIXED_AXIS):
		return true
	case (showAll || showVisible || showTranslate) && hasBoneFlag(bone, model.BONE_FLAG_CAN_TRANSLATE):
		return true
	case (showAll || showVisible || showRotate) && hasBoneFlag(bone, model.BONE_FLAG_CAN_ROTATE):
		return true
	case showAll && !hasBoneFlag(bone, model.BONE_FLAG_IS_VISIBLE):
		return true
	default:
		return false
	}
}

// projectToScreen はワールド座標をスクリーン座標へ変換する。
func projectToScreen(pos mmath.Vec3, view, projection mgl32.Mat4, width, height int) (float64, float64, bool) {
	if width == 0 || height == 0 {
		return 0, 0, false
	}
	glPos := mgl.NewGlVec3(&pos)
	clip := projection.Mul4(view).Mul4x1(mgl32.Vec4{glPos.X(), glPos.Y(), glPos.Z(), 1})
	if clip.W() == 0 {
		return 0, 0, false
	}
	ndc := clip.Mul(1.0 / clip.W())
	if ndc.Z() < -1.0 || ndc.Z() > 1.0 {
		return 0, 0, false
	}
	screenX := (float64(ndc.X()) + 1.0) * 0.5 * float64(width)
	screenY := (1.0 - float64(ndc.Y())) * 0.5 * float64(height)
	if screenX < 0 || screenY < 0 || screenX > float64(width) || screenY > float64(height) {
		return 0, 0, false
	}
	return screenX, screenY, true
}

// hasBoneFlag はボーンフラグの有無を判定する。
func hasBoneFlag(bone *model.Bone, flag model.BoneFlag) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&flag != 0
}

// modifierState は修飾キーの状態と倍率を返す。
func (vw *ViewerWindow) modifierState(base float64) (float64, bool, bool) {
	shift := vw.isShiftPressed()
	ctrl := vw.isCtrlPressed()
	ratio := base
	if shift {
		ratio *= 10.0
	} else if ctrl {
		ratio *= 0.1
	}
	return ratio, shift, ctrl
}

// isShiftPressed はShiftキーが押されているか判定する。
func (vw *ViewerWindow) isShiftPressed() bool {
	if vw == nil {
		return false
	}
	if vw.shiftPressed {
		return true
	}
	if vw.Window == nil {
		return false
	}
	return vw.Window.GetKey(glfw.KeyLeftShift) == glfw.Press ||
		vw.Window.GetKey(glfw.KeyRightShift) == glfw.Press
}

// isCtrlPressed はCtrlキーが押されているか判定する。
func (vw *ViewerWindow) isCtrlPressed() bool {
	if vw == nil {
		return false
	}
	if vw.ctrlPressed {
		return true
	}
	if vw.Window == nil {
		return false
	}
	return vw.Window.GetKey(glfw.KeyLeftControl) == glfw.Press ||
		vw.Window.GetKey(glfw.KeyRightControl) == glfw.Press
}

// updateCameraAngleByCursor はカーソル移動でカメラ角度を更新する。
func (vw *ViewerWindow) updateCameraAngleByCursor(xpos, ypos float64) {
	ratio, _, _ := vw.modifierState(0.1)
	xOffset := (xpos - vw.prevCursorPos.X) * ratio
	yOffset := (ypos - vw.prevCursorPos.Y) * ratio
	cam := vw.shader.Camera()
	if cam == nil {
		return
	}
	vw.markCameraManualOverride()
	cam.RotateOrbit(xOffset, yOffset)
	vw.shader.SetCamera(cam)
	vw.syncCameraToOthers()
}

// updateCameraPositionByCursor はカーソル移動でカメラ位置を更新する。
func (vw *ViewerWindow) updateCameraPositionByCursor(xpos, ypos float64) {
	ratio, _, _ := vw.modifierState(0.07)
	xOffset := (vw.prevCursorPos.X - xpos) * ratio
	yOffset := (vw.prevCursorPos.Y - ypos) * ratio

	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil {
		return
	}
	vw.markCameraManualOverride()

	right := cam.RightVector()
	up := cam.UpVector()

	upMovement := up.MuledScalar(-yOffset)
	rightMovement := right.MuledScalar(xOffset)
	movement := upMovement.Added(rightMovement)

	cam.Position.Add(movement)
	cam.LookAtCenter.Add(movement)
	vw.shader.SetCamera(cam)
	vw.syncCameraToOthers()
}

// syncCameraToOthers はカメラ状態を他ウィンドウへ同期する。
func (vw *ViewerWindow) syncCameraToOthers() {
	if !vw.list.shared.HasFlag(state.STATE_FLAG_CAMERA_SYNC) {
		return
	}
	currentCam := vw.shader.Camera()
	for _, other := range vw.list.windowList {
		if other.windowIndex == vw.windowIndex {
			continue
		}
		otherCam := other.shader.Camera()
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
		otherCam.Orientation = currentCam.Orientation
		other.shader.SetCamera(otherCam)
	}
}

// scrollCallback はホイール操作でカメラのFOVを調整する。
func (vw *ViewerWindow) scrollCallback(_ *glfw.Window, _ float64, yoff float64) {
	stepRatio, _, _ := vw.modifierState(1.0)
	step := float32(stepRatio)
	cam := vw.shader.Camera()
	if cam == nil {
		return
	}
	if vw.shouldSkipCameraManualOperation(messages.ViewerWindowKey004) {
		return
	}
	vw.markCameraManualOverride()
	if yoff > 0 {
		cam.FieldOfView -= step
		if cam.FieldOfView < 1.0 {
			cam.FieldOfView = 1.0
		}
	} else if yoff < 0 {
		cam.FieldOfView += step
	}
	vw.shader.SetCamera(cam)
	vw.syncCameraToOthers()
}

// isPlaying は再生中か判定する。
func (vw *ViewerWindow) isPlaying() bool {
	if vw == nil || vw.list == nil || vw.list.shared == nil {
		return false
	}
	return vw.list.shared.HasFlag(state.STATE_FLAG_PLAYING)
}

// markCameraManualOverride は再生中の手動カメラ操作でモーション上書きを抑止する。
func (vw *ViewerWindow) markCameraManualOverride() {
	if vw == nil || !vw.isPlaying() {
		return
	}
	cameraMotion := vw.resolveCameraMotion()
	if cameraMotion == nil || cameraMotion.CameraFrames == nil || cameraMotion.CameraFrames.Len() == 0 {
		return
	}
	vw.cameraManualOverride = true
}

// shouldSkipCameraManualOperation は再生中カメラモーション適用時に手動カメラ操作を抑止する。
func (vw *ViewerWindow) shouldSkipCameraManualOperation(operationKey string) bool {
	if vw == nil || !vw.isPlaying() {
		return false
	}
	cameraMotion := vw.resolveCameraMotion()
	if cameraMotion == nil || cameraMotion.CameraFrames == nil || cameraMotion.CameraFrames.Len() == 0 {
		return false
	}
	now := time.Now()
	if vw.lastCameraBlockedWarnAt.IsZero() || now.Sub(vw.lastCameraBlockedWarnAt) >= cameraOperationBlockedWarnCooldown {
		vw.lastCameraBlockedWarnAt = now
		operation := i18n.T(operationKey)
		warnMessage := i18n.T(messages.ViewerWindowKey005)
		logging.DefaultLogger().Warn(
			warnMessage,
			operation,
		)
	}
	return true
}

// focusCallback はフォーカス連動の通知を行う。
func (vw *ViewerWindow) focusCallback(_ *glfw.Window, focused bool) {
	if !vw.list.shared.IsFocusLinkEnabled() {
		return
	}
	if focused {
		vw.list.shared.TriggerLinkedFocus(vw.windowIndex)
	}
}

// sizeCallback はオーバーレイ表示時にサイズを同期する。
func (vw *ViewerWindow) sizeCallback(_ *glfw.Window, width, height int) {
	if !vw.list.shared.IsShowOverride() {
		return
	}
	for _, other := range vw.list.windowList {
		if other.windowIndex == vw.windowIndex {
			continue
		}
		other.SetSize(width, height)
	}
}

// iconifyCallback は最小化/復帰状態を同期する。
func (vw *ViewerWindow) iconifyCallback(_ *glfw.Window, iconified bool) {
	if iconified {
		vw.list.shared.SyncMinimize(vw.windowIndex)
	} else {
		vw.list.shared.SyncRestore(vw.windowIndex)
	}
}
