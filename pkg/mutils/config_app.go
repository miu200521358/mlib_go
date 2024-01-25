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

// LoadAppConfig アプリ設定ファイルの読み込み
func LoadAppConfig(resourceFiles embed.FS) AppConfig {
	fileData, err := fs.ReadFile(resourceFiles, "resources/app_config.json")
	if err != nil {
		return AppConfig{}
	}
	var appConfigData AppConfig
	json.Unmarshal(fileData, &appConfigData)
	return appConfigData
}

// LoadIconFile アイコンファイルの読み込み
func LoadIconFile(resourceFiles embed.FS) (image.Image, error) {
	return LoadImageFromResources(resourceFiles, "resources/app.png")
}

// ReadIconFile アイコンファイルの読み込み
func LoadImageFromResources(resourceFiles embed.FS, fileName string) (image.Image, error) {
	fileData, err := fs.ReadFile(resourceFiles, fileName)
	if err != nil {
		return nil, fmt.Errorf("image not found: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
}
