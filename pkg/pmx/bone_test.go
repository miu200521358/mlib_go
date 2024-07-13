package pmx

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/mmath"
)

func TestIkLink_Copy(t *testing.T) {
	ikLink := &IkLink{
		BoneIndex:          0,
		AngleLimit:         true,
		MinAngleLimit:      mmath.NewRotation(),
		MaxAngleLimit:      mmath.NewRotation(),
		LocalAngleLimit:    true,
		LocalMinAngleLimit: mmath.NewRotation(),
		LocalMaxAngleLimit: mmath.NewRotation(),
	}

	copied := ikLink.Copy()

	if copied == ikLink {
		t.Error("Expected Copy() to return a different instance")
	}

	if copied.BoneIndex != ikLink.BoneIndex {
		t.Error("Expected BoneIndex to match the original")
	}
	if copied.AngleLimit != ikLink.AngleLimit {
		t.Error("Expected AngleLimit to match the original")
	}
	if copied.MinAngleLimit.String() != ikLink.MinAngleLimit.String() {
		t.Error("Expected MinAngleLimit to match the original")
	}
	if copied.MaxAngleLimit.String() != ikLink.MaxAngleLimit.String() {
		t.Error("Expected MaxAngleLimit to match the original")
	}
	if copied.LocalAngleLimit != ikLink.LocalAngleLimit {
		t.Error("Expected LocalAngleLimit to match the original")
	}
	if copied.LocalMinAngleLimit.String() != ikLink.LocalMinAngleLimit.String() {
		t.Error("Expected LocalMinAngleLimit to match the original")
	}
	if copied.LocalMaxAngleLimit.String() != ikLink.LocalMaxAngleLimit.String() {
		t.Error("Expected LocalMaxAngleLimit to match the original")
	}
}

func TestIk_Copy(t *testing.T) {
	ik := &Ik{
		BoneIndex:    0,
		LoopCount:    1,
		UnitRotation: mmath.NewRotationFromDegrees(&mmath.MVec3{1, 2, 3}),
		Links: []*IkLink{
			{
				BoneIndex:          0,
				AngleLimit:         true,
				MinAngleLimit:      mmath.NewRotationFromDegrees(&mmath.MVec3{1, 2, 3}),
				MaxAngleLimit:      mmath.NewRotationFromDegrees(&mmath.MVec3{4, 5, 6}),
				LocalAngleLimit:    true,
				LocalMinAngleLimit: mmath.NewRotationFromDegrees(&mmath.MVec3{7, 8, 9}),
				LocalMaxAngleLimit: mmath.NewRotationFromDegrees(&mmath.MVec3{10, 11, 12}),
			},
		},
	}

	copied := ik.Copy()

	if copied.BoneIndex != ik.BoneIndex {
		t.Error("Expected BoneIndex to match the original")
	}

	if copied.LoopCount != ik.LoopCount {
		t.Error("Expected LoopCount to match the original")
	}

	if !copied.UnitRotation.GetDegrees().NearEquals(ik.UnitRotation.GetDegrees(), 1e-8) {
		t.Error("Expected UnitRotation to match the original")
	}

	if len(copied.Links) != len(ik.Links) {
		t.Error("Expected the length of Links to match the original")
	}

	if &copied.Links[0] == &ik.Links[0] {
		t.Error("Expected Links[0] to be a different instance")
	}
}

func TestBone_NormalizeFixedAxis(t *testing.T) {
	b := &Bone{}
	correctedFixedAxis := mmath.MVec3{1, 0, 0}
	b.NormalizeFixedAxis(&correctedFixedAxis)

	if !b.NormalizedFixedAxis.Equals(correctedFixedAxis.Normalize()) {
		t.Errorf("Expected NormalizedFixedAxis to be normalized")
	}
}

func TestBone_IsTailBone(t *testing.T) {
	b := &Bone{BoneFlag: BONE_FLAG_TAIL_IS_BONE}
	if !b.IsTailBone() {
		t.Errorf("Expected IsTailBone to return true")
	}

	b.BoneFlag = 0
	if b.IsTailBone() {
		t.Errorf("Expected IsTailBone to return false")
	}
}

func TestBone_IsLegD(t *testing.T) {
	b1 := NewBoneByName("左足D")
	if !b1.IsLegD() {
		t.Errorf("Expected IsLegD to return true")
	}

	b2 := NewBoneByName("右腕")
	if b2.IsLegD() {
		t.Errorf("Expected IsLegD to return false")
	}
}

