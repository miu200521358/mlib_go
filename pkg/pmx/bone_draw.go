//go:build windows
// +build windows

package pmx

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"github.com/miu200521358/mlib_go/pkg/mmath"
	"github.com/miu200521358/mlib_go/pkg/mview"
)

var BONE_COLORS_IK = []float32{1.0, 0.38, 0, 1.0}
var BONE_COLORS_IK_NORMAL = []float32{1.0, 0.58, 0.2, 0.7}
var BONE_COLORS_IK_TARGET = []float32{1.0, 0.57, 0.61, 1.0}
var BONE_COLORS_IK_TARGET_NORMAL = []float32{1.0, 0.77, 0.81, 0.7}
var BONE_COLORS_IK_LINK = []float32{1.0, 0.83, 0.49, 1.0}
var BONE_COLORS_IK_LINK_NORMAL = []float32{1.0, 1.0, 0.69, 0.7}
var BONE_COLORS_FIXED = []float32{0.72, 0.32, 1.0, 1.0}
var BONE_COLORS_FIXED_NORMAL = []float32{0.92, 0.52, 1.0, 0.7}
var BONE_COLORS_EFFECT = []float32{0.68, 0.64, 1.0, 1.0}
var BONE_COLORS_EFFECT_NORMAL = []float32{0.88, 0.84, 1.0, 0.7}
var BONE_COLORS_TRANSLATE = []float32{0.70, 1.0, 0.54, 1.0}
var BONE_COLORS_TRANSLATE_NORMAL = []float32{0.90, 1.0, 0.74, 0.7}
var BONE_COLORS_INVISIBLE = []float32{0.82, 0.82, 0.82, 1.0}
var BONE_COLORS_INVISIBLE_NORMAL = []float32{0.92, 0.92, 0.92, 0.7}
var BONE_COLORS_ROTATE = []float32{0.56, 0.78, 1.0, 1.0}
var BONE_COLORS_ROTATE_NORMAL = []float32{0.76, 0.98, 1.00, 0.7}

func (b *Bones) Draw(
	shader *mview.MShader,
	boneGlobalMatrixes []*mmath.MMat4,
	windowIndex int,
) {
	// ボーンをモデルメッシュの前面に描画するために深度テストを無効化
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.ALWAYS)

	// ブレンディングを有効にする
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	shader.UseBoneProgram()

	// ------------------------------

	// ボーン位置を設定
	positionVbo := make([]float32, 0)
	normalVbo := make([]float32, 0)

	for i, matrix := range boneGlobalMatrixes {
		bone := b.Get(i)

		posGl := matrix.Translation().GL()
		positionVbo = append(positionVbo, posGl[0], posGl[1], posGl[2])

		normalMatrix := matrix.Muled(bone.LocalMatrix)
		normalVbo = append(normalVbo, posGl[0], posGl[1], posGl[2])
		normalGl := normalMatrix.MulVec3(&mmath.MVec3{0, 0.6, 0}).GL()

		// ボーンの種類で色を変える
		if bone.IsIK() {
			// IKボーン
			positionVbo = append(positionVbo, BONE_COLORS_IK...)
			normalVbo = append(normalVbo, BONE_COLORS_IK_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_IK_NORMAL...)
		} else if len(bone.IkTargetBoneIndexes) > 0 {
			// IK先
			positionVbo = append(positionVbo, BONE_COLORS_IK_TARGET...)
			normalVbo = append(normalVbo, BONE_COLORS_IK_TARGET_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_IK_TARGET_NORMAL...)
		} else if len(bone.IkLinkBoneIndexes) > 0 {
			// IKリンク
			positionVbo = append(positionVbo, BONE_COLORS_IK_LINK...)
			normalVbo = append(normalVbo, BONE_COLORS_IK_LINK_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_IK_LINK_NORMAL...)
		} else if bone.HasFixedAxis() {
			// 軸制限
			positionVbo = append(positionVbo, BONE_COLORS_FIXED...)
			normalVbo = append(normalVbo, BONE_COLORS_FIXED_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_FIXED_NORMAL...)
		} else if bone.IsEffectorRotation() || bone.IsEffectorTranslation() {
			// 付与親
			positionVbo = append(positionVbo, BONE_COLORS_EFFECT...)
			normalVbo = append(normalVbo, BONE_COLORS_EFFECT_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_EFFECT_NORMAL...)
		} else if bone.CanTranslate() {
			// 移動可能
			positionVbo = append(positionVbo, BONE_COLORS_TRANSLATE...)
			normalVbo = append(normalVbo, BONE_COLORS_TRANSLATE_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_TRANSLATE_NORMAL...)
		} else if !bone.IsVisible() {
			// 非表示
			positionVbo = append(positionVbo, BONE_COLORS_INVISIBLE...)
			normalVbo = append(normalVbo, BONE_COLORS_INVISIBLE_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_INVISIBLE_NORMAL...)
		} else {
			// それ以外（回転）
			positionVbo = append(positionVbo, BONE_COLORS_ROTATE...)
			normalVbo = append(normalVbo, BONE_COLORS_ROTATE_NORMAL...)
			normalVbo = append(normalVbo, normalGl[0], normalGl[1], normalGl[2])
			normalVbo = append(normalVbo, BONE_COLORS_ROTATE_NORMAL...)
		}
	}

	positionVboGl := mview.NewVBOForBone(gl.Ptr(positionVbo), len(positionVbo))

	b.positionVao.Bind()
	positionVboGl.BindBone()
	b.positionIbo.Bind()

	gl.DrawElements(gl.LINES, b.positionIboCount, gl.UNSIGNED_INT, nil)

	b.positionIbo.Unbind()
	positionVboGl.Unbind()
	b.positionVao.Unbind()

	// ------------------------------

	normalVboGl := mview.NewVBOForBone(gl.Ptr(normalVbo), len(normalVbo))

	b.normalVao.Bind()
	normalVboGl.BindBone()
	b.normalIbo.Bind()

	gl.DrawElements(gl.LINES, b.normalIboCount, gl.UNSIGNED_INT, nil)

	b.normalIbo.Unbind()
	normalVboGl.Unbind()
	b.normalVao.Unbind()

	shader.Unuse()

	gl.Disable(gl.BLEND)
}

func (b *Bones) prepareDraw() {
	positionIbo := make([]uint32, 0, len(b.Data))
	normalIbo := make([]uint32, 0, len(b.Data))

	for i, bone := range b.GetSortedData() {
		positionIbo = append(positionIbo, uint32(bone.Index))
		positionIbo = append(positionIbo, uint32(bone.ParentIndex))

		normalIbo = append(normalIbo, uint32(i*2))
		normalIbo = append(normalIbo, uint32(i*2+1))
	}

	b.positionVao = mview.NewVAO()
	b.positionIbo = mview.NewIBO(gl.Ptr(positionIbo), len(positionIbo))
	b.positionIboCount = int32(len(positionIbo))

	b.normalVao = mview.NewVAO()
	b.normalIbo = mview.NewIBO(gl.Ptr(normalIbo), len(normalIbo))
	b.normalIboCount = int32(len(normalIbo))
}
