//go:build windows
// +build windows

package render

import (
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/rendering"
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

// デバッグ表示用のボーンカラー定義 (旧 bone.go での変数群)
var (
	boneColorsIK             = []float32{1.0, 0.38, 0, 1.0}
	boneColorsIKTarget       = []float32{1.0, 0.57, 0.61, 1.0}
	boneColorsIKLink         = []float32{1.0, 0.83, 0.49, 1.0}
	boneColorsFixed          = []float32{0.72, 0.32, 1.0, 1.0}
	boneColorsEffect         = []float32{0.68, 0.64, 1.0, 1.0}
	boneColorsEffectEffector = []float32{0.88, 0.84, 1.0, 0.7}
	boneColorsTranslate      = []float32{0.70, 1.0, 0.54, 1.0}
	boneColorsRotate         = []float32{0.56, 0.78, 1.0, 1.0}
	boneColorsInvisible      = []float32{0.82, 0.82, 0.82, 1.0}
	boneColorsHighlight      = []float32{1.0, 0.0, 0.0, 1.0} // ハイライト用赤色
)

// newBoneGl はボーンの先端位置データをGPU頂点用に詰める
func newBoneGl(bone *pmx.Bone) []float32 {
	// bone.Position は pmx.Bone の位置ベクトル
	p := mmath.NewGlVec3(bone.Position)
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(bone.Index()), 0, 0, 0, // デフォームボーンINDEX (4要素)
		1, 0, 0, 0, // ウェイト (4要素)
		0.0, 0.0, 0.0, 0.0, // 色 (4要素) - デバッグ描画用など
	}
}

// newTailBoneGl はボーンの Tail 部分の位置をGPU頂点用に詰める
func newTailBoneGl(bone *pmx.Bone) []float32 {
	// ボーン末端位置
	tailIndex := bone.Index()
	if bone.IsTailBone() && bone.TailIndex > 0 {
		tailIndex = bone.TailIndex
	}
	tailPos := bone.Position.Added(bone.ChildRelativePosition)

	p := mmath.NewGlVec3(tailPos)
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(tailIndex), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0.0, 0.0, 0.0, 0.0, // 色
	}
}

// getBoneDebugColor はボーンの種類(IK, 付与親, 固定軸, etc.)を見てカラーを返す
func getBoneDebugColor(bone *pmx.Bone, shared *state.SharedState, isHover bool) []float32 {
	// IK系
	if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneIk()) && bone.IsIK() {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsIK
	} else if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneIk()) &&
		len(bone.IkLinkBoneIndexes) > 0 {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsIKLink
	} else if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneIk()) &&
		len(bone.IkTargetBoneIndexes) > 0 {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsIKTarget
	}

	// 付与親系
	if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneEffector()) &&
		(bone.IsEffectorRotation() || bone.IsEffectorTranslation()) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsEffect
	} else if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneEffector()) &&
		len(bone.EffectiveBoneIndexes) > 0 {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsEffectEffector
	}

	// 固定軸
	if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneFixed()) &&
		bone.HasFixedAxis() {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsFixed
	}

	// 移動/回転
	if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneTranslate()) &&
		bone.CanTranslate() {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsTranslate
	} else if (shared.IsShowBoneAll() || shared.IsShowBoneVisible() || shared.IsShowBoneRotate()) &&
		bone.CanRotate() {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsRotate
	}

	// 非表示
	if shared.IsShowBoneAll() && !bone.IsVisible() {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsInvisible
	}

	// それ以外は色なし(透明)
	return []float32{0.0, 0.0, 0.0, 0.0}
}

// createBoneMatrixes は boneDeltas から ボーン行列テクスチャ用の float32配列を生成
func createBoneMatrixes(boneDeltas *delta.BoneDeltas) ([]float32, int, int, error) {
	numBones := boneDeltas.Length()
	// テクスチャサイズを計算
	texSize := int(math.Ceil(math.Sqrt(float64(numBones))))
	width := int(math.Ceil(float64(texSize)/4) * 4 * 4)
	height := int(math.Ceil((float64(numBones) * 4) / float64(width)))

	paddedMatrixes := make([]float32, height*width*4)
	boneDeltas.ForEach(func(index int, d *delta.BoneDelta) bool {
		var m mgl32.Mat4
		if d == nil {
			m = mgl32.Ident4()
		} else {
			m = mmath.NewGlMat4(d.FilledLocalMatrix())
		}
		copy(paddedMatrixes[index*16:], m[:])
		return true
	})

	return paddedMatrixes, width, height, nil
}

// bindBoneMatrixes : ボーン行列テクスチャをシェーダーに転送し、
// BoneMatrixTexture/BoneMatrixTextureWidth/BoneMatrixTextureHeight のユニフォームをセット
func bindBoneMatrixes(
	windowIndex int,
	paddedMatrixes []float32,
	width, height int,
	shader rendering.IShader,
	program uint32,
) {
	// テクスチャユニットをアクティブにする (windowIndexによってTEXTURE20〜TEXTURE22を割り当てる)
	switch windowIndex {
	case 0:
		gl.ActiveTexture(gl.TEXTURE20)
	case 1:
		gl.ActiveTexture(gl.TEXTURE21)
	case 2:
		gl.ActiveTexture(gl.TEXTURE22)
	}

	// ボーン用テクスチャをバインド
	gl.BindTexture(gl.TEXTURE_2D, shader.BoneTextureID())

	// 各種パラメータ設定
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// float配列をTEXTURE_2Dに貼り付け
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

	// シェーダーユニフォームにボーンテクスチャIDを送る
	boneUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderBoneMatrixTexture))
	switch windowIndex {
	case 0:
		gl.Uniform1i(boneUniform, 20)
	case 1:
		gl.Uniform1i(boneUniform, 21)
	case 2:
		gl.Uniform1i(boneUniform, 22)
	}

	modelWidthUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderBoneMatrixTextureWidth))
	gl.Uniform1i(modelWidthUniform, int32(width))

	modelHeightUniform := gl.GetUniformLocation(program, gl.Str(mgl.ShaderBoneMatrixTextureHeight))
	gl.Uniform1i(modelHeightUniform, int32(height))
}

// unbindBoneMatrixes : ボーンテクスチャのバインド解除
func unbindBoneMatrixes() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
