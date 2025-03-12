//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"image"
	"math"
	"sync"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
	"github.com/miu200521358/mlib_go/pkg/usecase/deform"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex         int                     // ウィンドウインデックス
	title               string                  // ウィンドウタイトル
	leftButtonPressed   bool                    // 左ボタン押下フラグ
	middleButtonPressed bool                    // 中ボタン押下フラグ
	rightButtonPressed  bool                    // 右ボタン押下フラグ
	shiftPressed        bool                    // Shiftキー押下フラグ
	ctrlPressed         bool                    // Ctrlキー押下フラグ
	updatedPrevCursor   bool                    // 前回のカーソル位置更新フラグ
	prevCursorPos       *mmath.MVec2            // 前回のカーソル位置
	yaw                 float64                 // カメラyaw
	pitch               float64                 // カメラpitch
	list                *ViewerList             // ビューワーリスト
	shader              rendering.IShader       // シェーダー
	physics             physics.IPhysics        // 物理エンジン
	modelRenderers      []*render.ModelRenderer // モデル描画オブジェクト
	motions             []*vmd.VmdMotion        // モーションデータ
	vmdDeltas           []*delta.VmdDeltas      // 変形情報
	nextVmdDeltas       []*delta.VmdDeltas      // 次フレーム用の変形情報
	calcMutex           sync.Mutex              // 並列計算用のミューテックス
	isCalculating       bool                    // 計算中フラグ
	nextFrame           float32                 // 次フレーム
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

	glWindow, err := glfw.CreateWindow(width, height, title, nil, mainWindow)
	if err != nil {
		return nil, err
	}

	glWindow.MakeContextCurrent()
	glWindow.SetInputMode(glfw.StickyKeysMode, glfw.True)
	glWindow.SetIcon([]image.Image{icon})

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		return nil, err
	}

	// エラーチェックを追加
	if err := gl.GetError(); err != gl.NO_ERROR {
		return nil, fmt.Errorf("OpenGL error before shader compilation: %v", err)
	}

	// シェーダー初期化
	shaderFactory := mgl.NewMShaderFactory()
	shader, err := shaderFactory.CreateShader(width, height)
	if err != nil {
		return nil, err
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	vw := &ViewWindow{
		Window:        glWindow,
		windowIndex:   windowIndex,
		title:         title,
		list:          list,
		shader:        shader,
		physics:       mbt.NewMPhysics(),
		prevCursorPos: mmath.NewMVec2(),
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

// キュー状態で次フレームの計算を開始する
func (vw *ViewWindow) startNextFrameCalculation(nextFrame float32) {
	vw.calcMutex.Lock()

	// すでに計算中または同じフレームの計算がキューにある場合はスキップ
	if vw.isCalculating || (vw.nextVmdDeltas != nil && vw.nextFrame == nextFrame) {
		vw.calcMutex.Unlock()
		return
	}

	vw.isCalculating = true
	vw.nextFrame = nextFrame
	vw.calcMutex.Unlock()

	go func() {
		// モデルが読み込まれていない場合は何もしない
		if len(vw.modelRenderers) == 0 || len(vw.motions) == 0 {
			vw.calcMutex.Lock()
			vw.isCalculating = false
			vw.calcMutex.Unlock()
			return
		}

		// 各モデルの次フレームのdeformBeforePhysicsを計算
		nextDeltas := make([]*delta.VmdDeltas, len(vw.modelRenderers))
		for i, modelRenderer := range vw.modelRenderers {
			if i < len(vw.motions) && vw.motions[i] != nil {
				// 現在のデルタが有効で、モデルとモーションのハッシュが一致する場合は再利用
				if i < len(vw.vmdDeltas) && vw.vmdDeltas[i] != nil &&
					vw.vmdDeltas[i].ModelHash() == modelRenderer.Model.Hash() &&
					vw.vmdDeltas[i].MotionHash() == vw.motions[i].Hash() {
					// deformBeforePhysicsのみ並列処理で行う
					nextDeltas[i] = deform.DeformBeforePhysics(
						modelRenderer.Model,
						vw.motions[i],
						vw.vmdDeltas[i],
						nextFrame,
					)
				} else {
					// 完全に新規計算が必要な場合
					nextDeltas[i] = deform.DeformBeforePhysics(
						modelRenderer.Model,
						vw.motions[i],
						nil,
						nextFrame,
					)
				}
			}
		}

		// 結果を保存
		vw.calcMutex.Lock()
		vw.nextVmdDeltas = nextDeltas
		vw.isCalculating = false
		vw.calcMutex.Unlock()
	}()
}

// 初期フレームの計算を実行
func (vw *ViewWindow) initializeFrameCalculation(frame float32) {
	if len(vw.modelRenderers) == 0 || len(vw.motions) == 0 {
		return
	}

	// 初期計算用に次フレーム予測を別途行う
	vw.calcMutex.Lock()
	vw.isCalculating = false // 進行中の計算があれば中止
	vw.calcMutex.Unlock()

	// 現在のフレームで同期計算を実行
	nextDeltas := make([]*delta.VmdDeltas, len(vw.modelRenderers))
	for i, modelRenderer := range vw.modelRenderers {
		if i < len(vw.motions) && vw.motions[i] != nil {
			// 現在のフレームで計算 (並列ではなく同期的に実行)
			nextDeltas[i] = deform.DeformBeforePhysics(
				modelRenderer.Model,
				vw.motions[i],
				vw.vmdDeltas[i],
				frame,
			)
		}
	}

	// 結果を即時反映
	vw.calcMutex.Lock()
	vw.nextVmdDeltas = nextDeltas
	vw.nextFrame = frame
	vw.calcMutex.Unlock()
}

func (vw *ViewWindow) Render(shared *state.SharedState, timeStep float32) {
	w, h := vw.GetSize()
	if w == 0 && h == 0 {
		// ウィンドウが最小化されている場合は描画しない
		return
	}

	mlog.IS("[%d] Start", vw.windowIndex)

	vw.MakeContextCurrent()

	// リサイズ（サイズが変わってなければ何もしない）
	vw.shader.Resize(w, h)

	// MSAAフレームバッファをバインド
	vw.shader.GetMsaa().Bind()

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

	// 前のフレームから変化がある場合にのみモデルとモーションを更新
	prevModelCount := len(vw.modelRenderers)
	prevMotionCount := len(vw.motions)

	vw.loadModelRenderers(shared)
	vw.loadMotions(shared)

	// モデル数かモーション数が変わった場合にのみ初期化
	modelMotionCountChanged := prevModelCount != len(vw.modelRenderers) || prevMotionCount != len(vw.motions)

	// 現在のフレームのデータがない場合か、モデル/モーション構成が変わった場合に初期化計算を実行
	needsInitCalculation := modelMotionCountChanged ||
		len(vw.vmdDeltas) != len(vw.modelRenderers) ||
		(len(vw.modelRenderers) > 0 && (vw.vmdDeltas[0] == nil || vw.vmdDeltas[0].Frame() != shared.Frame()))

	if needsInitCalculation {
		vw.initializeFrameCalculation(shared.Frame())
	}

	// 再生中かつ次フレームの予測計算が必要な場合のみ計算
	if shared.Playing() && !needsInitCalculation {
		nextFrame := shared.Frame() + 1
		if nextFrame < shared.MaxFrame() {
			vw.startNextFrameCalculation(nextFrame)
		}
	}

	// 事前計算された次フレームデータがあれば利用する
	vw.calcMutex.Lock()
	if vw.nextVmdDeltas != nil && len(vw.nextVmdDeltas) == len(vw.modelRenderers) &&
		vw.nextFrame == shared.Frame() {
		// 次フレームデータが現在のフレームと一致する場合は利用
		for i := range vw.nextVmdDeltas {
			if vw.nextVmdDeltas[i] != nil {
				vw.vmdDeltas[i] = vw.nextVmdDeltas[i]
			}
		}
		vw.nextVmdDeltas = nil
	}
	vw.calcMutex.Unlock()

	for i, modelRenderer := range vw.modelRenderers {
		// インデックスのチェックも追加
		if i >= len(vw.vmdDeltas) || vw.vmdDeltas[i] == nil || vw.vmdDeltas[i].Frame() != shared.Frame() ||
			vw.vmdDeltas[i].ModelHash() != modelRenderer.Model.Hash() ||
			(i < len(vw.motions) && vw.motions[i] != nil && vw.vmdDeltas[i].MotionHash() != vw.motions[i].Hash()) {
			// 計算されていないか、ハッシュが一致しない場合は再計算
			if i < len(vw.motions) && vw.motions[i] != nil {
				vw.vmdDeltas[i] = deform.Deform(shared, vw.physics, modelRenderer.Model, vw.motions[i], nil, timeStep)
			}
		} else {
			// 物理計算のみ実行
			vw.vmdDeltas[i] = deform.DeformPhysics(shared, vw.physics, modelRenderer.Model, vw.motions[i], vw.vmdDeltas[i], timeStep)
		}

		mlog.IS("[%d][%d] Deform Model", vw.windowIndex, i)

		modelRenderer.Render(vw.shader, shared, vw.vmdDeltas[i])

		mlog.IS("[%d][%d] Render Model", vw.windowIndex, i)
	}

	// 物理デバッグ描画
	vw.physics.DrawDebugLines(vw.shader, shared.IsShowRigidBodyFront() || shared.IsShowRigidBodyBack(),
		shared.IsShowJoint(), shared.IsShowRigidBodyFront())

	// 深度解決
	vw.shader.GetMsaa().Resolve()
	// TODO
	// if vw.appState.IsShowOverride() && vw.windowIndex == 0 {
	// 	vw.drawOverride()
	// }
	vw.shader.GetMsaa().Unbind()

	vw.SwapBuffers()
}
