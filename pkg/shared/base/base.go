// 指示: miu200521358
package base

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// BaseInitStage は初期化ステージを表す。
type BaseInitStage int

const (
	// BASE_INIT_CONFIG は設定初期化。
	BASE_INIT_CONFIG BaseInitStage = iota
	// BASE_INIT_I18N はi18n初期化。
	BASE_INIT_I18N
	// BASE_INIT_LOGGING はログ初期化。
	BASE_INIT_LOGGING
)

// BaseServices は基盤サービスの集合。
type BaseServices struct {
	ConfigStore  config.IConfigStore
	I18nService  i18n.II18n
	LoggerService logging.ILogger
}

// Config は設定ストアを返す。
func (b *BaseServices) Config() config.IConfigStore {
	if b == nil {
		return nil
	}
	return b.ConfigStore
}

// I18n はi18nサービスを返す。
func (b *BaseServices) I18n() i18n.II18n {
	if b == nil {
		return nil
	}
	return b.I18nService
}

// Logger はロガーを返す。
func (b *BaseServices) Logger() logging.ILogger {
	if b == nil {
		return nil
	}
	return b.LoggerService
}

// IBaseServices は基盤サービスI/F。
type IBaseServices interface {
	Config() config.IConfigStore
	I18n() i18n.II18n
	Logger() logging.ILogger
}

// DEFAULT_BASE_INIT_STAGES は既定の初期化順。
var DEFAULT_BASE_INIT_STAGES = []BaseInitStage{BASE_INIT_CONFIG, BASE_INIT_I18N, BASE_INIT_LOGGING}
