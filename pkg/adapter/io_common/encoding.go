// 指示: miu200521358
package io_common

import (
	"bytes"
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
)

// DecodePmxText はPMXの文字列をデコードする。
func DecodePmxText(primary encoding.Encoding, raw []byte) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	fallbacks := []encoding.Encoding{
		primary,
		japanese.ShiftJIS,
		unicode.UTF8,
		unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM),
	}
	var lastErr error
	for _, enc := range fallbacks {
		if enc == nil {
			continue
		}
		decoded, err := decodeText(enc, raw)
		if err == nil {
			return decoded, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("エンコードの判定に失敗しました")
	}
	return "", lastErr
}

// DecodeShiftJISFixed はShift-JIS固定長文字列をデコードする。
func DecodeShiftJISFixed(raw []byte) (string, error) {
	decoded, err := japanese.ShiftJIS.NewDecoder().Bytes(raw)
	if err != nil {
		return "", err
	}
	decoded = bytes.TrimRight(decoded, "\xfd")
	decoded = bytes.TrimRight(decoded, "\x00")
	decoded = bytes.ReplaceAll(decoded, []byte("\x00"), []byte(" "))
	return string(decoded), nil
}

// EncodeShiftJISFixed はShift-JISで固定長へエンコードする。
func EncodeShiftJISFixed(text string, size int) ([]byte, error) {
	if size <= 0 {
		return []byte{}, nil
	}
	encoder := japanese.ShiftJIS.NewEncoder()
	buf := make([]byte, 0, size)
	for _, r := range text {
		encoded, err := encoder.Bytes([]byte(string(r)))
		if err != nil {
			return nil, err
		}
		if len(buf)+len(encoded) > size {
			break
		}
		buf = append(buf, encoded...)
	}
	if len(buf) < size {
		padding := make([]byte, size-len(buf))
		buf = append(buf, padding...)
	}
	return buf, nil
}

// decodeText は指定エンコーディングでデコードする。
func decodeText(enc encoding.Encoding, raw []byte) (string, error) {
	decoded, err := enc.NewDecoder().Bytes(raw)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
