package pmx

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestIkLink_Copy(t *testing.T) {
	ikLink := &IkLink{
		BoneIndex:          0,
		AngleLimit:         true,
		MinAngleLimit:      mmath.NewMVec3(),
		MaxAngleLimit:      mmath.NewMVec3(),
		LocalAngleLimit:    true,
		LocalMinAngleLimit: mmath.NewMVec3(),
		LocalMaxAngleLimit: mmath.NewMVec3(),
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
		UnitRotation: &mmath.MVec3{X: 1, Y: 2, Z: 3},
		Links: []*IkLink{
			{
				BoneIndex:          0,
				AngleLimit:         true,
				MinAngleLimit:      &mmath.MVec3{X: 1, Y: 2, Z: 3},
				MaxAngleLimit:      &mmath.MVec3{X: 4, Y: 5, Z: 6},
				LocalAngleLimit:    true,
				LocalMinAngleLimit: &mmath.MVec3{X: 7, Y: 8, Z: 9},
				LocalMaxAngleLimit: &mmath.MVec3{X: 10, Y: 11, Z: 12},
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

	if !copied.UnitRotation.NearEquals(ik.UnitRotation, 1e-8) {
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
	b := NewBone()
	correctedFixedAxis := mmath.MVec3{X: 1, Y: 0, Z: 0}
	b.NormalizeFixedAxis(&correctedFixedAxis)

	if !b.Extend.NormalizedFixedAxis.Equals(correctedFixedAxis.Normalize()) {
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
			index:        0,
			name:         "Bone",
			englishName:  "Bone",
			Ik:           NewIk(),
			Position:     &mmath.MVec3{X: 1, Y: 2, Z: 3},
			TailPosition: &mmath.MVec3{X: 4, Y: 5, Z: 6},
			FixedAxis:    &mmath.MVec3{X: 7, Y: 8, Z: 9},
			LocalAxisX:   &mmath.MVec3{X: 10, Y: 11, Z: 12},
			LocalAxisZ:   &mmath.MVec3{X: 13, Y: 14, Z: 15},
			Extend: &BoneExtend{
				NormalizedLocalAxisZ:   &mmath.MVec3{X: 16, Y: 17, Z: 18},
				NormalizedLocalAxisX:   &mmath.MVec3{X: 19, Y: 20, Z: 21},
				NormalizedLocalAxisY:   &mmath.MVec3{X: 22, Y: 23, Z: 24},
				LocalAxis:              &mmath.MVec3{X: 25, Y: 26, Z: 27},
				ParentRelativePosition: &mmath.MVec3{X: 28, Y: 29, Z: 30},
				ChildRelativePosition:  &mmath.MVec3{X: 31, Y: 32, Z: 33},
				NormalizedFixedAxis:    &mmath.MVec3{X: 34, Y: 35, Z: 36},
				IkLinkBoneIndexes:      []int{1, 3, 5},
				IkTargetBoneIndexes:    []int{2, 4, 6},
				TreeBoneIndexes:        []int{3, 5, 7},
				RevertOffsetMatrix:     &mmath.MMat4{},
				OffsetMatrix:           &mmath.MMat4{},
				RelativeBoneIndexes:    []int{8, 9, 10},
				ChildBoneIndexes:       []int{10, 11, 12},
				EffectiveBoneIndexes:   []int{16, 17, 18},
				MinAngleLimit:          &mmath.MVec3{X: 1, Y: 2, Z: 3},
				MaxAngleLimit:          &mmath.MVec3{X: 5, Y: 6, Z: 7},
				LocalMinAngleLimit:     &mmath.MVec3{X: 10, Y: 11, Z: 12},
				LocalMaxAngleLimit:     &mmath.MVec3{X: 16, Y: 17, Z: 18},
			},
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
		if !copied.Extend.NormalizedLocalAxisZ.NearEquals(b.Extend.NormalizedLocalAxisZ, 1e-10) {
			t.Errorf("Expected copied NormalizedLocalXVector to be a deep copy of original NormalizedLocalXVector %s %s", copied.Extend.NormalizedLocalAxisZ.String(), b.Extend.NormalizedLocalAxisZ.String())
		}
		if &copied.Extend.NormalizedLocalAxisZ == &b.Extend.NormalizedLocalAxisZ {
			t.Errorf("Expected copied NormalizedLocalXVector to be a deep copy of original NormalizedLocalXVector %s %s", copied.Extend.NormalizedLocalAxisZ.String(), b.Extend.NormalizedLocalAxisZ.String())
		}
		if !copied.Extend.NormalizedLocalAxisX.NearEquals(b.Extend.NormalizedLocalAxisX, 1e-10) {
			t.Errorf("Expected copied NormalizedLocalYVector to be a deep copy of original NormalizedLocalYVector %s %s", copied.Extend.NormalizedLocalAxisX.String(), b.Extend.NormalizedLocalAxisX.String())
		}
		if &copied.Extend.NormalizedLocalAxisX == &b.Extend.NormalizedLocalAxisX {
			t.Errorf("Expected copied NormalizedLocalYVector to be a deep copy of original NormalizedLocalYVector %s %s", copied.Extend.NormalizedLocalAxisX.String(), b.Extend.NormalizedLocalAxisX.String())
		}
		if !copied.Extend.NormalizedLocalAxisY.NearEquals(b.Extend.NormalizedLocalAxisY, 1e-10) {
			t.Errorf("Expected copied NormalizedLocalZVector to be a deep copy of original NormalizedLocalZVector %s %s", copied.Extend.NormalizedLocalAxisY.String(), b.Extend.NormalizedLocalAxisY.String())
		}
		if &copied.Extend.NormalizedLocalAxisY == &b.Extend.NormalizedLocalAxisY {
			t.Errorf("Expected copied NormalizedLocalZVector to be a deep copy of original NormalizedLocalZVector %s %s", copied.Extend.NormalizedLocalAxisY.String(), b.Extend.NormalizedLocalAxisY.String())
		}
		if !copied.Extend.LocalAxis.NearEquals(b.Extend.LocalAxis, 1e-10) {
			t.Errorf("Expected copied LocalAxis to be a deep copy of original LocalAxis %s %s", copied.Extend.LocalAxis.String(), b.Extend.LocalAxis.String())
		}
		if &copied.Extend.LocalAxis == &b.Extend.LocalAxis {
			t.Errorf("Expected copied LocalAxis to be a deep copy of original LocalAxis %s %s", copied.Extend.LocalAxis.String(), b.Extend.LocalAxis.String())
		}
		if !copied.Extend.ParentRelativePosition.NearEquals(b.Extend.ParentRelativePosition, 1e-10) {
			t.Errorf("Expected copied ParentRelativePosition to be a deep copy of original ParentRelativePosition %s %s", copied.Extend.ParentRelativePosition.String(), b.Extend.ParentRelativePosition.String())
		}
		if &copied.Extend.ParentRelativePosition == &b.Extend.ParentRelativePosition {
			t.Errorf("Expected copied ParentRelativePosition to be a deep copy of original ParentRelativePosition %s %s", copied.Extend.ParentRelativePosition.String(), b.Extend.ParentRelativePosition.String())
		}
		if !copied.Extend.ChildRelativePosition.NearEquals(b.Extend.ChildRelativePosition, 1e-10) {
			t.Errorf("Expected copied TailRelativePosition to be a deep copy of original TailRelativePosition %s %s", copied.Extend.ChildRelativePosition.String(), b.Extend.ChildRelativePosition.String())
		}
		if &copied.Extend.ChildRelativePosition == &b.Extend.ChildRelativePosition {
			t.Errorf("Expected copied TailRelativePosition to be a deep copy of original TailRelativePosition %s %s", copied.Extend.ChildRelativePosition.String(), b.Extend.ChildRelativePosition.String())
		}
		if !copied.Extend.NormalizedFixedAxis.NearEquals(b.Extend.NormalizedFixedAxis, 1e-10) {
			t.Errorf("Expected copied NormalizedFixedAxis to be a deep copy of original NormalizedFixedAxis %s %s", copied.Extend.NormalizedFixedAxis.String(), b.Extend.NormalizedFixedAxis.String())
		}
		if &copied.Extend.NormalizedFixedAxis == &b.Extend.NormalizedFixedAxis {
			t.Errorf("Expected copied NormalizedFixedAxis to be a deep copy of original NormalizedFixedAxis %s %s", copied.Extend.NormalizedFixedAxis.String(), b.Extend.NormalizedFixedAxis.String())
		}
		if len(copied.Extend.IkLinkBoneIndexes) != len(b.Extend.IkLinkBoneIndexes) {
			t.Errorf("Expected copied IkLinkIndexes to have the same length as original")
		}
		if len(copied.Extend.IkTargetBoneIndexes) != len(b.Extend.IkTargetBoneIndexes) {
			t.Errorf("Expected copied IkTargetIndexes to have the same length as original")
		}
		if len(copied.Extend.TreeBoneIndexes) != len(b.Extend.TreeBoneIndexes) {
			t.Errorf("Expected copied TreeIndexes to have the same length as original")
		}
		if !copied.Extend.RevertOffsetMatrix.NearEquals(b.Extend.RevertOffsetMatrix, 1e-10) {
			t.Errorf("Expected copied ParentRevertMatrix to be a deep copy %s %s", copied.Extend.RevertOffsetMatrix.String(), b.Extend.RevertOffsetMatrix.String())
		}
		if &copied.Extend.RevertOffsetMatrix == &b.Extend.RevertOffsetMatrix {
			t.Errorf("Expected copied ParentRevertMatrix to be a deep copy")
		}
		if !copied.Extend.OffsetMatrix.NearEquals(b.Extend.OffsetMatrix, 1e-10) {
			t.Errorf("Expected copied OffsetMatrix to be a deep copy %s %s", copied.Extend.OffsetMatrix.String(), b.Extend.OffsetMatrix.String())
		}
		if &copied.Extend.OffsetMatrix == &b.Extend.OffsetMatrix {
			t.Errorf("Expected copied OffsetMatrix to be a deep copy")
		}
		if len(copied.Extend.RelativeBoneIndexes) != len(b.Extend.RelativeBoneIndexes) {
			t.Errorf("Expected copied RelativeBoneIndexes to have the same length as original")
		}
		if len(copied.Extend.ChildBoneIndexes) != len(b.Extend.ChildBoneIndexes) {
			t.Errorf("Expected copied ChildBoneIndexes to have the same length as original")
		}
		if len(copied.Extend.EffectiveBoneIndexes) != len(b.Extend.EffectiveBoneIndexes) {
			t.Errorf("Expected copied EffectiveBoneIndexes to have the same length as original")
		}
		if !copied.Extend.MinAngleLimit.NearEquals(b.Extend.MinAngleLimit, 1e-10) {
			t.Errorf("Expected copied MinAngleLimit to be a deep copy %s %s", copied.Extend.MinAngleLimit.String(), b.Extend.MinAngleLimit.String())
		}
		if &copied.Extend.MinAngleLimit == &b.Extend.MinAngleLimit {
			t.Errorf("Expected copied MinAngleLimit to be a deep copy")
		}
		if !copied.Extend.MaxAngleLimit.NearEquals(b.Extend.MaxAngleLimit, 1e-10) {
			t.Errorf("Expected copied MaxAngleLimit to be a deep copy %s %s", copied.Extend.MaxAngleLimit.String(), b.Extend.MaxAngleLimit.String())
		}
		if &copied.Extend.MaxAngleLimit == &b.Extend.MaxAngleLimit {
			t.Errorf("Expected copied MaxAngleLimit to be a deep copy")
		}
		if !copied.Extend.LocalMinAngleLimit.NearEquals(b.Extend.LocalMinAngleLimit, 1e-10) {
			t.Errorf("Expected copied LocalMinAngleLimit to be a deep copy %s %s", copied.Extend.LocalMinAngleLimit.String(), b.Extend.LocalMinAngleLimit.String())
		}
		if &copied.Extend.LocalMinAngleLimit == &b.Extend.LocalMinAngleLimit {
			t.Errorf("Expected copied LocalMinAngleLimit to be a deep copy")
		}
		if !copied.Extend.LocalMaxAngleLimit.NearEquals(b.Extend.LocalMaxAngleLimit, 1e-10) {
			t.Errorf("Expected copied LocalMaxAngleLimit to be a deep copy %s %s", copied.Extend.LocalMaxAngleLimit.String(), b.Extend.LocalMaxAngleLimit.String())
		}
		if &copied.Extend.LocalMaxAngleLimit == &b.Extend.LocalMaxAngleLimit {
			t.Errorf("Expected copied LocalMaxAngleLimit to be a deep copy")
		}
	})
}

func TestBones_Insert(t *testing.T) {
	bones := NewBones(0)

	// Create some bones to insert
	bone1 := NewBone()
	bone1.index = 0
	bone1.name = "Bone1"
	bone1.Layer = 0

	bone2 := NewBone()
	bone2.index = 1
	bone2.name = "Bone2"
	bone2.Layer = 0

	bone3 := NewBone()
	bone3.index = 2
	bone3.name = "Bone3"
	bone3.Layer = 0

	// Append initial bones
	bones.Append(bone1)
	bones.Append(bone2)
	bones.Setup()

	bone3.ParentIndex = bone1.Index()

	// Insert bone3 after bone1
	bones.Insert(bone3)

	{
		tests := []struct {
			name          string
			expectedLayer int
			expectedIndex int
		}{
			{"Bone1", 0, 0},
			{"Bone3", 0, 2},
			{"Bone2", 1, 1},
		}

		testGroup := "add3)"
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bone, _ := bones.GetByName(test.name)
				if bone == nil {
					t.Errorf("%s Expected %s to be found", testGroup, test.name)
					return
				}
				if bone.Layer != test.expectedLayer {
					t.Errorf("%s Expected %s Layer to be %d, got %d", testGroup, test.name, test.expectedLayer, bone.Layer)
				}
				if bone.Index() != test.expectedIndex {
					t.Errorf("%s Expected %s Index to be %d, got %d", testGroup, test.name, test.expectedIndex, bone.Index())
				}
			})
		}
	}

	// Insert bone3 after bone2
	bone4 := NewBone()
	bone4.index = 3
	bone4.name = "Bone4"
	bone4.Layer = 0

	bone4.ParentIndex = bone2.Index()

	bones.Insert(bone4)

	{
		tests := []struct {
			name          string
			expectedLayer int
			expectedIndex int
		}{
			{"Bone1", 0, 0},
			{"Bone3", 0, 2},
			{"Bone2", 1, 1},
			{"Bone4", 1, 3},
		}

		testGroup := "add4)"
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bone, _ := bones.GetByName(test.name)
				if bone == nil {
					t.Errorf("%s Expected %s to be found", testGroup, test.name)
					return
				}
				if bone.Layer != test.expectedLayer {
					t.Errorf("%s Expected %s Layer to be %d, got %d", testGroup, test.name, test.expectedLayer, bone.Layer)
				}
				if bone.Index() != test.expectedIndex {
					t.Errorf("%s Expected %s Index to be %d, got %d", testGroup, test.name, test.expectedIndex, bone.Index())
				}
			})
		}
	}

	// Insert bone5 at the end
	bone5 := NewBone()
	bone5.index = 4
	bone5.name = "Bone5"
	bone5.Layer = 0

	bone5.ParentIndex = bone4.Index()

	bones.Insert(bone5)

	{
		tests := []struct {
			name          string
			expectedLayer int
			expectedIndex int
		}{
			{"Bone1", 0, 0},
			{"Bone3", 0, 2},
			{"Bone2", 1, 1},
			{"Bone4", 1, 3},
			{"Bone5", 1, 4},
		}

		testGroup := "add5)"
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bone, _ := bones.GetByName(test.name)
				if bone == nil {
					t.Errorf("%s Expected %s to be found", testGroup, test.name)
					return
				}
				if bone.Layer != test.expectedLayer {
					t.Errorf("%s Expected %s Layer to be %d, got %d", testGroup, test.name, test.expectedLayer, bone.Layer)
				}
				if bone.Index() != test.expectedIndex {
					t.Errorf("%s Expected %s Index to be %d, got %d", testGroup, test.name, test.expectedIndex, bone.Index())
				}
			})
		}
	}

	// Insert bone6 after bone3
	bone6 := NewBone()
	bone6.index = 5
	bone6.name = "Bone6"
	bone6.Layer = 0

	bone6.ParentIndex = bone3.Index()

	bones.Insert(bone6)

	{
		tests := []struct {
			name          string
			expectedLayer int
			expectedIndex int
		}{
			{"Bone1", 0, 0},
			{"Bone3", 0, 2},
			{"Bone6", 1, 5},
			{"Bone2", 1, 1},
			{"Bone4", 1, 3},
			{"Bone5", 1, 4},
		}

		testGroup := "add6)"
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bone, _ := bones.GetByName(test.name)
				if bone == nil {
					t.Errorf("%s Expected %s to be found", testGroup, test.name)
					return
				}
				if bone.Layer != test.expectedLayer {
					t.Errorf("%s Expected %s Layer to be %d, got %d", testGroup, test.name, test.expectedLayer, bone.Layer)
				}
				if bone.Index() != test.expectedIndex {
					t.Errorf("%s Expected %s Index to be %d, got %d", testGroup, test.name, test.expectedIndex, bone.Index())
				}
			})
		}
	}

	// Insert bone7 after bone4
	bone7 := NewBone()
	bone7.index = 6
	bone7.name = "Bone7"
	bone7.Layer = 0

	bone7.ParentIndex = bone4.Index()
	bones.Insert(bone7)

	{
		tests := []struct {
			name          string
			expectedLayer int
			expectedIndex int
		}{
			{"Bone1", 0, 0},
			{"Bone3", 0, 2},
			{"Bone6", 1, 5},
			{"Bone2", 1, 1},
			{"Bone4", 1, 3},
			{"Bone7", 1, 6},
			{"Bone5", 1, 4},
		}

		testGroup := "add7)"
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bone, _ := bones.GetByName(test.name)
				if bone == nil {
					t.Errorf("%s Expected %s to be found", testGroup, test.name)
					return
				}
				if bone.Layer != test.expectedLayer {
					t.Errorf("%s Expected %s Layer to be %d, got %d", testGroup, test.name, test.expectedLayer, bone.Layer)
				}
				if bone.Index() != test.expectedIndex {
					t.Errorf("%s Expected %s Index to be %d, got %d", testGroup, test.name, test.expectedIndex, bone.Index())
				}
			})
		}
	}

	// Insert bone8 root
	bone8 := NewBone()
	bone8.index = 7
	bone8.name = "Bone8"
	bone8.Layer = 0

	bones.Insert(bone8)

	{
		tests := []struct {
			name          string
			expectedLayer int
			expectedIndex int
		}{
			{"Bone8", 0, 7},
			{"Bone1", 1, 0},
			{"Bone3", 0, 2},
			{"Bone6", 1, 5},
			{"Bone2", 2, 1},
			{"Bone4", 1, 3},
			{"Bone7", 1, 6},
			{"Bone5", 1, 4},
		}

		testGroup := "add8)"
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				bone, _ := bones.GetByName(test.name)
				if bone == nil {
					t.Errorf("%s Expected %s to be found", testGroup, test.name)
					return
				}
				if bone.Layer != test.expectedLayer {
					t.Errorf("%s Expected %s Layer to be %d, got %d", testGroup, test.name, test.expectedLayer, bone.Layer)
				}
				if bone.Index() != test.expectedIndex {
					t.Errorf("%s Expected %s Index to be %d, got %d", testGroup, test.name, test.expectedIndex, bone.Index())
				}
			})
		}
	}

}
