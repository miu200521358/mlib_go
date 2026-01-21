//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"math"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/deform"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/render"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

const (
	rightAngle = 89.9
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

	tooltipRenderer *mgl.TooltipRenderer
	boneHighlighter *mgl.BoneHighlighter

	modelRenderers []*render.ModelRenderer
	motions        []*motion.VmdMotion
	vmdDeltas      []*delta.VmdDeltas

	leftButtonPressed   bool
	middleButtonPressed bool
	rightButtonPressed  bool
	shiftPressed        bool
	ctrlPressed         bool
	updatedPrevCursor   bool
	prevCursorPos       mmath.Vec2
	cursorX             float64
	cursorY             float64
}

func newViewerWindow(windowIndex int, title string, width, height, positionX, positionY int,
	appConfig *config.AppConfig, mainWindow *glfw.Window, list *ViewerManager) (*ViewerWindow, error) {
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

	vw := &ViewerWindow{
		Window:          glWindow,
		windowIndex:     windowIndex,
		title:           title,
		list:            list,
		shader:          shader,
		tooltipRenderer: tooltipRenderer,
		boneHighlighter: mgl.NewBoneHighlighter(),
		modelRenderers:  make([]*render.ModelRenderer, 0),
		motions:         make([]*motion.VmdMotion, 0),
		vmdDeltas:       make([]*delta.VmdDeltas, 0),
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

// Title はタイトルを返す。
func (vw *ViewerWindow) Title() string {
	return vw.title
}

// SetTitle はタイトルを設定する。
func (vw *ViewerWindow) SetTitle(title string) {
	vw.title = title
	vw.Window.SetTitle(title)
}

func (vw *ViewerWindow) resetCameraPosition(yaw, pitch float64) {
	vw.shader.Camera().ResetPosition(yaw, pitch)
	vw.syncCameraToOthers()
}

func (vw *ViewerWindow) render(frame motion.Frame) {
	w, h := vw.GetSize()
	if w == 0 && h == 0 {
		return
	}
	vw.MakeContextCurrent()
	vw.shader.Resize(w, h)

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

	vw.loadModelRenderers()
	vw.loadMotions()

	vw.renderFloor()

	var debugBoneHover []*mgl.DebugBoneHover
	if vw.boneHighlighter != nil {
		debugBoneHover = vw.boneHighlighter.DebugBoneHoverInfo()
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		limit := 0
		if vw.list.appConfig != nil {
			limit = vw.list.appConfig.CursorPositionLimit
		}
		renderer.SetCursorPositionLimit(limit)
		renderer.SetModelIndex(i)

		motionData := motionFromSlice(vw.motions, i)
		vmdDeltas := vw.buildVmdDeltas(frame, renderer.Model, motionData)
		vw.vmdDeltas = ensureVmdDeltas(vw.vmdDeltas, i)
		vw.vmdDeltas[i] = vmdDeltas

		renderer.Render(vw.shader, vw.list.shared, vmdDeltas, debugBoneHover)
	}

	if vw.list.shared.IsAnyBoneVisible() && vw.boneHighlighter != nil {
		vw.boneHighlighter.CheckAndClearExpiredHighlight()
	}
	if vw.list.shared.IsAnyBoneVisible() && vw.tooltipRenderer != nil && vw.boneHighlighter != nil {
		if boneHover := vw.boneHighlighter.DebugBoneHoverInfo(); len(boneHover) > 0 {
			names := make([]string, 0, len(boneHover))
			for _, hover := range boneHover {
				if hover != nil && hover.Bone != nil {
					names = append(names, hover.Bone.Name())
				}
			}
			if len(names) > 0 {
				vw.tooltipRenderer.Render(strings.Join(names, ", "), float32(vw.cursorX), float32(vw.cursorY), w, h)
			}
		}
	}

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
}

func (vw *ViewerWindow) renderFloor() {
	vw.shader.UseProgram(graphics_api.ProgramTypeFloor)
	vw.shader.FloorRenderer().Bind()
	vw.shader.FloorRenderer().Render()
	vw.shader.FloorRenderer().Unbind()
	vw.shader.ResetProgram()
}

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

func motionFromSlice(motions []*motion.VmdMotion, index int) *motion.VmdMotion {
	if index < 0 || index >= len(motions) {
		return nil
	}
	return motions[index]
}

func ensureVmdDeltas(list []*delta.VmdDeltas, index int) []*delta.VmdDeltas {
	for index >= len(list) {
		list = append(list, nil)
	}
	return list
}

func (vw *ViewerWindow) buildVmdDeltas(frame motion.Frame, modelData *model.PmxModel, motionData *motion.VmdMotion) *delta.VmdDeltas {
	if modelData == nil {
		return nil
	}
	motionHash := ""
	if motionData != nil {
		motionHash = motionData.Hash()
	}
	vmdDeltas := delta.NewVmdDeltas(frame, modelData.Bones, modelData.Hash(), motionHash)
	boneDeltas, _ := deform.ComputeBoneDeltas(modelData, motionData, frame, nil, true, false, false)
	vmdDeltas.Bones = boneDeltas
	vmdDeltas.Morphs = deform.ComputeMorphDeltas(modelData, motionData, frame, nil)
	return vmdDeltas
}

func (vw *ViewerWindow) closeCallback(_ *glfw.Window) {
	vw.list.shared.SetClosed(true)
}

func (vw *ViewerWindow) keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
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
		vw.resetCameraPosition(preset.Yaw, preset.Pitch)
	}
}

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
			if vw.list.shared.IsAnyBoneVisible() {
				vw.selectBoneByCursor(vw.cursorX, vw.cursorY)
			} else if vw.boneHighlighter != nil {
				vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
			}
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = false
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = false
		}
	}
}

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
	if !vw.rightButtonPressed && !vw.middleButtonPressed {
		if vw.list.shared.IsAnyBoneVisible() {
			vw.selectBoneByCursor(xpos, ypos)
		} else if vw.boneHighlighter != nil {
			vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
		}
	}
	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos
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

	depth := vw.shader.Msaa().ReadDepthAt(int(xpos), int(ypos), width, height)
	world, err := mgl32.UnProject(
		mgl32.Vec3{float32(xpos), float32(height) - float32(ypos), depth},
		view,
		projection,
		0,
		0,
		width,
		height,
	)
	if err != nil {
		return
	}

	mouseWorld := mmath.NewVec3()
	mouseWorld.X = -float64(world.X())
	mouseWorld.Y = float64(world.Y())
	mouseWorld.Z = float64(world.Z())

	boneDistances := map[float64][]*mgl.DebugBoneHover{}
	for modelIndex, vmdDeltas := range vw.vmdDeltas {
		if vmdDeltas == nil || vmdDeltas.Bones == nil {
			continue
		}
		vmdDeltas.Bones.ForEach(func(_ int, boneDelta *delta.BoneDelta) bool {
			if boneDelta == nil || boneDelta.Bone == nil {
				return true
			}
			bonePos := boneDelta.FilledGlobalPosition()
			dist := mmath.Round(mouseWorld.Distance(bonePos), 0.01)
			boneDistances[dist] = append(boneDistances[dist], &mgl.DebugBoneHover{
				ModelIndex: modelIndex,
				Bone:       boneDelta.Bone,
				Distance:   dist,
			})
			return true
		})
	}
	if len(boneDistances) == 0 {
		vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
		return
	}

	closestDistance := math.MaxFloat64
	for dist := range boneDistances {
		if dist < closestDistance {
			closestDistance = dist
		}
	}
	closestBones := boneDistances[closestDistance]
	if len(closestBones) > 0 {
		vw.boneHighlighter.UpdateDebugHoverByBones(closestBones, true)
	} else {
		vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
	}
}

