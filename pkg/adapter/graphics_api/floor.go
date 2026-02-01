// 指示: miu200521358
package graphics_api

// IFloorRenderer は床描画の抽象。
type IFloorRenderer interface {
	// Bind は床描画をバインドする。
	Bind()
	// Unbind は床描画をアンバインドする。
	Unbind()
	// Render は床を描画する。
	Render()
	// Delete は床描画リソースを解放する。
	Delete()
}
