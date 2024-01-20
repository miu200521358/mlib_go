package reader

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"

	"github.com/miu200521358/mlib_go/pkg/core/hash_model"
	"github.com/miu200521358/mlib_go/pkg/math/mquaternion"
	"github.com/miu200521358/mlib_go/pkg/math/mvec2"
	"github.com/miu200521358/mlib_go/pkg/math/mvec3"
	"github.com/miu200521358/mlib_go/pkg/math/mvec4"
	"github.com/miu200521358/mlib_go/pkg/utils/util_file"
)

type ReaderInterface interface {
	ReadNameByFilepath(path string) (string, error)
	ReadByFilepath(path string) (hash_model.HashModelInterface, error)
}

type BaseReader[T hash_model.HashModelInterface] struct {
	file     *os.File
	reader   *bufio.Reader
	encoding encoding.Encoding
	ReadText func() string
}

func (r *BaseReader[T]) Open(path string) error {

	file, err := util_file.Open(path)
	if err != nil {
		return err
	}
	r.file = file
	r.reader = bufio.NewReader(r.file)

	return nil
}

func (r *BaseReader[T]) Close() {
	defer r.file.Close()
}

func (r *BaseReader[T]) ReadNameByFilepath(path string) (string, error) {
	panic("not implemented")
}

func (r *BaseReader[T]) ReadByFilepath(path string) (T, error) {
	panic("not implemented")
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

func (r *BaseReader[T]) decodeShiftJIS(fbytes []byte) (string, error) {
	decodedText, err := japanese.ShiftJIS.NewDecoder().Bytes(fbytes)
	if err != nil {
		return "", err
	}
	return string(decodedText), nil
}

func (r *BaseReader[T]) decodeBytes(fbytes []byte, encoding encoding.Encoding) (string, error) {
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

	value := int16(binary.LittleEndian.Uint16(chunk))
	return value, nil
}

// バイナリデータから uint16 を読み出す
func (r *BaseReader[T]) UnpackUShort() (uint16, error) {
	chunk, err := r.unpack(2)
	if err != nil {
		return 0, err
	}

	value := binary.LittleEndian.Uint16(chunk)
	return value, nil
}

// バイナリデータから int を読み出す
func (r *BaseReader[T]) UnpackInt() (int, error) {
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	value := int(binary.LittleEndian.Uint32(chunk))
	return value, nil
}

// バイナリデータから float64 を読み出す
func (r *BaseReader[T]) UnpackFloat() (float64, error) {
	// 単精度実数(4byte)なので、一旦uint32にしてからfloat32に変換する
	chunk, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	value := float64(math.Float32frombits(binary.LittleEndian.Uint32(chunk)))
	return value, nil
}

func (r *BaseReader[T]) UnpackVec2() (mvec2.T, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return mvec2.T{}, err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return mvec2.T{}, err
	}
	return mvec2.T{x, y}, nil
}

func (r *BaseReader[T]) UnpackVec3() (mvec3.T, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return mvec3.T{}, err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return mvec3.T{}, err
	}
	z, err := r.UnpackFloat()
	if err != nil {
		return mvec3.T{}, err
	}
	return mvec3.T{x, y, z}, nil
}

func (r *BaseReader[T]) UnpackVec4() (mvec4.T, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return mvec4.T{}, err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return mvec4.T{}, err
	}
	z, err := r.UnpackFloat()
	if err != nil {
		return mvec4.T{}, err
	}
	w, err := r.UnpackFloat()
	if err != nil {
		return mvec4.T{}, err
	}
	return mvec4.T{x, y, z, w}, nil
}

func (r *BaseReader[T]) UnpackQuaternion() (mquaternion.T, error) {
	x, err := r.UnpackFloat()
	if err != nil {
		return mquaternion.T{}, err
	}
	y, err := r.UnpackFloat()
	if err != nil {
		return mquaternion.T{}, err
	}
	z, err := r.UnpackFloat()
	if err != nil {
		return mquaternion.T{}, err
	}
	w, err := r.UnpackFloat()
	if err != nil {
		return mquaternion.T{}, err
	}
	return mquaternion.T{x, y, z, w}, nil
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
