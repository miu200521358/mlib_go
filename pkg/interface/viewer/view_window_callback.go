//go:build windows
// +build windows

package viewer

import (
	"fmt"
	"math"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/walk/pkg/walk"
)

const (
	// 16～32グループ運用なら 0xFFFF や 0xFFFFFFFF など、あなたの運用ビット幅に合わせる。
	collisionAllFilterMask = int(^uint32(0)) // 全ビット1
)

// 直角の定数値
const rightAngle = 89.9

// closeCallback はウィンドウのクローズイベントを処理する
func (vw *ViewWindow) closeCallback(w *glfw.Window) {
	// controllerStateを読み取り
	if !vw.list.appConfig.IsCloseConfirm() {
		vw.list.shared.SetClosed(true)
		return
	}
	if !vw.list.shared.IsClosed() {
		// ビューワーがまだ閉じていない場合のみ、確認ダイアログを表示
		if result := walk.MsgBox(
			nil,
			mi18n.T("終了確認"),
			mi18n.T("終了確認メッセージ"),
			walk.MsgBoxIconQuestion|walk.MsgBoxOKCancel,
		); result == walk.DlgCmdOK {
			vw.list.shared.SetClosed(true)
		}
	}
}

// CameraPreset はカメラの視点プリセットを定義
type CameraPreset struct {
	Name  string  // プリセット名（デバッグ用）
	Yaw   float64 // 水平方向の角度
	Pitch float64 // 垂直方向の角度
}

// カメラの視点プリセット定義
var cameraPresets = map[glfw.Key]CameraPreset{
	glfw.KeyKP1: {"Bottom", 0, -rightAngle}, // 下面から
	glfw.KeyKP2: {"Front", 0, 0},            // 正面から
	glfw.KeyKP4: {"Left", -rightAngle, 0},   // 左面から
	glfw.KeyKP5: {"Top", 0, rightAngle},     // 上面から
	glfw.KeyKP6: {"Right", rightAngle, 0},   // 右面から
	glfw.KeyKP8: {"Back", 180, 0},           // 背面から
}

// keyCallback はキーボードのイベントを処理する
func (vw *ViewWindow) keyCallback(
	w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey,
) {
	// 修飾キーの処理
	switch action {
	case glfw.Press:
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			vw.shiftPressed = true
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			vw.ctrlPressed = true
			return
		}
	case glfw.Release:
		switch key {
		case glfw.KeyLeftShift, glfw.KeyRightShift:
			vw.shiftPressed = false
			return
		case glfw.KeyLeftControl, glfw.KeyRightControl:
			vw.ctrlPressed = false
			return
		}
	}

	// カメラプリセットの適用
	if preset, exists := cameraPresets[key]; exists {
		vw.resetCameraPosition(preset.Yaw, preset.Pitch)
	}
}

// mouseCallback はマウスボタンのイベントを処理する
func (vw *ViewWindow) mouseCallback(
	w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey,
) {
	switch action {
	case glfw.Press:
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = true
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = true
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = true
		}
	case glfw.Release:
		switch button {
		case glfw.MouseButtonLeft:
			vw.leftButtonPressed = false
			if vw.list.shared.IsShowRigidBodyFront() || vw.list.shared.IsShowRigidBodyBack() {
				// 剛体デバッグ表示中なら剛体選択とハイライト
				vw.selectRigidBodyByCursor(vw.cursorX, vw.cursorY)
			} else if vw.list.shared.IsAnyBoneVisible() {
				// ボーン表示中ならボーン選択とハイライト
				vw.selectBoneByCursor(vw.cursorX, vw.cursorY)
			}
		case glfw.MouseButtonMiddle:
			vw.middleButtonPressed = false
		case glfw.MouseButtonRight:
			vw.rightButtonPressed = false
		}
	}
}