func TestBone_Copy(t *testing.T) {
	t.Run("Test Copy", func(t *testing.T) {
		b := &Bone{
			IndexNameModel:         &core.IndexNameModel{Index: 0, Name: "Bone"},
			Ik:                     NewIk(),
			Position:               &mmath.MVec3{1, 2, 3},
			TailPosition:           &mmath.MVec3{4, 5, 6},
			FixedAxis:              &mmath.MVec3{7, 8, 9},
			LocalAxisX:             &mmath.MVec3{10, 11, 12},
			LocalAxisZ:             &mmath.MVec3{13, 14, 15},
			NormalizedLocalAxisZ:   &mmath.MVec3{16, 17, 18},
			NormalizedLocalAxisX:   &mmath.MVec3{19, 20, 21},
			NormalizedLocalAxisY:   &mmath.MVec3{22, 23, 24},
			LocalAxis:              &mmath.MVec3{25, 26, 27},
			ParentRelativePosition: &mmath.MVec3{28, 29, 30},
			ChildRelativePosition:  &mmath.MVec3{31, 32, 33},
			NormalizedFixedAxis:    &mmath.MVec3{34, 35, 36},
			IkLinkBoneIndexes:      []int{1, 3, 5},
			IkTargetBoneIndexes:    []int{2, 4, 6},
			TreeBoneIndexes:        []int{3, 5, 7},
			RevertOffsetMatrix:     &mmath.MMat4{},
			OffsetMatrix:           &mmath.MMat4{},
			RelativeBoneIndexes:    []int{8, 9, 10},
			ChildBoneIndexes:       []int{10, 11, 12},
			EffectiveBoneIndexes:   []int{16, 17, 18},
			MinAngleLimit:          mmath.NewRotationFromRadians(&mmath.MVec3{1, 2, 3}),
			MaxAngleLimit:          mmath.NewRotationFromRadians(&mmath.MVec3{5, 6, 7}),
			LocalMinAngleLimit:     mmath.NewRotationFromRadians(&mmath.MVec3{10, 11, 12}),
			LocalMaxAngleLimit:     mmath.NewRotationFromRadians(&mmath.MVec3{16, 17, 18}),
		}

		copied := b.Copy().(*Bone)

		// Assert copied fields are not the same as original
		if copied == b {
			t.Errorf("Expected copied bone to be different from original")
		}

		// Assert copied fields are deep copies
		if copied.Ik == b.Ik {
			t.Errorf("Expected copied Ik to be a deep copy")
		}
		if !copied.Position.NearEquals(b.Position, 1e-10) {
			t.Errorf("Expected copied Position to be a deep copy of original Position %s %s", copied.Position.String(), b.Position.String())
		}
		if &copied.Position == &b.Position {
			t.Errorf("Expected copied Position to be a deep copy of original Position %s %s", copied.Position.String(), b.Position.String())
		}
		if !copied.TailPosition.NearEquals(b.TailPosition, 1e-10) {
			t.Errorf("Expected copied TailPosition to be a deep copy of original TailPosition %s %s", copied.TailPosition.String(), b.TailPosition.String())
		}
		if &copied.TailPosition == &b.TailPosition {
			t.Errorf("Expected copied TailPosition to be a deep copy of original TailPosition %s %s", copied.TailPosition.String(), b.TailPosition.String())
		}
		if !copied.FixedAxis.NearEquals(b.FixedAxis, 1e-10) {
			t.Errorf("Expected copied FixedAxis to be a deep copy of original FixedAxis %s %s", copied.FixedAxis.String(), b.FixedAxis.String())
		}
		if &copied.FixedAxis == &b.FixedAxis {
			t.Errorf("Expected copied FixedAxis to be a deep copy of original FixedAxis %s %s", copied.FixedAxis.String(), b.FixedAxis.String())
		}
		if !copied.LocalAxisX.NearEquals(b.LocalAxisX, 1e-10) {
			t.Errorf("Expected copied LocalXVector to be a deep copy of original LocalXVector %s %s", copied.LocalAxisX.String(), b.LocalAxisX.String())
		}
		if &copied.LocalAxisX == &b.LocalAxisX {
			t.Errorf("Expected copied LocalXVector to be a deep copy of original LocalXVector %s %s", copied.LocalAxisX.String(), b.LocalAxisX.String())
		}
		if !copied.LocalAxisZ.NearEquals(b.LocalAxisZ, 1e-10) {
			t.Errorf("Expected copied LocalZVector to be a deep copy of original LocalZVector %s %s", copied.LocalAxisZ.String(), b.LocalAxisZ.String())
		}
		if &copied.LocalAxisZ == &b.LocalAxisZ {
			t.Errorf("Expected copied LocalZVector to be a deep copy of original LocalZVector %s %s", copied.LocalAxisZ.String(), b.LocalAxisZ.String())
		}
		if !copied.NormalizedLocalAxisZ.NearEquals(b.NormalizedLocalAxisZ, 1e-10) {
			t.Errorf("Expected copied NormalizedLocalXVector to be a deep copy of original NormalizedLocalXVector %s %s", copied.NormalizedLocalAxisZ.String(), b.NormalizedLocalAxisZ.String())
		}
		if &copied.NormalizedLocalAxisZ == &b.NormalizedLocalAxisZ {
			t.Errorf("Expected copied NormalizedLocalXVector to be a deep copy of original NormalizedLocalXVector %s %s", copied.NormalizedLocalAxisZ.String(), b.NormalizedLocalAxisZ.String())
		}
		if !copied.NormalizedLocalAxisX.NearEquals(b.NormalizedLocalAxisX, 1e-10) {
			t.Errorf("Expected copied NormalizedLocalYVector to be a deep copy of original NormalizedLocalYVector %s %s", copied.NormalizedLocalAxisX.String(), b.NormalizedLocalAxisX.String())
		}
		if &copied.NormalizedLocalAxisX == &b.NormalizedLocalAxisX {
			t.Errorf("Expected copied NormalizedLocalYVector to be a deep copy of original NormalizedLocalYVector %s %s", copied.NormalizedLocalAxisX.String(), b.NormalizedLocalAxisX.String())
		}
		if !copied.NormalizedLocalAxisY.NearEquals(b.NormalizedLocalAxisY, 1e-10) {
			t.Errorf("Expected copied NormalizedLocalZVector to be a deep copy of original NormalizedLocalZVector %s %s", copied.NormalizedLocalAxisY.String(), b.NormalizedLocalAxisY.String())
		}
		if &copied.NormalizedLocalAxisY == &b.NormalizedLocalAxisY {
			t.Errorf("Expected copied NormalizedLocalZVector to be a deep copy of original NormalizedLocalZVector %s %s", copied.NormalizedLocalAxisY.String(), b.NormalizedLocalAxisY.String())
		}
		if !copied.LocalAxis.NearEquals(b.LocalAxis, 1e-10) {
			t.Errorf("Expected copied LocalAxis to be a deep copy of original LocalAxis %s %s", copied.LocalAxis.String(), b.LocalAxis.String())
		}
		if &copied.LocalAxis == &b.LocalAxis {
			t.Errorf("Expected copied LocalAxis to be a deep copy of original LocalAxis %s %s", copied.LocalAxis.String(), b.LocalAxis.String())
		}
		if !copied.ParentRelativePosition.NearEquals(b.ParentRelativePosition, 1e-10) {
			t.Errorf("Expected copied ParentRelativePosition to be a deep copy of original ParentRelativePosition %s %s", copied.ParentRelativePosition.String(), b.ParentRelativePosition.String())
		}
		if &copied.ParentRelativePosition == &b.ParentRelativePosition {
			t.Errorf("Expected copied ParentRelativePosition to be a deep copy of original ParentRelativePosition %s %s", copied.ParentRelativePosition.String(), b.ParentRelativePosition.String())
		}
		if !copied.ChildRelativePosition.NearEquals(b.ChildRelativePosition, 1e-10) {
			t.Errorf("Expected copied TailRelativePosition to be a deep copy of original TailRelativePosition %s %s", copied.ChildRelativePosition.String(), b.ChildRelativePosition.String())
		}
		if &copied.ChildRelativePosition == &b.ChildRelativePosition {
			t.Errorf("Expected copied TailRelativePosition to be a deep copy of original TailRelativePosition %s %s", copied.ChildRelativePosition.String(), b.ChildRelativePosition.String())
		}
		if !copied.NormalizedFixedAxis.NearEquals(b.NormalizedFixedAxis, 1e-10) {
			t.Errorf("Expected copied NormalizedFixedAxis to be a deep copy of original NormalizedFixedAxis %s %s", copied.NormalizedFixedAxis.String(), b.NormalizedFixedAxis.String())
		}
		if &copied.NormalizedFixedAxis == &b.NormalizedFixedAxis {
			t.Errorf("Expected copied NormalizedFixedAxis to be a deep copy of original NormalizedFixedAxis %s %s", copied.NormalizedFixedAxis.String(), b.NormalizedFixedAxis.String())
		}
		if len(copied.IkLinkBoneIndexes) != len(b.IkLinkBoneIndexes) {
			t.Errorf("Expected copied IkLinkIndexes to have the same length as original")
		}
		if len(copied.IkTargetBoneIndexes) != len(b.IkTargetBoneIndexes) {
			t.Errorf("Expected copied IkTargetIndexes to have the same length as original")
		}
		if len(copied.TreeBoneIndexes) != len(b.TreeBoneIndexes) {
			t.Errorf("Expected copied TreeIndexes to have the same length as original")
		}
		if !copied.RevertOffsetMatrix.NearEquals(b.RevertOffsetMatrix, 1e-10) {
			t.Errorf("Expected copied ParentRevertMatrix to be a deep copy %s %s", copied.RevertOffsetMatrix.String(), b.RevertOffsetMatrix.String())
		}
		if &copied.RevertOffsetMatrix == &b.RevertOffsetMatrix {
			t.Errorf("Expected copied ParentRevertMatrix to be a deep copy")
		}
		if !copied.OffsetMatrix.NearEquals(b.OffsetMatrix, 1e-10) {
			t.Errorf("Expected copied OffsetMatrix to be a deep copy %s %s", copied.OffsetMatrix.String(), b.OffsetMatrix.String())
		}
		if &copied.OffsetMatrix == &b.OffsetMatrix {
			t.Errorf("Expected copied OffsetMatrix to be a deep copy")
		}
		if len(copied.RelativeBoneIndexes) != len(b.RelativeBoneIndexes) {
			t.Errorf("Expected copied RelativeBoneIndexes to have the same length as original")
		}
		if len(copied.ChildBoneIndexes) != len(b.ChildBoneIndexes) {
			t.Errorf("Expected copied ChildBoneIndexes to have the same length as original")
		}
		if len(copied.EffectiveBoneIndexes) != len(b.EffectiveBoneIndexes) {
			t.Errorf("Expected copied EffectiveBoneIndexes to have the same length as original")
		}
		if !copied.MinAngleLimit.GetDegrees().NearEquals(b.MinAngleLimit.GetDegrees(), 1e-10) {
			t.Errorf("Expected copied MinAngleLimit to be a deep copy %s %s", copied.MinAngleLimit.GetDegrees().String(), b.MinAngleLimit.GetDegrees().String())
		}
		if &copied.MinAngleLimit == &b.MinAngleLimit {
			t.Errorf("Expected copied MinAngleLimit to be a deep copy")
		}
		if !copied.MaxAngleLimit.GetDegrees().NearEquals(b.MaxAngleLimit.GetDegrees(), 1e-10) {
			t.Errorf("Expected copied MaxAngleLimit to be a deep copy %s %s", copied.MaxAngleLimit.GetDegrees().String(), b.MaxAngleLimit.GetDegrees().String())
		}
		if &copied.MaxAngleLimit == &b.MaxAngleLimit {
			t.Errorf("Expected copied MaxAngleLimit to be a deep copy")
		}
		if !copied.LocalMinAngleLimit.GetDegrees().NearEquals(b.LocalMinAngleLimit.GetDegrees(), 1e-10) {
			t.Errorf("Expected copied LocalMinAngleLimit to be a deep copy %s %s", copied.LocalMinAngleLimit.GetDegrees().String(), b.LocalMinAngleLimit.GetDegrees().String())
		}
		if &copied.LocalMinAngleLimit == &b.LocalMinAngleLimit {
			t.Errorf("Expected copied LocalMinAngleLimit to be a deep copy")
		}
		if !copied.LocalMaxAngleLimit.GetDegrees().NearEquals(b.LocalMaxAngleLimit.GetDegrees(), 1e-10) {
			t.Errorf("Expected copied LocalMaxAngleLimit to be a deep copy %s %s", copied.LocalMaxAngleLimit.GetDegrees().String(), b.LocalMaxAngleLimit.GetDegrees().String())
		}
		if &copied.LocalMaxAngleLimit == &b.LocalMaxAngleLimit {
			t.Errorf("Expected copied LocalMaxAngleLimit to be a deep copy")
		}
	})
}
