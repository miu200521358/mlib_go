// 指示: miu200521358
package io_model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
)

func TestPmxJsonRepository_SaveAndLoad(t *testing.T) {
	pmxRep := pmx.NewPmxRepository()

	data, err := pmxRep.Load(testResourcePath("サンプルモデル_PMX読み取り確認用.pmx"))
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", data)
	}

	rep := NewPmxJsonRepository()
	jsonPath := filepath.Join(t.TempDir(), "sizing_model.json")
	if err := rep.Save(jsonPath, modelData, io_common.SaveOptions{IncludeSystem: false}); err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var savedData pmxJSON
	if err := json.Unmarshal(jsonData, &savedData); err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	if savedData.Name != modelData.Name() {
		t.Errorf("Expected model name to be '%s', got '%s'", modelData.Name(), savedData.Name)
	}

	loadedData, err := rep.Load(jsonPath)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err)
	}
	loadedModel, ok := loadedData.(*model.PmxModel)
	if !ok {
		t.Fatalf("Expected model type to be *PmxModel, got %T", loadedData)
	}

	if loadedModel.Name() != modelData.Name() {
		t.Errorf("Expected model name to be '%s', got '%s'", modelData.Name(), loadedModel.Name())
	}
}

// testResourcePath はテストリソースのパスを組み立てる。
func testResourcePath(name string) string {
	return filepath.Join("..", "..", "..", "internal", "test_resources", name)
}