func (vw *ViewerWindow) updateCameraAngleByCursor(xpos, ypos float64) {
	ratio := 0.1
	if vw.shiftPressed {
		ratio *= 3
	} else if vw.ctrlPressed {
		ratio *= 0.1
	}

	xOffset := (xpos - vw.prevCursorPos.X) * ratio
	yOffset := (ypos - vw.prevCursorPos.Y) * ratio
	cam := vw.shader.Camera()
	vw.resetCameraPosition(cam.Yaw+xOffset, cam.Pitch+yOffset)
}

func (vw *ViewerWindow) updateCameraPositionByCursor(xpos, ypos float64) {
	ratio := 0.07
	if vw.shiftPressed {
		ratio *= 3
	} else if vw.ctrlPressed {
		ratio *= 0.1
	}

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

func (vw *ViewerWindow) scrollCallback(_ *glfw.Window, _ float64, yoff float64) {
	step := float32(1.0)
	if vw.shiftPressed {
		step *= 5
	} else if vw.ctrlPressed {
		step *= 0.1
	}

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

func (vw *ViewerWindow) focusCallback(_ *glfw.Window, focused bool) {
	if !vw.list.shared.IsFocusLinkEnabled() {
		return
	}
	if focused {
		vw.list.shared.TriggerLinkedFocus(vw.windowIndex)
	}
}

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

func (vw *ViewerWindow) iconifyCallback(_ *glfw.Window, iconified bool) {
	if iconified {
		vw.list.shared.SyncMinimize(vw.windowIndex)
	} else {
		vw.list.shared.SyncRestore(vw.windowIndex)
	}
}
