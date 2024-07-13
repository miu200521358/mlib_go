//go:build windows
// +build windows

package pmx

import (
	"math"
	"slices"
	"sync"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/buffer"
	"github.com/miu200521358/mlib_go/pkg/infra/mgl"
)

type Meshes struct {
	meshes            []*Mesh
	vertices          []float32
	vao               *buffer.VAO
	vbo               *buffer.VBO
	normals           []float32
	normalVao         *buffer.VAO
	normalVbo         *buffer.VBO
	normalIbo         *buffer.IBO
	selectedVertexVao *buffer.VAO
	selectedVertexVbo *buffer.VBO
	selectedVertexIbo *buffer.IBO
	bones             []float32
	boneVao           *buffer.VAO
	boneVbo           *buffer.VBO
	boneIbo           *buffer.IBO
	boneIndexes       []int
	ssbo              uint32
	vertexCount       int
}

func NewMeshes(
	model *PmxModel,
	windowIndex int,
) *Meshes {
	// 頂点情報
	vertices := make([]float32, 0, len(model.Vertices.Data))
	normalVertices := make([]float32, 0, len(model.Vertices.Data)*2)
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
			vgl := vertex.GL()

			mu.Lock()
			vertices = append(vertices, vgl...)

			// 法線
			normalVertices = append(normalVertices, vgl...)
			normalVertices = append(normalVertices, vertex.NormalGL()...)
			normalFaces = append(normalFaces, uint32(n), uint32(n+1))

			// 選択頂点
			selectedVertices = append(selectedVertices, vertex.SelectedGL()...)
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

	// メッシュ情報の並列処理
	meshes := make([]*Mesh, len(model.Materials.Data))
	prevVerticesCount := 0

	// テクスチャの gl.GenTextures はスレッドセーフではないので、並列化しない
	for i, m := range model.Materials.Data {
		// テクスチャ
		var texture *Texture
		if m.TextureIndex != -1 && model.Textures.Contains(m.TextureIndex) {
			texture = model.Textures.Get(m.TextureIndex)
		}

		var toonTexture *Texture
		// 個別Toon
		if m.ToonSharingFlag == TOON_SHARING_INDIVIDUAL &&
			m.ToonTextureIndex != -1 &&
			model.Textures.Contains(m.ToonTextureIndex) {
			toonTexture = model.Textures.Get(m.ToonTextureIndex)
		}
		// 共有Toon
		if m.ToonSharingFlag == TOON_SHARING_SHARING &&
			m.ToonTextureIndex != -1 &&
			model.ToonTextures.Contains(m.ToonTextureIndex) {
			toonTexture = model.ToonTextures.Get(m.ToonTextureIndex)
		}

		var sphereTexture *Texture
		if m.SphereMode != SPHERE_MODE_INVALID &&
			m.SphereTextureIndex != -1 &&
			model.Textures.Contains(m.SphereTextureIndex) {
			sphereTexture = model.Textures.Get(m.SphereTextureIndex)
		}

		materialGl := m.GL(
			model.GetPath(),
			texture,
			toonTexture,
			sphereTexture,
			windowIndex,
			prevVerticesCount,
		)
		mesh := NewMesh(
			faces,
			materialGl,
			prevVerticesCount,
		)
		mu.Lock()
		meshes[i] = mesh
		mu.Unlock()

		prevVerticesCount += m.VerticesCount
	}

	// ボーン情報の並列処理
	bones := make([]float32, 0, len(model.Bones.Data)*4)
	boneFaces := make([]uint32, 0, len(model.Bones.Data)*4)
	boneIndexes := make([]int, 0, len(model.Bones.Data)*4)

	wg.Add(1)
	go func() {
		defer wg.Done()
		n := 0
		for _, bone := range model.Bones.Data {
			mu.Lock()
			bones = append(bones, bone.GL()...)
			bones = append(bones, bone.TailGL()...)
			boneFaces = append(boneFaces, uint32(n), uint32(n+1))
			boneIndexes = append(boneIndexes, bone.Index, bone.Index)
			mu.Unlock()

			n += 2

			if bone.ParentIndex >= 0 && model.Bones.Contains(bone.ParentIndex) &&
				!model.Bones.Get(bone.ParentIndex).Position.IsZero() {
				mu.Lock()
				bones = append(bones, bone.GL()...)
				bones = append(bones, bone.ParentGL()...)
				boneFaces = append(boneFaces, uint32(n), uint32(n+1))
				boneIndexes = append(boneIndexes, bone.Index, bone.ParentIndex)
				mu.Unlock()

				n += 2
			}
		}
	}()

	// WaitGroupの完了を待つ
	wg.Wait()

	// 以下の部分は並列化する必要がないのでそのままにする
	vao := buffer.NewVAO()
	vao.Bind()

	vbo := buffer.NewVBOForVertex(gl.Ptr(vertices), len(vertices))
	vbo.BindVertex(nil, nil)
	vbo.Unbind()
	vao.Unbind()

	normalVao := buffer.NewVAO()
	normalVao.Bind()
	normalVbo := buffer.NewVBOForVertex(gl.Ptr(normalVertices), len(normalVertices))
	normalVbo.BindVertex(nil, nil)
	normalIbo := buffer.NewIBO(gl.Ptr(normalFaces), len(normalFaces))
	normalIbo.Bind()
	normalIbo.Unbind()
	normalVbo.Unbind()
	normalVao.Unbind()

	boneVao := buffer.NewVAO()
	boneVao.Bind()
	boneVbo := buffer.NewVBOForVertex(gl.Ptr(bones), len(bones))
	boneVbo.BindVertex(nil, nil)
	boneIbo := buffer.NewIBO(gl.Ptr(boneFaces), len(boneFaces))
	boneIbo.Bind()
	boneIbo.Unbind()
	boneVbo.Unbind()
	boneVao.Unbind()

	selectedVertexVao := buffer.NewVAO()
	selectedVertexVao.Bind()
	selectedVertexVbo := buffer.NewVBOForVertex(gl.Ptr(selectedVertices), len(selectedVertices))
	selectedVertexVbo.BindVertex(nil, nil)
	selectedVertexIbo := buffer.NewIBO(gl.Ptr(selectedVertexFaces), len(selectedVertexFaces))
	selectedVertexIbo.Bind()
	selectedVertexIbo.Unbind()
	selectedVertexVbo.Unbind()
	selectedVertexVao.Unbind()

	// SSBOの作成
	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, model.Vertices.Len()*4*4, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, ssbo)

	return &Meshes{
		meshes:            meshes,
		vertices:          vertices,
		vao:               vao,
		vbo:               vbo,
		normals:           normalVertices,
		normalVao:         normalVao,
		normalVbo:         normalVbo,
		normalIbo:         normalIbo,
		selectedVertexVao: selectedVertexVao,
		selectedVertexVbo: selectedVertexVbo,
		selectedVertexIbo: selectedVertexIbo,
		bones:             bones,
		boneVao:           boneVao,
		boneVbo:           boneVbo,
		boneIbo:           boneIbo,
		boneIndexes:       boneIndexes,
		ssbo:              ssbo,
		vertexCount:       model.Vertices.Len(),
	}
}

