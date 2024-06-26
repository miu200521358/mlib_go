//go:build windows
// +build windows

package pmx

import (
	"embed"
	"math"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mview"
)

type Meshes struct {
	meshes      []*Mesh
	vertices    []float32
	vao         *mview.VAO
	vbo         *mview.VBO
	normals     []float32
	normalVao   *mview.VAO
	normalVbo   *mview.VBO
	normalIbo   *mview.IBO
	bones       []float32
	boneVao     *mview.VAO
	boneVbo     *mview.VBO
	boneIbo     *mview.IBO
	boneIndexes []int
}

func NewMeshes(
	model *PmxModel,
	windowIndex int,
	resourceFiles embed.FS,
) *Meshes {
	// 頂点情報
	vertices := make([]float32, 0, len(model.Vertices.Data))
	normalVertices := make([]float32, 0, len(model.Vertices.Data)*2)
	normalFaces := make([]uint32, 0, len(model.Vertices.Data)*2)

	n := 0
	for i := range len(model.Vertices.Data) {
		vertex := model.Vertices.Get(i)
		vertices = append(vertices, vertex.GL()...)

		normalVertices = append(normalVertices, vertex.GL()...)
		normalVertices = append(normalVertices, vertex.NormalGL()...)

		normalFaces = append(normalFaces, uint32(n), uint32(n+1))
		n += 2
	}

	// 面情報
	faces := make([]uint32, 0, len(model.Faces.Data)*3)
	for i := range len(model.Faces.Data) {
		vertices := model.Faces.Get(i).VertexIndexes
		faces = append(faces, uint32(vertices[2]), uint32(vertices[1]), uint32(vertices[0]))
	}

	meshes := make([]*Mesh, len(model.Materials.GetIndexes()))
	prevVerticesCount := 0
	for i := range model.Materials.Len() {
		m := model.Materials.Get(i)

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
			model.Path,
			texture,
			toonTexture,
			sphereTexture,
			windowIndex,
			prevVerticesCount,
			resourceFiles,
		)
		mesh := NewMesh(
			faces,
			materialGl,
			prevVerticesCount,
		)
		meshes[i] = mesh

		prevVerticesCount += m.VerticesCount
	}

	bones := make([]float32, 0, len(model.Bones.Data)*4)
	boneFaces := make([]uint32, 0, len(model.Bones.Data)*4)
	boneIndexes := make([]int, 0, len(model.Bones.Data)*4)

	n = 0
	for i := range len(model.Bones.Data) {
		bone := model.Bones.Get(i)
		bones = append(bones, bone.GL()...)
		bones = append(bones, bone.TailGL()...)

		boneFaces = append(boneFaces, uint32(n), uint32(n+1))
		boneIndexes = append(boneIndexes, bone.Index, bone.Index)

		n += 2

		if bone.ParentIndex >= 0 && model.Bones.Contains(bone.ParentIndex) &&
			!model.Bones.Get(bone.ParentIndex).Position.IsZero() {
			bones = append(bones, bone.GL()...)
			bones = append(bones, bone.ParentGL()...)

			boneFaces = append(boneFaces, uint32(n), uint32(n+1))
			boneIndexes = append(boneIndexes, bone.Index, bone.ParentIndex)

			n += 2
		}
	}

	vao := mview.NewVAO()
	vao.Bind()

	vbo := mview.NewVBOForVertex(gl.Ptr(vertices), len(vertices))
	vbo.BindVertex(nil, nil)
	vbo.Unbind()
	vao.Unbind()

	normalVao := mview.NewVAO()
	normalVao.Bind()
	normalVbo := mview.NewVBOForVertex(gl.Ptr(normalVertices), len(normalVertices))
	normalVbo.BindVertex(nil, nil)
	normalIbo := mview.NewIBO(gl.Ptr(normalFaces), len(normalFaces))
	normalIbo.Bind()
	normalIbo.Unbind()
	normalVbo.Unbind()
	normalVao.Unbind()

	boneVao := mview.NewVAO()
	boneVao.Bind()
	boneVbo := mview.NewVBOForVertex(gl.Ptr(bones), len(bones))
	boneVbo.BindVertex(nil, nil)
	boneIbo := mview.NewIBO(gl.Ptr(boneFaces), len(boneFaces))
	boneIbo.Bind()
	boneIbo.Unbind()
	boneVbo.Unbind()
	boneVao.Unbind()

	return &Meshes{
		meshes:      meshes,
		vertices:    vertices,
		vao:         vao,
		vbo:         vbo,
		normals:     normalVertices,
		normalVao:   normalVao,
		normalVbo:   normalVbo,
		normalIbo:   normalIbo,
		bones:       bones,
		boneVao:     boneVao,
		boneVbo:     boneVbo,
		boneIbo:     boneIbo,
		boneIndexes: boneIndexes,
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
	shader *mview.MShader,
	boneDeltas []mgl32.Mat4,
	vertexDeltas [][]float32,
	meshDeltas []*MeshDelta,
	windowIndex int,
	isDrawNormal bool,
	isDrawBones map[BoneFlag]bool,
	bones *Bones,
) {
	m.vao.Bind()
	defer m.vao.Unbind()

	m.vbo.BindVertex(m.vertices, vertexDeltas)
	defer m.vbo.Unbind()

	paddedMatrixes, matrixWidth, matrixHeight := m.createBoneMatrixes(boneDeltas)

	for i, mesh := range m.meshes {
		mesh.ibo.Bind()

		shader.Use(mview.PROGRAM_TYPE_MODEL)
		mesh.drawModel(shader, windowIndex, paddedMatrixes, matrixWidth, matrixHeight, meshDeltas[i])
		shader.Unuse()

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			shader.Use(mview.PROGRAM_TYPE_EDGE)
			mesh.drawEdge(shader, windowIndex, paddedMatrixes, matrixWidth, matrixHeight, meshDeltas[i])
			shader.Unuse()
		}

		mesh.ibo.Unbind()
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
	vertexDeltas = nil
	meshDeltas = nil
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
	shader *mview.MShader,
	paddedMatrixes []float32,
	width, height int,
	windowIndex int,
) {
	shader.Use(mview.PROGRAM_TYPE_NORMAL)

	m.normalVao.Bind()
	m.normalVbo.BindVertex(nil, nil)
	m.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, shader.NormalProgram)

	normalColor := mgl32.Vec4{0.3, 0.3, 0.7, 0.5}
	specularUniform := gl.GetUniformLocation(shader.NormalProgram, gl.Str(mview.SHADER_COLOR))
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

func (m *Meshes) drawBone(
	shader *mview.MShader,
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

	shader.Use(mview.PROGRAM_TYPE_BONE)

	m.boneVao.Bind()
	m.boneVbo.BindVertex(m.bones, m.fetchBoneDebugDeltas(bones, isDrawBones))
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
}

func (m *Meshes) fetchBoneDebugDeltas(bones *Bones, isDrawBones map[BoneFlag]bool) [][]float32 {
	vertexDeltas := make([][]float32, len(m.bones)/m.boneVbo.StrideSize)

	for i, boneIndex := range m.boneIndexes {
		bone := bones.Get(boneIndex)
		vertexDeltas[i] = bone.DeltaGL(isDrawBones)
	}

	return vertexDeltas
}