// cursorPosCallback はカーソル位置のイベントを処理する
func (vw *ViewWindow) cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	vw.cursorX = xpos
	vw.cursorY = ypos

	if !vw.updatedPrevCursor {
		vw.prevCursorPos.X = xpos
		vw.prevCursorPos.Y = ypos
		vw.updatedPrevCursor = true
		return
	}

	if vw.rightButtonPressed {
		// 右クリックはカメラの角度を更新
		vw.updateCameraAngleByCursor(xpos, ypos)
	} else if vw.middleButtonPressed {
		// 中クリックはカメラ位置と中心を移動
		vw.updateCameraPositionByCursor(xpos, ypos)
	}

	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos
}

// calculateRayFromTo はマウス座標からレイのFROMとTO座標を計算する共通関数
func (vw *ViewWindow) calculateRayFromTo(xpos, ypos float64) (*mmath.MVec3, *mmath.MVec3, error) {
	// ビューポートとカメラの検証
	width, height := vw.GetSize()
	if width == 0 || height == 0 {
		return nil, nil, fmt.Errorf("invalid viewport size: %dx%d", width, height)
	}
	cam := vw.shader.Camera()
	if cam == nil {
		return nil, nil, fmt.Errorf("camera is nil")
	}

	// 現在のカメラ状態を使用してマトリックスを取得
	projection := cam.GetProjectionMatrix(width, height)
	view := cam.GetViewMatrix()

	// NDC座標を計算
	ndcX := (2.0*float64(xpos))/float64(width) - 1.0
	ndcY := 1.0 - (2.0*float64(ypos))/float64(height)

	// NDCからワールド座標への変換（near平面とfar平面）
	nearWorld, err1 := mgl32.UnProject(
		mgl32.Vec3{float32(xpos), float32(height) - float32(ypos), 0.0},
		view, projection, 0, 0, width, height)

	farWorld, err2 := mgl32.UnProject(
		mgl32.Vec3{float32(xpos), float32(height) - float32(ypos), 1.0},
		view, projection, 0, 0, width, height)

	var rayFrom, rayTo *mmath.MVec3

	if err1 == nil && err2 == nil {
		// UnProjectが成功した場合：正確なレイを構築
		rayFrom = &mmath.MVec3{
			X: float64(nearWorld.X()),
			Y: float64(nearWorld.Y()),
			Z: float64(nearWorld.Z()),
		}
		rayTo = &mmath.MVec3{
			X: float64(farWorld.X()),
			Y: float64(farWorld.Y()),
			Z: float64(farWorld.Z()),
		}
		mlog.D("Camera-aware ray: from=%v to=%v (NDC: %.3f, %.3f)", rayFrom, rayTo, ndcX, ndcY)
		return rayFrom, rayTo, nil
	}

	// フォールバック：カメラからの直線レイ
	mlog.D("UnProject failed, using fallback ray casting")

	// 投影パラメータ
	aspect := float64(cam.AspectRatio)
	fovRad := mmath.DegToRad(float64(cam.FieldOfView))
	tanFov := math.Tan(fovRad * 0.5)

	// 視空間方向ベクトル
	dirCam := (&mmath.MVec3{
		X: ndcX * aspect * tanFov,
		Y: ndcY * tanFov,
		Z: -1.0,
	}).Normalized()

	// カメラの現在の方向基底を使用
	forward := cam.LookAtCenter.Subed(cam.Position).Normalize()
	right := forward.Cross(cam.Up).Normalize()
	up := right.Cross(forward).Normalize()

	// ワールド座標での方向ベクトル
	dirWorld := (&mmath.MVec3{
		X: dirCam.X*right.X + dirCam.Y*up.X + dirCam.Z*forward.X,
		Y: dirCam.X*right.Y + dirCam.Y*up.Y + dirCam.Z*forward.Y,
		Z: dirCam.X*right.Z + dirCam.Y*up.Z + dirCam.Z*forward.Z,
	}).Normalized()

	rayFrom = cam.Position.Added(dirWorld.MulScalar(float64(cam.NearPlane)))
	rayTo = cam.Position.Added(dirWorld.MulScalar(float64(cam.FarPlane)))
	return rayFrom, rayTo, nil
}

