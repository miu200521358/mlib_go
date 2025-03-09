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
	prevShowTime := glfw.GetTime()
	elapsedList := make([]float64, 0)

	for !vl.shared.IsClosed() {
		glfw.PollEvents()

		if vl.shared.IsWindowLinkage() && vl.shared.IsMovedControlWindow() {
			_, _, diffX, diffY := vl.shared.ControlWindowPosition()
			// mlog.IS("ViewerWindow linkage moving: diffX=%d, diffY=%d", diffX, diffY)
			// すべてのビューワーウィンドウの位置更新をメインスレッドで行う
			for _, viewWindow := range vl.viewerList {
				x, y := viewWindow.GetPos()
				viewWindow.SetPos(x+diffX, y+diffY)
			}
			vl.shared.SetMovedControlWindow(false)
		}

		if vl.shared.IsFocusViewWindow() {
			// mlog.IS("7) Run: SetFocusViewWindow(false)")
			// ビューワーのフォーカスが指示されている場合、全ビューワーを一旦フォーカスにする
			for _, viewWindow := range vl.viewerList {
				// mlog.IS("8) Run: SetFocusViewWindow[%d] (true)", viewWindow.windowIndex)
				viewWindow.Focus()
			}
			// mlog.IS("9) Run: SetFocusViewWindow(false)")
			vl.shared.SetFocusViewWindow(false)
		}

		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		var elapsed float32
		var timeStep float32
		if !vl.shared.IsEnabledFrameDrop() {
			// フレームドロップOFF
			// 物理fpsは60fps固定
			timeStep = physicsDefaultSpf
			// デフォームfpsはspf上限の経過時間
			elapsed = float32(mmath.Clamped(originalElapsed, 0.0, deformDefaultSpf))
		} else {
			// 物理fpsは経過時間
			timeStep = float32(originalElapsed)
			elapsed = float32(originalElapsed)
		}

		if elapsed < vl.shared.FrameInterval() {
			// fps制限は描画fpsにのみ依存

			// 待機時間(残り時間の9割)
			waitDuration := (vl.shared.FrameInterval() - elapsed) * 0.9

			// waitDurationが1ms以上なら、1ms未満になるまで待つ
			if waitDuration >= 0.001 {
				// あえて1000倍にしないで900倍にしているのは、time.Durationの最大値を超えないため
				time.Sleep(time.Duration(waitDuration*900) * time.Millisecond)
			}

			// 経過時間が1フレームの時間未満の場合はもう少し待つ
			continue
		}

		for _, viewWindow := range vl.viewerList {
			viewWindow.Render(vl.shared, timeStep)
		}

		if vl.shared.Playing() && !vl.shared.IsClosed() {
			// 再生中はフレームを進める
			frame := vl.shared.Frame() + float32(elapsed*deformDefaultFps)
			if frame > vl.shared.MaxFrame() {
				frame = 0
			}
			vl.shared.SetFrame(frame)
		}

		prevTime = frameTime

		// 描画にかかった時間を計測
		elapsedList = append(elapsedList, originalElapsed)

		if vl.shared.IsShowInfo() {
			prevShowTime, elapsedList = vl.showInfo(elapsedList, prevShowTime, timeStep)
		}
	}

	for _, viewWindow := range vl.viewerList {
		viewWindow.Destroy()
	}
}

func (vl *ViewerList) showInfo(elapsedList []float64, prevShowTime float64, timeStep float32) (float64, []float64) {
	nowShowTime := glfw.GetTime()

	// 1秒ごとにオリジナルの経過時間からFPSを表示
	if nowShowTime-prevShowTime >= 1.0 {
		elapsed := mmath.Mean(elapsedList)
		var suffixFps string
		if vl.appConfig.IsEnvProd() {
			// リリース版の場合、FPSの表示を簡略化
			suffixFps = fmt.Sprintf("%.2f fps", 1.0/elapsed)
		} else {
			// 開発版の場合、FPSの表示を詳細化
			suffixFps = fmt.Sprintf("d) %.2f / p) %.2f fps", 1.0/elapsed, 1.0/timeStep)
		}

		for _, viewWindow := range vl.viewerList {
			viewWindow.SetTitle(fmt.Sprintf("%s - %s", viewWindow.Title(), suffixFps))
		}

		return nowShowTime, make([]float64, 0)
	}

	return prevShowTime, elapsedList
}
