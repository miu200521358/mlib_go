//go:build windows
// +build windows

// 指示: miu200521358
package render

import (
	"math"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/mgl"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
)

// デバッグ表示用のボーンカラー定義
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
	boneColorsHighlight      = []float32{1.0, 0.0, 0.0, 1.0}
)

// newBoneGl はボーン位置データをGPU頂点用に詰める。
func newBoneGl(bone *model.Bone) []float32 {
	p := mgl.NewGlVec3(&bone.Position)
	return []float32{
		p[0], p[1], p[2],
		float32(bone.Index()), 0, 0, 0,
		1, 0, 0, 0,
		0.0, 0.0, 0.0, 0.0,
	}
}

// newTailBoneGl はボーン末端位置データをGPU頂点用に詰める。
func newTailBoneGl(bone *model.Bone, tailPos mmath.Vec3) []float32 {
	tailIndex := bone.Index()
	if isTailBone(bone) && bone.TailIndex >= 0 {
		tailIndex = bone.TailIndex
	}

	p := mgl.NewGlVec3(&tailPos)
	return []float32{
		p[0], p[1], p[2],
		float32(tailIndex), 0, 0, 0,
		1, 0, 0, 0,
		0.0, 0.0, 0.0, 0.0,
	}
}

// boneDebugInfo はボーンデバッグ判定用の情報を保持する。
type boneDebugInfo struct {
	ikTargets       map[int]struct{}
	ikLinks         map[int]struct{}
	effectorParents map[int]struct{}
}

// buildBoneDebugInfo はIK/付与関係の判定情報を構築する。
func buildBoneDebugInfo(bones *model.BoneCollection) boneDebugInfo {
	info := boneDebugInfo{
		ikTargets:       map[int]struct{}{},
		ikLinks:         map[int]struct{}{},
		effectorParents: map[int]struct{}{},
	}
	if bones == nil {
		return info
	}
	for _, bone := range bones.Values() {
		if bone == nil {
			continue
		}
		if bone.Ik != nil {
			if bone.Ik.BoneIndex >= 0 {
				info.ikTargets[bone.Ik.BoneIndex] = struct{}{}
			}
			for _, link := range bone.Ik.Links {
				if link.BoneIndex >= 0 {
					info.ikLinks[link.BoneIndex] = struct{}{}
				}
			}
		}
		if hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_ROTATION) || hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_TRANSLATION) {
			if bone.EffectIndex >= 0 {
				info.effectorParents[bone.EffectIndex] = struct{}{}
			}
		}
	}
	return info
}

// isIkTarget はIKターゲットか判定する。
func (info boneDebugInfo) isIkTarget(index int) bool {
	_, ok := info.ikTargets[index]
	return ok
}

// isIkLink はIKリンクか判定する。
func (info boneDebugInfo) isIkLink(index int) bool {
	_, ok := info.ikLinks[index]
	return ok
}

// isEffectorParent は付与親ボーンか判定する。
func (info boneDebugInfo) isEffectorParent(index int) bool {
	_, ok := info.effectorParents[index]
	return ok
}

// getBoneDebugColor はボーン種類に応じたカラーを返す。
func getBoneDebugColor(bone *model.Bone, shared *state.SharedState, info boneDebugInfo, isHover bool) []float32 {
	showAll := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_ALL)
	showVisible := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_VISIBLE)
	showIk := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_IK)
	showEffector := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_EFFECTOR)
	showFixed := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_FIXED)
	showRotate := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_ROTATE)
	showTranslate := shared.HasFlag(state.STATE_FLAG_SHOW_BONE_TRANSLATE)

	if (showAll || showVisible || showIk) && hasBoneFlag(bone, model.BONE_FLAG_IS_IK) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsIK
	} else if (showAll || showVisible || showIk) && info.isIkLink(bone.Index()) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsIKLink
	} else if (showAll || showVisible || showIk) && info.isIkTarget(bone.Index()) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsIKTarget
	}

	if (showAll || showVisible || showEffector) &&
		(hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_ROTATION) || hasBoneFlag(bone, model.BONE_FLAG_IS_EXTERNAL_TRANSLATION)) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsEffect
	} else if (showAll || showVisible || showEffector) && info.isEffectorParent(bone.Index()) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsEffectEffector
	}

	if (showAll || showVisible || showFixed) && hasBoneFlag(bone, model.BONE_FLAG_HAS_FIXED_AXIS) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsFixed
	}

	if (showAll || showVisible || showTranslate) && hasBoneFlag(bone, model.BONE_FLAG_CAN_TRANSLATE) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsTranslate
	} else if (showAll || showVisible || showRotate) && hasBoneFlag(bone, model.BONE_FLAG_CAN_ROTATE) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsRotate
	}

	if showAll && !hasBoneFlag(bone, model.BONE_FLAG_IS_VISIBLE) {
		if isHover {
			return boneColorsHighlight
		}
		return boneColorsInvisible
	}

	return []float32{0.0, 0.0, 0.0, 0.0}
}

