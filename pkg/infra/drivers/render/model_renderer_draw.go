//go:build windows
// +build windows

package render

import (
	"slices"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// ModelDrawer は、モデル全体の描画処理のうち、バッファ初期化以外の各描画処理を担当する。
// このファイルでは、法線描画、ボーン描画、選択頂点描画、カーソルライン描画などの処理を実装する。
type ModelDrawer struct {
	// 描画用リソース（頂点バッファハンドル）
	bufferHandle *mgl.VertexBufferHandle

	// 以下は初期化済みの VertexBufferHandle と関連バッファ
	normalBufferHandle *mgl.VertexBufferHandle
	normalIbo          *mgl.IndexBuffer
	normalVertices     []float32
	normalIndexCount   int

	boneLineBufferHandle *mgl.VertexBufferHandle
	boneLineIbo          *mgl.IndexBuffer
	boneLineCount        int
	boneLineIndexes      []int

	bonePointBufferHandle *mgl.VertexBufferHandle
	bonePointIbo          *mgl.IndexBuffer
	bonePointCount        int
	bonePointIndexes      []int

	selectedVertexBufferHandle *mgl.VertexBufferHandle
	selectedVertexIbo          *mgl.IndexBuffer
	selectedVertexCount        int

	cursorPositionBufferHandle *mgl.VertexBufferHandle

	// 頂点情報（SSBOから読み出す場合など）
	vertices []float32
	faces    []uint32

	// SSBO
	ssbo uint32
}

// delete はModelDrawerが保持するリソースを解放します
func (md *ModelDrawer) delete() {
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

// drawNormal 描画処理：法線表示
func (mr *ModelRenderer) drawNormal(windowIndex int, shader graphics_api.IShader, paddedMatrixes []float32, width, height int) {
	program := shader.Program(graphics_api.ProgramTypeNormal)
	gl.UseProgram(program)

	mr.normalBufferHandle.Bind()
	mr.normalIbo.Bind()

	// ボーン行列テクスチャ設定（共通関数）
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	normalColor := mgl32.Vec4{0.3, 0.3, 0.6, 0.4}
	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderColor))
	gl.Uniform4fv(colorUniform, 1, &normalColor[0])

	gl.DrawElements(
		gl.LINES,
		int32(mr.normalIndexCount),
		gl.UNSIGNED_INT,
		nil,
	)

	mr.normalIbo.Unbind()
	mr.normalBufferHandle.Unbind()

	gl.UseProgram(0)
}

// drawBone は、ボーン表示（ラインとポイント）の描画処理を行います。
func (mr *ModelRenderer) drawBone(windowIndex int, shader graphics_api.IShader, bones *model.BoneCollection, shared *state.SharedState, paddedMatrixes []float32, width, height int, debugBoneHover []*mgl.DebugBoneHover) {
	// モデルの前面にボーンを描画するため、深度テストの設定を変更
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(graphics_api.ProgramTypeBone)
	gl.UseProgram(program)

	// ボーン行列テクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	// --- ボーンライン描画 ---
	info := buildBoneDebugInfo(bones)
	mr.boneLineBufferHandle.Bind()
	// 取得したデバッグカラー情報で更新
	boneLineIndexes, boneLineDeltas := mr.fetchBoneLineDeltas(bones, shared, info, debugBoneHover)
	mr.boneLineBufferHandle.UpdateBoneDeltas(boneLineIndexes, boneLineDeltas)

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
	// 取得したデバッグカラー情報で更新
	bonePointIndexes, bonePointDeltas := mr.fetchBonePointDeltas(bones, shared, info, debugBoneHover)
	mr.bonePointBufferHandle.UpdateBoneDeltas(bonePointIndexes, bonePointDeltas)

	mr.bonePointIbo.Bind()

	// ボーンポイントの描画
	gl.PointSize(5.0)
	gl.DrawElements(
		gl.POINTS,
		int32(bones.Len()),
		gl.UNSIGNED_INT,
		nil,
	)

	mr.bonePointIbo.Unbind()
	mr.bonePointBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

// drawCursorLine は、カーソルの軌跡をラインで描画する処理です。
func (mr *ModelRenderer) drawCursorLine(shader graphics_api.IShader, cursorPositions []float32, vertexColor mgl32.Vec4) {
	// モデルの前面に描画するため深度テストを一時無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(graphics_api.ProgramTypeCursor)
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

// drawSelectedVertex は、選択頂点およびカーソルによる頂点選択の描画処理を行い、
// 更新後の選択頂点インデックスを返します。
func (mr *ModelRenderer) drawSelectedVertex(
	windowIndex int,
	vertices *model.VertexCollection,
	invisibleMaterialIndexes []int,
	nowSelectedVertexes []int,
	nowNoSelectedVertexes []int,
	shader graphics_api.IShader,
	paddedMatrixes []float32,
	width, height int,
	cursorPositions []float32,
	removeCursorPositions []float32,
) []int {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(graphics_api.ProgramTypeSelectedVertex)
	gl.UseProgram(program)

	// 選択頂点のモーフデルタ情報を生成
	selectedVertexDeltas := delta.NewVertexMorphDeltas(vertices.Len())
	for _, index := range nowNoSelectedVertexes {
		vd := delta.NewVertexMorphDelta(index)
		vd.Uv = &mmath.Vec2{X: -1, Y: 0}
		selectedVertexDeltas.Update(vd)
	}
	for _, index := range nowSelectedVertexes {
		vd := delta.NewVertexMorphDelta(index)
		vd.Uv = &mmath.Vec2{X: 1, Y: 0}
		selectedVertexDeltas.Update(vd)
	}

	// VBO の更新：従来の BindVertex の代替として、selectedVertexDeltas を []float32 に変換して更新
	vertexData := convertVertexMorphDeltasToFloat32(selectedVertexDeltas)

	mr.selectedVertexBufferHandle.Bind()
	mr.selectedVertexBufferHandle.VBO.Bind()
	mr.selectedVertexBufferHandle.VBO.BufferData(len(vertexData)*4, gl.Ptr(vertexData), graphics_api.BufferUsageStatic)
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

	const maxCursorPositions = 30
	cursorPositionsUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderCursorPositions))
	var cursorValues [maxCursorPositions * 3]float32
	srcCursorPositions := cursorPositions
	if removeCursorPositions != nil {
		srcCursorPositions = removeCursorPositions
	}
	if len(srcCursorPositions) > len(cursorValues) {
		srcCursorPositions = srcCursorPositions[:len(cursorValues)]
	}
	copy(cursorValues[:], srcCursorPositions)
	gl.Uniform3fv(cursorPositionsUniform, maxCursorPositions, &cursorValues[0])

	gl.DrawElements(
		gl.POINTS,
		int32(mr.selectedVertexCount),
		gl.UNSIGNED_INT,
		nil,
	)

	mr.selectedVertexIbo.Unbind()
	mr.selectedVertexBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	selectedSet := make(map[int]struct{}, len(nowSelectedVertexes))
	vertexCount := vertices.Len()
	for _, idx := range nowSelectedVertexes {
		if idx >= 0 && idx < vertexCount {
			selectedSet[idx] = struct{}{}
		}
	}
	for _, idx := range nowNoSelectedVertexes {
		delete(selectedSet, idx)
	}
	if vertexCount == 0 || mr.ssbo == 0 || (len(cursorPositions) == 0 && len(removeCursorPositions) == 0) {
		return selectedSetToSlice(selectedSet)
	}

	// シェーダ側のSSBO書き込み完了を待ってから読み出す。
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, mr.ssbo)
	positions := make([]float32, vertexCount*4)
	if len(positions) > 0 {
		gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, len(positions)*4, gl.Ptr(&positions[0]))
	}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	// w成分に距離が入るため、w >= 0 の頂点を選択対象とする。
	if removeCursorPositions != nil {
		for i := 0; i+3 < len(positions); i += 4 {
			if positions[i+3] >= 0 {
				delete(selectedSet, i/4)
			}
		}
	} else {
		for i := 0; i+3 < len(positions); i += 4 {
			if positions[i+3] >= 0 {
				selectedSet[i/4] = struct{}{}
			}
		}
	}

	// 選択頂点インデックスの更新結果を返す。
	return selectedSetToSlice(selectedSet)
}

