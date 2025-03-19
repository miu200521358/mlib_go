package rendering

type IOverrideRenderer interface {
	Bind()

	Unbind()

	// メインウィンドウでサブウィンドウの描画内容を書き込んだテクスチャを描画する
	Resolve()

	Resize(width, height int)

	Delete()

	SetSharedTextureID(sharedTextureID *uint32)

	SharedTextureIDPtr() *uint32

	TextureIDPtr() *uint32
}
