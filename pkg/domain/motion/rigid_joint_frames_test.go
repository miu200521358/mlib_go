// 指示: miu200521358
package motion

import "testing"

// TestRigidBodyFramesGetNil は未登録時にnilを返すことを確認する。
func TestRigidBodyFramesGetNil(t *testing.T) {
	frames := NewRigidBodyFrames()
	if frames.Get("rb") != nil {
		t.Fatalf("Get should return nil")
	}
	nameFrames := NewRigidBodyNameFrames("rb")
	frames.Update(nameFrames)
	if frames.Get("rb") == nil {
		t.Fatalf("Get should return frames")
	}
}

// TestRigidBodyFrameLerp は剛体補間を確認する。
func TestRigidBodyFrameLerp(t *testing.T) {
	prev := NewRigidBodyFrame(0)
	prev.Position = vec3Ptr(0, 0, 0)
	prev.Size = vec3Ptr(1, 1, 1)
	prevMass := 1.0
	prev.Mass = &prevMass
	next := NewRigidBodyFrame(10)
	next.Position = vec3Ptr(10, 0, 0)
	next.Size = vec3Ptr(2, 2, 2)
	nextMass := 3.0
	next.Mass = &nextMass
	out := next.lerpFrame(prev, 5)
	if out.Position == nil || out.Position.X < 4.9 || out.Position.X > 5.1 {
		t.Fatalf("Position lerp")
	}
	if out.Size == nil || out.Size.X < 1.4 || out.Size.X > 1.6 {
		t.Fatalf("Size lerp")
	}
	if out.Mass == nil || *out.Mass < 1.9 || *out.Mass > 2.1 {
		t.Fatalf("Mass lerp")
	}
}

// TestJointFrameLerp はジョイント補間を確認する。
func TestJointFrameLerp(t *testing.T) {
	prev := NewJointFrame(0)
	prev.TranslationLimitMin = vec3Ptr(0, 0, 0)
	next := NewJointFrame(10)
	next.TranslationLimitMin = vec3Ptr(10, 0, 0)
	out := next.lerpFrame(prev, 5)
	if out.TranslationLimitMin == nil || out.TranslationLimitMin.X < 4.9 || out.TranslationLimitMin.X > 5.1 {
		t.Fatalf("Joint lerp")
	}
}
