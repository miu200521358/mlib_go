//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"context"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// VertexSelectionRequest は選択頂点の更新要求をまとめる。
type VertexSelectionRequest struct {
	Mode                      state.SelectedVertexMode
	DepthMode                 state.SelectedVertexDepthMode
	Apply                     bool
	Remove                    bool
	CursorPositions           []float32
	CursorDepths              []float32
	RemoveCursorPositions     []float32
	RemoveCursorDepths        []float32
	CursorLinePositions       []float32
	RemoveCursorLinePositions []float32
	ScreenWidth               int
	ScreenHeight              int
	RectMin                   mmath.Vec2
	RectMax                   mmath.Vec2
	HasRect                   bool
}

// ModelRenderer は、PMXモデル全体の描画処理を統括する構造体です。
// バッファの初期化は model_renderer_buffer.go に、描画処理は model_renderer_draw.go に分割して実装します。
type ModelRenderer struct {
	*ModelDrawer

	// ウィンドウインデックス（複数ウィンドウ対応用）
	windowIndex int
	// モデルインデックス（選択材質などの参照用）
	modelIndex int

	// 描画対象のモデル（ドメイン層）
	Model *model.PmxModel
	// 共有状態の元モデル（再読込判定用）
	SourceModel *model.PmxModel

	// モデルのハッシュ値（モデル更新検出などに使用）
	hash string

	// 各材質ごとのメッシュ描画オブジェクト
	meshes []*MeshRenderer
	// ボーン行列テクスチャ用の再利用バッファ
	boneMatrixScratch []float32
	// テクスチャ管理（解放のため保持）
	textureManager *TextureManager
	// 選択材質のマスク
	selectedMaterialMask []bool
	// 選択材質の更新バージョン
	selectedMaterialVersion uint64
}

// ModelRenderBaseResult はベース描画の結果を表す。
type ModelRenderBaseResult struct {
	// ボーン行列テクスチャ用の行列
	PaddedMatrixes []float32
	// 行列テクスチャの幅
	MatrixWidth int
	// 行列テクスチャの高さ
	MatrixHeight int
	// 選択材質インデックス
	SelectedMaterialIndexes []int
}

// NewModelRendererEmpty は空のModelRendererを生成する。
func NewModelRendererEmpty() *ModelRenderer {
	return &ModelRenderer{
		ModelDrawer:             &ModelDrawer{},
		meshes:                  make([]*MeshRenderer, 0),
		selectedMaterialVersion: ^uint64(0),
	}
}

// NewModelRenderer は、新しい ModelRenderer を生成します。
// ここでは、モデルのバッファ初期化や各材質ごとの MeshRenderer の生成も行います。
func NewModelRenderer(windowIndex int, modelData *model.PmxModel) *ModelRenderer {
	if modelData == nil {
		return nil
	}
	// 法線/選択/ボーンは必要時に生成するため、基本バッファのみ準備する。
	bufferData, err := PrepareModelRendererBufferDataWithOptions(
		context.Background(),
		modelData,
		0,
		ModelRendererBufferOptions{},
	)
	if err != nil {
		logging.DefaultLogger().Warn("描画バッファ準備に失敗しました: %v", err)
	}
	return NewModelRendererWithPreparedData(windowIndex, modelData, bufferData)
}

// NewModelRendererWithPreparedData は準備済みのバッファデータからModelRendererを生成します。
func NewModelRendererWithPreparedData(windowIndex int, modelData *model.PmxModel, bufferData *ModelRendererBufferData) *ModelRenderer {
	if modelData == nil {
		return nil
	}
	mr := &ModelRenderer{
		ModelDrawer:             &ModelDrawer{},
		windowIndex:             windowIndex,
		Model:                   modelData,
		modelIndex:              0,
		selectedMaterialVersion: ^uint64(0),
	}

	textureTypes := resolveTextureTypesForRender(modelData)

	tm := NewTextureManager()
	if err := tm.LoadToonTextures(windowIndex); err != nil {
		logging.DefaultLogger().Warn("トゥーンテクスチャの読み込みに失敗しました: %v", err)
	}
	if err := tm.LoadAllTexturesWithTypes(windowIndex, modelData.Textures, modelData.Path(), textureTypes); err != nil {
		logging.DefaultLogger().Warn("テクスチャの読み込みに失敗しました: %v", err)
	}

	// メインのモデル描画用頂点バッファを初期化
	factory := mgl.NewBufferFactory()

	// バッファの初期化 (実装は model_renderer_buffer.go に記述)
	if bufferData == nil {
		mr.initializeBuffers(factory, modelData)
	} else {
		mr.initializeBuffersWithData(factory, modelData, bufferData)
	}

	// 各材質ごとに MeshRenderer を生成
	mr.meshes = make([]*MeshRenderer, modelData.Materials.Len())
	prevVerticesCount := 0
	for index, material := range modelData.Materials.Values() {
		// newMaterialGL は、model.Material から描画用拡張情報 materialGL を生成する関数
		materialExt := newMaterialGL(material, prevVerticesCount, tm)
		// MeshRenderer の生成 (実装は mesh_renderer.go)
		mr.meshes[index] = NewMeshRenderer(factory, mr.faces, materialExt, prevVerticesCount)
		prevVerticesCount += material.VerticesCount
	}

	mr.textureManager = tm

	// モデルのハッシュ値を設定
	mr.hash = modelData.Hash()

	return mr
}

