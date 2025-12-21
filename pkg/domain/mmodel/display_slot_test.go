package mmodel

import "testing"

func TestNewDisplaySlot(t *testing.T) {
	d := NewDisplaySlot()
	if d.Index() != -1 {
		t.Errorf("Index() = %v, want -1", d.Index())
	}
	if d.SpecialFlag != SPECIAL_FLAG_OFF {
		t.Errorf("SpecialFlag = %v, want SPECIAL_FLAG_OFF", d.SpecialFlag)
	}
	if len(d.References) != 0 {
		t.Errorf("References length = %v, want 0", len(d.References))
	}
}

func TestNewRootDisplaySlot(t *testing.T) {
	d := NewRootDisplaySlot()
	if d.Index() != 0 {
		t.Errorf("Index() = %v, want 0", d.Index())
	}
	if d.Name() != "Root" {
		t.Errorf("Name() = %v, want Root", d.Name())
	}
	if d.SpecialFlag != SPECIAL_FLAG_ON {
		t.Errorf("SpecialFlag = %v, want SPECIAL_FLAG_ON", d.SpecialFlag)
	}
}

func TestNewMorphDisplaySlot(t *testing.T) {
	d := NewMorphDisplaySlot()
	if d.Index() != 1 {
		t.Errorf("Index() = %v, want 1", d.Index())
	}
	if d.Name() != "表情" {
		t.Errorf("Name() = %v, want 表情", d.Name())
	}
	if d.SpecialFlag != SPECIAL_FLAG_ON {
		t.Errorf("SpecialFlag = %v, want SPECIAL_FLAG_ON", d.SpecialFlag)
	}
}

func TestDisplaySlot_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		d := NewDisplaySlot()
		if d.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("Root表示枠は有効", func(t *testing.T) {
		d := NewRootDisplaySlot()
		if !d.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestReference(t *testing.T) {
	r := NewReference()
	if r.DisplayType != DISPLAY_TYPE_BONE {
		t.Errorf("DisplayType = %v, want DISPLAY_TYPE_BONE", r.DisplayType)
	}
	if r.DisplayIndex != -1 {
		t.Errorf("DisplayIndex = %v, want -1", r.DisplayIndex)
	}

	r2 := NewReferenceByValues(DISPLAY_TYPE_MORPH, 5)
	if r2.DisplayType != DISPLAY_TYPE_MORPH {
		t.Errorf("DisplayType = %v, want DISPLAY_TYPE_MORPH", r2.DisplayType)
	}
	if r2.DisplayIndex != 5 {
		t.Errorf("DisplayIndex = %v, want 5", r2.DisplayIndex)
	}
}

func TestDisplaySlot_Copy(t *testing.T) {
	d := NewDisplaySlot()
	d.SetIndex(5)
	d.SetName("テスト")
	d.SpecialFlag = SPECIAL_FLAG_OFF
	d.References = append(d.References, NewReferenceByValues(DISPLAY_TYPE_BONE, 0))
	d.References = append(d.References, NewReferenceByValues(DISPLAY_TYPE_MORPH, 1))

	cp, err := d.Copy()
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	if cp.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", cp.Index())
	}
	if cp.Name() != "テスト" {
		t.Errorf("Copy() Name = %v, want テスト", cp.Name())
	}
	if len(cp.References) != 2 {
		t.Errorf("Copy() References length = %v, want 2", len(cp.References))
	}

	// 独立性確認
	d.References[0].DisplayIndex = 100
	if cp.References[0].DisplayIndex == 100 {
		t.Errorf("References should be independent")
	}
}
