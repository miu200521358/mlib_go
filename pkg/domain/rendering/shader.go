//go:build windows
// +build windows

package rendering

type ProgramType int

const (
	ProgramTypeModel          ProgramType = iota // モデル
	ProgramTypeEdge                              // エッジ
	ProgramTypeBone                              // ボーン
	ProgramTypePhysics                           // 物理剛体
	ProgramTypeNormal                            // 法線
	ProgramTypeFloor                             // 床
	ProgramTypeWire                              // ワイヤーフレーム
	ProgramTypeSelectedVertex                    // 選択頂点
	ProgramTypeOverride                          // ウィンドウを重ねて描画
	ProgramTypeCursor                            // カーソル
)

// IShader はシェーダー機能の抽象インターフェース
type IShader interface {
	// 基本操作
	Resize(width, height int)
	Cleanup()

	// プログラム関連
	Program(programType ProgramType) uint32
	UseProgram(programType ProgramType)
	ResetProgram()

	// テクスチャ関連
	BoneTextureID() uint32
	OverrideTextureID() uint32

	// カメラ設定
	Camera() *Camera
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
