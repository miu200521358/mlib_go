// 指示: miu200521358
package delta

import (
	"math"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/model/collection"
	"github.com/miu200521358/mlib_go/pkg/shared/contracts/mtime"
)

// newTestBone はテスト用ボーンを生成する。
func newTestBone(index int, name string) *model.Bone {
	bone := &model.Bone{}
	bone.SetIndex(index)
	bone.SetName(name)
	return bone
}

func vec3(x, y, z float64) mmath.Vec3 {
	v := mmath.NewVec3()
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// TestBoneDeltaTotals は総合値の計算を確認する。
func TestBoneDeltaTotals(t *testing.T) {
	bone := newTestBone(0, "bone")
	bone.BoneFlag = model.BONE_FLAG_HAS_FIXED_AXIS
	bone.FixedAxis = vec3(1, 0, 0)

	d := NewBoneDelta(bone, 0)
	rot := mmath.NewQuaternionFromAxisAngles(mmath.UNIT_Y_VEC3, math.Pi/2)
	d.FrameRotation = &rot
	pos := vec3(1, 2, 3)
	morphPos := vec3(2, 0, 0)
	d.FramePosition = &pos
	d.FrameMorphPosition = &morphPos
	scale := vec3(2, 2, 2)
	morphScale := vec3(0.5, 0.5, 0.5)
	d.FrameScale = &scale
	d.FrameMorphScale = &morphScale

	totalRot := d.FilledTotalRotation()
	axis, _ := totalRot.ToAxisAngle()
	if math.Abs(axis.Normalized().Dot(mmath.UNIT_X_VEC3)) < 0.999 {
		t.Fatalf("fixed axis not applied")
	}

	totalPos := d.FilledTotalPosition()
	if !totalPos.NearEquals(vec3(3, 2, 3), 1e-6) {
		t.Fatalf("total position mismatch: %v", totalPos)
	}

	totalScale := d.FilledTotalScale()
	if !totalScale.NearEquals(mmath.ONE_VEC3, 1e-6) {
		t.Fatalf("total scale mismatch: %v", totalScale)
	}
}

// TestNewBoneDeltaByGlobalMatrix はグローバル行列由来の差分を確認する。
func TestNewBoneDeltaByGlobalMatrix(t *testing.T) {
	parent := newTestBone(0, "parent")
	parent.Position = vec3(10, 0, 0)
	parentDelta := NewBoneDelta(parent, 0)
	parentGlobal := parent.Position.ToMat4()
	parentDelta.GlobalMatrix = &parentGlobal

	child := newTestBone(1, "child")
	child.Position = vec3(12, 0, 0)
	global := child.Position.ToMat4()

	delta := NewBoneDeltaByGlobalMatrix(child, mtime.Frame(1), global, parentDelta)
	if delta == nil {
		t.Fatalf("expected delta")
	}
	expectedUnit := parentGlobal.Inverted().Muled(global)
	if !delta.FilledUnitMatrix().NearEquals(expectedUnit, 1e-6) {
		t.Fatalf("unit matrix mismatch: %v", delta.FilledUnitMatrix())
	}
	expectedLocal := global.Muled(child.Position.Negated().ToMat4())
	if !delta.FilledLocalMatrix().NearEquals(expectedLocal, 1e-6) {
		t.Fatalf("local matrix mismatch: %v", delta.FilledLocalMatrix())
	}
	parentRelative := child.Position.Subed(parent.Position)
	expectedFramePos := expectedUnit.Translation().Subed(parentRelative)
	if delta.FramePosition == nil || !delta.FramePosition.NearEquals(expectedFramePos, 1e-6) {
		t.Fatalf("frame position mismatch: %v", delta.FramePosition)
	}
	if delta.FrameRotation == nil || !delta.FrameRotation.NearEquals(mmath.NewQuaternion(), 1e-6) {
		t.Fatalf("frame rotation mismatch: %v", delta.FrameRotation)
	}
}

// TestBoneDeltasLookup は検索系の動作を確認する。
func TestBoneDeltasLookup(t *testing.T) {
	bones := model.NewBoneCollection(0)
	bone := newTestBone(-1, "center")
	bones.Append(bone)
	deltas := NewBoneDeltas(bones)
	if deltas.GetByName("center") != nil {
		t.Fatalf("expected nil before update")
	}
	delta := NewBoneDelta(bone, 0)
	deltas.Update(delta)
	if deltas.Get(bone.Index()) == nil {
		t.Fatalf("expected delta")
	}
	if deltas.GetByName("center") == nil {
		t.Fatalf("expected delta by name")
	}
	if !deltas.Contains(bone.Index()) {
		t.Fatalf("expected contains")
	}
}

// TestVertexMorphDeltaIsZero はゼロ判定を確認する。
func TestVertexMorphDeltaIsZero(t *testing.T) {
	d := NewVertexMorphDelta(0)
	if !d.IsZero() {
		t.Fatalf("expected zero")
	}
	pos := vec3(1, 0, 0)
	d.Position = &pos
	if d.IsZero() {
		t.Fatalf("expected non-zero")
	}
}

// TestBoneMorphDeltaCopy はコピーの独立性を確認する。
func TestBoneMorphDeltaCopy(t *testing.T) {
	d := NewBoneMorphDelta(1)
	pos := vec3(1, 0, 0)
	d.FramePosition = &pos
	rot := mmath.NewQuaternionFromAxisAngles(mmath.UNIT_Y_VEC3, math.Pi/2)
	d.FrameRotation = &rot

	cp := d.Copy()
	pos.X = 2
	if cp.FramePosition == nil || cp.FramePosition.X != 1 {
		t.Fatalf("expected copied position")
	}
}

// TestMaterialMorphDeltaCalc は材質モーフ計算を確認する。
func TestMaterialMorphDeltaCalc(t *testing.T) {
	mat := model.NewMaterial()
	mat.SetIndex(0)
	mat.SetName("mat")
	delta := NewMaterialMorphDelta(mat)
	offset := &model.MaterialMorphOffset{
		Diffuse:             mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2},
		Specular:            mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1},
		Ambient:             vec3(1, 2, 3),
		Edge:                mmath.Vec4{X: 1, Y: 2, Z: 3, W: 4},
		EdgeSize:            2,
		TextureFactor:       mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2},
		SphereTextureFactor: mmath.Vec4{X: 3, Y: 3, Z: 3, W: 3},
		ToonTextureFactor:   mmath.Vec4{X: 4, Y: 4, Z: 4, W: 4},
	}
	delta.Mul(offset, 0.5)
	if !delta.MulMaterial.Diffuse.NearEquals(mmath.Vec4{X: 1.5, Y: 1.5, Z: 1.5, W: 1.5}, 1e-6) {
		t.Fatalf("mul diffuse mismatch: %v", delta.MulMaterial.Diffuse)
	}
	if !delta.MulMaterial.TextureFactor.NearEquals(mmath.Vec4{X: 1.5, Y: 1.5, Z: 1.5, W: 1.5}, 1e-6) {
		t.Fatalf("mul texture factor mismatch: %v", delta.MulMaterial.TextureFactor)
	}
	if !delta.MulMaterial.SphereTextureFactor.NearEquals(mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2}, 1e-6) {
		t.Fatalf("mul sphere factor mismatch: %v", delta.MulMaterial.SphereTextureFactor)
	}
	if !delta.MulMaterial.ToonTextureFactor.NearEquals(mmath.Vec4{X: 2.5, Y: 2.5, Z: 2.5, W: 2.5}, 1e-6) {
		t.Fatalf("mul toon factor mismatch: %v", delta.MulMaterial.ToonTextureFactor)
	}
	delta.Add(offset, 0.5)
	if !delta.AddMaterial.Diffuse.NearEquals(mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}, 1e-6) {
		t.Fatalf("add diffuse mismatch: %v", delta.AddMaterial.Diffuse)
	}
	if !delta.AddMaterial.TextureFactor.NearEquals(mmath.Vec4{X: 1, Y: 1, Z: 1, W: 1}, 1e-6) {
		t.Fatalf("add texture factor mismatch: %v", delta.AddMaterial.TextureFactor)
	}
	if !delta.AddMaterial.SphereTextureFactor.NearEquals(mmath.Vec4{X: 1.5, Y: 1.5, Z: 1.5, W: 1.5}, 1e-6) {
		t.Fatalf("add sphere factor mismatch: %v", delta.AddMaterial.SphereTextureFactor)
	}
	if !delta.AddMaterial.ToonTextureFactor.NearEquals(mmath.Vec4{X: 2, Y: 2, Z: 2, W: 2}, 1e-6) {
		t.Fatalf("add toon factor mismatch: %v", delta.AddMaterial.ToonTextureFactor)
	}
}

