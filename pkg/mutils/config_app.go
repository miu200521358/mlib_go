package mutils

import (
	"embed"
	"encoding/json"
	"image"
	"io/fs"

	"github.com/miu200521358/walk/pkg/walk"
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

// LoadImageFile 画像ファイルの読み込み
func LoadImageFile(resourceFiles embed.FS, imagePath string, dpi int) (walk.Image, error) {
	image, err := LoadImageFromResources(resourceFiles, imagePath)
	if err != nil {
		return nil, err
	}
	img, err := walk.NewIconFromImageForDPI(image, dpi)
	if err != nil {
		return nil, err
	}
	return img, nil
}
