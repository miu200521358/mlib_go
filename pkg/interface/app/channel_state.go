//go:build windows
// +build windows

package app

type channelState struct {
	frameChannel                 chan float32   // フレーム
	maxFrameChannel              chan float32   // 最大フレーム
	isEnabledFrameDropChannel    chan bool      // フレームドロップON/OFF
	isEnabledPhysicsChannel      chan bool      // 物理ON/OFF
	isPhysicsResetChannel        chan bool      // 物理リセット
	isShowNormalChannel          chan bool      // ボーンデバッグ表示
	isShowWireChannel            chan bool      // ワイヤーフレームデバッグ表示
	isShowOverrideChannel        chan bool      // オーバーライドデバッグ表示
	isShowSelectedVertexChannel  chan bool      // 選択頂点デバッグ表示
	isShowBoneAllChannel         chan bool      // 全ボーンデバッグ表示
	isShowBoneIkChannel          chan bool      // IKボーンデバッグ表示
	isShowBoneEffectorChannel    chan bool      // 付与親ボーンデバッグ表示
	isShowBoneFixedChannel       chan bool      // 軸制限ボーンデバッグ表示
	isShowBoneRotateChannel      chan bool      // 回転ボーンデバッグ表示
	isShowBoneTranslateChannel   chan bool      // 移動ボーンデバッグ表示
	isShowBoneVisibleChannel     chan bool      // 表示ボーンデバッグ表示
	isShowRigidBodyFrontChannel  chan bool      // 剛体デバッグ表示(前面)
	isShowRigidBodyBackChannel   chan bool      // 剛体デバッグ表示(埋め込み)
	isShowJointChannel           chan bool      // ジョイントデバッグ表示
	isShowInfoChannel            chan bool      // 情報デバッグ表示
	isLimitFps30Channel          chan bool      // 30FPS制限
	isLimitFps60Channel          chan bool      // 60FPS制限
	isUnLimitFpsChannel          chan bool      // FPS無制限
	isUnLimitFpsDeformChannel    chan bool      // デフォームFPS無制限
	isCameraSyncChannel          chan bool      // レンダリング同期
	isClosedChannel              chan bool      // ウィンドウクローズ
	playingChannel               chan bool      // 再生中フラグ
	frameIntervalChanel          chan float64   // FPS制限
	selectedVertexIndexesChannel chan [][][]int // 選択頂点インデックス
}

func newChannelState() *channelState {
	return &channelState{
		frameChannel:                 make(chan float32),
		maxFrameChannel:              make(chan float32),
		isEnabledFrameDropChannel:    make(chan bool),
		isEnabledPhysicsChannel:      make(chan bool),
		isPhysicsResetChannel:        make(chan bool),
		isShowNormalChannel:          make(chan bool),
		isShowWireChannel:            make(chan bool),
		isShowOverrideChannel:        make(chan bool),
		isShowSelectedVertexChannel:  make(chan bool),
		isShowBoneAllChannel:         make(chan bool),
		isShowBoneIkChannel:          make(chan bool),
		isShowBoneEffectorChannel:    make(chan bool),
		isShowBoneFixedChannel:       make(chan bool),
		isShowBoneRotateChannel:      make(chan bool),
		isShowBoneTranslateChannel:   make(chan bool),
		isShowBoneVisibleChannel:     make(chan bool),
		isShowRigidBodyFrontChannel:  make(chan bool),
		isShowRigidBodyBackChannel:   make(chan bool),
		isShowJointChannel:           make(chan bool),
		isShowInfoChannel:            make(chan bool),
		isLimitFps30Channel:          make(chan bool),
		isLimitFps60Channel:          make(chan bool),
		isUnLimitFpsChannel:          make(chan bool),
		isUnLimitFpsDeformChannel:    make(chan bool),
		isClosedChannel:              make(chan bool),
		playingChannel:               make(chan bool),
		frameIntervalChanel:          make(chan float64),
		selectedVertexIndexesChannel: make(chan [][][]int, 1),
	}
}

func (channelState *channelState) SetFrameChannel(v float32) {
	channelState.frameChannel <- v
}

func (channelState *channelState) SetMaxFrameChannel(v float32) {
	channelState.maxFrameChannel <- v
}

func (channelState *channelState) SetEnabledFrameDropChannel(v bool) {
	channelState.isEnabledFrameDropChannel <- v
}

func (channelState *channelState) SetEnabledPhysicsChannel(v bool) {
	channelState.isEnabledPhysicsChannel <- v
}

func (channelState *channelState) SetPhysicsResetChannel(v bool) {
	channelState.isPhysicsResetChannel <- v
}

func (channelState *channelState) SetShowNormalChannel(v bool) {
	channelState.isShowNormalChannel <- v
}

func (channelState *channelState) SetShowWireChannel(v bool) {
	channelState.isShowWireChannel <- v
}

func (channelState *channelState) SetShowOverrideChannel(v bool) {
	channelState.isShowOverrideChannel <- v
}

func (channelState *channelState) SetShowSelectedVertexChannel(v bool) {
	channelState.isShowSelectedVertexChannel <- v
}

func (channelState *channelState) SetShowBoneAllChannel(v bool) {
	channelState.isShowBoneAllChannel <- v
}

func (channelState *channelState) SetShowBoneIkChannel(v bool) {
	channelState.isShowBoneIkChannel <- v
}

func (channelState *channelState) SetShowBoneEffectorChannel(v bool) {
	channelState.isShowBoneEffectorChannel <- v
}

func (channelState *channelState) SetShowBoneFixedChannel(v bool) {
	channelState.isShowBoneFixedChannel <- v
}

func (channelState *channelState) SetShowBoneRotateChannel(v bool) {
	channelState.isShowBoneRotateChannel <- v
}

func (channelState *channelState) SetShowBoneTranslateChannel(v bool) {
	channelState.isShowBoneTranslateChannel <- v
}

func (channelState *channelState) SetShowBoneVisibleChannel(v bool) {
	channelState.isShowBoneVisibleChannel <- v
}

func (channelState *channelState) SetShowRigidBodyFrontChannel(v bool) {
	channelState.isShowRigidBodyFrontChannel <- v
}

func (channelState *channelState) SetShowRigidBodyBackChannel(v bool) {
	channelState.isShowRigidBodyBackChannel <- v
}

func (channelState *channelState) SetShowJointChannel(v bool) {
	channelState.isShowJointChannel <- v
}

func (channelState *channelState) SetShowInfoChannel(v bool) {
	channelState.isShowInfoChannel <- v
}

func (channelState *channelState) SetLimitFps30Channel(v bool) {
	channelState.isLimitFps30Channel <- v
}

func (channelState *channelState) SetLimitFps60Channel(v bool) {
	channelState.isLimitFps60Channel <- v
}

func (channelState *channelState) SetUnLimitFpsChannel(v bool) {
	channelState.isUnLimitFpsChannel <- v
}

func (channelState *channelState) SetUnLimitFpsDeformChannel(v bool) {
	channelState.isUnLimitFpsDeformChannel <- v
}

func (channelState *channelState) SetCameraSyncChannel(v bool) {
	channelState.isCameraSyncChannel <- v
}

func (channelState *channelState) SetClosedChannel(v bool) {
	channelState.isClosedChannel <- v
}

func (channelState *channelState) SetPlayingChannel(v bool) {
	channelState.playingChannel <- v
}

func (channelState *channelState) SetFrameIntervalChannel(v float64) {
	channelState.frameIntervalChanel <- v
}

func (channelState *channelState) SetSelectedVertexIndexesChannel(v [][][]int) {
	channelState.selectedVertexIndexesChannel <- v
}
