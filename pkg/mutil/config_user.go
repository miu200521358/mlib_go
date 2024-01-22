package mutil

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const USER_CONFIG_FILE_NAME = "user_config.json"

// 設定の保存
func SaveUserConfig(key string, value string, limit int) error {
	if value == "" {
		return nil
	}

	// UserConfigファイルをロードする
	existingValues, config := LoadUserConfig(key)
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
	configFilePath := filepath.Join(GetAppRootDir(), USER_CONFIG_FILE_NAME)

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
	configFilePath := filepath.Join(GetAppRootDir(), USER_CONFIG_FILE_NAME)
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
