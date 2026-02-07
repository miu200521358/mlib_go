//go:build windows
// +build windows

// 指示: miu200521358
package app

import (
	"embed"
	"image"

	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/infra/base/config"
	"github.com/miu200521358/mlib_go/pkg/infra/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/infra/base/mlogging"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
)

// AppConfigAdjuster はアプリ設定の調整関数を表す。
type AppConfigAdjuster func(appConfig *config.AppConfig)

// Result は起動時に必要な初期化結果をまとめる。
type Result struct {
	AppConfig    *config.AppConfig
	BaseServices *base.BaseServices
	Logger       *mlogging.Logger
	IconImage    image.Image
	AppIcon      *walk.Icon
}

// Init はアプリ起動に必要な設定・i18n・ログ・アイコンを初期化する。
func Init(appFiles embed.FS, i18nFiles embed.FS, adjuster AppConfigAdjuster) (*Result, error) {
	appConfig, err := config.LoadAppConfig(appFiles)
	if err != nil {
		return nil, err
	}
	if adjuster != nil {
		adjuster(appConfig)
	}

	userConfig := config.NewUserConfigStore()
	if err := i18n.InitI18n(i18nFiles, userConfig); err != nil {
		return &Result{AppConfig: appConfig}, err
	}

	logger := mlogging.NewLogger(i18n.Default())
	mlogging.SetDefaultLogger(logger)
	logging.SetDefaultLogger(logger)

	configStore := config.NewConfigStore(appConfig, userConfig)
	baseServices := &base.BaseServices{
		ConfigStore:   configStore,
		I18nService:   i18n.Default(),
		LoggerService: logger,
	}

	iconImage, appIcon := loadIcon(appFiles, appConfig, logger)

	return &Result{
		AppConfig:    appConfig,
		BaseServices: baseServices,
		Logger:       logger,
		IconImage:    iconImage,
		AppIcon:      appIcon,
	}, nil
}

// loadIcon はアプリアイコン画像とアイコンを読み込む。
func loadIcon(appFiles embed.FS, appConfig *config.AppConfig, logger *mlogging.Logger) (image.Image, *walk.Icon) {
	iconImage, iconErr := config.LoadAppIconImage(appFiles, appConfig)
	if iconErr != nil && logger != nil {
		logger.Error("アプリアイコンの読込に失敗しました: %s", iconErr.Error())
	}
	if iconImage == nil {
		return nil, nil
	}

	appIcon, iconErr := walk.NewIconFromImageForDPI(iconImage, 96)
	if iconErr != nil && logger != nil {
		logger.Error("アプリアイコンの生成に失敗しました: %s", iconErr.Error())
	}
	return iconImage, appIcon
}
