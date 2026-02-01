// 指示: miu200521358
package mmath

import (
	"math"
	"testing"
)

func TestMat4Basics(t *testing.T) {
	m := NewMat4()
	if !m.IsIdent() {
		t.Errorf("NewMat4")
	}
	m2 := NewMat4ByValues(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	if !m2.IsIdent() {
		t.Errorf("NewMat4ByValues")
	}
	m3 := NewMat4FromAxisAngle(Vec3{r3Vec(1, 0, 0)}, math.Pi/2)
	if m3.IsIdent() {
		t.Errorf("NewMat4FromAxisAngle")
	}
	m4 := NewMat4FromLookAt(Vec3{r3Vec(0, 0, 0)}, Vec3{r3Vec(0, 0, -1)}, Vec3{r3Vec(0, 1, 0)})
	if m4.IsZero() {
		t.Errorf("NewMat4FromLookAt")
	}
	if ZERO_MAT4.IsZero() == false || ZERO_MAT4.IsIdent() {
		t.Errorf("IsZero/IsIdent")
	}
	if m.String() == "" {
		t.Errorf("String")
	}
	if cp, err := m.Copy(); err != nil || cp != m {
		t.Errorf("Copy")
	}
	if !m.NearEquals(m2, 1e-10) {
		t.Errorf("NearEquals")
	}
	if m.Trace() != 4 || m.Trace3() != 3 {
		t.Errorf("Trace")
	}
	if m.MulVec3(Vec3{r3Vec(1, 2, 3)}) != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("MulVec3")
	}
	m5 := m
	m5[15] = 2
	if m5.MulVec3(Vec3{r3Vec(1, 0, 0)}).X != 0.5 {
		t.Errorf("MulVec3 w")
	}
	mt := m
	mt.Translate(Vec3{r3Vec(1, 2, 3)})
	if mt.Translation() != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("Translate")
	}
	if m.Translated(Vec3{r3Vec(1, 2, 3)}).Translation() != (Vec3{r3Vec(1, 2, 3)}) {
		t.Errorf("Translated")
	}
	ms := m
	ms.Scale(Vec3{r3Vec(2, 3, 4)})
	if ms.Scaling() != (Vec3{r3Vec(2, 3, 4)}) {
		t.Errorf("Scale/Scaling")
	}
	if m.Scaled(Vec3{r3Vec(2, 3, 4)}).Scaling() != (Vec3{r3Vec(2, 3, 4)}) {
		t.Errorf("Scaled")
	}
	mr := m
	mr.Rotate(NewQuaternion())
	if !mr.IsIdent() {
		t.Errorf("Rotate")
	}
	if !m.Rotated(NewQuaternion()).IsIdent() {
		t.Errorf("Rotated")
	}
	if !m.Quaternion().IsIdent() {
		t.Errorf("Quaternion")
	}
	m.Transpose()
	m.Transpose()
	if !m.IsIdent() {
		t.Errorf("Transpose")
	}
	m6 := NewMat4()
	m6.Mul(NewMat4())
	if !m6.IsIdent() {
		t.Errorf("Mul")
	}
	if m6.Muled(NewMat4()).IsIdent() == false {
		t.Errorf("Muled")
	}
	m7 := NewMat4()
	m7.Add(NewMat4())
	if m7.Trace() != 8 {
		t.Errorf("Add")
	}
	if m7.Added(IDENT_MAT4).Trace() != 12 {
		t.Errorf("Added")
	}
	m8 := NewMat4()
	m8.MulScalar(2)
	if m8.Trace() != 8 {
		t.Errorf("MulScalar")
	}
	if m8.MuledScalar(0.5).Trace() != 4 {
		t.Errorf("MuledScalar")
	}
	if NewMat4().Det() != 1 {
		t.Errorf("Det")
	}
	if ZERO_MAT4.Inverted() != NewMat4() {
		t.Errorf("Inverted singular")
	}
	inv := NewMat4().Inverted()
	if !inv.IsIdent() {
		t.Errorf("Inverted")
	}
	m9 := NewMat4()
	m9.Inverse()
	if !m9.IsIdent() {
		t.Errorf("Inverse")
	}
	m10 := Mat4{1e-7, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	m10.ClampIfVerySmall()
	if m10[0] != 0 {
		t.Errorf("ClampIfVerySmall")
	}
	if m.AxisX() == ZERO_VEC3 || m.AxisY() == ZERO_VEC3 || m.AxisZ() == ZERO_VEC3 {
		t.Errorf("Axis")
	}
}

func TestMat4QuaternionBranches(t *testing.T) {
	mx := NewMat4ByValues(
		1, 0, 0, 0,
		0, -1, 0, 0,
		0, 0, -1, 0,
		0, 0, 0, 1,
	)
	if mx.Quaternion() == quaternionZero {
		t.Errorf("Quaternion branch x")
	}
	my := NewMat4ByValues(
		-1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, -1, 0,
		0, 0, 0, 1,
	)
	if my.Quaternion() == quaternionZero {
		t.Errorf("Quaternion branch y")
	}
	mz := NewMat4ByValues(
		-1, 0, 0, 0,
		0, -1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	if mz.Quaternion() == quaternionZero {
		t.Errorf("Quaternion branch z")
	}
}

func TestMat4InvertedError(t *testing.T) {
	m := NewMat4()
	m[0] = math.NaN()
	if m.Inverted() != NewMat4() {
		t.Errorf("Inverted NaN")
	}
}
