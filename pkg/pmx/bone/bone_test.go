package bone

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/math/mmat4"
	"github.com/miu200521358/mlib_go/pkg/math/mrotation"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
)

func TestIkLink_Copy(t *testing.T) {
	ikLink := &IkLink{
		BoneIndex:          0,
		AngleLimit:         true,
		MinAngleLimit:      mrotation.T{},
		MaxAngleLimit:      mrotation.T{},
		LocalAngleLimit:    true,
		LocalMinAngleLimit: mrotation.T{},
		LocalMaxAngleLimit: mrotation.T{},
	}

	copied := ikLink.Copy()

	if copied == ikLink {
		t.Error("Expected Copy() to return a different instance")
	}

	if copied.BoneIndex != ikLink.BoneIndex ||
		copied.AngleLimit != ikLink.AngleLimit ||
		copied.MinAngleLimit != ikLink.MinAngleLimit ||
		copied.MaxAngleLimit != ikLink.MaxAngleLimit ||
		copied.LocalAngleLimit != ikLink.LocalAngleLimit ||
		copied.LocalMinAngleLimit != ikLink.LocalMinAngleLimit ||
		copied.LocalMaxAngleLimit != ikLink.LocalMaxAngleLimit {
		t.Error("Copied instance does not match the original")
	}
}

