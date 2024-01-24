package mwidget

import (
	"embed"
	"image"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type ModelSet struct {
	Model *pmx.PmxModel
}

type GlWindow struct {
	glfw.Window
	ModelSets   []ModelSet
	Shader      *mgl.MShader
	WindowIndex int
}

func NewGlWindow(
	title string,
	width int,
	height int,
	windowIndex int,
	resourceFiles embed.FS,
	mainWindow *GlWindow,
) (*GlWindow, error) {
	if mainWindow == nil {
		// GLFW の初期化(最初の一回だけ)
		if err := glfw.Init(); err != nil {
			return nil, err
		}
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var mw *glfw.Window = nil
	if mainWindow != nil {
		mw = &mainWindow.Window
	}

	// ウィンドウの作成
	w, err := glfw.CreateWindow(width, height, title, nil, mw)
	if err != nil {
		return nil, err
	}
	w.MakeContextCurrent()
	iconImg, err := mutils.ReadIconFile(resourceFiles)
	if err == nil {
		w.SetIcon([]image.Image{iconImg})
	}

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		return nil, err
	}

	shader, err := mgl.NewMShader(width, height, resourceFiles)
	if err != nil {
		return nil, err
	}

	return &GlWindow{
		Window:      *w,
		ModelSets:   make([]ModelSet, 0),
		Shader:      shader,
		WindowIndex: windowIndex,
	}, nil
}

func (w *GlWindow) AddData(pmxModel *pmx.PmxModel) {
	// OpenGLコンテキストをこのウィンドウに設定
	w.MakeContextCurrent()

	// TODO: モーションも追加する
	pmxModel.Meshes = pmx.NewMeshes(pmxModel, w.WindowIndex)
	w.ModelSets = append(w.ModelSets, ModelSet{Model: pmxModel})

	// // OpenGLコンテキストをこのウィンドウに設定
	// w.MakeContextCurrent()

	// // 背景色をクリア
	// gl.ClearColor(0.0, 0.0, 0.0, 1.0) // 黒色でクリア
	// gl.Clear(gl.COLOR_BUFFER_BIT)

	// model := mgl32.Ident4()

	// // Load the texture
	// texture, err := newTexture("grid.png")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // Configure the vertex data
	// var vao uint32
	// gl.GenVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)

	// var vbo uint32
	// gl.GenBuffers(1, &vbo)
	// gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	// vertAttrib := uint32(gl.GetAttribLocation(w.Shader.ModelProgram, gl.Str("vert\x00")))
	// gl.EnableVertexAttribArray(vertAttrib)
	// gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

	// texCoordAttrib := uint32(gl.GetAttribLocation(w.Shader.ModelProgram, gl.Str("vertTexCoord\x00")))
	// gl.EnableVertexAttribArray(texCoordAttrib)
	// gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

	// // Configure global settings
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)
	// gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// angle := 0.0
	// previousTime := glfw.GetTime()

	// for !w.ShouldClose() {
	// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 	// Update
	// 	time := glfw.GetTime()
	// 	elapsed := time - previousTime
	// 	previousTime = time

	// 	angle += elapsed
	// 	model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

	// 	// Render
	// 	gl.UseProgram(w.Shader.ModelProgram)
	// 	gl.UniformMatrix4fv(w.Shader.ModelUniform, 1, false, &model[0])

	// 	gl.BindVertexArray(vao)

	// 	gl.ActiveTexture(gl.TEXTURE0)
	// 	gl.BindTexture(gl.TEXTURE_2D, texture)

	// 	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

	// 	// Maintenance
	// 	w.SwapBuffers()
	// 	glfw.PollEvents()
	// }
	// w.Close()
}

func (w *GlWindow) Draw() {
	// OpenGLコンテキストをこのウィンドウに設定
	w.MakeContextCurrent()

	// 背景色をクリア
	gl.ClearColor(0.0, 0.0, 0.0, 1.0) // 黒色でクリア
	gl.Clear(gl.COLOR_BUFFER_BIT)

}

func (w *GlWindow) Size() walk.Size {
	x, y := w.Window.GetSize()
	return walk.Size{Width: x, Height: y}
}

func (w *GlWindow) Close() {
	w.Shader.Delete()
	w.Window.Destroy()
	glfw.Terminate()
}

func (w *GlWindow) Run() {
	for !w.Window.ShouldClose() {
		w.Draw()
		w.Window.SwapBuffers()
		glfw.PollEvents()
	}
	w.Close()
}
