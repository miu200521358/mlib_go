package mutils

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		input      string
		separators []string
		expected   []string
	}{
		{
			input:      "a,b,c",
			separators: []string{","},
			expected:   []string{"a", "b", "c"},
		},
		{
			input:      "a,b;c",
			separators: []string{",", ";"},
			expected:   []string{"a", "b", "c"},
		},
		{
			input:      "a,b;c.d",
			separators: []string{",", ";", "."},
			expected:   []string{"a", "b", "c", "d"},
		},
		{
			input:      "a,b|c",
			separators: []string{"|"},
			expected:   []string{"a,b", "c"},
		},
		{
			input:      "",
			separators: []string{","},
			expected:   []string{""},
		},
	}

	for _, test := range tests {
		result := Split(test.input, test.separators)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Split(%q, %v) = %v; want %v", test.input, test.separators, result, test.expected)
		}
	}
}
