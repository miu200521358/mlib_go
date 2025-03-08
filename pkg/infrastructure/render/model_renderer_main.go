//go:build windows
// +build windows

package render

import (
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/interface/state"
)

// ModelRenderer は、PMXモデル全体の描画処理を統括する構造体です。
// バッファの初期化は model_renderer_buffer.go に、描画処理は model_renderer_draw.go に分割して実装します。
type ModelRenderer struct {
	*ModelDrawer

	// ウィンドウインデックス（複数ウィンドウ対応用）
	windowIndex int

	// 描画対象のモデル（ドメイン層）
	Model *pmx.PmxModel

	// モデルのハッシュ値（モデル更新検出などに使用）
	Hash string

	// 各材質ごとのメッシュ描画オブジェクト
	meshes []*MeshRenderer

	// UI・選択情報管理（非表示材質、選択頂点など）
	invisibleMaterialIndexes map[int]struct{}
	selectedVertexes         map[int]struct{}
	noSelectedVertexes       map[int]struct{}
}

// NewModelRenderer は、新しい ModelRenderer を生成します。
// ここでは、モデルのバッファ初期化や各材質ごとの MeshRenderer の生成も行います。
func NewModelRenderer(windowIndex int, model *pmx.PmxModel) *ModelRenderer {
	mr := &ModelRenderer{
		ModelDrawer:              &ModelDrawer{},
		windowIndex:              windowIndex,
		Model:                    model,
		invisibleMaterialIndexes: make(map[int]struct{}),
		selectedVertexes:         make(map[int]struct{}),
		noSelectedVertexes:       make(map[int]struct{}),
	}

	tm := NewTextureManager()
	tm.LoadToonTextures(windowIndex)
	tm.LoadAllTextures(windowIndex, model.Textures, model.Path())

	// メインのモデル描画用頂点バッファを初期化
	factory := mgl.NewBufferFactory()

	// バッファの初期化 (実装は model_renderer_buffer.go に記述)
	mr.initializeBuffers(factory, model)

	// 各材質ごとに MeshRenderer を生成
	mr.meshes = make([]*MeshRenderer, model.Materials.Length())
	prevVerticesCount := 0
	for v := range model.Materials.Iterator() {
		i := v.Index
		m := v.Value
		// newMaterialGL は、pmx.Material から描画用拡張情報 materialGL を生成する関数
		materialExt := newMaterialGL(m, prevVerticesCount, tm)
		// MeshRenderer の生成 (実装は mesh_renderer.go)
		mr.meshes[i] = NewMeshRenderer(factory, mr.faces, materialExt, prevVerticesCount)
		prevVerticesCount += m.VerticesCount
	}

	// モデルのハッシュ値を設定
	mr.Hash = model.Hash()

	return mr
}

// Delete は、ModelRenderer によって生成されたすべての OpenGL リソースを解放します。
func (mr *ModelRenderer) Delete() {
	if mr.bufferHandle != nil {
		mr.bufferHandle.Delete()
	}
	for _, mesh := range mr.meshes {
		if mesh != nil {
			mesh.Delete()
		}
	}
	// ModelDrawerのリソースも解放
	mr.ModelDrawer.Delete()
}

// Render は、最新の変形情報 vmdDeltas とアプリケーション状態 appState に基づいてモデルを描画します。
// 描画前にバッファの更新処理を行い、その後各描画パス（メッシュ描画、法線、ボーン、選択頂点など）を呼び出します。
func (mr *ModelRenderer) Render(shader rendering.IShader, shared *state.SharedState, vmdDeltas *delta.VmdDeltas) {
	mr.bufferHandle.Bind()
	defer mr.bufferHandle.Unbind()

	// モーフの変更情報をもとに、頂点バッファを更新
	mr.bufferHandle.UpdateVertexDeltas(vmdDeltas.Morphs.Vertices)

	paddedMatrixes, matrixWidth, matrixHeight, err := createBoneMatrixes(vmdDeltas.Bones)
	if err != nil {
		return
	}

	// 各材質（メッシュ）ごとの描画
	for i, mesh := range mr.meshes {
		mesh.elemBuffer.Bind()

		// ここで、delta.NewMeshDelta を用いて材質ごとの変形情報を生成
		meshDelta := delta.NewMeshDelta(vmdDeltas.Morphs.Materials.Get(i))
		mesh.DrawModel(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, meshDelta)

		// エッジ描画（条件に応じて）
		if mesh.material.DrawFlag.IsDrawingEdge() {
			mesh.DrawEdge(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, meshDelta)
		}

		// ワイヤーフレーム描画（UI状態に応じて）
		if shared.IsShowWire() && !mr.ExistInvisibleMaterial(mesh.material.Index()) {
			mesh.DrawWire(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, false)
		}

		mesh.elemBuffer.Unbind()
	}

	// 法線描画
	if shared.IsShowNormal() {
		mr.DrawNormal(mr.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	// ボーン描画
	if shared.IsShowBoneAll() || shared.IsShowBoneEffector() || shared.IsShowBoneIk() ||
		shared.IsShowBoneFixed() || shared.IsShowBoneRotate() || shared.IsShowBoneTranslate() ||
		shared.IsShowBoneVisible() {
		mr.DrawBone(mr.windowIndex, shader, mr.Model.Bones, shared, paddedMatrixes, matrixWidth, matrixHeight)
	}

	// 選択頂点描画・カーソルライン描画は model_renderer_draw.go 内で実装済み
	// 例: selected := mr.drawSelectedVertex(...)

}

// ExistInvisibleMaterial は、指定された材質インデックスが非表示設定になっているかを返します。
func (mr *ModelRenderer) ExistInvisibleMaterial(index int) bool {
	_, ok := mr.invisibleMaterialIndexes[index]
	return ok
}
