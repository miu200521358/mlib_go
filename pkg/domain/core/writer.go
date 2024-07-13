package core

type IWriter interface {
	Save(data IHashModel, overridePath, overrideName string, includeSystem bool) error
}