// selectRigidBodyByCursor はカーソル位置に基づいて剛体を選択する
// カメラの現在位置・角度を考慮した正確なレイキャストを実行
func (vw *ViewWindow) selectRigidBodyByCursor(xpos, ypos float64) {
	if vw.physics == nil {
		return
	}

	// 共通レイ計算関数を使用
	rayFrom, rayTo, err := vw.calculateRayFromTo(xpos, ypos)
	if err != nil {
		mlog.W("レイ計算に失敗しました: %v", err)
		return
	}

	// NDC座標を計算（ログ用）
	width, height := vw.GetSize()
	ndcX := (2.0*float64(xpos))/float64(width) - 1.0
	ndcY := 1.0 - (2.0*float64(ypos))/float64(height)

	// レイテスト実行
	btRayFrom := bt.NewBtVector3(float32(rayFrom.X), float32(rayFrom.Y), float32(rayFrom.Z))
	defer bt.DeleteBtVector3(btRayFrom)
	btRayTo := bt.NewBtVector3(float32(rayTo.X), float32(rayTo.Y), float32(rayTo.Z))
	defer bt.DeleteBtVector3(btRayTo)

	cb := bt.NewBtClosestRayCallback(btRayFrom, btRayTo)
	defer bt.DeleteBtClosestRayCallback(cb)

	cb.SetCollisionFilterGroup(collisionAllFilterMask)
	cb.SetCollisionFilterMask(collisionAllFilterMask)

	vw.physics.GetWorld().RayTest(btRayFrom, btRayTo, cb)

	hasHit := cb.HasHit()
	frac := cb.GetHitFraction()
	hitObj := cb.GetCollisionObject()

	// 逆引きして剛体名を取得
	modelIdx, rb, ok := vw.physics.FindRigidBodyByCollisionHit(hitObj, hasHit)

	// ハイライト機能の更新
	if hasHit && ok && rb != nil {
		vw.rigidBodyHighlighter.UpdateDebugHoverByRigidBody(modelIdx, rb, true)
		mlog.V("pick: ndc=(%.3f,%.3f) from=%v to=%v hasHit=%v frac=%.5f model=%d name=%s",
			ndcX, ndcY, rayFrom, rayTo, hasHit, frac, modelIdx, rb.PmxRigidBody.Name())
	} else {
		vw.rigidBodyHighlighter.UpdateDebugHoverByRigidBody(0, nil, false)
		mlog.V("pick: ndc=(%.3f,%.3f) from=%v to=%v hasHit=%v frac=%.5f (no hit)",
			ndcX, ndcY, rayFrom, rayTo, hasHit, frac)
	}
}

