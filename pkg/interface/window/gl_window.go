//go:build windows
// +build windows

package window

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
	"github.com/miu200521358/mlib_go/pkg/interface/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

// 直角の定数値
const RIGHT_ANGLE = 89.9

type GlWindow struct {
	*glfw.Window
	UiState         *widget.UiState             // UI状態
	animationStates []*renderer.AnimationStates // アニメーション状態
	shader          *mgl.MShader                // シェーダー
	physics         *mbt.MPhysics               // 物理
	appConfig       *mconfig.AppConfig          // アプリケーション設定
	title           string                      // ウィンドウタイトル(fpsとか入ってないオリジナル)
	WindowIndex     int                         // ウィンドウインデックス
	prevCursorPos   *mmath.MVec2                // 前回のカーソル位置
	nowCursorPos    *mmath.MVec2                // 現在のカーソル位置
	yaw             float64                     // ウィンドウ操作yaw
	pitch           float64                     // ウィンドウ操作pitch
	prevFrame       int                         // 前回のフレーム
	width           int                         // ウィンドウ幅
	height          int                         // ウィンドウ高さ
	floor           *MFloor                     // 床
	worldPosFunc    func(prevXprevYFrontPos,
		prevXprevYBackPos,
		prevXnowYFrontPos,
		prevXnowYBackPos,
		nowXprevYFrontPos,
		nowXprevYBackPos,
		nowXnowYFrontPos,
		nowXnowYBackPos *mmath.MVec3,
		vmdDeltas []*delta.VmdDeltas) // 選択ポイントからのグローバル位置取得コールバック関数
	RemoveIndexChannel   chan int                    // モデル削除チャネル
	AppendModelChannel   chan *pmx.PmxModel          // モデル追加チャネル
	ReplaceModelChannel  chan map[int]*pmx.PmxModel  // モデル入替チャネル
	AppendMotionChannel  chan *vmd.VmdMotion         // モーション追加チャネル
	ReplaceMotionChannel chan map[int]*vmd.VmdMotion // モーション入替チャネル
	IsPlayingChannel     chan bool                   // 再生チャネル
	FrameChannel         chan int                    // フレームチャネル
	IsClosedChannel      chan bool                   // ウィンドウクローズチャネル
}

func NewGlWindow(
	width int,
	height int,
	windowIndex int,
	iconImg *image.Image,
	appConfig *mconfig.AppConfig,
	mainWindow *GlWindow,
	uiState *widget.UiState,
) (*GlWindow, error) {
	if windowIndex == 0 {
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
	title := mi18n.T("ビューワー")
	w, err := glfw.CreateWindow(width, height, title, nil, mw)
	if err != nil {
		return nil, err
	}
	w.MakeContextCurrent()
	w.SetInputMode(glfw.StickyKeysMode, glfw.True)
	w.SetIcon([]image.Image{*iconImg})

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
			mlog.V("[MEDIUM] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
		case gl.DEBUG_SEVERITY_LOW:
			mlog.V("[LOW] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
				source, glType, severity, message)
			// case gl.DEBUG_SEVERITY_NOTIFICATION:
			// 	mlog.D("[NOTIFICATION] GL CALLBACK: %v type = 0x%x, severity = 0x%x, message = %s\n",
			// 		source, glType, severity, message)
		}
	}, gl.Ptr(nil))
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS) // 同期的なデバッグ出力を有効にします。

	shader, err := mgl.NewMShader(width, height)
	if err != nil {
		return nil, err
	}

	glWindow := GlWindow{
		Window:               w,
		UiState:              uiState,
		animationStates:      make([]*renderer.AnimationStates, 0),
		shader:               shader,
		appConfig:            appConfig,
		title:                title,
		WindowIndex:          windowIndex,
		prevCursorPos:        mmath.NewMVec2(),
		nowCursorPos:         mmath.NewMVec2(),
		yaw:                  RIGHT_ANGLE,
		pitch:                0.0,
		physics:              mbt.NewMPhysics(),
		prevFrame:            0,
		width:                width,
		height:               height,
		floor:                newMFloor(),
		RemoveIndexChannel:   make(chan int, 1),
		AppendModelChannel:   make(chan *pmx.PmxModel, 1),
		ReplaceModelChannel:  make(chan map[int]*pmx.PmxModel, 1),
		AppendMotionChannel:  make(chan *vmd.VmdMotion, 1),
		ReplaceMotionChannel: make(chan map[int]*vmd.VmdMotion, 1),
		IsPlayingChannel:     make(chan bool, 1),
		FrameChannel:         make(chan int, 1),
	}

	w.SetScrollCallback(glWindow.handleScrollEvent)
	w.SetMouseButtonCallback(glWindow.handleMouseButtonEvent)
	w.SetCursorPosCallback(glWindow.handleCursorPosEvent)
	w.SetKeyCallback(glWindow.handleKeyEvent)
	w.SetCloseCallback(glWindow.TriggerClose)
	w.SetSizeCallback(glWindow.resize)
	w.SetFramebufferSizeCallback(glWindow.resizeBuffer)

	return &glWindow, nil
}

