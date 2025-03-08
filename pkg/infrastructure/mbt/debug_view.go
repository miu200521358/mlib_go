//go:build windows
// +build windows

package mbt

import (
	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

// DebugDrawOptions デバッグ描画のオプション
type DebugDrawOptions struct {
	VisibleRigidBody     bool
	VisibleJoint         bool
	DrawRigidBodyInFront bool
}

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
	debugBufferHandle *mgl.VertexBufferHandle
	vertices          []float32
}

// newMDebugDrawLiner は新しいデバッグドローラインを作成します
func newMDebugDrawLiner() *mDebugDrawLiner {
	ddl := &mDebugDrawLiner{
		vertices: make([]float32, 0),
	}

	// OpenGLのバッファを初期化
	ddl.debugBufferHandle = mgl.NewBufferFactory().CreateDebugBuffer()

	// Bulletのドローラインインターフェースを実装
	ddl.BtMDebugDrawLiner = bt.NewDirectorBtMDebugDrawLiner(ddl)

	return ddl
}

// DrawLine は3D空間に線を描画します
func (ddl *mDebugDrawLiner) DrawLine(from bt.BtVector3, to bt.BtVector3, color bt.BtVector3) {
	// 始点の座標と色
	ddl.vertices = append(ddl.vertices, from.GetX(), from.GetY(), from.GetZ())
	ddl.vertices = append(ddl.vertices, color.GetX(), color.GetY(), color.GetZ(), 0.6)

	// 終点の座標と色
	ddl.vertices = append(ddl.vertices, to.GetX(), to.GetY(), to.GetZ())
	ddl.vertices = append(ddl.vertices, color.GetX(), color.GetY(), color.GetZ(), 0.6)
}

// drawDebugLines はデバッグ線を描画します
func (ddl *mDebugDrawLiner) drawDebugLines(shader *mgl.MShader, isDrawRigidBodyFront bool) {
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
	ddl.vertices = make([]float32, 0)
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
	ddl.debugBufferHandle.Bind()
	ddl.debugBufferHandle.VBO.BufferData(len(ddl.vertices)*4, gl.Ptr(ddl.vertices), rendering.BufferUsageStatic)

	// ライン描画
	gl.DrawArrays(gl.LINES, 0, int32(len(ddl.vertices)/7))

	ddl.debugBufferHandle.Unbind()
}

// DrawDebugLines は物理デバッグ情報を描画します
func (physics *MPhysics) DrawDebugLines(
	shader *mgl.MShader, options DebugDrawOptions,
) {
	if !(options.VisibleRigidBody || options.VisibleJoint) {
		return
	}

	// デバッグモードの設定
	physics.configureDebugDrawMode(options)

	// デバッグ情報取得
	physics.world.DebugDrawWorld()

	// デバッグ描画
	physics.liner.drawDebugLines(shader, options.DrawRigidBodyInFront)
}

// configureDebugDrawMode はデバッグ描画モードを設定します
func (physics *MPhysics) configureDebugDrawMode(options DebugDrawOptions) {
	debugMode := 0

	if options.VisibleRigidBody {
		debugMode |= int(bt.BtIDebugDrawDBG_DrawWireframe | bt.BtIDebugDrawDBG_DrawContactPoints)
	} else {
		debugMode &= ^int(bt.BtIDebugDrawDBG_DrawWireframe | bt.BtIDebugDrawDBG_DrawContactPoints)
	}

	if options.VisibleJoint {
		debugMode |= int(bt.BtIDebugDrawDBG_DrawConstraints | bt.BtIDebugDrawDBG_DrawConstraintLimits)
	} else {
		debugMode &= ^int(bt.BtIDebugDrawDBG_DrawConstraints | bt.BtIDebugDrawDBG_DrawConstraintLimits)
	}

	physics.world.GetDebugDrawer().SetDebugMode(debugMode)
}
