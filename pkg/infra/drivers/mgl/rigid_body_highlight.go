// 指示: miu200521358
package mgl

import (
	"time"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
	"github.com/miu200521358/mlib_go/pkg/domain/mmath"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	sharedlogging "github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// DebugRigidBodyHover はデバッグカーソル下の剛体情報を保持する。
type DebugRigidBodyHover struct {
	RigidBody *model.RigidBody
	HitPoint  *mmath.Vec3
}

// RigidBodyHighlighter は剛体デバッグハイライトを管理する。
type RigidBodyHighlighter struct {
	debugHover              *DebugRigidBodyHover
	debugHoverStartTime     time.Time
	prevRigidBodyDebugState bool
}

// NewRigidBodyHighlighter はRigidBodyHighlighterを生成する。
func NewRigidBodyHighlighter() *RigidBodyHighlighter {
	return &RigidBodyHighlighter{}
}

// DebugHoverInfo は剛体デバッグホバー情報を返す。
func (mp *RigidBodyHighlighter) DebugHoverInfo() *DebugRigidBodyHover {
	return mp.debugHover
}

// UpdateDebugHoverByRigidBody は剛体ハイライト情報を更新する。
func (mp *RigidBodyHighlighter) UpdateDebugHoverByRigidBody(modelIndex int, rb *model.RigidBody, enable bool) {
	logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_PHYSICS, "剛体ハイライト更新: enable=%v, rigidBody=%v", enable, rb != nil)
	if !enable || rb == nil {
		mp.clearDebugHover()
		return
	}
	mp.debugHover = &DebugRigidBodyHover{RigidBody: rb}
	mp.debugHoverStartTime = time.Now()
}

// DrawDebugHighlight は剛体デバッグハイライトを描画する。
func (mp *RigidBodyHighlighter) DrawDebugHighlight(shader graphics_api.IShader, isDrawRigidBodyFront bool) {
	if mp.debugHover == nil {
		return
	}
	logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_PHYSICS, "剛体ハイライト描画は未実装のためスキップします")
}

// CheckAndClearHighlightOnDebugChange は剛体デバッグ状態変更時にハイライトをクリアする。
func (mp *RigidBodyHighlighter) CheckAndClearHighlightOnDebugChange(currentDebugState bool) {
	if mp.prevRigidBodyDebugState != currentDebugState {
		if !currentDebugState {
			mp.clearDebugHover()
		}
		mp.prevRigidBodyDebugState = currentDebugState
	}
}

// CheckAndClearExpiredHighlight は2秒経過したハイライトを自動的にクリアする。
func (mp *RigidBodyHighlighter) CheckAndClearExpiredHighlight() {
	if mp.debugHover == nil {
		return
	}
	if time.Since(mp.debugHoverStartTime) >= 2*time.Second {
		logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_PHYSICS, "剛体ハイライト自動クリア: 2秒経過しました")
		mp.clearDebugHover()
	}
}

// clearDebugHover は剛体デバッグホバー情報をクリアする。
func (mp *RigidBodyHighlighter) clearDebugHover() {
	mp.debugHover = nil
}
