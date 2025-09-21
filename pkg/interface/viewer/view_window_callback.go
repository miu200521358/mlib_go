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
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/walk/pkg/walk"
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

	// if vw.leftButtonPressed {
	// 	// 左クリックした場合、デバッグモードによって処理を分岐

	// 	// 左クリックはカーソル位置を取得
	// 	// if vw.ctrlPressed {
	// 	// 	vw.leftCursorRemoveWindowPositions[mgl32.Vec2{float32(xpos), float32(ypos)}] = 0.0
	// 	// } else {
	// 	// 	vw.leftCursorWindowPositions[mgl32.Vec2{float32(xpos), float32(ypos)}] = 0.0
	// 	// }
	// }

	if !vw.updatedPrevCursor {
		vw.prevCursorPos.X = xpos
		vw.prevCursorPos.Y = ypos
		vw.updatedPrevCursor = true
		return
	}

	if vw.rightButtonPressed {
		// 右クリックはカメラの角度を更新
		vw.updateCameraAngleByCursor(xpos, ypos)
		return
	} else if vw.middleButtonPressed {
		// 中クリックはカメラ位置と中心を移動
		vw.updateCameraPositionByCursor(xpos, ypos)
		return
	}

	vw.prevCursorPos.X = xpos
	vw.prevCursorPos.Y = ypos

	if vw.list.shared.IsShowRigidBodyFront() || vw.list.shared.IsShowRigidBodyBack() {
		// 剛体デバッグ中の場合、剛体選択を行う
		vw.selectRigidBodyByCursor(xpos, ypos)
	}

	// mlog.I("Cursor Position: (%.2f, %.2f)", xpos, ypos)
}

// selectRigidBodyByCursor はカーソル位置に基づいて剛体を選択する
// クリック位置 → レイ → ヒット剛体名をログ
func (vw *ViewWindow) selectRigidBodyByCursor(xpos, ypos float64) {
	if vw.physics == nil {
		return
	}

	// 実ビューポートとカメラ取得
	width, height := vw.GetSize()
	if width == 0 || height == 0 {
		return
	}
	cam := vw.shader.Camera()
	if cam == nil {
		return
	}

	var rayFrom, rayTo *mmath.MVec3
	var ndcX, ndcY float64

	// getWorldPositionを使用して実際のワールド座標を取得
	worldPos, _, _ := vw.getWorldPosition(int(xpos), int(ypos))

	if worldPos != nil {
		// 成功：実際の深度値を使用した精密なレイテスト
		// カメラからワールド座標への方向ベクトル
		direction := worldPos.Subed(cam.Position).Normalize()

		// Direction分手前から開始（100.0m手前）
		rayFrom = worldPos.Subed(direction.MulScalar(100.0))
		// Direction方向に奥行きを持たせる（100.0奥）
		rayTo = worldPos.Added(direction.MulScalar(100.0))

		// デバッグ用NDC計算
		ndcX = (2.0*float64(xpos))/float64(width) - 1.0
		ndcY = 1.0 - (2.0*float64(ypos))/float64(height)

		mlog.D("Using precise ray based on depth: worldPos=%v direction=%v", worldPos, direction)
	} else {
		// フォールバック：従来の実装を使用
		mlog.D("Fallback to traditional ray casting (no depth available)")

		// NDC
		ndcX = (2.0*float64(xpos))/float64(width) - 1.0
		ndcY = 1.0 - (2.0*float64(ypos))/float64(height)

		// 投影パラメータ
		aspect := float64(cam.AspectRatio)
		fovRad := mmath.DegToRad(float64(cam.FieldOfView))
		tanFov := math.Tan(fovRad * 0.5)

		// 視空間方向
		dirCam := (&mmath.MVec3{
			X: ndcX * aspect * tanFov,
			Y: ndcY * tanFov,
			Z: -1.0,
		}).Normalized()

		// カメラ基底 → 世界方向
		forward := cam.LookAtCenter.Subed(cam.Position).Normalize()
		right := forward.Cross(cam.Up).Normalize()
		up := right.Cross(forward).Normalize()
		dirWorld := (&mmath.MVec3{
			X: dirCam.X*right.X + dirCam.Y*up.X + dirCam.Z*forward.X,
			Y: dirCam.X*right.Y + dirCam.Y*up.Y + dirCam.Z*forward.Y,
			Z: dirCam.X*right.Z + dirCam.Y*up.Z + dirCam.Z*forward.Z,
		}).Normalized()

		// ニア面の少し先から撃つ（内側/数値不安定回避）
		rayFrom = cam.Position.Added(dirWorld.MulScalar(float64(cam.NearPlane)))
		rayTo = cam.Position.Added(dirWorld.MulScalar(float64(cam.FarPlane)))
	}

	// レイテスト実行
	btRayFrom := bt.NewBtVector3(float32(rayFrom.X), float32(rayFrom.Y), float32(rayFrom.Z))
	defer bt.DeleteBtVector3(btRayFrom)
	btRayTo := bt.NewBtVector3(float32(rayTo.X), float32(rayTo.Y), float32(rayTo.Z))
	defer bt.DeleteBtVector3(btRayTo)

	cb := bt.NewBtClosestRayCallback(btRayFrom, btRayTo)
	defer bt.DeleteBtClosestRayCallback(cb)

	vw.physics.GetWorld().RayTest(btRayFrom, btRayTo, cb) // レイキャスト実行

	hasHit := cb.HasHit()
	frac := cb.GetHitFraction()
	hitObj := cb.GetCollisionObject()

	// 逆引きして剛体名を取る
	modelIdx, pmxRB, ok := vw.physics.FindRigidBodyByCollisionHit(hitObj, hasHit)

	if hasHit && ok && pmxRB != nil {
		mlog.I("pick: ndc=(%.3f,%.3f) from=%v to=%v hasHit=%v frac=%.5f model=%d name=%s",
			ndcX, ndcY, rayFrom, rayTo, hasHit, frac, modelIdx, pmxRB.Name())
	} else {
		mlog.I("pick: ndc=(%.3f,%.3f) from=%v to=%v hasHit=%v frac=%.5f (hitObj=%v) (reverseLookup ok=%v)",
			ndcX, ndcY, rayFrom, rayTo, hasHit, frac, hitObj, ok)
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
