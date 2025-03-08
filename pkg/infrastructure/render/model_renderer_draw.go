//go:build windows
// +build windows

package render

import (
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
)

// ModelDrawer は、モデル全体の描画処理のうち、バッファ初期化以外の各描画処理を担当する。
// このファイルでは、法線描画、ボーン描画、選択頂点描画、カーソルライン描画などの処理を実装する。
type ModelDrawer struct {
	// 描画用リソース（頂点バッファハンドル）
	bufferHandle *mgl.VertexBufferHandle

	// 以下は初期化済みの VertexBufferHandle と関連バッファ
	normalBufferHandle *mgl.VertexBufferHandle
	normalIbo          *mgl.ElementBuffer
	normalVertices     []float32

	boneLineBufferHandle *mgl.VertexBufferHandle
	boneLineIbo          *mgl.ElementBuffer
	boneLineCount        int
	boneLineIndexes      []int

	bonePointBufferHandle *mgl.VertexBufferHandle
	bonePointIbo          *mgl.ElementBuffer
	bonePointCount        int
	bonePointIndexes      []int

	selectedVertexBufferHandle *mgl.VertexBufferHandle
	selectedVertexIbo          *mgl.ElementBuffer

	cursorPositionBufferHandle *mgl.VertexBufferHandle

	// 頂点情報（SSBOから読み出す場合など）
	vertices []float32
	faces    []uint32

	// SSBO
	ssbo uint32
}

// Delete はModelDrawerが保持するリソースを解放します
func (md *ModelDrawer) Delete() {
	if md.normalBufferHandle != nil {
		md.normalBufferHandle.Delete()
	}
	if md.normalIbo != nil {
		md.normalIbo.Delete()
	}

	if md.boneLineBufferHandle != nil {
		md.boneLineBufferHandle.Delete()
	}
	if md.boneLineIbo != nil {
		md.boneLineIbo.Delete()
	}

	if md.bonePointBufferHandle != nil {
		md.bonePointBufferHandle.Delete()
	}
	if md.bonePointIbo != nil {
		md.bonePointIbo.Delete()
	}

	if md.selectedVertexBufferHandle != nil {
		md.selectedVertexBufferHandle.Delete()
	}
	if md.selectedVertexIbo != nil {
		md.selectedVertexIbo.Delete()
	}

	if md.cursorPositionBufferHandle != nil {
		md.cursorPositionBufferHandle.Delete()
	}
}