// selectBoneByCursor は指定されたマウス座標から最も近いボーンを選択します
func (vw *ViewWindow) selectBoneByCursor(xpos, ypos float64) {
	mlog.V("ボーン選択処理開始: x=%.1f, y=%.1f", xpos, ypos)

	if len(vw.vmdDeltas) == 0 || vw.vmdDeltas[0] == nil {
		mlog.V("VmdDeltasが存在しない - ボーン選択中止")
		return
	}

	width, height := vw.GetSize()
	depth := vw.shader.Msaa().ReadDepthAt(int(xpos), int(ypos), width, height)

	// 現在のカメラ状態を使用してマトリックスを取得
	projection := vw.shader.Camera().GetProjectionMatrix(width, height)
	view := vw.shader.Camera().GetViewMatrix()

	// NDCからワールド座標への変換
	world, err := mgl32.UnProject(
		mgl32.Vec3{float32(xpos), float32(height) - float32(ypos), depth},
		view, projection, 0, 0, width, height)
	if err != nil {
		mlog.W("UnProject失敗: %v", err)
		return
	}

	mouseWorldPos := &mmath.MVec3{
		X: -float64(world.X()), // X軸反転
		Y: float64(world.Y()),
		Z: float64(world.Z()),
	}
	boneDistances := make(map[float64][]*mgl.DebugBoneHover) // 距離→ボーン配列マップ

	// 第1段階：最も近いボーンの距離を見つける
	for modelIndex, vmdDeltas := range vw.vmdDeltas {
		if vmdDeltas == nil || vmdDeltas.Bones == nil {
			continue
		}

		vmdDeltas.Bones.ForEach(func(index int, boneDelta *delta.BoneDelta) bool {
			if boneDelta == nil || boneDelta.Bone == nil {
				return true // 続行
			}

			bonePos := boneDelta.FilledGlobalPosition()

			// レイと線分の距離を計算
			distance := mmath.Round(mouseWorldPos.Distance(bonePos), 0.01)
			if _, ok := boneDistances[distance]; !ok {
				boneDistances[distance] = make([]*mgl.DebugBoneHover, 0)
			}

			boneDistances[distance] = append(boneDistances[distance], &mgl.DebugBoneHover{
				ModelIndex: modelIndex,
				Bone:       boneDelta.Bone,
				Distance:   distance,
			})

			return true
		})
	}

	// 最も近いボーンを特定（複数ボーンが同一位置にある場合は全件取得）
	var closestDistance float64 = math.MaxFloat64
	for dist := range boneDistances {
		if dist < closestDistance {
			closestDistance = dist
		}
	}

	closestBones := boneDistances[closestDistance] // 最も近い距離のボーン群

	if len(closestBones) > 0 {
		mlog.V("ボーン選択成功: %d個のボーン (最短距離=%.3f)", len(closestBones), closestDistance)
		vw.boneHighlighter.UpdateDebugHoverByBones(closestBones, true)
	} else {
		mlog.V("近いボーンが見つかりませんでした")
		vw.boneHighlighter.UpdateDebugHoverByBones(nil, false)
	}
}

// updateCameraAngleByCursor はカメラの角度をカーソル位置に基づいて更新する
func (vw *ViewWindow) updateCameraAngleByCursor(xpos, ypos float64) {
	ratio := 0.1
	if vw.shiftPressed {
		ratio *= 3
	} else if vw.ctrlPressed {
		ratio *= 0.1
	}

	// 右クリックはカメラ中心をそのままにカメラ位置を変える
	xOffset := (xpos - vw.prevCursorPos.X) * ratio
	yOffset := (ypos - vw.prevCursorPos.Y) * ratio

	// 方位角と仰角を更新
	vw.resetCameraPosition(vw.shader.Camera().Yaw+xOffset, vw.shader.Camera().Pitch+yOffset)
}

// updateCameraPositionByCursor はカメラ位置と中心をカーソル位置に基づいて更新する
func (vw *ViewWindow) updateCameraPositionByCursor(xpos float64, ypos float64) {
	// 中ボタンが押された場合の処理
	ratio := 0.07
	if vw.shiftPressed {
		ratio *= 3
	} else if vw.ctrlPressed {
		ratio *= 0.1
	}

	xOffset := (vw.prevCursorPos.X - xpos) * ratio
	yOffset := (vw.prevCursorPos.Y - ypos) * ratio

	cam := vw.shader.Camera()

	// カメラの向きに基づいて移動方向を計算
	forward := cam.LookAtCenter.Subed(cam.Position)
	right := forward.Cross(cam.Up).Normalize()
	up := right.Cross(forward.Normalize()).Normalize()

	// 上下移動のベクトルを計算
	upMovement := up.MulScalar(-yOffset)
	// 左右移動のベクトルを計算
	rightMovement := right.MulScalar(-xOffset)

	// 移動ベクトルを合成してカメラ位置と中心を更新
	movement := upMovement.Add(rightMovement)
	cam.Position.Add(movement)
	cam.LookAtCenter.Add(movement)

	vw.shader.SetCamera(cam)
	// カメラ同期が有効なら、他のウィンドウへも同じカメラ設定を反映
	vw.syncCameraToOthers()
}

