//go:build windows
// +build windows

// 指示: miu200521358
package controller

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
)

// MenuSeparatorKey はメニュー区切りを表す特殊キー。
const MenuSeparatorKey = "__separator__"

// MenuMessageItem はメニューの表示キーとメッセージキーを保持する。
type MenuMessageItem struct {
	TitleKey   string
	MessageKey string
}

// BuildMenuItems はキー一覧からメニュー項目を生成する。
func BuildMenuItems(translator i18n.II18n, logger logging.ILogger, keys []string) []declarative.MenuItem {
	if len(keys) == 0 {
		return nil
	}
	messageItems := make([]MenuMessageItem, 0, len(keys))
	for _, key := range keys {
		messageItems = append(messageItems, MenuMessageItem{TitleKey: key, MessageKey: key})
	}
	return BuildMenuItemsWithMessages(translator, logger, messageItems)
}

// BuildMenuItemsWithMessages は表示キーとメッセージキーの組からメニュー項目を生成する。
func BuildMenuItemsWithMessages(translator i18n.II18n, logger logging.ILogger, items []MenuMessageItem) []declarative.MenuItem {
	if len(items) == 0 {
		return nil
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	menuItems := make([]declarative.MenuItem, 0, len(items))
	for _, item := range items {
		if item.TitleKey == MenuSeparatorKey {
			menuItems = append(menuItems, declarative.Separator{})
			continue
		}
		title := i18n.TranslateOrMark(translator, item.TitleKey)
		message := title
		if item.MessageKey != "" {
			message = i18n.TranslateOrMark(translator, item.MessageKey)
		}
		menuItems = append(menuItems, declarative.Action{
			Text: title,
			OnTriggered: func() {
				logMenuMessage(logger, message)
			},
		})
	}
	return menuItems
}

// logMenuMessage はメニュー選択時のログを出力する。
func logMenuMessage(logger logging.ILogger, message string) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if lineLogger, ok := logger.(interface {
		InfoLine(msg string, params ...any)
	}); ok {
		lineLogger.InfoLine(message)
		return
	}
	logger.Info(message)
}
