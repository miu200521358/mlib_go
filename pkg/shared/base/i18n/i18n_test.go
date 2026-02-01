// 指示: miu200521358
package i18n

import "testing"

// TestDefaultLang は既定言語を確認する。
func TestDefaultLang(t *testing.T) {
	if DefaultLang != "ja" {
		t.Errorf("DefaultLang: got=%v", DefaultLang)
	}
}
