//go:build windows
// +build windows

package controller

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/renderer"
	"github.com/miu200521358/mlib_go/pkg/interface/core"
)

type controlState struct {
	appState                 core.IAppState                // アプリ状態
	motionPlayer             core.IPlayer                  // モーションプレイヤー
	controlWindow            core.IControlWindow           // コントロールウィンドウ
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

func NewControlState(appState core.IAppState) *controlState {
	u := &controlState{
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

func (s *controlState) Run() {
	go func() {
		prevTime := glfw.GetTime()

		for !s.appState.IsClosed() {
			frameTime := glfw.GetTime()
			elapsed := frameTime - prevTime

			if s.Playing() {
				// 再生中はフレームを進める
				// 経過秒数をキーフレームの進捗具合に合わせて調整
				if s.appState.SpfLimit() < -1 || elapsed >= s.appState.SpfLimit() {
					// デフォームFPS制限なしの場合、フレーム番号を常に進める
					if s.appState.IsEnabledFrameDrop() {
						// フレームドロップONの時、経過秒数分進める
						s.AddFrame(elapsed * 30)
					} else {
						// フレームドロップOFFの時、1だけ進める
						s.AddFrame(1)
					}
					prevTime = frameTime
				}
			} else {
				// 停止中はそのまま進める
				prevTime = frameTime
			}
		}
	}()

	go func() {
		for !s.appState.IsClosed() {
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

func (c *controlState) SetPlayer(mp core.IPlayer) {
	c.motionPlayer = mp
}

func (c *controlState) SetControlWindow(cw core.IControlWindow) {
	c.controlWindow = cw
}

func (c *controlState) SetAnimationState(state core.IAnimationState) {
	c.animationState <- state.(*renderer.AnimationState)
}

func (c *controlState) AnimationState() *renderer.AnimationState {
	return <-c.animationState
}

func (c *controlState) Frame() float64 {
	return c.motionPlayer.Frame()
}

func (c *controlState) SetFrame(frame float64) {
	c.motionPlayer.SetFrame(frame)
	c.frameChan <- frame
}

func (c *controlState) AddFrame(v float64) {
	f := c.Frame()
	c.SetFrame(f + v)
}

func (c *controlState) MaxFrame() int {
	return c.motionPlayer.MaxFrame()
}

func (c *controlState) SetMaxFrame(maxFrame int) {
	c.motionPlayer.SetMaxFrame(maxFrame)
	c.maxFrameChan <- maxFrame
}

func (c *controlState) UpdateMaxFrame(maxFrame int) {
	c.motionPlayer.UpdateMaxFrame(maxFrame)
	c.maxFrameChan <- maxFrame
}

func (c *controlState) PrevFrame() int {
	return c.motionPlayer.PrevFrame()
}

func (c *controlState) SetPrevFrame(prevFrame int) {
	c.motionPlayer.SetPrevFrame(prevFrame)
	c.prevFrameChan <- prevFrame
}

func (c *controlState) SetEnabledFrameDrop(enabled bool) {
	c.isEnabledFrameDropChan <- enabled
}

func (c *controlState) SetEnabledPhysics(enabled bool) {
	c.isEnabledPhysicsChan <- enabled
}

func (c *controlState) SetPhysicsReset(reset bool) {
	c.physicsResetChan <- reset
}

func (c *controlState) SetShowNormal(show bool) {
	c.isShowNormalChan <- show
}

func (c *controlState) SetShowWire(show bool) {
	c.isShowWireChan <- show
}

func (c *controlState) SetShowSelectedVertex(show bool) {
	c.isShowSelectedVertexChan <- show
}

func (c *controlState) SetShowBoneAll(show bool) {
	c.isShowBoneAllChan <- show
}

func (c *controlState) SetShowBoneIk(show bool) {
	c.isShowBoneIkChan <- show
}

func (c *controlState) SetShowBoneEffector(show bool) {
	c.isShowBoneEffectorChan <- show
}

func (c *controlState) SetShowBoneFixed(show bool) {
	c.isShowBoneFixedChan <- show
}

func (c *controlState) SetShowBoneRotate(show bool) {
	c.isShowBoneRotateChan <- show
}

func (c *controlState) SetShowBoneTranslate(show bool) {
	c.isShowBoneTranslateChan <- show
}

func (c *controlState) SetShowBoneVisible(show bool) {
	c.isShowBoneVisibleChan <- show
}

func (c *controlState) SetShowRigidBodyFront(show bool) {
	c.isShowRigidBodyFrontChan <- show
}

func (c *controlState) SetShowRigidBodyBack(show bool) {
	c.isShowRigidBodyBackChan <- show
}

func (c *controlState) SetShowJoint(show bool) {
	c.isShowJointChan <- show
}

func (c *controlState) SetShowInfo(show bool) {
	c.isShowInfoChan <- show
}

func (c *controlState) SetLimitFps30(limit bool) {
	c.isLimitFps30Chan <- limit
}

func (c *controlState) SetLimitFps60(limit bool) {
	c.isLimitFps60Chan <- limit
}

func (c *controlState) SetUnLimitFps(limit bool) {
	c.isUnLimitFpsChan <- limit
}

func (c *controlState) SetUnLimitFpsDeform(limit bool) {
	c.isUnLimitFpsDeformChan <- limit
}

func (c *controlState) SetLogLevelDebug(log bool) {
	c.isLogLevelDebugChan <- log
}

func (c *controlState) SetLogLevelVerbose(log bool) {
	c.isLogLevelVerboseChan <- log
}

func (c *controlState) SetLogLevelIkVerbose(log bool) {
	c.isLogLevelIkVerboseChan <- log
}

func (c *controlState) SetClosed(closed bool) {
	c.isClosedChan <- closed
}

func (c *controlState) Playing() bool {
	return c.motionPlayer != nil && c.motionPlayer.Playing()
}

func (c *controlState) TriggerPlay(p bool) {
	c.playingChan <- p
}

func (c *controlState) SetSpfLimit(spf float64) {
	c.spfLimitChan <- spf
}
