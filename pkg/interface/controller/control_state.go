//go:build windows
// +build windows

package controller

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/animation"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
)

type controlState struct {
	appState                 state.IAppState                // アプリ状態
	motionPlayer             app.IPlayer                    // モーションプレイヤー
	controlWindow            app.IControlWindow             // コントロールウィンドウ
	prevFrameChan            chan int                       // 前回フレーム
	frameChan                chan float64                   // フレーム
	maxFrameChan             chan int                       // 最大フレーム
	isEnabledFrameDropChan   chan bool                      // フレームドロップON/OFF
	isEnabledPhysicsChan     chan bool                      // 物理ON/OFF
	isPhysicsResetChan       chan bool                      // 物理リセット
	isShowNormalChan         chan bool                      // ボーンデバッグ表示
	isShowWireChan           chan bool                      // ワイヤーフレームデバッグ表示
	isShowOverrideChan       chan bool                      // オーバーライドデバッグ表示
	isShowSelectedVertexChan chan bool                      // 選択頂点デバッグ表示
	isShowBoneAllChan        chan bool                      // 全ボーンデバッグ表示
	isShowBoneIkChan         chan bool                      // IKボーンデバッグ表示
	isShowBoneEffectorChan   chan bool                      // 付与親ボーンデバッグ表示
	isShowBoneFixedChan      chan bool                      // 軸制限ボーンデバッグ表示
	isShowBoneRotateChan     chan bool                      // 回転ボーンデバッグ表示
	isShowBoneTranslateChan  chan bool                      // 移動ボーンデバッグ表示
	isShowBoneVisibleChan    chan bool                      // 表示ボーンデバッグ表示
	isShowRigidBodyFrontChan chan bool                      // 剛体デバッグ表示(前面)
	isShowRigidBodyBackChan  chan bool                      // 剛体デバッグ表示(埋め込み)
	isShowJointChan          chan bool                      // ジョイントデバッグ表示
	isShowInfoChan           chan bool                      // 情報デバッグ表示
	isLimitFps30Chan         chan bool                      // 30FPS制限
	isLimitFps60Chan         chan bool                      // 60FPS制限
	isUnLimitFpsChan         chan bool                      // FPS無制限
	isUnLimitFpsDeformChan   chan bool                      // デフォームFPS無制限
	isCameraSyncChan         chan bool                      // レンダリング同期
	isLogLevelDebugChan      chan bool                      // デバッグメッセージ表示
	isLogLevelVerboseChan    chan bool                      // 冗長メッセージ表示
	isLogLevelIkVerboseChan  chan bool                      // IK冗長メッセージ表示
	isClosedChan             chan bool                      // ウィンドウクローズ
	playingChan              chan bool                      // 再生中フラグ
	physicsResetChan         chan bool                      // 物理リセット
	spfLimitChan             chan float64                   // FPS制限
	animationState           chan *animation.AnimationState // アニメーションステート
}

func NewControlState(appState state.IAppState) *controlState {
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
		isShowOverrideChan:       make(chan bool),
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
		animationState:           make(chan *animation.AnimationState),
	}

	return u
}

