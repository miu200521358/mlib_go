//go:build windows
// +build windows

// 指示: miu200521358
package viewer

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/win"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/mlib_go/pkg/usecase/mdeform"
	"github.com/miu200521358/mlib_go/pkg/usecase/mphysics"
)

const (
	defaultFps                      = mtime.DefaultFps
	physicsInitialFrame mtime.Frame = 60
)

// ViewerManager はビューワー全体を管理する。
type ViewerManager struct {
	shared                *state.SharedState
	appConfig             *config.AppConfig
	userConfig            config.IUserConfig
	windowList            []*ViewerWindow
	physicsMotionsScratch []*motion.VmdMotion
	timeStepsScratch      []float32
	viewerProfileActive   bool
	viewerProfilePath     string
	viewerProfileFile     *os.File
}

// NewViewerManager はViewerManagerを生成する。
func NewViewerManager(shared *state.SharedState, baseServices base.IBaseServices) *ViewerManager {
	var appConfig *config.AppConfig
	var userConfig config.IUserConfig
	if baseServices != nil {
		if cfg := baseServices.Config(); cfg != nil {
			appConfig = cfg.AppConfig()
			userConfig = cfg.UserConfig()
		}
	}
	return &ViewerManager{
		shared:     shared,
		appConfig:  appConfig,
		userConfig: userConfig,
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
	prevShowTime := prevTime
	prevShowInfo := false
	elapsedList := make([]float64, 0, 1200)
	for !vl.shared.IsClosed() {
		vl.handleWindowLinkage()
		vl.handleWindowFocus()
		vl.handleVSync()
		glfw.PollEvents()
		vl.updateViewerProfile(logging.DefaultLogger())

		frameTime := glfw.GetTime()
		elapsed := frameTime - prevTime
		rendered, meanTimeStep := vl.processFrame(elapsed)
		showInfo := vl.shared.HasFlag(state.STATE_FLAG_SHOW_INFO)
		if rendered {
			if showInfo {
				elapsedList = append(elapsedList, elapsed)
				currentTime := glfw.GetTime()
				if currentTime-prevShowTime >= 1.0 {
					vl.updateFpsDisplay(mmath.Mean(elapsedList), meanTimeStep)
					prevShowTime = currentTime
					elapsedList = elapsedList[:0]
				}
			} else if prevShowInfo {
				vl.resetInfoDisplay()
				elapsedList = elapsedList[:0]
			}
			prevTime = frameTime
		}
		if prevShowInfo != showInfo {
			prevShowTime = glfw.GetTime()
			elapsedList = elapsedList[:0]
			prevShowInfo = showInfo
		}
	}

	for _, vw := range vl.windowList {
		vw.cleanupResources()
		vw.Destroy()
	}
	vl.stopViewerProfile(logging.DefaultLogger())
	glfw.Terminate()
}

// updateViewerProfile はビューワー冗長ログの状態に応じてpprofを開始/終了する。
func (vl *ViewerManager) updateViewerProfile(logger logging.ILogger) {
	if logger == nil {
		return
	}
	enabled := logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER)
	if enabled && !vl.viewerProfileActive {
		vl.startViewerProfile(logger)
		return
	}
	if !enabled && vl.viewerProfileActive {
		vl.stopViewerProfile(logger)
	}
}

// startViewerProfile はCPUプロファイルの計測を開始する。
func (vl *ViewerManager) startViewerProfile(logger logging.ILogger) {
	if vl.viewerProfileActive {
		return
	}
	file, displayPath, fullPath, err := vl.createViewerProfileFile()
	if err != nil {
		if logger != nil {
			logger.Error("ビューワープロファイル開始に失敗しました: %s", sanitizeProfileError(err, fullPath, displayPath))
		}
		return
	}
	if err := pprof.StartCPUProfile(file); err != nil {
		_ = file.Close()
		if logger != nil {
			logger.Error("ビューワープロファイル開始に失敗しました: %s", sanitizeProfileError(err, fullPath, displayPath))
		}
		return
	}
	vl.viewerProfileFile = file
	vl.viewerProfilePath = displayPath
	vl.viewerProfileActive = true
	if logger != nil {
		logger.Info("ビューワープロファイル開始: %s", displayPath)
	}
}

// stopViewerProfile はCPUプロファイルの計測を終了する。
func (vl *ViewerManager) stopViewerProfile(logger logging.ILogger) {
	if !vl.viewerProfileActive {
		return
	}
	pprof.StopCPUProfile()
	if vl.viewerProfileFile != nil {
		if err := vl.viewerProfileFile.Close(); err != nil && logger != nil {
			logger.Error("ビューワープロファイル保存に失敗しました: %s", sanitizeProfileError(err, "", vl.viewerProfilePath))
		}
		vl.viewerProfileFile = nil
	}
	if logger != nil && vl.viewerProfilePath != "" {
		logger.Info("ビューワープロファイル出力: %s", vl.viewerProfilePath)
	}
	vl.viewerProfileActive = false
	vl.viewerProfilePath = ""
}

