//go:build windows
// +build windows

package viewer

import (
	"image"
	"math"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
	"github.com/miu200521358/mlib_go/pkg/usecase/deform"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex         int                              // ウィンドウインデックス
	title               string                           // ウィンドウタイトル
	leftButtonPressed   bool                             // 左ボタン押下フラグ
	middleButtonPressed bool                             // 中ボタン押下フラグ
	rightButtonPressed  bool                             // 右ボタン押下フラグ
	shiftPressed        bool                             // Shiftキー押下フラグ
	ctrlPressed         bool                             // Ctrlキー押下フラグ
	updatedPrevCursor   bool                             // 前回のカーソル位置更新フラグ
	prevCursorPos       *mmath.MVec2                     // 前回のカーソル位置
	yaw                 float64                          // カメラyaw
	pitch               float64                          // カメラpitch
	list                *ViewerList                      // ビューワーリスト
	shader              rendering.IShader                // シェーダー
	physics             physics.IPhysics                 // 物理エンジン
	modelRenderers      []*render.ModelRenderer          // モデル描画オブジェクト
	motions             []*vmd.VmdMotion                 // モーションデータ
	vmdDeltas           []*delta.VmdDeltas               // 変形情報
	speculativeCaches   []*deform.SpeculativeDeformCache // 追加: 投機実行キャッシュ
	avgElapsedTime      float32                          // 追加: 平均経過時間
	elapsedSamples      int                              // 追加: 経過時間サンプル数
}

func newViewWindow(
	windowIndex int,
	title string,
	width, height, positionX, positionY int,
	icon image.Image,
	isProd bool,
	mainWindow *glfw.Window,
	list *ViewerList,
) (*ViewWindow, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	glWindow, err := glfw.CreateWindow(width, height, title, nil, mainWindow)
	if err != nil {
		return nil, err
	}

	glWindow.MakeContextCurrent()
	glWindow.SetInputMode(glfw.StickyKeysMode, glfw.True)
	glWindow.SetIcon([]image.Image{icon})
	glfw.SwapInterval(0) // 0=VSync無効, 1=VSync有効

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		return nil, err
	}

	// シェーダー初期化
	shaderFactory := mgl.NewMShaderFactory()
	shader, err := shaderFactory.CreateShader(width, height)
	if err != nil {
		return nil, err
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	vw := &ViewWindow{
		Window:            glWindow,
		windowIndex:       windowIndex,
		title:             title,
		list:              list,
		shader:            shader,
		physics:           mbt.NewMPhysics(),
		prevCursorPos:     mmath.NewMVec2(),
		speculativeCaches: make([]*deform.SpeculativeDeformCache, 0), // 投機実行キャッシュの初期化
		avgElapsedTime:    1.0 / 30.0,                                // デフォルトは30FPS想定
		elapsedSamples:    0,
	}

	glWindow.SetCloseCallback(vw.closeCallback)
	glWindow.SetScrollCallback(vw.scrollCallback)
	glWindow.SetKeyCallback(vw.keyCallback)
	glWindow.SetMouseButtonCallback(vw.mouseCallback)
	glWindow.SetCursorPosCallback(vw.cursorPosCallback)
	// glWindow.SetFocusCallback(vw.focusCallback)

	if !isProd {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                // 同期的なデバッグ出力有効
		gl.DebugMessageCallback(vw.debugMessageCallback, nil) // デバッグコールバック
	}

	// ウィンドウの位置を設定
	vw.SetPos(positionX, positionY)
	vw.list.shared.SetActivateViewWindow(windowIndex, true)
	vw.list.shared.SetInitializedViewWindow(windowIndex, true)

	return vw, nil
}

func (vw *ViewWindow) Title() string {
	return vw.title
}

func (vw *ViewWindow) resetCameraPosition(yaw, pitch float64) {
	vw.yaw = yaw
	vw.pitch = pitch

	// 球面座標系をデカルト座標系に変換
	radius := math.Abs(float64(rendering.InitialCameraPositionZ))

	// 四元数を使ってカメラの方向を計算
	yawRad := mgl64.DegToRad(yaw)
	pitchRad := mgl64.DegToRad(pitch)
	orientation := mmath.NewMQuaternionFromAxisAngles(mmath.MVec3UnitY, yawRad).Mul(
		mmath.NewMQuaternionFromAxisAngles(mmath.MVec3UnitX, pitchRad))
	forwardXYZ := orientation.MulVec3(mmath.MVec3UnitZNeg).MulScalar(radius)

	// カメラ位置を更新
	cam := vw.shader.Camera()
	cam.Position.X = forwardXYZ.X
	cam.Position.Y = rendering.InitialCameraPositionY + forwardXYZ.Y
	cam.Position.Z = forwardXYZ.Z
	vw.shader.SetCamera(cam)
}

