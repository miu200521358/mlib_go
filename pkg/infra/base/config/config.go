// 指示: miu200521358
package config

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
)

var (
	osExecutable     = os.Executable
	readFile         = os.ReadFile
	writeFile        = os.WriteFile
	loadUserConfigFn = loadUserConfig
	appRootDirFn     = func() (string, error) {
		exePath, err := osExecutable()
		if err != nil {
			return "", err
		}
		return filepath.Dir(exePath), nil
	}
)

// ConfigStore は設定ストアの実装。
type ConfigStore struct {
	appConfig  *config.AppConfig
	userConfig config.IUserConfig
}

// NewConfigStore は設定ストアを生成する。
func NewConfigStore(appConfig *config.AppConfig, userConfig config.IUserConfig) config.IConfigStore {
	return &ConfigStore{appConfig: appConfig, userConfig: userConfig}
}

// AppConfig はアプリ設定を返す。
func (cs *ConfigStore) AppConfig() *config.AppConfig {
	if cs == nil {
		return nil
	}
	return cs.appConfig
}

// UserConfig はユーザー設定を返す。
func (cs *ConfigStore) UserConfig() config.IUserConfig {
	if cs == nil {
		return nil
	}
	return cs.userConfig
}

var userConfigOnce sync.Once
var userConfigSingleton *UserConfigStore

// NewUserConfigStore はユーザー設定ストアを返す（シングルトン）。
func NewUserConfigStore() config.IUserConfig {
	userConfigOnce.Do(func() {
		userConfigSingleton = &UserConfigStore{}
	})
	return userConfigSingleton
}

// UserConfigStore はユーザー設定の実装。
type UserConfigStore struct{}

// Get はキーの値を返す。
func (u *UserConfigStore) Get(key string) (any, bool) {
	_, configMap := u.loadAll(key)
	val, ok := configMap[key]
	return val, ok
}

// Set はキーの値を保存する。
func (u *UserConfigStore) Set(key string, value any) error {
	switch v := value.(type) {
	case string:
		return u.saveValue(key, v, 1)
	case []string:
		limit := len(v)
		if limit == 0 {
			return nil
		}
		return u.SetStringSlice(key, v, limit)
	case bool:
		return u.SetBool(key, v)
	case int:
		return u.SetInt(key, v)
	default:
		return fmt.Errorf("unsupported config value type: %T", value)
	}
}

// GetStringSlice はキーのスライスを返す。
func (u *UserConfigStore) GetStringSlice(key string) []string {
	values, _ := u.loadAll(key)
	return values
}

// SetStringSlice はスライス値を保存する。
func (u *UserConfigStore) SetStringSlice(key string, values []string, limit int) error {
	return u.saveStringSlice(key, values, limit)
}

// GetBool はbool設定を返す。
func (u *UserConfigStore) GetBool(key string, defaultValue bool) bool {
	return u.boolValue(key, defaultValue)
}

// SetBool はbool設定を保存する。
func (u *UserConfigStore) SetBool(key string, value bool) error {
	return u.saveBool(key, value)
}

// GetInt はint設定を返す。
func (u *UserConfigStore) GetInt(key string, defaultValue int) int {
	return u.intValue(key, defaultValue)
}

// SetInt はint設定を保存する。
func (u *UserConfigStore) SetInt(key string, value int) error {
	return u.saveInt(key, value)
}

// GetAll はキーの値と全設定を返す。
func (u *UserConfigStore) GetAll(key string) ([]string, map[string]any) {
	values, configMap := u.loadAll(key)
	out := make(map[string]any, len(configMap))
	for k, v := range configMap {
		out[k] = v
	}
	return values, out
}

// AppRootDir はアプリルートを返す。
func (u *UserConfigStore) AppRootDir() (string, error) {
	return AppRootDir()
}

// saveValue は値を保存する。
func (u *UserConfigStore) saveValue(key, value string, limit int) error {
	if value == "" {
		return nil
	}
	return u.saveStringSlice(key, []string{value}, limit)
}

// boolValue はbool値を返す。
func (u *UserConfigStore) boolValue(key string, defaultValue bool) bool {
	values, _ := u.loadAll(key)
	if len(values) == 0 {
		return defaultValue
	}
	return values[0] == "ON"
}

