//go:build windows
// +build windows

package render

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

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

	// モデルのハッシュ値（モデル更新検出などに使用）
	hash string

	// 各材質ごとのメッシュ描画オブジェクト
	meshes []*MeshRenderer
}

// NewModelRendererEmpty は空のModelRendererを生成する。
func NewModelRendererEmpty() *ModelRenderer {
	return &ModelRenderer{
		ModelDrawer: &ModelDrawer{},
		meshes:      make([]*MeshRenderer, 0),
	}
}

// NewModelRenderer は、新しい ModelRenderer を生成します。
// ここでは、モデルのバッファ初期化や各材質ごとの MeshRenderer の生成も行います。
func NewModelRenderer(windowIndex int, modelData *model.PmxModel) *ModelRenderer {
	mr := &ModelRenderer{
		ModelDrawer: &ModelDrawer{},
		windowIndex: windowIndex,
		Model:       modelData,
		modelIndex:  0,
	}

	tm := NewTextureManager()
	if err := tm.LoadToonTextures(windowIndex); err != nil {
		logging.DefaultLogger().Warn("トゥーンテクスチャの読み込みに失敗しました: %v", err)
	}
	if err := tm.LoadAllTextures(windowIndex, modelData.Textures, modelData.Path()); err != nil {
		logging.DefaultLogger().Warn("テクスチャの読み込みに失敗しました: %v", err)
	}

	// メインのモデル描画用頂点バッファを初期化
	factory := mgl.NewBufferFactory()

	// バッファの初期化 (実装は model_renderer_buffer.go に記述)
	mr.initializeBuffers(factory, modelData)

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

	// モデルのハッシュ値を設定
	mr.hash = modelData.Hash()

	return mr
}

// Hash はモデルのハッシュ値を返す。
func (mr *ModelRenderer) Hash() string {
	return mr.hash
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
	// ModelDrawerのリソースも解放
	mr.ModelDrawer.delete()
}

// Render は、最新の変形情報 vmdDeltas とアプリケーション状態 appState に基づいてモデルを描画します。
// 描画前にバッファの更新処理を行い、その後各描画パス（メッシュ描画、法線、ボーン、選択頂点など）を呼び出します。
func (mr *ModelRenderer) Render(shader graphics_api.IShader, shared *state.SharedState, vmdDeltas *delta.VmdDeltas, debugBoneHover []*mgl.DebugBoneHover) {
	mr.bufferHandle.Bind()
	defer mr.bufferHandle.Unbind()

	// モーフの変更情報をもとに、頂点バッファを更新
	mr.bufferHandle.UpdateVertexDeltas(vmdDeltas.Morphs.Vertices())

	paddedMatrixes, matrixWidth, matrixHeight, err := createBoneMatrixes(vmdDeltas.Bones)
	if err != nil {
		return
	}

	selectedMaterialIndexes := shared.SelectedMaterialIndexes(mr.windowIndex, mr.modelIndex)

	// 各材質（メッシュ）ごとの描画
	for i, mesh := range mr.meshes {
		if mesh == nil || mesh.elemBuffer == nil || !slices.Contains(selectedMaterialIndexes, i) {
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
		mr.drawNormal(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	// ボーン描画
	if shared.IsAnyBoneVisible() {
		mr.drawBone(mr.windowIndex, shader, mr.Model.Bones, shared, paddedMatrixes, matrixWidth, matrixHeight, debugBoneHover)
	}
}