func TestIk_Copy(t *testing.T) {
	ik := &Ik{
		BoneIndex:    0,
		LoopCount:    1,
		UnitRotation: *mrotation.NewBaseRotationModelByRadians(&mvec3.T{0, 0, 0}),
		Links: []IkLink{
			{
				BoneIndex:          0,
				AngleLimit:         true,
				MinAngleLimit:      mrotation.T{},
				MaxAngleLimit:      mrotation.T{},
				LocalAngleLimit:    true,
				LocalMinAngleLimit: mrotation.T{},
				LocalMaxAngleLimit: mrotation.T{},
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

	if copied.UnitRotation != ik.UnitRotation {
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
	correctedFixedAxis := mvec3.T{1, 0, 0}
	b.NormalizeFixedAxis(correctedFixedAxis)

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
	b1 := &Bone{Name: "左足D"}
	if !b1.IsLegD() {
		t.Errorf("Expected IsLegD to return true")
	}

	b2 := &Bone{Name: "右腕"}
	if b2.IsLegD() {
		t.Errorf("Expected IsLegD to return false")
	}
}

func TestBone_Copy(t *testing.T) {
	t.Run("Test Copy", func(t *testing.T) {
		b := &Bone{
			Ik:                     &Ik{},
			Position:               mvec3.T{},
			TailPosition:           mvec3.T{},
			FixedAxis:              mvec3.T{},
			LocalAxisX:             mvec3.T{},
			LocalAxisZ:             mvec3.T{},
			NormalizedLocalAxisZ:   mvec3.T{},
			NormalizedLocalAxisX:   mvec3.T{},
			NormalizedLocalAxisY:   mvec3.T{},
			LocalAxis:              mvec3.T{},
			ParentRelativePosition: mvec3.T{},
			TailRelativePosition:   mvec3.T{},
			NormalizedFixedAxis:    mvec3.T{},
			IkLinkBoneIndexes:      []int{},
			IkTargetBoneIndexes:    []int{},
			TreeBoneIndexes:        []int{},
			ParentRevertMatrix:     mmat4.T{},
			OffsetMatrix:           mmat4.T{},
			RelativeBoneIndexes:    []int{},
			ChildBoneIndexes:       []int{},
			EffectiveBoneIndexes:   []int{},
			MinAngleLimit:          mrotation.T{},
			MaxAngleLimit:          mrotation.T{},
			LocalMinAngleLimit:     mrotation.T{},
			LocalMaxAngleLimit:     mrotation.T{},
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
		if copied.Position != b.Position {
			t.Errorf("Expected copied Position to be a deep copy of original Position %s %s", copied.Position.String(), b.Position.String())
		}
		if &copied.Position == &b.Position {
			t.Errorf("Expected copied Position to be a deep copy of original Position %s %s", copied.Position.String(), b.Position.String())
		}
		if copied.TailPosition != b.TailPosition {
			t.Errorf("Expected copied TailPosition to be a deep copy of original TailPosition %s %s", copied.TailPosition.String(), b.TailPosition.String())
		}
		if &copied.TailPosition == &b.TailPosition {
			t.Errorf("Expected copied TailPosition to be a deep copy of original TailPosition %s %s", copied.TailPosition.String(), b.TailPosition.String())
		}
		if copied.FixedAxis != b.FixedAxis {
			t.Errorf("Expected copied FixedAxis to be a deep copy of original FixedAxis %s %s", copied.FixedAxis.String(), b.FixedAxis.String())
		}
		if &copied.FixedAxis == &b.FixedAxis {
			t.Errorf("Expected copied FixedAxis to be a deep copy of original FixedAxis %s %s", copied.FixedAxis.String(), b.FixedAxis.String())
		}
		if copied.LocalAxisX != b.LocalAxisX {
			t.Errorf("Expected copied LocalXVector to be a deep copy of original LocalXVector %s %s", copied.LocalAxisX.String(), b.LocalAxisX.String())
		}
		if &copied.LocalAxisX == &b.LocalAxisX {
			t.Errorf("Expected copied LocalXVector to be a deep copy of original LocalXVector %s %s", copied.LocalAxisX.String(), b.LocalAxisX.String())
		}
		if copied.LocalAxisZ != b.LocalAxisZ {
			t.Errorf("Expected copied LocalZVector to be a deep copy of original LocalZVector %s %s", copied.LocalAxisZ.String(), b.LocalAxisZ.String())
		}
		if &copied.LocalAxisZ == &b.LocalAxisZ {
			t.Errorf("Expected copied LocalZVector to be a deep copy of original LocalZVector %s %s", copied.LocalAxisZ.String(), b.LocalAxisZ.String())
		}
		if copied.NormalizedLocalAxisZ != b.NormalizedLocalAxisZ {
			t.Errorf("Expected copied NormalizedLocalXVector to be a deep copy of original NormalizedLocalXVector %s %s", copied.NormalizedLocalAxisZ.String(), b.NormalizedLocalAxisZ.String())
		}
		if &copied.NormalizedLocalAxisZ == &b.NormalizedLocalAxisZ {
			t.Errorf("Expected copied NormalizedLocalXVector to be a deep copy of original NormalizedLocalXVector %s %s", copied.NormalizedLocalAxisZ.String(), b.NormalizedLocalAxisZ.String())
		}
		if copied.NormalizedLocalAxisX != b.NormalizedLocalAxisX {
			t.Errorf("Expected copied NormalizedLocalYVector to be a deep copy of original NormalizedLocalYVector %s %s", copied.NormalizedLocalAxisX.String(), b.NormalizedLocalAxisX.String())
		}
		if &copied.NormalizedLocalAxisX == &b.NormalizedLocalAxisX {
			t.Errorf("Expected copied NormalizedLocalYVector to be a deep copy of original NormalizedLocalYVector %s %s", copied.NormalizedLocalAxisX.String(), b.NormalizedLocalAxisX.String())
		}
		if copied.NormalizedLocalAxisY != b.NormalizedLocalAxisY {
			t.Errorf("Expected copied NormalizedLocalZVector to be a deep copy of original NormalizedLocalZVector %s %s", copied.NormalizedLocalAxisY.String(), b.NormalizedLocalAxisY.String())
		}
		if &copied.NormalizedLocalAxisY == &b.NormalizedLocalAxisY {
			t.Errorf("Expected copied NormalizedLocalZVector to be a deep copy of original NormalizedLocalZVector %s %s", copied.NormalizedLocalAxisY.String(), b.NormalizedLocalAxisY.String())
		}
		if copied.LocalAxis != b.LocalAxis {
			t.Errorf("Expected copied LocalAxis to be a deep copy of original LocalAxis %s %s", copied.LocalAxis.String(), b.LocalAxis.String())
		}
		if &copied.LocalAxis == &b.LocalAxis {
			t.Errorf("Expected copied LocalAxis to be a deep copy of original LocalAxis %s %s", copied.LocalAxis.String(), b.LocalAxis.String())
		}
		if copied.ParentRelativePosition != b.ParentRelativePosition {
			t.Errorf("Expected copied ParentRelativePosition to be a deep copy of original ParentRelativePosition %s %s", copied.ParentRelativePosition.String(), b.ParentRelativePosition.String())
		}
		if &copied.ParentRelativePosition == &b.ParentRelativePosition {
			t.Errorf("Expected copied ParentRelativePosition to be a deep copy of original ParentRelativePosition %s %s", copied.ParentRelativePosition.String(), b.ParentRelativePosition.String())
		}
		if copied.TailRelativePosition != b.TailRelativePosition {
			t.Errorf("Expected copied TailRelativePosition to be a deep copy of original TailRelativePosition %s %s", copied.TailRelativePosition.String(), b.TailRelativePosition.String())
		}
		if &copied.TailRelativePosition == &b.TailRelativePosition {
			t.Errorf("Expected copied TailRelativePosition to be a deep copy of original TailRelativePosition %s %s", copied.TailRelativePosition.String(), b.TailRelativePosition.String())
		}
		if copied.NormalizedFixedAxis != b.NormalizedFixedAxis {
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
		if copied.ParentRevertMatrix != b.ParentRevertMatrix {
			t.Errorf("Expected copied ParentRevertMatrix to be a deep copy %s %s", copied.ParentRevertMatrix.String(), b.ParentRevertMatrix.String())
		}
		if &copied.ParentRevertMatrix == &b.ParentRevertMatrix {
			t.Errorf("Expected copied ParentRevertMatrix to be a deep copy")
		}
		if copied.OffsetMatrix != b.OffsetMatrix {
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
		if copied.MinAngleLimit != b.MinAngleLimit {
			t.Errorf("Expected copied MinAngleLimit to be a deep copy %s %s", copied.MinAngleLimit.GetDegrees().String(), b.MinAngleLimit.GetDegrees().String())
		}
		if &copied.MinAngleLimit == &b.MinAngleLimit {
			t.Errorf("Expected copied MinAngleLimit to be a deep copy")
		}
		if copied.MaxAngleLimit != b.MaxAngleLimit {
			t.Errorf("Expected copied MaxAngleLimit to be a deep copy %s %s", copied.MaxAngleLimit.GetDegrees().String(), b.MaxAngleLimit.GetDegrees().String())
		}
		if &copied.MaxAngleLimit == &b.MaxAngleLimit {
			t.Errorf("Expected copied MaxAngleLimit to be a deep copy")
		}
		if copied.LocalMinAngleLimit != b.LocalMinAngleLimit {
			t.Errorf("Expected copied LocalMinAngleLimit to be a deep copy %s %s", copied.LocalMinAngleLimit.GetDegrees().String(), b.LocalMinAngleLimit.GetDegrees().String())
		}
		if &copied.LocalMinAngleLimit == &b.LocalMinAngleLimit {
			t.Errorf("Expected copied LocalMinAngleLimit to be a deep copy")
		}
		if copied.LocalMaxAngleLimit != b.LocalMaxAngleLimit {
			t.Errorf("Expected copied LocalMaxAngleLimit to be a deep copy %s %s", copied.LocalMaxAngleLimit.GetDegrees().String(), b.LocalMaxAngleLimit.GetDegrees().String())
		}
		if &copied.LocalMaxAngleLimit == &b.LocalMaxAngleLimit {
			t.Errorf("Expected copied LocalMaxAngleLimit to be a deep copy")
		}
	})
}
