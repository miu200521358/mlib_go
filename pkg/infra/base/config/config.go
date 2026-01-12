// 指示: miu200521358
package config

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	baseerr "github.com/miu200521358/mlib_go/pkg/shared/base/err"
)

var (
	osExecutable     = os.Executable
	readFile         = os.ReadFile
	writeFile        = os.WriteFile
	loadUserConfigFn = loadUserConfig
	appRootDirFn     = func() (string, error) {
		exePath, err := osExecutable()
		if err != nil {
			return "", newAppRootDirResolveFailed(baseerr.NewOsPackageError("os.Executableに失敗しました", err))
		}
		return filepath.Dir(exePath), nil
	}
)

const (
	appConfigLoadFailedErrorID         = "95201"
	userConfigSaveFailedErrorID        = "95202"
	appRootDirResolveFailedErrorID     = "95203"
	configValueTypeNotSupportedErrorID = "95204"
)

// newInternalError は内部エラーとして共通委譲エラーを生成する。
func newInternalError(id string, message string, cause error) error {
	return baseerr.NewCommonError(id, baseerr.ErrorKindInternal, message, cause)
}

// newAppConfigLoadFailed はアプリ設定の読込失敗エラーを生成する。
func newAppConfigLoadFailed(message string, cause error) error {
	return newInternalError(appConfigLoadFailedErrorID, message, cause)
}

// newUserConfigSaveFailed はユーザー設定の保存失敗エラーを生成する。
func newUserConfigSaveFailed(message string, cause error) error {
	return newInternalError(userConfigSaveFailedErrorID, message, cause)
}

// newAppRootDirResolveFailed はアプリルート取得失敗エラーを生成する。
func newAppRootDirResolveFailed(cause error) error {
	return newInternalError(appRootDirResolveFailedErrorID, "アプリルート取得に失敗しました", cause)
}

// newConfigValueTypeNotSupported は設定値の型未対応エラーを生成する。
func newConfigValueTypeNotSupported(message string) error {
	return newInternalError(configValueTypeNotSupportedErrorID, message, nil)
}

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

// Get はキーの値を返す（読込失敗時は error を返す）。
func (u *UserConfigStore) Get(key string) (any, bool, error) {
	_, configMap, err := u.loadAll(key)
	if err != nil {
		return nil, false, err
	}
	val, ok := configMap[key]
	return val, ok, nil
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
		return newConfigValueTypeNotSupported(fmt.Sprintf("未対応の設定値型です: %T", value))
	}
}

// GetStringSlice はキーのスライスを返す（読込失敗時は error を返す）。
func (u *UserConfigStore) GetStringSlice(key string) ([]string, error) {
	values, _, err := u.loadAll(key)
	if err != nil {
		return []string{}, err
	}
	return values, nil
}

// SetStringSlice はスライス値を保存する。
func (u *UserConfigStore) SetStringSlice(key string, values []string, limit int) error {
	return u.saveStringSlice(key, values, limit)
}

// GetBool はbool設定を返す（読込失敗時は error を返す）。
func (u *UserConfigStore) GetBool(key string, defaultValue bool) (bool, error) {
	return u.boolValue(key, defaultValue)
}

// SetBool はbool設定を保存する。
func (u *UserConfigStore) SetBool(key string, value bool) error {
	return u.saveBool(key, value)
}

// GetInt はint設定を返す（読込失敗時は error を返す）。
func (u *UserConfigStore) GetInt(key string, defaultValue int) (int, error) {
	return u.intValue(key, defaultValue)
}

// SetInt はint設定を保存する。
func (u *UserConfigStore) SetInt(key string, value int) error {
	return u.saveInt(key, value)
}

