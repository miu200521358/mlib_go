package mconfig

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/mutils"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
)

const USER_CONFIG_FILE_NAME = "user_config.json"
const USER_CONFIG_OLD_FILE_NAME = "history.json"

const (
	FRAME_DROP     = "FrameDrop"
	HIGH_SPEC_MODE = "HighSpecMode"
	PHYSICS        = "Physics"
)

// 設定の保存
func SaveUserConfig(key string, value string, limit int) error {
	if value == "" {
		return nil
	}

	// UserConfigファイルをロードする
	existingValues, config := LoadUserConfigAll(key)
	values := []string{value}

	// Determine the upper limit based on the smaller value between len(existingValues) and limit
	upperLimit := len(existingValues) + 1
	if limit < upperLimit {
		upperLimit = limit
	}

	// Remove the value if it already exists in existingValues
	for i := 0; i < (upperLimit - 1); i++ {
		if existingValues[i] != value {
			values = append(values, existingValues[i:i+1]...)
		}
	}

	// 同じ値があって、結果件数が変わらない場合、再設定
	if len(values) < upperLimit {
		upperLimit = len(values)
	}

	// Add key-value elements to the config map
	config[key] = values[:upperLimit]

	// Create a JSON representation of the config map without newlines and indentation
	updatedData, err := json.Marshal(config)
	if err != nil {
		updatedData = []byte("{}")
	}

	// ファイルのフルパスを取得
	configFilePath := filepath.Join(mutils.GetAppRootDir(), USER_CONFIG_FILE_NAME)

	// Overwrite the config.json file with the updated JSON data
	err = os.WriteFile(configFilePath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadUserConfig(key string) []string {
	existingValues, _ := LoadUserConfigAll(key)
	return existingValues
}

// 設定の読み込み
func LoadUserConfigAll(key string) ([]string, map[string]interface{}) {
	// Configファイルのフルパスを取得
	configFilePath := filepath.Join(mutils.GetAppRootDir(), USER_CONFIG_FILE_NAME)
	mlog.D("LoadConfig: %s: %s", key, configFilePath)

	// Read the config.json file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		// user_config.jsonがない場合、history.jsonを読み込む(次回以降はuser_config.jsonに保存される)
		configFilePath = filepath.Join(mutils.GetAppRootDir(), USER_CONFIG_OLD_FILE_NAME)
		data, err = os.ReadFile(configFilePath)
		if err != nil {
			data = []byte("{}")
		}
	}

	// Unmarshal the JSON data into a map
	config := make(map[string]interface{})
	err = json.Unmarshal(data, &config)
	if err != nil {
		return []string{}, config
	}

	// Check if the value is nil and initialize it as an empty slice of strings
	if config[key] == nil {
		return []string{}, config
	}

	// Convert the interface{} type to []interface{}
	values, ok := config[key].([]interface{})
	if !ok {
		return []string{}, config
	}

	// Convert each element to string
	result := make([]string, len(values))
	for i, v := range values {
		if str, ok := v.(string); ok {
			result[i] = str
		} else {
			return []string{}, config
		}
	}

	return result, config
}