// createViewerProfileFile はpprof出力用ファイルを生成する。
func (vl *ViewerManager) createViewerProfileFile() (*os.File, string, string, error) {
	dir, displayDir, err := vl.resolveViewerProfileDir()
	if err != nil {
		return nil, "", "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, "", "", err
	}
	fileName := viewerProfileFileName()
	fullPath := filepath.Join(dir, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, "", fullPath, err
	}
	displayPath := filepath.Join(displayDir, fileName)
	return file, displayPath, fullPath, nil
}

// resolveViewerProfileDir はpprofの保存先ディレクトリを決定する。
func (vl *ViewerManager) resolveViewerProfileDir() (string, string, error) {
	if vl.userConfig != nil {
		root, err := vl.userConfig.AppRootDir()
		if err == nil && root != "" {
			return filepath.Join(root, "logs"), "logs", nil
		}
	}
	return "logs", "logs", nil
}

// viewerProfileFileName はpprof出力ファイル名を生成する。
func viewerProfileFileName() string {
	stamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("pprof_viewer_%s.pprof", stamp)
}

// sanitizeProfileError はパスをマスクしたエラー文言を返す。
func sanitizeProfileError(err error, fullPath, displayPath string) string {
	if err == nil {
		return ""
	}
	text := err.Error()
	if fullPath != "" && displayPath != "" {
		text = strings.ReplaceAll(text, fullPath, displayPath)
		dir := filepath.Dir(fullPath)
		if dir != "" {
			text = strings.ReplaceAll(text, dir, filepath.Dir(displayPath))
		}
	}
	return text
}

// handleWindowLinkage はウィンドウ連動移動を処理する。
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

// handleWindowFocus はフォーカス連動の前面化を処理する。
func (vl *ViewerManager) handleWindowFocus() {
	if !vl.shared.IsFocusLinkEnabled() {
		return
	}
	if !vl.shared.IsControlWindowReady() || !vl.shared.IsAllViewerWindowsReady() {
		return
	}
	logger := logging.DefaultLogger()
	verboseEnabled := false
	if logger != nil {
		verboseEnabled = logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER)
	}
	if vl.shared.IsControlWindowFocused() {
		vl.bringViewerWindowsToFrontNoActivate()
		vl.bringControlWindowToFrontNoActivate()
		if verboseEnabled && logger != nil {
			logger.Verbose(logging.VERBOSE_INDEX_VIEWER, "フォーカス連動: コントローラー操作で他ウィンドウを前面化 -> 自分を再前面化")
		}
		vl.shared.KeepFocus()
		vl.shared.SetControlWindowFocused(false)
		return
	}
	for i := len(vl.windowList) - 1; i >= 0; i-- {
		if vl.shared.IsViewerWindowFocused(i) {
			vl.bringControlWindowToFrontNoActivate()
			vl.bringViewerWindowsToFrontNoActivateExcept(i)
			vl.bringWindowToFrontNoActivate(vl.shared.ViewerWindowHandle(i))
			vl.windowList[i].Focus()
			if verboseEnabled && logger != nil {
				logger.Verbose(logging.VERBOSE_INDEX_VIEWER, "フォーカス連動: ビューワー%d操作で他ウィンドウを前面化 -> 自分を再前面化", i)
			}
			vl.shared.KeepFocus()
			vl.shared.SetViewerWindowFocused(i, false)
			return
		}
	}
}

// bringViewerWindowsToFrontNoActivate はビューワーウィンドウを前面化する（フォーカスは奪わない）。
func (vl *ViewerManager) bringViewerWindowsToFrontNoActivate() {
	for i := 0; i < len(vl.windowList); i++ {
		handle := vl.shared.ViewerWindowHandle(i)
		vl.bringWindowToFrontNoActivate(handle)
	}
}

// bringViewerWindowsToFrontNoActivateExcept は指定したビューワー以外を前面化する。
func (vl *ViewerManager) bringViewerWindowsToFrontNoActivateExcept(excludeIndex int) {
	for i := 0; i < len(vl.windowList); i++ {
		if i == excludeIndex {
			continue
		}
		handle := vl.shared.ViewerWindowHandle(i)
		vl.bringWindowToFrontNoActivate(handle)
	}
}

// bringControlWindowToFrontNoActivate はコントロールウィンドウを前面化する（フォーカスは奪わない）。
func (vl *ViewerManager) bringControlWindowToFrontNoActivate() {
	vl.bringWindowToFrontNoActivate(vl.shared.ControlWindowHandle())
}

// bringWindowToFrontNoActivate はウィンドウを非アクティブのまま前面へ移動する。
func (vl *ViewerManager) bringWindowToFrontNoActivate(handle state.WindowHandle) {
	if handle == 0 {
		return
	}
	hwnd := win.HWND(uintptr(handle))
	flags := uint32(win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_NOACTIVATE)
	// 他アプリの前面でも可視になるよう、TOPMOSTの付け外しで前面化する。
	win.SetWindowPos(hwnd, win.HWND_TOPMOST, 0, 0, 0, 0, flags)
	win.SetWindowPos(hwnd, win.HWND_NOTOPMOST, 0, 0, 0, 0, flags)
}

