//go:build windows
// +build windows

package pmx

import (
	"embed"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mview"
)

type Meshes struct {
	meshes    []*Mesh
	vertices  []float32
	vao       *mview.VAO
	vbo       *mview.VBO
	normals   []float32
	normalVao *mview.VAO
	normalVbo *mview.VBO
	normalIbo *mview.IBO
	bones     []float32
	boneVao   *mview.VAO
	boneVbo   *mview.VBO
	boneIbo   *mview.IBO
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
		vertices = append(vertices, model.Vertices.Get(i).GL()...)

		normalVertices = append(normalVertices, model.Vertices.Get(i).GL()...)
		normalVertices = append(normalVertices, model.Vertices.Get(i).NormalGL()...)

		normalFaces = append(normalFaces, uint32(n))
		normalFaces = append(normalFaces, uint32(n+1))
		n += 2
	}

	// 面情報
	faces := make([]uint32, 0, len(model.Faces.Data)*3)
	for i := range len(model.Faces.Data) {
		vertices := model.Faces.Get(i).VertexIndexes
		faces = append(faces, uint32(vertices[2]))
		faces = append(faces, uint32(vertices[1]))
		faces = append(faces, uint32(vertices[0]))
	}

	meshes := make([]*Mesh, len(model.Materials.GetIndexes()))
	prevVerticesCount := 0
	for i, m := range model.Materials.GetSortedData() {
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

	bones := make([]float32, 0, len(model.Bones.Data))
	boneFaces := make([]uint32, 0, len(model.Bones.Data)*2)

	for i := range len(model.Bones.Data) {
		bone := model.Bones.Get(i)
		bones = append(bones, bone.GL()...)

		boneFaces = append(boneFaces, uint32(bone.Index))
		boneFaces = append(boneFaces, uint32(bone.ParentIndex))
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
		meshes:    meshes,
		vertices:  vertices,
		vao:       vao,
		vbo:       vbo,
		normals:   normalVertices,
		normalVao: normalVao,
		normalVbo: normalVbo,
		normalIbo: normalIbo,
		bones:     bones,
		boneVao:   boneVao,
		boneVbo:   boneVbo,
		boneIbo:   boneIbo,
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
	boneDeltas []*mgl32.Mat4,
	vertexDeltas [][]float32,
	materialDeltas []*Material,
	windowIndex int,
	isDrawNormal bool,
	isDrawBone bool,
) {
	// 隠面消去
	// https://learnopengl.com/Advanced-OpenGL/Depth-testing
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	shader.Msaa.Bind()
	m.vao.Bind()
	m.vbo.BindVertex(m.vertices, vertexDeltas)

	for i, mesh := range m.meshes {
		mesh.ibo.Bind()

		shader.Use(mview.PROGRAM_TYPE_MODEL)
		mesh.DrawModel(shader, windowIndex, boneDeltas, materialDeltas[i])
		shader.Unuse()

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			shader.Use(mview.PROGRAM_TYPE_EDGE)
			mesh.DrawEdge(shader, windowIndex, boneDeltas, materialDeltas[i])
			shader.Unuse()
		}

		mesh.ibo.Unbind()
	}

	m.vbo.Unbind()
	m.vao.Unbind()

	if isDrawNormal {
		m.drawNormal(shader, boneDeltas, windowIndex)
	}

	if isDrawBone {
		m.drawBone(shader, boneDeltas, windowIndex)
	}

	shader.Msaa.Unbind()

	gl.Disable(gl.BLEND)
	gl.Disable(gl.DEPTH_TEST)
}

func (m *Meshes) drawNormal(
	shader *mview.MShader,
	boneDeltas []*mgl32.Mat4,
	windowIndex int,
) {
	shader.Use(mview.PROGRAM_TYPE_NORMAL)

	m.normalVao.Bind()
	m.normalVbo.BindVertex(nil, nil)
	m.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	BindBoneMatrixes(boneDeltas, shader, shader.NormalProgram, windowIndex)

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
	boneDeltas []*mgl32.Mat4,
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
	m.boneVbo.BindVertex(nil, nil)
	m.boneIbo.Bind()

	// ボーンデフォームテクスチャ設定
	BindBoneMatrixes(boneDeltas, shader, shader.BoneProgram, windowIndex)

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