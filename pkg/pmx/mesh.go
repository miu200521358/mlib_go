package pmx

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type Mesh struct {
	material          Material
	texture           *Texture
	toonTexture       *Texture
	sphereTexture     *Texture
	prevVerticesCount int
	faceDtype         uint32
}

func NewMesh(
	material *Material,
	texture, toonTexture, sphereTexture *Texture,
	prevVerticesCount int,
	faceDtype uint32,
) *Mesh {

	return &Mesh{
		material:          *material,
		texture:           texture,
		toonTexture:       toonTexture,
		sphereTexture:     sphereTexture,
		prevVerticesCount: prevVerticesCount,
		faceDtype:         faceDtype,
	}
}

func (m *Mesh) DrawModel(
	boneMatrixes []mgl32.Mat4,
	shaderMaterial ShaderMaterial,
	shader mgl.MShader,
	ibo mgl.IBO,
	isShowBoneWeight bool,
	showBoneIndexes []int,
) {
	if m.material.DrawFlag == DRAW_FLAG_DOUBLE_SIDED_DRAWING {
		// 両面描画
		// カリングOFF
		gl.Disable(gl.CULL_FACE)
	} else {
		// 片面描画
		// カリングON
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	}

	// ボーンデフォームテクスチャ設定
	m.BindBoneMatrixes(boneMatrixes, shader, mgl.PROGRAM_TYPE_MODEL)

	// ------------------
	// 材質色設定
	// full.fx の AmbientColor相当
	diffuseColor := shaderMaterial.Diffuse()
	gl.Uniform3f(
		shader.DiffuseColorUniform[mgl.PROGRAM_TYPE_MODEL],
		diffuseColor[0],
		diffuseColor[1],
		diffuseColor[2],
	)
	gl.Uniform1f(
		shader.DiffuseAlphaUniform[mgl.PROGRAM_TYPE_MODEL],
		shaderMaterial.DiffuseAlpha(),
	)

	ambientColor := shaderMaterial.Ambient()
	gl.Uniform3f(
		shader.AmbientColorUniform[mgl.PROGRAM_TYPE_MODEL],
		ambientColor[0],
		ambientColor[1],
		ambientColor[2],
	)

	specularColor := shaderMaterial.Specular()
	gl.Uniform3f(
		shader.SpecularColorUniform[mgl.PROGRAM_TYPE_MODEL],
		specularColor[0],
		specularColor[1],
		specularColor[2],
	)
	gl.Uniform1f(
		shader.SpecularPowerUniform[mgl.PROGRAM_TYPE_MODEL],
		shaderMaterial.SpecularPower(),
	)

	// テクスチャ使用有無
	gl.Uniform1i(
		shader.UseTextureUniform[mgl.PROGRAM_TYPE_MODEL],
		mutils.BoolToInt(m.texture != nil && m.texture.Valid),
	)
	if m.texture != nil && m.texture.Valid {
		gl.Uniform1i(
			shader.TextureUniform[mgl.PROGRAM_TYPE_MODEL],
			int32(m.texture.TextureType),
		)
		textureFactor := shaderMaterial.TextureFactor()
		gl.Uniform4f(
			shader.TextureFactorUniform[mgl.PROGRAM_TYPE_MODEL],
			textureFactor[0],
			textureFactor[1],
			textureFactor[2],
			textureFactor[3],
		)
	}

	// Toon使用有無
	gl.Uniform1i(
		shader.UseToonUniform[mgl.PROGRAM_TYPE_MODEL],
		mutils.BoolToInt(m.toonTexture != nil && m.toonTexture.Valid),
	)
	if m.toonTexture != nil && m.toonTexture.Valid {
		gl.Uniform1i(
			shader.ToonUniform[mgl.PROGRAM_TYPE_MODEL],
			int32(m.toonTexture.TextureType),
		)
		toonTextureFactor := shaderMaterial.ToonTextureFactor()
		gl.Uniform4f(
			shader.ToonFactorUniform[mgl.PROGRAM_TYPE_MODEL],
			toonTextureFactor[0],
			toonTextureFactor[1],
			toonTextureFactor[2],
			toonTextureFactor[3],
		)
	}

	// Sphere使用有無
	gl.Uniform1i(
		shader.UseSphereUniform[mgl.PROGRAM_TYPE_MODEL],
		mutils.BoolToInt(m.sphereTexture != nil && m.sphereTexture.Valid && m.material.SphereMode != SPHERE_MODE_INVALID),
	)
	if m.sphereTexture != nil && m.sphereTexture.Valid {
		gl.Uniform1i(
			shader.SphereUniform[mgl.PROGRAM_TYPE_MODEL],
			int32(m.sphereTexture.TextureType),
		)
		sphereTextureFactor := shaderMaterial.SphereTextureFactor()
		gl.Uniform4f(
			shader.SphereFactorUniform[mgl.PROGRAM_TYPE_MODEL],
			sphereTextureFactor[0],
			sphereTextureFactor[1],
			sphereTextureFactor[2],
			sphereTextureFactor[3],
		)
	}

	// ウェイト描写
	gl.Uniform1i(
		shader.IsShowBoneWeightUniform[mgl.PROGRAM_TYPE_MODEL],
		mutils.BoolToInt(isShowBoneWeight),
	)
	gl.Uniform1iv(
		shader.ShowBoneIndexesUniform[mgl.PROGRAM_TYPE_MODEL],
		int32(len(showBoneIndexes)),
		(*int32)(unsafe.Pointer(&showBoneIndexes[0])),
	)

	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		uint32(ibo.Dtype),
		unsafe.Pointer(&m.prevVerticesCount),
	)

	if m.texture != nil && m.texture.Valid {
		m.texture.Unbind()
	}

	if m.toonTexture != nil && m.toonTexture.Valid {
		m.toonTexture.Unbind()
	}

	if m.sphereTexture != nil && m.sphereTexture.Valid {
		m.sphereTexture.Unbind()
	}

	m.UnbindBoneMatrixes()
}

