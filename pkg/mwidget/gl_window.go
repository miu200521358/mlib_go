//go:build windows
// +build windows

package mwidget

import (
	"fmt"
	"image"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/mphysics"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/mlib_go/pkg/mview"
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/mlib_go/pkg/vmd"
)

// 直角の定数値
const RIGHT_ANGLE = 89.9

type GlWindow struct {
	*glfw.Window
	mWindow                *MWindow              // walkウィンドウ
	modelSets              []*ModelSet           // モデルセット
	Shader                 *mview.MShader        // シェーダー
	appConfig              *mconfig.AppConfig    // アプリケーション設定
	title                  string                // ウィンドウタイトル(fpsとか入ってないオリジナル)
	WindowIndex            int                   // ウィンドウインデックス
	prevCursorPos          *mmath.MVec2          // 前回のカーソル位置
	nowCursorPos           *mmath.MVec2          // 現在のカーソル位置
	yaw                    float64               // ウィンドウ操作yaw
	pitch                  float64               // ウィンドウ操作pitch
	Physics                *mphysics.MPhysics    // 物理
	leftButtonPressed      bool                  // 左ボタン押下フラグ
	middleButtonPressed    bool                  // 中ボタン押下フラグ
	rightButtonPressed     bool                  // 右ボタン押下フラグ
	updatedPrev            bool                  // 前回のカーソル位置更新フラグ
	shiftPressed           bool                  // Shiftキー押下フラグ
	ctrlPressed            bool                  // Ctrlキー押下フラグ
	running                bool                  // 描画ループ中フラグ
	playing                bool                  // 再生中フラグ
	doResetPhysicsStart    bool                  // 物理リセット開始フラグ
	doResetPhysicsProgress bool                  // 物理リセット中フラグ
	doResetPhysicsCount    int                   // 物理リセット処理回数
	VisibleBones           map[pmx.BoneFlag]bool // ボーン表示フラグ
	VisibleNormal          bool                  // 法線表示フラグ
	VisibleWire            bool                  // ワイヤーフレーム表示フラグ
	VisibleSelectedVertex  bool                  // 選択頂点表示フラグ
	enablePhysics          bool                  // 物理有効フラグ
	EnableFrameDrop        bool                  // フレームドロップ有効フラグ
	isClosed               bool                  // walkウィンドウが閉じられたかどうか
	isShowInfo             bool                  // 情報表示フラグ
	spfLimit               float64               //fps制限
	frame                  float64               // 現在のフレーム
	prevFrame              int                   // 前回のフレーム
	isSaveDelta            bool                  // 前回デフォーム保存フラグ(walkウィンドウからの変更情報検知用)
	motionPlayer           *MotionPlayer         // 再生ウィジェット
	width                  int                   // ウィンドウ幅
	height                 int                   // ウィンドウ高さ
	floor                  *MFloor               // 床
	funcWorldPos           func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
		nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3,
		vmdDeltas []*vmd.VmdDeltas) // 選択ポイントからのグローバル位置取得コールバック関数
	AppendModelSetChannel      chan *ModelSet         // モデルセット追加チャネル
	RemoveModelSetIndexChannel chan int               // モデルセット削除チャネル
	ReplaceModelSetChannel     chan map[int]*ModelSet // モデルセット入替チャネル
	IsPlayingChannel           chan bool              // 再生チャネル
	FrameChannel               chan int               // フレームチャネル
	IsClosedChannel            chan bool              // ウィンドウクローズチャネル
}

