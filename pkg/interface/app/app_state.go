//go:build windows
// +build windows

package app

import (
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
)

type appState struct {
	frame                float32                   // フレーム
	maxFrame             float32                   // 最大フレーム
	isEnabledFrameDrop   bool                      // フレームドロップON/OFF
	isEnabledPhysics     bool                      // 物理ON/OFF
	isPhysicsReset       bool                      // 物理リセット
	isShowNormal         bool                      // ボーンデバッグ表示
	isShowWire           bool                      // ワイヤーフレームデバッグ表示
	isShowOverride       bool                      // オーバーライドデバッグ表示
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
	isCameraSync         bool                      // レンダーシンク
	isClosed             bool                      // ウィンドウクローズ
	playing              bool                      // 再生中フラグ
	frameInterval        float64                   // FPS制限
	funcGetModels        func() [][]*pmx.PmxModel  // モデル取得関数
	funcGetMotions       func() [][]*vmd.VmdMotion // モーション取得関数
}

func newAppState() *appState {
	return &appState{
		isEnabledPhysics:   true,       // 物理ON
		isEnabledFrameDrop: true,       // フレームドロップON
		isLimitFps30:       true,       // 30fps制限
		frameInterval:      1.0 / 30.0, // 30fps
		frame:              0.0,
		maxFrame:           1,
	}
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

func (appState *appState) FrameInterval() float64 {
	return appState.frameInterval
}

func (appState *appState) SetFrameInterval(spf float64) {
	appState.frameInterval = spf
}

func (appState *appState) SetCameraSync(sync bool) {
	appState.isCameraSync = sync
}

func (appState *appState) IsCameraSync() bool {
	return appState.isCameraSync
}

func (appState *appState) SetFuncGetModels(f func() [][]*pmx.PmxModel) {
	appState.funcGetModels = f
}

func (appState *appState) SetFuncGetMotions(f func() [][]*vmd.VmdMotion) {
	appState.funcGetMotions = f
}

func (appState *appState) GetModels() [][]*pmx.PmxModel {
	return appState.funcGetModels()
}

func (appState *appState) GetMotions() [][]*vmd.VmdMotion {
	return appState.funcGetMotions()
}
