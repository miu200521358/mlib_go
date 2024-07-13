package core

import (
	"bufio"
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

	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type IReader interface {
	ReadNameByFilepath(path string) (string, error)
	ReadByFilepath(path string) (IHashModel, error)
	ReadHashByFilePath(path string) (string, error)
}

type BaseReader[T IHashModel] struct {
	File     *os.File
	reader   *bufio.Reader
	encoding encoding.Encoding
	ReadText func() string
}

func (r *BaseReader[T]) Open(path string) error {

	file, err := mutils.Open(path)
	if err != nil {
		return err
	}
	r.File = file
	r.reader = bufio.NewReader(r.File)

	return nil
}

func (r *BaseReader[T]) Close() {
	defer r.File.Close()
}

func (r *BaseReader[T]) ReadNameByFilepath(path string) (string, error) {
	panic("not implemented")
}

func (r *BaseReader[T]) ReadByFilepath(path string) (T, error) {
	panic("not implemented")
}

func (r *BaseReader[T]) ReadHashByFilePath(path string) (string, error) {
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

func (r *BaseReader[T]) DefineEncoding(encoding encoding.Encoding) {
	r.encoding = encoding
	r.ReadText = r.defineReadText(encoding)
}

func (r *BaseReader[T]) defineReadText(encoding encoding.Encoding) func() string {
	return func() string {
		size, err := r.UnpackInt()
		if err != nil {
			return ""
		}
		fbytes, err := r.UnpackBytes(int(size))
		if err != nil {
			return ""
		}
		return r.DecodeText(encoding, fbytes)
	}
}

func (r *BaseReader[T]) DecodeText(mainEncoding encoding.Encoding, fbytes []byte) string {
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
			decodedText, err = r.DecodeShiftJIS(fbytes)
			if err != nil {
				continue
			}
		} else {
			// 変換できなかった文字は「?」に変換する
			decodedText, err = r.DecodeBytes(fbytes, targetEncoding)
			if err != nil {
				continue
			}
		}
		return decodedText
	}
	return ""
}

func (r *BaseReader[T]) DecodeShiftJIS(fbytes []byte) (string, error) {
	decodedText, err := japanese.ShiftJIS.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

func (r *BaseReader[T]) DecodeBytes(fbytes []byte, encoding encoding.Encoding) (string, error) {
	decodedText, err := encoding.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

// バイナリデータから bytes を読み出す
func (r *BaseReader[T]) UnpackBytes(size int) ([]byte, error) {
	chunk, err := r.unpack(size)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

// バイナリデータから byte を読み出す
func (r *BaseReader[T]) UnpackByte() (byte, error) {
	chunk, err := r.unpack(1)
	if err != nil {
		return 0, err
	}

	return chunk[0], nil
}

// バイナリデータから sbyte を読み出す
func (r *BaseReader[T]) UnpackSByte() (int8, error) {
	chunk, err := r.unpack(1)
	if err != nil {
		return 0, err
	}

	return int8(chunk[0]), nil
}

// バイナリデータから int16 を読み出す
func (r *BaseReader[T]) UnpackShort() (int16, error) {
	chunk, err := r.unpack(2)
	if err != nil {
		return 0, err
	}

	return int16(binary.LittleEndian.Uint16(chunk)), nil
}

// バイナリデータから uint16 を読み出す
func (r *BaseReader[T]) UnpackUShort() (uint16, error) {
	chunk, err := r.unpack(2)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(chunk), nil
}

// バイナリデータから uint を読み出す
func (r *BaseReader[T]) UnpackUInt() (uint, error) {
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return uint(binary.LittleEndian.Uint32(chunk)), nil
}

// バイナリデータから int を読み出す
func (r *BaseReader[T]) UnpackInt() (int, error) {
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return int(binary.LittleEndian.Uint32(chunk)), nil
}

// バイナリデータから float64 を読み出す
func (r *BaseReader[T]) UnpackFloat() (float64, error) {
	// 単精度実数(4byte)なので、一旦uint32にしてからfloat32に変換する
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return float64(math.Float32frombits(binary.LittleEndian.Uint32(chunk))), nil
}

func (r *BaseReader[T]) UnpackVec2() (mmath.MVec2, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec2(), err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec2(), err
	}
	return mmath.MVec2{x, y}, nil
}

func (r *BaseReader[T]) UnpackVec3(isConvertGl bool) (mmath.MVec3, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec3(), err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec3(), err
	}
	z, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec3(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return mmath.MVec3{x, y, z}, nil
}

func (r *BaseReader[T]) UnpackVec4(isConvertGl bool) (mmath.MVec4, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	z, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	w, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMVec4(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return mmath.MVec4{x, y, z, w}, nil
}

func (r *BaseReader[T]) UnpackQuaternion(isConvertGl bool) (mmath.MQuaternion, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	z, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	w, err := r.UnpackFloat()
	if err != nil {
		return *mmath.NewMQuaternion(), err
	}
	// if isConvertGl {
	// 	z = -z
	// }
	return *mmath.NewMQuaternionByValues(x, y, z, w), nil
}

func (r *BaseReader[T]) unpack(size int) ([]byte, error) {
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
