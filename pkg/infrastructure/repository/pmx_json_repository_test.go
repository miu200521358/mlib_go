package repository

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestPmxJsonRepository_Save1(t *testing.T) {
	pmxRep := NewPmxRepository(true)

	model, err := pmxRep.Load("C:/MMD/vmd_sizing_t3/pkg/usecase/base_model/model.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	rep := NewPmxJsonRepository()

	// Save the model
	jsonPath := "C:/MMD/vmd_sizing_t3/archive/sizing_model.json"
	err = rep.Save(jsonPath, model, false)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	// Read the saved JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data
	var savedData pmxJson
	err = json.Unmarshal(jsonData, &savedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Check if the saved data matches the original model
	if savedData.Name != model.Name() {
		t.Errorf("Expected model name to be '%s', got '%s'", model.Name(), savedData.Name)
	}
}

func TestPmxJsonRepository_Save2(t *testing.T) {
	pmxRep := NewPmxRepository(true)

	// model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/ゲーム/ウマ娘/ヨッシー式_メジロブライト/ヨッシー式メジロブライト.pmx")
	// model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/カイト.pmx")
	model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_あにまさ式/初音ミクVer2 準標準.pmx")
	// model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20240628/wa_129cm.pmx")
	// model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/055_鶯丸/鶯丸 さとく式 ver0.90/さとく式鶯丸ver0.90.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	rep := NewPmxJsonRepository()

	// Save the model
	jsonPath := strings.ReplaceAll(model.Path(), ".pmx", ".json")
	err = rep.Save(jsonPath, model, false)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	// Read the saved JSON file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data
	var savedData pmxJson
	err = json.Unmarshal(jsonData, &savedData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Check if the saved data matches the original model
	if savedData.Name != model.Name() {
		t.Errorf("Expected model name to be '%s', got '%s'", model.Name(), savedData.Name)
	}

	loadModel, err := rep.Load(jsonPath)
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	pmxRep.Save("C:/MMD/mlib_go/test_resources/json_pmx_loaded.pmx", loadModel, false)
}
