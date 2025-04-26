package repository

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
)

func (rep *PmxRepository) CreateSticky(
	model *pmx.PmxModel,
	insertBoneFunc func(vertices *pmx.Vertices, bones *pmx.Bones) error,
	insertDebugFunc func(bones *pmx.Bones, displaySlots *pmx.DisplaySlots) error,
) error {

	// 設定ボーンを追加
	if insertBoneFunc != nil {
		if err := insertBoneFunc(model.Vertices, model.Bones); err != nil {
			return err
		}
	}

	// デバッグボーンを追加
	if insertDebugFunc != nil {
		if err := insertDebugFunc(model.Bones, model.DisplaySlots); err != nil {
			return err
		}
	}

	// ボーン以外を削除
	model.SetName(fmt.Sprintf("棒人間_%s", model.Name()))
	model.Vertices = pmx.NewVertices(0)
	model.Faces = pmx.NewFaces(0)
	model.Materials = pmx.NewMaterials(0)
	model.RigidBodies = pmx.NewRigidBodies(0)
	model.Joints = pmx.NewJoints(0)

	model.Morphs.ForEach(func(index int, morph *pmx.Morph) bool {
		if !(morph.MorphType == pmx.MORPH_TYPE_GROUP || morph.MorphType == pmx.MORPH_TYPE_BONE) {
			// グループモーフとボーンモーフ以外は削除
			morph.Offsets = make([]pmx.IMorphOffset, 0)
			model.Morphs.Update(morph)
		}
		return true
	})

	// 棒人間情報を追加
	if err := rep.createStickBones(model); err != nil {
		return err
	}

	return nil
}

func (rep *PmxRepository) createStickBones(model *pmx.PmxModel) error {
	boneMaterial := rep.createStickMaterial()

	// 棒人間材質
	model.Materials.Append(boneMaterial)

	rootBone, err := model.Bones.GetRoot()
	if err != nil {
		return err
	}

	model.Bones.ForEach(func(index int, bone *pmx.Bone) bool {
		if (bone.Config() == nil && !bone.IsVisible()) ||
			(bone.RigidBody != nil && bone.RigidBody.PhysicsType != pmx.PHYSICS_TYPE_STATIC) ||
			bone.Name() == pmx.LEG_IK_PARENT.Left() || bone.Name() == pmx.LEG_IK_PARENT.Right() ||
			bone.Name() == pmx.LEG_IK.Left() || bone.Name() == pmx.LEG_IK.Right() {
			// 設定外の非表示ボーン、物理ボーン、IKはスルー
			// 操作と表示のフラグをOFFにする
			bone.BoneFlag &= ^pmx.BONE_FLAG_CAN_MANIPULATE
			bone.BoneFlag &= ^pmx.BONE_FLAG_IS_VISIBLE
			model.Bones.Update(bone)

			return true
		}

		// 原点に繋がるシステムボーンの場合、原点の直方体を作成
		if bone.IsSystem && bone.ParentBone != nil && bone.ParentBone.Index() == rootBone.Index() {
			rep.createStickRoot(model, bone)
			return true
		}

		if bone.ParentIndex < 0 {
			return true
		}

		parentBone, err := model.Bones.Get(bone.ParentIndex)
		if err != nil {
			return true
		}

		if parentBone == nil || parentBone.Index() == rootBone.Index() {
			return true
		}

		// ボーンの位置に合わせた直方体を作成
		rep.createStickBone(model, bone, parentBone, boneMaterial)

		return true
	})

	return nil
}

func (rep *PmxRepository) createStickRoot(model *pmx.PmxModel, bone *pmx.Bone) {
	// 直方体の頂点を作成
	v1 := pmx.NewVertex()
	v1.Position = &mmath.MVec3{X: -0.15, Y: -0.15, Z: -0.15}
	v1.Position.Add(bone.Position)
	v1.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v1)

	v2 := pmx.NewVertex()
	v2.Position = &mmath.MVec3{X: 0.15, Y: -0.15, Z: -0.15}
	v2.Position.Add(bone.Position)
	v2.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v2)

	v3 := pmx.NewVertex()
	v3.Position = &mmath.MVec3{X: -0.15, Y: 0.15, Z: -0.15}
	v3.Position.Add(bone.Position)
	v3.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v3)

	v4 := pmx.NewVertex()
	v4.Position = &mmath.MVec3{X: 0.15, Y: 0.15, Z: -0.15}
	v4.Position.Add(bone.Position)
	v4.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v4)

	v5 := pmx.NewVertex()
	v5.Position = &mmath.MVec3{X: -0.15, Y: -0.15, Z: 0.15}
	v5.Position.Add(bone.Position)
	v5.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v5)

	v6 := pmx.NewVertex()
	v6.Position = &mmath.MVec3{X: 0.15, Y: -0.15, Z: 0.15}
	v6.Position.Add(bone.Position)
	v6.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v6)

	v7 := pmx.NewVertex()
	v7.Position = &mmath.MVec3{X: -0.15, Y: 0.15, Z: 0.15}
	v7.Position.Add(bone.Position)
	v7.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v7)

	v8 := pmx.NewVertex()
	v8.Position = &mmath.MVec3{X: 0.15, Y: 0.15, Z: 0.15}
	v8.Position.Add(bone.Position)
	v8.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v8)

	// 立方体になるように面を指定
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v1.Index(), v2.Index(), v3.Index()}}) // 前左
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v2.Index(), v3.Index(), v4.Index()}}) // 前右
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v3.Index(), v8.Index(), v4.Index()}}) // 上前
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v3.Index(), v7.Index(), v8.Index()}}) // 上後
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v5.Index(), v3.Index(), v1.Index()}}) // 左下
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v5.Index(), v7.Index(), v3.Index()}}) // 左上
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v6.Index(), v7.Index(), v5.Index()}}) // 後右
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v6.Index(), v7.Index(), v8.Index()}}) // 後左
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v2.Index(), v5.Index(), v1.Index()}}) // 下前
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v2.Index(), v6.Index(), v5.Index()}}) // 下後
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v2.Index(), v8.Index(), v6.Index()}}) // 右前
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v2.Index(), v4.Index(), v8.Index()}}) // 右後

	material := pmx.NewMaterial()
	material.SetName(bone.Name())
	material.SetEnglishName(bone.EnglishName())

	material.Diffuse = randDiffuse()
	material.Specular = &mmath.MVec4{X: 0, Y: 0, Z: 0, W: 1}
	material.Ambient = &mmath.MVec3{X: 0.3, Y: 0.3, Z: 0.3}
	material.DrawFlag = pmx.DRAW_FLAG_DOUBLE_SIDED_DRAWING
	material.VerticesCount = 12 * 3

	model.Materials.Append(material)
}

