// 指示: miu200521358
package motion

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
)

// TestBoneFrameLerpLinear は線形補間の挙動を確認する。
func TestBoneFrameLerpLinear(t *testing.T) {
	prev := NewBoneFrame(0)
	prev.Position = vec3Ptr(0, 0, 0)
	prev.Rotation = ptrQuat(mmath.NewQuaternionFromDegrees(0, 0, 0))
	prev.UnitRotation = ptrQuat(mmath.NewQuaternionFromDegrees(0, 0, 0))
	prev.Scale = vec3Ptr(1, 1, 1)

	next := NewBoneFrame(10)
	next.Position = vec3Ptr(10, 0, 0)
	next.Rotation = ptrQuat(mmath.NewQuaternionFromDegrees(0, 90, 0))
	next.UnitRotation = ptrQuat(mmath.NewQuaternionFromDegrees(0, 90, 0))
	next.Scale = vec3Ptr(1, 1, 1)

	out := next.lerpFrame(prev, 5)
	if out == nil || out.Position == nil {
		t.Fatalf("lerpFrame nil")
	}
	if out.Position.X < 4.9 || out.Position.X > 5.1 {
		t.Fatalf("Position lerp: got=%v", out.Position.X)
	}
	if out.Rotation == nil || !out.Rotation.NearEquals(prev.Rotation.Slerp(*next.Rotation, 0.5), 1e-6) {
		t.Fatalf("Rotation lerp mismatch")
	}
}

// TestBoneFrameSplitCurve は曲線分割の有無を確認する。
func TestBoneFrameSplitCurve(t *testing.T) {
	prev := NewBoneFrame(0)
	next := NewBoneFrame(10)
	next.Curves = NewBoneCurves()
	insert := NewBoneFrame(5)
	insert.splitCurve(prev, next, 5)
	if insert.Curves == nil {
		t.Fatalf("splitCurve should set curves")
	}
}

// TestBoneNameFramesContainsActive は有効判定を確認する。
func TestBoneNameFramesContainsActive(t *testing.T) {
	frames := NewBoneNameFrames("b")
	frames.Append(NewBoneFrame(0))
	if frames.ContainsActive() {
		t.Fatalf("ContainsActive should be false")
	}

	bf := NewBoneFrame(1)
	bf.Position = vec3Ptr(1, 0, 0)
	frames.Append(bf)
	if !frames.ContainsActive() {
		t.Fatalf("ContainsActive should be true")
	}
}

// TestBoneFramesClean はCleanで削除されることを確認する。
func TestBoneFramesClean(t *testing.T) {
	bones := NewBoneFrames()
	bones.Get("b").Append(NewBoneFrame(0))
	bones.Clean()
	if len(bones.Names()) != 0 {
		t.Fatalf("Clean should remove empty frames")
	}
}

// TestBoneNameFramesReduceNoInflection は変曲点が少ない場合の挙動を確認する。
func TestBoneNameFramesReduceNoInflection(t *testing.T) {
	frames := NewBoneNameFrames("b")
	frames.Append(NewBoneFrame(0))
	frames.Append(NewBoneFrame(1))
	reduced, err := frames.Reduce()
	if err != nil {
		t.Fatalf("Reduce error: %v", err)
	}
	if reduced != frames {
		t.Fatalf("Reduce should return original when no inflection")
	}
}

// TestBoneNameFramesReduceRange は削減処理の区間成功を確認する。
func TestBoneNameFramesReduceRange(t *testing.T) {
	frames := NewBoneNameFrames("b")
	bf0 := NewBoneFrame(0)
	bf0.Position = vec3Ptr(0, 0, 0)
	bf2 := NewBoneFrame(2)
	bf2.Position = vec3Ptr(2, 0, 0)
	frames.Append(bf0)
	frames.Append(bf2)

	xs := []float64{0, 1, 2}
	ys := []float64{0, 0, 0}
	zs := []float64{0, 0, 0}
	quats := []mmath.Quaternion{mmath.NewQuaternion(), mmath.NewQuaternion(), mmath.NewQuaternion()}

	reduced := NewBoneNameFrames("b")
	end, err := frames.reduceRange(0, 1, 2, xs, ys, zs, quats, reduced)
	if err != nil {
		t.Fatalf("reduceRange error: %v", err)
	}
	if end != 2 {
		t.Fatalf("reduceRange end: got=%v", end)
	}
	if reduced.Len() == 0 {
		t.Fatalf("reduceRange should append")
	}
}

// TestFindInflectionFrames は変曲点抽出を確認する。
func TestFindInflectionFrames(t *testing.T) {
	frames := []Frame{0, 1, 2, 3}
	values := []float64{0, 1, 0, 1}
	found := findInflectionFrames(frames, values, 1e-4)
	if len(found) <= 2 {
		t.Fatalf("inflection not detected: %v", found)
	}
}

// ptrQuat はQuaternionのポインタを返す。
func ptrQuat(q mmath.Quaternion) *mmath.Quaternion {
	return &q
}
