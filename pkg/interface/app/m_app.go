//go:build windows
// +build windows

package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"
)

const physicsDefaultSpf = 1.0 / 60.0

type MApp struct {
	*appState                        // アプリ状態
	appConfig     *mconfig.AppConfig // アプリケーション設定
	viewWindows   []IViewWindow      // 描画ウィンドウリスト
	controlWindow IControlWindow     // 操作ウィンドウ
}

func NewMApp(appConfig *mconfig.AppConfig) *MApp {
	// GL初期化
	if err := glfw.Init(); err != nil {
		mlog.F("Failed to initialize GLFW: %v", err)
		return nil
	}

	app := &MApp{
		appState:    newAppState(),
		appConfig:   appConfig,
		viewWindows: make([]IViewWindow, 0),
	}

	return app
}

func (app *MApp) ControllerRun() {
	// 操作ウィンドウは別スレッドで起動している前提
	if app.appConfig.IsEnvProd() || app.appConfig.IsEnvDev() {
		defer app.recoverFromPanic()
	}
	app.controlWindow.SetEnabled(true)
	app.controlWindow.Run()
}

func (app *MApp) ViewerRun() {
	// 描画ウィンドウはメインスレッドで起動して描画し続ける
	app.appState.SetPrevFrame(0)
	app.appState.SetFrame(0)
	prevTime := glfw.GetTime()
	prevShowTime := glfw.GetTime()
	elapsedList := make([]float64, 0)

	for !app.IsClosed() {
		frameTime := glfw.GetTime()
		elapsed := frameTime - prevTime

		if elapsed < app.appState.SpfLimit() {
			// 1フレームの時間が経過していない場合はスキップ
			// fps制限は描画fpsにのみ依存
			continue
		}

		var timeStep float32
		if !app.appState.IsEnabledFrameDrop() {
			// フレームドロップOFF
			// 物理fpsは60fps固定
			timeStep = physicsDefaultSpf
		} else {
			// 物理fpsは経過時間
			timeStep = float32(elapsed)
		}

		// if app.IsEnabledPhysics() && app.IsPhysicsReset() {

		// 	// リセットフラグOFF
		// 	app.SetPhysicsReset(false)
		// }

		if app.IsPhysicsReset() {
			for i, w := range app.viewWindows {
				// 一旦アニメーション
				app.animationStates[i], app.nextAnimationStates[i] =
					w.Animate(app.animationStates[i], app.nextAnimationStates[i], timeStep)
			}

			for i, w := range app.viewWindows {
				w.ResetPhysics(app.animationStates[i])
			}

			// リセットが終わったらフラグを落とす
			app.SetPhysicsReset(false)
		}

		// カメラ同期(重複描画の場合はそっちで同期させる)
		if app.IsCameraSync() && !app.IsShowOverride() {
			for i := 1; i < len(app.viewWindows); i++ {
				app.viewWindows[i].UpdateViewerParameter(app.viewWindows[0].GetViewerParameter())
			}
		}

		// 重複描画
		if app.IsShowOverride() {
			for i := 1; i < len(app.viewWindows); i++ {
				app.viewWindows[i].SetOverrideTextureId(app.viewWindows[0].OverrideTextureId())
				// カメラの向きとか同期させる
				app.viewWindows[i].UpdateViewerParameter(app.viewWindows[0].GetViewerParameter())
			}
		} else {
			for i := 1; i < len(app.viewWindows); i++ {
				app.viewWindows[i].SetOverrideTextureId(0)
			}
		}

		for i := app.ViewerCount() - 1; i >= 0; i-- {
			// サブビューワーオーバーレイのため、逆順でレンダリング
			w := app.viewWindows[i]
			// アニメーション
			if !app.IsShowSelectedVertex() {
				for j := range app.animationStates[i] {
					app.animationStates[i][j].UpdateSelectedVertexIndexes(
						app.animationStates[i][j].SelectedVertexIndexes())
				}
			}
			app.animationStates[i], app.nextAnimationStates[i] =
				w.Animate(app.animationStates[i], app.nextAnimationStates[i], timeStep)
		}

		prevTime = frameTime

		// 描画が終わったらフレーム番号を更新
		app.prevFrame = app.frame

		// 描画にかかった時間を計測
		elapsedList = append(elapsedList, app.deformElapsed+elapsed)

		if app.IsShowInfo() {
			prevShowTime, elapsedList = app.showInfo(elapsedList, prevShowTime, timeStep)
		}
	}
	app.Close()
}

func (app *MApp) showInfo(elapsedList []float64, prevShowTime float64, timeStep float32) (float64, []float64) {
	nowShowTime := glfw.GetTime()

	// 1秒ごとにオリジナルの経過時間からFPSを表示
	if nowShowTime-prevShowTime >= 1.0 {
		elapsed := mmath.Avg(elapsedList)
		var suffixFps string
		if app.appConfig.IsEnvProd() {
			// リリース版の場合、FPSの表示を簡略化
			suffixFps = fmt.Sprintf("%.2f fps", 1.0/elapsed)
		} else {
			// 開発版の場合、FPSの表示を詳細化
			suffixFps = fmt.Sprintf("d) %.2f / p) %.2f fps", 1.0/elapsed, 1.0/timeStep)
		}

		for _, w := range app.viewWindows {
			w.GetWindow().SetTitle(fmt.Sprintf("%s - %s", w.Title(), suffixFps))
		}

		return nowShowTime, make([]float64, 0)
	}

	return prevShowTime, elapsedList
}

