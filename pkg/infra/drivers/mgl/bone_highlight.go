// 指示: miu200521358
package mgl

import (
	"strings"
	"time"

	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	sharedlogging "github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// DebugBoneHover はデバッグカーソル下のボーン情報を保持する。
type DebugBoneHover struct {
	ModelIndex int         // モデルインデックス
	Bone       *model.Bone // 検出されたボーン
	Distance   float64     // カーソルからボーンラインまでの最短距離
}

// BoneHighlighter はボーンのデバッグハイライトを管理する。
type BoneHighlighter struct {
	debugHover              *DebugRigidBodyHover // デバッグ用ホバー情報
	debugBoneHover          []*DebugBoneHover    // ボーンデバッグ用ホバー情報
	debugBoneHoverStartTime time.Time            // ボーンハイライト開始時刻（自動クリア用）
}

// NewBoneHighlighter はBoneHighlighterを生成する。
func NewBoneHighlighter() *BoneHighlighter {
	return &BoneHighlighter{}
}

// DebugBoneHoverInfo はボーンデバッグホバー情報を返す。
func (mp *BoneHighlighter) DebugBoneHoverInfo() []*DebugBoneHover {
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
		logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_VIEWER, "ハイライト自動クリア: 2秒経過しました")
		mp.clearDebugBoneHover()
	}
}

// UpdateDebugHoverByBones は複数ボーンによるハイライト情報を更新する。
func (mp *BoneHighlighter) UpdateDebugHoverByBones(closestBones []*DebugBoneHover, enable bool) {
	logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_VIEWER, "複数ボーンハイライト開始: enable=%v, bone数=%d", enable, len(closestBones))

	if !enable || len(closestBones) == 0 {
		logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_VIEWER, "複数ボーンハイライト無効またはbonesが空 - クリア")
		mp.clearDebugBoneHover()
		return
	}

	mp.debugBoneHover = closestBones
	mp.debugBoneHoverStartTime = time.Now() // タイマー開始

	var boneNames []string
	for _, bone := range closestBones {
		if bone != nil {
			boneNames = append(boneNames, bone.Bone.Name())
		}
	}
	logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_VIEWER, "複数ボーンハイライト設定完了: %d個のボーン [%s]", len(closestBones), strings.Join(boneNames, ", "))
}

// CheckAndClearBoneExpiredHighlight は2秒経過したボーンハイライトを自動的にクリアする。
func (mp *BoneHighlighter) CheckAndClearBoneExpiredHighlight() {
	if mp.debugBoneHover == nil {
		// ハイライトが設定されていない場合は何もしない
		return
	}

	// 2秒経過をチェック
	if time.Since(mp.debugBoneHoverStartTime) >= 2*time.Second {
		logging.DefaultLogger().Verbose(sharedlogging.VERBOSE_INDEX_VIEWER, "ボーンハイライト自動クリア: 2秒経過しました")
		mp.clearDebugBoneHover()
	}
}

// clearDebugBoneHover はボーンデバッグホバー情報をクリアする。
func (mp *BoneHighlighter) clearDebugBoneHover() {
	mp.debugBoneHover = nil
}