// Hash はモデルのハッシュ値を返す。
func (mr *ModelRenderer) Hash() string {
	return mr.hash
}

// SetCursorPositionLimit はカーソル位置の上限数を設定する。
func (mr *ModelRenderer) SetCursorPositionLimit(limit int) {
	if mr == nil {
		return
	}
	mr.cursorPositionLimit = limit
}

// SetModelIndex はモデルインデックスを設定する。
func (mr *ModelRenderer) SetModelIndex(index int) {
	if mr == nil {
		return
	}
	mr.modelIndex = index
}

// Delete は、ModelRenderer によって生成されたすべての OpenGL リソースを解放します。
func (mr *ModelRenderer) Delete() {
	if mr.bufferHandle != nil {
		mr.bufferHandle.Delete()
	}
	for _, mesh := range mr.meshes {
		if mesh != nil {
			mesh.delete()
		}
	}
	if mr.textureManager != nil {
		mr.textureManager.Delete()
	}
	// ModelDrawerのリソースも解放
	mr.ModelDrawer.delete()
}

// RenderBase はモデルのベース描画（メッシュ/法線/ボーン）を行う。
func (mr *ModelRenderer) RenderBase(
	shader graphics_api.IShader,
	shared *state.SharedState,
	vmdDeltas *delta.VmdDeltas,
	debugBoneHover []*graphics_api.DebugBoneHover,
) *ModelRenderBaseResult {
	if mr == nil || shader == nil || shared == nil || vmdDeltas == nil {
		return nil
	}
	if mr.bufferHandle == nil {
		return nil
	}
	mr.bufferHandle.Bind()
	defer mr.bufferHandle.Unbind()

	// モーフの変更情報をもとに、頂点バッファを更新
	mr.bufferHandle.UpdateVertexDeltas(vmdDeltas.Morphs.Vertices())

	paddedMatrixes, matrixWidth, matrixHeight, err := createBoneMatrixes(vmdDeltas.Bones, mr.boneMatrixScratch)
	if err != nil {
		return nil
	}
	mr.boneMatrixScratch = paddedMatrixes

	selectedMaterialIndexes, selectedMaterialVersion := shared.SelectedMaterialIndexesWithVersion(mr.windowIndex, mr.modelIndex)
	mr.updateSelectedMaterialMask(selectedMaterialIndexes, selectedMaterialVersion)
	selectedMaterialMask := mr.selectedMaterialMask

	// 各材質（メッシュ）ごとの描画
	for i, mesh := range mr.meshes {
		if mesh == nil || mesh.elemBuffer == nil {
			continue
		}
		if selectedMaterialMask == nil || i < 0 || i >= len(selectedMaterialMask) || !selectedMaterialMask[i] {
			// 頂点を持たない材質、選択されていない材質は描画しない
			continue
		}

		mesh.elemBuffer.Bind()

		// ここで、delta.NewMeshDelta を用いて材質ごとの変形情報を生成
		meshDelta := NewMeshDelta(vmdDeltas.Morphs.Materials().Get(i))
		mesh.drawModel(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, meshDelta)

		// エッジ描画（条件に応じて）
		if hasDrawFlag(mesh.material.DrawFlag, model.DRAW_FLAG_DRAWING_EDGE) {
			mesh.drawEdge(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, meshDelta)
		}

		// ワイヤーフレーム描画（UI状態に応じて）
		if shared.HasFlag(state.STATE_FLAG_SHOW_WIRE) {
			mesh.drawWire(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, false)
		}

		mesh.elemBuffer.Unbind()
	}

	// 法線描画
	if shared.HasFlag(state.STATE_FLAG_SHOW_NORMAL) {
		if mr.ensureNormalBuffers() {
			mr.drawNormal(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight)
		}
	}

	// ボーン描画
	if shared.IsAnyBoneVisible() {
		if mr.ensureBoneBuffers() {
			mr.drawBone(mr.windowIndex, shader, mr.Model.Bones, shared, paddedMatrixes, matrixWidth, matrixHeight, debugBoneHover)
		}
	}

	return &ModelRenderBaseResult{
		PaddedMatrixes:          paddedMatrixes,
		MatrixWidth:             matrixWidth,
		MatrixHeight:            matrixHeight,
		SelectedMaterialIndexes: selectedMaterialIndexes,
	}
}