func NewGlWindow(
	width int,
	height int,
	windowIndex int,
	iconImg *image.Image,
	appConfig *mconfig.AppConfig,
	mainWindow *GlWindow,
	fixViewWidget *FixViewWidget,
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

	shader, err := mview.NewMShader(width, height)
	if err != nil {
		return nil, err
	}

	glWindow := GlWindow{
		Window:                     w,
		modelSets:                  make([]*ModelSet, 0),
		Shader:                     shader,
		appConfig:                  appConfig,
		title:                      title,
		WindowIndex:                windowIndex,
		prevCursorPos:              &mmath.MVec2{0, 0},
		nowCursorPos:               &mmath.MVec2{0, 0},
		yaw:                        RIGHT_ANGLE,
		pitch:                      0.0,
		Physics:                    mphysics.NewMPhysics(shader),
		middleButtonPressed:        false,
		rightButtonPressed:         false,
		updatedPrev:                false,
		shiftPressed:               false,
		ctrlPressed:                false,
		VisibleBones:               make(map[pmx.BoneFlag]bool, 0),
		VisibleNormal:              false,
		VisibleWire:                false,
		VisibleSelectedVertex:      false,
		isClosed:                   false,
		isShowInfo:                 false,
		spfLimit:                   1.0 / 30.0,
		running:                    false,
		playing:                    false, // 最初は再生OFF
		enablePhysics:              true,  // 最初は物理ON
		EnableFrameDrop:            true,  // 最初はドロップON
		frame:                      0,
		prevFrame:                  0,
		isSaveDelta:                true,
		width:                      width,
		height:                     height,
		floor:                      newMFloor(),
		AppendModelSetChannel:      make(chan *ModelSet, 1),
		RemoveModelSetIndexChannel: make(chan int, 1),
		ReplaceModelSetChannel:     make(chan map[int]*ModelSet),
		IsPlayingChannel:           make(chan bool, 1),
		FrameChannel:               make(chan int, 1),
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

func (w *GlWindow) SetMWindow(mw *MWindow) {
	w.mWindow = mw
}

func (w *GlWindow) SetMotionPlayer(mp *MotionPlayer) {
	w.motionPlayer = mp
}

func (w *GlWindow) SetFuncWorldPos(f func(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos,
	nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos *mmath.MVec3, vmdDeltas []*vmd.VmdDeltas)) {
	w.funcWorldPos = f
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

func (w *GlWindow) TriggerPlay(p bool) {
	w.playing = p
	w.isSaveDelta = false
}

func (w *GlWindow) GetFrame() int {
	return int(w.frame)
}

func (w *GlWindow) SetFrame(f int) {
	w.frame = float64(f)
	w.prevFrame = f
	w.isSaveDelta = false
}

func (w *GlWindow) TriggerClose(window *glfw.Window) {
	w.running = false
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

	w.TriggerViewReset()

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
		if w.Shader.FieldOfViewAngle < 3.0 {
			w.Shader.FieldOfViewAngle = 3.0
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
	} else if button == glfw.MouseButtonLeft && w.funcWorldPos != nil && w.VisibleSelectedVertex {
		if action == glfw.Press {
			w.leftButtonPressed = true
			w.updatedPrev = false
		} else if action == glfw.Release {
			w.leftButtonPressed = false
			w.execWorldPos()
			w.nowCursorPos = &mmath.MVec2{0, 0}
		}
	}
}

func (gw *GlWindow) getWorldPosition(
	x, y, z float32,
) *mmath.MVec3 {
	// ウィンドウサイズを取得
	w, h := float32(gw.width), float32(gw.height)

	// プロジェクション行列の設定
	projection := mgl32.Perspective(mgl32.DegToRad(gw.Shader.FieldOfViewAngle), w/h, gw.Shader.NearPlane, gw.Shader.FarPlane)
	mlog.V("Projection: %s", projection.String())

	// カメラの位置と中心からビュー行列を計算
	cameraPosition := mgl.NewGlVec3FromMVec3(gw.Shader.CameraPosition)
	lookAtCenter := mgl.NewGlVec3FromMVec3(gw.Shader.LookAtCenterPosition)
	view := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})
	mlog.V("CameraPosition: %s, LookAtCenterPosition: %s", gw.Shader.CameraPosition.String(), gw.Shader.LookAtCenterPosition.String())
	mlog.V("View: %s", view.String())

	worldCoords, err := mgl32.UnProject(
		mgl32.Vec3{x, float32(gw.height) - y, z},
		view, projection, 0, 0, gw.width, gw.height)
	if err != nil {
		mlog.E("UnProject error: %v", err)
		return nil
	}

	worldPos := &mmath.MVec3{float64(-worldCoords.X()), float64(worldCoords.Y()), float64(worldCoords.Z())}
	mlog.D("WindowPos [x=%.3f, y=%.3f, z=%.7f] -> WorldPos [x=%.3f, y=%.3f, z=%.3f]",
		x, y, z, worldPos.GetX(), worldPos.GetY(), worldPos.GetZ())

	return worldPos
}

func (w GlWindow) execWorldPos() {
	prevX := w.prevCursorPos.GetX()
	prevY := w.prevCursorPos.GetY()
	nowX := w.nowCursorPos.GetX()
	nowY := w.nowCursorPos.GetY()

	vmdDeltas := make([]*vmd.VmdDeltas, len(w.modelSets))
	for i, modelSet := range w.modelSets {
		vmdDeltas[i] = modelSet.prevDeltas
	}

	if w.nowCursorPos.Length() == 0 || w.prevCursorPos.Distance(w.nowCursorPos) < 0.1 {
		// カーソルが動いていない場合は直近のだけ取得
		// x: prev, y: prevのワールド座標位置
		depth := w.Shader.Msaa.ReadDepthAt(int(prevX), int(prevY), w.width, w.height)
		worldPos := w.getWorldPosition(float32(prevX), float32(prevY), depth)

		w.funcWorldPos(worldPos, nil, nil, nil, nil, nil, nil, nil, vmdDeltas)
	} else {
		// x: prev, y: prevのワールド座標位置
		prevXprevYDepth := w.Shader.Msaa.ReadDepthAt(int(prevX), int(prevY), w.width, w.height)
		prevXnowYDepth := w.Shader.Msaa.ReadDepthAt(int(prevX), int(nowY), w.width, w.height)
		nowXprevYDepth := w.Shader.Msaa.ReadDepthAt(int(nowX), int(prevY), w.width, w.height)
		nowXnowYDepth := w.Shader.Msaa.ReadDepthAt(int(nowX), int(nowY), w.width, w.height)

		mlog.D("prevXprevYDepth=%.3f, prevXnowYDepth=%.3f, nowXprevYDepth=%.3f, nowXnowYDepth=%.3f",
			prevXprevYDepth, prevXnowYDepth, nowXprevYDepth, nowXnowYDepth)

		// 最も手前を基準とする
		depth := min(prevXprevYDepth, prevXnowYDepth, nowXprevYDepth, nowXnowYDepth)

		prevXprevYFrontPos := w.getWorldPosition(float32(prevX), float32(prevY), max(depth-1e-5, 0.0))
		prevXprevYBackPos := w.getWorldPosition(float32(prevX), float32(prevY), min(depth+1e-5, 1.0))

		// x: prev, y: nowのワールド座標位置
		prevXnowYFrontPos := w.getWorldPosition(float32(prevX), float32(nowY), max(depth-1e-5, 0.0))
		prevXnowYBackPos := w.getWorldPosition(float32(prevX), float32(nowY), min(depth+1e-5, 1.0))

		// x: now, y: prevのワールド座標位置
		nowXprevYFrontPos := w.getWorldPosition(float32(nowX), float32(prevY), max(depth-1e-5, 0.0))
		nowXprevYBackPos := w.getWorldPosition(float32(nowX), float32(prevY), min(depth+1e-5, 1.0))

		// x: now, y: nowのワールド座標位置
		nowXnowYFrontPos := w.getWorldPosition(float32(nowX), float32(nowY), max(depth-1e-5, 0.0))
		nowXnowYBackPos := w.getWorldPosition(float32(nowX), float32(nowY), min(depth+1e-5, 1.0))

		w.funcWorldPos(prevXprevYFrontPos, prevXprevYBackPos, prevXnowYFrontPos, prevXnowYBackPos, nowXprevYFrontPos, nowXprevYBackPos, nowXnowYFrontPos, nowXnowYBackPos, vmdDeltas)
	}
}

func (w *GlWindow) updateCameraAngle(xpos, ypos float64) {

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
	// mlog.D("xOffset %.7f, yOffset %.7f, CameraPosition: %s, LookAtCenterPosition: %s\n",
	// 	xOffset, yOffset, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())
}

func (w *GlWindow) updateCameraPosition(xpos, ypos float64) {

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

func (w *GlWindow) handleCursorPosEvent(window *glfw.Window, xpos float64, ypos float64) {
	// mlog.D("[start] yaw %.7f, pitch %.7f, CameraPosition: %s, LookAtCenterPosition: %s\n",
	// 	w.yaw, w.pitch, w.Shader.CameraPosition.String(), w.Shader.LookAtCenterPosition.String())

	if !w.updatedPrev {
		w.prevCursorPos.SetX(xpos)
		w.prevCursorPos.SetY(ypos)
		w.updatedPrev = true
		return
	}

	if w.leftButtonPressed {
		w.nowCursorPos.SetX(xpos)
		w.nowCursorPos.SetY(ypos)
		return
	} else if w.rightButtonPressed {
		// 右クリックはカメラの角度を更新
		w.updateCameraAngle(xpos, ypos)
	} else if w.middleButtonPressed {
		// 中クリックはカメラ位置と中心を移動
		w.updateCameraPosition(xpos, ypos)
	}

	w.prevCursorPos.SetX(xpos)
	w.prevCursorPos.SetY(ypos)
}

func (w *GlWindow) TriggerPhysicsEnabled(enabled bool) {
	w.enablePhysics = enabled
	for i := range w.modelSets {
		w.modelSets[i].prevDeltas = nil
	}
}

func (w *GlWindow) TriggerPhysicsReset() {
	if !w.doResetPhysicsProgress {
		w.doResetPhysicsStart = true
	}
}

func (w *GlWindow) resetPhysicsStart() {
	// 物理ON・まだリセット中ではないの時だけリセット処理を行う
	if w.enablePhysics && !w.doResetPhysicsProgress {
		// 一旦物理OFFにする
		w.TriggerPhysicsEnabled(false)
		w.Physics.ResetWorld()
		w.doResetPhysicsStart = false
		w.doResetPhysicsProgress = true
		w.doResetPhysicsCount = 0
	}
}

func (w *GlWindow) resetPhysicsFinish() {
	// 物理ONに戻してリセットフラグを落とす
	w.TriggerPhysicsEnabled(true)
	w.doResetPhysicsStart = false
	w.doResetPhysicsProgress = false
	w.doResetPhysicsCount = 0
}

func (w *GlWindow) TriggerViewReset() {
	// カメラとかリセット
	w.Shader.Reset()
	w.prevCursorPos = &mmath.MVec2{0, 0}
	w.yaw = RIGHT_ANGLE
	w.pitch = 0.0
	w.middleButtonPressed = false
	w.rightButtonPressed = false

}

func (w *GlWindow) Size() walk.Size {
	x, y := w.Window.GetSize()
	return walk.Size{Width: x, Height: y}
}

func (w *GlWindow) Run() {
	prevTime := glfw.GetTime()
	prevShowTime := glfw.GetTime()
	w.prevFrame = 0
	w.running = true

	go func() {
	channelLoop:
		for {
			select {
			case pair := <-w.AppendModelSetChannel:
				// 追加処理
				w.modelSets[len(w.modelSets)-1] = pair
				w.isSaveDelta = false
			case pairMap := <-w.ReplaceModelSetChannel:
				// 入替処理
				for i := range pairMap {
					// 変更が加えられている可能性があるので、セットアップ実施（変更がなければスルーされる）
					if pairMap[i].NextModel != nil {
						pairMap[i].NextModel.Setup()
					}

					for j := len(w.modelSets); j <= i; j++ {
						// 既存のモデルセットがない場合は追加
						w.modelSets = append(w.modelSets, NewModelSet())
					}

					w.modelSets[i].NextModel = pairMap[i].NextModel
					w.modelSets[i].NextMotion = pairMap[i].NextMotion
					w.modelSets[i].NextSelectedVertexIndexes = pairMap[i].NextSelectedVertexIndexes
					w.modelSets[i].NextInvisibleMaterialIndexes = pairMap[i].NextInvisibleMaterialIndexes
				}
				w.isSaveDelta = false
			case index := <-w.RemoveModelSetIndexChannel:
				// 削除処理
				if index < len(w.modelSets) {
					if w.modelSets[index].Model != nil {
						w.modelSets[index].Model.Delete()
					}
					w.modelSets[index] = NewModelSet()
				}
				w.isSaveDelta = false
			case isPlaying := <-w.IsPlayingChannel:
				// 再生設定
				w.TriggerPlay(isPlaying)
			case frame := <-w.FrameChannel:
				// フレーム設定
				w.SetFrame(frame)
			case isClosed := <-w.IsClosedChannel:
				// ウィンドウが閉じられた場合
				w.isClosed = isClosed
				break channelLoop
			}
		}
	}()

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

		if w.playing && w.motionPlayer != nil && w.frame >= w.motionPlayer.FrameEdit.MaxValue() {
			// 再生中に最後までいったら最初にループして戻る
			w.TriggerPhysicsReset()
			w.SetFrame(0)

			go func() {
				w.motionPlayer.SetValue(int(w.frame))
			}()
		}

		frameTime := glfw.GetTime()
		originalElapsed := frameTime - prevTime

		var elapsed float64
		var timeStep float32
		if !w.EnableFrameDrop {
			// フレームドロップOFF
			// 物理fpsは60fps固定
			timeStep = w.Physics.PhysicsSpf
			// デフォームfpsは30fps上限の経過時間
			elapsed = mmath.ClampedFloat(originalElapsed, 0.0, float64(w.Physics.DeformSpf))
		} else {
			// 物理fpsは経過時間
			timeStep = float32(originalElapsed)
			elapsed = originalElapsed
		}

		if elapsed < w.spfLimit {
			// 1フレームの時間が経過していない場合はスキップ
			// fps制限は描画fpsにのみ依存
			continue
		}

		w.MakeContextCurrent()

		// MSAAフレームバッファをバインド
		w.Shader.Msaa.Bind()

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
		cameraPosition := mgl.NewGlVec3FromMVec3(w.Shader.CameraPosition)

		// カメラの中心
		lookAtCenter := mgl.NewGlVec3FromMVec3(w.Shader.LookAtCenterPosition)
		camera := mgl32.LookAtV(cameraPosition, lookAtCenter, mgl32.Vec3{0, 1, 0})

		for _, program := range w.Shader.GetPrograms() {
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

		// 床平面を描画
		w.drawFloor()

		if w.playing {
			// 経過秒数をキーフレームの進捗具合に合わせて調整
			if w.spfLimit < -1 {
				// デフォームFPS制限なしの場合、フレーム番号を常に進める
				w.frame += 1.0
			} else {
				w.frame += elapsed * float64(w.Physics.DeformFps)
			}
		}

		if w.doResetPhysicsStart {
			// 物理リセット開始
			w.resetPhysicsStart()
		}

		// デフォーム
		w.modelSets = deformsAll(
			w.Physics, w.modelSets, int(w.frame), w.prevFrame,
			timeStep, w.enablePhysics, w.doResetPhysicsProgress)

		// 描画
		for k := range w.modelSets {
			if w.modelSets[k].Model != nil && w.modelSets[k].prevDeltas != nil {
				w.modelSets[k].prevDeltas = draw(
					w.Physics, w.modelSets[k].Model, w.Shader, w.modelSets[k].prevDeltas,
					w.modelSets[k].InvisibleMaterialIndexes, w.modelSets[k].NextInvisibleMaterialIndexes,
					w.WindowIndex, w.VisibleNormal, w.VisibleWire, w.VisibleSelectedVertex, w.VisibleBones)
			}

			if w.modelSets[k].NextMotion != nil {
				w.TriggerPhysicsReset()

				w.modelSets[k].Motion = w.modelSets[k].NextMotion
				w.modelSets[k].NextMotion = nil
				w.isSaveDelta = false
			}

			// モデルが変わっている場合は最新の情報を取得する
			if w.modelSets[k].NextModel != nil {
				w.TriggerPhysicsReset()

				// 次のモデルが指定されている場合、初期化して入替
				if w.modelSets[k].Model != nil && w.modelSets[k].Model.DrawInitialized {
					// 既存モデルが描画初期化されてたら削除
					w.modelSets[k].Model.Delete()
					w.modelSets[k].Model = nil
				}
				if !w.modelSets[k].NextModel.DrawInitialized {
					w.modelSets[k].NextModel.Index = k
					w.modelSets[k].NextModel.DrawInitialize(w.WindowIndex, w.Physics)
				}
				w.modelSets[k].Model = w.modelSets[k].NextModel
				w.modelSets[k].NextModel = nil
				w.isSaveDelta = false
			}

			if w.modelSets[k].NextInvisibleMaterialIndexes != nil {
				w.modelSets[k].InvisibleMaterialIndexes = w.modelSets[k].NextInvisibleMaterialIndexes
				w.modelSets[k].NextInvisibleMaterialIndexes = nil
			}

			if w.modelSets[k].NextSelectedVertexIndexes != nil {
				w.modelSets[k].SelectedVertexIndexes = w.modelSets[k].NextSelectedVertexIndexes
				w.modelSets[k].NextSelectedVertexIndexes = nil
			}

			// キーフレの手動変更がなかった場合のみ前回デフォームとして保持
			if !w.isSaveDelta {
				w.modelSets[k].prevDeltas = nil
			}
			w.isSaveDelta = true
		}

		if w.doResetPhysicsProgress {
			if w.doResetPhysicsCount > 1 {
				// 0: 物理リセット開始
				// 1: 物理リセット中(リセット状態で物理更新)
				// 2: 物理リセット完了
				// 物理リセット完了
				w.resetPhysicsFinish()
			} else {
				w.doResetPhysicsCount++
			}
		}

		prevTime = frameTime

		if w.playing && int(w.frame) > w.prevFrame {
			// フレーム番号上書き
			w.prevFrame = int(w.frame)
			if w.playing && w.motionPlayer != nil {
				go func() {
					w.motionPlayer.SetValue(int(w.frame))
				}()
			}
		}

		w.Shader.Msaa.Resolve()
		w.Shader.Msaa.Unbind()

		// 物理デバッグ表示（要不要は中で見ている）MSAAには含めない
		w.Physics.DrawDebugLines()

		w.SwapBuffers()

		if w.isShowInfo {
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

		// if w.frame > 100 {
		// 	goto closeApp
		// }
	}

closeApp:
	w.Shader.Delete()
	for i := range w.modelSets {
		w.modelSets[i].Model.Delete()
	}
	if w.WindowIndex == 0 {
		glfw.Terminate()
		walk.App().Exit(0)
	}
}

func (w *GlWindow) IsRunning() bool {
	return !w.isClosed && // walkウィンドウ側が閉じられたか
		w.running && // GLウィンドウ側が閉じられたか
		!CheckOpenGLError() && !w.ShouldClose() &&
		((w.mWindow != nil && !w.mWindow.IsDisposed()) || w.mWindow == nil)
}

// 床描画 ------------------

type MFloor struct {
	vao   *mview.VAO
	vbo   *mview.VBO
	count int32
}

func newMFloor() *MFloor {
	mf := &MFloor{}

	mf.vao = mview.NewVAO()
	mf.vao.Bind()
	mf.vbo, mf.count = mview.NewVBOForFloor()
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

	gl.DrawArrays(gl.LINES, 0, w.floor.count)

	w.floor.vbo.Unbind()
	w.floor.vao.Unbind()

	w.Shader.Unuse()
}
