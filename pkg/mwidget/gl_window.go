package mwidget

import (
	"embed"
	"fmt"
	"image"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
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
	ModelSets           []ModelSet
	Shader              *mgl.MShader
	WindowIndex         int
	resourceFiles       embed.FS
	prevCursorPos       *mmath.MVec2
	yaw                 float64
	pitch               float64
	middleButtonPressed bool
	rightButtonPressed  bool
	updatedPrev         bool
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
	iconImg, err := mutils.LoadIconFile(resourceFiles)
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

	glWindow := GlWindow{
		Window:              *w,
		ModelSets:           make([]ModelSet, 0),
		Shader:              shader,
		WindowIndex:         windowIndex,
		resourceFiles:       resourceFiles,
		prevCursorPos:       &mmath.MVec2{0, 0},
		yaw:                 89.0,
		pitch:               0.0,
		middleButtonPressed: false,
		rightButtonPressed:  false,
		updatedPrev:         false,
	}

	w.SetScrollCallback(glWindow.handleScrollEvent)
	w.SetMouseButtonCallback(glWindow.handleMouseButtonEvent)
	w.SetCursorPosCallback(glWindow.handleCursorPosEvent)
	w.SetKeyCallback(glWindow.handleKeyEvent)
	w.SetCloseCallback(glWindow.Close)
	glWindow.Draw()

	return &glWindow, nil
}

func (w *GlWindow) Close(window *glfw.Window) {
	w.Shader.Delete()
	for _, modelSet := range w.ModelSets {
		modelSet.Model.Meshes.Delete()
	}
	w.Window.Destroy()
	if w.WindowIndex == 0 {
		glfw.Terminate()
	}
}

func (w *GlWindow) handleKeyEvent(
	window *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if action != glfw.Press ||
		!(key == glfw.KeyKP0 ||
			key == glfw.KeyKP2 ||
			key == glfw.KeyKP4 ||
			key == glfw.KeyKP5 ||
			key == glfw.KeyKP6 ||
			key == glfw.KeyKP8) {
		return
	}

	w.Reset()

	switch key {
	case glfw.KeyKP0: // 下面から
		w.yaw = 90
		w.pitch = 90
	case glfw.KeyKP2: // 正面から
		w.yaw = 90
		w.pitch = 0
	case glfw.KeyKP4: // 左面から
		w.yaw = 180
		w.pitch = 0
	case glfw.KeyKP5: // 上面から
		w.yaw = 90
		w.pitch = -90
	case glfw.KeyKP6: // 右面から
		w.yaw = 0
		w.pitch = 0
	case glfw.KeyKP8: // 背面から
		w.yaw = -90
		w.pitch = 0
	default:
		return
	}

	// カメラの新しい位置を計算
	radius := mgl.INITIAL_CAMERA_POSITION_Z

	// 球面座標系をデカルト座標系に変換
	cameraX := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw))
	cameraY := radius * math.Sin(mgl64.DegToRad(w.pitch))
	cameraZ := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw))

	// カメラ位置を更新
	w.Shader.CameraPosition.SetX(cameraX)
	w.Shader.CameraPosition.SetY(mgl.INITIAL_CAMERA_POSITION_Y + cameraY)
	w.Shader.CameraPosition.SetZ(cameraZ)
}

func (w *GlWindow) handleScrollEvent(window *glfw.Window, xoff float64, yoff float64) {
	if yoff > 0 {
		w.Shader.FieldOfViewAngle -= 1.0
		if w.Shader.FieldOfViewAngle < 5.0 {
			w.Shader.FieldOfViewAngle = 5.0
		}
	} else if yoff < 0 {
		w.Shader.FieldOfViewAngle += 1.0
	}
}

func (w *GlWindow) handleMouseButtonEvent(
	window *glfw.Window,
	button glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey,
) {
	if button == glfw.MouseButtonMiddle {
		if action == glfw.Press {
			w.middleButtonPressed = true
			w.updatedPrev = false
		} else if action == glfw.Release {
			w.middleButtonPressed = false
		}
	} else if button == glfw.MouseButtonRight {
		if action == glfw.Press {
			w.rightButtonPressed = true
			w.updatedPrev = false
		} else if action == glfw.Release {
			w.rightButtonPressed = false
		}
	}
}

