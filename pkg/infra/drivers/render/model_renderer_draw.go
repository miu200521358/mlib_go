//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"math"
	"slices"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"gonum.org/v1/gonum/spatial/r3"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
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
	// selectedVertexVertices は選択頂点用VBOの元データを保持し、ベースコピーの参照先を維持する。
	selectedVertexVertices []float32

	// SSBO
	ssbo uint32

	// cursorPositionLimit はカーソル位置の上限数。
	cursorPositionLimit int
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
	if md.ssbo != 0 {
		gl.DeleteBuffers(1, &md.ssbo)
		md.ssbo = 0
	}
}

// drawNormal 描画処理：法線表示
func (mr *ModelRenderer) drawNormal(windowIndex int, shader graphics_api.IShader, paddedMatrixes []float32, width, height int) {
	program := shader.Program(graphics_api.ProgramTypeNormal)
	gl.UseProgram(program)

	// 法線表示は深度を書き換えない。
	gl.DepthMask(false)
	defer gl.DepthMask(true)

	mr.normalBufferHandle.Bind()
	mr.normalIbo.Bind()

	// ボーン行列テクスチャ設定（共通関数）
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	normalColor := mgl32.Vec4{0.3, 0.3, 0.6, 0.4}
	colorUniform := mgl.GetUniformLocation(program, mgl.ShaderColor)
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
func (mr *ModelRenderer) drawBone(windowIndex int, shader graphics_api.IShader, bones *model.BoneCollection, shared *state.SharedState, paddedMatrixes []float32, width, height int, debugBoneHover []*graphics_api.DebugBoneHover) {
	// モデルの前面にボーンを描画するため、深度テストの設定を変更
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)
	// ボーン表示は深度を書き換えない。
	gl.DepthMask(false)
	defer gl.DepthMask(true)

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
		int32(mr.bonePointCount),
		gl.UNSIGNED_INT,
		nil,
	)

	mr.bonePointIbo.Unbind()
	mr.bonePointBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

// buildCursorLineVertices はカーソル位置の点列からライン描画用の頂点配列を生成する。
func buildCursorLineVertices(cursorPositions []float32, vertexColor mgl32.Vec4, limit int) []float32 {
	pointCount := len(cursorPositions) / 3
	if pointCount < 2 {
		return nil
	}
	if limit > 0 && pointCount > limit {
		start := (pointCount - limit) * 3
		cursorPositions = cursorPositions[start:]
		pointCount = limit
	}
	if pointCount < 2 {
		return nil
	}
	const cursorLineFieldCount = 7
	vertexCount := (pointCount - 1) * 2
	out := make([]float32, 0, vertexCount*cursorLineFieldCount)
	for i := 0; i < pointCount-1; i++ {
		idx0 := i * 3
		idx1 := (i + 1) * 3
		out = append(
			out,
			cursorPositions[idx0], cursorPositions[idx0+1], cursorPositions[idx0+2],
			vertexColor[0], vertexColor[1], vertexColor[2], vertexColor[3],
			cursorPositions[idx1], cursorPositions[idx1+1], cursorPositions[idx1+2],
			vertexColor[0], vertexColor[1], vertexColor[2], vertexColor[3],
		)
	}
	return out
}

