//go:build windows
// +build windows

package app

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type appState struct {
	frame                        float32                   // フレーム
	maxFrame                     float32                   // 最大フレーム
	isEnabledFrameDrop           bool                      // フレームドロップON/OFF
	isEnabledPhysics             bool                      // 物理ON/OFF
	isPhysicsReset               bool                      // 物理リセット
	isShowNormal                 bool                      // ボーンデバッグ表示
	isShowWire                   bool                      // ワイヤーフレームデバッグ表示
	isShowOverride               bool                      // オーバーライドデバッグ表示
	isShowSelectedVertex         bool                      // 選択頂点デバッグ表示
	isShowBoneAll                bool                      // 全ボーンデバッグ表示
	isShowBoneIk                 bool                      // IKボーンデバッグ表示
	isShowBoneEffector           bool                      // 付与親ボーンデバッグ表示
	isShowBoneFixed              bool                      // 軸制限ボーンデバッグ表示
	isShowBoneRotate             bool                      // 回転ボーンデバッグ表示
	isShowBoneTranslate          bool                      // 移動ボーンデバッグ表示
	isShowBoneVisible            bool                      // 表示ボーンデバッグ表示
	isShowRigidBodyFront         bool                      // 剛体デバッグ表示(前面)
	isShowRigidBodyBack          bool                      // 剛体デバッグ表示(埋め込み)
	isShowJoint                  bool                      // ジョイントデバッグ表示
	isShowInfo                   bool                      // 情報デバッグ表示
	isLimitFps30                 bool                      // 30FPS制限
	isLimitFps60                 bool                      // 60FPS制限
	isUnLimitFps                 bool                      // FPS無制限
	isUnLimitFpsDeform           bool                      // デフォームFPS無制限
	isCameraSync                 bool                      // レンダーシンク
	isClosed                     bool                      // ウィンドウクローズ
	playing                      bool                      // 再生中フラグ
	spfLimit                     float64                   // FPS制限
	frameChannel                 chan float32              // フレーム
	maxFrameChannel              chan float32              // 最大フレーム
	isEnabledFrameDropChannel    chan bool                 // フレームドロップON/OFF
	isEnabledPhysicsChannel      chan bool                 // 物理ON/OFF
	isPhysicsResetChannel        chan bool                 // 物理リセット
	isShowNormalChannel          chan bool                 // ボーンデバッグ表示
	isShowWireChannel            chan bool                 // ワイヤーフレームデバッグ表示
	isShowOverrideChannel        chan bool                 // オーバーライドデバッグ表示
	isShowSelectedVertexChannel  chan bool                 // 選択頂点デバッグ表示
	isShowBoneAllChannel         chan bool                 // 全ボーンデバッグ表示
	isShowBoneIkChannel          chan bool                 // IKボーンデバッグ表示
	isShowBoneEffectorChannel    chan bool                 // 付与親ボーンデバッグ表示
	isShowBoneFixedChannel       chan bool                 // 軸制限ボーンデバッグ表示
	isShowBoneRotateChannel      chan bool                 // 回転ボーンデバッグ表示
	isShowBoneTranslateChannel   chan bool                 // 移動ボーンデバッグ表示
	isShowBoneVisibleChannel     chan bool                 // 表示ボーンデバッグ表示
	isShowRigidBodyFrontChannel  chan bool                 // 剛体デバッグ表示(前面)
	isShowRigidBodyBackChannel   chan bool                 // 剛体デバッグ表示(埋め込み)
	isShowJointChannel           chan bool                 // ジョイントデバッグ表示
	isShowInfoChannel            chan bool                 // 情報デバッグ表示
	isLimitFps30Channel          chan bool                 // 30FPS制限
	isLimitFps60Channel          chan bool                 // 60FPS制限
	isUnLimitFpsChannel          chan bool                 // FPS無制限
	isUnLimitFpsDeformChannel    chan bool                 // デフォームFPS無制限
	isCameraSyncChannel          chan bool                 // レンダリング同期
	isClosedChannel              chan bool                 // ウィンドウクローズ
	playingChannel               chan bool                 // 再生中フラグ
	physicsResetChannel          chan bool                 // 物理リセット
	spfLimitChanel               chan float64              // FPS制限
	selectedVertexIndexesChannel chan [][][]int            // 選択頂点インデックス
	funcGetModels                func() [][]*pmx.PmxModel  // モデル取得関数
	funcGetMotions               func() [][]*vmd.VmdMotion // モーション取得関数
}

func newAppState() *appState {
	u := &appState{
		isEnabledPhysics:             true,       // 物理ON
		isEnabledFrameDrop:           true,       // フレームドロップON
		isLimitFps30:                 true,       // 30fps制限
		spfLimit:                     1.0 / 30.0, // 30fps
		frame:                        0.0,
		maxFrame:                     1,
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
		physicsResetChannel:          make(chan bool),
		spfLimitChanel:               make(chan float64),
		selectedVertexIndexesChannel: make(chan [][][]int, 1),
	}

	return u
}