// handleVSync はVSyncの切り替えを処理する。
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

// processFrame は1フレーム分の更新と描画を実行する。
func (vl *ViewerManager) processFrame(elapsed float64) (bool, float32) {
	if elapsed < 0 {
		return false, 0
	}
	if len(vl.windowList) == 0 {
		return false, 0
	}

	logger := logging.DefaultLogger()
	frame := vl.shared.Frame()
	maxFrame := vl.shared.MaxFrame()
	playing := vl.shared.HasFlag(state.STATE_FLAG_PLAYING)
	frameDrop := vl.shared.HasFlag(state.STATE_FLAG_FRAME_DROP) || !playing
	defaultSpf := mtime.FpsToSpf(defaultFps)
	physicsEnabled := vl.shared.HasFlag(state.STATE_FLAG_PHYSICS_ENABLED)

	if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
			"物理フレーム開始: frame=%v playing=%t physics=%t windowCount=%d elapsed=%.4f",
			frame, playing, physicsEnabled, len(vl.windowList), elapsed)
	}

	physicsWorldMotions := vl.physicsMotionsScratch
	if cap(physicsWorldMotions) < len(vl.windowList) {
		physicsWorldMotions = make([]*motion.VmdMotion, len(vl.windowList))
	} else {
		physicsWorldMotions = physicsWorldMotions[:len(vl.windowList)]
		clear(physicsWorldMotions)
	}
	vl.physicsMotionsScratch = physicsWorldMotions

	timeSteps := vl.timeStepsScratch
	if cap(timeSteps) < len(vl.windowList) {
		timeSteps = make([]float32, len(vl.windowList))
	} else {
		timeSteps = timeSteps[:len(vl.windowList)]
		clear(timeSteps)
	}
	vl.timeStepsScratch = timeSteps

	// 再生中はフレーム落ちを抑えるため、経過時間を制限する。
	frameElapsed := float32(elapsed)
	if !frameDrop {
		frameElapsed = float32(mmath.Clamped(elapsed, 0, float64(defaultSpf)))
	}

	// 物理刻みと描画タイミングを決定する。
	needRender := false
	frameInterval := vl.shared.FrameInterval()
	for i := range vl.windowList {
		physicsWorldMotions[i] = resolvePhysicsWorldMotion(vl.shared, i)
		if frameDrop {
			timeSteps[i] = float32(elapsed)
		} else {
			timeSteps[i] = resolveFixedTimeStep(physicsWorldMotions[i], motion.Frame(frame))
		}
		if frameInterval <= 0 || mtime.Seconds(frameElapsed) >= frameInterval {
			needRender = true
		}
	}
	if !needRender {
		waitDuration := frameInterval - mtime.Seconds(frameElapsed)
		if waitDuration >= 0.001 {
			time.Sleep(time.Duration(waitDuration*900) * time.Millisecond)
		}
		return false, 0
	}

	// 物理リセット種別をモーションと共有状態から集約する。
	physicsResetType := vl.shared.PhysicsResetType()
	for i := range physicsWorldMotions {
		physicsResetType = maxResetType(physicsResetType, resolveResetTypeFromMotion(physicsWorldMotions[i], motion.Frame(frame)))
	}
	if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
			"物理リセット種別: frame=%v resetType=%d", frame, physicsResetType)
	}

	// 描画前にモデル/モーションを同期する。
	for _, vw := range vl.windowList {
		vw.prepareFrame()
	}

	// 物理差分生成と物理前変形を行う。
	physicsDeltasByWindow := make([][]*delta.PhysicsDeltas, len(vl.windowList))
	for i, vw := range vl.windowList {
		physicsDeltasByWindow[i] = vl.buildPhysicsDeltas(vw, motion.Frame(frame))
		maxSubSteps := resolveMaxSubSteps(physicsWorldMotions[i], motion.Frame(frame))
		fixedTimeStep := resolveFixedTimeStep(physicsWorldMotions[i], motion.Frame(frame))
		if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
			logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
				"物理前変形: window=%d frame=%v timeStep=%.5f maxSubSteps=%d fixed=%.5f",
				vw.windowIndex, frame, timeSteps[i], maxSubSteps, fixedTimeStep)
		}
		vl.deformWindow(
			vw,
			vw.motions,
			motion.Frame(frame),
			timeSteps[i],
			maxSubSteps,
			fixedTimeStep,
			physicsResetType,
			physicsDeltasByWindow[i],
		)
	}

	for i := len(vl.windowList); i > 0; i-- {
		vl.windowList[i-1].render(motion.Frame(frame))
	}

	if physicsResetType != state.PHYSICS_RESET_TYPE_NONE {
		// リセット種別がある場合は物理ワールドを再構築する。
		for i, vw := range vl.windowList {
			gravity := resolveGravity(physicsWorldMotions[i], motion.Frame(frame))
			maxSubSteps := resolveMaxSubSteps(physicsWorldMotions[i], motion.Frame(frame))
			fixedTimeStep := resolveFixedTimeStep(physicsWorldMotions[i], motion.Frame(frame))
			if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
				logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
					"物理リセット開始: window=%d frame=%v resetType=%d gravity=%s",
					vw.windowIndex, frame, physicsResetType, gravity.StringByDigits(4))
			}
			vl.resetPhysics(
				vw,
				motion.Frame(frame),
				timeSteps[i],
				maxSubSteps,
				fixedTimeStep,
				gravity,
				physicsResetType,
				physicsDeltasByWindow[i],
			)
		}
		vl.shared.SetPhysicsResetType(state.PHYSICS_RESET_TYPE_NONE)
	}

	if playing && !vl.shared.IsClosed() {
		// 再生中はフレームを進め、ループ時にリセット種別を追加する。
		deltaFrame := mtime.SecondsToFrames(mtime.Seconds(frameElapsed), defaultFps)
		if deltaFrame > 0 {
			frame += deltaFrame
			if maxFrame > 0 && frame > maxFrame {
				deltaSaveDone := false
				for i := range vl.windowList {
					if vl.shared.IsDeltaSaveEnabled(i) {
						vl.shared.SetDeltaSaveEnabled(i, false)
						deltaSaveDone = true
					}
				}
				if deltaSaveDone {
					frame = maxFrame
					playing = false
					vl.shared.DisableFlag(state.STATE_FLAG_PLAYING)
				} else {
					frame = 0
					vl.shared.SetPhysicsResetType(maxResetType(vl.shared.PhysicsResetType(), state.PHYSICS_RESET_TYPE_START_FRAME))
				}
			}
			if playing && len(physicsWorldMotions) > 0 {
				motionReset := resolveResetTypeFromMotion(physicsWorldMotions[0], motion.Frame(frame))
				vl.shared.SetPhysicsResetType(maxResetType(vl.shared.PhysicsResetType(), motionReset))
			}
			vl.shared.SetFrame(frame)
		}
	}

	if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
			"物理フレーム終了: frame=%v playing=%t physics=%t", frame, playing, physicsEnabled)
	}
	meanTimeStep := float32(mmath.Mean(timeSteps))
	return true, meanTimeStep
}

