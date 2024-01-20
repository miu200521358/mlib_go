package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/utils/util_file"

)

const USER_CONFIG_FILE_NAME = "user_config.json"

// 設定の保存
func SaveUserConfig(key string, value string, limit int) error {

	// UserConfigファイルをロードする
	existingValues, config := LoadUserConfig(key)

	// Add key-value elements to the config map
	values := append([]string{value}, existingValues...)

	// newSliceの長さがlimit以下であることを確認し、必要に応じて調整
	if limit > len(values) {
		limit = len(values)
	}

	// Add key-value elements to the config map
	config[key] = values[:limit]

	// Create a JSON representation of the config map without newlines and indentation
	updatedData, err := json.Marshal(config)
	if err != nil {
		updatedData = []byte("{}")
	}

	// ファイルのフルパスを取得
	configFilePath := filepath.Join(util_file.GetAppRootDir(), USER_CONFIG_FILE_NAME)

	// Overwrite the config.json file with the updated JSON data
	err = os.WriteFile(configFilePath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// 設定の読み込み
func LoadUserConfig(key string) ([]string, map[string]interface{}) {
	// Configファイルのフルパスを取得
	configFilePath := filepath.Join(util_file.GetAppRootDir(), USER_CONFIG_FILE_NAME)
	println("LoadConfig: ", key, configFilePath)

	// Read the config.json file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		data = []byte("{}")
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
