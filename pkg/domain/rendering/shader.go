//go:build windows
// +build windows

package rendering

// IShader はシェーダー機能の抽象インターフェース
type IShader interface {
	// 基本操作
	Resize(width, height int)
	Cleanup()

	// プログラム関連
	GetProgram(programType ProgramType) uint32
	UseProgram(programType ProgramType)
	ResetProgram()

	// テクスチャ関連
	GetBoneTextureID() uint32
	GetOverrideTextureID() uint32

	// カメラ設定
	GetCamera() *Camera
	SetCamera(*Camera)
	UpdateCamera()

	// 床描画機能
	DrawFloor()

	// MSAA関連
	GetMsaa() IMsaa
}

// IShaderFactory はシェーダー生成の抽象ファクトリー
type IShaderFactory interface {
	CreateShader(width, height int) (IShader, error)
}
