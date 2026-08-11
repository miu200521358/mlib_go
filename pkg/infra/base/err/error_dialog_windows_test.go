//go:build windows
// +build windows

// 指示: miu200521358
package err

import (
	"errors"
	"strings"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
)

// stackedError は walk.Error と同じ Stack() []byte / Inner() / Message() を持つ疑似エラー。
type stackedError struct {
	inner   error
	message string
	stack   []byte
}

func (e *stackedError) Error() string { return e.Message() + "\n\nStack:\n" + string(e.stack) }

// Message は walk.Error と同じく、自分の文が無ければ内側へ委譲する。
func (e *stackedError) Message() string {
	if e.message != "" {
		return e.message
	}
	if e.inner != nil {
		return e.inner.Error()
	}
	return ""
}
func (e *stackedError) Stack() []byte { return e.stack }
func (e *stackedError) Inner() error  { return e.inner }

// TestBuildErrorTextIncludesCauseStack は原因連鎖のスタックが本文へ出ることを確認する。
func TestBuildErrorTextIncludesCauseStack(t *testing.T) {
	cause := &stackedError{message: "EM_SETCUEBANNER failed", stack: []byte("goroutine 1 [running]:\nwalk.setCueBanner()\n")}
	err := merr.NewCommonError("99901", merr.ErrorKindInternal, "コントローラーウィンドウの初期化に失敗しました", cause)

	text := buildErrorText(nil, err)

	if !strings.Contains(text, "EM_SETCUEBANNER failed") {
		t.Fatalf("原因メッセージがありません: %s", text)
	}
	if !strings.Contains(text, "Stack:\ngoroutine 1 [running]:") {
		t.Fatalf("原因のスタックがありません: %s", text)
	}
	if !strings.Contains(text, "------------") {
		t.Fatalf("技術情報の区切りがありません: %s", text)
	}
	// walk 相当の Error() はスタックを含むため、メッセージ行の組み立てに
	// Error() を使うとスタックが重複する。1 回だけ出ることを固定する。
	if strings.Count(text, "goroutine 1 [running]:") != 1 {
		t.Fatalf("スタックが重複しています: %s", text)
	}
}

// TestBuildErrorTextFallsBackToDisplayStack はスタックを持たない連鎖でも必ずスタックが出ることを確認する。
func TestBuildErrorTextFallsBackToDisplayStack(t *testing.T) {
	err := merr.NewCommonError("99902", merr.ErrorKindInternal, "初期化に失敗しました", errors.New("plain failure"))

	text := buildErrorText(nil, err)

	if !strings.Contains(text, "plain failure") {
		t.Fatalf("原因メッセージがありません: %s", text)
	}
	if !strings.Contains(text, "\nStack:\n") {
		t.Fatalf("フォールバックのスタックがありません: %s", text)
	}
	if !strings.Contains(text, "goroutine") {
		t.Fatalf("表示時点のスタック本文がありません: %s", text)
	}
}

// TestBuildTechnicalTextFollowsInnerChain は Unwrap を持たない Inner() 連鎖をたどることを確認する。
func TestBuildTechnicalTextFollowsInnerChain(t *testing.T) {
	leaf := errors.New("root cause")
	wrapped := &stackedError{inner: leaf, message: "outer failure", stack: []byte("goroutine 7 [running]:\n")}

	text := buildTechnicalText(wrapped)

	if !strings.Contains(text, "outer failure") {
		t.Fatalf("外側メッセージがありません: %s", text)
	}
	if !strings.Contains(text, "root cause") {
		t.Fatalf("Inner 連鎖の根本原因がありません: %s", text)
	}
	if !strings.Contains(text, "goroutine 7 [running]:") {
		t.Fatalf("Inner 連鎖のスタックがありません: %s", text)
	}
}

// TestBuildTechnicalTextDeduplicatesConsecutiveLines は同文の連鎖段が重複表示されないことを確認する。
func TestBuildTechnicalTextDeduplicatesConsecutiveLines(t *testing.T) {
	// message 未設定の walk 相当エラーは Message() が内側へ委譲するため、
	// 隣接する 2 段が同じ文字列になる。
	leaf := errors.New("same message")
	wrapped := &stackedError{inner: leaf, stack: []byte("goroutine 9 [running]:\n")}

	text := buildTechnicalText(wrapped)

	if strings.Count(text, "same message") != 1 {
		t.Fatalf("連鎖の同文が重複しています: %s", text)
	}
}
