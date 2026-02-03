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
	if cfg.CursorPositionLimit != 100 {
		t.Errorf("CursorPositionLimit default: got=%v", cfg.CursorPositionLimit)
	}
}

// TestAppConfigUnmarshalAppName はAppName優先とViewerWindowSize優先を確認する。
func TestAppConfigUnmarshalAppName(t *testing.T) {
	payload := `{"AppName":"MainApp","Name":"Legacy","Env":"stg","ViewerWindowSize":{"Width":3,"Height":4},"ViewWindowSize":{"Width":7,"Height":8},"CloseConfirm":true,"CursorPositionLimit":5}`
	var cfg AppConfig
	if err := json.Unmarshal([]byte(payload), &cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if cfg.AppName != "MainApp" {
		t.Errorf("AppName: got=%v", cfg.AppName)
	}
	if cfg.ViewerWindowSize.Width != 3 || cfg.ViewerWindowSize.Height != 4 {
		t.Errorf("ViewerWindowSize: got=%v", cfg.ViewerWindowSize)
	}
	if cfg.Env() != APP_ENV_STG || !cfg.IsStg() {
		t.Errorf("Env STG: got=%v", cfg.Env())
	}
	if !cfg.IsCloseConfirmEnabled() {
		t.Errorf("CloseConfirm: expected true")
	}
	if cfg.CursorPositionLimit != 5 {
		t.Errorf("CursorPositionLimit: got=%v", cfg.CursorPositionLimit)
	}
}

// TestAppConfigUnmarshalInvalid はJSONエラーを確認する。
func TestAppConfigUnmarshalInvalid(t *testing.T) {
	var cfg AppConfig
	if err := cfg.UnmarshalJSON([]byte("{invalid")); err == nil {
		t.Errorf("UnmarshalJSON expected error")
	}
}

// TestAppConfigEnvVariants はEnv判定を確認する。
func TestAppConfigEnvVariants(t *testing.T) {
	var nilCfg *AppConfig
	if nilCfg.Env() != APP_ENV_DEV || !nilCfg.IsDev() {
		t.Errorf("nil Env default: got=%v", nilCfg.Env())
	}
	if nilCfg.IsCloseConfirmEnabled() {
		t.Errorf("nil CloseConfirm should be false")
	}
	cfg := &AppConfig{}
	if cfg.Env() != APP_ENV_DEV || !cfg.IsDev() {
		t.Errorf("empty Env default: got=%v", cfg.Env())
	}
	cfg.EnvValue = APP_ENV_PROD
	if cfg.Env() != APP_ENV_PROD || !cfg.IsProd() || cfg.IsStg() {
		t.Errorf("Env PROD: got=%v", cfg.Env())
	}
	if cfg.IsCloseConfirmEnabled() {
		t.Errorf("CloseConfirm default: expected false")
	}
}

// TestApplyBuildEnv はビルド環境値の反映を確認する。
func TestApplyBuildEnv(t *testing.T) {
	ApplyBuildEnv(nil, "prod")

	tests := []struct {
		name     string
		initial  AppEnv
		buildEnv string
		want     AppEnv
	}{
		{name: "空は維持", initial: APP_ENV_STG, buildEnv: "", want: APP_ENV_STG},
		{name: "空白は維持", initial: APP_ENV_STG, buildEnv: "  ", want: APP_ENV_STG},
		{name: "dev設定", initial: APP_ENV_PROD, buildEnv: "dev", want: APP_ENV_DEV},
		{name: "debugはdev", initial: APP_ENV_PROD, buildEnv: "debug", want: APP_ENV_DEV},
		{name: "大文字prod", initial: APP_ENV_DEV, buildEnv: "PrOd", want: APP_ENV_PROD},
		{name: "stg設定", initial: APP_ENV_DEV, buildEnv: "stg", want: APP_ENV_STG},
		{name: "不明は維持", initial: APP_ENV_STG, buildEnv: "unknown", want: APP_ENV_STG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{EnvValue: tt.initial}
			ApplyBuildEnv(cfg, tt.buildEnv)
			if cfg.EnvValue != tt.want {
				t.Errorf("ApplyBuildEnv: got=%v want=%v", cfg.EnvValue, tt.want)
			}
		})
	}
}
