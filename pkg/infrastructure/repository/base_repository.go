package repository

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
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
	LoadHash(path string) (string, error)
	Save(overridePath string, data core.IHashModel, includeSystem bool) error
}

type baseRepository[T core.IHashModel] struct {
	file     *os.File
	reader   *bufio.Reader
	encoding encoding.Encoding
	readText func() string
	newFunc  func(path string) T
}

func (r *baseRepository[T]) open(path string) error {

	file, err := mutils.Open(path)
	if err != nil {
		return err
	}
	r.file = file
	r.reader = bufio.NewReader(r.file)

	return nil
}

func (r *baseRepository[T]) close() {
	defer r.file.Close()
}

func (r *baseRepository[T]) LoadName(path string) (string, error) {
	panic("not implemented")
}

func (r *baseRepository[T]) Load(path string) (T, error) {
	panic("not implemented")
}

func (r *baseRepository[T]) LoadHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sha1Hash := sha1.New()
	if _, err := io.Copy(sha1Hash, file); err != nil {
		return "", err
	}

	// ファイルパスをハッシュに含める
	sha1Hash.Write([]byte(path))

	return hex.EncodeToString(sha1Hash.Sum(nil)), nil

}

func (r *baseRepository[T]) defineEncoding(encoding encoding.Encoding) {
	r.encoding = encoding
	r.readText = r.defineReadText(encoding)
}

func (r *baseRepository[T]) defineReadText(encoding encoding.Encoding) func() string {
	return func() string {
		size, err := r.unpackInt()
		if err != nil {
			return ""
		}
		fbytes, err := r.unpackBytes(int(size))
		if err != nil {
			return ""
		}
		return r.decodeText(encoding, fbytes)
	}
}

func (r *baseRepository[T]) decodeText(mainEncoding encoding.Encoding, fbytes []byte) string {
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
			decodedText, err = r.decodeShiftJIS(fbytes)
			if err != nil {
				continue
			}
		} else {
			// 変換できなかった文字は「?」に変換する
			decodedText, err = r.decodeBytes(fbytes, targetEncoding)
			if err != nil {
				continue
			}
		}
		return decodedText
	}
	return ""
}

func (r *baseRepository[T]) decodeShiftJIS(fbytes []byte) (string, error) {
	decodedText, err := japanese.ShiftJIS.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

func (r *baseRepository[T]) decodeBytes(fbytes []byte, encoding encoding.Encoding) (string, error) {
	decodedText, err := encoding.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

// バイナリデータから bytes を読み出す
func (r *baseRepository[T]) unpackBytes(size int) ([]byte, error) {
	chunk, err := r.unpack(size)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// バイナリデータから byte を読み出す
func (r *baseRepository[T]) unpackByte() (byte, error) {
	chunk, err := r.unpack(1)
	if err != nil {
		return 0, err
	}

	return chunk[0], nil
}

// バイナリデータから sbyte を読み出す
func (r *baseRepository[T]) unpackSByte() (int8, error) {
	chunk, err := r.unpack(1)
	if err != nil {
		return 0, err
	}

	return int8(chunk[0]), nil
}

// バイナリデータから int16 を読み出す
func (r *baseRepository[T]) unpackShort() (int16, error) {
	chunk, err := r.unpack(2)
	if err != nil {
		return 0, err
	}

	return int16(binary.LittleEndian.Uint16(chunk)), nil
}

// バイナリデータから uint16 を読み出す
func (r *baseRepository[T]) unpackUShort() (uint16, error) {
	chunk, err := r.unpack(2)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(chunk), nil
}

// バイナリデータから uint を読み出す
func (r *baseRepository[T]) unpackUInt() (uint, error) {
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return uint(binary.LittleEndian.Uint32(chunk)), nil
}

// バイナリデータから int を読み出す
func (r *baseRepository[T]) unpackInt() (int, error) {
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return int(binary.LittleEndian.Uint32(chunk)), nil
}

// バイナリデータから float64 を読み出す
func (r *baseRepository[T]) unpackFloat() (float64, error) {
	// 単精度実数(4byte)なので、一旦uint32にしてからfloat32に変換する
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return float64(math.Float32frombits(binary.LittleEndian.Uint32(chunk))), nil
}

func (r *baseRepository[T]) unpackVec2() (mmath.MVec2, error) {
	x, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec2(), err
	}
	y, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec2(), err
	}
	return mmath.MVec2{x, y}, nil
}

func (r *baseRepository[T]) unpackVec3() (mmath.MVec3, error) {
	x, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec3(), err
	}
	y, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec3(), err
	}
	z, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec3(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return mmath.MVec3{x, y, z}, nil
}

func (r *baseRepository[T]) unpackVec4() (mmath.MVec4, error) {
	x, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	y, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	z, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	w, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return mmath.MVec4{x, y, z, w}, nil
}

func (r *baseRepository[T]) unpackQuaternion() (mmath.MQuaternion, error) {
	x, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	y, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	z, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	w, err := r.unpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return *mmath.NewMQuaternionByValues(x, y, z, w), nil
}

func (r *baseRepository[T]) unpack(size int) ([]byte, error) {
	if r.reader == nil {
		return nil, fmt.Errorf("file is not opened")
	}

	chunk := make([]byte, size)
	_, err := io.ReadFull(r.reader, chunk)
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

func (r *baseRepository[T]) writeNumber(
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
		return r.writeDefaultNumber(fout, valType, defaultValue)
	}

	// ファイルへの書き込み
	_, err = fout.Write(buf.Bytes())
	if err != nil {
		return r.writeDefaultNumber(fout, valType, defaultValue)
	}
	return nil
}

func (r *baseRepository[T]) writeDefaultNumber(fout *os.File, valType binaryType, defaultValue float64) error {
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

func (r *baseRepository[T]) writeBool(fout *os.File, val bool) error {
	var buf bytes.Buffer
	var err error

	err = binary.Write(&buf, binary.LittleEndian, byte(mmath.BoolToInt(val)))

	if err != nil {
		return err
	}

	_, err = fout.Write(buf.Bytes())
	return err
}

func (r *baseRepository[T]) writeByte(fout *os.File, val int, isUnsigned bool) error {
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

func (r *baseRepository[T]) writeShort(fout *os.File, val uint16) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, val)
	if err != nil {
		return err
	}
	_, err = fout.Write(buf.Bytes())
	return err
}