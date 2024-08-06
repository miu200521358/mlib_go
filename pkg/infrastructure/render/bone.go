//go:build windows
// +build windows

package render

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

var bone_colors_ik = []float32{1.0, 0.38, 0, 1.0}
var bone_colors_ik_target = []float32{1.0, 0.57, 0.61, 1.0}
var bone_colors_ik_link = []float32{1.0, 0.83, 0.49, 1.0}
var bone_colors_fixed = []float32{0.72, 0.32, 1.0, 1.0}
var bone_colors_effect = []float32{0.68, 0.64, 1.0, 1.0}
var bone_colors_effect_effector = []float32{0.88, 0.84, 1.0, 0.7}
var bone_colors_translate = []float32{0.70, 1.0, 0.54, 1.0}
var bone_colors_rotate = []float32{0.56, 0.78, 1.0, 1.0}
var bone_colors_invisible = []float32{0.82, 0.82, 0.82, 1.0}

func newBoneGl(bone *pmx.Bone) []float32 {
	p := mgl.NewGlVec3(bone.Position)
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(bone.Index()), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0.0, 0.0, 0.0, 0.0, // 色
	}
}

func newTailBoneGl(bone *pmx.Bone) []float32 {
	p := mgl.NewGlVec3(bone.Position.Added(bone.Extend.ChildRelativePosition))
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(bone.Index()), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0.0, 0.0, 0.0, 0.0, // 色
	}
}

func getBoneDebugColor(bone *pmx.Bone, appState state.IAppState) []float32 {
	// IK
	if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneIk()) && bone.IsIK() {
		// IK
		return bone_colors_ik
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneIk()) &&
		len(bone.Extend.IkLinkBoneIndexes) > 0 {
		// IKリンク
		return bone_colors_ik_link
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneIk()) &&
		len(bone.Extend.IkTargetBoneIndexes) > 0 {
		// IKターゲット
		return bone_colors_ik_target
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneEffector()) &&
		(bone.IsEffectorRotation() || bone.IsEffectorTranslation()) {
		// 付与親
		return bone_colors_effect
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneEffector()) &&
		len(bone.Extend.EffectiveBoneIndexes) > 0 {
		// 付与親の付与元
		return bone_colors_effect_effector
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneFixed()) &&
		bone.HasFixedAxis() {
		// 軸固定
		return bone_colors_fixed
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneTranslate()) &&
		bone.CanTranslate() {
		// 移動
		return bone_colors_translate
	} else if (appState.IsShowBoneAll() || appState.IsShowBoneVisible() || appState.IsShowBoneRotate()) &&
		bone.CanRotate() {
		// 回転
		return bone_colors_rotate
	} else if appState.IsShowBoneAll() && !bone.IsVisible() {
		// 非表示
		return bone_colors_invisible
	}

	return []float32{0.0, 0.0, 0.0, 0.0}
}

func createBoneMatrixes(boneDeltas *delta.BoneDeltas) ([]float32, int, int, error) {
	// テクスチャのサイズを計算する
	numBones := len(boneDeltas.Data)
	texSize := int(math.Ceil(math.Sqrt(float64(numBones))))
	width := int(math.Ceil(float64(texSize)/4) * 4 * 4)
	height := int(math.Ceil((float64(numBones) * 4) / float64(width)))

	paddedMatrixes := make([]float32, height*width*4)
	for i, d := range boneDeltas.Data {
		var m mgl32.Mat4
		if d == nil {
			return nil, 0, 0, fmt.Errorf("boneDeltas[%d] is nil", i)
		} else {
			m = mgl.NewGlMat4(d.FilledLocalMatrix())
		}
		copy(paddedMatrixes[i*16:], m[:])
	}

	return paddedMatrixes, width, height, nil
}

func bindBoneMatrixes(
	windowIndex int,
	paddedMatrixes []float32,
	width, height int,
	shader mgl.IShader,
	program uint32,
) {
	// テクスチャをアクティブにする
	switch windowIndex {
	case 0:
		gl.ActiveTexture(gl.TEXTURE20)
	case 1:
		gl.ActiveTexture(gl.TEXTURE21)
	case 2:
		gl.ActiveTexture(gl.TEXTURE22)
	}

	// テクスチャをバインドする
	gl.BindTexture(gl.TEXTURE_2D, shader.BoneTextureId())

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

	boneUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_BONE_MATRIX_TEXTURE))
	switch windowIndex {
	case 0:
		gl.Uniform1i(boneUniform, 20)
	case 1:
		gl.Uniform1i(boneUniform, 21)
	case 2:
		gl.Uniform1i(boneUniform, 22)
	}

	modelWidthUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_BONE_MATRIX_TEXTURE_WIDTH))
	gl.Uniform1i(modelWidthUniform, int32(width))

	modelHeightUniform := gl.GetUniformLocation(program, gl.Str(mgl.SHADER_BONE_MATRIX_TEXTURE_HEIGHT))
	gl.Uniform1i(modelHeightUniform, int32(height))
}

func unbindBoneMatrixes() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
