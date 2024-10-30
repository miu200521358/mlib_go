//go:build windows
// +build windows

package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"time"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/deform"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
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

// RunViewerToControlChannel ビューアからコントロールへのチャンネルを監視
func (app *MApp) RunViewerToControlChannel() {
	go func() {
		for !app.IsClosed() {
			select {
			case frame := <-app.viewerToControlChannel.frameChannel:
				app.controlWindow.SetFrame(frame)
			case selectedVertexes := <-app.viewerToControlChannel.selectedVertexesChannel:
				app.controlWindow.SetSelectedVertexes(selectedVertexes)
			case closed := <-app.viewerToControlChannel.isClosedChannel:
				app.controlWindow.SetClosed(closed)
				return
			default:
				continue
			}
		}
	}()
}

// RunControlToViewerChannel コントロールからビューアへのチャンネルを監視
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
				return
			case playing := <-app.controlToViewerChannel.playingChannel:
				app.SetPlaying(playing)
			case frameInterval := <-app.controlToViewerChannel.frameIntervalChanel:
				app.SetFrameInterval(frameInterval)
			case invisibleMaterials := <-app.controlToViewerChannel.invisibleMaterialsChannel:
				app.SetInvisibleMaterials(invisibleMaterials)
			case selectedVertexes := <-app.controlToViewerChannel.selectedVertexesChannel:
				app.SetSelectedVertexes(selectedVertexes)
			case noSelectedVertexes := <-app.controlToViewerChannel.noSelectedVertexesChannel:
				app.SetNoSelectedVertexes(noSelectedVertexes)
			default:
				continue
			}
		}
	}()
}

