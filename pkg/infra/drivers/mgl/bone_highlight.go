//go:build windows
// +build windows

// 指示: miu200521358
package mgl

import (
	"time"

	"github.com/miu200521358/mlib_go/pkg/adapter/graphics_api"
)

// BoneHighlighter はボーンのデバッグハイライトを管理する。
type BoneHighlighter struct {
	debugHover              *DebugRigidBodyHover           // デバッグ用ホバー情報
	debugBoneHover          []*graphics_api.DebugBoneHover // ボーンデバッグ用ホバー情報
	debugBoneHoverStartTime time.Time                      // ボーンハイライト開始時刻（自動クリア用）
}

// NewBoneHighlighter はBoneHighlighterを生成する。
func NewBoneHighlighter() *BoneHighlighter {
	return &BoneHighlighter{}
}

// DebugBoneHoverInfo はボーンデバッグホバー情報を返す。
func (mp *BoneHighlighter) DebugBoneHoverInfo() []*graphics_api.DebugBoneHover {
	return mp.debugBoneHover
}

// CheckAndClearExpiredHighlight は2秒経過したハイライトを自動的にクリアする。
func (mp *BoneHighlighter) CheckAndClearExpiredHighlight() {
	if mp.debugBoneHover == nil {
		// ハイライトが設定されていない場合は何もしない
		return
	}

	// 2秒経過をチェック
	if time.Since(mp.debugBoneHoverStartTime) >= 2*time.Second {
		mp.clearDebugBoneHover()
	}
}

// UpdateDebugHoverByBones は複数ボーンによるハイライト情報を更新する。
func (mp *BoneHighlighter) UpdateDebugHoverByBones(closestBones []*graphics_api.DebugBoneHover, enable bool) {
	if !enable || len(closestBones) == 0 {
		mp.clearDebugBoneHover()
		return
	}

	mp.debugBoneHover = closestBones
	mp.debugBoneHoverStartTime = time.Now() // タイマー開始
}

// CheckAndClearBoneExpiredHighlight は2秒経過したボーンハイライトを自動的にクリアする。
func (mp *BoneHighlighter) CheckAndClearBoneExpiredHighlight() {
	if mp.debugBoneHover == nil {
		// ハイライトが設定されていない場合は何もしない
		return
	}

	// 2秒経過をチェック
	if time.Since(mp.debugBoneHoverStartTime) >= 2*time.Second {
		mp.clearDebugBoneHover()
	}
}

// clearDebugBoneHover はボーンデバッグホバー情報をクリアする。
func (mp *BoneHighlighter) clearDebugBoneHover() {
	mp.debugBoneHover = nil
}