// updateFpsDisplay は情報表示ON時のFPS表示を更新する。
func (vl *ViewerManager) updateFpsDisplay(meanElapsed float64, meanTimeStep float32) {
	deformFps := 0.0
	if meanElapsed > 0 {
		deformFps = 1.0 / meanElapsed
	}
	suffix := ""
	if vl.appConfig == nil || vl.appConfig.IsProd() {
		suffix = formatFpsSimple(deformFps)
	} else {
		physicsFps := 0.0
		if meanTimeStep > 0 {
			physicsFps = 1.0 / float64(meanTimeStep)
		}
		suffix = formatFpsDetail(deformFps, physicsFps)
	}

	for _, vw := range vl.windowList {
		if vw == nil {
			continue
		}
		vw.Window.SetTitle(vw.Title() + " - " + suffix)
	}
}

// resetInfoDisplay は情報表示OFF時にタイトルを元へ戻す。
func (vl *ViewerManager) resetInfoDisplay() {
	for _, vw := range vl.windowList {
		if vw == nil {
			continue
		}
		vw.Window.SetTitle(vw.Title())
	}
}

// formatFpsSimple はFPS表示を生成する。
func formatFpsSimple(fps float64) string {
	return fmt.Sprintf("%.2f fps", fps)
}

// formatFpsDetail はFPS表示（詳細）を生成する。
func formatFpsDetail(deformFps, physicsFps float64) string {
	return fmt.Sprintf("d) %.2f / p) %.2f fps", deformFps, physicsFps)
}

