// 指示: miu200521358
package graphics_api

// IMsaa はマルチサンプルアンチエイリアシングの抽象。
type IMsaa interface {
	// Bind はMSAAフレームバッファをバインドする。
	Bind()
	// Unbind はMSAAフレームバッファをアンバインドする。
	Unbind()
	// Resolve はMSAAの結果を解決する。
	Resolve()
	// ReadDepthAt は指定座標の深度値を読み取る。
	ReadDepthAt(x, y, width, height int) float32
	// Delete はMSAAリソースを解放する。
	Delete()
	// Resize はMSAAバッファのサイズを変更する。
	Resize(width, height int)
}
