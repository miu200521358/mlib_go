//go:build windows
// +build windows

package viewer

import (
	"fmt"
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

	vw.loadModelRenderers(vw.list.shared)

	vw.loadMotions(vw.list.shared)

	for i, modelRenderer := range vw.modelRenderers {
		vw.vmdDeltas[i] = deform.Deform(vw.list.shared, vw.physics, modelRenderer.Model, vw.motions[i], vw.vmdDeltas[i], timeStep)

		// mlog.IS("[%d][%d] Deform [%f]", vw.windowIndex, i, vw.list.shared.Frame())

		modelRenderer.Render(vw.shader, vw.list.shared, vw.vmdDeltas[i])

		// mlog.IS("[%d][%d] Render [%f]", vw.windowIndex, i, vw.list.shared.Frame())
	}

	// 物理デバッグ描画
	vw.physics.DrawDebugLines(vw.shader, vw.list.shared.IsShowRigidBodyFront() || vw.list.shared.IsShowRigidBodyBack(),
		vw.list.shared.IsShowJoint(), vw.list.shared.IsShowRigidBodyFront())

	// 深度解決
	vw.shader.GetMsaa().Resolve()
	// TODO
	// if vw.appState.IsShowOverride() && vw.windowIndex == 0 {
	// 	vw.drawOverride()
	// }
	vw.shader.GetMsaa().Unbind()

	vw.SwapBuffers()
}
