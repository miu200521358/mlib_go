package mbt

import (
	"strings"
	"time"

	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/physics"
)

// DebugBoneHoverInfo はボーンデバッグホバー情報を返します
func (mp *MPhysics) DebugBoneHoverInfo() []*physics.DebugBoneHover {
	return mp.debugBoneHover
}

// UpdateDebugHoverByBones は複数ボーンによるハイライト情報を更新します
func (mp *MPhysics) UpdateDebugHoverByBones(closestBones []*physics.DebugBoneHover, enable bool) {
	mlog.V("複数ボーンハイライト開始: enable=%v, bone数=%d", enable, len(closestBones))

	if !enable || len(closestBones) == 0 {
		mlog.V("複数ボーンハイライト無効またはbonesが空 - クリア")
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
	mlog.V("複数ボーンハイライト設定完了: %d個のボーン [%s]", len(closestBones), strings.Join(boneNames, ", "))
}

// CheckAndClearBoneExpiredHighlight は2秒経過したボーンハイライトを自動的にクリアします
func (mp *MPhysics) CheckAndClearBoneExpiredHighlight() {
	if mp.debugBoneHover == nil {
		// ハイライトが設定されていない場合は何もしない
		return
	}

	// 2秒経過をチェック
	if time.Since(mp.debugBoneHoverStartTime) >= 2*time.Second {
		mlog.V("ボーンハイライト自動クリア: 2秒経過しました")
		mp.clearDebugBoneHover()
	}
}

// clearDebugBoneHover はボーンデバッグホバー情報をクリアします
func (mp *MPhysics) clearDebugBoneHover() {
	mp.debugBoneHover = nil
}
