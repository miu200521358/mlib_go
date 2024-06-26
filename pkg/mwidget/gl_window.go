//go:build windows
// +build windows

package mwidget

import (
	"embed"
	"image"
	"math"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mview"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

type ModelSet struct {
	model      *pmx.PmxModel
	motion     *vmd.VmdMotion
	prevDeltas *vmd.VmdDeltas
}

// 直角の定数値
const RIGHT_ANGLE = 89.9

type GlWindow struct {
	*glfw.Window
	ModelSets           []ModelSet
	Shader              *mview.MShader
	WindowIndex         int
	resourceFiles       embed.FS
	prevCursorPos       *mmath.MVec2
	yaw                 float64
	pitch               float64
	Physics             *mphysics.MPhysics
	middleButtonPressed bool
	rightButtonPressed  bool
	updatedPrev         bool
	shiftPressed        bool
	ctrlPressed         bool
	running             bool
	drawing             bool
	playing             bool
	VisibleBones        map[pmx.BoneFlag]bool
	VisibleNormal       bool
	EnablePhysics       bool
	EnableFrameDrop     bool
	frame               float64
	prevFrame           int
	motionPlayer        *MotionPlayer
	width               int
	height              int
	floor               *MFloor
}

func NewGlWindow(
	title string,
	width int,
	height int,
	windowIndex int,
	resourceFiles embed.FS,
	mainWindow *GlWindow,
	motionPlayer *MotionPlayer,
) (*GlWindow, error) {
	if mainWindow == nil {
		// GLFW の初期化(最初の一回だけ)
		if err := glfw.Init(); err != nil {
			return nil, err
		}
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var mw *glfw.Window = nil
	if mainWindow != nil {
		mw = mainWindow.Window
	}

	// ウィンドウの作成
	w, err := glfw.CreateWindow(width, height, title, nil, mw)
	if err != nil {
		return nil, err
	}
	w.MakeContextCurrent()
	w.SetInputMode(glfw.StickyKeysMode, glfw.True)

	iconImg, err := mconfig.LoadIconFile(resourceFiles)
	if err == nil {
		w.SetIcon([]image.Image{*iconImg})
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
			mlog.E("[HIGH] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
			panic("critical OpenGL error")
		case gl.DEBUG_SEVERITY_MEDIUM:
			mlog.D("[MEDIUM] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
		case gl.DEBUG_SEVERITY_LOW:
			mlog.D("[LOW] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
			// case gl.DEBUG_SEVERITY_NOTIFICATION:
			// 	mlog.D("[NOTIFICATION] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
			// 		source, glType, severity, message)
		}
	}, gl.Ptr(nil))
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS) // 同期的なデバッグ出力を有効にします。

	shader, err := mview.NewMShader(width, height, resourceFiles)
	if err != nil {
		return nil, err
	}

	glWindow := GlWindow{
		Window:              w,
		ModelSets:           make([]ModelSet, 0),
		Shader:              shader,
		WindowIndex:         windowIndex,
		resourceFiles:       resourceFiles,
		prevCursorPos:       &mmath.MVec2{0, 0},
		yaw:                 RIGHT_ANGLE,
		pitch:               0.0,
		Physics:             mphysics.NewMPhysics(shader),
		middleButtonPressed: false,
		rightButtonPressed:  false,
		updatedPrev:         false,
		shiftPressed:        false,
		ctrlPressed:         false,
		VisibleBones:        make(map[pmx.BoneFlag]bool, 0),
		VisibleNormal:       false,
		running:             false,
		drawing:             false,
		playing:             false, // 最初は再生OFF
		EnablePhysics:       true,  // 最初は物理ON
		EnableFrameDrop:     true,  // 最初はドロップON
		frame:               0,
		prevFrame:           0,
		motionPlayer:        motionPlayer,
		width:               width,
		height:              height,
		floor:               newMFloor(),
	}

	w.SetScrollCallback(glWindow.handleScrollEvent)
	w.SetMouseButtonCallback(glWindow.handleMouseButtonEvent)
	w.SetCursorPosCallback(glWindow.handleCursorPosEvent)
	w.SetKeyCallback(glWindow.handleKeyEvent)
	w.SetCloseCallback(glWindow.Close)
	w.SetSizeCallback(glWindow.resize)
	w.SetFramebufferSizeCallback(glWindow.resizeBuffer)

	return &glWindow, nil
}

func (w *GlWindow) resizeBuffer(window *glfw.Window, width int, height int) {
	w.width = width
	w.height = height
	if width > 0 && height > 0 {
		gl.Viewport(0, 0, int32(width), int32(height))
	}
}

func (w *GlWindow) resize(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	w.Shader.Resize(width, height)
}

func (w *GlWindow) Play(p bool) {
	w.playing = p
}

func (w *GlWindow) GetFrame() int {
	return int(w.frame)
}

func (w *GlWindow) SetFrame(f int) {
	w.frame = float64(f)
	w.prevFrame = f

	for i := range w.ModelSets {
		// 前のデフォーム情報をクリア
		w.ModelSets[i].prevDeltas = nil
	}
}

func (w *GlWindow) Close(window *glfw.Window) {
	w.running = false
	w.drawing = false
	w.Shader.Delete()
	for _, modelSet := range w.ModelSets {
		modelSet.model.Delete(w.Physics)
	}
	if w.WindowIndex == 0 {
		defer glfw.Terminate()
		defer walk.App().Exit(0)
	}
}

func (w *GlWindow) handleKeyEvent(
	window *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if !(action == glfw.Press || action == glfw.Repeat) ||
		!(key == glfw.KeyKP0 ||
			key == glfw.KeyKP2 ||
			key == glfw.KeyKP4 ||
			key == glfw.KeyKP5 ||
			key == glfw.KeyKP6 ||
			key == glfw.KeyKP8 ||
			key == glfw.KeyLeftShift ||
			key == glfw.KeyRightShift ||
			key == glfw.KeyLeftControl ||
			key == glfw.KeyRightControl ||
			key == glfw.KeyLeft ||
			key == glfw.KeyRight ||
			key == glfw.KeyUp ||
			key == glfw.KeyDown) {
		return
	}

	if w.motionPlayer != nil {
		if key == glfw.KeyRight || key == glfw.KeyUp {
			w.motionPlayer.SetValue(w.GetFrame() + 1)
			return
		} else if key == glfw.KeyLeft || key == glfw.KeyDown {
			w.motionPlayer.SetValue(w.GetFrame() - 1)
			return
		}
	}

	if key == glfw.KeyLeftShift || key == glfw.KeyRightShift {
		if action == glfw.Press {
			w.shiftPressed = true
		} else if action == glfw.Release {
			w.shiftPressed = false
		}
		return
	} else if key == glfw.KeyLeftControl || key == glfw.KeyRightControl {
		if action == glfw.Press {
			w.ctrlPressed = true
		} else if action == glfw.Release {
			w.ctrlPressed = false
		}
		return
	}

	w.Reset()

	switch key {
	case glfw.KeyKP0: // 下面から
		w.yaw = RIGHT_ANGLE
		w.pitch = RIGHT_ANGLE
	case glfw.KeyKP2: // 正面から
		w.yaw = RIGHT_ANGLE
		w.pitch = 0
	case glfw.KeyKP4: // 左面から
		w.yaw = 180
		w.pitch = 0
	case glfw.KeyKP5: // 上面から
		w.yaw = RIGHT_ANGLE
		w.pitch = -RIGHT_ANGLE
	case glfw.KeyKP6: // 右面から
		w.yaw = 0
		w.pitch = 0
	case glfw.KeyKP8: // 背面から
		w.yaw = -RIGHT_ANGLE
		w.pitch = 0
	default:
		return
	}

	// カメラの新しい位置を計算
	radius := mview.INITIAL_CAMERA_POSITION_Z

	// 球面座標系をデカルト座標系に変換
	cameraX := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw))
	cameraY := radius * math.Sin(mgl64.DegToRad(w.pitch))
	cameraZ := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw))

	// カメラ位置を更新
	w.Shader.CameraPosition.SetX(cameraX)
	w.Shader.CameraPosition.SetY(mview.INITIAL_CAMERA_POSITION_Y + cameraY)
	w.Shader.CameraPosition.SetZ(cameraZ)
}

