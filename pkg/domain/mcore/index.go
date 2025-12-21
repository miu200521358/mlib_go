// Package mcore はドメイン層で共通して使用するインターフェースとコレクション型を提供します。
package mcore

// IIndexable はインデックスを持つモデルの基本インターフェースです。
type IIndexable interface {
	// Index はモデルのインデックスを返します。
	Index() int
	// SetIndex はモデルのインデックスを設定します。
	SetIndex(index int)
}

// IValidatable は有効性を検証できるモデルのインターフェースです。
type IValidatable interface {
	// IsValid はモデルが有効かどうかを返します。
	IsValid() bool
}

// IIndexModel はインデックスを持つモデルの基本インターフェースです。
// PMXモデルやVMDモーションの各要素が実装します。
type IIndexModel interface {
	IIndexable
	IValidatable
}
