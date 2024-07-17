package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/window"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"
)

type MApp struct {
	*appState                           // UI状態
	appConfig     *mconfig.AppConfig    // アプリケーション設定
	ViewWindows   []window.IViewWindow  // 描画ウィンドウリスト
	ControlWindow window.IControlWindow // 操作ウィンドウ
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
		ViewWindows: make([]window.IViewWindow, 0),
	}

	return app
}

func (a *MApp) ControllerRun() {
	go func() {
		// 操作ウィンドウは別スレッドで起動
		a.ControlWindow.Run()
	}()
}

func (a *MApp) ViewerRun() {
	// 描画ウィンドウはメインスレッドで起動
	for !a.isClosed {
		for _, w := range a.ViewWindows {
			w.Render()
		}
	}
	a.Close()
}

func (a *MApp) Close() {
	for _, w := range a.ViewWindows {
		w.Close()
	}
	a.ControlWindow.Close()

	glfw.Terminate()
	walk.App().Exit(0)
}

// エラー監視
func (a *MApp) recoverFromPanic() {
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
						a.Close()
						os.Exit(1)
					},
				},
			},
		}).Run(); err != nil {
			walk.MsgBox(nil, mi18n.T("エラーダイアログ起動失敗"), string(stackTrace), walk.MsgBoxIconError)
		}

		a.Close()
	}
}

func (a *MApp) SetControlWindow(controlWindow window.IControlWindow) {
	a.ControlWindow = controlWindow

	if a.appConfig.IsEnvProd() || a.appConfig.IsEnvDev() {
		defer a.recoverFromPanic()
	}
}

func (a *MApp) AddViewWindow(viewWindow window.IViewWindow) {
	a.ViewWindows = append(a.ViewWindows, viewWindow)
}

func (a *MApp) Dispose() {
	for _, w := range a.ViewWindows {
		w.Dispose()
	}
	a.ControlWindow.Dispose()
}

func (a *MApp) ResetPhysics() {
	// 物理ON・まだリセット中ではないの時だけリセット処理を行う
	if a.IsEnabledPhysics() {
		for _, w := range a.ViewWindows {
			w.ResetPhysicsStart()
		}
	}
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

func (a *MApp) Center() {
	// スクリーンの解像度を取得
	screenWidth := getSystemMetrics(SM_CXSCREEN)
	screenHeight := getSystemMetrics(SM_CYSCREEN)

	// ウィンドウのサイズを取得
	mWidth, mHeight := a.ControlWindow.Size()

	viewWindowWidth := 0
	viewWindowHeight := 0
	for _, w := range a.ViewWindows {
		gWidth, gHeight := w.Size()
		viewWindowWidth += gWidth
		viewWindowHeight += gHeight
	}

	// ウィンドウを中央に配置
	if a.appConfig.Horizontal {
		centerX := (screenWidth - (mWidth + viewWindowWidth)) / 2
		centerY := (screenHeight - mHeight) / 2

		centerX += viewWindowWidth
		a.ControlWindow.SetPosition(centerX, centerY)

		for _, w := range a.ViewWindows {
			gWidth, _ := w.Size()
			centerX -= gWidth
			w.SetPosition(centerX, centerY)
		}
	} else {
		centerX := (screenWidth - mWidth) / 2
		centerY := (screenHeight - (mHeight + viewWindowHeight)) / 2

		centerY += mHeight
		a.ControlWindow.SetPosition(centerX, centerY)

		for _, w := range a.ViewWindows {
			_, gHeight := w.Size()
			centerY -= gHeight
			w.SetPosition(centerX, centerY)
		}
	}
}
