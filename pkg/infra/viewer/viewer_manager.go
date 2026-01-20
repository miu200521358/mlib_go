//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

const (
	defaultFps = mtime.DefaultFps
)

// ViewerManager はビューワー全体を管理する。
type ViewerManager struct {
	shared     *state.SharedState
	appConfig  *config.AppConfig
	windowList []*ViewerWindow
}

// NewViewerManager はViewerManagerを生成する。
func NewViewerManager(shared *state.SharedState, baseServices base.IBaseServices) *ViewerManager {
	var appConfig *config.AppConfig
	if baseServices != nil {
		if cfg := baseServices.Config(); cfg != nil {
			appConfig = cfg.AppConfig()
		}
	}
	return &ViewerManager{
		shared:     shared,
		appConfig:  appConfig,
		windowList: make([]*ViewerWindow, 0),
	}
}

// AddWindow はウィンドウを追加する。
func (vl *ViewerManager) AddWindow(title string, width, height, positionX, positionY int) error {
	var mainWindow *glfw.Window
	if len(vl.windowList) > 0 {
		mainWindow = vl.windowList[0].Window
	}
	vw, err := newViewerWindow(
		len(vl.windowList),
		title,
		width,
		height,
		positionX,
		positionY,
		vl.appConfig,
		mainWindow,
		vl,
	)
	if err != nil {
		return err
	}
	vl.windowList = append(vl.windowList, vw)
	return nil
}

// InitOverlay はオーバーレイ合成を初期化する。
func (vl *ViewerManager) InitOverlay() {
	if len(vl.windowList) > 1 {
		main := vl.windowList[0]
		sub := vl.windowList[1]
		main.shader.OverrideRenderer().SetSharedTextureID(sub.shader.OverrideRenderer().TextureIDPtr())
	}
}

// Run は描画ループを実行する。
func (vl *ViewerManager) Run() {
	prevTime := glfw.GetTime()
	for !vl.shared.IsClosed() {
		vl.handleWindowLinkage()
		vl.handleWindowFocus()
		vl.handleVSync()
		glfw.PollEvents()

		frameTime := glfw.GetTime()
		elapsed := frameTime - prevTime
		if vl.processFrame(elapsed) {
			prevTime = frameTime
		}
	}

	for _, vw := range vl.windowList {
		vw.Destroy()
	}
	glfw.Terminate()
}

func (vl *ViewerManager) handleWindowLinkage() {
	if !vl.shared.HasFlag(state.STATE_FLAG_WINDOW_LINKAGE) {
		return
	}
	if !vl.shared.IsControlWindowMoving() {
		return
	}
	pos := vl.shared.ControlWindowPosition()
	if pos.DiffX == 0 && pos.DiffY == 0 {
		vl.shared.SetControlWindowMoving(false)
		return
	}
	for _, vw := range vl.windowList {
		x, y := vw.GetPos()
		vw.SetPos(x+pos.DiffX, y+pos.DiffY)
	}
	vl.shared.SetControlWindowMoving(false)
}

func (vl *ViewerManager) handleWindowFocus() {
	if !vl.shared.IsFocusLinkEnabled() {
		return
	}
	if !vl.shared.IsControlWindowReady() || !vl.shared.IsAllViewerWindowsReady() {
		return
	}
	for i := len(vl.windowList) - 1; i >= 0; i-- {
		if vl.shared.IsViewerWindowFocused(i) {
			vl.windowList[i].Focus()
			vl.shared.KeepFocus()
			vl.shared.SetViewerWindowFocused(i, false)
			return
		}
	}
}

func (vl *ViewerManager) handleVSync() {
	if !vl.shared.IsFpsLimitTriggered() {
		return
	}
	if vl.shared.FrameInterval() < 0 {
		glfw.SwapInterval(0)
	} else {
		glfw.SwapInterval(1)
	}
	vl.shared.SetFpsLimitTriggered(false)
}

func (vl *ViewerManager) processFrame(elapsed float64) bool {
	frame := vl.shared.Frame()
	maxFrame := vl.shared.MaxFrame()
	if elapsed < 0 {
		return false
	}
	if vl.shared.HasFlag(state.STATE_FLAG_PLAYING) {
		spf := vl.shared.FrameInterval()
		if spf <= 0 {
			spf = mtime.FpsToSpf(defaultFps)
		}
		deltaFrame := mtime.Frame(float32(elapsed) / float32(spf))
		if deltaFrame > 0 {
			frame += deltaFrame
			if maxFrame > 0 && frame > maxFrame {
				frame = 0
			}
			vl.shared.SetFrame(frame)
		}
	}

	for _, vw := range vl.windowList {
		vw.render(frame)
	}
	return true
}