func (w *GlWindow) handleCursorPosEvent(window *glfw.Window, xpos float64, ypos float64) {
	fmt.Printf("[start] yaw %.8f, pitch %.8f, CameraPosition: %s, LookAtCenterPosition: %s\n",
		w.yaw, w.pitch, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())

	if !w.updatedPrev {
		w.prevCursorPos.SetX(xpos)
		w.prevCursorPos.SetY(ypos)
		w.updatedPrev = true
		return
	}

	if w.rightButtonPressed {
		// 右クリックはカメラ中心をそのままにカメラ位置を変える
		xOffset := (w.prevCursorPos.GetX() - xpos) * 0.1
		yOffset := (w.prevCursorPos.GetY() - ypos) * 0.1

		// 方位角と仰角を更新
		w.yaw += xOffset
		w.pitch += yOffset

		// 仰角の制限（水平面より上下に行き過ぎないようにする）
		if w.pitch > 89.9 {
			w.pitch = 89.9
		} else if w.pitch < -89.9 {
			w.pitch = -89.9
		}

		// 方位角の制限（360度を超えないようにする）
		if w.yaw > 360.0 {
			w.yaw -= 360.0
		} else if w.yaw < -360.0 {
			w.yaw += 360.0
		}

		// 球面座標系をデカルト座標系に変換
		radius := float64(-w.Shader.CameraPosition.Sub(w.Shader.LookAtCenterPosition).Length())
		cameraX := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw))
		cameraY := radius * math.Sin(mgl64.DegToRad(w.pitch))
		cameraZ := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw))

		// カメラ位置を更新
		w.Shader.CameraPosition.SetX(cameraX)
		w.Shader.CameraPosition.SetY(mgl.INITIAL_CAMERA_POSITION_Y + cameraY)
		w.Shader.CameraPosition.SetZ(cameraZ)
		fmt.Printf("xOffset %.8f, yOffset %.8f, CameraPosition: %s, LookAtCenterPosition: %s\n",
			xOffset, yOffset, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())
	} else if w.middleButtonPressed {
		// 中クリックはカメラ中心とカメラ位置を一緒に動かす
		// カメラの方向ベクトルと上方ベクトルを計算
		cameraDirection := mgl64.Vec3{
			math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw)),
			math.Sin(mgl64.DegToRad(w.pitch)),
			math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw)),
		}.Normalize()
		upVector := mgl64.Vec3{0, 1, 0}

		// カメラの横方向と縦方向のベクトルを計算
		rightVector := cameraDirection.Cross(upVector).Normalize()
		upCameraVector := rightVector.Cross(cameraDirection).Normalize()

		// カメラがモデルの側面を向いているかを確認
		var horizontalSign float64 = -1
		if (w.yaw > 160 && w.yaw < 200) || (w.yaw > 340 && w.yaw < 380) {
			horizontalSign = 1
		}

		xOffset := (w.prevCursorPos.GetX() - xpos) * 0.07 * horizontalSign
		yOffset := (w.prevCursorPos.GetY() - ypos) * 0.07

		// カメラの位置を更新（カメラの方向を考慮）
		w.Shader.CameraPosition.SetX(
			w.Shader.CameraPosition.GetX() +
				float64(rightVector.X())*xOffset - float64(upCameraVector.X())*yOffset)
		w.Shader.CameraPosition.SetY(
			w.Shader.CameraPosition.GetY() +
				float64(rightVector.Y())*xOffset - float64(upCameraVector.Y())*yOffset)
		w.Shader.CameraPosition.SetZ(
			w.Shader.CameraPosition.GetZ() +
				float64(rightVector.Z())*xOffset - float64(upCameraVector.Z())*yOffset)

		// カメラの中心位置も同様に更新
		w.Shader.LookAtCenterPosition.SetX(
			w.Shader.LookAtCenterPosition.GetX() +
				float64(rightVector.X())*xOffset -
				float64(upCameraVector.X())*yOffset)
		w.Shader.LookAtCenterPosition.SetY(
			w.Shader.LookAtCenterPosition.GetY() +
				float64(rightVector.Y())*xOffset -
				float64(upCameraVector.Y())*yOffset)
		w.Shader.LookAtCenterPosition.SetZ(
			w.Shader.LookAtCenterPosition.GetZ() +
				float64(rightVector.Z())*xOffset -
				float64(upCameraVector.Z())*yOffset)

		fmt.Printf("xOffset %.8f, yOffset %.8f, CameraPosition: %s, LookAtCenterPosition: %s\n",
			xOffset, yOffset, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())
	}

	w.prevCursorPos.SetX(xpos)
	w.prevCursorPos.SetY(ypos)
}

func (w *GlWindow) Reset() {
	// カメラとかリセット
	w.Shader.Reset()
	w.prevCursorPos = &mmath.MVec2{0, 0}
	w.yaw = 89.0
	w.pitch = 0.0
	w.middleButtonPressed = false
	w.rightButtonPressed = false

}

func (w *GlWindow) AddData(pmxModel *pmx.PmxModel) {
	// OpenGLコンテキストをこのウィンドウに設定
	w.MakeContextCurrent()
	w.Reset()

	// TODO: モーションも追加する
	pmxModel.InitializeDraw(w.WindowIndex, w.resourceFiles)
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

func (w *GlWindow) Run() {
	for !w.ShouldClose() {
		// 深度バッファのクリア
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, program := range w.Shader.GetPrograms() {
			// プログラムの切り替え
			gl.UseProgram(program)
			// カメラの再計算
			projection := mgl32.Perspective(
				mgl32.DegToRad(w.Shader.FieldOfViewAngle),
				float32(w.Shader.Width)/float32(w.Shader.Height),
				w.Shader.NearPlane,
				w.Shader.FarPlane,
			)
			projectionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_PROJECTION_MATRIX))
			gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

			// カメラの位置
			cameraPosition := w.Shader.CameraPosition.Mgl()
			cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CAMERA_POSITION))
			gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

			// カメラの中心
			lookAtCenter := w.Shader.LookAtCenterPosition.Mgl()
			camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})
			cameraUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_MATRIX))
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		}

		w.Draw()

		// Maintenance
		w.SwapBuffers()
		glfw.PollEvents()
	}
	w.Close(&w.Window)
}