// RenderSelection は選択頂点の描画と選択結果を返す。
func (mr *ModelRenderer) RenderSelection(
	shader graphics_api.IShader,
	shared *state.SharedState,
	selectedVertexIndexes []int,
	selectionRequest *VertexSelectionRequest,
	base *ModelRenderBaseResult,
) ([]int, int) {
	if mr == nil || shader == nil || shared == nil || base == nil {
		return selectedVertexIndexes, -1
	}
	if !shared.HasFlag(state.STATE_FLAG_SHOW_SELECTED_VERTEX) {
		return selectedVertexIndexes, -1
	}
	if !mr.ensureSelectedVertexBuffers() {
		return selectedVertexIndexes, -1
	}
	return mr.drawSelectedVertex(
		mr.windowIndex,
		mr.Model.Vertices,
		base.SelectedMaterialIndexes,
		selectedVertexIndexes,
		nil,
		shader,
		base.PaddedMatrixes,
		base.MatrixWidth,
		base.MatrixHeight,
		selectionRequest,
	)
}

// updateSelectedMaterialMask は選択材質マスクを更新する。
func (mr *ModelRenderer) updateSelectedMaterialMask(selectedMaterialIndexes []int, version uint64) {
	if mr == nil {
		return
	}
	materialCount := len(mr.meshes)
	if materialCount == 0 {
		mr.selectedMaterialMask = nil
		mr.selectedMaterialVersion = version
		return
	}
	if mr.selectedMaterialVersion == version && len(mr.selectedMaterialMask) == materialCount {
		return
	}

	mask := mr.selectedMaterialMask
	if mask == nil || len(mask) != materialCount {
		mask = make([]bool, materialCount)
	} else {
		clear(mask)
	}
	for _, idx := range selectedMaterialIndexes {
		if idx < 0 || idx >= materialCount {
			continue
		}
		mask[idx] = true
	}
	mr.selectedMaterialMask = mask
	mr.selectedMaterialVersion = version
}

// Render は、最新の変形情報 vmdDeltas とアプリケーション状態 appState に基づいてモデルを描画します。
// 描画前にバッファの更新処理を行い、その後各描画パス（メッシュ描画、法線、ボーン、選択頂点など）を呼び出します。
func (mr *ModelRenderer) Render(
	shader graphics_api.IShader,
	shared *state.SharedState,
	vmdDeltas *delta.VmdDeltas,
	debugBoneHover []*graphics_api.DebugBoneHover,
	selectedVertexIndexes []int,
	selectionRequest *VertexSelectionRequest,
) ([]int, int) {
	base := mr.RenderBase(shader, shared, vmdDeltas, debugBoneHover)
	if base == nil {
		return selectedVertexIndexes, -1
	}
	return mr.RenderSelection(shader, shared, selectedVertexIndexes, selectionRequest, base)
}

// resolveTextureTypesForRender は材質参照に基づいてテクスチャ種別を決定する。
func resolveTextureTypesForRender(modelData *model.PmxModel) map[int]model.TextureType {
	if modelData == nil || modelData.Materials == nil || modelData.Textures == nil {
		return nil
	}

	normal := make(map[int]struct{})
	sphere := make(map[int]struct{})
	toon := make(map[int]struct{})

	for _, material := range modelData.Materials.Values() {
		if material == nil {
			continue
		}
		if material.TextureIndex >= 0 {
			normal[material.TextureIndex] = struct{}{}
		}
		if material.SphereMode != model.SPHERE_MODE_INVALID && material.SphereTextureIndex >= 0 {
			sphere[material.SphereTextureIndex] = struct{}{}
		}
		if material.ToonSharingFlag == model.TOON_SHARING_INDIVIDUAL && material.ToonTextureIndex >= 0 {
			toon[material.ToonTextureIndex] = struct{}{}
		}
	}

	types := make(map[int]model.TextureType)
	for _, texture := range modelData.Textures.Values() {
		if texture == nil || !texture.IsValid() {
			continue
		}
		idx := texture.Index()
		if idx < 0 {
			continue
		}
		if _, ok := normal[idx]; ok {
			types[idx] = model.TEXTURE_TYPE_TEXTURE
			continue
		}
		if _, ok := sphere[idx]; ok {
			types[idx] = model.TEXTURE_TYPE_SPHERE
			continue
		}
		if _, ok := toon[idx]; ok {
			types[idx] = model.TEXTURE_TYPE_TOON
		}
	}
	return types
}
