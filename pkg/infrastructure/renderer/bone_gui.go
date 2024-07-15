//go:build windows
// +build windows

package renderer

import (
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mgl"
)

var bone_colors_ik = []float32{1.0, 0.38, 0, 1.0}
var bone_colors_ik_normal = []float32{1.0, 0.58, 0.2, 0.7}
var bone_colors_ik_target = []float32{1.0, 0.57, 0.61, 1.0}
var bone_colors_ik_target_normal = []float32{1.0, 0.77, 0.81, 0.7}
var bone_colors_ik_link = []float32{1.0, 0.83, 0.49, 1.0}
var bone_colors_ik_link_normal = []float32{1.0, 1.0, 0.69, 0.7}
var bone_colors_fixed = []float32{0.72, 0.32, 1.0, 1.0}
var bone_colors_fixed_normal = []float32{0.92, 0.52, 1.0, 0.7}
var bone_colors_effect = []float32{0.68, 0.64, 1.0, 1.0}
var bone_colors_effect_normal = []float32{0.88, 0.84, 1.0, 0.7}
var bone_colors_translate = []float32{0.70, 1.0, 0.54, 1.0}
var bone_colors_translate_normal = []float32{0.90, 1.0, 0.74, 0.7}
var bone_colors_invisible = []float32{0.82, 0.82, 0.82, 1.0}
var bone_colors_invisible_normal = []float32{0.92, 0.92, 0.92, 0.7}
var bone_colors_rotate = []float32{0.56, 0.78, 1.0, 1.0}
var bone_colors_rotate_normal = []float32{0.76, 0.98, 1.00, 0.7}

func color(b *pmx.Bone, isNormal bool) []float32 {
	// ボーンの種類で色を変える
	if b.IsIK() {
		// IKボーン
		if isNormal {
			return bone_colors_ik_normal
		}
		return bone_colors_ik
	} else if len(b.Extend.IkTargetBoneIndexes) > 0 {
		// IK先
		if isNormal {
			return bone_colors_ik_target_normal
		}
		return bone_colors_ik_target
	} else if len(b.Extend.IkLinkBoneIndexes) > 0 {
		// IKリンク
		if isNormal {
			return bone_colors_ik_link_normal
		}
		return bone_colors_ik_link
	} else if b.HasFixedAxis() {
		// 軸制限
		if isNormal {
			return bone_colors_fixed_normal
		}
		return bone_colors_fixed
	} else if b.IsEffectorRotation() || b.IsEffectorTranslation() {
		// 付与親
		if isNormal {
			return bone_colors_effect_normal
		}
		return bone_colors_effect
	} else if b.CanTranslate() {
		// 移動可能
		if isNormal {
			return bone_colors_translate_normal
		}
		return bone_colors_translate
	} else if !b.IsVisible() {
		// 非表示
		if isNormal {
			return bone_colors_invisible_normal
		}
		return bone_colors_invisible
	}

	// それ以外（回転）
	if isNormal {
		return bone_colors_rotate_normal
	}
	return bone_colors_rotate
}

