//go:build windows
// +build windows

package render

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

type Mesh struct {
	material          materialGL               // 描画用材質
	prevVerticesCount int                      // 前の頂点数
	ibo               rendering.IElementBuffer // 頂点インデックスバッファ
}

func newMesh(
	allFaces []uint32,
	material *materialGL,
	prevVerticesCount int,
) *Mesh {
	faces := allFaces[prevVerticesCount:(prevVerticesCount + material.VerticesCount)]
	var ibo rendering.IElementBuffer
	if len(faces) == 0 {
		ebo := mgl.NewElementBuffer()
		ebo.BufferData(0, nil, rendering.BufferUsageStatic)
		ibo = ebo
	} else {
		ebo := mgl.NewElementBuffer()
		ebo.BufferData(len(faces)*4, gl.Ptr(faces), rendering.BufferUsageStatic)
		ibo = ebo
	}

	return &Mesh{
		material:          *material,
		prevVerticesCount: prevVerticesCount,
		ibo:               ibo,
	}
}

func (mesh *Mesh) delete() {
	mesh.ibo.Delete()
}

func (mesh *Mesh) drawModel(
	windowIndex int,
	shader rendering.IShader,
	paddedMatrixes []float32,
	width, height int,
	meshDelta *delta.MeshDelta,
) {
	modelProgram := shader.Program(rendering.ProgramTypeModel)
	gl.UseProgram(modelProgram)

	if mesh.material.DrawFlag.IsDoubleSidedDrawing() {
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
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, modelProgram)
	defer unbindBoneMatrixes()

	// ------------------
	// 材質色設定
	// full.fx の AmbientColor相当
	diffuseUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderDiffuse))
	gl.Uniform4fv(diffuseUniform, 1, &meshDelta.Diffuse[0])

	ambientUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderAmbient))
	gl.Uniform3fv(ambientUniform, 1, &meshDelta.Ambient[0])

	specularUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSpecular))
	gl.Uniform4fv(specularUniform, 1, &meshDelta.Specular[0])

	emissiveUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderEmissive))
	gl.Uniform3fv(emissiveUniform, 1, &meshDelta.Emissive[0])

	lightDiffuseUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderLightDiffuse))
	gl.Uniform3f(lightDiffuseUniform, delta.LIGHT_DIFFUSE, delta.LIGHT_DIFFUSE, delta.LIGHT_DIFFUSE)

	lightSpecularUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderLightSpecular))
	gl.Uniform3f(lightSpecularUniform, delta.LIGHT_SPECULAR, delta.LIGHT_SPECULAR, delta.LIGHT_SPECULAR)

	lightAmbientUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderLightAmbient))
	gl.Uniform3f(lightAmbientUniform, delta.LIGHT_AMBIENT, delta.LIGHT_AMBIENT, delta.LIGHT_AMBIENT)

	// テクスチャ使用有無
	useTextureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseTexture))
	gl.Uniform1i(useTextureUniform, int32(mmath.BoolToInt(mesh.material.texture != nil)))
	if mesh.material.texture != nil {
		mesh.material.texture.bind()
		defer mesh.material.texture.unbind()
		textureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderTextureSampler))
		gl.Uniform1i(textureUniform, int32(mesh.material.texture.TextureUnitNo))
	}

	// Toon使用有無
	useToonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseToon))
	gl.Uniform1i(useToonUniform, int32(mmath.BoolToInt(mesh.material.toonTexture != nil)))
	if mesh.material.toonTexture != nil {
		mesh.material.toonTexture.bind()
		defer mesh.material.toonTexture.unbind()
		toonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderToonSampler))
		gl.Uniform1i(toonUniform, int32(mesh.material.toonTexture.TextureUnitNo))
	}

	// Sphere使用有無
	useSphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseSphere))
	gl.Uniform1i(useSphereUniform,
		int32(mmath.BoolToInt(mesh.material.sphereTexture != nil && mesh.material.sphereTexture.Initialized)))
	if mesh.material.sphereTexture != nil && mesh.material.sphereTexture.Initialized {
		mesh.material.sphereTexture.bind()
		defer mesh.material.sphereTexture.unbind()
		sphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSphereSampler))
		gl.Uniform1i(sphereUniform, int32(mesh.material.sphereTexture.TextureUnitNo))
	}

	// SphereMode
	sphereModeUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSphereMode))
	gl.Uniform1i(sphereModeUniform, int32(mesh.material.SphereMode))

	// 頂点インデックスバッファをバインド
	mesh.ibo.Bind()
	defer mesh.ibo.Unbind()

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(mesh.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

func (mesh *Mesh) drawEdge(
	windowIndex int,
	shader rendering.IShader,
	paddedMatrixes []float32,
	width, height int,
	meshDelta *delta.MeshDelta,
) {
	program := shader.Program(rendering.ProgramTypeEdge)
	gl.UseProgram(program)

	gl.Enable(gl.CULL_FACE)
	defer gl.Disable(gl.CULL_FACE)

	gl.CullFace(gl.FRONT)

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	// ------------------
	// エッジ色設定
	edgeColorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderEdgeColor))
	gl.Uniform4fv(edgeColorUniform, 1, &meshDelta.Edge[0])

	// エッジサイズ
	edgeSizeUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderEdgeSize))
	gl.Uniform1f(edgeSizeUniform, meshDelta.EdgeSize)

	// 頂点インデックスバッファをバインド
	mesh.ibo.Bind()
	defer mesh.ibo.Unbind()

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(mesh.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

func (mesh *Mesh) drawWire(
	windowIndex int,
	shader rendering.IShader,
	paddedMatrixes []float32,
	width, height int,
	invisibleMesh bool,
) {
	program := shader.Program(rendering.ProgramTypeWire)
	gl.UseProgram(program)

	// カリングOFF
	gl.Disable(gl.CULL_FACE)

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	var wireColor mgl32.Vec4
	if invisibleMesh {
		wireColor = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	} else {
		wireColor = mgl32.Vec4{0.2, 0.6, 0.2, 0.5}
	}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderColor))
	gl.Uniform4fv(specularUniform, 1, &wireColor[0])

	// 描画モードをワイヤーフレームに変更
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// 頂点インデックスバッファをバインド
	mesh.ibo.Bind()
	defer mesh.ibo.Unbind()

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(mesh.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	// 描画モードを元に戻す
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	gl.UseProgram(0)
}
