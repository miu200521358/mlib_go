package repository

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"

	"github.com/miu200521358/mlib_go/pkg/domain/core"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IRepository interface {
	LoadName(path string) (string, error)
	Load(path string) (core.IHashModel, error)
	Save(overridePath string, data core.IHashModel, includeSystem bool) error
}

type baseRepository[T core.IHashModel] struct {
	file     fs.File
	reader   *bufio.Reader
	encoding encoding.Encoding
	readText func() string
	newFunc  func(path string) T
}

func (rep *baseRepository[T]) open(path string) error {

	file, err := mutils.Open(path)
	if err != nil {
		return err
	}
	rep.file = file
	rep.reader = bufio.NewReader(rep.file)

	return nil
}

func (rep *baseRepository[T]) close() {
	defer rep.file.Close()
}

func (rep *baseRepository[T]) LoadName(path string) (string, error) {
	panic("not implemented")
}

func (rep *baseRepository[T]) Load(path string) (T, error) {
	panic("not implemented")
}

func (rep *baseRepository[T]) defineEncoding(encoding encoding.Encoding) {
	rep.encoding = encoding
	rep.readText = rep.defineReadText(encoding)
}

func (rep *baseRepository[T]) defineReadText(encoding encoding.Encoding) func() string {
	return func() string {
		size, err := rep.unpackInt()
		if err != nil {
			return ""
		}
		fbytes, err := rep.unpackBytes(int(size))
		if err != nil {
			return ""
		}
		return rep.decodeText(encoding, fbytes)
	}
}

func (rep *baseRepository[T]) decodeText(mainEncoding encoding.Encoding, fbytes []byte) string {
	// 基本のエンコーディングを第一候補でデコードして、ダメなら順次テスト
	for _, targetEncoding := range []encoding.Encoding{
		mainEncoding,
		japanese.ShiftJIS,
		unicode.UTF8,
		unicode.UTF16(unicode.LittleEndian, unicode.UseBOM),
	} {
		var decodedText string
		var err error
		if targetEncoding == japanese.ShiftJIS {
			// shift-jisは一旦cp932に変換してもう一度戻したので返す
			decodedText, err = rep.decodeShiftJIS(fbytes)
			if err != nil {
				continue
			}
		} else {
			// 変換できなかった文字は「?」に変換する
			decodedText, err = rep.decodeBytes(fbytes, targetEncoding)
			if err != nil {
				continue
			}
		}
		return decodedText
	}
	return ""
}

