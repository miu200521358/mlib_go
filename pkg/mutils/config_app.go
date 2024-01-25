package mutils

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
)

type AppConfig struct {
	AppName    string `json:"AppName"`
	AppVersion string `json:"AppVersion"`
}

// ReadAppConfig アプリ設定ファイルの読み込み
func ReadAppConfig(resourceFiles embed.FS) AppConfig {
	fileData, err := fs.ReadFile(resourceFiles, "resources/app_config.json")
	if err != nil {
		return AppConfig{}
	}
	var appConfigData AppConfig
	json.Unmarshal(fileData, &appConfigData)
	return appConfigData
}

// ReadIconFile アイコンファイルの読み込み
func ReadIconFile(resourceFiles embed.FS) (image.Image, error) {
	fileData, err := fs.ReadFile(resourceFiles, "resources/app.png")
	if err != nil {
		return nil, fmt.Errorf("app icon not found: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
}
