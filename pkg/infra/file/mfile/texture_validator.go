// 指示: miu200521358
package mfile

// TextureValidator はテクスチャ検証処理を表す。
type TextureValidator struct{}

// NewTextureValidator はTextureValidatorを生成する。
func NewTextureValidator() *TextureValidator {
	return &TextureValidator{}
}

// ExistsFile はファイルの存在を判定する。
func (v *TextureValidator) ExistsFile(path string) (bool, error) {
	return ExistsFile(path)
}

// ValidateImage は画像として読み込めるか検証する。
func (v *TextureValidator) ValidateImage(path string) error {
	_, err := LoadImage(path)
	return err
}