func (w *GlWindow) SetTitle(title string) {
	w.title = title
	w.Window.SetTitle(title)
}

func (w *GlWindow) SetWorldPosFunc(f func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
	nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3, vmdDeltas []*delta.VmdDeltas)) {
	w.worldPosFunc = f
}

func (w *GlWindow) TriggerClose(window *glfw.Window) {
	w.UiState.IsGlRunning = false
}

func (gw *GlWindow) getWorldPosition(
	x, y, z float32,
) *mmath.MVec3 {
	// ウィンドウサイズを取得
	w, h := float32(gw.width), float32(gw.height)

	// プロジェクション行列の設定
	projection := mgl32.Perspective(mgl32.DegToRad(gw.shader.FieldOfViewAngle), w/h, gw.shader.NearPlane, gw.shader.FarPlane)
	mlog.V("Projection: %s", projection.String())

	// カメラの位置と中心からビュー行列を計算
	cameraPosition := mgl.NewGlVec3(gw.shader.CameraPosition)
	lookAtCenter := mgl.NewGlVec3(gw.shader.LookAtCenterPosition)
	view := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})
	mlog.V("CameraPosition: %s, LookAtCenterPosition: %s", gw.shader.CameraPosition.String(), gw.shader.LookAtCenterPosition.String())
	mlog.V("View: %s", view.String())

	worldCoords, err := mgl32.UnProject(
		mgl32.Vec3{x, float32(gw.height) - y, z},
		view, projection, 0, 0, gw.width, gw.height)
	if err != nil {
		mlog.E("UnProject error: %v", err)
		return nil
	}

	worldPos := &mmath.MVec3{X: float64(-worldCoords.X()), Y: float64(worldCoords.Y()), Z: float64(worldCoords.Z())}
	mlog.D("WindowPos [x=%.3f, y=%.3f, z=%.7f] -> WorldPos [x=%.3f, y=%.3f, z=%.3f]",
		x, y, z, worldPos.X, worldPos.Y, worldPos.Z)

	return worldPos
}

