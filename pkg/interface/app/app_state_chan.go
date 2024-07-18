package app

import (
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
)

type appStateChannels struct {
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
	animationStates          [][]*renderer.AnimationStates // アニメーションステート
}

func newAppStateChannels() *appStateChannels {
	u := &appStateChannels{
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
		animationStates:          make([][]*renderer.AnimationStates, 0),
	}

	return u
}

func (s *appStateChannels) SetAnimationState(state *renderer.AnimationState) {
	windowIndex := state.WindowIndex
	modelIndex := state.ModelIndex
	for len(s.animationStates) <= windowIndex {
		s.animationStates = append(s.animationStates, make([]*renderer.AnimationStates, 0))
	}
	for len(s.animationStates[windowIndex]) <= modelIndex {
		s.animationStates[windowIndex] = append(s.animationStates[windowIndex], renderer.NewAnimationStates())
	}
	s.animationStates[windowIndex][modelIndex].Next = state
}

func (s *appStateChannels) Frame() float64 {
	return <-s.frameChan
}

func (s *appStateChannels) SetFrame(frame float64) {
	s.frameChan <- frame
}

func (s *appStateChannels) AddFrame(v float64) {
	f := s.Frame()
	s.frameChan <- v + f
}

func (s *appStateChannels) MaxFrame() int {
	return <-s.maxFrameChan
}

func (s *appStateChannels) SetMaxFrame(maxFrame int) {
	s.maxFrameChan <- maxFrame
}

func (s *appStateChannels) PrevFrame() int {
	return <-s.prevFrameChan
}

func (s *appStateChannels) SetPrevFrame(prevFrame int) {
	s.prevFrameChan <- prevFrame
}

func (s *appStateChannels) IsEnabledFrameDrop() bool {
	return <-s.isEnabledFrameDropChan
}

func (s *appStateChannels) SetEnabledFrameDrop(enabled bool) {
	s.isEnabledFrameDropChan <- enabled
}

func (s *appStateChannels) IsEnabledPhysics() bool {
	return <-s.isEnabledPhysicsChan
}

func (s *appStateChannels) SetEnabledPhysics(enabled bool) {
	s.isEnabledPhysicsChan <- enabled
}

func (s *appStateChannels) IsPhysicsReset() bool {
	return <-s.isPhysicsResetChan
}

func (s *appStateChannels) SetPhysicsReset(reset bool) {
	s.physicsResetChan <- reset
}

func (s *appStateChannels) IsShowNormal() bool {
	return <-s.isShowNormalChan
}

func (s *appStateChannels) SetShowNormal(show bool) {
	s.isShowNormalChan <- show
}

func (s *appStateChannels) IsShowWire() bool {
	return <-s.isShowWireChan
}

func (s *appStateChannels) SetShowWire(show bool) {
	s.isShowWireChan <- show
}

func (s *appStateChannels) IsShowSelectedVertex() bool {
	return <-s.isShowSelectedVertexChan
}

func (s *appStateChannels) SetShowSelectedVertex(show bool) {
	s.isShowSelectedVertexChan <- show
}

func (s *appStateChannels) IsShowBoneAll() bool {
	return <-s.isShowBoneAllChan
}

func (s *appStateChannels) SetShowBoneAll(show bool) {
	s.isShowBoneAllChan <- show
}

func (s *appStateChannels) IsShowBoneIk() bool {
	return <-s.isShowBoneIkChan
}

func (s *appStateChannels) SetShowBoneIk(show bool) {
	s.isShowBoneIkChan <- show
}

func (s *appStateChannels) IsShowBoneEffector() bool {
	return <-s.isShowBoneEffectorChan
}

func (s *appStateChannels) SetShowBoneEffector(show bool) {
	s.isShowBoneEffectorChan <- show
}

func (s *appStateChannels) IsShowBoneFixed() bool {
	return <-s.isShowBoneFixedChan
}

func (s *appStateChannels) SetShowBoneFixed(show bool) {
	s.isShowBoneFixedChan <- show
}

