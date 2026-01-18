//go:build windows
// +build windows

package render

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
)

// ライト定数。
const (
	LightDiffuse  float32 = 0
	LightSpecular float32 = 154.0 / 255.0
	LightAmbient  float32 = 154.0 / 255.0
)

// MeshDelta は材質モーフ差分を描画用に変換した値を表す。
type MeshDelta struct {
	Diffuse  mgl32.Vec4
	Specular mgl32.Vec4
	Ambient  mgl32.Vec3
	Emissive mgl32.Vec3
	Edge     mgl32.Vec4
	EdgeSize float32
}

// NewMeshDelta は材質モーフ差分を描画向けに変換する。
func NewMeshDelta(materialDelta *delta.MaterialMorphDelta) *MeshDelta {
	if materialDelta == nil {
		materialDelta = delta.NewMaterialMorphDelta(nil)
	}

	base := materialDelta.Material
	add := materialDelta.AddMaterial
	mul := materialDelta.MulMaterial

	diffuse := base.Diffuse.Muled(mul.Diffuse).Added(add.Diffuse)
	specular := base.Specular.Muled(mul.Specular).Added(add.Specular)
	ambient := base.Diffuse.XYZ().Muled(mul.Diffuse.XYZ()).Added(add.Diffuse.XYZ())
	emissive := base.Ambient.Muled(mul.Ambient).Added(add.Ambient)
	edge := base.Edge.Muled(mul.Edge).Added(add.Edge)
	edgeSize := float32(base.EdgeSize*mul.EdgeSize + add.EdgeSize)

	return &MeshDelta{
		Diffuse:  vec4ToMgl(diffuse),
		Specular: vec4ToMgl(specular),
		Ambient:  vec3ToMgl(ambient),
		Emissive: vec3ToMgl(emissive),
		Edge:     vec4ToMgl(edge),
		EdgeSize: edgeSize,
	}
}

// vec3ToMgl はmmath.Vec3をmgl32.Vec3へ変換する。
func vec3ToMgl(v mmath.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v.X), float32(v.Y), float32(v.Z)}
}

// vec4ToMgl はmmath.Vec4をmgl32.Vec4へ変換する。
func vec4ToMgl(v mmath.Vec4) mgl32.Vec4 {
	return mgl32.Vec4{float32(v.X), float32(v.Y), float32(v.Z), float32(v.W)}
}

// MeshRenderer はメッシュ描画専用の構造体。
type MeshRenderer struct {
	// 描画用マテリアル（材質とOpenGL拡張情報）
	material materialGL

	// 前の材質までの頂点数（頂点配列のオフセット計算用）
	prevVerticesCount int

	// IndexBuffer（旧 IBO 相当）: インデックス配列を保持
	elemBuffer *mgl.IndexBuffer
}

// NewMeshRenderer はメッシュ描画用の MeshRenderer を生成する。
func NewMeshRenderer(
	factory *mgl.BufferFactory,
	allFaces []uint32,
	material *materialGL,
	prevVerticesCount int,
) *MeshRenderer {
	faces := allFaces[prevVerticesCount : prevVerticesCount+material.VerticesCount]

	var elemBuf *mgl.IndexBuffer
	if len(faces) > 0 {
		elemBuf = factory.NewIndexBuffer(gl.Ptr(faces), len(faces))
	}

	return &MeshRenderer{
		material:          *material,
		prevVerticesCount: prevVerticesCount,
		elemBuffer:        elemBuf,
	}
}

// delete はメッシュに紐づく IndexBuffer を解放する。
func (mr *MeshRenderer) delete() {
	if mr.elemBuffer != nil {
		mr.elemBuffer.Delete()
	}
}

