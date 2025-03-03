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
