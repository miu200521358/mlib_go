package model

// TextureType はテクスチャ種別を表す。
type TextureType int

const (
	// TEXTURE_TYPE_TEXTURE は通常テクスチャ。
	TEXTURE_TYPE_TEXTURE TextureType = iota
	// TEXTURE_TYPE_TOON はトゥーンテクスチャ。
	TEXTURE_TYPE_TOON
	// TEXTURE_TYPE_SPHERE はスフィアテクスチャ。
	TEXTURE_TYPE_SPHERE
)

// Texture はテクスチャ要素を表す。
type Texture struct {
	index       int
	name        string
	EnglishName string
	TextureType TextureType
	valid       bool
}

// NewTexture は既定値で Texture を生成する。
func NewTexture() *Texture {
	return &Texture{
		index:       -1,
		name:        "",
		EnglishName: "",
		TextureType: TEXTURE_TYPE_TEXTURE,
		valid:       false,
	}
}

// Index はテクスチャ index を返す。
func (t *Texture) Index() int {
	return t.index
}

// SetIndex はテクスチャ index を設定する。
func (t *Texture) SetIndex(index int) {
	t.index = index
}

// Name はテクスチャ名を返す。
func (t *Texture) Name() string {
	return t.name
}

// SetName はテクスチャ名を設定する。
func (t *Texture) SetName(name string) {
	t.name = name
}

// IsValid はテクスチャが有効か判定する。
func (t *Texture) IsValid() bool {
	return t != nil && t.valid
}

// SetValid は有効フラグを設定する。
func (t *Texture) SetValid(valid bool) {
	t.valid = valid
}