func (w GlWindow) execWorldPos() {
	// prevX := w.prevCursorPos.X
	// prevY := w.prevCursorPos.Y
	// nowX := w.nowCursorPos.X
	// nowY := w.nowCursorPos.Y

	// vmdDeltas := make([]*delta.VmdDeltas, len(w.animateStates))
	// for i, modelSet := range w.animateStates {
	// 	vmdDeltas[i] = modelSet.PrevDeltas
	// }

	// if w.nowCursorPos.Length() == 0 || w.prevCursorPos.Distance(w.nowCursorPos) < 0.1 {
	// 	// カーソルが動いていない場合は直近のだけ取得
	// 	// x: prev, y: prevのワールド座標位置
	// 	depth := w.shader.Msaa.ReadDepthAt(int(prevX), int(prevY), w.width, w.height)
	// 	worldPos := w.getWorldPosition(float32(prevX), float32(prevY), depth)

	// 	w.worldPosFunc(worldPos, nil, nil, nil, nil, nil, nil, nil, vmdDeltas)
	// } else {
	// 	// x: prev, y: prevのワールド座標位置
	// 	prevXprevYDepth := w.shader.Msaa.ReadDepthAt(int(prevX), int(prevY), w.width, w.height)
	// 	prevXnowYDepth := w.shader.Msaa.ReadDepthAt(int(prevX), int(nowY), w.width, w.height)
	// 	nowXprevYDepth := w.shader.Msaa.ReadDepthAt(int(nowX), int(prevY), w.width, w.height)
	// 	nowXnowYDepth := w.shader.Msaa.ReadDepthAt(int(nowX), int(nowY), w.width, w.height)

	// 	mlog.D("prevXprevYDepth=%.3f, prevXnowYDepth=%.3f, nowXprevYDepth=%.3f, nowXnowYDepth=%.3f",
	// 		prevXprevYDepth, prevXnowYDepth, nowXprevYDepth, nowXnowYDepth)

	// 	// 最も手前を基準とする
	// 	depth := min(prevXprevYDepth, prevXnowYDepth, nowXprevYDepth, nowXnowYDepth)

	// 	prevXprevYFrontPos := w.getWorldPosition(float32(prevX), float32(prevY), max(depth-1e-5, 0.0))
	// 	prevXprevYBackPos := w.getWorldPosition(float32(prevX), float32(prevY), min(depth+1e-5, 1.0))

	// 	// x: prev, y: nowのワールド座標位置
	// 	prevXnowYFrontPos := w.getWorldPosition(float32(prevX), float32(nowY), max(depth-1e-5, 0.0))
	// 	prevXnowYBackPos := w.getWorldPosition(float32(prevX), float32(nowY), min(depth+1e-5, 1.0))

	// 	// x: now, y: prevのワールド座標位置
	// 	nowXprevYFrontPos := w.getWorldPosition(float32(nowX), float32(prevY), max(depth-1e-5, 0.0))
	// 	nowXprevYBackPos := w.getWorldPosition(float32(nowX), float32(prevY), min(depth+1e-5, 1.0))

	// 	// x: now, y: nowのワールド座標位置
	// 	nowXnowYFrontPos := w.getWorldPosition(float32(nowX), float32(nowY), max(depth-1e-5, 0.0))
	// 	nowXnowYBackPos := w.getWorldPosition(float32(nowX), float32(nowY), min(depth+1e-5, 1.0))

	// 	w.worldPosFunc(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos, nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos, vmdDeltas)
	// }
}

func (w *GlWindow) resetPhysicsStart() {
	// 物理ON・まだリセット中ではないの時だけリセット処理を行う
	if w.UiState.EnabledPhysics && !w.UiState.DoResetPhysicsProgress {
		// 一旦物理OFFにする
		w.UiState.TriggerPhysicsEnabled(false)
		w.physics.ResetWorld()
		w.UiState.DoResetPhysicsStart = false
		w.UiState.DoResetPhysicsProgress = true
		w.UiState.DoResetPhysicsCount = 0
	}
}

func (w *GlWindow) resetPhysicsFinish() {
	// 物理ONに戻してリセットフラグを落とす
	w.UiState.TriggerPhysicsEnabled(true)
	w.UiState.DoResetPhysicsStart = false
	w.UiState.DoResetPhysicsProgress = false
	w.UiState.DoResetPhysicsCount = 0
}

func (w *GlWindow) TriggerViewReset() {
	// カメラとかリセット
	w.shader.Reset()
	w.prevCursorPos = mmath.NewMVec2()
	w.nowCursorPos = mmath.NewMVec2()
	w.yaw = RIGHT_ANGLE
	w.pitch = 0.0
	w.UiState.MiddleButtonPressed = false
	w.UiState.RightButtonPressed = false

}