// --- 内部ヘルパー関数 ---

// selectedSetToSlice は選択頂点インデックス集合をスライス化する。
func selectedSetToSlice(selectedSet map[int]struct{}) []int {
	if len(selectedSet) == 0 {
		return []int{}
	}
	out := make([]int, 0, len(selectedSet))
	for idx := range selectedSet {
		out = append(out, idx)
	}
	slices.Sort(out)
	return out
}

// fetchBoneLineDeltas は、ボーンライン描画用のデバッグカラー情報を取得します。
func (mr *ModelRenderer) fetchBoneLineDeltas(bones *model.BoneCollection, shared *state.SharedState, info boneDebugInfo, debugBoneHover []*mgl.DebugBoneHover) ([]int, [][]float32) {
	indexes := make([]int, len(mr.boneLineIndexes))
	deltas := make([][]float32, len(mr.boneLineIndexes))
	for i, boneIndex := range mr.boneLineIndexes {
		indexes[i] = i
		if bone, err := bones.Get(boneIndex); err == nil {
			isHover := false
			for _, hover := range debugBoneHover {
				if hover.Bone.Index() == bone.Index() && mr.modelIndex == hover.ModelIndex {
					isHover = true
					break
				}
			}
			deltas[i] = getBoneDebugColor(bone, shared, info, isHover)
		}
	}
	return indexes, deltas
}

// fetchBonePointDeltas は、ボーンポイント描画用のデバッグカラー情報を取得します。
func (mr *ModelRenderer) fetchBonePointDeltas(bones *model.BoneCollection, shared *state.SharedState, info boneDebugInfo, debugBoneHover []*mgl.DebugBoneHover) ([]int, [][]float32) {
	indexes := make([]int, len(mr.bonePointIndexes))
	deltas := make([][]float32, len(mr.bonePointIndexes))
	for i, boneIndex := range mr.bonePointIndexes {
		indexes[i] = i
		if bone, err := bones.Get(boneIndex); err == nil {
			isHover := false
			for _, hover := range debugBoneHover {
				if hover.Bone.Index() == bone.Index() && mr.modelIndex == hover.ModelIndex {
					isHover = true
					break
				}
			}
			deltas[i] = getBoneDebugColor(bone, shared, info, isHover)
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

// convertVertexMorphDeltasToFloat32 は、delta.VertexMorphDeltas から []float32 を生成します。
// ここでは、各 VertexMorphDelta に対して newVertexMorphDeltaGl を呼び出し、その結果を連結します。
func convertVertexMorphDeltasToFloat32(deltas *delta.VertexMorphDeltas) []float32 {
	var result []float32
	deltas.ForEach(func(index int, vd *delta.VertexMorphDelta) bool {
		// newVertexMorphDeltaGl は、1つの VertexMorphDelta を []float32 に変換する既存関数です。
		data := newVertexMorphDeltaGl(vd)
		result = append(result, data...)
		return true
	})
	return result
}
