// 指示: miu200521358
package base

import "testing"

// TestDefaultBaseInitStages は初期化順序を確認する。
func TestDefaultBaseInitStages(t *testing.T) {
	if len(DEFAULT_BASE_INIT_STAGES) != 3 {
		t.Errorf("DEFAULT_BASE_INIT_STAGES length: got=%v", len(DEFAULT_BASE_INIT_STAGES))
	}
	if DEFAULT_BASE_INIT_STAGES[0] != BASE_INIT_CONFIG || DEFAULT_BASE_INIT_STAGES[1] != BASE_INIT_I18N || DEFAULT_BASE_INIT_STAGES[2] != BASE_INIT_LOGGING {
		t.Errorf("DEFAULT_BASE_INIT_STAGES order: got=%v", DEFAULT_BASE_INIT_STAGES)
	}
}

// TestBaseServicesAccessors はnil安全性を確認する。
func TestBaseServicesAccessors(t *testing.T) {
	var b *BaseServices
	if b.Config() != nil || b.I18n() != nil || b.Logger() != nil {
		t.Errorf("BaseServices nil accessors should return nil")
	}
}

// TestBaseServicesAccessorsNonNil は非nilの返却を確認する。
func TestBaseServicesAccessorsNonNil(t *testing.T) {
	b := &BaseServices{}
	if b.Config() != nil || b.I18n() != nil || b.Logger() != nil {
		t.Errorf("BaseServices empty fields should return nil")
	}
}
