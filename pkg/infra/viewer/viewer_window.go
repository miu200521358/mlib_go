//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"fmt"
	"image"
	"math"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/physics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
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
)

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
	physics     physics_api.IPhysics

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

	modelRenderers     []*render.ModelRenderer
	motions            []*motion.VmdMotion
	vmdDeltas          []*delta.VmdDeltas
	physicsModelHashes []string

	leftButtonPressed            bool
	middleButtonPressed          bool
	rightButtonPressed           bool
	shiftPressed                 bool
	ctrlPressed                  bool
	updatedPrevCursor            bool
	prevCursorPos                mmath.Vec2
	cursorX                      float64
	cursorY                      float64
	leftCursorWindowPositions             map[mmath.Vec2]float32
	leftCursorRemoveWindowPositions       map[mmath.Vec2]float32
	leftCursorWorldHistoryPositions       []*mmath.Vec3
	leftCursorRemoveWorldHistoryPositions []*mmath.Vec3
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
	physics := physics_api.NewPhysics(&gravity)

	vw := &ViewerWindow{
		Window:                        glWindow,
		windowIndex:                   windowIndex,
		title:                         title,
		list:                          list,
		shader:                        shader,
		physics:                       physics,
		tooltipRenderer:               tooltipRenderer,
		boneHighlighter:               mgl.NewBoneHighlighter(),
		rigidBodyHighlighter:          mgl.NewRigidBodyHighlighter(),
		modelRenderers:                make([]*render.ModelRenderer, 0),
		motions:                       make([]*motion.VmdMotion, 0),
		vmdDeltas:                     make([]*delta.VmdDeltas, 0),
		physicsModelHashes:            make([]string, 0),
		selectedVertexHoverIndex:      -1,
		selectedVertexHoverModelIndex: -1,
		leftCursorWindowPositions:             make(map[mmath.Vec2]float32),
		leftCursorRemoveWindowPositions:       make(map[mmath.Vec2]float32),
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
	cam.ResetPosition(yaw, pitch)
	if cam.FieldOfView == 0 {
		cam.FieldOfView = graphics_api.FieldOfViewAngle
	}
	vw.syncCameraToOthers()
}

// resetCameraPositionForPreset はカメラの視点とFOVを既定値へ戻して同期する。
func (vw *ViewerWindow) resetCameraPositionForPreset(yaw, pitch float64) {
	cam := vw.shader.Camera()
	if cam == nil {
		return
	}
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
	vw.MakeContextCurrent()
	vw.shader.Resize(fbW, fbH)

	showSelectedVertex := vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX)
	var leftCursorWorldPositions []*mmath.Vec3
	var leftCursorRemoveWorldPositions []*mmath.Vec3
	if showSelectedVertex {
		leftCursorWorldPositions, leftCursorRemoveWorldPositions = vw.updateCursorPositions()
	}

	limit := 0
	if vw.list.appConfig != nil {
		limit = vw.list.appConfig.CursorPositionLimit
	}
	var hoverCursorPositions []float32
	if showSelectedVertex {
		hoverCursorPositions = vw.buildCursorPositions(vw.cursorX, vw.cursorY, limit)
	}
	selectedCursorPositions := flattenCursorPositions(leftCursorWorldPositions, limit)
	removeSelectedCursorPositions := flattenCursorPositions(leftCursorRemoveWorldPositions, limit)
	applySelection := len(selectedCursorPositions) > 0 || len(removeSelectedCursorPositions) > 0
	removeSelection := len(removeSelectedCursorPositions) > 0

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

	vw.shader.UpdateCamera()

	vw.renderFloor()

	vw.updateSelectedVertexHoverExpire()
	vw.updateBoneHoverExpire()
	vw.updateRigidBodyHoverExpire()
	vw.updateJointHoverExpire()
	var debugBoneHover []*graphics_api.DebugBoneHover
	if vw.boneHighlighter != nil {
		debugBoneHover = vw.boneHighlighter.DebugBoneHoverInfo()
	}

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
		selectedVertexIndexes := []int{}
		if showSelectedVertex {
			selectedVertexIndexes = vw.list.shared.SelectedVertexIndexes(vw.windowIndex, i)
		}
		cursorPositions := hoverCursorPositions
		removeCursorPositions := []float32(nil)
		applySelectionForModel := false
		if showSelectedVertex && i == 0 && applySelection {
			if removeSelection {
				cursorPositions = removeSelectedCursorPositions
				removeCursorPositions = cursorPositions
			} else {
				cursorPositions = selectedCursorPositions
			}
			applySelectionForModel = true
		}
		updatedSelected, hoverIndex := renderer.Render(
			vw.shader,
			vw.list.shared,
			vmdDeltas,
			debugBoneHover,
			selectedVertexIndexes,
			cursorPositions,
			removeCursorPositions,
			applySelectionForModel,
		)
		if applySelectionForModel {
			vw.list.shared.SetSelectedVertexIndexes(vw.windowIndex, i, updatedSelected)
		}
		if showSelectedVertex && i == 0 {
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

	vw.SwapBuffers()
	if showSelectedVertex && !vw.leftButtonPressed {
		vw.leftCursorWindowPositions = make(map[mmath.Vec2]float32)
		vw.leftCursorRemoveWindowPositions = make(map[mmath.Vec2]float32)
		vw.leftCursorWorldHistoryPositions = make([]*mmath.Vec3, 0)
		vw.leftCursorRemoveWorldHistoryPositions = make([]*mmath.Vec3, 0)
	}
}

