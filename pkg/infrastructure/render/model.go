//go:build windows
// +build windows

package render

import (
	"slices"
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

type RenderModel struct {
	windowIndex              int              // ウィンドウインデックス
	invisibleMaterialIndexes map[int]struct{} // 非表示材質インデックス
	selectedVertexIndexes    map[int]struct{} // 選択頂点インデックス
	noSelectedVertexIndexes  map[int]struct{} // 非選択頂点インデックス
	model                    *pmx.PmxModel    // 元モデル
	hash                     string           // ハッシュ
	meshes                   []*Mesh          // メッシュ
	textures                 []*textureGl     // テクスチャ
	toonTextures             []*textureGl     // トゥーンテクスチャ
	vertices                 []float32        // 頂点情報
	vao                      *mgl.VAO         // 頂点VAO
	vbo                      *mgl.VBO         // 頂点VBO
	normalVertices           []float32        // 法線情報
	normalVao                *mgl.VAO         // 法線VAO
	normalVbo                *mgl.VBO         // 法線VBO
	normalIbo                *mgl.IBO         // 法線IBO
	selectedVertexVao        *mgl.VAO         // 選択頂点VAO
	selectedVertexVbo        *mgl.VBO         // 選択頂点VBO
	selectedVertexIbo        *mgl.IBO         // 選択頂点IBO
	boneLineCount            int              // ボーン情報
	boneLineVao              *mgl.VAO         // ボーンVAO
	boneLineVbo              *mgl.VBO         // ボーンVBO
	boneLineIbo              *mgl.IBO         // ボーンIBO
	boneLineIndexes          []int            // ボーンインデックス
	bonePointCount           int              // ボーン情報
	bonePointVao             *mgl.VAO         // ボーンVAO
	bonePointVbo             *mgl.VBO         // ボーンVBO
	bonePointIbo             *mgl.IBO         // ボーンIBO
	bonePointIndexes         []int            // ボーンインデックス
	cursorPositionVao        *mgl.VAO         // カーソルラインVAO
	cursorPositionVbo        *mgl.VBO         // カーソルラインVBO
	ssbo                     uint32           // SSBO
	vertexCount              int              // 頂点数
}

func NewRenderModel(windowIndex int, model *pmx.PmxModel) *RenderModel {
	m := &RenderModel{
		windowIndex: windowIndex,
		model:       model,
		vertexCount: model.Vertices.Len(),
	}
	m.initToonTexturesGl(windowIndex)
	m.initTexturesGl(windowIndex, model.Textures, model.Path())

	m.initializeBuffer(model)
	m.hash = model.Hash()

	return m
}

func (renderModel *RenderModel) Model() *pmx.PmxModel {
	return renderModel.model
}

func (renderModel *RenderModel) Hash() string {
	return renderModel.hash
}

func (renderModel *RenderModel) Delete() {
	renderModel.vao.Delete()
	renderModel.vbo.Delete()
	renderModel.normalVao.Delete()
	renderModel.normalVbo.Delete()
	renderModel.normalIbo.Delete()
	renderModel.selectedVertexVao.Delete()
	renderModel.selectedVertexVbo.Delete()
	renderModel.selectedVertexIbo.Delete()
	renderModel.boneLineVao.Delete()
	renderModel.boneLineVbo.Delete()
	renderModel.boneLineIbo.Delete()
	gl.DeleteBuffers(1, &renderModel.ssbo)
	for _, mesh := range renderModel.meshes {
		if mesh != nil {
			mesh.delete()
		}
	}
	for _, texture := range renderModel.textures {
		if texture != nil {
			texture.delete()
		}
	}
	for _, texture := range renderModel.toonTextures {
		if texture != nil {
			texture.delete()
		}
	}
}

func (renderModel *RenderModel) initializeBuffer(
	model *pmx.PmxModel,
) {
	// 頂点情報
	renderModel.vertices = make([]float32, 0, len(model.Vertices.Data))
	renderModel.normalVertices = make([]float32, 0, len(model.Vertices.Data)*2)
	normalFaces := make([]uint32, 0, len(model.Vertices.Data)*2)
	selectedVertices := make([]float32, 0, len(model.Vertices.Data))
	selectedVertexFaces := make([]uint32, 0, len(model.Vertices.Data))

	// WaitGroupを用いて並列処理を管理
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 頂点情報の並列処理
	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for i, vertex := range model.Vertices.Data {
			vgl := newVertexGl(vertex)

			mu.Lock()
			renderModel.vertices = append(renderModel.vertices, vgl...)

			// 法線
			renderModel.normalVertices = append(renderModel.normalVertices, vgl...)
			renderModel.normalVertices = append(renderModel.normalVertices, newVertexNormalGl(vertex)...)
			normalFaces = append(normalFaces, uint32(n), uint32(n+1))

			// 選択頂点
			selectedVertices = append(selectedVertices, newSelectedVertexGl(vertex)...)
			selectedVertexFaces = append(selectedVertexFaces, uint32(i))

			mu.Unlock()
			n += 2
		}
	}()

	// 面情報の並列処理
	faces := make([]uint32, 0, len(model.Faces.Data)*3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for _, face := range model.Faces.Data {
			vertices := face.VertexIndexes
			mu.Lock()
			faces = append(faces, uint32(vertices[2]), uint32(vertices[1]), uint32(vertices[0]))
			mu.Unlock()

			n += 3
		}
	}()

	// ボーン情報の並列処理
	bones := make([]float32, 0)
	boneFaces := make([]uint32, 0, len(model.Bones.Data)*2)
	boneIndexes := make([]int, len(model.Bones.Data)*2)

	bonePoints := make([]float32, 0)
	bonePointFaces := make([]uint32, len(model.Bones.Data))
	bonePointIndexes := make([]int, len(model.Bones.Data))

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for _, bone := range model.Bones.Data {
			bones = append(bones, newBoneGl(bone)...)
			bones = append(bones, newTailBoneGl(bone)...)
			boneFaces = append(boneFaces, uint32(n), uint32(n+1))
			boneIndexes[n] = bone.Index()
			boneIndexes[n+1] = bone.Index()

			bonePoints = append(bonePoints, newBoneGl(bone)...)
			bonePointFaces[bone.Index()] = uint32(bone.Index())
			bonePointIndexes[bone.Index()] = bone.Index()

			n += 2
		}
	}()

	// メッシュ情報
	renderModel.meshes = make([]*Mesh, len(model.Materials.Data))
	prevVerticesCount := 0

	for i, m := range model.Materials.Data {
		// テクスチャ
		var texGl *textureGl
		if m.TextureIndex != -1 && model.Textures.Contains(m.TextureIndex) {
			texGl = renderModel.textures[m.TextureIndex]
		}

		var toonTexGl *textureGl
		if m.ToonSharingFlag == pmx.TOON_SHARING_INDIVIDUAL &&
			m.ToonTextureIndex != -1 &&
			model.Textures.Contains(m.ToonTextureIndex) {
			// 個別Toon
			toonTexGl = renderModel.textures[m.ToonTextureIndex]
		} else if m.ToonSharingFlag == pmx.TOON_SHARING_SHARING &&
			m.ToonTextureIndex != -1 {
			// 共有Toon
			toonTexGl = renderModel.toonTextures[m.ToonTextureIndex]
		}

		var sphereTexGl *textureGl
		if m.SphereMode != pmx.SPHERE_MODE_INVALID &&
			m.SphereTextureIndex != -1 &&
			model.Textures.Contains(m.SphereTextureIndex) {
			sphereTexGl = renderModel.textures[m.SphereTextureIndex]
		}

		materialGl := &materialGL{
			Material:          m,
			texture:           texGl,
			sphereTexture:     sphereTexGl,
			toonTexture:       toonTexGl,
			prevVerticesCount: prevVerticesCount,
		}

		mesh := newMesh(
			faces,
			materialGl,
			prevVerticesCount,
		)

		renderModel.meshes[i] = mesh
		prevVerticesCount += m.VerticesCount
	}

	// WaitGroupの完了を待つ
	wg.Wait()

	// GLオブジェクトの生成は並列化しない
	renderModel.vao = mgl.NewVAO()
	renderModel.vao.Bind()

	renderModel.vbo = mgl.NewVBOForVertex(gl.Ptr(renderModel.vertices), len(renderModel.vertices))
	renderModel.vbo.BindVertex(nil)
	renderModel.vbo.Unbind()
	renderModel.vao.Unbind()

	renderModel.normalVao = mgl.NewVAO()
	renderModel.normalVao.Bind()
	renderModel.normalVbo = mgl.NewVBOForVertex(gl.Ptr(renderModel.normalVertices), len(renderModel.normalVertices))
	renderModel.normalVbo.BindVertex(nil)
	renderModel.normalIbo = mgl.NewIBO(gl.Ptr(normalFaces), len(normalFaces))
	renderModel.normalIbo.Bind()
	renderModel.normalIbo.Unbind()
	renderModel.normalVbo.Unbind()
	renderModel.normalVao.Unbind()

	renderModel.boneLineCount = len(bones)
	renderModel.boneLineIndexes = boneIndexes
	renderModel.boneLineVao = mgl.NewVAO()
	renderModel.boneLineVao.Bind()
	renderModel.boneLineVbo = mgl.NewVBOForBone(gl.Ptr(bones), len(bones))
	renderModel.boneLineVbo.BindBone(nil, nil)
	renderModel.boneLineIbo = mgl.NewIBO(gl.Ptr(boneFaces), len(boneFaces))
	renderModel.boneLineIbo.Bind()
	renderModel.boneLineIbo.Unbind()
	renderModel.boneLineVbo.Unbind()
	renderModel.boneLineVao.Unbind()

	renderModel.bonePointCount = len(bonePoints)
	renderModel.bonePointIndexes = bonePointIndexes
	renderModel.bonePointVao = mgl.NewVAO()
	renderModel.bonePointVao.Bind()
	renderModel.bonePointVbo = mgl.NewVBOForBone(gl.Ptr(bonePoints), len(bonePoints))
	renderModel.bonePointVbo.BindBone(nil, nil)
	renderModel.bonePointIbo = mgl.NewIBO(gl.Ptr(bonePointFaces), len(bonePointFaces))
	renderModel.bonePointIbo.Bind()
	renderModel.bonePointIbo.Unbind()
	renderModel.bonePointVbo.Unbind()
	renderModel.bonePointVao.Unbind()

	renderModel.selectedVertexVao = mgl.NewVAO()
	renderModel.selectedVertexVao.Bind()
	renderModel.selectedVertexVbo = mgl.NewVBOForVertex(gl.Ptr(selectedVertices), len(selectedVertices))
	renderModel.selectedVertexVbo.BindVertex(nil)
	renderModel.selectedVertexIbo = mgl.NewIBO(gl.Ptr(selectedVertexFaces), len(selectedVertexFaces))
	renderModel.selectedVertexIbo.Bind()
	renderModel.selectedVertexIbo.Unbind()
	renderModel.selectedVertexVbo.Unbind()
	renderModel.selectedVertexVao.Unbind()

	cursorPositions := []float32{0, 0, 0}
	renderModel.cursorPositionVao = mgl.NewVAO()
	renderModel.cursorPositionVao.Bind()
	renderModel.cursorPositionVbo = mgl.NewVBOForCursorPosition()
	renderModel.cursorPositionVbo.BindCursorPosition(cursorPositions)
	renderModel.cursorPositionVbo.Unbind()
	renderModel.cursorPositionVao.Unbind()

	// SSBOの作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, model.Vertices.Len()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	renderModel.ssbo = ssbo
}

