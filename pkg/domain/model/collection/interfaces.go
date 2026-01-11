package collection

// IIndexable は可変な index を持つ要素を表す。
type IIndexable interface {
	Index() int
	SetIndex(index int)
	IsValid() bool
}

// INameable は名前と index を持つ要素を表す。
type INameable interface {
	IIndexable
	Name() string
	SetName(name string)
}
