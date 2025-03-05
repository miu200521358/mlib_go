package rendering

// IMsaa はマルチサンプルアンチエイリアシング機能のインターフェース
type IMsaa interface {
	// Bind はMSAAフレームバッファをバインドする
	Bind()

	// Unbind はMSAAフレームバッファをアンバインドする
	Unbind()

	// Resolve はMSAAの結果をデフォルトフレームバッファに解決する
	Resolve()

	// ReadDepthAt は指定座標の深度値を読み取る
	ReadDepthAt(x, y int) float32

	// SetOverrideTargetTexture はオーバーライドターゲットテクスチャを設定する
	SetOverrideTargetTexture(textureID uint32)

	// OverrideTargetTexture はオーバーライドターゲットテクスチャのIDを取得する
	OverrideTargetTexture() uint32

	// BindOverrideTexture はオーバーライドテクスチャをバインドする
	BindOverrideTexture(windowIndex int, programID uint32)

	// UnbindOverrideTexture はオーバーライドテクスチャをアンバインドする
	UnbindOverrideTexture()

	// Delete はMSAAリソースを解放する
	Delete()

	// Resize はMSAAバッファのサイズを変更する
	Resize(width, height int)
}

// MSAAConfig はMSAA設定を保持する構造体
type MSAAConfig struct {
	Width       int
	Height      int
	SampleCount int
}
