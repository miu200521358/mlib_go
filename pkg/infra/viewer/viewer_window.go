//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"math"
	"strings"
	"time"
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
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

const (
	rightAngle = 89.9
	// boneHoverMaxScreenDistance はボーンホバー判定の最大距離（ピクセル）。
	boneHoverMaxScreenDistance = 20.0
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
	boneHoverActive bool
	lastBoneHoverAt time.Time

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

	vw.updateBoneHoverExpire()
	var debugBoneHover []*mgl.DebugBoneHover
	if vw.boneHighlighter != nil {
		debugBoneHover = vw.boneHighlighter.DebugBoneHoverInfo()
	}

	logger := logging.DefaultLogger()
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

		if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER) {
			vw.logBonePositions(i, vmdDeltas)
		}

		renderer.Render(vw.shader, vw.list.shared, vmdDeltas, debugBoneHover)
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
				vw.boneHoverActive = false
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
			vw.boneHoverActive = false
		}
	}
	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos
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
	closestBones := make([]*mgl.DebugBoneHover, 0)

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
				closestBones = append(closestBones, &mgl.DebugBoneHover{
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

// logBonePositions はボーンの3D位置を冗長ログへ出力する。
func (vw *ViewerWindow) logBonePositions(modelIndex int, vmdDeltas *delta.VmdDeltas) {
	if vmdDeltas == nil || vmdDeltas.Bones == nil {
		return
	}
	logger := logging.DefaultLogger()
	vmdDeltas.Bones.ForEach(func(_ int, boneDelta *delta.BoneDelta) bool {
		if boneDelta == nil || boneDelta.Bone == nil {
			return true
		}
		pos := boneDelta.FilledGlobalPosition()
		logger.Verbose(logging.VERBOSE_INDEX_VIEWER, "ボーン位置: model=%d bone=%s index=%d pos=%s",
			modelIndex, boneDelta.Bone.Name(), boneDelta.Bone.Index(), pos.StringByDigits(4))
		return true
	})
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
