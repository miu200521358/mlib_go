package mcore

// INamed は名前を持つモデルのインターフェースです。
type INamed interface {
	// Name はモデルの名前を返します。
	Name() string
	// SetName はモデルの名前を設定します。
	SetName(name string)
	// EnglishName はモデルの英語名を返します。
	EnglishName() string
	// SetEnglishName はモデルの英語名を設定します。
	SetEnglishName(englishName string)
}

// IIndexNameModel は名前とインデックスを持つモデルのインターフェースです。
// Bone, Morph, Material などが実装します。
type IIndexNameModel interface {
	IIndexModel
	INamed
}
