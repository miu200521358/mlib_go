package app

type appState struct {
	frame                float64 // フレーム
	prevFrame            int     // 前回のフレーム
	maxFrame             int     // 最大フレーム
	isEnabledFrameDrop   bool    // フレームドロップON/OFF
	isEnabledPhysics     bool    // 物理ON/OFF
	isPhysicsReset       bool    // 物理リセット
	isShowNormal         bool    // ボーンデバッグ表示
	isShowWire           bool    // ワイヤーフレームデバッグ表示
	isShowSelectedVertex bool    // 選択頂点デバッグ表示
	isShowBoneAll        bool    // 全ボーンデバッグ表示
	isShowBoneIk         bool    // IKボーンデバッグ表示
	isShowBoneEffector   bool    // 付与親ボーンデバッグ表示
	isShowBoneFixed      bool    // 軸制限ボーンデバッグ表示
	isShowBoneRotate     bool    // 回転ボーンデバッグ表示
	isShowBoneTranslate  bool    // 移動ボーンデバッグ表示
	isShowBoneVisible    bool    // 表示ボーンデバッグ表示
	isShowRigidBodyFront bool    // 剛体デバッグ表示(前面)
	isShowRigidBodyBack  bool    // 剛体デバッグ表示(埋め込み)
	isShowJoint          bool    // ジョイントデバッグ表示
	isShowInfo           bool    // 情報デバッグ表示
	isLimitFps30         bool    // 30FPS制限
	isLimitFps60         bool    // 60FPS制限
	isUnLimitFps         bool    // FPS無制限
	isUnLimitFpsDeform   bool    // デフォームFPS無制限
	isLogLevelDebug      bool    // デバッグメッセージ表示
	isLogLevelVerbose    bool    // 冗長メッセージ表示
	isLogLevelIkVerbose  bool    // IK冗長メッセージ表示
	isClosed             bool    // ウィンドウクローズ
	playing              bool    // 再生中フラグ
	spfLimit             float64 // FPS制限
}

func newAppState() *appState {
	u := &appState{
		isEnabledPhysics:   true, // 最初は物理ON
		isEnabledFrameDrop: true, // 最初はフレームドロップON
		isLimitFps30:       true, // 最初は30fps制限
	}

	return u
}

func (u *appState) Frame() float64 {
	return u.frame
}

func (u *appState) SetFrame(frame float64) {
	u.frame = frame
}

func (u *appState) ChangeFrame(frame float64) {
	u.frame = frame
}

func (u *appState) AddFrame(v float64) {
	u.frame += v
}

func (u *appState) MaxFrame() int {
	return u.maxFrame
}

func (u *appState) SetMaxFrame(maxFrame int) {
	u.maxFrame = maxFrame
}

func (u *appState) PrevFrame() int {
	return u.prevFrame
}

func (u *appState) SetPrevFrame(prevFrame int) {
	u.prevFrame = prevFrame
}

func (u *appState) IsEnabledFrameDrop() bool {
	return u.isEnabledFrameDrop
}

func (u *appState) SetEnabledFrameDrop(enabled bool) {
	u.isEnabledFrameDrop = enabled
}

func (u *appState) IsEnabledPhysics() bool {
	return u.isEnabledPhysics
}

func (u *appState) SetEnabledPhysics(enabled bool) {
	u.isEnabledPhysics = enabled
}

func (u *appState) IsPhysicsReset() bool {
	return u.isPhysicsReset
}

func (u *appState) SetPhysicsReset(reset bool) {
	u.isPhysicsReset = reset
}

func (u *appState) IsShowNormal() bool {
	return u.isShowNormal
}

func (u *appState) SetShowNormal(show bool) {
	u.isShowNormal = show
}

func (u *appState) IsShowWire() bool {
	return u.isShowWire
}

func (u *appState) SetShowWire(show bool) {
	u.isShowWire = show
}

func (u *appState) IsShowSelectedVertex() bool {
	return u.isShowSelectedVertex
}

func (u *appState) SetShowSelectedVertex(show bool) {
	u.isShowSelectedVertex = show
}

