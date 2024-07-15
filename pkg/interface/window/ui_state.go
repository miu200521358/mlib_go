package window

import "github.com/miu200521358/mlib_go/pkg/domain/pmx"

type UiState struct {
	Frame                  int                   // フレーム
	MaxFrame               int                   // 最大フレーム
	EnabledFrameDrop       bool                  // フレームドロップON/OFF
	EnabledPhysics         bool                  // 物理ON/OFF
	DoPhysicsReset         bool                  // 物理リセット
	IsShowNormal           bool                  // ボーンデバッグ表示
	IsShowWire             bool                  // ワイヤーフレームデバッグ表示
	IsShowSelectedVertex   bool                  // 選択頂点デバッグ表示
	IsShowBones            map[pmx.BoneFlag]bool // ボーン表示フラグ
	IsShowRigidBodyFront   bool                  // 剛体デバッグ表示(前面)
	IsShowRigidBodyBack    bool                  // 剛体デバッグ表示(埋め込み)
	IsShowJoint            bool                  // ジョイントデバッグ表示
	IsShowInfo             bool                  // 情報デバッグ表示
	IsLimitFps30           bool                  // 30FPS制限
	IsLimitFps60           bool                  // 60FPS制限
	IsUnLimitFps           bool                  // FPS無制限
	IsUnLimitFpsDeform     bool                  // デフォームFPS無制限
	IsLogLevelDebug        bool                  // デバッグメッセージ表示
	IsLogLevelVerbose      bool                  // 冗長メッセージ表示
	IsLogLevelIkVerbose    bool                  // IK冗長メッセージ表示
	SpfLimit               float64               //fps制限
	LeftButtonPressed      bool                  // 左ボタン押下フラグ
	MiddleButtonPressed    bool                  // 中ボタン押下フラグ
	RightButtonPressed     bool                  // 右ボタン押下フラグ
	UpdatedPrev            bool                  // 前回のカーソル位置更新フラグ
	ShiftPressed           bool                  // Shiftキー押下フラグ
	CtrlPressed            bool                  // Ctrlキー押下フラグ
	IsGlRunning            bool                  // 描画ループ中フラグ
	IsWalkRunning          bool                  // walkウィンドウが閉じられたかどうか
	DoResetPhysicsStart    bool                  // 物理リセット開始フラグ
	DoResetPhysicsProgress bool                  // 物理リセット中フラグ
	DoResetPhysicsCount    int                   // 物理リセット処理回数
	IsSaveDelta            bool                  // 前回デフォーム保存フラグ(walkウィンドウからの変更情報検知用)
	Playing                bool                  // 再生中フラグ
}

func NewUiState() *UiState {
	return &UiState{
		IsShowBones:         make(map[pmx.BoneFlag]bool),
		SpfLimit:            1.0 / 30.0,
		DoResetPhysicsCount: 0,
		IsGlRunning:         true,
		IsWalkRunning:       true,
	}
}

func (u *UiState) TriggerPlay(p bool) {
	u.Playing = p
	u.IsSaveDelta = false
}

func (u *UiState) TriggerPhysicsEnabled(enabled bool) {
	u.EnabledPhysics = enabled
	u.IsSaveDelta = false
}

func (u *UiState) TriggerPhysicsReset() {
	if !u.DoResetPhysicsProgress {
		u.DoResetPhysicsStart = true
		u.IsSaveDelta = false
	}
}