func (renderModel *RenderModel) Render(
	shader mgl.IShader, appState state.IAppState, vmdDeltas *delta.VmdDeltas,
	leftCursorPositions, leftCursorRemovePositions, leftCursorWorldHistoryPositions,
	leftCursorRemoveWorldHistoryPositions []*mgl32.Vec3,
) {
	renderModel.vao.Bind()
	defer renderModel.vao.Unbind()

	renderModel.vbo.BindVertex(vmdDeltas.Morphs.Vertices)
	defer renderModel.vbo.Unbind()

	paddedMatrixes, matrixWidth, matrixHeight, err := createBoneMatrixes(vmdDeltas.Bones)
	if err != nil {
		return
	}

	for i, mesh := range renderModel.meshes {
		mesh.ibo.Bind()

		md := delta.NewMeshDelta(vmdDeltas.Morphs.Materials.Get(i))
		mesh.drawModel(renderModel.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, md)

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			mesh.drawEdge(renderModel.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, md)
		}

		if appState.IsShowWire() && !renderModel.ExistInvisibleMaterialIndex(mesh.material.Index()) {
			mesh.drawWire(renderModel.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight, false)
		}

		mesh.ibo.Unbind()
	}

	if appState.IsShowNormal() {
		renderModel.drawNormal(renderModel.windowIndex, shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	if appState.IsShowBoneAll() || appState.IsShowBoneEffector() || appState.IsShowBoneIk() ||
		appState.IsShowBoneFixed() || appState.IsShowBoneRotate() ||
		appState.IsShowBoneTranslate() || appState.IsShowBoneVisible() {
		renderModel.drawBone(renderModel.windowIndex,
			shader, renderModel.model.Bones, appState,
			paddedMatrixes, matrixWidth, matrixHeight)
	}

	if appState.IsShowSelectedVertex() {
		var cursorPositions []float32
		var removeCursorPositions []float32

		if len(leftCursorPositions) > 0 {
			cursorPositions = make([]float32, len(leftCursorPositions)*3)
			for i, pos := range leftCursorPositions {
				cursorPositions[i*3] = pos[0]
				cursorPositions[i*3+1] = pos[1]
				cursorPositions[i*3+2] = pos[2]
			}
		}

		if len(leftCursorRemovePositions) > 0 {
			removeCursorPositions = make([]float32, len(leftCursorRemovePositions)*3)
			for i, pos := range leftCursorRemovePositions {
				removeCursorPositions[i*3] = pos[0]
				removeCursorPositions[i*3+1] = pos[1]
				removeCursorPositions[i*3+2] = pos[2]
			}
		}

		for i := len(leftCursorPositions); i < 30; i++ {
			// 追加固定で30個までのカーソル位置を確保
			cursorPositions = append(cursorPositions, 0, 0, 0)
		}

		if len(leftCursorRemovePositions) > 0 {
			for i := len(leftCursorRemovePositions); i < 30; i++ {
				// 追加固定で30個までのカーソル位置を確保
				removeCursorPositions = append(removeCursorPositions, 0, 0, 0)
			}
		}

		selectedVertexIndexes := renderModel.drawSelectedVertex(
			renderModel.windowIndex, renderModel.InvisibleMaterialIndexes(),
			renderModel.SelectedVertexIndexes(), renderModel.NoSelectedVertexIndexes(),
			shader, paddedMatrixes, matrixWidth, matrixHeight, cursorPositions, removeCursorPositions)
		if len(removeCursorPositions) > 0 {
			renderModel.UpdateNoSelectedVertexIndexes(selectedVertexIndexes)
		} else if len(leftCursorPositions) > 0 {
			renderModel.UpdateSelectedVertexIndexes(selectedVertexIndexes)
		}

		// カーソルの軌跡
		if len(leftCursorWorldHistoryPositions) > 1 {
			cursorPositions = make([]float32, len(leftCursorWorldHistoryPositions)*3)
			for i, pos := range leftCursorWorldHistoryPositions {
				cursorPositions[i*3] = pos[0]
				cursorPositions[i*3+1] = pos[1]
				cursorPositions[i*3+2] = pos[2]
			}
			renderModel.drawCursorLine(shader, cursorPositions, mgl32.Vec4{1, 0, 0, 1})
		}
		if len(leftCursorRemoveWorldHistoryPositions) > 1 {
			removeCursorPositions = make([]float32, len(leftCursorRemoveWorldHistoryPositions)*3)
			for i, pos := range leftCursorRemoveWorldHistoryPositions {
				removeCursorPositions[i*3] = pos[0]
				removeCursorPositions[i*3+1] = pos[1]
				removeCursorPositions[i*3+2] = pos[2]
			}
			renderModel.drawCursorLine(shader, removeCursorPositions, mgl32.Vec4{0, 0, 1, 1})
		}
	}
}

func (renderModel *RenderModel) drawNormal(
	windowIndex int,
	shader mgl.IShader,
	paddedMatrixes []float32,
	width, height int,
) {
	program := shader.Program(mgl.PROGRAM_TYPE_NORMAL)
	gl.UseProgram(program)

	renderModel.normalVao.Bind()
	renderModel.normalVbo.BindVertex(nil)
	renderModel.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	normalColor := mgl32.Vec4{0.3, 0.3, 0.7, 0.5}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &normalColor[0])

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(renderModel.normalVertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.normalIbo.Unbind()
	renderModel.normalVbo.Unbind()
	renderModel.normalVao.Unbind()

	gl.UseProgram(0)
}

func (renderModel *RenderModel) drawCursorLine(
	shader mgl.IShader, cursorPositions []float32, vertexColor mgl32.Vec4,
) {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(mgl.PROGRAM_TYPE_CURSOR)
	gl.UseProgram(program)

	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(colorUniform, 1, &vertexColor[0])

	renderModel.cursorPositionVao.Bind()
	renderModel.cursorPositionVbo.BindCursorPosition(cursorPositions)

	// ライン描画
	gl.DrawArrays(gl.LINES, 0, int32(len(cursorPositions)/3))

	renderModel.cursorPositionVbo.Unbind()
	renderModel.cursorPositionVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

func (renderModel *RenderModel) drawSelectedVertex(
	windowIndex int,
	invisibleMaterialIndexes []int,
	nowSelectedVertexIndexes []int,
	nowNoSelectedVertexIndexes []int,
	shader mgl.IShader,
	paddedMatrixes []float32, width, height int,
	cursorPositions []float32,
	removeCursorPositions []float32,
) []int {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(mgl.PROGRAM_TYPE_SELECTED_VERTEX)
	gl.UseProgram(program)

	selectedVertexDeltas := delta.NewVertexMorphDeltas()
	for _, index := range nowNoSelectedVertexIndexes {
		// 選択されていない頂点のUVXを-1にして表示しない
		vd := delta.NewVertexMorphDelta(index)
		vd.Uv = &mmath.MVec2{X: -1, Y: 0}
		selectedVertexDeltas.Update(vd)
	}
	for _, index := range nowSelectedVertexIndexes {
		// 選択されている頂点のUVXを＋にして（フラグをたてて）表示する
		vd := delta.NewVertexMorphDelta(index)
		vd.Uv = &mmath.MVec2{X: 1, Y: 0}
		selectedVertexDeltas.Update(vd)
	}

	renderModel.selectedVertexVao.Bind()
	renderModel.selectedVertexVbo.BindVertex(selectedVertexDeltas)
	renderModel.selectedVertexIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	vertexColor := mgl32.Vec4{1.0, 0.4, 0.0, 0.7}
	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(colorUniform, 1, &vertexColor[0])
	gl.PointSize(5.0) // 選択頂点のサイズ

	// 閾値に視野角を追加
	thresholdUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CURSOR_THRESHOLD))
	gl.Uniform1f(thresholdUniform, shader.ViewOfAngle())

	cursorPositionsUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_CURSOR_POSITIONS))
	if removeCursorPositions != nil {
		gl.Uniform3fv(cursorPositionsUniform, int32(len(removeCursorPositions)/3), &removeCursorPositions[0])
	} else {
		gl.Uniform3fv(cursorPositionsUniform, int32(len(cursorPositions)/3), &cursorPositions[0])
	}

	// 点描画
	gl.DrawElements(
		gl.POINTS,
		int32(len(renderModel.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	// SSBOからデータを読み込む
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, renderModel.ssbo)
	ptr := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)
	selectedVertexPositions := make(map[int]mgl32.Vec4)
	selectedVertexIndexes := make([]int, 0, renderModel.vertexCount)

	if ptr != nil {
		// SSBOから読み取り
		for i := range renderModel.vertexCount {
			vertexPositions := make([]float32, 4)
			for j := range 4 {
				vertexPositions[j] = *(*float32)(unsafe.Pointer(uintptr(ptr) + uintptr((i*4+j)*4)))
			}
			isInvisible := false
			for _, materialIndex := range renderModel.model.Vertices.Get(i).MaterialIndexes {
				if slices.Contains(invisibleMaterialIndexes, materialIndex) {
					isInvisible = true
					break
				}
			}
			if !isInvisible && (vertexPositions[0] != -1 || vertexPositions[1] != -1 ||
				vertexPositions[2] != -1 || vertexPositions[3] != -1) {
				// 表示対象でいずれが-1でない場合はマップに保持
				selectedVertexIndexes = append(selectedVertexIndexes, i)
				selectedVertexPositions[i] = mgl32.Vec4{
					vertexPositions[0], vertexPositions[1], vertexPositions[2], vertexPositions[3]}
			}
		}
	}

	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	renderModel.selectedVertexIbo.Unbind()
	renderModel.selectedVertexVbo.Unbind()
	renderModel.selectedVertexVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	if mlog.IsDebug() {
		if cursorPositions[0] != 0 {
			mlog.L()
			for i := 0; i < 90; i += 3 {
				if cursorPositions[i] == 0 {
					break
				}
				mlog.D("CURSOR index: %d, position: %v", i/3, cursorPositions[i:i+3])
			}
			for i, position := range selectedVertexPositions {
				mlog.D("VERTEX index: %d, position: %v", i, position)
			}
		}
	}

	return selectedVertexIndexes
}

