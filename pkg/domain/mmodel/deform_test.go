package mmodel

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

func TestBdef1(t *testing.T) {
	t.Run("生成とPacked", func(t *testing.T) {
		b := NewBdef1(5)
		if b.Type() != DEFORM_BDEF1 {
			t.Errorf("Type() = %v, want DEFORM_BDEF1", b.Type())
		}
		packed := b.Packed()
		if packed[0] != 5 || packed[4] != 1.0 {
			t.Errorf("Packed() = %v", packed)
		}
	})

	t.Run("Indexes/Weights", func(t *testing.T) {
		b := NewBdef1(3)
		if len(b.Indexes()) != 1 || b.Indexes()[0] != 3 {
			t.Errorf("Indexes() = %v", b.Indexes())
		}
		if len(b.Weights()) != 1 || b.Weights()[0] != 1.0 {
			t.Errorf("Weights() = %v", b.Weights())
		}
	})
}

func TestBdef2(t *testing.T) {
	t.Run("生成とPacked", func(t *testing.T) {
		b := NewBdef2(1, 2, 0.7)
		if b.Type() != DEFORM_BDEF2 {
			t.Errorf("Type() = %v, want DEFORM_BDEF2", b.Type())
		}
		packed := b.Packed()
		if packed[0] != 1 || packed[1] != 2 {
			t.Errorf("Packed() indexes = %v", packed)
		}
		if packed[4] != 0.7 || packed[5] != 0.3 {
			t.Errorf("Packed() weights = %v", packed)
		}
	})
}

func TestBdef4(t *testing.T) {
	t.Run("生成とPacked", func(t *testing.T) {
		b := NewBdef4(0, 1, 2, 3, 0.4, 0.3, 0.2, 0.1)
		if b.Type() != DEFORM_BDEF4 {
			t.Errorf("Type() = %v, want DEFORM_BDEF4", b.Type())
		}
		packed := b.Packed()
		if packed[0] != 0 || packed[1] != 1 || packed[2] != 2 || packed[3] != 3 {
			t.Errorf("Packed() indexes = %v", packed)
		}
	})
}

func TestSdef(t *testing.T) {
	t.Run("生成", func(t *testing.T) {
		c := mmath.NewVec3ByValues(1, 2, 3)
		r0 := mmath.NewVec3ByValues(4, 5, 6)
		r1 := mmath.NewVec3ByValues(7, 8, 9)
		s := NewSdef(0, 1, 0.6, c, r0, r1)
		if s.Type() != DEFORM_SDEF {
			t.Errorf("Type() = %v, want DEFORM_SDEF", s.Type())
		}
		if s.C != c || s.R0 != r0 || s.R1 != r1 {
			t.Errorf("SDEF params mismatch")
		}
	})
}

func TestDeform_Index(t *testing.T) {
	b := NewBdef4(10, 20, 30, 40, 0.25, 0.25, 0.25, 0.25)
	if b.Index(20) != 1 {
		t.Errorf("Index(20) = %v, want 1", b.Index(20))
	}
	if b.Index(99) != -1 {
		t.Errorf("Index(99) = %v, want -1", b.Index(99))
	}
}

func TestDeform_IndexWeight(t *testing.T) {
	b := NewBdef2(5, 10, 0.8)
	w5 := b.IndexWeight(5)
	if w5 < 0.79 || w5 > 0.81 {
		t.Errorf("IndexWeight(5) = %v, want ~0.8", w5)
	}
	w10 := b.IndexWeight(10)
	if w10 < 0.19 || w10 > 0.21 {
		t.Errorf("IndexWeight(10) = %v, want ~0.2", w10)
	}
	if b.IndexWeight(99) != 0 {
		t.Errorf("IndexWeight(99) = %v, want 0", b.IndexWeight(99))
	}
}

func TestDeform_IndexesByWeight(t *testing.T) {
	b := NewBdef4(0, 1, 2, 3, 0.5, 0.3, 0.15, 0.05)
	indexes := b.IndexesByWeight(0.2)
	if len(indexes) != 2 {
		t.Errorf("IndexesByWeight(0.2) length = %v, want 2", len(indexes))
	}
}

func TestDeform_WeightsByWeight(t *testing.T) {
	b := NewBdef4(0, 1, 2, 3, 0.5, 0.3, 0.15, 0.05)
	weights := b.WeightsByWeight(0.2)
	if len(weights) != 2 {
		t.Errorf("WeightsByWeight(0.2) length = %v, want 2", len(weights))
	}
}

func TestDeform_Normalize(t *testing.T) {
	t.Run("単純正規化", func(t *testing.T) {
		b := NewBdef2(0, 1, 0.6)
		b.Normalize(false)
		sum := b.Weights()[0] + b.Weights()[1]
		if sum < 0.99 || sum > 1.01 {
			t.Errorf("Normalize sum = %v, want 1.0", sum)
		}
	})
}

func TestDeform_SetIndexes(t *testing.T) {
	b := NewBdef1(0)
	b.SetIndexes([]int{5, 6})
	if b.Indexes()[0] != 5 || b.Indexes()[1] != 6 {
		t.Errorf("SetIndexes failed")
	}
}