func (w *GlWindow) Size() walk.Size {
	x, y := w.Window.GetSize()
	return walk.Size{Width: x, Height: y}
}

func (w *GlWindow) RunChannel() {

	go func() {
	channelLoop:
		for {
			select {
			case model := <-w.AppendModelChannel:
				// モデル追加処理
				w.animationStates[len(w.animationStates)-1].Next.Model = model
				w.UiState.IsSaveDelta = false
			case models := <-w.ReplaceModelChannel:
				// モデル入替処理
				for i := range models {
					// 変更が加えられている可能性があるので、セットアップ実施（変更がなければスルーされる）
					for j := len(w.animationStates); j <= i; j++ {
						// 既存のモデルセットがない場合は追加
						w.animationStates = append(w.animationStates, renderer.NewAnimationStates())
					}

					w.animationStates[i].Next.Model = models[i]
					w.animationStates[i].Next.Model.Setup()
				}
				w.UiState.IsSaveDelta = false
			case motion := <-w.AppendMotionChannel:
				// モーション追加処理
				w.animationStates[len(w.animationStates)-1].Next.Motion = motion
				w.UiState.IsSaveDelta = false
			case motions := <-w.ReplaceMotionChannel:
				// モーション入替処理
				for i := range motions {
					// 変更が加えられている可能性があるので、セットアップ実施（変更がなければスルーされる）
					for j := len(w.animationStates); j <= i; j++ {
						// 既存のモデルセットがない場合は追加
						w.animationStates = append(w.animationStates, renderer.NewAnimationStates())
					}

					w.animationStates[i].Next.Motion = motions[i]
				}
				w.UiState.IsSaveDelta = false
			case index := <-w.RemoveIndexChannel:
				// 削除処理
				if index < len(w.animationStates) {
					if w.animationStates[index].Now.Model != nil {
						w.physics.DeleteModel(index)
						if w.animationStates[index].Now.RenderModel != nil {
							w.animationStates[index].Now.RenderModel.Delete()
						}
					}
					w.animationStates[index] = renderer.NewAnimationStates()
				}
				w.UiState.IsSaveDelta = false
			case isPlaying := <-w.IsPlayingChannel:
				// 再生設定
				w.UiState.TriggerPlay(isPlaying)
			case frame := <-w.FrameChannel:
				// フレーム設定
				w.UiState.SetFrame(float64(frame))
			case isClosed := <-w.IsClosedChannel:
				// ウィンドウが閉じられた場合
				w.UiState.IsWalkRunning = !isClosed
				break channelLoop
			}
		}
	}()

	w.TriggerClose(w.Window)
}