func randDiffuse() *mmath.MVec4 {
	// ランダムな色相を生成（0.0～1.0）
	h := rand.Float64()
	s := min(1.0, rand.Float64()*3)
	v := min(1.0, rand.Float64()*3)

	// HSV を RGB に変換
	r, g, b := hsvToRgb(h, s, v)

	return &mmath.MVec4{X: r, Y: g, Z: b, W: 1}
}

// hsvToRgb は HSV 値 (h, s, v) を RGB 値 (r, g, b) に変換する関数です。
// h は [0,1] の範囲で表し、s と v も同じく [0,1] の値とします。
func hsvToRgb(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v // 彩度が0ならグレースケール
	}

	h = h * 6          // h の範囲を [0,6) に変換
	i := math.Floor(h) // セクション番号を取得
	f := h - i         // セクション内の比率
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch int(i) % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	}

	return 0, 0, 0 // 通常はここに到達しない
}

func (rep *PmxRepository) createStickBone(model *pmx.PmxModel, bone, parentBone *pmx.Bone, boneMaterial *pmx.Material) {
	// 直方体の頂点を作成
	v1 := pmx.NewVertex()
	v1.Position = parentBone.Position.Copy()
	v1.Deform = pmx.NewBdef1(parentBone.Index())
	model.Vertices.Append(v1)

	v2 := pmx.NewVertex()
	v2.Position = bone.Position.Copy()
	v2.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v2)

	v3 := pmx.NewVertex()
	v3.Position = parentBone.Position.Copy()
	v3.Position.X += 0.05
	v3.Deform = pmx.NewBdef1(parentBone.Index())
	model.Vertices.Append(v3)

	v4 := pmx.NewVertex()
	v4.Position = bone.Position.Copy()
	v4.Position.X += 0.05
	v4.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v4)

	v5 := pmx.NewVertex()
	v5.Position = parentBone.Position.Copy()
	v5.Position.X += 0.05
	v5.Position.Y += 0.05
	v5.Deform = pmx.NewBdef1(parentBone.Index())
	model.Vertices.Append(v5)

	v6 := pmx.NewVertex()
	v6.Position = bone.Position.Copy()
	v6.Position.X += 0.05
	v6.Position.Y += 0.05
	v6.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v6)

	v7 := pmx.NewVertex()
	v7.Position = parentBone.Position.Copy()
	v7.Position.Z += 0.05
	v7.Deform = pmx.NewBdef1(parentBone.Index())
	model.Vertices.Append(v7)

	v8 := pmx.NewVertex()
	v8.Position = bone.Position.Copy()
	v8.Position.Z += 0.05
	v8.Deform = pmx.NewBdef1(bone.Index())
	model.Vertices.Append(v8)

	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v1.Index(), v2.Index(), v3.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v3.Index(), v2.Index(), v4.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v3.Index(), v4.Index(), v5.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v5.Index(), v4.Index(), v6.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v5.Index(), v6.Index(), v7.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v7.Index(), v6.Index(), v8.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v7.Index(), v8.Index(), v1.Index()}})
	model.Faces.Append(&pmx.Face{VertexIndexes: [3]int{v1.Index(), v8.Index(), v2.Index()}})

	boneMaterial.VerticesCount += 8 * 3
}

func (rep *PmxRepository) createStickMaterial() *pmx.Material {
	// マテリアルを追加
	material := pmx.NewMaterial()
	material.SetName("棒人間")
	material.SetEnglishName("Stick")
	material.Diffuse = &mmath.MVec4{X: 0, Y: 0, Z: 1, W: 1}
	material.Specular = &mmath.MVec4{X: 0, Y: 0, Z: 0, W: 1}
	material.Ambient = &mmath.MVec3{X: 0.3, Y: 0.3, Z: 0.3}
	material.DrawFlag = pmx.DRAW_FLAG_DOUBLE_SIDED_DRAWING

	return material
}
