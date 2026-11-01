//go:build windows
// +build windows

// 指示: miu200521358
package controller

import (
	"testing"

	"github.com/miu200521358/mlib_go/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
	"github.com/miu200521358/walk/pkg/declarative"
)

type focusLinkUserConfigStub struct {
	config.IUserConfig
	getBoolCalls   int
	getBoolKey     string
	getBoolDefault bool
	getBoolValue   bool
	setBoolCalls   int
	setBoolKey     string
	setBoolValue   bool
}

// GetBool はフォーカス連動設定の読み取りを記録する。
func (s *focusLinkUserConfigStub) GetBool(key string, defaultValue bool) (bool, error) {
	s.getBoolCalls++
	s.getBoolKey = key
	s.getBoolDefault = defaultValue
	return s.getBoolValue, nil
}

// SetBool はフォーカス連動設定の保存を記録する。
func (s *focusLinkUserConfigStub) SetBool(key string, value bool) error {
	s.setBoolCalls++
	s.setBoolKey = key
	s.setBoolValue = value
	return nil
}

// hasTopLevelActionText はトップレベルのアクションに指定文言があるか調べる。
func hasTopLevelActionText(items []declarative.MenuItem, text string) bool {
	for _, item := range items {
		if action, ok := item.(declarative.Action); ok && action.Text == text {
			return true
		}
	}
	return false
}

func TestNonDevBuildHidesAndIgnoresFocusLinkSetting(t *testing.T) {
	shared := state.NewSharedState(1)
	userConfig := &focusLinkUserConfigStub{getBoolValue: false}
	cw := &ControlWindow{
		shared:     shared,
		appConfig:  &config.AppConfig{EnvValue: config.APP_ENV_PROD},
		userConfig: userConfig,
	}

	cw.applyFocusLinkUserConfig()
	items := cw.buildControllerMenuItems()
	if hasTopLevelActionText(items, cw.t(messages.ControlWindowKey110)) {
		t.Fatalf("通常ビルドに最前面化メニューが表示されています")
	}
	if userConfig.getBoolCalls != 0 {
		t.Fatalf("通常ビルドで focus_link を読み取りました: %d 回", userConfig.getBoolCalls)
	}
	if !shared.IsFocusLinkEnabled() {
		t.Fatalf("通常ビルドのフォーカス連動は既定で有効であるべきです")
	}

	cw.SetFocusLinkEnabled(false)
	if userConfig.setBoolCalls != 0 {
		t.Fatalf("通常ビルドで focus_link を保存しました: %d 回", userConfig.setBoolCalls)
	}
	if !shared.IsFocusLinkEnabled() {
		t.Fatalf("通常ビルドで最前面化連動が無効になりました")
	}
}

func TestDevBuildLoadsAndSavesFocusLinkSetting(t *testing.T) {
	shared := state.NewSharedState(1)
	userConfig := &focusLinkUserConfigStub{getBoolValue: false}
	cw := &ControlWindow{
		shared:     shared,
		appConfig:  &config.AppConfig{EnvValue: config.APP_ENV_DEV},
		userConfig: userConfig,
	}

	cw.applyFocusLinkUserConfig()
	if shared.IsFocusLinkEnabled() {
		t.Fatalf("開発ビルドの保存済み OFF が反映されていません")
	}
	if userConfig.getBoolCalls != 1 || userConfig.getBoolKey != config.UserConfigKeyFocusLink || !userConfig.getBoolDefault {
		t.Fatalf("開発ビルドの設定読み取りが不正です: %+v", userConfig)
	}
	items := cw.buildControllerMenuItems()
	if !hasTopLevelActionText(items, cw.t(messages.ControlWindowKey110)) {
		t.Fatalf("開発ビルドに最前面化メニューがありません")
	}

	cw.SetFocusLinkEnabled(true)
	if userConfig.setBoolCalls != 1 || userConfig.setBoolKey != config.UserConfigKeyFocusLink || !userConfig.setBoolValue {
		t.Fatalf("開発ビルドの設定保存が不正です: %+v", userConfig)
	}
}

func TestWindowLinkageAndFocusLinkRemainIndependent(t *testing.T) {
	shared := state.NewSharedState(1)
	cw := &ControlWindow{shared: shared}
	shared.SetFocusLinkEnabled(false)
	cw.SetWindowLinkageEnabled(true)
	if !shared.HasFlag(state.STATE_FLAG_WINDOW_LINKAGE) || shared.IsFocusLinkEnabled() {
		t.Fatalf("画面移動連動の変更がフォーカス連動へ漏れています")
	}

	shared.SetFocusLinkEnabled(true)
	cw.SetWindowLinkageEnabled(false)
	if shared.HasFlag(state.STATE_FLAG_WINDOW_LINKAGE) || !shared.IsFocusLinkEnabled() {
		t.Fatalf("フォーカス連動の変更が画面移動連動へ漏れています")
	}
}
