//go:build windows
// +build windows

package renderer

import (
	"unsafe"

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

func NewMesh(
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
	defer UnbindBoneMatrixes()

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
	gl.Uniform1i(useTextureUniform, int32(mmath.BoolToInt(m.material.Texture != nil)))
	if m.material.Texture != nil {
		m.material.Texture.Bind()
		defer m.material.Texture.Unbind()
		textureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_TEXTURE_SAMPLER))
		gl.Uniform1i(textureUniform, int32(m.material.Texture.TextureUnitNo))
	}

	// Toon使用有無
	useToonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_USE_TOON))
	gl.Uniform1i(useToonUniform, int32(mmath.BoolToInt(m.material.ToonTexture != nil)))
	if m.material.ToonTexture != nil {
		m.material.ToonTexture.Bind()
		defer m.material.ToonTexture.Unbind()
		toonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_TOON_SAMPLER))
		gl.Uniform1i(toonUniform, int32(m.material.ToonTexture.TextureUnitNo))
	}

	// Sphere使用有無
	useSphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_USE_SPHERE))
	gl.Uniform1i(useSphereUniform,
		int32(mmath.BoolToInt(m.material.SphereTexture != nil && m.material.SphereTexture.Initialized)))
	if m.material.SphereTexture != nil && m.material.SphereTexture.Initialized {
		m.material.SphereTexture.Bind()
		defer m.material.SphereTexture.Unbind()
		sphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.SHADER_SPHERE_SAMPLER))
		gl.Uniform1i(sphereUniform, int32(m.material.SphereTexture.TextureUnitNo))
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
	defer UnbindBoneMatrixes()

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
	defer UnbindBoneMatrixes()

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

func bindBoneMatrixes(
	paddedMatrixes []float32,
	width, height int,
	shader *mgl.MShader,
	program uint32,
) {
	// テクスチャをアクティブにする
	gl.ActiveTexture(gl.TEXTURE20)

	// テクスチャをバインドする
	gl.BindTexture(gl.TEXTURE_2D, shader.BoneTextureId)

	// テクスチャのパラメーターの設定
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// テクスチャをシェーダーに渡す
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA32F,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.FLOAT,
		unsafe.Pointer(&paddedMatrixes[0]),
	)

	modelUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_BONE_MATRIX_TEXTURE))
	gl.Uniform1i(modelUniform, 20)

	modelWidthUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_BONE_MATRIX_TEXTURE_WIDTH))
	gl.Uniform1i(modelWidthUniform, int32(width))

	modelHeightUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_BONE_MATRIX_TEXTURE_HEIGHT))
	gl.Uniform1i(modelHeightUniform, int32(height))
}

func UnbindBoneMatrixes() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
