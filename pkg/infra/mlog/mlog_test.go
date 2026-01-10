// 指示: miu200521358
package mlog

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

type stubI18n struct {
	messages map[string]string
}

// Lang は言語を返す。
func (s *stubI18n) Lang() i18n.LangCode { return "ja" }

// SetLang は未使用のため固定値を返す。
func (s *stubI18n) SetLang(lang i18n.LangCode) i18n.LangChangeAction {
	return i18n.LANG_CHANGE_RESTART_REQUIRED
}

// IsReady は初期化済み扱いを返す。
func (s *stubI18n) IsReady() bool { return true }

// T はキー変換を返す。
func (s *stubI18n) T(key string) string { return s.messages[key] }

// TWithLang はキー変換を返す。
func (s *stubI18n) TWithLang(lang i18n.LangCode, key string) string { return s.messages[key] }

type stubVerboseSink struct {
	lines []string
}

// WriteLine は行を追加する。
func (s *stubVerboseSink) WriteLine(text string) error {
	s.lines = append(s.lines, text)
	return nil
}

// Close は未使用のためnilを返す。
func (s *stubVerboseSink) Close() error { return nil }

// TestLoggerLevelSuppression はレベル抑制を確認する。
func TestLoggerLevelSuppression(t *testing.T) {
	logger := NewLogger(nil)
	logger.SetLevel(logging.LOG_LEVEL_INFO)
	logger.MessageBuffer().Clear()

	logger.Debug("dbg")
	logger.Info("info")

	lines := logger.MessageBuffer().Lines()
	if len(lines) != 1 || lines[0] != "info" {
		t.Errorf("lines: got=%v", lines)
	}
}

// TestLoggerVerbose は冗長ログの出力条件を確認する。
func TestLoggerVerbose(t *testing.T) {
	logger := NewLogger(nil)
	logger.SetLevel(logging.LOG_LEVEL_VERBOSE)
	logger.EnableVerbose(logging.VERBOSE_INDEX_IK)
	sink := &stubVerboseSink{}
	logger.AttachVerboseSink(logging.VERBOSE_INDEX_IK, sink)

	logger.Verbose(logging.VERBOSE_INDEX_IK, "ik-%d", 1)
	logger.Verbose(logging.VERBOSE_INDEX_MOTION, "motion")

	if len(sink.lines) != 1 || sink.lines[0] != "ik-1" {
		t.Errorf("verbose lines: got=%v", sink.lines)
	}
}

// TestLoggerTranslation はINFO以上の翻訳を確認する。
func TestLoggerTranslation(t *testing.T) {
	translator := &stubI18n{messages: map[string]string{"key": "保存開始: %s"}}
	logger := NewLogger(translator)
	logger.SetLevel(logging.LOG_LEVEL_INFO)
	logger.MessageBuffer().Clear()

	logger.Info("key", "Vmd")

	lines := logger.MessageBuffer().Lines()
	if len(lines) != 1 || lines[0] != "保存開始: Vmd" {
		t.Errorf("translated line: got=%v", lines)
	}
}
