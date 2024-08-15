package repository

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPmxJsonRepository_Save1(t *testing.T) {
	pmxRep := NewPmxRepository()

	model, err := pmxRep.Load("C:/MMD/vmd_sizing_t3/pkg/usecase/model/mannequin.pmx")
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
	pmxRep := NewPmxRepository()

	model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/003_三日月宗近/三日月宗近 わち式 （刀ミュインナーβ）/わち式三日月宗近（刀ミュインナーβ）.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	rep := NewPmxJsonRepository()

	// Save the model
	jsonPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/刀剣乱舞/003_三日月宗近/三日月宗近 わち式 （刀ミュインナーβ）/わち式三日月宗近（刀ミュインナーβ）.json"
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

func TestPmxJsonRepository_Save3(t *testing.T) {
	pmxRep := NewPmxRepository()

	model, err := pmxRep.Load("D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20240628/wa_129cm.pmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	}

	rep := NewPmxJsonRepository()

	// Save the model
	jsonPath := "D:/MMD/MikuMikuDance_v926x64/UserFile/Model/_VMDサイジング/wa_129cm 20240628/wa_129cm.json"
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