func (w *GlWindow) Run() {
	prevTime := glfw.GetTime()
	prevShowTime := glfw.GetTime()
	w.prevFrame = 0
	w.UiState.IsGlRunning = true

	for {
		glfw.PollEvents()

		if !w.IsRunning() {
			goto closeApp
		}

		if w.width == 0 || w.height == 0 {
			// ウィンドウが最小化されている場合は描画をスキップ(フレームも進めない)
			prevTime = glfw.GetTime()
			continue
		}

		if w.UiState.Playing && w.UiState.Frame() >= float64(w.UiState.MaxFrame) {
			// 再生中に最後までいったら最初にループして戻る
			w.UiState.TriggerPhysicsReset()
			w.UiState.SetFrame(0)
		}

		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		var elapsed float64
		var timeStep float32
		if !w.UiState.EnabledFrameDrop {
			// フレームドロップOFF
			// 物理fpsは60fps固定
			timeStep = w.physics.PhysicsSpf
			// デフォームfpsは30fps上限の経過時間
			elapsed = mmath.ClampedFloat(originalElapsed, 0.0, float64(w.physics.DeformSpf))
		} else {
			// 物理fpsは経過時間
			timeStep = float32(originalElapsed)
			elapsed = originalElapsed
		}

		if elapsed < w.UiState.SpfLimit {
			// 1フレームの時間が経過していない場合はスキップ
			// fps制限は描画fpsにのみ依存
			continue
		}

		w.MakeContextCurrent()

		// MSAAフレームバッファをバインド
		w.shader.Msaa.Bind()

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
			mgl32.DegToRad(w.shader.FieldOfViewAngle),
			float32(w.shader.Width)/float32(w.shader.Height),
			w.shader.NearPlane,
			w.shader.FarPlane,
		)

		// カメラの位置
		cameraPosition := mgl.NewGlVec3(w.shader.CameraPosition)

		// カメラの中心
		lookAtCenter := mgl.NewGlVec3(w.shader.LookAtCenterPosition)
		camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})

		for _, program := range w.shader.GetPrograms() {
			// プログラムの切り替え
			gl.UseProgram(program)

			// カメラの再計算
			projectionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_PROJECTION_MATRIX))
			gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

			// カメラの位置
			cameraPositionUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CAMERA_POSITION))
			gl.Uniform3fv(cameraPositionUniform, 1, &cameraPosition[0])

			// カメラの中心
			cameraUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_MODEL_VIEW_MATRIX))
			gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

			gl.UseProgram(0)
		}

		// 床平面を描画
		w.drawFloor()

		if w.UiState.Playing {
			// 経過秒数をキーフレームの進捗具合に合わせて調整
			if w.UiState.SpfLimit < -1 {
				// デフォームFPS制限なしの場合、フレーム番号を常に進める
				w.UiState.AddFrame(1.0)
			} else {
				w.UiState.AddFrame(elapsed * float64(w.physics.DeformFps))
			}
		}

		if w.UiState.DoResetPhysicsStart {
			// 物理リセット開始
			w.resetPhysicsStart()
		}

		// デフォーム
		w.animationStates = renderer.Animate(
			w.physics, w.animationStates, int(w.UiState.Frame()), timeStep, w.UiState.EnabledPhysics, w.UiState.DoResetPhysicsProgress)

		// 描画
		for k := range w.animationStates {
			if w.animationStates[k].Now.RenderModel != nil && w.animationStates[k].Now.VmdDeltas != nil {
				w.animationStates[k].Now.RenderModel.Render(
					w.shader, w.animationStates[k], w.WindowIndex, w.UiState.IsShowNormal, w.UiState.IsShowWire,
					w.UiState.IsShowSelectedVertex, w.UiState.IsShowBones)
			}

			if w.animationStates[k].Next.Motion != nil {
				w.animationStates[k].Now.Motion = w.animationStates[k].Next.Motion
				w.animationStates[k].Next.Motion = nil
				w.UiState.IsSaveDelta = false
			}

			// モデルが変わっている場合は最新の情報を取得する
			if w.animationStates[k].Next.Model != nil {
				// 次のモデルが指定されている場合、初期化して入替
				if w.animationStates[k].Now.RenderModel != nil {
					// 既存モデルが描画初期化されてたら削除
					w.physics.DeleteModel(k)
					w.animationStates[k].Now.RenderModel.Delete()
					w.animationStates[k].Now.Model = nil
				}
				w.animationStates[k].Now.Model = w.animationStates[k].Next.Model
				w.animationStates[k].Next.Model = nil
				w.animationStates[k].Now.RenderModel =
					renderer.NewRenderModel(w.WindowIndex, w.animationStates[k].Now.Model)
				w.physics.AddModel(k, w.animationStates[k].Now.Model)

				w.UiState.TriggerPhysicsReset()
			}

			if w.animationStates[k].Next.InvisibleMaterialIndexes != nil {
				w.animationStates[k].Now.InvisibleMaterialIndexes = w.animationStates[k].Next.InvisibleMaterialIndexes
				w.animationStates[k].Next.InvisibleMaterialIndexes = nil
			}

			if w.animationStates[k].Next.SelectedVertexIndexes != nil {
				w.animationStates[k].Now.SelectedVertexIndexes = w.animationStates[k].Next.SelectedVertexIndexes
				w.animationStates[k].Next.SelectedVertexIndexes = nil
			}

			// キーフレの手動変更がなかった場合のみ前回デフォームとして保持
			if !w.UiState.IsSaveDelta {
				w.animationStates[k].Now.VmdDeltas = nil
			}
			w.UiState.IsSaveDelta = true
		}

		if w.UiState.DoResetPhysicsProgress {
			if w.UiState.DoResetPhysicsCount > 1 {
				// 0: 物理リセット開始
				// 1: 物理リセット中(リセット状態で物理更新)
				// 2: 物理リセット完了
				// 物理リセット完了
				w.resetPhysicsFinish()
			} else {
				w.UiState.DoResetPhysicsCount++
			}
		}

		prevTime = frameTime

		if w.UiState.Playing && int(w.UiState.Frame()) > w.prevFrame {
			// フレーム番号上書き
			w.prevFrame = int(w.UiState.Frame())
			if w.UiState.Playing {
				w.UiState.SetFrame(w.UiState.Frame())
			}
		}

		w.shader.Msaa.Resolve()
		w.shader.Msaa.Unbind()

		// 物理デバッグ表示
		w.physics.DrawDebugLines(
			w.shader,
			w.UiState.IsShowRigidBodyFront || w.UiState.IsShowRigidBodyBack,
			w.UiState.IsShowJoint, w.UiState.IsShowRigidBodyFront)

		w.SwapBuffers()

		if w.UiState.IsShowInfo {
			nowShowTime := glfw.GetTime()
			// 1秒ごとにオリジナルの経過時間からFPSを表示
			if nowShowTime-prevShowTime >= 1.0 {
				var suffixFps string
				if w.appConfig.IsEnvProd() {
					// リリース版の場合、FPSの表示を簡略化
					suffixFps = fmt.Sprintf("%.2f fps", 1.0/elapsed)
				} else {
					// 開発版の場合、FPSの表示を詳細化
					suffixFps = fmt.Sprintf("d) %.2f (%.2f) / p) %.2f fps", 1.0/elapsed, 1/originalElapsed, 1.0/timeStep)
				}

				w.Window.SetTitle(fmt.Sprintf("%s - %s", w.title, suffixFps))
				prevShowTime = nowShowTime
			}
		} else {
			w.Window.SetTitle(w.title)
		}

		// if w.UiState.Frame() > 100 {
		// 	goto closeApp
		// }
	}

