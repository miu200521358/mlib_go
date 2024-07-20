package controller

import (
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
)

type ControlState struct {
	appState                 state.IAppState               // アプリ状態
	motionPlayer             state.IPlayer                 // モーションプレイヤー
	controlWindow            state.IControlWindow          // コントロールウィンドウ
	prevFrameChan            chan int                      // 前回フレーム
	frameChan                chan float64                  // フレーム
	maxFrameChan             chan int                      // 最大フレーム
	isEnabledFrameDropChan   chan bool                     // フレームドロップON/OFF
	isEnabledPhysicsChan     chan bool                     // 物理ON/OFF
	isPhysicsResetChan       chan bool                     // 物理リセット
	isShowNormalChan         chan bool                     // ボーンデバッグ表示
	isShowWireChan           chan bool                     // ワイヤーフレームデバッグ表示
	isShowSelectedVertexChan chan bool                     // 選択頂点デバッグ表示
	isShowBoneAllChan        chan bool                     // 全ボーンデバッグ表示
	isShowBoneIkChan         chan bool                     // IKボーンデバッグ表示
	isShowBoneEffectorChan   chan bool                     // 付与親ボーンデバッグ表示
	isShowBoneFixedChan      chan bool                     // 軸制限ボーンデバッグ表示
	isShowBoneRotateChan     chan bool                     // 回転ボーンデバッグ表示
	isShowBoneTranslateChan  chan bool                     // 移動ボーンデバッグ表示
	isShowBoneVisibleChan    chan bool                     // 表示ボーンデバッグ表示
	isShowRigidBodyFrontChan chan bool                     // 剛体デバッグ表示(前面)
	isShowRigidBodyBackChan  chan bool                     // 剛体デバッグ表示(埋め込み)
	isShowJointChan          chan bool                     // ジョイントデバッグ表示
	isShowInfoChan           chan bool                     // 情報デバッグ表示
	isLimitFps30Chan         chan bool                     // 30FPS制限
	isLimitFps60Chan         chan bool                     // 60FPS制限
	isUnLimitFpsChan         chan bool                     // FPS無制限
	isUnLimitFpsDeformChan   chan bool                     // デフォームFPS無制限
	isLogLevelDebugChan      chan bool                     // デバッグメッセージ表示
	isLogLevelVerboseChan    chan bool                     // 冗長メッセージ表示
	isLogLevelIkVerboseChan  chan bool                     // IK冗長メッセージ表示
	isClosedChan             chan bool                     // ウィンドウクローズ
	playingChan              chan bool                     // 再生中フラグ
	physicsResetChan         chan bool                     // 物理リセット
	spfLimitChan             chan float64                  // FPS制限
	animationState           chan *renderer.AnimationState // アニメーションステート
}

func NewControlState(appState state.IAppState) *ControlState {
	u := &ControlState{
		appState:                 appState,
		prevFrameChan:            make(chan int),
		frameChan:                make(chan float64),
		maxFrameChan:             make(chan int),
		isEnabledFrameDropChan:   make(chan bool),
		isEnabledPhysicsChan:     make(chan bool),
		isPhysicsResetChan:       make(chan bool),
		isShowNormalChan:         make(chan bool),
		isShowWireChan:           make(chan bool),
		isShowSelectedVertexChan: make(chan bool),
		isShowBoneAllChan:        make(chan bool),
		isShowBoneIkChan:         make(chan bool),
		isShowBoneEffectorChan:   make(chan bool),
		isShowBoneFixedChan:      make(chan bool),
		isShowBoneRotateChan:     make(chan bool),
		isShowBoneTranslateChan:  make(chan bool),
		isShowBoneVisibleChan:    make(chan bool),
		isShowRigidBodyFrontChan: make(chan bool),
		isShowRigidBodyBackChan:  make(chan bool),
		isShowJointChan:          make(chan bool),
		isShowInfoChan:           make(chan bool),
		isLimitFps30Chan:         make(chan bool),
		isLimitFps60Chan:         make(chan bool),
		isUnLimitFpsChan:         make(chan bool),
		isUnLimitFpsDeformChan:   make(chan bool),
		isLogLevelDebugChan:      make(chan bool),
		isLogLevelVerboseChan:    make(chan bool),
		isLogLevelIkVerboseChan:  make(chan bool),
		isClosedChan:             make(chan bool),
		playingChan:              make(chan bool),
		physicsResetChan:         make(chan bool),
		spfLimitChan:             make(chan float64),
		animationState:           make(chan *renderer.AnimationState),
	}

	return u
}

