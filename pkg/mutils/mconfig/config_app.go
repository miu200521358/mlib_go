package mconfig

import (
	"embed"
	"encoding/json"
	"image"
	"io/fs"

	"github.com/miu200521358/mlib_go/pkg/mutils"
)

type Size struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

type AppConfig struct {
	Name              string `json:"Name"`
	Version           string `json:"Version"`
	Horizontal        bool   `json:"Horizontal"`
	ControlWindowSize Size   `json:"ControlWindowSize"`
	ViewWindowSize    Size   `json:"ViewWindowSize"`
	Env               string
	IconImage         *image.Image
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
func LoadAppConfig(appFiles embed.FS) *AppConfig {
	fileData, err := fs.ReadFile(appFiles, "app/app_config.json")
	if err != nil {
		return &AppConfig{}
	}
	var appConfigData AppConfig
	json.Unmarshal(fileData, &appConfigData)

	appConfigData.IconImage, _ = loadIconFile(appFiles)

	return &appConfigData
}

// LoadIconFile アイコンファイルの読み込み
func loadIconFile(resources embed.FS) (*image.Image, error) {
	return mutils.LoadImageFromResources(resources, "app/app.png")
}