closeApp:
	for i := range w.animationStates {
		if w.animationStates[i].Now.RenderModel != nil {
			w.physics.DeleteModel(i)
			w.animationStates[i].Now.RenderModel.Delete()
		}
	}
	w.shader.Delete()
	if w.WindowIndex == 0 {
		glfw.Terminate()
		walk.App().Exit(0)
	}
}

func (w *GlWindow) IsRunning() bool {
	return w.UiState.IsWalkRunning && // walkウィンドウ側が閉じられたか
		w.UiState.IsGlRunning && // GLウィンドウ側が閉じられたか
		!mgl.CheckOpenGLError() && !w.ShouldClose()
}

// 床描画 ------------------

type MFloor struct {
	vao   *buffer.VAO
	vbo   *buffer.VBO
	count int32
}

func newMFloor() *MFloor {
	mf := &MFloor{}

	mf.vao = buffer.NewVAO()
	mf.vao.Bind()
	mf.vbo, mf.count = buffer.NewVBOForFloor()
	mf.vbo.Unbind()
	mf.vao.Unbind()

	return mf
}

func (w *GlWindow) drawFloor() {
	// mlog.D("MFloor.DrawLine")
	program := w.shader.GetProgram(mgl.PROGRAM_TYPE_FLOOR)
	gl.UseProgram(program)

	// 平面を引く
	w.floor.vao.Bind()
	w.floor.vbo.BindFloor()

	gl.DrawArrays(gl.LINES, 0, w.floor.count)

	w.floor.vbo.Unbind()
	w.floor.vao.Unbind()

	gl.UseProgram(0)
}
