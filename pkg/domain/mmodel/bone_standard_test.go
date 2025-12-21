package mmodel

import "testing"

func TestBoneDirection(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		if BONE_DIRECTION_LEFT.String() != "左" {
			t.Errorf("String() = %v, want 左", BONE_DIRECTION_LEFT.String())
		}
		if BONE_DIRECTION_RIGHT.String() != "右" {
			t.Errorf("String() = %v, want 右", BONE_DIRECTION_RIGHT.String())
		}
		if BONE_DIRECTION_TRUNK.String() != "" {
			t.Errorf("String() = %v, want empty", BONE_DIRECTION_TRUNK.String())
		}
	})

	t.Run("Sign", func(t *testing.T) {
		if BONE_DIRECTION_LEFT.Sign() != -1.0 {
			t.Errorf("Sign() = %v, want -1", BONE_DIRECTION_LEFT.Sign())
		}
		if BONE_DIRECTION_RIGHT.Sign() != 1.0 {
			t.Errorf("Sign() = %v, want 1", BONE_DIRECTION_RIGHT.Sign())
		}
		if BONE_DIRECTION_TRUNK.Sign() != 0.0 {
			t.Errorf("Sign() = %v, want 0", BONE_DIRECTION_TRUNK.Sign())
		}
	})
}

func TestStandardBoneName(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		if BONE_CENTER.String() != "センター" {
			t.Errorf("String() = %v, want センター", BONE_CENTER.String())
		}
	})

	t.Run("StringFromDirection", func(t *testing.T) {
		name := BONE_ARM.StringFromDirection(BONE_DIRECTION_LEFT)
		if name != "左腕" {
			t.Errorf("StringFromDirection() = %v, want 左腕", name)
		}

		name = BONE_ARM.StringFromDirection(BONE_DIRECTION_RIGHT)
		if name != "右腕" {
			t.Errorf("StringFromDirection() = %v, want 右腕", name)
		}
	})

	t.Run("Left", func(t *testing.T) {
		if BONE_ARM.Left() != "左腕" {
			t.Errorf("Left() = %v, want 左腕", BONE_ARM.Left())
		}
		if BONE_ELBOW.Left() != "左ひじ" {
			t.Errorf("Left() = %v, want 左ひじ", BONE_ELBOW.Left())
		}
	})

	t.Run("Right", func(t *testing.T) {
		if BONE_ARM.Right() != "右腕" {
			t.Errorf("Right() = %v, want 右腕", BONE_ARM.Right())
		}
		if BONE_WRIST.Right() != "右手首" {
			t.Errorf("Right() = %v, want 右手首", BONE_WRIST.Right())
		}
	})

	t.Run("体幹ボーン", func(t *testing.T) {
		// 体幹ボーンはプレースホルダなし
		if BONE_CENTER.String() != "センター" {
			t.Errorf("String() = %v, want センター", BONE_CENTER.String())
		}
		if BONE_UPPER.String() != "上半身" {
			t.Errorf("String() = %v, want 上半身", BONE_UPPER.String())
		}
	})
}

func TestBoneCategory(t *testing.T) {
	t.Run("iota値確認", func(t *testing.T) {
		if CATEGORY_ROOT != 0 {
			t.Errorf("CATEGORY_ROOT = %v, want 0", CATEGORY_ROOT)
		}
		if CATEGORY_TRUNK != 1 {
			t.Errorf("CATEGORY_TRUNK = %v, want 1", CATEGORY_TRUNK)
		}
	})
}