func (s *appStateChannels) IsShowBoneRotate() bool {
	return <-s.isShowBoneRotateChan
}

func (s *appStateChannels) SetShowBoneRotate(show bool) {
	s.isShowBoneRotateChan <- show
}

func (s *appStateChannels) IsShowBoneTranslate() bool {
	return <-s.isShowBoneTranslateChan
}

func (s *appStateChannels) SetShowBoneTranslate(show bool) {
	s.isShowBoneTranslateChan <- show
}

func (s *appStateChannels) IsShowBoneVisible() bool {
	return <-s.isShowBoneVisibleChan
}

func (s *appStateChannels) SetShowBoneVisible(show bool) {
	s.isShowBoneVisibleChan <- show
}

func (s *appStateChannels) IsShowRigidBodyFront() bool {
	return <-s.isShowRigidBodyFrontChan
}

func (s *appStateChannels) SetShowRigidBodyFront(show bool) {
	s.isShowRigidBodyFrontChan <- show
}

func (s *appStateChannels) IsShowRigidBodyBack() bool {
	return <-s.isShowRigidBodyBackChan
}

func (s *appStateChannels) SetShowRigidBodyBack(show bool) {
	s.isShowRigidBodyBackChan <- show
}

func (s *appStateChannels) IsShowJoint() bool {
	return <-s.isShowJointChan
}

func (s *appStateChannels) SetShowJoint(show bool) {
	s.isShowJointChan <- show
}

func (s *appStateChannels) IsShowInfo() bool {
	return <-s.isShowInfoChan
}

func (s *appStateChannels) SetShowInfo(show bool) {
	s.isShowInfoChan <- show
}

func (s *appStateChannels) IsLimitFps30() bool {
	return <-s.isLimitFps30Chan
}

func (s *appStateChannels) SetLimitFps30(limit bool) {
	s.isLimitFps30Chan <- limit
}

func (s *appStateChannels) IsLimitFps60() bool {
	return <-s.isLimitFps60Chan
}

func (s *appStateChannels) SetLimitFps60(limit bool) {
	s.isLimitFps60Chan <- limit
}

func (s *appStateChannels) IsUnLimitFps() bool {
	return <-s.isUnLimitFpsChan
}

func (s *appStateChannels) SetUnLimitFps(limit bool) {
	s.isUnLimitFpsChan <- limit
}

func (s *appStateChannels) IsUnLimitFpsDeform() bool {
	return <-s.isUnLimitFpsDeformChan
}

func (s *appStateChannels) SetUnLimitFpsDeform(limit bool) {
	s.isUnLimitFpsDeformChan <- limit
}

func (s *appStateChannels) IsLogLevelDebug() bool {
	return <-s.isLogLevelDebugChan
}

func (s *appStateChannels) SetLogLevelDebug(log bool) {
	s.isLogLevelDebugChan <- log
}

func (s *appStateChannels) IsLogLevelVerbose() bool {
	return <-s.isLogLevelVerboseChan
}

func (s *appStateChannels) SetLogLevelVerbose(log bool) {
	s.isLogLevelVerboseChan <- log
}

func (s *appStateChannels) IsLogLevelIkVerbose() bool {
	return <-s.isLogLevelIkVerboseChan
}

func (s *appStateChannels) SetLogLevelIkVerbose(log bool) {
	s.isLogLevelIkVerboseChan <- log
}

func (s *appStateChannels) IsClosed() bool {
	return <-s.isClosedChan
}

func (s *appStateChannels) SetClosed(closed bool) {
	s.isClosedChan <- closed
}

func (s *appStateChannels) Playing() bool {
	return <-s.playingChan
}

func (s *appStateChannels) TriggerPlay(p bool) {
	s.playingChan <- p
}

func (s *appStateChannels) SpfLimit() float64 {
	return <-s.spfLimitChan
}

func (s *appStateChannels) SetSpfLimit(spf float64) {
	s.spfLimitChan <- spf
}
