//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mbt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/render"
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
	cursorX             float64                 // カーソルX位置
	cursorY             float64                 // カーソルY位置
	tooltipRenderer     *mgl.TooltipRenderer    // ツールチップ描画オブジェクト
	list                *ViewerList             // ビューワーリスト
	shader              rendering.IShader       // シェーダー
	physics             physics.IPhysics        // 物理エンジン
	modelRenderers      []*render.ModelRenderer // モデル描画オブジェクト
	motions             []*vmd.VmdMotion        // モーションデータ
	vmdDeltas           []*delta.VmdDeltas      // 変形情報
	overrideOffset      *mmath.MVec3            // オーバーライド補正オフセット
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
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 4)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	glWindow, err := glfw.CreateWindow(width, height, title, nil, mainWindow)
	if err != nil {
		return nil, err
	}

	glWindow.MakeContextCurrent()
	glWindow.SetInputMode(glfw.StickyKeysMode, glfw.True)
	glWindow.SetIcon([]image.Image{icon})
	glfw.SwapInterval(0) // VSync無効

	// OpenGL の初期化
	if err := gl.Init(); err != nil {
		return nil, err
	}

	// シェーダー初期化
	shaderFactory := mgl.NewMShaderFactory()
	shader, err := shaderFactory.CreateShader(windowIndex, width, height)
	if err != nil {
		return nil, err
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	gravity := list.shared.LoadPhysicsWorldMotion(windowIndex).GravityFrames.Get(0).Gravity

	vw := &ViewWindow{
		Window:         glWindow,
		windowIndex:    windowIndex,
		title:          title,
		list:           list,
		shader:         shader,
		physics:        mbt.NewMPhysics(gravity),
		prevCursorPos:  mmath.NewMVec2(),
		overrideOffset: mmath.NewMVec3(),
		modelRenderers: make([]*render.ModelRenderer, 0),
		motions:        make([]*vmd.VmdMotion, 0),
		vmdDeltas:      make([]*delta.VmdDeltas, 0),
	}

	tooltipRenderer, err := mgl.NewTooltipRenderer()
	if err != nil {
		return nil, err
	}
	vw.tooltipRenderer = tooltipRenderer

	glWindow.SetCloseCallback(vw.closeCallback)
	glWindow.SetScrollCallback(vw.scrollCallback)
	glWindow.SetKeyCallback(vw.keyCallback)
	glWindow.SetMouseButtonCallback(vw.mouseCallback)
	glWindow.SetCursorPosCallback(vw.cursorPosCallback)
	glWindow.SetFocusCallback(vw.focusCallback)
	glWindow.SetIconifyCallback(vw.iconifyCallback)
	glWindow.SetSizeCallback(vw.sizeCallback)

	if !isProd {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)                // 同期的なデバッグ出力有効
		gl.DebugMessageCallback(vw.debugMessageCallback, nil) // デバッグコールバック
	}

	// ウィンドウの位置を設定
	vw.SetPos(positionX, positionY)
	// 設定保持
	vw.list.shared.SetInitializedViewWindow(windowIndex, true)
	// ウィンドウハンドルを保持
	handle := int32(uintptr(unsafe.Pointer(glfw.GetCurrentContext().GetWin32Window())))
	vw.list.shared.SetViewerWindowHandle(windowIndex, handle)

	return vw, nil
}

func (vw *ViewWindow) Title() string {
	return vw.title
}

func (vw *ViewWindow) SetTitle(title string) {
	vw.Window.SetTitle(title)
}

func (vw *ViewWindow) resetCameraPosition(yaw, pitch float64) {
	vw.shader.Camera().ResetPosition(yaw, pitch)

	// カメラ同期が有効なら、他のウィンドウへも同じカメラ設定を反映
	vw.syncCameraToOthers()
}

