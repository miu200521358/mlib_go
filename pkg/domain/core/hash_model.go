package core

type IHashModel interface {
	GetName() string
	SetName(name string)
	IsNotEmpty() bool
	IsEmpty() bool
	GetHash() string
	SetHash(hash string)
	GetPath() string
	SetPath(path string)
	Delete()
}

type HashModel struct {
	path string
	hash string
}

func NewHashModel(path string) *HashModel {
	return &HashModel{
		path: path,
		hash: "",
	}
}

func (m *HashModel) GetPath() string {
	return m.path
}

func (m *HashModel) SetPath(path string) {
	m.path = path
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
	return m.hash
}

func (m *HashModel) SetHash(hash string) {
	m.hash = hash
}

func (m *HashModel) IsNotEmpty() bool {
	// パスが定義されていたら、中身入り
	return len(m.path) > 0
}

func (m *HashModel) IsEmpty() bool {
	// パスが定義されていなかったら、空
	return len(m.path) == 0
}

func (m *HashModel) Delete() {
}
