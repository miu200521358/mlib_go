package mconfig

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/fs"

	"github.com/miu200521358/walk/pkg/walk"
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
	CloseConfirm      bool   `json:"CloseConfirm"`
	Env               string
	Icon              *walk.Icon
	IconImage         image.Image
}

func (ac *AppConfig) IsSetEnv() bool {
	if ac == nil {
		return false
	}
	return ac.Env != ""
}

func (ac *AppConfig) IsEnvProd() bool {
	if ac == nil {
		return false
	}
	return ac.Env == "prod"
}

func (ac *AppConfig) IsEnvDev() bool {
	if ac == nil {
		// ローカル起動時もDevとみなす
		return true
	}
	return ac.Env == "dev" || ac.Env == "debug" || ac.Env == ""
}

func (ac *AppConfig) IsEnvStg() bool {
	if ac == nil {
		return false
	}
	return ac.Env == "stg"
}

func (ac *AppConfig) IsCloseConfirm() bool {
	if ac == nil {
		return false
	}
	return ac.CloseConfirm
}

// LoadAppConfig アプリ設定ファイルの読み込み
func LoadAppConfig(appFiles embed.FS) *AppConfig {
	fileData, err := fs.ReadFile(appFiles, "app/app_config.json")
	if err != nil {
		return &AppConfig{}
	}
	var appConfigData AppConfig
	err = json.Unmarshal(fileData, &appConfigData)
	if err != nil {
		return &AppConfig{}
	}

	err = appConfigData.loadImageFile(appFiles)
	if err != nil {
		return &AppConfig{}
	}

	return &appConfigData
}

// LoadIconFile アイコンファイルの読み込み
func (ac *AppConfig) loadImageFile(resources embed.FS) error {
	fileData, err := fs.ReadFile(resources, "app/app.png")
	if err != nil {
		return fmt.Errorf("image not found: %v", err)
	}
	file := bytes.NewReader(fileData)

	ac.IconImage, err = png.Decode(file)
	if err != nil {
		return err
	}

	ac.Icon, err = walk.NewIconFromImageForDPI(ac.IconImage, 96)
	if err != nil {
		return err
	}

	return nil
}