// resolvePhysicsWorldMotion は物理ワールドモーションを取得する。
func resolvePhysicsWorldMotion(shared state.ISharedState, viewerIndex int) *motion.VmdMotion {
	if shared == nil {
		return nil
	}
	if raw := shared.PhysicsWorldMotion(viewerIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// resolvePhysicsModelMotion は物理モデルモーションを取得する。
func resolvePhysicsModelMotion(shared state.ISharedState, viewerIndex, modelIndex int) *motion.VmdMotion {
	if shared == nil {
		return nil
	}
	if raw := shared.PhysicsModelMotion(viewerIndex, modelIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// resolveWindMotion は風モーションを取得する。
func resolveWindMotion(shared state.ISharedState, viewerIndex int) *motion.VmdMotion {
	if shared == nil {
		return nil
	}
	if raw := shared.WindMotion(viewerIndex); raw != nil {
		if m, ok := raw.(*motion.VmdMotion); ok {
			return m
		}
	}
	return nil
}

// resolveMaxSubSteps は最大演算回数を取得する。
func resolveMaxSubSteps(physicsWorldMotion *motion.VmdMotion, frame motion.Frame) int {
	if physicsWorldMotion == nil || physicsWorldMotion.MaxSubStepsFrames == nil {
		return 2
	}
	maxFrame := physicsWorldMotion.MaxSubStepsFrames.Get(frame)
	if maxFrame == nil || maxFrame.MaxSubSteps <= 0 {
		return 2
	}
	return maxFrame.MaxSubSteps
}

// resolveFixedTimeStep は固定タイムステップを取得する。
func resolveFixedTimeStep(physicsWorldMotion *motion.VmdMotion, frame motion.Frame) float32 {
	if physicsWorldMotion == nil || physicsWorldMotion.FixedTimeStepFrames == nil {
		return float32(1.0 / 60.0)
	}
	fixedFrame := physicsWorldMotion.FixedTimeStepFrames.Get(frame)
	if fixedFrame == nil {
		return float32(1.0 / 60.0)
	}
	return float32(fixedFrame.FixedTimeStep())
}

// resolveGravity は重力ベクトルを取得する。
func resolveGravity(physicsWorldMotion *motion.VmdMotion, frame motion.Frame) *mmath.Vec3 {
	fallback := mmath.UNIT_Y_NEG_VEC3.MuledScalar(9.8)
	if physicsWorldMotion == nil || physicsWorldMotion.GravityFrames == nil {
		return &fallback
	}
	gravityFrame := physicsWorldMotion.GravityFrames.Get(frame)
	if gravityFrame == nil || gravityFrame.Gravity == nil {
		return &fallback
	}
	gravity := *gravityFrame.Gravity
	return &gravity
}

// resolveResetTypeFromMotion はモーション由来のリセット種別を取得する。
func resolveResetTypeFromMotion(physicsWorldMotion *motion.VmdMotion, frame motion.Frame) state.PhysicsResetType {
	if physicsWorldMotion == nil || physicsWorldMotion.PhysicsResetFrames == nil {
		return state.PHYSICS_RESET_TYPE_NONE
	}
	resetFrame := physicsWorldMotion.PhysicsResetFrames.Get(frame)
	if resetFrame == nil {
		return state.PHYSICS_RESET_TYPE_NONE
	}
	return state.PhysicsResetType(resetFrame.PhysicsResetType)
}

// maxResetType は大きい方のリセット種別を返す。
func maxResetType(a, b state.PhysicsResetType) state.PhysicsResetType {
	if a >= b {
		return a
	}
	return b
}

// buildPhysicsDeltas は物理差分を生成する。
func (vl *ViewerManager) buildPhysicsDeltas(vw *ViewerWindow, frame motion.Frame) []*delta.PhysicsDeltas {
	if vw == nil {
		return nil
	}
	deltas := make([]*delta.PhysicsDeltas, len(vw.modelRenderers))
	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		physicsMotion := resolvePhysicsModelMotion(vl.shared, vw.windowIndex, i)
		deltas[i] = mphysics.BuildPhysicsDeltas(renderer.Model, physicsMotion, frame)
	}
	return deltas
}

// deformWindow はビューワー単位で変形・物理・スキニングを実行する。
func (vl *ViewerManager) deformWindow(
	vw *ViewerWindow,
	motions []*motion.VmdMotion,
	frame motion.Frame,
	timeStep float32,
	maxSubSteps int,
	fixedTimeStep float32,
	resetType state.PhysicsResetType,
	physicsDeltas []*delta.PhysicsDeltas,
) {
	if vw == nil {
		return
	}
	if len(motions) == 0 {
		motions = nil
	}

	logger := logging.DefaultLogger()
	physicsEnabled := vl.shared.HasFlag(state.STATE_FLAG_PHYSICS_ENABLED)
	if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) && vw.physics == nil {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS, "物理未初期化: window=%d frame=%v", vw.windowIndex, frame)
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			if i < len(vw.vmdDeltas) {
				vw.vmdDeltas[i] = nil
			}
			continue
		}
		motionData := motionFromIndex(motions, i)
		vw.vmdDeltas[i] = mdeform.BuildBeforePhysics(
			renderer.Model,
			motionData,
			vw.vmdDeltas[i],
			frame,
			&mdeform.DeformOptions{EnableIK: true},
		)
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil {
			continue
		}
		physicsDelta := physicsDeltaFromIndex(physicsDeltas, i)
		if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
			hasDelta := false
			if physicsDelta != nil {
				hasDelta = true
			}
			logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
				"物理同期: window=%d model=%d frame=%v physicsEnabled=%t resetType=%d hasDelta=%t",
				vw.windowIndex, i, frame, physicsEnabled, resetType, hasDelta)
		}
		vw.syncPhysicsModel(i, renderer.Model, vw.vmdDeltas[i], physicsDelta)
		if vw.physics != nil && resetType == state.PHYSICS_RESET_TYPE_CONTINUE_FRAME && physicsDelta != nil {
			vw.physics.UpdatePhysicsSelectively(i, renderer.Model, physicsDelta)
		}
		vw.vmdDeltas[i] = mdeform.BuildForPhysics(
			vw.physics,
			i,
			renderer.Model,
			vw.vmdDeltas[i],
			physicsDelta,
			physicsEnabled,
			resetType,
		)
	}

	if (physicsEnabled || resetType != state.PHYSICS_RESET_TYPE_NONE) && vw.physics != nil {
		if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
			logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
				"物理シミュレーション: window=%d frame=%v timeStep=%.5f maxSubSteps=%d fixed=%.5f",
				vw.windowIndex, frame, timeStep, maxSubSteps, fixedTimeStep)
		}
		vl.updateWind(vw, frame)
		vw.physics.StepSimulation(timeStep, maxSubSteps, fixedTimeStep)
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil {
			continue
		}
		motionData := motionFromIndex(motions, i)
		vw.vmdDeltas[i] = mdeform.BuildAfterPhysics(
			vw.physics,
			physicsEnabled,
			i,
			renderer.Model,
			motionData,
			vw.vmdDeltas[i],
			frame,
		)
		mdeform.ApplySkinning(renderer.Model, vw.vmdDeltas[i].Bones, vw.vmdDeltas[i].Morphs)
	}
}