// syncCameraToOthers は、現在のウィンドウのカメラ設定を他のウィンドウに反映する
func (vw *ViewWindow) syncCameraToOthers() {
	if !vw.list.shared.IsCameraSync() {
		return
	}

	currentCam := vw.shader.Camera()
	for _, otherVW := range vw.list.windowList {
		if otherVW.windowIndex != vw.windowIndex {
			otherCam := otherVW.shader.Camera()
			otherCam.Position.X = currentCam.Position.X
			otherCam.Position.Y = currentCam.Position.Y
			otherCam.Position.Z = currentCam.Position.Z
			otherCam.LookAtCenter.X = currentCam.LookAtCenter.X
			otherCam.LookAtCenter.Y = currentCam.LookAtCenter.Y
			otherCam.LookAtCenter.Z = currentCam.LookAtCenter.Z
			otherCam.Up.X = currentCam.Up.X
			otherCam.Up.Y = currentCam.Up.Y
			otherCam.Up.Z = currentCam.Up.Z
			otherCam.FieldOfView = currentCam.FieldOfView
			otherCam.AspectRatio = currentCam.AspectRatio
			otherCam.NearPlane = currentCam.NearPlane
			otherCam.FarPlane = currentCam.FarPlane
			otherCam.Yaw = currentCam.Yaw
			otherCam.Pitch = currentCam.Pitch
			otherVW.shader.SetCamera(otherCam)
		}
	}
}

// scrollCallback はマウスホイールのスクロールイベントを処理する
func (vw *ViewWindow) scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	step := float32(1.0)
	if vw.shiftPressed {
		step *= 5
	} else if vw.ctrlPressed {
		step *= 0.1
	}

	cam := vw.shader.Camera()

	if yoff > 0 {
		cam.FieldOfView -= step
		if cam.FieldOfView < 1.0 {
			cam.FieldOfView = 1.0
		}
	} else if yoff < 0 {
		cam.FieldOfView += step
	}

	vw.shader.SetCamera(cam)
	// カメラ同期が有効なら、他のウィンドウへも同じカメラ設定を反映
	vw.syncCameraToOthers()
}

func (vw *ViewWindow) focusCallback(w *glfw.Window, focused bool) {
	if !vw.list.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	if focused && !vw.list.shared.IsLinkingFocus() {
		// ユーザー操作等でウィンドウが前面になった場合に連動フォーカスを発火
		vw.list.shared.TriggerLinkedFocus(vw.windowIndex)
	}
}

func (vw *ViewWindow) sizeCallback(w *glfw.Window, width int, height int) {
	if !vw.list.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	vw.syncSizeToOthers(width, height)
}

// syncCameraToOthers は、現在のウィンドウのカメラ設定を他のウィンドウに反映する
func (vw *ViewWindow) syncSizeToOthers(width, height int) {
	if !vw.list.shared.IsShowOverride() {
		// オーバーライドが無効ならスルー
		return
	}

	for _, otherVW := range vw.list.windowList {
		if otherVW.windowIndex != vw.windowIndex {
			otherVW.SetSize(width, height)
		}
	}
}

func (vw *ViewWindow) iconifyCallback(w *glfw.Window, iconified bool) {
	if !vw.list.shared.IsInitializedAllWindows() {
		// 初期化が終わってない場合、スルー
		return
	}

	if iconified {
		vw.list.shared.SyncMinimize(vw.windowIndex)
	} else {
		vw.list.shared.SyncRestore(vw.windowIndex)
	}
}

// debugMessageCallback はOpenGLのデバッグメッセージを処理する
func (vw *ViewWindow) debugMessageCallback(
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
		// 対象のエラーメッセージを検出
		if strings.Contains(message, "glBlitFramebuffer failed because the framebuffer configurations require that the source and destination sample counts match") {
			// 現在のウィンドウサイズを取得
			width, height := vw.GetSize()
			// middleMsaa を生成
			newMsaa := mgl.NewMiddleMsaa(width, height)
			// shader.SetMsaa() で差し替え
			vw.shader.SetMsaa(newMsaa)
			mlog.W(mi18n.T("中間MSAA差し替え"), message)
			return
		}

		panic(fmt.Errorf("[HIGH] GL CRITICAL ERROR: %v type = 0x%x, severity = 0x%x, message = %s",
			source, glType, severity, message))
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
}
