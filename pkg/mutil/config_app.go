package mutil

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"os"
)

type AppConfig struct {
	AppName    string `json:"AppName"`
	AppVersion string `json:"AppVersion"`
	IconFile   string `json:"IconFile"`
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
	imgFile, err := os.Open("resources/app.png")
	if err != nil {
		return nil, fmt.Errorf("app icon not found on disk: %v", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	return img, nil
}
