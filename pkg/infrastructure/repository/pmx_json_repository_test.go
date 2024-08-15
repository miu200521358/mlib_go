package repository

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPmxJsonRepository_Save(t *testing.T) {
	pmxRep := NewPmxRepository()

	model, err := pmxRep.Load("C:/MMD/vmd_sizing_t3/pkg/ui/model/mannequin.pmx")
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