func (u *appState) IsShowBoneAll() bool {
	return u.isShowBoneAll
}

func (u *appState) SetShowBoneAll(show bool) {
	u.isShowBoneAll = show
}

func (u *appState) IsShowBoneIk() bool {
	return u.isShowBoneIk
}

func (u *appState) SetShowBoneIk(show bool) {
	u.isShowBoneIk = show
}

func (u *appState) IsShowBoneEffector() bool {
	return u.isShowBoneEffector
}

func (u *appState) SetShowBoneEffector(show bool) {
	u.isShowBoneEffector = show
}

func (u *appState) IsShowBoneFixed() bool {
	return u.isShowBoneFixed
}

func (u *appState) SetShowBoneFixed(show bool) {
	u.isShowBoneFixed = show
}

func (u *appState) IsShowBoneRotate() bool {
	return u.isShowBoneRotate
}

func (u *appState) SetShowBoneRotate(show bool) {
	u.isShowBoneRotate = show
}

func (u *appState) IsShowBoneTranslate() bool {
	return u.isShowBoneTranslate
}

func (u *appState) SetShowBoneTranslate(show bool) {
	u.isShowBoneTranslate = show
}

func (u *appState) IsShowBoneVisible() bool {
	return u.isShowBoneVisible
}

func (u *appState) SetShowBoneVisible(show bool) {
	u.isShowBoneVisible = show
}

func (u *appState) IsShowRigidBodyFront() bool {
	return u.isShowRigidBodyFront
}

func (u *appState) SetShowRigidBodyFront(show bool) {
	u.isShowRigidBodyFront = show
}

func (u *appState) IsShowRigidBodyBack() bool {
	return u.isShowRigidBodyBack
}

func (u *appState) SetShowRigidBodyBack(show bool) {
	u.isShowRigidBodyBack = show
}

func (u *appState) IsShowJoint() bool {
	return u.isShowJoint
}

func (u *appState) SetShowJoint(show bool) {
	u.isShowJoint = show
}

func (u *appState) IsShowInfo() bool {
	return u.isShowInfo
}

func (u *appState) SetShowInfo(show bool) {
	u.isShowInfo = show
}

func (u *appState) IsLimitFps30() bool {
	return u.isLimitFps30
}

func (u *appState) SetLimitFps30(limit bool) {
	u.isLimitFps30 = limit
}

func (u *appState) IsLimitFps60() bool {
	return u.isLimitFps60
}

func (u *appState) SetLimitFps60(limit bool) {
	u.isLimitFps60 = limit
}

func (u *appState) IsUnLimitFps() bool {
	return u.isUnLimitFps
}

func (u *appState) SetUnLimitFps(limit bool) {
	u.isUnLimitFps = limit
}

func (u *appState) IsUnLimitFpsDeform() bool {
	return u.isUnLimitFpsDeform
}

func (u *appState) SetUnLimitFpsDeform(limit bool) {
	u.isUnLimitFpsDeform = limit
}

func (u *appState) IsLogLevelDebug() bool {
	return u.isLogLevelDebug
}

func (u *appState) SetLogLevelDebug(log bool) {
	u.isLogLevelDebug = log
}

func (u *appState) IsLogLevelVerbose() bool {
	return u.isLogLevelVerbose
}

func (u *appState) SetLogLevelVerbose(log bool) {
	u.isLogLevelVerbose = log
}

func (u *appState) IsLogLevelIkVerbose() bool {
	return u.isLogLevelIkVerbose
}

func (u *appState) SetLogLevelIkVerbose(log bool) {
	u.isLogLevelIkVerbose = log
}

func (u *appState) IsClosed() bool {
	return u.isClosed
}

func (u *appState) SetClosed(closed bool) {
	u.isClosed = closed
}

func (u *appState) Playing() bool {
	return u.playing
}

func (u *appState) TriggerPlay(p bool) {
	u.playing = p
}

func (u *appState) SpfLimit() float64 {
	return u.spfLimit
}

func (u *appState) SetSpfLimit(spf float64) {
	u.spfLimit = spf
}
