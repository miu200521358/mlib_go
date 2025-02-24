//go:build windows
// +build windows

package viewer

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type ViewWindow struct {
	*glfw.Window
	windowIndex int    // ウィンドウインデックス
	title       string // ウィンドウタイトル
	list        *ViewerList
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

	viewWindow := &ViewWindow{
		Window:      glWindow,
		windowIndex: windowIndex,
		title:       title,
		list:        list,
	}

	glWindow.SetCloseCallback(viewWindow.closeCallback)
	// glWindow.SetScrollCallback(viewWindow.scrollCallback)
	// glWindow.SetKeyCallback(viewWindow.keyCallback)
	// glWindow.SetMouseButtonCallback(viewWindow.mouseCallback)
	// glWindow.SetCursorPosCallback(viewWindow.cursorPosCallback)
	// glWindow.SetSizeCallback(viewWindow.resizeCallback)
	// glWindow.SetFramebufferSizeCallback(viewWindow.resizeCallback)

	if !isProd {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                        // 同期的なデバッグ出力有効
		gl.DebugMessageCallback(viewWindow.debugMessageCallback, nil) // デバッグコールバック
	}

	viewWindow.SetPos(positionX, positionY)

	return viewWindow, nil
}
