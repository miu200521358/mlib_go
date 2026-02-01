// 指示: miu200521358
package logging

import "testing"

// TestLogLevelConstants はログレベル定数を確認する。
func TestLogLevelConstants(t *testing.T) {
	if LOG_LEVEL_VERBOSE != 0 || LOG_LEVEL_DEBUG != 10 || LOG_LEVEL_INFO != 20 || LOG_LEVEL_WARN != 30 || LOG_LEVEL_ERROR != 40 || LOG_LEVEL_FATAL != 50 {
		t.Errorf("LogLevel constants mismatch")
	}
}

// TestVerboseIndexConstants は冗長ログ定数を確認する。
func TestVerboseIndexConstants(t *testing.T) {
	if VERBOSE_INDEX_MOTION != 230 || VERBOSE_INDEX_IK != 320 || VERBOSE_INDEX_PHYSICS != 330 || VERBOSE_INDEX_VIEWER != 540 {
		t.Errorf("VerboseIndex constants mismatch")
	}
}