func (s *ControlState) Run() {
	go func() {
		for {
			select {
			case prevFrame := <-s.prevFrameChan:
				s.appState.SetPrevFrame(prevFrame)
			case frame := <-s.frameChan:
				s.appState.SetFrame(frame)
			case maxFrame := <-s.maxFrameChan:
				s.appState.UpdateMaxFrame(maxFrame)
			case enabledFrameDrop := <-s.isEnabledFrameDropChan:
				s.appState.SetEnabledFrameDrop(enabledFrameDrop)
			case enabledPhysics := <-s.isEnabledPhysicsChan:
				s.appState.SetEnabledPhysics(enabledPhysics)
			case resetPhysics := <-s.physicsResetChan:
				s.appState.SetPhysicsReset(resetPhysics)
			case showNormal := <-s.isShowNormalChan:
				s.appState.SetShowNormal(showNormal)
			case showWire := <-s.isShowWireChan:
				s.appState.SetShowWire(showWire)
			case showSelectedVertex := <-s.isShowSelectedVertexChan:
				s.appState.SetShowSelectedVertex(showSelectedVertex)
			case showBoneAll := <-s.isShowBoneAllChan:
				s.appState.SetShowBoneAll(showBoneAll)
			case showBoneIk := <-s.isShowBoneIkChan:
				s.appState.SetShowBoneIk(showBoneIk)
			case showBoneEffector := <-s.isShowBoneEffectorChan:
				s.appState.SetShowBoneEffector(showBoneEffector)
			case showBoneFixed := <-s.isShowBoneFixedChan:
				s.appState.SetShowBoneFixed(showBoneFixed)
			case showBoneRotate := <-s.isShowBoneRotateChan:
				s.appState.SetShowBoneRotate(showBoneRotate)
			case showBoneTranslate := <-s.isShowBoneTranslateChan:
				s.appState.SetShowBoneTranslate(showBoneTranslate)
			case showBoneVisible := <-s.isShowBoneVisibleChan:
				s.appState.SetShowBoneVisible(showBoneVisible)
			case showRigidBodyFront := <-s.isShowRigidBodyFrontChan:
				s.appState.SetShowRigidBodyFront(showRigidBodyFront)
			case showRigidBodyBack := <-s.isShowRigidBodyBackChan:
				s.appState.SetShowRigidBodyBack(showRigidBodyBack)
			case showJoint := <-s.isShowJointChan:
				s.appState.SetShowJoint(showJoint)
			case showInfo := <-s.isShowInfoChan:
				s.appState.SetShowInfo(showInfo)
			case spfLimit := <-s.spfLimitChan:
				s.appState.SetSpfLimit(spfLimit)
			case closed := <-s.isClosedChan:
				s.appState.SetClosed(closed)
			case playing := <-s.playingChan:
				s.appState.TriggerPlay(playing)
			case spfLimit := <-s.spfLimitChan:
				s.appState.SetSpfLimit(spfLimit)
			case animationState := <-s.animationState:
				s.appState.SetAnimationState(animationState)
			}
		}
	}()
}

func (c *ControlState) SetPlayer(mp state.IPlayer) {
	c.motionPlayer = mp
}

func (c *ControlState) SetControlWindow(cw state.IControlWindow) {
	c.controlWindow = cw
}

func (c *ControlState) SetAnimationState(state state.IAnimationState) {
	c.animationState <- state.(*renderer.AnimationState)
}

func (c *ControlState) AnimationState() *renderer.AnimationState {
	return <-c.animationState
}

func (c *ControlState) Frame() float64 {
	return c.motionPlayer.Frame()
}

func (c *ControlState) SetFrame(frame float64) {
	c.motionPlayer.SetFrame(frame)
	c.frameChan <- frame
}

func (c *ControlState) AddFrame(v float64) {
	f := c.Frame()
	c.SetFrame(f + v)
}

func (c *ControlState) MaxFrame() int {
	return c.motionPlayer.MaxFrame()
}

func (c *ControlState) SetMaxFrame(maxFrame int) {
	c.motionPlayer.SetMaxFrame(maxFrame)
	c.maxFrameChan <- maxFrame
}

func (c *ControlState) UpdateMaxFrame(maxFrame int) {
	c.motionPlayer.UpdateMaxFrame(maxFrame)
	c.maxFrameChan <- maxFrame
}

func (c *ControlState) PrevFrame() int {
	return c.motionPlayer.PrevFrame()
}

func (c *ControlState) SetPrevFrame(prevFrame int) {
	c.motionPlayer.SetPrevFrame(prevFrame)
	c.prevFrameChan <- prevFrame
}

func (c *ControlState) IsEnabledFrameDrop() bool {
	return c.controlWindow.IsEnabledFrameDrop()
}

func (c *ControlState) SetEnabledFrameDrop(enabled bool) {
	c.isEnabledFrameDropChan <- enabled
}

func (c *ControlState) IsEnabledPhysics() bool {
	return <-c.isEnabledPhysicsChan
}

func (c *ControlState) SetEnabledPhysics(enabled bool) {
	c.isEnabledPhysicsChan <- enabled
}

func (c *ControlState) IsPhysicsReset() bool {
	return <-c.isPhysicsResetChan
}

func (c *ControlState) SetPhysicsReset(reset bool) {
	c.physicsResetChan <- reset
}

func (c *ControlState) IsShowNormal() bool {
	return <-c.isShowNormalChan
}

func (c *ControlState) SetShowNormal(show bool) {
	c.isShowNormalChan <- show
}

