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
