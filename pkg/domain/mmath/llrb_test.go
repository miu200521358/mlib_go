// 指示: miu200521358
package mmath

import (
	"testing"
)

func TestIntIndexes_Prev(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		index    int
		expected int
	}{
		{
			name:     "Empty tree",
			elements: []int{},
			index:    5,
			expected: 0,
		},
		{
			name:     "Single element, index not in tree",
			elements: []int{10},
			index:    5,
			expected: 10,
		},
		{
			name:     "Single element, index in tree",
			elements: []int{10},
			index:    10,
			expected: 10,
		},
		{
			name:     "Multiple elements, index not in tree",
			elements: []int{1, 3, 5, 7, 9},
			index:    6,
			expected: 5,
		},
		{
			name:     "Multiple elements, index in tree",
			elements: []int{1, 3, 5, 7, 9},
			index:    7,
			expected: 5,
		},
		{
			name:     "Index is the smallest element",
			elements: []int{1, 3, 5, 7, 9},
			index:    1,
			expected: 1,
		},
		{
			name:     "Index is the largest element",
			elements: []int{1, 3, 5, 7, 9},
			index:    9,
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexes := NewLlrbIndexes[int]()
			for _, elem := range tt.elements {
				indexes.InsertNoReplace(NewLlrbItem(elem))
			}

			result := indexes.Prev(tt.index)
			if result != tt.expected {
				t.Errorf("Prev(%d) = %d; expected %d", tt.index, result, tt.expected)
			}
		})
	}
}

func TestIntIndexes_Next(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		index    int
		expected int
	}{
		{
			name:     "Empty tree",
			elements: []int{},
			index:    5,
			expected: 5,
		},
		{
			name:     "Single element, index not in tree",
			elements: []int{10},
			index:    5,
			expected: 10,
		},
		{
			name:     "Single element, index in tree",
			elements: []int{10},
			index:    10,
			expected: 10,
		},
		{
			name:     "Multiple elements, index not in tree",
			elements: []int{1, 3, 5, 7, 9},
			index:    6,
			expected: 7,
		},
		{
			name:     "Multiple elements, index in tree",
			elements: []int{1, 3, 5, 7, 9},
			index:    7,
			expected: 9,
		},
		{
			name:     "Index is the smallest element",
			elements: []int{1, 3, 5, 7, 9},
			index:    1,
			expected: 3,
		},
		{
			name:     "Index is the largest element",
			elements: []int{1, 3, 5, 7, 9},
			index:    9,
			expected: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexes := NewLlrbIndexes[int]()
			for _, elem := range tt.elements {
				indexes.InsertNoReplace(NewLlrbItem(elem))
			}

			result := indexes.Next(tt.index)
			if result != tt.expected {
				t.Errorf("Next(%d) = %d; expected %d", tt.index, result, tt.expected)
			}
		})
	}
}