func (c *ControlState) IsShowWire() bool {
	return <-c.isShowWireChan
}

func (c *ControlState) SetShowWire(show bool) {
	c.isShowWireChan <- show
}

func (c *ControlState) IsShowSelectedVertex() bool {
	return <-c.isShowSelectedVertexChan
}

func (c *ControlState) SetShowSelectedVertex(show bool) {
	c.isShowSelectedVertexChan <- show
}

func (c *ControlState) IsShowBoneAll() bool {
	return <-c.isShowBoneAllChan
}

func (c *ControlState) SetShowBoneAll(show bool) {
	c.isShowBoneAllChan <- show
}

func (c *ControlState) IsShowBoneIk() bool {
	return <-c.isShowBoneIkChan
}

func (c *ControlState) SetShowBoneIk(show bool) {
	c.isShowBoneIkChan <- show
}

func (c *ControlState) IsShowBoneEffector() bool {
	return <-c.isShowBoneEffectorChan
}

func (c *ControlState) SetShowBoneEffector(show bool) {
	c.isShowBoneEffectorChan <- show
}

func (c *ControlState) IsShowBoneFixed() bool {
	return <-c.isShowBoneFixedChan
}

func (c *ControlState) SetShowBoneFixed(show bool) {
	c.isShowBoneFixedChan <- show
}

func (c *ControlState) IsShowBoneRotate() bool {
	return <-c.isShowBoneRotateChan
}

func (c *ControlState) SetShowBoneRotate(show bool) {
	c.isShowBoneRotateChan <- show
}

func (c *ControlState) IsShowBoneTranslate() bool {
	return <-c.isShowBoneTranslateChan
}

func (c *ControlState) SetShowBoneTranslate(show bool) {
	c.isShowBoneTranslateChan <- show
}

func (c *ControlState) IsShowBoneVisible() bool {
	return <-c.isShowBoneVisibleChan
}

func (c *ControlState) SetShowBoneVisible(show bool) {
	c.isShowBoneVisibleChan <- show
}

func (c *ControlState) IsShowRigidBodyFront() bool {
	return <-c.isShowRigidBodyFrontChan
}

func (c *ControlState) SetShowRigidBodyFront(show bool) {
	c.isShowRigidBodyFrontChan <- show
}

func (c *ControlState) IsShowRigidBodyBack() bool {
	return <-c.isShowRigidBodyBackChan
}

func (c *ControlState) SetShowRigidBodyBack(show bool) {
	c.isShowRigidBodyBackChan <- show
}

func (c *ControlState) IsShowJoint() bool {
	return <-c.isShowJointChan
}

func (c *ControlState) SetShowJoint(show bool) {
	c.isShowJointChan <- show
}

func (c *ControlState) IsShowInfo() bool {
	return <-c.isShowInfoChan
}

func (c *ControlState) SetShowInfo(show bool) {
	c.isShowInfoChan <- show
}

func (c *ControlState) IsLimitFps30() bool {
	return <-c.isLimitFps30Chan
}

func (c *ControlState) SetLimitFps30(limit bool) {
	c.isLimitFps30Chan <- limit
}

func (c *ControlState) IsLimitFps60() bool {
	return <-c.isLimitFps60Chan
}

func (c *ControlState) SetLimitFps60(limit bool) {
	c.isLimitFps60Chan <- limit
}

func (c *ControlState) IsUnLimitFps() bool {
	return <-c.isUnLimitFpsChan
}

func (c *ControlState) SetUnLimitFps(limit bool) {
	c.isUnLimitFpsChan <- limit
}

func (c *ControlState) IsUnLimitFpsDeform() bool {
	return <-c.isUnLimitFpsDeformChan
}

func (c *ControlState) SetUnLimitFpsDeform(limit bool) {
	c.isUnLimitFpsDeformChan <- limit
}

func (c *ControlState) IsLogLevelDebug() bool {
	return <-c.isLogLevelDebugChan
}

func (c *ControlState) SetLogLevelDebug(log bool) {
	c.isLogLevelDebugChan <- log
}

func (c *ControlState) IsLogLevelVerbose() bool {
	return <-c.isLogLevelVerboseChan
}

func (c *ControlState) SetLogLevelVerbose(log bool) {
	c.isLogLevelVerboseChan <- log
}

func (c *ControlState) IsLogLevelIkVerbose() bool {
	return <-c.isLogLevelIkVerboseChan
}

func (c *ControlState) SetLogLevelIkVerbose(log bool) {
	c.isLogLevelIkVerboseChan <- log
}

func (c *ControlState) IsClosed() bool {
	return <-c.isClosedChan
}

func (c *ControlState) SetClosed(closed bool) {
	c.isClosedChan <- closed
}

func (c *ControlState) Playing() bool {
	return <-c.playingChan
}

func (c *ControlState) TriggerPlay(p bool) {
	c.playingChan <- p
}

func (c *ControlState) SpfLimit() float64 {
	return <-c.spfLimitChan
}

func (c *ControlState) SetSpfLimit(spf float64) {
	c.spfLimitChan <- spf
}
