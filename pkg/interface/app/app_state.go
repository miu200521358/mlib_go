//go:build windows
// +build windows

package app

import (
	"github.com/miu200521358/mlib_go/pkg/infrastructure/animation"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/state"
)

type appState struct {
	frame                float64                   // フレーム
	prevFrame            int                       // 前回のフレーム
	maxFrame             int                       // 最大フレーム
	isEnabledFrameDrop   bool                      // フレームドロップON/OFF
	isEnabledPhysics     bool                      // 物理ON/OFF
	isPhysicsReset       bool                      // 物理リセット
	isShowNormal         bool                      // ボーンデバッグ表示
	isShowWire           bool                      // ワイヤーフレームデバッグ表示
	isShowSelectedVertex bool                      // 選択頂点デバッグ表示
	isShowBoneAll        bool                      // 全ボーンデバッグ表示
	isShowBoneIk         bool                      // IKボーンデバッグ表示
	isShowBoneEffector   bool                      // 付与親ボーンデバッグ表示
	isShowBoneFixed      bool                      // 軸制限ボーンデバッグ表示
	isShowBoneRotate     bool                      // 回転ボーンデバッグ表示
	isShowBoneTranslate  bool                      // 移動ボーンデバッグ表示
	isShowBoneVisible    bool                      // 表示ボーンデバッグ表示
	isShowRigidBodyFront bool                      // 剛体デバッグ表示(前面)
	isShowRigidBodyBack  bool                      // 剛体デバッグ表示(埋め込み)
	isShowJoint          bool                      // ジョイントデバッグ表示
	isShowInfo           bool                      // 情報デバッグ表示
	isLimitFps30         bool                      // 30FPS制限
	isLimitFps60         bool                      // 60FPS制限
	isUnLimitFps         bool                      // FPS無制限
	isUnLimitFpsDeform   bool                      // デフォームFPS無制限
	isLogLevelDebug      bool                      // デバッグメッセージ表示
	isLogLevelVerbose    bool                      // 冗長メッセージ表示
	isLogLevelIkVerbose  bool                      // IK冗長メッセージ表示
	isClosed             bool                      // ウィンドウクローズ
	playing              bool                      // 再生中フラグ
	spfLimit             float64                   // FPS制限
	animationStates      [][]state.IAnimationState // アニメーションステート
	nextState            state.IAnimationState     // 次のアニメーションステート
}

func newAppState() *appState {
	u := &appState{
		isEnabledPhysics:   true,       // 物理ON
		isEnabledFrameDrop: true,       // フレームドロップON
		isLimitFps30:       true,       // 30fps制限
		spfLimit:           1.0 / 30.0, // 30fps
		prevFrame:          -1,
		frame:              0,
		maxFrame:           1,
		animationStates:    make([][]state.IAnimationState, 0),
	}

	return u
}

func (appState *appState) SetAnimationState(nextState state.IAnimationState) {
	windowIndex := nextState.WindowIndex()
	modelIndex := nextState.ModelIndex()

	// 必要に応じて拡張
	for len(appState.animationStates) <= windowIndex {
		appState.animationStates = append(appState.animationStates, make([]state.IAnimationState, 0))
	}
	for i := 0; i < len(appState.animationStates); i++ {
		for len(appState.animationStates[i]) <= modelIndex {
			appState.animationStates[i] = append(appState.animationStates[i],
				animation.NewAnimationState(windowIndex, len(appState.animationStates[i])))
		}
	}

	// モーションが指定されてたらセット
	if nextState.Motion() != nil {
		appState.animationStates[windowIndex][modelIndex].SetMotion(nextState.Motion())
		model := appState.animationStates[windowIndex][modelIndex].Model()
		if model != nil {
			vmdDeltas, renderDeltas :=
				appState.animationStates[windowIndex][modelIndex].DeformBeforePhysics(appState, model)
			appState.animationStates[windowIndex][modelIndex].SetVmdDeltas(vmdDeltas)
			appState.animationStates[windowIndex][modelIndex].SetRenderDeltas(renderDeltas)
		}
	}

	if nextState.Model() != nil {
		// 次のステータスとしてモデル情報だけセット
		appState.nextState = animation.NewAnimationState(windowIndex, modelIndex)
		appState.nextState.SetModel(nextState.Model())
	}
}

func (appState *appState) Frame() float64 {
	return appState.frame
}

func (appState *appState) SetFrame(frame float64) {
	appState.frame = frame
	for i := 0; i < len(appState.animationStates); i++ {
		for j := 0; j < len(appState.animationStates[i]); j++ {
			if appState.animationStates[i][j] != nil && appState.animationStates[i][j].Model() != nil {
				vmdDeltas, renderDeltas := appState.animationStates[i][j].DeformBeforePhysics(
					appState, appState.animationStates[i][j].Model())
				appState.animationStates[i][j].SetVmdDeltas(vmdDeltas)
				appState.animationStates[i][j].SetRenderDeltas(renderDeltas)
			}
		}
	}
}

func (a *appState) AddFrame(v float64) {
	a.SetFrame(a.frame + v)
}

func (a *appState) MaxFrame() int {
	return a.maxFrame
}

func (a *appState) UpdateMaxFrame(maxFrame int) {
	if a.maxFrame < maxFrame {
		a.maxFrame = maxFrame
	}
}

func (a *appState) SetMaxFrame(maxFrame int) {
	a.maxFrame = maxFrame
}