func (vw *ViewWindow) Render(timeStep float32) {
	w, h := vw.GetSize()
	if w == 0 && h == 0 {
		// ウィンドウが最小化されている場合は描画しない
		return
	}

	vw.MakeContextCurrent()

	// リサイズ（サイズが変わってなければ何もしない）
	vw.shader.Resize(w, h)

	// MSAAフレームバッファをバインド
	vw.shader.Msaa().Bind()

	// 深度バッファのクリア
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 隠面消去
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// シェーダーのカメラ設定更新
	vw.shader.UpdateCamera()

	// 床描画
	vw.shader.DrawFloor()

	// vw.loadModelRenderers(vw.list.shared)
	// vw.loadMotions(vw.list.shared)

	// for i := range vw.modelRenderers {
	// 	if vw.list.shared.IsChangedEnableDropFrame() {
	// 		// ドロップフレーム設定が変わったらキャッシュもリセット
	// 		vw.speculativeCaches[i].Reset()
	// 		vw.list.shared.SetChangedEnableDropFrame(false)
	// 		vw.avgElapsedTime = 1.0 / 30.0
	// 		vw.elapsedSamples = 0
	// 	}
	// }

	// // 経過時間の平均計算を更新
	// vw.updateAverageElapsedTime(float32(timeStep))

	// currentFrame := vw.list.shared.Frame()

	// // 改良された予測関数を使用
	// nextFrame := vw.predictNextFrame(currentFrame, timeStep)

	// vw.list.deform(vw.windowIndex, vw, timeStep)

	for i, modelRenderer := range vw.modelRenderers {
		if modelRenderer == nil || vw.vmdDeltas[i] == nil {
			continue
		}

		// // スキップ条件を先に確認
		// if i >= len(vw.speculativeCaches) || vw.speculativeCaches[i] == nil {
		// 	// 通常通り計算
		// 	vw.vmdDeltas[i] = deform.Deform(vw.list.shared, vw.physics, modelRenderer.Model, vw.motions[i], vw.vmdDeltas[i], timeStep)
		// } else {
		// 	// キャッシュから結果を取得を試みる
		// 	modelHash := modelRenderer.Model.Hash()
		// 	motionHash := vw.motions[i].Hash()
		// 	cachedDeltas := vw.speculativeCaches[i].GetResult(currentFrame, modelHash, motionHash)

		// 	if cachedDeltas != nil {
		// 		// キャッシュヒット - キャッシュされた結果を使用
		// 		vw.vmdDeltas[i] = cachedDeltas
		// 		vw.vmdDeltas[i] = deform.DeformPhysics(vw.list.shared, vw.physics, modelRenderer.Model, vw.motions[i], vw.vmdDeltas[i], timeStep)
		// 	} else {
		// 		// キャッシュミス - 通常通り計算
		// 		vw.vmdDeltas[i] = deform.Deform(vw.list.shared, vw.physics, modelRenderer.Model, vw.motions[i], vw.vmdDeltas[i], timeStep)
		// 	}
		// }

		// モデルをレンダリング
		modelRenderer.Render(vw.shader, vw.list.shared, vw.vmdDeltas[i])

		// // 次のフレームの投機的計算を開始
		// deform.SpeculativeDeform(
		// 	vw.list.shared,
		// 	vw.physics,
		// 	modelRenderer.Model,
		// 	vw.motions[i],
		// 	vw.vmdDeltas[i],
		// 	vw.speculativeCaches[i],
		// 	timeStep,
		// 	nextFrame,
		// )
	}

	// 物理デバッグ描画
	vw.physics.DrawDebugLines(vw.shader, vw.list.shared.IsShowRigidBodyFront() || vw.list.shared.IsShowRigidBodyBack(),
		vw.list.shared.IsShowJoint(), vw.list.shared.IsShowRigidBodyFront())

	// 深度解決
	vw.shader.Msaa().Resolve()
	// TODO
	// if vw.appState.IsShowOverride() && vw.windowIndex == 0 {
	// 	vw.drawOverride()
	// }
	vw.shader.Msaa().Unbind()

	vw.SwapBuffers()
}

// predictNextFrame フレーム予測部分を修正
func (vw *ViewWindow) predictNextFrame(currentFrame, timeStep float32) float32 {
	if !vw.list.shared.Playing() {
		return currentFrame // 停止中は同じフレーム
	}

	// 直近のフレーム間隔を使用して次のフレームを予測
	frameAdvanceRate := float32(30.0) // 基準フレームレート

	// 実際のフレームレートに対する適応処理
	actualFrameRate := 1.0 / vw.avgElapsedTime
	frameRateRatio := actualFrameRate / 30.0

	// フレームレート比に基づいて調整（より正確なフレーム進行を実現）
	adjustedAdvance := (timeStep * frameAdvanceRate) * frameRateRatio

	// 整数フレームに合わせる補正（オプション）
	nextFrameRaw := currentFrame + adjustedAdvance
	nextFrameInt := float32(int(nextFrameRaw))
	nextFrameFrac := nextFrameRaw - nextFrameInt

	// 端数が0.5に近づかないように補正
	var nextFrame float32
	if nextFrameFrac < 0.3 {
		nextFrame = nextFrameInt
	} else if nextFrameFrac > 0.7 {
		nextFrame = nextFrameInt + 1.0
	} else {
		nextFrame = nextFrameRaw
	}

	// 最大フレームを超えないようにする
	if nextFrame > vw.list.shared.MaxFrame() {
		nextFrame = 0.0
	}

	return nextFrame
}

// 平均経過時間を更新するヘルパーメソッド
// 平均経過時間更新メソッドを改良

func (vw *ViewWindow) updateAverageElapsedTime(elapsed float32) {
	const maxSamples = 30          // サンプル数を減らしてより最近の値に重みを
	const minElapsed = 1.0 / 120.0 // 最小値制限（異常な値を除外）
	const maxElapsed = 1.0 / 15.0  // 最大値制限（極端な遅延を除外）

	// 異常値のフィルタリング
	if elapsed < minElapsed || elapsed > maxElapsed {
		// 極端な値は平均計算から除外
		return
	}

	if vw.elapsedSamples < maxSamples {
		// 初期サンプル集計
		vw.avgElapsedTime = ((vw.avgElapsedTime * float32(vw.elapsedSamples)) + elapsed) / float32(vw.elapsedSamples+1)
		vw.elapsedSamples++
	} else {
		// 動的な重み係数を使用（より迅速に変化に対応）
		alpha := float32(0.2) // 新しいサンプルの重み（大きくすると反応が早く、小さくすると安定）
		vw.avgElapsedTime = (alpha * elapsed) + ((1.0 - alpha) * vw.avgElapsedTime)
	}
}
