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

// MeshRenderer はメッシュ描画専用の構造体です。
// 描画用マテリアル情報（materialGL）と、
// インデックスバッファ（ElementBuffer）を保持し、
// 各種描画メソッド（DrawModel, DrawEdge, DrawWire）を提供します。
type MeshRenderer struct {
	// 描画用マテリアル（pmx.Material と OpenGL拡張情報を保持）
	material materialGL

	// 前の材質までの頂点数（頂点配列のオフセット計算用）
	prevVerticesCount int

	// ElementBuffer（旧 IBO 相当）: インデックス配列を保持
	elemBuffer *mgl.ElementBuffer
}

// NewMeshRenderer はメッシュ描画用の MeshRenderer を生成します。
// allFaces : 全頂点インデックスの配列
// material : 対象材質の描画用情報（materialGL）
// prevVerticesCount : 前材質までの頂点数（オフセット）
func NewMeshRenderer(
	allFaces []uint32,
	material *materialGL,
	prevVerticesCount int,
) *MeshRenderer {
	// 対象材質の頂点数だけインデックス配列から切り出す
	faces := allFaces[prevVerticesCount : prevVerticesCount+material.VerticesCount]

	// ElementBuffer を生成して、faces のデータを転送する
	elemBuf := mgl.NewElementBuffer()
	if len(faces) > 0 {
		elemBuf.Bind()
		// 各インデックスは uint32 (4バイト)
		size := len(faces) * 4
		elemBuf.BufferData(size, gl.Ptr(faces), rendering.BufferUsageStatic)
		elemBuf.Unbind()
	}

	return &MeshRenderer{
		material:          *material,
		prevVerticesCount: prevVerticesCount,
		elemBuffer:        elemBuf,
	}
}

// Delete はメッシュに紐づく ElementBuffer を解放します。
func (mr *MeshRenderer) Delete() {
	mr.elemBuffer.Delete()
}

// DrawModel は通常描画（テクスチャ・ライティングあり）を行います。
func (mr *MeshRenderer) DrawModel(
	windowIndex int,
	shader rendering.IShader,
	paddedMatrixes []float32, // ボーン行列データ
	width, height int, // ボーン行列テクスチャの幅・高さ
	meshDelta *delta.MeshDelta,
) {
	modelProgram := shader.Program(rendering.ProgramTypeModel)
	gl.UseProgram(modelProgram)

	// 両面描画フラグに応じたカリング設定
	if mr.material.DrawFlag.IsDoubleSidedDrawing() {
		gl.Disable(gl.CULL_FACE)
	} else {
		gl.Enable(gl.CULL_FACE)
		defer gl.Disable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	}

	// ボーン行列テクスチャをシェーダーへバインド
	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, modelProgram)
	defer unbindBoneMatrixes()

	// 材質色などのシェーダーユニフォーム設定
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

	// 通常テクスチャ
	useTextureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseTexture))
	hasTexture := (mr.material.texture != nil)
	gl.Uniform1i(useTextureUniform, int32(mmath.BoolToInt(hasTexture)))
	if hasTexture {
		mr.material.texture.bind()
		defer mr.material.texture.unbind()
		texUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderTextureSampler))
		gl.Uniform1i(texUniform, int32(mr.material.texture.TextureUnitNo))
	}

	// Toonテクスチャ
	useToonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseToon))
	hasToon := (mr.material.toonTexture != nil)
	gl.Uniform1i(useToonUniform, int32(mmath.BoolToInt(hasToon)))
	if hasToon {
		mr.material.toonTexture.bind()
		defer mr.material.toonTexture.unbind()
		toonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderToonSampler))
		gl.Uniform1i(toonUniform, int32(mr.material.toonTexture.TextureUnitNo))
	}

	// スフィアテクスチャ
	useSphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseSphere))
	hasSphere := (mr.material.sphereTexture != nil && mr.material.sphereTexture.Initialized)
	gl.Uniform1i(useSphereUniform, int32(mmath.BoolToInt(hasSphere)))
	if hasSphere {
		mr.material.sphereTexture.bind()
		defer mr.material.sphereTexture.unbind()
		sphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSphereSampler))
		gl.Uniform1i(sphereUniform, int32(mr.material.sphereTexture.TextureUnitNo))
	}

	// スフィアモード
	sphereModeUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSphereMode))
	gl.Uniform1i(sphereModeUniform, int32(mr.material.SphereMode))

	// 描画：インデックスバッファをバインドして描画
	mr.elemBuffer.Bind()
	defer mr.elemBuffer.Unbind()

	gl.DrawElements(
		gl.TRIANGLES,
		int32(mr.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

// DrawEdge はエッジ（輪郭）描画を行います。
func (mr *MeshRenderer) DrawEdge(
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

	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	edgeColorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderEdgeColor))
	gl.Uniform4fv(edgeColorUniform, 1, &meshDelta.Edge[0])

	edgeSizeUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderEdgeSize))
	gl.Uniform1f(edgeSizeUniform, meshDelta.EdgeSize)

	mr.elemBuffer.Bind()
	defer mr.elemBuffer.Unbind()

	gl.DrawElements(
		gl.TRIANGLES,
		int32(mr.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

// DrawWire はワイヤーフレーム描画を行います。
func (mr *MeshRenderer) DrawWire(
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

	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, program)
	defer unbindBoneMatrixes()

	var wireColor mgl32.Vec4
	if invisibleMesh {
		wireColor = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	} else {
		wireColor = mgl32.Vec4{0.2, 0.6, 0.2, 0.5}
	}
	colorUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderColor))
	gl.Uniform4fv(colorUniform, 1, &wireColor[0])

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	mr.elemBuffer.Bind()
	defer mr.elemBuffer.Unbind()

	gl.DrawElements(
		gl.TRIANGLES,
		int32(mr.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	gl.UseProgram(0)
}