// GetAll はキーの値と全設定を返す（読込失敗時は error を返す）。
func (u *UserConfigStore) GetAll(key string) ([]string, map[string]any, error) {
	values, configMap, err := u.loadAll(key)
	if err != nil {
		return []string{}, map[string]any{}, err
	}
	out := make(map[string]any, len(configMap))
	for k, v := range configMap {
		out[k] = v
	}
	return values, out, nil
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
func (u *UserConfigStore) boolValue(key string, defaultValue bool) (bool, error) {
	values, err := u.GetStringSlice(key)
	if err != nil {
		return defaultValue, err
	}
	if len(values) == 0 {
		return defaultValue, nil
	}
	return values[0] == "ON", nil
}

// saveBool はbool値を保存する。
func (u *UserConfigStore) saveBool(key string, value bool) error {
	if value {
		return u.saveValue(key, "ON", 1)
	}
	return u.saveValue(key, "OFF", 1)
}

// intValue はint値を返す。
func (u *UserConfigStore) intValue(key string, defaultValue int) (int, error) {
	values, err := u.GetStringSlice(key)
	if err != nil {
		return defaultValue, err
	}
	if len(values) == 0 {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(values[0])
	if err != nil {
		return defaultValue, nil
	}
	return parsed, nil
}

// saveInt はint値を保存する。
func (u *UserConfigStore) saveInt(key string, value int) error {
	return u.saveValue(key, strconv.Itoa(value), 1)
}

// loadAll は設定を読み込み、指定キーの値と全体を返す。
func (u *UserConfigStore) loadAll(key string) ([]string, map[string]any, error) {
	configMap, err := loadUserConfigFn()
	if err != nil {
		return []string{}, map[string]any{}, err
	}
	values, ok := configMap[key]
	if !ok {
		return []string{}, configMap, nil
	}

	switch list := values.(type) {
	case []any:
		out := make([]string, 0, len(list))
		for _, v := range list {
			str, ok := v.(string)
			if !ok {
				return []string{}, map[string]any{}, newConfigValueTypeNotSupported("user_config.jsonの値が未対応です: " + key)
			}
			out = append(out, str)
		}
		return out, configMap, nil
	case []string:
		return list, configMap, nil
	default:
		return []string{}, map[string]any{}, newConfigValueTypeNotSupported("user_config.jsonの値が未対応です: " + key)
	}
}

// saveStringSlice は文字列スライスを保存する。
func (u *UserConfigStore) saveStringSlice(key string, values []string, limit int) error {
	if len(values) == 0 {
		return nil
	}
	current, configMap, err := u.loadAll(key)
	if err != nil {
		return newUserConfigSaveFailed("user_config.jsonの保存に失敗しました", err)
	}
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
		cause := baseerr.NewJsonPackageError("user_config.jsonの保存用JSON生成に失敗しました", err)
		return newUserConfigSaveFailed("user_config.jsonの保存に失敗しました", cause)
	}

	root, err := AppRootDir()
	if err != nil {
		return newUserConfigSaveFailed("user_config.jsonの保存に失敗しました", err)
	}
	path := filepath.Join(root, config.UserConfigFileName)
	if err := writeFile(path, data, 0644); err != nil {
		cause := baseerr.NewOsPackageError("user_config.jsonの書き込みに失敗しました: "+path, err)
		return newUserConfigSaveFailed("user_config.jsonの保存に失敗しました", cause)
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
		if !errors.Is(err, os.ErrNotExist) {
			return map[string]any{}, baseerr.NewOsPackageError("user_config.jsonの読込に失敗しました: "+path, err)
		}
		path = filepath.Join(root, config.UserConfigLegacyFileName)
		data, err = readFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return map[string]any{}, nil
			}
			return map[string]any{}, baseerr.NewOsPackageError("history.jsonの読込に失敗しました: "+path, err)
		}
	}

	configMap := make(map[string]any)
	if err := json.Unmarshal(data, &configMap); err != nil {
		return map[string]any{}, baseerr.NewJsonPackageError("設定JSONの解析に失敗しました: "+path, err)
	}
	return configMap, nil
}

// AppRootDir はアプリルートディレクトリを返す。
func AppRootDir() (string, error) {
	root, err := appRootDirFn()
	if err != nil {
		if ce, ok := err.(*baseerr.CommonError); ok && ce.ErrorID() == appRootDirResolveFailedErrorID {
			return "", err
		}
		return "", newAppRootDirResolveFailed(err)
	}
	return root, nil
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
		cause := baseerr.NewFsPackageError("app_config.jsonの読込に失敗しました: "+config.AppConfigFilePath, err)
		return nil, newAppConfigLoadFailed("app_config.jsonの読込に失敗しました", cause)
	}
	var cfg config.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		cause := baseerr.NewJsonPackageError("app_config.jsonの解析に失敗しました: "+config.AppConfigFilePath, err)
		return nil, newAppConfigLoadFailed("app_config.jsonの解析に失敗しました", cause)
	}

	iconPath := cfg.IconImagePath
	if iconPath == "" {
		iconPath = config.AppIconImagePath
		cfg.IconImagePath = iconPath
	}
	if iconPath != "" {
		iconBytes, err := fs.ReadFile(appFiles, iconPath)
		if err != nil {
			cause := baseerr.NewFsPackageError("アプリアイコンの読込に失敗しました: "+iconPath, err)
			return nil, newAppConfigLoadFailed("アプリアイコンの読込に失敗しました", cause)
		}
		if _, err := png.Decode(bytes.NewReader(iconBytes)); err != nil {
			cause := baseerr.NewImagePackageError("アプリアイコンのデコードに失敗しました: "+iconPath, err)
			return nil, newAppConfigLoadFailed("アプリアイコンのデコードに失敗しました", cause)
		}
	}
	return &cfg, nil
}
