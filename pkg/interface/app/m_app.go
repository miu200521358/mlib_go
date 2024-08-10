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
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
	"golang.org/x/sys/windows"
)

const physics_default_spf = 1.0 / 60.0
const deform_default_spf = 1.0 / 30.0
const deform_default_fps = 30.0

type MApp struct {
	*appState                                   // アプリ状態
	viewerToControlChannel *channelState        // ビューアからコントロールへのチャンネル
	controlToViewerChannel *channelState        // コントロールからビューアへのチャンネル
	appConfig              *mconfig.AppConfig   // アプリケーション設定
	viewWindows            []state.IViewWindow  // 描画ウィンドウリスト
	controlWindow          state.IControlWindow // 操作ウィンドウ
}

func NewMApp(appConfig *mconfig.AppConfig) *MApp {
	// GL初期化
	if err := glfw.Init(); err != nil {
		mlog.F("Failed to initialize GLFW: %v", err)
		return nil
	}

	app := &MApp{
		appState:               newAppState(),
		viewerToControlChannel: newChannelState(),
		controlToViewerChannel: newChannelState(),
		appConfig:              appConfig,
		viewWindows:            make([]state.IViewWindow, 0),
	}

	return app
}

func (app *MApp) AppState() state.IAppState {
	return app.appState
}

func (app *MApp) ViewerToControlChannel() state.IChannelState {
	return app.viewerToControlChannel
}

func (app *MApp) ControlToViewerChannel() state.IChannelState {
	return app.controlToViewerChannel
}

func (app *MApp) RunController() {
	// 操作ウィンドウは別スレッドで起動している前提
	if app.appConfig.IsEnvProd() || app.appConfig.IsEnvDev() {
		defer app.recoverFromPanic()
	}
	app.controlWindow.Run()
}

func (app *MApp) RunViewerToControlChannel() {
	go func() {
		for !app.IsClosed() {
			select {
			case frame := <-app.viewerToControlChannel.frameChannel:
				app.controlWindow.SetFrame(frame)
			case selectedVertexIndexes := <-app.viewerToControlChannel.selectedVertexIndexesChannel:
				app.controlWindow.UpdateSelectedVertexIndexes(selectedVertexIndexes)
			case closed := <-app.viewerToControlChannel.isClosedChannel:
				app.controlWindow.SetClosed(closed)
			default:
				continue
			}
		}
	}()
}

func (app *MApp) RunControlToViewerChannel() {
	go func() {
		for !app.IsClosed() {
			select {
			case frame := <-app.controlToViewerChannel.frameChannel:
				app.SetFrame(frame)
			case maxFrame := <-app.controlToViewerChannel.maxFrameChannel:
				app.UpdateMaxFrame(maxFrame)
			case enabledFrameDrop := <-app.controlToViewerChannel.isEnabledFrameDropChannel:
				app.SetEnabledFrameDrop(enabledFrameDrop)
			case enabledPhysics := <-app.controlToViewerChannel.isEnabledPhysicsChannel:
				app.SetEnabledPhysics(enabledPhysics)
			case resetPhysics := <-app.controlToViewerChannel.isPhysicsResetChannel:
				app.SetPhysicsReset(resetPhysics)
			case showNormal := <-app.controlToViewerChannel.isShowNormalChannel:
				app.SetShowNormal(showNormal)
			case showWire := <-app.controlToViewerChannel.isShowWireChannel:
				app.SetShowWire(showWire)
			case showOverride := <-app.controlToViewerChannel.isShowOverrideChannel:
				app.SetShowOverride(showOverride)
			case showSelectedVertex := <-app.controlToViewerChannel.isShowSelectedVertexChannel:
				app.SetShowSelectedVertex(showSelectedVertex)
			case showBoneAll := <-app.controlToViewerChannel.isShowBoneAllChannel:
				app.SetShowBoneAll(showBoneAll)
			case showBoneIk := <-app.controlToViewerChannel.isShowBoneIkChannel:
				app.SetShowBoneIk(showBoneIk)
			case showBoneEffector := <-app.controlToViewerChannel.isShowBoneEffectorChannel:
				app.SetShowBoneEffector(showBoneEffector)
			case showBoneFixed := <-app.controlToViewerChannel.isShowBoneFixedChannel:
				app.SetShowBoneFixed(showBoneFixed)
			case showBoneRotate := <-app.controlToViewerChannel.isShowBoneRotateChannel:
				app.SetShowBoneRotate(showBoneRotate)
			case showBoneTranslate := <-app.controlToViewerChannel.isShowBoneTranslateChannel:
				app.SetShowBoneTranslate(showBoneTranslate)
			case showBoneVisible := <-app.controlToViewerChannel.isShowBoneVisibleChannel:
				app.SetShowBoneVisible(showBoneVisible)
			case showRigidBodyFront := <-app.controlToViewerChannel.isShowRigidBodyFrontChannel:
				app.SetShowRigidBodyFront(showRigidBodyFront)
			case showRigidBodyBack := <-app.controlToViewerChannel.isShowRigidBodyBackChannel:
				app.SetShowRigidBodyBack(showRigidBodyBack)
			case showJoint := <-app.controlToViewerChannel.isShowJointChannel:
				app.SetShowJoint(showJoint)
			case showInfo := <-app.controlToViewerChannel.isShowInfoChannel:
				app.SetShowInfo(showInfo)
			case frameInterval := <-app.controlToViewerChannel.frameIntervalChanel:
				app.SetFrameInterval(frameInterval)
			case cameraSync := <-app.controlToViewerChannel.isCameraSyncChannel:
				app.SetCameraSync(cameraSync)
			case closed := <-app.controlToViewerChannel.isClosedChannel:
				app.SetClosed(closed)
			case playing := <-app.controlToViewerChannel.playingChannel:
				app.SetPlaying(playing)
			case frameInterval := <-app.controlToViewerChannel.frameIntervalChanel:
				app.SetFrameInterval(frameInterval)
			default:
				continue
			}
		}
	}()
}

