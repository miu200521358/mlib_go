package state

type IControllerState interface {
	IsEnabledFrameDrop() bool
	SetEnabledFrameDrop(enabled bool)
	IsEnabledPhysics() bool
	SetEnabledPhysics(enabled bool)
	IsPhysicsReset() bool
	SetPhysicsReset(reset bool)
	IsShowNormal() bool
	SetShowNormal(show bool)
	IsShowWire() bool
	SetShowWire(show bool)
	IsShowOverride() bool
	SetShowOverride(show bool)
	IsShowSelectedVertex() bool
	SetShowSelectedVertex(show bool)
	IsShowBoneAll() bool
	SetShowBoneAll(show bool)
	IsShowBoneIk() bool
	SetShowBoneIk(show bool)
	IsShowBoneEffector() bool
	SetShowBoneEffector(show bool)
	IsShowBoneFixed() bool
	SetShowBoneFixed(show bool)
	IsShowBoneRotate() bool
	SetShowBoneRotate(show bool)
	IsShowBoneTranslate() bool
	SetShowBoneTranslate(show bool)
	IsShowBoneVisible() bool
	SetShowBoneVisible(show bool)
	IsShowRigidBodyFront() bool
	SetShowRigidBodyFront(show bool)
	IsShowRigidBodyBack() bool
	SetShowRigidBodyBack(show bool)
	IsShowJoint() bool
	SetShowJoint(show bool)
	IsShowInfo() bool
	SetShowInfo(show bool)
	IsLimitFps30() bool
	SetLimitFps30(limit bool)
	IsLimitFps60() bool
	SetLimitFps60(limit bool)
	IsUnLimitFps() bool
	SetUnLimitFps(limit bool)
	IsCameraSync() bool
	SetCameraSync(sync bool)
	IsClosed() bool
	SetClosed(closed bool)
}

type controllerState struct {
	isEnabledFrameDrop   bool // フレームドロップON/OFF
	isEnabledPhysics     bool // 物理ON/OFF
	isPhysicsReset       bool // 物理リセット
	isShowNormal         bool // ボーンデバッグ表示
	isShowWire           bool // ワイヤーフレームデバッグ表示
	isShowOverride       bool // オーバーライドデバッグ表示
	isShowSelectedVertex bool // 選択頂点デバッグ表示
	isShowBoneAll        bool // 全ボーンデバッグ表示
	isShowBoneIk         bool // IKボーンデバッグ表示
	isShowBoneEffector   bool // 付与親ボーンデバッグ表示
	isShowBoneFixed      bool // 軸制限ボーンデバッグ表示
	isShowBoneRotate     bool // 回転ボーンデバッグ表示
	isShowBoneTranslate  bool // 移動ボーンデバッグ表示
	isShowBoneVisible    bool // 表示ボーンデバッグ表示
	isShowRigidBodyFront bool // 剛体デバッグ表示(前面)
	isShowRigidBodyBack  bool // 剛体デバッグ表示(埋め込み)
	isShowJoint          bool // ジョイントデバッグ表示
	isShowInfo           bool // 情報デバッグ表示
	isLimitFps30         bool // 30FPS制限
	isLimitFps60         bool // 60FPS制限
	isUnLimitFps         bool // FPS無制限
	isCameraSync         bool // カメラ同期
	isClosed             bool // 操作ウィンドウクローズ
}

func NewControllerState() IControllerState {
	return &controllerState{
		isEnabledPhysics:   true, // 物理ON
		isEnabledFrameDrop: true, // フレームドロップON
		isLimitFps30:       true, // 30fps制限
	}
}

func (controllerState *controllerState) IsEnabledFrameDrop() bool {
	return controllerState.isEnabledFrameDrop
}

func (controllerState *controllerState) SetEnabledFrameDrop(enabled bool) {
	controllerState.isEnabledFrameDrop = enabled
}

func (controllerState *controllerState) IsEnabledPhysics() bool {
	return controllerState.isEnabledPhysics
}

func (controllerState *controllerState) SetEnabledPhysics(enabled bool) {
	controllerState.isEnabledPhysics = enabled
}

func (controllerState *controllerState) IsPhysicsReset() bool {
	return controllerState.isPhysicsReset
}

func (controllerState *controllerState) SetPhysicsReset(reset bool) {
	controllerState.isPhysicsReset = reset
}

func (controllerState *controllerState) IsShowNormal() bool {
	return controllerState.isShowNormal
}

