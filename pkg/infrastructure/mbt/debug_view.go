//go:build windows
// +build windows

package mbt

import (
	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

// newConstBtMDefaultColors はデバッグ描画の色を定義します
func newConstBtMDefaultColors() bt.BtMDefaultColors {
	return bt.NewBtMDefaultColors(
		bt.NewBtVector3(float32(0.0), float32(0.5), float32(0.5)), // activeObject	(物理剛体)
		bt.NewBtVector3(float32(0.5), float32(0.5), float32(0.0)), // deactivatedObject
		bt.NewBtVector3(float32(0.5), float32(0.0), float32(0.5)), // wantsDeactivationObject
		bt.NewBtVector3(float32(1.0), float32(0.0), float32(0.0)), // disabledDeactivationObject
		bt.NewBtVector3(float32(0.0), float32(1.0), float32(0.0)), // disabledSimulationObject	(ボーン追従剛体)
		bt.NewBtVector3(float32(1.0), float32(1.0), float32(0.0)), // aabb
		bt.NewBtVector3(float32(0.0), float32(0.0), float32(1.0)), // contactPoint
	)
}

// mDebugDrawLiner はデバッグ描画のためのラインレンダラーです
type mDebugDrawLiner struct {
	bt.BtMDebugDrawLiner
	vertices          []float32
	debugBufferHandle *mgl.VertexBufferHandle
}

// newMDebugDrawLiner は新しいデバッグドローラインを作成します
func newMDebugDrawLiner() *mDebugDrawLiner {
	ddl := &mDebugDrawLiner{
		vertices:          make([]float32, 0),
		debugBufferHandle: nil,
	}

	// Bulletのドローラインインターフェースを実装
	ddl.BtMDebugDrawLiner = bt.NewDirectorBtMDebugDrawLiner(ddl)

	return ddl
}

// DrawLine は3D空間に線を描画します
func (ddl *mDebugDrawLiner) DrawLine(from bt.BtVector3, to bt.BtVector3, color bt.BtVector3) {
	// 始点の座標と色
	ddl.vertices = append(ddl.vertices, from.GetX(), from.GetY(), from.GetZ(),
		color.GetX(), color.GetY(), color.GetZ(), 0.6)

	// 終点の座標と色
	ddl.vertices = append(ddl.vertices, to.GetX(), to.GetY(), to.GetZ(),
		color.GetX(), color.GetY(), color.GetZ(), 0.6)
}

// drawDebugLines はデバッグ線を描画します
func (ddl *mDebugDrawLiner) drawDebugLines(shader rendering.IShader, isDrawRigidBodyFront bool) {
	if len(ddl.vertices) == 0 {
		return
	}

	// シェーダープログラムの使用開始
	program := shader.Program(rendering.ProgramTypePhysics)
	gl.UseProgram(program)

	// 深度テストの設定
	ddl.configureDepthTest(isDrawRigidBodyFront)

	// 描画処理
	ddl.renderLines()

	// 深度テストを元に戻す
	ddl.restoreDepthTest(isDrawRigidBodyFront)

	// シェーダープログラムの使用終了
	gl.UseProgram(0)

	// 頂点データをクリア
	ddl.vertices = ddl.vertices[:0]
}

// configureDepthTest は深度テストを設定します
func (ddl *mDebugDrawLiner) configureDepthTest(isDrawRigidBodyFront bool) {
	if isDrawRigidBodyFront {
		// モデルメッシュの前面に描画するために深度テストを無効化
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.ALWAYS)
	}
}

// restoreDepthTest は深度テストの設定を元に戻します
func (ddl *mDebugDrawLiner) restoreDepthTest(isDrawRigidBodyFront bool) {
	if isDrawRigidBodyFront {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)
	}
}

// renderLines は線を描画します
func (ddl *mDebugDrawLiner) renderLines() {
	// OpenGLのバッファを初期化
	if ddl.debugBufferHandle == nil {
		// 初回のみバッファを生成
		ddl.debugBufferHandle = mgl.NewBufferFactory().CreateDebugBuffer(gl.Ptr(&ddl.vertices[0]), len(ddl.vertices))
		ddl.debugBufferHandle.Bind()
	} else {
		// 2回目以降はバッファを更新
		ddl.debugBufferHandle.Bind()
		ddl.debugBufferHandle.UpdateDebugBuffer(ddl.vertices)
	}

	// ライン描画
	gl.DrawArrays(gl.LINES, 0, int32(len(ddl.vertices)/7))

	ddl.debugBufferHandle.Unbind()
}

// DrawDebugLines は物理デバッグ情報を描画します
func (mp *MPhysics) DrawDebugLines(
	shader rendering.IShader, visibleRigidBody, visibleJoint, isDrawRigidBodyFront bool,
) {
	if !(visibleRigidBody || visibleJoint) {
		return
	}

	// デバッグモードの設定
	mp.configureDebugDrawMode(visibleRigidBody, visibleJoint)

	// デバッグ情報取得
	mp.world.DebugDrawWorld()

	// デバッグ描画
	mp.liner.drawDebugLines(shader, isDrawRigidBodyFront)
}

// configureDebugDrawMode はデバッグ描画モードを設定します
func (mp *MPhysics) configureDebugDrawMode(visibleRigidBody, visibleJoint bool) {
	debugMode := mp.world.GetDebugDrawer().GetDebugMode()

	if visibleRigidBody {
		debugMode |= int(bt.BtIDebugDrawDBG_DrawWireframe | bt.BtIDebugDrawDBG_DrawContactPoints)
	} else {
		debugMode &= ^int(bt.BtIDebugDrawDBG_DrawWireframe | bt.BtIDebugDrawDBG_DrawContactPoints)
	}

	if visibleJoint {
		debugMode |= int(bt.BtIDebugDrawDBG_DrawConstraints | bt.BtIDebugDrawDBG_DrawConstraintLimits)
	} else {
		debugMode &= ^int(bt.BtIDebugDrawDBG_DrawConstraints | bt.BtIDebugDrawDBG_DrawConstraintLimits)
	}

	mp.world.GetDebugDrawer().SetDebugMode(debugMode)
}
