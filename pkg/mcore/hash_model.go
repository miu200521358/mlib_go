package mcore

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

type HashModelInterface interface {
	GetName() string
	SetName(name string)
	IsNotEmpty() bool
	IsEmpty() bool
	GetDigest() string
	UpdateDigest() error
}

type HashModel struct {
	Path   string
	Digest string
}

func NewHashModel(path string) *HashModel {
	return &HashModel{
		Path:   path,
		Digest: "",
	}
}

func (m *HashModel) GetName() string {
	// モデル内の名前に相当する値を返す
	panic("not implemented")
}

func (m *HashModel) SetName(name string) {
	// モデル内の名前に相当する値を設定する
	panic("not implemented")
}

func (m *HashModel) GetDigest() string {
	return m.Digest
}

func (m *HashModel) UpdateDigest() error {
	file, err := os.Open(m.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	sha1Hash := sha1.New()
	if _, err := io.Copy(sha1Hash, file); err != nil {
		return err
	}

	// ファイルパスをハッシュに含める
	sha1Hash.Write([]byte(m.Path))

	m.Digest = hex.EncodeToString(sha1Hash.Sum(nil))

	return nil
}

func (m *HashModel) IsNotEmpty() bool {
	// パスが定義されていたら、中身入り
	return len(m.Path) > 0
}

func (m *HashModel) IsEmpty() bool {
	// パスが定義されていなかったら、空
	return len(m.Path) == 0
}
