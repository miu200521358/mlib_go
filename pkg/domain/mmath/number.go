package mmath

// Number は数値型の制約を定義するインターフェースです。
// ジェネリック関数で使用されます。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// SignedNumber は符号付き数値型の制約を定義するインターフェースです。
type SignedNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Float は浮動小数点数型の制約を定義するインターフェースです。
type Float interface {
	~float32 | ~float64
}