// DrawNormal 描画処理：法線表示
func (mr *ModelRenderer) DrawNormal(windowIndex int, shader rendering.IShader, paddedMatrixes []float32, width, height int) {
	program := shader.Program(rendering.ProgramTypeModel)
	gl.UseProgram(program)

	mr.normalBufferHandle.Bind()
	mr.normalIbo.Bind()

	// ボーン行列テクスチャ設定（共通関数）
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	normalColor := mgl32.Vec4{0.3, 0.3, 0.7, 0.5}
	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderColor))
	gl.Uniform4fv(colorUniform, 1, &normalColor[0])

	gl.DrawElements(
		gl.LINES,
		int32(len(mr.normalVertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	mr.normalIbo.Unbind()
	mr.normalBufferHandle.Unbind()

	gl.UseProgram(0)
}

// DrawBone は、ボーン表示（ラインとポイント）の描画処理を行います。
func (mr *ModelRenderer) DrawBone(windowIndex int, shader rendering.IShader, bones *pmx.Bones, shared *state.SharedState, paddedMatrixes []float32, width, height int) {
	// モデルの前面にボーンを描画するため、深度テストの設定を変更
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(rendering.ProgramTypeBone)
	gl.UseProgram(program)

	// ボーン行列テクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	// --- ボーンライン描画 ---
	mr.boneLineBufferHandle.Bind()
	// 取得したデバッグカラー情報
	_, deltas := mr.fetchBoneLineDeltas(bones, shared)
	// 平坦化して VBO を更新
	boneLineData := flattenFloat32Matrix(deltas)
	mr.boneLineBufferHandle.VBO.Bind()
	mr.boneLineBufferHandle.VBO.BufferData(len(boneLineData)*4, gl.Ptr(boneLineData), rendering.BufferUsageStatic)
	mr.boneLineBufferHandle.VBO.Unbind()
	mr.boneLineIbo.Bind()
	gl.DrawElements(
		gl.LINES,
		int32(mr.boneLineCount),
		gl.UNSIGNED_INT,
		nil,
	)
	mr.boneLineIbo.Unbind()
	mr.boneLineBufferHandle.Unbind()

	// --- ボーンポイント描画 ---
	mr.bonePointBufferHandle.Bind()
	_, bonePointDeltas := mr.fetchBonePointDeltas(bones, shared)
	bonePointData := flattenFloat32Matrix(bonePointDeltas)
	mr.bonePointBufferHandle.VBO.Bind()
	mr.bonePointBufferHandle.VBO.BufferData(len(bonePointData)*4, gl.Ptr(bonePointData), rendering.BufferUsageStatic)
	mr.bonePointBufferHandle.VBO.Unbind()
	mr.bonePointIbo.Bind()
	gl.PointSize(5.0)
	gl.DrawElements(
		gl.POINTS,
		int32(bones.Length()),
		gl.UNSIGNED_INT,
		nil,
	)
	mr.bonePointIbo.Unbind()
	mr.bonePointBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

// DrawCursorLine は、カーソルの軌跡をラインで描画する処理です。
func (mr *ModelRenderer) DrawCursorLine(shader rendering.IShader, cursorPositions []float32, vertexColor mgl32.Vec4) {
	// モデルの前面に描画するため深度テストを一時無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(rendering.ProgramTypeCursor)
	gl.UseProgram(program)

	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderColor))
	gl.Uniform4fv(colorUniform, 1, &vertexColor[0])

	mr.cursorPositionBufferHandle.Bind()
	gl.DrawArrays(gl.LINES, 0, int32(len(cursorPositions)/3))
	mr.cursorPositionBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

// DrawSelectedVertex は、選択頂点およびカーソルによる頂点選択の描画処理を行い、
// 更新後の選択頂点インデックスを返します。
func (mr *ModelRenderer) DrawSelectedVertex(
	windowIndex int,
	invisibleMaterialIndexes []int,
	nowSelectedVertexes []int,
	nowNoSelectedVertexes []int,
	shader rendering.IShader,
	paddedMatrixes []float32,
	width, height int,
	cursorPositions []float32,
	removeCursorPositions []float32,
) []int {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(rendering.ProgramTypeSelectedVertex)
	gl.UseProgram(program)

	// 選択頂点のモーフデルタ情報を生成
	selectedVertexDeltas := delta.NewVertexMorphDeltas()
	for _, index := range nowNoSelectedVertexes {
		vd := delta.NewVertexMorphDelta(index)
		vd.Uv = &mmath.MVec2{X: -1, Y: 0}
		selectedVertexDeltas.Update(vd)
	}
	for _, index := range nowSelectedVertexes {
		vd := delta.NewVertexMorphDelta(index)
		vd.Uv = &mmath.MVec2{X: 1, Y: 0}
		selectedVertexDeltas.Update(vd)
	}

	// VBO の更新：従来の BindVertex の代替として、selectedVertexDeltas を []float32 に変換して更新
	vertexData := convertVertexMorphDeltasToFloat32(selectedVertexDeltas)

	mr.selectedVertexBufferHandle.Bind()
	mr.selectedVertexBufferHandle.VBO.Bind()
	mr.selectedVertexBufferHandle.VBO.BufferData(len(vertexData)*4, gl.Ptr(vertexData), rendering.BufferUsageStatic)
	mr.selectedVertexBufferHandle.VBO.Unbind()
	mr.selectedVertexIbo.Bind()

	// ボーン行列テクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	vertexColor := mgl32.Vec4{1.0, 0.4, 0.0, 0.7}
	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderColor))
	gl.Uniform4fv(colorUniform, 1, &vertexColor[0])
	gl.PointSize(5.0)

	thresholdUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderCursorThreshold))
	gl.Uniform1f(thresholdUniform, shader.Camera().FieldOfView)

	cursorPositionsUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderCursorPositions))
	if removeCursorPositions != nil {
		gl.Uniform3fv(cursorPositionsUniform, int32(len(removeCursorPositions)/3), &removeCursorPositions[0])
	} else {
		gl.Uniform3fv(cursorPositionsUniform, int32(len(cursorPositions)/3), &cursorPositions[0])
	}

	gl.DrawElements(
		gl.POINTS,
		int32(len(mr.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	mr.selectedVertexIbo.Unbind()
	mr.selectedVertexBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// 選択頂点インデックスの更新結果を返す（ここではサンプルとして空スライス）
	return []int{}
}

// --- 内部ヘルパー関数 ---

// fetchBoneLineDeltas は、ボーンライン描画用のデバッグカラー情報を取得します。
func (mr *ModelRenderer) fetchBoneLineDeltas(bones *pmx.Bones, shared *state.SharedState) ([]int, [][]float32) {
	indexes := make([]int, len(mr.boneLineIndexes))
	deltas := make([][]float32, len(mr.boneLineIndexes))
	for i, boneIndex := range mr.boneLineIndexes {
		indexes[i] = i
		if bone, err := bones.Get(boneIndex); err == nil {
			deltas[i] = getBoneDebugColor(bone, shared)
		}
	}
	return indexes, deltas
}

// fetchBonePointDeltas は、ボーンポイント描画用のデバッグカラー情報を取得します。
func (mr *ModelRenderer) fetchBonePointDeltas(bones *pmx.Bones, shared *state.SharedState) ([]int, [][]float32) {
	indexes := make([]int, len(mr.bonePointIndexes))
	deltas := make([][]float32, len(mr.bonePointIndexes))
	for i, boneIndex := range mr.bonePointIndexes {
		indexes[i] = i
		if bone, err := bones.Get(boneIndex); err == nil {
			deltas[i] = getBoneDebugColor(bone, shared)
		}
	}
	return indexes, deltas
}

// flattenFloat32Matrix は [][]float32 の各要素を1次元の []float32 にまとめて返します。
func flattenFloat32Matrix(data [][]float32) []float32 {
	var flat []float32
	for _, row := range data {
		flat = append(flat, row...)
	}
	return flat
}

// updateVertexBuffer は、vbo に data の内容を BufferData を用いて更新します。
func updateVertexBuffer(vbo *mgl.VertexBuffer, data []float32) {
	vbo.Bind()
	if data != nil && len(data) > 0 {
		// float32は4バイト
		size := len(data) * 4
		vbo.BufferData(size, unsafe.Pointer(&data[0]), rendering.BufferUsageStatic)
	}
	vbo.Unbind()
}

// convertVertexMorphDeltasToFloat32 は、delta.VertexMorphDeltas から []float32 を生成します。
// ここでは、各 VertexMorphDelta に対して newVertexMorphDeltaGl を呼び出し、その結果を連結します。
func convertVertexMorphDeltasToFloat32(deltas *delta.VertexMorphDeltas) []float32 {
	var result []float32
	for v := range deltas.Iterator() {
		vd := v.Value
		// newVertexMorphDeltaGl は、1つの VertexMorphDelta を []float32 に変換する既存関数です。
		data := newVertexMorphDeltaGl(vd)
		result = append(result, data...)
	}
	return result
}
