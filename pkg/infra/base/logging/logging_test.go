// 指示: miu200521358
package logging

import (
	"bytes"
	"errors"
	"runtime"
	"strings"
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
func (s *stubI18n) SetLang(lang i18n.LangCode) (i18n.LangChangeAction, error) {
	return i18n.LANG_CHANGE_RESTART_REQUIRED, nil
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
func (s *stubVerboseSink) WriteLine(text string) {
	s.lines = append(s.lines, text)
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

type stubConsoleSink struct {
	buf bytes.Buffer
}

// Write はログを記録する。
func (s *stubConsoleSink) Write(p []byte) (int, error) {
	return s.buf.Write(p)
}

// TestDefaultLoggerAndConsoleSink は既定ロガーの切替と出力を確認する。
func TestDefaultLoggerAndConsoleSink(t *testing.T) {
	prev := DefaultLogger()
	logger := NewLogger(nil)
	SetDefaultLogger(logger)
	t.Cleanup(func() { SetDefaultLogger(prev) })

	SetDefaultLogger(nil)
	if DefaultLogger() != logger {
		t.Errorf("SetDefaultLogger nil should keep current")
	}

	sink := &stubConsoleSink{}
	SetConsoleSink(sink)
	logger.SetLevel(logging.LOG_LEVEL_INFO)
	logger.Info("hello")

	if !strings.Contains(ConsoleText(), "hello") {
		t.Errorf("ConsoleText missing message: %v", ConsoleText())
	}
	if sink.buf.Len() == 0 {
		t.Errorf("console sink did not receive output")
	}
}

// TestLoggerVerboseDisabled は冗長ログの無効時を確認する。
func TestLoggerVerboseDisabled(t *testing.T) {
	logger := NewLogger(nil)
	logger.SetLevel(logging.LOG_LEVEL_INFO)
	sink := &stubVerboseSink{}
	logger.AttachVerboseSink(logging.VERBOSE_INDEX_VIEWER, sink)
	logger.EnableVerbose(logging.VERBOSE_INDEX_VIEWER)
	logger.DisableVerbose(logging.VERBOSE_INDEX_VIEWER)

	logger.Verbose(logging.VERBOSE_INDEX_VIEWER, "skip")
	if len(sink.lines) != 0 {
		t.Errorf("Verbose should be skipped: got=%v", sink.lines)
	}
	if logger.IsVerboseEnabled(logging.VERBOSE_INDEX_VIEWER) {
		t.Errorf("Verbose should be disabled")
	}
}

// TestLoggerTitlesAndErrors はタイトル付き出力を確認する。
func TestLoggerTitlesAndErrors(t *testing.T) {
	logger := NewLogger(nil)
	logger.SetLevel(logging.LOG_LEVEL_DEBUG)
	logger.MessageBuffer().Clear()

	logger.Line()
	logger.InfoLine("info-line")
	logger.InfoStamp("info-stamp")
	logger.InfoTitle("Title", "info-title")
	logger.InfoLineTitle("Title2", "info-line-title")
	logger.WarnTitle("Warn", "warn-msg")
	logger.ErrorTitle("Err", errors.New("boom"), "err-msg")
	logger.FatalTitle("Fatal", errors.New("boom"), "fatal-msg")

	text := logger.ConsoleText()
	if !strings.Contains(text, "info-line") || !strings.Contains(text, "■■■■■") ||
		!strings.Contains(text, "~~~~~~~~~~") || !strings.Contains(text, "**********") || !strings.Contains(text, "!!!!!!!!!!") {
		t.Errorf("ConsoleText missing title outputs: %v", text)
	}
}

// TestLoggerDebugFormat はDEBUG時の翻訳抑制を確認する。
func TestLoggerDebugFormat(t *testing.T) {
	translator := &stubI18n{messages: map[string]string{"key": "translated"}}
	logger := NewLogger(translator)
	logger.SetLevel(logging.LOG_LEVEL_VERBOSE)
	logger.MessageBuffer().Clear()

	logger.Debug("key")
	lines := logger.MessageBuffer().Lines()
	if len(lines) == 0 || lines[0] != "key" {
		t.Errorf("Debug should not translate: got=%v", lines)
	}
}

// TestMessageBufferAppendAndTrim はバッファの正規化と上限を確認する。
func TestMessageBufferAppendAndTrim(t *testing.T) {
	var buf messageBuffer
	buf.appendText("a\r\nb\r\n")
	if got := buf.Lines(); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("appendText normalize: got=%v", got)
	}

	var builder strings.Builder
	for i := 0; i < maxConsoleLines+1; i++ {
		builder.WriteString("x\n")
	}
	buf.appendText(builder.String())
	if len(buf.Lines()) != maxConsoleLines {
		t.Errorf("appendText trim: got=%v", len(buf.Lines()))
	}
	buf.Clear()
	if len(buf.Lines()) != 0 {
		t.Errorf("Clear failed: got=%v", buf.Lines())
	}
}

// TestLogSinkWrite は書き込みとバッファ更新を確認する。
func TestLogSinkWrite(t *testing.T) {
	sink := &stubConsoleSink{}
	buffer := &messageBuffer{}
	ls := &logSink{sink: sink, buffer: buffer}
	if _, err := ls.Write([]byte("line1\n")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if !strings.Contains(sink.buf.String(), "line1") {
		t.Errorf("sink output missing: %v", sink.buf.String())
	}
	if got := buffer.Lines(); len(got) == 0 || got[0] != "line1" {
		t.Errorf("buffer lines: got=%v", got)
	}

	ls = &logSink{sink: nil, buffer: buffer}
	ls.textWrite([]byte("line2\n"))
	if got := buffer.Lines(); len(got) < 2 || got[len(got)-1] != "line2" {
		t.Errorf("textWrite buffer: got=%v", got)
	}
}

// TestLoggerMemory はメモリ出力の分岐を確認する。
func TestLoggerMemory(t *testing.T) {
	logger := NewLogger(nil)
	logger.SetLevel(logging.LOG_LEVEL_INFO)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	expected := bToMb(m.Alloc)

	prevMem = expected
	logger.Memory("mem")

	prevMem = expected ^ 1
	logger.Memory("mem")
}

// TestLoggerSuppressionBranches は抑制分岐と各メソッドを確認する。
func TestLoggerSuppressionBranches(t *testing.T) {
	logger := NewLogger(nil)
	logger.SetLevel(logging.LOG_LEVEL_ERROR)
	if logger.Level() != logging.LOG_LEVEL_ERROR {
		t.Errorf("Level mismatch")
	}

	logger.Line()
	logger.InfoStamp("skip")
	logger.InfoTitle("skip", "skip")
	logger.WarnTitle("skip", "skip")
	logger.ErrorTitle("skip", nil, "")
	logger.FatalTitle("skip", nil, "")

	logger.Warn("warn")
	logger.Error("err")
	logger.Fatal("fatal")

	logger.SetLevel(logging.LogLevel(100))
	logger.ErrorTitle("skip", nil, "")
	logger.FatalTitle("skip", nil, "")
}