func (a *appState) PrevFrame() int {
	return a.prevFrame
}

func (a *appState) SetPrevFrame(prevFrame int) {
	a.prevFrame = prevFrame
}

func (a *appState) IsEnabledFrameDrop() bool {
	return a.isEnabledFrameDrop
}

func (a *appState) SetEnabledFrameDrop(enabled bool) {
	a.isEnabledFrameDrop = enabled
}

func (a *appState) IsEnabledPhysics() bool {
	return a.isEnabledPhysics
}

func (a *appState) SetEnabledPhysics(enabled bool) {
	a.isEnabledPhysics = enabled
}

func (a *appState) IsPhysicsReset() bool {
	return a.isPhysicsReset
}

func (a *appState) SetPhysicsReset(reset bool) {
	a.isPhysicsReset = reset
}

func (a *appState) IsShowNormal() bool {
	return a.isShowNormal
}

func (a *appState) SetShowNormal(show bool) {
	a.isShowNormal = show
}

func (a *appState) IsShowWire() bool {
	return a.isShowWire
}

func (a *appState) SetShowWire(show bool) {
	a.isShowWire = show
}

func (a *appState) IsShowSelectedVertex() bool {
	return a.isShowSelectedVertex
}

func (a *appState) SetShowSelectedVertex(show bool) {
	a.isShowSelectedVertex = show
}

func (a *appState) IsShowBoneAll() bool {
	return a.isShowBoneAll
}

func (a *appState) SetShowBoneAll(show bool) {
	a.isShowBoneAll = show
}

func (a *appState) IsShowBoneIk() bool {
	return a.isShowBoneIk
}

func (a *appState) SetShowBoneIk(show bool) {
	a.isShowBoneIk = show
}

func (a *appState) IsShowBoneEffector() bool {
	return a.isShowBoneEffector
}

func (a *appState) SetShowBoneEffector(show bool) {
	a.isShowBoneEffector = show
}

func (a *appState) IsShowBoneFixed() bool {
	return a.isShowBoneFixed
}

func (a *appState) SetShowBoneFixed(show bool) {
	a.isShowBoneFixed = show
}

func (a *appState) IsShowBoneRotate() bool {
	return a.isShowBoneRotate
}

func (a *appState) SetShowBoneRotate(show bool) {
	a.isShowBoneRotate = show
}

func (a *appState) IsShowBoneTranslate() bool {
	return a.isShowBoneTranslate
}

func (a *appState) SetShowBoneTranslate(show bool) {
	a.isShowBoneTranslate = show
}

func (a *appState) IsShowBoneVisible() bool {
	return a.isShowBoneVisible
}

func (a *appState) SetShowBoneVisible(show bool) {
	a.isShowBoneVisible = show
}

func (a *appState) IsShowRigidBodyFront() bool {
	return a.isShowRigidBodyFront
}

func (a *appState) SetShowRigidBodyFront(show bool) {
	a.isShowRigidBodyFront = show
}

func (a *appState) IsShowRigidBodyBack() bool {
	return a.isShowRigidBodyBack
}

func (a *appState) SetShowRigidBodyBack(show bool) {
	a.isShowRigidBodyBack = show
}

func (a *appState) IsShowJoint() bool {
	return a.isShowJoint
}

func (a *appState) SetShowJoint(show bool) {
	a.isShowJoint = show
}

func (a *appState) IsShowInfo() bool {
	return a.isShowInfo
}

func (a *appState) SetShowInfo(show bool) {
	a.isShowInfo = show
}

func (a *appState) IsLimitFps30() bool {
	return a.isLimitFps30
}

func (a *appState) SetLimitFps30(limit bool) {
	a.isLimitFps30 = limit
}

func (a *appState) IsLimitFps60() bool {
	return a.isLimitFps60
}

func (a *appState) SetLimitFps60(limit bool) {
	a.isLimitFps60 = limit
}

func (a *appState) IsUnLimitFps() bool {
	return a.isUnLimitFps
}

func (a *appState) SetUnLimitFps(limit bool) {
	a.isUnLimitFps = limit
}

func (a *appState) IsUnLimitFpsDeform() bool {
	return a.isUnLimitFpsDeform
}

func (a *appState) SetUnLimitFpsDeform(limit bool) {
	a.isUnLimitFpsDeform = limit
}

func (a *appState) IsLogLevelDebug() bool {
	return a.isLogLevelDebug
}

func (a *appState) SetLogLevelDebug(log bool) {
	a.isLogLevelDebug = log
}

func (a *appState) IsLogLevelVerbose() bool {
	return a.isLogLevelVerbose
}

func (a *appState) SetLogLevelVerbose(log bool) {
	a.isLogLevelVerbose = log
}

func (a *appState) IsLogLevelIkVerbose() bool {
	return a.isLogLevelIkVerbose
}

func (a *appState) SetLogLevelIkVerbose(log bool) {
	a.isLogLevelIkVerbose = log
}

func (a *appState) IsClosed() bool {
	return a.isClosed
}

func (a *appState) SetClosed(closed bool) {
	a.isClosed = closed
}

func (a *appState) Playing() bool {
	return a.playing
}

func (a *appState) TriggerPlay(p bool) {
	a.playing = p
}

func (a *appState) SpfLimit() float64 {
	return a.spfLimit
}

func (a *appState) SetSpfLimit(spf float64) {
	a.spfLimit = spf
}
