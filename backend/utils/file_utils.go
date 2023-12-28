package mlib

import (
	"encoding/json"
	"fmt"
	"os"

)

const configFilePath = "config.json"

// 設定の保存
func SaveConfig(key string, values []string, limit int) error {
	// Unmarshal the JSON data into a map
	println(configFilePath)
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		data = []byte("{}")
	}

	// Unmarshal the JSON data into a map
	config := make(map[string]interface{})
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Check if the value is nil and initialize it as an empty slice of strings
	if config[key] == nil {
		config[key] = []string{}
	}

	// Convert the interface{} type to []interface{}
	existingValues, ok := config[key].([]interface{})
	if !ok {
		existingValues = make([]interface{}, 0)
	}

	// Convert each element to string
	existingValuesStrList := make([]string, len(existingValues))
	for i, v := range existingValues {
		existingValuesStrList[i] = fmt.Sprintf("%v", v)
	}

	// Add key-value elements to the config map
	config[key] = append(values, existingValuesStrList...)[:limit]

	// Create a JSON representation of the config map without newlines and indentation
	updatedData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// Overwrite the config.json file with the updated JSON data
	err = os.WriteFile(configFilePath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// 設定の読み込み
func LoadConfig(key string) []string {
	// Read the config.json file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		data = []byte("{}")
	}

	// Unmarshal the JSON data into a map
	config := make(map[string]interface{})
	err = json.Unmarshal(data, &config)
	if err != nil {
		return []string{}
	}

	// Check if the value is nil and initialize it as an empty slice of strings
	if config[key] == nil {
		return []string{}
	}

	// Convert the interface{} type to []interface{}
	values, ok := config[key].([]interface{})
	if !ok {
		return []string{}
	}

	// Convert each element to string
	result := make([]string, len(values))
	for i, v := range values {
		if str, ok := v.(string); ok {
			result[i] = str
		} else {
			return []string{}
		}
	}

	return result
}
