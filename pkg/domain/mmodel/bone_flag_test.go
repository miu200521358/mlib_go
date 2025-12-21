package mmodel

import "testing"

func TestBoneFlag(t *testing.T) {
	t.Run("フラグ設定と取得", func(t *testing.T) {
		f := BONE_FLAG_NONE

		f = f.SetCanRotate(true)
		if !f.CanRotate() {
			t.Errorf("CanRotate() = false, want true")
		}

		f = f.SetCanTranslate(true)
		if !f.CanTranslate() {
			t.Errorf("CanTranslate() = false, want true")
		}

		f = f.SetIsVisible(true)
		if !f.IsVisible() {
			t.Errorf("IsVisible() = false, want true")
		}

		f = f.SetCanRotate(false)
		if f.CanRotate() {
			t.Errorf("CanRotate() = true, want false")
		}
		if !f.CanTranslate() {
			t.Errorf("CanTranslate() should still be true")
		}
	})

	t.Run("IKフラグ", func(t *testing.T) {
		f := BONE_FLAG_NONE

		f = f.SetIsIK(true)
		if !f.IsIK() {
			t.Errorf("IsIK() = false, want true")
		}
	})

	t.Run("付与フラグ", func(t *testing.T) {
		f := BONE_FLAG_NONE

		f = f.SetExternalRotation(true)
		if !f.IsExternalRotation() {
			t.Errorf("IsExternalRotation() = false, want true")
		}

		f = f.SetExternalTranslation(true)
		if !f.IsExternalTranslation() {
			t.Errorf("IsExternalTranslation() = false, want true")
		}
	})

	t.Run("軸フラグ", func(t *testing.T) {
		f := BONE_FLAG_NONE

		f = f.SetHasFixedAxis(true)
		if !f.HasFixedAxis() {
			t.Errorf("HasFixedAxis() = false, want true")
		}

		f = f.SetHasLocalAxis(true)
		if !f.HasLocalAxis() {
			t.Errorf("HasLocalAxis() = false, want true")
		}
	})

	t.Run("物理後変形フラグ", func(t *testing.T) {
		f := BONE_FLAG_NONE

		f = f.SetAfterPhysicsDeform(true)
		if !f.IsAfterPhysicsDeform() {
			t.Errorf("IsAfterPhysicsDeform() = false, want true")
		}
	})

	t.Run("複合フラグ", func(t *testing.T) {
		f := BONE_FLAG_CAN_ROTATE | BONE_FLAG_CAN_TRANSLATE | BONE_FLAG_IS_VISIBLE | BONE_FLAG_CAN_MANIPULATE
		if !f.CanRotate() {
			t.Errorf("CanRotate() = false, want true")
		}
		if !f.CanTranslate() {
			t.Errorf("CanTranslate() = false, want true")
		}
		if !f.IsVisible() {
			t.Errorf("IsVisible() = false, want true")
		}
		if !f.CanManipulate() {
			t.Errorf("CanManipulate() = false, want true")
		}
		if f.IsIK() {
			t.Errorf("IsIK() = true, want false")
		}
	})
}