func (contState *controlState) Run() {
	go func() {
		prevTime := glfw.GetTime()

		for !contState.appState.IsClosed() {
			frameTime := glfw.GetTime()
			elapsed := frameTime - prevTime

			if contState.Playing() {
				// 再生中はフレームを進める
				// 経過秒数をキーフレームの進捗具合に合わせて調整
				if elapsed >= contState.appState.SpfLimit() {
					// デフォームFPS制限なしの場合、フレーム番号を常に進める
					if contState.appState.IsEnabledFrameDrop() {
						// フレームドロップONの時、経過秒数分進める
						contState.AddFrame(elapsed * 30)
					} else {
						// フレームドロップOFFもしくはデフォーム無制限の時、1だけ進める
						contState.AddFrame(1)
					}

					if contState.Frame() > float64(contState.MaxFrame()) {
						// 最後まで行ったら物理リセットフラグを立てて、最初に戻す
						contState.appState.SetPhysicsReset(true)
						contState.appState.SetFrame(0)
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
		for !contState.appState.IsClosed() {
			select {
			case prevFrame := <-contState.prevFrameChan:
				contState.appState.SetPrevFrame(prevFrame)
			case frame := <-contState.frameChan:
				contState.appState.SetFrame(frame)
			case maxFrame := <-contState.maxFrameChan:
				contState.appState.UpdateMaxFrame(maxFrame)
			case enabledFrameDrop := <-contState.isEnabledFrameDropChan:
				contState.appState.SetEnabledFrameDrop(enabledFrameDrop)
			case enabledPhysics := <-contState.isEnabledPhysicsChan:
				contState.appState.SetEnabledPhysics(enabledPhysics)
			case resetPhysics := <-contState.physicsResetChan:
				contState.appState.SetPhysicsReset(resetPhysics)
			case showNormal := <-contState.isShowNormalChan:
				contState.appState.SetShowNormal(showNormal)
			case showWire := <-contState.isShowWireChan:
				contState.appState.SetShowWire(showWire)
			case showOverride := <-contState.isShowOverrideChan:
				contState.appState.SetShowOverride(showOverride)
			case showSelectedVertex := <-contState.isShowSelectedVertexChan:
				contState.appState.SetShowSelectedVertex(showSelectedVertex)
			case showBoneAll := <-contState.isShowBoneAllChan:
				contState.appState.SetShowBoneAll(showBoneAll)
			case showBoneIk := <-contState.isShowBoneIkChan:
				contState.appState.SetShowBoneIk(showBoneIk)
			case showBoneEffector := <-contState.isShowBoneEffectorChan:
				contState.appState.SetShowBoneEffector(showBoneEffector)
			case showBoneFixed := <-contState.isShowBoneFixedChan:
				contState.appState.SetShowBoneFixed(showBoneFixed)
			case showBoneRotate := <-contState.isShowBoneRotateChan:
				contState.appState.SetShowBoneRotate(showBoneRotate)
			case showBoneTranslate := <-contState.isShowBoneTranslateChan:
				contState.appState.SetShowBoneTranslate(showBoneTranslate)
			case showBoneVisible := <-contState.isShowBoneVisibleChan:
				contState.appState.SetShowBoneVisible(showBoneVisible)
			case showRigidBodyFront := <-contState.isShowRigidBodyFrontChan:
				contState.appState.SetShowRigidBodyFront(showRigidBodyFront)
			case showRigidBodyBack := <-contState.isShowRigidBodyBackChan:
				contState.appState.SetShowRigidBodyBack(showRigidBodyBack)
			case showJoint := <-contState.isShowJointChan:
				contState.appState.SetShowJoint(showJoint)
			case showInfo := <-contState.isShowInfoChan:
				contState.appState.SetShowInfo(showInfo)
			case spfLimit := <-contState.spfLimitChan:
				contState.appState.SetSpfLimit(spfLimit)
			case cameraSync := <-contState.isCameraSyncChan:
				contState.appState.SetCameraSync(cameraSync)
			case closed := <-contState.isClosedChan:
				contState.appState.SetClosed(closed)
			case playing := <-contState.playingChan:
				contState.appState.TriggerPlay(playing)
			case spfLimit := <-contState.spfLimitChan:
				contState.appState.SetSpfLimit(spfLimit)
			case animationState := <-contState.animationState:
				contState.appState.SetAnimationState(animationState)
			}
		}
	}()
}

func (contState *controlState) SetPlayer(mp app.IPlayer) {
	contState.motionPlayer = mp
}

func (contState *controlState) SetControlWindow(cw app.IControlWindow) {
	contState.controlWindow = cw
}

func (contState *controlState) SetAnimationState(state state.IAnimationState) {
	contState.animationState <- state.(*animation.AnimationState)
}

func (contState *controlState) AnimationState() *animation.AnimationState {
	return <-contState.animationState
}

func (contState *controlState) Frame() float64 {
	return contState.motionPlayer.Frame()
}

func (contState *controlState) SetFrame(frame float64) {
	contState.motionPlayer.SetFrame(frame)
	contState.frameChan <- frame
}

func (contState *controlState) AddFrame(v float64) {
	f := contState.Frame()
	contState.SetFrame(f + v)
}

func (contState *controlState) MaxFrame() int {
	return contState.motionPlayer.MaxFrame()
}

func (contState *controlState) SetMaxFrame(maxFrame int) {
	contState.motionPlayer.SetMaxFrame(maxFrame)
	contState.maxFrameChan <- maxFrame
}

func (contState *controlState) UpdateMaxFrame(maxFrame int) {
	contState.motionPlayer.UpdateMaxFrame(maxFrame)
	contState.maxFrameChan <- maxFrame
}

func (contState *controlState) PrevFrame() int {
	return contState.motionPlayer.PrevFrame()
}

func (contState *controlState) SetPrevFrame(prevFrame int) {
	contState.motionPlayer.SetPrevFrame(prevFrame)
	contState.prevFrameChan <- prevFrame
}

func (contState *controlState) SetEnabledFrameDrop(enabled bool) {
	contState.isEnabledFrameDropChan <- enabled
}

func (contState *controlState) SetEnabledPhysics(enabled bool) {
	contState.isEnabledPhysicsChan <- enabled
}

func (contState *controlState) SetPhysicsReset(reset bool) {
	contState.physicsResetChan <- reset
}

func (contState *controlState) SetShowNormal(show bool) {
	contState.isShowNormalChan <- show
}

func (contState *controlState) SetShowWire(show bool) {
	contState.isShowWireChan <- show
}

func (contState *controlState) SetShowOverride(show bool) {
	contState.isShowOverrideChan <- show
}

func (contState *controlState) SetShowSelectedVertex(show bool) {
	contState.isShowSelectedVertexChan <- show
}

func (contState *controlState) SetShowBoneAll(show bool) {
	contState.isShowBoneAllChan <- show
}

func (contState *controlState) SetShowBoneIk(show bool) {
	contState.isShowBoneIkChan <- show
}

func (contState *controlState) SetShowBoneEffector(show bool) {
	contState.isShowBoneEffectorChan <- show
}

func (contState *controlState) SetShowBoneFixed(show bool) {
	contState.isShowBoneFixedChan <- show
}

func (contState *controlState) SetShowBoneRotate(show bool) {
	contState.isShowBoneRotateChan <- show
}

func (contState *controlState) SetShowBoneTranslate(show bool) {
	contState.isShowBoneTranslateChan <- show
}

func (contState *controlState) SetShowBoneVisible(show bool) {
	contState.isShowBoneVisibleChan <- show
}

func (contState *controlState) SetShowRigidBodyFront(show bool) {
	contState.isShowRigidBodyFrontChan <- show
}

func (contState *controlState) SetShowRigidBodyBack(show bool) {
	contState.isShowRigidBodyBackChan <- show
}

func (contState *controlState) SetShowJoint(show bool) {
	contState.isShowJointChan <- show
}

func (contState *controlState) SetShowInfo(show bool) {
	contState.isShowInfoChan <- show
}

func (contState *controlState) SetLimitFps30(limit bool) {
	contState.isLimitFps30Chan <- limit
}

func (contState *controlState) SetLimitFps60(limit bool) {
	contState.isLimitFps60Chan <- limit
}

func (contState *controlState) SetUnLimitFps(limit bool) {
	contState.isUnLimitFpsChan <- limit
}

func (contState *controlState) SetUnLimitFpsDeform(limit bool) {
	contState.isUnLimitFpsDeformChan <- limit
}

func (contState *controlState) SetLogLevelDebug(log bool) {
	contState.isLogLevelDebugChan <- log
}

func (contState *controlState) SetLogLevelVerbose(log bool) {
	contState.isLogLevelVerboseChan <- log
}

func (contState *controlState) SetLogLevelIkVerbose(log bool) {
	contState.isLogLevelIkVerboseChan <- log
}

func (contState *controlState) SetClosed(closed bool) {
	contState.isClosedChan <- closed
}

func (contState *controlState) Playing() bool {
	return contState.motionPlayer != nil && contState.motionPlayer.Playing()
}

func (contState *controlState) TriggerPlay(p bool) {
	contState.playingChan <- p
}

func (contState *controlState) SetSpfLimit(spf float64) {
	contState.spfLimitChan <- spf
}

func (contState *controlState) IsCameraSync() bool {
	return contState.appState.IsCameraSync()
}

func (contState *controlState) SetCameraSync(sync bool) {
	contState.appState.SetCameraSync(sync)
}