// updateWind は風パラメータを物理エンジンへ適用する。
func (vl *ViewerManager) updateWind(vw *ViewerWindow, frame motion.Frame) {
	if vw == nil || vw.physics == nil {
		return
	}
	windMotion := resolveWindMotion(vl.shared, vw.windowIndex)
	if windMotion == nil {
		return
	}
	enabledFrame := windMotion.WindEnabledFrames.Get(frame)
	directionFrame := windMotion.WindDirectionFrames.Get(frame)
	speedFrame := windMotion.WindSpeedFrames.Get(frame)
	liftCoeffFrame := windMotion.WindLiftCoeffFrames.Get(frame)
	dragCoeffFrame := windMotion.WindDragCoeffFrames.Get(frame)
	randomnessFrame := windMotion.WindRandomnessFrames.Get(frame)
	turbulenceFrame := windMotion.WindTurbulenceFreqHzFrames.Get(frame)

	if enabledFrame != nil {
		vw.physics.EnableWind(enabledFrame.Enabled)
	}
	if directionFrame != nil && speedFrame != nil && randomnessFrame != nil {
		vw.physics.SetWind(
			directionFrame.Direction,
			float32(speedFrame.Speed),
			float32(randomnessFrame.Randomness),
		)
	}
	if liftCoeffFrame != nil && dragCoeffFrame != nil && turbulenceFrame != nil {
		vw.physics.SetWindAdvanced(
			float32(dragCoeffFrame.DragCoeff),
			float32(liftCoeffFrame.LiftCoeff),
			float32(turbulenceFrame.TurbulenceFreqHz),
		)
	}
}

// updatePhysicsSelectively は継続フレーム用の選択的更新を行う。
func (vl *ViewerManager) updatePhysicsSelectively(
	vw *ViewerWindow,
	frame motion.Frame,
	physicsDeltas []*delta.PhysicsDeltas,
) {
	if vw == nil {
		return
	}
	physicsEnabled := vl.shared.HasFlag(state.STATE_FLAG_PHYSICS_ENABLED)

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		motionData := motionFromIndex(vw.motions, i)
		vw.vmdDeltas[i] = mdeform.BuildBeforePhysics(
			renderer.Model,
			motionData,
			vw.vmdDeltas[i],
			frame,
			&mdeform.DeformOptions{EnableIK: true},
		)
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		physicsDelta := physicsDeltaFromIndex(physicsDeltas, i)
		if vw.physics != nil && physicsDelta != nil {
			vw.physics.UpdatePhysicsSelectively(i, renderer.Model, physicsDelta)
		}
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil {
			continue
		}
		vw.vmdDeltas[i] = mdeform.BuildForPhysics(
			vw.physics,
			i,
			renderer.Model,
			vw.vmdDeltas[i],
			physicsDelta,
			physicsEnabled,
			state.PHYSICS_RESET_TYPE_CONTINUE_FRAME,
		)
	}
}

// resetPhysics は物理リセット処理を実行する。
func (vl *ViewerManager) resetPhysics(
	vw *ViewerWindow,
	frame motion.Frame,
	timeStep float32,
	maxSubSteps int,
	fixedTimeStep float32,
	gravity *mmath.Vec3,
	resetType state.PhysicsResetType,
	physicsDeltas []*delta.PhysicsDeltas,
) {
	if vw == nil || vw.physics == nil {
		return
	}
	logger := logging.DefaultLogger()
	if resetType == state.PHYSICS_RESET_TYPE_CONTINUE_FRAME {
		if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
			logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
				"物理リセット(継続): window=%d frame=%v", vw.windowIndex, frame)
		}
		vl.updatePhysicsSelectively(vw, frame, physicsDeltas)
		return
	}

	iterationFinishFrame, resetMotions := vl.deformForReset(vw, frame, resetType)
	if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_PHYSICS) {
		logger.Verbose(logging.VERBOSE_INDEX_PHYSICS,
			"物理リセット準備: window=%d frame=%v resetType=%d finishFrame=%v",
			vw.windowIndex, frame, resetType, iterationFinishFrame)
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		vw.physics.DeleteModel(i)
	}

	vw.physics.ResetWorld(gravity)

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil || vw.vmdDeltas[i].Bones == nil {
			continue
		}
		physicsDelta := physicsDeltaFromIndex(physicsDeltas, i)
		vw.physics.AddModelByDeltas(i, renderer.Model, vw.vmdDeltas[i].Bones, physicsDelta)
		vw.vmdDeltas[i] = mdeform.BuildForPhysics(
			vw.physics,
			i,
			renderer.Model,
			vw.vmdDeltas[i],
			physicsDelta,
			vl.shared.HasFlag(state.STATE_FLAG_PHYSICS_ENABLED),
			resetType,
		)
	}

	if resetType == state.PHYSICS_RESET_TYPE_START_FIT_FRAME {
		resetEnd := iterationFinishFrame + motion.Frame(physicsInitialFrame) + 3
		for f := motion.Frame(0); f < resetEnd; f++ {
			vl.deformWindow(
				vw,
				resetMotions,
				f,
				fixedTimeStep,
				maxSubSteps,
				fixedTimeStep,
				resetType,
				physicsDeltas,
			)
			if len(vl.windowList) > 0 && vw.windowIndex == 0 {
				vl.windowList[0].render(f)
			}
		}
	}
}

