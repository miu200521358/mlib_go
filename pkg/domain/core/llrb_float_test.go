package core

import (
	"testing"
)

func TestFloatIndexes_Prev(t *testing.T) {
	tests := []struct {
		name     string
		elements []float32
		index    float32
		expected float32
	}{
		{
			name:     "Prev element exists",
			elements: []float32{1.0, 2.0, 3.0, 4.0, 5.0},
			index:    3.0,
			expected: 2.0,
		},
		{
			name:     "Prev element does not exist",
			elements: []float32{1.0, 2.0, 3.0, 4.0, 5.0},
			index:    1.0,
			expected: 1.0, // Min value
		},
		{
			name:     "Empty tree",
			elements: []float32{},
			index:    3.0,
			expected: 0.0, // Min value of empty tree
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			floatIndexes := NewFloatIndexes()
			for _, elem := range tt.elements {
				floatIndexes.LLRB.ReplaceOrInsert(Float(elem))
			}

			result := floatIndexes.Prev(tt.index)
			if result != tt.expected {
				t.Errorf("Prev(%v) = %v; expected %v", tt.index, result, tt.expected)
			}
		})
	}
}

func TestFloatIndexes_Next(t *testing.T) {
	tests := []struct {
		name     string
		elements []float32
		index    float32
		expected float32
	}{
		{
			name:     "Next element exists",
			elements: []float32{1.0, 2.0, 3.0, 4.0, 5.0},
			index:    3.0,
			expected: 4.0,
		},
		{
			name:     "Next element does not exist",
			elements: []float32{1.0, 2.0, 3.0, 4.0, 5.0},
			index:    5.0,
			expected: 5.0, // Max value
		},
		{
			name:     "Empty tree",
			elements: []float32{},
			index:    3.0,
			expected: 3.0, // Index itself as tree is empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			floatIndexes := NewFloatIndexes()
			for _, elem := range tt.elements {
				floatIndexes.LLRB.ReplaceOrInsert(Float(elem))
			}

			result := floatIndexes.Next(tt.index)
			if result != tt.expected {
				t.Errorf("Next(%v) = %v; expected %v", tt.index, result, tt.expected)
			}
		})
	}
}
