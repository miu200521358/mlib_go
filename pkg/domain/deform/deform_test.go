// 指示: miu200521358
package deform

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/delta"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

func vec3(x, y, z float64) mmath.Vec3 {
	v := mmath.NewVec3()
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// newTestModel はテスト用の最小モデルを生成する。
func newTestModel() *model.PmxModel {
	m := model.NewPmxModel()
	vertex := &model.Vertex{Position: vec3(1, 2, 3), Normal: vec3(0, 1, 0), Deform: model.NewBdef1(0)}
	m.Vertices.Append(vertex)
	bone := &model.Bone{Position: mmath.NewVec3()}
	bone.SetName("bone")
	m.Bones.Append(bone)
	material := model.NewMaterial()
	material.SetName("mat")
	m.Materials.Append(material)
	return m
}

// TestRecomputeSdef はSDEF再計算を確認する。
func TestRecomputeSdef(t *testing.T) {
	c, r0, r1 := RecomputeSdef(
		vec3(0, 0, 0),
		vec3(2, 0, 0),
		vec3(1, 0, 0),
	)
	if !c.NearEquals(vec3(1, 0, 0), 1e-6) {
		t.Fatalf("sdef c mismatch: %v", c)
	}
	if !r0.NearEquals(vec3(0.5, 0, 0), 1e-6) {
		t.Fatalf("sdef r0 mismatch: %v", r0)
	}
	if !r1.NearEquals(vec3(1.5, 0, 0), 1e-6) {
		t.Fatalf("sdef r1 mismatch: %v", r1)
	}

	invalid, _, _ := RecomputeSdef(
		vec3(math.NaN(), 0, 0),
		vec3(0, 0, 0),
		vec3(1, 0, 0),
	)
	if !invalid.NearEquals(vec3(1, 0, 0), 1e-6) {
		t.Fatalf("fallback mismatch: %v", invalid)
	}
}

// TestApplySkinningBdef1 はBDEF1スキニングを確認する。
func TestApplySkinningBdef1(t *testing.T) {
	m := newTestModel()
	boneDeltas := delta.NewBoneDeltas(m.Bones)
	boneDelta := delta.NewBoneDelta(m.Bones.Values()[0], 0)
	boneDeltas.Update(boneDelta)

	ApplySkinning(m.Vertices, boneDeltas, nil)
	vertex, _ := m.Vertices.Get(0)
	if !vertex.Position.NearEquals(vec3(1, 2, 3), 1e-6) {
		t.Fatalf("position mismatch: %v", vertex.Position)
	}
	if !vertex.Normal.NearEquals(vec3(0, 1, 0), 1e-6) {
		t.Fatalf("normal mismatch: %v", vertex.Normal)
	}
}

// TestComputeMorphDeltasGroupMaterial はグループ/材質モーフを確認する。
func TestComputeMorphDeltasGroupMaterial(t *testing.T) {
	m := newTestModel()
	vertexMorph := &model.Morph{MorphType: model.MORPH_TYPE_VERTEX}
	vertexMorph.SetName("v")
	vertexMorph.Offsets = []model.MorphOffset{&model.VertexMorphOffset{VertexIndex: 0, Position: vec3(0, 1, 0)}}
	m.Morphs.Append(vertexMorph)

	groupMorph := &model.Morph{MorphType: model.MORPH_TYPE_GROUP}
	groupMorph.SetName("g")
	groupMorph.Offsets = []model.MorphOffset{&model.GroupMorphOffset{MorphIndex: 0, MorphFactor: 0.5}}
	m.Morphs.Append(groupMorph)

	materialMorph := &model.Morph{MorphType: model.MORPH_TYPE_MATERIAL}
	materialMorph.SetName("m")
	materialMorph.Offsets = []model.MorphOffset{
		&model.MaterialMorphOffset{MaterialIndex: 0, CalcMode: model.CALC_MODE_MULTIPLICATION, Diffuse: mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2}},
		&model.MaterialMorphOffset{MaterialIndex: 0, CalcMode: model.CALC_MODE_ADDITION, Diffuse: mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}},
	}
	m.Morphs.Append(materialMorph)

	motionData := motion.NewVmdMotion("")
	motionData.MorphFrames.Get("g").Append(&motion.MorphFrame{BaseFrame: motion.NewBaseFrame(0), Ratio: 2})
	motionData.MorphFrames.Get("m").Append(&motion.MorphFrame{BaseFrame: motion.NewBaseFrame(0), Ratio: 1})

	deltas := ComputeMorphDeltas(m, motionData, 0, nil)
	vDelta := deltas.Vertices().Get(0)
	if vDelta == nil || vDelta.Position == nil || !vDelta.Position.NearEquals(vec3(0, 1, 0), 1e-6) {
		t.Fatalf("vertex delta mismatch: %v", vDelta)
	}
	mDelta := deltas.Materials().Get(0)
	if mDelta == nil || !mDelta.MulMaterial.Diffuse.NearEquals(mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2}, 1e-6) {
		t.Fatalf("material mul mismatch: %v", mDelta)
	}
	if !mDelta.AddMaterial.Diffuse.NearEquals(mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}, 1e-6) {
		t.Fatalf("material add mismatch: %v", mDelta)
	}
}

