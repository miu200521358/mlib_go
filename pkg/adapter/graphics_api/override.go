// 指示: miu200521358
package graphics_api

// IOverrideRenderer はオーバーライド描画の抽象。
type IOverrideRenderer interface {
	// Bind は描画先をバインドする。
	Bind()
	// Unbind は描画先をアンバインドする。
	Unbind()
	// Resolve はサブウィンドウの描画結果を合成する。
	Resolve()
	// Resize は描画先サイズを変更する。
	Resize(width, height int)
	// Delete はリソースを解放する。
	Delete()
	// SetSharedTextureID は共有テクスチャIDを設定する。
	SetSharedTextureID(sharedTextureID *uint32)
	// SharedTextureIDPtr は共有テクスチャIDの参照を返す。
	SharedTextureIDPtr() *uint32
	// TextureIDPtr はテクスチャIDの参照を返す。
	TextureIDPtr() *uint32
}
