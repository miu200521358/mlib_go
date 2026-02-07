// 指示: miu200521358
package motion

import "testing"

// TestVmdMotionIsVpd はVPD判定を確認する。
func TestVmdMotionIsVpd(t *testing.T) {
	motion := NewVmdMotion("sample.VPD")
	if !motion.IsVpd() {
		t.Fatalf("IsVpd should be true")
	}
}

// TestVmdMotionUpdateHash はハッシュ更新条件を確認する。
func TestVmdMotionUpdateHash(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.SetName("n")
	motion.UpdateHash()
	baseHash := motion.Hash()

	motion.SetFileModTime(123)
	motion.UpdateHash()
	if motion.Hash() == baseHash {
		t.Fatalf("Hash should change by mod time")
	}

	motion.BoneFrames.Get("b").Append(NewBoneFrame(0))
	motion.UpdateHash()
	hashWithFrames := motion.Hash()
	motion.MorphFrames.Get("m").Append(NewMorphFrame(1))
	motion.UpdateHash()
	if motion.Hash() == hashWithFrames {
		t.Fatalf("Hash should change by hash parts")
	}
}

// TestVmdMotionIndexes はIndexesの対象を確認する。
func TestVmdMotionIndexes(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.BoneFrames.Get("b").Append(NewBoneFrame(0))
	motion.MorphFrames.Get("m").Append(NewMorphFrame(10))
	motion.CameraFrames.Append(NewCameraFrame(20))
	indexes := motion.Indexes()
	if len(indexes) != 3 || indexes[0] != 0 || indexes[1] != 10 || indexes[2] != 20 {
		t.Fatalf("Indexes: got=%v", indexes)
	}
}

// TestVmdMotionMinFrame は最小フレームの判定を確認する。
func TestVmdMotionMinFrame(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.BoneFrames.Get("b").Append(NewBoneFrame(0))
	motion.MorphFrames.Get("m").Append(NewMorphFrame(5))
	if motion.MinFrame() != 0 {
		t.Fatalf("MinFrame: got=%v", motion.MinFrame())
	}
}

// TestVmdMotionMaxFrameWithCameraOnly はカメラのみの最大フレームを確認する。
func TestVmdMotionMaxFrameWithCameraOnly(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.CameraFrames.Append(NewCameraFrame(120))
	if motion.MaxFrame() != 120 {
		t.Fatalf("MaxFrame(camera only): got=%v", motion.MaxFrame())
	}
}

// TestVmdMotionMinFrameWithCameraOnly はカメラのみの最小フレームを確認する。
func TestVmdMotionMinFrameWithCameraOnly(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.CameraFrames.Append(NewCameraFrame(12))
	if motion.MinFrame() != 12 {
		t.Fatalf("MinFrame(camera only): got=%v", motion.MinFrame())
	}
}

// TestVmdMotionCopy はCopyの挙動を確認する。
func TestVmdMotionCopy(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.BoneFrames.Get("b").Append(NewBoneFrame(0))
	motion.UpdateHash()
	copied, err := motion.Copy()
	if err != nil {
		t.Fatalf("Copy error: %v", err)
	}
	if copied.Hash() == "" {
		t.Fatalf("Copy hash empty")
	}
	copied.BoneFrames.Get("b").Append(NewBoneFrame(1))
	if motion.BoneFrames.Len() == copied.BoneFrames.Len() {
		t.Fatalf("Copy should be deep")
	}
}

// TestVmdMotionClean はCleanの挙動を確認する。
func TestVmdMotionClean(t *testing.T) {
	motion := NewVmdMotion("path.vmd")
	motion.CameraFrames.Append(NewCameraFrame(0))
	motion.Clean()
	if motion.CameraFrames.Len() != 0 {
		t.Fatalf("Clean should clear default camera")
	}
}
