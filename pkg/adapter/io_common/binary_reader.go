// 指示: miu200521358
package io_common

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"gonum.org/v1/gonum/spatial/r3"
)

// BinaryReader はバイナリ読み込みを補助する。
type BinaryReader struct {
	reader *bufio.Reader
}

// NewBinaryReader はバイナリリーダーを生成する。
func NewBinaryReader(r io.Reader) *BinaryReader {
	return &BinaryReader{reader: bufio.NewReader(r)}
}

// ReadBytes は指定サイズのバイト列を読み込む。
func (b *BinaryReader) ReadBytes(size int) ([]byte, error) {
	if size <= 0 {
		return []byte{}, nil
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(b.reader, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// DiscardAll は残りのデータを破棄する。
func (b *BinaryReader) DiscardAll() error {
	if b == nil || b.reader == nil {
		return nil
	}
	_, err := io.Copy(io.Discard, b.reader)
	return err
}

// ReadUint8 はuint8を読み込む。
func (b *BinaryReader) ReadUint8() (uint8, error) {
	buf, err := b.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadInt8 はint8を読み込む。
func (b *BinaryReader) ReadInt8() (int8, error) {
	buf, err := b.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return int8(buf[0]), nil
}

// ReadUint16 はuint16を読み込む。
func (b *BinaryReader) ReadUint16() (uint16, error) {
	buf, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// ReadInt16 はint16を読み込む。
func (b *BinaryReader) ReadInt16() (int16, error) {
	buf, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(buf)), nil
}

// ReadUint32 はuint32を読み込む。
func (b *BinaryReader) ReadUint32() (uint32, error) {
	buf, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// ReadInt32 はint32を読み込む。
func (b *BinaryReader) ReadInt32() (int32, error) {
	buf, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(buf)), nil
}

// ReadFloat32 はfloat32を読み込んでfloat64で返す。
func (b *BinaryReader) ReadFloat32() (float64, error) {
	value, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return float64(math.Float32frombits(value)), nil
}

// ReadFloat32s はfloat32配列を読み込む。
func (b *BinaryReader) ReadFloat32s(values []float64) ([]float64, error) {
	for i := range values {
		v, err := b.ReadFloat32()
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

// ReadVec2 はVec2を読み込む。
func (b *BinaryReader) ReadVec2() (mmath.Vec2, error) {
	x, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec2{}, err
	}
	y, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec2{}, err
	}
	return mmath.Vec2{X: x, Y: y}, nil
}

// ReadVec3 はVec3を読み込む。
func (b *BinaryReader) ReadVec3() (mmath.Vec3, error) {
	x, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec3{}, err
	}
	y, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec3{}, err
	}
	z, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec3{}, err
	}
	return mmath.Vec3{Vec: r3.Vec{X: x, Y: y, Z: z}}, nil
}

// ReadVec4 はVec4を読み込む。
func (b *BinaryReader) ReadVec4() (mmath.Vec4, error) {
	x, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	y, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	z, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	w, err := b.ReadFloat32()
	if err != nil {
		return mmath.Vec4{}, err
	}
	return mmath.Vec4{X: x, Y: y, Z: z, W: w}, nil
}
