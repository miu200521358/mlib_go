package mmodel

import "testing"

func TestNewRigidBody(t *testing.T) {
	r := NewRigidBody()
	if r.Index() != -1 {
		t.Errorf("Index() = %v, want -1", r.Index())
	}
	if r.BoneIndex != -1 {
		t.Errorf("BoneIndex = %v, want -1", r.BoneIndex)
	}
	if r.ShapeType != SHAPE_BOX {
		t.Errorf("ShapeType = %v, want SHAPE_BOX", r.ShapeType)
	}
	if r.PhysicsType != PHYSICS_TYPE_STATIC {
		t.Errorf("PhysicsType = %v, want PHYSICS_TYPE_STATIC", r.PhysicsType)
	}
}

func TestRigidBody_IsValid(t *testing.T) {
	t.Run("新規作成は無効", func(t *testing.T) {
		r := NewRigidBody()
		if r.IsValid() {
			t.Errorf("IsValid() = true, want false")
		}
	})

	t.Run("インデックス設定後は有効", func(t *testing.T) {
		r := NewRigidBody()
		r.SetIndex(0)
		if !r.IsValid() {
			t.Errorf("IsValid() = false, want true")
		}
	})
}

func TestRigidBody_AsDynamic(t *testing.T) {
	tests := []struct {
		name        string
		physicsType PhysicsType
		expected    bool
	}{
		{"STATIC", PHYSICS_TYPE_STATIC, false},
		{"DYNAMIC", PHYSICS_TYPE_DYNAMIC, true},
		{"DYNAMIC_BONE", PHYSICS_TYPE_DYNAMIC_BONE, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRigidBody()
			r.PhysicsType = tt.physicsType
			if r.AsDynamic() != tt.expected {
				t.Errorf("AsDynamic() = %v, want %v", r.AsDynamic(), tt.expected)
			}
		})
	}
}

func TestRigidBodyParam(t *testing.T) {
	p := NewRigidBodyParam()
	if p.Mass != 1 {
		t.Errorf("Mass = %v, want 1", p.Mass)
	}
	if p.LinearDamping != 0.5 {
		t.Errorf("LinearDamping = %v, want 0.5", p.LinearDamping)
	}

	s := p.String()
	if s == "" {
		t.Errorf("String() should not be empty")
	}
}

func TestCollisionGroup(t *testing.T) {
	t.Run("全衝突グループ", func(t *testing.T) {
		cg := NewCollisionGroupAll()
		if len(cg.IsCollisions) != 16 {
			t.Errorf("IsCollisions length = %v, want 16", len(cg.IsCollisions))
		}
	})

	t.Run("スライスから生成", func(t *testing.T) {
		slice := []uint16{1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		cg := NewCollisionGroupFromSlice(slice)
		if len(cg.IsCollisions) != 16 {
			t.Errorf("IsCollisions length = %v, want 16", len(cg.IsCollisions))
		}
	})
}

func TestRigidBody_Copy(t *testing.T) {
	r := NewRigidBody()
	r.SetIndex(5)
	r.SetName("テスト剛体")
	r.BoneIndex = 10
	r.PhysicsType = PHYSICS_TYPE_DYNAMIC
	r.Bone = NewBone()

	cp, err := r.Copy()
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	if cp.Index() != 5 {
		t.Errorf("Copy() Index = %v, want 5", cp.Index())
	}
	if cp.Name() != "テスト剛体" {
		t.Errorf("Copy() Name = %v, want テスト剛体", cp.Name())
	}
	if cp.BoneIndex != 10 {
		t.Errorf("Copy() BoneIndex = %v, want 10", cp.BoneIndex)
	}
	if cp.Bone != nil {
		t.Errorf("Copy() Bone should be nil")
	}
}