// drawModel は通常描画（テクスチャ・ライティングあり）を行う。
func (mr *MeshRenderer) drawModel(
	windowIndex int,
	shader graphics_api.IShader,
	paddedMatrixes []float32,
	width, height int,
	meshDelta *MeshDelta,
) {
	modelProgram := shader.Program(graphics_api.ProgramTypeModel)
	gl.UseProgram(modelProgram)

	if hasDrawFlag(mr.material.DrawFlag, model.DRAW_FLAG_DOUBLE_SIDED_DRAWING) {
		gl.Disable(gl.CULL_FACE)
	} else {
		gl.Enable(gl.CULL_FACE)
		defer gl.Disable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	}

	bindBoneMatrixes(windowIndex, paddedMatrixes, width, height, shader, modelProgram)
	defer unbindBoneMatrixes()

	diffuseUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderDiffuse))
	gl.Uniform4fv(diffuseUniform, 1, &meshDelta.Diffuse[0])

	ambientUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderAmbient))
	gl.Uniform3fv(ambientUniform, 1, &meshDelta.Ambient[0])

	specularUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSpecular))
	gl.Uniform4fv(specularUniform, 1, &meshDelta.Specular[0])

	emissiveUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderEmissive))
	gl.Uniform3fv(emissiveUniform, 1, &meshDelta.Emissive[0])

	lightDiffuseUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderLightDiffuse))
	gl.Uniform3f(lightDiffuseUniform, LightDiffuse, LightDiffuse, LightDiffuse)

	lightSpecularUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderLightSpecular))
	gl.Uniform3f(lightSpecularUniform, LightSpecular, LightSpecular, LightSpecular)

	lightAmbientUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderLightAmbient))
	gl.Uniform3f(lightAmbientUniform, LightAmbient, LightAmbient, LightAmbient)

	useTextureUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseTexture))
	hasTexture := (mr.material.texture != nil)
	gl.Uniform1i(useTextureUniform, int32(mmath.BoolToInt(hasTexture)))
	if hasTexture {
		mr.material.texture.bind()
		defer mr.material.texture.unbind()
		texUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderTextureSampler))
		gl.Uniform1i(texUniform, int32(mr.material.texture.TextureUnitNo))
	}

	useToonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseToon))
	hasToon := (mr.material.toonTexture != nil)
	gl.Uniform1i(useToonUniform, int32(mmath.BoolToInt(hasToon)))
	if hasToon {
		mr.material.toonTexture.bind()
		defer mr.material.toonTexture.unbind()
		toonUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderToonSampler))
		gl.Uniform1i(toonUniform, int32(mr.material.toonTexture.TextureUnitNo))
	}

	useSphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderUseSphere))
	hasSphere := (mr.material.sphereTexture != nil && mr.material.sphereTexture.Initialized)
	gl.Uniform1i(useSphereUniform, int32(mmath.BoolToInt(hasSphere)))
	if hasSphere {
		mr.material.sphereTexture.bind()
		defer mr.material.sphereTexture.unbind()
		sphereUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSphereSampler))
		gl.Uniform1i(sphereUniform, int32(mr.material.sphereTexture.TextureUnitNo))
	}

	sphereModeUniform := gl.GetUniformLocation(modelProgram, gl.Str(mgl.ShaderSphereMode))
	gl.Uniform1i(sphereModeUniform, int32(mr.material.SphereMode))

	gl.DrawElements(
		gl.TRIANGLES,
		int32(mr.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	gl.UseProgram(0)
}

// drawEdge はエッジ（輪郭）描画を行う。
func (mr *MeshRenderer) drawEdge(
	windowIndex int,
	shader graphics_api.IShader,
	paddedMatrixes []float32,
	width, height int,
	meshDelta *MeshDelta,
) {
	program := shader.Program(graphics_api.ProgramTypeEdge)
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

// drawWire はワイヤーフレーム描画を行う。
func (mr *MeshRenderer) drawWire(
	windowIndex int,
	shader graphics_api.IShader,
	paddedMatrixes []float32,
	width, height int,
	invisibleMesh bool,
) {
	program := shader.Program(graphics_api.ProgramTypeWire)
	gl.UseProgram(program)

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

// hasDrawFlag は描画フラグの有無を判定する。
func hasDrawFlag(value model.DrawFlag, flag model.DrawFlag) bool {
	return value&flag != 0
}