// TestApplyBoneMatrices は行列合成を確認する。
func TestApplyBoneMatrices(t *testing.T) {
	m := model.NewPmxModel()
	bone := &model.Bone{Position: vec3(1, 0, 0)}
	bone.ParentIndex = -1
	bone.SetName("bone")
	m.Bones.Append(bone)
	boneDeltas := delta.NewBoneDeltas(m.Bones)
	d := delta.NewBoneDelta(bone, 0)
	pos := vec3(2, 0, 0)
	d.FramePosition = &pos
	boneDeltas.Update(d)

	ApplyBoneMatrices(m, boneDeltas)
	updated := boneDeltas.Get(0)
	if updated == nil || updated.UnitMatrix == nil {
		t.Fatalf("missing unit matrix")
	}
	if !updated.UnitMatrix.Translation().NearEquals(vec3(3, 0, 0), 1e-6) {
		t.Fatalf("unit translation mismatch: %v", updated.UnitMatrix.Translation())
	}
	if updated.LocalMatrix == nil || !updated.LocalMatrix.Translation().NearEquals(vec3(2, 0, 0), 1e-6) {
		t.Fatalf("local translation mismatch: %v", updated.LocalMatrix)
	}
}

// TestIkAxisValue はIK角度制限の反射挙動を確認する。
func TestIkAxisValue(t *testing.T) {
	minLimit := -0.5
	maxLimit := 0.5
	v := getIkAxisValue(-1.0, minLimit, maxLimit, 0, 4)
	if math.Abs(v) > 1e-6 {
		t.Fatalf("expected reflected value")
	}
	v = getIkAxisValue(-1.0, minLimit, maxLimit, 3, 4)
	if math.Abs(v-minLimit) > 1e-6 {
		t.Fatalf("expected clamped value")
	}
}

// TestGetLinkAxis は単軸制限の軸選択を確認する。
func TestGetLinkAxis(t *testing.T) {
	link := model.IkLink{AngleLimit: true, MinAngleLimit: vec3(-1, 0, 0), MaxAngleLimit: vec3(1, 0, 0)}
	axis := getLinkAxis(link, vec3(0, 1, 0), vec3(0, 0, 1))
	if !axis.NearEquals(mmath.UNIT_X_VEC3, 1e-6) {
		t.Fatalf("axis mismatch: %v", axis)
	}
}

// TestMaterialMorphApply は材質差分適用を確認する。
func TestMaterialMorphApply(t *testing.T) {
	m := newTestModel()
	mat := m.Materials.Values()[0]
	mat.Diffuse = mmath.ONE_VEC4
	deltas := delta.NewMorphDeltas(m.Vertices, m.Materials, m.Bones)
	mDelta := delta.NewMaterialMorphDelta(m.Materials.Values()[0])
	mDelta.MulMaterial.Diffuse = mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2}
	mDelta.AddMaterial.Diffuse = mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}
	deltas.Materials().Update(mDelta)

	ApplyMorphDeltas(m, deltas)
	mat, _ = m.Materials.Get(0)
	if !mat.Diffuse.NearEquals(mmath.Vec4{X: 3, Y: 3, Z: 3, W: 3}, 1e-6) {
		t.Fatalf("material apply mismatch: %v", mat.Diffuse)
	}
}

// TestVertexMorphApply は頂点差分適用を確認する。
func TestVertexMorphApply(t *testing.T) {
	m := newTestModel()
	deltas := delta.NewMorphDeltas(m.Vertices, m.Materials, m.Bones)
	pos := vec3(1, 0, 0)
	vDelta := delta.NewVertexMorphDelta(0)
	vDelta.Position = &pos
	deltas.Vertices().Update(vDelta)

	ApplyMorphDeltas(m, deltas)
	vertex, _ := m.Vertices.Get(0)
	if !vertex.Position.NearEquals(vec3(2, 2, 3), 1e-6) {
		t.Fatalf("vertex apply mismatch: %v", vertex.Position)
	}
}

// TestApplyMorphUv1 は追加UV1の適用を確認する。
func TestApplyMorphUv1(t *testing.T) {
	m := newTestModel()
	vertex, _ := m.Vertices.Get(0)
	vertex.ExtendedUvs = []mmath.Vec4{{}}
	deltas := delta.NewMorphDeltas(m.Vertices, m.Materials, m.Bones)
	uv := mmath.Vec2{X: 0.5}
	vDelta := delta.NewVertexMorphDelta(0)
	vDelta.Uv1 = &uv
	deltas.Vertices().Update(vDelta)

	ApplyMorphDeltas(m, deltas)
	if math.Abs(vertex.ExtendedUvs[0].X-0.5) > 1e-6 {
		t.Fatalf("uv1 mismatch: %v", vertex.ExtendedUvs[0])
	}
}

// TestApplySkinningAfterVertex は後頂点差分を確認する。
func TestApplySkinningAfterVertex(t *testing.T) {
	m := newTestModel()
	boneDeltas := delta.NewBoneDeltas(m.Bones)
	boneDelta := delta.NewBoneDelta(m.Bones.Values()[0], 0)
	boneDeltas.Update(boneDelta)

	deltas := delta.NewMorphDeltas(m.Vertices, m.Materials, m.Bones)
	after := vec3(0, 0, 1)
	vDelta := delta.NewVertexMorphDelta(0)
	vDelta.AfterPosition = &after
	deltas.Vertices().Update(vDelta)

	ApplySkinning(m.Vertices, boneDeltas, deltas)
	vertex, _ := m.Vertices.Get(0)
	if !vertex.Position.NearEquals(vec3(1, 2, 4), 1e-6) {
		t.Fatalf("after position mismatch: %v", vertex.Position)
	}
}

// TestApplyMorphMaterialCollection は材質コレクションのno-opを確認する。
func TestApplyMorphMaterialCollection(t *testing.T) {
	materials := collection.NewNamedCollection[*model.Material](0)
	m := model.NewPmxModel()
	m.Materials = materials
	deltas := delta.NewMorphDeltas(m.Vertices, m.Materials, m.Bones)
	ApplyMorphDeltas(m, deltas)
}
