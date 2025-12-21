package mmodel

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestNewIkLink(t *testing.T) {
	l := NewIkLink()
	if l.BoneIndex != -1 {
		t.Errorf("BoneIndex = %v, want -1", l.BoneIndex)
	}
	if l.AngleLimit {
		t.Errorf("AngleLimit = true, want false")
	}
	if l.MinAngleLimit == nil || l.MaxAngleLimit == nil {
		t.Errorf("angle limits should not be nil")
	}
}

func TestIkLink_Copy(t *testing.T) {
	t.Run("基本コピー", func(t *testing.T) {
		l := NewIkLink()
		l.BoneIndex = 5
		l.AngleLimit = true
		l.MinAngleLimit = mmath.NewVec3ByValues(-1, -2, -3)
		l.MaxAngleLimit = mmath.NewVec3ByValues(1, 2, 3)

		cp, err := l.Copy()
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}
		if cp.BoneIndex != 5 {
			t.Errorf("BoneIndex = %v, want 5", cp.BoneIndex)
		}
		if !cp.AngleLimit {
			t.Errorf("AngleLimit = false, want true")
		}
	})

	t.Run("別オブジェクト確認", func(t *testing.T) {
		l := NewIkLink()
		l.MinAngleLimit = mmath.NewVec3ByValues(-1, -2, -3)

		cp, _ := l.Copy()

		if l.MinAngleLimit == cp.MinAngleLimit {
			t.Errorf("MinAngleLimit pointer should be different")
		}

		l.MinAngleLimit.X = 100
		if cp.MinAngleLimit.X == 100 {
			t.Errorf("MinAngleLimit should be independent")
		}
	})
}

func TestNewIk(t *testing.T) {
	ik := NewIk()
	if ik.BoneIndex != -1 {
		t.Errorf("BoneIndex = %v, want -1", ik.BoneIndex)
	}
	if ik.LoopCount != 0 {
		t.Errorf("LoopCount = %v, want 0", ik.LoopCount)
	}
	if ik.UnitRotation == nil {
		t.Errorf("UnitRotation should not be nil")
	}
	if len(ik.Links) != 0 {
		t.Errorf("Links length = %v, want 0", len(ik.Links))
	}
}

func TestIk_Copy(t *testing.T) {
	t.Run("基本コピー", func(t *testing.T) {
		ik := NewIk()
		ik.BoneIndex = 10
		ik.LoopCount = 20
		ik.UnitRotation = mmath.NewVec3ByValues(0.1, 0.2, 0.3)
		ik.Links = append(ik.Links, NewIkLink())
		ik.Links[0].BoneIndex = 5

		cp, err := ik.Copy()
		if err != nil {
			t.Fatalf("Copy() error = %v", err)
		}
		if cp.BoneIndex != 10 {
			t.Errorf("BoneIndex = %v, want 10", cp.BoneIndex)
		}
		if cp.LoopCount != 20 {
			t.Errorf("LoopCount = %v, want 20", cp.LoopCount)
		}
		if len(cp.Links) != 1 {
			t.Fatalf("Links length = %v, want 1", len(cp.Links))
		}
		if cp.Links[0].BoneIndex != 5 {
			t.Errorf("Links[0].BoneIndex = %v, want 5", cp.Links[0].BoneIndex)
		}
	})

	t.Run("別オブジェクト確認_UnitRotation", func(t *testing.T) {
		ik := NewIk()
		ik.UnitRotation = mmath.NewVec3ByValues(0.1, 0.2, 0.3)

		cp, _ := ik.Copy()

		if ik.UnitRotation == cp.UnitRotation {
			t.Errorf("UnitRotation pointer should be different")
		}
	})

	t.Run("別オブジェクト確認_Links", func(t *testing.T) {
		ik := NewIk()
		ik.Links = append(ik.Links, NewIkLink())

		cp, _ := ik.Copy()

		if ik.Links[0] == cp.Links[0] {
			t.Errorf("Links[0] pointer should be different")
		}

		ik.Links[0].BoneIndex = 999
		if cp.Links[0].BoneIndex == 999 {
			t.Errorf("Links should be independent")
		}
	})
}
