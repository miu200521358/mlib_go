//go:build windows
// +build windows

package app

import (
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
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
	isLogLevelDebug      bool                      // デバッグメッセージ表示
	isLogLevelVerbose    bool                      // 冗長メッセージ表示
	isLogLevelIkVerbose  bool                      // IK冗長メッセージ表示
	isClosed             bool                      // ウィンドウクローズ
	playing              bool                      // 再生中フラグ
	spfLimit             float64                   // FPS制限
	animationStates      [][]state.IAnimationState // アニメーションステート
	nextAnimationStates  [][]state.IAnimationState // 次のアニメーションステート
	mu                   sync.Mutex
}

func newAppState() *appState {
	u := &appState{
		isEnabledPhysics:    true,       // 物理ON
		isEnabledFrameDrop:  true,       // フレームドロップON
		isLimitFps30:        true,       // 30fps制限
		spfLimit:            1.0 / 30.0, // 30fps
		prevFrame:           -1,
		frame:               0,
		maxFrame:            1,
		animationStates:     make([][]state.IAnimationState, 0),
		nextAnimationStates: make([][]state.IAnimationState, 0),
	}

	return u
}

func (appState *appState) ExtendAnimationState(windowIndex, modelIndex int) {
	for len(appState.animationStates) <= windowIndex {
		appState.animationStates = append(appState.animationStates, make([]state.IAnimationState, 0))
		appState.nextAnimationStates = append(appState.nextAnimationStates, make([]state.IAnimationState, 0))
	}
	for len(appState.animationStates[windowIndex]) <= modelIndex {
		appState.animationStates[windowIndex] =
			append(appState.animationStates[windowIndex], animation.NewAnimationState(windowIndex, modelIndex))
		appState.nextAnimationStates[windowIndex] =
			append(appState.nextAnimationStates[windowIndex], animation.NewAnimationState(windowIndex, modelIndex))
	}
}

func (appState *appState) SetAnimationState(animationState state.IAnimationState) {
	windowIndex := animationState.WindowIndex()
	modelIndex := animationState.ModelIndex()
	appState.ExtendAnimationState(windowIndex, modelIndex)

	if animationState.Model() != nil {
		appState.nextAnimationStates[windowIndex][modelIndex] = animation.NewAnimationState(windowIndex, modelIndex)
		appState.nextAnimationStates[windowIndex][modelIndex].SetModel(animationState.Model())

		motion := appState.animationStates[windowIndex][modelIndex].Motion()
		if motion == nil {
			motion = vmd.NewVmdMotion("")
		}
		animationState.SetMotion(motion)

		vmdDeltas, renderDeltas := animationState.DeformBeforePhysics(appState, animationState.Model())
		appState.nextAnimationStates[windowIndex][modelIndex].SetMotion(animationState.Motion())
		appState.nextAnimationStates[windowIndex][modelIndex].SetVmdDeltas(vmdDeltas)
		appState.nextAnimationStates[windowIndex][modelIndex].SetRenderDeltas(renderDeltas)

	} else if animationState.Motion() != nil {
		// モーションが指定されてたらセット
		model := appState.animationStates[windowIndex][modelIndex].Model()
		if model != nil {
			appState.nextAnimationStates[windowIndex][modelIndex] = animation.NewAnimationState(windowIndex, modelIndex)

			vmdDeltas, renderDeltas := animationState.DeformBeforePhysics(appState, model)
			appState.nextAnimationStates[windowIndex][modelIndex].SetMotion(animationState.Motion())
			appState.nextAnimationStates[windowIndex][modelIndex].SetVmdDeltas(vmdDeltas)
			appState.nextAnimationStates[windowIndex][modelIndex].SetRenderDeltas(renderDeltas)
		}
	}

	if animationState.InvisibleMaterialIndexes() != nil {
		appState.nextAnimationStates[windowIndex][modelIndex].SetInvisibleMaterialIndexes(
			animationState.InvisibleMaterialIndexes())
	}
}

func (appState *appState) Frame() float64 {
	return appState.frame
}
func (appState *appState) SetFrame(frame float64) {
	appState.frame = frame

	var wg sync.WaitGroup
	for i := range appState.animationStates {
		for j := range appState.animationStates[i] {
			if appState.animationStates[i][j] != nil && appState.animationStates[i][j].Model() != nil {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()

					vmdDeltas, renderDeltas := appState.animationStates[i][j].DeformBeforePhysics(
						appState, appState.animationStates[i][j].Model())

					nextState := animation.NewAnimationState(i, j)
					nextState.SetVmdDeltas(vmdDeltas)
					nextState.SetRenderDeltas(renderDeltas)

					appState.mu.Lock()
					defer appState.mu.Unlock()
					appState.nextAnimationStates[i][j] = nextState
				}(i, j)
			}
		}
	}

	// すべてのGoroutineが終了するのを待つ
	wg.Wait()
}

func (appState *appState) AddFrame(v float64) {
	appState.SetFrame(appState.frame + v)
}

func (appState *appState) MaxFrame() int {
	return appState.maxFrame
}

func (appState *appState) UpdateMaxFrame(maxFrame int) {
	if appState.maxFrame < maxFrame {
		appState.maxFrame = maxFrame
	}
}

func (appState *appState) SetMaxFrame(maxFrame int) {
	appState.maxFrame = maxFrame
}

func (appState *appState) PrevFrame() int {
	return appState.prevFrame
}

func (appState *appState) SetPrevFrame(prevFrame int) {
	appState.prevFrame = prevFrame
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

func (appState *appState) IsLogLevelDebug() bool {
	return appState.isLogLevelDebug
}

func (appState *appState) SetLogLevelDebug(log bool) {
	appState.isLogLevelDebug = log
}

func (appState *appState) IsLogLevelVerbose() bool {
	return appState.isLogLevelVerbose
}

func (appState *appState) SetLogLevelVerbose(log bool) {
	appState.isLogLevelVerbose = log
}

func (appState *appState) IsLogLevelIkVerbose() bool {
	return appState.isLogLevelIkVerbose
}

func (appState *appState) SetLogLevelIkVerbose(log bool) {
	appState.isLogLevelIkVerbose = log
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

func (appState *appState) TriggerPlay(p bool) {
	appState.playing = p
}

func (appState *appState) SpfLimit() float64 {
	return appState.spfLimit
}

func (appState *appState) SetSpfLimit(spf float64) {
	appState.spfLimit = spf
}