func (vw *ViewWindow) render() {
	w, h := vw.GetSize()
	if w == 0 && h == 0 {
		// ウィンドウが最小化されている場合は描画しない
		return
	}

	vw.MakeContextCurrent()

	// リサイズ（サイズが変わってなければ何もしない）
	vw.shader.Resize(w, h)

	// override が有効かつサブウィンドウの場合、カメラを調整してオーバーライド描画
	if len(vw.list.windowList) > 1 && vw.list.shared.IsShowOverride() && vw.windowIndex != 0 {
		// サブウィンドウ側のカメラを調整（調整後の状態でレンダリングする）
		vw.syncSizeToOthers(vw.list.windowList[0].GetSize())

		vw.adjustCameraForOverride()
		vw.shader.OverrideRenderer().Bind()
	} else {
		// メインウィンドウや override が無効の場合は MSAA FBO へ描画
		vw.shader.Msaa().Bind()
	}

	// 深度バッファのクリア
	gl.ClearColor(0.7, 0.7, 0.7, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 隠面消去
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// マルチサンプル有効化
	gl.Enable(gl.MULTISAMPLE)

	// シェーダーのカメラ設定更新
	vw.shader.UpdateCamera()
	drawRigidBodyFront := vw.list.shared.IsShowRigidBodyFront()
	drawRigidBodyBack := vw.list.shared.IsShowRigidBodyBack()
	highlightEnabled := drawRigidBodyFront || drawRigidBodyBack

	// 床描画
	vw.renderFloor()

	for i, modelRenderer := range vw.modelRenderers {
		if modelRenderer == nil || vw.vmdDeltas[i] == nil {
			continue
		}

		// モデルをレンダリング
		modelRenderer.Render(vw.shader, vw.list.shared, vw.vmdDeltas[i])
	}

	// 物理デバッグ描画
	vw.physics.DrawDebugLines(vw.shader, highlightEnabled, vw.list.shared.IsShowJoint(), drawRigidBodyFront)

	// 剛体デバッグが有効な場合のみハイライト描画
	if drawRigidBodyFront || drawRigidBodyBack {
		vw.physics.DrawDebugHighlight(vw.shader, drawRigidBodyFront)
	}

	// 剛体デバッグ状態変更の検出とハイライトクリア
	vw.physics.CheckAndClearHighlightOnDebugChange(drawRigidBodyFront || drawRigidBodyBack)

	// 剛体デバッグが有効な場合のみ2秒経過ハイライト自動クリアをチェック
	if drawRigidBodyFront || drawRigidBodyBack {
		vw.physics.CheckAndClearExpiredHighlight()
	}

	// 剛体デバッグが有効な場合のみツールチップを表示
	if (drawRigidBodyFront || drawRigidBodyBack) && vw.tooltipRenderer != nil {
		if hover := vw.physics.DebugHoverInfo(); hover != nil && hover.RigidBody != nil {
			text := fmt.Sprintf("%s(G%d)", hover.RigidBody.Name(), int(hover.RigidBody.CollisionGroup))
			vw.tooltipRenderer.Render(text, float32(vw.cursorX), float32(vw.cursorY), w, h)
		}
	}

	// 描画終了後のFBO解除
	if len(vw.list.windowList) > 1 && vw.list.shared.IsShowOverride() && vw.windowIndex != 0 {
		// サブウィンドウの場合、override FBO のアンバインド後にその内容をファイル出力
		vw.shader.OverrideRenderer().Unbind()
	} else {
		// メインウィンドウの場合は MSAA FBO の解決とアンバインド
		vw.shader.Msaa().Resolve()
		vw.shader.Msaa().Unbind()
	}

	// メインウィンドウでは、サブウィンドウの描画内容（overrideテクスチャ）を半透明合成して描画
	if len(vw.list.windowList) > 1 && vw.list.shared.IsShowOverride() &&
		vw.windowIndex == 0 && vw.shader.OverrideRenderer().SharedTextureIDPtr() != nil {
		vw.shader.OverrideRenderer().Resolve()
	}

	vw.SwapBuffers()
}

func (vw *ViewWindow) renderFloor() {
	vw.shader.UseProgram(rendering.ProgramTypeFloor)
	vw.shader.FloorRenderer().Bind()
	vw.shader.FloorRenderer().Render()
	vw.shader.FloorRenderer().Unbind()

	gl.UseProgram(0)
}

// adjustCameraForOverride は、サブウィンドウのカメラを
// メインウィンドウ側の人物モデルの NECK_ROOT と TRUNK_ROOT の位置に合わせる補正を行います。
// 体格が異なる場合でも、各骨間距離に応じたスケール補正を導入しています。
func (vw *ViewWindow) adjustCameraForOverride() {
	// サブウィンドウのみ対象
	if vw.windowIndex == 0 {
		return
	}
	mainVW := vw.list.windowList[0]
	if len(mainVW.vmdDeltas) == 0 || len(vw.vmdDeltas) == 0 {
		return
	}

	// 合わせるボーン名
	targetBoneNames := []string{pmx.TRUNK_ROOT.String()}
	if vw.list.shared.IsShowOverrideUpper() {
		// 上半身合わせ
		targetBoneNames = append(targetBoneNames, pmx.NECK_ROOT.String(), pmx.ARM.Left(), pmx.ARM.Right(), pmx.WRIST.Left(), pmx.WRIST.Right())
	} else if vw.list.shared.IsShowOverrideLower() {
		// 下半身合わせ
		targetBoneNames = append(targetBoneNames, pmx.LEG_CENTER.String(), pmx.LEG.Left(), pmx.LEG.Right(), pmx.ANKLE.Left(), pmx.ANKLE.Right())
	} else {
		// 合わせない場合、そのまま返す
		return
	}

	// 合わせる対象のボーンが1つでもなかった場合は処理しない
	for _, boneName := range targetBoneNames {
		if len(mainVW.vmdDeltas) == 0 || mainVW.vmdDeltas[0] == nil || !mainVW.vmdDeltas[0].Bones.ContainsByName(boneName) ||
			len(vw.vmdDeltas) == 0 || vw.vmdDeltas[0] == nil || !vw.vmdDeltas[0].Bones.ContainsByName(boneName) {
			return
		}
	}

	boneProjections := make([][]*mmath.MVec3, 0, 2)
	boneNDCs := make([][]*mmath.MVec3, 0, 2)
	boneDistances := make([][]float64, 0, 2)

	for n, w := range []*ViewWindow{mainVW, vw} {
		boneProjections = append(boneProjections, make([]*mmath.MVec3, 0, len(targetBoneNames)))
		boneNDCs = append(boneNDCs, make([]*mmath.MVec3, 0, len(targetBoneNames)))
		boneDistances = append(boneDistances, make([]float64, 0, len(targetBoneNames)-1))

		// ウィンドウサイズを取得
		width, height := w.GetSize()

		for _, boneName := range targetBoneNames {
			// ボーンのワールド座標を取得
			bonePos := w.vmdDeltas[0].Bones.GetByName(boneName).FilledGlobalPosition()
			// ボーン位置を NDC に変換
			projectionPoint, ndcPoint := projectPoint(bonePos, w.shader.Camera(), width, height)
			boneProjections[n] = append(boneProjections[n], projectionPoint)
			boneNDCs[n] = append(boneNDCs[n], ndcPoint)
		}

		boneDistances[n] = boneNDCs[n][0].Distances(boneNDCs[n][1:])
	}

	// 縮尺の中央値
	boneScales := make([]float64, 0, len(boneDistances[0]))
	for m := range len(boneDistances[0]) {
		if !mmath.NearEquals(boneDistances[0][m], 0.0, 1e-3) {
			boneScales = append(boneScales, boneDistances[1][m]/boneDistances[0][m])
		}
	}

	scaleRatio := mmath.Median(boneScales)

	// ボーン間の差分を取得
	boneDiffXs := make([]float64, 0, len(targetBoneNames))
	boneDiffYs := make([]float64, 0, len(targetBoneNames))
	for m := range len(boneNDCs[0]) {
		boneDiffXs = append(boneDiffXs, boneNDCs[0][m].X-boneNDCs[1][m].X)
		boneDiffYs = append(boneDiffYs, boneNDCs[0][m].Y-boneNDCs[1][m].Y)
	}

	// 差分の中央値を取る
	avgDiffX := mmath.Median(boneDiffXs)
	avgDiffY := mmath.Median(boneDiffYs)

	// 差分が十分小さければ調整は不要
	if mmath.NearEquals(avgDiffX, 0.0, 1e-5) && mmath.NearEquals(avgDiffY, 0.0, 1e-5) {
		return
	}

	// カメラ設定を取得
	cam := vw.shader.Camera()

	// 1. ベースとなるカメラ設定をメインウィンドウと同期
	mainCam := mainVW.shader.Camera()
	cam.FieldOfView += mainCam.FieldOfView - max(mainCam.FieldOfView/float32(scaleRatio), 1.0)

	// カメラの視点ベクトルを取得
	viewVector := cam.LookAtCenter.Subed(cam.Position).Normalize()
	// 右方向ベクトルを取得
	rightVector := viewVector.Cross(cam.Up).Normalize()
	// 上方向ベクトルを取得
	upVector := rightVector.Cross(viewVector).Normalize()

	// 右方向と上方向への移動量を計算
	rightMove := rightVector.MulScalar(float64(avgDiffX) * scaleRatio)
	upMove := upVector.MulScalar(-float64(avgDiffY) * scaleRatio)

	// カメラ位置と注視点を調整
	cam.Position.Add(rightMove).Add(upMove)
	cam.LookAtCenter.Add(rightMove).Add(upMove)

	// カメラ位置を角度から再計算
	cam.ResetPosition(mainCam.Yaw, mainCam.Pitch)

	// 更新したカメラ設定を適用
	vw.shader.SetCamera(cam)
}

// projectPoint は、与えられたワールド座標 point を、指定されたカメラ(cam)とウィンドウサイズ(w,h)に基づき
// 正規化デバイス座標（NDC）に変換して返します。
func projectPoint(point *mmath.MVec3, cam *rendering.Camera, w, h int) (projectionPoint, ndcPoint *mmath.MVec3) {
	// プロジェクション行列とビュー行列を取得（mgl32.Mat4）
	proj := cam.GetProjectionMatrix(w, h)
	view := cam.GetViewMatrix()

	// mgl64.Vec3 を mgl32.Vec4 に変換（w=1）
	p := mgl32.Vec4{
		float32(point.X),
		float32(point.Y),
		float32(point.Z),
		1.0,
	}
	// クリップ座標を計算
	clip := proj.Mul4(view).Mul4x1(p)
	// パースペクティブ除算により NDC を算出
	ndc := clip.Mul(1.0 / clip.W())

	projectionPoint = &mmath.MVec3{X: float64(ndc.X()) * float64(w), Y: float64(ndc.Y()) * float64(h), Z: float64(ndc.Z())}
	ndcPoint = &mmath.MVec3{X: float64(ndc.X()), Y: float64(ndc.Y()), Z: float64(ndc.Z())}
	return projectionPoint, ndcPoint
}

func (vw *ViewWindow) saveDeltaMotions(frame float32) {
	deltaIndex := vw.list.shared.SaveDeltaIndex(vw.windowIndex)

	for n := range vw.modelRenderers {
		if vw.vmdDeltas[n] == nil {
			continue
		}

		deltaMotion := vw.list.shared.LoadDeltaMotion(vw.windowIndex, n, deltaIndex)

		vw.vmdDeltas[n].Bones.ForEach(func(index int, value *delta.BoneDelta) bool {
			if value == nil || value.Bone == nil {
				return true // 続行
			}

			// 変形情報をモーションに保存
			bf := vmd.NewBoneFrame(frame)
			bf.Position = value.FramePosition
			bf.Rotation = value.FrameRotation
			bf.UnitRotation = value.FilledUnitMatrix().Quaternion() // (モーフ含む)トータル回転も保存
			deltaMotion.AppendBoneFrame(value.Bone.Name(), bf)

			return true // 続行
		})

		vw.list.shared.StoreDeltaMotion(vw.windowIndex, n, deltaIndex, deltaMotion)
	}
}

func (vw *ViewWindow) updateWind(frame float32) {
	windMotion := vw.list.shared.LoadWindMotion(vw.windowIndex)

	if windMotion == nil {
		return
	}

	enabledF := windMotion.WindEnabledFrames.Get(frame)
	directionF := windMotion.WindDirectionFrames.Get(frame)
	speedF := windMotion.WindSpeedFrames.Get(frame)
	liftCoeffF := windMotion.WindLiftCoeffFrames.Get(frame)
	dragCoeffF := windMotion.WindDragCoeffFrames.Get(frame)
	randomnessF := windMotion.WindRandomnessFrames.Get(frame)
	turbulenceFreqHzF := windMotion.WindTurbulenceFreqHzFrames.Get(frame)

	vw.physics.EnableWind(enabledF.Enabled)
	vw.physics.SetWind(directionF.Direction, speedF.Speed, randomnessF.Randomness)
	vw.physics.SetWindAdvanced(dragCoeffF.DragCoeff, liftCoeffF.LiftCoeff, turbulenceFreqHzF.TurbulenceFreqHz)
}

// getWorldPosition は指定されたマウス座標からワールド座標位置を取得します
func (vw *ViewWindow) getWorldPosition(x, y int) (*mmath.MVec3, []*delta.VmdDeltas, *mmath.MMat4) {
	mlog.V("x=%d, y=%d\n", x, y)

	// ウィンドウサイズを取得
	w, h := vw.GetSize()

	// プロジェクション行列とビュー行列を取得
	projection := vw.shader.Camera().GetProjectionMatrix(w, h)
	view := vw.shader.Camera().GetViewMatrix()

	mlog.V("Projection: %s\n", projection.String())
	mlog.V("CameraPosition: %s, LookAtCenterPosition: %s\n",
		vw.shader.Camera().Position.String(), vw.shader.Camera().LookAtCenter.String())
	mlog.V("View: %s\n", view.String())

	// MSAAから深度値を読み取る
	depth := vw.shader.Msaa().ReadDepthAt(x, y, w, h)

	// スクリーン座標からワールド座標に変換
	worldCoords, err := mgl32.UnProject(
		mgl32.Vec3{float32(x), float32(h) - float32(y), depth},
		view, projection, 0, 0, w, h)
	if err != nil {
		mlog.V("UnProject error: %v\n", err)
		return nil, nil, nil
	}

	// mmath.MVec3形式に変換
	worldPos := &mmath.MVec3{X: float64(worldCoords.X()), Y: float64(worldCoords.Y()), Z: float64(worldCoords.Z())}
	mlog.V("WorldPosResult: x=%.7f, y=%.7f, z=%.7f (%.7f)\n", worldPos.X, worldPos.Y, worldPos.Z, depth)

	// ビュー行列の逆行列を計算してmmath.MMat4形式に変換
	viewInv := view.Inv()
	viewMat := &mmath.MMat4{
		float64(viewInv[0]), float64(viewInv[1]), float64(viewInv[2]), float64(viewInv[3]),
		float64(viewInv[4]), float64(viewInv[5]), float64(viewInv[6]), float64(viewInv[7]),
		float64(viewInv[8]), float64(viewInv[9]), float64(viewInv[10]), float64(viewInv[11]),
		float64(viewInv[12]), float64(viewInv[13]), float64(viewInv[14]), float64(viewInv[15]),
	}

	// 現在のVmdDeltasをコピー
	vmdDeltas := make([]*delta.VmdDeltas, len(vw.vmdDeltas))
	copy(vmdDeltas, vw.vmdDeltas)

	return worldPos, vmdDeltas, viewMat
}
