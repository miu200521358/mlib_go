package pmx

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/mgl"
	"github.com/miu200521358/mlib_go/pkg/mutils"

)

type Mesh struct {
	material          MaterialGL // 描画用材質
	prevVerticesCount int        // 前の頂点数
	ibo               *mgl.IBO   // 頂点インデックスバッファ
}

func NewMesh(
	allFaces []uint32,
	material *MaterialGL,
	prevVerticesCount int,
) *Mesh {
	faces := allFaces[prevVerticesCount:(prevVerticesCount + material.VerticesCount)]
	// println(
	// 	"faces",
	// 	mutils.JoinSlice(mutils.ConvertUint32ToInterfaceSlice(faces)),
	// 	"prevVerticesCount",
	// 	prevVerticesCount,
	// 	"material.VerticesCount",
	// 	material.VerticesCount,
	// )

	ibo := mgl.NewIBO(gl.Ptr(faces), len(faces))

	return &Mesh{
		material:          *material,
		prevVerticesCount: prevVerticesCount,
		ibo:               ibo,
	}
}

func (m *Mesh) DrawModel(
	shader *mgl.MShader,
	windowIndex int,
	boneMatrixes []mgl32.Mat4,
) {
	m.ibo.Bind()

	if m.material.DrawFlag.IsDoubleSidedDrawing() {
		// 両面描画
		// カリングOFF
		gl.Disable(gl.CULL_FACE)
	} else { // 片面描画
		// カリングON
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)
	}

	// // ボーンデフォームテクスチャ設定
	// m.BindBoneMatrixes(boneMatrixes, shader, mgl.PROGRAM_TYPE_MODEL, windowIndex)

	// ------------------
	// 材質色設定
	// full.fx の AmbientColor相当
	diffuse := m.material.Diffuse()
	diffuseUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_DIFFUSE))
	gl.Uniform4fv(diffuseUniform, 1, &diffuse[0])

	ambient := m.material.Ambient()
	ambientUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_AMBIENT))
	gl.Uniform3fv(ambientUniform, 1, &ambient[0])

	specular := m.material.Specular()
	specularUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_SPECULAR))
	gl.Uniform4fv(specularUniform, 1, &specular[0])

	// テクスチャ使用有無
	useTextureUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_USE_TEXTURE))
	gl.Uniform1i(useTextureUniform, mutils.BoolToInt(m.material.Texture != nil))
	if m.material.Texture != nil {
		m.material.Texture.Bind()
		textureUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_TEXTURE_SAMPLER))
		gl.Uniform1i(textureUniform, int32(m.material.Texture.TextureUnitNo))
	}

	// Toon使用有無
	useToonUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_USE_TOON))
	gl.Uniform1i(useToonUniform, mutils.BoolToInt(m.material.ToonTexture != nil))
	if m.material.ToonTexture != nil {
		m.material.ToonTexture.Bind()
		toonUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_TOON_SAMPLER))
		gl.Uniform1i(toonUniform, int32(m.material.ToonTexture.TextureUnitNo))
	}

	// Sphere使用有無
	useSphereUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_USE_SPHERE))
	gl.Uniform1i(useSphereUniform, mutils.BoolToInt(m.material.SphereTexture != nil && m.material.SphereTexture.Valid))
	if m.material.SphereTexture != nil && m.material.SphereTexture.Valid {
		m.material.SphereTexture.Bind()
		sphereUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_SPHERE_SAMPLER))
		gl.Uniform1i(sphereUniform, int32(m.material.SphereTexture.TextureUnitNo))
	}

	// SphereMode
	sphereModeUniform := gl.GetUniformLocation(shader.ModelProgram, gl.Str(mgl.SHADER_SPHERE_MODE))
	gl.Uniform1i(sphereModeUniform, int32(m.material.SphereMode))

	// // ウェイト描写
	// gl.Uniform1i(
	// 	shader.IsShowBoneWeightUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	mutils.BoolToInt(isShowBoneWeight),
	// )
	// gl.Uniform1iv(
	// 	shader.ShowBoneIndexesUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	int32(len(showBoneIndexes)),
	// 	(*int32)(unsafe.Pointer(&showBoneIndexes[0])),
	// )

	// prevVerticesSize := m.prevVerticesCount * int(ibo.Dtype)

	// if err := mutils.CheckGLError(); err != nil {
	// 	panic(fmt.Errorf("Mesh DrawElements failed: %v", err))
	// }

	// gl.DrawElements(
	// 	gl.TRIANGLES,
	// 	int32(m.material.VerticesCount),
	// 	ibo.Dtype,
	// 	unsafe.Pointer(&prevVerticesCount),
	// )

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)

	if m.material.Texture != nil {
		m.material.Texture.Unbind()
	}

	if m.material.ToonTexture != nil {
		m.material.ToonTexture.Unbind()
	}

	if m.material.SphereTexture != nil {
		m.material.SphereTexture.Unbind()
	}

	// m.UnbindBoneMatrixes()
	m.ibo.Unbind()
}

func (m *Mesh) Delete() {
	m.ibo.Delete()
	if m.material.Texture != nil {
		m.material.Texture.Delete()
	}
}

// func (m *Mesh) BindBoneMatrixes(
// 	matrixes []mgl32.Mat4,
// 	shader *mgl.MShader,
// 	programType mgl.ProgramType,
// 	windowIndex int,
// ) {
// 	// テクスチャをアクティブにする
// 	gl.ActiveTexture(gl.TEXTURE20)

// 	// テクスチャをバインドする
// 	gl.BindTexture(
// 		gl.TEXTURE_2D, shader.BoneMatrixTextureId[programType],
// 	)

// 	// テクスチャのパラメーターの設定
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
// 	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

// 	// テクスチャのサイズを計算する
// 	numBones := len(matrixes)
// 	texSize := int(math.Ceil(math.Sqrt(float64(numBones))))
// 	width := int(math.Ceil(float64(texSize)/4) * 4 * 4)
// 	height := int(math.Ceil((float64(numBones) * 4) / float64(width)))

// 	paddedMatrixes := make([]float32, height*width*4)
// 	for i, matrix := range matrixes {
// 		copy(paddedMatrixes[i*16:], matrix[:])
// 	}

// 	// テクスチャをシェーダーに渡す
// 	gl.TexImage2D(
// 		gl.TEXTURE_2D,
// 		0,
// 		gl.RGBA32F,
// 		int32(width),
// 		int32(height),
// 		0,
// 		gl.RGBA,
// 		gl.FLOAT,
// 		unsafe.Pointer(&paddedMatrixes[0]),
// 	)

// 	gl.Uniform1i(shader.BoneMatrixTextureUniform[programType], 20)
// 	gl.Uniform1i(shader.BoneMatrixTextureWidth[programType], int32(width))
// 	gl.Uniform1i(shader.BoneMatrixTextureHeight[programType], int32(height))
// }

// func (m *Mesh) UnbindBoneMatrixes() {
// 	gl.BindTexture(gl.TEXTURE_2D, 0)
// }