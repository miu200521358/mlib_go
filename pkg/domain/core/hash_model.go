package core

type IHashModel interface {
	Name() string
	SetName(name string)
	IsNotEmpty() bool
	IsEmpty() bool
	Hash() string
	SetHash(hash string)
	Path() string
	SetPath(path string)
	Delete()
}

type HashModel struct {
	name string
	path string
	hash string
}

func NewHashModel(path string) *HashModel {
	return &HashModel{
		path: path,
		hash: "",
	}
}

func (m *HashModel) Path() string {
	return m.path
}

func (m *HashModel) SetPath(path string) {
	m.path = path
}

// モデル内の名前に相当する値を返す
func (m *HashModel) Name() string {
	return m.name
}

// モデル内の名前に相当する値を設定する
func (m *HashModel) SetName(name string) {
	m.name = name
}

func (m *HashModel) Hash() string {
	return m.hash
}

func (m *HashModel) SetHash(hash string) {
	m.hash = hash
}

// パスが定義されていたら、中身入り
func (m *HashModel) IsNotEmpty() bool {
	return len(m.path) > 0
}

// パスが定義されていなかったら、空
func (m *HashModel) IsEmpty() bool {
	return len(m.path) == 0
}

func (m *HashModel) Delete() {
}