func (appState *appState) Frame() float32 {
	return appState.frame
}

func (appState *appState) SetFrame(frame float32) {
	appState.frame = frame
}

func (appState *appState) MaxFrame() float32 {
	return appState.maxFrame
}

func (appState *appState) UpdateMaxFrame(maxFrame float32) {
	if appState.maxFrame < maxFrame {
		appState.maxFrame = maxFrame
	}
}

func (appState *appState) SetMaxFrame(maxFrame float32) {
	appState.maxFrame = maxFrame
}

func (appState *appState) IsEnabledFrameDrop() bool {
	return appState.isEnabledFrameDrop
}

func (appState *appState) SetEnabledFrameDrop(enabled bool) {
	appState.isEnabledFrameDrop = enabled
}

func (appState *appState) IsEnabledPhysics() bool {
	return appState.isEnabledPhysics
}

func (appState *appState) SetEnabledPhysics(enabled bool) {
	appState.SetFrame(appState.frame)
	appState.isEnabledPhysics = enabled
}

func (appState *appState) IsPhysicsReset() bool {
	return appState.isPhysicsReset
}

func (appState *appState) SetPhysicsReset(reset bool) {
	appState.SetFrame(appState.frame)
	appState.isPhysicsReset = reset
}

func (appState *appState) IsShowNormal() bool {
	return appState.isShowNormal
}

func (appState *appState) SetShowNormal(show bool) {
	appState.isShowNormal = show
}

func (appState *appState) IsShowWire() bool {
	return appState.isShowWire
}

func (appState *appState) SetShowWire(show bool) {
	appState.isShowWire = show
}

func (appState *appState) IsShowOverride() bool {
	return appState.isShowOverride
}

func (appState *appState) SetShowOverride(show bool) {
	appState.isShowOverride = show
}

func (appState *appState) IsShowSelectedVertex() bool {
	return appState.isShowSelectedVertex
}

func (appState *appState) SetShowSelectedVertex(show bool) {
	appState.isShowSelectedVertex = show
}

func (appState *appState) IsShowBoneAll() bool {
	return appState.isShowBoneAll
}

func (appState *appState) SetShowBoneAll(show bool) {
	appState.isShowBoneAll = show
}

func (appState *appState) IsShowBoneIk() bool {
	return appState.isShowBoneIk
}

func (appState *appState) SetShowBoneIk(show bool) {
	appState.isShowBoneIk = show
}

func (appState *appState) IsShowBoneEffector() bool {
	return appState.isShowBoneEffector
}

func (appState *appState) SetShowBoneEffector(show bool) {
	appState.isShowBoneEffector = show
}

func (appState *appState) IsShowBoneFixed() bool {
	return appState.isShowBoneFixed
}

func (appState *appState) SetShowBoneFixed(show bool) {
	appState.isShowBoneFixed = show
}

func (appState *appState) IsShowBoneRotate() bool {
	return appState.isShowBoneRotate
}

func (appState *appState) SetShowBoneRotate(show bool) {
	appState.isShowBoneRotate = show
}

func (appState *appState) IsShowBoneTranslate() bool {
	return appState.isShowBoneTranslate
}

func (appState *appState) SetShowBoneTranslate(show bool) {
	appState.isShowBoneTranslate = show
}

func (appState *appState) IsShowBoneVisible() bool {
	return appState.isShowBoneVisible
}

func (appState *appState) SetShowBoneVisible(show bool) {
	appState.isShowBoneVisible = show
}

func (appState *appState) IsShowRigidBodyFront() bool {
	return appState.isShowRigidBodyFront
}

func (appState *appState) SetShowRigidBodyFront(show bool) {
	appState.isShowRigidBodyFront = show
}

func (appState *appState) IsShowRigidBodyBack() bool {
	return appState.isShowRigidBodyBack
}

func (appState *appState) SetShowRigidBodyBack(show bool) {
	appState.isShowRigidBodyBack = show
}

func (appState *appState) IsShowJoint() bool {
	return appState.isShowJoint
}

func (appState *appState) SetShowJoint(show bool) {
	appState.isShowJoint = show
}

func (appState *appState) IsShowInfo() bool {
	return appState.isShowInfo
}

func (appState *appState) SetShowInfo(show bool) {
	appState.isShowInfo = show
}

func (appState *appState) IsLimitFps30() bool {
	return appState.isLimitFps30
}

func (appState *appState) SetLimitFps30(limit bool) {
	appState.isLimitFps30 = limit
}

func (appState *appState) IsLimitFps60() bool {
	return appState.isLimitFps60
}

func (appState *appState) SetLimitFps60(limit bool) {
	appState.isLimitFps60 = limit
}

func (appState *appState) IsUnLimitFps() bool {
	return appState.isUnLimitFps
}

func (appState *appState) SetUnLimitFps(limit bool) {
	appState.isUnLimitFps = limit
}

