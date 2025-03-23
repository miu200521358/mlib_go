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
	"github.com/miu200521358/mlib_go/pkg/usecase/deform"
)

type ViewerList struct {
	shared     *state.SharedState // SharedState への参照
	appConfig  *mconfig.AppConfig // アプリケーション設定
	windowList []*ViewWindow
}

func NewViewerList(shared *state.SharedState, appConfig *mconfig.AppConfig) *ViewerList {
	return &ViewerList{
		shared:     shared,
		appConfig:  appConfig,
		windowList: make([]*ViewWindow, 0),
	}
}

// Add は ViewerList に ViewerWindow を追加します。
func (vl *ViewerList) Add(title string, width, height, positionX, positionY int) error {
	var mainViewerWindow *glfw.Window
	if len(vl.windowList) > 0 {
		mainViewerWindow = vl.windowList[0].Window
	}

	vw, err := newViewWindow(
		len(vl.windowList),
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

	vl.windowList = append(vl.windowList, vw)

	return nil
}

const (
	physicsDefaultSpf = float32(1.0 / 60.0) // デフォルトの物理spf
	deformDefaultSpf  = 1.0 / 30.0          // デフォルトのデフォームspf
	deformDefaultFps  = float32(30.0)       // デフォルトのデフォームfps
)

func (vl *ViewerList) InitOverride() {
	if len(vl.windowList) > 1 {
		vl.windowList[0].shader.OverrideRenderer().SetSharedTextureID(
			vl.windowList[1].shader.OverrideRenderer().TextureIDPtr())
	}
}

func (vl *ViewerList) Run() {
	prevTime := glfw.GetTime()
	prevShowTime := prevTime

	elapsedList := make([]float64, 0, 1200)

	for !vl.shared.IsClosed() {
		// ウィンドウリンケージ処理
		vl.handleWindowLinkage()

		// ウィンドウフォーカス処理
		vl.handleWindowFocus()

		// イベント処理
		glfw.PollEvents()

		// フレームタイミング計算
		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		// フレームレート制御と描画処理
		if isRendered, timeStep := vl.processFrame(originalElapsed); isRendered {
			// 描画にかかった時間を計測
			elapsedList = append(elapsedList, originalElapsed)

			// 情報表示処理
			if vl.shared.IsShowInfo() {
				currentTime := glfw.GetTime()
				if currentTime-prevShowTime >= 1.0 {
					vl.updateFpsDisplay(mmath.Mean(elapsedList), timeStep)
					prevShowTime = currentTime
					elapsedList = elapsedList[:0]
				}
			}

			prevTime = frameTime
		}
	}

	// クリーンアップ
	for _, vw := range vl.windowList {
		vw.Destroy()
	}
}

// ウィンドウリンケージ処理を
func (vl *ViewerList) handleWindowLinkage() {
	if vl.shared.IsWindowLinkage() && vl.shared.IsMovedControlWindow() {
		_, _, diffX, diffY := vl.shared.ControlWindowPosition()
		for _, vw := range vl.windowList {
			x, y := vw.GetPos()
			vw.SetPos(x+diffX, y+diffY)
		}
		vl.shared.SetMovedControlWindow(false)
	}
}

// ウィンドウフォーカス処理
func (vl *ViewerList) handleWindowFocus() {
	if !vl.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	for i, vw := range vl.windowList {
		if vl.shared.IsFocusViewWindow(i) {
			vw.Focus()
			vl.shared.KeepFocus()
			vl.shared.SetFocusViewWindow(i, false)
		}
	}
}

// processFrame フレーム処理ロジック
func (vl *ViewerList) processFrame(originalElapsed float64) (isRendered bool, timeStep float32) {
	var elapsed float32

	if vl.shared.IsEnabledFrameDrop() {
		// フレームドロップON
		// 物理fpsは経過時間
		timeStep = float32(originalElapsed)
		elapsed = float32(originalElapsed)
	} else {
		// フレームドロップOFF
		// 物理fpsは60fps固定
		timeStep = physicsDefaultSpf
		// デフォームfpsはspf上限の経過時間
		elapsed = float32(mmath.Clamped(originalElapsed, 0.0, deformDefaultSpf))
	}

	if vl.shared.FrameInterval() > 0 && elapsed < vl.shared.FrameInterval() {
		// fps制限は描画fpsにのみ依存

		// 待機時間(残り時間の9割)
		waitDuration := (vl.shared.FrameInterval() - elapsed) * 0.9

		// waitDurationが1ms以上なら、1ms未満になるまで待つ
		if waitDuration >= 0.001 {
			// あえて1000倍にしないで900倍にしているのは、time.Durationの最大値を超えないため
			time.Sleep(time.Duration(waitDuration*900) * time.Millisecond)
		}
		// mlog.I("Elapsed: %.5f, timeStep: %.5f", originalElapsed, timeStep)

		// 経過時間が1フレームの時間未満の場合はもう少し待つ
		return false, timeStep
	}

	for _, vw := range vl.windowList {
		// デフォーム処理
		vl.deform(vw, timeStep)
	}

	// レンダリング処理
	for n := len(vl.windowList); n > 0; n-- {
		// サブビューワーオーバーレイのため、逆順でレンダリング
		vl.windowList[n-1].render()
	}

	if vl.shared.IsPhysicsReset() {
		// 物理リセット
		for _, vw := range vl.windowList {
			vl.resetPhysics(vw)
		}

		// リセット完了
		vl.shared.SetPhysicsReset(false)
	}

	// フレーム更新
	if vl.shared.Playing() && !vl.shared.IsClosed() {
		frame := vl.shared.Frame() + (elapsed * deformDefaultFps)
		if frame > vl.shared.MaxFrame() {
			frame = 0.0
			// 物理リセットON
			vl.shared.SetPhysicsReset(true)
		}
		vl.shared.SetFrame(frame)
	}

	return true, timeStep
}

func (vl *ViewerList) resetPhysics(vw *ViewWindow) {
	// 物理リセット用のデフォーム処理
	vl.deformForReset(vw)

	for _, model := range vw.modelRenderers {
		// モデルの物理削除
		vw.physics.DeleteModel(model.Model.Index())
	}

	// ワールド作り直し
	vw.physics.ResetWorld()

	for n, model := range vw.modelRenderers {
		if model == nil || vw.vmdDeltas[n] == nil {
			continue
		}

		// モデルの物理追加
		vw.physics.AddModelByBoneDeltas(n, model.Model, vw.vmdDeltas[n].Bones)

		// 物理再設定
		vw.vmdDeltas[n] = deform.DeformForPhysics(
			vl.shared,
			vw.physics,
			vw.modelRenderers[n].Model,
			vw.vmdDeltas[n],
		)
	}
}

func (vl *ViewerList) deformForReset(vw *ViewWindow) {
	vw.MakeContextCurrent()

	vw.loadModelRenderers(vl.shared)
	vw.loadMotions(vl.shared)

	frame := vl.shared.Frame()

	// デフォーム処理
	for n := range vw.modelRenderers {
		// 物理前変形
		vw.vmdDeltas[n] = deform.DeformBeforePhysicsReset(
			vw.modelRenderers[n].Model,
			vw.motions[n],
			vw.vmdDeltas[n],
			frame,
		)
	}
}

func (vl *ViewerList) deform(vw *ViewWindow, timeStep float32) {
	vw.MakeContextCurrent()

	vw.loadModelRenderers(vl.shared)
	vw.loadMotions(vl.shared)

	frame := vl.shared.Frame()

	// デフォーム処理
	for n := range vw.modelRenderers {
		// 物理前変形
		vw.vmdDeltas[n] = deform.DeformBeforePhysics(
			vw.modelRenderers[n].Model,
			vw.motions[n],
			vw.vmdDeltas[n],
			frame,
		)

		// 物理変形のための事前処理
		vw.vmdDeltas[n] = deform.DeformForPhysics(
			vl.shared,
			vw.physics,
			vw.modelRenderers[n].Model,
			vw.vmdDeltas[n],
		)
	}

	if vl.shared.IsEnabledPhysics() || vl.shared.IsPhysicsReset() {
		// 物理更新
		vw.physics.StepSimulation(timeStep)
	}

	for n := range vw.modelRenderers {
		// 物理後変形
		vw.vmdDeltas[n] = deform.DeformAfterPhysics(
			vl.shared,
			vw.physics,
			vw.modelRenderers[n].Model,
			vw.motions[n],
			vw.vmdDeltas[n],
			frame,
		)
	}
}

// updateFpsDisplay FPS表示を更新する処理
func (vl *ViewerList) updateFpsDisplay(meanElapsed float64, timeStep float32) {
	deformFps := 1.0 / meanElapsed
	var suffixFps string

	if vl.appConfig.IsEnvProd() {
		// リリース版の場合、FPSの表示を簡略化
		suffixFps = fmt.Sprintf("%.2f fps", deformFps)
	} else {
		// 開発版の場合、FPSの表示を詳細化
		physicsFps := 1.0 / timeStep
		suffixFps = fmt.Sprintf("d) %.2f / p) %.2f fps", deformFps, physicsFps)
	}

	for _, vw := range vl.windowList {
		vw.SetTitle(fmt.Sprintf("%s - %s", vw.Title(), suffixFps))
	}
}
