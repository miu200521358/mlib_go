package config

import (
	"embed"
	"encoding/json"
	"io/fs"

)

type AppConfig struct {
	AppName    string `json:"AppName"`
	AppVersion string `json:"AppVersion"`
	IconFile   string `json:"IconFile"`
}

// ReadAppConfig アプリ設定ファイルの読み込み
func ReadAppConfig(appConfig embed.FS) AppConfig {
	fileData, err := fs.ReadFile(appConfig, "resources/app_config.json")
	if err != nil {
		return AppConfig{}
	}
	var appConfigData AppConfig
	json.Unmarshal(fileData, &appConfigData)
	return appConfigData
}

// ReadIconFile アイコンファイルの読み込み
func ReadIconFile(appConfig embed.FS) ([]byte, error) {
	fileData, err := fs.ReadFile(appConfig, "resources/icon.ico")
	if err != nil {
		return nil, err
	}
	return fileData, nil
}