func (w *GlWindow) handleScrollEvent(window *glfw.Window, xoff float64, yoff float64) {
	ratio := float32(1.0)
	if w.shiftPressed {
		ratio *= 10
	} else if w.ctrlPressed {
		ratio *= 0.1
	}

	if yoff > 0 {
		w.Shader.FieldOfViewAngle -= ratio
		if w.Shader.FieldOfViewAngle < 5.0 {
			w.Shader.FieldOfViewAngle = 5.0
		}
	} else if yoff < 0 {
		w.Shader.FieldOfViewAngle += ratio
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
	// mlog.D("[start] yaw %.8f, pitch %.8f, CameraPosition: %s, LookAtCenterPosition: %s\n",
	// 	w.yaw, w.pitch, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())

	if !w.updatedPrev {
		w.prevCursorPos.SetX(xpos)
		w.prevCursorPos.SetY(ypos)
		w.updatedPrev = true
		return
	}

	if w.rightButtonPressed {
		ratio := 0.1
		if w.shiftPressed {
			ratio *= 10
		} else if w.ctrlPressed {
			ratio *= 0.1
		}

		// 右クリックはカメラ中心をそのままにカメラ位置を変える
		xOffset := (w.prevCursorPos.GetX() - xpos) * ratio
		yOffset := (w.prevCursorPos.GetY() - ypos) * ratio

		// 方位角と仰角を更新
		w.yaw += xOffset
		w.pitch += yOffset

		// 仰角の制限（水平面より上下に行き過ぎないようにする）
		if w.pitch > RIGHT_ANGLE {
			w.pitch = RIGHT_ANGLE
		} else if w.pitch < -RIGHT_ANGLE {
			w.pitch = -RIGHT_ANGLE
		}

		// 方位角の制限（360度を超えないようにする）
		if w.yaw > 360.0 {
			w.yaw -= 360.0
		} else if w.yaw < -360.0 {
			w.yaw += 360.0
		}

		// 球面座標系をデカルト座標系に変換
		// radius := float64(-w.Shader.CameraPosition.Sub(w.Shader.LookAtCenterPosition).Length())
		radius := float64(mview.INITIAL_CAMERA_POSITION_Z)
		cameraX := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Cos(mgl64.DegToRad(w.yaw))
		cameraY := radius * math.Sin(mgl64.DegToRad(w.pitch))
		cameraZ := radius * math.Cos(mgl64.DegToRad(w.pitch)) * math.Sin(mgl64.DegToRad(w.yaw))

		// カメラ位置を更新
		w.Shader.CameraPosition.SetX(cameraX)
		w.Shader.CameraPosition.SetY(mview.INITIAL_CAMERA_POSITION_Y + cameraY)
		w.Shader.CameraPosition.SetZ(cameraZ)
		// mlog.D("xOffset %.8f, yOffset %.8f, CameraPosition: %s, LookAtCenterPosition: %s\n",
		// 	xOffset, yOffset, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())
	} else if w.middleButtonPressed {
		ratio := 0.07
		if w.shiftPressed {
			ratio *= 10
		} else if w.ctrlPressed {
			ratio *= 0.1
		}
		// 中ボタンが押された場合の処理
		if w.middleButtonPressed {
			ratio := 0.07
			if w.shiftPressed {
				ratio *= 10
			} else if w.ctrlPressed {
				ratio *= 0.1
			}

			xOffset := (w.prevCursorPos.GetX() - xpos) * ratio
			yOffset := (w.prevCursorPos.GetY() - ypos) * ratio

			// カメラの向きに基づいて移動方向を計算
			forward := w.Shader.LookAtCenterPosition.Subed(w.Shader.CameraPosition)
			right := forward.Cross(mmath.MVec3UnitY).Normalize()
			up := right.Cross(forward.Normalize()).Normalize()

			// 上下移動のベクトルを計算
			upMovement := up.MulScalar(-yOffset)
			// 左右移動のベクトルを計算
			rightMovement := right.MulScalar(-xOffset)

			// 移動ベクトルを合成してカメラ位置と中心を更新
			movement := upMovement.Add(rightMovement)
			w.Shader.CameraPosition.Add(movement)
			w.Shader.LookAtCenterPosition.Add(movement)
		}
	}

	w.prevCursorPos.SetX(xpos)
	w.prevCursorPos.SetY(ypos)
}

func (w *GlWindow) ResetPhysics() {
	for _, modelSet := range w.ModelSets {
		modelSet.model.DeletePhysics(w.Physics)
		modelSet.model.InitPhysics(w.Physics)
	}
}

func (w *GlWindow) Reset() {
	// カメラとかリセット
	w.Shader.Reset()
	w.prevCursorPos = &mmath.MVec2{0, 0}
	w.yaw = RIGHT_ANGLE
	w.pitch = 0.0
	w.middleButtonPressed = false
	w.rightButtonPressed = false

}

func (w *GlWindow) AddData(pmxModel *pmx.PmxModel, vmdMotion *vmd.VmdMotion) {
	w.drawing = true
	pmxModel.InitDraw(len(w.ModelSets), w.resourceFiles)
	pmxModel.InitPhysics(w.Physics)
	w.ModelSets = append(w.ModelSets, ModelSet{model: pmxModel, motion: vmdMotion, prevDeltas: nil})
}

func (w *GlWindow) ClearData() {
	for _, modelSet := range w.ModelSets {
		modelSet.model.DeletePhysics(w.Physics)
	}
	w.ModelSets = make([]ModelSet, 0)
	w.frame = 0
}

func (w *GlWindow) Size() walk.Size {
	x, y := w.Window.GetSize()
	return walk.Size{Width: x, Height: y}
}

func (w *GlWindow) Run() {
	if w == nil || w.running {
		return
	}

	prevTime := glfw.GetTime()
	w.prevFrame = 0
	w.running = true

	for w.IsRunning() {
		w.MakeContextCurrent()

		if w.width == 0 || w.height == 0 {
			// ウィンドウが最小化されている場合は描画をスキップ(フレームも進めない)
			prevTime = glfw.GetTime()

			glfw.PollEvents()
			continue
		}

		if w.playing && w.motionPlayer != nil && w.frame >= w.motionPlayer.FrameEdit.MaxValue() {
			w.frame = 0
			w.prevFrame = 0
			w.motionPlayer.SetValue(int(w.frame))
		}

		// 深度バッファのクリア
		gl.ClearColor(0.7, 0.7, 0.7, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// 隠面消去
		// https://learnopengl.com/Advanced-OpenGL/Depth-testing
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)

		// ブレンディングを有効にする
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		// カメラの再計算
		projection := mgl32.Perspective(
			mgl32.DegToRad(w.Shader.FieldOfViewAngle),
			float32(w.Shader.Width)/float32(w.Shader.Height),
			w.Shader.NearPlane,
			w.Shader.FarPlane,
		)

		// カメラの位置
		cameraPosition := w.Shader.CameraPosition.GL()

		// カメラの中心
		lookAtCenter := w.Shader.LookAtCenterPosition.GL()
		camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})

		for _, program := range w.Shader.GetPrograms() {
			if !w.IsRunning() {
				break
			}

			// プログラムの切り替え
			gl.UseProgram(program)

			// カメラの再計算
			projectionUniform := gl.GetUniformLocation(program, gl.Str(mview.SHADER_MODEL_VIEW_PROJECTION_MATRIX))
			gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

			// カメラの位置
			cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(mview.SHADER_CAMERA_POSITION))
			gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

			// カメラの中心
			cameraUniform := gl.GetUniformLocation(program, gl.Str(mview.SHADER_MODEL_VIEW_MATRIX))
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

			gl.UseProgram(0)
		}

		frameTime := glfw.GetTime()
		elapsed := float32(frameTime - prevTime)

		if !w.EnableFrameDrop {
			// フレームドロップOFFの場合、最大1Fずつ
			elapsed = mmath.ClampFloat32(elapsed, 0, w.Physics.Spf)
		}

		if !w.IsRunning() {
			break
		}

		if w.playing {
			// 経過秒数をキーフレームの進捗具合に合わせて調整
			w.frame += float64(elapsed * w.Physics.Fps)
			mlog.V("previousTime=%.8f, time=%.8f, elapsed=%.8f, frame=%.8f", prevTime, frameTime, elapsed, w.frame)
		}

		prevTime = frameTime

		isDeform := false
		if int(w.frame) > w.prevFrame {
			isDeform = true
			// フレーム番号上書き
			w.prevFrame = int(w.frame)
			if w.playing && w.motionPlayer != nil {
				w.motionPlayer.SetValue(int(w.frame))
			}
		}

		// mlog.Memory("GL.Run[4]")
		w.Shader.Msaa.Bind()

		// 床平面を描画
		w.drawFloor()

		// 描画
		for i, modelSet := range w.ModelSets {
			if !w.IsRunning() {
				break
			}

			w.ModelSets[i].prevDeltas = draw(w.Physics, modelSet.model, modelSet.motion, w.Shader,
				modelSet.prevDeltas, i, int(w.frame), elapsed, isDeform, w.EnablePhysics, w.VisibleNormal, w.VisibleBones)
		}

		w.Shader.Msaa.Unbind()

		if !w.IsRunning() {
			break
		}

		w.SwapBuffers()

		if !w.IsRunning() {
			break
		}

		glfw.PollEvents()

		if !w.IsRunning() {
			break
		}

		// 60fpsに制限するための処理
		frameTime = glfw.GetTime()
		elapsed64 := frameTime - prevTime
		if elapsed64 < w.Physics.PhysicsSpf {
			// mlog.I("PhysicsSleep frame=%.8f, elapsed64=%.8f, PhysicsSpf=%.8f", w.frame, elapsed64, w.Physics.PhysicsSpf)
			time.Sleep(time.Duration((elapsed64 - w.Physics.PhysicsSpf) * float64(time.Second)))
		}
	}
	if w.WindowIndex == 0 {
		defer walk.App().Exit(0)
	}
}

func (w *GlWindow) IsRunning() bool {
	return w.drawing && w.running && !CheckOpenGLError() && !w.ShouldClose()
}

// 床描画 ------------------

type MFloor struct {
	vao *mview.VAO
	vbo *mview.VBO
}

func newMFloor() *MFloor {
	mf := &MFloor{}

	mf.vao = mview.NewVAO()
	mf.vao.Bind()
	mf.vbo = mview.NewVBOForFloor()
	mf.vbo.Unbind()
	mf.vao.Unbind()

	return mf
}

func (w *GlWindow) drawFloor() {
	// mlog.D("MFloor.DrawLine")
	w.Shader.Use(mview.PROGRAM_TYPE_FLOOR)

	// 平面を引く
	w.floor.vao.Bind()
	w.floor.vbo.BindFloor()

	gl.DrawArrays(gl.LINES, 0, 240)

	w.floor.vbo.Unbind()
	w.floor.vao.Unbind()

	w.Shader.Unuse()
}
