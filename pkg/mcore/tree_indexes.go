package mcore

// Comparer は比較可能な型に対するインターフェースです。
// このインターフェースを満たす任意の型を二分木に格納することができます。
type Comparer[T any] interface {
	// Compare は自身と引数の値を比較し、
	// 自身が小さい場合は負、等しい場合は0、大きい場合は正の整数を返します。
	Compare(T) int
}

// TreeIndexes は二分木のノードを表します。
type TreeIndexes[T Comparer[T]] struct {
	Value T
	Left  *TreeIndexes[T]
	Right *TreeIndexes[T]
}

func NewTreeIndexes[T Comparer[T]]() *TreeIndexes[T] {
	return &TreeIndexes[T]{}
}

func (n *TreeIndexes[T]) IsEmpty() bool {
	return n == nil
}

// GetValues は二分木の値をスライスにして返します。
func (n *TreeIndexes[T]) GetValues() []T {
	if n == nil {
		return nil
	}

	return append(append(n.Left.GetValues(), n.Value), n.Right.GetValues()...)
}

// Insert は二分木に新しい要素を挿入します。
func (n *TreeIndexes[T]) Insert(value T) {
	if n.Value.Compare(value) > 0 {
		if n.Left == nil {
			n.Left = &TreeIndexes[T]{Value: value}
		} else {
			n.Left.Insert(value)
		}
	} else {
		if n.Right == nil {
			n.Right = &TreeIndexes[T]{Value: value}
		} else {
			n.Right.Insert(value)
		}
	}
}

// Contains は二分木に値が含まれているかどうかを返します。
func (n *TreeIndexes[T]) Contains(value T) bool {
	if n == nil {
		return false
	}

	switch n.Value.Compare(value) {
	case 0:
		return true
	case 1:
		return n.Left.Contains(value)
	default:
		return n.Right.Contains(value)
	}
}

// Search は二分木から値を検索します。
func (n *TreeIndexes[T]) Search(value T) *TreeIndexes[T] {
	if n == nil {
		return nil
	}

	switch n.Value.Compare(value) {
	case 0:
		return n
	case 1:
		return n.Left.Search(value)
	default:
		return n.Right.Search(value)
	}
}

// SearchLeft は二分木から指定した値以下の最大の値を検索します。
func (n *TreeIndexes[T]) SearchLeft(value T) *TreeIndexes[T] {
	if n == nil {
		return nil
	}

	if n.Value.Compare(value) > 0 {
		return n.Left.SearchLeft(value)
	}

	if n.Right == nil || n.Right.Value.Compare(value) > 0 {
		return n
	}

	return n.Right.SearchLeft(value)
}

// SearchRight は二分木から指定した値以上の最小の値を検索します。
func (n *TreeIndexes[T]) SearchRight(value T) *TreeIndexes[T] {
	if n == nil {
		return nil
	}

	if n.Value.Compare(value) < 0 {
		return n.Right.SearchRight(value)
	}

	if n.Left == nil || n.Left.Value.Compare(value) < 0 {
		return n
	}

	return n.Left.SearchRight(value)
}

func (n *TreeIndexes[T]) GetMax() T {
	if n.Right == nil {
		return n.Value
	}
	return n.Right.GetMax()
}

func (n *TreeIndexes[T]) GetMin() T {
	if n.Left == nil {
		return n.Value
	}
	return n.Left.GetMin()
}

// Int は int 型のラッパーで、Comparer インターフェースを実装します。
type Int int

func NewInt(value int) Int {
	return Int(value)
}

func (a Int) Compare(b Int) int {
	return int(a) - int(b)
}

// Float64 は float64 型のラッパーで、Comparer インターフェースを実装します。
type Float64 float64

func NewFloat64(value float64) Float64 {
	return Float64(value)
}

func (a Float64) Compare(b Float64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// Float64 は float64 型のラッパーで、Comparer インターフェースを実装します。
type Float32 float64

var Float32Epsilon = 1e-6
var Float32Zero = Float32(float32(0))

func NewFloat32(value float32) Float32 {
	return Float32(value)
}

func (a Float32) Compare(b Float32) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
