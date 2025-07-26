//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
	"github.com/miu200521358/mlib_go/pkg/config/mproc"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
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
	physicsInitialFrame = float32(120.0)
	deformDefaultSpf    = 1.0 / 30.0    // デフォルトのデフォームspf
	deformDefaultFps    = float32(30.0) // デフォルトのデフォームfps
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
		isPhysicsReset := vl.shared.IsPhysicsReset()

		// ウィンドウリンケージ処理
		vl.handleWindowLinkage()

		// ウィンドウフォーカス処理
		vl.handleWindowFocus()

		// FPS制限処理
		vl.handleVSync()

		// イベント処理
		glfw.PollEvents()

		// フレームタイミング計算
		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		// フレームレート制御と描画処理
		if isRendered, timeStep := vl.processFrame(originalElapsed, isPhysicsReset); isRendered {
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

// handleVSync VSync処理
func (vl *ViewerList) handleVSync() {
	if vl.shared.IsTriggeredFpsLimit() {
		// FPS制限を変更したタイミングで、VSyncを再設定
		if vl.shared.FrameInterval() < 0 {
			// FPS無制限でVSync無効
			glfw.SwapInterval(0)
			mproc.SetMaxProcess(true)
		} else {
			// FPS制限でVSync有効
			glfw.SwapInterval(1)
			mproc.SetMaxProcess(false)
		}
		vl.shared.SetTriggeredFpsLimit(false)
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

	for i := len(vl.windowList) - 1; i >= 0; i-- {
		vw := vl.windowList[i]
		if vl.shared.IsFocusViewWindow(i) {
			vw.Focus()
			vl.shared.KeepFocus()
			vl.shared.SetFocusViewWindow(i, false)
		}
	}
}

// processFrame フレーム処理ロジック
func (vl *ViewerList) processFrame(originalElapsed float64, isPhysicsReset bool) (isRendered bool, timeStep float32) {
	var elapsed float32

	if vl.shared.IsEnabledFrameDrop() || !vl.shared.Playing() {
		// フレームドロップON (再生なし時は常にフレームドロップON)
		// 物理fpsは経過時間
		timeStep = float32(originalElapsed)
		elapsed = float32(originalElapsed)
	} else {
		// フレームドロップOFF
		// 物理fpsは固定時間ステップ
		timeStep = vl.shared.FixedTimeStep()
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

		// 経過時間が1フレームの時間未満の場合はもう少し待つ
		return false, timeStep
	}

	for _, vw := range vl.windowList {
		vw.loadMotions(vl.shared)
		// デフォーム処理
		vl.deform(vw, vw.motions, vl.shared.Frame(), timeStep, isPhysicsReset)
	}

	// レンダリング処理
	for n := len(vl.windowList); n > 0; n-- {
		// サブビューワーオーバーレイのため、逆順でレンダリング
		vl.windowList[n-1].render()
	}

	if isPhysicsReset {
		// 物理リセット
		for _, vw := range vl.windowList {
			vl.resetPhysics(vw, isPhysicsReset)
		}

		// リセット完了
		vl.shared.SetPhysicsReset(false)
	}

	// フレーム更新
	if vl.shared.Playing() && !vl.shared.IsClosed() {
		frame := vl.shared.Frame() + (elapsed * deformDefaultFps)
		if frame > vl.shared.MaxFrame() {
			// フレームが最大フレームを超えた場合、かつ変形情報保存中はINDEXを増やす
			for windowIndex, vw := range vl.windowList {
				if vl.shared.IsSaveDelta(windowIndex) && vl.shared.MaxFrame() > 1.0 {
					// 変形情報のインデックスを増やす
					vl.shared.SetSaveDeltaIndex(vw.windowIndex, vw.list.shared.SaveDeltaIndex(vw.windowIndex)+1)
				}
			}

			frame = 0.0
			// 物理リセットON
			vl.shared.SetPhysicsReset(true)
		}
		vl.shared.SetFrame(frame)
	}

	return true, timeStep
}

func (vl *ViewerList) resetPhysics(vw *ViewWindow, isPhysicsReset bool) {
	// 物理リセット用のデフォーム処理
	iterationFinishFrame, physicsResetMotions := vl.deformForReset(vw)

	for _, model := range vw.modelRenderers {
		if model == nil || model.Model == nil {
			continue
		}

		// モデルの物理削除
		vw.physics.DeleteModel(model.Model.Index())
	}

	// ワールド作り直し
	gravity := vl.shared.Gravity()
	vw.physics.ResetWorld(gravity)

	for n, model := range vw.modelRenderers {
		if model == nil || model.Model == nil || vw.vmdDeltas[n] == nil {
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

	// 物理リセット変形を適用（描画は変更しない）
	if vw.list.shared.IsSaveDelta(vw.windowIndex) && vl.shared.IsEnabledPhysics() {
		for frame := float32(0); frame < iterationFinishFrame+physicsInitialFrame+10; frame++ {
			vl.deform(vw, physicsResetMotions, frame, vl.shared.FixedTimeStep(), isPhysicsReset)

			// // レンダリング処理
			// for n := len(vl.windowList); n > 0; n-- {
			// 	// サブビューワーオーバーレイのため、逆順でレンダリング
			// 	vl.windowList[n-1].render()
			// }
		}
	}
}

func (vl *ViewerList) deformForReset(vw *ViewWindow) (float32, []*vmd.VmdMotion) {
	vw.MakeContextCurrent()

	vw.loadModelRenderers(vl.shared)
	vw.loadMotions(vl.shared)

	frame := vl.shared.Frame()

	// 物理リセット変形用モーション
	physicsResetMotions := make([]*vmd.VmdMotion, len(vw.modelRenderers))

	if !(vw.list.shared.IsSaveDelta(vw.windowIndex) && vl.shared.IsEnabledPhysics()) {
		// 焼き込みしない場合は、そのままデフォームして開始する

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

		return 0, physicsResetMotions
	}

	// モデルごとに0F目の変形量を保持
	deformMaxTranslations := make([][]float64, len(vw.modelRenderers))
	deformMaxRotations := make([][]float64, len(vw.modelRenderers))

	for n := range vw.modelRenderers {
		model := vw.modelRenderers[n].Model
		if model == nil {
			continue
		}

		if deformMaxRotations[n] == nil {
			// 各ボーンの変形量を初期化
			deformMaxRotations[n] = make([]float64, model.Bones.Length())
			deformMaxTranslations[n] = make([]float64, model.Bones.Length())

			model.Bones.ForEach(func(boneIndex int, bone *pmx.Bone) bool {
				// リセットフレームの変形量を取得
				bf := vw.motions[n].BoneFrames.Get(bone.Name()).Get(frame)

				if bf.Position != nil {
					deformMaxTranslations[n][boneIndex] = bf.Position.Length()
				}
				if bf.Rotation != nil {
					deformMaxRotations[n][boneIndex] = bf.Rotation.ToDegree()
				}

				return true
			})
		}
	}

	// 変形量(各フレームの最大移動量を0.5、最大回転量を2度に制限した場合の変形用反復回数)
	iterationFinishFrame := float32(max(
		mmath.Max(mmath.Flatten(deformMaxTranslations))/0.5,
		mmath.Max(mmath.Flatten(deformMaxRotations))/2.0,
		60.0)) // 60.0はデフォルトの反復回数

	// 物理リセット変形用モーションを作成する
	for n := range vw.modelRenderers {
		model := vw.modelRenderers[n].Model
		if model == nil {
			continue
		}

		if physicsResetMotions[n] == nil {
			// モーションが未設定の場合、空のモーションを作成
			physicsResetMotions[n] = vmd.NewVmdMotion("")
		}

		model.Bones.ForEach(func(boneIndex int, bone *pmx.Bone) bool {
			// 0F目の変形量をリセット変形用モーションに全部初期化
			bf := vmd.NewBoneFrame(0)
			if bf.Position == nil {
				bf.Position = mmath.NewMVec3()
			}
			if bf.Rotation == nil {
				bf.Rotation = mmath.NewMQuaternion()
			}

			physicsResetMotions[n].AppendBoneFrame(bone.Name(), bf)
			return true
		})

		// モデルに右腕と左腕がある場合、Yスタンスに変形させる
		for _, direction := range []pmx.BoneDirection{pmx.BONE_DIRECTION_RIGHT, pmx.BONE_DIRECTION_LEFT} {
			armBone, err := model.Bones.GetArm(direction)
			if err != nil {
				continue
			}

			// 腕の現在のベクトルを取得
			armVector := armBone.ChildRelativePosition.Normalized()

			// Yスタンスに変形させるためのベクトルを計算
			yStanceVector := &mmath.MVec3{X: -1 * direction.Sign(), Y: 1.3, Z: 0}

			// モーションに回転情報を追加
			bf := vmd.NewBoneFrame(0)
			// 腕のベクトルをYスタンスに変形させる回転情報を追加
			bf.Rotation = mmath.NewMQuaternionRotate(armVector, yStanceVector.Normalized())
			physicsResetMotions[n].AppendBoneFrame(armBone.Name(), bf)
		}

		model.Bones.ForEach(func(boneIndex int, bone *pmx.Bone) bool {
			{
				// 初期位置を保持して物理を動かす
				bf := physicsResetMotions[n].BoneFrames.Get(bone.Name()).Get(physicsInitialFrame)
				physicsResetMotions[n].AppendBoneFrame(bone.Name(), bf)
			}

			{
				// リセットタイミングフレームの変形を保持して、物理変形を適用させる
				v := vw.motions[n].BoneFrames.Get(bone.Name()).Get(frame).Copy()
				bf := v.(*vmd.BoneFrame)
				bf.SetIndex(physicsInitialFrame + iterationFinishFrame)
				if bf.Position == nil {
					bf.Position = mmath.NewMVec3()
				}
				if bf.Rotation == nil {
					bf.Rotation = mmath.NewMQuaternion()
				}
				physicsResetMotions[n].AppendBoneFrame(bone.Name(), bf)
			}
			return true
		})
	}

	for n := range vw.modelRenderers {
		// 初回の物理前変形
		vw.vmdDeltas[n] = deform.DeformBeforePhysicsReset(
			vw.modelRenderers[n].Model,
			physicsResetMotions[n],
			nil,
			0.0,
		)
	}

	return iterationFinishFrame, physicsResetMotions
}

func (vl *ViewerList) deform(vw *ViewWindow, motions []*vmd.VmdMotion, frame, timeStep float32, isPhysicsReset bool) {
	vw.MakeContextCurrent()

	vw.loadModelRenderers(vl.shared)

	// デフォーム処理
	for n := range vw.modelRenderers {
		if vw.modelRenderers[n] == nil {
			vw.vmdDeltas[n] = nil
			continue
		}

		// 物理前変形
		vw.vmdDeltas[n] = deform.DeformBeforePhysics(
			vw.modelRenderers[n].Model,
			motions[n],
			vw.vmdDeltas[n],
			frame,
		)
	}

	for n := range vw.modelRenderers {
		if vw.modelRenderers[n] == nil {
			continue
		}

		// 物理変形のための事前処理
		vw.vmdDeltas[n] = deform.DeformForPhysics(
			vl.shared,
			vw.physics,
			vw.modelRenderers[n].Model,
			vw.vmdDeltas[n],
		)
	}

	if vl.shared.IsEnabledPhysics() || isPhysicsReset {
		// 物理更新
		vw.physics.StepSimulation(timeStep, vl.shared.MaxSubSteps(), vl.shared.FixedTimeStep())
	}

	for n := range vw.modelRenderers {
		if vw.modelRenderers[n] == nil {
			continue
		}

		// 物理後変形
		vw.vmdDeltas[n] = deform.DeformAfterPhysics(
			vl.shared,
			vw.physics,
			vw.modelRenderers[n].Model,
			motions[n],
			vw.vmdDeltas[n],
			frame,
		)

		if vw.list.shared.IsSaveDelta(vw.windowIndex) && vl.shared.IsEnabledPhysics() && !isPhysicsReset {
			// モデルのデフォーム更新
			vw.saveDeltaMotions(frame)
		}
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
