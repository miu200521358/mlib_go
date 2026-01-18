// 指示: miu200521358
package x

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	xHeaderSize      = 16
	xCompressedStart = 20
)

// decompressMSZip はMSZIP圧縮バイナリを解凍する。
func decompressMSZip(data []byte) ([]byte, error) {
	if len(data) < xCompressedStart {
		return nil, fmt.Errorf("X圧縮データが不足しています")
	}
	finalSize := binary.LittleEndian.Uint32(data[xHeaderSize:xCompressedStart])
	if finalSize < xHeaderSize {
		return nil, fmt.Errorf("X圧縮データの最終サイズが不正です")
	}
	out := make([]byte, int(finalSize))
	copy(out[:xHeaderSize], data[:xHeaderSize])

	outPos := xHeaderSize
	inPos := xCompressedStart
	for inPos < len(data) && outPos < int(finalSize) {
		if inPos+4 > len(data) {
			return nil, fmt.Errorf("X圧縮ブロックのサイズが不足しています")
		}
		uncompressedSize := binary.LittleEndian.Uint16(data[inPos:])
		compressedSize := binary.LittleEndian.Uint16(data[inPos+2:])
		inPos += 4
		if inPos+int(compressedSize) > len(data) {
			return nil, fmt.Errorf("X圧縮ブロックが不足しています")
		}
		block := data[inPos : inPos+int(compressedSize)]
		inPos += int(compressedSize)
		payload := block
		if len(block) >= 2 && block[0] == 'C' && block[1] == 'K' {
			payload = block[2:]
		}

		dict := buildMSZipDict(out[:outPos])
		var reader io.ReadCloser
		if len(dict) > 0 {
			reader = flate.NewReaderDict(bytes.NewReader(payload), dict)
		} else {
			reader = flate.NewReader(bytes.NewReader(payload))
		}
		decoded, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			return nil, fmt.Errorf("X圧縮ブロックの解凍に失敗しました")
		}
		if int(uncompressedSize) != len(decoded) {
			return nil, fmt.Errorf("X圧縮ブロックの伸長サイズが不正です")
		}
		if outPos+len(decoded) > int(finalSize) {
			return nil, fmt.Errorf("X圧縮ブロックの展開先が不足しています")
		}
		copy(out[outPos:], decoded)
		outPos += len(decoded)
	}
	return out[:outPos], nil
}

// buildMSZipDict は直近32KBの辞書を返す。
func buildMSZipDict(data []byte) []byte {
	const dictSize = 32 * 1024
	if len(data) <= xHeaderSize {
		return nil
	}
	payload := data[xHeaderSize:]
	if len(payload) <= dictSize {
		return payload
	}
	return payload[len(payload)-dictSize:]
}
