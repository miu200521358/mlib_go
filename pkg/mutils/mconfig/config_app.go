package mconfig

import (
	"embed"
	"encoding/json"
	"image"
	"io/fs"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type AppConfig struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	Env     string `json:"Env"`
}

func (ac *AppConfig) IsEnvProd() bool {
	if ac == nil {
		return false
	}
	return ac.Env == "prod"
}

func (ac *AppConfig) IsEnvDev() bool {
	if ac == nil {
		return false
	}
	return ac.Env == "dev"
}

// LoadAppConfig アプリ設定ファイルの読み込み
func LoadAppConfig(resourceFiles embed.FS) *AppConfig {
	fileData, err := fs.ReadFile(resourceFiles, "resources/app_config.json")
	if err != nil {
		return &AppConfig{}
	}
	var appConfigData AppConfig
	json.Unmarshal(fileData, &appConfigData)
	return &appConfigData
}

// LoadIconFile アイコンファイルの読み込み
func LoadIconFile(resourceFiles embed.FS) (*image.Image, error) {
	return mutils.LoadImageFromResources(resourceFiles, "resources/app.png")
}