func (m *Meshes) delete() {
	for _, mesh := range m.meshes {
		mesh.delete()
	}
	m.vao.Delete()
	m.vbo.Delete()
}

func (m *Meshes) Draw(
	shader *mgl.MShader, boneDeltas []mgl32.Mat4,
	vertexMorphIndexes []int, vertexMorphDeltas [][]float32,
	selectedVertexIndexes []int, selectedVertexDeltas [][]float32, meshDeltas []*MeshDelta,
	invisibleMaterialIndexes []int, nextInvisibleMaterialIndexes []int, windowIndex int,
	isDrawNormal, isDrawWire, isDrawSelectedVertex bool, isDrawBones map[BoneFlag]bool, bones *Bones,
) [][]float32 {
	m.vao.Bind()
	defer m.vao.Unbind()

	m.vbo.BindVertex(vertexMorphIndexes, vertexMorphDeltas)
	defer m.vbo.Unbind()

	paddedMatrixes, matrixWidth, matrixHeight := m.createBoneMatrixes(boneDeltas)

	for i, mesh := range m.meshes {
		mesh.ibo.Bind()

		shader.Use(mgl.PROGRAM_TYPE_MODEL)
		mesh.drawModel(shader, windowIndex, paddedMatrixes, matrixWidth, matrixHeight, meshDeltas[i])
		shader.Unuse()

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			shader.Use(mgl.PROGRAM_TYPE_EDGE)
			mesh.drawEdge(shader, windowIndex, paddedMatrixes, matrixWidth, matrixHeight, meshDeltas[i])
			shader.Unuse()
		}

		if isDrawWire {
			shader.Use(mgl.PROGRAM_TYPE_WIRE)
			mesh.drawWire(shader, windowIndex, paddedMatrixes, matrixWidth, matrixHeight, ((len(invisibleMaterialIndexes) > 0 && len(nextInvisibleMaterialIndexes) == 0 &&
				slices.Contains(invisibleMaterialIndexes, mesh.material.Index)) ||
				slices.Contains(nextInvisibleMaterialIndexes, mesh.material.Index)))
			shader.Unuse()
		}

		mesh.ibo.Unbind()
	}

	// if isDrawWire {
	// 	m.drawWire(wireVertexDeltas, shader, paddedMatrixes, matrixWidth, matrixHeight, windowIndex)
	// }

	vertexPositions := make([][]float32, 0)
	if isDrawSelectedVertex {
		vertexPositions = m.drawSelectedVertex(selectedVertexIndexes, selectedVertexDeltas,
			shader, paddedMatrixes, matrixWidth, matrixHeight, windowIndex)
	}

	if isDrawNormal {
		m.drawNormal(shader, paddedMatrixes, matrixWidth, matrixHeight, windowIndex)
	}

	isDrawBone := false
	for _, drawBone := range isDrawBones {
		if drawBone {
			isDrawBone = true
			break
		}
	}

	if isDrawBone {
		m.drawBone(shader, bones, isDrawBones, paddedMatrixes, matrixWidth, matrixHeight, windowIndex)
	}

	paddedMatrixes = nil
	boneDeltas = nil
	meshDeltas = nil

	return vertexPositions
}