func (rep *baseRepository[T]) decodeShiftJIS(fbytes []byte) (string, error) {
	decodedText, err := japanese.ShiftJIS.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

func (rep *baseRepository[T]) decodeBytes(fbytes []byte, encoding encoding.Encoding) (string, error) {
	decodedText, err := encoding.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

// バイナリデータから bytes を読み出す
func (rep *baseRepository[T]) unpackBytes(size int) ([]byte, error) {
	chunk, err := rep.unpack(size)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// バイナリデータから byte を読み出す
func (rep *baseRepository[T]) unpackByte() (byte, error) {
	chunk, err := rep.unpack(1)
	if err != nil {
		return 0, err
	}

	return chunk[0], nil
}

// バイナリデータから sbyte を読み出す
func (rep *baseRepository[T]) unpackSByte() (int8, error) {
	chunk, err := rep.unpack(1)
	if err != nil {
		return 0, err
	}

	return int8(chunk[0]), nil
}

// バイナリデータから sbyte を読み出す
func (rep *baseRepository[T]) unpackSBytes(count int) ([]int8, error) {
	values := make([]int8, count)
	chunk, err := rep.unpack(count)
	if err != nil {
		return values, err
	}

	for i := 0; i < count; i++ {
		values[i] = int8(chunk[i])
	}

	return values, nil
}

// バイナリデータから int16 を読み出す
func (rep *baseRepository[T]) unpackShort() (int16, error) {
	chunk, err := rep.unpack(2)
	if err != nil {
		return 0, err
	}

	return int16(binary.LittleEndian.Uint16(chunk)), nil
}

// バイナリデータから int16 を読み出す
func (rep *baseRepository[T]) unpackShorts(count int) ([]int16, error) {
	values := make([]int16, count)
	chunk, err := rep.unpack(2 * count)
	if err != nil {
		return values, err
	}

	for i := 0; i < count; i++ {
		values[i] = int16(binary.LittleEndian.Uint16(chunk[i*2 : (i+1)*2]))
	}

	return values, nil
}

// バイナリデータから uint16 を読み出す
func (rep *baseRepository[T]) unpackUShort() (uint16, error) {
	chunk, err := rep.unpack(2)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(chunk), nil
}

// バイナリデータから uint16 を読み出す
func (rep *baseRepository[T]) unpackUShorts(count int) ([]uint16, error) {
	values := make([]uint16, count)
	chunk, err := rep.unpack(2 * count)
	if err != nil {
		return values, err
	}

	for i := 0; i < count; i++ {
		values[i] = binary.LittleEndian.Uint16(chunk[i*2 : (i+1)*2])
	}

	return values, nil
}

// バイナリデータから uint を読み出す
func (rep *baseRepository[T]) unpackUInt() (uint, error) {
	chunk, err := rep.unpack(4)
	if err != nil {
		return 0, err
	}

	return uint(binary.LittleEndian.Uint32(chunk)), nil
}

// バイナリデータから uint を読み出す
func (rep *baseRepository[T]) unpackUInts(count int) ([]uint, error) {
	values := make([]uint, count)
	chunk, err := rep.unpack(4 * count)
	if err != nil {
		return values, err
	}

	for i := 0; i < count; i++ {
		values[i] = uint(binary.LittleEndian.Uint32(chunk[i*4 : (i+1)*4]))
	}

	return values, nil
}

// バイナリデータから int を読み出す
func (rep *baseRepository[T]) unpackInt() (int, error) {
	chunk, err := rep.unpack(4)
	if err != nil {
		return 0, err
	}

	return int(binary.LittleEndian.Uint32(chunk)), nil
}

// バイナリデータから int を読み出す
func (rep *baseRepository[T]) unpackInts(count int) ([]int, error) {
	values := make([]int, count)
	chunk, err := rep.unpack(4 * count)
	if err != nil {
		return values, err
	}

	for i := 0; i < count; i++ {
		values[i] = int(binary.LittleEndian.Uint32(chunk[i*4 : (i+1)*4]))
	}

	return values, nil
}

// バイナリデータから float64 を読み出す
func (rep *baseRepository[T]) unpackFloat() (float64, error) {
	// 単精度実数(4byte)なので、一旦uint32にしてからfloat32に変換する
	chunk, err := rep.unpack(4)
	if err != nil {
		return 0, err
	}

	return float64(math.Float32frombits(binary.LittleEndian.Uint32(chunk))), nil
}

// バイナリデータから float64 を複数個一気に読み出す
func (rep *baseRepository[T]) unpackFloats(values []float64, count int) ([]float64, error) {
	// 単精度実数(4byte)なので、一旦uint32にしてからfloat32に変換する
	chunk, err := rep.unpack(4 * count)
	if err != nil {
		return values, err
	}

	for i := 0; i < count; i++ {
		values[i] = float64(math.Float32frombits(binary.LittleEndian.Uint32(chunk[i*4 : (i+1)*4])))
	}

	return values, nil
}

func (rep *baseRepository[T]) unpackVec2() (*mmath.MVec2, error) {
	x, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec2(), err
	}
	y, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec2(), err
	}
	return &mmath.MVec2{X: x, Y: y}, nil
}

func (rep *baseRepository[T]) unpackVec3() (*mmath.MVec3, error) {
	x, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec3(), err
	}
	y, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec3(), err
	}
	z, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec3(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return &mmath.MVec3{X: x, Y: y, Z: z}, nil
}

func (rep *baseRepository[T]) unpackVec4() (*mmath.MVec4, error) {
	x, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec4(), err
	}
	y, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec4(), err
	}
	z, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec4(), err
	}
	w, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMVec4(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return &mmath.MVec4{X: x, Y: y, Z: z, W: w}, nil
}