// drawCursorLine は、カーソルの軌跡をラインで描画する処理です。
func (mr *ModelRenderer) drawCursorLine(shader graphics_api.IShader, cursorPositions []float32, vertexColor mgl32.Vec4) {
	if mr.cursorPositionBufferHandle == nil {
		return
	}
	lineVertices := buildCursorLineVertices(cursorPositions, vertexColor, 0)
	if len(lineVertices) == 0 {
		return
	}

	// モデルの前面に描画するため深度テストを一時無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)
	// カーソル線は深度を書き換えない。
	gl.DepthMask(false)
	defer gl.DepthMask(true)

	program := shader.Program(graphics_api.ProgramTypeCursor)
	gl.UseProgram(program)

	colorUniform := mgl.GetUniformLocation(program, mgl.ShaderColor)
	gl.Uniform4fv(colorUniform, 1, &vertexColor[0])

	const cursorLineFieldCount = 7
	vertexCount := len(lineVertices) / cursorLineFieldCount
	mr.cursorPositionBufferHandle.Bind()
	// デバッグバッファ更新はバインド後に行う。
	mr.cursorPositionBufferHandle.UpdateDebugBuffer(lineVertices)
	gl.DrawArrays(gl.LINES, 0, int32(vertexCount))
	mr.cursorPositionBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

// drawSelectedVertex は、選択頂点およびカーソルによる頂点選択の描画処理を行い、
// 更新後の選択頂点インデックスとホバー対象の頂点を返します。
func (mr *ModelRenderer) drawSelectedVertex(
	windowIndex int,
	vertices *model.VertexCollection,
	selectedMaterialIndexes []int,
	nowSelectedVertexes []int,
	nowNoSelectedVertexes []int,
	shader graphics_api.IShader,
	paddedMatrixes []float32,
	width, height int,
	selectionRequest *VertexSelectionRequest,
) ([]int, int) {
	selectionMode := state.SELECTED_VERTEX_MODE_POINT
	depthMode := state.SELECTED_VERTEX_DEPTH_MODE_ALL
	applySelection := false
	removeSelection := false
	cursorPositions := []float32(nil)
	cursorDepths := []float32(nil)
	removeCursorPositions := []float32(nil)
	removeCursorDepths := []float32(nil)
	cursorLinePositions := []float32(nil)
	removeCursorLinePositions := []float32(nil)
	screenWidth := 0
	screenHeight := 0
	rectMin := mmath.Vec2{}
	rectMax := mmath.Vec2{}
	hasRect := false
	if selectionRequest != nil {
		selectionMode = selectionRequest.Mode
		depthMode = selectionRequest.DepthMode
		applySelection = selectionRequest.Apply
		removeSelection = selectionRequest.Remove
		cursorPositions = selectionRequest.CursorPositions
		cursorDepths = selectionRequest.CursorDepths
		removeCursorPositions = selectionRequest.RemoveCursorPositions
		removeCursorDepths = selectionRequest.RemoveCursorDepths
		cursorLinePositions = selectionRequest.CursorLinePositions
		removeCursorLinePositions = selectionRequest.RemoveCursorLinePositions
		screenWidth = selectionRequest.ScreenWidth
		screenHeight = selectionRequest.ScreenHeight
		rectMin = selectionRequest.RectMin
		rectMax = selectionRequest.RectMax
		hasRect = selectionRequest.HasRect
	}

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

	// 選択状態の更新は頂点デルタのみ反映し、元の頂点情報（位置やボーン等）を維持する。
	mr.selectedVertexBufferHandle.UpdateVertexDeltas(selectedVertexDeltas)

	mr.selectedVertexBufferHandle.Bind()
	mr.selectedVertexIbo.Bind()

	// ボーン行列テクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	vertexColor := mgl32.Vec4{1.0, 0.4, 0.0, 0.7}
	colorUniform := mgl.GetUniformLocation(program, mgl.ShaderColor)
	gl.Uniform4fv(colorUniform, 1, &vertexColor[0])
	gl.PointSize(5.0)

	thresholdUniform := mgl.GetUniformLocation(program, mgl.ShaderCursorThreshold)
	cursorThreshold := float32(0)
	if shader.Camera() != nil {
		cursorThreshold = shader.Camera().FieldOfView
	}
	if selectionMode == state.SELECTED_VERTEX_MODE_BOX && applySelection && hasRect {
		cursorThreshold = 1.0e9
	}
	gl.Uniform1f(thresholdUniform, cursorThreshold)

	const maxCursorPositions = 100
	effectiveLimit := mr.effectiveCursorPositionLimit()
	cursorPositionsUniform := mgl.GetUniformLocation(program, mgl.ShaderCursorPositions)
	var cursorValues [maxCursorPositions * 3]float32
	srcCursorPositions := cursorPositions
	maxValueCount := effectiveLimit * 3
	if maxValueCount > len(cursorValues) {
		maxValueCount = len(cursorValues)
	}
	if len(srcCursorPositions) > maxValueCount {
		srcCursorPositions = srcCursorPositions[:maxValueCount]
	}
	if len(srcCursorPositions) > len(cursorValues) {
		srcCursorPositions = srcCursorPositions[:len(cursorValues)]
	}
	cursorDepths = truncateCursorDepths(cursorDepths, len(srcCursorPositions)/3)
	removeCursorDepths = truncateCursorDepths(removeCursorDepths, len(srcCursorPositions)/3)
	copy(cursorValues[:], srcCursorPositions)
	gl.Uniform3fv(cursorPositionsUniform, maxCursorPositions, &cursorValues[0])

	// 選択頂点描画は深度バッファを書き換えずに行う。
	gl.DepthMask(false)
	gl.DrawElements(
		gl.POINTS,
		int32(mr.selectedVertexCount),
		gl.UNSIGNED_INT,
		nil,
	)
	gl.DepthMask(true)

	mr.selectedVertexIbo.Unbind()
	mr.selectedVertexBufferHandle.Unbind()

	gl.UseProgram(0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	if selectionMode == state.SELECTED_VERTEX_MODE_POINT || selectionMode == state.SELECTED_VERTEX_MODE_BOX {
		// カーソル軌跡は淡い黄色で表示する。
		cursorLineColor := mgl32.Vec4{0.95, 1.0, 0.75, 0.8}
		if len(cursorLinePositions) > 0 {
			mr.drawCursorLine(shader, cursorLinePositions, cursorLineColor)
		}
		if selectionMode == state.SELECTED_VERTEX_MODE_POINT && len(removeCursorLinePositions) > 0 {
			mr.drawCursorLine(shader, removeCursorLinePositions, cursorLineColor)
		}
	}

	visibleVertexFlags := buildVisibleVertexFlags(vertices, mr.Model.Materials, mr.ModelDrawer.faces, selectedMaterialIndexes)
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
	if len(selectedMaterialIndexes) == 0 {
		// 材質が全て未選択の場合、頂点の選択/ホバーは無効化する。
		return selectedSetToSlice(selectedSet), -1
	}
	truncatedCursorPositions := cursorPositions
	if len(truncatedCursorPositions) > maxValueCount {
		truncatedCursorPositions = truncatedCursorPositions[:maxValueCount]
	}
	needCursorPositions := len(truncatedCursorPositions) > 0
	if selectionMode == state.SELECTED_VERTEX_MODE_BOX && applySelection && hasRect {
		needCursorPositions = true
	}
	if vertexCount == 0 || mr.ssbo == 0 || !needCursorPositions {
		return selectedSetToSlice(selectedSet), -1
	}

	// シェーダ側のSSBO書き込み完了を待ってから読み出す。
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, mr.ssbo)
	positions := make([]float32, vertexCount*4)
	if len(positions) > 0 {
		gl.GetBufferSubData(gl.SHADER_STORAGE_BUFFER, 0, len(positions)*4, gl.Ptr(&positions[0]))
	}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	if applySelection {
		switch {
		case selectionMode == state.SELECTED_VERTEX_MODE_BOX && hasRect:
			minX, minY, maxX, maxY, ok := normalizeSelectionRect(rectMin, rectMax, screenWidth, screenHeight)
			if ok {
				view, projection, ok := selectionViewProjection(shader, screenWidth, screenHeight)
				if ok {
					rectX0, rectY0, rectW, rectH, rectOk := selectionRectToPixels(minX, minY, maxX, maxY, screenWidth, screenHeight)
					depthMap := []float32(nil)
					depthModeFront := depthMode == state.SELECTED_VERTEX_DEPTH_MODE_FRONT
					if depthModeFront && rectOk && shader != nil && shader.Msaa() != nil {
						if reader, ok := shader.Msaa().(interface {
							ReadDepthRegion(x, y, width, height, framebufferHeight int) []float32
						}); ok {
							depthMap = reader.ReadDepthRegion(rectX0, rectY0, rectW, rectH, screenHeight)
						}
					}
					depthTolerance := depthToleranceFromBuffer()
					logger := logging.DefaultLogger()
					for i := 0; i+3 < len(positions); i += 4 {
						if positions[i+3] < 0 {
							continue
						}
						idx := i / 4
						if visibleVertexFlags != nil && !visibleVertexFlags[idx] {
							continue
						}
						pos := mmath.Vec3{
							Vec: r3.Vec{
								X: float64(positions[i]), Y: float64(positions[i+1]), Z: float64(positions[i+2]),
							},
						}
						screenX, screenY, depth, ok := projectToScreenForSelection(pos, view, projection, screenWidth, screenHeight)
						if !ok {
							continue
						}
						if screenX < minX || screenX > maxX || screenY < minY || screenY > maxY {
							continue
						}
						depthAt := float32(-1)
						if depthModeFront {
							if !rectOk || len(depthMap) != rectW*rectH {
								continue
							}
							px := int(math.Floor(screenX))
							py := int(math.Floor(screenY))
							if px < rectX0 || py < rectY0 || px >= rectX0+rectW || py >= rectY0+rectH {
								continue
							}
							ix := px - rectX0
							iy := (rectH - 1) - (py - rectY0)
							depthAt = depthMap[iy*rectW+ix]
							if depthAt <= 0.0 || depthAt >= 1.0 {
								continue
							}
							// 深度バッファ上で手前にある頂点のみ選択対象とする。
							if depth > depthAt+depthTolerance {
								continue
							}
						}
						if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER) {
							logger.Verbose(
								logging.VERBOSE_INDEX_VIEWER,
								"ボックス選択:頂点 index=%d pos=[x=%.4f, y=%.4f, z=%.4f] screen=(%.2f, %.2f) depth=%.6f depthAt=%.6f remove=%t",
								idx, pos.X, pos.Y, pos.Z, screenX, screenY, depth, depthAt, removeSelection,
							)
						}
						if removeSelection {
							delete(selectedSet, idx)
						} else {
							selectedSet[idx] = struct{}{}
						}
					}
				}
			}
		default:
			view, projection, ok := selectionViewProjection(shader, screenWidth, screenHeight)
			if ok {
				type candidate struct {
					Index    int
					ScreenX  float64
					ScreenY  float64
					Depth    float32
					Distance float32
				}
				candidates := make([]candidate, 0)
				minX := math.MaxFloat64
				minY := math.MaxFloat64
				maxX := -math.MaxFloat64
				maxY := -math.MaxFloat64
				for i := 0; i+3 < len(positions); i += 4 {
					distance := positions[i+3]
					if distance < 0 {
						continue
					}
					idx := i / 4
					if visibleVertexFlags != nil && !visibleVertexFlags[idx] {
						continue
					}
					pos := mmath.Vec3{
						Vec: r3.Vec{
							X: float64(positions[i]), Y: float64(positions[i+1]), Z: float64(positions[i+2]),
						},
					}
					screenX, screenY, depth, ok := projectToScreenForSelection(pos, view, projection, screenWidth, screenHeight)
					if !ok {
						continue
					}
					if screenX < minX {
						minX = screenX
					}
					if screenY < minY {
						minY = screenY
					}
					if screenX > maxX {
						maxX = screenX
					}
					if screenY > maxY {
						maxY = screenY
					}
					candidates = append(candidates, candidate{
						Index:    idx,
						ScreenX:  screenX,
						ScreenY:  screenY,
						Depth:    depth,
						Distance: distance,
					})
				}
				if len(candidates) > 0 {
					if depthMode == state.SELECTED_VERTEX_DEPTH_MODE_ALL {
						for _, candidate := range candidates {
							if removeSelection {
								delete(selectedSet, candidate.Index)
								continue
							}
							selectedSet[candidate.Index] = struct{}{}
						}
					} else {
						activeCursorDepths := cursorDepths
						if removeSelection && len(removeCursorDepths) > 0 {
							activeCursorDepths = removeCursorDepths
						}
						rectX0, rectY0, rectW, rectH, rectOk := selectionRectToPixels(minX, minY, maxX, maxY, screenWidth, screenHeight)
						depthMap := []float32(nil)
						if rectOk && shader != nil && shader.Msaa() != nil {
							if reader, ok := shader.Msaa().(interface {
								ReadDepthRegion(x, y, width, height, framebufferHeight int) []float32
							}); ok {
								depthMap = reader.ReadDepthRegion(rectX0, rectY0, rectW, rectH, screenHeight)
							}
						}
						depthTolerance := depthToleranceFromBuffer()
						frontIndex := -1
						frontDepthDiff := float32(math.MaxFloat32)
						frontDepth := float32(math.MaxFloat32)
						frontDistance := float32(math.MaxFloat32)
						useCursorDepth := len(activeCursorDepths) > 0
						// 画素深度で可視面を絞り込み、カーソル深度との差分が最小の頂点を選択する。
						for _, candidate := range candidates {
							if !rectOk || len(depthMap) != rectW*rectH {
								continue
							}
							px := int(math.Floor(candidate.ScreenX))
							py := int(math.Floor(candidate.ScreenY))
							if px < rectX0 || py < rectY0 || px >= rectX0+rectW || py >= rectY0+rectH {
								continue
							}
							ix := px - rectX0
							iy := (rectH - 1) - (py - rectY0)
							depthAt := depthMap[iy*rectW+ix]
							if depthAt <= 0.0 || depthAt >= 1.0 {
								continue
							}
							// 深度バッファ上で手前にある頂点のみ選択対象とする。
							if candidate.Depth > depthAt+depthTolerance {
								continue
							}
							if useCursorDepth {
								depthDiff, ok := minCursorDepthDifference(candidate.Depth, activeCursorDepths)
								if !ok {
									continue
								}
								if depthDiff < frontDepthDiff ||
									(depthDiff == frontDepthDiff && candidate.Depth < frontDepth) ||
									(depthDiff == frontDepthDiff && candidate.Depth == frontDepth && candidate.Distance < frontDistance) {
									frontDepthDiff = depthDiff
									frontDepth = candidate.Depth
									frontDistance = candidate.Distance
									frontIndex = candidate.Index
								}
								continue
							}
							// 深度取得できない場合は、最前面(最小深度)を優先する。
							if candidate.Depth < frontDepth ||
								(candidate.Depth == frontDepth && candidate.Distance < frontDistance) {
								frontDepth = candidate.Depth
								frontDistance = candidate.Distance
								frontIndex = candidate.Index
							}
						}
						if frontIndex >= 0 {
							if removeCursorPositions != nil {
								delete(selectedSet, frontIndex)
							} else {
								selectedSet[frontIndex] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	hoverIndex := -1
	hoverDistance := float32(math.MaxFloat32)
	for i := 0; i+3 < len(positions); i += 4 {
		distance := positions[i+3]
		if distance < 0 {
			continue
		}
		idx := i / 4
		if visibleVertexFlags != nil && !visibleVertexFlags[idx] {
			continue
		}
		if _, ok := selectedSet[idx]; !ok {
			continue
		}
		if distance < hoverDistance {
			hoverDistance = distance
			hoverIndex = idx
		}
	}

	// 選択頂点インデックスの更新結果を返す。
	return selectedSetToSlice(selectedSet), hoverIndex
}

// truncateCursorDepths はカーソル深度配列をカーソル位置数に合わせて切り詰める。
func truncateCursorDepths(depths []float32, positionCount int) []float32 {
	if len(depths) == 0 || positionCount <= 0 {
		return nil
	}
	if positionCount >= len(depths) {
		return depths
	}
	return depths[:positionCount]
}

// depthToleranceFromBuffer は深度バッファの量子化幅から比較用の許容値を算出する。
// OpenGL定数が参照できない環境があるため、ここでは24bitを既定値として扱う。
func depthToleranceFromBuffer() float32 {
	const depthBits = uint(24)
	const maxValue = float32((1 << depthBits) - 1)
	step := 1.0 / maxValue
	return step * 0.5
}

// minCursorDepthDifference は候補深度とカーソル深度の差分最小値を返す。
func minCursorDepthDifference(depth float32, cursorDepths []float32) (float32, bool) {
	if len(cursorDepths) == 0 {
		return 0, false
	}
	minDiff := float32(math.MaxFloat32)
	for _, cursorDepth := range cursorDepths {
		if cursorDepth <= 0.0 || cursorDepth >= 1.0 {
			continue
		}
		diff := float32(math.Abs(float64(depth - cursorDepth)))
		if diff < minDiff {
			minDiff = diff
		}
	}
	if minDiff == float32(math.MaxFloat32) {
		return 0, false
	}
	return minDiff, true
}

// selectionViewProjection は選択判定用のビュー・射影行列を構築する。
func selectionViewProjection(shader graphics_api.IShader, width, height int) (mgl32.Mat4, mgl32.Mat4, bool) {
	if shader == nil || width <= 0 || height <= 0 {
		return mgl32.Mat4{}, mgl32.Mat4{}, false
	}
	cam := shader.Camera()
	if cam == nil || cam.Position == nil || cam.LookAtCenter == nil || cam.Up == nil {
		return mgl32.Mat4{}, mgl32.Mat4{}, false
	}
	projection := mgl32.Perspective(
		mgl32.DegToRad(cam.FieldOfView),
		float32(width)/float32(height),
		cam.NearPlane,
		cam.FarPlane,
	)
	view := mgl32.LookAtV(
		mgl.NewGlVec3(cam.Position),
		mgl.NewGlVec3(cam.LookAtCenter),
		mgl.NewGlVec3(cam.Up),
	)
	return view, projection, true
}

// projectToScreenForSelection はワールド座標をスクリーン座標へ変換し、深度値を返す。
func projectToScreenForSelection(pos mmath.Vec3, view, projection mgl32.Mat4, width, height int) (float64, float64, float32, bool) {
	if width <= 0 || height <= 0 {
		return 0, 0, 0, false
	}
	// SSBOの座標はOpenGL座標系なので、そのまま投影する。
	clip := projection.Mul4(view).Mul4x1(mgl32.Vec4{float32(pos.X), float32(pos.Y), float32(pos.Z), 1})
	if clip.W() == 0 {
		return 0, 0, 0, false
	}
	ndc := clip.Mul(1.0 / clip.W())
	if ndc.Z() < -1.0 || ndc.Z() > 1.0 {
		return 0, 0, 0, false
	}
	screenX := (float64(ndc.X()) + 1.0) * 0.5 * float64(width)
	screenY := (1.0 - float64(ndc.Y())) * 0.5 * float64(height)
	depth := float32((ndc.Z() + 1.0) * 0.5)
	return screenX, screenY, depth, true
}

// selectionRectToPixels は選択矩形をピクセル単位に変換する。
func selectionRectToPixels(minX, minY, maxX, maxY float64, width, height int) (int, int, int, int, bool) {
	if width <= 0 || height <= 0 {
		return 0, 0, 0, 0, false
	}
	x0 := int(math.Floor(minX))
	y0 := int(math.Floor(minY))
	x1 := int(math.Ceil(maxX))
	y1 := int(math.Ceil(maxY))
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 >= width {
		x1 = width - 1
	}
	if y1 >= height {
		y1 = height - 1
	}
	if x1 < x0 || y1 < y0 {
		return 0, 0, 0, 0, false
	}
	rectW := x1 - x0 + 1
	rectH := y1 - y0 + 1
	if rectW <= 0 || rectH <= 0 {
		return 0, 0, 0, 0, false
	}
	return x0, y0, rectW, rectH, true
}

// normalizeSelectionRect は選択矩形の座標を正規化する。
func normalizeSelectionRect(minPos, maxPos mmath.Vec2, width, height int) (float64, float64, float64, float64, bool) {
	if width <= 0 || height <= 0 {
		return 0, 0, 0, 0, false
	}
	minX := math.Min(minPos.X, maxPos.X)
	minY := math.Min(minPos.Y, maxPos.Y)
	maxX := math.Max(minPos.X, maxPos.X)
	maxY := math.Max(minPos.Y, maxPos.Y)
	if maxX < 0 || maxY < 0 || minX > float64(width) || minY > float64(height) {
		return 0, 0, 0, 0, false
	}
	minX = max(minX, 0)
	minY = max(minY, 0)
	maxX = min(maxX, float64(width))
	maxY = min(maxY, float64(height))
	if minX > maxX || minY > maxY {
		return 0, 0, 0, 0, false
	}
	return minX, minY, maxX, maxY, true
}

// buildVisibleVertexFlags は選択中の材質に属する頂点だけを可視判定用に抽出する。
func buildVisibleVertexFlags(
	vertices *model.VertexCollection,
	materials *model.MaterialCollection,
	faces []uint32,
	selectedMaterialIndexes []int,
) []bool {
	if vertices == nil || materials == nil || len(faces) == 0 {
		return nil
	}
	vertexCount := vertices.Len()
	if vertexCount == 0 {
		return nil
	}
	if len(selectedMaterialIndexes) == 0 {
		return nil
	}

	selectedSet := make(map[int]struct{}, len(selectedMaterialIndexes))
	for _, idx := range selectedMaterialIndexes {
		selectedSet[idx] = struct{}{}
	}

	flags := make([]bool, vertexCount)
	offset := 0
	for materialIndex, material := range materials.Values() {
		count := 0
		if material != nil {
			count = material.VerticesCount
		}
		if count < 0 {
			count = 0
		}
		if _, ok := selectedSet[materialIndex]; ok {
			end := offset + count
			if end > len(faces) {
				end = len(faces)
			}
			for _, faceIndex := range faces[offset:end] {
				if int(faceIndex) >= 0 && int(faceIndex) < vertexCount {
					flags[int(faceIndex)] = true
				}
			}
		}
		offset += count
		if offset >= len(faces) {
			break
		}
	}
	return flags
}

// --- 内部ヘルパー関数 ---

// effectiveCursorPositionLimit はカーソル位置の上限値を取得する。
func (mr *ModelRenderer) effectiveCursorPositionLimit() int {
	const defaultCursorPositionLimit = 100
	const maxCursorPositions = 100
	limit := mr.cursorPositionLimit
	if limit <= 0 {
		limit = defaultCursorPositionLimit
	}
	if limit > maxCursorPositions {
		return maxCursorPositions
	}
	return limit
}

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
func (mr *ModelRenderer) fetchBoneLineDeltas(bones *model.BoneCollection, shared *state.SharedState, info boneDebugInfo, debugBoneHover []*graphics_api.DebugBoneHover) ([]int, [][]float32) {
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
func (mr *ModelRenderer) fetchBonePointDeltas(bones *model.BoneCollection, shared *state.SharedState, info boneDebugInfo, debugBoneHover []*graphics_api.DebugBoneHover) ([]int, [][]float32) {
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
