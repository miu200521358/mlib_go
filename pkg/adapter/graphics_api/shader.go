// 指示: miu200521358
package graphics_api

// ProgramType はシェーダープログラム種別を表す。
type ProgramType int

const (
	// ProgramTypeModel はモデル用。
	ProgramTypeModel ProgramType = iota
	// ProgramTypeEdge はエッジ用。
	ProgramTypeEdge
	// ProgramTypeBone はボーン用。
	ProgramTypeBone
	// ProgramTypePhysics は物理剛体用。
	ProgramTypePhysics
	// ProgramTypeNormal は法線用。
	ProgramTypeNormal
	// ProgramTypeFloor は床用。
	ProgramTypeFloor
	// ProgramTypeWire はワイヤーフレーム用。
	ProgramTypeWire
	// ProgramTypeSelectedVertex は選択頂点用。
	ProgramTypeSelectedVertex
	// ProgramTypeOverride はオーバーライド用。
	ProgramTypeOverride
	// ProgramTypeCursor はカーソル用。
	ProgramTypeCursor
)

// IShader はシェーダー機能の抽象インターフェース。
type IShader interface {
	// Resize はビューポートと付随リソースを更新する。
	Resize(width, height int)
	// Cleanup はシェーダーと関連リソースを解放する。
	Cleanup()
	// Program はプログラムIDを返す。
	Program(programType ProgramType) uint32
	// UseProgram は指定プログラムを使用する。
	UseProgram(programType ProgramType)
	// ResetProgram はプログラム利用を解除する。
	ResetProgram()
	// BoneTextureID はボーン行列テクスチャIDを返す。
	BoneTextureID() uint32
	// OverrideTextureID はオーバーライド用テクスチャIDを返す。
	OverrideTextureID() uint32
	// Camera はカメラを返す。
	Camera() *Camera
	// SetCamera はカメラを設定する。
	SetCamera(camera *Camera)
	// UpdateCamera はシェーダーへカメラを反映する。
	UpdateCamera()
	// Msaa はMSAA実装を返す。
	Msaa() IMsaa
	// SetMsaa はMSAA実装を設定する。
	SetMsaa(msaa IMsaa)
	// FloorRenderer は床描画を返す。
	FloorRenderer() IFloorRenderer
	// OverrideRenderer はオーバーライド描画を返す。
	OverrideRenderer() IOverrideRenderer
}

// IShaderFactory はシェーダー生成の抽象ファクトリー。
type IShaderFactory interface {
	// CreateShader はシェーダーを生成する。
	CreateShader(width, height int) (IShader, error)
}
