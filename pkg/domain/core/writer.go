package core

type IWriter interface {
	Save(overridePath string) error
}