// calcTailPosition はボーン末端位置を算出する。
func calcTailPosition(bone *model.Bone, bones *model.BoneCollection) mmath.Vec3 {
	if bone == nil {
		return mmath.NewVec3()
	}
	if isTailBone(bone) && bone.TailIndex >= 0 && bones != nil {
		if tail, err := bones.Get(bone.TailIndex); err == nil && tail != nil {
			return tail.Position
		}
	}
	// 末端指定が無い場合は、tail位置→子ボーン→親方向の順で補完する。
	if !isTailBone(bone) && bone.TailPosition.Length() > 0 {
		return bone.Position.Added(bone.TailPosition)
	}
	if pos, ok := findChildPosition(bone.Index(), bones); ok {
		return pos
	}
	if bones != nil && bone.ParentIndex >= 0 {
		if parent, err := bones.Get(bone.ParentIndex); err == nil && parent != nil {
			dir := bone.Position.Subed(parent.Position)
			if dir.Length() > 0 {
				return bone.Position.Added(dir)
			}
		}
	}
	return bone.Position
}

// findChildPosition は親に紐づく最初の子ボーン位置を返す。
func findChildPosition(boneIndex int, bones *model.BoneCollection) (mmath.Vec3, bool) {
	if bones == nil {
		return mmath.NewVec3(), false
	}
	for _, child := range bones.Values() {
		if child == nil {
			continue
		}
		if child.ParentIndex == boneIndex {
			return child.Position, true
		}
	}
	return mmath.NewVec3(), false
}

// createBoneMatrixes は boneDeltas から ボーン行列テクスチャ用の float32配列を生成する。
func createBoneMatrixes(boneDeltas *delta.BoneDeltas) ([]float32, int, int, error) {
	numBones := boneDeltas.Len()
	texSize := int(math.Ceil(math.Sqrt(float64(numBones))))
	width := int(math.Ceil(float64(texSize)/4) * 4 * 4)
	height := int(math.Ceil((float64(numBones) * 4) / float64(width)))

	paddedMatrixes := make([]float32, height*width*4)
	boneDeltas.ForEach(func(index int, d *delta.BoneDelta) bool {
		var m mgl32.Mat4
		if d == nil {
			m = mgl32.Ident4()
		} else {
			m = mgl.NewGlMat4(d.FilledLocalMatrix())
		}
		copy(paddedMatrixes[index*16:], m[:])
		return true
	})

	return paddedMatrixes, width, height, nil
}

// bindBoneMatrixes はボーン行列テクスチャをシェーダーに転送する。
func bindBoneMatrixes(
	windowIndex int,
	paddedMatrixes []float32,
	width, height int,
	shader graphics_api.IShader,
	program uint32,
) {
	switch windowIndex {
	case 0:
		gl.ActiveTexture(gl.TEXTURE20)
	case 1:
		gl.ActiveTexture(gl.TEXTURE21)
	case 2:
		gl.ActiveTexture(gl.TEXTURE22)
	}

	gl.BindTexture(gl.TEXTURE_2D, shader.BoneTextureID())

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

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

// unbindBoneMatrixes はボーンテクスチャのバインドを解除する。
func unbindBoneMatrixes() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// hasBoneFlag はボーンフラグの有無を判定する。
func hasBoneFlag(bone *model.Bone, flag model.BoneFlag) bool {
	if bone == nil {
		return false
	}
	return bone.BoneFlag&flag != 0
}

// isTailBone は接続先ボーン指定か判定する。
func isTailBone(bone *model.Bone) bool {
	return hasBoneFlag(bone, model.BONE_FLAG_TAIL_IS_BONE)
}
