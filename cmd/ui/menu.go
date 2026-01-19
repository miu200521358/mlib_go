//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/infra/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
)

// NewMenuItems はサンプル用のメニュー項目を生成する。
func NewMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text: i18n.T("&サンプルメニュー"),
			OnTriggered: func() {
				logging.DefaultLogger().InfoLine(i18n.T("サンプルヘルプ"))
			},
		},
	}
}
