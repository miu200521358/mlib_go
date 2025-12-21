package mmodel

import "testing"

func TestNewJoint(t *testing.T) {
	j := NewJoint()
	if j.Index() != -1 {
		t.Errorf("Index() = %v, want -1", j.Index())
	}
	if j.JointType != 0 {
		t.Errorf("JointType = %v, want 0", j.JointType)
	}
	if j.RigidBodyIndexA != -1 {
		t.Errorf("RigidBodyIndexA = %v, want -1", j.RigidBodyIndexA)
	}
	if j.RigidBodyIndexB != -1 {
		t.Errorf("RigidBodyIndexB = %v, want -1", j.RigidBodyIndexB)
	}
	if j.JointParam == nil {
		t.Errorf("JointParam should not be nil")
	}
}

func TestNewJointByName(t *testing.T) {
	j := NewJointByName("テストジョイント")
	if j.Name() != "テストジョイント" {
		t.Errorf("Name() = %v, want テストジョイント", j.Name())
	}
}

func TestJoint_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		j := NewJoint()
		if j.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		j := NewJoint()
		j.SetIndex(0)
		if !j.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestJointParam(t *testing.T) {
	p := NewJointParam()
	if p.TranslationLimitMin == nil {
		t.Errorf("TranslationLimitMin should not be nil")
	}
	if p.RotationLimitMax == nil {
		t.Errorf("RotationLimitMax should not be nil")
	}

	s := p.String()
	if s == "" {
		t.Errorf("String() should not be empty")
	}
}

func TestJoint_Copy(t *testing.T) {
	j := NewJoint()
	j.SetIndex(5)
	j.SetName("テスト")
	j.RigidBodyIndexA = 1
	j.RigidBodyIndexB = 2

	cp, err := j.Copy()
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	if cp.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", cp.Index())
	}
	if cp.Name() != "テスト" {
		t.Errorf("Copy() Name = %v, want テスト", cp.Name())
	}
	if cp.RigidBodyIndexA != 1 {
		t.Errorf("Copy() RigidBodyIndexA = %v, want 1", cp.RigidBodyIndexA)
	}

	// 独立性確認
	j.JointParam.TranslationLimitMin.X = 100
	if cp.JointParam.TranslationLimitMin.X == 100 {
		t.Errorf("JointParam should be independent")
	}
}