func (app *MApp) ViewerCount() int {
	return len(app.viewWindows)
}

func (app *MApp) MainViewWindow() IViewWindow {
	return app.viewWindows[0]
}

func (app *MApp) Close() {
	for _, w := range app.viewWindows {
		w.Close()
	}
	if app.controlWindow != nil {
		app.controlWindow.Close()
	}

	glfw.Terminate()
	walk.App().Exit(0)
}

// エラー監視
func (app *MApp) recoverFromPanic() {
	if r := recover(); r != nil {
		stackTrace := debug.Stack()

		var errMsg string
		// パニックの値がerror型である場合、エラーメッセージを取得
		if err, ok := r.(error); ok {
			errMsg = err.Error()
		} else {
			// それ以外の型の場合は、文字列に変換
			errMsg = fmt.Sprintf("%v", r)
		}

		var errT *walk.TextEdit
		if _, err := (declarative.MainWindow{
			Title:   mi18n.T("予期せぬエラーが発生しました"),
			Size:    declarative.Size{Width: 500, Height: 400},
			MinSize: declarative.Size{Width: 500, Height: 400},
			MaxSize: declarative.Size{Width: 500, Height: 400},
			Layout:  declarative.VBox{},
			Children: []declarative.Widget{
				declarative.TextLabel{
					Text: mi18n.T("予期せぬエラーヘッダー"),
				},
				declarative.TextEdit{
					Text: fmt.Sprintf("GLError: %d", gl.GetError()) +
						string("\r\n------------\r\n") +
						fmt.Sprintf("Error: %s", errMsg) +
						string("\r\n------------\r\n") +
						string(bytes.ReplaceAll(stackTrace, []byte("\n"), []byte("\r\n"))),
					ReadOnly: true,
					AssignTo: &errT,
					VScroll:  true,
					HScroll:  true,
				},
				declarative.PushButton{
					Text:      mi18n.T("コミュニティ報告"),
					Alignment: declarative.AlignHFarVNear,
					OnClicked: func() {
						if err := walk.Clipboard().SetText(errT.Text()); err != nil {
							walk.MsgBox(nil, mi18n.T("クリップボードコピー失敗"),
								string(stackTrace), walk.MsgBoxIconError)
						}
						exec.Command("cmd", "/c", "start", "https://com.nicovideo.jp/community/co5387214").Start()
					},
				},
				declarative.PushButton{
					Text: mi18n.T("アプリを終了"),
					OnClicked: func() {
						app.Close()
						os.Exit(1)
					},
				},
			},
		}).Run(); err != nil {
			walk.MsgBox(nil, mi18n.T("エラーダイアログ起動失敗"), string(stackTrace), walk.MsgBoxIconError)
		}

		app.Close()
	}
}

func (app *MApp) SetControlWindow(controlWindow IControlWindow) {
	app.controlWindow = controlWindow
}

func (app *MApp) AddViewWindow(viewWindow IViewWindow) {
	app.viewWindows = append(app.viewWindows, viewWindow)
	app.animationStates = append(app.animationStates, make([]state.IAnimationState, 0))
	app.nextAnimationStates = append(app.nextAnimationStates, make([]state.IAnimationState, 0))
}

func (app *MApp) Dispose() {
	for _, w := range app.viewWindows {
		w.Dispose()
	}
	app.controlWindow.Dispose()
}

// ----------------------

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func getSystemMetrics(nIndex int) int {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(nIndex))
	return int(ret)
}

func GetWindowSize(width int, height int) declarative.Size {
	screenWidth := getSystemMetrics(SM_CXSCREEN)
	screenHeight := getSystemMetrics(SM_CYSCREEN)

	if width > screenWidth-50 {
		width = screenWidth - 50
	}
	if height > screenHeight-50 {
		height = screenHeight - 50
	}

	return declarative.Size{Width: width, Height: height}
}

func (app *MApp) Center() {
	go func() {
		for {
			if app.controlWindow != nil && len(app.viewWindows) > 0 {
				break
			}
		}

		// スクリーンの解像度を取得
		screenWidth := getSystemMetrics(SM_CXSCREEN)
		screenHeight := getSystemMetrics(SM_CYSCREEN)

		// ウィンドウのサイズを取得
		mWidth, mHeight := app.controlWindow.Size()

		viewWindowWidth := 0
		viewWindowHeight := 0
		for _, w := range app.viewWindows {
			gWidth, gHeight := w.Size()
			viewWindowWidth += gWidth
			viewWindowHeight += gHeight
		}

		// ウィンドウを中央に配置
		if app.appConfig.Horizontal {
			centerX := (screenWidth - (mWidth + viewWindowWidth)) / 2
			centerY := (screenHeight - mHeight) / 2

			centerX += viewWindowWidth
			app.controlWindow.SetPosition(centerX, centerY)

			for _, w := range app.viewWindows {
				gWidth, _ := w.Size()
				centerX -= gWidth
				w.SetPosition(centerX, centerY)
			}
		} else {
			centerX := (screenWidth - mWidth) / 2
			centerY := (screenHeight - (mHeight + viewWindowHeight)) / 2

			centerY += mHeight
			app.controlWindow.SetPosition(centerX, centerY)

			for _, w := range app.viewWindows {
				_, gHeight := w.Size()
				centerY -= gHeight
				w.SetPosition(centerX, centerY)
			}
		}
	}()
}
