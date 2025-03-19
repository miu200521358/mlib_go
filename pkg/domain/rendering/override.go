package rendering

type IOverrideRenderer interface {
	Bind()

	Unbind()

	// サブウィンドウの描画内容をテクスチャに書き込む
	Render()

	// メインウィンドウでサブウィンドウの描画内容を書き込んだテクスチャを描画する
	Resolve()

	Resize(width, height int)

	Delete()

	SetSharedTextureID(sharedTextureID *uint32)

	TextureID() uint32
}
