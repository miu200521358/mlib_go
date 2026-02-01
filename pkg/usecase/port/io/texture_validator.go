// 指示: miu200521358
package io

// ITextureValidator はテクスチャ検証の契約を表す。
type ITextureValidator interface {
	// ExistsFile はファイルの存在を判定する。
	ExistsFile(path string) (bool, error)
	// ValidateImage は画像として読み込めるか検証する。
	ValidateImage(path string) error
}
