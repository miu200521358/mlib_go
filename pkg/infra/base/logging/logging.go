// 指示: miu200521358
package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

const maxConsoleLines = 10000

// ILogSink はログ出力先I/F。
type ILogSink interface {
	Write(p []byte) (int, error)
}

// Logger はinfraログ実装。
type Logger struct {
	mu             sync.Mutex
	level          logging.LogLevel
	translator     i18n.II18n
	buffer         *messageBuffer
	consoleSink    *logSink
	verboseEnabled map[logging.VerboseIndex]bool
	verboseSinks   map[logging.VerboseIndex]logging.IVerboseSink
}

var defaultLogger = NewLogger(nil)

// NewLogger はロガーを生成する。
func NewLogger(translator i18n.II18n) *Logger {
	buf := &messageBuffer{}
	ls := &logSink{sink: os.Stderr, buffer: buf}
	log.SetFlags(0)
	log.SetOutput(ls)
	return &Logger{
		level:          logging.LOG_LEVEL_INFO,
		translator:     translator,
		buffer:         buf,
		consoleSink:    ls,
		verboseEnabled: map[logging.VerboseIndex]bool{},
		verboseSinks:   map[logging.VerboseIndex]logging.IVerboseSink{},
	}
}

// DefaultLogger は既定ロガーを返す。
func DefaultLogger() *Logger {
	return defaultLogger
}

// SetDefaultLogger は既定ロガーを差し替える。
func SetDefaultLogger(logger *Logger) {
	if logger == nil {
		return
	}
	defaultLogger = logger
}

// SetConsoleSink はコンソール出力先を設定する。
func SetConsoleSink(writer ILogSink) {
	defaultLogger.SetConsoleSink(writer)
}

// ConsoleText はメッセージ欄の全文を返す。
func ConsoleText() string {
	return defaultLogger.ConsoleText()
}

// Level は現在のログレベルを返す。
func (l *Logger) Level() logging.LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetLevel はログレベルを設定する。
func (l *Logger) SetLevel(level logging.LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// EnableVerbose は冗長ログを有効化する。
func (l *Logger) EnableVerbose(idx logging.VerboseIndex) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.verboseEnabled[idx] = true
}

// DisableVerbose は冗長ログを無効化する。
func (l *Logger) DisableVerbose(idx logging.VerboseIndex) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.verboseEnabled[idx] = false
}

// IsVerboseEnabled は冗長ログの有効可否を返す。
func (l *Logger) IsVerboseEnabled(idx logging.VerboseIndex) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.verboseEnabled[idx]
}

// MessageBuffer はメッセージ欄バッファを返す。
func (l *Logger) MessageBuffer() logging.IMessageBuffer {
	return l.buffer
}

// AttachVerboseSink は冗長ログの出力先を設定する。
func (l *Logger) AttachVerboseSink(idx logging.VerboseIndex, sink logging.IVerboseSink) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.verboseSinks[idx] = sink
}

// Verbose は冗長ログを出力する。
func (l *Logger) Verbose(idx logging.VerboseIndex, msg string, params ...any) {
	l.mu.Lock()
	level := l.level
	enabled := l.verboseEnabled[idx]
	sink := l.verboseSinks[idx]
	l.mu.Unlock()

	if level != logging.LOG_LEVEL_VERBOSE || !enabled || sink == nil {
		return
	}
	text := fmt.Sprintf(msg, params...)
	sink.WriteLine(text)
}

// Debug はデバッグログを出力する。
func (l *Logger) Debug(msg string, params ...any) {
	l.logWithLevel(logging.LOG_LEVEL_DEBUG, msg, params...)
}

// Info は情報ログを出力する。
func (l *Logger) Info(msg string, params ...any) {
	l.logWithLevel(logging.LOG_LEVEL_INFO, msg, params...)
}

// Warn は警告ログを出力する。
func (l *Logger) Warn(msg string, params ...any) {
	l.logWithLevel(logging.LOG_LEVEL_WARN, msg, params...)
}

// Error はエラーログを出力する。
func (l *Logger) Error(msg string, params ...any) {
	l.logWithLevel(logging.LOG_LEVEL_ERROR, msg, params...)
}

// Fatal は致命ログを出力する。
func (l *Logger) Fatal(msg string, params ...any) {
	l.logWithLevel(logging.LOG_LEVEL_FATAL, msg, params...)
}

// Line は区切り線を出力する。
func (l *Logger) Line() {
	if !l.shouldOutput(logging.LOG_LEVEL_INFO) {
		return
	}
	log.Printf("---------------------------------")
}

// InfoLine は区切り線付き情報ログを出力する。
func (l *Logger) InfoLine(msg string, params ...any) {
	l.Line()
	l.Info(msg, params...)
}

// InfoStamp はタイムスタンプ付き情報ログを出力する。
func (l *Logger) InfoStamp(msg string, params ...any) {
	if !l.shouldOutput(logging.LOG_LEVEL_INFO) {
		return
	}
	stamp := time.Now().Format("15:04:05.999999999")
	formatted := l.formatMessage(logging.LOG_LEVEL_INFO, msg)
	log.Printf("[%s]"+formatted, append([]any{stamp}, params...)...)
}

// InfoTitle はタイトル付き情報ログを出力する。
func (l *Logger) InfoTitle(title, msg string, params ...any) {
	if !l.shouldOutput(logging.LOG_LEVEL_INFO) {
		return
	}
	log.Printf("■■■■■ %s ■■■■■", title)
	log.Printf(l.formatMessage(logging.LOG_LEVEL_INFO, msg), params...)
}