func (app *MApp) RunViewer() {
	// 描画ウィンドウはメインスレッドで起動して描画し続ける
	app.SetFrame(0)
	prevTime := glfw.GetTime()
	prevShowTime := glfw.GetTime()
	elapsedList := make([]float64, 0)

	for !app.IsClosed() {

		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		var elapsed float64
		var timeStep float32
		if !app.IsEnabledFrameDrop() {
			// フレームドロップOFF
			// 物理fpsは60fps固定
			timeStep = physics_default_spf
			// デフォームfpsはspf上限の経過時間
			elapsed = mmath.ClampedFloat(originalElapsed, 0.0, deform_default_spf)
		} else {
			// 物理fpsは経過時間
			timeStep = float32(originalElapsed)
			elapsed = originalElapsed
		}

		if elapsed < app.FrameInterval() {
			// 1フレームの時間が経過していない場合はスキップ
			// fps制限は描画fpsにのみ依存
			continue
		}

		// デフォーム
		vmdDeltas := make([][]*delta.VmdDeltas, len(app.viewWindows))
		models := app.GetModels()
		motions := app.GetMotions()

		for i, w := range app.viewWindows {
			w.LoadModels(models[i])
		}

		for i, w := range app.viewWindows {
			vmdDeltas[i] = deform.Deform(w.Physics(), app, timeStep, models[i], motions[i])
			vmdDeltas[i] = deform.DeformPhysics(w.Physics(), app, timeStep, models[i], motions[i], vmdDeltas[i])
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
			w.Render(models[i], vmdDeltas[i])
		}

		// if app.IsShowSelectedVertex() {
		// 	// 頂点選択機能が有効の場合、選択頂点インデックスを更新
		// 	selectedVertexIndexes := make([][][]int, len(app.animationStates))
		// 	for i := range app.animationStates {
		// 		selectedVertexIndexes[i] = make([][]int, len(app.animationStates[i]))
		// 		for j := range app.animationStates[i] {
		// 			selectedVertexIndexes[i][j] = app.animationStates[i][j].SelectedVertexIndexes()
		// 		}
		// 	}
		// 	app.selectedVertexIndexesChan <- selectedVertexIndexes
		// }

		if app.IsPhysicsReset() {
			// 物理リセット
			app.resetPhysics(models, motions, vmdDeltas, timeStep)
		}

		if app.Playing() {
			// 再生中はフレームを進める
			f := app.Frame() + float32(elapsed*deform_default_fps)
			if f > app.MaxFrame() {
				f = 0
			}
			app.viewerToControlChannel.SetFrameChannel(f)
			app.frame = f
		}

		prevTime = frameTime

		// 描画にかかった時間を計測
		elapsedList = append(elapsedList, originalElapsed)

		if app.IsShowInfo() {
			prevShowTime, elapsedList = app.showInfo(elapsedList, prevShowTime, timeStep)
		}
	}
	app.Close()
}

func (app *MApp) resetPhysics(
	models [][]*pmx.PmxModel, motions [][]*vmd.VmdMotion, vmdDeltas [][]*delta.VmdDeltas, timeStep float32,
) {
	// 物理リセット
	for i, w := range app.viewWindows {
		vmdDeltas[i] = deform.Deform(w.Physics(), app, timeStep, models[i], motions[i])

		// 物理削除
		for _, m := range models[i] {
			if m == nil {
				continue
			}
			w.Physics().DeleteModel(m.Index())
		}

		// ワールド作り直し
		w.Physics().ResetWorld()

		// 物理追加
		for j, m := range models[i] {
			if m == nil || vmdDeltas[i][j] == nil {
				continue
			}
			w.Physics().AddModelByBoneDeltas(m.Index(), m, vmdDeltas[i][j].Bones)
		}

		// 物理再設定
		for j, m := range models[i] {
			if m == nil || vmdDeltas[i][j] == nil {
				continue
			}
			deform.DeformPhysicsByBone(app, models[i][j], vmdDeltas[i][j], w.Physics())
		}
	}
	// リセット完了
	app.SetPhysicsReset(false)
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

func (app *MApp) MainViewWindow() state.IViewWindow {
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
					Text: fmt.Sprintf("ToolName: %s, Version: %s", app.appConfig.Name, app.appConfig.Version) +
						string("\r\n------------\r\n") +
						fmt.Sprintf("GLError: %d", gl.GetError()) +
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
						exec.Command("cmd", "/c", "start", "https://discord.gg/MW2Bn47aCN").Start()
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

func (app *MApp) SetControlWindow(controlWindow state.IControlWindow) {
	app.controlWindow = controlWindow
}

func (app *MApp) AddViewWindow(viewWindow state.IViewWindow) {
	app.viewWindows = append(app.viewWindows, viewWindow)
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