func (renderModel *RenderModel) drawBone(
	windowIndex int,
	shader mgl.IShader,
	bones *pmx.Bones,
	appState state.IAppState,
	paddedMatrixes []float32,
	width, height int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.Program(mgl.PROGRAM_TYPE_BONE)
	gl.UseProgram(program)

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	renderModel.boneLineVao.Bind()
	renderModel.boneLineVbo.BindBone(renderModel.fetchBoneLineDeltas(bones, appState))
	renderModel.boneLineIbo.Bind()

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(renderModel.boneLineCount),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.boneLineIbo.Unbind()
	renderModel.boneLineVbo.Unbind()
	renderModel.boneLineVao.Unbind()

	renderModel.bonePointVao.Bind()
	renderModel.bonePointVbo.BindBone(renderModel.fetchBonePointDeltas(bones, appState))
	renderModel.bonePointIbo.Bind()

	gl.PointSize(5.0) // 選択頂点のサイズ

	// 点描画
	gl.DrawElements(
		gl.POINTS,
		int32(len(renderModel.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.bonePointIbo.Unbind()
	renderModel.bonePointVbo.Unbind()
	renderModel.bonePointVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

}

func (renderModel *RenderModel) fetchBoneLineDeltas(
	bones *pmx.Bones, appState state.IAppState,
) ([]int, [][]float32) {
	indexes := make([]int, len(renderModel.boneLineIndexes))
	deltas := make([][]float32, len(renderModel.boneLineIndexes))

	for i, boneIndex := range renderModel.boneLineIndexes {
		indexes[i] = i
		deltas[i] = getBoneDebugColor(bones.Get(boneIndex), appState)
	}

	return indexes, deltas
}

func (renderModel *RenderModel) fetchBonePointDeltas(
	bones *pmx.Bones, appState state.IAppState,
) ([]int, [][]float32) {
	indexes := make([]int, len(renderModel.bonePointIndexes))
	deltas := make([][]float32, len(renderModel.bonePointIndexes))

	for i, boneIndex := range renderModel.bonePointIndexes {
		indexes[i] = i
		deltas[i] = getBoneDebugColor(bones.Get(boneIndex), appState)
	}

	return indexes, deltas
}

func (renderModel *RenderModel) InvisibleMaterialIndexes() []int {
	if renderModel.invisibleMaterialIndexes == nil {
		return nil
	}

	indexes := make([]int, 0, len(renderModel.invisibleMaterialIndexes))
	for i := range renderModel.invisibleMaterialIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (renderModel *RenderModel) SetInvisibleMaterialIndexes(indexes []int) {
	if len(indexes) == 0 {
		renderModel.invisibleMaterialIndexes = nil
		return
	}
	renderModel.invisibleMaterialIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		renderModel.invisibleMaterialIndexes[i] = struct{}{}
	}
}

func (renderModel *RenderModel) ExistInvisibleMaterialIndex(index int) bool {
	_, ok := renderModel.invisibleMaterialIndexes[index]
	return ok
}

func (renderModel *RenderModel) SelectedVertexIndexes() []int {
	indexes := make([]int, 0, len(renderModel.selectedVertexIndexes))
	for i := range renderModel.selectedVertexIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (renderModel *RenderModel) NoSelectedVertexIndexes() []int {
	indexes := make([]int, 0, len(renderModel.noSelectedVertexIndexes))
	for i := range renderModel.noSelectedVertexIndexes {
		indexes = append(indexes, i)
	}
	return indexes
}

func (renderModel *RenderModel) SetSelectedVertexIndexes(indexes []int) {
	renderModel.selectedVertexIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		renderModel.selectedVertexIndexes[i] = struct{}{}
	}
}

func (renderModel *RenderModel) SetNoSelectedVertexIndexes(indexes []int) {
	renderModel.noSelectedVertexIndexes = make(map[int]struct{}, len(indexes))
	for _, i := range indexes {
		renderModel.noSelectedVertexIndexes[i] = struct{}{}
	}
}

func (renderModel *RenderModel) ClearSelectedVertexIndexes() {
	renderModel.UpdateNoSelectedVertexIndexes(renderModel.SelectedVertexIndexes())
	renderModel.selectedVertexIndexes = make(map[int]struct{})
}

func (renderModel *RenderModel) UpdateSelectedVertexIndexes(indexes []int) {
	for _, index := range indexes {
		renderModel.selectedVertexIndexes[index] = struct{}{}
	}
}

func (renderModel *RenderModel) UpdateNoSelectedVertexIndexes(indexes []int) {
	if renderModel.noSelectedVertexIndexes == nil {
		renderModel.noSelectedVertexIndexes = make(map[int]struct{}, len(indexes))
	}
	for _, index := range indexes {
		renderModel.noSelectedVertexIndexes[index] = struct{}{}
		delete(renderModel.selectedVertexIndexes, index)
	}
}
