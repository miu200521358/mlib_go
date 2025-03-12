//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
)

type ViewerList struct {
	shared     *state.SharedState // SharedState への参照
	appConfig  *mconfig.AppConfig // アプリケーション設定
	viewerList []*ViewWindow
}

func NewViewerList(shared *state.SharedState, appConfig *mconfig.AppConfig) *ViewerList {
	return &ViewerList{
		shared:     shared,
		appConfig:  appConfig,
		viewerList: make([]*ViewWindow, 0),
	}
}

// Add は ViewerList に ViewerWindow を追加します。
func (vl *ViewerList) Add(title string, width, height, positionX, positionY int) error {
	var mainViewerWindow *glfw.Window
	if len(vl.viewerList) > 0 {
		mainViewerWindow = vl.viewerList[0].Window
	}

	viewWindow, err := newViewWindow(
		len(vl.viewerList),
		title,
		width,
		height,
		positionX,
		positionY,
		vl.appConfig.IconImage,
		vl.appConfig.IsEnvProd(),
		mainViewerWindow,
		vl,
	)

	if err != nil {
		return err
	}

	vl.viewerList = append(vl.viewerList, viewWindow)

	return nil
}

const (
	physicsDefaultSpf = 1.0 / 60.0 // デフォルトの物理spf
	deformDefaultSpf  = 1.0 / 30.0 // デフォルトのデフォームspf
	deformDefaultFps  = 30.0       // デフォルトのデフォームfps
)

func (vl *ViewerList) Run() {
	prevTime := glfw.GetTime()
	prevShowTime := prevTime

	elapsedList := make([]float32, 0, 120)

	for !vl.shared.IsClosed() {
		// イベント処理
		glfw.PollEvents()

		// ウィンドウリンケージ処理
		vl.handleWindowLinkage()

		// ウィンドウフォーカス処理
		vl.handleWindowFocus()

		// フレームタイミング計算
		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		// フレームレート制御と描画処理
		if isRendered, elapsed, timeStep := vl.processFrame(originalElapsed); isRendered {
			// 描画にかかった時間を計測
			elapsedList = append(elapsedList, elapsed)

			// 情報表示処理
			if vl.shared.IsShowInfo() {
				currentTime := glfw.GetTime()
				if currentTime-prevShowTime >= 1.0 {
					vl.updateFpsDisplay(float32(mmath.Mean(elapsedList)), timeStep)
					prevShowTime = currentTime
					elapsedList = elapsedList[:0]
				}
			}

			prevTime = frameTime
		}

	}

	// クリーンアップ
	for _, viewWindow := range vl.viewerList {
		viewWindow.Destroy()
	}
}

// ウィンドウリンケージ処理を分離
func (vl *ViewerList) handleWindowLinkage() {
	if vl.shared.IsWindowLinkage() && vl.shared.IsMovedControlWindow() {
		_, _, diffX, diffY := vl.shared.ControlWindowPosition()
		for _, viewWindow := range vl.viewerList {
			x, y := viewWindow.GetPos()
			viewWindow.SetPos(x+diffX, y+diffY)
		}
		vl.shared.SetMovedControlWindow(false)
	}
}

// ウィンドウフォーカス処理を分離
func (vl *ViewerList) handleWindowFocus() {
	if vl.shared.IsFocusViewWindow() {
		for _, viewWindow := range vl.viewerList {
			viewWindow.Focus()
		}
		vl.shared.SetFocusViewWindow(false)
	}
}

// processFrame フレーム処理ロジック
func (vl *ViewerList) processFrame(originalElapsed float64) (isRendered bool, elapsed float32, timeStep float32) {

	if !vl.shared.IsEnabledFrameDrop() {
		timeStep = physicsDefaultSpf
		elapsed = float32(mmath.Clamped(originalElapsed, 0.0, deformDefaultSpf))
	} else {
		timeStep = float32(originalElapsed)
		elapsed = float32(originalElapsed)
	}

	// FPS制限処理
	if elapsed < vl.shared.FrameInterval() {
		waitDuration := (vl.shared.FrameInterval() - elapsed) * 0.9
		if waitDuration >= 0.001 {
			// あえて1000倍にしないで900倍にしているのは、time.Durationの最大値を超えないため
			sleepDur := time.Duration(waitDuration*900) * time.Millisecond
			// 経過時間が1フレームの時間未満の場合はもう少し待つ
			time.Sleep(sleepDur)
		}
		return false, 0, 0
	}

	// レンダリング処理
	for _, viewWindow := range vl.viewerList {
		viewWindow.Render(vl.shared, timeStep)
	}

	// フレーム更新
	if vl.shared.Playing() && !vl.shared.IsClosed() {
		frame := vl.shared.Frame() + float32(elapsed*deformDefaultFps)
		if frame > vl.shared.MaxFrame() {
			vl.shared.SetClosed(true)
			frame = 0
		}
		vl.shared.SetFrame(frame)
	}

	return true, float32(originalElapsed), timeStep
}

// updateFpsDisplay FPS表示を更新する処理
func (vl *ViewerList) updateFpsDisplay(avgElapsed, timeStep float32) {
	fps := 1.0 / avgElapsed
	var suffixFps string

	if vl.appConfig.IsEnvProd() {
		// リリース版の場合、FPSの表示を簡略化
		suffixFps = fmt.Sprintf("%.2f fps", fps)
	} else {
		// 開発版の場合、FPSの表示を詳細化
		suffixFps = fmt.Sprintf("d) %.2f / p) %.2f fps", fps, 1.0/timeStep)
	}

	for _, viewWindow := range vl.viewerList {
		viewWindow.SetTitle(fmt.Sprintf("%s - %s", viewWindow.Title(), suffixFps))
	}
}
