package core

type IHashModel interface {
	Name() string
	SetName(name string)
	Hash() string
	SetHash(hash string)
	Path() string
	SetPath(path string)
	UpdateHash()
}