// TestVmdDeltasAccessors はVmdDeltasのアクセサを確認する。
func TestVmdDeltasAccessors(t *testing.T) {
	bones := model.NewBoneCollection(0)
	v := NewVmdDeltas(1, bones, "model", "motion")
	if v.Frame() != 1 {
		t.Fatalf("frame mismatch")
	}
	v.SetFrame(2)
	if v.Frame() != 2 {
		t.Fatalf("frame update mismatch")
	}
	if v.ModelHash() != "model" || v.MotionHash() != "motion" {
		t.Fatalf("hash mismatch")
	}
	v.SetModelHash("m2")
	v.SetMotionHash("v2")
	if v.ModelHash() != "m2" || v.MotionHash() != "v2" {
		t.Fatalf("hash update mismatch")
	}
}

// TestPhysicsDeltasAccessors はPhysicsDeltasのアクセサを確認する。
func TestPhysicsDeltasAccessors(t *testing.T) {
	rigidBodies := collection.NewNamedCollection[*model.RigidBody](0)
	joints := collection.NewNamedCollection[*model.Joint](0)
	p := NewPhysicsDeltas(1, rigidBodies, joints, "model", "motion")
	if p.Frame() != 1 {
		t.Fatalf("frame mismatch")
	}
	p.SetFrame(2)
	if p.Frame() != 2 {
		t.Fatalf("frame update mismatch")
	}
	if p.ModelHash() != "model" || p.MotionHash() != "motion" {
		t.Fatalf("hash mismatch")
	}
	p.SetModelHash("m2")
	p.SetMotionHash("v2")
	if p.ModelHash() != "m2" || p.MotionHash() != "v2" {
		t.Fatalf("hash update mismatch")
	}
}

// TestRigidBodyDeltasLookup は剛体差分の検索を確認する。
func TestRigidBodyDeltasLookup(t *testing.T) {
	rigidBodies := collection.NewNamedCollection[*model.RigidBody](0)
	body := &model.RigidBody{}
	body.SetName("rb")
	rigidBodies.Append(body)
	deltas := NewRigidBodyDeltas(rigidBodies)
	delta := NewRigidBodyDelta(body, 0)
	deltas.Update(delta)
	if deltas.GetByName("rb") == nil {
		t.Fatalf("expected rigid body delta")
	}
}

// TestJointDeltasLookup はジョイント差分の検索を確認する。
func TestJointDeltasLookup(t *testing.T) {
	joints := collection.NewNamedCollection[*model.Joint](0)
	joint := &model.Joint{}
	joint.SetName("j")
	joints.Append(joint)
	deltas := NewJointDeltas(joints)
	delta := NewJointDelta(joint, mtime.Frame(1))
	deltas.Update(delta)
	if deltas.GetByName("j") == nil {
		t.Fatalf("expected joint delta")
	}
}
