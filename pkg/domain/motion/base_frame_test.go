// 指示: miu200521358
package motion

import "testing"

// TestBaseFrameMethods はBaseFrameの基本動作を確認する。
func TestBaseFrameMethods(t *testing.T) {
	f := NewBaseFrame(1)
	if f.Index() != 1 {
		t.Fatalf("Index: got=%v", f.Index())
	}
	copied, err := f.Copy()
	if err != nil {
		t.Fatalf("Copy error: %v", err)
	}
	if copied.Index() != 1 || copied.Read != f.Read {
		t.Fatalf("Copy mismatch")
	}
}
