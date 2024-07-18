package controller

import (
	"github.com/miu200521358/mlib_go/pkg/domain/state"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
)

type controlState struct {
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

func NewControlState() *controlState {
	u := &controlState{
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

func (c *controlState) SetAnimationState(state state.IAnimationState) {
	c.animationState <- state.(*renderer.AnimationState)
}

func (c *controlState) AnimationState() *renderer.AnimationState {
	return <-c.animationState
}

func (c *controlState) Frame() float64 {
	return <-c.frameChan
}

func (c *controlState) SetFrame(frame float64) {
	c.frameChan <- frame
}

func (c *controlState) AddFrame(v float64) {
	f := c.Frame()
	c.frameChan <- v + f
}

func (c *controlState) MaxFrame() int {
	return <-c.maxFrameChan
}

func (c *controlState) SetMaxFrame(maxFrame int) {
	c.maxFrameChan <- maxFrame
}

func (c *controlState) PrevFrame() int {
	return <-c.prevFrameChan
}

func (c *controlState) SetPrevFrame(prevFrame int) {
	c.prevFrameChan <- prevFrame
}

func (c *controlState) IsEnabledFrameDrop() bool {
	return <-c.isEnabledFrameDropChan
}

func (c *controlState) SetEnabledFrameDrop(enabled bool) {
	c.isEnabledFrameDropChan <- enabled
}

func (c *controlState) IsEnabledPhysics() bool {
	return <-c.isEnabledPhysicsChan
}

func (c *controlState) SetEnabledPhysics(enabled bool) {
	c.isEnabledPhysicsChan <- enabled
}

func (c *controlState) IsPhysicsReset() bool {
	return <-c.isPhysicsResetChan
}

func (c *controlState) SetPhysicsReset(reset bool) {
	c.physicsResetChan <- reset
}

func (c *controlState) IsShowNormal() bool {
	return <-c.isShowNormalChan
}

func (c *controlState) SetShowNormal(show bool) {
	c.isShowNormalChan <- show
}

func (c *controlState) IsShowWire() bool {
	return <-c.isShowWireChan
}

func (c *controlState) SetShowWire(show bool) {
	c.isShowWireChan <- show
}

func (c *controlState) IsShowSelectedVertex() bool {
	return <-c.isShowSelectedVertexChan
}

func (c *controlState) SetShowSelectedVertex(show bool) {
	c.isShowSelectedVertexChan <- show
}

func (c *controlState) IsShowBoneAll() bool {
	return <-c.isShowBoneAllChan
}

func (c *controlState) SetShowBoneAll(show bool) {
	c.isShowBoneAllChan <- show
}

func (c *controlState) IsShowBoneIk() bool {
	return <-c.isShowBoneIkChan
}

func (c *controlState) SetShowBoneIk(show bool) {
	c.isShowBoneIkChan <- show
}

func (c *controlState) IsShowBoneEffector() bool {
	return <-c.isShowBoneEffectorChan
}

func (c *controlState) SetShowBoneEffector(show bool) {
	c.isShowBoneEffectorChan <- show
}

func (c *controlState) IsShowBoneFixed() bool {
	return <-c.isShowBoneFixedChan
}

func (c *controlState) SetShowBoneFixed(show bool) {
	c.isShowBoneFixedChan <- show
}

func (c *controlState) IsShowBoneRotate() bool {
	return <-c.isShowBoneRotateChan
}

func (c *controlState) SetShowBoneRotate(show bool) {
	c.isShowBoneRotateChan <- show
}

func (c *controlState) IsShowBoneTranslate() bool {
	return <-c.isShowBoneTranslateChan
}

func (c *controlState) SetShowBoneTranslate(show bool) {
	c.isShowBoneTranslateChan <- show
}

func (c *controlState) IsShowBoneVisible() bool {
	return <-c.isShowBoneVisibleChan
}

func (c *controlState) SetShowBoneVisible(show bool) {
	c.isShowBoneVisibleChan <- show
}

func (c *controlState) IsShowRigidBodyFront() bool {
	return <-c.isShowRigidBodyFrontChan
}

func (c *controlState) SetShowRigidBodyFront(show bool) {
	c.isShowRigidBodyFrontChan <- show
}

func (c *controlState) IsShowRigidBodyBack() bool {
	return <-c.isShowRigidBodyBackChan
}

func (c *controlState) SetShowRigidBodyBack(show bool) {
	c.isShowRigidBodyBackChan <- show
}

func (c *controlState) IsShowJoint() bool {
	return <-c.isShowJointChan
}

func (c *controlState) SetShowJoint(show bool) {
	c.isShowJointChan <- show
}

func (c *controlState) IsShowInfo() bool {
	return <-c.isShowInfoChan
}

func (c *controlState) SetShowInfo(show bool) {
	c.isShowInfoChan <- show
}

func (c *controlState) IsLimitFps30() bool {
	return <-c.isLimitFps30Chan
}

func (c *controlState) SetLimitFps30(limit bool) {
	c.isLimitFps30Chan <- limit
}

func (c *controlState) IsLimitFps60() bool {
	return <-c.isLimitFps60Chan
}

func (c *controlState) SetLimitFps60(limit bool) {
	c.isLimitFps60Chan <- limit
}

func (c *controlState) IsUnLimitFps() bool {
	return <-c.isUnLimitFpsChan
}

func (c *controlState) SetUnLimitFps(limit bool) {
	c.isUnLimitFpsChan <- limit
}

func (c *controlState) IsUnLimitFpsDeform() bool {
	return <-c.isUnLimitFpsDeformChan
}

func (c *controlState) SetUnLimitFpsDeform(limit bool) {
	c.isUnLimitFpsDeformChan <- limit
}

func (c *controlState) IsLogLevelDebug() bool {
	return <-c.isLogLevelDebugChan
}

func (c *controlState) SetLogLevelDebug(log bool) {
	c.isLogLevelDebugChan <- log
}

func (c *controlState) IsLogLevelVerbose() bool {
	return <-c.isLogLevelVerboseChan
}

func (c *controlState) SetLogLevelVerbose(log bool) {
	c.isLogLevelVerboseChan <- log
}

func (c *controlState) IsLogLevelIkVerbose() bool {
	return <-c.isLogLevelIkVerboseChan
}

func (c *controlState) SetLogLevelIkVerbose(log bool) {
	c.isLogLevelIkVerboseChan <- log
}

func (c *controlState) IsClosed() bool {
	return <-c.isClosedChan
}

func (c *controlState) SetClosed(closed bool) {
	c.isClosedChan <- closed
}

func (c *controlState) Playing() bool {
	return <-c.playingChan
}

func (c *controlState) TriggerPlay(p bool) {
	c.playingChan <- p
}

func (c *controlState) SpfLimit() float64 {
	return <-c.spfLimitChan
}

func (c *controlState) SetSpfLimit(spf float64) {
	c.spfLimitChan <- spf
}
