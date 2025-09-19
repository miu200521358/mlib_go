//go:build windows
// +build windows

package mbt

import (
	"github.com/go-gl/gl/v4.4-core/gl"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/bt"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

// mDebugDrawHighlighter は選択された剛体のハイライト描画を担当します
type mDebugDrawHighlighter struct {
	selectedRigidBody     *RigidBodyHit
	highlightVertices     []float32
	highlightBufferHandle *mgl.VertexBufferHandle
}

// newMDebugDrawHighlighter は新しいハイライト描画器を作成します
func newMDebugDrawHighlighter() *mDebugDrawHighlighter {
	return &mDebugDrawHighlighter{
		selectedRigidBody:     nil,
		highlightVertices:     make([]float32, 0),
		highlightBufferHandle: nil,
	}
}

// SetSelectedRigidBody は選択された剛体を設定します
func (ddh *mDebugDrawHighlighter) SetSelectedRigidBody(hit *RigidBodyHit) {
	ddh.selectedRigidBody = hit
	ddh.generateHighlightVertices()
}

// ClearSelection は選択を解除します
func (ddh *mDebugDrawHighlighter) ClearSelection() {
	ddh.selectedRigidBody = nil
	ddh.highlightVertices = ddh.highlightVertices[:0]
}

// generateHighlightVertices は選択された剛体の面描画用頂点を生成します
func (ddh *mDebugDrawHighlighter) generateHighlightVertices() {
	ddh.highlightVertices = ddh.highlightVertices[:0]

	if ddh.selectedRigidBody == nil || ddh.selectedRigidBody.RigidBody == nil {
		return
	}

	rigidBody := ddh.selectedRigidBody.RigidBody

	// 剛体の形状に応じて頂点を生成
	switch rigidBody.ShapeType {
	case pmx.SHAPE_BOX:
		ddh.generateBoxVertices(rigidBody)
	case pmx.SHAPE_SPHERE:
		ddh.generateSphereVertices(rigidBody)
	case pmx.SHAPE_CAPSULE:
		ddh.generateCapsuleVertices(rigidBody)
	}
}

// generateBoxVertices はボックス形状の頂点を生成します
func (ddh *mDebugDrawHighlighter) generateBoxVertices(rigidBody *pmx.RigidBody) {
	// ボックスのサイズ（半径）を float32 に変換
	sx := float32(rigidBody.Size.X)
	sy := float32(rigidBody.Size.Y)
	sz := float32(rigidBody.Size.Z)

	// ボックスの8つの頂点
	vertices := [][]float32{
		{-sx, -sy, -sz}, {sx, -sy, -sz}, {sx, sy, -sz}, {-sx, sy, -sz}, // 底面
		{-sx, -sy, sz}, {sx, -sy, sz}, {sx, sy, sz}, {-sx, sy, sz}, // 上面
	}

	// 面を構成する三角形のインデックス（各面2つの三角形）
	faces := [][3]int{
		// 底面 (-Z)
		{0, 1, 2}, {0, 2, 3},
		// 上面 (+Z)
		{4, 7, 6}, {4, 6, 5},
		// 前面 (-Y)
		{0, 4, 5}, {0, 5, 1},
		// 後面 (+Y)
		{2, 6, 7}, {2, 7, 3},
		// 左面 (-X)
		{0, 3, 7}, {0, 7, 4},
		// 右面 (+X)
		{1, 5, 6}, {1, 6, 2},
	}

	// 半透明のハイライト色 (青系、アルファ0.3)
	color := []float32{0.2, 0.6, 1.0, 0.3}

	// 各三角形の頂点を追加
	for _, face := range faces {
		for _, vertexIndex := range face {
			vertex := vertices[vertexIndex]
			// 位置 + 色
			ddh.highlightVertices = append(ddh.highlightVertices,
				vertex[0]+float32(rigidBody.Position.X),
				vertex[1]+float32(rigidBody.Position.Y),
				vertex[2]+float32(rigidBody.Position.Z))
			ddh.highlightVertices = append(ddh.highlightVertices, color...)
		}
	}
}

// generateSphereVertices は球形状の頂点を生成します（簡易実装）
func (ddh *mDebugDrawHighlighter) generateSphereVertices(rigidBody *pmx.RigidBody) {
	// 簡易的にボックスとして描画
	ddh.generateBoxVertices(rigidBody)
}

// generateCapsuleVertices はカプセル形状の頂点を生成します（簡易実装）
func (ddh *mDebugDrawHighlighter) generateCapsuleVertices(rigidBody *pmx.RigidBody) {
	// 簡易的にボックスとして描画
	ddh.generateBoxVertices(rigidBody)
}

// drawHighlight は選択された剛体のハイライトを描画します
func (ddh *mDebugDrawHighlighter) drawHighlight(shader rendering.IShader, isDrawRigidBodyFront bool) {
	if ddh.selectedRigidBody == nil || len(ddh.highlightVertices) == 0 {
		return
	}

	// シェーダープログラムの使用開始
	program := shader.Program(rendering.ProgramTypePhysics)
	gl.UseProgram(program)

	// ブレンド有効化（半透明描画のため）
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// 深度テストの設定
	ddh.configureDepthTest(isDrawRigidBodyFront)

	// 描画処理
	ddh.renderHighlight()

	// 深度テストを元に戻す
	ddh.restoreDepthTest(isDrawRigidBodyFront)

	// ブレンド無効化
	gl.Disable(gl.BLEND)

	// シェーダープログラムの使用終了
	gl.UseProgram(0)
}

// configureDepthTest は深度テストを設定します
func (ddh *mDebugDrawHighlighter) configureDepthTest(isDrawRigidBodyFront bool) {
	if isDrawRigidBodyFront {
		// モデルメッシュの前面に描画するために深度テストを調整
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)
	}
}

// restoreDepthTest は深度テストの設定を元に戻します
func (ddh *mDebugDrawHighlighter) restoreDepthTest(isDrawRigidBodyFront bool) {
	if isDrawRigidBodyFront {
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LEQUAL)
	}
}

// renderHighlight はハイライトを描画します
func (ddh *mDebugDrawHighlighter) renderHighlight() {
	if len(ddh.highlightVertices) == 0 {
		return
	}

	// OpenGLのバッファを初期化
	if ddh.highlightBufferHandle == nil {
		// 初回のみバッファを生成
		ddh.highlightBufferHandle = mgl.NewBufferFactory().CreateDebugBuffer(gl.Ptr(&ddh.highlightVertices[0]), len(ddh.highlightVertices))
		ddh.highlightBufferHandle.Bind()
	} else {
		// 2回目以降はバッファを更新
		ddh.highlightBufferHandle.Bind()
		ddh.highlightBufferHandle.UpdateDebugBuffer(ddh.highlightVertices)
	}

	// 三角形描画
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(ddh.highlightVertices)/7))

	ddh.highlightBufferHandle.Unbind()
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