// InfoLineTitle は区切り線とタイトル付き情報ログを出力する。
func (l *Logger) InfoLineTitle(title, msg string, params ...any) {
	l.Line()
	l.InfoTitle(title, msg, params...)
}

// WarnTitle はタイトル付き警告ログを出力する。
func (l *Logger) WarnTitle(title, msg string, params ...any) {
	if !l.shouldOutput(logging.LOG_LEVEL_WARN) {
		return
	}
	log.Printf("~~~~~~~~~~ %s ~~~~~~~~~~", title)
	log.Printf(l.formatMessage(logging.LOG_LEVEL_WARN, msg), params...)
}

// ErrorTitle はタイトル付きエラーログを出力する。
func (l *Logger) ErrorTitle(title string, err error, msg string, params ...any) {
	if !l.shouldOutput(logging.LOG_LEVEL_ERROR) {
		return
	}
	log.Printf("********** %s **********", title)
	if err != nil {
		log.Printf("Error Message: %s", err.Error())
		log.Printf("Stack Trace:\n%s", dumpAllGoroutines())
	}
	if msg != "" {
		log.Printf(l.formatMessage(logging.LOG_LEVEL_ERROR, msg), params...)
	}
}

// FatalTitle はタイトル付き致命ログを出力する。
func (l *Logger) FatalTitle(title string, err error, msg string, params ...any) {
	if !l.shouldOutput(logging.LOG_LEVEL_FATAL) {
		return
	}
	log.Printf("!!!!!!!!!! %s !!!!!!!!!!", title)
	if err != nil {
		log.Printf("Error Message: %s", err.Error())
		log.Printf("Stack Trace:\n%s", dumpAllGoroutines())
	}
	if msg != "" {
		log.Printf(l.formatMessage(logging.LOG_LEVEL_FATAL, msg), params...)
	}
}

// Memory はメモリ使用量を出力する。
func (l *Logger) Memory(prefix string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mem := bToMb(m.Alloc)
	if mem == prevMem {
		return
	}
	log.Printf("[%s] Alloc = %v -> %v MiB, HeapAlloc = %v MiB, HeapSys = %v MiB, HeapIdle = %v MiB, "+
		"HeapInuse = %v MiB, HeapReleased = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
		prefix, prevMem, mem, bToMb(m.HeapAlloc), bToMb(m.HeapSys), bToMb(m.HeapIdle),
		bToMb(m.HeapInuse), bToMb(m.HeapReleased), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	prevMem = mem
}

// SetConsoleSink はコンソール出力先を設定する。
func (l *Logger) SetConsoleSink(writer ILogSink) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.consoleSink.sink = writer
	log.SetOutput(l.consoleSink)
}

// ConsoleText はメッセージ欄の全文を返す。
func (l *Logger) ConsoleText() string {
	return l.buffer.Text()
}

// logWithLevel はレベル判定付きで出力する。
func (l *Logger) logWithLevel(level logging.LogLevel, msg string, params ...any) {
	if !l.shouldOutput(level) {
		return
	}
	log.Printf(l.formatMessage(level, msg), params...)
}

// shouldOutput は出力可否を判定する。
func (l *Logger) shouldOutput(level logging.LogLevel) bool {
	l.mu.Lock()
	current := l.level
	l.mu.Unlock()
	if current == logging.LOG_LEVEL_VERBOSE {
		current = logging.LOG_LEVEL_DEBUG
	}
	return current <= level
}

// formatMessage はレベルに応じた翻訳済みメッセージを返す。
func (l *Logger) formatMessage(level logging.LogLevel, msg string) string {
	if level < logging.LOG_LEVEL_INFO {
		return msg
	}
	l.mu.Lock()
	translator := l.translator
	l.mu.Unlock()
	if translator == nil {
		return msg
	}
	return translator.T(msg)
}

// messageBuffer はログバッファを保持する。
type messageBuffer struct {
	mu    sync.Mutex
	lines []string
}

// Text は全文を返す。
func (b *messageBuffer) Text() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return strings.Join(b.lines, "\r\n")
}

// Lines は行配列を返す。
func (b *messageBuffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]string, len(b.lines))
	copy(out, b.lines)
	return out
}

// Clear はバッファをクリアする。
func (b *messageBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = nil
}

// appendText はログテキストを追記する。
func (b *messageBuffer) appendText(text string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	parts := strings.Split(normalized, "\n")
	for i, part := range parts {
		if i == len(parts)-1 && part == "" {
			continue
		}
		b.lines = append(b.lines, part)
	}
	if len(b.lines) > maxConsoleLines {
		b.lines = b.lines[len(b.lines)-maxConsoleLines:]
	}
}

// logSink はログ出力とバッファ更新を担当する。
type logSink struct {
	mu     sync.Mutex
	sink   ILogSink
	buffer *messageBuffer
}

// Write はログを書き込む。
func (s *logSink) Write(p []byte) (int, error) {
	s.textWrite(p)
	return len(p), nil
}

// textWrite は出力とバッファ更新を行う。
func (s *logSink) textWrite(p []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sink != nil {
		_, _ = s.sink.Write(p)
	}
	if s.buffer != nil {
		s.buffer.appendText(string(p))
	}
}

// dumpAllGoroutines は全ゴルーチンのスタックを取得する。
func dumpAllGoroutines() string {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	return string(bytes.ReplaceAll(buf[:n], []byte("\n"), []byte("\r\n")))
}

var prevMem uint64

// bToMb はバイトをMiBへ変換する。
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