func (app *MApp) RunViewer() {
	if app.appConfig.IsEnvProd() || app.appConfig.IsEnvDev() {
		defer app.recoverFromPanic()
	}

	// 描画ウィンドウはメインスレッドで起動して描画し続ける
	app.SetFrame(0)
	prevTime := glfw.GetTime()
	prevShowTime := glfw.GetTime()
	elapsedList := make([]float64, 0)
	vmdDeltas := make([][]*delta.VmdDeltas, app.ViewerCount())
	viewerParameters := make([]*state.ViewerParameter, app.ViewerCount())

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
			// fps制限は描画fpsにのみ依存

			// 待機時間(残り時間の9割)
			waitDuration := (app.FrameInterval() - elapsed) * 0.9

			// waitDurationがapp.FrameIntervalの9割以下ならsleep
			if waitDuration <= app.FrameInterval()*0.9 {
				time.Sleep(time.Duration(waitDuration * float64(time.Second)))
			}

			// 経過時間が1フレームの時間未満の場合はもう少し待つ
			continue
		}

		// デフォーム
		models := app.GetModels()
		motions := app.GetMotions()

		for len(models) < app.ViewerCount() {
			models = append(models, make([]*pmx.PmxModel, 0))
		}
		for len(motions) < app.ViewerCount() {
			motions = append(motions, make([]*vmd.VmdMotion, 0))
		}

		for i, window := range app.viewWindows {
			window.LoadModels(models[i])
		}

		for i, window := range app.viewWindows {
			var err error
			vmdDeltas[i] = deform.Deform(window.Physics(), app, timeStep, models[i], motions[i], vmdDeltas[i])
			vmdDeltas[i], err = deform.DeformPhysics(window.Physics(), app, timeStep, models[i], motions[i], vmdDeltas[i])
			if err != nil {
				mlog.V("DeformPhysics Error: %v", err)
			}
		}

		// 重複描画
		if app.IsShowOverride() {
			for i := 1; i < app.ViewerCount(); i++ {
				app.viewWindows[i].SetOverrideTextureId(app.viewWindows[0].OverrideTextureId())
			}
		} else {
			for i := 1; i < app.ViewerCount(); i++ {
				app.viewWindows[i].SetOverrideTextureId(0)
			}
		}

		for _, window := range app.viewWindows {
			window.PollEvents()
		}

		// カメラの向きとか同期させる
		app.syncViewer(viewerParameters)

		selectedVertexes := make([][][]int, app.ViewerCount())

		for i := app.ViewerCount() - 1; i >= 0; i-- {
			// サブビューワーオーバーレイのため、逆順でレンダリング
			var invisibleMaterials [][]int
			if i < len(app.invisibleMaterials) {
				invisibleMaterials = app.invisibleMaterials[i]
			}
			var windowSelectedVertexes [][]int
			if i < len(app.selectedVertexes) {
				windowSelectedVertexes = app.selectedVertexes[i]
			}
			var windowNoSelectedVertexes [][]int
			if i < len(app.noSelectedVertexes) {
				windowNoSelectedVertexes = app.noSelectedVertexes[i]
			}
			selectedVertexes[i] = app.viewWindows[i].Render(
				models[i], vmdDeltas[i], invisibleMaterials, windowSelectedVertexes, windowNoSelectedVertexes,
				viewerParameters[i])
		}

		if app.IsShowSelectedVertex() && !app.IsClosed() {
			// 頂点選択機能が有効の場合、選択頂点インデックスを更新
			app.viewerToControlChannel.SetSelectedVertexesChannel(selectedVertexes)
			app.selectedVertexes = selectedVertexes
		}

		if app.IsPhysicsReset() {
			// 物理リセット
			app.resetPhysics(models, motions, timeStep)
		}

		if app.Playing() && !app.IsClosed() {
			// 再生中はフレームを進める
			frame := app.Frame() + float32(elapsed*deform_default_fps)
			if frame > app.MaxFrame() {
				frame = 0
			}
			app.viewerToControlChannel.SetFrameChannel(frame)
			app.frame = frame
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

func (app *MApp) syncViewer(viewerParameters []*state.ViewerParameter) {
	// ビューワー冗長ログ出力時は表示固定
	if mlog.IsViewerVerbose() {
		for i := range app.ViewerCount() {
			if viewerParameters[i] == nil {
				viewerParameters[i] = &state.ViewerParameter{}
			}
			viewerParameters[i].Yaw = mgl.VIEWER_VERBOSE_YAW
			viewerParameters[i].Pitch = mgl.VIEWER_VERBOSE_PITCH
			viewerParameters[i].FieldOfViewAngle = mgl.VIEWER_VERBOSE_FIELD_OF_VIEW_ANGLE
			viewerParameters[i].Size = mgl.VIEWER_VERBOSE_WINDOW_SIZE
			viewerParameters[i].CameraPos = mgl.VIEWER_VERBOSE_CAMERA_POSITION
			viewerParameters[i].CameraUp = mgl.VIEWER_VERBOSE_CAMERA_UP
			viewerParameters[i].LookAtCenter = mgl.VIEWER_VERBOSE_LOOK_AT_CENTER
		}
		return
	}

	if !(app.IsCameraSync() || app.IsShowOverride() || mlog.IsViewerVerbose()) {
		for i := range app.ViewerCount() {
			viewerParameters[i] = nil
		}
		return
	}

	if viewerParameters[0] == nil {
		// 初回は入ってないので、メインビューアのパラメータを採用
		for i := range app.ViewerCount() {
			viewerParameters[i] = app.viewWindows[0].GetViewerParameter()
		}
	} else {
		// 2回目以降は変更があったウィンドウのパラメーターを他に適用する
		changedIndex := 0
		for i, window := range app.viewWindows {
			nowViewerParameter := window.GetViewerParameter()
			if !viewerParameters[i].Equals(nowViewerParameter) {
				changedIndex = i
				break
			}
		}
		// 変更があったウィンドウのパラメータを他に適用する
		nowViewerParam := app.viewWindows[changedIndex].GetViewerParameter()
		for i := range app.ViewerCount() {
			viewerParameters[i] = nowViewerParam
		}
	}
}

func (app *MApp) resetPhysics(
	models [][]*pmx.PmxModel, motions [][]*vmd.VmdMotion, timeStep float32,
) {
	vmdDeltas := make([][]*delta.VmdDeltas, app.ViewerCount())

	// 物理リセット
	for i, window := range app.viewWindows {
		// リセット用デフォーム
		vmdDeltas[i] = deform.DeformForReset(window.Physics(), app, timeStep, models[i], motions[i], nil)

		// 物理削除
		for _, model := range models[i] {
			if model == nil {
				continue
			}
			window.Physics().DeleteModel(model.Index())
		}

		// ワールド作り直し
		window.Physics().ResetWorld()

		// 物理追加
		for j, model := range models[i] {
			if model == nil || vmdDeltas[i][j] == nil {
				continue
			}
			window.Physics().AddModelByBoneDeltas(j, model, vmdDeltas[i][j].Bones)
		}

		// 物理再設定
		for j, model := range models[i] {
			if model == nil || vmdDeltas[i][j] == nil {
				continue
			}
			if err := deform.DeformPhysicsByBone(app, models[i][j], vmdDeltas[i][j], window.Physics()); err != nil {
				mlog.V("DeformPhysicsByBone Error: %v", err)
			}
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
	for _, window := range app.viewWindows {
		window.Close()
	}
	if app.controlWindow != nil {
		app.controlWindow.SetClosed(true)
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
			if app.controlWindow != nil && app.ViewerCount() > 0 {
				break
			}
		}

		// スクリーンの解像度を取得
		screenWidth := getSystemMetrics(SM_CXSCREEN)
		screenHeight := getSystemMetrics(SM_CYSCREEN)

		// ウィンドウのサイズを取得
		mWidth, mHeight := app.controlWindow.WindowSize()

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