// deformForReset は物理リセット用の変形を準備する。
func (vl *ViewerManager) deformForReset(
	vw *ViewerWindow,
	frame motion.Frame,
	resetType state.PhysicsResetType,
) (motion.Frame, []*motion.VmdMotion) {
	if vw == nil {
		return 0, nil
	}
	vw.MakeContextCurrent()
	vw.loadModelRenderers()
	vw.loadMotions()
	vw.ensurePhysicsModelSlots()

	resetMotions := make([]*motion.VmdMotion, len(vw.modelRenderers))

	if resetType == state.PHYSICS_RESET_TYPE_START_FIT_FRAME {
		return vl.deformForResetStartFit(vw, frame, resetMotions)
	}

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		motionData := motionFromIndex(vw.motions, i)
		vw.vmdDeltas[i] = mdeform.RebuildBeforePhysics(
			renderer.Model,
			motionData,
			vw.vmdDeltas[i],
			frame,
			&mdeform.DeformOptions{EnableIK: true},
		)
	}

	return 0, resetMotions
}

// deformForResetStartFit はSTART_FIT用のリセットモーションを生成する。
func (vl *ViewerManager) deformForResetStartFit(
	vw *ViewerWindow,
	frame motion.Frame,
	resetMotions []*motion.VmdMotion,
) (motion.Frame, []*motion.VmdMotion) {
	// 変位・回転の最大量を集計してリセット移行フレーム数を見積もる。
	deformMaxTranslations := make([][]float64, len(vw.modelRenderers))
	deformMaxRotations := make([][]float64, len(vw.modelRenderers))

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		modelData := renderer.Model
		motionData := motionFromIndex(vw.motions, i)
		if deformMaxTranslations[i] == nil {
			deformMaxTranslations[i] = make([]float64, modelData.Bones.Len())
			deformMaxRotations[i] = make([]float64, modelData.Bones.Len())
			if motionData != nil && motionData.BoneFrames != nil {
				for _, bone := range modelData.Bones.Values() {
					if bone == nil {
						continue
					}
					boneIndex := bone.Index()
					if boneIndex < 0 || boneIndex >= len(deformMaxTranslations[i]) {
						continue
					}
					bf := motionData.BoneFrames.Get(bone.Name()).Get(frame)
					if bf != nil && bf.Position != nil {
						deformMaxTranslations[i][boneIndex] = bf.Position.Length()
					}
					if bf != nil && bf.Rotation != nil {
						deformMaxRotations[i][boneIndex] = bf.Rotation.ToDegree()
					}
				}
			}
		}
	}

	maxTranslation := mmath.Max(mmath.Flatten(deformMaxTranslations))
	maxRotation := mmath.Max(mmath.Flatten(deformMaxRotations))
	iterationFinish := math.Max(math.Max(maxTranslation/0.5, maxRotation/1.0), 60.0)
	iterationFinishFrame := motion.Frame(iterationFinish)

	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		modelData := renderer.Model
		motionData := motionFromIndex(vw.motions, i)
		if resetMotions[i] == nil {
			resetMotions[i] = motion.NewVmdMotion("")
		}
		// 初期フレームをゼロ姿勢で埋める。
		for _, bone := range modelData.Bones.Values() {
			if bone == nil {
				continue
			}
			resetMotions[i].AppendBoneFrame(bone.Name(), newZeroBoneFrame(0))
		}

		vl.applyYStanceRotation(resetMotions[i], modelData)

		// 初期→移行フレームの補間を追加する。
		for _, bone := range modelData.Bones.Values() {
			if bone == nil {
				continue
			}
			baseFrame := resetMotions[i].BoneFrames.Get(bone.Name()).Get(physicsInitialFrame)
			resetMotions[i].AppendBoneFrame(bone.Name(), baseFrame)

			srcFrame := resolveBoneFrameAt(motionData, bone.Name(), frame)
			if srcFrame == nil {
				continue
			}
			nextFrame := copyBoneFrameWithIndex(srcFrame, physicsInitialFrame+iterationFinishFrame)
			ensureBoneFrameDefaults(nextFrame)
			resetMotions[i].AppendBoneFrame(bone.Name(), nextFrame)
		}
	}

	// リセット用モーションで物理前差分を構築する。
	for i, renderer := range vw.modelRenderers {
		if renderer == nil || renderer.Model == nil {
			continue
		}
		vw.vmdDeltas[i] = mdeform.RebuildBeforePhysics(
			renderer.Model,
			resetMotions[i],
			nil,
			0,
			&mdeform.DeformOptions{EnableIK: true},
		)
	}

	return iterationFinishFrame, resetMotions
}