func (controllerState *controllerState) SetShowNormal(show bool) {
	controllerState.isShowNormal = show
}

func (controllerState *controllerState) IsShowWire() bool {
	return controllerState.isShowWire
}

func (controllerState *controllerState) SetShowWire(show bool) {
	controllerState.isShowWire = show
}

func (controllerState *controllerState) IsShowOverride() bool {
	return controllerState.isShowOverride
}

func (controllerState *controllerState) SetShowOverride(show bool) {
	controllerState.isShowOverride = show
}

func (controllerState *controllerState) IsShowSelectedVertex() bool {
	return controllerState.isShowSelectedVertex
}

func (controllerState *controllerState) SetShowSelectedVertex(show bool) {
	controllerState.isShowSelectedVertex = show
}

func (controllerState *controllerState) IsShowBoneAll() bool {
	return controllerState.isShowBoneAll
}

func (controllerState *controllerState) SetShowBoneAll(show bool) {
	controllerState.isShowBoneAll = show
}

func (controllerState *controllerState) IsShowBoneIk() bool {
	return controllerState.isShowBoneIk
}

func (controllerState *controllerState) SetShowBoneIk(show bool) {
	controllerState.isShowBoneIk = show
}

func (controllerState *controllerState) IsShowBoneEffector() bool {
	return controllerState.isShowBoneEffector
}

func (controllerState *controllerState) SetShowBoneEffector(show bool) {
	controllerState.isShowBoneEffector = show
}

func (controllerState *controllerState) IsShowBoneFixed() bool {
	return controllerState.isShowBoneFixed
}

func (controllerState *controllerState) SetShowBoneFixed(show bool) {
	controllerState.isShowBoneFixed = show
}

func (controllerState *controllerState) IsShowBoneRotate() bool {
	return controllerState.isShowBoneRotate
}

func (controllerState *controllerState) SetShowBoneRotate(show bool) {
	controllerState.isShowBoneRotate = show
}

func (controllerState *controllerState) IsShowBoneTranslate() bool {
	return controllerState.isShowBoneTranslate
}

func (controllerState *controllerState) SetShowBoneTranslate(show bool) {
	controllerState.isShowBoneTranslate = show
}

func (controllerState *controllerState) IsShowBoneVisible() bool {
	return controllerState.isShowBoneVisible
}

func (controllerState *controllerState) SetShowBoneVisible(show bool) {
	controllerState.isShowBoneVisible = show
}

func (controllerState *controllerState) IsShowRigidBodyFront() bool {
	return controllerState.isShowRigidBodyFront
}

func (controllerState *controllerState) SetShowRigidBodyFront(show bool) {
	controllerState.isShowRigidBodyFront = show
}

func (controllerState *controllerState) IsShowRigidBodyBack() bool {
	return controllerState.isShowRigidBodyBack
}

func (controllerState *controllerState) SetShowRigidBodyBack(show bool) {
	controllerState.isShowRigidBodyBack = show
}

func (controllerState *controllerState) IsShowJoint() bool {
	return controllerState.isShowJoint
}

func (controllerState *controllerState) SetShowJoint(show bool) {
	controllerState.isShowJoint = show
}

func (controllerState *controllerState) IsShowInfo() bool {
	return controllerState.isShowInfo
}

func (controllerState *controllerState) SetShowInfo(show bool) {
	controllerState.isShowInfo = show
}

func (controllerState *controllerState) IsLimitFps30() bool {
	return controllerState.isLimitFps30
}

func (controllerState *controllerState) SetLimitFps30(limit bool) {
	controllerState.isLimitFps30 = limit
}

func (controllerState *controllerState) IsLimitFps60() bool {
	return controllerState.isLimitFps60
}

func (controllerState *controllerState) SetLimitFps60(limit bool) {
	controllerState.isLimitFps60 = limit
}

func (controllerState *controllerState) IsUnLimitFps() bool {
	return controllerState.isUnLimitFps
}

func (controllerState *controllerState) SetUnLimitFps(limit bool) {
	controllerState.isUnLimitFps = limit
}

func (controllerState *controllerState) IsCameraSync() bool {
	return controllerState.isCameraSync
}

func (controllerState *controllerState) SetCameraSync(sync bool) {
	controllerState.isCameraSync = sync
}

func (controllerState *controllerState) IsClosed() bool {
	return controllerState.isClosed
}

func (controllerState *controllerState) SetClosed(closed bool) {
	controllerState.isClosed = closed
}
