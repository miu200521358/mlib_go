//go:build windows
// +build windows

package renderer

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl/buffer"
)

type Mesh struct {
	material          materialGL  // 描画用材質
	prevVerticesCount int         // 前の頂点数
	ibo               *buffer.IBO // 頂点インデックスバッファ
}

func newMesh(
	allFaces []uint32,
	material *materialGL,
	prevVerticesCount int,
) *Mesh {
	faces := allFaces[prevVerticesCount:(prevVerticesCount + material.VerticesCount)]
	ibo := buffer.NewIBO(gl.Ptr(faces), len(faces))

	return &Mesh{
		material:          *material,
		prevVerticesCount: prevVerticesCount,
		ibo:               ibo,
	}
}

func (m *Mesh) delete() {
	m.ibo.Delete()
}

func (m *Mesh) drawModel(
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
	meshDelta *MeshDelta,
) {
	modelProgram := shader.GetProgram(mgl.PROGRAM_TYPE_MODEL)
	gl.UseProgram(modelProgram)

	if m.material.DrawFlag.IsDoubleSidedDrawing() {
		// 両面描画
		// カリングOFF
		gl.Disable(gl.CULL_FACE)
	} else {
		// 片面描画
		// カリングON
		gl.Enable(gl.CULL_FACE)
		defer gl.Disable(gl.CULL_FACE)

		gl.CullFace(gl.BACK)
	}

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, modelProgram)
	defer unbindBoneMatrixes()

	// ------------------
	// 材質色設定
	// full.fx の AmbientColor相当
	diffuseUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_DIFFUSE))
	gl.Uniform4fv(diffuseUniform, 1, &meshDelta.Diffuse[0])

	ambientUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_AMBIENT))
	gl.Uniform3fv(ambientUniform, 1, &meshDelta.Ambient[0])

	specularUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_SPECULAR))
	gl.Uniform4fv(specularUniform, 1, &meshDelta.Specular[0])

	// テクスチャ使用有無
	useTextureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_USE_TEXTURE))
	gl.Uniform1i(useTextureUniform, int32(mmath.BoolToInt(m.material.texture != nil)))
	if m.material.texture != nil {
		m.material.texture.bind()
		defer m.material.texture.unbind()
		textureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_TEXTURE_SAMPLER))
		gl.Uniform1i(textureUniform, int32(m.material.texture.TextureUnitNo))
	}

	// Toon使用有無
	useToonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_USE_TOON))
	gl.Uniform1i(useToonUniform, int32(mmath.BoolToInt(m.material.toonTexture != nil)))
	if m.material.toonTexture != nil {
		m.material.toonTexture.bind()
		defer m.material.toonTexture.unbind()
		toonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_TOON_SAMPLER))
		gl.Uniform1i(toonUniform, int32(m.material.toonTexture.TextureUnitNo))
	}

	// Sphere使用有無
	useSphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_USE_SPHERE))
	gl.Uniform1i(useSphereUniform,
		int32(mmath.BoolToInt(m.material.sphereTexture != nil && m.material.sphereTexture.Initialized)))
	if m.material.sphereTexture != nil && m.material.sphereTexture.Initialized {
		m.material.sphereTexture.bind()
		defer m.material.sphereTexture.unbind()
		sphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_SPHERE_SAMPLER))
		gl.Uniform1i(sphereUniform, int32(m.material.sphereTexture.TextureUnitNo))
	}

	// SphereMode
	sphereModeUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_SPHERE_MODE))
	gl.Uniform1i(sphereModeUniform, int32(m.material.SphereMode))

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

func (m *Mesh) drawEdge(
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
	meshDelta *MeshDelta,
) {
	program := shader.GetProgram(mgl.PROGRAM_TYPE_EDGE)
	gl.UseProgram(program)

	gl.Enable(gl.CULL_FACE)
	defer gl.Disable(gl.CULL_FACE)

	gl.CullFace(gl.FRONT)

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	// ------------------
	// エッジ色設定
	edgeColorUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_EDGE_COLOR))
	gl.Uniform4fv(edgeColorUniform, 1, &meshDelta.Edge[0])

	// エッジサイズ
	edgeSizeUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_EDGE_SIZE))
	gl.Uniform1f(edgeSizeUniform, meshDelta.EdgeSize)

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

func (m *Mesh) drawWire(
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
	invisibleMesh bool,
) {
	program := shader.GetProgram(mgl.PROGRAM_TYPE_WIRE)
	gl.UseProgram(program)

	// カリングOFF
	gl.Disable(gl.CULL_FACE)

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	var wireColor mgl32.Vec4
	if invisibleMesh {
		wireColor = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	} else {
		wireColor = mgl32.Vec4{0.2, 0.6, 0.2, 0.5}
	}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &wireColor[0])

	// 描画モードをワイヤーフレームに変更
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	// 描画モードを元に戻す
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	gl.UseProgram(0)
}