func BoneGL(b *pmx.Bone) []float32 {
	p := mgl.NewGlVec3(b.Position)
	c := color(b, false)
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(b.BoneFlag), 0.0, 0.0, // 法線
		float32(0), float32(0), // UV
		float32(0), float32(0), // 追加UV
		float32(0),                // エッジ倍率
		float32(b.Index), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0,       // SDEFであるか否か
		0, 0, 0, // SDEF-C
		0, 0, 0, // SDEF-R0
		0, 0, 0, // SDEF-R1
		0.0, 0.0, 0.0, // 頂点モーフ
		0.0, 0.0, 0.0, 0.0, // UVモーフ
		c[0], c[1], c[2], c[3], // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}

func ParentGL(b *pmx.Bone) []float32 {
	p := mgl.NewGlVec3(b.Position.Subed(b.Extend.ParentRelativePosition))
	c := color(b, false)
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(b.BoneFlag), 0.0, 0.0, // 法線
		float32(0), float32(0), // UV
		float32(0), float32(0), // 追加UV
		float32(0),                // エッジ倍率
		float32(b.Index), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0,       // SDEFであるか否か
		0, 0, 0, // SDEF-C
		0, 0, 0, // SDEF-R0
		0, 0, 0, // SDEF-R1
		0.0, 0.0, 0.0, // 頂点モーフ
		0.0, 0.0, 0.0, 0.0, // UVモーフ
		c[0], c[1], c[2], c[3], // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}

func TailGL(b *pmx.Bone) []float32 {
	p := mgl.NewGlVec3(b.Position.Added(b.Extend.ChildRelativePosition))
	c := color(b, false)
	return []float32{
		p[0], p[1], p[2], // 位置
		float32(b.BoneFlag), 0.0, 0.0, // 法線
		float32(0), float32(0), // UV
		float32(0), float32(0), // 追加UV
		float32(0),                // エッジ倍率
		float32(b.Index), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0,       // SDEFであるか否か
		0, 0, 0, // SDEF-C
		0, 0, 0, // SDEF-R0
		0, 0, 0, // SDEF-R1
		0.0, 0.0, 0.0, // 頂点モーフ
		0.0, 0.0, 0.0, 0.0, // UVモーフ
		c[0], c[1], c[2], c[3], // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}

func DeltaGL(b *pmx.Bone, isDrawBones map[pmx.BoneFlag]bool) []float32 {
	c := color(b, false)

	ikAlpha := float32(1.0)
	fixedAlpha := float32(1.0)
	effectorRotateAlpha := float32(1.0)
	effectorTranslateAlpha := float32(1.0)
	rotateAlpha := float32(1.0)
	translateAlpha := float32(1.0)
	visibleAlpha := float32(1.0)
	// IK
	if (!isDrawBones[pmx.BONE_FLAG_IS_IK] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_IS_IK] && !(b.IsIK() || len(b.Extend.IkLinkBoneIndexes) > 0 || len(b.Extend.IkTargetBoneIndexes) > 0)) {
		ikAlpha = float32(0.0)
	}
	// 付与親回転
	if (!isDrawBones[pmx.BONE_FLAG_IS_EXTERNAL_ROTATION] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_IS_EXTERNAL_ROTATION] && !(b.IsEffectorRotation() || len(b.Extend.EffectiveBoneIndexes) > 0)) {
		effectorRotateAlpha = float32(0.0)
	}
	// 付与親移動
	if (!isDrawBones[pmx.BONE_FLAG_IS_EXTERNAL_TRANSLATION] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_IS_EXTERNAL_TRANSLATION] && !(b.IsEffectorTranslation() || len(b.Extend.EffectiveBoneIndexes) > 0)) {
		effectorTranslateAlpha = float32(0.0)
	}
	// 軸固定
	if (!isDrawBones[pmx.BONE_FLAG_HAS_FIXED_AXIS] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_HAS_FIXED_AXIS] && !b.HasFixedAxis()) {
		fixedAlpha = float32(0.0)
	}
	// 回転
	if (!isDrawBones[pmx.BONE_FLAG_CAN_ROTATE] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_CAN_ROTATE] && !b.CanRotate()) {
		rotateAlpha = float32(0.0)
	}
	// 移動
	if (!isDrawBones[pmx.BONE_FLAG_CAN_TRANSLATE] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_CAN_TRANSLATE] && !b.CanTranslate()) {
		translateAlpha = float32(0.0)
	}
	// 表示
	if (!isDrawBones[pmx.BONE_FLAG_IS_VISIBLE] && !isDrawBones[pmx.BONE_FLAG_NONE]) ||
		(isDrawBones[pmx.BONE_FLAG_IS_VISIBLE] && !b.IsVisible()) {
		visibleAlpha = float32(0.0)
	}

	// それぞれのボーン種別による透明度最大値を採用
	alpha := max(ikAlpha, fixedAlpha, effectorRotateAlpha, effectorTranslateAlpha,
		rotateAlpha, translateAlpha, visibleAlpha)
	return []float32{
		0.0, 0.0, 0.0, // 頂点モーフ
		0.0, 0.0, 0.0, 0.0, // UVモーフ
		c[0], c[1], c[2], c[3] * alpha, // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}

func NormalGL(b *pmx.Bone) []float32 {
	p := mgl.NewGlVec3(b.Extend.LocalMatrix.MulVec3(&mmath.MVec3{0, 0.6, 0}))
	c := color(b, true)
	return []float32{
		p[0], p[1], p[2], // 位置
		0.0, 0.0, 0.0, // 法線
		float32(0), float32(0), // UV
		float32(0), float32(0), // 追加UV
		float32(0),                // エッジ倍率
		float32(b.Index), 0, 0, 0, // デフォームボーンINDEX
		1, 0, 0, 0, // デフォームボーンウェイト
		0,       // SDEFであるか否か
		0, 0, 0, // SDEF-C
		0, 0, 0, // SDEF-R0
		0, 0, 0, // SDEF-R1
		0.0, 0.0, 0.0, // 頂点モーフ
		0.0, 0.0, 0.0, 0.0, // UVモーフ
		c[0], c[1], c[2], c[3], // 追加UV1モーフ
		0.0, 0.0, 0.0, // 変形後頂点モーフ
	}
}