// saveBool はbool値を保存する。
func (u *UserConfigStore) saveBool(key string, value bool) error {
	if value {
		return u.saveValue(key, "ON", 1)
	}
	return u.saveValue(key, "OFF", 1)
}

// intValue はint値を返す。
func (u *UserConfigStore) intValue(key string, defaultValue int) int {
	values, _ := u.loadAll(key)
	if len(values) == 0 {
		return defaultValue
	}
	parsed, err := strconv.Atoi(values[0])
	if err != nil {
		return defaultValue
	}
	return parsed
}

// saveInt はint値を保存する。
func (u *UserConfigStore) saveInt(key string, value int) error {
	return u.saveValue(key, strconv.Itoa(value), 1)
}

// loadAll は設定を読み込み、指定キーの値と全体を返す。
func (u *UserConfigStore) loadAll(key string) ([]string, map[string]any) {
	configMap, _ := loadUserConfigFn()
	values, ok := configMap[key]
	if !ok {
		return []string{}, configMap
	}

	switch list := values.(type) {
	case []any:
		out := make([]string, 0, len(list))
		for _, v := range list {
			str, ok := v.(string)
			if !ok {
				return []string{}, configMap
			}
			out = append(out, str)
		}
		return out, configMap
	case []string:
		return list, configMap
	default:
		return []string{}, configMap
	}
}

// saveStringSlice は文字列スライスを保存する。
func (u *UserConfigStore) saveStringSlice(key string, values []string, limit int) error {
	if len(values) == 0 {
		return nil
	}
	current, configMap := u.loadAll(key)
	seen := make(map[string]struct{}, len(values)+len(current))
	merged := make([]string, 0, len(values)+len(current))

	for _, v := range values {
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		merged = append(merged, v)
	}
	for _, v := range current {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		merged = append(merged, v)
	}

	if limit > 0 && len(merged) > limit {
		merged = merged[:limit]
	}

	configMap[key] = merged
	data, err := json.Marshal(configMap)
	if err != nil {
		return err
	}

	root, err := AppRootDir()
	if err != nil {
		return err
	}
	path := filepath.Join(root, config.UserConfigFileName)
	if err := writeFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

// loadUserConfig はユーザー設定を読み込む。
func loadUserConfig() (map[string]any, error) {
	root, err := AppRootDir()
	if err != nil {
		return map[string]any{}, err
	}
	path := filepath.Join(root, config.UserConfigFileName)
	data, err := readFile(path)
	if err != nil {
		path = filepath.Join(root, config.UserConfigLegacyFileName)
		data, err = readFile(path)
		if err != nil {
			return map[string]any{}, nil
		}
	}

	configMap := make(map[string]any)
	if err := json.Unmarshal(data, &configMap); err != nil {
		return map[string]any{}, nil
	}
	return configMap, nil
}

// AppRootDir はアプリルートディレクトリを返す。
func AppRootDir() (string, error) {
	return appRootDirFn()
}

// MustAppRootDir はアプリルートをpanic付きで返す。
func MustAppRootDir() string {
	root, err := AppRootDir()
	if err != nil {
		panic(err)
	}
	return root
}

// LoadAppConfig は埋め込み設定を読み込む。
func LoadAppConfig(appFiles embed.FS) (*config.AppConfig, error) {
	return loadAppConfigFS(appFiles)
}

// loadAppConfigFS はFSから設定を読み込む。
func loadAppConfigFS(appFiles fs.FS) (*config.AppConfig, error) {
	data, err := fs.ReadFile(appFiles, config.AppConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("app config read failed: %w", err)
	}
	var cfg config.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("app config parse failed: %w", err)
	}

	iconPath := cfg.IconImagePath
	if iconPath == "" {
		iconPath = config.AppIconImagePath
		cfg.IconImagePath = iconPath
	}
	if iconPath != "" {
		iconBytes, err := fs.ReadFile(appFiles, iconPath)
		if err != nil {
			return nil, fmt.Errorf("app icon read failed: %w", err)
		}
		if _, err := png.Decode(bytes.NewReader(iconBytes)); err != nil {
			return nil, fmt.Errorf("app icon decode failed: %w", err)
		}
	}
	return &cfg, nil
}
