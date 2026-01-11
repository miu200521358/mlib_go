package collection

// NameIndex は name -> index の参照を提供する。
type NameIndex[T INameable] struct {
	indexByName map[string]int
}

// NewNameIndex は NameIndex を生成する。
func NewNameIndex[T INameable]() *NameIndex[T] {
	return &NameIndex[T]{indexByName: make(map[string]int)}
}

// GetByName は name に対応する index を返す。
func (n *NameIndex[T]) GetByName(name string) (int, bool) {
	if n == nil {
		return 0, false
	}
	idx, ok := n.indexByName[name]
	return idx, ok
}

// Names はユニークな名前一覧を返す（順序は保証しない）。
func (n *NameIndex[T]) Names() []string {
	if n == nil {
		return nil
	}
	out := make([]string, 0, len(n.indexByName))
	for name := range n.indexByName {
		out = append(out, name)
	}
	return out
}

// SetIfAbsent は name が存在しない場合のみ name/index を追加する。
func (n *NameIndex[T]) SetIfAbsent(name string, index int) bool {
	if n == nil {
		return false
	}
	if _, ok := n.indexByName[name]; ok {
		return false
	}
	n.indexByName[name] = index
	return true
}

// Rebuild は先勝ちを維持しながら index マップを再構築する。
func (n *NameIndex[T]) Rebuild(values []T) {
	if n == nil {
		return
	}
	n.indexByName = make(map[string]int)
	for _, v := range values {
		if !v.IsValid() {
			continue
		}
		name := v.Name()
		if _, ok := n.indexByName[name]; ok {
			continue
		}
		n.indexByName[name] = v.Index()
	}
}
