package io_common

import (
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
)

func TestDecodePmxTextFallbackFromUTF8ToShiftJIS(t *testing.T) {
	want := "右ひざ"
	raw, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(want))
	if err != nil {
		t.Fatalf("shift-jis encode failed: %v", err)
	}

	got, err := DecodePmxText(unicode.UTF8, raw)
	if err != nil {
		t.Fatalf("DecodePmxText returned error: %v", err)
	}
	if got != want {
		t.Fatalf("DecodePmxText mismatch: got=%q want=%q", got, want)
	}
}

func TestDecodePmxTextUTF8(t *testing.T) {
	want := "左スカート前"
	raw := []byte(want)

	got, err := DecodePmxText(unicode.UTF8, raw)
	if err != nil {
		t.Fatalf("DecodePmxText returned error: %v", err)
	}
	if got != want {
		t.Fatalf("DecodePmxText mismatch: got=%q want=%q", got, want)
	}
}
