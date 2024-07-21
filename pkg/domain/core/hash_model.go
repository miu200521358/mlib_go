package core

type IHashModel interface {
	Name() string
	SetName(name string)
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

func (hModel *HashModel) Path() string {
	return hModel.path
}

func (hModel *HashModel) SetPath(path string) {
	hModel.path = path
}

// モデル内の名前に相当する値を返す
func (hModel *HashModel) Name() string {
	return hModel.name
}

// モデル内の名前に相当する値を設定する
func (hModel *HashModel) SetName(name string) {
	hModel.name = name
}

func (hModel *HashModel) Hash() string {
	return hModel.hash
}

func (hModel *HashModel) SetHash(hash string) {
	hModel.hash = hash
}

func (hModel *HashModel) Delete() {
}
