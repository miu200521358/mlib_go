//go:build windows
// +build windows

package renderer

import (
	"math"
	"slices"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

func (renderModel *RenderModel) Draw(
	shader *mgl.MShader, boneDeltas []mgl32.Mat4,
	vertexMorphIndexes []int, vertexMorphDeltas [][]float32,
	selectedVertexIndexes []int, selectedVertexDeltas [][]float32, meshDeltas []*MeshDelta,
	invisibleMaterialIndexes []int, nextInvisibleMaterialIndexes []int, windowIndex int,
	isDrawNormal, isDrawWire, isDrawSelectedVertex bool, isDrawBones map[pmx.BoneFlag]bool, bones *pmx.Bones,
) [][]float32 {
	renderModel.vao.Bind()
	defer renderModel.vao.Unbind()

	renderModel.vbo.BindVertex(vertexMorphIndexes, vertexMorphDeltas)
	defer renderModel.vbo.Unbind()

	paddedMatrixes, matrixWidth, matrixHeight := renderModel.createBoneMatrixes(boneDeltas)

	for i, mesh := range renderModel.meshes {
		mesh.ibo.Bind()

		mesh.drawModel(shader, paddedMatrixes, matrixWidth, matrixHeight, meshDeltas[i])

		if mesh.material.DrawFlag.IsDrawingEdge() {
			// エッジ描画
			mesh.drawEdge(shader, paddedMatrixes, matrixWidth, matrixHeight, meshDeltas[i])
		}

		if isDrawWire {
			mesh.drawWire(shader, paddedMatrixes, matrixWidth, matrixHeight, ((len(invisibleMaterialIndexes) > 0 && len(nextInvisibleMaterialIndexes) == 0 &&
				slices.Contains(invisibleMaterialIndexes, mesh.material.Index)) ||
				slices.Contains(nextInvisibleMaterialIndexes, mesh.material.Index)))
		}

		mesh.ibo.Unbind()
	}

	vertexPositions := make([][]float32, 0)
	if isDrawSelectedVertex {
		vertexPositions = renderModel.drawSelectedVertex(selectedVertexIndexes, selectedVertexDeltas,
			shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	if isDrawNormal {
		renderModel.drawNormal(shader, paddedMatrixes, matrixWidth, matrixHeight)
	}

	isDrawBone := false
	for _, drawBone := range isDrawBones {
		if drawBone {
			isDrawBone = true
			break
		}
	}

	if isDrawBone {
		renderModel.drawBone(shader, bones, isDrawBones, paddedMatrixes, matrixWidth, matrixHeight)
	}

	paddedMatrixes = nil
	boneDeltas = nil
	meshDeltas = nil

	return vertexPositions
}

func (renderModel *RenderModel) createBoneMatrixes(matrixes []mgl32.Mat4) ([]float32, int, int) {
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

func (renderModel *RenderModel) drawNormal(
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
) {
	program := shader.GetProgram(mgl.PROGRAM_TYPE_NORMAL)
	gl.UseProgram(program)

	renderModel.normalVao.Bind()
	renderModel.normalVbo.BindVertex(nil, nil)
	renderModel.normalIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)

	normalColor := mgl32.Vec4{0.3, 0.3, 0.7, 0.5}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &normalColor[0])

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(renderModel.normalVertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.normalIbo.Unbind()
	renderModel.normalVbo.Unbind()
	renderModel.normalVao.Unbind()

	gl.UseProgram(0)
}

func (renderModel *RenderModel) drawSelectedVertex(
	selectedVertexMorphIndexes []int,
	selectedVertexDeltas [][]float32,
	shader *mgl.MShader,
	paddedMatrixes []float32,
	width, height int,
) [][]float32 {
	// モデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	program := shader.GetProgram(mgl.PROGRAM_TYPE_SELECTED_VERTEX)
	gl.UseProgram(program)

	renderModel.selectedVertexVao.Bind()
	renderModel.selectedVertexVbo.BindVertex(selectedVertexMorphIndexes, selectedVertexDeltas)
	renderModel.selectedVertexIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)

	vertexColor := mgl32.Vec4{1.0, 0.4, 0.0, 0.7}
	specularUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_COLOR))
	gl.Uniform4fv(specularUniform, 1, &vertexColor[0])
	gl.PointSize(5.0) // 選択頂点のサイズ

	// 点描画
	gl.DrawElements(
		gl.POINTS,
		int32(len(renderModel.vertices)),
		gl.UNSIGNED_INT,
		nil,
	)

	// SSBOからデータを読み込む
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, renderModel.ssbo)
	ptr := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)
	results := make([][]float32, renderModel.Model.Vertices.Len())

	if ptr != nil {
		// 頂点数を取得
		for i := range results {
			results[i] = make([]float32, 4)
		}

		// SSBOから読み取り
		for i := 0; i < renderModel.Model.Vertices.Len(); i++ {
			for j := 0; j < 4; j++ {
				results[i][j] = *(*float32)(unsafe.Pointer(uintptr(ptr) + uintptr((i*4+j)*4)))
			}
		}
	}

	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	renderModel.selectedVertexIbo.Unbind()
	renderModel.selectedVertexVbo.Unbind()
	renderModel.selectedVertexVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	return results
}

func (renderModel *RenderModel) drawBone(
	shader *mgl.MShader,
	bones *pmx.Bones,
	isDrawBones map[pmx.BoneFlag]bool,
	paddedMatrixes []float32,
	width, height int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	program := shader.GetProgram(mgl.PROGRAM_TYPE_BONE)
	gl.UseProgram(program)

	renderModel.boneVao.Bind()
	renderModel.boneVbo.BindVertex(renderModel.fetchBoneDebugDeltas(bones, isDrawBones))
	renderModel.boneIbo.Bind()

	// ボーンデフォームテクスチャ設定
	bindBoneMatrixes(paddedMatrixes, width, height, shader, program)

	// ライン描画
	gl.DrawElements(
		gl.LINES,
		int32(len(renderModel.bones)),
		gl.UNSIGNED_INT,
		nil,
	)

	renderModel.boneIbo.Unbind()
	renderModel.boneVbo.Unbind()
	renderModel.boneVao.Unbind()

	gl.UseProgram(0)

	// 深度テストを有効に戻す
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

}

func (renderModel *RenderModel) fetchBoneDebugDeltas(bones *pmx.Bones, isDrawBones map[pmx.BoneFlag]bool) ([]int, [][]float32) {
	indexes := make([]int, 0)
	deltas := make([][]float32, 0)

	for _, boneIndex := range renderModel.boneIndexes {
		bone := bones.Get(boneIndex)
		indexes = append(indexes, boneIndex)
		deltas = append(deltas, newBoneDebugAlphaGl(bone, isDrawBones))
	}

	return indexes, deltas
}
