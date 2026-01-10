// 指示: miu200521358
package errorregistry

import "testing"

// TestLoadDefaultRegistry は埋め込みCSVの読込を確認する。
func TestLoadDefaultRegistry(t *testing.T) {
	recs, err := LoadDefaultRegistry()
	if err != nil {
		t.Fatalf("LoadDefaultRegistry failed: %v", err)
	}
	if len(recs) == 0 {
		t.Fatalf("LoadDefaultRegistry returned empty")
	}
	if recs[0].ID == "" {
		t.Errorf("first record ID empty")
	}
}
