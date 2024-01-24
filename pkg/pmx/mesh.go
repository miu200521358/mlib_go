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

	// if m.material.DrawFlag == DRAW_FLAG_DOUBLE_SIDED_DRAWING {
	// 両面描画
	// カリングOFF
	gl.Disable(gl.CULL_FACE)
	// } else {
	// 	// 片面描画
	// 	// カリングON
	// 	gl.Enable(gl.CULL_FACE)
	// 	gl.CullFace(gl.BACK)
	// }

	// // ボーンデフォームテクスチャ設定
	// m.BindBoneMatrixes(boneMatrixes, shader, mgl.PROGRAM_TYPE_MODEL, windowIndex)

	// ------------------
	// 材質色設定
	// full.fx の AmbientColor相当
	gl.Uniform4f(
		shader.DiffuseUniform[mgl.PROGRAM_TYPE_MODEL],
		m.material.Diffuse[0],
		m.material.Diffuse[1],
		m.material.Diffuse[2],
		m.material.Diffuse[3],
	)
	mutils.CheckGLError()

	// 三角形描画
	gl.DrawElements(
		gl.TRIANGLES,
		int32(m.material.VerticesCount),
		gl.UNSIGNED_INT,
		nil,
	)
	mutils.CheckGLError()

	// ambient := m.material.Ambient
	// gl.Uniform3f(
	// 	shader.AmbientUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	ambient[0],
	// 	ambient[1],
	// 	ambient[2],
	// )

	// specular := m.material.Specular
	// gl.Uniform4f(
	// 	shader.SpecularUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	specular[0],
	// 	specular[1],
	// 	specular[2],
	// 	specular[3],
	// )

	// // テクスチャ使用有無
	// gl.Uniform1i(
	// 	shader.UseTextureUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	mutils.BoolToInt(m.material.Texture != nil && m.material.Texture.Valid),
	// )
	// if m.material.Texture != nil && m.material.Texture.Valid {
	// 	gl.Uniform1i(
	// 		shader.TextureUniform[mgl.PROGRAM_TYPE_MODEL],
	// 		int32(m.material.Texture.textureType),
	// 	)
	// 	textureFactor := shaderMaterial.TextureFactor()
	// 	gl.Uniform4f(
	// 		shader.TextureFactorUniform[mgl.PROGRAM_TYPE_MODEL],
	// 		textureFactor[0],
	// 		textureFactor[1],
	// 		textureFactor[2],
	// 		textureFactor[3],
	// 	)
	// }

	// // Toon使用有無
	// gl.Uniform1i(
	// 	shader.UseToonUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	mutils.BoolToInt(m.material.ToonTexture != nil && m.material.ToonTexture.Valid),
	// )
	// if m.material.ToonTexture != nil && m.material.ToonTexture.Valid {
	// 	gl.Uniform1i(
	// 		shader.ToonUniform[mgl.PROGRAM_TYPE_MODEL],
	// 		int32(m.material.ToonTexture.textureType),
	// 	)
	// 	toonTextureFactor := shaderMaterial.ToonTextureFactor()
	// 	gl.Uniform4f(
	// 		shader.ToonFactorUniform[mgl.PROGRAM_TYPE_MODEL],
	// 		toonTextureFactor[0],
	// 		toonTextureFactor[1],
	// 		toonTextureFactor[2],
	// 		toonTextureFactor[3],
	// 	)
	// }

	// // Sphere使用有無
	// gl.Uniform1i(
	// 	shader.UseSphereUniform[mgl.PROGRAM_TYPE_MODEL],
	// 	mutils.BoolToInt(m.material.SphereTexture != nil && m.material.SphereTexture.Valid && m.material.SphereMode != SPHERE_MODE_INVALID),
	// )
	// if m.material.SphereTexture != nil && m.material.SphereTexture.Valid {
	// 	gl.Uniform1i(
	// 		shader.SphereUniform[mgl.PROGRAM_TYPE_MODEL],
	// 		int32(m.material.SphereTexture.textureType),
	// 	)
	// 	sphereTextureFactor := shaderMaterial.SphereTextureFactor()
	// 	gl.Uniform4f(
	// 		shader.SphereFactorUniform[mgl.PROGRAM_TYPE_MODEL],
	// 		sphereTextureFactor[0],
	// 		sphereTextureFactor[1],
	// 		sphereTextureFactor[2],
	// 		sphereTextureFactor[3],
	// 	)
	// }

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

	// if m.material.Texture != nil {
	// 	m.material.Texture.Unbind()
	// }

	// if m.material.ToonTexture != nil {
	// 	m.material.ToonTexture.Unbind()
	// }

	// if m.material.SphereTexture != nil {
	// 	m.material.SphereTexture.Unbind()
	// }

	// m.UnbindBoneMatrixes()
	m.ibo.Unbind()
}

func (m *Mesh) Delete() {
	m.ibo.Delete()
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