// cleanupResources はビューワーが保持するOpenGLリソースを解放する。
func (vw *ViewerWindow) cleanupResources() {
	if vw == nil {
		return
	}
	if vw.Window != nil {
		vw.MakeContextCurrent()
	}

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

// prepareFrame は描画前にモデル/モーションを同期する。
func (vw *ViewerWindow) prepareFrame() {
	vw.MakeContextCurrent()
	vw.loadModelRenderers()
	vw.loadMotions()
	vw.ensurePhysicsModelSlots()
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

// loadModelRenderers はモデルレンダラを共有状態から同期する。
func (vw *ViewerWindow) loadModelRenderers() {
	modelCount := vw.list.shared.ModelCount(vw.windowIndex)
	motionCount := vw.list.shared.MotionCount(vw.windowIndex)
	count := max(modelCount, motionCount)
	for i := 0; i < count; i++ {
		if i >= len(vw.modelRenderers) {
			vw.modelRenderers = append(vw.modelRenderers, nil)
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
			continue
		}
		if existing == nil || existing.Hash() != modelData.Hash() || existing.Model != modelData {
			if existing != nil {
				existing.Delete()
			}
			vw.modelRenderers[i] = render.NewModelRenderer(vw.windowIndex, modelData)
			vw.vmdDeltas[i] = nil
		}
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

	if preset, ok := cameraPresets[key]; ok {
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
	switch action {
	case glfw.Press:
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = true
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = true
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = true
		}
	case glfw.Release:
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = false
			if vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
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

	if !vw.updatedPrevCursor {
		vw.prevCursorPos.X = xpos
		vw.prevCursorPos.Y = ypos
		vw.updatedPrevCursor = true
		return
	}
	if vw.rightButtonPressed {
		vw.updateCameraAngleByCursor(xpos, ypos)
	} else if vw.middleButtonPressed {
		vw.updateCameraPositionByCursor(xpos, ypos)
	}
	if vw.leftButtonPressed && vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
		screenPos := mmath.Vec2{X: xpos, Y: ypos}
		if vw.isCtrlPressed() {
			vw.leftCursorRemoveWindowPositions[screenPos] = 0
		} else {
			vw.leftCursorWindowPositions[screenPos] = 0
		}
	}
	if !vw.rightButtonPressed && !vw.middleButtonPressed {
		if vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
			vw.clearNonVertexHovers()
			vw.prevCursorPos.X = xpos
			vw.prevCursorPos.Y = ypos
			return
		}
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
	}
	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos
}

// updateCursorPositions は左クリックで記録したスクリーン座標をワールド座標に変換する。
func (vw *ViewerWindow) updateCursorPositions() ([]*mmath.Vec3, []*mmath.Vec3) {
	if vw == nil || vw.shader == nil || vw.list == nil || vw.list.shared == nil {
		return nil, nil
	}
	if !vw.list.shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
		return nil, nil
	}
	leftCursorWorldPositions := make([]*mmath.Vec3, 0)
	leftCursorRemoveWorldPositions := make([]*mmath.Vec3, 0)

	winW, winH := vw.GetSize()
	fbW, fbH := vw.GetFramebufferSize()
	if fbW <= 0 || fbH <= 0 {
		return leftCursorWorldPositions, leftCursorRemoveWorldPositions
	}
	scaleX, scaleY := 1.0, 1.0
	if winW > 0 && winH > 0 {
		scaleX = float64(fbW) / float64(winW)
		scaleY = float64(fbH) / float64(winH)
	}

	for screenPos, depth := range vw.leftCursorWindowPositions {
		if depth != 0 {
			continue
		}
		for range 20 {
			x := float32(screenPos.X * scaleX)
			y := float32(screenPos.Y * scaleY)
			newDepth := vw.shader.Msaa().ReadDepthAt(int(x), int(y), fbW, fbH)
			if newDepth > 0.0 && newDepth < 1.0 {
				vw.leftCursorWindowPositions[screenPos] = newDepth
				worldPos := vw.getWorldPosition(x, y, newDepth, fbW, fbH)
				if worldPos != nil {
					leftCursorWorldPositions = append(leftCursorWorldPositions, worldPos)
					vw.leftCursorWorldHistoryPositions = append(vw.leftCursorWorldHistoryPositions, worldPos)
				}
				break
			}
		}
	}
	for screenPos, depth := range vw.leftCursorRemoveWindowPositions {
		if depth != 0 {
			continue
		}
		for range 20 {
			x := float32(screenPos.X * scaleX)
			y := float32(screenPos.Y * scaleY)
			newDepth := vw.shader.Msaa().ReadDepthAt(int(x), int(y), fbW, fbH)
			if newDepth > 0.0 && newDepth < 1.0 {
				vw.leftCursorRemoveWindowPositions[screenPos] = newDepth
				worldPos := vw.getWorldPosition(x, y, newDepth, fbW, fbH)
				if worldPos != nil {
					leftCursorRemoveWorldPositions = append(leftCursorRemoveWorldPositions, worldPos)
					vw.leftCursorRemoveWorldHistoryPositions = append(vw.leftCursorRemoveWorldHistoryPositions, worldPos)
				}
				break
			}
		}
	}

	return leftCursorWorldPositions, leftCursorRemoveWorldPositions
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

// queueSelectedVertexSelection は選択頂点の更新要求をキューに積む。
func (vw *ViewerWindow) queueSelectedVertexSelection(xpos, ypos float64, remove bool) {
	screenPos := mmath.Vec2{X: xpos, Y: ypos}
	if remove {
		vw.leftCursorRemoveWindowPositions[screenPos] = 0
		return
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
	for range 20 {
		depth := vw.shader.Msaa().ReadDepthAt(int(xpos), int(ypos), framebufferWidth, framebufferHeight)
		if depth <= 0.0 || depth >= 1.0 {
			continue
		}
		world := vw.getWorldPosition(float32(xpos), float32(ypos), depth, framebufferWidth, framebufferHeight)
		if world != nil {
			return world, true
		}
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
	newYaw := cam.Yaw + xOffset
	newPitch := cam.Pitch + yOffset
	vw.resetCameraPosition(newYaw, newPitch)
}

// updateCameraPositionByCursor はカーソル移動でカメラ位置を更新する。
func (vw *ViewerWindow) updateCameraPositionByCursor(xpos, ypos float64) {
	ratio, _, _ := vw.modifierState(0.07)
	xOffset := (vw.prevCursorPos.X - xpos) * ratio
	yOffset := (vw.prevCursorPos.Y - ypos) * ratio

	cam := vw.shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return
	}

	forward := cam.LookAtCenter.Subed(*cam.Position)
	right := forward.Cross(*cam.Up).Normalized()
	up := right.Cross(forward.Normalized()).Normalized()

	upMovement := up.MuledScalar(-yOffset)
	rightMovement := right.MuledScalar(-xOffset)
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
		otherCam.Yaw = currentCam.Yaw
		otherCam.Pitch = currentCam.Pitch
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