func (m *Mesh) DrawEdge(
	boneMatrixes []mgl32.Mat4,
	shaderMaterial ShaderMaterial,
	shader mgl.MShader,
	ibo mgl.IBO,
) {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)

	// ボーンデフォームテクスチャ設定
	m.BindBoneMatrixes(boneMatrixes, shader, mgl.PROGRAM_TYPE_EDGE)

	// ------------------
	// エッジ設定
	edgeColor := shaderMaterial.EdgeColor()
	gl.Uniform3f(
		shader.EdgeColorUniform[mgl.PROGRAM_TYPE_EDGE],
		edgeColor[0],
		edgeColor[1],
		edgeColor[2],
	)
	gl.Uniform1f(
		shader.EdgeAlphaUniform[mgl.PROGRAM_TYPE_EDGE],
		shaderMaterial.EdgeAlpha(),
	)
	gl.Uniform1f(
		shader.EdgeSizeUniform[mgl.PROGRAM_TYPE_EDGE],
		shaderMaterial.EdgeSize(),
	)

	// 前の頂点からのオフセット
	prevVertexPtr := m.prevVerticesCount * int(m.faceDtype)
	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		ibo.Dtype,
		unsafe.Pointer(&prevVertexPtr),
	)

	errorCode := gl.GetError()
	if errorCode != gl.NO_ERROR {
		panic(fmt.Sprintf("Mesh draw_edge Failure\n%d", errorCode))
	}

	m.UnbindBoneMatrixes()
}

func (m *Mesh) BindBoneMatrixes(
	matrixes []mgl32.Mat4,
	shader mgl.MShader,
	programType mgl.ProgramType,
) {
	// テクスチャをアクティブにする
	gl.ActiveTexture(gl.TEXTURE3)

	// テクスチャをバインドする
	gl.BindTexture(
		gl.TEXTURE_2D, shader.BoneMatrixTextureId[programType],
	)

	// テクスチャのパラメーターの設定
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// テクスチャのサイズを計算する
	numBones := len(matrixes)
	texSize := int(math.Ceil(math.Sqrt(float64(numBones))))
	width := int(math.Ceil(float64(texSize)/4) * 4 * 4)
	height := int(math.Ceil((float64(numBones) * 4) / float64(width)))

	paddedMatrixes := make([]float32, height*width*4)
	for i, matrix := range matrixes {
		copy(paddedMatrixes[i*16:], matrix[:])
	}

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

	gl.Uniform1i(shader.BoneMatrixTextureUniform[programType], 3)
	gl.Uniform1i(shader.BoneMatrixTextureWidth[programType], int32(width))
	gl.Uniform1i(shader.BoneMatrixTextureHeight[programType], int32(height))
}

func (m *Mesh) UnbindBoneMatrixes() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
