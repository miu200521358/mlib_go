//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"image"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex int    // ウィンドウインデックス
	title       string // ウィンドウタイトル
	list        *ViewerList
	shader      rendering.IShader
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
	// glWindow.SetKeyCallback(vw.keyCallback)
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

func (vw *ViewWindow) Render(shared *state.SharedState) {
	w, h := vw.GetSize()
	if w == 0 && h == 0 {
		// ウィンドウが最小化されている場合は描画しない
		return
	}

	vw.MakeContextCurrent()

	// 深度バッファのクリア
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 隠面消去
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// 床描画
	vw.shader.DrawFloor()

	vw.SwapBuffers()
}
