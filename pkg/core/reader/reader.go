package reader

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"

	"github.com/miu200521358/mlib_go/pkg/core/hash_model"
)

type ReaderInterface interface {
	ReadNameByFilepath(path string) (string, error)
	ReadHashByFilepath(path string) (string, error)
	ReadByFilepath(path string) (hash_model.HashModelInterface, error)
	CreateModel(path string) hash_model.HashModelInterface
	ReadHeader(hashData hash_model.HashModelInterface)
	ReadData(hashData hash_model.HashModelInterface)
	DefineEncoding(encoding encoding.Encoding)
}

type BaseReader[T hash_model.HashModelInterface] struct {
	reader   *bufio.Reader
	encoding encoding.Encoding
	readText func() string
}

// 指定されたパスのファイルから該当名称を読み込む
func (r *BaseReader[T]) ReadNameByFilepath(path string) (string, error) {
	if !fileExists(path) {
		return "", nil
	}

	// モデルを新規作成
	model := r.createModel(path)

	err := r.open(path)
	if err != nil {
		return "", err
	}

	r.ReadHeader(model)

	return model.GetName(), nil
}

// 指定されたパスのファイルからハッシュデータを読み込む
func (r *BaseReader[T]) ReadHashByFilepath(path string) string {
	if !fileExists(path) {
		return ""
	}

	// モデルを新規作成
	model := r.createModel(path)

	err := model.UpdateDigest()
	if err != nil {
		return ""
	}

	// ハッシュデータを読み取り
	return model.GetDigest()
}

// 指定されたパスのファイルからデータを読み込む
func (r *BaseReader[T]) ReadByFilepath(path string) (T, error) {
	// モデルを新規作成
	hashData := r.createModel(path)

	if !fileExists(path) {
		return hashData, nil
	}

	err := r.open(path)
	if err != nil {
		return hashData, err
	}

	err = hashData.UpdateDigest()
	if err != nil {
		return hashData, err
	}

	r.ReadHeader(hashData)
	r.ReadData(hashData)

	return hashData, nil
}

// 指定されたパスがファイルとして存在しているか
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (r *BaseReader[T]) open(path string) error {
	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	r.reader = bufio.NewReader(file)

	return nil
}

func (r *BaseReader[T]) createModel(path string) T {
	panic("not implemented")
}

func (r *BaseReader[T]) ReadHeader(hashData T) {
	panic("not implemented")
}

func (r *BaseReader[T]) ReadData(hashData T) {
	panic("not implemented")
}

func (r *BaseReader[T]) DefineEncoding(encoding encoding.Encoding) {
	r.encoding = encoding
	r.readText = r.defineReadText(encoding)
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

// unpackInt はバイナリデータから int32 を読み出す
func (r *BaseReader[T]) UnpackBytes(size int) ([]byte, error) {
	v, err := r.unpack(size)
	if err != nil {
		return nil, err
	}

	return v.([]byte), nil
}

// UnpackInt はバイナリデータから int32 を読み出す
func (r *BaseReader[T]) UnpackInt() (int32, error) {
	v, err := r.unpack(4)
	if err != nil {
		return 0, err
	}

	return v.(int32), nil
}

// unpack はバイナリデータから any を読み出す
func (r *BaseReader[T]) unpack(size int) (any, error) {
	chunk := make([]byte, size)
	n, err := r.reader.Read(chunk)
	if err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}
	}
	if n < size {
		return nil, err
	}

	// バイナリデータから Block 構造体への変換
	buffer := bytes.NewReader(chunk)
	var v any

	// データ の読み出し
	if err := binary.Read(buffer, binary.LittleEndian, &v); err != nil {
		return nil, err
	}

	return v, nil
}
