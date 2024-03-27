package mcore

type IHashModel interface {
	GetName() string
	SetName(name string)
	IsNotEmpty() bool
	IsEmpty() bool
	GetHash() string
	GetPath() string
	SetPath(path string)
}

type HashModel struct {
	Path string
	Hash string
}

func NewHashModel(path string) *HashModel {
	return &HashModel{
		Path: path,
		Hash: "",
	}
}

func (m *HashModel) GetPath() string {
	return m.Path
}

func (m *HashModel) SetPath(path string) {
	m.Path = path
}

func (m *HashModel) GetName() string {
	// モデル内の名前に相当する値を返す
	panic("not implemented")
}

func (m *HashModel) SetName(name string) {
	// モデル内の名前に相当する値を設定する
	panic("not implemented")
}

func (m *HashModel) GetHash() string {
	return m.Hash
}

func (m *HashModel) IsNotEmpty() bool {
	// パスが定義されていたら、中身入り
	return len(m.Path) > 0
}

func (m *HashModel) IsEmpty() bool {
	// パスが定義されていなかったら、空
	return len(m.Path) == 0
}