func (appState *appState) IsUnLimitFpsDeform() bool {
	return appState.isUnLimitFpsDeform
}

func (appState *appState) SetUnLimitFpsDeform(limit bool) {
	appState.isUnLimitFpsDeform = limit
}

func (appState *appState) IsClosed() bool {
	return appState.isClosed
}

func (appState *appState) SetClosed(closed bool) {
	appState.isClosed = closed
}

func (appState *appState) Playing() bool {
	return appState.playing
}

func (appState *appState) SetPlaying(p bool) {
	appState.playing = p
}

func (appState *appState) SpfLimit() float64 {
	return appState.spfLimit
}

func (appState *appState) SetSpfLimit(spf float64) {
	appState.spfLimit = spf
}

func (appState *appState) SetCameraSync(sync bool) {
	appState.isCameraSync = sync
}

func (appState *appState) IsCameraSync() bool {
	return appState.isCameraSync
}

func (appState *appState) SetFrameChannel(v float32) {
	appState.frameChannel <- v
}

func (appState *appState) SetMaxFrameChannel(v float32) {
	appState.maxFrameChannel <- v
}

func (appState *appState) SetEnabledFrameDropChannel(v bool) {
	appState.isEnabledFrameDropChannel <- v
}

func (appState *appState) SetEnabledPhysicsChannel(v bool) {
	appState.isEnabledPhysicsChannel <- v
}

func (appState *appState) SetPhysicsResetChannel(v bool) {
	appState.isPhysicsResetChannel <- v
}

func (appState *appState) SetShowNormalChannel(v bool) {
	appState.isShowNormalChannel <- v
}

func (appState *appState) SetShowWireChannel(v bool) {
	appState.isShowWireChannel <- v
}

func (appState *appState) SetShowOverrideChannel(v bool) {
	appState.isShowOverrideChannel <- v
}

func (appState *appState) SetShowSelectedVertexChannel(v bool) {
	appState.isShowSelectedVertexChannel <- v
}

func (appState *appState) SetShowBoneAllChannel(v bool) {
	appState.isShowBoneAllChannel <- v
}

func (appState *appState) SetShowBoneIkChannel(v bool) {
	appState.isShowBoneIkChannel <- v
}

func (appState *appState) SetShowBoneEffectorChannel(v bool) {
	appState.isShowBoneEffectorChannel <- v
}

func (appState *appState) SetShowBoneFixedChannel(v bool) {
	appState.isShowBoneFixedChannel <- v
}

func (appState *appState) SetShowBoneRotateChannel(v bool) {
	appState.isShowBoneRotateChannel <- v
}

func (appState *appState) SetShowBoneTranslateChannel(v bool) {
	appState.isShowBoneTranslateChannel <- v
}

func (appState *appState) SetShowBoneVisibleChannel(v bool) {
	appState.isShowBoneVisibleChannel <- v
}

func (appState *appState) SetShowRigidBodyFrontChannel(v bool) {
	appState.isShowRigidBodyFrontChannel <- v
}

func (appState *appState) SetShowRigidBodyBackChannel(v bool) {
	appState.isShowRigidBodyBackChannel <- v
}

func (appState *appState) SetShowJointChannel(v bool) {
	appState.isShowJointChannel <- v
}

func (appState *appState) SetShowInfoChannel(v bool) {
	appState.isShowInfoChannel <- v
}

func (appState *appState) SetLimitFps30Channel(v bool) {
	appState.isLimitFps30Channel <- v
}

func (appState *appState) SetLimitFps60Channel(v bool) {
	appState.isLimitFps60Channel <- v
}

func (appState *appState) SetUnLimitFpsChannel(v bool) {
	appState.isUnLimitFpsChannel <- v
}

func (appState *appState) SetUnLimitFpsDeformChannel(v bool) {
	appState.isUnLimitFpsDeformChannel <- v
}

func (appState *appState) SetCameraSyncChannel(v bool) {
	appState.isCameraSyncChannel <- v
}

func (appState *appState) SetClosedChannel(v bool) {
	appState.isClosedChannel <- v
}

func (appState *appState) SetPlayingChannel(v bool) {
	appState.playingChannel <- v
}

func (appState *appState) SetSpfLimitChannel(v float64) {
	appState.spfLimitChanel <- v
}

func (appState *appState) SetSelectedVertexIndexesChannel(v [][][]int) {
	appState.selectedVertexIndexesChannel <- v
}

func (appState *appState) SetGetModels(f func() [][]*pmx.PmxModel) {
	appState.funcGetModels = f
}

func (appState *appState) SetGetMotions(f func() [][]*vmd.VmdMotion) {
	appState.funcGetMotions = f
}

func (appState *appState) GetModels() [][]*pmx.PmxModel {
	return appState.funcGetModels()
}

func (appState *appState) GetMotions() [][]*vmd.VmdMotion {
	return appState.funcGetMotions()
}
