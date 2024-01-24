package mwidget

import (
	"embed"
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/pmx"
)

type ModelSet struct {
	Model *pmx.PmxModel
}

func (ms *ModelSet) Draw(shader *mgl.MShader, windowIndex int) {
	// TODO: モーション計算
	boneMatrixes := []mgl32.Mat4{
		mgl32.Ident4(),
	}
	ms.Model.Draw(shader, boneMatrixes, windowIndex)
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

	// デバッグコールバックを設定します。
	gl.DebugMessageCallback(func(
		source uint32,
		glType uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer,
	) {
		switch severity {
		case gl.DEBUG_SEVERITY_HIGH:
			fmt.Printf("[HIGH] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
			panic("critical OpenGL error")
		case gl.DEBUG_SEVERITY_MEDIUM:
			fmt.Printf("[MEDIUM] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
		case gl.DEBUG_SEVERITY_LOW:
			fmt.Printf("[LOW] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
			// case gl.DEBUG_SEVERITY_NOTIFICATION:
			// 	fmt.Printf("[NOTIFICATION] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
			// 		source, glType, severity, message)
		}
	}, gl.Ptr(nil))
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS) // 同期的なデバッグ出力を有効にします。

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
	pmxModel.InitializeDraw(w.WindowIndex)
	w.ModelSets = append(w.ModelSets, ModelSet{Model: pmxModel})
}

func (w *GlWindow) ClearData() {
	for _, modelSet := range w.ModelSets {
		modelSet.Model.Meshes.Delete()
	}
	w.ModelSets = make([]ModelSet, 0)
}

func (w *GlWindow) Draw() {
	// OpenGLコンテキストをこのウィンドウに設定
	w.MakeContextCurrent()

	// 背景色をクリア
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// モデル描画
	for _, modelSet := range w.ModelSets {
		modelSet.Draw(w.Shader, w.WindowIndex)
	}
}

func (w *GlWindow) Size() walk.Size {
	x, y := w.Window.GetSize()
	return walk.Size{Width: x, Height: y}
}

func (w *GlWindow) Close() {
	w.Shader.Delete()
	for _, modelSet := range w.ModelSets {
		modelSet.Model.Meshes.Delete()
	}
	w.Window.Destroy()
	glfw.Terminate()
}

func (w *GlWindow) Run() {
	angle := 0.0
	previousTime := glfw.GetTime()
	modelUniform := gl.GetUniformLocation(w.Shader.ModelProgram, gl.Str("model\x00"))

	for !w.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		model := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		gl.UseProgram(w.Shader.ModelProgram)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		w.Draw()

		// Maintenance
		w.SwapBuffers()
		glfw.PollEvents()
	}
	w.Close()
}

// func (w *GlWindow) Run2() {
// 	// OpenGLコンテキストをこのウィンドウに設定
// 	w.MakeContextCurrent()

// 	w.Shader.UseModelProgram()
// 	program := w.Shader.ModelProgram

// 	model := mgl32.Ident4()
// 	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
// 	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

// 	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
// 	gl.Uniform1i(textureUniform, 0)

// 	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

// 	// Load the texture
// 	texture, err := newTexture("grid.png")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// Configure the vertex data
// 	var vao uint32
// 	gl.GenVertexArrays(1, &vao)
// 	gl.BindVertexArray(vao)

// 	var vbo uint32
// 	gl.GenBuffers(1, &vbo)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(w.faces)*4, gl.Ptr(w.faces), gl.STATIC_DRAW)

// 	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
// 	gl.EnableVertexAttribArray(vertAttrib)
// 	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 5*4, 0)

// 	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
// 	gl.EnableVertexAttribArray(texCoordAttrib)
// 	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 5*4, 3*4)

// 	// Configure global settings
// 	gl.Enable(gl.DEPTH_TEST)
// 	gl.DepthFunc(gl.LESS)
// 	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

// 	angle := 0.0
// 	previousTime := glfw.GetTime()

// 	for !w.ShouldClose() {
// 		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

// 		// Update
// 		time := glfw.GetTime()
// 		elapsed := time - previousTime
// 		previousTime = time

// 		angle += elapsed
// 		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

// 		// Render
// 		gl.UseProgram(program)
// 		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

// 		gl.BindVertexArray(vao)

// 		gl.ActiveTexture(gl.TEXTURE0)
// 		gl.BindTexture(gl.TEXTURE_2D, texture)

// 		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(w.faces)*3))

// 		// Maintenance
// 		w.SwapBuffers()
// 		glfw.PollEvents()
// 	}
// 	w.Close()
// }

// func newTexture(file string) (uint32, error) {
// 	imgFile, err := os.Open(file)
// 	if err != nil {
// 		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
// 	}
// 	img, _, err := image.Decode(imgFile)
// 	if err != nil {
// 		return 0, err
// 	}

// 	rgba := image.NewRGBA(img.Bounds())
// 	if rgba.Stride != rgba.Rect.Size().X*4 {
// 		return 0, fmt.Errorf("unsupported stride")
// 	}
// 	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

// 	var texture uint32
// 	gl.GenTextures(1, &texture)
// 	gl.ActiveTexture(gl.TEXTURE0)
// 	gl.BindTexture(gl.TEXTURE_2D, texture)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
// 	gl.TexImage2D(
// 		gl.TEXTURE_2D,
// 		0,
// 		gl.RGBA,
// 		int32(rgba.Rect.Size().X),
// 		int32(rgba.Rect.Size().Y),
// 		0,
// 		gl.RGBA,
// 		gl.UNSIGNED_BYTE,
// 		gl.Ptr(rgba.Pix))

// 	return texture, nil
// }