// motionFromIndex はモーション配列から該当モーションを返す。
func motionFromIndex(motions []*motion.VmdMotion, index int) *motion.VmdMotion {
	if index < 0 || index >= len(motions) {
		return nil
	}
	return motions[index]
}

// physicsDeltaFromIndex は物理差分配列から該当差分を返す。
func physicsDeltaFromIndex(deltas []*delta.PhysicsDeltas, index int) *delta.PhysicsDeltas {
	if index < 0 || index >= len(deltas) {
		return nil
	}
	return deltas[index]
}

// newZeroBoneFrame はゼロ値のボーンフレームを生成する。
func newZeroBoneFrame(index motion.Frame) *motion.BoneFrame {
	bf := motion.NewBoneFrame(index)
	pos := mmath.ZERO_VEC3
	rot := mmath.NewQuaternion()
	bf.Position = &pos
	bf.Rotation = &rot
	return bf
}

// ensureBoneFrameDefaults はPosition/Rotationがnilの場合に初期化する。
func ensureBoneFrameDefaults(bf *motion.BoneFrame) {
	if bf == nil {
		return
	}
	if bf.Position == nil {
		pos := mmath.ZERO_VEC3
		bf.Position = &pos
	}
	if bf.Rotation == nil {
		rot := mmath.NewQuaternion()
		bf.Rotation = &rot
	}
}

// resolveBoneFrameAt はモーションからボーンフレームを取得する。
func resolveBoneFrameAt(motionData *motion.VmdMotion, boneName string, frame motion.Frame) *motion.BoneFrame {
	if motionData == nil || motionData.BoneFrames == nil {
		return nil
	}
	return motionData.BoneFrames.Get(boneName).Get(frame)
}

// copyBoneFrameWithIndex はフレーム番号を差し替えて複製する。
func copyBoneFrameWithIndex(src *motion.BoneFrame, index motion.Frame) *motion.BoneFrame {
	if src == nil {
		return motion.NewBoneFrame(index)
	}
	copied, _ := src.Copy()
	bf, ok := copied.(*motion.BoneFrame)
	if !ok || bf == nil {
		return motion.NewBoneFrame(index)
	}
	read := false
	if src.BaseFrame != nil {
		read = src.Read
	}
	bf.BaseFrame = motion.NewBaseFrame(index)
	bf.Read = read
	return bf
}

// applyYStanceRotation は腕ボーンのYスタンス補正を追加する。
func (vl *ViewerManager) applyYStanceRotation(resetMotion *motion.VmdMotion, modelData *model.PmxModel) {
	if resetMotion == nil || modelData == nil || modelData.Bones == nil {
		return
	}
	for _, dir := range []string{"右", "左"} {
		arm := findArmBoneByPrefix(modelData, dir)
		if arm == nil {
			continue
		}
		armVec := boneChildVector(modelData, arm)
		if armVec.IsZero() {
			continue
		}
		sign := directionSign(dir)
		target := mmath.NewVec3()
		target.X = -1 * sign
		target.Y = 1.3
		target.Z = 0
		target = target.Normalized()
		rot := mmath.NewQuaternionRotate(armVec.Normalized(), target)
		bf := motion.NewBoneFrame(0)
		bf.Rotation = &rot
		resetMotion.AppendBoneFrame(arm.Name(), bf)
	}
}

// findArmBoneByPrefix は腕ボーンを名前から取得する。
func findArmBoneByPrefix(modelData *model.PmxModel, direction string) *model.Bone {
	if modelData == nil || modelData.Bones == nil {
		return nil
	}
	name := direction + "腕"
	bone, err := modelData.Bones.GetByName(name)
	if err != nil {
		return nil
	}
	return bone
}

// directionSign は方向文字列から符号を返す。
func directionSign(direction string) float64 {
	switch direction {
	case "左":
		return -1.0
	case "右":
		return 1.0
	}
	return 0
}

// boneChildVector はボーンの子方向ベクトルを推定する。
func boneChildVector(modelData *model.PmxModel, bone *model.Bone) mmath.Vec3 {
	if modelData == nil || bone == nil {
		return mmath.ZERO_VEC3
	}
	if bone.TailIndex >= 0 && modelData.Bones != nil {
		if child, err := modelData.Bones.Get(bone.TailIndex); err == nil && child != nil {
			return child.Position.Subed(bone.Position)
		}
	}
	return bone.TailPosition
}