func (rep *baseRepository[T]) unpackQuaternion() (*mmath.MQuaternion, error) {
	x, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMQuaternion(), err
	}
	y, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMQuaternion(), err
	}
	z, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMQuaternion(), err
	}
	w, err := rep.unpackFloat()
	if err != nil {
		return mmath.NewMQuaternion(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return mmath.NewMQuaternionByValues(x, y, z, w), nil
}

func (rep *baseRepository[T]) unpack(size int) ([]byte, error) {
	if rep.reader == nil {
		return nil, fmt.Errorf("file is not opened")
	}

	chunk := make([]byte, size)
	_, err := io.ReadFull(rep.reader, chunk)
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("EOF")
		}
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	return chunk, nil
}

type binaryType string

const (
	binaryType_float         binaryType = "<f"
	binaryType_byte          binaryType = "<b"
	binaryType_unsignedByte  binaryType = "<B"
	binaryType_short         binaryType = "<h"
	binaryType_unsignedShort binaryType = "<H"
	binaryType_int           binaryType = "<i"
	binaryType_unsignedInt   binaryType = "<I"
	binaryType_long          binaryType = "<l"
	binaryType_unsignedLong  binaryType = "<L"
)

func (rep *baseRepository[T]) writeNumber(
	fout *os.File, valType binaryType, val float64, defaultValue float64, isPositiveOnly bool) error {
	// 値の検証と修正
	if math.IsNaN(val) || math.IsInf(val, 0) {
		val = defaultValue
	}
	if isPositiveOnly && val < 0 {
		val = 0
	}

	// バイナリデータの作成
	var buf bytes.Buffer
	var err error
	switch valType {
	case binaryType_float:
		err = binary.Write(&buf, binary.LittleEndian, float32(val))
	case binaryType_unsignedInt:
		err = binary.Write(&buf, binary.LittleEndian, uint32(val))
	case binaryType_unsignedByte:
		err = binary.Write(&buf, binary.LittleEndian, uint8(val))
	case binaryType_unsignedShort:
		err = binary.Write(&buf, binary.LittleEndian, uint16(val))
	case binaryType_byte:
		err = binary.Write(&buf, binary.LittleEndian, int8(val))
	case binaryType_short:
		err = binary.Write(&buf, binary.LittleEndian, int16(val))
	default:
		err = binary.Write(&buf, binary.LittleEndian, int32(val))
	}
	if err != nil {
		return rep.writeDefaultNumber(fout, valType, defaultValue)
	}

	// ファイルへの書き込み
	_, err = fout.Write(buf.Bytes())
	if err != nil {
		return rep.writeDefaultNumber(fout, valType, defaultValue)
	}
	return nil
}

func (rep *baseRepository[T]) writeDefaultNumber(fout *os.File, valType binaryType, defaultValue float64) error {
	var buf bytes.Buffer
	var err error
	switch valType {
	case binaryType_float:
		err = binary.Write(&buf, binary.LittleEndian, float32(defaultValue))
	default:
		err = binary.Write(&buf, binary.LittleEndian, int32(defaultValue))
	}
	if err != nil {
		return err
	}
	_, err = fout.Write(buf.Bytes())
	return err
}

func (rep *baseRepository[T]) writeBool(fout *os.File, val bool) error {
	var buf bytes.Buffer
	var err error

	err = binary.Write(&buf, binary.LittleEndian, byte(mmath.BoolToInt(val)))

	if err != nil {
		return err
	}

	_, err = fout.Write(buf.Bytes())
	return err
}

func (rep *baseRepository[T]) writeByte(fout *os.File, val int, isUnsigned bool) error {
	var buf bytes.Buffer
	var err error

	if isUnsigned {
		err = binary.Write(&buf, binary.LittleEndian, uint8(val))
	} else {
		err = binary.Write(&buf, binary.LittleEndian, int8(val))
	}

	if err != nil {
		return err
	}

	_, err = fout.Write(buf.Bytes())
	return err
}

func (rep *baseRepository[T]) writeShort(fout *os.File, val uint16) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, val)
	if err != nil {
		return err
	}
	_, err = fout.Write(buf.Bytes())
	return err
}
