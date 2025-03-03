//go:build windows
// +build windows

package rendering

type ProgramType int

const (
	PROGRAM_TYPE_MODEL           ProgramType = iota // モデル
	PROGRAM_TYPE_EDGE                               // エッジ
	PROGRAM_TYPE_BONE                               // ボーン
	PROGRAM_TYPE_PHYSICS                            // 物理
	PROGRAM_TYPE_NORMAL                             // 法線
	PROGRAM_TYPE_FLOOR                              // 床
	PROGRAM_TYPE_WIRE                               // ワイヤー
	PROGRAM_TYPE_SELECTED_VERTEX                    // 選択頂点
	PROGRAM_TYPE_OVERRIDE                           // 重ねて描画
	PROGRAM_TYPE_CURSOR                             // 頂点選択時のライン
)
