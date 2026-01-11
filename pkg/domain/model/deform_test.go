package model

import "testing"

func TestDeformConstructors(t *testing.T) {
	b1 := NewBdef1(-1)
	if b1.DeformType() != BDEF1 {
		t.Fatalf("Bdef1 type = %v", b1.DeformType())
	}
	if len(b1.Indexes()) != 1 || b1.Indexes()[0] != -1 {
		t.Fatalf("Bdef1 indexes = %v", b1.Indexes())
	}
	if len(b1.Weights()) != 1 || b1.Weights()[0] != 1.0 {
		t.Fatalf("Bdef1 weights = %v", b1.Weights())
	}

	b2 := NewBdef2(1, 2, 0.25)
	if b2.DeformType() != BDEF2 {
		t.Fatalf("Bdef2 type = %v", b2.DeformType())
	}
	if got := b2.Indexes(); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("Bdef2 indexes = %v", got)
	}
	if got := b2.Weights(); len(got) != 2 || got[0] != 0.25 || got[1] != 0.75 {
		t.Fatalf("Bdef2 weights = %v", got)
	}

	b4 := NewBdef4([4]int{1, 2, 3, 4}, [4]float64{0.1, 0.2, 0.3, 0.4})
	if b4.DeformType() != BDEF4 {
		t.Fatalf("Bdef4 type = %v", b4.DeformType())
	}
	if got := b4.Indexes(); len(got) != 4 || got[3] != 4 {
		t.Fatalf("Bdef4 indexes = %v", got)
	}
	if got := b4.Weights(); len(got) != 4 || got[3] != 0.4 {
		t.Fatalf("Bdef4 weights = %v", got)
	}

	s := NewSdef(7, 8, 0.6)
	if s.DeformType() != SDEF {
		t.Fatalf("Sdef type = %v", s.DeformType())
	}
	if got := s.Indexes(); len(got) != 2 || got[0] != 7 || got[1] != 8 {
		t.Fatalf("Sdef indexes = %v", got)
	}
	if got := s.Weights(); len(got) != 2 || got[0] != 0.6 || got[1] != 0.4 {
		t.Fatalf("Sdef weights = %v", got)
	}
}
