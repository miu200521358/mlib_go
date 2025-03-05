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
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex  int               // ウィンドウインデックス
	title        string            // ウィンドウタイトル
	shiftPressed bool              // Shiftキー押下フラグ
	ctrlPressed  bool              // Ctrlキー押下フラグ
	list         *ViewerList       // ビューワーリスト
	shader       rendering.IShader // シェーダー
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

	shaderFactory := mgl.NewMShaderFactory()

	// シェーダー初期化
	shader, err := shaderFactory.CreateShader(width, height)
	if err != nil {
		return nil, err
	}

	gl.Viewport(0, 0, int32(width), int32(height))

	vw := &ViewWindow{
		Window:      glWindow,
		windowIndex: windowIndex,
		title:       title,
		list:        list,
		shader:      shader,
	}

	glWindow.SetCloseCallback(vw.closeCallback)
	// glWindow.SetScrollCallback(vw.scrollCallback)
	glWindow.SetKeyCallback(vw.keyCallback)
	// glWindow.SetMouseButtonCallback(vw.mouseCallback)
	// glWindow.SetCursorPosCallback(vw.cursorPosCallback)
	// glWindow.SetSizeCallback(vw.resizeCallback)
	// glWindow.SetFramebufferSizeCallback(vw.resizeCallback)

	if !isProd {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                // 同期的なデバッグ出力有効
		gl.DebugMessageCallback(vw.debugMessageCallback, nil) // デバッグコールバック
	}

	vw.SetPos(positionX, positionY)

	return vw, nil
}

func (vw *ViewWindow) resetCameraPosition(yaw, pitch float64) {
	// 球面座標系をデカルト座標系に変換
	radius := math.Abs(float64(rendering.InitialCameraPositionZ))

	// 四元数を使ってカメラの方向を計算
	yawRad := mgl64.DegToRad(yaw)
	pitchRad := mgl64.DegToRad(pitch)
	orientation := mmath.NewMQuaternionFromAxisAngles(mmath.MVec3UnitY, yawRad).Mul(
		mmath.NewMQuaternionFromAxisAngles(mmath.MVec3UnitX, pitchRad))
	forwardXYZ := orientation.MulVec3(mmath.MVec3UnitZNeg).MulScalar(radius)

	// カメラ位置を更新
	cam := rendering.NewDefaultCamera(vw.GetSize())
	cam.Position.X = forwardXYZ.X
	cam.Position.Y = rendering.InitialCameraPositionY + forwardXYZ.Y
	cam.Position.Z = forwardXYZ.Z
	vw.shader.SetCamera(cam)
}

func (vw *ViewWindow) Render(shared *state.SharedState) {
	w, h := vw.GetSize()
	if w == 0 && h == 0 {
		// ウィンドウが最小化されている場合は描画しない
		return
	}

	vw.MakeContextCurrent()

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
	vw.shader.UpdateCamera(vw.shader.GetCamera())

	// 床描画
	vw.shader.DrawFloor()

	// 深度解決
	vw.shader.GetMsaa().Resolve()
	// if vw.appState.IsShowOverride() && vw.windowIndex == 0 {
	// 	vw.drawOverride()
	// }
	vw.shader.GetMsaa().Unbind()

	vw.SwapBuffers()
}