func (m *Meshes) createBoneMatrixes(matrixes []mgl32.Mat4) ([]float32, int, int) {
	// テクスチャのサイズを計算する
	numBones := len(matrixes)
	texSize := int(math.Ceil(math.Sqrt(float64(numBones))))
	width := int(math.Ceil(float64(texSize)/4) * 4 * 4)
	height := int(math.Ceil((float64(numBones) * 4) / float64(width)))

	paddedMatrixes := make([]float32, height*width*4)
	for i, matrix := range matrixes {
		copy(paddedMatrixes[i*16:], matrix[:])
	}

	return paddedMatrixes, width, height
}

func (m *Meshes) drawNormal(
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
	windowIndex int,
) {
	shader.Use(mgl.PROGRAM_TYPE_NORMAL)

	m.normalVao.Bind()
	m.normalVbo.BindVertex(nil, nil)
	m.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, shader.NormalProgram)

	normalColor := mgl32.Vec4{0.3, 0.3, 0.7, 0.5}
	specularUniform := gl.GetUniformLocation(shader.NormalProgram, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &normalColor[0])

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(m.normals)),
		gl.UNSIGNED_INT,
		nil,
	)

	m.normalIbo.Unbind()
	m.normalVbo.Unbind()
	m.normalVao.Unbind()

	shader.Unuse()
}

func (m *Meshes) drawSelectedVertex(
	selectedVertexMorphIndexes []int,
	selectedVertexDeltas [][]float32,
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
	windowIndex int,
) [][]float32 {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	shader.Use(mgl.PROGRAM_TYPE_SELECTED_VERTEX)

	m.selectedVertexVao.Bind()
	m.selectedVertexVbo.BindVertex(selectedVertexMorphIndexes, selectedVertexDeltas)
	m.selectedVertexIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, shader.SelectedVertexProgram)

	vertexColor := mgl32.Vec4{1.0, 0.4, 0.0, 0.7}
	specularUniform := gl.GetUniformLocation(shader.SelectedVertexProgram, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &vertexColor[0])
	gl.PointSize(5.0) // 選択頂点のサイズ

	// 点描画
	gl.DrawElements(
		gl.POINTS,
		int32(len(m.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	// SSBOからデータを読み込む
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, m.ssbo)
	ptr := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)
	results := make([][]float32, m.vertexCount)

	if ptr != nil {
		// 頂点数を取得
		for i := range results {
			results[i] = make([]float32, 4)
		}

		// SSBOから読み取り
		for i := 0; i < m.vertexCount; i++ {
			for j := 0; j < 4; j++ {
				results[i][j] = *(*float32)(unsafe.Pointer(uintptr(ptr) + uintptr((i*4+j)*4)))
			}
		}
	}

	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	m.selectedVertexIbo.Unbind()
	m.selectedVertexVbo.Unbind()
	m.selectedVertexVao.Unbind()

	shader.Unuse()

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	return results
}

func (m *Meshes) drawBone(
	shader *mgl.MShader,
	bones *Bones,
	isDrawBones map[BoneFlag]bool,
	paddedMatrixes []float32,
	width, height int,
	windowIndex int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	shader.Use(mgl.PROGRAM_TYPE_BONE)

	m.boneVao.Bind()
	m.boneVbo.BindVertex(m.fetchBoneDebugDeltas(bones, isDrawBones))
	m.boneIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, shader.BoneProgram)

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(m.bones)),
		gl.UNSIGNED_INT,
		nil,
	)

	m.boneIbo.Unbind()
	m.boneVbo.Unbind()
	m.boneVao.Unbind()

	shader.Unuse()

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

}

func (m *Meshes) fetchBoneDebugDeltas(bones *Bones, isDrawBones map[BoneFlag]bool) ([]int, [][]float32) {
	indexes := make([]int, 0)
	deltas := make([][]float32, 0)

	for _, boneIndex := range m.boneIndexes {
		bone := bones.Get(boneIndex)
		indexes = append(indexes, boneIndex)
		deltas = append(deltas, bone.DeltaGL(isDrawBones))
	}

	return indexes, deltas
}
