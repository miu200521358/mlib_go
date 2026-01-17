// 指示: miu200521358
package io_common

import (
	"encoding/binary"
	"io"
	"math"
)

// BinaryWriter はバイナリ書き込みを補助する。
type BinaryWriter struct {
	writer io.Writer
}

// NewBinaryWriter はバイナリライターを生成する。
func NewBinaryWriter(w io.Writer) *BinaryWriter {
	return &BinaryWriter{writer: w}
}

// WriteBytes はバイト列を書き込む。
func (b *BinaryWriter) WriteBytes(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	_, err := b.writer.Write(data)
	return err
}

// WriteUint8 はuint8を書き込む。
func (b *BinaryWriter) WriteUint8(value uint8) error {
	return b.WriteBytes([]byte{value})
}

// WriteInt8 はint8を書き込む。
func (b *BinaryWriter) WriteInt8(value int8) error {
	return b.WriteBytes([]byte{byte(value)})
}

// WriteUint16 はuint16を書き込む。
func (b *BinaryWriter) WriteUint16(value uint16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, value)
	return b.WriteBytes(buf)
}

// WriteInt16 はint16を書き込む。
func (b *BinaryWriter) WriteInt16(value int16) error {
	return b.WriteUint16(uint16(value))
}

// WriteUint32 はuint32を書き込む。
func (b *BinaryWriter) WriteUint32(value uint32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	return b.WriteBytes(buf)
}

// WriteInt32 はint32を書き込む。
func (b *BinaryWriter) WriteInt32(value int32) error {
	return b.WriteUint32(uint32(value))
}

// WriteFloat32 はfloat32を書き込む。
func (b *BinaryWriter) WriteFloat32(value float64, defaultVal float64, positiveOnly bool) error {
	sanitized := sanitizeFloat(value, defaultVal, positiveOnly)
	bits := math.Float32bits(float32(sanitized))
	return b.WriteUint32(bits)
}

// sanitizeFloat はNaNや無限大を既定値に置換し、必要なら負値を抑制する。
func sanitizeFloat(value float64, defaultVal float64, positiveOnly bool) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return defaultVal
	}
	if positiveOnly && value < 0 {
		return 0
	}
	return value
}
