// 指示: miu200521358
package config

import (
	"encoding/json"
	"testing"
)

// TestConfigConstants は設定定数を確認する。
func TestConfigConstants(t *testing.T) {
	if AppConfigFilePath != "app/app_config.json" {
		t.Errorf("AppConfigFilePath: got=%v", AppConfigFilePath)
	}
	if UserConfigFileName != "user_config.json" {
		t.Errorf("UserConfigFileName: got=%v", UserConfigFileName)
	}
	if UserConfigLegacyFileName != "history.json" {
		t.Errorf("UserConfigLegacyFileName: got=%v", UserConfigLegacyFileName)
	}
}

// TestAppConfigUnmarshal は旧キーの変換を確認する。
func TestAppConfigUnmarshal(t *testing.T) {
	payload := `{"Name":"TestApp","Version":"1","Horizontal":true,"ViewWindowSize":{"Width":1,"Height":2}}`
	var cfg AppConfig
	if err := json.Unmarshal([]byte(payload), &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if cfg.AppName != "TestApp" {
		t.Errorf("AppName: got=%v", cfg.AppName)
	}
	if cfg.ViewerWindowSize.Width != 1 || cfg.ViewerWindowSize.Height != 2 {
		t.Errorf("ViewerWindowSize: got=%v", cfg.ViewerWindowSize)
	}
	if cfg.Env() != APP_ENV_DEV {
		t.Errorf("Env default: got=%v", cfg.Env())
	}
}
