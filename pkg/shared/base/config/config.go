// 指示: miu200521358
package config

import "encoding/json"

// WindowSize はウィンドウサイズを表す。
type WindowSize struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

// AppEnv はアプリ環境を表す。
type AppEnv string

const (
	// APP_ENV_DEV は開発環境を表す。
	APP_ENV_DEV AppEnv = "dev"
	// APP_ENV_STG はステージング環境を表す。
	APP_ENV_STG AppEnv = "stg"
	// APP_ENV_PROD は本番環境を表す。
	APP_ENV_PROD AppEnv = "prod"
)

// AppConfig はアプリ設定を表す。
type AppConfig struct {
	AppName           string     `json:"AppName"`
	Version           string     `json:"Version"`
	EnvValue          AppEnv     `json:"Env"`
	Horizontal        bool       `json:"Horizontal"`
	ControlWindowSize WindowSize `json:"ControlWindowSize"`
	ViewerWindowSize  WindowSize `json:"ViewerWindowSize"`
	CloseConfirm      bool       `json:"CloseConfirm"`
	IconPath          string     `json:"IconPath"`
	IconImagePath     string     `json:"IconImagePath"`
	CursorPositionLimit int      `json:"CursorPositionLimit"`
}

// UnmarshalJSON は旧キーを吸収して設定を取り込む。
func (ac *AppConfig) UnmarshalJSON(data []byte) error {
	var raw struct {
		AppName           string     `json:"AppName"`
		Name              string     `json:"Name"`
		Version           string     `json:"Version"`
		Env               AppEnv     `json:"Env"`
		Horizontal        bool       `json:"Horizontal"`
		ControlWindowSize WindowSize `json:"ControlWindowSize"`
		ViewerWindowSize  WindowSize `json:"ViewerWindowSize"`
		ViewWindowSize    WindowSize `json:"ViewWindowSize"`
		CloseConfirm      bool       `json:"CloseConfirm"`
		IconPath          string     `json:"IconPath"`
		IconImagePath     string     `json:"IconImagePath"`
		CursorPositionLimit int      `json:"CursorPositionLimit"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	ac.AppName = raw.AppName
	if ac.AppName == "" {
		ac.AppName = raw.Name
	}
	ac.Version = raw.Version
	ac.EnvValue = raw.Env
	ac.Horizontal = raw.Horizontal
	ac.ControlWindowSize = raw.ControlWindowSize
	ac.ViewerWindowSize = raw.ViewerWindowSize
	if ac.ViewerWindowSize.Width == 0 && ac.ViewerWindowSize.Height == 0 {
		ac.ViewerWindowSize = raw.ViewWindowSize
	}
	ac.CloseConfirm = raw.CloseConfirm
	ac.IconPath = raw.IconPath
	ac.IconImagePath = raw.IconImagePath
	ac.CursorPositionLimit = raw.CursorPositionLimit
	if ac.CursorPositionLimit <= 0 {
		ac.CursorPositionLimit = 100
	}
	return nil
}

// Env は環境値を返す。
func (ac *AppConfig) Env() AppEnv {
	if ac == nil {
		return APP_ENV_DEV
	}
	if ac.EnvValue == "" {
		return APP_ENV_DEV
	}
	return ac.EnvValue
}

// IsDev は開発環境か判定する。
func (ac *AppConfig) IsDev() bool {
	return ac.Env() == APP_ENV_DEV
}

// IsStg はステージング環境か判定する。
func (ac *AppConfig) IsStg() bool {
	return ac.Env() == APP_ENV_STG
}

// IsProd は本番環境か判定する。
func (ac *AppConfig) IsProd() bool {
	return ac.Env() == APP_ENV_PROD
}

// IsCloseConfirmEnabled は終了確認の有効可否を返す。
func (ac *AppConfig) IsCloseConfirmEnabled() bool {
	if ac == nil {
		return false
	}
	return ac.CloseConfirm
}

// IUserConfig はユーザー設定のI/F。
type IUserConfig interface {
	Get(key string) (any, bool, error)
	Set(key string, value any) error
	GetStringSlice(key string) ([]string, error)
	SetStringSlice(key string, values []string, limit int) error
	GetBool(key string, defaultValue bool) (bool, error)
	SetBool(key string, value bool) error
	GetInt(key string, defaultValue int) (int, error)
	SetInt(key string, value int) error
	GetAll(key string) ([]string, map[string]any, error)
	AppRootDir() (string, error)
}

// IConfigStore は設定ストアのI/F。
type IConfigStore interface {
	AppConfig() *AppConfig
	UserConfig() IUserConfig
}

const (
	// AppConfigFilePath はアプリ設定ファイルのパス。
	AppConfigFilePath = "app/app_config.json"
	// AppIconImagePath はアプリアイコン画像のパス。
	AppIconImagePath = "app/app.png"
	// UserConfigFileName はユーザー設定ファイル名。
	UserConfigFileName = "user_config.json"
	// UserConfigLegacyFileName は旧ユーザー設定ファイル名。
	UserConfigLegacyFileName = "history.json"
	// UserConfigKeyFpsLimit はFPS制限のキー。
	UserConfigKeyFpsLimit = "fps_limit"
	// UserConfigKeyLang は言語キー。
	UserConfigKeyLang = "lang"
	// UserConfigKeyWindowLinkage は画面連動キー。
	UserConfigKeyWindowLinkage = "window_linkage"
	// UserConfigKeyFrameDrop はフレームドロップキー。
	UserConfigKeyFrameDrop = "frame_drop"
)
